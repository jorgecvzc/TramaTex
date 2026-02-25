-- ============================================================================
-- Migration: v2_006_init_mes.sql
-- Description: Initialize MES (Manufacturing Execution System) module
-- Date: 2026-02-25
-- Modules: Tasks, Positions, Service Groups, MES Works
-- ============================================================================

BEGIN;

-- ============================================================================
-- MASTER DATA TABLES
-- ============================================================================

-- Tasks (master data for operations)
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_tasks_is_active ON tasks(is_active);

COMMENT ON TABLE tasks IS 'Master data for manufacturing tasks/operations';

-- ============================================================================
-- Positions (work stations)
CREATE TABLE IF NOT EXISTS positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_positions_code ON positions(code);
CREATE INDEX idx_positions_is_active ON positions(is_active);

COMMENT ON TABLE positions IS 'Work positions/stations in production floor';

-- ============================================================================
-- Service Groups (template for service workflows)
CREATE TABLE IF NOT EXISTS service_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    product_group_id UUID,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_service_groups_product_group FOREIGN KEY (product_group_id) REFERENCES product_groups(id) ON DELETE SET NULL
);

CREATE INDEX idx_service_groups_is_active ON service_groups(is_active);
CREATE INDEX idx_service_groups_product_group_id ON service_groups(product_group_id);

COMMENT ON TABLE service_groups IS 'Templates for service workflows (e.g., "Embroidery Process")';

-- ============================================================================
-- Service Group Tasks (many-to-many with sequence)
CREATE TABLE IF NOT EXISTS service_group_tasks (
    service_group_id UUID NOT NULL,
    task_id UUID NOT NULL,
    sequence INT NOT NULL,
    
    PRIMARY KEY (service_group_id, task_id),
    CONSTRAINT fk_service_group_tasks_service_group FOREIGN KEY (service_group_id) REFERENCES service_groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_service_group_tasks_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CONSTRAINT chk_service_group_tasks_sequence_positive CHECK (sequence > 0)
);

CREATE INDEX idx_service_group_tasks_service_group_id ON service_group_tasks(service_group_id);
CREATE INDEX idx_service_group_tasks_sequence ON service_group_tasks(service_group_id, sequence);

COMMENT ON TABLE service_group_tasks IS 'Ordered tasks within a service group workflow';

-- ============================================================================
-- MES WORK ORDERS
-- ============================================================================

-- MES Works (work orders)
CREATE TABLE IF NOT EXISTS mes_works (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_number VARCHAR(50) NOT NULL UNIQUE,
    work_name VARCHAR(200) NOT NULL,
    party_id VARCHAR(36) NOT NULL,
    tangible_group_id UUID NOT NULL,
    garment_notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    priority VARCHAR(20) NOT NULL DEFAULT 'NORMAL',
    start_date TIMESTAMP WITH TIME ZONE,
    due_date TIMESTAMP WITH TIME ZONE,
    completed_date TIMESTAMP WITH TIME ZONE,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_mes_works_party FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE RESTRICT,
    CONSTRAINT fk_mes_works_tangible_group FOREIGN KEY (tangible_group_id) REFERENCES product_groups(id) ON DELETE RESTRICT,
    CONSTRAINT chk_mes_works_status CHECK (status IN ('DRAFT', 'PENDING', 'IN_PROGRESS', 'ON_HOLD', 'COMPLETED', 'CANCELLED')),
    CONSTRAINT chk_mes_works_priority CHECK (priority IN ('LOW', 'NORMAL', 'HIGH', 'URGENT'))
);

CREATE INDEX idx_mes_works_work_number ON mes_works(work_number);
CREATE INDEX idx_mes_works_party_id ON mes_works(party_id);
CREATE INDEX idx_mes_works_status ON mes_works(status);
CREATE INDEX idx_mes_works_priority ON mes_works(priority);
CREATE INDEX idx_mes_works_due_date ON mes_works(due_date);

COMMENT ON TABLE mes_works IS 'Manufacturing work orders (e.g., embroider 50 t-shirts)';

-- ============================================================================
-- MES Work Service Groups (service instances per work order)
CREATE TABLE IF NOT EXISTS mes_work_service_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mes_work_id UUID NOT NULL,
    service_group_id UUID NOT NULL,
    position_id UUID NOT NULL,
    design_file_path VARCHAR(500),
    notes TEXT,
    sequence INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_mes_work_service_groups_work FOREIGN KEY (mes_work_id) REFERENCES mes_works(id) ON DELETE CASCADE,
    CONSTRAINT fk_mes_work_service_groups_service_group FOREIGN KEY (service_group_id) REFERENCES service_groups(id) ON DELETE RESTRICT,
    CONSTRAINT fk_mes_work_service_groups_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE RESTRICT,
    CONSTRAINT chk_mes_work_service_groups_sequence_positive CHECK (sequence > 0)
);

CREATE INDEX idx_mes_work_service_groups_work ON mes_work_service_groups(mes_work_id);
CREATE INDEX idx_mes_work_service_groups_sequence ON mes_work_service_groups(mes_work_id, sequence);

COMMENT ON TABLE mes_work_service_groups IS 'Service workflows applied to a specific work order';

-- ============================================================================
-- MES Work Tasks (task instances with tracking)
CREATE TABLE IF NOT EXISTS mes_work_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mes_work_service_group_id UUID NOT NULL,
    task_id UUID NOT NULL,
    sequence INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    assigned_to UUID,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_mes_work_tasks_group FOREIGN KEY (mes_work_service_group_id) REFERENCES mes_work_service_groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_mes_work_tasks_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE RESTRICT,
    CONSTRAINT fk_mes_work_tasks_assigned_to FOREIGN KEY (assigned_to) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_mes_work_tasks_sequence_positive CHECK (sequence > 0),
    CONSTRAINT chk_mes_work_tasks_status CHECK (status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'BLOCKED', 'SKIPPED'))
);

CREATE INDEX idx_mes_work_tasks_group ON mes_work_tasks(mes_work_service_group_id);
CREATE INDEX idx_mes_work_tasks_status ON mes_work_tasks(status);
CREATE INDEX idx_mes_work_tasks_assigned_to ON mes_work_tasks(assigned_to);

COMMENT ON TABLE mes_work_tasks IS 'Individual task instances with execution tracking';

-- ============================================================================
-- FOREIGN KEYS FROM SALES MODULE
-- ============================================================================
-- Add foreign key constraints from Sales module to MES Works
-- (The columns were already defined in v2_005_init_sales.sql)

DO $$ BEGIN
    ALTER TABLE quote_line_items
        ADD CONSTRAINT fk_quote_line_items_mes_work
        FOREIGN KEY (mes_work_id) REFERENCES mes_works(id) ON DELETE SET NULL;
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    ALTER TABLE order_line_items
        ADD CONSTRAINT fk_order_line_items_mes_work
        FOREIGN KEY (mes_work_id) REFERENCES mes_works(id) ON DELETE SET NULL;
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE INDEX IF NOT EXISTS idx_quote_line_items_mes_work_id ON quote_line_items(mes_work_id);
CREATE INDEX IF NOT EXISTS idx_order_line_items_mes_work_id ON order_line_items(mes_work_id);

-- ============================================================================
-- TRIGGERS
-- ============================================================================
DROP TRIGGER IF EXISTS trg_tasks_updated_at ON tasks;
DROP TRIGGER IF EXISTS trg_positions_updated_at ON positions;
DROP TRIGGER IF EXISTS trg_service_groups_updated_at ON service_groups;
DROP TRIGGER IF EXISTS trg_mes_works_updated_at ON mes_works;
DROP TRIGGER IF EXISTS trg_mes_work_service_groups_updated_at ON mes_work_service_groups;
DROP TRIGGER IF EXISTS trg_mes_work_tasks_updated_at ON mes_work_tasks;

CREATE TRIGGER trg_tasks_updated_at BEFORE UPDATE ON tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_positions_updated_at BEFORE UPDATE ON positions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_service_groups_updated_at BEFORE UPDATE ON service_groups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_mes_works_updated_at BEFORE UPDATE ON mes_works FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_mes_work_service_groups_updated_at BEFORE UPDATE ON mes_work_service_groups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_mes_work_tasks_updated_at BEFORE UPDATE ON mes_work_tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;

-- ============================================================================
-- END OF MIGRATION: v2_006_init_mes.sql
-- ============================================================================
