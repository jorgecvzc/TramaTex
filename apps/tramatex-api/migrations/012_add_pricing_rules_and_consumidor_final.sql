-- ============================================================================
-- Migration: 012_add_pricing_rules_and_consumidor_final.sql
-- Description: Add missing pricing_rules and sales_discount_rules tables,
--              seed CONSUMIDOR_FINAL party for simplified invoices/tickets
-- Date: 2026-03-09
-- ============================================================================

BEGIN;

-- ============================================================================
-- PRICING RULES TABLE (used by PricingRuleDataModel in Go)
-- ============================================================================
CREATE TABLE IF NOT EXISTS pricing_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    product_variant_id UUID,
    party_category VARCHAR(50),
    markup_percentage NUMERIC(8,4) NOT NULL,
    min_quantity INT NOT NULL DEFAULT 0,
    max_quantity INT,
    effective_from TIMESTAMP WITH TIME ZONE NOT NULL,
    effective_to TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pricing_rules_variant ON pricing_rules(product_variant_id);
CREATE INDEX IF NOT EXISTS idx_pricing_rules_effective ON pricing_rules(effective_from, effective_to);

COMMENT ON TABLE pricing_rules IS 'Pricing markup rules for product variants (quantity-based)';

-- ============================================================================
-- SALES DISCOUNT RULES TABLE (used by SalesDiscountRuleDataModel in Go)
-- ============================================================================
CREATE TABLE IF NOT EXISTS sales_discount_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    client_id UUID,
    product_variant_id UUID,
    min_quantity INT,
    discount_type VARCHAR(20) NOT NULL,
    percentage_value NUMERIC(8,4),
    fixed_amount NUMERIC(12,2),
    currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    priority INT NOT NULL DEFAULT 0,
    effective_from TIMESTAMP WITH TIME ZONE NOT NULL,
    effective_to TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sales_discount_rules_client ON sales_discount_rules(client_id);
CREATE INDEX IF NOT EXISTS idx_sales_discount_rules_variant ON sales_discount_rules(product_variant_id);
CREATE INDEX IF NOT EXISTS idx_sales_discount_rules_priority ON sales_discount_rules(priority);
CREATE INDEX IF NOT EXISTS idx_sales_discount_rules_effective ON sales_discount_rules(effective_from, effective_to);

COMMENT ON TABLE sales_discount_rules IS 'Sales discount rules for clients/variants';

-- ============================================================================
-- TRIGGERS
-- ============================================================================
DROP TRIGGER IF EXISTS trg_pricing_rules_updated_at ON pricing_rules;
DROP TRIGGER IF EXISTS trg_sales_discount_rules_updated_at ON sales_discount_rules;

CREATE TRIGGER trg_pricing_rules_updated_at BEFORE UPDATE ON pricing_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_sales_discount_rules_updated_at BEFORE UPDATE ON sales_discount_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SEED: CONSUMIDOR FINAL party (for simplified invoices / tickets)
-- UUID: 00000000-0000-0000-0000-000000000001
-- ============================================================================
INSERT INTO parties (id, status, created_by, modified_by, default_discount_percentage)
VALUES ('00000000-0000-0000-0000-000000000001', 'ACTIVE', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 0)
ON CONFLICT (id) DO NOTHING;

INSERT INTO organization_profiles (party_id, name, tax_id, tax_id_type)
VALUES ('00000000-0000-0000-0000-000000000001', 'CONSUMIDOR FINAL', NULL, NULL)
ON CONFLICT (party_id) DO NOTHING;

INSERT INTO party_roles (party_id, role)
VALUES ('00000000-0000-0000-0000-000000000001', 'CLIENT')
ON CONFLICT (party_id, role) DO NOTHING;

COMMIT;

-- ============================================================================
-- END OF MIGRATION: 012_add_pricing_rules_and_consumidor_final.sql
-- ============================================================================
