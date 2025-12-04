'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

const links = [
  { href: '/admin', label: 'Вход' },
  { href: '/admin/dashboard', label: 'Панель' },
];

export function AdminNav() {
  const pathname = usePathname();

  return (
    <nav className="mt-8 flex flex-wrap gap-3 text-sm">
      {links.map((link) => {
        const active = pathname === link.href || pathname.startsWith(`${link.href}/`);
        return (
          <Link
            key={link.href}
            href={link.href}
            className={`rounded-full border px-4 py-2 transition ${
              active
                ? 'border-white/60 bg-white/10 text-white'
                : 'border-white/20 text-white/70 hover:border-white/50 hover:text-white'
            }`}
          >
            {link.label}
          </Link>
        );
      })}
    </nav>
  );
}
