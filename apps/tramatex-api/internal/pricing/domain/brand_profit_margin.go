package domain

import (
	"time"

	"github.com/google/uuid"
)

type BrandProfitMargin struct {
	ID            uuid.UUID
	BrandID       uuid.UUID
	Percentage    *Percentage
	FixedAmount   *Money
	EffectiveFrom time.Time
	EffectiveTo   *time.Time
	IsActive      bool
}

func NewBrandProfitMargin(
	brandID uuid.UUID,
	percentage *Percentage,
	fixedAmount *Money,
	effectiveFrom time.Time,
	effectiveTo *time.Time,
) (*BrandProfitMargin, error) {
	if brandID == uuid.Nil {
		return nil, NewValidationError("brand ID cannot be empty")
	}
	if percentage == nil && fixedAmount == nil {
		return nil, NewValidationError("either percentage or fixed amount must be provided")
	}
	if effectiveTo != nil && effectiveTo.Before(effectiveFrom) {
		return nil, NewValidationError("effectiveTo cannot be before effectiveFrom")
	}

	return &BrandProfitMargin{
		ID:            uuid.New(),
		BrandID:       brandID,
		Percentage:    percentage,
		FixedAmount:   fixedAmount,
		EffectiveFrom: effectiveFrom,
		EffectiveTo:   effectiveTo,
		IsActive:      true,
	}, nil
}

func (m *BrandProfitMargin) AppliesTo(brandID uuid.UUID, at time.Time) bool {
	if !m.IsActive {
		return false
	}
	if m.BrandID != brandID {
		return false
	}
	if at.Before(m.EffectiveFrom) {
		return false
	}
	if m.EffectiveTo != nil && at.After(*m.EffectiveTo) {
		return false
	}
	return true
}
