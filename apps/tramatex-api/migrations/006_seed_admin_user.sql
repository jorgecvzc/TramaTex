-- Session 19: Seed initial admin user if not exists
-- Date: 2026-02-01
-- Status: Ready for migration

INSERT INTO users (id, email, password, role, is_active)
VALUES (
  'f47ac10b-58cc-4372-a567-0e02b2c3d479',
  'admin@tramatex.local',
  '$2a$10$vDyHnw2RDOyyVDOWHX6sVuH.NSj/V4rUoeWhK3tYejgbrrKIwQIdS',
  'admin',
  true
)
ON CONFLICT (email) DO NOTHING;
