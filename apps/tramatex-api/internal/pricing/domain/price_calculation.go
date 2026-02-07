package domain

import (
	"time"

	"github.com/google/uuid"
)

type PriceCalculation struct {
	ID               uuid.UUID
	ProductVariantID uuid.UUID
	ClientID         uuid.UUID
	Quantity         int
	BaseCost         Money
	FinalPrice       Money
	AppliedRules     []string
	CalculatedAt     time.Time
}

func NewPriceCalculation(
	variantID uuid.UUID,
	clientID uuid.UUID,
	quantity int,
	baseCost Money,
	finalPrice Money,
	appliedRules []string,
) (*PriceCalculation, error) {
	if variantID == uuid.Nil {
		return nil, NewValidationError("product variant ID cannot be empty")
	}
	if clientID == uuid.Nil {
		return nil, NewValidationError("client ID cannot be empty")
	}
	if quantity <= 0 {
		return nil, NewValidationError("quantity must be greater than zero")
	}

	return &PriceCalculation{
		ID:               uuid.New(),
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         quantity,
		BaseCost:         baseCost,
		FinalPrice:       finalPrice,
		AppliedRules:     appliedRules,
		CalculatedAt:     time.Now(),
	}, nil
}
