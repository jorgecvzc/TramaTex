package application

import (
	"time"

	"github.com/google/uuid"
)

type CreatePricingRuleCommand struct {
	Name             string     `json:"name"`
	ProductVariantID *uuid.UUID `json:"product_variant_id"`
	PartyCategory    *string    `json:"party_category"`
	MarkupPercentage float64    `json:"margin_percentage"`
	MinQuantity      int        `json:"min_quantity"`
	MaxQuantity      *int       `json:"max_quantity"`
	EffectiveFrom    time.Time  `json:"effective_from"`
	EffectiveTo      *time.Time `json:"effective_to"`
}

type CalculatePriceCommand struct {
	ProductVariantID uuid.UUID `json:"product_variant_id"`
	ClientID         uuid.UUID `json:"client_id"`
	Quantity         int       `json:"quantity"`
}

type CreateClientPricingCommand struct {
	ClientID         uuid.UUID  `json:"client_id"`
	ProductVariantID uuid.UUID  `json:"product_variant_id"`
	FixedPrice       float64    `json:"fixed_price"`
	Currency         string     `json:"currency"`
	EffectiveFrom    time.Time  `json:"effective_from"`
	EffectiveTo      *time.Time `json:"effective_to"`
}
