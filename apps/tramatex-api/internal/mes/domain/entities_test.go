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

func TestNewServiceGroup(t *testing.T) {
	taskID := uuid.New()
	group, err := NewServiceGroup("Serigrafía", "1 color", nil, true, []ServiceGroupTask{{TaskID: taskID, Sequence: 1}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if group.ID == uuid.Nil {
		t.Fatal("expected generated id")
	}
	if len(group.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(group.Tasks))
	}

	_, err = NewServiceGroup("", "", nil, true, nil)
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}

	_, err = NewServiceGroup("Serigrafía", "", nil, true, []ServiceGroupTask{{TaskID: uuid.Nil, Sequence: 1}})
	if err == nil {
		t.Fatal("expected validation error for nil task id")
	}

	_, err = NewServiceGroup("Serigrafía", "", nil, true, []ServiceGroupTask{{TaskID: taskID, Sequence: 0}})
	if err == nil {
		t.Fatal("expected validation error for non-positive sequence")
	}
}

func TestNewMESWork(t *testing.T) {
	serviceGroupID := uuid.New()
	positionID := uuid.New()
	taskID := uuid.New()

	work, err := NewMESWork(
		"MES-2026-001",
		"Trabajo A",
		"party-1",
		uuid.New(),
		"observaciones",
		ProductionStatusDraft,
		WorkPriorityNormal,
		nil,
		nil,
		nil,
		[]MESWorkServiceGroup{
			{
				ID:             uuid.New(),
				ServiceGroupID: serviceGroupID,
				PositionID:     positionID,
				Sequence:       1,
				Tasks: []MESWorkTask{{
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
	if work.ID == uuid.Nil {
		t.Fatal("expected generated work id")
	}

	_, err = NewMESWork("", "Trabajo", "party-1", uuid.New(), "", ProductionStatusDraft, WorkPriorityNormal, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected validation error for empty work number")
	}

	_, err = NewMESWork("MES-1", "", "party-1", uuid.New(), "", ProductionStatusDraft, WorkPriorityNormal, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected validation error for empty work name")
	}

	_, err = NewMESWork("MES-1", "Trabajo", "", uuid.New(), "", ProductionStatusDraft, WorkPriorityNormal, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected validation error for empty party id")
	}

	_, err = NewMESWork("MES-1", "Trabajo", "party-1", uuid.Nil, "", ProductionStatusDraft, WorkPriorityNormal, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected validation error for nil tangible group id")
	}
}
