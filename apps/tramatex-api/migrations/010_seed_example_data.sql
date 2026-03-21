-- ============================================================================
-- Migration: 010_seed_example_data.sql
-- Description: Populate database with realistic example data for TramaTex ERP
-- Date: 2026-03-06
-- Usage: Runs automatically via migrator, OR manually:
--        docker exec tramatex_db psql -U tramatex -d tramatex -f /app/migrations/010_seed_example_data.sql
-- Re-run: To re-seed, delete the migration record first:
--         DELETE FROM schema_migrations WHERE filename = '010_seed_example_data.sql';
-- ============================================================================

BEGIN;

-- ============================================================================
-- FIXED UUIDs (for repeatability and cross-reference)
-- ============================================================================

-- IAM
-- admin user already exists from 001: f47ac10b-58cc-4372-a567-0e02b2c3d479
-- Additional users
DO $$
BEGIN
    INSERT INTO users (id, email, password, role, is_active)
    VALUES
        ('a0000000-0000-0000-0000-000000000002', 'gerente@tramatex.local',
         '$2a$10$Y/OQADKTTu/8dOpA3BPc3eqOAQYREPJ03JWLuNKWpYApnFm5rl1oe', 'manager', true),
        ('a0000000-0000-0000-0000-000000000003', 'operario@tramatex.local',
         '$2a$10$Y/OQADKTTu/8dOpA3BPc3eqOAQYREPJ03JWLuNKWpYApnFm5rl1oe', 'operator', true),
        ('a0000000-0000-0000-0000-000000000004', 'cajera@tramatex.local',
         '$2a$10$Y/OQADKTTu/8dOpA3BPc3eqOAQYREPJ03JWLuNKWpYApnFm5rl1oe', 'cashier', true)
    ON CONFLICT (email) DO NOTHING;
END $$;

-- ============================================================================
-- PARTY MODULE
-- ============================================================================

-- Party: Organización cliente — Confecciones López S.L.
INSERT INTO parties (id, status, created_by, modified_by, default_discount_percentage)
VALUES ('b1000000-0000-0000-0000-000000000001', 'ACTIVE',
        'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 5.00)
ON CONFLICT (id) DO NOTHING;

INSERT INTO organization_profiles (party_id, name, tax_id, tax_id_type, phone, email)
VALUES ('b1000000-0000-0000-0000-000000000001', 'Confecciones López S.L.', 'B12345678', 'CIF',
        '+34 961 234 567', 'info@confeccioneslopez.es')
ON CONFLICT (party_id) DO NOTHING;

INSERT INTO party_roles (party_id, role) VALUES
    ('b1000000-0000-0000-0000-000000000001', 'CLIENT')
ON CONFLICT (party_id, role) DO NOTHING;

INSERT INTO party_addresses (id, party_id, street, city, province, postal_code, country, is_primary,
    created_by, modified_by)
VALUES ('b1a00000-0000-0000-0000-000000000001',
        'b1000000-0000-0000-0000-000000000001',
        'Calle Mayor 42', 'Valencia', 'Valencia', '46001', 'España', true,
        'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479')
ON CONFLICT (id) DO NOTHING;

-- Party: Organización proveedor — Textiles Martínez S.A.
INSERT INTO parties (id, status, created_by, modified_by, default_discount_percentage)
VALUES ('b2000000-0000-0000-0000-000000000001', 'ACTIVE',
        'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 0.00)
ON CONFLICT (id) DO NOTHING;

INSERT INTO organization_profiles (party_id, name, tax_id, tax_id_type, phone, email)
VALUES ('b2000000-0000-0000-0000-000000000001', 'Textiles Martínez S.A.', 'A87654321', 'CIF',
        '+34 963 456 789', 'compras@textilesmartinez.com')
ON CONFLICT (party_id) DO NOTHING;

INSERT INTO party_roles (party_id, role) VALUES
    ('b2000000-0000-0000-0000-000000000001', 'SUPPLIER')
ON CONFLICT (party_id, role) DO NOTHING;

-- Party: Organización cliente+proveedor — Moda Ibérica S.L.
INSERT INTO parties (id, status, created_by, modified_by, default_discount_percentage)
VALUES ('b3000000-0000-0000-0000-000000000001', 'ACTIVE',
        'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 10.00)
ON CONFLICT (id) DO NOTHING;

INSERT INTO organization_profiles (party_id, name, tax_id, tax_id_type, phone, email)
VALUES ('b3000000-0000-0000-0000-000000000001', 'Moda Ibérica S.L.', 'B11223344', 'CIF',
        '+34 915 678 901', 'pedidos@modaiberica.es')
ON CONFLICT (party_id) DO NOTHING;

INSERT INTO party_roles (party_id, role) VALUES
    ('b3000000-0000-0000-0000-000000000001', 'CLIENT'),
    ('b3000000-0000-0000-0000-000000000001', 'SUPPLIER')
ON CONFLICT (party_id, role) DO NOTHING;

INSERT INTO party_addresses (id, party_id, street, city, province, postal_code, country, is_primary,
    created_by, modified_by)
VALUES ('b3a00000-0000-0000-0000-000000000001',
        'b3000000-0000-0000-0000-000000000001',
        'Avda. de la Constitución 15', 'Madrid', 'Madrid', '28001', 'España', true,
        'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479')
ON CONFLICT (id) DO NOTHING;

-- Party: Persona física — Juan García (contacto de Confecciones López)
INSERT INTO parties (id, status, created_by, modified_by, default_discount_percentage)
VALUES ('b4000000-0000-0000-0000-000000000001', 'ACTIVE',
        'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 0.00)
ON CONFLICT (id) DO NOTHING;

INSERT INTO person_profiles (party_id, first_name, last_name, phone, email)
VALUES ('b4000000-0000-0000-0000-000000000001', 'Juan', 'García Pérez',
        '+34 612 345 678', 'jgarcia@confeccioneslopez.es')
ON CONFLICT (party_id) DO NOTHING;

-- Relación: Juan García es contacto de Confecciones López
INSERT INTO party_relationships (id, from_party_id, to_party_id, type, created_by)
VALUES ('b4r00000-0000-0000-0000-000000000001',
        'b1000000-0000-0000-0000-000000000001',
        'b4000000-0000-0000-0000-000000000001',
        'CONTACT_FOR',
        'f47ac10b-58cc-4372-a567-0e02b2c3d479')
ON CONFLICT (id) DO NOTHING;

-- Party: Persona física — María López (empleada)
INSERT INTO parties (id, status, created_by, modified_by, default_discount_percentage)
VALUES ('b5000000-0000-0000-0000-000000000001', 'ACTIVE',
        'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 0.00)
ON CONFLICT (id) DO NOTHING;

INSERT INTO person_profiles (party_id, first_name, last_name, phone, email)
VALUES ('b5000000-0000-0000-0000-000000000001', 'María', 'López Ruiz',
        '+34 655 987 321', 'mlopez@tramatex.local')
ON CONFLICT (party_id) DO NOTHING;

INSERT INTO party_roles (party_id, role) VALUES
    ('b5000000-0000-0000-0000-000000000001', 'EMPLOYEE')
ON CONFLICT (party_id, role) DO NOTHING;

-- ============================================================================
-- PRODUCT MODULE — BRANDS
-- ============================================================================

INSERT INTO brands (id, name, default_markup_percentage, is_active) VALUES
    ('c1000000-0000-0000-0000-000000000001', 'FyR Algodón',    0.10, true),
    ('c1000000-0000-0000-0000-000000000002', 'Seritex',         0.00, true),
    ('c1000000-0000-0000-0000-000000000003', 'ValenciaSport',  0.15, true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- PRODUCT MODULE — PRODUCT GROUPS
-- ============================================================================

INSERT INTO product_groups (id, name, parent_group_id, group_type, is_active) VALUES
    -- Top-level tangible groups
    ('c2000000-0000-0000-0000-000000000001', 'Camisetas',     NULL, 'TANGIBLE', true),
    ('c2000000-0000-0000-0000-000000000002', 'Pantalones',    NULL, 'TANGIBLE', true),
    ('c2000000-0000-0000-0000-000000000003', 'Ropa Deportiva', NULL, 'TANGIBLE', true),
    -- Top-level service groups
    ('c2000000-0000-0000-0000-000000000010', 'Serigrafía',    NULL, 'SERVICE', true),
    ('c2000000-0000-0000-0000-000000000011', 'Bordado',       NULL, 'SERVICE', true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- PRODUCT MODULE — ATTRIBUTES & VALUES
-- ============================================================================

-- Attribute: Talla
INSERT INTO attributes (id, name, code, sort_order) VALUES
    ('c3000000-0000-0000-0000-000000000001', 'Talla', 'T', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO attribute_values (id, attribute_id, value, code, has_price_modifier, modifier_type, modifier_amount) VALUES
    ('c3100000-0000-0000-0000-000000000001', 'c3000000-0000-0000-0000-000000000001', 'S',   'S',   false, NULL, NULL),
    ('c3100000-0000-0000-0000-000000000002', 'c3000000-0000-0000-0000-000000000001', 'M',   'M',   false, NULL, NULL),
    ('c3100000-0000-0000-0000-000000000003', 'c3000000-0000-0000-0000-000000000001', 'L',   'L',   false, NULL, NULL),
    ('c3100000-0000-0000-0000-000000000004', 'c3000000-0000-0000-0000-000000000001', 'XL',  'XL',  true,  'PERCENTAGE', 5.00),
    ('c3100000-0000-0000-0000-000000000005', 'c3000000-0000-0000-0000-000000000001', '2XL', '2XL', true,  'PERCENTAGE', 8.00),
    ('c3100000-0000-0000-0000-000000000006', 'c3000000-0000-0000-0000-000000000001', '3XL', '3XL', true,  'PERCENTAGE', 11.00)
ON CONFLICT (id) DO NOTHING;

-- Attribute: Color
INSERT INTO attributes (id, name, code, sort_order) VALUES
    ('c3000000-0000-0000-0000-000000000002', 'Color', 'C', 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO attribute_values (id, attribute_id, value, code, has_price_modifier, modifier_type, modifier_amount) VALUES
    ('c3200000-0000-0000-0000-000000000001', 'c3000000-0000-0000-0000-000000000002', 'Blanco',  'B',   true,  'FIXED', -0.20),
    ('c3200000-0000-0000-0000-000000000002', 'c3000000-0000-0000-0000-000000000002', 'Negro',   'N',   false, NULL,    NULL),
    ('c3200000-0000-0000-0000-000000000003', 'c3000000-0000-0000-0000-000000000002', 'Rojo',    'R',   true,  'FIXED', 0.30),
    ('c3200000-0000-0000-0000-000000000004', 'c3000000-0000-0000-0000-000000000002', 'Azul',    'AZ',  false, NULL,    NULL),
    ('c3200000-0000-0000-0000-000000000005', 'c3000000-0000-0000-0000-000000000002', 'Verde',   'VD',  false, NULL,    NULL)
ON CONFLICT (id) DO NOTHING;

-- Attribute: Material (para pantalones)
INSERT INTO attributes (id, name, code, sort_order) VALUES
    ('c3000000-0000-0000-0000-000000000003', 'Material', 'MAT', 3)
ON CONFLICT (id) DO NOTHING;

INSERT INTO attribute_values (id, attribute_id, value, code, has_price_modifier, modifier_type, modifier_amount) VALUES
    ('c3300000-0000-0000-0000-000000000001', 'c3000000-0000-0000-0000-000000000003', 'Algodón',   'ALG', false, NULL,    NULL),
    ('c3300000-0000-0000-0000-000000000002', 'c3000000-0000-0000-0000-000000000003', 'Poliéster', 'POL', true,  'FIXED', -0.50),
    ('c3300000-0000-0000-0000-000000000003', 'c3000000-0000-0000-0000-000000000003', 'Lycra',     'LYC', true,  'FIXED', 1.00)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- PRODUCT MODULE — PRODUCTS
-- ============================================================================

-- Producto 1: Camiseta Corta (FyR Algodón) — base_price 2.00€
INSERT INTO products (id, sku, name, long_name, description, product_type, brand_id, base_price, tax_rate,
    group_ids, direct_attribute_ids, is_active)
VALUES (
    'c4000000-0000-0000-0000-000000000001',
    'FYR2020',
    'Camiseta Corta',
    'Camiseta Corta Algodón 100%',
    'Camiseta manga corta de algodón 100%, ideal para serigrafía e impresión.',
    'TANGIBLE',
    'c1000000-0000-0000-0000-000000000001',
    2.00,
    21.00,
    ARRAY['c2000000-0000-0000-0000-000000000001']::UUID[],
    ARRAY['c3000000-0000-0000-0000-000000000001', 'c3000000-0000-0000-0000-000000000002']::UUID[],
    true
) ON CONFLICT (id) DO NOTHING;

-- Producto 2: Camiseta Larga (FyR Algodón) — base_price 3.50€
INSERT INTO products (id, sku, name, long_name, description, product_type, brand_id, base_price, tax_rate,
    group_ids, direct_attribute_ids, is_active)
VALUES (
    'c4000000-0000-0000-0000-000000000002',
    'FYR3050',
    'Camiseta Larga',
    'Camiseta Manga Larga Algodón 100%',
    'Camiseta manga larga de algodón premium.',
    'TANGIBLE',
    'c1000000-0000-0000-0000-000000000001',
    3.50,
    21.00,
    ARRAY['c2000000-0000-0000-0000-000000000001']::UUID[],
    ARRAY['c3000000-0000-0000-0000-000000000001', 'c3000000-0000-0000-0000-000000000002']::UUID[],
    true
) ON CONFLICT (id) DO NOTHING;

-- Producto 3: Pantalón Básico (ValenciaSport) — base_price 8.00€
INSERT INTO products (id, sku, name, long_name, description, product_type, brand_id, base_price, tax_rate,
    group_ids, direct_attribute_ids, is_active)
VALUES (
    'c4000000-0000-0000-0000-000000000003',
    'VAL4',
    'Pantalón Básico',
    'Pantalón de Chándal Básico',
    'Pantalón de chándal ligero para uso deportivo.',
    'TANGIBLE',
    'c1000000-0000-0000-0000-000000000003',
    8.00,
    21.00,
    ARRAY['c2000000-0000-0000-0000-000000000002', 'c2000000-0000-0000-0000-000000000003']::UUID[],
    ARRAY['c3000000-0000-0000-0000-000000000001', 'c3000000-0000-0000-0000-000000000002', 'c3000000-0000-0000-0000-000000000003']::UUID[],
    true
) ON CONFLICT (id) DO NOTHING;

-- Producto 4: Serigrafía (servicio — sin variantes físicas) — precio base por unidad
INSERT INTO products (id, sku, name, long_name, description, product_type, brand_id, base_price, tax_rate,
    group_ids, direct_attribute_ids, is_active)
VALUES (
    'c4000000-0000-0000-0000-000000000004',
    'SER',
    'Serigrafía',
    'Servicio de Serigrafía',
    'Servicio de impresión por serigrafía sobre textil.',
    'SERVICE',
    'c1000000-0000-0000-0000-000000000002',
    1.50,
    21.00,
    ARRAY['c2000000-0000-0000-0000-000000000010']::UUID[],
    ARRAY[]::UUID[],
    true
) ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- PRODUCT MODULE — PRODUCT VARIANTS
-- ============================================================================

-- Camiseta Corta (FYR2020): variantes representativas
INSERT INTO product_variants (id, product_id, sku, status, attribute_values, is_active) VALUES
    -- Talla S, colores
    ('c5000000-0000-0000-0000-000000000001', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.S-C.B',   'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000001', 'c3200000-0000-0000-0000-000000000001']::UUID[], true),
    ('c5000000-0000-0000-0000-000000000002', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.S-C.N',   'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000001', 'c3200000-0000-0000-0000-000000000002']::UUID[], true),
    ('c5000000-0000-0000-0000-000000000003', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.S-C.R',   'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000001', 'c3200000-0000-0000-0000-000000000003']::UUID[], true),
    -- Talla M
    ('c5000000-0000-0000-0000-000000000004', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.M-C.B',   'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000002', 'c3200000-0000-0000-0000-000000000001']::UUID[], true),
    ('c5000000-0000-0000-0000-000000000005', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.M-C.N',   'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000002', 'c3200000-0000-0000-0000-000000000002']::UUID[], true),
    ('c5000000-0000-0000-0000-000000000006', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.M-C.AZ',  'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000002', 'c3200000-0000-0000-0000-000000000004']::UUID[], true),
    -- Talla L
    ('c5000000-0000-0000-0000-000000000007', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.L-C.N',   'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000003', 'c3200000-0000-0000-0000-000000000002']::UUID[], true),
    ('c5000000-0000-0000-0000-000000000008', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.L-C.R',   'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000003', 'c3200000-0000-0000-0000-000000000003']::UUID[], true),
    -- Talla XL (5% surcharge)
    ('c5000000-0000-0000-0000-000000000009', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.XL-C.N',  'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000004', 'c3200000-0000-0000-0000-000000000002']::UUID[], true),
    -- Talla 3XL (11% surcharge) + Color Blanco (-0.20€) → baseCost = 2.02
    ('c5000000-0000-0000-0000-000000000010', 'c4000000-0000-0000-0000-000000000001', 'FYR2020-T.3XL-C.B', 'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000006', 'c3200000-0000-0000-0000-000000000001']::UUID[], true)
ON CONFLICT (id) DO NOTHING;

-- Camiseta Larga (FYR3050): unas pocas variantes
INSERT INTO product_variants (id, product_id, sku, status, attribute_values, is_active) VALUES
    ('c5000000-0000-0000-0000-000000000020', 'c4000000-0000-0000-0000-000000000002', 'FYR3050-T.M-C.N',   'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000002', 'c3200000-0000-0000-0000-000000000002']::UUID[], true),
    ('c5000000-0000-0000-0000-000000000021', 'c4000000-0000-0000-0000-000000000002', 'FYR3050-T.L-C.AZ',  'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000003', 'c3200000-0000-0000-0000-000000000004']::UUID[], true),
    ('c5000000-0000-0000-0000-000000000022', 'c4000000-0000-0000-0000-000000000002', 'FYR3050-T.XL-C.R',  'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000004', 'c3200000-0000-0000-0000-000000000003']::UUID[], true)
ON CONFLICT (id) DO NOTHING;

-- Pantalón Básico (VAL4): variantes con material
INSERT INTO product_variants (id, product_id, sku, status, attribute_values, is_active) VALUES
    ('c5000000-0000-0000-0000-000000000030', 'c4000000-0000-0000-0000-000000000003', 'VAL4-T.M-C.N-MAT.ALG',  'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000002', 'c3200000-0000-0000-0000-000000000002', 'c3300000-0000-0000-0000-000000000001']::UUID[], true),
    ('c5000000-0000-0000-0000-000000000031', 'c4000000-0000-0000-0000-000000000003', 'VAL4-T.L-C.AZ-MAT.POL', 'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000003', 'c3200000-0000-0000-0000-000000000004', 'c3300000-0000-0000-0000-000000000002']::UUID[], true),
    ('c5000000-0000-0000-0000-000000000032', 'c4000000-0000-0000-0000-000000000003', 'VAL4-T.XL-C.N-MAT.LYC', 'CONFIRMED',
     ARRAY['c3100000-0000-0000-0000-000000000004', 'c3200000-0000-0000-0000-000000000002', 'c3300000-0000-0000-0000-000000000003']::UUID[], true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- PRICING MODULE
-- ============================================================================

-- Base sales price rule: FyR Algodón — 10% markup (applies at brand level)
INSERT INTO base_sales_price_rules (id, name, brand_id, value_type, percentage_value, is_active)
VALUES ('d1000000-0000-0000-0000-000000000001',
        'Margen FyR Algodón 10%',
        'c1000000-0000-0000-0000-000000000001',
        'PERCENTAGE_MARKUP', 0.1000, true)
ON CONFLICT (id) DO NOTHING;

-- Base sales price rule: ValenciaSport — 15% markup
INSERT INTO base_sales_price_rules (id, name, brand_id, value_type, percentage_value, is_active)
VALUES ('d1000000-0000-0000-0000-000000000002',
        'Margen ValenciaSport 15%',
        'c1000000-0000-0000-0000-000000000003',
        'PERCENTAGE_MARKUP', 0.1500, true)
ON CONFLICT (id) DO NOTHING;

-- Sale modification rule: 5% discount for orders over 500€
INSERT INTO sale_modification_rules (id, name, value_type, percentage_value,
    min_order_total_amount, priority, effective_from, is_active)
VALUES ('d2000000-0000-0000-0000-000000000001',
        'Dto. 5% pedidos +500€',
        'APPLY_PERCENTAGE_DISCOUNT', 0.0500,
        500.00, 10, '2026-01-01', true)
ON CONFLICT (id) DO NOTHING;

-- Client pricing override: Moda Ibérica — precio fijo Camiseta Corta M Negro
INSERT INTO client_pricing_overrides (id, client_id, product_variant_id, fixed_price, currency,
    effective_from, is_active)
VALUES ('d3000000-0000-0000-0000-000000000001',
        'b3000000-0000-0000-0000-000000000001',
        'c5000000-0000-0000-0000-000000000005',
        1.75, 'EUR', '2026-01-01', true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- MES MODULE — MASTER DATA
-- ============================================================================

-- Tasks
INSERT INTO tasks (id, name, description, is_active) VALUES
    ('e1000000-0000-0000-0000-000000000001', 'Preparar marco',     'Preparar marco de serigrafía con fotolito', true),
    ('e1000000-0000-0000-0000-000000000002', 'Imprimir',           'Aplicar tinta sobre prenda',                true),
    ('e1000000-0000-0000-0000-000000000003', 'Secar',              'Secado en túnel de flash o plancha',        true),
    ('e1000000-0000-0000-0000-000000000004', 'Inspección calidad', 'Revisión visual de impresión',              true),
    ('e1000000-0000-0000-0000-000000000005', 'Digitalizar diseño', 'Digitalizar diseño para bordadora',         true),
    ('e1000000-0000-0000-0000-000000000006', 'Bordar',             'Ejecutar bordado en máquina',               true),
    ('e1000000-0000-0000-0000-000000000007', 'Empaquetar',         'Doblar, etiquetar y empaquetar',            true)
ON CONFLICT (id) DO NOTHING;

-- Positions (zonas de la prenda donde se aplica el marcado/arreglo)
INSERT INTO positions (id, name, code, description, is_active) VALUES
    ('e2000000-0000-0000-0000-000000000001', 'Pecho Izquierdo',   'PI',   'Zona pecho lado izquierdo',               true),
    ('e2000000-0000-0000-0000-000000000002', 'Pecho Derecho',     'PD',   'Zona pecho lado derecho',                 true),
    ('e2000000-0000-0000-0000-000000000003', 'Espalda',           'E',    'Zona central de la espalda',              true),
    ('e2000000-0000-0000-0000-000000000004', 'Manga Izquierda',   'MI',   'Manga lado izquierdo',                    true),
    ('e2000000-0000-0000-0000-000000000005', 'Manga Derecha',     'MD',   'Manga lado derecho',                      true),
    ('e2000000-0000-0000-0000-000000000006', 'Pernera Izquierda', 'PLI',  'Pernera lado izquierdo del pantalón',     true),
    ('e2000000-0000-0000-0000-000000000007', 'Pernera Derecha',   'PLD',  'Pernera lado derecho del pantalón',       true),
    ('e2000000-0000-0000-0000-000000000008', 'Bajos',             'BAJ',  'Bajos de pantalón o falda (dobladillo)',   true),
    ('e2000000-0000-0000-0000-000000000009', 'Cuello',            'CU',   'Zona del cuello o solapa',                true),
    ('e2000000-0000-0000-0000-000000000010', 'Cintura',           'CIN',  'Zona de cintura (ajustes, elásticos)',     true)
ON CONFLICT (id) DO NOTHING;

-- Service Groups (workflow templates)
INSERT INTO service_groups (id, name, description, product_group_id, is_active) VALUES
    ('e3000000-0000-0000-0000-000000000001', 'Proceso Serigrafía', 'Flujo completo de impresión por serigrafía',
     'c2000000-0000-0000-0000-000000000010', true),
    ('e3000000-0000-0000-0000-000000000002', 'Proceso Bordado',    'Flujo completo de bordado industrial',
     'c2000000-0000-0000-0000-000000000011', true)
ON CONFLICT (id) DO NOTHING;

-- Service Group Tasks (ordered steps)
INSERT INTO service_group_tasks (service_group_id, task_id, sequence) VALUES
    -- Serigrafía: Preparar → Imprimir → Secar → Inspección → Empaquetar
    ('e3000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000001', 1),
    ('e3000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000002', 2),
    ('e3000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000003', 3),
    ('e3000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000004', 4),
    ('e3000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000007', 5),
    -- Bordado: Digitalizar → Bordar → Inspección → Empaquetar
    ('e3000000-0000-0000-0000-000000000002', 'e1000000-0000-0000-0000-000000000005', 1),
    ('e3000000-0000-0000-0000-000000000002', 'e1000000-0000-0000-0000-000000000006', 2),
    ('e3000000-0000-0000-0000-000000000002', 'e1000000-0000-0000-0000-000000000004', 3),
    ('e3000000-0000-0000-0000-000000000002', 'e1000000-0000-0000-0000-000000000007', 4)
ON CONFLICT (service_group_id, task_id) DO NOTHING;

-- ============================================================================
-- MES MODULE — WORK ORDERS
-- ============================================================================

-- MES Work 1: Serigrafía para Confecciones López — 200 camisetas cortas
INSERT INTO mes_works (id, work_number, work_name, party_id, tangible_group_id, garment_notes,
    status, priority, start_date, due_date)
VALUES ('e4000000-0000-0000-0000-000000000001',
        'MES-2026-001',
        'Serigrafía camisetas López Primavera',
        'b1000000-0000-0000-0000-000000000001',
        'c2000000-0000-0000-0000-000000000001',
        'Logo frontal 20x15cm, 3 colores. Diseño adjunto en carpeta compartida.',
        'PENDING', 'NORMAL',
        '2026-03-10', '2026-03-25')
ON CONFLICT (id) DO NOTHING;

-- Work service group + tasks for MES-2026-001
INSERT INTO mes_work_service_groups (id, mes_work_id, service_group_id, position_id, notes, sequence)
VALUES ('e5000000-0000-0000-0000-000000000001',
        'e4000000-0000-0000-0000-000000000001',
        'e3000000-0000-0000-0000-000000000001',
        'e2000000-0000-0000-0000-000000000001',
        'Usar tintas base agua, secado a 160°C',
        1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO mes_work_tasks (id, mes_work_service_group_id, task_id, sequence, status) VALUES
    ('e6000000-0000-0000-0000-000000000001', 'e5000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000001', 1, 'COMPLETED'),
    ('e6000000-0000-0000-0000-000000000002', 'e5000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000002', 2, 'IN_PROGRESS'),
    ('e6000000-0000-0000-0000-000000000003', 'e5000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000003', 3, 'PENDING'),
    ('e6000000-0000-0000-0000-000000000004', 'e5000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000004', 4, 'PENDING'),
    ('e6000000-0000-0000-0000-000000000005', 'e5000000-0000-0000-0000-000000000001', 'e1000000-0000-0000-0000-000000000007', 5, 'PENDING')
ON CONFLICT (id) DO NOTHING;

-- MES Work 2: Bordado para Moda Ibérica — polos
INSERT INTO mes_works (id, work_number, work_name, party_id, tangible_group_id, garment_notes,
    status, priority, start_date, due_date)
VALUES ('e4000000-0000-0000-0000-000000000002',
        'MES-2026-002',
        'Bordado polos Moda Ibérica',
        'b3000000-0000-0000-0000-000000000001',
        'c2000000-0000-0000-0000-000000000001',
        'Logo pecho izquierdo 8x5cm, hilo dorado.',
        'DRAFT', 'HIGH',
        '2026-03-15', '2026-04-01')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- SALES MODULE — QUOTES
-- ============================================================================

-- Quote 1: Presupuesto para Confecciones López — 200 camisetas
INSERT INTO quotes (id, quote_number, party_id, quote_date, expiration_date, status, notes,
    subtotal_amount, subtotal_currency, tax_amount, tax_currency, total_amount, total_currency,
    mes_work_ids, created_at, updated_at)
VALUES ('f1000000-0000-0000-0000-000000000001',
        'PRE-2026-0001',
        'b1000000-0000-0000-0000-000000000001',
        '2026-03-01', '2026-03-31',
        'BORRADOR',
        'Presupuesto serigrafía camisetas primavera. Incluye trabajo MES-2026-001.',
        440.00, 'EUR',
        92.40, 'EUR',
        532.40, 'EUR',
        ARRAY['e4000000-0000-0000-0000-000000000001']::UUID[],
        NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Quote line items
INSERT INTO quote_line_items (id, quote_id, product_variant_id, quantity,
    calculated_unit_price_amount, calculated_unit_price_currency,
    final_unit_price_amount, final_unit_price_currency,
    tax_rate,
    final_discount_per_unit_amount, final_discount_per_unit_currency,
    subtotal_amount, subtotal_currency,
    tax_amount)
VALUES
    -- 100x Camiseta Corta M Negro @ 2.20€ (baseCost 2.00, brand 10% = 2.20)
    ('f1100000-0000-0000-0000-000000000001',
     'f1000000-0000-0000-0000-000000000001',
     'c5000000-0000-0000-0000-000000000005',
     100,
     2.20, 'EUR',
     2.20, 'EUR',
     21.00,
     0.00, 'EUR',
     220.00, 'EUR',
     46.20),
    -- 100x Camiseta Corta M Azul @ 2.20€
    ('f1100000-0000-0000-0000-000000000002',
     'f1000000-0000-0000-0000-000000000001',
     'c5000000-0000-0000-0000-000000000006',
     100,
     2.20, 'EUR',
     2.20, 'EUR',
     21.00,
     0.00, 'EUR',
     220.00, 'EUR',
     46.20)
ON CONFLICT (id) DO NOTHING;

-- Quote 2: Presupuesto para Moda Ibérica — pantalones
INSERT INTO quotes (id, quote_number, party_id, quote_date, expiration_date, status, notes,
    subtotal_amount, subtotal_currency, tax_amount, tax_currency, total_amount, total_currency,
    created_at, updated_at)
VALUES ('f1000000-0000-0000-0000-000000000002',
        'PRE-2026-0002',
        'b3000000-0000-0000-0000-000000000001',
        '2026-03-05', '2026-04-04',
        'EMITIDA',
        'Presupuesto pantalones deportivos temporada.',
        345.00, 'EUR',
        72.45, 'EUR',
        417.45, 'EUR',
        NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO quote_line_items (id, quote_id, product_variant_id, quantity,
    calculated_unit_price_amount, calculated_unit_price_currency,
    final_unit_price_amount, final_unit_price_currency,
    tax_rate,
    final_discount_per_unit_amount, final_discount_per_unit_currency,
    subtotal_amount, subtotal_currency,
    tax_amount)
VALUES
    -- 50x Pantalón M Negro Algodón @ 9.20€ (baseCost 8.00, brand 15% = 9.20)
    ('f1100000-0000-0000-0000-000000000010',
     'f1000000-0000-0000-0000-000000000002',
     'c5000000-0000-0000-0000-000000000030',
     50,
     9.20, 'EUR',
     9.20, 'EUR',
     21.00,
     0.00, 'EUR',
     460.00, 'EUR',
     96.60)
ON CONFLICT (id) DO NOTHING;

COMMIT;

-- ============================================================================
-- SEED DATA REFERENCE
-- ============================================================================
-- Users: admin@tramatex.local / gerente / operario / cajera (password: admin123)
--
-- Parties:
--   b1... = Confecciones López S.L.    (CLIENT, 5% discount)
--   b2... = Textiles Martínez S.A.     (SUPPLIER)
--   b3... = Moda Ibérica S.L.          (CLIENT+SUPPLIER, 10% discount)
--   b4... = Juan García Pérez          (persona física, contacto de López)
--   b5... = María López Ruiz           (persona física, EMPLOYEE)
--
-- Brands: FyR Algodón (10%), Seritex (0%), ValenciaSport (15%)
--
-- Products:
--   FYR2020  = Camiseta Corta     (2.00€, Talla+Color)       → 10 variantes
--   FYR3050  = Camiseta Larga     (3.50€, Talla+Color)       → 3 variantes
--   VAL4     = Pantalón Básico    (8.00€, Talla+Color+Mat.)  → 3 variantes
--   SER      = Serigrafía (serv.) (1.50€, sin variantes)
--
-- Pricing rules:
--   FyR Algodón  → +10% markup
--   ValenciaSport → +15% markup
--   Descuento 5% pedidos > 500€
--   Override: Moda Ibérica → Camiseta M Negro = 1.75€ fijo
--
-- MES:
--   MES-2026-001 = Serigrafía López (PENDING, con tareas activas)
--   MES-2026-002 = Bordado Moda Ibérica (DRAFT)
--
-- Quotes:
--   PRE-2026-0001 = López, 200 camisetas (DRAFT, vinculado a MES-2026-001)
--   PRE-2026-0002 = Moda Ibérica, 50 pantalones (SENT)
-- ============================================================================
