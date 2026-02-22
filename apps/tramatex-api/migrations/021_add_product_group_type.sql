-- Migration: 020_add_product_group_type.sql
-- Created at: 2026-02-18
-- Description: Adds group_type column to product_groups table to distinguish between tangible products and services.

-- Create enum type for product group classification
DO $$ BEGIN
    CREATE TYPE product_group_type AS ENUM ('TANGIBLE', 'SERVICE');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Add group_type column with default value TANGIBLE
ALTER TABLE "product_groups" 
ADD COLUMN IF NOT EXISTS "group_type" product_group_type NOT NULL DEFAULT 'TANGIBLE';

-- Add index for filtering by group type
CREATE INDEX IF NOT EXISTS "idx_product_groups_group_type" ON "product_groups" ("group_type");

-- Update existing records to explicitly set TANGIBLE (for clarity)
UPDATE "product_groups" SET "group_type" = 'TANGIBLE' WHERE "group_type" IS NULL;

-- Add comment to document the column purpose
COMMENT ON COLUMN "product_groups"."group_type" IS 'Classification of product group: TANGIBLE for physical products, SERVICE for service-based products';
