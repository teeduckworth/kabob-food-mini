'use client';

import Image from 'next/image';
import { useCartStore } from '@/store/cart';
import type { Product } from '@/types/api';

interface Props {
  product: Product;
}

export function ProductCard({ product }: Props) {
  const addItem = useCartStore((state) => state.addItem);

  return (
    <div className="rounded-2xl border p-4 shadow-sm flex flex-col gap-2 bg-white">
      {product.image_url && (
        <div className="relative w-full h-40 overflow-hidden rounded-xl">
          <Image src={product.image_url} alt={product.name} fill className="object-cover" />
        </div>
      )}
      <div className="flex-1">
        <h3 className="font-semibold text-lg">{product.name}</h3>
        {product.description && (
          <p className="text-sm text-gray-500 line-clamp-2">{product.description}</p>
        )}
      </div>
      <div className="flex items-center justify-between">
        <div>
          <p className="font-bold text-lg">{product.price.toFixed(0)} ₽</p>
          {product.old_price && (
            <p className="text-xs text-gray-400 line-through">{product.old_price.toFixed(0)} ₽</p>
          )}
        </div>
        <button onClick={() => addItem(product)} className="px-4 py-2 rounded-full btn-primary text-sm">
          Добавить
        </button>
      </div>
    </div>
  );
}
