-- Migration: 015_create_sales_tables.sql
-- Created at: 2026-02-07
-- Description: Creates tables for Sales module (quotes, orders, delivery notes, invoices).

-- Enums
CREATE TYPE quote_status AS ENUM (
    'BORRADOR',
    'ENVIADA',
    'APROBADA',
    'RECHAZADA',
    'EXPIRADA',
    'CONVERTIDA_A_PEDIDO'
);

CREATE TYPE sales_order_status AS ENUM (
    'PENDIENTE',
    'EN_PREPARACION',
    'ENTREGADO_PARCIALMENTE',
    'ENTREGADO',
    'CANCELADO',
    'FACTURADO_PARCIALMENTE',
    'FACTURADO_COMPLETAMENTE'
);

CREATE TYPE delivery_note_status AS ENUM (
    'PENDIENTE',
    'ENTREGADO',
    'CANCELADO'
);

CREATE TYPE invoice_status AS ENUM (
    'BORRADOR',
    'EMITIDA',
    'PAGADA',
    'VENCIDA',
    'ANULADA'
);

-- Quotes
CREATE TABLE IF NOT EXISTS "quotes" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "quote_number" VARCHAR(50) NOT NULL,
    "party_id" UUID NOT NULL,
    "quote_date" TIMESTAMP WITH TIME ZONE NOT NULL,
    "expiration_date" TIMESTAMP WITH TIME ZONE NOT NULL,
    "status" quote_status NOT NULL,
    "subtotal_amount" NUMERIC(12,2) NOT NULL,
    "subtotal_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "tax_amount" NUMERIC(12,2) NOT NULL,
    "tax_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "total_amount" NUMERIC(12,2) NOT NULL,
    "total_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "notes" TEXT,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_quotes_number" ON "quotes" ("quote_number");
CREATE INDEX IF NOT EXISTS "idx_quotes_party_id" ON "quotes" ("party_id");
CREATE INDEX IF NOT EXISTS "idx_quotes_status" ON "quotes" ("status");
CREATE INDEX IF NOT EXISTS "idx_quotes_quote_date" ON "quotes" ("quote_date");

CREATE TABLE IF NOT EXISTS "quote_line_items" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "quote_id" UUID NOT NULL,
    "product_variant_id" UUID NOT NULL,
    "quantity" INT NOT NULL,
    "calculated_unit_price_amount" NUMERIC(12,2) NOT NULL,
    "calculated_unit_price_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "manual_unit_price_amount" NUMERIC(12,2),
    "manual_unit_price_currency" VARCHAR(3),
    "final_unit_price_amount" NUMERIC(12,2) NOT NULL,
    "final_unit_price_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "calculated_discount_per_unit_amount" NUMERIC(12,2),
    "calculated_discount_per_unit_currency" VARCHAR(3),
    "manual_discount_per_unit_amount" NUMERIC(12,2),
    "manual_discount_per_unit_currency" VARCHAR(3),
    "final_discount_per_unit_amount" NUMERIC(12,2) NOT NULL DEFAULT 0,
    "final_discount_per_unit_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "subtotal_amount" NUMERIC(12,2) NOT NULL,
    "subtotal_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_quote_line_items_quote" FOREIGN KEY ("quote_id") REFERENCES "quotes" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_quote_line_items_variant" FOREIGN KEY ("product_variant_id") REFERENCES "product_variants" ("id") ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS "idx_quote_line_items_quote_id" ON "quote_line_items" ("quote_id");

-- Sales Orders
CREATE TABLE IF NOT EXISTS "sales_orders" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "order_number" VARCHAR(50) NOT NULL,
    "quote_id" UUID,
    "party_id" UUID NOT NULL,
    "order_date" TIMESTAMP WITH TIME ZONE NOT NULL,
    "delivery_date" TIMESTAMP WITH TIME ZONE NOT NULL,
    "status" sales_order_status NOT NULL,
    "subtotal_amount" NUMERIC(12,2) NOT NULL,
    "subtotal_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "tax_amount" NUMERIC(12,2) NOT NULL,
    "tax_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "total_amount" NUMERIC(12,2) NOT NULL,
    "total_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "notes" TEXT,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_sales_orders_quote" FOREIGN KEY ("quote_id") REFERENCES "quotes" ("id") ON DELETE SET NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_sales_orders_number" ON "sales_orders" ("order_number");
CREATE INDEX IF NOT EXISTS "idx_sales_orders_party_id" ON "sales_orders" ("party_id");
CREATE INDEX IF NOT EXISTS "idx_sales_orders_status" ON "sales_orders" ("status");
CREATE INDEX IF NOT EXISTS "idx_sales_orders_order_date" ON "sales_orders" ("order_date");

CREATE TABLE IF NOT EXISTS "order_line_items" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "sales_order_id" UUID NOT NULL,
    "product_variant_id" UUID NOT NULL,
    "quantity" INT NOT NULL,
    "calculated_unit_price_amount" NUMERIC(12,2) NOT NULL,
    "calculated_unit_price_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "manual_unit_price_amount" NUMERIC(12,2),
    "manual_unit_price_currency" VARCHAR(3),
    "final_unit_price_amount" NUMERIC(12,2) NOT NULL,
    "final_unit_price_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "calculated_discount_per_unit_amount" NUMERIC(12,2),
    "calculated_discount_per_unit_currency" VARCHAR(3),
    "manual_discount_per_unit_amount" NUMERIC(12,2),
    "manual_discount_per_unit_currency" VARCHAR(3),
    "final_discount_per_unit_amount" NUMERIC(12,2) NOT NULL DEFAULT 0,
    "final_discount_per_unit_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "subtotal_amount" NUMERIC(12,2) NOT NULL,
    "subtotal_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_order_line_items_order" FOREIGN KEY ("sales_order_id") REFERENCES "sales_orders" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_order_line_items_variant" FOREIGN KEY ("product_variant_id") REFERENCES "product_variants" ("id") ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS "idx_order_line_items_order_id" ON "order_line_items" ("sales_order_id");

-- Delivery Notes
CREATE TABLE IF NOT EXISTS "delivery_notes" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "delivery_note_number" VARCHAR(50) NOT NULL,
    "sales_order_id" UUID NOT NULL,
    "party_id" UUID NOT NULL,
    "delivery_date" TIMESTAMP WITH TIME ZONE NOT NULL,
    "status" delivery_note_status NOT NULL,
    "notes" TEXT,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_delivery_notes_order" FOREIGN KEY ("sales_order_id") REFERENCES "sales_orders" ("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_delivery_notes_number" ON "delivery_notes" ("delivery_note_number");
CREATE INDEX IF NOT EXISTS "idx_delivery_notes_sales_order" ON "delivery_notes" ("sales_order_id");
CREATE INDEX IF NOT EXISTS "idx_delivery_notes_party_id" ON "delivery_notes" ("party_id");

CREATE TABLE IF NOT EXISTS "delivery_note_line_items" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "delivery_note_id" UUID NOT NULL,
    "sales_order_line_item_id" UUID NOT NULL,
    "product_variant_id" UUID NOT NULL,
    "delivered_quantity" INT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_delivery_note_line_items_note" FOREIGN KEY ("delivery_note_id") REFERENCES "delivery_notes" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_delivery_note_line_items_order_item" FOREIGN KEY ("sales_order_line_item_id") REFERENCES "order_line_items" ("id") ON DELETE RESTRICT,
    CONSTRAINT "fk_delivery_note_line_items_variant" FOREIGN KEY ("product_variant_id") REFERENCES "product_variants" ("id") ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS "idx_delivery_note_line_items_note_id" ON "delivery_note_line_items" ("delivery_note_id");

-- Invoices
CREATE TABLE IF NOT EXISTS "invoices" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "invoice_number" VARCHAR(50) NOT NULL,
    "party_id" UUID NOT NULL,
    "invoice_date" TIMESTAMP WITH TIME ZONE NOT NULL,
    "due_date" TIMESTAMP WITH TIME ZONE NOT NULL,
    "status" invoice_status NOT NULL,
    "payment_terms" TEXT,
    "subtotal_amount" NUMERIC(12,2) NOT NULL,
    "subtotal_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "tax_amount" NUMERIC(12,2) NOT NULL,
    "tax_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "total_amount" NUMERIC(12,2) NOT NULL,
    "total_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_invoices_number" ON "invoices" ("invoice_number");
CREATE INDEX IF NOT EXISTS "idx_invoices_party_id" ON "invoices" ("party_id");
CREATE INDEX IF NOT EXISTS "idx_invoices_status" ON "invoices" ("status");
CREATE INDEX IF NOT EXISTS "idx_invoices_invoice_date" ON "invoices" ("invoice_date");

CREATE TABLE IF NOT EXISTS "invoice_line_items" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "invoice_id" UUID NOT NULL,
    "sales_order_line_item_id" UUID,
    "product_variant_id" UUID NOT NULL,
    "quantity" INT NOT NULL,
    "unit_price_amount" NUMERIC(12,2) NOT NULL,
    "unit_price_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "discount_amount" NUMERIC(12,2),
    "discount_currency" VARCHAR(3),
    "subtotal_amount" NUMERIC(12,2) NOT NULL,
    "subtotal_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "tax_amount" NUMERIC(12,2),
    "tax_currency" VARCHAR(3),
    "total_amount" NUMERIC(12,2) NOT NULL,
    "total_currency" VARCHAR(3) NOT NULL DEFAULT 'EUR',
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT "fk_invoice_line_items_invoice" FOREIGN KEY ("invoice_id") REFERENCES "invoices" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_invoice_line_items_order_item" FOREIGN KEY ("sales_order_line_item_id") REFERENCES "order_line_items" ("id") ON DELETE SET NULL,
    CONSTRAINT "fk_invoice_line_items_variant" FOREIGN KEY ("product_variant_id") REFERENCES "product_variants" ("id") ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS "idx_invoice_line_items_invoice_id" ON "invoice_line_items" ("invoice_id");

-- Update triggers
DROP TRIGGER IF EXISTS update_quotes_updated_at ON "quotes";
DROP TRIGGER IF EXISTS update_quote_line_items_updated_at ON "quote_line_items";
DROP TRIGGER IF EXISTS update_sales_orders_updated_at ON "sales_orders";
DROP TRIGGER IF EXISTS update_order_line_items_updated_at ON "order_line_items";
DROP TRIGGER IF EXISTS update_delivery_notes_updated_at ON "delivery_notes";
DROP TRIGGER IF EXISTS update_delivery_note_line_items_updated_at ON "delivery_note_line_items";
DROP TRIGGER IF EXISTS update_invoices_updated_at ON "invoices";
DROP TRIGGER IF EXISTS update_invoice_line_items_updated_at ON "invoice_line_items";

CREATE TRIGGER update_quotes_updated_at BEFORE UPDATE ON "quotes" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_quote_line_items_updated_at BEFORE UPDATE ON "quote_line_items" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_sales_orders_updated_at BEFORE UPDATE ON "sales_orders" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_order_line_items_updated_at BEFORE UPDATE ON "order_line_items" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_delivery_notes_updated_at BEFORE UPDATE ON "delivery_notes" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_delivery_note_line_items_updated_at BEFORE UPDATE ON "delivery_note_line_items" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_invoices_updated_at BEFORE UPDATE ON "invoices" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_invoice_line_items_updated_at BEFORE UPDATE ON "invoice_line_items" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
