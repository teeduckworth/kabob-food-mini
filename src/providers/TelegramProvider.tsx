'use client';

import { useEffect } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import WebApp from '@twa-dev/sdk';

export function TelegramProvider({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();

  useEffect(() => {
    if (!WebApp || typeof window === 'undefined') return;
    try {
      WebApp.ready();
      WebApp.expand();
    } catch (err) {
      console.warn('Telegram WebApp init skipped', err);
    }

    const applyTheme = () => {
      if (!WebApp?.themeParams) return;
      const root = document.documentElement;
      const { bg_color, text_color, button_color, button_text_color } = WebApp.themeParams;
      if (bg_color) root.style.setProperty('--background', bg_color);
      if (text_color) root.style.setProperty('--foreground', text_color);
      if (button_color) root.style.setProperty('--tg-button', button_color);
      if (button_text_color)
        root.style.setProperty('--tg-button-text', button_text_color);
    };

    applyTheme();

    const handler = () => applyTheme();
    WebApp.onEvent?.('themeChanged', handler);

    return () => {
      WebApp.offEvent?.('themeChanged', handler);
    };
  }, []);

  useEffect(() => {
    if (!WebApp?.BackButton || typeof window === 'undefined') return;

    const handleBack = () => {
      if (pathname === '/') {
        WebApp.close?.();
      } else {
        router.back();
      }
    };

    if (pathname === '/') {
      WebApp.BackButton.hide();
      WebApp.BackButton.offClick(handleBack);
    } else {
      WebApp.BackButton.show();
      WebApp.BackButton.onClick(handleBack);
    }

    return () => {
      WebApp.BackButton.offClick(handleBack);
    };
  }, [pathname, router]);

  return <>{children}</>;
}
