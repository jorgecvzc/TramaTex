package application

import (
	"time"

	"github.com/google/uuid"
)

type QuoteLineItemInput struct {
	ProductVariantID uuid.UUID `json:"productVariantId"`
	Quantity         int       `json:"quantity"`
	UnitPrice        *MoneyDTO `json:"unitPrice,omitempty"`
	DiscountPercent  *float64  `json:"discountPercent,omitempty"`
}

type MesWorkRefInput struct {
	MesWorkID    uuid.UUID `json:"mesWorkId"`
	Observations string    `json:"observations"`
}

type CreateQuoteCommand struct {
	PartyID        uuid.UUID            `json:"partyId"`
	ExpirationDate time.Time            `json:"expirationDate"`
	MesWorkRefs    []MesWorkRefInput    `json:"mesWorkRefs,omitempty"`
	Notes          *string              `json:"notes"`
	Items          []QuoteLineItemInput `json:"items"`
}

type UpdateQuoteCommand struct {
	QuoteID        uuid.UUID            `json:"-"`
	ExpirationDate *time.Time           `json:"expirationDate"`
	MesWorkRefs    []MesWorkRefInput    `json:"mesWorkRefs,omitempty"`
	Notes          *string              `json:"notes"`
	Items          []QuoteLineItemInput `json:"items"`
}

type ChangeQuoteStatusCommand struct {
	QuoteID   uuid.UUID `json:"-"`
	NewStatus string    `json:"newStatus"`
}

type DeleteQuoteCommand struct {
	QuoteID uuid.UUID `json:"-"`
}

type PreviewQuoteCommand struct {
	PartyID uuid.UUID            `json:"partyId"`
	Items   []QuoteLineItemInput `json:"items"`
}

type ConvertQuoteToOrderCommand struct {
	QuoteID      uuid.UUID `json:"quoteId"`
	DeliveryDate time.Time `json:"deliveryDate"`
}

type AcceptAndConvertQuoteCommand struct {
	QuoteID      uuid.UUID `json:"-"`
	DeliveryDate time.Time `json:"deliveryDate"`
}

type OrderLineItemInput struct {
	ProductVariantID uuid.UUID `json:"productVariantId"`
	Quantity         int       `json:"quantity"`
	UnitPrice        *MoneyDTO `json:"unitPrice,omitempty"`
	DiscountPercent  *float64  `json:"discountPercent,omitempty"`
}

type PreviewOrderCommand struct {
	PartyID uuid.UUID            `json:"partyId"`
	Items   []OrderLineItemInput `json:"items"`
}

type CreateOrderCommand struct {
	PartyID      uuid.UUID            `json:"partyId"`
	QuoteID      *uuid.UUID           `json:"quoteId"`
	DeliveryDate time.Time            `json:"deliveryDate"`
	MesWorkRefs  []MesWorkRefInput    `json:"mesWorkRefs,omitempty"`
	Notes        *string              `json:"notes"`
	Items        []OrderLineItemInput `json:"items"`
}

type UpdateOrderDetailsCommand struct {
	OrderID      uuid.UUID         `json:"-"`
	PartyID      *uuid.UUID        `json:"partyId"`
	DeliveryDate *time.Time        `json:"deliveryDate"`
	Notes        *string           `json:"notes"`
	MesWorkRefs  []MesWorkRefInput `json:"mesWorkRefs,omitempty"`
}

type ChangeOrderStatusCommand struct {
	OrderID   uuid.UUID `json:"-"`
	NewStatus string    `json:"newStatus"`
}

type ChangeInvoiceStatusCommand struct {
	InvoiceID uuid.UUID `json:"-"`
	NewStatus string    `json:"newStatus"`
}

type ChangeDeliveryNoteStatusCommand struct {
	DeliveryNoteID uuid.UUID `json:"-"`
	NewStatus      string    `json:"newStatus"`
}

type AddOrderLineItemCommand struct {
	OrderID uuid.UUID          `json:"-"`
	Item    OrderLineItemInput `json:"item"`
}

type UpdateOrderLineItemCommand struct {
	OrderID         uuid.UUID `json:"-"`
	LineItemID      uuid.UUID `json:"-"`
	Quantity        *int      `json:"quantity"`
	UnitPrice       *MoneyDTO `json:"unitPrice,omitempty"`
	DiscountPercent *float64  `json:"discountPercent,omitempty"`
}

type RemoveOrderLineItemCommand struct {
	OrderID    uuid.UUID `json:"-"`
	LineItemID uuid.UUID `json:"-"`
}

type DeliveryNoteLineItemInput struct {
	SalesOrderLineItemID uuid.UUID `json:"salesOrderLineItemId"`
	DeliveredQuantity    int       `json:"deliveredQuantity"`
}

type CreateDeliveryNoteCommand struct {
	SalesOrderID uuid.UUID                   `json:"salesOrderId"`
	DeliveryDate time.Time                   `json:"deliveryDate"`
	Notes        *string                     `json:"notes"`
	Items        []DeliveryNoteLineItemInput `json:"items"`
}

type CreateInvoiceCommand struct {
	PartyID         uuid.UUID   `json:"partyId"`
	SalesOrderIDs   []uuid.UUID `json:"salesOrderIds"`
	DeliveryNoteIDs []uuid.UUID `json:"deliveryNoteIds"`
	InvoiceDate     time.Time   `json:"invoiceDate"`
	DueDate         time.Time   `json:"dueDate"`
	PaymentTerms    *string     `json:"paymentTerms"`
}

// OrderLineItemInputSimplified is a simplified version for quick ticket creation
type OrderLineItemInputSimplified struct {
	ProductVariantID uuid.UUID `json:"productVariantId"`
	Quantity         int       `json:"quantity"`
	DiscountPercent  float64   `json:"discountPercent"` // 0-100, manual discount for ticket lines
}

// CreateSimplifiedInvoiceCommand creates a ticket (factura simplificada) for retail sales < 3,000 EUR
// This is optimized for fast TPV/POS workflow without requiring full order creation
type CreateSimplifiedInvoiceCommand struct {
	PartyID     uuid.UUID                      `json:"partyId"` // Can be CONSUMIDOR_FINAL generic party
	InvoiceDate time.Time                      `json:"invoiceDate"`
	Items       []OrderLineItemInputSimplified `json:"items"` // Simple: just variant ID + quantity
	// Series defaults to "FT" (Factura de Ticket)
}
