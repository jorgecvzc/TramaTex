package domain

import "fmt"

const DefaultCurrency = "EUR"

type Money struct {
	amount   float64
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
	return Money{amount: amount, currency: currency}, nil
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
