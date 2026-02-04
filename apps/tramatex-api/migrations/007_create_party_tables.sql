-- Migration: 007_create_party_tables
-- Description: Create Party tables aligned to ADR-012 (Party, Profiles, Roles, Relationships, ContactDetails)
-- Created: 2026-02-02

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

CREATE INDEX idx_parties_status ON parties(status);

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

CREATE INDEX idx_organization_profiles_tax_id ON organization_profiles(tax_id);

-- Party Roles (many per party)
CREATE TABLE IF NOT EXISTS party_roles (
    party_id VARCHAR(36) NOT NULL,
    role VARCHAR(30) NOT NULL,
    PRIMARY KEY (party_id, role),
    CONSTRAINT fk_party_roles_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE
);

CREATE INDEX idx_party_roles_role ON party_roles(role);

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

CREATE INDEX idx_party_relationships_from ON party_relationships(from_party_id);
CREATE INDEX idx_party_relationships_to ON party_relationships(to_party_id);
CREATE INDEX idx_party_relationships_type ON party_relationships(type);

-- Contact Details (for organization profile)
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

CREATE INDEX idx_contact_details_org_party ON contact_details(organization_party_id);

-- Addresses (kept for MVP, attached to party)
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

CREATE INDEX idx_party_addresses_party_id ON party_addresses(party_id);
CREATE INDEX idx_party_addresses_is_primary ON party_addresses(is_primary);

COMMIT;
