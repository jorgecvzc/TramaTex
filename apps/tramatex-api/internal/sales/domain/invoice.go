package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID            uuid.UUID
	InvoiceNumber InvoiceNumber
	Type          InvoiceType   // COMPLETA | SIMPLIFICADA (ticket)
	Series        InvoiceSeries // Series de numeración (A, TKT, etc.)
	PartyID       uuid.UUID
	InvoiceDate   time.Time
	DueDate       time.Time
	Status        InvoiceStatus
	LineItems     []InvoiceLineItem
	Subtotal      Money
	TaxAmount     Money
	Total         Money
	PaymentTerms  string
}

type InvoiceLineItem struct {
	ID                   uuid.UUID
	SalesOrderLineItemID *uuid.UUID
	ProductVariantID     uuid.UUID
	Quantity             int
	UnitPrice            Money
	TaxRate              float64 // Tax rate as percentage (e.g., 21.0 = 21%)
	DiscountAmount       *Money
	Subtotal             Money
	TaxAmount            *Money
	Total                Money
}

func NewInvoice(
	number InvoiceNumber,
	invoiceType InvoiceType,
	series InvoiceSeries,
	partyID uuid.UUID,
	invoiceDate time.Time,
	dueDate time.Time,
	lineItems []InvoiceLineItem,
	taxAmount Money,
	paymentTerms string,
) (*Invoice, error) {
	if partyID == uuid.Nil {
		return nil, NewValidationError("party ID cannot be empty")
	}
	if err := invoiceType.IsValid(); err != nil {
		return nil, err
	}
	if dueDate.Before(invoiceDate) {
		return nil, NewValidationError("due date cannot be before invoice date")
	}
	if err := InvoiceStatusDraft.IsValid(); err != nil {
		return nil, err
	}

	subtotal, err := sumInvoiceLineItemSubtotals(lineItems)
	if err != nil {
		return nil, err
	}
	calculatedTax, err := sumInvoiceLineItemTaxAmounts(lineItems)
	if err != nil {
		return nil, err
	}
	total, err := subtotal.Add(calculatedTax)
	if err != nil {
		return nil, err
	}

	invoice := &Invoice{
		ID:            uuid.New(),
		InvoiceNumber: number,
		Type:          invoiceType,
		Series:        series,
		PartyID:       partyID,
		InvoiceDate:   invoiceDate,
		DueDate:       dueDate,
		Status:        InvoiceStatusDraft,
		LineItems:     lineItems,
		Subtotal:      subtotal,
		TaxAmount:     calculatedTax,
		Total:         total,
		PaymentTerms:  paymentTerms,
	}

	// Validate legal limits for simplified invoices (tickets)
	if err := invoice.ValidateLegalLimits(); err != nil {
		return nil, err
	}

	return invoice, nil
}

func NewInvoiceLineItem(
	productVariantID uuid.UUID,
	quantity int,
	unitPrice Money,
	discountAmount *Money,
	taxAmount *Money,
	taxRateOptional ...float64,
) (InvoiceLineItem, error) {
	if productVariantID == uuid.Nil {
		return InvoiceLineItem{}, NewValidationError("product variant ID cannot be empty")
	}
	if quantity <= 0 {
		return InvoiceLineItem{}, NewValidationError("quantity must be greater than zero")
	}
	taxRate := 0.0
	if len(taxRateOptional) > 0 {
		taxRate = taxRateOptional[0]
	}

	if taxRate < 0 || taxRate > 100 {
		return InvoiceLineItem{}, NewValidationError("tax rate must be between 0 and 100")
	}
	if discountAmount != nil && discountAmount.Currency() != unitPrice.Currency() {
		return InvoiceLineItem{}, NewValidationError("discount currency mismatch")
	}
	if taxAmount != nil && taxAmount.Currency() != unitPrice.Currency() {
		return InvoiceLineItem{}, NewValidationError("tax currency mismatch")
	}

	discount := Money{}
	if discountAmount != nil {
		discount = *discountAmount
	} else {
		zero, err := NewMoney(0, unitPrice.Currency())
		if err != nil {
			return InvoiceLineItem{}, err
		}
		discount = zero
	}

	subtotal, err := calculateLineSubtotal(unitPrice, discount, quantity)
	if err != nil {
		return InvoiceLineItem{}, err
	}

	lineTax := Money{}
	if taxAmount != nil {
		lineTax = *taxAmount
	} else {
		// Calculate tax from subtotal if not provided
		lineTax, err = calculateTaxAmount(subtotal, taxRate)
		if err != nil {
			return InvoiceLineItem{}, err
		}
	}

	total, err := subtotal.Add(lineTax)
	if err != nil {
		return InvoiceLineItem{}, err
	}

	return InvoiceLineItem{
		ID:               uuid.New(),
		ProductVariantID: productVariantID,
		Quantity:         quantity,
		UnitPrice:        unitPrice,
		TaxRate:          taxRate,
		DiscountAmount:   discountAmount,
		Subtotal:         subtotal,
		TaxAmount:        &lineTax,
		Total:            total,
	}, nil
}

func (i *Invoice) ChangeStatus(newStatus InvoiceStatus) error {
	if err := newStatus.IsValid(); err != nil {
		return err
	}
	if !canTransitionInvoice(i.Status, newStatus) {
		return NewConflictError("invalid invoice status transition")
	}
	i.Status = newStatus
	return nil
}

func (i *Invoice) RecalculateTotals() error {
	subtotal, err := sumInvoiceLineItemSubtotals(i.LineItems)
	if err != nil {
		return err
	}
	i.Subtotal = subtotal
	calculatedTax, err := sumInvoiceLineItemTaxAmounts(i.LineItems)
	if err != nil {
		return err
	}
	i.TaxAmount = calculatedTax
	total, err := i.Subtotal.Add(i.TaxAmount)
	if err != nil {
		return err
	}
	i.Total = total

	// Re-validate legal limits after recalculation
	if err := i.ValidateLegalLimits(); err != nil {
		return err
	}

	return nil
}

// ValidateLegalLimits validates that simplified invoices (tickets) comply with Spanish legislation
// SIMPLIFICADA invoices must have Total < 3,000 EUR according to Real Decreto 1619/2012
func (i *Invoice) ValidateLegalLimits() error {
	if i.Type == InvoiceTypeSimplified {
		// Spanish legislation: simplified invoices (tickets) must be < 3,000 EUR
		const maxSimplifiedAmount = 3000.0
		if i.Total.Amount() >= maxSimplifiedAmount {
			return NewValidationError(fmt.Sprintf(
				"simplified invoice (ticket) total %.2f EUR exceeds legal limit of %.2f EUR",
				i.Total.Amount(),
				maxSimplifiedAmount,
			))
		}
	}
	return nil
}

func sumInvoiceLineItemTaxAmounts(items []InvoiceLineItem) (Money, error) {
	var amounts []Money
	for _, item := range items {
		if item.TaxAmount != nil {
			amounts = append(amounts, *item.TaxAmount)
		}
	}
	return SumAmounts(amounts)
}

func sumInvoiceLineItemSubtotals(items []InvoiceLineItem) (Money, error) {
	amounts := make([]Money, len(items))
	for i, item := range items {
		amounts[i] = item.Subtotal
	}
	return SumAmounts(amounts)
}
