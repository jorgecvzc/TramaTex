-- Migration: 024_add_mes_work_reference_to_sales_line_items.sql
-- Created at: 2026-02-21
-- Description: Adds optional MES work reference to sales quote/order line items.

BEGIN;

DO $$
BEGIN
    IF to_regclass('public.quote_line_items') IS NOT NULL THEN
        ALTER TABLE "quote_line_items"
        ADD COLUMN IF NOT EXISTS "mes_work_id" UUID NULL;
    END IF;

    IF to_regclass('public.order_line_items') IS NOT NULL THEN
        ALTER TABLE "order_line_items"
        ADD COLUMN IF NOT EXISTS "mes_work_id" UUID NULL;
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('public.quote_line_items') IS NOT NULL
       AND to_regclass('public.mes_works') IS NOT NULL
       AND NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_quote_line_items_mes_work'
    ) THEN
        ALTER TABLE "quote_line_items"
        ADD CONSTRAINT "fk_quote_line_items_mes_work"
        FOREIGN KEY ("mes_work_id") REFERENCES "mes_works" ("id") ON DELETE SET NULL;
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('public.order_line_items') IS NOT NULL
       AND to_regclass('public.mes_works') IS NOT NULL
       AND NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_order_line_items_mes_work'
    ) THEN
        ALTER TABLE "order_line_items"
        ADD CONSTRAINT "fk_order_line_items_mes_work"
        FOREIGN KEY ("mes_work_id") REFERENCES "mes_works" ("id") ON DELETE SET NULL;
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('public.quote_line_items') IS NOT NULL THEN
        CREATE INDEX IF NOT EXISTS "idx_quote_line_items_mes_work_id" ON "quote_line_items" ("mes_work_id");
    END IF;

    IF to_regclass('public.order_line_items') IS NOT NULL THEN
        CREATE INDEX IF NOT EXISTS "idx_order_line_items_mes_work_id" ON "order_line_items" ("mes_work_id");
    END IF;
END $$;

COMMIT;
