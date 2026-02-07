package application

import "github.com/google/uuid"

type ListPricingRulesQuery struct{}

type GetPricingHistoryQuery struct {
	ProductVariantID uuid.UUID
}
