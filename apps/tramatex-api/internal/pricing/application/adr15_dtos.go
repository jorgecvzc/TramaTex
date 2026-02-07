package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type RuleValueDTO struct {
	Type            string         `json:"type"`
	PercentageValue *PercentageDTO `json:"percentageValue,omitempty"`
	MoneyValue      *MoneyDTO      `json:"moneyValue,omitempty"`
}

type BaseSalesPriceRuleDTO struct {
	ID             uuid.UUID    `json:"id"`
	Name           string       `json:"name"`
	BrandID        *uuid.UUID   `json:"brandId,omitempty"`
	ProductGroupID *uuid.UUID   `json:"productGroupId,omitempty"`
	ProductID      *uuid.UUID   `json:"productId,omitempty"`
	VariantID      *uuid.UUID   `json:"variantId,omitempty"`
	Value          RuleValueDTO `json:"value"`
	IsActive       bool         `json:"isActive"`
}

type SaleModificationRuleDTO struct {
	ID                  uuid.UUID    `json:"id"`
	Name                string       `json:"name"`
	ClientIDs           []uuid.UUID  `json:"clientIds,omitempty"`
	ProductGroupID      *uuid.UUID   `json:"productGroupId,omitempty"`
	MinOrderTotalAmount *MoneyDTO    `json:"minOrderTotalAmount,omitempty"`
	Value               RuleValueDTO `json:"value"`
	Priority            int          `json:"priority"`
	IsActive            bool         `json:"isActive"`
	EffectiveFrom       time.Time    `json:"effectiveFrom"`
	EffectiveTo         *time.Time   `json:"effectiveTo,omitempty"`
}

type CalculateBaseSalesPriceRequest struct {
	ProductID uuid.UUID `json:"productId"`
	VariantID uuid.UUID `json:"variantId"`
}

type CalculatedBaseSalesPriceResponse struct {
	VariantID      uuid.UUID `json:"variantId"`
	BaseSalesPrice MoneyDTO  `json:"baseSalesPrice"`
}

type SaleItemRequest struct {
	ProductVariantID uuid.UUID `json:"productVariantId"`
	Quantity         int       `json:"quantity"`
}

type CalculateFinalSalePriceRequest struct {
	SaleItems []SaleItemRequest `json:"saleItems"`
	ClientID  uuid.UUID         `json:"clientId"`
	SaleDate  time.Time         `json:"saleDate"`
}

type CalculatedSaleItemResponse struct {
	ProductVariantID uuid.UUID `json:"productVariantId"`
	Quantity         int       `json:"quantity"`
	BaseSalesPrice   MoneyDTO  `json:"baseSalesPrice"`
	FinalPrice       MoneyDTO  `json:"finalPrice"`
}

type CalculateFinalSalePriceResponse struct {
	CalculatedItems []CalculatedSaleItemResponse `json:"calculatedItems"`
	SaleTotal       MoneyDTO                     `json:"saleTotal"`
}

func NewRuleValueDTO(value domain.RuleValue) RuleValueDTO {
	var percentage *PercentageDTO
	if value.PercentageValue != nil {
		percentage = &PercentageDTO{Value: value.PercentageValue.Value()}
	}
	var money *MoneyDTO
	if value.MoneyValue != nil {
		m := NewMoneyDTO(*value.MoneyValue)
		money = &m
	}

	return RuleValueDTO{
		Type:            string(value.Type),
		PercentageValue: percentage,
		MoneyValue:      money,
	}
}

func NewBaseSalesPriceRuleDTO(rule *domain.BaseSalesPriceRule) BaseSalesPriceRuleDTO {
	return BaseSalesPriceRuleDTO{
		ID:             rule.ID,
		Name:           rule.Name,
		BrandID:        rule.BrandID,
		ProductGroupID: rule.ProductGroupID,
		ProductID:      rule.ProductID,
		VariantID:      rule.VariantID,
		Value:          NewRuleValueDTO(rule.Value),
		IsActive:       rule.IsActive,
	}
}

func NewSaleModificationRuleDTO(rule *domain.SaleModificationRule) SaleModificationRuleDTO {
	var minOrder *MoneyDTO
	if rule.MinOrderTotal != nil {
		dto := NewMoneyDTO(*rule.MinOrderTotal)
		minOrder = &dto
	}

	return SaleModificationRuleDTO{
		ID:                  rule.ID,
		Name:                rule.Name,
		ClientIDs:           rule.ClientIDs,
		ProductGroupID:      rule.ProductGroupID,
		MinOrderTotalAmount: minOrder,
		Value:               NewRuleValueDTO(rule.Value),
		Priority:            rule.Priority,
		IsActive:            rule.IsActive,
		EffectiveFrom:       rule.EffectiveFrom,
		EffectiveTo:         rule.EffectiveTo,
	}
}
