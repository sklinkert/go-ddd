-- Create UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Sellers table
CREATE TABLE sellers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Products table
CREATE TABLE products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    seller_id UUID NOT NULL REFERENCES sellers(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Idempotency records table
CREATE TABLE idempotency_records (
    id UUID PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    request TEXT NOT NULL,
    response TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Indexes
CREATE INDEX idx_products_seller_id ON products(seller_id);
CREATE UNIQUE INDEX idx_idempotency_key ON idempotency_records(key);