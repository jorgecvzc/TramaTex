-- Migration: 018_add_invoice_type_and_series.sql
-- Author: AI Assistant
-- Description: Adds invoice type (COMPLETA/SIMPLIFICADA) and series fields to invoices table
-- for Spanish legislation compliance (ADR-020: Tickets and Invoice Series).
-- Date: 2026-02-14

-- Create invoice_type enum
CREATE TYPE invoice_type AS ENUM (
    'COMPLETA',
    'SIMPLIFICADA'
);

-- Add invoice type and series columns
ALTER TABLE "invoices" 
    ADD COLUMN IF NOT EXISTS "type" invoice_type NOT NULL DEFAULT 'COMPLETA',
    ADD COLUMN IF NOT EXISTS "series_code" VARCHAR(10) NOT NULL DEFAULT 'A',
    ADD COLUMN IF NOT EXISTS "series_year" INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM NOW()),
    ADD COLUMN IF NOT EXISTS "series_prefix" VARCHAR(10) NOT NULL DEFAULT 'A';

-- Create index for filtering by type (e.g., list all tickets)
CREATE INDEX IF NOT EXISTS "idx_invoices_type" ON "invoices" ("type");

-- Create composite index for series management (code + year)
CREATE INDEX IF NOT EXISTS "idx_invoices_series" ON "invoices" ("series_code", "series_year");

-- Add comment describing the purpose
COMMENT ON COLUMN "invoices"."type" IS 'Type of invoice: COMPLETA (full B2B invoice) or SIMPLIFICADA (ticket for retail < 3,000 EUR)';
COMMENT ON COLUMN "invoices"."series_code" IS 'Invoice series code (e.g., A, TKT, B)';
COMMENT ON COLUMN "invoices"."series_year" IS 'Fiscal year for the invoice series';
COMMENT ON COLUMN "invoices"."series_prefix" IS 'Prefix used for invoice number formatting';
