package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSalesDiscountRuleValidation(t *testing.T) {
	p, _ := NewPercentage(0.1)
	m, _ := NewMoney(5, DefaultCurrency)
	clientID := uuid.New()

	_, err := NewSalesDiscountRule("", nil, nil, nil, DiscountTypePercentage, &p, nil, 0, time.Now(), nil)
	require.Error(t, err)

	_, err = NewSalesDiscountRule("rule", nil, nil, nil, DiscountTypePercentage, nil, nil, -1, time.Now(), nil)
	require.Error(t, err)

	_, err = NewSalesDiscountRule("rule", nil, nil, nil, DiscountTypePercentage, nil, nil, 0, time.Now(), nil)
	require.Error(t, err)

	_, err = NewSalesDiscountRule("rule", nil, nil, nil, DiscountTypeFixed, nil, nil, 0, time.Now(), nil)
	require.Error(t, err)

	end := time.Now()
	start := end.Add(time.Hour)
	_, err = NewSalesDiscountRule("rule", &clientID, nil, nil, DiscountTypeFixed, nil, &m, 0, start, &end)
	require.Error(t, err)
}

func TestSalesDiscountRuleAppliesTo(t *testing.T) {
	p, _ := NewPercentage(0.1)
	clientID := uuid.New()
	variantID := uuid.New()
	minQty := 5
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	rule, err := NewSalesDiscountRule("rule", &clientID, &variantID, &minQty, DiscountTypePercentage, &p, nil, 0, start, &end)
	require.NoError(t, err)

	require.False(t, rule.AppliesTo(uuid.New(), variantID, 5, time.Now()))
	require.False(t, rule.AppliesTo(clientID, uuid.New(), 5, time.Now()))
	require.False(t, rule.AppliesTo(clientID, variantID, 4, time.Now()))
	require.False(t, rule.AppliesTo(clientID, variantID, 5, start.Add(-time.Minute)))
	require.False(t, rule.AppliesTo(clientID, variantID, 5, end.Add(time.Minute)))
	require.True(t, rule.AppliesTo(clientID, variantID, 5, time.Now()))

	rule.IsActive = false
	require.False(t, rule.AppliesTo(clientID, variantID, 5, time.Now()))
}

func TestSalesDiscountRuleAppliesToGlobalRule(t *testing.T) {
	p, _ := NewPercentage(0.1)
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	rule, err := NewSalesDiscountRule("rule", nil, nil, nil, DiscountTypePercentage, &p, nil, 0, start, &end)
	require.NoError(t, err)

	require.True(t, rule.AppliesTo(uuid.New(), uuid.New(), 1, time.Now()))
}
