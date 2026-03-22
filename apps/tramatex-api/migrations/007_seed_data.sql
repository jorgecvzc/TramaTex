-- ============================================================================
-- Migration: 007_seed_data.sql
-- Description: Sample/development data for TramaTex
-- Generated from DB state: 2026-03-22
-- Replaces: 008_seed_test_data.sql
-- ============================================================================
-- NOTE: product_groups use fixed UUIDs referenced by products (group_ids array).
--       UUIDs are stable so re-seeding is idempotent (all ON CONFLICT DO NOTHING).
-- ============================================================================

BEGIN;

-- ============================================================================
-- IAM — Admin user (also seeded in 001, ON CONFLICT is idempotent)
-- Password: admin123
-- ============================================================================
INSERT INTO users (id, email, password, role, is_active)
VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'admin@tramatex.local',
    '$2a$10$Y/OQADKTTu/8dOpA3BPc3eqOAQYREPJ03JWLuNKWpYApnFm5rl1oe',
    'admin',
    true
)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PRODUCT GROUPS
-- Fixed UUIDs matching the group_ids arrays in products below.
-- Root tangible groups first (self-FK), then children.
-- ============================================================================
-- Root: Prendas Textiles
INSERT INTO product_groups (id, name, parent_group_id, group_type, is_active)
VALUES ('a1111111-0001-0001-0001-000000000001', 'Prendas Textiles', NULL, 'TANGIBLE', true)
ON CONFLICT DO NOTHING;

-- Children of Prendas Textiles
INSERT INTO product_groups (id, name, parent_group_id, group_type, is_active)
VALUES ('99dc48c8-cc5a-431d-9e56-0e60edcd8005', 'Camisetas', 'a1111111-0001-0001-0001-000000000001', 'TANGIBLE', true)
ON CONFLICT DO NOTHING;

INSERT INTO product_groups (id, name, parent_group_id, group_type, is_active)
VALUES ('7f62b3ee-2772-4b52-83ff-2c54e92dda5c', 'Pantalones', 'a1111111-0001-0001-0001-000000000001', 'TANGIBLE', true)
ON CONFLICT DO NOTHING;

INSERT INTO product_groups (id, name, parent_group_id, group_type, is_active)
VALUES ('a1111111-0001-0001-0001-000000000002', 'Polos', 'a1111111-0001-0001-0001-000000000001', 'TANGIBLE', true)
ON CONFLICT DO NOTHING;

-- Root: Marcados (service)
INSERT INTO product_groups (id, name, parent_group_id, group_type, is_active)
VALUES ('75e68c9d-f08f-47c5-b8bd-5076660c4ada', 'Marcados', NULL, 'SERVICE', true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- BRANDS
-- ============================================================================
INSERT INTO brands (id, name, default_markup_percentage, is_active) VALUES
    ('29a2f611-a70c-4b4f-80c5-fd8f54626a3c', 'Fall&Ress', 25.00, true),
    ('30a3e38e-fcad-4a9e-8fd8-e7bd2051d708', 'Valross',   55.00, true),
    ('0ad268be-b248-4657-ac94-bdb1a05ecbd8', 'TramaTex',   0.00, true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- ATTRIBUTES (no sort_order column — dropped in 003)
-- ============================================================================
INSERT INTO attributes (id, name, code) VALUES
    ('1418245c-50a5-49ed-8400-bd826438d1a2', 'Color',       'C'),
    ('d1cff03a-b8b2-4027-b310-755dae65f2c0', 'Talla',       'T'),
    ('54c542c1-6ddb-4c0b-bf78-b5c90a3a336b', 'Tipo Tejido', 'TJ')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- ATTRIBUTE VALUES
-- ============================================================================
-- Color values
INSERT INTO attribute_values (id, attribute_id, value, code, has_price_modifier, modifier_type, modifier_amount) VALUES
    ('684f9aef-5ea6-4864-8c88-dc4a0cbf30c0', '1418245c-50a5-49ed-8400-bd826438d1a2', 'Blanco',      'BL', true,  'FIXED',      -0.20),
    ('2eae3a1d-de46-4932-b4bc-dc8987928917', '1418245c-50a5-49ed-8400-bd826438d1a2', 'Rojo',        'RJ', false, NULL,          NULL),
    ('c927a195-ab0e-4573-9035-840fe9b5cc65', '1418245c-50a5-49ed-8400-bd826438d1a2', 'Azul',        'AZ', false, NULL,          NULL),
    ('52974604-a98f-4b53-8255-00e34c1cd59c', '1418245c-50a5-49ed-8400-bd826438d1a2', 'Verde kelly', 'VK', false, NULL,          NULL),
    ('b553875b-787e-4e2b-9088-5bed91191446', '1418245c-50a5-49ed-8400-bd826438d1a2', 'Negro',       'NG', false, NULL,          NULL)
ON CONFLICT DO NOTHING;

-- Talla values
INSERT INTO attribute_values (id, attribute_id, value, code, has_price_modifier, modifier_type, modifier_amount) VALUES
    ('ec2889f5-7902-496c-972f-d2299da913b7', 'd1cff03a-b8b2-4027-b310-755dae65f2c0', '2XS', '2XS', false, NULL,         NULL),
    ('4d2fcfba-49d5-4597-bdbc-dd786efe2c82', 'd1cff03a-b8b2-4027-b310-755dae65f2c0', 'XS',  'XS',  false, NULL,         NULL),
    ('ced66ae9-7c55-4557-82ca-5b1e9cebdd8b', 'd1cff03a-b8b2-4027-b310-755dae65f2c0', 'S',   'S',   false, NULL,         NULL),
    ('39854811-247a-4947-a947-ba20841fba5c', 'd1cff03a-b8b2-4027-b310-755dae65f2c0', 'M',   'M',   false, NULL,         NULL),
    ('72378d06-a594-4373-a145-cf1bc2f3a96d', 'd1cff03a-b8b2-4027-b310-755dae65f2c0', 'L',   'L',   false, NULL,         NULL),
    ('2dd7dcb5-55bc-4f29-a8bc-222395ac615f', 'd1cff03a-b8b2-4027-b310-755dae65f2c0', 'XL',  'XL',  false, NULL,         NULL),
    ('1158a4b7-0e08-4b36-9a97-373a4c15b0a0', 'd1cff03a-b8b2-4027-b310-755dae65f2c0', '2XL', '2XL', false, NULL,         NULL),
    ('e1f77b4d-6204-465c-ab7e-64a12bbd6375', 'd1cff03a-b8b2-4027-b310-755dae65f2c0', '3XL', '3XL', true,  'PERCENTAGE', 5.00),
    ('887f0897-fdcb-4054-b162-12aa9b6476b5', 'd1cff03a-b8b2-4027-b310-755dae65f2c0', '4XL', '4XL', true,  'FIXED',      8.00)
ON CONFLICT DO NOTHING;

-- Tipo Tejido values
INSERT INTO attribute_values (id, attribute_id, value, code, has_price_modifier, modifier_type, modifier_amount) VALUES
    ('ee934d4d-b8f7-4a6f-8309-d6c8e4affca8', '54c542c1-6ddb-4c0b-bf78-b5c90a3a336b', 'Algodón',  'AL', false, NULL, NULL),
    ('98c36a68-2915-4567-8add-9bc4011d4e9a', '54c542c1-6ddb-4c0b-bf78-b5c90a3a336b', 'Poliéster','PO', false, NULL, NULL),
    ('2f6c99cf-070d-486f-b113-84c695cdea82', '54c542c1-6ddb-4c0b-bf78-b5c90a3a336b', 'Sarga',    'SA', false, NULL, NULL),
    ('0ac79f3e-ffa9-4ed4-9e4c-69823b0b2f3e', '54c542c1-6ddb-4c0b-bf78-b5c90a3a336b', 'Popelín',  'PP', false, NULL, NULL)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PARTIES + PROFILES + ROLES
-- NOTE: CONSUMIDOR FINAL (00000000-...) also seeded in 002, ON CONFLICT is idempotent.
-- ============================================================================
INSERT INTO parties (id, status, default_discount_percentage, created_by, modified_by)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'ACTIVE',  0.00, 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479'),
    ('d4a5e6f7-0001-0001-0001-000000000001', 'ACTIVE',  0.00, 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479'),
    ('d4a5e6f7-0002-0002-0002-000000000002', 'ACTIVE', 15.00, 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479'),
    ('d4a5e6f7-0003-0003-0003-000000000003', 'ACTIVE',  0.00, 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479')
ON CONFLICT DO NOTHING;

INSERT INTO organization_profiles (party_id, name, tax_id, tax_id_type, website, phone, email, notes) VALUES
    ('00000000-0000-0000-0000-000000000001', 'CONSUMIDOR FINAL',             NULL,        NULL,  '',                    NULL, NULL, ''),
    ('d4a5e6f7-0001-0001-0001-000000000001', 'Juan Ramón Orzola',            '19230061N', 'NIF', '',                    NULL, NULL, ''),
    ('d4a5e6f7-0002-0002-0002-000000000002', 'TexPrendar SL',                'B12345768', 'CIF', 'https://texprendar.com', NULL, NULL, 'Cliente Principal'),
    ('d4a5e6f7-0003-0003-0003-000000000003', 'Pedro Ramón',                  '19652658M', 'NIF', '',                    NULL, NULL, '')
ON CONFLICT DO NOTHING;

INSERT INTO party_roles (party_id, role) VALUES
    ('00000000-0000-0000-0000-000000000001', 'CLIENT'),
    ('d4a5e6f7-0001-0001-0001-000000000001', 'CLIENT'),
    ('d4a5e6f7-0002-0002-0002-000000000002', 'CLIENT'),
    ('d4a5e6f7-0002-0002-0002-000000000002', 'SUPPLIER'),
    ('d4a5e6f7-0003-0003-0003-000000000003', 'CONTACT')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PRODUCTS
-- group_ids use the stable UUIDs defined in product_groups above.
-- direct_attribute_ids define attribute order for SKU generation.
-- ============================================================================
INSERT INTO products (id, sku, name, long_name, description, product_type, brand_id, base_price, tax_rate, group_ids, direct_attribute_ids, is_active)
VALUES
    -- Fall&Ress camisetas (Camisetas group)
    ('d93a499a-472c-460f-905b-1b446250a941', 'FYR2020', 'Camiseta Clásica MC', 'Camiseta Clásica 100% Algodón Manga Corta',
     'Camiseta de tejido algodón, corte clásico y manga corta',
     'TANGIBLE', '29a2f611-a70c-4b4f-80c5-fd8f54626a3c', 2.21, 21.00,
     '{99dc48c8-cc5a-431d-9e56-0e60edcd8005}',
     '{1418245c-50a5-49ed-8400-bd826438d1a2,d1cff03a-b8b2-4027-b310-755dae65f2c0}',
     true),
    ('158097ea-e2ca-40dd-9584-a3bd42ed0c0a', 'FYR2021', 'Camiseta Clásica ML', 'Camiseta Clásica 100% Algodón ML',
     'Camiseta de tejido algodón, corte clásico y manga larga',
     'TANGIBLE', '29a2f611-a70c-4b4f-80c5-fd8f54626a3c', 0.00, 21.00,
     '{99dc48c8-cc5a-431d-9e56-0e60edcd8005}',
     '{1418245c-50a5-49ed-8400-bd826438d1a2,d1cff03a-b8b2-4027-b310-755dae65f2c0}',
     true),
    -- Fall&Ress polos (Polos group)
    ('e18b77c8-56f9-4d91-bbe0-92fa931bac48', 'FYR2040', 'Polo Básico MC', 'Polo Básico MC',
     '',
     'TANGIBLE', '29a2f611-a70c-4b4f-80c5-fd8f54626a3c', 4.00, 21.00,
     '{a1111111-0001-0001-0001-000000000002}',
     '{1418245c-50a5-49ed-8400-bd826438d1a2,d1cff03a-b8b2-4027-b310-755dae65f2c0}',
     true),
    -- Valross pantalones (Pantalones group)
    ('91ff4c67-2810-49fb-96b0-9a8111ff8baa', 'VALPALB', 'Pantalón Laboral Multitejido', 'Pantalón Laboral Multitejido',
     'Pantalón laboral básico varios tejidos.',
     'TANGIBLE', '30a3e38e-fcad-4a9e-8fd8-e7bd2051d708', 6.00, 21.00,
     '{7f62b3ee-2772-4b52-83ff-2c54e92dda5c}',
     '{54c542c1-6ddb-4c0b-bf78-b5c90a3a336b,1418245c-50a5-49ed-8400-bd826438d1a2,d1cff03a-b8b2-4027-b310-755dae65f2c0}',
     true),
    -- TramaTex servicios (Marcados group)
    ('2f094614-4244-451f-8daf-0507df6eb281', 'SER', 'Serigrafía', 'Serigrafía',
     '',
     'SERVICE', '0ad268be-b248-4657-ac94-bdb1a05ecbd8', 1.50, 21.00,
     '{75e68c9d-f08f-47c5-b8bd-5076660c4ada}',
     '{}',
     true),
    ('25729ded-f598-43df-978c-1a8f65348c74', 'BOR', 'Bordado', 'Bordado',
     '',
     'SERVICE', '0ad268be-b248-4657-ac94-bdb1a05ecbd8', 3.50, 21.00,
     '{75e68c9d-f08f-47c5-b8bd-5076660c4ada}',
     '{}',
     true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PRODUCT VARIANTS
-- attribute_values arrays preserve the original creation order.
-- SKUs prefixed with C (Color) then T (Talla) per direct_attribute_ids order.
-- ============================================================================

-- FYR2020 variants (Camiseta Clásica MC — base_price 2.21 EUR)
INSERT INTO product_variants (id, product_id, sku, status, attribute_values, is_active) VALUES
    ('f298ef9d-31b4-4736-a8cb-3abb47de2a58', 'd93a499a-472c-460f-905b-1b446250a941', 'FYR2020-T.M-C.RJ',   'PROVISIONAL', '{2eae3a1d-de46-4932-b4bc-dc8987928917,39854811-247a-4947-a947-ba20841fba5c}', true),
    ('00104c1d-b541-412e-a09e-6716c9a0e8e9', 'd93a499a-472c-460f-905b-1b446250a941', 'FYR2020-C.AZ-T.XL',  'PROVISIONAL', '{2dd7dcb5-55bc-4f29-a8bc-222395ac615f,c927a195-ab0e-4573-9035-840fe9b5cc65}', true),
    ('ff911867-580e-443c-9868-deee73625c6c', 'd93a499a-472c-460f-905b-1b446250a941', 'FYR2020-C.BL-T.M',   'PROVISIONAL', '{39854811-247a-4947-a947-ba20841fba5c,684f9aef-5ea6-4864-8c88-dc4a0cbf30c0}', true),
    ('c674368b-b870-4e66-94d1-ec7d9843210b', 'd93a499a-472c-460f-905b-1b446250a941', 'FYR2020-C.NG-T.3XL', 'PROVISIONAL', '{b553875b-787e-4e2b-9088-5bed91191446,e1f77b4d-6204-465c-ab7e-64a12bbd6375}', true),
    ('35f1f8cd-f3a7-4bf3-9b90-d35a427c60c9', 'd93a499a-472c-460f-905b-1b446250a941', 'FYR2020-C.AZ-T.S',   'PROVISIONAL', '{c927a195-ab0e-4573-9035-840fe9b5cc65,ced66ae9-7c55-4557-82ca-5b1e9cebdd8b}', true),
    ('187ca801-71f4-416f-af22-063932620852', 'd93a499a-472c-460f-905b-1b446250a941', 'FYR2020-C.RJ-T.L',   'PROVISIONAL', '{2eae3a1d-de46-4932-b4bc-dc8987928917,72378d06-a594-4373-a145-cf1bc2f3a96d}', true),
    ('ff08f0db-c3c4-4b58-94ed-9086a993487b', 'd93a499a-472c-460f-905b-1b446250a941', 'FYR2020-C.BL-T.3XL', 'PROVISIONAL', '{684f9aef-5ea6-4864-8c88-dc4a0cbf30c0,e1f77b4d-6204-465c-ab7e-64a12bbd6375}', true),
    ('813a0857-8279-417a-b4cb-4cae5e3eae3c', 'd93a499a-472c-460f-905b-1b446250a941', 'FYR2020-C.BL-T.L',   'PROVISIONAL', '{684f9aef-5ea6-4864-8c88-dc4a0cbf30c0,72378d06-a594-4373-a145-cf1bc2f3a96d}', true)
ON CONFLICT DO NOTHING;

-- FYR2021 variants (Camiseta Clásica ML — base_price 0.00 EUR)
INSERT INTO product_variants (id, product_id, sku, status, attribute_values, is_active) VALUES
    ('149ad077-b8be-4e8f-bc9c-4566eb83d919', '158097ea-e2ca-40dd-9584-a3bd42ed0c0a', 'FYR2021-C.RJ-T.3XL', 'PROVISIONAL', '{2eae3a1d-de46-4932-b4bc-dc8987928917,e1f77b4d-6204-465c-ab7e-64a12bbd6375}', true)
ON CONFLICT DO NOTHING;

-- FYR2040 variants (Polo Básico MC — base_price 4.00 EUR)
INSERT INTO product_variants (id, product_id, sku, status, attribute_values, is_active) VALUES
    ('7f3bb2bd-5b07-4e80-bf9e-28af747981a1', 'e18b77c8-56f9-4d91-bbe0-92fa931bac48', 'FYR2040-C.RJ-T.L',   'CONFIRMED',   '{2eae3a1d-de46-4932-b4bc-dc8987928917,72378d06-a594-4373-a145-cf1bc2f3a96d}', true),
    ('aa0f7b0a-7873-45fb-b848-96684dcfbac1', 'e18b77c8-56f9-4d91-bbe0-92fa931bac48', 'FYR2040-C.BL-T.3XL', 'CONFIRMED',   '{684f9aef-5ea6-4864-8c88-dc4a0cbf30c0,e1f77b4d-6204-465c-ab7e-64a12bbd6375}', true),
    ('5f843cf2-7d8f-45c3-9593-dbb344985fdc', 'e18b77c8-56f9-4d91-bbe0-92fa931bac48', 'FYR2040-C.BL-T.4XL', 'CONFIRMED',   '{684f9aef-5ea6-4864-8c88-dc4a0cbf30c0,887f0897-fdcb-4054-b162-12aa9b6476b5}', false),
    ('6b9643f5-2cb8-4c9e-a79c-c7824c316ce9', 'e18b77c8-56f9-4d91-bbe0-92fa931bac48', 'FYR2040-C.NG-T.L',   'PROVISIONAL', '{72378d06-a594-4373-a145-cf1bc2f3a96d,b553875b-787e-4e2b-9088-5bed91191446}', true),
    ('9675e14d-3e79-4c45-a9da-7c5a2d38bd0a', 'e18b77c8-56f9-4d91-bbe0-92fa931bac48', 'FYR2040-C.NG-T.XL',  'PROVISIONAL', '{2dd7dcb5-55bc-4f29-a8bc-222395ac615f,b553875b-787e-4e2b-9088-5bed91191446}', true),
    ('ca33621c-d5a1-4961-abee-8898c3169262', 'e18b77c8-56f9-4d91-bbe0-92fa931bac48', 'FYR2040-C.NG-T.2XL', 'PROVISIONAL', '{1158a4b7-0e08-4b36-9a97-373a4c15b0a0,b553875b-787e-4e2b-9088-5bed91191446}', true)
ON CONFLICT DO NOTHING;

-- VALPALB variants (Pantalón Laboral — Tipo Tejido × Color × Talla)
INSERT INTO product_variants (id, product_id, sku, status, attribute_values, is_active) VALUES
    ('edc06fe5-29b1-4ba6-bcae-f3e3eb8a5b2a', '91ff4c67-2810-49fb-96b0-9a8111ff8baa', 'VALPALB-C.RJ-T.M-TJ.PO', 'PROVISIONAL', '{2eae3a1d-de46-4932-b4bc-dc8987928917,39854811-247a-4947-a947-ba20841fba5c,98c36a68-2915-4567-8add-9bc4011d4e9a}', true)
ON CONFLICT DO NOTHING;

-- Servicios (no-attribute default variants — CONFIRMED)
INSERT INTO product_variants (id, product_id, sku, status, attribute_values, is_active) VALUES
    ('aea24acd-26b8-48a2-bc30-35c5f37baecc', '2f094614-4244-451f-8daf-0507df6eb281', 'SER', 'CONFIRMED', '{}', true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- SALES — PRESUPUESTOS
-- ============================================================================
INSERT INTO quotes (id, quote_number, party_id, quote_date, expiration_date, status,
                    subtotal_amount, tax_amount, total_amount, notes)
VALUES
    ('dfb4debf-eeaf-4149-b871-10241be572b3', 'PRE-2026-0001',
     'd4a5e6f7-0001-0001-0001-000000000001',
     '2026-03-21 19:21:38+00', '2026-04-25 00:00:00+00',
     'BORRADOR', 26.50, 5.57, 32.07, 'Presupuesto de camisetas urgente'),
    ('dc4fd259-9013-40e0-9048-cb9c80b524fb', 'PRE-2026-0002',
     'd4a5e6f7-0002-0002-0002-000000000002',
     '2026-03-21 20:18:20+00', '2026-05-10 00:00:00+00',
     'BORRADOR', 169.75, 35.64, 205.39, '')
ON CONFLICT DO NOTHING;

-- Quote line items
INSERT INTO quote_line_items (id, quote_id, product_variant_id, quantity,
    list_unit_price_amount, unit_price_amount, discount_per_unit_amount, discount_percent,
    tax_rate, tax_amount, subtotal_amount)
VALUES
    -- PRE-2026-0001: 10 × FYR2020-C.RJ-T.L @ 2.76 EUR (-4% over list)
    ('ca73e1d3-6498-487a-a612-6d6fa78758f2',
     'dfb4debf-eeaf-4149-b871-10241be572b3',
     '187ca801-71f4-416f-af22-063932620852',
     10, 2.76, 2.76, 0.00, 4.00, 21.00, 5.57, 26.50),
    -- PRE-2026-0002: 10 × FYR2040-C.NG-T.L @ 5.00 (-15%)
    ('afc68b6f-a810-4c05-8feb-d5d8d6b9e8d8',
     'dc4fd259-9013-40e0-9048-cb9c80b524fb',
     '9675e14d-3e79-4c45-a9da-7c5a2d38bd0a',
     10, 5.00, 5.00, 0.75, 15.00, 21.00, NULL, 42.50),
    -- PRE-2026-0002: 50 × SER @ 1.50 (-15%)
    ('475d811c-41b6-4656-b254-98098fd1b195',
     'dc4fd259-9013-40e0-9048-cb9c80b524fb',
     'aea24acd-26b8-48a2-bc30-35c5f37baecc',
     50, 1.50, 1.50, 0.23, 15.00, 21.00, NULL, 63.50),
    -- PRE-2026-0002: 5 × FYR2040-C.NG-T.2XL @ 5.00 (-15%)
    ('982849c1-36e6-4981-bbee-b7d3fea079fe',
     'dc4fd259-9013-40e0-9048-cb9c80b524fb',
     'ca33621c-d5a1-4961-abee-8898c3169262',
     5, 5.00, 5.00, 0.75, 15.00, 21.00, NULL, 21.25),
    -- PRE-2026-0002: 10 × FYR2040-C.NG-T.L @ 5.00 (-15%)  [same product, second line]
    ('3a28bb0f-5fd0-43ae-a88f-27f2640cdd32',
     'dc4fd259-9013-40e0-9048-cb9c80b524fb',
     '6b9643f5-2cb8-4c9e-a79c-c7824c316ce9',
     10, 5.00, 5.00, 0.75, 15.00, 21.00, NULL, 42.50)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- SALES — PEDIDOS
-- ============================================================================
INSERT INTO sales_orders (id, order_number, quote_id, party_id, order_date, delivery_date,
                          status, subtotal_amount, tax_amount, total_amount, notes)
VALUES
    ('c787e99d-a1ac-4377-b907-e59d74eb6b11', 'PED-2026-0002',
     NULL,
     'd4a5e6f7-0001-0001-0001-000000000001',
     '2026-03-21 00:00:00+00', '2026-03-26 00:00:00+00',
     'EN_PREPARACION', 26.50, 5.57, 32.07, 'Urgente')
ON CONFLICT DO NOTHING;

-- Order line items
INSERT INTO order_line_items (id, sales_order_id, product_variant_id, quantity,
    list_unit_price_amount, unit_price_amount, discount_per_unit_amount, discount_percent,
    tax_rate, tax_amount, subtotal_amount)
VALUES
    ('4b697d86-0760-48a9-8e19-a92db1b2023c',
     'c787e99d-a1ac-4377-b907-e59d74eb6b11',
     '187ca801-71f4-416f-af22-063932620852',
     10, 2.76, 2.76, 0.00, 4.00, 21.00, 5.57, 26.50)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- MES — WORK ORDERS
-- ============================================================================
INSERT INTO work_orders (id, work_number, work_name, party_id, work_setup_id, notes, status, priority)
VALUES
    ('92827368-9a10-4dd3-aacd-31fd14bb1e44',
     'MES-2026-001',
     'Bordadas en pecho logo empresa y serigrafía espalda nombre y logo',
     'd4a5e6f7-0001-0001-0001-000000000001',
     NULL,
     'Urgente',
     'PENDING',
     'NORMAL')
ON CONFLICT DO NOTHING;

-- Order ↔ Work Order link
INSERT INTO order_work_setups (id, order_id, work_setup_id, work_order_id, sequence, description)
VALUES
    ('083852ad-cf3f-4e06-acca-81e9a704733c',
     'c787e99d-a1ac-4377-b907-e59d74eb6b11',
     NULL,
     '92827368-9a10-4dd3-aacd-31fd14bb1e44',
     1,
     'Bordadas en pecho logo empresa y serigrafía espalda nombre y logo')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- MES — WORK TYPES
-- ============================================================================
INSERT INTO work_types (id, name, description, reference, is_active) VALUES
    ('d043dd11-d364-4dc5-9448-56c847ad4865', 'Bordado',             NULL,                                   'BOR',  true),
    ('30da7d74-c855-4c17-a066-f81d24f6d329', 'Serigrafía 1 Color',  'Serigrafía logo a 1 color',            'SER1', true),
    ('00e7758f-4c98-447d-9b26-9eda205851a1', 'Serigrafía 2 Colores','Serigrafía logos a dos colores',       'SER2', true),
    ('a129b90f-f0e0-4978-ab56-cca115e9eaae', 'Sublimación',         'Marcado mediante sublimación',         'SUB',  true),
    ('be125823-2699-499b-a262-eb1b4c179049', 'Vinilo Corte',        NULL,                                   'VIC',  true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- MES — TASKS
-- ============================================================================
INSERT INTO tasks (id, name, description, reference, is_active) VALUES
    ('28e769ec-6993-4c9f-a438-ac72d8220be8', 'Diseñar',       'Diseño y preparación de logos',                                                   'D',  true),
    ('f22bbe86-a9a4-4fc8-9274-d31d7f03bc20', 'Bordar',        'Bordado de logo',                                                                 'B',  true),
    ('3bbcaffa-4c20-4670-be1b-562bef35d6d1', 'Empaquetar',    'Preparar prenda para su entrega',                                                 'Q',  true),
    ('7777dbe4-065e-4307-93ad-2cb844236710', 'Imprimir',      NULL,                                                                              'I',  true),
    ('a92e378e-31a3-4258-a1a8-70d4e586a8ae', 'Marcar',        'Marcar directamente sobre la prenda',                                             'M',  true),
    ('3cdc00d7-d1d6-4530-82eb-7d93f8abc99d', 'Pelar',         'Para los vinilos hay que extraer el vinilo sobrante después del corte.',           'L',  true),
    ('040f8e1e-b289-4b47-945c-3672c12896e0', 'Planchar',      'Imprimar logo en prenda mediante planchado',                                      'H',  true),
    ('e2eadf54-ed33-41cb-b79e-99d7601821fc', 'Preparar Marco','Plasmar logos en la plancha de serigrafía',                                       'PM', true),
    ('8db6bfa1-b38f-43aa-8e8e-e6b7c495c9a0', 'Pulimerizar',   'Fijar serigrafía mediante planchado',                                             'Z',  true)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- MES — WORK TYPE ↔ TASK SEQUENCES
-- ============================================================================
INSERT INTO work_type_tasks (work_type_id, task_id, sequence) VALUES
    -- Bordado: Diseñar → Bordar → Empaquetar
    ('d043dd11-d364-4dc5-9448-56c847ad4865', '28e769ec-6993-4c9f-a438-ac72d8220be8', 1),
    ('d043dd11-d364-4dc5-9448-56c847ad4865', 'f22bbe86-a9a4-4fc8-9274-d31d7f03bc20', 2),
    ('d043dd11-d364-4dc5-9448-56c847ad4865', '3bbcaffa-4c20-4670-be1b-562bef35d6d1', 3),
    -- Serigrafía 1 Color: Diseñar → Imprimir → Preparar Marco → Marcar → Pulimerizar → Empaquetar
    ('30da7d74-c855-4c17-a066-f81d24f6d329', '28e769ec-6993-4c9f-a438-ac72d8220be8', 1),
    ('30da7d74-c855-4c17-a066-f81d24f6d329', '7777dbe4-065e-4307-93ad-2cb844236710', 2),
    ('30da7d74-c855-4c17-a066-f81d24f6d329', 'e2eadf54-ed33-41cb-b79e-99d7601821fc', 3),
    ('30da7d74-c855-4c17-a066-f81d24f6d329', 'a92e378e-31a3-4258-a1a8-70d4e586a8ae', 4),
    ('30da7d74-c855-4c17-a066-f81d24f6d329', '8db6bfa1-b38f-43aa-8e8e-e6b7c495c9a0', 5),
    ('30da7d74-c855-4c17-a066-f81d24f6d329', '3bbcaffa-4c20-4670-be1b-562bef35d6d1', 6),
    -- Serigrafía 2 Colores: Diseñar → Imprimir → Preparar Marco → Marcar → Pulimerizar → Empaquetar
    ('00e7758f-4c98-447d-9b26-9eda205851a1', '28e769ec-6993-4c9f-a438-ac72d8220be8', 1),
    ('00e7758f-4c98-447d-9b26-9eda205851a1', '7777dbe4-065e-4307-93ad-2cb844236710', 2),
    ('00e7758f-4c98-447d-9b26-9eda205851a1', 'e2eadf54-ed33-41cb-b79e-99d7601821fc', 3),
    ('00e7758f-4c98-447d-9b26-9eda205851a1', 'a92e378e-31a3-4258-a1a8-70d4e586a8ae', 4),
    ('00e7758f-4c98-447d-9b26-9eda205851a1', '8db6bfa1-b38f-43aa-8e8e-e6b7c495c9a0', 5),
    ('00e7758f-4c98-447d-9b26-9eda205851a1', '3bbcaffa-4c20-4670-be1b-562bef35d6d1', 6),
    -- Sublimación: Diseñar → Imprimir → Planchar → Empaquetar
    ('a129b90f-f0e0-4978-ab56-cca115e9eaae', '28e769ec-6993-4c9f-a438-ac72d8220be8', 1),
    ('a129b90f-f0e0-4978-ab56-cca115e9eaae', '7777dbe4-065e-4307-93ad-2cb844236710', 2),
    ('a129b90f-f0e0-4978-ab56-cca115e9eaae', '040f8e1e-b289-4b47-945c-3672c12896e0', 3),
    ('a129b90f-f0e0-4978-ab56-cca115e9eaae', '3bbcaffa-4c20-4670-be1b-562bef35d6d1', 4),
    -- Vinilo Corte: Diseñar → Imprimir → Pelar → Planchar → Empaquetar
    ('be125823-2699-499b-a262-eb1b4c179049', '28e769ec-6993-4c9f-a438-ac72d8220be8', 1),
    ('be125823-2699-499b-a262-eb1b4c179049', '7777dbe4-065e-4307-93ad-2cb844236710', 2),
    ('be125823-2699-499b-a262-eb1b4c179049', '3cdc00d7-d1d6-4530-82eb-7d93f8abc99d', 3),
    ('be125823-2699-499b-a262-eb1b4c179049', '040f8e1e-b289-4b47-945c-3672c12896e0', 4),
    ('be125823-2699-499b-a262-eb1b4c179049', '3bbcaffa-4c20-4670-be1b-562bef35d6d1', 5)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- DOCUMENT SEQUENCES
-- Counters as of 2026-03-21: 2 presupuestos, 2 pedidos issued (sequence at 2)
-- ============================================================================
INSERT INTO document_sequences (prefix, year, current_value) VALUES
    ('PRE', 2026, 2),
    ('PED', 2026, 2)
ON CONFLICT (prefix, year) DO UPDATE
    SET current_value = GREATEST(document_sequences.current_value, EXCLUDED.current_value);

COMMIT;

-- ============================================================================
-- END OF MIGRATION: 007_seed_data.sql
-- 5 product groups · 3 brands · 3 attributes · 26 attribute values
-- 4 parties · 6 products · 17 variants
-- 2 presupuestos · 1 pedido · 1 work order
-- 5 work types · 9 tasks · 24 work type task sequences
-- ============================================================================
