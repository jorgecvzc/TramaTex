package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type SalesDiscountRuleDataModel struct {
	gorm.Model
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;"`
	Name             string     `gorm:"not null"`
	ClientID         *uuid.UUID `gorm:"type:uuid"`
	ProductVariantID *uuid.UUID `gorm:"type:uuid"`
	MinQuantity      *int       `gorm:""`
	DiscountType     string     `gorm:"type:varchar(20);not null"`
	PercentageValue  *float64   `gorm:"type:numeric(8,4)"`
	FixedAmount      *float64   `gorm:"type:numeric(12,2)"`
	Currency         string     `gorm:"type:varchar(3);not null"`
	Priority         int        `gorm:"not null;default:0"`
	EffectiveFrom    time.Time  `gorm:"not null"`
	EffectiveTo      *time.Time
	IsActive         bool `gorm:"not null;default:true"`
}

func (SalesDiscountRuleDataModel) TableName() string {
	return "sales_discount_rules"
}

func (m *SalesDiscountRuleDataModel) ToDomain() (*domain.SalesDiscountRule, error) {
	var percentage *domain.Percentage
	if m.PercentageValue != nil {
		p, err := domain.NewPercentage(*m.PercentageValue)
		if err != nil {
			return nil, err
		}
		percentage = &p
	}

	var fixed *domain.Money
	if m.FixedAmount != nil {
		money, err := domain.NewMoney(*m.FixedAmount, m.Currency)
		if err != nil {
			return nil, err
		}
		fixed = &money
	}

	return &domain.SalesDiscountRule{
		ID:               m.ID,
		Name:             m.Name,
		ClientID:         m.ClientID,
		ProductVariantID: m.ProductVariantID,
		MinQuantity:      m.MinQuantity,
		DiscountType:     domain.DiscountType(m.DiscountType),
		Percentage:       percentage,
		FixedAmount:      fixed,
		Priority:         m.Priority,
		EffectiveFrom:    m.EffectiveFrom,
		EffectiveTo:      m.EffectiveTo,
		IsActive:         m.IsActive,
	}, nil
}

func SalesDiscountRuleFromDomain(rule *domain.SalesDiscountRule) *SalesDiscountRuleDataModel {
	var percentageValue *float64
	if rule.Percentage != nil {
		value := rule.Percentage.Value()
		percentageValue = &value
	}

	var fixedAmount *float64
	currency := domain.DefaultCurrency
	if rule.FixedAmount != nil {
		amount := rule.FixedAmount.Amount()
		fixedAmount = &amount
		currency = rule.FixedAmount.Currency()
	}

	return &SalesDiscountRuleDataModel{
		ID:               rule.ID,
		Name:             rule.Name,
		ClientID:         rule.ClientID,
		ProductVariantID: rule.ProductVariantID,
		MinQuantity:      rule.MinQuantity,
		DiscountType:     string(rule.DiscountType),
		PercentageValue:  percentageValue,
		FixedAmount:      fixedAmount,
		Currency:         currency,
		Priority:         rule.Priority,
		EffectiveFrom:    rule.EffectiveFrom,
		EffectiveTo:      rule.EffectiveTo,
		IsActive:         rule.IsActive,
	}
}
