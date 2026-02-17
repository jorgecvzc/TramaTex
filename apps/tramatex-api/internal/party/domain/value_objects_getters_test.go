package domain

import "testing"

func TestValueObjectGetters(t *testing.T) {
	email, _ := NewEmail("USER@EXAMPLE.COM")
	if email.Value() != "user@example.com" {
		t.Fatalf("expected lowercase email value")
	}

	phone, _ := NewPhone("+34 600 111 222")
	if phone.Value() == "" {
		t.Fatalf("expected phone value")
	}

	taxID, _ := NewTaxID("B12345678", "CIF")
	if taxID.Value() != "B12345678" {
		t.Fatalf("expected tax ID value")
	}

	address, err := NewAddress("Calle 1", "Madrid", "Madrid", "28001", "Spain")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if address.Province() != "Madrid" {
		t.Fatalf("expected province Madrid")
	}
	if address.PostalCode() != "28001" {
		t.Fatalf("expected postal code 28001")
	}
	if address.Country() != "Spain" {
		t.Fatalf("expected country Spain")
	}
}
