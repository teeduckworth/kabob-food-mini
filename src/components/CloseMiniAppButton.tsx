'use client';

import { useState } from 'react';
import WebApp from '@twa-dev/sdk';

type TelegramWindow = Window & {
  Telegram?: {
    WebApp?: unknown;
  };
};

export function CloseMiniAppButton() {
  const [isVisible] = useState(() => {
    if (typeof window === 'undefined') return false;
    const telegramWindow = window as TelegramWindow;
    return Boolean(telegramWindow.Telegram?.WebApp || WebApp?.initData);
  });

  if (!isVisible) return null;

  return (
    <button
      type="button"
      aria-label="Закрыть мини-апп"
      onClick={() => WebApp.close?.()}
      className="absolute right-0 top-0 text-xs uppercase tracking-wide border border-gray-200 px-3 py-1 rounded-full text-gray-500"
    >
      Закрыть
    </button>
  );
}
