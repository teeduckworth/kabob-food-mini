export interface TelegramThemeParams {
  bg_color?: string;
  text_color?: string;
  button_color?: string;
  button_text_color?: string;
}

export interface TelegramBackButton {
  show: () => void;
  hide: () => void;
  onClick: (cb: () => void) => void;
  offClick: (cb: () => void) => void;
}

export interface TelegramWebApp {
  ready: () => void;
  expand?: () => void;
  close?: () => void;
  initData?: string;
  themeParams?: TelegramThemeParams;
  onEvent?: (event: string, cb: () => void) => void;
  offEvent?: (event: string, cb: () => void) => void;
  BackButton?: TelegramBackButton;
}

type TelegramLikeGlobal = typeof globalThis & {
  Telegram?: {
    WebApp?: TelegramWebApp;
  };
};

export function getTelegramWebApp(): TelegramWebApp | null {
  const globalRef: TelegramLikeGlobal | undefined =
    typeof globalThis !== 'undefined' ? (globalThis as TelegramLikeGlobal) : undefined;
  if (!globalRef) {
    return null;
  }
  return globalRef.Telegram?.WebApp ?? null;
}
