package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type ClientPricingDataModel struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;primary_key;"`
	ClientID         uuid.UUID `gorm:"type:uuid;not null"`
	ProductVariantID uuid.UUID `gorm:"type:uuid;not null"`
	FixedPrice       float64   `gorm:"type:numeric(12,2);not null"`
	Currency         string    `gorm:"type:varchar(3);not null"`
	EffectiveFrom    time.Time `gorm:"not null"`
	EffectiveTo      *time.Time
	IsActive         bool `gorm:"not null;default:true"`
}

func (ClientPricingDataModel) TableName() string {
	return "client_pricing_overrides"
}

func (m *ClientPricingDataModel) ToDomain() (*domain.ClientPricing, error) {
	price, err := domain.NewMoney(m.FixedPrice, m.Currency)
	if err != nil {
		return nil, err
	}

	return &domain.ClientPricing{
		ID:               m.ID,
		ClientID:         m.ClientID,
		ProductVariantID: m.ProductVariantID,
		FixedPrice:       price,
		EffectiveFrom:    m.EffectiveFrom,
		EffectiveTo:      m.EffectiveTo,
		IsActive:         m.IsActive,
	}, nil
}

func ClientPricingFromDomain(override *domain.ClientPricing) *ClientPricingDataModel {
	return &ClientPricingDataModel{
		ID:               override.ID,
		ClientID:         override.ClientID,
		ProductVariantID: override.ProductVariantID,
		FixedPrice:       override.FixedPrice.Amount(),
		Currency:         override.FixedPrice.Currency(),
		EffectiveFrom:    override.EffectiveFrom,
		EffectiveTo:      override.EffectiveTo,
		IsActive:         override.IsActive,
	}
}
