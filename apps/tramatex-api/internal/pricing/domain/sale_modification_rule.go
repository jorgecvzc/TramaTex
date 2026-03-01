package domain

import (
	"time"

	"github.com/google/uuid"
)

type SaleModificationRule struct {
	ID             uuid.UUID
	Name           string
	ClientIDs      []string
	ProductGroupID *uuid.UUID
	MinOrderTotal  *Money
	Value          RuleValue
	Priority       int
	EffectiveFrom  time.Time
	EffectiveTo    *time.Time
	IsActive       bool
}

func NewSaleModificationRule(
	name string,
	clientIDs []string,
	productGroupID *uuid.UUID,
	minOrderTotal *Money,
	value RuleValue,
	priority int,
	effectiveFrom time.Time,
	effectiveTo *time.Time,
) (*SaleModificationRule, error) {
	if name == "" {
		return nil, NewValidationError("rule name cannot be empty")
	}
	if priority < 0 {
		return nil, NewValidationError("priority cannot be negative")
	}
	if effectiveTo != nil && effectiveTo.Before(effectiveFrom) {
		return nil, NewValidationError("effectiveTo cannot be before effectiveFrom")
	}

	return &SaleModificationRule{
		ID:             uuid.New(),
		Name:           name,
		ClientIDs:      clientIDs,
		ProductGroupID: productGroupID,
		MinOrderTotal:  minOrderTotal,
		Value:          value,
		Priority:       priority,
		EffectiveFrom:  effectiveFrom,
		EffectiveTo:    effectiveTo,
		IsActive:       true,
	}, nil
}

func (r *SaleModificationRule) AppliesTo(clientID string, productGroupID *uuid.UUID, orderTotal Money, at time.Time) bool {
	if !r.IsActive {
		return false
	}
	if at.Before(r.EffectiveFrom) {
		return false
	}
	if r.EffectiveTo != nil && at.After(*r.EffectiveTo) {
		return false
	}
	if r.ProductGroupID != nil {
		if productGroupID == nil || *r.ProductGroupID != *productGroupID {
			return false
		}
	}
	if len(r.ClientIDs) > 0 {
		matched := false
		for _, id := range r.ClientIDs {
			if id == clientID {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if r.MinOrderTotal != nil {
		if orderTotal.Currency() != r.MinOrderTotal.Currency() {
			return false
		}
		if orderTotal.Amount() < r.MinOrderTotal.Amount() {
			return false
		}
	}
	return true
}
