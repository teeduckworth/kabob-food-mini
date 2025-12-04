'use client';

import { useCallback, useEffect, useState } from 'react';

const STORAGE_KEY = 'kabobfood_admin_token';

export function useAdminSession() {
  const [token, setTokenState] = useState<string | null>(null);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    const saved = typeof window !== 'undefined' ? window.localStorage.getItem(STORAGE_KEY) : null;
    if (saved) {
      setTokenState(saved);
    }
    setReady(true);
  }, []);

  const setToken = useCallback((value: string | null) => {
    setTokenState(value);
    if (typeof window === 'undefined') return;
    if (value) {
      window.localStorage.setItem(STORAGE_KEY, value);
    } else {
      window.localStorage.removeItem(STORAGE_KEY);
    }
  }, []);

  const logout = useCallback(() => {
    setToken(null);
  }, [setToken]);

  return { token, setToken, logout, ready };
}
