-- ============================================================================
-- Migration: 006_init_mes.sql
-- Module: MES (Tasks, Positions, Work Types, Work Setups, Work Orders)
-- Also creates quote_work_setups and order_work_setups (Sales join tables
-- that depend on MES work_setups and work_orders).
-- Date: 2026-04-14
-- ============================================================================


-- ============================================================================
-- TASKS (master data for manufacturing operations)
-- ============================================================================
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    reference VARCHAR(255) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_is_active ON tasks(is_active);

COMMENT ON TABLE tasks IS 'Master data for manufacturing tasks/operations';

-- ============================================================================
-- POSITIONS (work stations)
-- ============================================================================
CREATE TABLE IF NOT EXISTS positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_positions_code ON positions(code);
CREATE INDEX IF NOT EXISTS idx_positions_is_active ON positions(is_active);

COMMENT ON TABLE positions IS 'Work positions/stations in production floor';

-- ============================================================================
-- WORK TYPES (workflow templates)
-- ============================================================================
CREATE TABLE IF NOT EXISTS work_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    reference VARCHAR(255) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_work_types_is_active ON work_types(is_active);

COMMENT ON TABLE work_types IS 'Templates for service workflows (e.g., Embroidery Process)';

-- ============================================================================
-- WORK TYPE TASKS (ordered tasks within a work type)
-- ============================================================================
CREATE TABLE IF NOT EXISTS work_type_tasks (
    work_type_id UUID NOT NULL,
    task_id UUID NOT NULL,
    sequence INT NOT NULL,
    
    PRIMARY KEY (work_type_id, task_id),
    CONSTRAINT fk_work_type_tasks_work_type FOREIGN KEY (work_type_id) REFERENCES work_types(id) ON DELETE CASCADE,
    CONSTRAINT fk_work_type_tasks_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT chk_work_type_tasks_sequence_positive CHECK (sequence > 0)
);

CREATE INDEX IF NOT EXISTS idx_work_type_tasks_work_type_id ON work_type_tasks(work_type_id);
CREATE INDEX IF NOT EXISTS idx_work_type_tasks_sequence ON work_type_tasks(work_type_id, sequence);

COMMENT ON TABLE work_type_tasks IS 'Ordered tasks within a work type workflow';

-- ============================================================================
-- WORK SETUPS (predefined work configurations)
-- ============================================================================
CREATE TABLE IF NOT EXISTS work_setups (
    id              UUID PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    party_id        VARCHAR(255) NOT NULL,
    tangible_group_id UUID,
    description     TEXT NOT NULL DEFAULT '',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_work_setups_party ON work_setups(party_id);

COMMENT ON TABLE work_setups IS 'Predefined work configurations (party + product group + lines)';

-- ============================================================================
-- WORK SETUP LINES
-- ============================================================================
CREATE TABLE IF NOT EXISTS work_setup_lines (
    id              UUID PRIMARY KEY,
    work_setup_id   UUID NOT NULL REFERENCES work_setups(id) ON DELETE CASCADE,
    work_type_id    UUID NOT NULL REFERENCES work_types(id),
    position_id     UUID NOT NULL REFERENCES positions(id),
    design_file_path TEXT NOT NULL DEFAULT '',
    notes           TEXT NOT NULL DEFAULT '',
    sequence        INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_work_setup_lines_setup ON work_setup_lines(work_setup_id);

COMMENT ON TABLE work_setup_lines IS 'Lines in a work setup (work type + position pairs)';

-- ============================================================================
-- WORK ORDERS
-- ============================================================================
CREATE TABLE IF NOT EXISTS work_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_number VARCHAR(50) NOT NULL UNIQUE,
    work_name VARCHAR(200) NOT NULL,
    party_id VARCHAR(36) NOT NULL,
    work_setup_id UUID,
    notes TEXT DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    priority VARCHAR(20) NOT NULL DEFAULT 'NORMAL',
    start_date TIMESTAMP WITH TIME ZONE,
    due_date TIMESTAMP WITH TIME ZONE,
    completed_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_work_orders_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE RESTRICT,
    CONSTRAINT fk_work_orders_work_setup FOREIGN KEY (work_setup_id) REFERENCES work_setups(id) ON DELETE SET NULL,
    CONSTRAINT chk_work_orders_status CHECK (status IN ('PENDING', 'IN_PROGRESS', 'ON_HOLD', 'COMPLETED', 'CANCELLED', 'SUSPENDED')),
    CONSTRAINT chk_work_orders_priority CHECK (priority IN ('LOW', 'NORMAL', 'HIGH', 'URGENT'))
);

CREATE INDEX IF NOT EXISTS idx_work_orders_work_number ON work_orders(work_number);
CREATE INDEX IF NOT EXISTS idx_work_orders_party_id ON work_orders(party_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_status ON work_orders(status);
CREATE INDEX IF NOT EXISTS idx_work_orders_priority ON work_orders(priority);
CREATE INDEX IF NOT EXISTS idx_work_orders_due_date ON work_orders(due_date);
CREATE INDEX IF NOT EXISTS idx_work_orders_work_setup_id ON work_orders(work_setup_id);

COMMENT ON TABLE work_orders IS 'Manufacturing work orders (e.g., embroider 50 t-shirts)';

-- ============================================================================
-- WORK ORDER LINES
-- ============================================================================
CREATE TABLE IF NOT EXISTS work_order_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_order_id UUID NOT NULL,
    work_type_id UUID NOT NULL,
    position_id UUID NOT NULL,
    design_file_path VARCHAR(500),
    notes TEXT,
    sequence INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_work_order_lines_work_order FOREIGN KEY (work_order_id) REFERENCES work_orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_work_order_lines_work_type FOREIGN KEY (work_type_id) REFERENCES work_types(id) ON DELETE RESTRICT,
    CONSTRAINT fk_work_order_lines_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE RESTRICT,
    CONSTRAINT chk_work_order_lines_sequence_positive CHECK (sequence > 0)
);

CREATE INDEX IF NOT EXISTS idx_work_order_lines_work_order ON work_order_lines(work_order_id);
CREATE INDEX IF NOT EXISTS idx_work_order_lines_sequence ON work_order_lines(work_order_id, sequence);

COMMENT ON TABLE work_order_lines IS 'Service workflows applied to a specific work order';

-- ============================================================================
-- WORK ORDER TASKS
-- ============================================================================
CREATE TABLE IF NOT EXISTS work_order_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_order_line_id UUID NOT NULL,
    task_id UUID NOT NULL,
    sequence INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    assigned_to UUID,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_work_order_tasks_line FOREIGN KEY (work_order_line_id) REFERENCES work_order_lines(id) ON DELETE CASCADE,
    CONSTRAINT fk_work_order_tasks_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE RESTRICT,
    CONSTRAINT fk_work_order_tasks_assigned_to FOREIGN KEY (assigned_to) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_work_order_tasks_sequence_positive CHECK (sequence > 0),
    CONSTRAINT chk_work_order_tasks_status CHECK (status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'BLOCKED', 'SKIPPED'))
);

CREATE INDEX IF NOT EXISTS idx_work_order_tasks_line ON work_order_tasks(work_order_line_id);
CREATE INDEX IF NOT EXISTS idx_work_order_tasks_status ON work_order_tasks(status);
CREATE INDEX IF NOT EXISTS idx_work_order_tasks_assigned_to ON work_order_tasks(assigned_to);

COMMENT ON TABLE work_order_tasks IS 'Individual task instances with execution tracking';

-- ============================================================================
-- QUOTE WORK SETUPS (Sales join table — depends on work_setups)
-- ============================================================================
CREATE TABLE IF NOT EXISTS quote_work_setups (
    id            UUID PRIMARY KEY,
    quote_id      UUID NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    work_setup_id UUID REFERENCES work_setups(id),
    sequence      INT  NOT NULL DEFAULT 1,
    description   TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_quote_work_setups_quote ON quote_work_setups(quote_id);

COMMENT ON TABLE quote_work_setups IS 'MES work setup references linked to sales quotes';

-- ============================================================================
-- ORDER WORK SETUPS (Sales join table — depends on work_setups + work_orders)
-- ============================================================================
CREATE TABLE IF NOT EXISTS order_work_setups (
    id            UUID PRIMARY KEY,
    order_id      UUID NOT NULL REFERENCES sales_orders(id) ON DELETE CASCADE,
    work_setup_id UUID REFERENCES work_setups(id),
    work_order_id UUID REFERENCES work_orders(id),
    sequence      INT  NOT NULL DEFAULT 1,
    description   TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_order_work_setups_order ON order_work_setups(order_id);

COMMENT ON TABLE order_work_setups IS 'MES work setup references linked to sales orders';

-- ============================================================================
-- TRIGGERS
-- ============================================================================
DROP TRIGGER IF EXISTS trg_tasks_updated_at ON tasks;
DROP TRIGGER IF EXISTS trg_positions_updated_at ON positions;
DROP TRIGGER IF EXISTS trg_work_types_updated_at ON work_types;
DROP TRIGGER IF EXISTS trg_work_orders_updated_at ON work_orders;
DROP TRIGGER IF EXISTS trg_work_order_lines_updated_at ON work_order_lines;
DROP TRIGGER IF EXISTS trg_work_order_tasks_updated_at ON work_order_tasks;

CREATE TRIGGER trg_tasks_updated_at BEFORE UPDATE ON tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_positions_updated_at BEFORE UPDATE ON positions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_work_types_updated_at BEFORE UPDATE ON work_types FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_work_orders_updated_at BEFORE UPDATE ON work_orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_work_order_lines_updated_at BEFORE UPDATE ON work_order_lines FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_work_order_tasks_updated_at BEFORE UPDATE ON work_order_tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

