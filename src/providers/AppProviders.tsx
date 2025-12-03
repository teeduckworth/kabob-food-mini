'use client';

import type { ReactNode } from 'react';
import { AuthProvider } from './AuthProvider';
import { TelegramProvider } from './TelegramProvider';

export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <TelegramProvider>
      <AuthProvider>{children}</AuthProvider>
    </TelegramProvider>
  );
}
