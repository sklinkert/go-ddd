ALTER TABLE products ADD COLUMN price DECIMAL(10,2);
UPDATE products SET price = price_cents / 100.0;
ALTER TABLE products ALTER COLUMN price SET NOT NULL;
ALTER TABLE products DROP COLUMN price_cents;
ALTER TABLE products DROP COLUMN currency;
