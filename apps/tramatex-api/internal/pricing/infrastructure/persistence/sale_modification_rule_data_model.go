package persistence

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type SaleModificationRuleDataModel struct {
	gorm.Model
	ID                    uuid.UUID      `gorm:"type:uuid;primary_key;"`
	Name                  string         `gorm:"not null"`
	ClientIDs             pq.StringArray `gorm:"type:uuid[]"`
	ProductGroupID        *uuid.UUID     `gorm:"type:uuid"`
	MinOrderTotalAmount   *float64       `gorm:"type:numeric(12,2)"`
	MinOrderTotalCurrency string         `gorm:"type:varchar(3);not null"`
	ValueType             string         `gorm:"type:varchar(50);not null"`
	PercentageValue       *float64       `gorm:"type:numeric(8,4)"`
	MoneyValueAmount      *float64       `gorm:"type:numeric(12,2)"`
	MoneyValueCurrency    string         `gorm:"type:varchar(3);not null"`
	Priority              int            `gorm:"not null;default:0"`
	EffectiveFrom         time.Time      `gorm:"not null"`
	EffectiveTo           *time.Time
	IsActive              bool `gorm:"not null;default:true"`
}

func (SaleModificationRuleDataModel) TableName() string {
	return "sale_modification_rules"
}

func (m *SaleModificationRuleDataModel) ToDomain() (*domain.SaleModificationRule, error) {
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

	var minOrder *domain.Money
	if m.MinOrderTotalAmount != nil {
		orderValue, err := domain.NewMoney(*m.MinOrderTotalAmount, m.MinOrderTotalCurrency)
		if err != nil {
			return nil, err
		}
		minOrder = &orderValue
	}

	return &domain.SaleModificationRule{
		ID:             m.ID,
		Name:           m.Name,
		ClientIDs:      stringArrayToUUIDs(m.ClientIDs),
		ProductGroupID: m.ProductGroupID,
		MinOrderTotal:  minOrder,
		Value:          value,
		Priority:       m.Priority,
		EffectiveFrom:  m.EffectiveFrom,
		EffectiveTo:    m.EffectiveTo,
		IsActive:       m.IsActive,
	}, nil
}

func SaleModificationRuleFromDomain(rule *domain.SaleModificationRule) *SaleModificationRuleDataModel {
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

	var minOrderAmount *float64
	minOrderCurrency := domain.DefaultCurrency
	if rule.MinOrderTotal != nil {
		amount := rule.MinOrderTotal.Amount()
		minOrderAmount = &amount
		minOrderCurrency = rule.MinOrderTotal.Currency()
	}

	return &SaleModificationRuleDataModel{
		ID:                    rule.ID,
		Name:                  rule.Name,
		ClientIDs:             uuidArrayToStringArray(rule.ClientIDs),
		ProductGroupID:        rule.ProductGroupID,
		MinOrderTotalAmount:   minOrderAmount,
		MinOrderTotalCurrency: minOrderCurrency,
		ValueType:             string(rule.Value.Type),
		PercentageValue:       percentageValue,
		MoneyValueAmount:      moneyValueAmount,
		MoneyValueCurrency:    currency,
		Priority:              rule.Priority,
		EffectiveFrom:         rule.EffectiveFrom,
		EffectiveTo:           rule.EffectiveTo,
		IsActive:              rule.IsActive,
	}
}
