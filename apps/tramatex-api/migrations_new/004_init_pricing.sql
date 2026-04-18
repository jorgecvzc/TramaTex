-- ============================================================================
-- Migration: 004_init_pricing.sql
-- Module: Pricing (Base Sales Price Rules, Sale Modifications, Client Overrides)
-- Date: 2026-04-14
-- ============================================================================

BEGIN;

-- ============================================================================
-- RULE VALUE TYPES CATALOG
-- ============================================================================
CREATE TABLE IF NOT EXISTS rule_value_types (
    type VARCHAR(50) PRIMARY KEY,
    description VARCHAR(255) NOT NULL
);

INSERT INTO rule_value_types (type, description) VALUES
    ('PERCENTAGE_MARKUP', 'Increase by percentage'),
    ('FIXED_AMOUNT_INCREASE', 'Increase by fixed amount'),
    ('SET_TO_FIXED_PRICE', 'Set to fixed price'),
    ('APPLY_PERCENTAGE_DISCOUNT', 'Discount by percentage'),
    ('APPLY_FIXED_AMOUNT_DISCOUNT', 'Discount by fixed amount'),
    ('SET_TO_FIXED_DISCOUNTED_PRICE', 'Set to fixed discounted price')
ON CONFLICT (type) DO NOTHING;

COMMENT ON TABLE rule_value_types IS 'Catalog of pricing rule value types';

-- ============================================================================
-- BASE SALES PRICE RULES (ADR-015)
-- ============================================================================
CREATE TABLE IF NOT EXISTS base_sales_price_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    brand_id UUID,
    product_group_id UUID,
    product_id UUID,
    variant_id UUID,
    value_type VARCHAR(50) NOT NULL,
    percentage_value NUMERIC(8,4),
    money_value_amount NUMERIC(12,2),
    money_value_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_base_sales_price_rules_value_type FOREIGN KEY (value_type) REFERENCES rule_value_types(type)
);

CREATE INDEX idx_base_sales_price_rules_variant ON base_sales_price_rules(variant_id);
CREATE INDEX idx_base_sales_price_rules_brand ON base_sales_price_rules(brand_id);
CREATE INDEX idx_base_sales_price_rules_product_group ON base_sales_price_rules(product_group_id);
CREATE INDEX idx_base_sales_price_rules_product ON base_sales_price_rules(product_id);

COMMENT ON TABLE base_sales_price_rules IS 'Rules for calculating base sales prices from variant base costs';

-- ============================================================================
-- SALE MODIFICATION RULES (ADR-015)
-- ============================================================================
CREATE TABLE IF NOT EXISTS sale_modification_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    client_ids UUID[],
    product_group_id UUID,
    min_order_total_amount NUMERIC(12,2),
    min_order_total_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    value_type VARCHAR(50) NOT NULL,
    percentage_value NUMERIC(8,4),
    money_value_amount NUMERIC(12,2),
    money_value_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    priority INT NOT NULL DEFAULT 0,
    effective_from TIMESTAMP WITH TIME ZONE NOT NULL,
    effective_to TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_sale_modification_rules_value_type FOREIGN KEY (value_type) REFERENCES rule_value_types(type)
);

CREATE INDEX idx_sale_modification_rules_group ON sale_modification_rules(product_group_id);
CREATE INDEX idx_sale_modification_rules_priority ON sale_modification_rules(priority);
CREATE INDEX idx_sale_modification_rules_effective ON sale_modification_rules(effective_from, effective_to);

COMMENT ON TABLE sale_modification_rules IS 'Rules for modifying sales prices (discounts, bulk pricing, etc.)';

-- ============================================================================
-- CLIENT PRICING OVERRIDES
-- ============================================================================
CREATE TABLE IF NOT EXISTS client_pricing_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL,
    product_variant_id UUID NOT NULL,
    fixed_price NUMERIC(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    effective_from TIMESTAMP WITH TIME ZONE NOT NULL,
    effective_to TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_client_pricing_override ON client_pricing_overrides(client_id, product_variant_id);
CREATE INDEX idx_client_pricing_effective ON client_pricing_overrides(effective_from, effective_to);

COMMENT ON TABLE client_pricing_overrides IS 'Client-specific fixed price overrides for product variants';

-- ============================================================================
-- PRICE CALCULATIONS (Audit log)
-- ============================================================================
CREATE TABLE IF NOT EXISTS price_calculations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_variant_id UUID NOT NULL,
    client_id UUID NOT NULL,
    quantity INT NOT NULL,
    base_cost NUMERIC(12,2) NOT NULL,
    final_price NUMERIC(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    applied_rules JSONB NOT NULL DEFAULT '[]',
    calculated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_price_calculations_variant ON price_calculations(product_variant_id);
CREATE INDEX idx_price_calculations_client ON price_calculations(client_id);
CREATE INDEX idx_price_calculations_calculated_at ON price_calculations(calculated_at);

COMMENT ON TABLE price_calculations IS 'Audit log of price calculations with applied rules';

-- ============================================================================
-- TRIGGERS
-- ============================================================================
DROP TRIGGER IF EXISTS trg_base_sales_price_rules_updated_at ON base_sales_price_rules;
DROP TRIGGER IF EXISTS trg_sale_modification_rules_updated_at ON sale_modification_rules;
DROP TRIGGER IF EXISTS trg_client_pricing_overrides_updated_at ON client_pricing_overrides;
DROP TRIGGER IF EXISTS trg_price_calculations_updated_at ON price_calculations;

CREATE TRIGGER trg_base_sales_price_rules_updated_at BEFORE UPDATE ON base_sales_price_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_sale_modification_rules_updated_at BEFORE UPDATE ON sale_modification_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_client_pricing_overrides_updated_at BEFORE UPDATE ON client_pricing_overrides FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_price_calculations_updated_at BEFORE UPDATE ON price_calculations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
