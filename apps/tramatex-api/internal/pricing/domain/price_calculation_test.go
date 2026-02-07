package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPriceCalculationValidation(t *testing.T) {
	price, _ := NewMoney(10, DefaultCurrency)
	variantID := uuid.New()
	clientID := uuid.New()

	_, err := NewPriceCalculation(uuid.Nil, clientID, 1, price, price, nil)
	require.Error(t, err)

	_, err = NewPriceCalculation(variantID, uuid.Nil, 1, price, price, nil)
	require.Error(t, err)

	_, err = NewPriceCalculation(variantID, clientID, 0, price, price, nil)
	require.Error(t, err)
}

func TestPriceCalculationSuccess(t *testing.T) {
	price, _ := NewMoney(10, DefaultCurrency)
	variantID := uuid.New()
	clientID := uuid.New()

	calc, err := NewPriceCalculation(variantID, clientID, 1, price, price, []string{"rule"})
	require.NoError(t, err)
	require.Equal(t, variantID, calc.ProductVariantID)
}
