package persistence

import (
	"time"

	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type BrandProfitMarginDataModel struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;"`
	BrandID         uuid.UUID `gorm:"type:uuid;not null"`
	PercentageValue *float64  `gorm:"type:numeric(8,4)"`
	FixedAmount     *float64  `gorm:"type:numeric(12,2)"`
	Currency        string    `gorm:"type:varchar(3);not null"`
	EffectiveFrom   time.Time `gorm:"not null"`
	EffectiveTo     *time.Time
	IsActive        bool `gorm:"not null;default:true"`
}

func (BrandProfitMarginDataModel) TableName() string {
	return "brand_profit_margins"
}

func (m *BrandProfitMarginDataModel) ToDomain() (*domain.BrandProfitMargin, error) {
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

	return &domain.BrandProfitMargin{
		ID:            m.ID,
		BrandID:       m.BrandID,
		Percentage:    percentage,
		FixedAmount:   fixed,
		EffectiveFrom: m.EffectiveFrom,
		EffectiveTo:   m.EffectiveTo,
		IsActive:      m.IsActive,
	}, nil
}

func BrandProfitMarginFromDomain(margin *domain.BrandProfitMargin) *BrandProfitMarginDataModel {
	var percentageValue *float64
	if margin.Percentage != nil {
		value := margin.Percentage.Value()
		percentageValue = &value
	}

	var fixedAmount *float64
	currency := domain.DefaultCurrency
	if margin.FixedAmount != nil {
		amount := margin.FixedAmount.Amount()
		fixedAmount = &amount
		currency = margin.FixedAmount.Currency()
	}

	return &BrandProfitMarginDataModel{
		ID:              margin.ID,
		BrandID:         margin.BrandID,
		PercentageValue: percentageValue,
		FixedAmount:     fixedAmount,
		Currency:        currency,
		EffectiveFrom:   margin.EffectiveFrom,
		EffectiveTo:     margin.EffectiveTo,
		IsActive:        margin.IsActive,
	}
}
