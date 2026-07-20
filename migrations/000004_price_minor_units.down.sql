ALTER TABLE products RENAME COLUMN price_minor_units TO price_cents;

UPDATE outbox_events
SET payload = (payload - 'PriceMinorUnits') || jsonb_build_object('PriceCents', payload->'PriceMinorUnits')
WHERE event_name = 'product.created' AND payload ? 'PriceMinorUnits';
