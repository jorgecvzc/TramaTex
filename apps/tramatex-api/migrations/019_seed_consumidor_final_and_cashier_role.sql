-- Migration: 019_seed_consumidor_final_and_cashier_role
-- Description: Seeds CONSUMIDOR_FINAL party for retail sales and adds cashier role to IAM
-- Created: 2026-02-14
-- Related: ADR-020 (Facturación simplificada)

BEGIN;

-- 1. Add cashier role to users table constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_role;
ALTER TABLE users ADD CONSTRAINT chk_role CHECK (role IN ('admin', 'commercial', 'designer', 'workshop', 'cashier'));

-- Update comments
COMMENT ON COLUMN users.role IS 'User role: admin, commercial, designer, workshop, cashier';

-- 2. Seed cashier role in roles table
INSERT INTO roles (name) VALUES ('cashier') ON CONFLICT (name) DO NOTHING;

-- 3. Seed CONSUMIDOR_FINAL party for retail sales without identified customer
-- This party is used for simplified invoices (tickets) per Spanish legislation (Real Decreto 1619/2012)
-- UUID: 00000000-0000-0000-0000-000000000001
-- NIF: 99999999R (reserved generic NIF for non-identified customers)

INSERT INTO parties (id, status, created_by, modified_by, created_at, modified_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'ACTIVE',
    'f47ac10b-58cc-4372-a567-0e02b2c3d479', -- admin user UUID from migration 006
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (id) DO NOTHING;

-- 4. Create organization profile for CONSUMIDOR_FINAL
INSERT INTO organization_profiles (party_id, name, tax_id, tax_id_type)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'CONSUMIDOR FINAL',
    '99999999R',
    'NIF'
)
ON CONFLICT (party_id) DO NOTHING;

-- 5. Assign CLIENT role to CONSUMIDOR_FINAL
INSERT INTO party_roles (party_id, role)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'CLIENT'
)
ON CONFLICT (party_id, role) DO NOTHING;

-- 6. Add comment for documentation
COMMENT ON TABLE parties IS 'Parties table - stores clients, suppliers, employees. UUID 00000000-0000-0000-0000-000000000001 is reserved for CONSUMIDOR_FINAL (generic retail customer per ADR-020)';

COMMIT;
