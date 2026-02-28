-- ============================================================================
-- Migration: v2_001_init_iam.sql
-- Description: Initialize IAM (Identity and Access Management) module
-- Date: 2026-02-25
-- Modules: Users, Roles, Permissions, Authentication
-- ============================================================================

BEGIN;

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- USERS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'operator',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT chk_role CHECK (role IN ('admin', 'manager', 'operator', 'cashier')),
    CONSTRAINT chk_email_not_empty CHECK (email <> ''),
    CONSTRAINT chk_password_not_empty CHECK (password <> '')
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_is_active ON users(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

COMMENT ON TABLE users IS 'Users table - stores user accounts with authentication data';
COMMENT ON COLUMN users.id IS 'Unique identifier (UUID)';
COMMENT ON COLUMN users.email IS 'User email address - unique, used for login';
COMMENT ON COLUMN users.password IS 'Bcrypt hashed password (cost=10)';
COMMENT ON COLUMN users.role IS 'User role: admin, manager, operator, cashier';
COMMENT ON COLUMN users.is_active IS 'Account active status';
COMMENT ON COLUMN users.created_at IS 'Record creation timestamp (UTC)';
COMMENT ON COLUMN users.updated_at IS 'Record last update timestamp (UTC)';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp (NULL = active)';

-- ===========================================================================
-- ROLES & PERMISSIONS TABLES
-- ============================================================================
CREATE TABLE IF NOT EXISTS permissions (
    id VARCHAR(100) PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_permissions_category ON permissions(category);

COMMENT ON TABLE permissions IS 'System permissions catalog';

-- ============================================================================
CREATE TABLE IF NOT EXISTS role_permissions (
    role VARCHAR(50) NOT NULL,
    permission_id VARCHAR(100) NOT NULL,
    granted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (role, permission_id),
    CONSTRAINT fk_role_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE INDEX idx_role_permissions_role ON role_permissions(role);

COMMENT ON TABLE role_permissions IS 'Maps roles to permissions';

-- ============================================================================
-- REVOKED TOKENS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS revoked_tokens (
    token_id VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    
    CONSTRAINT fk_revoked_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_revoked_tokens_user_id ON revoked_tokens(user_id);
CREATE INDEX idx_revoked_tokens_expires_at ON revoked_tokens(expires_at);

COMMENT ON TABLE revoked_tokens IS 'Stores revoked JWT tokens for logout functionality';

-- ============================================================================
-- TRIGGERS
-- ============================================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_users_update_timestamp ON users;
CREATE TRIGGER trg_users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SEED DATA: Admin User
-- ============================================================================
-- Password: admin123
-- Hash generated with bcrypt cost=10
INSERT INTO users (id, email, password, role, is_active)
VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'admin@tramatex.local',
    '$2a$10$Y/OQADKTTu/8dOpA3BPc3eqOAQYREPJ03JWLuNKWpYApnFm5rl1oe',
    'admin',
    true
)
ON CONFLICT (email) DO UPDATE
SET
    password = EXCLUDED.password,
    role = 'admin',
    is_active = true,
    deleted_at = NULL;

COMMIT;

-- ============================================================================
-- END OF MIGRATION: v2_001_init_iam.sql
-- ============================================================================
