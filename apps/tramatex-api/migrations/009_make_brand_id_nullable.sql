-- Migration to make brand_id nullable in products and other tables
-- This allows products to be created without a mandatory brand association.

ALTER TABLE products ALTER COLUMN brand_id DROP NOT NULL;

-- Also check brand_profit_margins if it should be nullable (keeping it NOT NULL if it defines margins per brand)
-- In 004_init_pricing.sql: brand_id UUID NOT NULL
-- If a product has no brand, it just won''t have a brand-specific margin. That''s fine.

-- base_sales_price_rules already has nullable brand_id in some places but not others?
-- Let''s make it nullable where it makes sense.
