package application

import (
	"time"

	"github.com/google/uuid"
)

type QuoteLineItemInput struct {
	ProductVariantID      uuid.UUID `json:"productVariantId"`
	Quantity              int       `json:"quantity"`
	ManualUnitPrice       *MoneyDTO `json:"manualUnitPrice,omitempty"`
	ManualDiscountPerUnit *MoneyDTO `json:"manualDiscountPerUnit,omitempty"`
}

type CreateQuoteCommand struct {
	PartyID        uuid.UUID            `json:"partyId"`
	ExpirationDate time.Time            `json:"expirationDate"`
	Notes          *string              `json:"notes"`
	Items          []QuoteLineItemInput `json:"items"`
}

type UpdateQuoteCommand struct {
	QuoteID        uuid.UUID            `json:"-"`
	ExpirationDate *time.Time           `json:"expirationDate"`
	Notes          *string              `json:"notes"`
	Items          []QuoteLineItemInput `json:"items"`
}

type ChangeQuoteStatusCommand struct {
	QuoteID   uuid.UUID `json:"-"`
	NewStatus string    `json:"newStatus"`
}

type ConvertQuoteToOrderCommand struct {
	QuoteID      uuid.UUID `json:"quoteId"`
	DeliveryDate time.Time `json:"deliveryDate"`
}

type OrderLineItemInput struct {
	ProductVariantID      uuid.UUID `json:"productVariantId"`
	Quantity              int       `json:"quantity"`
	ManualUnitPrice       *MoneyDTO `json:"manualUnitPrice,omitempty"`
	ManualDiscountPerUnit *MoneyDTO `json:"manualDiscountPerUnit,omitempty"`
}

type CreateOrderCommand struct {
	PartyID      uuid.UUID            `json:"partyId"`
	QuoteID      *uuid.UUID           `json:"quoteId"`
	DeliveryDate time.Time            `json:"deliveryDate"`
	Notes        *string              `json:"notes"`
	Items        []OrderLineItemInput `json:"items"`
}

type UpdateOrderDetailsCommand struct {
	OrderID      uuid.UUID  `json:"-"`
	PartyID      *uuid.UUID `json:"partyId"`
	DeliveryDate *time.Time `json:"deliveryDate"`
	Notes        *string    `json:"notes"`
}

type ChangeOrderStatusCommand struct {
	OrderID   uuid.UUID `json:"-"`
	NewStatus string    `json:"newStatus"`
}

type AddOrderLineItemCommand struct {
	OrderID uuid.UUID          `json:"-"`
	Item    OrderLineItemInput `json:"item"`
}

type UpdateOrderLineItemCommand struct {
	OrderID               uuid.UUID `json:"-"`
	LineItemID            uuid.UUID `json:"-"`
	Quantity              *int      `json:"quantity"`
	ManualUnitPrice       *MoneyDTO `json:"manualUnitPrice,omitempty"`
	ManualDiscountPerUnit *MoneyDTO `json:"manualDiscountPerUnit,omitempty"`
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
