package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBaseSalesPriceRuleValidation(t *testing.T) {
	p, _ := NewPercentage(0.1)
	value, _ := NewRuleValue(RuleValuePercentageMarkup, &p, nil)

	_, err := NewBaseSalesPriceRule("", nil, nil, nil, nil, value)
	require.Error(t, err)
}

func TestBaseSalesPriceRuleSuccess(t *testing.T) {
	brandID := uuid.New()
	p, _ := NewPercentage(0.1)
	value, _ := NewRuleValue(RuleValuePercentageMarkup, &p, nil)

	rule, err := NewBaseSalesPriceRule("rule", &brandID, nil, nil, nil, value)
	require.NoError(t, err)
	require.Equal(t, brandID, *rule.BrandID)
}
