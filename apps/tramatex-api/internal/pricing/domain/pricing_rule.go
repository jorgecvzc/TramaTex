package domain

import (
	"time"

	"github.com/google/uuid"
)

type PricingRule struct {
	ID               uuid.UUID
	Name             string
	ProductVariantID *uuid.UUID
	PartyCategory    *string
	Markup           Percentage
	MinQuantity      int
	MaxQuantity      *int
	EffectiveFrom    time.Time
	EffectiveTo      *time.Time
	IsActive         bool
}

func NewPricingRule(
	name string,
	productVariantID *uuid.UUID,
	partyCategory *string,
	markup Percentage,
	minQuantity int,
	maxQuantity *int,
	effectiveFrom time.Time,
	effectiveTo *time.Time,
) (*PricingRule, error) {
	if name == "" {
		return nil, NewValidationError("rule name cannot be empty")
	}
	if minQuantity < 0 {
		return nil, NewValidationError("min quantity cannot be negative")
	}
	if maxQuantity != nil && *maxQuantity < minQuantity {
		return nil, NewValidationError("max quantity cannot be less than min quantity")
	}
	if effectiveTo != nil && effectiveTo.Before(effectiveFrom) {
		return nil, NewValidationError("effectiveTo cannot be before effectiveFrom")
	}

	return &PricingRule{
		ID:               uuid.New(),
		Name:             name,
		ProductVariantID: productVariantID,
		PartyCategory:    partyCategory,
		Markup:           markup,
		MinQuantity:      minQuantity,
		MaxQuantity:      maxQuantity,
		EffectiveFrom:    effectiveFrom,
		EffectiveTo:      effectiveTo,
		IsActive:         true,
	}, nil
}

func (r *PricingRule) AppliesTo(variantID uuid.UUID, quantity int, at time.Time) bool {
	if !r.IsActive {
		return false
	}
	if r.ProductVariantID != nil && *r.ProductVariantID != variantID {
		return false
	}
	if quantity < r.MinQuantity {
		return false
	}
	if r.MaxQuantity != nil && quantity > *r.MaxQuantity {
		return false
	}
	if at.Before(r.EffectiveFrom) {
		return false
	}
	if r.EffectiveTo != nil && at.After(*r.EffectiveTo) {
		return false
	}
	return true
}
