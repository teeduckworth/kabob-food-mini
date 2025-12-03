'use client';

import { useState, useEffect } from 'react';
import clsx from 'clsx';
import type { MenuCategory } from '@/types/api';

interface Props {
  categories: MenuCategory[];
  onSelect: (id: number | 'all') => void;
}

export function CategoryTabs({ categories, onSelect }: Props) {
  const [active, setActive] = useState<number | 'all'>('all');

  useEffect(() => {
    onSelect(active);
  }, [active, onSelect]);

  const options: Array<number | 'all'> = ['all', ...categories.map((c) => c.id)];

  return (
    <div className="flex gap-2 overflow-x-auto pb-2">
      {options.map((id) => {
        const label = id === 'all' ? 'Все' : categories.find((c) => c.id === id)?.name;
        return (
          <button
            key={id}
            onClick={() => setActive(id)}
            className={clsx(
              'px-4 py-2 rounded-full border text-sm flex-shrink-0',
              active === id ? 'bg-black text-white border-black' : 'bg-white border-gray-200'
            )}
          >
            {label}
          </button>
        );
      })}
    </div>
  );
}
