package domain

import (
	"testing"
)

func TestInvoiceType_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		it      InvoiceType
		wantErr bool
	}{
		{"valid complete", InvoiceTypeComplete, false},
		{"valid simplified", InvoiceTypeSimplified, false},
		{"invalid empty", InvoiceType(""), true},
		{"invalid unknown", InvoiceType("UNKNOWN"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.it.IsValid()
			if (err != nil) != tt.wantErr {
				t.Errorf("InvoiceType.IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInvoiceType_String(t *testing.T) {
	tests := []struct {
		name string
		it   InvoiceType
		want string
	}{
		{"complete", InvoiceTypeComplete, "COMPLETE"},
		{"simplified", InvoiceTypeSimplified, "SIMPLIFIED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.it.String(); got != tt.want {
				t.Errorf("InvoiceType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
