package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMoneyValidation(t *testing.T) {
	_, err := NewMoney(-1, DefaultCurrency)
	require.Error(t, err)

	_, err = NewMoney(10, "USD")
	require.Error(t, err)

	money, err := NewMoney(10, "")
	require.NoError(t, err)
	require.Equal(t, DefaultCurrency, money.Currency())
}

func TestMoneyOperations(t *testing.T) {
	money, err := NewMoney(10, DefaultCurrency)
	require.NoError(t, err)

	other, err := NewMoney(5, DefaultCurrency)
	require.NoError(t, err)

	sum, err := money.Add(other)
	require.NoError(t, err)
	require.Equal(t, 15.0, sum.Amount())

	diff, err := money.Subtract(other)
	require.NoError(t, err)
	require.Equal(t, 5.0, diff.Amount())

	_, err = other.Subtract(money)
	require.Error(t, err)

	product, err := money.Multiply(2)
	require.NoError(t, err)
	require.Equal(t, 20.0, product.Amount())

	_, err = money.Multiply(-1)
	require.Error(t, err)
}

func TestMoneyCurrencyMismatch(t *testing.T) {
	moneyEUR := Money{amount: 10, currency: "EUR"}
	moneyUSD := Money{amount: 5, currency: "USD"}

	_, err := moneyEUR.Add(moneyUSD)
	require.Error(t, err)

	_, err = moneyEUR.Subtract(moneyUSD)
	require.Error(t, err)
}
