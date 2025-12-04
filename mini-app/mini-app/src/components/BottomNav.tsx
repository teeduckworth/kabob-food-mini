'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import clsx from 'clsx';

const links = [
  { href: '/', label: 'Меню' },
  { href: '/checkout', label: 'Оформление' },
  { href: '/profile', label: 'Профиль' },
];

export function BottomNav() {
  const pathname = usePathname();
  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-white border-t shadow-lg py-3">
      <div className="max-w-2xl mx-auto px-6 flex justify-around">
        {links.map((link) => (
          <Link
            key={link.href}
            href={link.href}
            className={clsx(
              'text-sm font-medium',
              pathname === link.href ? 'text-black' : 'text-gray-400'
            )}
          >
            {link.label}
          </Link>
        ))}
      </div>
    </nav>
  );
}
