-- 027: Refactor mes_works to reference work_setups instead of tangible_groups
-- Work orders now derive their lines/tasks from a WorkSetup at creation time.
-- tangible_group_id and garment_notes are replaced by work_setup_id and notes.

BEGIN;

-- 1. Add new columns
ALTER TABLE mes_works ADD COLUMN IF NOT EXISTS work_setup_id UUID;
ALTER TABLE mes_works ADD COLUMN IF NOT EXISTS notes TEXT DEFAULT '';

-- 2. Migrate data: copy garment_notes to notes
UPDATE mes_works SET notes = COALESCE(garment_notes, '') WHERE notes = '' OR notes IS NULL;

-- 3. Drop old FK and columns
ALTER TABLE mes_works DROP CONSTRAINT IF EXISTS fk_mes_works_tangible_group;
ALTER TABLE mes_works DROP COLUMN IF EXISTS tangible_group_id;
ALTER TABLE mes_works DROP COLUMN IF EXISTS garment_notes;

-- 4. Add FK for work_setup_id (nullable for legacy rows without a setup)
ALTER TABLE mes_works ADD CONSTRAINT fk_mes_works_work_setup
    FOREIGN KEY (work_setup_id) REFERENCES work_setups(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_mes_works_work_setup_id ON mes_works(work_setup_id);

-- 5. Update status CHECK to remove DRAFT (already done in 025, but ensure consistency)
ALTER TABLE mes_works DROP CONSTRAINT IF EXISTS chk_mes_works_status;
ALTER TABLE mes_works ADD CONSTRAINT chk_mes_works_status
    CHECK (status IN ('PENDING', 'IN_PROGRESS', 'ON_HOLD', 'COMPLETED', 'CANCELLED'));

COMMIT;
