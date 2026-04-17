package application

import "github.com/google/uuid"

type GetPricingHistoryQuery struct {
	ProductVariantID uuid.UUID
}
