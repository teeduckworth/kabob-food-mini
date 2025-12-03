import type {
  Address,
  AddressInput,
  AddressesResponse,
  AuthResponse,
  CreateOrderPayload,
  MenuResponse,
  Order,
  Profile,
  RegionsResponse,
} from '@/types/api';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers || {}),
    },
    ...init,
    cache: 'no-store',
  });

  if (!res.ok) {
    throw new Error(`API error: ${res.status}`);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  const text = await res.text();
  if (!text) {
    return undefined as T;
  }

  return JSON.parse(text) as T;
}

export const api = {
  authTelegram: (initData: string) =>
    request<AuthResponse>('/auth/telegram', {
      method: 'POST',
      body: JSON.stringify({ init_data: initData }),
    }),
  getMenu: () => request<MenuResponse>('/menu'),
  getRegions: () => request<RegionsResponse>('/regions'),
  getProfile: (token: string) =>
    request<Profile>('/profile', {
      headers: { Authorization: `Bearer ${token}` },
    }),
  getAddresses: (token: string) =>
    request<AddressesResponse>('/addresses', {
      headers: { Authorization: `Bearer ${token}` },
    }),
  createAddress: (token: string, body: AddressInput) =>
    request<Address>('/addresses', {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(body),
    }),
  updateAddress: (token: string, id: number, body: AddressInput) =>
    request<Address>(`/addresses/${id}`, {
      method: 'PUT',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(body),
    }),
  deleteAddress: (token: string, id: number) =>
    request<void>(`/addresses/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    }),
  createOrder: (token: string, body: CreateOrderPayload) =>
    request<Order>('/orders', {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify(body),
    }),
};

export type ApiClient = typeof api;
