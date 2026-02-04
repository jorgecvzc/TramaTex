-- Migration: 008_migrate_party_data
-- Description: Migrate Party data (organizations, persons, addresses) to Party schema
-- Created: 2026-02-02

BEGIN;

-- 1) Migrate organizations to parties
INSERT INTO parties (id, status, created_by, created_at, modified_by, modified_at)
SELECT
    o.id,
    o.status,
    o.created_by,
    o.created_at,
    o.modified_by,
    o.modified_at
FROM organizations o
ON CONFLICT (id) DO NOTHING;

-- 2) Organization profiles
INSERT INTO organization_profiles (party_id, name, tax_id, tax_id_type, website)
SELECT
    o.id,
    o.name,
    o.tax_id,
    o.tax_id_type,
    o.website
FROM organizations o
ON CONFLICT (party_id) DO UPDATE SET
    name = EXCLUDED.name,
    tax_id = EXCLUDED.tax_id,
    tax_id_type = EXCLUDED.tax_id_type,
    website = EXCLUDED.website;

-- 3) Party roles from organization roles
-- CLIENT
INSERT INTO party_roles (party_id, role)
SELECT o.id, 'CLIENT'
FROM organizations o
WHERE o.role IN ('CLIENT', 'BOTH')
ON CONFLICT DO NOTHING;

-- SUPPLIER
INSERT INTO party_roles (party_id, role)
SELECT o.id, 'SUPPLIER'
FROM organizations o
WHERE o.role IN ('SUPPLIER', 'BOTH')
ON CONFLICT DO NOTHING;

-- 4) Migrate persons to parties with person profiles
INSERT INTO parties (id, status, created_by, created_at, modified_by, modified_at)
SELECT
    p.id,
    'ACTIVE',
    p.created_by,
    p.created_at,
    p.modified_by,
    p.modified_at
FROM persons p
ON CONFLICT (id) DO NOTHING;

INSERT INTO person_profiles (party_id, first_name, last_name)
SELECT
    p.id,
    p.first_name,
    p.last_name
FROM persons p
ON CONFLICT (party_id) DO UPDATE SET
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name;

-- Optional: give persons an EMPLOYEE role by default
INSERT INTO party_roles (party_id, role)
SELECT p.id, 'EMPLOYEE'
FROM persons p
ON CONFLICT DO NOTHING;

-- 5) Relationships: person -> organization
INSERT INTO party_relationships (id, from_party_id, to_party_id, type, created_at)
SELECT
    'rel-' || p.id,
    p.id,
    p.organization_id,
    'IS_EMPLOYEE_OF',
    COALESCE(p.created_at, CURRENT_TIMESTAMP)
FROM persons p
ON CONFLICT (id) DO NOTHING;

-- 6) Contact details for organization profile (from persons)
INSERT INTO contact_details (id, organization_party_id, type_description, phone, email, related_party_id)
SELECT
    'contact-' || p.id,
    p.organization_id,
    COALESCE(NULLIF(p.job_title, ''), 'Contacto'),
    p.phone,
    p.email,
    p.id
FROM persons p
ON CONFLICT (id) DO UPDATE SET
    organization_party_id = EXCLUDED.organization_party_id,
    type_description = EXCLUDED.type_description,
    phone = EXCLUDED.phone,
    email = EXCLUDED.email,
    related_party_id = EXCLUDED.related_party_id;

-- 7) Addresses -> party_addresses (attached to organization party)
INSERT INTO party_addresses (
    id, party_id, street, city, province, postal_code, country, is_primary,
    created_by, created_at, modified_by, modified_at
)
SELECT
    a.id,
    a.organization_id,
    a.street,
    a.city,
    a.province,
    a.postal_code,
    COALESCE(NULLIF(a.country, ''), 'Spain'),
    a.is_primary,
    a.created_by,
    a.created_at,
    a.modified_by,
    a.modified_at
FROM addresses a
ON CONFLICT (id) DO UPDATE SET
    party_id = EXCLUDED.party_id,
    street = EXCLUDED.street,
    city = EXCLUDED.city,
    province = EXCLUDED.province,
    postal_code = EXCLUDED.postal_code,
    country = EXCLUDED.country,
    is_primary = EXCLUDED.is_primary,
    modified_by = EXCLUDED.modified_by,
    modified_at = EXCLUDED.modified_at;

COMMIT;
