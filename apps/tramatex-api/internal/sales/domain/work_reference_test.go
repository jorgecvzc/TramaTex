package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ===== WorkReference Struct Tests =====

func TestWorkReference_FieldsPopulated(t *testing.T) {
	setupID := uuid.New()
	ws := WorkReference{
		ID:          uuid.New(),
		WorkSetupID: &setupID,
		WorkOrderID: nil,
		Description: "Logo en pecho izquierdo, tintas especiales",
		Sequence:    1,
	}

	assert.NotEqual(t, uuid.Nil, ws.ID)
	assert.Equal(t, &setupID, ws.WorkSetupID)
	assert.Nil(t, ws.WorkOrderID)
	assert.Equal(t, "Logo en pecho izquierdo, tintas especiales", ws.Description)
	assert.Equal(t, 1, ws.Sequence)
}

func TestWorkReference_WithoutWorkSetupID(t *testing.T) {
	ws := WorkReference{
		ID:          uuid.New(),
		WorkSetupID: nil,
		Description: "Bordado personalizado, ver indicaciones adjuntas",
		Sequence:    1,
	}

	assert.Nil(t, ws.WorkSetupID)
	assert.Equal(t, "Bordado personalizado, ver indicaciones adjuntas", ws.Description)
}
