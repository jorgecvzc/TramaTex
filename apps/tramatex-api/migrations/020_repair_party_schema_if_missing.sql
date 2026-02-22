-- Migration: 020_repair_party_schema_if_missing
-- Description: Repairs Party ADR schema when older environments marked migration 007/008 as executed but tables are missing.
-- Created: 2026-02-19

BEGIN;

-- Parties (Aggregate Root)
CREATE TABLE IF NOT EXISTS parties (
    id VARCHAR(36) PRIMARY KEY,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_by UUID NOT NULL,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_parties_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_parties_modified_by FOREIGN KEY (modified_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_parties_status ON parties(status);

-- Person Profile (optional, 1:1 with party)
CREATE TABLE IF NOT EXISTS person_profiles (
    party_id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    CONSTRAINT fk_person_profiles_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE
);

-- Organization Profile (optional, 1:1 with party)
CREATE TABLE IF NOT EXISTS organization_profiles (
    party_id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    tax_id VARCHAR(50),
    tax_id_type VARCHAR(20),
    website VARCHAR(255),
    CONSTRAINT fk_organization_profiles_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_organization_profiles_tax_id ON organization_profiles(tax_id);

-- Party Roles (many per party)
CREATE TABLE IF NOT EXISTS party_roles (
    party_id VARCHAR(36) NOT NULL,
    role VARCHAR(30) NOT NULL,
    creation_identifier VARCHAR(255) NULL,
    PRIMARY KEY (party_id, role),
    CONSTRAINT fk_party_roles_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE
);

ALTER TABLE party_roles
ADD COLUMN IF NOT EXISTS creation_identifier VARCHAR(255) NULL;

CREATE INDEX IF NOT EXISTS idx_party_roles_role ON party_roles(role);

-- Party Relationships
CREATE TABLE IF NOT EXISTS party_relationships (
    id VARCHAR(36) PRIMARY KEY,
    from_party_id VARCHAR(36) NOT NULL,
    to_party_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_party_relationships_from FOREIGN KEY (from_party_id) REFERENCES parties(id) ON DELETE CASCADE,
    CONSTRAINT fk_party_relationships_to FOREIGN KEY (to_party_id) REFERENCES parties(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_party_relationships_from ON party_relationships(from_party_id);
CREATE INDEX IF NOT EXISTS idx_party_relationships_to ON party_relationships(to_party_id);
CREATE INDEX IF NOT EXISTS idx_party_relationships_type ON party_relationships(type);

-- Contact Details
CREATE TABLE IF NOT EXISTS contact_details (
    id VARCHAR(36) PRIMARY KEY,
    organization_party_id VARCHAR(36) NOT NULL,
    type_description VARCHAR(100) NOT NULL,
    phone VARCHAR(30),
    email VARCHAR(255),
    related_party_id VARCHAR(36),
    CONSTRAINT fk_contact_details_org_party FOREIGN KEY (organization_party_id) REFERENCES parties(id) ON DELETE CASCADE,
    CONSTRAINT fk_contact_details_related_party FOREIGN KEY (related_party_id) REFERENCES parties(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_contact_details_org_party ON contact_details(organization_party_id);

-- Addresses (attached to party)
CREATE TABLE IF NOT EXISTS party_addresses (
    id VARCHAR(36) PRIMARY KEY,
    party_id VARCHAR(36) NOT NULL,
    street VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100),
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_by UUID NOT NULL,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_party_addresses_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE,
    CONSTRAINT fk_party_addresses_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_party_addresses_modified_by FOREIGN KEY (modified_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_party_addresses_party_id ON party_addresses(party_id);
CREATE INDEX IF NOT EXISTS idx_party_addresses_is_primary ON party_addresses(is_primary);

-- Backfill from legacy organizations/persons/addresses when they still exist.
DO $$
BEGIN
    IF to_regclass('public.organizations') IS NOT NULL THEN
        INSERT INTO parties (id, status, created_by, created_at, modified_by, modified_at)
        SELECT
            o.id,
            o.status,
            CASE
                WHEN EXISTS (SELECT 1 FROM users u WHERE u.id = o.created_by) THEN o.created_by
                ELSE 'f47ac10b-58cc-4372-a567-0e02b2c3d479'::uuid
            END,
            o.created_at,
            CASE
                WHEN EXISTS (SELECT 1 FROM users u WHERE u.id = o.modified_by) THEN o.modified_by
                ELSE 'f47ac10b-58cc-4372-a567-0e02b2c3d479'::uuid
            END,
            o.modified_at
        FROM organizations o
        ON CONFLICT (id) DO NOTHING;

        INSERT INTO organization_profiles (party_id, name, tax_id, tax_id_type, website)
        SELECT o.id, o.name, o.tax_id, o.tax_id_type, o.website
        FROM organizations o
        ON CONFLICT (party_id) DO UPDATE SET
            name = EXCLUDED.name,
            tax_id = EXCLUDED.tax_id,
            tax_id_type = EXCLUDED.tax_id_type,
            website = EXCLUDED.website;

        INSERT INTO party_roles (party_id, role)
        SELECT o.id, 'CLIENT'
        FROM organizations o
        WHERE o.role IN ('CLIENT', 'BOTH')
        ON CONFLICT DO NOTHING;

        INSERT INTO party_roles (party_id, role)
        SELECT o.id, 'SUPPLIER'
        FROM organizations o
        WHERE o.role IN ('SUPPLIER', 'BOTH')
        ON CONFLICT DO NOTHING;
    END IF;

    IF to_regclass('public.persons') IS NOT NULL THEN
        INSERT INTO parties (id, status, created_by, created_at, modified_by, modified_at)
        SELECT
            p.id,
            'ACTIVE',
            CASE
                WHEN EXISTS (SELECT 1 FROM users u WHERE u.id = p.created_by) THEN p.created_by
                ELSE 'f47ac10b-58cc-4372-a567-0e02b2c3d479'::uuid
            END,
            p.created_at,
            CASE
                WHEN EXISTS (SELECT 1 FROM users u WHERE u.id = p.modified_by) THEN p.modified_by
                ELSE 'f47ac10b-58cc-4372-a567-0e02b2c3d479'::uuid
            END,
            p.modified_at
        FROM persons p
        ON CONFLICT (id) DO NOTHING;

        INSERT INTO person_profiles (party_id, first_name, last_name)
        SELECT p.id, p.first_name, p.last_name
        FROM persons p
        ON CONFLICT (party_id) DO UPDATE SET
            first_name = EXCLUDED.first_name,
            last_name = EXCLUDED.last_name;

        INSERT INTO party_roles (party_id, role)
        SELECT p.id, 'EMPLOYEE'
        FROM persons p
        ON CONFLICT DO NOTHING;

        INSERT INTO party_relationships (id, from_party_id, to_party_id, type, created_at)
        SELECT 'rel-' || p.id, p.id, p.organization_id, 'IS_EMPLOYEE_OF', COALESCE(p.created_at, CURRENT_TIMESTAMP)
        FROM persons p
        ON CONFLICT (id) DO NOTHING;

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
    END IF;

    IF to_regclass('public.addresses') IS NOT NULL THEN
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
            CASE
                WHEN EXISTS (SELECT 1 FROM users u WHERE u.id = a.created_by) THEN a.created_by
                ELSE 'f47ac10b-58cc-4372-a567-0e02b2c3d479'::uuid
            END,
            a.created_at,
            CASE
                WHEN EXISTS (SELECT 1 FROM users u WHERE u.id = a.modified_by) THEN a.modified_by
                ELSE 'f47ac10b-58cc-4372-a567-0e02b2c3d479'::uuid
            END,
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
    END IF;
END $$;

COMMIT;
