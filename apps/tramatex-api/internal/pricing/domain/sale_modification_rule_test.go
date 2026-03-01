package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaleModificationRuleValidation(t *testing.T) {
	p, _ := NewPercentage(0.1)
	value, _ := NewRuleValue(RuleValuePercentageMarkup, &p, nil)

	_, err := NewSaleModificationRule("", nil, nil, nil, value, 0, time.Now(), nil)
	require.Error(t, err)

	_, err = NewSaleModificationRule("rule", nil, nil, nil, value, -1, time.Now(), nil)
	require.Error(t, err)

	end := time.Now()
	start := end.Add(time.Hour)
	_, err = NewSaleModificationRule("rule", nil, nil, nil, value, 0, start, &end)
	require.Error(t, err)
}

func TestSaleModificationRuleAppliesTo(t *testing.T) {
	p, _ := NewPercentage(0.1)
	value, _ := NewRuleValue(RuleValuePercentageMarkup, &p, nil)
	clientID := uuid.New().String()
	groupID := uuid.New()
	minOrder, _ := NewMoney(100, DefaultCurrency)
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	rule, err := NewSaleModificationRule("rule", []string{clientID}, &groupID, &minOrder, value, 1, start, &end)
	require.NoError(t, err)

	orderTotal := Money{amount: 200, currency: DefaultCurrency}
	require.False(t, rule.AppliesTo(uuid.New().String(), &groupID, orderTotal, time.Now()))
	require.False(t, rule.AppliesTo(clientID, nil, orderTotal, time.Now()))

	lowOrder := Money{amount: 50, currency: DefaultCurrency}
	require.False(t, rule.AppliesTo(clientID, &groupID, lowOrder, time.Now()))
	require.False(t, rule.AppliesTo(clientID, &groupID, orderTotal, start.Add(-time.Minute)))
	require.False(t, rule.AppliesTo(clientID, &groupID, orderTotal, end.Add(time.Minute)))
	require.True(t, rule.AppliesTo(clientID, &groupID, orderTotal, time.Now()))

	rule.IsActive = false
	require.False(t, rule.AppliesTo(clientID, &groupID, orderTotal, time.Now()))
}

func TestSaleModificationRuleAppliesToCurrencyMismatch(t *testing.T) {
	p, _ := NewPercentage(0.1)
	value, _ := NewRuleValue(RuleValuePercentageMarkup, &p, nil)
	minOrder := Money{amount: 100, currency: DefaultCurrency}
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	rule, err := NewSaleModificationRule("rule", nil, nil, &minOrder, value, 1, start, &end)
	require.NoError(t, err)

	orderTotal := Money{amount: 200, currency: "USD"}
	require.False(t, rule.AppliesTo(uuid.New().String(), nil, orderTotal, time.Now()))
}

func TestSaleModificationRuleAppliesToNoFilters(t *testing.T) {
	p, _ := NewPercentage(0.1)
	value, _ := NewRuleValue(RuleValuePercentageMarkup, &p, nil)
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	rule, err := NewSaleModificationRule("rule", nil, nil, nil, value, 1, start, &end)
	require.NoError(t, err)

	orderTotal := Money{amount: 1, currency: DefaultCurrency}
	require.True(t, rule.AppliesTo(uuid.New().String(), nil, orderTotal, time.Now()))
}
