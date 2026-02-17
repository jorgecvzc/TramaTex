package application

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// PRODUCTION RECIPE COMMANDS
// ============================================================================

// CreateProductionRecipeCommand creates a new production recipe
type CreateProductionRecipeCommand struct {
	Name            string                `json:"name" binding:"required"`
	ClientID        uuid.UUID             `json:"clientId" binding:"required"`
	ProductID       uuid.UUID             `json:"productId" binding:"required"`
	RecipeType      string                `json:"recipeType" binding:"required"` // STANDARD, CUSTOM
	TaskDefinitions []TaskDefinitionInput `json:"taskDefinitions"`
}

// UpdateProductionRecipeCommand updates an existing recipe
type UpdateProductionRecipeCommand struct {
	ID              uuid.UUID             `json:"id" binding:"required"`
	Name            string                `json:"name" binding:"required"`
	TaskDefinitions []TaskDefinitionInput `json:"taskDefinitions"`
}

// TaskDefinitionInput represents a task definition in a recipe
type TaskDefinitionInput struct {
	Name                string     `json:"name" binding:"required"`
	Description         string     `json:"description"`
	SequenceOrder       int        `json:"sequenceOrder" binding:"required"`
	EstimatedDurationHs float64    `json:"estimatedDurationHs" binding:"required"`
	WorkCenterID        *uuid.UUID `json:"workCenterId"`
}

// ============================================================================
// PRODUCTION ORDER COMMANDS
// ============================================================================

// CreateProductionOrderCommand creates a new production order
type CreateProductionOrderCommand struct {
	SalesOrderID uuid.UUID `json:"salesOrderId" binding:"required"`
	RecipeID     uuid.UUID `json:"recipeId" binding:"required"`
	ProductID    uuid.UUID `json:"productId" binding:"required"`
	Quantity     int       `json:"quantity" binding:"required,min=1"`
	StartDate    time.Time `json:"startDate" binding:"required"`
	EndDate      time.Time `json:"endDate" binding:"required"`
}

// UpdateProductionOrderStatusCommand changes order status
type UpdateProductionOrderStatusCommand struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	NewStatus string    `json:"newStatus" binding:"required"` // PENDING, IN_PROGRESS, COMPLETED, CANCELLED
}

// AssignWorkCenterCommand assigns a work center to a production order
type AssignWorkCenterCommand struct {
	ProductionOrderID uuid.UUID `json:"productionOrderId"  binding:"required"`
	WorkCenterID      uuid.UUID `json:"workCenterId" binding:"required"`
}

// ============================================================================
// TASK INSTANCE COMMANDS
// ============================================================================

// UpdateTaskStatusCommand updates a task instance status
type UpdateTaskStatusCommand struct {
	ProductionOrderID uuid.UUID `json:"productionOrderId" binding:"required"`
	TaskInstanceID    uuid.UUID `json:"taskInstanceId" binding:"required"`
	NewStatus         string    `json:"newStatus" binding:"required"` // PENDING, IN_PROGRESS, COMPLETED, BLOCKED
}

// AssignOperatorToTaskCommand assigns an operator to a task
type AssignOperatorToTaskCommand struct {
	ProductionOrderID uuid.UUID `json:"productionOrderId" binding:"required"`
	TaskInstanceID    uuid.UUID `json:"taskInstanceId" binding:"required"`
	OperatorID        uuid.UUID `json:"operatorId" binding:"required"`
}

// RecordTaskProgressCommand records actual time/completion for a task
type RecordTaskProgressCommand struct {
	ProductionOrderID uuid.UUID  `json:"productionOrderId" binding:"required"`
	TaskInstanceID    uuid.UUID  `json:"taskInstanceId" binding:"required"`
	ActualStartTime   *time.Time `json:"actualStartTime"`
	ActualEndTime     *time.Time `json:"actualEndTime"`
	Notes             string     `json:"notes"`
}

// ============================================================================
// WORK CENTER COMMANDS
// ============================================================================

// CreateWorkCenterCommand creates a new work center
type CreateWorkCenterCommand struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

// UpdateWorkCenterCommand updates an existing work center
type UpdateWorkCenterCommand struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
}
