'use client';

import useSWR from 'swr';
import { useMemo, useState } from 'react';
import { CategoryTabs } from '@/components/CategoryTabs';
import { ProductCard } from '@/components/ProductCard';
import { CartDrawer } from '@/components/CartDrawer';
import { api } from '@/lib/api';
import type { MenuCategory } from '@/types/api';

const fetcher = () => api.getMenu();

export default function HomePage() {
  const { data, isLoading } = useSWR('menu', fetcher);
  const [filter, setFilter] = useState<number | 'all'>('all');

  const categories = useMemo(() => data?.categories ?? [], [data]);

  const filtered: MenuCategory[] = useMemo(() => {
    if (filter === 'all') return categories;
    return categories.filter((cat) => cat.id === filter);
  }, [categories, filter]);

  return (
    <main className="space-y-6">
      <header className="space-y-2">
        <p className="text-sm text-amber-600 font-semibold">KabobFood</p>
        <h1 className="text-2xl font-bold">Выбирайте любимые блюда и заказывайте за пару тапов</h1>
        <p className="text-gray-500 text-sm">Быстрая доставка, понятное меню и очень вкусно.</p>
      </header>

      {isLoading && <p className="text-sm text-gray-500">Загружаем меню...</p>}

      {categories.length > 0 && (
        <CategoryTabs categories={categories} onSelect={setFilter} />
      )}

      <section className="space-y-6">
        {filtered.map((category) => (
          <div key={category.id} className="space-y-3">
            <h2 className="text-xl font-semibold flex items-center gap-2">
              {category.emoji && <span>{category.emoji}</span>}
              {category.name}
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {category.products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>
          </div>
        ))}
      </section>

      <CartDrawer />
    </main>
  );
}
