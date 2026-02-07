package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProductionPlanningService defines the interface for creating and orchestrating ProductionOrders.
type ProductionPlanningService interface {
	CreateProductionOrder(ctx context.Context, salesOrder SalesOrderRef, recipe *ProductionRecipe, quantity int, createdBy string) (*ProductionOrder, error)
	PlanProductionTasks(ctx context.Context, order *ProductionOrder) error
}

// TaskManagementService defines the interface for managing the lifecycle and dependencies of TaskInstances.
type TaskManagementService interface {
	AssignTask(ctx context.Context, task *TaskInstance, assigneeID uuid.UUID, workCenterID uuid.UUID, modifiedBy string) error
	StartTask(ctx context.Context, task *TaskInstance, modifiedBy string) error
	CompleteTask(ctx context.Context, task *TaskInstance, modifiedBy string) error
	FailTask(ctx context.Context, task *TaskInstance, reason string, modifiedBy string) error
	UpdateTaskParameters(ctx context.Context, task *TaskInstance, parameters map[string]string, modifiedBy string) error
	// Potentially other methods for managing task dependencies.
}

// QualityControlService defines the interface for performing and evaluating QualityChecks.
type QualityControlService interface {
	RecordQualityCheck(ctx context.Context, check *QualityCheck, createdBy string) error
	EvaluateQualityImpact(ctx context.Context, check *QualityCheck) error // e.g., blocks order, allows rework
}

// SalesOrderRef is a reference type to the Sales module's Order.
// This allows the MES domain to reference SalesOrder without knowing its full internal structure.
type SalesOrderRef struct {
	ID        uuid.UUID
	ProductID uuid.UUID // Product being ordered
	Quantity  int
	ClientID  uuid.UUID
}