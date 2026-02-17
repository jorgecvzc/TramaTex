-- Seed master data for Product Module testing
-- Date: 2026-02-10
-- Purpose: Provide initial data for brands, groups, and attributes

BEGIN;

-- ============================================================================
-- BRANDS
-- ============================================================================
INSERT INTO brands (id, name, is_active, created_at, updated_at) VALUES
('11111111-1111-1111-1111-111111111111', 'Nike', true, NOW(), NOW()),
('22222222-2222-2222-2222-222222222222', 'Adidas', true, NOW(), NOW()),
('33333333-3333-3333-3333-333333333333', 'Puma', true, NOW(), NOW()),
('44444444-4444-4444-4444-444444444444', 'Reebok', true, NOW(), NOW()),
('55555555-5555-5555-5555-555555555555', 'Under Armour', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- PRODUCT GROUPS (Categories)
-- ============================================================================
INSERT INTO product_groups (id, name, is_active, created_at, updated_at) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Calzado Deportivo', true, NOW(), NOW()),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Ropa Deportiva', true, NOW(), NOW()),
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'Accesorios', true, NOW(), NOW()),
('dddddddd-dddd-dddd-dddd-dddddddddddd', 'Equipamiento', true, NOW(), NOW()),
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Textiles', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- ATTRIBUTES (Configurable options) 
-- ============================================================================

-- Generic attributes (no scope - apply to all)
INSERT INTO attributes (id, name, code, sort_order, scope_brand_id, scope_group_id, created_at, updated_at) VALUES
('a0000000-0000-0000-0000-000000000001', 'Talla', 'SIZE', 1, NULL, NULL, NOW(), NOW()),
('a0000000-0000-0000-0000-000000000002', 'Color', 'COLOR', 2, NULL, NULL, NOW(), NOW()),
('a0000000-0000-0000-0000-000000000003', 'Material', 'MATERIAL', 3, NULL, NULL, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Brand-specific attributes (Nike only)
INSERT INTO attributes (id, name, code, sort_order, scope_brand_id, scope_group_id, created_at, updated_at) VALUES
('a0000000-0000-0000-0000-000000000010', 'Tecnología Nike', 'TECH_NIKE', 10, '11111111-1111-1111-1111-111111111111', NULL, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Group-specific attributes (Calzado only)
INSERT INTO attributes (id, name, code, sort_order, scope_brand_id, scope_group_id, created_at, updated_at) VALUES
('a0000000-0000-0000-0000-000000000020', 'Tipo de Suela', 'SOLE_TYPE', 20, NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW(), NOW()),
('a0000000-0000-0000-0000-000000000021', 'Amortiguación', 'CUSHIONING', 21, NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Brand + Group specific (Nike Calzado)
INSERT INTO attributes (id, name, code, sort_order, scope_brand_id, scope_group_id, created_at, updated_at) VALUES
('a0000000-0000-0000-0000-000000000030', 'Línea Nike Running', 'NIKE_RUNNING_LINE', 30, '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- ATTRIBUTE VALUES
-- ============================================================================

-- Talla values
INSERT INTO attribute_values (id, attribute_id, value, code, created_at, updated_at) VALUES
('10000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'XS', 'XS', NOW(), NOW()),
('10000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'S', 'S', NOW(), NOW()),
('10000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'M', 'M', NOW(), NOW()),
('10000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'L', 'L', NOW(), NOW()),
('10000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 'XL', 'XL', NOW(), NOW()),
('10000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000001', 'XXL', 'XXL', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Color values
INSERT INTO attribute_values (id, attribute_id, value, code, created_at, updated_at) VALUES
('20000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000002', 'Negro', 'BLACK', NOW(), NOW()),
('20000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000002', 'Blanco', 'WHITE', NOW(), NOW()),
('20000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000002', 'Rojo', 'RED', NOW(), NOW()),
('20000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000002', 'Azul', 'BLUE', NOW(), NOW()),
('20000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000002', 'Verde', 'GREEN', NOW(), NOW()),
('20000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000002', 'Amarillo', 'YELLOW', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Material values
INSERT INTO attribute_values (id, attribute_id, value, code, created_at, updated_at) VALUES
('30000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000003', 'Algodón', 'COTTON', NOW(), NOW()),
('30000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000003', 'Poliéster', 'POLYESTER', NOW(), NOW()),
('30000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000003', 'Cuero', 'LEATHER', NOW(), NOW()),
('30000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000003', 'Sintético', 'SYNTHETIC', NOW(), NOW()),
('30000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000003', 'Mesh', 'MESH', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Tecnología Nike values
INSERT INTO attribute_values (id, attribute_id, value, code, created_at, updated_at) VALUES
('40000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000010', 'Air Max', 'AIR_MAX', NOW(), NOW()),
('40000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000010', 'Zoom', 'ZOOM', NOW(), NOW()),
('40000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000010', 'React', 'REACT', NOW(), NOW()),
('40000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000010', 'Flyknit', 'FLYKNIT', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Tipo de Suela values
INSERT INTO attribute_values (id, attribute_id, value, code, created_at, updated_at) VALUES
('50000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000020', 'Goma', 'RUBBER', NOW(), NOW()),
('50000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000020', 'EVA', 'EVA', NOW(), NOW()),
('50000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000020', 'Caucho', 'CAOUTCHOUC', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Amortiguación values
INSERT INTO attribute_values (id, attribute_id, value, code, created_at, updated_at) VALUES
('60000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000021', 'Baja', 'LOW', NOW(), NOW()),
('60000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000021', 'Media', 'MEDIUM', NOW(), NOW()),
('60000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000021', 'Alta', 'HIGH', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Línea Nike Running values
INSERT INTO attribute_values (id, attribute_id, value, code, created_at, updated_at) VALUES
('70000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000030', 'Pegasus', 'PEGASUS', NOW(), NOW()),
('70000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000030', 'Vaporfly', 'VAPORFLY', NOW(), NOW()),
('70000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000030', 'Structure', 'STRUCTURE', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

COMMIT;

-- Verification queries (commented out)
-- SELECT 'Brands:', COUNT(*) FROM brands;
-- SELECT 'Product Groups:', COUNT(*) FROM product_groups;
-- SELECT 'Attributes:', COUNT(*) FROM attributes;
-- SELECT 'Attribute Values:', COUNT(*) FROM attribute_values;
