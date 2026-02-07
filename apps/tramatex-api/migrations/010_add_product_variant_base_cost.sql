-- Migration: 010_add_product_variant_base_cost.sql
-- Created at: 2026-02-06
-- Description: Adds base cost to product variants.

ALTER TABLE "product_variants"
ADD COLUMN IF NOT EXISTS "base_cost" NUMERIC(12,2) NOT NULL DEFAULT 0;
