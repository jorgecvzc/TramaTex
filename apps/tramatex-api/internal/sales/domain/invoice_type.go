package domain

import "fmt"

// InvoiceType represents the type of invoice according to Spanish legislation
// COMPLETA: Full invoice for B2B transactions (requires complete Party fiscal data)
// SIMPLIFICADA: Simplified invoice/ticket for retail sales (< 3,000 EUR, can use generic CONSUMIDOR_FINAL)
type InvoiceType string

const (
	InvoiceTypeComplete   InvoiceType = "COMPLETA"
	InvoiceTypeSimplified InvoiceType = "SIMPLIFICADA"
)

func (it InvoiceType) IsValid() error {
	switch it {
	case InvoiceTypeComplete, InvoiceTypeSimplified:
		return nil
	default:
		return NewValidationError(fmt.Sprintf("invalid invoice type: %s", it))
	}
}

func (it InvoiceType) String() string {
	return string(it)
}
