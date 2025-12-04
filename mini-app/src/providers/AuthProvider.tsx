'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { api, setAuthToken } from '@/lib/api';
import type { Profile } from '@/types/api';

interface AuthContextValue {
  token: string | null;
  profile: Profile | null;
  status: 'loading' | 'ready' | 'error';
  error: string | null;
  refreshProfile: () => Promise<void>;
  reauthenticate: () => Promise<void>;
  setProfile: (profile: Profile) => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);
const TOKEN_STORAGE_KEY = 'auth_token';

type BrowserWindow = Window & typeof globalThis;

type MaybeBrowser = typeof globalThis & { window?: BrowserWindow };

function getBrowserWindow(): BrowserWindow | null {
  if (typeof globalThis === 'undefined') return null;
  const maybeBrowser = globalThis as MaybeBrowser;
  return maybeBrowser.window ?? null;
}

function getLocalStorage(): Storage | null {
  const win = getBrowserWindow();
  return win?.localStorage ?? null;
}

function getQueryParam(name: string): string | null {
  const win = getBrowserWindow();
  if (!win) return null;
  const search = win.location?.search ?? '';
  if (!search) return null;
  return new URLSearchParams(search).get(name);
}

function removeQueryParam(name: string) {
  const win = getBrowserWindow();
  if (!win) return;
  const url = new URL(win.location.href);
  if (!url.searchParams.has(name)) return;
  url.searchParams.delete(name);
  win.history.replaceState({}, '', `${url.pathname}${url.search}${url.hash}`);
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(null);
  const [profile, setProfile] = useState<Profile | null>(null);
  const [status, setStatus] = useState<'loading' | 'ready' | 'error'>('loading');
  const [error, setError] = useState<string | null>(null);

  const loadProfile = useCallback(async () => {
    const data = await api.getProfile();
    setProfile(data);
  }, []);

  const finishAuth = useCallback(
    async (jwt: string) => {
      const storage = getLocalStorage();
      storage?.setItem(TOKEN_STORAGE_KEY, jwt);
      setAuthToken(jwt);
      setToken(jwt);
      await loadProfile();
      setStatus('ready');
      setError(null);
    },
    [loadProfile]
  );

  const authenticate = useCallback(async () => {
    setStatus('loading');
    setError(null);

    const storage = getLocalStorage();
    if (!storage) {
      setAuthToken(null);
      setToken(null);
      setStatus('error');
      setError('localStorage недоступен. Откройте мини-апп в браузере или через Telegram.');
      return;
    }

    const urlToken = getQueryParam('token');
    const devToken = process.env.NEXT_PUBLIC_DEV_TOKEN || null;
    const incomingToken = urlToken || devToken;
    if (incomingToken) {
      storage.setItem(TOKEN_STORAGE_KEY, incomingToken);
      if (urlToken) {
        removeQueryParam('token');
      }
    }

    const storedToken = storage.getItem(TOKEN_STORAGE_KEY);
    if (!storedToken) {
      setAuthToken(null);
      setToken(null);
      setStatus('error');
      setError('Нужна авторизация. Перейдите по ссылке из Telegram-бота.');
      return;
    }

    try {
      await finishAuth(storedToken);
    } catch (err) {
      console.error('Stored token is invalid', err);
      storage.removeItem(TOKEN_STORAGE_KEY);
      setAuthToken(null);
      setToken(null);
      setStatus('error');
      setError('Токен недействителен или истёк. Получите новую ссылку в боте.');
    }
  }, [finishAuth]);

  useEffect(() => {
    authenticate();
  }, [authenticate]);

  const refreshProfile = useCallback(async () => {
    try {
      if (!token) return;
      await loadProfile();
    } catch (err) {
      console.error(err);
    }
  }, [loadProfile, token]);

  const value = useMemo<AuthContextValue>(
    () => ({
      token,
      profile,
      status,
      error,
      refreshProfile,
      reauthenticate: authenticate,
      setProfile,
    }),
    [authenticate, error, profile, refreshProfile, status, token]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return ctx;
}
