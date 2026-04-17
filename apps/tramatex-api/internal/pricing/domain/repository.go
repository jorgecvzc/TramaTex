package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ClientPricingRepository interface {
	Save(ctx context.Context, override *ClientPricing) error
	FindApplicable(ctx context.Context, clientID uuid.UUID, variantID uuid.UUID, at time.Time) (*ClientPricing, error)
	FindApplicableBulk(ctx context.Context, clientID uuid.UUID, variantIDs []uuid.UUID, at time.Time) (map[uuid.UUID]*ClientPricing, error)
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
	ListActive(ctx context.Context, at time.Time) ([]*SaleModificationRule, error)
	ListApplicable(ctx context.Context, clientID string, productGroupID *uuid.UUID, orderTotal Money, at time.Time) ([]*SaleModificationRule, error)
}
