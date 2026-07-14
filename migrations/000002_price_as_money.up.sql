-- Store money as integer minor units plus an ISO 4217 currency code.
-- Floating point (and implicit currency) is how money bugs are born.
ALTER TABLE products ADD COLUMN price_cents BIGINT;
UPDATE products SET price_cents = ROUND(price * 100);
ALTER TABLE products ALTER COLUMN price_cents SET NOT NULL;
ALTER TABLE products DROP COLUMN price;

ALTER TABLE products ADD COLUMN currency TEXT;
UPDATE products SET currency = 'USD';
ALTER TABLE products ALTER COLUMN currency SET NOT NULL;
