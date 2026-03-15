package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ===== SalesWorkStatus Validation Tests =====

func TestSalesWorkStatus_IsValid(t *testing.T) {
	validStatuses := []SalesWorkStatus{
		SalesWorkStatusDraft,
		SalesWorkStatusPending,
		SalesWorkStatusInProgress,
		SalesWorkStatusCompleted,
		SalesWorkStatusCanceled,
	}
	for _, s := range validStatuses {
		assert.NoError(t, s.IsValid(), "expected %s to be valid", s)
	}
}

func TestSalesWorkStatus_IsValid_Invalid(t *testing.T) {
	invalid := SalesWorkStatus("INVALID")
	err := invalid.IsValid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sales work status")
}

// ===== SalesWorkStatus Transition Tests =====

func TestSalesWorkStatus_Transitions_DraftToPending(t *testing.T) {
	assert.True(t, canTransitionSalesWork(SalesWorkStatusDraft, SalesWorkStatusPending))
}

func TestSalesWorkStatus_Transitions_DraftToCanceled(t *testing.T) {
	assert.True(t, canTransitionSalesWork(SalesWorkStatusDraft, SalesWorkStatusCanceled))
}

func TestSalesWorkStatus_Transitions_PendingToInProgress(t *testing.T) {
	assert.True(t, canTransitionSalesWork(SalesWorkStatusPending, SalesWorkStatusInProgress))
}

func TestSalesWorkStatus_Transitions_PendingToCanceled(t *testing.T) {
	assert.True(t, canTransitionSalesWork(SalesWorkStatusPending, SalesWorkStatusCanceled))
}

func TestSalesWorkStatus_Transitions_InProgressToCompleted(t *testing.T) {
	assert.True(t, canTransitionSalesWork(SalesWorkStatusInProgress, SalesWorkStatusCompleted))
}

func TestSalesWorkStatus_Transitions_InProgressToCanceled(t *testing.T) {
	assert.True(t, canTransitionSalesWork(SalesWorkStatusInProgress, SalesWorkStatusCanceled))
}

func TestSalesWorkStatus_Transitions_InvalidCompletedToAnything(t *testing.T) {
	assert.False(t, canTransitionSalesWork(SalesWorkStatusCompleted, SalesWorkStatusDraft))
	assert.False(t, canTransitionSalesWork(SalesWorkStatusCompleted, SalesWorkStatusPending))
	assert.False(t, canTransitionSalesWork(SalesWorkStatusCompleted, SalesWorkStatusCanceled))
}

func TestSalesWorkStatus_Transitions_InvalidCanceledToAnything(t *testing.T) {
	assert.False(t, canTransitionSalesWork(SalesWorkStatusCanceled, SalesWorkStatusDraft))
	assert.False(t, canTransitionSalesWork(SalesWorkStatusCanceled, SalesWorkStatusPending))
}

func TestSalesWorkStatus_Transitions_InvalidSkipStates(t *testing.T) {
	assert.False(t, canTransitionSalesWork(SalesWorkStatusDraft, SalesWorkStatusInProgress))
	assert.False(t, canTransitionSalesWork(SalesWorkStatusDraft, SalesWorkStatusCompleted))
	assert.False(t, canTransitionSalesWork(SalesWorkStatusPending, SalesWorkStatusCompleted))
}

// ===== SalesWorkSetup Struct Tests =====

func TestSalesWorkSetup_FieldsPopulated(t *testing.T) {
	setupID := uuid.New()
	ws := SalesWorkSetup{
		ID:           uuid.New(),
		WorkSetupID:  &setupID,
		WorkOrderID:  nil,
		Name:         "Serigrafía camisetas",
		Observations: "Logo en pecho izquierdo, tintas especiales",
		Status:       SalesWorkStatusDraft,
		Sequence:     1,
	}

	assert.NotEqual(t, uuid.Nil, ws.ID)
	assert.Equal(t, &setupID, ws.WorkSetupID)
	assert.Nil(t, ws.WorkOrderID)
	assert.Equal(t, "Serigrafía camisetas", ws.Name)
	assert.Equal(t, "Logo en pecho izquierdo, tintas especiales", ws.Observations)
	assert.Equal(t, SalesWorkStatusDraft, ws.Status)
	assert.Equal(t, 1, ws.Sequence)
}

func TestSalesWorkSetup_WithoutWorkSetupID(t *testing.T) {
	ws := SalesWorkSetup{
		ID:           uuid.New(),
		WorkSetupID:  nil,
		Name:         "Trabajo nuevo desde comercial",
		Observations: "Bordado personalizado, ver indicaciones adjuntas",
		Status:       SalesWorkStatusDraft,
		Sequence:     1,
	}

	assert.Nil(t, ws.WorkSetupID)
	assert.Equal(t, "Bordado personalizado, ver indicaciones adjuntas", ws.Observations)
}
