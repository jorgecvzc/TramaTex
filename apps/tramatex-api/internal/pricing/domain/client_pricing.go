package domain

import (
	"time"

	"github.com/google/uuid"
)

type ClientPricing struct {
	ID               uuid.UUID
	ClientID         uuid.UUID
	ProductVariantID uuid.UUID
	FixedPrice       Money
	EffectiveFrom    time.Time
	EffectiveTo      *time.Time
	IsActive         bool
}

func NewClientPricing(
	clientID uuid.UUID,
	productVariantID uuid.UUID,
	fixedPrice Money,
	effectiveFrom time.Time,
	effectiveTo *time.Time,
) (*ClientPricing, error) {
	if clientID == uuid.Nil {
		return nil, NewValidationError("client ID cannot be empty")
	}
	if productVariantID == uuid.Nil {
		return nil, NewValidationError("product variant ID cannot be empty")
	}
	if effectiveTo != nil && effectiveTo.Before(effectiveFrom) {
		return nil, NewValidationError("effectiveTo cannot be before effectiveFrom")
	}

	return &ClientPricing{
		ID:               uuid.New(),
		ClientID:         clientID,
		ProductVariantID: productVariantID,
		FixedPrice:       fixedPrice,
		EffectiveFrom:    effectiveFrom,
		EffectiveTo:      effectiveTo,
		IsActive:         true,
	}, nil
}

func (c *ClientPricing) AppliesTo(clientID uuid.UUID, variantID uuid.UUID, at time.Time) bool {
	if !c.IsActive {
		return false
	}
	if c.ClientID != clientID || c.ProductVariantID != variantID {
		return false
	}
	if at.Before(c.EffectiveFrom) {
		return false
	}
	if c.EffectiveTo != nil && at.After(*c.EffectiveTo) {
		return false
	}
	return true
}
