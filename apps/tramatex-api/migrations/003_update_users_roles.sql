-- Session 16: Update IAM roles to match project documentation
-- Date: 2026-02-01
-- Status: Ready for migration

-- Map legacy roles to new set
UPDATE users SET role = 'commercial' WHERE role = 'manager';
UPDATE users SET role = 'workshop' WHERE role = 'operator';

-- Update role constraint and default
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_role;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'commercial';
ALTER TABLE users ADD CONSTRAINT chk_role CHECK (role IN ('admin', 'commercial', 'designer', 'workshop'));

-- Comments update
COMMENT ON COLUMN users.role IS 'User role: admin, commercial, designer, workshop';
