'use client';

import { useMemo, useState } from 'react';
import useSWR from 'swr';
import { api } from '@/lib/api';
import { AddressForm } from '@/components/AddressForm';
import { useAuth } from '@/providers/AuthProvider';
import type { Address, AddressInput } from '@/types/api';

const regionsFetcher = () => api.getRegions();

export default function ProfilePage() {
  const { profile, status, error, token, refreshProfile } = useAuth();
  const { data: regionsData, isLoading: loadingRegions } = useSWR('regions', regionsFetcher);
  const regions = regionsData?.regions ?? [];

  const [formVisible, setFormVisible] = useState(false);
  const [editing, setEditing] = useState<Address | null>(null);
  const [mutationError, setMutationError] = useState<string | null>(null);
  const [mutating, setMutating] = useState(false);

  const addresses = useMemo(() => profile?.addresses ?? [], [profile]);

  const sortedAddresses = useMemo(() => {
    return [...addresses].sort((a, b) => Number(b.is_default) - Number(a.is_default));
  }, [addresses]);

  const canMutate = Boolean(token);

  const closeForm = () => {
    setFormVisible(false);
    setEditing(null);
    setMutationError(null);
  };

  const toInput = (addr: Address): AddressInput => ({
    region_id: addr.region_id,
    street: addr.street,
    house: addr.house,
    entrance: addr.entrance,
    flat: addr.flat,
    comment: addr.comment,
    is_default: addr.is_default,
  });

  const saveAddress = async (values: AddressInput) => {
    if (!token) return;
    setMutating(true);
    setMutationError(null);
    try {
      if (editing) {
        await api.updateAddress(editing.id, values);
      } else {
        await api.createAddress(values);
      }
      await refreshProfile();
      closeForm();
    } catch (err) {
      console.error(err);
      setMutationError('Не удалось сохранить адрес. Попробуйте ещё раз.');
    } finally {
      setMutating(false);
    }
  };

  const deleteAddress = async (addr: Address) => {
    if (!token) return;
    setMutating(true);
    setMutationError(null);
    try {
      await api.deleteAddress(addr.id);
      await refreshProfile();
    } catch (err) {
      console.error(err);
      setMutationError('Не удалось удалить адрес.');
    } finally {
      setMutating(false);
    }
  };

  const setDefault = async (addr: Address) => {
    if (!token) return;
    setMutating(true);
    setMutationError(null);
    try {
      await api.updateAddress(addr.id, { ...toInput(addr), is_default: true });
      await refreshProfile();
    } catch (err) {
      console.error(err);
      setMutationError('Не удалось обновить адрес.');
    } finally {
      setMutating(false);
    }
  };

  if (status === 'loading') {
    return (
      <main className="space-y-4">
        <p className="text-sm text-gray-500">Загружаем профиль…</p>
      </main>
    );
  }

  if (status === 'error') {
    return (
      <main className="space-y-4">
        <p className="text-sm text-red-500">{error ?? 'Не удалось загрузить профиль'}</p>
      </main>
    );
  }

  if (!profile || !token) {
    return (
      <main className="space-y-4">
        <p className="text-sm text-gray-500">Профиль появится после перехода по ссылке с токеном из бота.</p>
      </main>
    );
  }

  return (
    <main className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold">Профиль</h1>
        <p className="text-sm text-gray-500">Данные синхронизированы с информацией, которую вы отправили боту.</p>
      </div>

      <section className="space-y-3">
        <div className="rounded-2xl border p-4 bg-white">
          <p className="text-sm text-gray-500">Имя</p>
          <p className="font-semibold">
            {[profile.user.first_name, profile.user.last_name].filter(Boolean).join(' ') || '—'}
          </p>
        </div>
        <div className="rounded-2xl border p-4 bg-white">
          <p className="text-sm text-gray-500">Телефон</p>
          <p className="font-semibold">{profile.user.phone || '—'}</p>
        </div>
      </section>

      <section className="space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-semibold">Адреса доставки</h2>
          <button
            className="text-sm font-semibold text-amber-600"
            onClick={() => {
              setEditing(null);
              setFormVisible(true);
            }}
            disabled={!canMutate || mutating}
          >
            + Добавить
          </button>
        </div>

        {sortedAddresses.length === 0 && (
          <p className="text-sm text-gray-500">Сохранённых адресов пока нет.</p>
        )}

        <div className="space-y-3">
          {sortedAddresses.map((addr) => (
            <div key={addr.id} className="rounded-2xl border p-4 bg-white space-y-2">
              <div className="flex items-center justify-between">
                <p className="font-semibold">
                  {addr.street}, {addr.house}
                  {addr.flat && `, кв. ${addr.flat}`}
                </p>
                {addr.is_default && <span className="text-xs px-2 py-1 rounded-full bg-emerald-100 text-emerald-700">По умолчанию</span>}
              </div>
              {addr.comment && <p className="text-sm text-gray-500">{addr.comment}</p>}
              <div className="flex gap-3 text-sm text-gray-500">
                <button
                  onClick={() => {
                    setEditing(addr);
                    setFormVisible(true);
                  }}
                  disabled={mutating}
                >
                  Изменить
                </button>
                <button onClick={() => deleteAddress(addr)} disabled={mutating}>
                  Удалить
                </button>
                {!addr.is_default && (
                  <button onClick={() => setDefault(addr)} disabled={mutating}>
                    Сделать основным
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
        {mutationError && !(formVisible || editing) && (
          <p className="text-sm text-red-500">{mutationError}</p>
        )}
      </section>

      {(formVisible || editing) && (
        <section className="rounded-2xl border p-4 bg-white">
          {loadingRegions && <p className="text-sm text-gray-500">Загружаем регионы…</p>}
          {regions.length === 0 && !loadingRegions ? (
            <div className="space-y-3">
              <p className="text-sm text-red-500">Нет регионов доставки. Создайте их через админку.</p>
              <button
                className="w-full rounded-full border py-3 text-sm"
                type="button"
                onClick={closeForm}
              >
                Вернуться
              </button>
            </div>
          ) : (
            <AddressForm
              regions={regions}
              defaultValues={editing ? toInput(editing) : undefined}
              onSubmit={saveAddress}
              onCancel={closeForm}
              isSubmitting={mutating}
              title={editing ? 'Редактировать адрес' : 'Новый адрес'}
              error={mutationError}
            />
          )}
        </section>
      )}
    </main>
  );
}
