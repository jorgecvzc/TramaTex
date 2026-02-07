package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPricingRuleValidation(t *testing.T) {
	p, _ := NewPercentage(0.1)
	_, err := NewPricingRule("", nil, nil, p, 0, nil, time.Now(), nil)
	require.Error(t, err)

	_, err = NewPricingRule("rule", nil, nil, p, -1, nil, time.Now(), nil)
	require.Error(t, err)

	max := 1
	_, err = NewPricingRule("rule", nil, nil, p, 2, &max, time.Now(), nil)
	require.Error(t, err)

	past := time.Now().Add(2 * time.Hour)
	future := time.Now()
	_, err = NewPricingRule("rule", nil, nil, p, 0, nil, past, &future)
	require.Error(t, err)
}

func TestPricingRuleAppliesTo(t *testing.T) {
	p, _ := NewPercentage(0.1)
	variantID := uuid.New()
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	rule, err := NewPricingRule("rule", &variantID, nil, p, 5, nil, start, &end)
	require.NoError(t, err)

	require.False(t, rule.AppliesTo(uuid.New(), 5, time.Now()))
	require.False(t, rule.AppliesTo(variantID, 4, time.Now()))
	require.False(t, rule.AppliesTo(variantID, 5, start.Add(-time.Minute)))
	require.False(t, rule.AppliesTo(variantID, 5, end.Add(time.Minute)))
	require.True(t, rule.AppliesTo(variantID, 5, time.Now()))

	rule.IsActive = false
	require.False(t, rule.AppliesTo(variantID, 5, time.Now()))
}
