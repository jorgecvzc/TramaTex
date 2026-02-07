-- Migration: 011_create_pricing_tables.sql
-- Created at: 2026-02-06
-- Description: Creates tables for the Pricing module (ADR-014).

CREATE TABLE IF NOT EXISTS "pricing_rules" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "product_variant_id" UUID,
    "party_category" VARCHAR(50),
    "markup_percentage" NUMERIC(8,4) NOT NULL,
    "min_quantity" INT NOT NULL DEFAULT 0,
    "max_quantity" INT,
    "effective_from" TIMESTAMP WITH TIME ZONE NOT NULL,
    "effective_to" TIMESTAMP WITH TIME ZONE,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS "idx_pricing_rules_variant" ON "pricing_rules" ("product_variant_id");

CREATE TABLE IF NOT EXISTS "client_pricing_overrides" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "client_id" UUID NOT NULL,
    "product_variant_id" UUID NOT NULL,
    "fixed_price" NUMERIC(12,2) NOT NULL,
    "currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "effective_from" TIMESTAMP WITH TIME ZONE NOT NULL,
    "effective_to" TIMESTAMP WITH TIME ZONE,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS "idx_client_pricing_override" ON "client_pricing_overrides" ("client_id", "product_variant_id");

CREATE TABLE IF NOT EXISTS "brand_profit_margins" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "brand_id" UUID NOT NULL,
    "percentage_value" NUMERIC(8,4),
    "fixed_amount" NUMERIC(12,2),
    "currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "effective_from" TIMESTAMP WITH TIME ZONE NOT NULL,
    "effective_to" TIMESTAMP WITH TIME ZONE,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS "idx_brand_profit_margin" ON "brand_profit_margins" ("brand_id");

CREATE TABLE IF NOT EXISTS "sales_discount_rules" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "client_id" UUID,
    "product_variant_id" UUID,
    "min_quantity" INT,
    "discount_type" VARCHAR(20) NOT NULL,
    "percentage_value" NUMERIC(8,4),
    "fixed_amount" NUMERIC(12,2),
    "currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "priority" INT NOT NULL DEFAULT 0,
    "effective_from" TIMESTAMP WITH TIME ZONE NOT NULL,
    "effective_to" TIMESTAMP WITH TIME ZONE,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS "idx_sales_discount_rules_variant" ON "sales_discount_rules" ("product_variant_id");

CREATE TABLE IF NOT EXISTS "price_calculations" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "product_variant_id" UUID NOT NULL,
    "client_id" UUID NOT NULL,
    "quantity" INT NOT NULL,
    "base_cost" NUMERIC(12,2) NOT NULL,
    "final_price" NUMERIC(12,2) NOT NULL,
    "currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "applied_rules" JSONB NOT NULL DEFAULT '[]',
    "calculated_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS "idx_price_calculations_variant" ON "price_calculations" ("product_variant_id");
