'use client';

import { useEffect } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { getTelegramWebApp } from '@/lib/telegram';

export function TelegramProvider({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();

  useEffect(() => {
    const webApp = getTelegramWebApp();
    if (!webApp) return;
    try {
      webApp.ready();
      webApp.expand?.();
    } catch (err) {
      console.warn('Telegram WebApp init skipped', err);
    }

    const applyTheme = () => {
      if (!webApp.themeParams) return;
      const root = document.documentElement;
      const { bg_color, text_color, button_color, button_text_color } = webApp.themeParams;
      if (bg_color) root.style.setProperty('--background', bg_color);
      if (text_color) root.style.setProperty('--foreground', text_color);
      if (button_color) root.style.setProperty('--tg-button', button_color);
      if (button_text_color)
        root.style.setProperty('--tg-button-text', button_text_color);
    };

    applyTheme();

    const handler = () => applyTheme();
    webApp.onEvent?.('themeChanged', handler);

    return () => {
      webApp.offEvent?.('themeChanged', handler);
    };
  }, []);

  useEffect(() => {
    const webApp = getTelegramWebApp();
    if (!webApp?.BackButton) return;

    const handleBack = () => {
      if (pathname === '/') {
        webApp.close?.();
      } else {
        router.back();
      }
    };

    if (pathname === '/') {
      webApp.BackButton.hide();
      webApp.BackButton.offClick(handleBack);
    } else {
      webApp.BackButton.show();
      webApp.BackButton.onClick(handleBack);
    }

    return () => {
      webApp.BackButton?.offClick(handleBack);
    };
  }, [pathname, router]);

  return <>{children}</>;
}
