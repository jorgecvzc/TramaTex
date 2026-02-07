package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProductionRecipe (Receta de Producción / Job Definition):
// Represents the reusable definition of a work type for a specific product/client.
type ProductionRecipe struct {
	ID                 uuid.UUID
	Name               string
	ClientID           uuid.UUID // Reference to the Party module's client
	ProductID          uuid.UUID // Reference to the Product module's product (or variant)
	TaskDefinitions    []TaskDefinition
	RecipeType         RecipeType
	Version            int // For tracking modifications
	IsMaster           bool
}

// NewProductionRecipe creates a new ProductionRecipe.
func NewProductionRecipe(id uuid.UUID, name string, clientID, productID uuid.UUID, recipeType RecipeType) (*ProductionRecipe, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("production recipe ID cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("production recipe name cannot be empty")
	}
	if clientID == uuid.Nil {
		return nil, fmt.Errorf("client ID cannot be empty")
	}
	if productID == uuid.Nil {
		return nil, fmt.Errorf("product ID cannot be empty")
	}
	if !recipeType.IsValid() {
		return nil, fmt.Errorf("invalid recipe type: %s", recipeType)
	}

	return &ProductionRecipe{
		ID:              id,
		Name:            name,
		ClientID:        clientID,
		ProductID:       productID,
		TaskDefinitions: make([]TaskDefinition, 0),
		RecipeType:      recipeType,
		Version:         1,
		IsMaster:        true,
	}, nil
}

// AddTaskDefinition adds a TaskDefinition to the ProductionRecipe.
func (pr *ProductionRecipe) AddTaskDefinition(taskDef TaskDefinition) {
	pr.TaskDefinitions = append(pr.TaskDefinitions, taskDef)
}

// ProductionOrder (Orden de Producción / Job Instance):
// A concrete instance of a ProductionRecipe for a specific quantity and sales order.
type ProductionOrder struct {
	ID                    uuid.UUID
	SalesOrderID          uuid.UUID // Reference to the Sales module's Order
	RecipeID              uuid.UUID // Reference to the ProductionRecipe that defines this job
	ProductID             uuid.UUID // Copy or reference of the Product original.
	Quantity              int
	Status                ProductionStatus
	TaskInstances         []TaskInstance // The instances of tasks to execute for this order.
	StartDate             time.Time
	EndDate               time.Time
	AssignedToWorkCenterID *uuid.UUID // Main work center assigned (optional)
}

// NewProductionOrder creates a new ProductionOrder.
func NewProductionOrder(id, salesOrderID, recipeID, productID uuid.UUID, quantity int, status ProductionStatus) (*ProductionOrder, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("production order ID cannot be empty")
	}
	if salesOrderID == uuid.Nil {
		return nil, fmt.Errorf("sales order ID cannot be empty")
	}
	if recipeID == uuid.Nil {
		return nil, fmt.Errorf("recipe ID cannot be empty")
	}
	if productID == uuid.Nil {
		return nil, fmt.Errorf("product ID cannot be empty")
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid production status: %s", status)
	}

	return &ProductionOrder{
		ID:            id,
		SalesOrderID:  salesOrderID,
		RecipeID:      recipeID,
		ProductID:     productID,
		Quantity:      quantity,
		Status:        status,
		TaskInstances: make([]TaskInstance, 0),
		StartDate:     time.Now(),
	}, nil
}

// TaskDefinition (Definición de Tarea):
// Defines a generic or specific task that is part of a ProductionRecipe.
type TaskDefinition struct {
	ID                   uuid.UUID
	Name                 string
	Description          string
	TaskType             TaskType // ONE_TIME or RECURRENT
	RequiredParameters   map[string]string
	ExpectedDurationMinutes int
	ResourceRequirements map[string]int // e.g., "machine_type_A": 1, "operator_role_B": 1
}

// NewTaskDefinition creates a new TaskDefinition.
func NewTaskDefinition(id uuid.UUID, name, description string, taskType TaskType, expectedDurationMinutes int) (*TaskDefinition, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("task definition ID cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("task definition name cannot be empty")
	}
	if !taskType.IsValid() {
		return nil, fmt.Errorf("invalid task type: %s", taskType)
	}
	if expectedDurationMinutes < 0 {
		return nil, fmt.Errorf("expected duration cannot be negative")
	}

	return &TaskDefinition{
		ID:                   id,
		Name:                 name,
		Description:          description,
		TaskType:             taskType,
		RequiredParameters:   make(map[string]string),
		ExpectedDurationMinutes: expectedDurationMinutes,
		ResourceRequirements: make(map[string]int),
	}, nil
}

// TaskInstance (Instancia de Tarea):
// Represents the concrete execution of a TaskDefinition for a specific ProductionOrder.
type TaskInstance struct {
	ID                  uuid.UUID
	TaskDefinitionID    uuid.UUID // Reference to the TaskDefinition.
	ProductionOrderID   uuid.UUID // Reference to the ProductionOrder parent.
	Status              TaskStatus
	AssignedToWorkCenterID *uuid.UUID // ID of the WorkCenter assigned
	AssignedToUserID    *uuid.UUID // ID of the User assigned (from IAM module)
	ActualParameters    map[string]string
	Notes               string
	ActualStartTime     *time.Time
	ActualEndTime       *time.Time
}

// NewTaskInstance creates a new TaskInstance from a TaskDefinition.
func NewTaskInstance(id uuid.UUID, taskDefinitionID, productionOrderID uuid.UUID, status TaskStatus, initialParameters map[string]string) (*TaskInstance, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("task instance ID cannot be empty")
	}
	if taskDefinitionID == uuid.Nil {
		return nil, fmt.Errorf("task definition ID cannot be empty")
	}
	if productionOrderID == uuid.Nil {
		return nil, fmt.Errorf("production order ID cannot be empty")
	}
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid task status: %s", status)
	}

	if initialParameters == nil {
		initialParameters = make(map[string]string)
	}

	return &TaskInstance{
		ID:                id,
		TaskDefinitionID:  taskDefinitionID,
		ProductionOrderID: productionOrderID,
		Status:            status,
		ActualParameters:  initialParameters,
	}, nil
}

// WorkCenter (Centro de Trabajo):
// Represents a station or resource where a part of the production is performed.
type WorkCenter struct {
	ID              uuid.UUID
	Name            string
	Description     string
	Capacity        int // e.g., number of items, hours, etc.
	AssignedTaskTypes []TaskType // What types of tasks it can perform.
}

// NewWorkCenter creates a new WorkCenter.
func NewWorkCenter(id uuid.UUID, name string, capacity int) (*WorkCenter, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("work center ID cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("work center name cannot be empty")
	}
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity must be greater than zero")
	}

	return &WorkCenter{
		ID:                id,
		Name:              name,
		Description:       description,
		Capacity:          capacity,
		AssignedTaskTypes: make([]TaskType, 0),
	}, nil
}

// QualityCheck (Control de Calidad):
// Record of a quality verification performed on a ProductionOrder or a TaskInstance.
type QualityCheck struct {
	ID                uuid.UUID
	ProductionOrderID uuid.UUID // Reference to the ProductionOrder
	TaskInstanceID    *uuid.UUID // If the check is specific to a task.
	InspectorUserID   uuid.UUID // Reference to the IAM module's User ID
	CheckDate         time.Time
	Result            QualityStatus
	Notes             string
	AttachedFileIDs   []uuid.UUID // Reference to files of evidence (NAS).
}

// NewQualityCheck creates a new QualityCheck.
func NewQualityCheck(id, productionOrderID, inspectorUserID uuid.UUID, checkDate time.Time, result QualityStatus) (*QualityCheck, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("quality check ID cannot be empty")
	}
	if productionOrderID == uuid.Nil {
		return nil, fmt.Errorf("production order ID cannot be empty")
	}
	if inspectorUserID == uuid.Nil {
		return nil, fmt.Errorf("inspector user ID cannot be empty")
	}
	if checkDate.IsZero() {
		return nil, fmt.Errorf("check date cannot be empty")
	}
	if !result.IsValid() {
		return nil, fmt.Errorf("invalid quality status: %s", result)
	}

	return &QualityCheck{
		ID:                id,
		ProductionOrderID: productionOrderID,
		InspectorUserID:   inspectorUserID,
		CheckDate:         checkDate,
		Result:            result,
		Notes:             "",
		AttachedFileIDs:   make([]uuid.UUID, 0),
	}, nil
}