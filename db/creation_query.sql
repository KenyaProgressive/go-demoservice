CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS customers (
  id BIGSERIAL PRIMARY KEY,
  customer_id TEXT NOT NULL UNIQUE,
  email VARCHAR(100),
  name VARCHAR(50),
  -- максимальная длина номера по Миру в 15 цифр и знак +
  phone VARCHAR(16)
);

CREATE TABLE IF NOT EXISTS order_info(
    order_uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id TEXT REFERENCES customers(customer_id),
    track_number TEXT,
    entry TEXT NOT NULL,
    -- locale: максимальная длина номера по Миру в 15 цифр и знак +
    locale VARCHAR(2) CHECK (locale ~ '^[a-z]{2}$'),
    delivery_service TEXT NOT NULL,
    shardkey TEXT,
    sm_id INT,
    oof_shard TEXT,
    internal_signature TEXT,
    date_created TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS payments(
    id BIGSERIAL PRIMARY KEY,
    transaction UUID REFERENCES order_info(order_uid),
    request_id TEXT,
    provider TEXT NOT NULL,
    bank TEXT NOT NULL,
    amount BIGINT NOT NULL,
    -- currency: валюта - USD, RUB и т.п.
    currency CHAR(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    payment_dt BIGINT NOT NULL,
    delivery_cost BIGINT NOT NULL,
    custom_fee BIGINT,
    goods_total INT NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries (
  id BIGSERIAL PRIMARY KEY,
  order_uid UUID REFERENCES order_info(order_uid),
  name VARCHAR(50) NOT NULL,
  phone VARCHAR(16) NOT NULL,
  -- zip: по самой большой длине zip-кода в Мире
  zip VARCHAR(8),
  -- city: по самой большой длине названия города в Мире
  city VARCHAR(170) NOT NULL,
  address VARCHAR(100) NOT NULL,
  region VARCHAR(100) NOT NULL,
  email VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS order_items (
  id BIGSERIAL PRIMARY KEY,
  order_uid UUID REFERENCES order_info(order_uid),
  nm_id BIGINT,
  chrt_id BIGINT,
  rid TEXT,
  name TEXT,
  brand TEXT,
  size TEXT,
  price BIGINT,
  sale INT,
  total_price INT,
  status INT
);
