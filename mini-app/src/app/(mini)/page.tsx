'use client';

import Link from 'next/link';
import { useEffect, useMemo, useState } from 'react';
import useSWR from 'swr';
import { api } from '@/lib/api';
import { CartDrawer } from '@/components/CartDrawer';
import { AddressForm } from '@/components/AddressForm';
import { useAuth } from '@/providers/AuthProvider';
import { useCartStore } from '@/store/cart';
import { PremiumProductCard } from '@/components/PremiumProductCard';
import { ChevronDown, LocationIcon } from '@/components/PremiumIcons';
import type { Address, AddressInput, RegionsResponse } from '@/types/api';

const menuFetcher = () => api.getMenu();
const regionsFetcher = () => api.getRegions();
const LAST_ADDRESS_STORAGE_KEY = 'selected_address_id';

function getLocalStorage(): Storage | null {
  if (typeof window === 'undefined') return null;
  return window.localStorage ?? null;
}

export default function HomePage() {
  const { data: menuData, isLoading } = useSWR('menu', menuFetcher);
  const { data: regionsData } = useSWR<RegionsResponse>('regions', regionsFetcher);
  const { profile, token, refreshProfile } = useAuth();
  const addItem = useCartStore((state) => state.addItem);
  const categories = useMemo(() => menuData?.categories ?? [], [menuData]);
  const [activeCategoryId, setActiveCategoryId] = useState<number | null>(null);
  const addresses = useMemo(() => profile?.addresses ?? [], [profile]);
  const [selectedAddressId, setSelectedAddressId] = useState<number | null>(null);
  const [addressSheetOpen, setAddressSheetOpen] = useState(false);
  const [showAddressForm, setShowAddressForm] = useState(false);
  const [addressError, setAddressError] = useState<string | null>(null);
  const [savingAddress, setSavingAddress] = useState(false);

  useEffect(() => {
    if (!categories.length) {
      setActiveCategoryId(null);
      return;
    }
    if (!activeCategoryId || !categories.some((category) => category.id === activeCategoryId)) {
      setActiveCategoryId(categories[0].id);
    }
  }, [categories, activeCategoryId]);

  useEffect(() => {
    if (!addresses.length) {
      setSelectedAddressId(null);
      return;
    }
    const storage = getLocalStorage();
    const storedId = storage?.getItem(LAST_ADDRESS_STORAGE_KEY);
    const parsedId = storedId ? Number(storedId) : null;
    if (parsedId && addresses.some((addr) => addr.id === parsedId)) {
      setSelectedAddressId(parsedId);
      return;
    }
    if (!selectedAddressId || !addresses.some((addr) => addr.id === selectedAddressId)) {
      setSelectedAddressId(addresses[0].id);
    }
  }, [addresses, selectedAddressId]);

  const activeCategory = categories.find((category) => category.id === activeCategoryId) ?? categories[0];
  const fallbackRegionName = regionsData?.regions?.[0]?.name ?? 'Центр';
  const registrationLocation = useMemo(() => {
    const lat = profile?.user.latitude;
    const lon = profile?.user.longitude;
    if (typeof lat === 'number' && typeof lon === 'number') {
      return `${lat.toFixed(3)}, ${lon.toFixed(3)}`;
    }
    return fallbackRegionName;
  }, [fallbackRegionName, profile?.user.latitude, profile?.user.longitude]);
  const activeAddress = addresses.find((addr) => addr.id === selectedAddressId) ?? addresses[0];
  const locationLabel = activeAddress ? `${activeAddress.street}, ${activeAddress.house}` : registrationLocation;
  const locationHint = activeAddress ? 'Сохранённый адрес' : 'Локация при регистрации';
  const heroEmoji = activeCategory?.emoji || '🌭';
  const heroTitle = activeCategory?.name || 'Меню';

  const handleAddressSelect = (addr: Address) => {
    setSelectedAddressId(addr.id);
    getLocalStorage()?.setItem(LAST_ADDRESS_STORAGE_KEY, String(addr.id));
    setAddressSheetOpen(false);
  };

  const handleCreateAddress = async (values: AddressInput) => {
    if (!token) {
      setAddressError('Получите ссылку с токеном в боте, чтобы сохранять адреса.');
      return;
    }
    setSavingAddress(true);
    setAddressError(null);
    try {
      const created = await api.createAddress(values);
      await refreshProfile();
      setSelectedAddressId(created.id);
      getLocalStorage()?.setItem(LAST_ADDRESS_STORAGE_KEY, String(created.id));
      setShowAddressForm(false);
      setAddressSheetOpen(false);
    } catch (err) {
      console.error(err);
      setAddressError('Не удалось сохранить адрес. Попробуйте позже.');
    } finally {
      setSavingAddress(false);
    }
  };

  const renderAddressSheet = () => {
    if (!addressSheetOpen) return null;
    return (
      <div
        className="fixed inset-0 z-50 flex flex-col justify-end bg-black/40 px-4 pb-6 pt-20"
        onClick={() => {
          setAddressSheetOpen(false);
          setShowAddressForm(false);
        }}
      >
        <div
          className="mx-auto w-full max-w-md rounded-3xl bg-white p-4 shadow-2xl"
          onClick={(event) => event.stopPropagation()}
        >
          {!showAddressForm ? (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-xs uppercase tracking-wide text-slate-400">Доставка</p>
                  <p className="text-base font-semibold text-slate-800">Выберите адрес</p>
                </div>
                <button className="text-sm text-amber-600" onClick={() => setAddressSheetOpen(false)}>
                  Закрыть
                </button>
              </div>

              {addresses.length === 0 && (
                <p className="text-sm text-slate-500">
                  Сохранённых адресов нет.
                  {!token && ' Нужна авторизация через бота, чтобы добавить адрес.'}
                </p>
              )}

              <div className="space-y-2">
                {addresses.map((addr) => {
                  const active = addr.id === selectedAddressId;
                  return (
                    <button
                      key={addr.id}
                      onClick={() => handleAddressSelect(addr)}
                      className={`w-full rounded-2xl border px-4 py-3 text-left text-sm transition ${
                        active ? 'border-slate-900 bg-slate-900/5' : 'border-slate-200'
                      }`}
                    >
                      <p className="font-semibold text-slate-900">
                        {addr.street}, {addr.house}
                      </p>
                      <p className="text-xs text-slate-500">
                        {addr.flat && `кв. ${addr.flat} · `}
                        {addr.comment || 'Без комментария'}
                      </p>
                    </button>
                  );
                })}
              </div>

              <button
                type="button"
                onClick={() => {
                  if (!token) {
                    setAddressError('Сначала пройдите регистрацию через бота.');
                    return;
                  }
                  setAddressError(null);
                  setShowAddressForm(true);
                }}
                className="w-full rounded-2xl border border-dashed border-slate-300 py-3 text-sm font-semibold text-slate-700"
              >
                + Добавить адрес
              </button>

              {addressError && <p className="text-sm text-red-500">{addressError}</p>}
            </div>
          ) : (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <p className="text-base font-semibold">Новый адрес</p>
                <button
                  className="text-sm text-slate-500"
                  onClick={() => {
                    setShowAddressForm(false);
                    setAddressError(null);
                  }}
                >
                  Назад
                </button>
              </div>
              <AddressForm
                regions={regionsData?.regions ?? []}
                onSubmit={handleCreateAddress}
                onCancel={() => {
                  setShowAddressForm(false);
                  setAddressError(null);
                }}
                isSubmitting={savingAddress}
                title="Добавить адрес"
                error={addressError}
              />
            </div>
          )}
        </div>
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-[#f7f7f7] px-4 pb-32 pt-6">
      <div className="mx-auto max-w-md space-y-6">
        <header className="flex items-center gap-3">
          <button
            type="button"
            onClick={() => {
              setAddressSheetOpen(true);
              setShowAddressForm(false);
            }}
            className="flex flex-1 items-center gap-3 rounded-[32px] bg-white px-5 py-3 text-left shadow-[0_8px_25px_rgba(15,23,42,0.08)]"
          >
            <span className="flex h-10 w-10 items-center justify-center rounded-full bg-amber-50 text-amber-600">
              <LocationIcon className="h-5 w-5" />
            </span>
            <div className="flex-1 overflow-hidden">
              <p className="text-xs uppercase tracking-wide text-slate-400">{locationHint}</p>
              <p className="truncate text-sm font-semibold text-slate-800">{locationLabel}</p>
            </div>
            <ChevronDown className="h-4 w-4 text-slate-400" />
          </button>
          <Link
            href="/profile"
            className="group relative flex h-12 w-12 items-center justify-center rounded-2xl border border-white/50 bg-gradient-to-br from-slate-900 to-slate-700 text-base font-semibold text-white shadow-[0_12px_30px_rgba(15,23,42,0.25)] transition hover:scale-105"
            aria-label="Мой профиль"
          >
            <span>Я</span>
            <span className="absolute -bottom-1 right-0 h-2.5 w-2.5 rounded-full bg-emerald-400 ring-2 ring-slate-900/80" aria-hidden />
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
      {renderAddressSheet()}
    </div>
  );
}
