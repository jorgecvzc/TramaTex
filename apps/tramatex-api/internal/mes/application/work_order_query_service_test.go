package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/mes/domain"
)

// --- In-memory mock repositories for testing ---

type mockTaskRepo struct {
	tasks map[uuid.UUID]*domain.Task
}

func newMockTaskRepo(tasks ...*domain.Task) *mockTaskRepo {
	m := &mockTaskRepo{tasks: make(map[uuid.UUID]*domain.Task)}
	for _, t := range tasks {
		m.tasks[t.ID] = t
	}
	return m
}

func (r *mockTaskRepo) Save(_ context.Context, task *domain.Task) error { return nil }
func (r *mockTaskRepo) FindByID(_ context.Context, id uuid.UUID) (*domain.Task, error) {
	return r.tasks[id], nil
}
func (r *mockTaskRepo) FindAll(_ context.Context, _ *domain.TaskFilters) ([]*domain.Task, error) {
	return nil, nil
}
func (r *mockTaskRepo) Delete(_ context.Context, _ uuid.UUID) error         { return nil }
func (r *mockTaskRepo) Exists(_ context.Context, _ uuid.UUID) (bool, error) { return true, nil }

type mockWorkOrderRepo struct {
	orders map[uuid.UUID]*domain.WorkOrder
}

func newMockWorkOrderRepo(orders ...*domain.WorkOrder) *mockWorkOrderRepo {
	m := &mockWorkOrderRepo{orders: make(map[uuid.UUID]*domain.WorkOrder)}
	for _, o := range orders {
		m.orders[o.ID] = o
	}
	return m
}

func (r *mockWorkOrderRepo) Save(_ context.Context, _ *domain.WorkOrder) error { return nil }
func (r *mockWorkOrderRepo) FindByID(_ context.Context, id uuid.UUID) (*domain.WorkOrder, error) {
	return r.orders[id], nil
}
func (r *mockWorkOrderRepo) FindAll(_ context.Context, _ *domain.WorkOrderFilters) ([]*domain.WorkOrder, error) {
	return nil, nil
}
func (r *mockWorkOrderRepo) CountByYear(_ context.Context, _ int) (int64, error) { return 0, nil }

// --- Tests ---

func TestGetWorkOrderProgress_FullyCompleted(t *testing.T) {
	taskDesign := &domain.Task{ID: uuid.New(), Name: "Diseñar", IsActive: true}
	taskPrint := &domain.Task{ID: uuid.New(), Name: "Imprimir", IsActive: true}

	now := time.Now()
	wo := &domain.WorkOrder{
		ID:              uuid.New(),
		OrderNumber:     "WO-2026-001",
		OrderName:       "Test Order",
		PartyID:         "party-1",
		WorkSetupID: uuidPtr(uuid.New()),
		Status:          domain.ProductionStatusCompleted,
		Priority:        domain.WorkPriorityNormal,
		CompletedDate:   &now,
		Lines: []domain.WorkOrderLine{
			{
				ID:         uuid.New(),
				WorkTypeID: uuid.New(),
				PositionID: uuid.New(),
				Sequence:   1,
				Tasks: []domain.WorkOrderTask{
					{ID: uuid.New(), TaskID: taskDesign.ID, Sequence: 1, Status: domain.TaskStatusCompleted},
					{ID: uuid.New(), TaskID: taskPrint.ID, Sequence: 2, Status: domain.TaskStatusCompleted},
				},
			},
		},
	}

	svc := NewWorkOrderQueryService(
		newMockWorkOrderRepo(wo),
		newMockTaskRepo(taskDesign, taskPrint),
	)

	progress, err := svc.GetWorkOrderProgress(context.Background(), wo.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.WorkOrderID != wo.ID {
		t.Errorf("expected work order ID %s, got %s", wo.ID, progress.WorkOrderID)
	}
	if progress.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", progress.Status)
	}
	if progress.TotalTasks != 2 {
		t.Errorf("expected 2 total tasks, got %d", progress.TotalTasks)
	}
	if progress.CompletedTasks != 2 {
		t.Errorf("expected 2 completed tasks, got %d", progress.CompletedTasks)
	}
	if len(progress.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(progress.Lines))
	}
	if progress.Lines[0].TotalTasks != 2 {
		t.Errorf("expected line with 2 tasks, got %d", progress.Lines[0].TotalTasks)
	}
	if progress.Lines[0].CompletedTasks != 2 {
		t.Errorf("expected line with 2 completed, got %d", progress.Lines[0].CompletedTasks)
	}
}

func TestGetWorkOrderProgress_PartiallyCompleted(t *testing.T) {
	taskDesign := &domain.Task{ID: uuid.New(), Name: "Diseñar", IsActive: true}
	taskPrint := &domain.Task{ID: uuid.New(), Name: "Imprimir", IsActive: true}
	taskMark := &domain.Task{ID: uuid.New(), Name: "Marcar", IsActive: true}

	wo := &domain.WorkOrder{
		ID:              uuid.New(),
		OrderNumber:     "WO-2026-002",
		OrderName:       "Partial Order",
		PartyID:         "party-1",
		WorkSetupID: uuidPtr(uuid.New()),
		Status:          domain.ProductionStatusInProgress,
		Priority:        domain.WorkPriorityHigh,
		Lines: []domain.WorkOrderLine{
			{
				ID: uuid.New(), WorkTypeID: uuid.New(), PositionID: uuid.New(), Sequence: 1,
				Tasks: []domain.WorkOrderTask{
					{ID: uuid.New(), TaskID: taskDesign.ID, Sequence: 1, Status: domain.TaskStatusCompleted},
					{ID: uuid.New(), TaskID: taskPrint.ID, Sequence: 2, Status: domain.TaskStatusInProgress},
				},
			},
			{
				ID: uuid.New(), WorkTypeID: uuid.New(), PositionID: uuid.New(), Sequence: 2,
				Tasks: []domain.WorkOrderTask{
					{ID: uuid.New(), TaskID: taskDesign.ID, Sequence: 1, Status: domain.TaskStatusCompleted},
					{ID: uuid.New(), TaskID: taskMark.ID, Sequence: 2, Status: domain.TaskStatusPending},
				},
			},
		},
	}

	svc := NewWorkOrderQueryService(
		newMockWorkOrderRepo(wo),
		newMockTaskRepo(taskDesign, taskPrint, taskMark),
	)

	progress, err := svc.GetWorkOrderProgress(context.Background(), wo.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.TotalTasks != 4 {
		t.Errorf("expected 4 total tasks, got %d", progress.TotalTasks)
	}
	if progress.CompletedTasks != 2 {
		t.Errorf("expected 2 completed tasks, got %d", progress.CompletedTasks)
	}

	// Line 1: 1/2 completed
	if progress.Lines[0].CompletedTasks != 1 {
		t.Errorf("line 1: expected 1 completed, got %d", progress.Lines[0].CompletedTasks)
	}
	// Line 2: 1/2 completed
	if progress.Lines[1].CompletedTasks != 1 {
		t.Errorf("line 2: expected 1 completed, got %d", progress.Lines[1].CompletedTasks)
	}
}

func TestGetWorkOrderProgress_TaskNamesEnriched(t *testing.T) {
	taskDesign := &domain.Task{ID: uuid.New(), Name: "Diseñar", IsActive: true}

	wo := &domain.WorkOrder{
		ID:              uuid.New(),
		OrderNumber:     "WO-2026-003",
		OrderName:       "Names Test",
		PartyID:         "party-1",
		WorkSetupID: uuidPtr(uuid.New()),
		Status:          domain.ProductionStatusPending,
		Priority:        domain.WorkPriorityNormal,
		Lines: []domain.WorkOrderLine{
			{
				ID: uuid.New(), WorkTypeID: uuid.New(), PositionID: uuid.New(), Sequence: 1,
				Tasks: []domain.WorkOrderTask{
					{ID: uuid.New(), TaskID: taskDesign.ID, Sequence: 1, Status: domain.TaskStatusPending},
				},
			},
		},
	}

	svc := NewWorkOrderQueryService(
		newMockWorkOrderRepo(wo),
		newMockTaskRepo(taskDesign),
	)

	progress, err := svc.GetWorkOrderProgress(context.Background(), wo.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.Lines[0].Tasks[0].TaskName != "Diseñar" {
		t.Errorf("expected task name 'Diseñar', got '%s'", progress.Lines[0].Tasks[0].TaskName)
	}
}

func TestGetWorkOrderProgress_NotFound(t *testing.T) {
	svc := NewWorkOrderQueryService(
		newMockWorkOrderRepo(),
		newMockTaskRepo(),
	)

	_, err := svc.GetWorkOrderProgress(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for missing work order")
	}
}

func TestGetWorkOrdersProgress_MultipleOrders(t *testing.T) {
	taskID := uuid.New()
	task := &domain.Task{ID: taskID, Name: "Plegar", IsActive: true}

	wo1 := &domain.WorkOrder{
		ID: uuid.New(), OrderNumber: "WO-001", OrderName: "Order 1",
		PartyID: "p1", WorkSetupID: uuidPtr(uuid.New()),
		Status: domain.ProductionStatusCompleted, Priority: domain.WorkPriorityNormal,
		Lines: []domain.WorkOrderLine{{
			ID: uuid.New(), WorkTypeID: uuid.New(), PositionID: uuid.New(), Sequence: 1,
			Tasks: []domain.WorkOrderTask{
				{ID: uuid.New(), TaskID: taskID, Sequence: 1, Status: domain.TaskStatusCompleted},
			},
		}},
	}
	wo2 := &domain.WorkOrder{
		ID: uuid.New(), OrderNumber: "WO-002", OrderName: "Order 2",
		PartyID: "p1", WorkSetupID: uuidPtr(uuid.New()),
		Status: domain.ProductionStatusInProgress, Priority: domain.WorkPriorityNormal,
		Lines: []domain.WorkOrderLine{{
			ID: uuid.New(), WorkTypeID: uuid.New(), PositionID: uuid.New(), Sequence: 1,
			Tasks: []domain.WorkOrderTask{
				{ID: uuid.New(), TaskID: taskID, Sequence: 1, Status: domain.TaskStatusPending},
			},
		}},
	}

	svc := NewWorkOrderQueryService(
		newMockWorkOrderRepo(wo1, wo2),
		newMockTaskRepo(task),
	)

	results, err := svc.GetWorkOrdersProgress(context.Background(), []uuid.UUID{wo1.ID, wo2.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].CompletedTasks != 1 {
		t.Errorf("wo1: expected 1 completed, got %d", results[0].CompletedTasks)
	}
	if results[1].CompletedTasks != 0 {
		t.Errorf("wo2: expected 0 completed, got %d", results[1].CompletedTasks)
	}
}

func TestGetWorkOrdersProgress_Empty(t *testing.T) {
	svc := NewWorkOrderQueryService(
		newMockWorkOrderRepo(),
		newMockTaskRepo(),
	)

	results, err := svc.GetWorkOrdersProgress(context.Background(), []uuid.UUID{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestGetWorkOrdersProgress_SkipsMissing(t *testing.T) {
	taskID := uuid.New()
	task := &domain.Task{ID: taskID, Name: "Test", IsActive: true}

	wo := &domain.WorkOrder{
		ID: uuid.New(), OrderNumber: "WO-001", OrderName: "Existing",
		PartyID: "p1", WorkSetupID: uuidPtr(uuid.New()),
		Status: domain.ProductionStatusPending, Priority: domain.WorkPriorityNormal,
		Lines: []domain.WorkOrderLine{{
			ID: uuid.New(), WorkTypeID: uuid.New(), PositionID: uuid.New(), Sequence: 1,
			Tasks: []domain.WorkOrderTask{
				{ID: uuid.New(), TaskID: taskID, Sequence: 1, Status: domain.TaskStatusPending},
			},
		}},
	}

	svc := NewWorkOrderQueryService(
		newMockWorkOrderRepo(wo),
		newMockTaskRepo(task),
	)

	missingID := uuid.New()
	results, err := svc.GetWorkOrdersProgress(context.Background(), []uuid.UUID{wo.ID, missingID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result (skipping missing), got %d", len(results))
	}
}
