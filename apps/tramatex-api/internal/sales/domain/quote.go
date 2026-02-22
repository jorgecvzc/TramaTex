package domain

import (
	"time"

	"github.com/google/uuid"
)

type Quote struct {
	ID             uuid.UUID
	QuoteNumber    QuoteNumber
	PartyID        uuid.UUID
	QuoteDate      time.Time
	ExpirationDate time.Time
	Status         QuoteStatus
	LineItems      []QuoteLineItem
	Subtotal       Money
	TaxAmount      Money
	Total          Money
	Notes          string
}

type QuoteLineItem struct {
	ID                        uuid.UUID
	MESWorkID                 *uuid.UUID
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
	calculatedUnitPrice Money,
	manualUnitPrice *Money,
	calculatedDiscountPerUnit *Money,
	manualDiscountPerUnit *Money,
) (QuoteLineItem, error) {
	if productVariantID == uuid.Nil {
		return QuoteLineItem{}, NewValidationError("product variant ID cannot be empty")
	}
	if quantity <= 0 {
		return QuoteLineItem{}, NewValidationError("quantity must be greater than zero")
	}

	finalUnitPrice := calculatedUnitPrice
	if manualUnitPrice != nil {
		if manualUnitPrice.Currency() != calculatedUnitPrice.Currency() {
			return QuoteLineItem{}, NewValidationError("unit price currency mismatch")
		}
		finalUnitPrice = *manualUnitPrice
	}

	finalDiscount, err := resolveDiscount(calculatedDiscountPerUnit, manualDiscountPerUnit, calculatedUnitPrice.Currency())
	if err != nil {
		return QuoteLineItem{}, err
	}

	subtotal, err := calculateLineSubtotal(finalUnitPrice, finalDiscount, quantity)
	if err != nil {
		return QuoteLineItem{}, err
	}

	return QuoteLineItem{
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

func (q *Quote) RecalculateTotals() error {
	subtotal, err := sumLineItemSubtotals(q.LineItems)
	if err != nil {
		return err
	}
	q.Subtotal = subtotal
	total, err := q.Subtotal.Add(q.TaxAmount)
	if err != nil {
		return err
	}
	q.Total = total
	return nil
}

func resolveDiscount(calculated *Money, manual *Money, currency string) (Money, error) {
	if manual != nil {
		if manual.Currency() != currency {
			return Money{}, NewValidationError("discount currency mismatch")
		}
		return *manual, nil
	}
	if calculated != nil {
		if calculated.Currency() != currency {
			return Money{}, NewValidationError("discount currency mismatch")
		}
		return *calculated, nil
	}
	return NewMoney(0, currency)
}

func calculateLineSubtotal(unitPrice Money, discount Money, quantity int) (Money, error) {
	netUnit, err := unitPrice.Subtract(discount)
	if err != nil {
		return Money{}, err
	}
	return netUnit.Multiply(float64(quantity))
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
