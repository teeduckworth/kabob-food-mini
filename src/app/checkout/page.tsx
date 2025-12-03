'use client';

import { useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import useSWR from 'swr';
import { v4 as uuid } from 'uuid';
import { useCartStore } from '@/store/cart';
import { api } from '@/lib/api';
import { useAuth } from '@/providers/AuthProvider';
import type { AddressInput } from '@/types/api';

const regionsFetcher = () => api.getRegions();

type AddressMode = 'saved' | 'new';

interface CheckoutForm {
  name: string;
  phone: string;
  type: 'delivery' | 'pickup';
  street?: string;
  house?: string;
  entrance?: string;
  flat?: string;
  comment?: string;
  payment_method: string;
  region_id?: number;
}

export default function CheckoutPage() {
  const { items, clear, total } = useCartStore();
  const { token, profile, status: authStatus, error: authError, refreshProfile } = useAuth();
  const { data: regionsData, isLoading: loadingRegions } = useSWR('regions', regionsFetcher);
  const regions = useMemo(() => regionsData?.regions ?? [], [regionsData]);
  const addresses = useMemo(() => profile?.addresses ?? [], [profile]);

  const [addrMode, setAddrMode] = useState<AddressMode>(addresses.length > 0 ? 'saved' : 'new');
  const [selectedAddressId, setSelectedAddressId] = useState<number | null>(addresses[0]?.id ?? null);
  const [submitState, setSubmitState] = useState<'idle' | 'submitting' | 'success' | 'error'>('idle');
  const [submitError, setSubmitError] = useState('');

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
  } = useForm<CheckoutForm>({
    shouldUnregister: true,
    defaultValues: {
      type: 'delivery',
      payment_method: 'cash',
      region_id: regions[0]?.id,
    },
  });

  useEffect(() => {
    if (regions.length) {
      setValue('region_id', regions[0].id);
    }
  }, [regions, setValue]);

  useEffect(() => {
    if (addresses.length > 0) {
      setAddrMode('saved');
      setSelectedAddressId(addresses[0].id);
    } else {
      setAddrMode('new');
      setSelectedAddressId(null);
    }
  }, [addresses]);

  const deliveryType = watch('type');

  const fallbackRegionId = useMemo(() => {
    return addresses[0]?.region_id || regions[0]?.id || 1;
  }, [addresses, regions]);

  const onSubmit = handleSubmit(async (values) => {
    if (!token) {
      setSubmitError('Авторизуйтесь в Telegram, чтобы оформить заказ.');
      return;
    }
    if (items.length === 0) return;

    setSubmitState('submitting');
    setSubmitError('');

    try {
      let addressId: number | undefined;
      let regionId = fallbackRegionId;

      if (values.type === 'delivery') {
        if (addrMode === 'saved') {
          if (!selectedAddressId) {
            throw new Error('Выберите сохранённый адрес.');
          }
          addressId = selectedAddressId;
          regionId = addresses.find((addr) => addr.id === selectedAddressId)?.region_id || fallbackRegionId;
        } else {
          if (!values.street || !values.house || !values.region_id) {
            throw new Error('Заполните адрес доставки.');
          }
          const newAddress: AddressInput = {
            region_id: values.region_id,
            street: values.street,
            house: values.house,
            entrance: values.entrance,
            flat: values.flat,
            comment: values.comment,
          };
          const created = await api.createAddress(token, newAddress);
          addressId = created.id;
          regionId = created.region_id;
          await refreshProfile();
        }
      } else {
        regionId = values.region_id || fallbackRegionId;
      }

      const payload = {
        client_request_id: uuid(),
        type: values.type,
        region_id: regionId,
        address_id: addressId,
        payment_method: values.payment_method,
        customer_name: values.name,
        customer_phone: values.phone,
        comment: values.comment,
        items: items.map((item) => ({
          product_id: item.product.id,
          qty: item.qty,
        })),
      };

      await api.createOrder(token, payload);
      setSubmitState('success');
      clear();
    } catch (err) {
      console.error(err);
      setSubmitState('error');
      setSubmitError(err instanceof Error ? err.message : 'Не удалось оформить заказ.');
    }
  });

  if (submitState === 'success') {
    return (
      <main className="space-y-4">
        <h1 className="text-2xl font-semibold">Спасибо!</h1>
        <p className="text-gray-600">Мы уже начинаем готовить ваш заказ. Курьер свяжется перед доставкой.</p>
      </main>
    );
  }

  if (authStatus === 'loading') {
    return (
      <main className="space-y-4">
        <p className="text-sm text-gray-500">Подготавливаем мини-апп…</p>
      </main>
    );
  }

  if (authStatus === 'error') {
    return (
      <main className="space-y-4">
        <p className="text-sm text-red-500">{authError ?? 'Авторизация недоступна'}</p>
      </main>
    );
  }

  return (
    <main className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold">Оформление заказа</h1>
        <p className="text-gray-500 text-sm">Всего к оплате {total().toFixed(0)} ₽</p>
      </div>

      <form className="space-y-4" onSubmit={onSubmit}>
        <div className="space-y-1">
          <label className="text-sm text-gray-500">Имя*</label>
          <input
            className="w-full rounded-xl border p-3"
            {...register('name', { required: 'Укажите имя' })}
            placeholder="Ваше имя"
          />
          {errors.name && <p className="text-xs text-red-500">{errors.name.message}</p>}
        </div>

        <div className="space-y-1">
          <label className="text-sm text-gray-500">Телефон*</label>
          <input
            className="w-full rounded-xl border p-3"
            {...register('phone', { required: 'Укажите телефон' })}
            placeholder="+7 999 123-45-67"
          />
          {errors.phone && <p className="text-xs text-red-500">{errors.phone.message}</p>}
        </div>

        <div className="space-y-1">
          <label className="text-sm text-gray-500">Формат получения</label>
          <select className="w-full rounded-xl border p-3" {...register('type')}>
            <option value="delivery">Доставка</option>
            <option value="pickup">Самовывоз</option>
          </select>
        </div>

        {deliveryType === 'delivery' && (
          <div className="space-y-3">
            {addresses.length > 0 && (
              <div className="flex gap-2 bg-gray-100 rounded-full p-1 text-sm">
                <button
                  type="button"
                  onClick={() => setAddrMode('saved')}
                  className={`flex-1 rounded-full py-2 ${addrMode === 'saved' ? 'bg-white font-semibold' : ''}`}
                >
                  Сохранённые
                </button>
                <button
                  type="button"
                  onClick={() => setAddrMode('new')}
                  className={`flex-1 rounded-full py-2 ${addrMode === 'new' ? 'bg-white font-semibold' : ''}`}
                >
                  Новый
                </button>
              </div>
            )}

            {addrMode === 'saved' && addresses.length > 0 && (
              <div className="space-y-2">
                {addresses.map((addr) => (
                  <label
                    key={addr.id}
                    className={`rounded-2xl border p-3 flex items-start gap-3 ${
                      selectedAddressId === addr.id ? 'border-black' : 'border-gray-200'
                    }`}
                  >
                    <input
                      type="radio"
                      name="saved-address"
                      checked={selectedAddressId === addr.id}
                      onChange={() => setSelectedAddressId(addr.id)}
                    />
                    <div>
                      <p className="font-semibold">
                        {addr.street}, {addr.house}
                        {addr.flat && `, кв. ${addr.flat}`}
                      </p>
                      {addr.comment && <p className="text-sm text-gray-500">{addr.comment}</p>}
                    </div>
                  </label>
                ))}
              </div>
            )}

            {(addrMode === 'new' || addresses.length === 0) && (
              <div className="space-y-4">
                <div className="space-y-1">
                  <label className="text-sm text-gray-500">Регион</label>
                  <select
                    className="w-full rounded-xl border p-3"
                    {...register('region_id', { required: 'Выберите регион', valueAsNumber: true })}
                    disabled={loadingRegions}
                  >
                    {regions.map((region) => (
                      <option key={region.id} value={region.id}>
                        {region.name}
                      </option>
                    ))}
                  </select>
                  {errors.region_id && <p className="text-xs text-red-500">{errors.region_id.message}</p>}
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-1">
                    <label className="text-sm text-gray-500">Улица*</label>
                    <input
                      className="w-full rounded-xl border p-3"
                      {...register('street', { required: 'Укажите улицу' })}
                    />
                    {errors.street && <p className="text-xs text-red-500">{errors.street.message}</p>}
                  </div>
                  <div className="space-y-1">
                    <label className="text-sm text-gray-500">Дом*</label>
                    <input
                      className="w-full rounded-xl border p-3"
                      {...register('house', { required: 'Укажите дом' })}
                    />
                    {errors.house && <p className="text-xs text-red-500">{errors.house.message}</p>}
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-1">
                    <label className="text-sm text-gray-500">Подъезд</label>
                    <input className="w-full rounded-xl border p-3" {...register('entrance')} />
                  </div>
                  <div className="space-y-1">
                    <label className="text-sm text-gray-500">Квартира</label>
                    <input className="w-full rounded-xl border p-3" {...register('flat')} />
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        <div className="space-y-1">
          <label className="text-sm text-gray-500">Оплата</label>
          <select className="w-full rounded-xl border p-3" {...register('payment_method')}>
            <option value="cash">Наличные</option>
            <option value="terminal">Терминал</option>
            <option value="transfer">Перевод</option>
          </select>
        </div>

        <div className="space-y-1">
          <label className="text-sm text-gray-500">Комментарий</label>
          <textarea className="w-full rounded-xl border p-3" rows={3} {...register('comment')} />
        </div>

        {submitError && <p className="text-sm text-red-500">{submitError}</p>}

        <button
          type="submit"
          disabled={submitState === 'submitting' || !token}
          className="w-full rounded-full btn-primary py-3 font-semibold"
        >
          {submitState === 'submitting' ? 'Отправляем…' : 'Подтвердить заказ'}
        </button>
      </form>
    </main>
  );
}
