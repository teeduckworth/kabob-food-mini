import type { MenuResponse, Product } from '@/types/api';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

async function request<T>(path: string, init?: RequestInit, token?: string): Promise<T> {
  const { headers: initHeaders, ...rest } = init || {};
  const headers = new Headers(initHeaders || {});
  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const res = await fetch(`${API_URL}${path}`, {
    ...rest,
    headers,
    cache: 'no-store',
  });

  const text = await res.text();
  if (!res.ok) {
    const message = text || `Request failed with status ${res.status}`;
    throw new Error(message);
  }

  if (!text) {
    return undefined as T;
  }

  return JSON.parse(text) as T;
}

export interface AdminCategoryInput {
  name: string;
  emoji?: string;
  sort_order: number;
  is_active: boolean;
}

export interface AdminCategory {
  id: number;
  name: string;
  emoji?: string;
  sort_order: number;
  is_active: boolean;
}

export interface AdminProductInput {
  category_id: number;
  name: string;
  description?: string;
  price: number;
  old_price?: number;
  image_url?: string;
  is_active: boolean;
  sort_order: number;
}

export const adminApi = {
  login: (username: string, password: string) =>
    request<{ token: string }>(
      '/admin/login',
      {
        method: 'POST',
        body: JSON.stringify({ username, password }),
      }
    ),
  getMenu: () => request<MenuResponse>('/menu'),
  createCategory: (token: string, body: AdminCategoryInput) =>
    request<AdminCategory>(
      '/admin/categories',
      {
        method: 'POST',
        body: JSON.stringify(body),
      },
      token
    ),
  createProduct: (token: string, body: AdminProductInput) =>
    request<Product>(
      '/admin/products',
      {
        method: 'POST',
        body: JSON.stringify(body),
      },
      token
    ),
  updateProduct: (token: string, id: number, body: AdminProductInput) =>
    request<Product>(
      `/admin/products/${id}`,
      {
        method: 'PUT',
        body: JSON.stringify(body),
      },
      token
    ),
};
