export interface Product {
  id: number;
  category_id: number;
  name: string;
  description?: string;
  price: number;
  old_price?: number;
  image_url?: string;
  is_active: boolean;
  sort_order: number;
}

export interface MenuCategory {
  id: number;
  name: string;
  emoji?: string;
  sort_order: number;
  products: Product[];
}

export interface MenuResponse {
  categories: MenuCategory[];
}

export interface Region {
  id: number;
  name: string;
  delivery_price: number;
  is_active: boolean;
}

export interface RegionsResponse {
  regions: Region[];
}

export interface Address {
  id: number;
  region_id: number;
  street: string;
  house: string;
  entrance?: string;
  flat?: string;
  comment?: string;
  is_default: boolean;
}

export interface AddressInput {
  region_id: number;
  street: string;
  house: string;
  entrance?: string;
  flat?: string;
  comment?: string;
  is_default?: boolean;
}

export interface AddressesResponse {
  addresses: Address[];
}

export interface Profile {
  user: {
    id: number;
    first_name?: string;
    last_name?: string;
    phone?: string;
  };
  addresses: Address[];
}

export interface CreateOrderPayload {
  client_request_id: string;
  type: 'delivery' | 'pickup';
  region_id: number;
  address_id?: number;
  payment_method: string;
  customer_name: string;
  customer_phone: string;
  comment?: string;
  items: Array<{ product_id: number; qty: number }>;
}

export interface Order extends CreateOrderPayload {
  id: number;
  status: string;
  items_total: number;
  total_price: number;
}

export interface AuthResponse {
  token: string;
  profile: Profile;
}
