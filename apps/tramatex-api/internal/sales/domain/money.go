package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

const DefaultCurrency = "EUR"

type Money struct {
	amount   decimal.Decimal
	currency string
}

func NewMoney(amount float64, currency string) (Money, error) {
	if currency == "" {
		currency = DefaultCurrency
	}
	if currency != DefaultCurrency {
		return Money{}, NewValidationError(fmt.Sprintf("unsupported currency: %s", currency))
	}
	if amount < 0 {
		return Money{}, NewValidationError("amount cannot be negative")
	}
	return Money{amount: decimal.NewFromFloat(amount).Round(2), currency: currency}, nil
}

func NewMoneyFromDecimal(amount decimal.Decimal, currency string) (Money, error) {
	if currency == "" {
		currency = DefaultCurrency
	}
	if currency != DefaultCurrency {
		return Money{}, NewValidationError(fmt.Sprintf("unsupported currency: %s", currency))
	}
	if amount.IsNegative() {
		return Money{}, NewValidationError("amount cannot be negative")
	}
	return Money{amount: amount.Round(2), currency: currency}, nil
}

func (m Money) Amount() float64 {
	f, _ := m.amount.Float64()
	return f
}

func (m Money) Decimal() decimal.Decimal {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, NewValidationError("currency mismatch")
	}
	return NewMoneyFromDecimal(m.amount.Add(other.amount), m.currency)
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, NewValidationError("currency mismatch")
	}
	if m.amount.LessThan(other.amount) {
		return Money{}, NewValidationError("resulting amount cannot be negative")
	}
	return NewMoneyFromDecimal(m.amount.Sub(other.amount), m.currency)
}

func (m Money) Multiply(factor float64) (Money, error) {
	if factor < 0 {
		return Money{}, NewValidationError("factor cannot be negative")
	}
	return NewMoneyFromDecimal(m.amount.Mul(decimal.NewFromFloat(factor)), m.currency)
}
