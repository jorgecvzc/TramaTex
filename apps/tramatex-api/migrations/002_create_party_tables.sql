-- Migration: 002_create_party_tables
-- Description: Create tables for Party module (organizations, persons, contacts, addresses)
-- Created: 2026-01-18

BEGIN;

-- Organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('CLIENT', 'SUPPLIER', 'BOTH')),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    tax_id VARCHAR(50),
    tax_id_type VARCHAR(20),
    website VARCHAR(255),
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_by UUID NOT NULL,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_organizations_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_organizations_modified_by FOREIGN KEY (modified_by) REFERENCES users(id)
);

CREATE INDEX idx_organizations_role ON organizations(role);
CREATE INDEX idx_organizations_status ON organizations(status);
CREATE INDEX idx_organizations_tax_id ON organizations(tax_id);

-- Persons table (contacts within organizations)
CREATE TABLE IF NOT EXISTS persons (
    id VARCHAR(36) PRIMARY KEY,
    organization_id VARCHAR(36) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    job_title VARCHAR(100),
    is_primary_contact BOOLEAN DEFAULT FALSE,
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_by UUID NOT NULL,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_persons_organization FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_persons_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_persons_modified_by FOREIGN KEY (modified_by) REFERENCES users(id)
);

CREATE INDEX idx_persons_organization_id ON persons(organization_id);
CREATE INDEX idx_persons_email ON persons(email);
CREATE INDEX idx_persons_is_primary_contact ON persons(is_primary_contact);

-- Addresses table
CREATE TABLE IF NOT EXISTS addresses (
    id VARCHAR(36) PRIMARY KEY,
    organization_id VARCHAR(36) NOT NULL,
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
    CONSTRAINT fk_addresses_organization FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_addresses_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_addresses_modified_by FOREIGN KEY (modified_by) REFERENCES users(id)
);

CREATE INDEX idx_addresses_organization_id ON addresses(organization_id);
CREATE INDEX idx_addresses_is_primary ON addresses(is_primary);

COMMIT;
