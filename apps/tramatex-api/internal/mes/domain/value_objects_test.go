package domain

import "testing"

func TestProductionStatusValidation(t *testing.T) {
	valid := []ProductionStatus{
		ProductionStatusDraft,
		ProductionStatusPending,
		ProductionStatusInProgress,
		ProductionStatusOnHold,
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

func TestQualityAndRecipeAndTaskTypeValidation(t *testing.T) {
	if !QualityStatusPassed.IsValid() || !QualityStatusPending.IsValid() {
		t.Fatal("expected quality statuses to be valid")
	}
	if QualityStatus("NOPE").IsValid() {
		t.Fatal("expected invalid quality status")
	}

	if !RecipeTypePhysicalProduct.IsValid() || !RecipeTypeServiceProduct.IsValid() {
		t.Fatal("expected recipe types to be valid")
	}
	if RecipeType("OTHER").IsValid() {
		t.Fatal("expected invalid recipe type")
	}

	if !TaskTypeOneTime.IsValid() || !TaskTypeRecurrent.IsValid() {
		t.Fatal("expected task types to be valid")
	}
	if TaskType("TEMP").IsValid() {
		t.Fatal("expected invalid task type")
	}
}
