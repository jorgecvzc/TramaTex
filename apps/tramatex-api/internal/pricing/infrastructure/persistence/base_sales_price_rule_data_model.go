package persistence

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type BaseSalesPriceRuleDataModel struct {
	gorm.Model
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key;"`
	Name               string     `gorm:"not null"`
	BrandID            *uuid.UUID `gorm:"type:uuid"`
	ProductGroupID     *uuid.UUID `gorm:"type:uuid"`
	ProductID          *uuid.UUID `gorm:"type:uuid"`
	VariantID          *uuid.UUID `gorm:"type:uuid"`
	ValueType          string     `gorm:"type:varchar(50);not null"`
	PercentageValue    *float64   `gorm:"type:numeric(8,4)"`
	MoneyValueAmount   *float64   `gorm:"type:numeric(12,2)"`
	MoneyValueCurrency string     `gorm:"type:varchar(3);not null"`
	IsActive           bool       `gorm:"not null;default:true"`
}

func (BaseSalesPriceRuleDataModel) TableName() string {
	return "base_sales_price_rules"
}

func (m *BaseSalesPriceRuleDataModel) ToDomain() (*domain.BaseSalesPriceRule, error) {
	var percentage *domain.Percentage
	if m.PercentageValue != nil {
		p, err := domain.NewPercentage(*m.PercentageValue)
		if err != nil {
			return nil, err
		}
		percentage = &p
	}

	var money *domain.Money
	if m.MoneyValueAmount != nil {
		value, err := domain.NewMoney(*m.MoneyValueAmount, m.MoneyValueCurrency)
		if err != nil {
			return nil, err
		}
		money = &value
	}

	value, err := domain.NewRuleValue(domain.RuleValueType(m.ValueType), percentage, money)
	if err != nil {
		return nil, err
	}

	return &domain.BaseSalesPriceRule{
		ID:             m.ID,
		Name:           m.Name,
		BrandID:        m.BrandID,
		ProductGroupID: m.ProductGroupID,
		ProductID:      m.ProductID,
		VariantID:      m.VariantID,
		Value:          value,
		IsActive:       m.IsActive,
	}, nil
}

func BaseSalesPriceRuleFromDomain(rule *domain.BaseSalesPriceRule) *BaseSalesPriceRuleDataModel {
	var percentageValue *float64
	if rule.Value.PercentageValue != nil {
		value := rule.Value.PercentageValue.Value()
		percentageValue = &value
	}

	var moneyValueAmount *float64
	currency := domain.DefaultCurrency
	if rule.Value.MoneyValue != nil {
		amount := rule.Value.MoneyValue.Amount()
		moneyValueAmount = &amount
		currency = rule.Value.MoneyValue.Currency()
	}

	return &BaseSalesPriceRuleDataModel{
		ID:                 rule.ID,
		Name:               rule.Name,
		BrandID:            rule.BrandID,
		ProductGroupID:     rule.ProductGroupID,
		ProductID:          rule.ProductID,
		VariantID:          rule.VariantID,
		ValueType:          string(rule.Value.Type),
		PercentageValue:    percentageValue,
		MoneyValueAmount:   moneyValueAmount,
		MoneyValueCurrency: currency,
		IsActive:           rule.IsActive,
	}
}
