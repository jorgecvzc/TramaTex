package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/mes/domain"
)

// WorkOrderQueryService provides read-only queries about WorkOrder execution state.
// Designed to be consumed by other modules (e.g., Sales) via adapter interfaces.
type WorkOrderQueryService struct {
	workOrderRepo domain.WorkOrderRepository
	taskRepo      domain.TaskRepository
}

// NewWorkOrderQueryService creates a new query service for work order progress.
func NewWorkOrderQueryService(
	workOrderRepo domain.WorkOrderRepository,
	taskRepo domain.TaskRepository,
) *WorkOrderQueryService {
	return &WorkOrderQueryService{
		workOrderRepo: workOrderRepo,
		taskRepo:      taskRepo,
	}
}

// GetWorkOrderProgress returns the full execution progress of a single WorkOrder,
// including per-line and per-task breakdown. All progress computation happens here
// in MES — consuming modules receive a pre-computed read model.
func (s *WorkOrderQueryService) GetWorkOrderProgress(ctx context.Context, workOrderID uuid.UUID) (*WorkOrderProgressDTO, error) {
	wo, err := s.workOrderRepo.FindByID(ctx, workOrderID)
	if err != nil {
		return nil, fmt.Errorf("find work order: %w", err)
	}
	if wo == nil {
		return nil, fmt.Errorf("work order not found: %s", workOrderID)
	}

	// Build task name lookup for enrichment
	taskNames, err := s.buildTaskNameMap(ctx, wo)
	if err != nil {
		return nil, err
	}

	return s.toProgressDTO(wo, taskNames), nil
}

// GetWorkOrdersProgress returns progress for multiple WorkOrders at once.
// Useful when Sales needs to check status of all work setups in an order.
func (s *WorkOrderQueryService) GetWorkOrdersProgress(ctx context.Context, workOrderIDs []uuid.UUID) ([]WorkOrderProgressDTO, error) {
	if len(workOrderIDs) == 0 {
		return []WorkOrderProgressDTO{}, nil
	}

	// Build a combined task name map across all orders
	allTaskNames := make(map[uuid.UUID]string)

	results := make([]WorkOrderProgressDTO, 0, len(workOrderIDs))
	for _, id := range workOrderIDs {
		wo, err := s.workOrderRepo.FindByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("find work order %s: %w", id, err)
		}
		if wo == nil {
			continue // skip missing orders silently
		}

		taskNames, err := s.buildTaskNameMap(ctx, wo)
		if err != nil {
			return nil, err
		}
		for k, v := range taskNames {
			allTaskNames[k] = v
		}

		results = append(results, *s.toProgressDTO(wo, taskNames))
	}

	return results, nil
}

// buildTaskNameMap loads task names for all tasks referenced in a WorkOrder.
func (s *WorkOrderQueryService) buildTaskNameMap(ctx context.Context, wo *domain.WorkOrder) (map[uuid.UUID]string, error) {
	taskIDs := make(map[uuid.UUID]struct{})
	for _, line := range wo.Lines {
		for _, task := range line.Tasks {
			taskIDs[task.TaskID] = struct{}{}
		}
	}

	names := make(map[uuid.UUID]string, len(taskIDs))
	for id := range taskIDs {
		task, err := s.taskRepo.FindByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("find task %s: %w", id, err)
		}
		if task != nil {
			names[id] = task.Name
		}
	}
	return names, nil
}

// toProgressDTO converts a domain WorkOrder into the cross-module progress DTO.
func (s *WorkOrderQueryService) toProgressDTO(wo *domain.WorkOrder, taskNames map[uuid.UUID]string) *WorkOrderProgressDTO {
	totalTasks := 0
	completedTasks := 0

	lines := make([]WorkOrderLineProgressDTO, 0, len(wo.Lines))
	for _, line := range wo.Lines {
		lineCompleted := 0
		tasks := make([]WorkOrderTaskProgressDTO, 0, len(line.Tasks))
		for _, task := range line.Tasks {
			if task.Status == domain.TaskStatusCompleted {
				lineCompleted++
			}
			tasks = append(tasks, WorkOrderTaskProgressDTO{
				TaskID:      task.TaskID,
				TaskName:    taskNames[task.TaskID],
				Sequence:    task.Sequence,
				Status:      string(task.Status),
				AssignedTo:  task.AssignedTo,
				StartedAt:   task.StartedAt,
				CompletedAt: task.CompletedAt,
			})
		}

		totalTasks += len(line.Tasks)
		completedTasks += lineCompleted

		lines = append(lines, WorkOrderLineProgressDTO{
			LineID:         line.ID,
			WorkTypeID:     line.WorkTypeID,
			PositionID:     line.PositionID,
			Sequence:       line.Sequence,
			Tasks:          tasks,
			TotalTasks:     len(line.Tasks),
			CompletedTasks: lineCompleted,
		})
	}

	return &WorkOrderProgressDTO{
		WorkOrderID:    wo.ID,
		OrderNumber:    wo.OrderNumber,
		OrderName:      wo.OrderName,
		Status:         string(wo.Status),
		Priority:       string(wo.Priority),
		StartDate:      wo.StartDate,
		DueDate:        wo.DueDate,
		CompletedDate:  wo.CompletedDate,
		Lines:          lines,
		TotalTasks:     totalTasks,
		CompletedTasks: completedTasks,
	}
}
