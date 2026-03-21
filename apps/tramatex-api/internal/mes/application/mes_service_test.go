package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/mes/domain"
)

func uuidPtr(id uuid.UUID) *uuid.UUID { return &id }

type fakeTaskRepo struct {
	saved   *domain.Task
	byID    *domain.Task
	all     []*domain.Task
	deleted uuid.UUID
}

func (r *fakeTaskRepo) Save(_ context.Context, task *domain.Task) error {
	r.saved = task
	return nil
}

func (r *fakeTaskRepo) FindByID(_ context.Context, _ uuid.UUID) (*domain.Task, error) {
	return r.byID, nil
}

func (r *fakeTaskRepo) FindAll(_ context.Context, _ *domain.TaskFilters) ([]*domain.Task, error) {
	return r.all, nil
}

func (r *fakeTaskRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.deleted = id
	return nil
}

func (r *fakeTaskRepo) Exists(_ context.Context, _ uuid.UUID) (bool, error) {
	return true, nil
}

type fakePositionRepo struct {
	saved   *domain.Position
	byID    *domain.Position
	all     []*domain.Position
	deleted uuid.UUID
}

func (r *fakePositionRepo) Save(_ context.Context, position *domain.Position) error {
	r.saved = position
	return nil
}
func (r *fakePositionRepo) FindByID(_ context.Context, _ uuid.UUID) (*domain.Position, error) {
	return r.byID, nil
}
func (r *fakePositionRepo) FindAll(_ context.Context, _ *domain.PositionFilters) ([]*domain.Position, error) {
	return r.all, nil
}
func (r *fakePositionRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.deleted = id
	return nil
}

type fakeWorkTypeRepo struct {
	saved   *domain.WorkType
	byID    *domain.WorkType
	all     []*domain.WorkType
	deleted uuid.UUID
}

func (r *fakeWorkTypeRepo) Save(_ context.Context, workType *domain.WorkType) error {
	r.saved = workType
	return nil
}

func (r *fakeWorkTypeRepo) FindByID(_ context.Context, _ uuid.UUID) (*domain.WorkType, error) {
	return r.byID, nil
}

func (r *fakeWorkTypeRepo) FindAll(_ context.Context, _ *domain.WorkTypeFilters) ([]*domain.WorkType, error) {
	return r.all, nil
}

func (r *fakeWorkTypeRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.deleted = id
	return nil
}

type fakeWorkOrderRepo struct {
	saved     *domain.WorkOrder
	byID      *domain.WorkOrder
	all       []*domain.WorkOrder
	yearCount int64
}

func (r *fakeWorkOrderRepo) Save(_ context.Context, workOrder *domain.WorkOrder) error {
	r.saved = workOrder
	return nil
}

func (r *fakeWorkOrderRepo) FindByID(_ context.Context, _ uuid.UUID) (*domain.WorkOrder, error) {
	return r.byID, nil
}

func (r *fakeWorkOrderRepo) FindAll(_ context.Context, _ *domain.WorkOrderFilters) ([]*domain.WorkOrder, error) {
	return r.all, nil
}

func (r *fakeWorkOrderRepo) CountByYear(_ context.Context, _ int) (int64, error) {
	return r.yearCount, nil
}

type fakeWorkSetupRepo struct {
	saved *domain.WorkSetup
	byID  *domain.WorkSetup
	all   []*domain.WorkSetup
}

func (r *fakeWorkSetupRepo) Save(_ context.Context, ws *domain.WorkSetup) error {
	r.saved = ws
	return nil
}

func (r *fakeWorkSetupRepo) FindByID(_ context.Context, _ uuid.UUID) (*domain.WorkSetup, error) {
	return r.byID, nil
}

func (r *fakeWorkSetupRepo) FindAll(_ context.Context, _ *domain.WorkSetupFilters) ([]*domain.WorkSetup, error) {
	return r.all, nil
}

func (r *fakeWorkSetupRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

func TestCreateTask_AppliesDefaultsAndPersists(t *testing.T) {
	taskRepo := &fakeTaskRepo{}
	service := NewMESService(taskRepo, &fakePositionRepo{}, &fakeWorkTypeRepo{}, &fakeWorkOrderRepo{}, nil)

	result, err := service.CreateTask(context.Background(), CreateTaskCommand{Name: "Diseñar"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected task dto")
	}
	if result.Name != "Diseñar" {
		t.Fatalf("expected name Diseñar, got %s", result.Name)
	}
	if !result.IsActive {
		t.Fatal("expected default is_active true")
	}
	if result.Description != "" {
		t.Fatalf("expected empty description default, got %s", result.Description)
	}
	if taskRepo.saved == nil {
		t.Fatal("expected repository save call")
	}
}

func TestCreateWorkType_MapsTaskAssignments(t *testing.T) {
	repo := &fakeWorkTypeRepo{}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, repo, &fakeWorkOrderRepo{}, nil)

	taskID := uuid.New()
	result, err := service.CreateWorkType(context.Background(), CreateWorkTypeCommand{
		Name: "Serigrafía",
		TaskAssignments: []WorkTypeTaskInput{
			{TaskID: taskID, Sequence: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected work type dto")
	}
	if len(result.Tasks) != 1 {
		t.Fatalf("expected 1 mapped task, got %d", len(result.Tasks))
	}
	if repo.saved == nil || len(repo.saved.Tasks) != 1 {
		t.Fatal("expected saved entity with mapped assignments")
	}
	if repo.saved.Tasks[0].TaskID != taskID || repo.saved.Tasks[0].Sequence != 1 {
		t.Fatal("expected mapped task assignment to match input")
	}
}

func TestUpdateWorkType_DoesNotOverrideTasksWhenNilAssignments(t *testing.T) {
	existingTaskID := uuid.New()
	repo := &fakeWorkTypeRepo{
		byID: &domain.WorkType{
			ID:       uuid.New(),
			Name:     "Bordado",
			IsActive: true,
			Tasks: []domain.WorkTypeTask{
				{TaskID: existingTaskID, Sequence: 1},
			},
		},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, repo, &fakeWorkOrderRepo{}, nil)

	newName := "Bordado premium"
	result, err := service.UpdateWorkType(context.Background(), UpdateWorkTypeCommand{
		ID:   repo.byID.ID,
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != newName {
		t.Fatalf("expected updated name, got %s", result.Name)
	}
	if len(repo.saved.Tasks) != 1 || repo.saved.Tasks[0].TaskID != existingTaskID {
		t.Fatal("expected existing tasks to remain unchanged when task_assignments is nil")
	}
}

func TestCreateWorkOrder_GeneratesTasksFromWorkSetup(t *testing.T) {
	taskID := uuid.New()
	workTypeID := uuid.New()
	positionID := uuid.New()
	workSetupID := uuid.New()
	workRepo := &fakeWorkOrderRepo{yearCount: 2}
	workTypeRepo := &fakeWorkTypeRepo{
		byID: &domain.WorkType{
			ID:       workTypeID,
			Name:     "Serigrafía",
			IsActive: true,
			Tasks: []domain.WorkTypeTask{
				{TaskID: taskID, Sequence: 1},
			},
		},
	}
	workSetupRepo := &fakeWorkSetupRepo{
		byID: &domain.WorkSetup{
			ID:      workSetupID,
			Name:    "Conf. Cliente A",
			PartyID: "party-001",
			Lines: []domain.WorkSetupLine{
				{WorkTypeID: workTypeID, PositionID: positionID, Sequence: 1},
			},
		},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, workTypeRepo, workRepo, workSetupRepo)

	result, err := service.CreateWorkOrder(context.Background(), CreateWorkOrderCommand{
		WorkName:    "Trabajo Cliente A",
		PartyID:     "party-001",
		WorkSetupID: &workSetupID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected work order dto")
	}
	if result.WorkNumber == "" {
		t.Fatal("expected generated work number")
	}
	if workRepo.saved == nil {
		t.Fatal("expected work order persistence")
	}
	if len(workRepo.saved.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(workRepo.saved.Lines))
	}
	if len(workRepo.saved.Lines[0].Tasks) != 1 {
		t.Fatalf("expected 1 generated task, got %d", len(workRepo.saved.Lines[0].Tasks))
	}
	if workRepo.saved.Lines[0].Tasks[0].TaskID != taskID {
		t.Fatal("expected generated task to match work type template task")
	}
}

func TestTaskFlows_ListUpdateDelete(t *testing.T) {
	taskID := uuid.New()
	taskRepo := &fakeTaskRepo{
		byID: &domain.Task{ID: taskID, Name: "Diseñar", Description: "old", IsActive: true},
		all:  []*domain.Task{{ID: taskID, Name: "Diseñar", IsActive: true}},
	}
	service := NewMESService(taskRepo, &fakePositionRepo{}, &fakeWorkTypeRepo{}, &fakeWorkOrderRepo{}, nil)

	list, err := service.ListTasks(context.Background(), ListTasksQuery{})
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list tasks success, err=%v len=%d", err, len(list))
	}

	name := "Diseñar Premium"
	isActive := false
	updated, err := service.UpdateTask(context.Background(), UpdateTaskCommand{ID: taskID, Name: &name, IsActive: &isActive})
	if err != nil {
		t.Fatalf("expected update success, got %v", err)
	}
	if updated.Name != name || updated.IsActive != isActive {
		t.Fatal("expected updated task values")
	}

	if err := service.DeleteTask(context.Background(), DeleteTaskCommand{ID: taskID}); err != nil {
		t.Fatalf("expected delete success, got %v", err)
	}
	if taskRepo.deleted != taskID {
		t.Fatal("expected deleted id to match")
	}
}

func TestPositionFlows_CreateGetListUpdateDelete(t *testing.T) {
	positionID := uuid.New()
	positionRepo := &fakePositionRepo{
		byID: &domain.Position{ID: positionID, Name: "Espalda", Code: "BACK", IsActive: true},
		all:  []*domain.Position{{ID: positionID, Name: "Espalda", Code: "BACK", IsActive: true}},
	}
	service := NewMESService(&fakeTaskRepo{}, positionRepo, &fakeWorkTypeRepo{}, &fakeWorkOrderRepo{}, nil)

	created, err := service.CreatePosition(context.Background(), CreatePositionCommand{Name: "Pecho", Code: "CHEST"})
	if err != nil || created == nil {
		t.Fatalf("expected create position success, err=%v", err)
	}

	got, err := service.GetPositionByID(context.Background(), GetPositionByIDQuery{ID: positionID})
	if err != nil || got == nil {
		t.Fatalf("expected get position success, err=%v", err)
	}

	list, err := service.ListPositions(context.Background(), ListPositionsQuery{})
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list positions success, err=%v len=%d", err, len(list))
	}

	name := "Espalda Completa"
	updated, err := service.UpdatePosition(context.Background(), UpdatePositionCommand{ID: positionID, Name: &name})
	if err != nil || updated.Name != name {
		t.Fatalf("expected update position success, err=%v", err)
	}

	if err := service.DeletePosition(context.Background(), DeletePositionCommand{ID: positionID}); err != nil {
		t.Fatalf("expected delete position success, got %v", err)
	}
	if positionRepo.deleted != positionID {
		t.Fatal("expected deleted position id to match")
	}
}

func TestWorkTypeFlows_ListGetDelete(t *testing.T) {
	id := uuid.New()
	repo := &fakeWorkTypeRepo{
		byID: &domain.WorkType{ID: id, Name: "Bordado", IsActive: true},
		all:  []*domain.WorkType{{ID: id, Name: "Bordado", IsActive: true}},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, repo, &fakeWorkOrderRepo{}, nil)

	got, err := service.GetWorkTypeByID(context.Background(), GetWorkTypeByIDQuery{ID: id})
	if err != nil || got == nil {
		t.Fatalf("expected get work type success, err=%v", err)
	}

	list, err := service.ListWorkTypes(context.Background(), ListWorkTypesQuery{})
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list work types success, err=%v len=%d", err, len(list))
	}

	if err := service.DeleteWorkType(context.Background(), DeleteWorkTypeCommand{ID: id}); err != nil {
		t.Fatalf("expected delete work type success, got %v", err)
	}
	if repo.deleted != id {
		t.Fatal("expected deleted work type id to match")
	}
}

func TestWorkOrderFlows_ListAndGet(t *testing.T) {
	id := uuid.New()
	workRepo := &fakeWorkOrderRepo{
		byID: &domain.WorkOrder{ID: id, OrderNumber: "MES-2026-001", OrderName: "Trabajo", PartyID: "party-1", WorkSetupID: uuidPtr(uuid.New()), Status: domain.ProductionStatusPending, Priority: domain.WorkPriorityNormal},
		all:  []*domain.WorkOrder{{ID: id, OrderNumber: "MES-2026-001", OrderName: "Trabajo", PartyID: "party-1", WorkSetupID: uuidPtr(uuid.New()), Status: domain.ProductionStatusPending, Priority: domain.WorkPriorityNormal}},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeWorkTypeRepo{}, workRepo, nil)

	list, err := service.ListWorkOrders(context.Background(), ListWorkOrdersQuery{})
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list work orders success, err=%v len=%d", err, len(list))
	}

	got, err := service.GetWorkOrderByID(context.Background(), GetWorkOrderByIDQuery{ID: id})
	if err != nil || got == nil {
		t.Fatalf("expected get work order success, err=%v", err)
	}
}

func TestNotFoundBranches_ReturnError(t *testing.T) {
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeWorkTypeRepo{}, &fakeWorkOrderRepo{}, nil)

	if _, err := service.GetTaskByID(context.Background(), GetTaskByIDQuery{ID: uuid.New()}); err == nil {
		t.Fatal("expected task not found error")
	}

	if _, err := service.GetPositionByID(context.Background(), GetPositionByIDQuery{ID: uuid.New()}); err == nil {
		t.Fatal("expected position not found error")
	}

	if _, err := service.GetWorkTypeByID(context.Background(), GetWorkTypeByIDQuery{ID: uuid.New()}); err == nil {
		t.Fatal("expected work type not found error")
	}

	if _, err := service.GetWorkOrderByID(context.Background(), GetWorkOrderByIDQuery{ID: uuid.New()}); err == nil {
		t.Fatal("expected work order not found error")
	}
}

func TestWorkOrderDashboardStatsAndOverdue(t *testing.T) {
	today := time.Now().UTC()
	yesterday := today.Add(-24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	overdueWork := &domain.WorkOrder{
		ID:          uuid.New(),
		OrderNumber: "MES-2026-010",
		OrderName:   "Trabajo vencido",
		PartyID:     "party-1",
		WorkSetupID: uuidPtr(uuid.New()),
		Status:      domain.ProductionStatusInProgress,
		Priority:    domain.WorkPriorityHigh,
		DueDate:     &yesterday,
	}

	dueTodayWork := &domain.WorkOrder{
		ID:          uuid.New(),
		OrderNumber: "MES-2026-011",
		OrderName:   "Trabajo hoy",
		PartyID:     "party-1",
		WorkSetupID: uuidPtr(uuid.New()),
		Status:      domain.ProductionStatusPending,
		Priority:    domain.WorkPriorityNormal,
		DueDate:     &today,
	}

	notOverdueCompleted := &domain.WorkOrder{
		ID:          uuid.New(),
		OrderNumber: "MES-2026-012",
		OrderName:   "Trabajo completado",
		PartyID:     "party-1",
		WorkSetupID: uuidPtr(uuid.New()),
		Status:      domain.ProductionStatusCompleted,
		Priority:    domain.WorkPriorityLow,
		DueDate:     &yesterday,
	}

	futureWork := &domain.WorkOrder{
		ID:          uuid.New(),
		OrderNumber: "MES-2026-013",
		OrderName:   "Trabajo futuro",
		PartyID:     "party-1",
		WorkSetupID: uuidPtr(uuid.New()),
		Status:      domain.ProductionStatusPending,
		Priority:    domain.WorkPriorityNormal,
		DueDate:     &tomorrow,
	}

	workRepo := &fakeWorkOrderRepo{
		all: []*domain.WorkOrder{overdueWork, dueTodayWork, notOverdueCompleted, futureWork},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeWorkTypeRepo{}, workRepo, nil)

	stats, err := service.GetWorkOrderDashboardStats(context.Background())
	if err != nil {
		t.Fatalf("expected dashboard stats success, got %v", err)
	}
	if stats.Total != 4 {
		t.Fatalf("expected total 4, got %d", stats.Total)
	}
	if stats.Overdue != 1 {
		t.Fatalf("expected overdue 1, got %d", stats.Overdue)
	}
	if stats.ByStatus[string(domain.ProductionStatusInProgress)] != 1 {
		t.Fatal("expected IN_PROGRESS count 1")
	}

	overdue, err := service.ListOverdueWorkOrders(context.Background(), ListOverdueWorkOrdersQuery{Limit: 5})
	if err != nil {
		t.Fatalf("expected overdue list success, got %v", err)
	}
	if len(overdue) != 1 {
		t.Fatalf("expected 1 overdue work, got %d", len(overdue))
	}
	if overdue[0].ID != overdueWork.ID {
		t.Fatal("expected overdue work id to match")
	}
}

func TestUpdateWorkOrderTaskStatus_StartAndCompleteFlow(t *testing.T) {
	taskID := uuid.New()
	workID := uuid.New()
	lineID := uuid.New()

	workRepo := &fakeWorkOrderRepo{
		byID: &domain.WorkOrder{
			ID:          workID,
			OrderNumber: "MES-2026-020",
			OrderName:   "Trabajo tablet",
			PartyID:     "party-1",
			WorkSetupID: uuidPtr(uuid.New()),
			Status:      domain.ProductionStatusPending,
			Priority:    domain.WorkPriorityNormal,
			Lines: []domain.WorkOrderLine{
				{
					ID:         lineID,
					WorkTypeID: uuid.New(),
					PositionID: uuid.New(),
					Sequence:   1,
					Tasks: []domain.WorkOrderTask{
						{ID: taskID, TaskID: uuid.New(), Sequence: 1, Status: domain.TaskStatusPending},
					},
				},
			},
		},
	}

	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeWorkTypeRepo{}, workRepo, nil)

	startResult, err := service.UpdateWorkOrderTaskStatus(context.Background(), UpdateWorkOrderTaskStatusCommand{
		WorkID: workID,
		TaskID: taskID,
		Action: "START",
	})
	if err != nil {
		t.Fatalf("expected start action success, got %v", err)
	}
	if startResult.Status != string(domain.ProductionStatusInProgress) {
		t.Fatalf("expected work status IN_PROGRESS after start, got %s", startResult.Status)
	}
	if startResult.Lines[0].Tasks[0].StartedAt == nil {
		t.Fatal("expected started_at to be set on START")
	}

	completeResult, err := service.UpdateWorkOrderTaskStatus(context.Background(), UpdateWorkOrderTaskStatusCommand{
		WorkID: workID,
		TaskID: taskID,
		Action: "COMPLETE",
	})
	if err != nil {
		t.Fatalf("expected complete action success, got %v", err)
	}
	if completeResult.Status != string(domain.ProductionStatusCompleted) {
		t.Fatalf("expected work status COMPLETED after complete, got %s", completeResult.Status)
	}
	if completeResult.CompletedDate == nil {
		t.Fatal("expected work completed_date after task completion")
	}
	if workRepo.saved == nil {
		t.Fatal("expected work order saved after status update")
	}
}

// ===== SuspendWorkOrders Tests =====

func TestSuspendWorkOrders_SuspendsPendingAndInProgress(t *testing.T) {
	woID1 := uuid.New()
	woID2 := uuid.New()

	wo1 := &domain.WorkOrder{ID: woID1, Status: domain.ProductionStatusPending, OrderNumber: "MES-1", OrderName: "A", PartyID: "p", Priority: domain.WorkPriorityNormal}
	wo2 := &domain.WorkOrder{ID: woID2, Status: domain.ProductionStatusInProgress, OrderNumber: "MES-2", OrderName: "B", PartyID: "p", Priority: domain.WorkPriorityNormal}

	var saved []*domain.WorkOrder
	workRepo := &fakeWorkOrderRepo{}
	origSave := workRepo.Save
	_ = origSave
	workRepo.byID = nil // will be overridden per call

	// Use a custom repo to track multiple saves and lookups
	repo := &multiWorkOrderRepo{
		orders: map[uuid.UUID]*domain.WorkOrder{woID1: wo1, woID2: wo2},
		saved:  &saved,
	}

	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeWorkTypeRepo{}, repo, nil)
	err := service.SuspendWorkOrders(context.Background(), []uuid.UUID{woID1, woID2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if wo1.Status != domain.ProductionStatusSuspended {
		t.Fatalf("expected PENDING → SUSPENDED, got %s", wo1.Status)
	}
	if wo2.Status != domain.ProductionStatusSuspended {
		t.Fatalf("expected IN_PROGRESS → SUSPENDED, got %s", wo2.Status)
	}
}

func TestSuspendWorkOrders_SkipsCompletedAndCancelled(t *testing.T) {
	woID1 := uuid.New()
	woID2 := uuid.New()

	wo1 := &domain.WorkOrder{ID: woID1, Status: domain.ProductionStatusCompleted, OrderNumber: "MES-1", OrderName: "A", PartyID: "p", Priority: domain.WorkPriorityNormal}
	wo2 := &domain.WorkOrder{ID: woID2, Status: domain.ProductionStatusCancelled, OrderNumber: "MES-2", OrderName: "B", PartyID: "p", Priority: domain.WorkPriorityNormal}

	repo := &multiWorkOrderRepo{
		orders: map[uuid.UUID]*domain.WorkOrder{woID1: wo1, woID2: wo2},
		saved:  &[]*domain.WorkOrder{},
	}

	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeWorkTypeRepo{}, repo, nil)
	err := service.SuspendWorkOrders(context.Background(), []uuid.UUID{woID1, woID2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if wo1.Status != domain.ProductionStatusCompleted {
		t.Fatalf("expected COMPLETED unchanged, got %s", wo1.Status)
	}
	if wo2.Status != domain.ProductionStatusCancelled {
		t.Fatalf("expected CANCELLED unchanged, got %s", wo2.Status)
	}
}

// ===== ReactivateWorkOrders Tests =====

func TestReactivateWorkOrders_ReactivatesSuspendedAndCancelled(t *testing.T) {
	woID1 := uuid.New()
	woID2 := uuid.New()

	wo1 := &domain.WorkOrder{ID: woID1, Status: domain.ProductionStatusSuspended, OrderNumber: "MES-1", OrderName: "A", PartyID: "p", Priority: domain.WorkPriorityNormal}
	wo2 := &domain.WorkOrder{ID: woID2, Status: domain.ProductionStatusCancelled, OrderNumber: "MES-2", OrderName: "B", PartyID: "p", Priority: domain.WorkPriorityNormal}

	repo := &multiWorkOrderRepo{
		orders: map[uuid.UUID]*domain.WorkOrder{woID1: wo1, woID2: wo2},
		saved:  &[]*domain.WorkOrder{},
	}

	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeWorkTypeRepo{}, repo, nil)
	err := service.ReactivateWorkOrders(context.Background(), []uuid.UUID{woID1, woID2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if wo1.Status != domain.ProductionStatusPending {
		t.Fatalf("expected SUSPENDED → PENDING, got %s", wo1.Status)
	}
	if wo2.Status != domain.ProductionStatusPending {
		t.Fatalf("expected CANCELLED → PENDING, got %s", wo2.Status)
	}
}

func TestReactivateWorkOrders_SkipsCompletedAndInProgress(t *testing.T) {
	woID1 := uuid.New()
	woID2 := uuid.New()

	wo1 := &domain.WorkOrder{ID: woID1, Status: domain.ProductionStatusCompleted, OrderNumber: "MES-1", OrderName: "A", PartyID: "p", Priority: domain.WorkPriorityNormal}
	wo2 := &domain.WorkOrder{ID: woID2, Status: domain.ProductionStatusInProgress, OrderNumber: "MES-2", OrderName: "B", PartyID: "p", Priority: domain.WorkPriorityNormal}

	repo := &multiWorkOrderRepo{
		orders: map[uuid.UUID]*domain.WorkOrder{woID1: wo1, woID2: wo2},
		saved:  &[]*domain.WorkOrder{},
	}

	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeWorkTypeRepo{}, repo, nil)
	err := service.ReactivateWorkOrders(context.Background(), []uuid.UUID{woID1, woID2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if wo1.Status != domain.ProductionStatusCompleted {
		t.Fatalf("expected COMPLETED unchanged, got %s", wo1.Status)
	}
	if wo2.Status != domain.ProductionStatusInProgress {
		t.Fatalf("expected IN_PROGRESS unchanged, got %s", wo2.Status)
	}
}

// ===== RecalculateWorkStatus guard for SUSPENDED =====

func TestRecalculateWorkStatus_DoesNotOverrideSuspended(t *testing.T) {
	taskID := uuid.New()
	workID := uuid.New()

	work := &domain.WorkOrder{
		ID:          workID,
		OrderNumber: "MES-2026-001",
		OrderName:   "Test",
		PartyID:     "party-1",
		Status:      domain.ProductionStatusSuspended,
		Priority:    domain.WorkPriorityNormal,
		Lines: []domain.WorkOrderLine{
			{
				ID:         uuid.New(),
				WorkTypeID: uuid.New(),
				PositionID: uuid.New(),
				Sequence:   1,
				Tasks: []domain.WorkOrderTask{
					{ID: uuid.New(), TaskID: taskID, Sequence: 1, Status: domain.TaskStatusPending},
				},
			},
		},
	}

	workRepo := &multiWorkOrderRepo{
		orders: map[uuid.UUID]*domain.WorkOrder{workID: work},
		saved:  &[]*domain.WorkOrder{},
	}

	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeWorkTypeRepo{}, workRepo, nil)

	// Try to complete a task — recalculateWorkStatus should NOT override SUSPENDED
	_, err := service.UpdateWorkOrderTaskStatus(context.Background(), UpdateWorkOrderTaskStatusCommand{
		WorkID: workID,
		TaskID: work.Lines[0].Tasks[0].ID,
		Action: "COMPLETE",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if work.Status != domain.ProductionStatusSuspended {
		t.Fatalf("expected SUSPENDED to be preserved, got %s", work.Status)
	}
}

// multiWorkOrderRepo supports lookups by ID for multiple work orders.
type multiWorkOrderRepo struct {
	orders map[uuid.UUID]*domain.WorkOrder
	saved  *[]*domain.WorkOrder
}

func (r *multiWorkOrderRepo) Save(_ context.Context, wo *domain.WorkOrder) error {
	*r.saved = append(*r.saved, wo)
	return nil
}

func (r *multiWorkOrderRepo) FindByID(_ context.Context, id uuid.UUID) (*domain.WorkOrder, error) {
	return r.orders[id], nil
}

func (r *multiWorkOrderRepo) FindAll(_ context.Context, _ *domain.WorkOrderFilters) ([]*domain.WorkOrder, error) {
	return nil, nil
}

func (r *multiWorkOrderRepo) CountByYear(_ context.Context, _ int) (int64, error) {
	return 0, nil
}
