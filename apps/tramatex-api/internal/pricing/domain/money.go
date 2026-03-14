package domain

import (
	"fmt"
	"math"
)

const DefaultCurrency = "EUR"

type Money struct {
	amount   float64
	currency string
}

// roundTo2Decimals rounds a float64 to 2 decimal places using round-half-up
// (commercial rounding). This avoids banker's rounding where .5 rounds to even,
// which caused subtle discrepancies in discount calculations (e.g. 5.025 → 5.02
// instead of the expected 5.03).
func roundTo2Decimals(amount float64) float64 {
	return math.Floor(amount*100+0.5) / 100
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
	return Money{amount: roundTo2Decimals(amount), currency: currency}, nil
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, NewValidationError("currency mismatch")
	}
	return NewMoney(m.amount+other.amount, m.currency)
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, NewValidationError("currency mismatch")
	}
	if m.amount < other.amount {
		return Money{}, NewValidationError("resulting amount cannot be negative")
	}
	return NewMoney(m.amount-other.amount, m.currency)
}

func (m Money) Multiply(factor float64) (Money, error) {
	if factor < 0 {
		return Money{}, NewValidationError("factor cannot be negative")
	}
	return NewMoney(m.amount*factor, m.currency)
}
