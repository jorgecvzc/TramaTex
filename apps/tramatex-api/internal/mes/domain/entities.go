package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// --- Datos Maestros ---

// Task represents an atomic, indivisible process within the garment customization flow.
// Examples: Diseñar, Imprimir, Marcar, Plegar, Embolsar.
type Task struct {
	ID          uuid.UUID
	Name        string
	Reference   string
	Description string
	IsActive    bool
}

func NewTask(name, reference, description string, isActive bool) (*Task, error) {
	if name == "" {
		return nil, NewValidationError("task name is required")
	}

	return &Task{
		ID:          uuid.New(),
		Name:        name,
		Reference:   reference,
		Description: description,
		IsActive:    isActive,
	}, nil
}

// Position represents a zone on the garment where customization is applied.
// Examples: Pecho izquierdo, Espalda, Bajo pantalón.
type Position struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Description string
	IsActive    bool
}

func NewPosition(name, code, description string, isActive bool) (*Position, error) {
	if name == "" {
		return nil, NewValidationError("position name is required")
	}
	if code == "" {
		return nil, NewValidationError("position code is required")
	}

	return &Position{
		ID:          uuid.New(),
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    isActive,
	}, nil
}

// WorkTypeTask is a value object linking a Task to a WorkType with execution order.
type WorkTypeTask struct {
	TaskID   uuid.UUID
	Sequence int
}

// WorkType defines an ordered sequence of tasks for a specific type of marking/customization.
// Examples: "Marcado por vinilo" = Diseñar → Imprimir → Marcar → Plegar → Embolsar.
type WorkType struct {
	ID          uuid.UUID
	Name        string
	Reference   string
	Description string
	IsActive    bool
	Tasks       []WorkTypeTask
}

func NewWorkType(name, reference, description string, isActive bool, tasks []WorkTypeTask) (*WorkType, error) {
	if name == "" {
		return nil, NewValidationError("work type name is required")
	}

	for _, task := range tasks {
		if task.TaskID == uuid.Nil {
			return nil, NewValidationError("task ID is required")
		}
		if task.Sequence <= 0 {
			return nil, NewValidationError("task sequence must be greater than zero")
		}
	}

	return &WorkType{
		ID:          uuid.New(),
		Name:        name,
		Reference:   reference,
		Description: description,
		IsActive:    isActive,
		Tasks:       tasks,
	}, nil
}

// --- Configuración por Cliente ---

// WorkSetupLine defines one customization operation: a WorkType applied at a Position.
type WorkSetupLine struct {
	ID             uuid.UUID
	WorkTypeID     uuid.UUID
	PositionID     uuid.UUID
	DesignFilePath string
	Notes          string
	Sequence       int
}

// WorkSetup is a reusable template defining the complete customization for a garment type
// for a specific customer. Example: Confecciones López / Camisetas → Serigrafía en Pecho + Vinilo en Espalda.
type WorkSetup struct {
	ID              uuid.UUID
	Name            string
	PartyID         string
	TangibleGroupID *uuid.UUID
	Description     string
	IsActive        bool
	Lines           []WorkSetupLine
}

func NewWorkSetup(
	name, partyID string,
	tangibleGroupID *uuid.UUID,
	description string,
	isActive bool,
	lines []WorkSetupLine,
) (*WorkSetup, error) {
	if name == "" {
		return nil, NewValidationError("work setup name is required")
	}
	if partyID == "" {
		return nil, NewValidationError("party id is required")
	}

	for _, line := range lines {
		if line.WorkTypeID == uuid.Nil {
			return nil, NewValidationError("work type id is required in setup line")
		}
		if line.PositionID == uuid.Nil {
			return nil, NewValidationError("position id is required in setup line")
		}
		if line.Sequence <= 0 {
			return nil, NewValidationError("setup line sequence must be greater than zero")
		}
	}

	return &WorkSetup{
		ID:              uuid.New(),
		Name:            name,
		PartyID:         partyID,
		TangibleGroupID: tangibleGroupID,
		Description:     description,
		IsActive:        isActive,
		Lines:           lines,
	}, nil
}

// --- Ejecución ---

// WorkOrderTask is an executable task instance with status tracking, operator assignment and timestamps.
type WorkOrderTask struct {
	ID          uuid.UUID
	TaskID      uuid.UUID
	Sequence    int
	Status      TaskStatus
	AssignedTo  *uuid.UUID
	StartedAt   *time.Time
	CompletedAt *time.Time
	Notes       string
}

// WorkOrderLine is an executable line (WorkType at a Position) within a WorkOrder.
type WorkOrderLine struct {
	ID             uuid.UUID
	WorkTypeID     uuid.UUID
	PositionID     uuid.UUID
	DesignFilePath string
	Notes          string
	Sequence       int
	Tasks          []WorkOrderTask
}

// WorkOrder is a real production order derived from a WorkSetup configuration.
// Its lines and tasks are auto-generated from the WorkSetup at creation time.
// WorkSetupID is optional — when created from Sales on order confirmation,
// the WorkOrder may not be backed by a WorkSetup.
type WorkOrder struct {
	ID            uuid.UUID
	OrderNumber   string
	OrderName     string
	PartyID       string
	WorkSetupID   *uuid.UUID
	Notes         string
	Status        ProductionStatus
	Priority      WorkPriority
	StartDate     *time.Time
	DueDate       *time.Time
	CompletedDate *time.Time
	Lines         []WorkOrderLine
}

func NewWorkOrder(
	orderNumber, orderName, partyID string,
	workSetupID *uuid.UUID,
	notes string,
	priority WorkPriority,
	dueDate *time.Time,
	lines []WorkOrderLine,
) (*WorkOrder, error) {
	if orderNumber == "" {
		return nil, fmt.Errorf("order number is required")
	}
	if orderName == "" {
		return nil, fmt.Errorf("order name is required")
	}
	if partyID == "" {
		return nil, fmt.Errorf("party id is required")
	}
	if !priority.IsValid() {
		return nil, fmt.Errorf("invalid work priority")
	}

	for _, line := range lines {
		if line.WorkTypeID == uuid.Nil {
			return nil, fmt.Errorf("work type id is required in order line")
		}
		if line.PositionID == uuid.Nil {
			return nil, fmt.Errorf("position id is required in order line")
		}
		if line.Sequence <= 0 {
			return nil, fmt.Errorf("order line sequence must be greater than zero")
		}
		for _, task := range line.Tasks {
			if task.TaskID == uuid.Nil {
				return nil, fmt.Errorf("task id is required in order task")
			}
			if task.Sequence <= 0 {
				return nil, fmt.Errorf("order task sequence must be greater than zero")
			}
			if !task.Status.IsValid() {
				return nil, fmt.Errorf("invalid order task status")
			}
		}
	}

	return &WorkOrder{
		ID:          uuid.New(),
		OrderNumber: orderNumber,
		OrderName:   orderName,
		PartyID:     partyID,
		WorkSetupID: workSetupID,
		Notes:       notes,
		Status:      ProductionStatusPending,
		Priority:    priority,
		DueDate:     dueDate,
		Lines:       lines,
	}, nil
}
