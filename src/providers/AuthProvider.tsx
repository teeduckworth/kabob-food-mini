'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import WebApp from '@twa-dev/sdk';
import { api } from '@/lib/api';
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

const isBrowser = typeof window !== 'undefined';

function getDevToken(): string | null {
  if (!isBrowser) return null;
  const params = new URLSearchParams(window.location.search);
  return params.get('token') || process.env.NEXT_PUBLIC_DEV_TOKEN || null;
}

function getInitData(): string | null {
  if (!isBrowser) return null;
  if (WebApp?.initData) {
    return WebApp.initData;
  }
  const params = new URLSearchParams(window.location.search);
  return params.get('tgWebAppData') || process.env.NEXT_PUBLIC_DEV_INIT_DATA || null;
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(null);
  const [profile, setProfile] = useState<Profile | null>(null);
  const [status, setStatus] = useState<'loading' | 'ready' | 'error'>('loading');
  const [error, setError] = useState<string | null>(null);

  const loadProfile = useCallback(
    async (jwt: string) => {
      const data = await api.getProfile(jwt);
      setProfile(data);
    },
    []
  );

  const finishAuth = useCallback(
    async (jwt: string, profileData?: Profile) => {
      sessionStorage.setItem(TOKEN_STORAGE_KEY, jwt);
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
    if (!isBrowser) return;
    setStatus('loading');
    setError(null);

    const storedToken = sessionStorage.getItem(TOKEN_STORAGE_KEY);
    if (storedToken) {
      try {
        await finishAuth(storedToken);
        return;
      } catch (err) {
        console.error('Stored token is invalid, re-authenticating', err);
        sessionStorage.removeItem(TOKEN_STORAGE_KEY);
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
    if (!isBrowser) return;
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
