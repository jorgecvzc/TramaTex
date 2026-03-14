package domain

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryNote struct {
	ID                 uuid.UUID
	DeliveryNoteNumber DeliveryNoteNumber
	SalesOrderID       uuid.UUID
	PartyID            uuid.UUID
	DeliveryDate       time.Time
	Status             DeliveryNoteStatus
	LineItems          []DeliveryNoteLineItem
	Notes              string
}

type DeliveryNoteLineItem struct {
	ID                   uuid.UUID
	SalesOrderLineItemID uuid.UUID
	ProductVariantID     uuid.UUID
	DeliveredQuantity    int
	InvoiceLineItemID    *uuid.UUID
}

func NewDeliveryNote(
	number DeliveryNoteNumber,
	salesOrderID uuid.UUID,
	partyID uuid.UUID,
	deliveryDate time.Time,
	lineItems []DeliveryNoteLineItem,
	notes string,
) (*DeliveryNote, error) {
	if salesOrderID == uuid.Nil {
		return nil, NewValidationError("sales order ID cannot be empty")
	}
	if partyID == uuid.Nil {
		return nil, NewValidationError("party ID cannot be empty")
	}
	if err := DeliveryNoteStatusPending.IsValid(); err != nil {
		return nil, err
	}

	return &DeliveryNote{
		ID:                 uuid.New(),
		DeliveryNoteNumber: number,
		SalesOrderID:       salesOrderID,
		PartyID:            partyID,
		DeliveryDate:       deliveryDate,
		Status:             DeliveryNoteStatusPending,
		LineItems:          lineItems,
		Notes:              notes,
	}, nil
}

func NewDeliveryNoteLineItem(
	salesOrderLineItemID uuid.UUID,
	productVariantID uuid.UUID,
	deliveredQuantity int,
) (DeliveryNoteLineItem, error) {
	if salesOrderLineItemID == uuid.Nil {
		return DeliveryNoteLineItem{}, NewValidationError("sales order line item ID cannot be empty")
	}
	if productVariantID == uuid.Nil {
		return DeliveryNoteLineItem{}, NewValidationError("product variant ID cannot be empty")
	}
	if deliveredQuantity <= 0 {
		return DeliveryNoteLineItem{}, NewValidationError("delivered quantity must be greater than zero")
	}

	return DeliveryNoteLineItem{
		ID:                   uuid.New(),
		SalesOrderLineItemID: salesOrderLineItemID,
		ProductVariantID:     productVariantID,
		DeliveredQuantity:    deliveredQuantity,
	}, nil
}

func (d *DeliveryNote) ChangeStatus(newStatus DeliveryNoteStatus) error {
	if err := newStatus.IsValid(); err != nil {
		return err
	}
	if !canTransitionDeliveryNote(d.Status, newStatus) {
		return NewConflictError("invalid delivery note status transition")
	}
	d.Status = newStatus
	return nil
}
