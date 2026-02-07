package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSellingPriceCalculatorService(t *testing.T) {
	service := NewSellingPriceCalculatorService()
	base, _ := NewMoney(100, DefaultCurrency)
	p, _ := NewPercentage(0.1)
	fixed, _ := NewMoney(5, DefaultCurrency)

	price, err := service.CalculateSellingPrice(base, p, nil)
	require.NoError(t, err)
	require.InDelta(t, 110.0, price.Amount(), 0.0001)

	price, err = service.CalculateSellingPrice(base, p, &fixed)
	require.NoError(t, err)
	require.InDelta(t, 115.0, price.Amount(), 0.0001)
}

func TestSalesDiscountCalculatorService(t *testing.T) {
	service := NewSalesDiscountCalculatorService()
	base, _ := NewMoney(100, DefaultCurrency)
	p, _ := NewPercentage(0.1)
	fixed, _ := NewMoney(5, DefaultCurrency)

	rules := []*SalesDiscountRule{
		{DiscountType: DiscountTypeFixed, FixedAmount: &fixed, Priority: 1, Name: "fixed"},
		{DiscountType: DiscountTypePercentage, Percentage: &p, Priority: 2, Name: "percent"},
	}

	final, discounts, err := service.ApplyDiscounts(base, rules)
	require.NoError(t, err)
	require.Len(t, discounts, 2)
	require.Equal(t, "percent", discounts[0].Name)
	require.InDelta(t, 85.0, final.Amount(), 0.0001)
}

func TestSalesDiscountCalculatorServiceInvalidRule(t *testing.T) {
	service := NewSalesDiscountCalculatorService()
	base, _ := NewMoney(100, DefaultCurrency)
	rules := []*SalesDiscountRule{{DiscountType: DiscountType("UNKNOWN")}}

	_, _, err := service.ApplyDiscounts(base, rules)
	require.Error(t, err)
}

func TestSalesDiscountCalculatorServiceNoRules(t *testing.T) {
	service := NewSalesDiscountCalculatorService()
	base, _ := NewMoney(100, DefaultCurrency)

	final, discounts, err := service.ApplyDiscounts(base, nil)
	require.NoError(t, err)
	require.Empty(t, discounts)
	require.InDelta(t, 100.0, final.Amount(), 0.0001)
}
