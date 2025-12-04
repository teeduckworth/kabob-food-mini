import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';

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

export default function RootLayout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="ru">
			<body className={`${geistSans.variable} ${geistMono.variable} antialiased bg-[radial-gradient(circle_at_top,#f8fafc, #e2e8f0)] text-slate-900`}>
				{children}
			</body>
		</html>
	);
}
