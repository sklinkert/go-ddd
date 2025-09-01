-- Drop indexes
DROP INDEX IF EXISTS idx_idempotency_key;
DROP INDEX IF EXISTS idx_products_seller_id;

-- Drop tables
DROP TABLE IF EXISTS idempotency_records;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS sellers;

-- Drop extension (only if no other objects depend on it)
DROP EXTENSION IF EXISTS "uuid-ossp";