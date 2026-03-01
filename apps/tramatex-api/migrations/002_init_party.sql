-- ============================================================================
-- Migration: v2_002_init_party.sql
-- Description: Initialize Party module
-- Date: 2026-02-25
-- Modules: Parties, Person/Organization Profiles, Roles, Relationships, Contact Details, Addresses
-- ============================================================================

BEGIN;

-- ============================================================================
-- PARTIES TABLE (Aggregate Root)
-- ============================================================================
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

COMMENT ON TABLE parties IS 'Parties aggregate root - represents any entity (person or organization)';
COMMENT ON COLUMN parties.status IS 'Party status: ACTIVE, INACTIVE';

-- ============================================================================
-- PERSON PROFILE (Optional, 1:1 with party)
-- ============================================================================
CREATE TABLE IF NOT EXISTS person_profiles (
    party_id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    
    CONSTRAINT fk_person_profiles_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE
);

COMMENT ON TABLE person_profiles IS 'Person-specific profile data (1:1 with party)';

-- ============================================================================
-- ORGANIZATION PROFILE (Optional, 1:1 with party)
-- ============================================================================
CREATE TABLE IF NOT EXISTS organization_profiles (
    party_id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    tax_id VARCHAR(50),
    tax_id_type VARCHAR(20),
    website VARCHAR(255),
    
    CONSTRAINT fk_organization_profiles_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_organization_profiles_tax_id ON organization_profiles(tax_id);

COMMENT ON TABLE organization_profiles IS 'Organization-specific profile data (1:1 with party)';

-- ============================================================================
-- PARTY ROLES (Many per party)
-- ============================================================================
CREATE TABLE IF NOT EXISTS party_roles (
    party_id VARCHAR(36) NOT NULL,
    role VARCHAR(30) NOT NULL,
    creation_identifier VARCHAR(255) NULL,
    
    PRIMARY KEY (party_id, role),
    CONSTRAINT fk_party_roles_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_party_roles_role ON party_roles(role);

COMMENT ON TABLE party_roles IS 'Party roles (CLIENT, SUPPLIER, EMPLOYEE, etc.)';
COMMENT ON COLUMN party_roles.creation_identifier IS 'Identifier from the source system that created this role';

-- ============================================================================
-- PARTY RELATIONSHIPS
-- ============================================================================
CREATE TABLE IF NOT EXISTS party_relationships (
    id VARCHAR(36) PRIMARY KEY,
    from_party_id VARCHAR(36) NOT NULL,
    to_party_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    modified_at TIMESTAMP,
    modified_by UUID,
    
    CONSTRAINT fk_party_relationships_from FOREIGN KEY (from_party_id) REFERENCES parties(id) ON DELETE CASCADE,
    CONSTRAINT fk_party_relationships_to FOREIGN KEY (to_party_id) REFERENCES parties(id) ON DELETE CASCADE,
    CONSTRAINT fk_party_relationships_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_party_relationships_modified_by FOREIGN KEY (modified_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_party_relationships_from ON party_relationships(from_party_id);
CREATE INDEX IF NOT EXISTS idx_party_relationships_to ON party_relationships(to_party_id);
CREATE INDEX IF NOT EXISTS idx_party_relationships_type ON party_relationships(type);

COMMENT ON TABLE party_relationships IS 'Relationships between parties (EMPLOYEE_OF, CONTACT_FOR, etc.)';

-- ============================================================================
-- CONTACT DETAILS (For organization contacts)
-- ============================================================================
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

COMMENT ON TABLE contact_details IS 'Contact information for organization parties';

-- ============================================================================
-- PARTY ADDRESSES
-- ============================================================================
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

COMMENT ON TABLE party_addresses IS 'Physical addresses for parties';

-- ============================================================================
-- PARTY SERVICE CONFIGURATIONS
-- ============================================================================
CREATE TABLE IF NOT EXISTS party_service_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    party_id UUID NOT NULL,
    service_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    configuration_details JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_party_service_configurations_party_id ON party_service_configurations(party_id);

DROP TRIGGER IF EXISTS trg_party_service_configurations_updated_at ON party_service_configurations;
CREATE TRIGGER trg_party_service_configurations_updated_at
BEFORE UPDATE ON party_service_configurations
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE party_service_configurations IS 'Service-specific configurations for parties';

COMMIT;

-- ============================================================================
-- END OF MIGRATION: v2_002_init_party.sql
-- ============================================================================
