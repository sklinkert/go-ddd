-- Rename the amount column to ISO 4217 terminology: "cents" is only
-- correct for exponent-2 currencies, "minor units" is currency-neutral.
ALTER TABLE products RENAME COLUMN price_cents TO price_minor_units;

-- Rewrite stored event payloads to the renamed field so pending outbox
-- rows publish with the new schema instead of a silently-zero price.
UPDATE outbox_events
SET payload = (payload - 'PriceCents') || jsonb_build_object('PriceMinorUnits', payload->'PriceCents')
WHERE event_name = 'product.created' AND payload ? 'PriceCents';
