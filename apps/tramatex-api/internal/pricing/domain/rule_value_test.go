package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuleValueValidation(t *testing.T) {
	p, _ := NewPercentage(0.1)
	m, _ := NewMoney(5, DefaultCurrency)

	_, err := NewRuleValue(RuleValuePercentageMarkup, nil, nil)
	require.Error(t, err)

	_, err = NewRuleValue(RuleValueFixedAmountIncrease, nil, nil)
	require.Error(t, err)

	_, err = NewRuleValue(RuleValueType("UNKNOWN"), &p, &m)
	require.Error(t, err)
}

func TestRuleValueApply(t *testing.T) {
	base, _ := NewMoney(100, DefaultCurrency)
	p, _ := NewPercentage(0.1)
	m5, _ := NewMoney(5, DefaultCurrency)
	m50, _ := NewMoney(50, DefaultCurrency)
	m40, _ := NewMoney(40, DefaultCurrency)

	value, _ := NewRuleValue(RuleValuePercentageMarkup, &p, nil)
	result, err := value.Apply(base)
	require.NoError(t, err)
	require.InDelta(t, 110.0, result.Amount(), 0.0001)

	value, _ = NewRuleValue(RuleValueFixedAmountIncrease, nil, &m5)
	result, err = value.Apply(base)
	require.NoError(t, err)
	require.InDelta(t, 105.0, result.Amount(), 0.0001)

	value, _ = NewRuleValue(RuleValueSetToFixedPrice, nil, &m50)
	result, err = value.Apply(base)
	require.NoError(t, err)
	require.InDelta(t, 50.0, result.Amount(), 0.0001)

	value, _ = NewRuleValue(RuleValueApplyPercentageDiscount, &p, nil)
	result, err = value.Apply(base)
	require.NoError(t, err)
	require.InDelta(t, 90.0, result.Amount(), 0.0001)

	value, _ = NewRuleValue(RuleValueApplyFixedAmountDiscount, nil, &m5)
	result, err = value.Apply(base)
	require.NoError(t, err)
	require.InDelta(t, 95.0, result.Amount(), 0.0001)

	value, _ = NewRuleValue(RuleValueSetToFixedDiscountedPrice, nil, &m40)
	result, err = value.Apply(base)
	require.NoError(t, err)
	require.InDelta(t, 40.0, result.Amount(), 0.0001)
}

func TestRuleValueApplyInvalidType(t *testing.T) {
	base, _ := NewMoney(100, DefaultCurrency)
	value := RuleValue{Type: RuleValueType("UNKNOWN")}
	_, err := value.Apply(base)
	require.Error(t, err)
}
