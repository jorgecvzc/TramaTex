package domain

import "testing"

func TestProductionStatusValidation(t *testing.T) {
	valid := []ProductionStatus{
		ProductionStatusPending,
		ProductionStatusInProgress,
		ProductionStatusOnHold,
		ProductionStatusSuspended,
		ProductionStatusCompleted,
		ProductionStatusCancelled,
	}
	for _, status := range valid {
		if !status.IsValid() {
			t.Fatalf("expected status %s to be valid", status)
		}
	}

	if ProductionStatus("UNKNOWN").IsValid() {
		t.Fatal("expected UNKNOWN production status to be invalid")
	}
}

func TestTaskStatusValidation(t *testing.T) {
	valid := []TaskStatus{
		TaskStatusPending,
		TaskStatusInProgress,
		TaskStatusCompleted,
		TaskStatusBlocked,
		TaskStatusSkipped,
	}
	for _, status := range valid {
		if !status.IsValid() {
			t.Fatalf("expected task status %s to be valid", status)
		}
	}

	if TaskStatus("UNKNOWN").IsValid() {
		t.Fatal("expected unknown task status to be invalid")
	}
}

func TestWorkPriorityValidation(t *testing.T) {
	valid := []WorkPriority{WorkPriorityLow, WorkPriorityNormal, WorkPriorityHigh, WorkPriorityUrgent}
	for _, priority := range valid {
		if !priority.IsValid() {
			t.Fatalf("expected priority %s to be valid", priority)
		}
	}

	if WorkPriority("ASAP").IsValid() {
		t.Fatal("expected invalid work priority")
	}
}
