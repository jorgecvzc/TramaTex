package domain

import (
	"time"

	"github.com/google/uuid"
)

// SalesWorkSetup represents a work/personalization job associated with a sales document.
// It bridges Sales and MES: the commercial team defines what work is needed,
// and the workshop picks it up for execution.
type SalesWorkSetup struct {
	ID           uuid.UUID
	WorkSetupID  *uuid.UUID      // Optional reference to an existing MES WorkSetup template
	WorkOrderID  *uuid.UUID      // Populated when MES creates the WorkOrder for execution
	Name         string          // Descriptive name (e.g. "Serigrafía camisetas Confecciones López")
	Observations string          // Free-text: work characteristics, special instructions, notes
	Status       SalesWorkStatus // Lifecycle within Sales (BORRADOR → PENDIENTE → EN_PROCESO → COMPLETADO)
	Sequence     int             // Order within the document
}

type Quote struct {
	ID              uuid.UUID
	QuoteNumber     QuoteNumber
	PartyID         uuid.UUID
	QuoteDate       time.Time
	ExpirationDate  time.Time
	Status          QuoteStatus
	SalesWorkSetups []SalesWorkSetup // Document-level MES work references with observations
	LineItems       []QuoteLineItem
	Subtotal        Money
	TaxAmount       Money
	Total           Money
	Notes           string
}

type QuoteLineItem struct {
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

func NewQuote(
	number QuoteNumber,
	partyID uuid.UUID,
	quoteDate time.Time,
	expirationDate time.Time,
	lineItems []QuoteLineItem,
	taxAmount Money,
	notes string,
) (*Quote, error) {
	if partyID == uuid.Nil {
		return nil, NewValidationError("party ID cannot be empty")
	}
	if expirationDate.Before(quoteDate) {
		return nil, NewValidationError("expiration date cannot be before quote date")
	}
	if err := QuoteStatusDraft.IsValid(); err != nil {
		return nil, err
	}
	subtotal, err := sumLineItemSubtotals(lineItems)
	if err != nil {
		return nil, err
	}
	total, err := subtotal.Add(taxAmount)
	if err != nil {
		return nil, err
	}

	return &Quote{
		ID:             uuid.New(),
		QuoteNumber:    number,
		PartyID:        partyID,
		QuoteDate:      quoteDate,
		ExpirationDate: expirationDate,
		Status:         QuoteStatusDraft,
		LineItems:      lineItems,
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		Total:          total,
		Notes:          notes,
	}, nil
}

func NewQuoteLineItem(
	productVariantID uuid.UUID,
	quantity int,
	listUnitPrice Money,
	unitPrice *Money,
	discountPercent float64,
	taxRateOptional ...float64,
) (QuoteLineItem, error) {
	if productVariantID == uuid.Nil {
		return QuoteLineItem{}, NewValidationError("product variant ID cannot be empty")
	}
	if quantity <= 0 {
		return QuoteLineItem{}, NewValidationError("quantity must be greater than zero")
	}
	taxRate := 0.0
	if len(taxRateOptional) > 0 {
		taxRate = taxRateOptional[0]
	}

	if taxRate < 0 || taxRate > 100 {
		return QuoteLineItem{}, NewValidationError("tax rate must be between 0 and 100")
	}
	if discountPercent < 0 || discountPercent > 100 {
		return QuoteLineItem{}, NewValidationError("discount percent must be between 0 and 100")
	}

	finalUnitPrice := listUnitPrice
	if unitPrice != nil {
		if unitPrice.Currency() != listUnitPrice.Currency() {
			return QuoteLineItem{}, NewValidationError("unit price currency mismatch")
		}
		finalUnitPrice = *unitPrice
	}

	discount, err := resolveDiscountFromPercent(discountPercent, finalUnitPrice)
	if err != nil {
		return QuoteLineItem{}, err
	}

	subtotal, err := calculateLineSubtotal(finalUnitPrice, discount, quantity)
	if err != nil {
		return QuoteLineItem{}, err
	}

	taxAmount, err := calculateTaxAmount(subtotal, taxRate)
	if err != nil {
		return QuoteLineItem{}, err
	}

	return QuoteLineItem{
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

func (q *Quote) ChangeStatus(newStatus QuoteStatus) error {
	if err := newStatus.IsValid(); err != nil {
		return err
	}
	if !canTransitionQuote(q.Status, newStatus) {
		return NewConflictError("invalid quote status transition")
	}
	q.Status = newStatus
	return nil
}

// NewQuoteLineItemFromCalculated creates a line item from pre-calculated values
// provided by the Pricing engine. This ensures a single source of truth for all
// monetary calculations (discount, subtotal, tax).
func NewQuoteLineItemFromCalculated(
	productVariantID uuid.UUID,
	quantity int,
	listUnitPrice Money,
	unitPrice Money,
	discountPercent float64,
	discountPerUnit Money,
	subtotal Money,
	taxRate float64,
	taxAmount Money,
) (QuoteLineItem, error) {
	if productVariantID == uuid.Nil {
		return QuoteLineItem{}, NewValidationError("product variant ID cannot be empty")
	}
	if quantity <= 0 {
		return QuoteLineItem{}, NewValidationError("quantity must be greater than zero")
	}
	if taxRate < 0 || taxRate > 100 {
		return QuoteLineItem{}, NewValidationError("tax rate must be between 0 and 100")
	}
	if discountPercent < 0 || discountPercent > 100 {
		return QuoteLineItem{}, NewValidationError("discount percent must be between 0 and 100")
	}
	return QuoteLineItem{
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

func (q *Quote) RecalculateTotals() error {
	subtotal, err := sumLineItemSubtotals(q.LineItems)
	if err != nil {
		return err
	}
	q.Subtotal = subtotal

	// Recalculate tax from line items
	taxTotal, err := sumLineItemTaxAmounts(q.LineItems)
	if err != nil {
		return err
	}
	if taxTotal.Amount() == 0 && q.TaxAmount.Amount() > 0 {
		taxTotal = q.TaxAmount
	}
	q.TaxAmount = taxTotal

	total, err := q.Subtotal.Add(q.TaxAmount)
	if err != nil {
		return err
	}
	q.Total = total
	return nil
}

func resolveDiscountFromPercent(percent float64, unitPrice Money) (Money, error) {
	if percent <= 0 {
		return NewMoney(0, unitPrice.Currency())
	}
	discountAmount := unitPrice.Amount() * percent / 100
	return NewMoney(discountAmount, unitPrice.Currency())
}

func calculateLineSubtotal(unitPrice Money, discount Money, quantity int) (Money, error) {
	netUnit, err := unitPrice.Subtract(discount)
	if err != nil {
		return Money{}, err
	}
	return netUnit.Multiply(float64(quantity))
}

// calculateTaxAmount calculates tax amount from subtotal and tax rate percentage
func calculateTaxAmount(subtotal Money, taxRate float64) (Money, error) {
	if taxRate == 0 {
		return NewMoney(0, subtotal.Currency())
	}
	// taxAmount = subtotal * (taxRate / 100)
	taxMultiplier := taxRate / 100.0
	return subtotal.Multiply(taxMultiplier)
}
func sumLineItemSubtotals(items []QuoteLineItem) (Money, error) {
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

func sumLineItemTaxAmounts(items []QuoteLineItem) (Money, error) {
	taxTotal, err := NewMoney(0, DefaultCurrency)
	if err != nil {
		return Money{}, err
	}
	for _, item := range items {
		taxTotal, err = taxTotal.Add(item.TaxAmount)
		if err != nil {
			return Money{}, err
		}
	}
	return taxTotal, nil
}
