package domain

// legacy_aliases.go provides backward-compatible type aliases for the old MES entity names.
// These allow the existing application and infrastructure layers to compile
// while they are incrementally migrated to use the new names.
// TODO: Remove this file once the full application layer refactor is complete.

import (
	"time"

	"github.com/google/uuid"
)

// --- Entity Aliases ---

// ServiceGroup is the old name for WorkType.
type ServiceGroup = WorkType

// ServiceGroupTask is the old name for WorkTypeTask.
type ServiceGroupTask = WorkTypeTask

// MESWork is the old name for WorkOrder.
type MESWork = WorkOrder

// MESWorkServiceGroup is the old name for WorkOrderLine.
type MESWorkServiceGroup = WorkOrderLine

// MESWorkTask is the old name for WorkOrderTask.
type MESWorkTask = WorkOrderTask

// --- Filter Aliases ---

// ServiceGroupFilters is the old name for WorkTypeFilters.
type ServiceGroupFilters = WorkTypeFilters

// MESWorkFilters is the old name for WorkOrderFilters.
type MESWorkFilters = WorkOrderFilters

// --- Repository Aliases ---

// ServiceGroupRepository is the old name for WorkTypeRepository.
type ServiceGroupRepository = WorkTypeRepository

// MESWorkRepository is the old name for WorkOrderRepository.
type MESWorkRepository = WorkOrderRepository

// --- Constructor Aliases ---

// NewServiceGroup adapts the old constructor signature to NewWorkType.
// The old API included a productGroupID parameter that was removed in the refactor.
func NewServiceGroup(name, description string, _ *uuid.UUID, isActive bool, tasks []WorkTypeTask) (*WorkType, error) {
	return NewWorkType(name, description, isActive, tasks)
}

// NewMESWork adapts the old constructor signature to NewWorkOrder.
func NewMESWork(
	workNumber, workName, partyID string,
	tangibleGroupID uuid.UUID,
	garmentNotes string,
	status ProductionStatus,
	priority WorkPriority,
	startDate, dueDate, completedDate *time.Time,
	groups []WorkOrderLine,
) (*WorkOrder, error) {
	return NewWorkOrder(workNumber, workName, partyID, tangibleGroupID, nil, garmentNotes, status, priority, startDate, dueDate, completedDate, groups)
}
