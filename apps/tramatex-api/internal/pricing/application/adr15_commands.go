package application

import (
	"time"

	"github.com/google/uuid"
)

type CreateBaseSalesPriceRuleCommand struct {
	Name           string       `json:"name"`
	BrandID        *uuid.UUID   `json:"brandId"`
	ProductGroupID *uuid.UUID   `json:"productGroupId"`
	ProductID      *uuid.UUID   `json:"productId"`
	VariantID      *uuid.UUID   `json:"variantId"`
	Value          RuleValueDTO `json:"value"`
	IsActive       *bool        `json:"isActive"`
}

type UpdateBaseSalesPriceRuleCommand struct {
	ID             uuid.UUID
	Name           *string       `json:"name"`
	BrandID        *uuid.UUID    `json:"brandId"`
	ProductGroupID *uuid.UUID    `json:"productGroupId"`
	ProductID      *uuid.UUID    `json:"productId"`
	VariantID      *uuid.UUID    `json:"variantId"`
	Value          *RuleValueDTO `json:"value"`
	IsActive       *bool         `json:"isActive"`
}

type CreateSaleModificationRuleCommand struct {
	Name                string       `json:"name"`
	ClientIDs           []uuid.UUID  `json:"clientIds"`
	ProductGroupID      *uuid.UUID   `json:"productGroupId"`
	MinOrderTotalAmount *MoneyDTO    `json:"minOrderTotalAmount"`
	Value               RuleValueDTO `json:"value"`
	Priority            int          `json:"priority"`
	IsActive            *bool        `json:"isActive"`
	EffectiveFrom       time.Time    `json:"effectiveFrom"`
	EffectiveTo         *time.Time   `json:"effectiveTo"`
}

type UpdateSaleModificationRuleCommand struct {
	ID                  uuid.UUID
	Name                *string       `json:"name"`
	ClientIDs           []uuid.UUID   `json:"clientIds"`
	ProductGroupID      *uuid.UUID    `json:"productGroupId"`
	MinOrderTotalAmount *MoneyDTO     `json:"minOrderTotalAmount"`
	Value               *RuleValueDTO `json:"value"`
	Priority            *int          `json:"priority"`
	IsActive            *bool         `json:"isActive"`
	EffectiveFrom       *time.Time    `json:"effectiveFrom"`
	EffectiveTo         *time.Time    `json:"effectiveTo"`
}
