'use client';

import Link from 'next/link';
import { useMemo } from 'react';
import { useCartStore } from '@/store/cart';

export function CartDrawer() {
  const items = useCartStore((state) => state.items);
  const total = useCartStore((state) => state.total);
  const totalItems = useMemo(() => items.reduce((sum, it) => sum + it.qty, 0), [items]);

  if (items.length === 0) return null;

  return (
    <div className="fixed bottom-20 left-1/2 -translate-x-1/2 w-[90%] max-w-md bg-white border shadow-lg rounded-full px-4 py-3 flex items-center justify-between">
      <div>
        <p className="text-sm text-gray-500">{totalItems} товаров</p>
        <p className="font-semibold">Итого {total().toFixed(0)} ₽</p>
      </div>
      <Link href="/checkout" className="btn-primary rounded-full px-4 py-2 text-sm">
        Оформить
      </Link>
    </div>
  );
}
