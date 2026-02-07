package domain

import "fmt"

// ProductionStatus represents the status of a ProductionOrder.
type ProductionStatus string

const (
	ProductionStatusPending    ProductionStatus = "PENDING"
	ProductionStatusInProgress ProductionStatus = "IN_PROGRESS"
	ProductionStatusCompleted  ProductionStatus = "COMPLETED"
	ProductionStatusOnHold     ProductionStatus = "ON_HOLD"
)

// IsValid checks if the ProductionStatus is one of the predefined valid statuses.
func (ps ProductionStatus) IsValid() bool {
	switch ps {
	case ProductionStatusPending, ProductionStatusInProgress, ProductionStatusCompleted, ProductionStatusOnHold:
		return true
	}
	return false
}

// TaskStatus represents the status of a TaskInstance.
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusBlocked    TaskStatus = "BLOCKED"
	TaskStatusSkipped    TaskStatus = "SKIPPED" // For ONE_TIME tasks already done
)

// IsValid checks if the TaskStatus is one of the predefined valid statuses.
func (ts TaskStatus) IsValid() bool {
	switch ts {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted, TaskStatusBlocked, TaskStatusSkipped:
		return true
	}
	return false
}


// QualityStatus represents the result of a QualityCheck.
type QualityStatus string

const (
	QualityStatusPassed  QualityStatus = "PASSED"
	QualityStatusFailed  QualityStatus = "FAILED"
	QualityStatusRework  QualityStatus = "REWORK"
	QualityStatusPending QualityStatus = "PENDING"
)

// IsValid checks if the QualityStatus is one of the predefined valid statuses.
func (qs QualityStatus) IsValid() bool {
	switch qs {
	case QualityStatusPassed, QualityStatusFailed, QualityStatusRework, QualityStatusPending:
		return true
	}
	return false
}

// RecipeType categorizes the type of product managed by the ProductionRecipe.
type RecipeType string

const (
	RecipeTypePhysicalProduct RecipeType = "PHYSICAL_PRODUCT"
	RecipeTypeServiceProduct  RecipeType = "SERVICE_PRODUCT"
)

// IsValid checks if the RecipeType is one of the predefined valid types.
func (rt RecipeType) IsValid() bool {
	switch rt {
	case RecipeTypePhysicalProduct, RecipeTypeServiceProduct:
		return true
	}
	return false
}

// TaskType differentiates between one-time and recurrent tasks.
type TaskType string

const (
	TaskTypeOneTime   TaskType = "ONE_TIME"
	TaskTypeRecurrent TaskType = "RECURRENT"
)

// IsValid checks if the TaskType is one of the predefined valid types.
func (tt TaskType) IsValid() bool {
	switch tt {
	case TaskTypeOneTime, TaskTypeRecurrent:
		return true
	}
	return false
}