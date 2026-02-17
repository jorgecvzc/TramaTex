package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===== NewDeliveryNote Tests =====

func TestNewDeliveryNote_Success(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	partyID := uuid.New()
	deliveryDate := time.Now()

	lineItem, err := NewDeliveryNoteLineItem(
		uuid.New(),
		uuid.New(),
		10,
	)
	require.NoError(t, err)

	dn, err := NewDeliveryNote(
		number,
		salesOrderID,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{lineItem},
		"Test delivery note",
	)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, dn.ID)
	assert.Equal(t, number, dn.DeliveryNoteNumber)
	assert.Equal(t, salesOrderID, dn.SalesOrderID)
	assert.Equal(t, partyID, dn.PartyID)
	assert.Equal(t, DeliveryNoteStatusPending, dn.Status)
	assert.Equal(t, 1, len(dn.LineItems))
	assert.Equal(t, "Test delivery note", dn.Notes)
}

func TestNewDeliveryNote_EmptySalesOrderID(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	partyID := uuid.New()
	deliveryDate := time.Now()

	dn, err := NewDeliveryNote(
		number,
		uuid.Nil,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{},
		"",
	)

	assert.Error(t, err)
	assert.Nil(t, dn)
	assert.Contains(t, err.Error(), "sales order ID cannot be empty")
}

func TestNewDeliveryNote_EmptyPartyID(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	deliveryDate := time.Now()

	dn, err := NewDeliveryNote(
		number,
		salesOrderID,
		uuid.Nil,
		deliveryDate,
		[]DeliveryNoteLineItem{},
		"",
	)

	assert.Error(t, err)
	assert.Nil(t, dn)
	assert.Contains(t, err.Error(), "party ID cannot be empty")
}

func TestNewDeliveryNote_WithMultipleLineItems(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	partyID := uuid.New()
	deliveryDate := time.Now()

	lineItem1, _ := NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 5)
	lineItem2, _ := NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 10)
	lineItem3, _ := NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 3)

	dn, err := NewDeliveryNote(
		number,
		salesOrderID,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{lineItem1, lineItem2, lineItem3},
		"Multiple items delivery",
	)

	require.NoError(t, err)
	assert.Equal(t, 3, len(dn.LineItems))
	assert.Equal(t, 5, dn.LineItems[0].DeliveredQuantity)
	assert.Equal(t, 10, dn.LineItems[1].DeliveredQuantity)
	assert.Equal(t, 3, dn.LineItems[2].DeliveredQuantity)
}

func TestNewDeliveryNote_EmptyLineItems(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	partyID := uuid.New()
	deliveryDate := time.Now()

	dn, err := NewDeliveryNote(
		number,
		salesOrderID,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{},
		"",
	)

	require.NoError(t, err)
	assert.Equal(t, 0, len(dn.LineItems))
}

// ===== NewDeliveryNoteLineItem Tests =====

func TestNewDeliveryNoteLineItem_Success(t *testing.T) {
	orderLineItemID := uuid.New()
	variantID := uuid.New()
	quantity := 15

	lineItem, err := NewDeliveryNoteLineItem(
		orderLineItemID,
		variantID,
		quantity,
	)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, lineItem.ID)
	assert.Equal(t, orderLineItemID, lineItem.SalesOrderLineItemID)
	assert.Equal(t, variantID, lineItem.ProductVariantID)
	assert.Equal(t, quantity, lineItem.DeliveredQuantity)
}

func TestNewDeliveryNoteLineItem_EmptySalesOrderLineItemID(t *testing.T) {
	variantID := uuid.New()

	_, err := NewDeliveryNoteLineItem(
		uuid.Nil,
		variantID,
		10,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sales order line item ID cannot be empty")
}

func TestNewDeliveryNoteLineItem_EmptyProductVariantID(t *testing.T) {
	orderLineItemID := uuid.New()

	_, err := NewDeliveryNoteLineItem(
		orderLineItemID,
		uuid.Nil,
		10,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product variant ID cannot be empty")
}

func TestNewDeliveryNoteLineItem_ZeroQuantity(t *testing.T) {
	orderLineItemID := uuid.New()
	variantID := uuid.New()

	_, err := NewDeliveryNoteLineItem(
		orderLineItemID,
		variantID,
		0,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delivered quantity must be greater than zero")
}

func TestNewDeliveryNoteLineItem_NegativeQuantity(t *testing.T) {
	orderLineItemID := uuid.New()
	variantID := uuid.New()

	_, err := NewDeliveryNoteLineItem(
		orderLineItemID,
		variantID,
		-5,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delivered quantity must be greater than zero")
}

// ===== ChangeStatus Tests =====

func TestDeliveryNote_ChangeStatus_PendingToDelivered_Success(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	partyID := uuid.New()
	deliveryDate := time.Now()

	lineItem, _ := NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 10)

	dn, _ := NewDeliveryNote(
		number,
		salesOrderID,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{lineItem},
		"",
	)

	err := dn.ChangeStatus(DeliveryNoteStatusDelivered)

	assert.NoError(t, err)
	assert.Equal(t, DeliveryNoteStatusDelivered, dn.Status)
}

func TestDeliveryNote_ChangeStatus_PendingToCanceled_Success(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	partyID := uuid.New()
	deliveryDate := time.Now()

	lineItem, _ := NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 10)

	dn, _ := NewDeliveryNote(
		number,
		salesOrderID,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{lineItem},
		"",
	)

	err := dn.ChangeStatus(DeliveryNoteStatusCanceled)

	assert.NoError(t, err)
	assert.Equal(t, DeliveryNoteStatusCanceled, dn.Status)
}

func TestDeliveryNote_ChangeStatus_DeliveredToAnything_Fail(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	partyID := uuid.New()
	deliveryDate := time.Now()

	lineItem, _ := NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 10)

	dn, _ := NewDeliveryNote(
		number,
		salesOrderID,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{lineItem},
		"",
	)

	// First transition to Delivered
	_ = dn.ChangeStatus(DeliveryNoteStatusDelivered)

	// Try to transition to Canceled (should fail)
	err := dn.ChangeStatus(DeliveryNoteStatusCanceled)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid delivery note status transition")
	assert.Equal(t, DeliveryNoteStatusDelivered, dn.Status) // Status remains unchanged
}

func TestDeliveryNote_ChangeStatus_CanceledToAnything_Fail(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	partyID := uuid.New()
	deliveryDate := time.Now()

	lineItem, _ := NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 10)

	dn, _ := NewDeliveryNote(
		number,
		salesOrderID,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{lineItem},
		"",
	)

	// First transition to Canceled
	_ = dn.ChangeStatus(DeliveryNoteStatusCanceled)

	// Try to transition to Delivered (should fail)
	err := dn.ChangeStatus(DeliveryNoteStatusDelivered)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid delivery note status transition")
	assert.Equal(t, DeliveryNoteStatusCanceled, dn.Status) // Status remains unchanged
}

func TestDeliveryNote_ChangeStatus_InvalidStatus_Fail(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	partyID := uuid.New()
	deliveryDate := time.Now()

	lineItem, _ := NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 10)

	dn, _ := NewDeliveryNote(
		number,
		salesOrderID,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{lineItem},
		"",
	)

	err := dn.ChangeStatus(DeliveryNoteStatus("INVALID_STATUS"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid delivery note status")
	assert.Equal(t, DeliveryNoteStatusPending, dn.Status) // Status remains unchanged
}

func TestDeliveryNote_ChangeStatus_PendingToPending_Fail(t *testing.T) {
	number, _ := NewDeliveryNoteNumber("DN/2026/0001")
	salesOrderID := uuid.New()
	partyID := uuid.New()
	deliveryDate := time.Now()

	lineItem, _ := NewDeliveryNoteLineItem(uuid.New(), uuid.New(), 10)

	dn, _ := NewDeliveryNote(
		number,
		salesOrderID,
		partyID,
		deliveryDate,
		[]DeliveryNoteLineItem{lineItem},
		"",
	)

	err := dn.ChangeStatus(DeliveryNoteStatusPending)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid delivery note status transition")
	assert.Equal(t, DeliveryNoteStatusPending, dn.Status)
}
