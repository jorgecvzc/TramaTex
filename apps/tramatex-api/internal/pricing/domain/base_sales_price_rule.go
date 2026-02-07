package domain

import "github.com/google/uuid"

type BaseSalesPriceRule struct {
	ID             uuid.UUID
	Name           string
	BrandID        *uuid.UUID
	ProductGroupID *uuid.UUID
	ProductID      *uuid.UUID
	VariantID      *uuid.UUID
	Value          RuleValue
	IsActive       bool
}

func NewBaseSalesPriceRule(
	name string,
	brandID *uuid.UUID,
	productGroupID *uuid.UUID,
	productID *uuid.UUID,
	variantID *uuid.UUID,
	value RuleValue,
) (*BaseSalesPriceRule, error) {
	if name == "" {
		return nil, NewValidationError("rule name cannot be empty")
	}

	return &BaseSalesPriceRule{
		ID:             uuid.New(),
		Name:           name,
		BrandID:        brandID,
		ProductGroupID: productGroupID,
		ProductID:      productID,
		VariantID:      variantID,
		Value:          value,
		IsActive:       true,
	}, nil
}
