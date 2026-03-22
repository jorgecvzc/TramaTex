package domain

import (
	"time"

	"github.com/google/uuid"
)

type SalesOrder struct {
	ID             uuid.UUID
	OrderNumber    OrderNumber
	QuoteID        *uuid.UUID
	PartyID        uuid.UUID
	OrderDate      time.Time
	DeliveryDate   time.Time
	Status         SalesOrderStatus
	WorkReferences []WorkReference // Document-level MES work references with observations
	LineItems      []OrderLineItem
	Subtotal       Money
	TaxAmount      Money
	Total          Money
	Notes          string
}

type OrderLineItem struct {
	ID               uuid.UUID
	ProductVariantID uuid.UUID
	Quantity         int
	ListUnitPrice    Money   // Precio de tarifa (from pricing engine)
	UnitPrice        Money   // Precio de venta (defaults to list, user can override)
	TaxRate          float64 // Tax rate as percentage (e.g., 21.0 = 21%)
	DiscountPercent  float64 // Discount percentage entered by user (source of truth)
	DiscountPerUnit  Money   // Discount amount per unit (derived from percent)
	Subtotal         Money
	TaxAmount        Money // Tax calculated from Subtotal * TaxRate
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
	listUnitPrice Money,
	unitPrice *Money,
	discountPercent float64,
	taxRateOptional ...float64,
) (OrderLineItem, error) {
	if productVariantID == uuid.Nil {
		return OrderLineItem{}, NewValidationError("product variant ID cannot be empty")
	}
	if quantity <= 0 {
		return OrderLineItem{}, NewValidationError("quantity must be greater than zero")
	}
	taxRate := 0.0
	if len(taxRateOptional) > 0 {
		taxRate = taxRateOptional[0]
	}

	if taxRate < 0 || taxRate > 100 {
		return OrderLineItem{}, NewValidationError("tax rate must be between 0 and 100")
	}
	if discountPercent < 0 || discountPercent > 100 {
		return OrderLineItem{}, NewValidationError("discount percent must be between 0 and 100")
	}

	finalUnitPrice := listUnitPrice
	if unitPrice != nil {
		if unitPrice.Currency() != listUnitPrice.Currency() {
			return OrderLineItem{}, NewValidationError("unit price currency mismatch")
		}
		finalUnitPrice = *unitPrice
	}

	discount, err := resolveDiscountFromPercent(discountPercent, finalUnitPrice)
	if err != nil {
		return OrderLineItem{}, err
	}

	subtotal, err := calculateLineSubtotal(finalUnitPrice, discount, quantity)
	if err != nil {
		return OrderLineItem{}, err
	}

	taxAmount, err := calculateTaxAmount(subtotal, taxRate)
	if err != nil {
		return OrderLineItem{}, err
	}

	return OrderLineItem{
		ID:               uuid.New(),
		ProductVariantID: productVariantID,
		Quantity:         quantity,
		ListUnitPrice:    listUnitPrice,
		UnitPrice:        finalUnitPrice,
		TaxRate:          taxRate,
		DiscountPercent:  discountPercent,
		DiscountPerUnit:  discount,
		Subtotal:         subtotal,
		TaxAmount:        taxAmount,
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

// NewOrderLineItemFromCalculated creates a line item from pre-calculated values
// provided by the Pricing engine. This ensures a single source of truth for all
// monetary calculations (discount, subtotal, tax).
func NewOrderLineItemFromCalculated(
	productVariantID uuid.UUID,
	quantity int,
	listUnitPrice Money,
	unitPrice Money,
	discountPercent float64,
	discountPerUnit Money,
	subtotal Money,
	taxRate float64,
	taxAmount Money,
) (OrderLineItem, error) {
	if productVariantID == uuid.Nil {
		return OrderLineItem{}, NewValidationError("product variant ID cannot be empty")
	}
	if quantity <= 0 {
		return OrderLineItem{}, NewValidationError("quantity must be greater than zero")
	}
	if taxRate < 0 || taxRate > 100 {
		return OrderLineItem{}, NewValidationError("tax rate must be between 0 and 100")
	}
	if discountPercent < 0 || discountPercent > 100 {
		return OrderLineItem{}, NewValidationError("discount percent must be between 0 and 100")
	}
	return OrderLineItem{
		ID:               uuid.New(),
		ProductVariantID: productVariantID,
		Quantity:         quantity,
		ListUnitPrice:    listUnitPrice,
		UnitPrice:        unitPrice,
		TaxRate:          taxRate,
		DiscountPercent:  discountPercent,
		DiscountPerUnit:  discountPerUnit,
		Subtotal:         subtotal,
		TaxAmount:        taxAmount,
	}, nil
}

func (o *SalesOrder) RecalculateTotals() error {
	subtotal, err := sumOrderLineItemSubtotals(o.LineItems)
	if err != nil {
		return err
	}
	o.Subtotal = subtotal

	// Recalculate tax from line items
	taxTotal, err := sumOrderLineItemTaxAmounts(o.LineItems)
	if err != nil {
		return err
	}
	if taxTotal.Amount() == 0 && o.TaxAmount.Amount() > 0 {
		taxTotal = o.TaxAmount
	}
	o.TaxAmount = taxTotal

	total, err := o.Subtotal.Add(o.TaxAmount)
	if err != nil {
		return err
	}
	o.Total = total
	return nil
}

func sumOrderLineItemSubtotals(items []OrderLineItem) (Money, error) {
	amounts := make([]Money, len(items))
	for i, item := range items {
		amounts[i] = item.Subtotal
	}
	return SumAmounts(amounts)
}

func sumOrderLineItemTaxAmounts(items []OrderLineItem) (Money, error) {
	amounts := make([]Money, len(items))
	for i, item := range items {
		amounts[i] = item.TaxAmount
	}
	return SumAmounts(amounts)
}
