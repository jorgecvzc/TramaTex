-- ============================================================================
-- Migration: v2_005_init_sales.sql
-- Description: Initialize Sales module
-- Date: 2026-02-25
-- Modules: Quotes, Sales Orders, Delivery Notes, Invoices
-- ============================================================================

BEGIN;

-- ============================================================================
-- ENUMS
-- ============================================================================
DO $$ BEGIN
    CREATE TYPE quote_status AS ENUM (
        'BORRADOR',
        'ENVIADA',
        'APROBADA',
        'RECHAZADA',
        'EXPIRADA',
        'CONVERTIDA_A_PEDIDO'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE sales_order_status AS ENUM (
        'PENDIENTE',
        'EN_PREPARACION',
        'ENTREGADO_PARCIALMENTE',
        'ENTREGADO',
        'CANCELADO',
        'FACTURADO_PARCIALMENTE',
        'FACTURADO_COMPLETAMENTE'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE delivery_note_status AS ENUM (
        'PENDIENTE',
        'ENTREGADO',
        'CANCELADO'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE invoice_status AS ENUM (
        'BORRADOR',
        'EMITIDA',
        'PAGADA',
        'VENCIDA',
        'ANULADA'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE invoice_type AS ENUM (
        'COMPLETA',
        'SIMPLIFICADA'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- ============================================================================
-- QUOTES
-- ============================================================================
CREATE TABLE IF NOT EXISTS quotes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quote_number VARCHAR(50) NOT NULL,
    party_id UUID NOT NULL,
    quote_date TIMESTAMP WITH TIME ZONE NOT NULL,
    expiration_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status quote_status NOT NULL,
    subtotal_amount NUMERIC(12,2) NOT NULL,
    subtotal_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    tax_amount NUMERIC(12,2) NOT NULL,
    tax_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    total_amount NUMERIC(12,2) NOT NULL,
    total_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_quotes_number ON quotes(quote_number);
CREATE INDEX idx_quotes_party_id ON quotes(party_id);
CREATE INDEX idx_quotes_status ON quotes(status);
CREATE INDEX idx_quotes_quote_date ON quotes(quote_date);

COMMENT ON TABLE quotes IS 'Sales quotes (presupuestos)';

-- ============================================================================
CREATE TABLE IF NOT EXISTS quote_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quote_id UUID NOT NULL,
    product_variant_id UUID NOT NULL,
    quantity INT NOT NULL,
    calculated_unit_price_amount NUMERIC(12,2) NOT NULL,
    calculated_unit_price_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    manual_unit_price_amount NUMERIC(12,2),
    manual_unit_price_currency VARCHAR(3),
    final_unit_price_amount NUMERIC(12,2) NOT NULL,
    final_unit_price_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    calculated_discount_per_unit_amount NUMERIC(12,2),
    calculated_discount_per_unit_currency VARCHAR(3),
    manual_discount_per_unit_amount NUMERIC(12,2),
    manual_discount_per_unit_currency VARCHAR(3),
    final_discount_per_unit_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    final_discount_per_unit_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    tax_rate NUMERIC(5,2) NOT NULL DEFAULT 21.00,
    tax_amount NUMERIC(10,2),
    subtotal_amount NUMERIC(12,2) NOT NULL,
    subtotal_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    mes_work_id UUID NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_quote_line_items_quote FOREIGN KEY (quote_id) REFERENCES quotes(id) ON DELETE CASCADE,
    CONSTRAINT fk_quote_line_items_variant FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE RESTRICT,
    CONSTRAINT chk_quote_line_items_tax_rate CHECK (tax_rate >= 0 AND tax_rate <= 100)
);

CREATE INDEX idx_quote_line_items_quote_id ON quote_line_items(quote_id);
CREATE INDEX idx_quote_line_items_tax_rate ON quote_line_items(tax_rate);

COMMENT ON TABLE quote_line_items IS 'Line items in sales quotes';
COMMENT ON COLUMN quote_line_items.mes_work_id IS 'Optional reference to MES work (for service products)';

-- ============================================================================
-- SALES ORDERS
-- ============================================================================
CREATE TABLE IF NOT EXISTS sales_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) NOT NULL,
    quote_id UUID,
    party_id UUID NOT NULL,
    order_date TIMESTAMP WITH TIME ZONE NOT NULL,
    delivery_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status sales_order_status NOT NULL,
    subtotal_amount NUMERIC(12,2) NOT NULL,
    subtotal_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    tax_amount NUMERIC(12,2) NOT NULL,
    tax_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    total_amount NUMERIC(12,2) NOT NULL,
    total_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_sales_orders_quote FOREIGN KEY (quote_id) REFERENCES quotes(id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX idx_sales_orders_number ON sales_orders(order_number);
CREATE INDEX idx_sales_orders_party_id ON sales_orders(party_id);
CREATE INDEX idx_sales_orders_status ON sales_orders(status);
CREATE INDEX idx_sales_orders_order_date ON sales_orders(order_date);

COMMENT ON TABLE sales_orders IS 'Sales orders (pedidos)';

-- ============================================================================
CREATE TABLE IF NOT EXISTS order_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sales_order_id UUID NOT NULL,
    product_variant_id UUID NOT NULL,
    quantity INT NOT NULL,
    calculated_unit_price_amount NUMERIC(12,2) NOT NULL,
    calculated_unit_price_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    manual_unit_price_amount NUMERIC(12,2),
    manual_unit_price_currency VARCHAR(3),
    final_unit_price_amount NUMERIC(12,2) NOT NULL,
    final_unit_price_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    calculated_discount_per_unit_amount NUMERIC(12,2),
    calculated_discount_per_unit_currency VARCHAR(3),
    manual_discount_per_unit_amount NUMERIC(12,2),
    manual_discount_per_unit_currency VARCHAR(3),
    final_discount_per_unit_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    final_discount_per_unit_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    tax_rate NUMERIC(5,2) NOT NULL DEFAULT 21.00,
    tax_amount NUMERIC(10,2),
    subtotal_amount NUMERIC(12,2) NOT NULL,
    subtotal_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    mes_work_id UUID NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_order_line_items_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_order_line_items_variant FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE RESTRICT,
    CONSTRAINT chk_order_line_items_tax_rate CHECK (tax_rate >= 0 AND tax_rate <= 100)
);

CREATE INDEX idx_order_line_items_order_id ON order_line_items(sales_order_id);
CREATE INDEX idx_order_line_items_tax_rate ON order_line_items(tax_rate);

COMMENT ON TABLE order_line_items IS 'Line items in sales orders';
COMMENT ON COLUMN order_line_items.mes_work_id IS 'Optional reference to MES work (for service products)';

-- ============================================================================
-- DELIVERY NOTES
-- ============================================================================
CREATE TABLE IF NOT EXISTS delivery_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    delivery_note_number VARCHAR(50) NOT NULL,
    sales_order_id UUID NOT NULL,
    party_id UUID NOT NULL,
    delivery_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status delivery_note_status NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_delivery_notes_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_delivery_notes_number ON delivery_notes(delivery_note_number);
CREATE INDEX idx_delivery_notes_sales_order ON delivery_notes(sales_order_id);
CREATE INDEX idx_delivery_notes_party_id ON delivery_notes(party_id);

COMMENT ON TABLE delivery_notes IS 'Delivery notes (albaranes)';

-- ============================================================================
CREATE TABLE IF NOT EXISTS delivery_note_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    delivery_note_id UUID NOT NULL,
    sales_order_line_item_id UUID NOT NULL,
    product_variant_id UUID NOT NULL,
    delivered_quantity INT NOT NULL,
    tax_rate NUMERIC(5,2) NOT NULL DEFAULT 21.00,
    tax_amount NUMERIC(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_delivery_note_line_items_note FOREIGN KEY (delivery_note_id) REFERENCES delivery_notes(id) ON DELETE CASCADE,
    CONSTRAINT fk_delivery_note_line_items_order_item FOREIGN KEY (sales_order_line_item_id) REFERENCES order_line_items(id) ON DELETE RESTRICT,
    CONSTRAINT fk_delivery_note_line_items_variant FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE RESTRICT,
    CONSTRAINT chk_delivery_note_line_items_tax_rate CHECK (tax_rate >= 0 AND tax_rate <= 100)
);

CREATE INDEX idx_delivery_note_line_items_note_id ON delivery_note_line_items(delivery_note_id);
CREATE INDEX idx_delivery_note_line_items_tax_rate ON delivery_note_line_items(tax_rate);

COMMENT ON TABLE delivery_note_line_items IS 'Line items in delivery notes';

-- ============================================================================
-- INVOICES
-- ============================================================================
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50) NOT NULL,
    type invoice_type NOT NULL DEFAULT 'COMPLETA',
    series_code VARCHAR(10) NOT NULL DEFAULT 'A',
    series_year INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM NOW()),
    series_prefix VARCHAR(10) NOT NULL DEFAULT 'A',
    party_id UUID NOT NULL,
    invoice_date TIMESTAMP WITH TIME ZONE NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status invoice_status NOT NULL,
    payment_terms TEXT,
    subtotal_amount NUMERIC(12,2) NOT NULL,
    subtotal_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    tax_amount NUMERIC(12,2) NOT NULL,
    tax_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    total_amount NUMERIC(12,2) NOT NULL,
    total_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_invoices_number ON invoices(invoice_number);
CREATE INDEX idx_invoices_party_id ON invoices(party_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_type ON invoices(type);
CREATE INDEX idx_invoices_series ON invoices(series_code, series_year);
CREATE INDEX idx_invoices_invoice_date ON invoices(invoice_date);

COMMENT ON TABLE invoices IS 'Sales invoices (facturas)';
COMMENT ON COLUMN invoices.type IS 'COMPLETA (full B2B invoice) or SIMPLIFICADA (ticket < 3,000 EUR)';
COMMENT ON COLUMN invoices.series_code IS 'Invoice series code (A, TKT, B, etc.)';

-- ============================================================================
CREATE TABLE IF NOT EXISTS invoice_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    sales_order_line_item_id UUID,
    product_variant_id UUID NOT NULL,
    quantity INT NOT NULL,
    unit_price_amount NUMERIC(12,2) NOT NULL,
    unit_price_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    discount_amount NUMERIC(12,2),
    discount_currency VARCHAR(3),
    tax_rate NUMERIC(5,2) NOT NULL DEFAULT 21.00,
    tax_amount NUMERIC(10,2),
    subtotal_amount NUMERIC(12,2) NOT NULL,
    subtotal_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    total_amount NUMERIC(12,2) NOT NULL,
    total_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_invoice_line_items_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
    CONSTRAINT fk_invoice_line_items_order_item FOREIGN KEY (sales_order_line_item_id) REFERENCES order_line_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_invoice_line_items_variant FOREIGN KEY (product_variant_id) REFERENCES product_variants(id) ON DELETE RESTRICT,
    CONSTRAINT chk_invoice_line_items_tax_rate CHECK (tax_rate >= 0 AND tax_rate <= 100)
);

CREATE INDEX idx_invoice_line_items_invoice_id ON invoice_line_items(invoice_id);
CREATE INDEX idx_invoice_line_items_tax_rate ON invoice_line_items(tax_rate);

COMMENT ON TABLE invoice_line_items IS 'Line items in invoices';

-- ============================================================================
-- TRIGGERS
-- ============================================================================
DROP TRIGGER IF EXISTS trg_quotes_updated_at ON quotes;
DROP TRIGGER IF EXISTS trg_quote_line_items_updated_at ON quote_line_items;
DROP TRIGGER IF EXISTS trg_sales_orders_updated_at ON sales_orders;
DROP TRIGGER IF EXISTS trg_order_line_items_updated_at ON order_line_items;
DROP TRIGGER IF EXISTS trg_delivery_notes_updated_at ON delivery_notes;
DROP TRIGGER IF EXISTS trg_delivery_note_line_items_updated_at ON delivery_note_line_items;
DROP TRIGGER IF EXISTS trg_invoices_updated_at ON invoices;
DROP TRIGGER IF EXISTS trg_invoice_line_items_updated_at ON invoice_line_items;

CREATE TRIGGER trg_quotes_updated_at BEFORE UPDATE ON quotes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_quote_line_items_updated_at BEFORE UPDATE ON quote_line_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_sales_orders_updated_at BEFORE UPDATE ON sales_orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_order_line_items_updated_at BEFORE UPDATE ON order_line_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_delivery_notes_updated_at BEFORE UPDATE ON delivery_notes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_delivery_note_line_items_updated_at BEFORE UPDATE ON delivery_note_line_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_invoices_updated_at BEFORE UPDATE ON invoices FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_invoice_line_items_updated_at BEFORE UPDATE ON invoice_line_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;

-- ============================================================================
-- END OF MIGRATION: v2_005_init_sales.sql
-- ============================================================================
