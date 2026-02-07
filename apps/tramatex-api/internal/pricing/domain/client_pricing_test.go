package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestClientPricingValidation(t *testing.T) {
	price, _ := NewMoney(10, DefaultCurrency)
	clientID := uuid.New()
	variantID := uuid.New()

	_, err := NewClientPricing(uuid.Nil, variantID, price, time.Now(), nil)
	require.Error(t, err)

	_, err = NewClientPricing(clientID, uuid.Nil, price, time.Now(), nil)
	require.Error(t, err)

	end := time.Now()
	start := end.Add(time.Hour)
	_, err = NewClientPricing(clientID, variantID, price, start, &end)
	require.Error(t, err)
}

func TestClientPricingAppliesTo(t *testing.T) {
	price, _ := NewMoney(10, DefaultCurrency)
	clientID := uuid.New()
	variantID := uuid.New()
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	override, err := NewClientPricing(clientID, variantID, price, start, &end)
	require.NoError(t, err)

	require.False(t, override.AppliesTo(uuid.New(), variantID, time.Now()))
	require.False(t, override.AppliesTo(clientID, uuid.New(), time.Now()))
	require.False(t, override.AppliesTo(clientID, variantID, start.Add(-time.Minute)))
	require.False(t, override.AppliesTo(clientID, variantID, end.Add(time.Minute)))
	require.True(t, override.AppliesTo(clientID, variantID, time.Now()))

	override.IsActive = false
	require.False(t, override.AppliesTo(clientID, variantID, time.Now()))
}
