package application

import (
	"time"

	"github.com/google/uuid"
)

type CreateClientPricingCommand struct {
	ClientID         uuid.UUID  `json:"client_id"`
	ProductVariantID uuid.UUID  `json:"product_variant_id"`
	FixedPrice       float64    `json:"fixed_price"`
	Currency         string     `json:"currency"`
	EffectiveFrom    time.Time  `json:"effective_from"`
	EffectiveTo      *time.Time `json:"effective_to"`
}
