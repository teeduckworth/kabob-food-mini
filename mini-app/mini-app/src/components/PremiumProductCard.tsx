'use client';

import Image from 'next/image';
import type { Product } from '@/types/api';

interface Props {
  product: Product;
  onAdd: () => void;
  accentEmoji?: string;
}

export function PremiumProductCard({ product, onAdd, accentEmoji = '🌭' }: Props) {
  return (
    <div className="flex flex-col gap-4 rounded-[28px] bg-[#fcfcfc] p-4 shadow-[0_20px_45px_rgba(15,23,42,0.08)]">
      <div
        className="relative w-full overflow-hidden rounded-3xl bg-gradient-to-br from-slate-100 to-slate-200"
        style={{ aspectRatio: '4/3' }}
      >
        {product.image_url ? (
          <Image src={product.image_url} alt={product.name} fill className="object-cover" />
        ) : (
          <div className="flex h-full items-center justify-center text-4xl text-slate-300">
            {accentEmoji}
          </div>
        )}
      </div>
      <div className="space-y-1">
        <p className="text-2xl font-semibold text-slate-900">{product.price.toLocaleString('ru-RU')} ₽</p>
        <p className="text-sm text-slate-500">
          <span className="mr-1">{accentEmoji}</span>
          {product.name}
        </p>
      </div>
      <button
        type="button"
        onClick={onAdd}
        className="flex items-center justify-center gap-2 rounded-full bg-white py-3 text-sm font-semibold text-slate-900 shadow-[0_10px_30px_rgba(15,23,42,0.12)]"
      >
        <PlusIcon className="h-4 w-4" /> Добавить
      </button>
    </div>
  );
}

function PlusIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className={className}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M12 5v14M5 12h14" />
    </svg>
  );
}
