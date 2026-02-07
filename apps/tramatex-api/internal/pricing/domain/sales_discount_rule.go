package domain

import (
	"time"

	"github.com/google/uuid"
)

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "PERCENTAGE"
	DiscountTypeFixed      DiscountType = "FIXED_AMOUNT"
)

type SalesDiscountRule struct {
	ID               uuid.UUID
	Name             string
	ClientID         *uuid.UUID
	ProductVariantID *uuid.UUID
	MinQuantity      *int
	DiscountType     DiscountType
	Percentage       *Percentage
	FixedAmount      *Money
	Priority         int
	EffectiveFrom    time.Time
	EffectiveTo      *time.Time
	IsActive         bool
}

func NewSalesDiscountRule(
	name string,
	clientID *uuid.UUID,
	productVariantID *uuid.UUID,
	minQuantity *int,
	discountType DiscountType,
	percentage *Percentage,
	fixedAmount *Money,
	priority int,
	effectiveFrom time.Time,
	effectiveTo *time.Time,
) (*SalesDiscountRule, error) {
	if name == "" {
		return nil, NewValidationError("discount name cannot be empty")
	}
	if priority < 0 {
		return nil, NewValidationError("priority cannot be negative")
	}
	if discountType == DiscountTypePercentage && percentage == nil {
		return nil, NewValidationError("percentage discount requires percentage value")
	}
	if discountType == DiscountTypeFixed && fixedAmount == nil {
		return nil, NewValidationError("fixed discount requires money value")
	}
	if effectiveTo != nil && effectiveTo.Before(effectiveFrom) {
		return nil, NewValidationError("effectiveTo cannot be before effectiveFrom")
	}

	return &SalesDiscountRule{
		ID:               uuid.New(),
		Name:             name,
		ClientID:         clientID,
		ProductVariantID: productVariantID,
		MinQuantity:      minQuantity,
		DiscountType:     discountType,
		Percentage:       percentage,
		FixedAmount:      fixedAmount,
		Priority:         priority,
		EffectiveFrom:    effectiveFrom,
		EffectiveTo:      effectiveTo,
		IsActive:         true,
	}, nil
}

func (r *SalesDiscountRule) AppliesTo(clientID uuid.UUID, variantID uuid.UUID, quantity int, at time.Time) bool {
	if !r.IsActive {
		return false
	}
	if r.ClientID != nil && *r.ClientID != clientID {
		return false
	}
	if r.ProductVariantID != nil && *r.ProductVariantID != variantID {
		return false
	}
	if r.MinQuantity != nil && quantity < *r.MinQuantity {
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
