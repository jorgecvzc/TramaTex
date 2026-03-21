-- 022: Add reference field to service_groups (WorkType)
ALTER TABLE service_groups ADD COLUMN reference VARCHAR(255) NOT NULL DEFAULT '';
