package domain

// SumAmounts sums a slice of Money values, all in the same currency.
// Returns zero in DefaultCurrency if the slice is empty.
func SumAmounts(amounts []Money) (Money, error) {
	total, err := NewMoney(0, DefaultCurrency)
	if err != nil {
		return Money{}, err
	}
	for _, amount := range amounts {
		total, err = total.Add(amount)
		if err != nil {
			return Money{}, err
		}
	}
	return total, nil
}
