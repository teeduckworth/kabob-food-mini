'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { api } from '@/lib/api';
import { getTelegramWebApp } from '@/lib/telegram';
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
const TOKEN_STORAGE_KEY = 'kabobfood_token';

type BrowserWindow = Window & typeof globalThis;

type MaybeBrowser = typeof globalThis & { window?: BrowserWindow };

function getBrowserWindow(): BrowserWindow | null {
  if (typeof globalThis === 'undefined') return null;
  const maybeBrowser = globalThis as MaybeBrowser;
  return maybeBrowser.window ?? null;
}

function getSessionStorage(): Storage | null {
  const win = getBrowserWindow();
  return win?.sessionStorage ?? null;
}

function getQueryParam(name: string): string | null {
  const win = getBrowserWindow();
  if (!win) return null;
  const search = win.location?.search ?? '';
  if (!search) return null;
  return new URLSearchParams(search).get(name);
}

function getDevToken(): string | null {
  return getQueryParam('token') || process.env.NEXT_PUBLIC_DEV_TOKEN || null;
}

function getInitData(): string | null {
  const telegram = getTelegramWebApp();
  if (telegram?.initData) {
    return telegram.initData;
  }
  return getQueryParam('tgWebAppData') || process.env.NEXT_PUBLIC_DEV_INIT_DATA || null;
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(null);
  const [profile, setProfile] = useState<Profile | null>(null);
  const [status, setStatus] = useState<'loading' | 'ready' | 'error'>('loading');
  const [error, setError] = useState<string | null>(null);

  const loadProfile = useCallback(async (jwt: string) => {
    const data = await api.getProfile(jwt);
    setProfile(data);
  }, []);

  const finishAuth = useCallback(
    async (jwt: string, profileData?: Profile) => {
      getSessionStorage()?.setItem(TOKEN_STORAGE_KEY, jwt);
      setToken(jwt);
      if (profileData) {
        setProfile(profileData);
      } else {
        await loadProfile(jwt);
      }
      setStatus('ready');
      setError(null);
    },
    [loadProfile]
  );

  const authenticate = useCallback(async () => {
    const storage = getSessionStorage();
    if (!storage) {
      setStatus('error');
      setError('Нет доступа к sessionStorage (мини-апп нужно открыть в Telegram).');
      return;
    }

    setStatus('loading');
    setError(null);

    const storedToken = storage.getItem(TOKEN_STORAGE_KEY);
    if (storedToken) {
      try {
        await finishAuth(storedToken);
        return;
      } catch (err) {
        console.error('Stored token is invalid, re-authenticating', err);
        storage.removeItem(TOKEN_STORAGE_KEY);
      }
    }

    try {
      const directToken = getDevToken();
      if (directToken) {
        await finishAuth(directToken);
        return;
      }

      const initData = getInitData();
      if (!initData) {
        throw new Error('initData missing');
      }

      const { token: jwt, profile: profileData } = await api.authTelegram(initData);
      await finishAuth(jwt, profileData);
    } catch (err) {
      console.error(err);
      setStatus('error');
      setError('Не удалось авторизоваться. Откройте мини-апп внутри Telegram.');
    }
  }, [finishAuth]);

  useEffect(() => {
    authenticate();
  }, [authenticate]);

  const refreshProfile = useCallback(async () => {
    if (!token) return;
    try {
      await loadProfile(token);
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
