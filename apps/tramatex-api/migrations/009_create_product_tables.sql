-- Migration: 009_create_product_tables.sql
-- Created at: 2026-02-03 19:00:00
-- Description: Creates the initial tables for the Product module based on the finalized domain model.

-- 1. brands
CREATE TABLE IF NOT EXISTS "brands" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_brands_name" ON "brands" ("name");

-- 2. product_groups
CREATE TABLE IF NOT EXISTS "product_groups" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "parent_group_id" UUID,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_parent_group" FOREIGN KEY ("parent_group_id") REFERENCES "product_groups" ("id") ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS "idx_product_groups_parent_group_id" ON "product_groups" ("parent_group_id");

-- 3. attributes
CREATE TABLE IF NOT EXISTS "attributes" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "code" VARCHAR(50) NOT NULL,
    "sort_order" INT NOT NULL DEFAULT 0,
    "scope_brand_id" UUID,
    "scope_group_id" UUID,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_scope_brand" FOREIGN KEY ("scope_brand_id") REFERENCES "brands" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_scope_product_group" FOREIGN KEY ("scope_group_id") REFERENCES "product_groups" ("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_attributes_code" ON "attributes" ("code");
CREATE INDEX IF NOT EXISTS "idx_attributes_scope" ON "attributes" ("scope_brand_id", "scope_group_id");

-- 4. attribute_values
CREATE TABLE IF NOT EXISTS "attribute_values" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "attribute_id" UUID NOT NULL,
    "value" VARCHAR(255) NOT NULL,
    "code" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES "attributes" ("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_attribute_values_attribute_id_code" ON "attribute_values" ("attribute_id", "code");

-- 5. products
CREATE TYPE product_type AS ENUM ('TANGIBLE', 'SERVICE');
CREATE TABLE IF NOT EXISTS "products" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "sku" VARCHAR(255) UNIQUE,
    "name" VARCHAR(100) NOT NULL,
    "long_name" VARCHAR(255),
    "barcode" VARCHAR(255) UNIQUE,
    "description" TEXT,
    "product_type" product_type NOT NULL,
    "brand_id" UUID NOT NULL,
    "group_ids" UUID[],
    "direct_attribute_ids" UUID[],
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_products_brand" FOREIGN KEY ("brand_id") REFERENCES "brands" ("id") ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS "idx_products_brand_id" ON "products" ("brand_id");

-- 6. product_variants
CREATE TYPE variant_status AS ENUM ('PROVISIONAL', 'CONFIRMED');
CREATE TABLE IF NOT EXISTS "product_variants" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "product_id" UUID NOT NULL,
    "sku" VARCHAR(255) NOT NULL,
    "barcode" VARCHAR(255),
    "status" variant_status NOT NULL DEFAULT 'PROVISIONAL',
    "attribute_values" UUID[],
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_product_variants_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_product_variants_sku" ON "product_variants" ("sku");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_product_variants_barcode" ON "product_variants" ("barcode") WHERE barcode IS NOT NULL;


-- Timestamps update trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Drop existing triggers before creating new ones to avoid errors on re-run
DROP TRIGGER IF EXISTS update_brands_updated_at ON "brands";
DROP TRIGGER IF EXISTS update_product_groups_updated_at ON "product_groups";
DROP TRIGGER IF EXISTS update_attributes_updated_at ON "attributes";
DROP TRIGGER IF EXISTS update_attribute_values_updated_at ON "attribute_values";
DROP TRIGGER IF EXISTS update_products_updated_at ON "products";
DROP TRIGGER IF EXISTS update_product_variants_updated_at ON "product_variants";

-- Create triggers
CREATE TRIGGER update_brands_updated_at BEFORE UPDATE ON "brands" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_product_groups_updated_at BEFORE UPDATE ON "product_groups" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_attributes_updated_at BEFORE UPDATE ON "attributes" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_attribute_values_updated_at BEFORE UPDATE ON "attribute_values" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON "products" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_product_variants_updated_at BEFORE UPDATE ON "product_variants" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
