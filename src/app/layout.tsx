import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';
import { BottomNav } from '@/components/BottomNav';
import { CloseMiniAppButton } from '@/components/CloseMiniAppButton';
import { AppProviders } from '@/providers/AppProviders';

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
});

export const metadata: Metadata = {
  title: 'KabobFood Mini App',
  description: 'Order delicious kebabs inside Telegram',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ru">
      <body className={`${geistSans.variable} ${geistMono.variable} bg-zinc-50 min-h-screen pb-24`}>
        <AppProviders>
          <div className="relative max-w-2xl mx-auto px-4 pb-32 pt-6">
            <CloseMiniAppButton />
            {children}
          </div>

          <BottomNav />
        </AppProviders>
      </body>
    </html>
  );
}
