package domain

import (
	"time"

	"github.com/google/uuid"
)

type SalesOrder struct {
	ID           uuid.UUID
	OrderNumber  OrderNumber
	QuoteID      *uuid.UUID
	PartyID      uuid.UUID
	OrderDate    time.Time
	DeliveryDate time.Time
	Status       SalesOrderStatus
	LineItems    []OrderLineItem
	Subtotal     Money
	TaxAmount    Money
	Total        Money
	Notes        string
}

type OrderLineItem struct {
	ID                        uuid.UUID
	ProductVariantID          uuid.UUID
	Quantity                  int
	CalculatedUnitPrice       Money
	ManualUnitPrice           *Money
	FinalUnitPrice            Money
	CalculatedDiscountPerUnit *Money
	ManualDiscountPerUnit     *Money
	FinalDiscountPerUnit      Money
	Subtotal                  Money
}

func NewSalesOrder(
	number OrderNumber,
	partyID uuid.UUID,
	orderDate time.Time,
	deliveryDate time.Time,
	lineItems []OrderLineItem,
	taxAmount Money,
	notes string,
) (*SalesOrder, error) {
	if partyID == uuid.Nil {
		return nil, NewValidationError("party ID cannot be empty")
	}
	if deliveryDate.Before(orderDate) {
		return nil, NewValidationError("delivery date cannot be before order date")
	}
	if err := SalesOrderStatusPending.IsValid(); err != nil {
		return nil, err
	}
	subtotal, err := sumOrderLineItemSubtotals(lineItems)
	if err != nil {
		return nil, err
	}
	total, err := subtotal.Add(taxAmount)
	if err != nil {
		return nil, err
	}

	return &SalesOrder{
		ID:           uuid.New(),
		OrderNumber:  number,
		PartyID:      partyID,
		OrderDate:    orderDate,
		DeliveryDate: deliveryDate,
		Status:       SalesOrderStatusPending,
		LineItems:    lineItems,
		Subtotal:     subtotal,
		TaxAmount:    taxAmount,
		Total:        total,
		Notes:        notes,
	}, nil
}

func NewOrderLineItem(
	productVariantID uuid.UUID,
	quantity int,
	calculatedUnitPrice Money,
	manualUnitPrice *Money,
	calculatedDiscountPerUnit *Money,
	manualDiscountPerUnit *Money,
) (OrderLineItem, error) {
	if productVariantID == uuid.Nil {
		return OrderLineItem{}, NewValidationError("product variant ID cannot be empty")
	}
	if quantity <= 0 {
		return OrderLineItem{}, NewValidationError("quantity must be greater than zero")
	}

	finalUnitPrice := calculatedUnitPrice
	if manualUnitPrice != nil {
		if manualUnitPrice.Currency() != calculatedUnitPrice.Currency() {
			return OrderLineItem{}, NewValidationError("unit price currency mismatch")
		}
		finalUnitPrice = *manualUnitPrice
	}

	finalDiscount, err := resolveDiscount(calculatedDiscountPerUnit, manualDiscountPerUnit, calculatedUnitPrice.Currency())
	if err != nil {
		return OrderLineItem{}, err
	}

	subtotal, err := calculateLineSubtotal(finalUnitPrice, finalDiscount, quantity)
	if err != nil {
		return OrderLineItem{}, err
	}

	return OrderLineItem{
		ID:                        uuid.New(),
		ProductVariantID:          productVariantID,
		Quantity:                  quantity,
		CalculatedUnitPrice:       calculatedUnitPrice,
		ManualUnitPrice:           manualUnitPrice,
		FinalUnitPrice:            finalUnitPrice,
		CalculatedDiscountPerUnit: calculatedDiscountPerUnit,
		ManualDiscountPerUnit:     manualDiscountPerUnit,
		FinalDiscountPerUnit:      finalDiscount,
		Subtotal:                  subtotal,
	}, nil
}

func (o *SalesOrder) ChangeStatus(newStatus SalesOrderStatus) error {
	if err := newStatus.IsValid(); err != nil {
		return err
	}
	if !canTransitionOrder(o.Status, newStatus) {
		return NewConflictError("invalid sales order status transition")
	}
	o.Status = newStatus
	return nil
}

func (o *SalesOrder) RecalculateTotals() error {
	subtotal, err := sumOrderLineItemSubtotals(o.LineItems)
	if err != nil {
		return err
	}
	o.Subtotal = subtotal
	total, err := o.Subtotal.Add(o.TaxAmount)
	if err != nil {
		return err
	}
	o.Total = total
	return nil
}

func sumOrderLineItemSubtotals(items []OrderLineItem) (Money, error) {
	subtotal, err := NewMoney(0, DefaultCurrency)
	if err != nil {
		return Money{}, err
	}
	for _, item := range items {
		subtotal, err = subtotal.Add(item.Subtotal)
		if err != nil {
			return Money{}, err
		}
	}
	return subtotal, nil
}
