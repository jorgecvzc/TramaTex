-- Migration 017: Sequential document numbering
-- Replaces timestamp+UUID numbering with sequential counters per prefix+year.
-- Format: PREFIX-YEAR-NNNN (e.g., PRE-2026-0001, PED-2026-0001, ALB-2026-0001, FV-2026-0001, FT-2026-0001)

CREATE TABLE IF NOT EXISTS document_sequences (
    prefix       VARCHAR(10) NOT NULL,
    year         INTEGER     NOT NULL,
    current_value INTEGER    NOT NULL DEFAULT 0,
    PRIMARY KEY (prefix, year)
);

COMMENT ON TABLE document_sequences IS 'Sequential counters for document numbering per prefix and year';
COMMENT ON COLUMN document_sequences.prefix IS 'Document prefix: PRE (presupuestos), PED (pedidos), ALB (albaranes), FV (facturas venta), FT (facturas ticket)';
COMMENT ON COLUMN document_sequences.year IS 'Fiscal year for the sequence';
COMMENT ON COLUMN document_sequences.current_value IS 'Last assigned sequential number';

-- Seed initial counters based on existing documents so new numbers don't collide.
-- Each INSERT is wrapped in a DO block to gracefully skip if the source table does not exist.

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'quotes') THEN
    INSERT INTO document_sequences (prefix, year, current_value)
    SELECT 'PRE', EXTRACT(YEAR FROM quote_date)::INTEGER, COUNT(*)
    FROM quotes
    WHERE deleted_at IS NULL
    GROUP BY EXTRACT(YEAR FROM quote_date)
    ON CONFLICT (prefix, year) DO UPDATE SET current_value = GREATEST(document_sequences.current_value, EXCLUDED.current_value);
  END IF;
END $$;

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'sales_orders') THEN
    INSERT INTO document_sequences (prefix, year, current_value)
    SELECT 'PED', EXTRACT(YEAR FROM order_date)::INTEGER, COUNT(*)
    FROM sales_orders
    WHERE deleted_at IS NULL
    GROUP BY EXTRACT(YEAR FROM order_date)
    ON CONFLICT (prefix, year) DO UPDATE SET current_value = GREATEST(document_sequences.current_value, EXCLUDED.current_value);
  END IF;
END $$;

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'delivery_notes') THEN
    INSERT INTO document_sequences (prefix, year, current_value)
    SELECT 'ALB', EXTRACT(YEAR FROM delivery_date)::INTEGER, COUNT(*)
    FROM delivery_notes
    WHERE deleted_at IS NULL
    GROUP BY EXTRACT(YEAR FROM delivery_date)
    ON CONFLICT (prefix, year) DO UPDATE SET current_value = GREATEST(document_sequences.current_value, EXCLUDED.current_value);
  END IF;
END $$;

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'invoices') THEN
    INSERT INTO document_sequences (prefix, year, current_value)
    SELECT
        CASE WHEN series_code = 'TKT' THEN 'FT' ELSE 'FV' END,
        EXTRACT(YEAR FROM invoice_date)::INTEGER,
        COUNT(*)
    FROM invoices
    WHERE deleted_at IS NULL
    GROUP BY CASE WHEN series_code = 'TKT' THEN 'FT' ELSE 'FV' END, EXTRACT(YEAR FROM invoice_date)
    ON CONFLICT (prefix, year) DO UPDATE SET current_value = GREATEST(document_sequences.current_value, EXCLUDED.current_value);

    UPDATE invoices SET series_code = 'FV', series_prefix = 'FV' WHERE series_code = 'A';
    UPDATE invoices SET series_code = 'FT', series_prefix = 'FT' WHERE series_code = 'TKT';
  END IF;
END $$;
