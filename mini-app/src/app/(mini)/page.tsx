'use client';

import Link from 'next/link';
import { useEffect, useMemo, useState } from 'react';
import useSWR from 'swr';
import { api } from '@/lib/api';
import { CartDrawer } from '@/components/CartDrawer';
import { useCartStore } from '@/store/cart';
import { PremiumProductCard } from '@/components/PremiumProductCard';
import { ChevronDown, LocationIcon } from '@/components/PremiumIcons';
import type { RegionsResponse } from '@/types/api';

const menuFetcher = () => api.getMenu();
const regionsFetcher = () => api.getRegions();

export default function HomePage() {
  const { data: menuData, isLoading } = useSWR('menu', menuFetcher);
  const { data: regionsData } = useSWR<RegionsResponse>('regions', regionsFetcher);
  const addItem = useCartStore((state) => state.addItem);
  const categories = useMemo(() => menuData?.categories ?? [], [menuData]);
  const [activeCategoryId, setActiveCategoryId] = useState<number | null>(null);

  useEffect(() => {
    if (!categories.length) {
      setActiveCategoryId(null);
      return;
    }
    if (!activeCategoryId || !categories.some((category) => category.id === activeCategoryId)) {
      setActiveCategoryId(categories[0].id);
    }
  }, [categories, activeCategoryId]);

  const activeCategory = categories.find((category) => category.id === activeCategoryId) ?? categories[0];
  const regionName = regionsData?.regions?.[0]?.name ?? 'Центр';
  const heroEmoji = activeCategory?.emoji || '🌭';
  const heroTitle = activeCategory?.name || 'Меню';

  return (
    <div className="min-h-screen bg-[#f7f7f7] px-4 pb-32 pt-6">
      <div className="mx-auto max-w-md space-y-8">
        <header className="flex items-center justify-between gap-4">
          <button
            type="button"
            className="flex flex-1 items-center gap-2 rounded-[24px] bg-white px-5 py-3 text-sm font-medium text-slate-800 shadow-[0_10px_35px_rgba(15,23,42,0.08)]"
          >
            <LocationIcon className="h-5 w-5 text-amber-500" />
            <span className="truncate">{regionName}</span>
            <ChevronDown className="h-4 w-4 text-slate-400" />
          </button>
          <Link
            href="/profile"
            className="group relative flex h-14 w-14 items-center justify-center rounded-full border border-white/40 bg-gradient-to-br from-slate-900 to-slate-700 text-lg font-semibold text-white shadow-[0_15px_40px_rgba(15,23,42,0.35)] transition hover:scale-105"
            aria-label="Мой профиль"
          >
            <span>Я</span>
            <span className="absolute -bottom-1 right-0 h-3 w-3 rounded-full bg-emerald-400 ring-2 ring-slate-900/80" aria-hidden />
          </Link>
        </header>

        <section className="-mx-4 overflow-x-auto px-4">
          <div className="flex gap-4 pb-2">
            {categories.map((category) => {
              const active = category.id === activeCategory?.id;
              return (
                <button
                  key={category.id}
                  type="button"
                  onClick={() => setActiveCategoryId(category.id)}
                  className={`flex flex-shrink-0 items-center gap-2 rounded-full px-5 py-3 text-sm font-semibold transition ${
                    active
                      ? 'bg-slate-900 text-white shadow-[0_10px_25px_rgba(15,23,42,0.25)]'
                      : 'bg-white text-slate-500 shadow-[0_5px_20px_rgba(15,23,42,0.08)]'
                  }`}
                >
                  <span className="text-lg">{category.emoji}</span>
                  <span>{category.name}</span>
                </button>
              );
            })}
          </div>
        </section>

        <section className="flex items-center justify-between">
          <div>
            <p className="text-sm uppercase tracking-[0.3em] text-slate-400">Лучшее сегодня</p>
            <div className="mt-3 flex items-center gap-3">
              <span className="text-3xl" aria-hidden>
                {heroEmoji}
              </span>
              <h1 className="text-3xl font-semibold text-slate-900">{heroTitle}</h1>
            </div>
          </div>
          <div className="rounded-full bg-white/70 px-3 py-1 text-xs text-slate-500 shadow">{categories.length} категорий</div>
        </section>

        {isLoading && (
          <div className="rounded-[28px] bg-white/60 p-6 text-center text-sm text-slate-500 shadow-inner">
            Загружаем меню…
          </div>
        )}

        {!isLoading && !activeCategory && (
          <div className="rounded-[28px] bg-white/80 p-6 text-center text-sm text-slate-500 shadow">
            Меню пока пустое.
          </div>
        )}

        {!isLoading && activeCategory && (
          <section className="grid grid-cols-2 gap-4">
            {activeCategory.products.map((product) => (
              <PremiumProductCard
                key={product.id}
                product={product}
                onAdd={() => addItem(product)}
                accentEmoji={activeCategory.emoji || '🌭'}
              />
            ))}
          </section>
        )}
      </div>

      <CartDrawer />
    </div>
  );
}
