-- Migration 015: Add discount_percent column to line item tables
-- The discount percentage is the source of truth entered by the user.
-- The money amount (discount_per_unit) is derived from it.

ALTER TABLE quote_line_items
    ADD COLUMN discount_percent NUMERIC(5,2) NOT NULL DEFAULT 0;

ALTER TABLE order_line_items
    ADD COLUMN discount_percent NUMERIC(5,2) NOT NULL DEFAULT 0;

-- Backfill existing data: derive percent from money amounts
UPDATE quote_line_items
SET discount_percent = CASE
    WHEN unit_price_amount > 0
    THEN ROUND((discount_per_unit_amount / unit_price_amount) * 100, 2)
    ELSE 0
END
WHERE discount_per_unit_amount > 0;

UPDATE order_line_items
SET discount_percent = CASE
    WHEN unit_price_amount > 0
    THEN ROUND((discount_per_unit_amount / unit_price_amount) * 100, 2)
    ELSE 0
END
WHERE discount_per_unit_amount > 0;
