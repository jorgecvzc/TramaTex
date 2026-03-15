package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTask(t *testing.T) {
	task, err := NewTask("Diseñar", "Diseño inicial", true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task.ID == uuid.Nil {
		t.Fatal("expected generated id")
	}
	if task.Name != "Diseñar" {
		t.Fatalf("expected name Diseñar, got %s", task.Name)
	}

	_, err = NewTask("", "", true)
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}
}

func TestNewPosition(t *testing.T) {
	position, err := NewPosition("Espalda", "BACK", "Espalda completa", true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if position.ID == uuid.Nil {
		t.Fatal("expected generated id")
	}

	_, err = NewPosition("", "BACK", "", true)
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}

	_, err = NewPosition("Espalda", "", "", true)
	if err == nil {
		t.Fatal("expected validation error for empty code")
	}
}

func TestNewWorkType(t *testing.T) {
	taskID := uuid.New()
	wt, err := NewWorkType("Serigrafía", "1 color", true, []WorkTypeTask{{TaskID: taskID, Sequence: 1}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if wt.ID == uuid.Nil {
		t.Fatal("expected generated id")
	}
	if len(wt.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(wt.Tasks))
	}

	_, err = NewWorkType("", "", true, nil)
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}

	_, err = NewWorkType("Serigrafía", "", true, []WorkTypeTask{{TaskID: uuid.Nil, Sequence: 1}})
	if err == nil {
		t.Fatal("expected validation error for nil task id")
	}

	_, err = NewWorkType("Serigrafía", "", true, []WorkTypeTask{{TaskID: taskID, Sequence: 0}})
	if err == nil {
		t.Fatal("expected validation error for non-positive sequence")
	}
}

func TestNewWorkSetup(t *testing.T) {
	workTypeID := uuid.New()
	positionID := uuid.New()

	setup, err := NewWorkSetup(
		"Camisetas López",
		"party-1",
		uuid.New(),
		"Personalización estándar",
		true,
		[]WorkSetupLine{
			{
				ID:         uuid.New(),
				WorkTypeID: workTypeID,
				PositionID: positionID,
				Sequence:   1,
			},
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if setup.ID == uuid.Nil {
		t.Fatal("expected generated id")
	}
	if len(setup.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(setup.Lines))
	}

	_, err = NewWorkSetup("", "party-1", uuid.New(), "", true, nil)
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}

	_, err = NewWorkSetup("Setup", "", uuid.New(), "", true, nil)
	if err == nil {
		t.Fatal("expected validation error for empty party id")
	}

	_, err = NewWorkSetup("Setup", "party-1", uuid.Nil, "", true, nil)
	if err == nil {
		t.Fatal("expected validation error for nil tangible group id")
	}

	_, err = NewWorkSetup("Setup", "party-1", uuid.New(), "", true, []WorkSetupLine{{ID: uuid.New(), WorkTypeID: uuid.Nil, PositionID: positionID, Sequence: 1}})
	if err == nil {
		t.Fatal("expected validation error for nil work type id in line")
	}

	_, err = NewWorkSetup("Setup", "party-1", uuid.New(), "", true, []WorkSetupLine{{ID: uuid.New(), WorkTypeID: workTypeID, PositionID: uuid.Nil, Sequence: 1}})
	if err == nil {
		t.Fatal("expected validation error for nil position id in line")
	}
}

func TestNewWorkOrder(t *testing.T) {
	workTypeID := uuid.New()
	positionID := uuid.New()
	taskID := uuid.New()

	order, err := NewWorkOrder(
		"OT-2026-001",
		"Orden A",
		"party-1",
		uuid.New(),
		nil,
		"observaciones",
		ProductionStatusDraft,
		WorkPriorityNormal,
		nil,
		nil,
		nil,
		[]WorkOrderLine{
			{
				ID:         uuid.New(),
				WorkTypeID: workTypeID,
				PositionID: positionID,
				Sequence:   1,
				Tasks: []WorkOrderTask{{
					ID:       uuid.New(),
					TaskID:   taskID,
					Sequence: 1,
					Status:   TaskStatusPending,
				}},
			},
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.ID == uuid.Nil {
		t.Fatal("expected generated order id")
	}

	_, err = NewWorkOrder("", "Orden", "party-1", uuid.New(), nil, "", ProductionStatusDraft, WorkPriorityNormal, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected validation error for empty order number")
	}

	_, err = NewWorkOrder("OT-1", "", "party-1", uuid.New(), nil, "", ProductionStatusDraft, WorkPriorityNormal, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected validation error for empty order name")
	}

	_, err = NewWorkOrder("OT-1", "Orden", "", uuid.New(), nil, "", ProductionStatusDraft, WorkPriorityNormal, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected validation error for empty party id")
	}

	_, err = NewWorkOrder("OT-1", "Orden", "party-1", uuid.Nil, nil, "", ProductionStatusDraft, WorkPriorityNormal, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected validation error for nil tangible group id")
	}
}
