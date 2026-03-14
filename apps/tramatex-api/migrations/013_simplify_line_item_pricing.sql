-- Migration 013: Simplify line item pricing model
-- Replaces 3-tier (Calculated/Manual/Final) with 2-tier (List/Sale) pricing.
-- Tables are empty at this point, so we can safely restructure columns.

-- ============================================================================
-- QUOTE LINE ITEMS
-- ============================================================================

-- Rename calculated → list (precio de tarifa)
ALTER TABLE quote_line_items RENAME COLUMN calculated_unit_price_amount TO list_unit_price_amount;
ALTER TABLE quote_line_items RENAME COLUMN calculated_unit_price_currency TO list_unit_price_currency;

-- Rename final → unit (precio de venta)
ALTER TABLE quote_line_items RENAME COLUMN final_unit_price_amount TO unit_price_amount;
ALTER TABLE quote_line_items RENAME COLUMN final_unit_price_currency TO unit_price_currency;

-- Drop manual price columns (no longer needed)
ALTER TABLE quote_line_items DROP COLUMN IF EXISTS manual_unit_price_amount;
ALTER TABLE quote_line_items DROP COLUMN IF EXISTS manual_unit_price_currency;

-- Consolidate discount: rename final_discount → discount_per_unit, make NOT NULL with default
ALTER TABLE quote_line_items RENAME COLUMN final_discount_per_unit_amount TO discount_per_unit_amount;
ALTER TABLE quote_line_items RENAME COLUMN final_discount_per_unit_currency TO discount_per_unit_currency;

-- Drop calculated and manual discount columns
ALTER TABLE quote_line_items DROP COLUMN IF EXISTS calculated_discount_per_unit_amount;
ALTER TABLE quote_line_items DROP COLUMN IF EXISTS calculated_discount_per_unit_currency;
ALTER TABLE quote_line_items DROP COLUMN IF EXISTS manual_discount_per_unit_amount;
ALTER TABLE quote_line_items DROP COLUMN IF EXISTS manual_discount_per_unit_currency;

-- ============================================================================
-- ORDER LINE ITEMS
-- ============================================================================

-- Rename calculated → list (precio de tarifa)
ALTER TABLE order_line_items RENAME COLUMN calculated_unit_price_amount TO list_unit_price_amount;
ALTER TABLE order_line_items RENAME COLUMN calculated_unit_price_currency TO list_unit_price_currency;

-- Rename final → unit (precio de venta)
ALTER TABLE order_line_items RENAME COLUMN final_unit_price_amount TO unit_price_amount;
ALTER TABLE order_line_items RENAME COLUMN final_unit_price_currency TO unit_price_currency;

-- Drop manual price columns
ALTER TABLE order_line_items DROP COLUMN IF EXISTS manual_unit_price_amount;
ALTER TABLE order_line_items DROP COLUMN IF EXISTS manual_unit_price_currency;

-- Consolidate discount: rename final_discount → discount_per_unit
ALTER TABLE order_line_items RENAME COLUMN final_discount_per_unit_amount TO discount_per_unit_amount;
ALTER TABLE order_line_items RENAME COLUMN final_discount_per_unit_currency TO discount_per_unit_currency;

-- Drop calculated and manual discount columns
ALTER TABLE order_line_items DROP COLUMN IF EXISTS calculated_discount_per_unit_amount;
ALTER TABLE order_line_items DROP COLUMN IF EXISTS calculated_discount_per_unit_currency;
ALTER TABLE order_line_items DROP COLUMN IF EXISTS manual_discount_per_unit_amount;
ALTER TABLE order_line_items DROP COLUMN IF EXISTS manual_discount_per_unit_currency;
