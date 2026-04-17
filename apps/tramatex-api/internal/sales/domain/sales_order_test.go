package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===== NewSalesOrder Tests =====

func TestNewSalesOrder_Success(t *testing.T) {
	partyID := uuid.New()
	orderDate := time.Now()
	deliveryDate := orderDate.Add(7 * 24 * time.Hour)
	money, _ := NewMoney(100.0, "EUR")
	taxAmount, _ := NewMoney(21.0, "EUR")
	number, _ := NewOrderNumber("O/2026/0001")

	lineItem, err := NewOrderLineItem(
		uuid.New(),
		2,
		money,
		nil,
		0,
	)
	require.NoError(t, err)

	order, err := NewSalesOrder(
		number,
		partyID,
		orderDate,
		deliveryDate,
		[]OrderLineItem{lineItem},
		taxAmount,
		"Test order",
	)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, order.ID)
	assert.Equal(t, number, order.OrderNumber)
	assert.Equal(t, partyID, order.PartyID)
	assert.Equal(t, SalesOrderStatusPending, order.Status)
	assert.Equal(t, 1, len(order.LineItems))
	assert.Equal(t, "Test order", order.Notes)

	// Verify totals calculation
	assert.Equal(t, 200.0, order.Subtotal.Amount()) // 100 * 2
	assert.Equal(t, 21.0, order.TaxAmount.Amount())
	assert.Equal(t, 221.0, order.Total.Amount()) // 200 + 21
}

func TestNewSalesOrder_EmptyPartyID(t *testing.T) {
	orderDate := time.Now()
	deliveryDate := orderDate.Add(7 * 24 * time.Hour)
	taxAmount, _ := NewMoney(0, "EUR")
	number, _ := NewOrderNumber("O/2026/0001")

	order, err := NewSalesOrder(
		number,
		uuid.Nil,
		orderDate,
		deliveryDate,
		[]OrderLineItem{},
		taxAmount,
		"",
	)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "party ID cannot be empty")
}

func TestNewSalesOrder_DeliveryBeforeOrderDate(t *testing.T) {
	partyID := uuid.New()
	orderDate := time.Now()
	deliveryDate := orderDate.Add(-5 * 24 * time.Hour) // 5 days before
	taxAmount, _ := NewMoney(0, "EUR")
	number, _ := NewOrderNumber("O/2026/0001")

	order, err := NewSalesOrder(
		number,
		partyID,
		orderDate,
		deliveryDate,
		[]OrderLineItem{},
		taxAmount,
		"",
	)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "delivery date cannot be before order date")
}

func TestNewSalesOrder_MultipleLineItems(t *testing.T) {
	partyID := uuid.New()
	orderDate := time.Now()
	deliveryDate := orderDate.Add(7 * 24 * time.Hour)
	taxAmount, _ := NewMoney(15.0, "EUR")
	number, _ := NewOrderNumber("O/2026/0001")

	price1, _ := NewMoney(60.0, "EUR")
	price2, _ := NewMoney(40.0, "EUR")

	lineItem1, _ := NewOrderLineItem(uuid.New(), 2, price1, nil, 0)
	lineItem2, _ := NewOrderLineItem(uuid.New(), 3, price2, nil, 0)

	order, err := NewSalesOrder(
		number,
		partyID,
		orderDate,
		deliveryDate,
		[]OrderLineItem{lineItem1, lineItem2},
		taxAmount,
		"",
	)

	require.NoError(t, err)
	assert.Equal(t, 2, len(order.LineItems))
	// Subtotal = (60*2) + (40*3) = 120 + 120 = 240
	assert.Equal(t, 240.0, order.Subtotal.Amount())
	// Total = 240 + 15 = 255
	assert.Equal(t, 255.0, order.Total.Amount())
}

// ===== NewOrderLineItem Tests =====

func TestNewOrderLineItem_Success(t *testing.T) {
	variantID := uuid.New()
	price, _ := NewMoney(75.0, "EUR")

	lineItem, err := NewOrderLineItem(
		variantID,
		4,
		price,
		nil,
		0,
	)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, lineItem.ID)
	assert.Equal(t, variantID, lineItem.ProductVariantID)
	assert.Equal(t, 4, lineItem.Quantity)
	assert.Equal(t, 75.0, lineItem.ListUnitPrice.Amount())
	assert.Equal(t, 75.0, lineItem.UnitPrice.Amount())
	assert.Equal(t, 0.0, lineItem.DiscountPerUnit.Amount())
	assert.Equal(t, 300.0, lineItem.Subtotal.Amount()) // 75 * 4
}

func TestNewOrderLineItem_WithOverridePrice(t *testing.T) {
	variantID := uuid.New()
	listPrice, _ := NewMoney(100.0, "EUR")
	overridePrice, _ := NewMoney(85.0, "EUR")

	lineItem, err := NewOrderLineItem(
		variantID,
		3,
		listPrice,
		&overridePrice,
		0,
	)

	require.NoError(t, err)
	assert.Equal(t, 100.0, lineItem.ListUnitPrice.Amount())
	assert.Equal(t, 85.0, lineItem.UnitPrice.Amount()) // Override applied
	assert.Equal(t, 255.0, lineItem.Subtotal.Amount()) // 85 * 3
}

func TestNewOrderLineItem_WithDiscount(t *testing.T) {
	variantID := uuid.New()
	price, _ := NewMoney(100.0, "EUR")

	lineItem, err := NewOrderLineItem(
		variantID,
		2,
		price,
		nil,
		15.0,
	)

	require.NoError(t, err)
	assert.Equal(t, 100.0, lineItem.UnitPrice.Amount())
	assert.Equal(t, 15.0, lineItem.DiscountPerUnit.Amount())
	// Subtotal = (100 - 15) * 2 = 85 * 2 = 170
	assert.Equal(t, 170.0, lineItem.Subtotal.Amount())
}

func TestNewOrderLineItem_EmptyVariantID(t *testing.T) {
	price, _ := NewMoney(100.0, "EUR")

	_, err := NewOrderLineItem(
		uuid.Nil,
		1,
		price,
		nil,
		0,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product variant ID cannot be empty")
}

func TestNewOrderLineItem_ZeroQuantity(t *testing.T) {
	price, _ := NewMoney(100.0, "EUR")

	_, err := NewOrderLineItem(
		uuid.New(),
		0,
		price,
		nil,
		0,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantity must be greater than zero")
}

func TestNewOrderLineItem_CurrencyMismatch(t *testing.T) {
	variantID := uuid.New()
	listPrice, _ := NewMoney(100.0, "EUR")
	overridePrice, _ := NewMoney(90.0, "USD") // Different currency

	_, err := NewOrderLineItem(
		variantID,
		1,
		listPrice,
		&overridePrice,
		0,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "currency mismatch")
}

// ===== ChangeStatus Tests =====

func TestSalesOrder_ChangeStatus_PendingToInPreparation(t *testing.T) {
	order := createValidOrder(t)

	err := order.ChangeStatus(SalesOrderStatusInPreparation)

	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusInPreparation, order.Status)
}

func TestSalesOrder_ChangeStatus_PendingToCanceled(t *testing.T) {
	order := createValidOrder(t)

	err := order.ChangeStatus(SalesOrderStatusCancelled)

	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusCancelled, order.Status)
}

func TestSalesOrder_ChangeStatus_InPreparationToPartiallyDelivered(t *testing.T) {
	order := createValidOrder(t)
	order.Status = SalesOrderStatusInPreparation

	err := order.ChangeStatus(SalesOrderStatusPartiallyDelivered)

	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusPartiallyDelivered, order.Status)
}

func TestSalesOrder_ChangeStatus_InPreparationToDelivered(t *testing.T) {
	order := createValidOrder(t)
	order.Status = SalesOrderStatusInPreparation

	err := order.ChangeStatus(SalesOrderStatusDelivered)

	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusDelivered, order.Status)
}

func TestSalesOrder_ChangeStatus_PartiallyDeliveredToDelivered(t *testing.T) {
	order := createValidOrder(t)
	order.Status = SalesOrderStatusPartiallyDelivered

	err := order.ChangeStatus(SalesOrderStatusDelivered)

	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusDelivered, order.Status)
}

func TestSalesOrder_ChangeStatus_DeliveredToPartiallyInvoiced(t *testing.T) {
	order := createValidOrder(t)
	order.Status = SalesOrderStatusDelivered

	err := order.ChangeStatus(SalesOrderStatusPartiallyInvoiced)

	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusPartiallyInvoiced, order.Status)
}

func TestSalesOrder_ChangeStatus_DeliveredToInvoiced(t *testing.T) {
	order := createValidOrder(t)
	order.Status = SalesOrderStatusDelivered

	err := order.ChangeStatus(SalesOrderStatusInvoiced)

	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusInvoiced, order.Status)
}

func TestSalesOrder_ChangeStatus_PartiallyInvoicedToInvoiced(t *testing.T) {
	order := createValidOrder(t)
	order.Status = SalesOrderStatusPartiallyInvoiced

	err := order.ChangeStatus(SalesOrderStatusInvoiced)

	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusInvoiced, order.Status)
}

func TestSalesOrder_ChangeStatus_InvalidTransition_PendingToDelivered(t *testing.T) {
	order := createValidOrder(t)

	err := order.ChangeStatus(SalesOrderStatusDelivered)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sales order status transition")
	assert.Equal(t, SalesOrderStatusPending, order.Status) // Status unchanged
}

func TestSalesOrder_ChangeStatus_InvalidTransition_DeliveredToPending(t *testing.T) {
	order := createValidOrder(t)
	order.Status = SalesOrderStatusDelivered

	err := order.ChangeStatus(SalesOrderStatusPending)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sales order status transition")
}

func TestSalesOrder_ChangeStatus_ValidTransition_CanceledToPending(t *testing.T) {
	order := createValidOrder(t)
	order.Status = SalesOrderStatusCancelled

	err := order.ChangeStatus(SalesOrderStatusPending)

	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusPending, order.Status)
}

func TestSalesOrder_ChangeStatus_ReadyForProduction(t *testing.T) {
	order := createValidOrder(t)

	// From Pending to ReadyForProduction
	err := order.ChangeStatus(SalesOrderStatusReadyForProduction)
	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusReadyForProduction, order.Status)

	// Reset to InPreparation
	order.Status = SalesOrderStatusInPreparation
	err = order.ChangeStatus(SalesOrderStatusReadyForProduction)
	assert.NoError(t, err)
	assert.Equal(t, SalesOrderStatusReadyForProduction, order.Status)
}

func TestSalesOrder_ChangeStatus_InvalidStatus(t *testing.T) {
	order := createValidOrder(t)

	err := order.ChangeStatus(SalesOrderStatus("INVALID"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sales order status")
}

// ===== RecalculateTotals Tests =====

func TestSalesOrder_RecalculateTotals_SingleItem(t *testing.T) {
	order := createValidOrder(t)

	// Modify line item subtotal
	newPrice, _ := NewMoney(250.0, "EUR")
	order.LineItems[0].Subtotal = newPrice

	err := order.RecalculateTotals()

	require.NoError(t, err)
	assert.Equal(t, 250.0, order.Subtotal.Amount())
	// Total = 250 (subtotal) + 20 (tax)
	assert.Equal(t, 270.0, order.Total.Amount())
}

func TestSalesOrder_RecalculateTotals_MultipleItems(t *testing.T) {
	partyID := uuid.New()
	orderDate := time.Now()
	deliveryDate := orderDate.Add(7 * 24 * time.Hour)
	taxAmount, _ := NewMoney(25.0, "EUR")
	number, _ := NewOrderNumber("O/2026/0001")

	price1, _ := NewMoney(50.0, "EUR")
	price2, _ := NewMoney(30.0, "EUR")

	lineItem1, _ := NewOrderLineItem(uuid.New(), 3, price1, nil, 0)
	lineItem2, _ := NewOrderLineItem(uuid.New(), 2, price2, nil, 0)

	order, _ := NewSalesOrder(
		number,
		partyID,
		orderDate,
		deliveryDate,
		[]OrderLineItem{lineItem1, lineItem2},
		taxAmount,
		"",
	)

	// Modify line items
	newSubtotal1, _ := NewMoney(180.0, "EUR")
	newSubtotal2, _ := NewMoney(70.0, "EUR")
	order.LineItems[0].Subtotal = newSubtotal1
	order.LineItems[1].Subtotal = newSubtotal2

	err := order.RecalculateTotals()

	require.NoError(t, err)
	// Subtotal = 180 + 70 = 250
	assert.Equal(t, 250.0, order.Subtotal.Amount())
	// Total = 250 + 25 = 275
	assert.Equal(t, 275.0, order.Total.Amount())
}

// ===== Helper Functions =====

func createValidOrder(t *testing.T) *SalesOrder {
	partyID := uuid.New()
	orderDate := time.Now()
	deliveryDate := orderDate.Add(7 * 24 * time.Hour)
	price, _ := NewMoney(100.0, "EUR")
	taxAmount, _ := NewMoney(20.0, "EUR")
	number, _ := NewOrderNumber("O/2026/0001")

	lineItem, _ := NewOrderLineItem(
		uuid.New(),
		1,
		price,
		nil,
		0,
	)

	order, err := NewSalesOrder(
		number,
		partyID,
		orderDate,
		deliveryDate,
		[]OrderLineItem{lineItem},
		taxAmount,
		"",
	)
	require.NoError(t, err)
	return order
}
