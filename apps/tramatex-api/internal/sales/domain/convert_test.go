package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuote_ConvertToOrder_GeneratesNewWorkReferenceIDs(t *testing.T) {
	quote := createValidQuote(t)
	quote.Status = QuoteStatusApproved

	ws1 := uuid.New()
	ws2 := uuid.New()
	quote.WorkReferences = []WorkReference{
		{
			ID:          uuid.New(),
			WorkSetupID: &ws1,
			Sequence:    1,
			Description: "Bordado frontal",
		},
		{
			ID:          uuid.New(),
			WorkSetupID: &ws2,
			Sequence:    2,
			Description: "Serigrafia trasera",
		},
	}

	originalIDs := map[uuid.UUID]struct{}{
		quote.WorkReferences[0].ID: {},
		quote.WorkReferences[1].ID: {},
	}

	orderNumber, err := NewOrderNumber("PED-2026-0001")
	require.NoError(t, err)

	order, err := quote.ConvertToOrder(orderNumber, time.Now().Add(24*time.Hour))
	require.NoError(t, err)
	require.Len(t, order.WorkReferences, 2)

	seenOrderIDs := make(map[uuid.UUID]struct{}, len(order.WorkReferences))
	for i, wr := range order.WorkReferences {
		assert.NotEqual(t, uuid.Nil, wr.ID)
		_, existsInQuote := originalIDs[wr.ID]
		assert.False(t, existsInQuote, "order work reference ID must not reuse quote work reference ID")

		_, duplicateInOrder := seenOrderIDs[wr.ID]
		assert.False(t, duplicateInOrder, "order work reference IDs must be unique")
		seenOrderIDs[wr.ID] = struct{}{}

		assert.Equal(t, quote.WorkReferences[i].WorkSetupID, wr.WorkSetupID)
		assert.Equal(t, quote.WorkReferences[i].Sequence, wr.Sequence)
		assert.Equal(t, quote.WorkReferences[i].Description, wr.Description)
	}
}
