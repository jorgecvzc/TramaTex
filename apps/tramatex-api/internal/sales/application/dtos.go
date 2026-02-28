package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

type MoneyDTO struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type QuoteLineItemDTO struct {
	ID                        uuid.UUID  `json:"id"`
	MesWorkID                 *uuid.UUID `json:"mesWorkId,omitempty"`
	ProductVariantID          uuid.UUID  `json:"productVariantId"`
	Quantity                  int        `json:"quantity"`
	CalculatedUnitPrice       MoneyDTO   `json:"calculatedUnitPrice"`
	ManualUnitPrice           *MoneyDTO  `json:"manualUnitPrice,omitempty"`
	FinalUnitPrice            MoneyDTO   `json:"finalUnitPrice"`
	TaxRate                   float64    `json:"taxRate"`
	CalculatedDiscountPerUnit *MoneyDTO  `json:"calculatedDiscountPerUnit,omitempty"`
	ManualDiscountPerUnit     *MoneyDTO  `json:"manualDiscountPerUnit,omitempty"`
	FinalDiscountPerUnit      MoneyDTO   `json:"finalDiscountPerUnit"`
	Subtotal                  MoneyDTO   `json:"subtotal"`
	TaxAmount                 MoneyDTO   `json:"taxAmount"`
}

type QuoteDTO struct {
	ID             uuid.UUID          `json:"id"`
	QuoteNumber    string             `json:"quoteNumber"`
	PartyID        uuid.UUID          `json:"partyId"`
	QuoteDate      time.Time          `json:"quoteDate"`
	ExpirationDate time.Time          `json:"expirationDate"`
	Status         string             `json:"status"`
	LineItems      []QuoteLineItemDTO `json:"lineItems"`
	Subtotal       MoneyDTO           `json:"subtotal"`
	TaxAmount      MoneyDTO           `json:"taxAmount"`
	Total          MoneyDTO           `json:"total"`
	Notes          string             `json:"notes"`
}

type OrderLineItemDTO struct {
	ID                        uuid.UUID  `json:"id"`
	MesWorkID                 *uuid.UUID `json:"mesWorkId,omitempty"`
	ProductVariantID          uuid.UUID  `json:"productVariantId"`
	Quantity                  int        `json:"quantity"`
	CalculatedUnitPrice       MoneyDTO   `json:"calculatedUnitPrice"`
	ManualUnitPrice           *MoneyDTO  `json:"manualUnitPrice,omitempty"`
	FinalUnitPrice            MoneyDTO   `json:"finalUnitPrice"`
	TaxRate                   float64    `json:"taxRate"`
	CalculatedDiscountPerUnit *MoneyDTO  `json:"calculatedDiscountPerUnit,omitempty"`
	ManualDiscountPerUnit     *MoneyDTO  `json:"manualDiscountPerUnit,omitempty"`
	FinalDiscountPerUnit      MoneyDTO   `json:"finalDiscountPerUnit"`
	Subtotal                  MoneyDTO   `json:"subtotal"`
	TaxAmount                 MoneyDTO   `json:"taxAmount"`
}

type SalesOrderDTO struct {
	ID           uuid.UUID          `json:"id"`
	OrderNumber  string             `json:"orderNumber"`
	QuoteID      *uuid.UUID         `json:"quoteId,omitempty"`
	PartyID      uuid.UUID          `json:"partyId"`
	OrderDate    time.Time          `json:"orderDate"`
	DeliveryDate time.Time          `json:"deliveryDate"`
	Status       string             `json:"status"`
	LineItems    []OrderLineItemDTO `json:"lineItems"`
	Subtotal     MoneyDTO           `json:"subtotal"`
	TaxAmount    MoneyDTO           `json:"taxAmount"`
	Total        MoneyDTO           `json:"total"`
	Notes        string             `json:"notes"`
}

type DeliveryNoteLineItemDTO struct {
	ID                   uuid.UUID `json:"id"`
	SalesOrderLineItemID uuid.UUID `json:"salesOrderLineItemId"`
	ProductVariantID     uuid.UUID `json:"productVariantId"`
	DeliveredQuantity    int       `json:"deliveredQuantity"`
}

type DeliveryNoteDTO struct {
	ID                 uuid.UUID                 `json:"id"`
	DeliveryNoteNumber string                    `json:"deliveryNoteNumber"`
	SalesOrderID       uuid.UUID                 `json:"salesOrderId"`
	PartyID            uuid.UUID                 `json:"partyId"`
	DeliveryDate       time.Time                 `json:"deliveryDate"`
	Status             string                    `json:"status"`
	LineItems          []DeliveryNoteLineItemDTO `json:"lineItems"`
	Notes              string                    `json:"notes"`
}

type InvoiceLineItemDTO struct {
	ID                   uuid.UUID  `json:"id"`
	SalesOrderLineItemID *uuid.UUID `json:"salesOrderLineItemId,omitempty"`
	ProductVariantID     uuid.UUID  `json:"productVariantId"`
	Quantity             int        `json:"quantity"`
	UnitPrice            MoneyDTO   `json:"unitPrice"`
	TaxRate              float64    `json:"taxRate"`
	DiscountAmount       *MoneyDTO  `json:"discountAmount,omitempty"`
	Subtotal             MoneyDTO   `json:"subtotal"`
	TaxAmount            *MoneyDTO  `json:"taxAmount,omitempty"`
	Total                MoneyDTO   `json:"total"`
}

type InvoiceDTO struct {
	ID              uuid.UUID            `json:"id"`
	InvoiceNumber   string               `json:"invoiceNumber"`
	PartyID         uuid.UUID            `json:"partyId"`
	InvoiceDate     time.Time            `json:"invoiceDate"`
	DueDate         time.Time            `json:"dueDate"`
	Status          string               `json:"status"`
	LineItems       []InvoiceLineItemDTO `json:"lineItems"`
	RelatedOrderIDs []uuid.UUID          `json:"relatedOrderIds,omitempty"`
	Subtotal        MoneyDTO             `json:"subtotal"`
	TaxAmount       MoneyDTO             `json:"taxAmount"`
	Total           MoneyDTO             `json:"total"`
	PaymentTerms    string               `json:"paymentTerms"`
}

func NewMoneyDTO(m domain.Money) MoneyDTO {
	return MoneyDTO{Amount: m.Amount(), Currency: m.Currency()}
}

func NewQuoteDTO(q *domain.Quote) *QuoteDTO {
	items := make([]QuoteLineItemDTO, 0, len(q.LineItems))
	for _, item := range q.LineItems {
		items = append(items, NewQuoteLineItemDTO(item))
	}
	return &QuoteDTO{
		ID:             q.ID,
		QuoteNumber:    q.QuoteNumber.String(),
		PartyID:        q.PartyID,
		QuoteDate:      q.QuoteDate,
		ExpirationDate: q.ExpirationDate,
		Status:         string(q.Status),
		LineItems:      items,
		Subtotal:       NewMoneyDTO(q.Subtotal),
		TaxAmount:      NewMoneyDTO(q.TaxAmount),
		Total:          NewMoneyDTO(q.Total),
		Notes:          q.Notes,
	}
}

func NewQuoteLineItemDTO(item domain.QuoteLineItem) QuoteLineItemDTO {
	return QuoteLineItemDTO{
		ID:                        item.ID,
		MesWorkID:                 item.MESWorkID,
		ProductVariantID:          item.ProductVariantID,
		Quantity:                  item.Quantity,
		CalculatedUnitPrice:       NewMoneyDTO(item.CalculatedUnitPrice),
		ManualUnitPrice:           toMoneyDTOPtr(item.ManualUnitPrice),
		FinalUnitPrice:            NewMoneyDTO(item.FinalUnitPrice),
		TaxRate:                   item.TaxRate,
		CalculatedDiscountPerUnit: toMoneyDTOPtr(item.CalculatedDiscountPerUnit),
		ManualDiscountPerUnit:     toMoneyDTOPtr(item.ManualDiscountPerUnit),
		FinalDiscountPerUnit:      NewMoneyDTO(item.FinalDiscountPerUnit),
		Subtotal:                  NewMoneyDTO(item.Subtotal),
		TaxAmount:                 NewMoneyDTO(item.TaxAmount),
	}
}

func NewSalesOrderDTO(order *domain.SalesOrder) *SalesOrderDTO {
	items := make([]OrderLineItemDTO, 0, len(order.LineItems))
	for _, item := range order.LineItems {
		items = append(items, NewOrderLineItemDTO(item))
	}
	return &SalesOrderDTO{
		ID:           order.ID,
		OrderNumber:  order.OrderNumber.String(),
		QuoteID:      order.QuoteID,
		PartyID:      order.PartyID,
		OrderDate:    order.OrderDate,
		DeliveryDate: order.DeliveryDate,
		Status:       string(order.Status),
		LineItems:    items,
		Subtotal:     NewMoneyDTO(order.Subtotal),
		TaxAmount:    NewMoneyDTO(order.TaxAmount),
		Total:        NewMoneyDTO(order.Total),
		Notes:        order.Notes,
	}
}

func NewOrderLineItemDTO(item domain.OrderLineItem) OrderLineItemDTO {
	return OrderLineItemDTO{
		ID:                        item.ID,
		MesWorkID:                 item.MESWorkID,
		ProductVariantID:          item.ProductVariantID,
		Quantity:                  item.Quantity,
		CalculatedUnitPrice:       NewMoneyDTO(item.CalculatedUnitPrice),
		ManualUnitPrice:           toMoneyDTOPtr(item.ManualUnitPrice),
		FinalUnitPrice:            NewMoneyDTO(item.FinalUnitPrice),
		TaxRate:                   item.TaxRate,
		CalculatedDiscountPerUnit: toMoneyDTOPtr(item.CalculatedDiscountPerUnit),
		ManualDiscountPerUnit:     toMoneyDTOPtr(item.ManualDiscountPerUnit),
		FinalDiscountPerUnit:      NewMoneyDTO(item.FinalDiscountPerUnit),
		Subtotal:                  NewMoneyDTO(item.Subtotal),
		TaxAmount:                 NewMoneyDTO(item.TaxAmount),
	}
}

func NewDeliveryNoteDTO(note *domain.DeliveryNote) *DeliveryNoteDTO {
	items := make([]DeliveryNoteLineItemDTO, 0, len(note.LineItems))
	for _, item := range note.LineItems {
		items = append(items, DeliveryNoteLineItemDTO{
			ID:                   item.ID,
			SalesOrderLineItemID: item.SalesOrderLineItemID,
			ProductVariantID:     item.ProductVariantID,
			DeliveredQuantity:    item.DeliveredQuantity,
		})
	}
	return &DeliveryNoteDTO{
		ID:                 note.ID,
		DeliveryNoteNumber: note.DeliveryNoteNumber.String(),
		SalesOrderID:       note.SalesOrderID,
		PartyID:            note.PartyID,
		DeliveryDate:       note.DeliveryDate,
		Status:             string(note.Status),
		LineItems:          items,
		Notes:              note.Notes,
	}
}

func NewInvoiceDTO(invoice *domain.Invoice, relatedOrderIDs []uuid.UUID) *InvoiceDTO {
	items := make([]InvoiceLineItemDTO, 0, len(invoice.LineItems))
	for _, item := range invoice.LineItems {
		items = append(items, NewInvoiceLineItemDTO(item))
	}
	return &InvoiceDTO{
		ID:              invoice.ID,
		InvoiceNumber:   invoice.InvoiceNumber.String(),
		PartyID:         invoice.PartyID,
		InvoiceDate:     invoice.InvoiceDate,
		DueDate:         invoice.DueDate,
		Status:          string(invoice.Status),
		LineItems:       items,
		RelatedOrderIDs: relatedOrderIDs,
		Subtotal:        NewMoneyDTO(invoice.Subtotal),
		TaxAmount:       NewMoneyDTO(invoice.TaxAmount),
		Total:           NewMoneyDTO(invoice.Total),
		PaymentTerms:    invoice.PaymentTerms,
	}
}

func NewInvoiceLineItemDTO(item domain.InvoiceLineItem) InvoiceLineItemDTO {
	return InvoiceLineItemDTO{
		ID:                   item.ID,
		SalesOrderLineItemID: item.SalesOrderLineItemID,
		ProductVariantID:     item.ProductVariantID,
		Quantity:             item.Quantity,
		UnitPrice:            NewMoneyDTO(item.UnitPrice),
		TaxRate:              item.TaxRate,
		DiscountAmount:       toMoneyDTOPtr(item.DiscountAmount),
		Subtotal:             NewMoneyDTO(item.Subtotal),
		TaxAmount:            toMoneyDTOPtr(item.TaxAmount),
		Total:                NewMoneyDTO(item.Total),
	}
}

func toMoneyDTOPtr(m *domain.Money) *MoneyDTO {
	if m == nil {
		return nil
	}
	dto := NewMoneyDTO(*m)
	return &dto
}
