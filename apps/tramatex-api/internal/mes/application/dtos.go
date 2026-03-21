package application

import (
	"time"

	"github.com/google/uuid"
)

type TaskDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Reference   string    `json:"reference,omitempty"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
}

type PositionDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
}

type WorkTypeTaskDTO struct {
	TaskID   uuid.UUID `json:"task_id"`
	Sequence int       `json:"sequence"`
}

type WorkTypeDTO struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Reference   string            `json:"reference,omitempty"`
	Description string            `json:"description,omitempty"`
	IsActive    bool              `json:"is_active"`
	Tasks       []WorkTypeTaskDTO `json:"tasks"`
}

type WorkOrderTaskDTO struct {
	ID          uuid.UUID  `json:"id"`
	TaskID      uuid.UUID  `json:"task_id"`
	Sequence    int        `json:"sequence"`
	Status      string     `json:"status"`
	AssignedTo  *uuid.UUID `json:"assigned_to,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Notes       string     `json:"notes,omitempty"`
}

type WorkOrderLineDTO struct {
	ID             uuid.UUID          `json:"id"`
	WorkTypeID     uuid.UUID          `json:"work_type_id"`
	PositionID     uuid.UUID          `json:"position_id"`
	DesignFilePath string             `json:"design_file_path,omitempty"`
	Notes          string             `json:"notes,omitempty"`
	Sequence       int                `json:"sequence"`
	Tasks          []WorkOrderTaskDTO `json:"tasks"`
}

type WorkOrderDTO struct {
	ID               uuid.UUID          `json:"id"`
	WorkNumber       string             `json:"work_number"`
	WorkName         string             `json:"work_name"`
	PartyID          string             `json:"party_id"`
	WorkSetupID      *uuid.UUID         `json:"work_setup_id,omitempty"`
	Notes            string             `json:"notes,omitempty"`
	Status           string             `json:"status"`
	Priority         string             `json:"priority"`
	StartDate        *time.Time         `json:"start_date,omitempty"`
	DueDate          *time.Time         `json:"due_date,omitempty"`
	CompletedDate    *time.Time         `json:"completed_date,omitempty"`
	Lines            []WorkOrderLineDTO `json:"lines"`
	SalesOrderID     *uuid.UUID         `json:"sales_order_id,omitempty"`
	SalesOrderNumber string             `json:"sales_order_number,omitempty"`
}

type WorkOrderDashboardStatsDTO struct {
	Total    int            `json:"total"`
	ByStatus map[string]int `json:"by_status"`
	Overdue  int            `json:"overdue"`
	DueToday int            `json:"due_today"`
}

// --- WorkSetup DTOs ---

type WorkSetupLineDTO struct {
	ID             uuid.UUID `json:"id"`
	WorkTypeID     uuid.UUID `json:"work_type_id"`
	PositionID     uuid.UUID `json:"position_id"`
	DesignFilePath string    `json:"design_file_path,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	Sequence       int       `json:"sequence"`
}

type WorkSetupDTO struct {
	ID              uuid.UUID          `json:"id"`
	Name            string             `json:"name"`
	PartyID         string             `json:"party_id"`
	TangibleGroupID *uuid.UUID         `json:"tangible_group_id,omitempty"`
	Description     string             `json:"description,omitempty"`
	IsActive        bool               `json:"is_active"`
	Lines           []WorkSetupLineDTO `json:"lines"`
}

// --- Work Order Progress (consumed by Sales and other modules) ---

// WorkOrderTaskProgressDTO summarizes the execution state of one task within a line.
type WorkOrderTaskProgressDTO struct {
	TaskID      uuid.UUID  `json:"task_id"`
	TaskName    string     `json:"task_name,omitempty"`
	Sequence    int        `json:"sequence"`
	Status      string     `json:"status"`
	AssignedTo  *uuid.UUID `json:"assigned_to,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// WorkOrderLineProgressDTO summarizes progress of one line (WorkType @ Position).
type WorkOrderLineProgressDTO struct {
	LineID         uuid.UUID                  `json:"line_id"`
	WorkTypeID     uuid.UUID                  `json:"work_type_id"`
	PositionID     uuid.UUID                  `json:"position_id"`
	Sequence       int                        `json:"sequence"`
	Tasks          []WorkOrderTaskProgressDTO `json:"tasks"`
	TotalTasks     int                        `json:"total_tasks"`
	CompletedTasks int                        `json:"completed_tasks"`
}

// WorkOrderProgressDTO is the cross-module read model that other modules (Sales)
// can consume to understand the execution state of a WorkOrder without knowing
// MES internals. All progress logic is computed by MES.
type WorkOrderProgressDTO struct {
	WorkOrderID    uuid.UUID                  `json:"work_order_id"`
	OrderNumber    string                     `json:"order_number"`
	OrderName      string                     `json:"order_name"`
	Status         string                     `json:"status"`
	Priority       string                     `json:"priority"`
	StartDate      *time.Time                 `json:"start_date,omitempty"`
	DueDate        *time.Time                 `json:"due_date,omitempty"`
	CompletedDate  *time.Time                 `json:"completed_date,omitempty"`
	Lines          []WorkOrderLineProgressDTO `json:"lines"`
	TotalTasks     int                        `json:"total_tasks"`
	CompletedTasks int                        `json:"completed_tasks"`
}

// PendingWorkSetupDTO represents a Sales order work setup config that does not
// yet have a MES WorkOrder. Populated by the PendingSetupProvider adapter
// (Sales infrastructure) and returned through MES API.
type PendingWorkSetupDTO struct {
	ID           uuid.UUID  `json:"id"`
	WorkSetupID  *uuid.UUID `json:"workSetupId,omitempty"`
	Description  string     `json:"description"`
	OrderID      uuid.UUID  `json:"orderId"`
	OrderNumber  string     `json:"orderNumber"`
	DeliveryDate time.Time  `json:"deliveryDate"`
	PartyID      uuid.UUID  `json:"partyId"`
}
