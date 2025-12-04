import type { ReactNode } from 'react';
import { AdminNav } from './AdminNav';

export const metadata = {
  title: 'KabobFood Admin',
  description: 'Premium control panel for KabobFood operations',
};

export default function AdminLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(253,230,138,0.3),_transparent_60%)]" />
      <div className="relative z-10 mx-auto flex min-h-screen max-w-6xl flex-col px-6 py-10">
        <header className="border-b border-white/10 pb-6">
          <p className="text-sm uppercase tracking-[0.3em] text-amber-200/70">KabobFood Executive</p>
          <div className="mt-2 flex flex-col gap-2 md:flex-row md:items-end md:justify-between">
            <h1 className="text-3xl font-semibold text-white">Панель управления</h1>
            <p className="text-sm text-white/60">Контролируйте меню, регионы и операционные метрики в один клик.</p>
          </div>
          <AdminNav />
        </header>
        <main className="flex-1 py-10">{children}</main>
      </div>
    </div>
  );
}
