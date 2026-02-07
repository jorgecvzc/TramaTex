-- Migration: 012_create_pricing_rules_v2.sql
-- Created at: 2026-02-06
-- Description: Creates ADR-015 pricing rule tables.

CREATE TABLE IF NOT EXISTS "rule_value_types" (
    "type" VARCHAR(50) PRIMARY KEY,
    "description" VARCHAR(255) NOT NULL
);

INSERT INTO "rule_value_types" ("type", "description") VALUES
    ('PERCENTAGE_MARKUP', 'Increase by percentage'),
    ('FIXED_AMOUNT_INCREASE', 'Increase by fixed amount'),
    ('SET_TO_FIXED_PRICE', 'Set to fixed price'),
    ('APPLY_PERCENTAGE_DISCOUNT', 'Discount by percentage'),
    ('APPLY_FIXED_AMOUNT_DISCOUNT', 'Discount by fixed amount'),
    ('SET_TO_FIXED_DISCOUNTED_PRICE', 'Set to fixed discounted price')
ON CONFLICT ("type") DO NOTHING;

CREATE TABLE IF NOT EXISTS "base_sales_price_rules" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "brand_id" UUID,
    "product_group_id" UUID,
    "product_id" UUID,
    "variant_id" UUID,
    "value_type" VARCHAR(50) NOT NULL,
    "percentage_value" NUMERIC(8,4),
    "money_value_amount" NUMERIC(12,2),
    "money_value_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS "idx_base_sales_price_rules_variant" ON "base_sales_price_rules" ("variant_id");

CREATE TABLE IF NOT EXISTS "sale_modification_rules" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "client_ids" UUID[],
    "product_group_id" UUID,
    "min_order_total_amount" NUMERIC(12,2),
    "min_order_total_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "value_type" VARCHAR(50) NOT NULL,
    "percentage_value" NUMERIC(8,4),
    "money_value_amount" NUMERIC(12,2),
    "money_value_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "priority" INT NOT NULL DEFAULT 0,
    "effective_from" TIMESTAMP WITH TIME ZONE NOT NULL,
    "effective_to" TIMESTAMP WITH TIME ZONE,
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS "idx_sale_modification_rules_group" ON "sale_modification_rules" ("product_group_id");
