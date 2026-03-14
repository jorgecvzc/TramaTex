-- ============================================================================
-- Migration: v2_003_init_product.sql
-- Description: Initialize Product module
-- Date: 2026-02-25
-- Modules: Brands, Product Groups, Attributes, Attribute Values, Products, Product Variants
-- ============================================================================

BEGIN;

-- ============================================================================
-- ENUMS
-- ============================================================================
DO $$ BEGIN
    CREATE TYPE product_type AS ENUM ('TANGIBLE', 'SERVICE');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE product_group_type AS ENUM ('TANGIBLE', 'SERVICE');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE variant_status AS ENUM ('PROVISIONAL', 'CONFIRMED');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE modifier_type AS ENUM ('FIXED', 'PERCENTAGE');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- ============================================================================
-- BRANDS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    default_markup_percentage NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_brands_name ON brands(name);

COMMENT ON TABLE brands IS 'Product brands catalog';
COMMENT ON COLUMN brands.default_markup_percentage IS 'Default markup percentage applied to product variants (e.g., 30.0 = 30%)';

-- ============================================================================
-- PRODUCT GROUPS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS product_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    parent_group_id UUID,
    group_type product_group_type NOT NULL DEFAULT 'TANGIBLE',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_parent_group FOREIGN KEY (parent_group_id) REFERENCES product_groups(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_product_groups_parent_group_id ON product_groups(parent_group_id);
CREATE INDEX IF NOT EXISTS idx_product_groups_group_type ON product_groups(group_type);

COMMENT ON TABLE product_groups IS 'Hierarchical product groups';
COMMENT ON COLUMN product_groups.group_type IS 'Classification: TANGIBLE for physical products, SERVICE for service-based products';

-- ============================================================================
-- ATTRIBUTES TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS attributes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    scope_brand_id UUID,
    scope_group_id UUID,
    created_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    modified_by VARCHAR(255),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_scope_brand FOREIGN KEY (scope_brand_id) REFERENCES brands(id) ON DELETE CASCADE,
    CONSTRAINT fk_scope_product_group FOREIGN KEY (scope_group_id) REFERENCES product_groups(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_attributes_code ON attributes(code);
CREATE INDEX IF NOT EXISTS idx_attributes_scope ON attributes(scope_brand_id, scope_group_id);

COMMENT ON TABLE attributes IS 'Product attributes (Size, Color, Material, etc.)';
COMMENT ON COLUMN attributes.scope_brand_id IS 'Optional: restrict attribute to specific brand';
COMMENT ON COLUMN attributes.scope_group_id IS 'Optional: restrict attribute to specific product group';

-- ============================================================================
-- ATTRIBUTE VALUES TABLE (with price modifiers)
-- ============================================================================
CREATE TABLE IF NOT EXISTS attribute_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attribute_id UUID NOT NULL,
    value VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    has_price_modifier BOOLEAN NOT NULL DEFAULT FALSE,
    modifier_type modifier_type,
    modifier_amount NUMERIC(10,2),
    created_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    modified_by VARCHAR(255),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_attribute FOREIGN KEY (attribute_id) REFERENCES attributes(id) ON DELETE CASCADE,
    CONSTRAINT chk_price_modifier_consistency CHECK (
        (has_price_modifier = FALSE AND modifier_type IS NULL AND modifier_amount IS NULL) OR
        (has_price_modifier = TRUE AND modifier_type IS NOT NULL AND modifier_amount IS NOT NULL)
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_attribute_values_attribute_id_code ON attribute_values(attribute_id, code);
CREATE INDEX IF NOT EXISTS idx_attribute_values_price_modifier ON attribute_values(has_price_modifier) WHERE has_price_modifier = TRUE;

COMMENT ON TABLE attribute_values IS 'Values for product attributes with optional price modifiers';
COMMENT ON COLUMN attribute_values.has_price_modifier IS 'Whether this value affects product variant price';
COMMENT ON COLUMN attribute_values.modifier_type IS 'Type: FIXED (absolute ±EUR) or PERCENTAGE (±%)';
COMMENT ON COLUMN attribute_values.modifier_amount IS 'Modifier value (5.00 = +5 EUR or +5%)';

-- ============================================================================
-- PRODUCTS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(255) UNIQUE,
    name VARCHAR(100) NOT NULL,
    long_name VARCHAR(255),
    barcode VARCHAR(255) UNIQUE,
    description TEXT,
    product_type product_type NOT NULL,
    brand_id UUID NOT NULL,
    base_price NUMERIC(12,2) NOT NULL DEFAULT 0,
    tax_rate NUMERIC(5,2) NOT NULL DEFAULT 21.00,
    group_ids UUID[],
    direct_attribute_ids UUID[],
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    modified_by VARCHAR(255),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_products_brand FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_products_brand_id ON products(brand_id);

COMMENT ON TABLE products IS 'Products catalog (templates for variants)';
COMMENT ON COLUMN products.base_price IS 'Base cost/price - source of truth for variant pricing';
COMMENT ON COLUMN products.tax_rate IS 'Tax rate (VAT/IVA) percentage (e.g., 21.0 = 21% IVA)';

-- ============================================================================
-- PRODUCT VARIANTS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL,
    sku VARCHAR(255) NOT NULL,
    barcode VARCHAR(255),
    status variant_status NOT NULL DEFAULT 'PROVISIONAL',
    attribute_values UUID[],
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    modified_by VARCHAR(255),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_product_variants_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_variants_sku ON product_variants(sku);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_variants_barcode ON product_variants(barcode) WHERE barcode IS NOT NULL;

COMMENT ON TABLE product_variants IS 'Product variants - concrete products with specific attributes';
COMMENT ON COLUMN product_variants.attribute_values IS 'Array of attribute_value IDs that define this variant';
COMMENT ON COLUMN product_variants.status IS 'PROVISIONAL (JIT-created) or CONFIRMED (pre-generated)';

-- ============================================================================
-- TRIGGERS
-- ============================================================================
DROP TRIGGER IF EXISTS trg_brands_updated_at ON brands;
DROP TRIGGER IF EXISTS trg_product_groups_updated_at ON product_groups;
DROP TRIGGER IF EXISTS trg_attributes_updated_at ON attributes;
DROP TRIGGER IF EXISTS trg_attribute_values_updated_at ON attribute_values;
DROP TRIGGER IF EXISTS trg_products_updated_at ON products;
DROP TRIGGER IF EXISTS trg_product_variants_updated_at ON product_variants;

CREATE TRIGGER trg_brands_updated_at BEFORE UPDATE ON brands FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_product_groups_updated_at BEFORE UPDATE ON product_groups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_attributes_updated_at BEFORE UPDATE ON attributes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_attribute_values_updated_at BEFORE UPDATE ON attribute_values FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_product_variants_updated_at BEFORE UPDATE ON product_variants FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;

-- ============================================================================
-- END OF MIGRATION: v2_003_init_product.sql
-- ============================================================================
