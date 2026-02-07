package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PricingRuleRepository interface {
	Save(ctx context.Context, rule *PricingRule) error
	FindByID(ctx context.Context, id uuid.UUID) (*PricingRule, error)
	List(ctx context.Context) ([]*PricingRule, error)
	FindApplicable(ctx context.Context, variantID uuid.UUID, quantity int, at time.Time) ([]*PricingRule, error)
}

type ClientPricingRepository interface {
	Save(ctx context.Context, override *ClientPricing) error
	FindApplicable(ctx context.Context, clientID uuid.UUID, variantID uuid.UUID, at time.Time) (*ClientPricing, error)
}

type BrandProfitMarginRepository interface {
	Save(ctx context.Context, margin *BrandProfitMargin) error
	FindApplicable(ctx context.Context, brandID uuid.UUID, at time.Time) (*BrandProfitMargin, error)
}

type SalesDiscountRuleRepository interface {
	Save(ctx context.Context, rule *SalesDiscountRule) error
	FindApplicable(ctx context.Context, clientID uuid.UUID, variantID uuid.UUID, quantity int, at time.Time) ([]*SalesDiscountRule, error)
}

type PriceCalculationRepository interface {
	Save(ctx context.Context, calc *PriceCalculation) error
	ListByProductVariantID(ctx context.Context, variantID uuid.UUID) ([]*PriceCalculation, error)
}

type BaseSalesPriceRuleRepository interface {
	Save(ctx context.Context, rule *BaseSalesPriceRule) error
	FindByID(ctx context.Context, id uuid.UUID) (*BaseSalesPriceRule, error)
	List(ctx context.Context) ([]*BaseSalesPriceRule, error)
}

type SaleModificationRuleRepository interface {
	Save(ctx context.Context, rule *SaleModificationRule) error
	FindByID(ctx context.Context, id uuid.UUID) (*SaleModificationRule, error)
	ListApplicable(ctx context.Context, clientID uuid.UUID, productGroupID *uuid.UUID, orderTotal Money, at time.Time) ([]*SaleModificationRule, error)
}
