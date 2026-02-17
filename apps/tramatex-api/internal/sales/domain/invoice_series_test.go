package domain

import (
	"testing"
)

func TestNewInvoiceSeries(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		year    int
		wantErr bool
	}{
		{"valid series A/2026", "A", 2026, false},
		{"valid series TKT/2026", "TKT", 2026, false},
		{"valid series B/2025", "B", 2025, false},
		{"lowercase code normalized", "tkt", 2026, false},
		{"empty code", "", 2026, true},
		{"code too long", "VERYLONGCODE123", 2026, true},
		{"year too old", "A", 1999, true},
		{"year too far", "A", 2101, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInvoiceSeries(tt.code, tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInvoiceSeries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Code() == "" {
					t.Error("NewInvoiceSeries() returned empty code")
				}
				if got.Year() != tt.year {
					t.Errorf("NewInvoiceSeries() year = %v, want %v", got.Year(), tt.year)
				}
				// Default prefix should match code (uppercased)
				if got.Prefix() != got.Code() {
					t.Errorf("NewInvoiceSeries() prefix = %v, want %v", got.Prefix(), got.Code())
				}
			}
		})
	}
}

func TestNewInvoiceSeriesWithPrefix(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		year    int
		prefix  string
		wantErr bool
	}{
		{"valid with custom prefix", "A", 2026, "FACT", false},
		{"prefix normalized", "TKT", 2026, "ticket", false},
		{"empty prefix", "A", 2026, "", true},
		{"prefix too long", "A", 2026, "VERYLONGPREFIX123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInvoiceSeriesWithPrefix(tt.code, tt.year, tt.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInvoiceSeriesWithPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Prefix() == "" {
				t.Error("NewInvoiceSeriesWithPrefix() returned empty prefix")
			}
		})
	}
}

func TestInvoiceSeries_FormatNumber(t *testing.T) {
	tests := []struct {
		name   string
		series InvoiceSeries
		number int
		want   string
	}{
		{
			"series A number 1",
			InvoiceSeries{code: "A", year: 2026, prefix: "A"},
			1,
			"A/00001/2026",
		},
		{
			"series TKT number 123",
			InvoiceSeries{code: "TKT", year: 2026, prefix: "TKT"},
			123,
			"TKT/00123/2026",
		},
		{
			"series with custom prefix",
			InvoiceSeries{code: "A", year: 2026, prefix: "FACT"},
			456,
			"FACT/00456/2026",
		},
		{
			"large number",
			InvoiceSeries{code: "B", year: 2025, prefix: "B"},
			99999,
			"B/99999/2025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.series.FormatNumber(tt.number); got != tt.want {
				t.Errorf("InvoiceSeries.FormatNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvoiceSeries_String(t *testing.T) {
	series, _ := NewInvoiceSeries("TKT", 2026)
	want := "TKT/2026"
	if got := series.String(); got != want {
		t.Errorf("InvoiceSeries.String() = %v, want %v", got, want)
	}
}
