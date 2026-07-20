-- Rename the amount column to ISO 4217 terminology: "cents" is only
-- correct for exponent-2 currencies, "minor units" is currency-neutral.
ALTER TABLE products RENAME COLUMN price_cents TO price_minor_units;
