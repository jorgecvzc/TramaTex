package domain

// ProductionStatus represents the status of a WorkOrder.
type ProductionStatus string

const (
	ProductionStatusDraft      ProductionStatus = "DRAFT"
	ProductionStatusPending    ProductionStatus = "PENDING"
	ProductionStatusInProgress ProductionStatus = "IN_PROGRESS"
	ProductionStatusCompleted  ProductionStatus = "COMPLETED"
	ProductionStatusOnHold     ProductionStatus = "ON_HOLD"
	ProductionStatusSuspended  ProductionStatus = "SUSPENDED"
	ProductionStatusCancelled  ProductionStatus = "CANCELLED"
)

func (ps ProductionStatus) IsValid() bool {
	switch ps {
	case ProductionStatusDraft, ProductionStatusPending, ProductionStatusInProgress, ProductionStatusCompleted, ProductionStatusOnHold, ProductionStatusSuspended, ProductionStatusCancelled:
		return true
	}
	return false
}

// WorkPriority represents the urgency of a WorkOrder.
type WorkPriority string

const (
	WorkPriorityLow    WorkPriority = "LOW"
	WorkPriorityNormal WorkPriority = "NORMAL"
	WorkPriorityHigh   WorkPriority = "HIGH"
	WorkPriorityUrgent WorkPriority = "URGENT"
)

func (wp WorkPriority) IsValid() bool {
	switch wp {
	case WorkPriorityLow, WorkPriorityNormal, WorkPriorityHigh, WorkPriorityUrgent:
		return true
	}
	return false
}

// TaskStatus represents the status of a WorkOrderTask.
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusBlocked    TaskStatus = "BLOCKED"
	TaskStatusSkipped    TaskStatus = "SKIPPED"
)

func (ts TaskStatus) IsValid() bool {
	switch ts {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted, TaskStatusBlocked, TaskStatusSkipped:
		return true
	}
	return false
}
