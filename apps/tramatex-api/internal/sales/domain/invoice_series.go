package domain

import (
	"fmt"
	"strings"
)

// InvoiceSeries represents a Value Object for invoice numbering series
// It encapsulates the series code, year, and prefix to generate unique invoice numbers
// Example format: "TKT/00123/2026" (series: TKT, number: 123, year: 2026)
type InvoiceSeries struct {
	code   string // e.g., "A" for B2B invoices, "TKT" for tickets, "B" for rectificativas
	year   int    // Fiscal year (e.g., 2026)
	prefix string // Optional prefix (default: same as code)
}

// NewInvoiceSeries creates a new invoice series with validation
func NewInvoiceSeries(code string, year int) (InvoiceSeries, error) {
	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		return InvoiceSeries{}, NewValidationError("invoice series code cannot be empty")
	}
	if len(code) > 10 {
		return InvoiceSeries{}, NewValidationError("invoice series code too long (max 10 characters)")
	}
	if year < 2000 || year > 2100 {
		return InvoiceSeries{}, NewValidationError("invalid year for invoice series")
	}

	return InvoiceSeries{
		code:   code,
		year:   year,
		prefix: code, // Default prefix is the same as code
	}, nil
}

// NewInvoiceSeriesWithPrefix creates a new invoice series with a custom prefix
func NewInvoiceSeriesWithPrefix(code string, year int, prefix string) (InvoiceSeries, error) {
	series, err := NewInvoiceSeries(code, year)
	if err != nil {
		return InvoiceSeries{}, err
	}
	prefix = strings.TrimSpace(strings.ToUpper(prefix))
	if prefix == "" {
		return InvoiceSeries{}, NewValidationError("invoice series prefix cannot be empty")
	}
	if len(prefix) > 10 {
		return InvoiceSeries{}, NewValidationError("invoice series prefix too long (max 10 characters)")
	}
	series.prefix = prefix
	return series, nil
}

// Code returns the series code
func (s InvoiceSeries) Code() string {
	return s.code
}

// Year returns the fiscal year
func (s InvoiceSeries) Year() int {
	return s.year
}

// Prefix returns the series prefix
func (s InvoiceSeries) Prefix() string {
	return s.prefix
}

// FormatNumber generates the formatted invoice number with series
// Example: series "TKT", number 123, year 2026 → "TKT/00123/2026"
func (s InvoiceSeries) FormatNumber(number int) string {
	return fmt.Sprintf("%s/%05d/%d", s.prefix, number, s.year)
}

// String returns the series identifier (code/year)
func (s InvoiceSeries) String() string {
	return fmt.Sprintf("%s/%d", s.code, s.year)
}
