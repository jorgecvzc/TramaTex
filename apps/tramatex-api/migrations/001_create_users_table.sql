-- Session 14: Create users table with indexes and triggers
-- Date: 2026-01-13
-- Status: Ready for migration

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'operator',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT chk_role CHECK (role IN ('admin', 'manager', 'operator')),
    CONSTRAINT chk_email_not_empty CHECK (email != ''),
    CONSTRAINT chk_password_hash_not_empty CHECK (password_hash != '')
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Create function for automatic updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for updated_at
DROP TRIGGER IF EXISTS trg_users_update_timestamp ON users;
CREATE TRIGGER trg_users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE users IS 'Users table - stores user accounts with authentication data';
COMMENT ON COLUMN users.id IS 'Unique identifier (UUID)';
COMMENT ON COLUMN users.email IS 'User email address - unique, used for login';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password (cost=10)';
COMMENT ON COLUMN users.role IS 'User role: admin, manager, operator';
COMMENT ON COLUMN users.is_active IS 'Soft delete flag';
COMMENT ON COLUMN users.created_at IS 'Record creation timestamp (UTC)';
COMMENT ON COLUMN users.updated_at IS 'Record last update timestamp (UTC, auto-updated)';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp (NULL = active)';
