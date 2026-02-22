package persistence

import (
	"time"

	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type PricingRuleDataModel struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;"`
	Name             string     `gorm:"not null"`
	ProductVariantID *uuid.UUID `gorm:"type:uuid"`
	PartyCategory    *string    `gorm:"type:varchar(50)"`
	MarkupPercentage float64    `gorm:"type:numeric(8,4);not null"`
	MinQuantity      int        `gorm:"not null;default:0"`
	MaxQuantity      *int       `gorm:""`
	EffectiveFrom    time.Time  `gorm:"not null"`
	EffectiveTo      *time.Time `gorm:""`
	IsActive         bool       `gorm:"not null;default:true"`
}

func (PricingRuleDataModel) TableName() string {
	return "pricing_rules"
}

func (m *PricingRuleDataModel) ToDomain() (*domain.PricingRule, error) {
	percentage, err := domain.NewPercentage(m.MarkupPercentage)
	if err != nil {
		return nil, err
	}

	return &domain.PricingRule{
		ID:               m.ID,
		Name:             m.Name,
		ProductVariantID: m.ProductVariantID,
		PartyCategory:    m.PartyCategory,
		Markup:           percentage,
		MinQuantity:      m.MinQuantity,
		MaxQuantity:      m.MaxQuantity,
		EffectiveFrom:    m.EffectiveFrom,
		EffectiveTo:      m.EffectiveTo,
		IsActive:         m.IsActive,
	}, nil
}

func PricingRuleFromDomain(rule *domain.PricingRule) *PricingRuleDataModel {
	return &PricingRuleDataModel{
		ID:               rule.ID,
		Name:             rule.Name,
		ProductVariantID: rule.ProductVariantID,
		PartyCategory:    rule.PartyCategory,
		MarkupPercentage: rule.Markup.Value(),
		MinQuantity:      rule.MinQuantity,
		MaxQuantity:      rule.MaxQuantity,
		EffectiveFrom:    rule.EffectiveFrom,
		EffectiveTo:      rule.EffectiveTo,
		IsActive:         rule.IsActive,
	}
}
