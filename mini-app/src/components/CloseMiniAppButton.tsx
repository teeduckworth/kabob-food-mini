'use client';

import { useEffect, useState } from 'react';
import { getTelegramWebApp } from '@/lib/telegram';

export function CloseMiniAppButton() {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    setIsVisible(Boolean(getTelegramWebApp()));
  }, []);

  if (!isVisible) return null;

  const handleClose = () => {
    getTelegramWebApp()?.close?.();
  };

  return (
    <button
      type="button"
      aria-label="Закрыть мини-апп"
      onClick={handleClose}
      className="absolute right-0 top-0 text-xs uppercase tracking-wide border border-gray-200 px-3 py-1 rounded-full text-gray-500"
    >
      Закрыть
    </button>
  );
}
