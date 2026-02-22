-- Migration: 025_add_product_base_price.sql
-- Created at: 2026-02-19
-- Description: Adds base_price field to products table.
-- This base price is used as the default cost for product variants when created JIT or pre-generated.

ALTER TABLE "products"
ADD COLUMN IF NOT EXISTS "base_price" NUMERIC(12,2) NOT NULL DEFAULT 0;

COMMENT ON COLUMN "products"."base_price" IS 'Base cost/price of the product, used as default for variants';
