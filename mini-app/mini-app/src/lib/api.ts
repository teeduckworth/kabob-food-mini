import type {
  Address,
  AddressInput,
  AddressesResponse,
  CreateOrderPayload,
  MenuResponse,
  Order,
  Profile,
  RegionsResponse,
} from '@/types/api';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

type RequestOptions = RequestInit & { auth?: boolean };

let authToken: string | null = null;

export function setAuthToken(token: string | null) {
  authToken = token;
}

async function request<T>(path: string, init?: RequestOptions): Promise<T> {
  const { auth, headers: initHeaders, ...rest } = init || {};
  const headers = new Headers(initHeaders ?? {});
  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }
  if (auth) {
    if (!authToken) {
      throw new Error('auth token missing');
    }
    headers.set('Authorization', `Bearer ${authToken}`);
  }

  const res = await fetch(`${API_URL}${path}`, {
    ...rest,
    headers,
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
  getMenu: () => request<MenuResponse>('/menu'),
  getRegions: () => request<RegionsResponse>('/regions'),
  getProfile: () => request<Profile>('/profile', { auth: true }),
  getAddresses: () => request<AddressesResponse>('/addresses', { auth: true }),
  createAddress: (body: AddressInput) =>
    request<Address>('/addresses', {
      method: 'POST',
      body: JSON.stringify(body),
      auth: true,
    }),
  updateAddress: (id: number, body: AddressInput) =>
    request<Address>(`/addresses/${id}`, {
      method: 'PUT',
      body: JSON.stringify(body),
      auth: true,
    }),
  deleteAddress: (id: number) => request<void>(`/addresses/${id}`, { method: 'DELETE', auth: true }),
  createOrder: (body: CreateOrderPayload) =>
    request<Order>('/orders', {
      method: 'POST',
      body: JSON.stringify(body),
      auth: true,
    }),
};

export type ApiClient = typeof api;
