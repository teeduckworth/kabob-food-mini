'use client';

import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import type { AddressInput, Region } from '@/types/api';

interface AddressFormProps {
  regions: Region[];
  defaultValues?: AddressInput;
  onSubmit: (values: AddressInput) => Promise<void>;
  onCancel: () => void;
  isSubmitting: boolean;
  title: string;
  error?: string | null;
}

export function AddressForm({
  regions,
  defaultValues,
  onSubmit,
  onCancel,
  isSubmitting,
  title,
  error,
}: AddressFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<AddressInput>({
    defaultValues: {
      region_id: defaultValues?.region_id || regions[0]?.id || 0,
      street: defaultValues?.street || '',
      house: defaultValues?.house || '',
      entrance: defaultValues?.entrance || '',
      flat: defaultValues?.flat || '',
      comment: defaultValues?.comment || '',
      is_default: defaultValues?.is_default ?? false,
    },
  });

  useEffect(() => {
    reset({
      region_id: defaultValues?.region_id || regions[0]?.id || 0,
      street: defaultValues?.street || '',
      house: defaultValues?.house || '',
      entrance: defaultValues?.entrance || '',
      flat: defaultValues?.flat || '',
      comment: defaultValues?.comment || '',
      is_default: defaultValues?.is_default ?? false,
    });
  }, [defaultValues, regions, reset]);

  const submit = handleSubmit(async (values) => {
    await onSubmit({ ...values, region_id: Number(values.region_id) });
  });

  return (
    <form className="space-y-4" onSubmit={submit}>
      <div>
        <h2 className="text-lg font-semibold">{title}</h2>
        <p className="text-sm text-gray-500">Укажите адрес, который будем предлагать при оформлении заказа.</p>
      </div>

      <div className="space-y-1">
        <label className="text-sm text-gray-500">Регион</label>
        <select
          className="w-full rounded-xl border p-3"
          {...register('region_id', { required: 'Выберите регион', valueAsNumber: true })}
        >
          {regions.map((region) => (
            <option key={region.id} value={region.id}>
              {region.name}
            </option>
          ))}
        </select>
        {errors.region_id && <p className="text-xs text-red-500">{errors.region_id.message}</p>}
        {!regions.length && (
          <p className="text-xs text-amber-600">Сначала настройте регионы доставки.</p>
        )}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-1">
          <label className="text-sm text-gray-500">Улица*</label>
          <input
            className="w-full rounded-xl border p-3"
            {...register('street', { required: 'Укажите улицу' })}
            placeholder="Тверская"
          />
          {errors.street && <p className="text-xs text-red-500">{errors.street.message}</p>}
        </div>
        <div className="space-y-1">
          <label className="text-sm text-gray-500">Дом*</label>
          <input
            className="w-full rounded-xl border p-3"
            {...register('house', { required: 'Укажите дом' })}
            placeholder="15"
          />
          {errors.house && <p className="text-xs text-red-500">{errors.house.message}</p>}
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-1">
          <label className="text-sm text-gray-500">Подъезд</label>
          <input className="w-full rounded-xl border p-3" {...register('entrance')} placeholder="1" />
        </div>
        <div className="space-y-1">
          <label className="text-sm text-gray-500">Квартира</label>
          <input className="w-full rounded-xl border p-3" {...register('flat')} placeholder="45" />
        </div>
      </div>

      <div className="space-y-1">
        <label className="text-sm text-gray-500">Комментарий</label>
        <textarea className="w-full rounded-xl border p-3" rows={3} {...register('comment')} />
      </div>

      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" {...register('is_default')} />
        По умолчанию
      </label>

      {error && <p className="text-sm text-red-500">{error}</p>}

      <div className="flex gap-3">
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 rounded-full border border-gray-200 py-3 font-medium"
          disabled={isSubmitting}
        >
          Отмена
        </button>
        <button
          type="submit"
          disabled={isSubmitting || !regions.length}
          className="flex-1 rounded-full btn-primary py-3 font-semibold"
        >
          {isSubmitting ? 'Сохраняем…' : 'Сохранить'}
        </button>
      </div>
    </form>
  );
}
