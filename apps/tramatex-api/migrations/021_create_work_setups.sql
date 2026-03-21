-- 021: Create work_setups and work_setup_lines tables
CREATE TABLE IF NOT EXISTS work_setups (
    id              UUID PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    party_id        VARCHAR(255) NOT NULL,
    tangible_group_id UUID NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS work_setup_lines (
    id              UUID PRIMARY KEY,
    work_setup_id   UUID NOT NULL REFERENCES work_setups(id) ON DELETE CASCADE,
    work_type_id    UUID NOT NULL REFERENCES service_groups(id),
    position_id     UUID NOT NULL REFERENCES positions(id),
    design_file_path TEXT NOT NULL DEFAULT '',
    notes           TEXT NOT NULL DEFAULT '',
    sequence        INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_work_setups_party ON work_setups(party_id);
CREATE INDEX IF NOT EXISTS idx_work_setup_lines_setup ON work_setup_lines(work_setup_id);
