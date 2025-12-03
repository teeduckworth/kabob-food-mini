CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    first_name TEXT,
    last_name TEXT,
    username TEXT,
    phone TEXT,
    language TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS regions (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    delivery_price NUMERIC(10, 2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    region_id BIGINT REFERENCES regions(id),
    street TEXT NOT NULL,
    house TEXT NOT NULL,
    entrance TEXT,
    flat TEXT,
    comment TEXT,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    emoji TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    old_price NUMERIC(10, 2),
    image_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    client_request_id UUID NOT NULL UNIQUE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    address_id BIGINT REFERENCES addresses(id),
    type TEXT NOT NULL CHECK (type IN ('delivery', 'pickup')),
    payment_method TEXT NOT NULL,
    status TEXT NOT NULL,
    region_id BIGINT REFERENCES regions(id),
    delivery_price NUMERIC(10, 2) NOT NULL DEFAULT 0,
    items_total NUMERIC(10, 2) NOT NULL DEFAULT 0,
    total_price NUMERIC(10, 2) NOT NULL DEFAULT 0,
    comment TEXT,
    customer_name TEXT NOT NULL,
    customer_phone TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS orders_user_id_idx ON orders(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS orders_client_request_id_idx ON orders(client_request_id);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    product_name TEXT NOT NULL,
    qty INT NOT NULL CHECK (qty > 0),
    price NUMERIC(10, 2) NOT NULL,
    total NUMERIC(10, 2) NOT NULL
);

CREATE INDEX IF NOT EXISTS order_items_order_id_idx ON order_items(order_id);
