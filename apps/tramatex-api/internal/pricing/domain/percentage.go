package domain

import "fmt"

type Percentage struct {
	value float64
}

func NewPercentage(value float64) (Percentage, error) {
	if value < 0 || value > 1 {
		return Percentage{}, NewValidationError(fmt.Sprintf("percentage out of range: %v", value))
	}
	return Percentage{value: value}, nil
}

func (p Percentage) Value() float64 {
	return p.value
}

func (p Percentage) Apply(amount float64) float64 {
	return amount * p.value
}
