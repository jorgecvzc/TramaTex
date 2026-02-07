package persistence

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type PriceCalculationDataModel struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;primary_key;"`
	ProductVariantID uuid.UUID `gorm:"type:uuid;not null"`
	ClientID         uuid.UUID `gorm:"type:uuid;not null"`
	Quantity         int       `gorm:"not null"`
	BaseCost         float64   `gorm:"type:numeric(12,2);not null"`
	FinalPrice       float64   `gorm:"type:numeric(12,2);not null"`
	Currency         string    `gorm:"type:varchar(3);not null"`
	AppliedRules     string    `gorm:"type:jsonb;not null"`
	CalculatedAt     time.Time `gorm:"not null"`
}

func (PriceCalculationDataModel) TableName() string {
	return "price_calculations"
}

func (m *PriceCalculationDataModel) ToDomain() (*domain.PriceCalculation, error) {
	baseCost, err := domain.NewMoney(m.BaseCost, m.Currency)
	if err != nil {
		return nil, err
	}
	finalPrice, err := domain.NewMoney(m.FinalPrice, m.Currency)
	if err != nil {
		return nil, err
	}

	var rules []string
	if err := json.Unmarshal([]byte(m.AppliedRules), &rules); err != nil {
		return nil, err
	}

	return &domain.PriceCalculation{
		ID:               m.ID,
		ProductVariantID: m.ProductVariantID,
		ClientID:         m.ClientID,
		Quantity:         m.Quantity,
		BaseCost:         baseCost,
		FinalPrice:       finalPrice,
		AppliedRules:     rules,
		CalculatedAt:     m.CalculatedAt,
	}, nil
}

func PriceCalculationFromDomain(calc *domain.PriceCalculation) (*PriceCalculationDataModel, error) {
	rulesPayload, err := json.Marshal(calc.AppliedRules)
	if err != nil {
		return nil, err
	}

	return &PriceCalculationDataModel{
		ID:               calc.ID,
		ProductVariantID: calc.ProductVariantID,
		ClientID:         calc.ClientID,
		Quantity:         calc.Quantity,
		BaseCost:         calc.BaseCost.Amount(),
		FinalPrice:       calc.FinalPrice.Amount(),
		Currency:         calc.FinalPrice.Currency(),
		AppliedRules:     string(rulesPayload),
		CalculatedAt:     calc.CalculatedAt,
	}, nil
}
