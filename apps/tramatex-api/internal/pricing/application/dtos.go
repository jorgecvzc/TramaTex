package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type MoneyDTO struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type PercentageDTO struct {
	Value float64 `json:"value"`
}

type ClientPricingDTO struct {
	ID               uuid.UUID  `json:"id"`
	ClientID         uuid.UUID  `json:"clientId"`
	ProductVariantID uuid.UUID  `json:"productVariantId"`
	FixedPrice       MoneyDTO   `json:"fixedPrice"`
	EffectiveFrom    time.Time  `json:"effectiveFrom"`
	EffectiveTo      *time.Time `json:"effectiveTo,omitempty"`
	IsActive         bool       `json:"isActive"`
}

type PriceCalculationDTO struct {
	ID               uuid.UUID `json:"id"`
	ProductVariantID uuid.UUID `json:"productVariantId"`
	ClientID         uuid.UUID `json:"clientId"`
	Quantity         int       `json:"quantity"`
	BaseCost         MoneyDTO  `json:"baseCost"`
	FinalPrice       MoneyDTO  `json:"finalPrice"`
	AppliedRules     []string  `json:"appliedRules"`
	CalculatedAt     time.Time `json:"calculatedAt"`
}

func NewMoneyDTO(m domain.Money) MoneyDTO {
	return MoneyDTO{Amount: m.Amount(), Currency: m.Currency()}
}

func NewClientPricingDTO(override *domain.ClientPricing) ClientPricingDTO {
	return ClientPricingDTO{
		ID:               override.ID,
		ClientID:         override.ClientID,
		ProductVariantID: override.ProductVariantID,
		FixedPrice:       NewMoneyDTO(override.FixedPrice),
		EffectiveFrom:    override.EffectiveFrom,
		EffectiveTo:      override.EffectiveTo,
		IsActive:         override.IsActive,
	}
}

func NewPriceCalculationDTO(calc *domain.PriceCalculation) PriceCalculationDTO {
	return PriceCalculationDTO{
		ID:               calc.ID,
		ProductVariantID: calc.ProductVariantID,
		ClientID:         calc.ClientID,
		Quantity:         calc.Quantity,
		BaseCost:         NewMoneyDTO(calc.BaseCost),
		FinalPrice:       NewMoneyDTO(calc.FinalPrice),
		AppliedRules:     calc.AppliedRules,
		CalculatedAt:     calc.CalculatedAt,
	}
}
