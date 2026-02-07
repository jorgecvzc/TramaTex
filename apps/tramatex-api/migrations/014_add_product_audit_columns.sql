-- Migration: 014_add_product_audit_columns.sql
-- Created at: 2026-02-07
-- Description: Adds audit columns to Product module tables.

ALTER TABLE "attributes"
ADD COLUMN IF NOT EXISTS "created_by" VARCHAR(255),
ADD COLUMN IF NOT EXISTS "modified_by" VARCHAR(255);

ALTER TABLE "attribute_values"
ADD COLUMN IF NOT EXISTS "created_by" VARCHAR(255),
ADD COLUMN IF NOT EXISTS "modified_by" VARCHAR(255);

ALTER TABLE "products"
ADD COLUMN IF NOT EXISTS "created_by" VARCHAR(255),
ADD COLUMN IF NOT EXISTS "modified_by" VARCHAR(255);

ALTER TABLE "product_variants"
ADD COLUMN IF NOT EXISTS "created_by" VARCHAR(255),
ADD COLUMN IF NOT EXISTS "modified_by" VARCHAR(255);
