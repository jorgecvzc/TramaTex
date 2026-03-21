-- 020: Add reference column to tasks
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS reference VARCHAR(255) NOT NULL DEFAULT '';
