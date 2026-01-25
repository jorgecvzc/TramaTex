package domain

import (
	"testing"
)

// Test Email Value Object
func TestNewEmail_Valid(t *testing.T) {
	email, err := NewEmail("test@example.com")
	if err != nil {
		t.Errorf("NewEmail should not error for valid email, got: %v", err)
	}
	if email.String() != "test@example.com" {
		t.Errorf("Email value should be 'test@example.com', got: %s", email.String())
	}
}

func TestNewEmail_CaseInsensitive(t *testing.T) {
	email, err := NewEmail("Test@Example.COM")
	if err != nil {
		t.Errorf("NewEmail should normalize case, got error: %v", err)
	}
	if email.String() != "test@example.com" {
		t.Errorf("Email should be lowercase, got: %s", email.String())
	}
}

func TestNewEmail_Invalid(t *testing.T) {
	tests := []string{
		"",
		"invalid.email",
		"@example.com",
		"test@",
		"test@.com",
	}

	for _, invalidEmail := range tests {
		_, err := NewEmail(invalidEmail)
		if err == nil {
			t.Errorf("NewEmail should error for invalid email: %s", invalidEmail)
		}
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("test@example.com")
	email2, _ := NewEmail("test@example.com")
	email3, _ := NewEmail("other@example.com")

	if !email1.Equals(email2) {
		t.Error("Equal emails should compare as equal")
	}
	if email1.Equals(email3) {
		t.Error("Different emails should not compare as equal")
	}
	if email1.Equals(nil) {
		t.Error("Email should not equal nil")
	}
}

// Test Phone Value Object
func TestNewPhone_Valid(t *testing.T) {
	tests := []string{
		"666123456",
		"+34 666 123456",
		"+34-666-123456",
		"(34) 666 123456",
	}

	for _, validPhone := range tests {
		phone, err := NewPhone(validPhone)
		if err != nil {
			t.Errorf("NewPhone should not error for valid phone %s, got: %v", validPhone, err)
		}
		if phone == nil {
			t.Errorf("Phone should not be nil for valid input: %s", validPhone)
		}
	}
}

func TestNewPhone_Invalid(t *testing.T) {
	tests := []string{
		"",
		"123", // too short
		"abc", // non-numeric
		"!!!", // invalid characters
	}

	for _, invalidPhone := range tests {
		_, err := NewPhone(invalidPhone)
		if err == nil {
			t.Errorf("NewPhone should error for invalid phone: %s", invalidPhone)
		}
	}
}

func TestPhone_Equals(t *testing.T) {
	phone1, _ := NewPhone("666123456")
	phone2, _ := NewPhone("666123456")
	phone3, _ := NewPhone("777123456")

	if !phone1.Equals(phone2) {
		t.Error("Equal phones should compare as equal")
	}
	if phone1.Equals(phone3) {
		t.Error("Different phones should not compare as equal")
	}
	if phone1.Equals(nil) {
		t.Error("Phone should not equal nil")
	}
}

// Test TaxID Value Object
func TestNewTaxID_Valid(t *testing.T) {
	taxID, err := NewTaxID("A12345678", "CIF")
	if err != nil {
		t.Errorf("NewTaxID should not error for valid tax ID, got: %v", err)
	}
	if taxID.String() != "A12345678" {
		t.Errorf("TaxID value should be 'A12345678', got: %s", taxID.String())
	}
	if taxID.Type() != "CIF" {
		t.Errorf("TaxID type should be 'CIF', got: %s", taxID.Type())
	}
}

func TestNewTaxID_CaseInsensitive(t *testing.T) {
	taxID, err := NewTaxID("a12345678", "cif")
	if err != nil {
		t.Errorf("NewTaxID should normalize case, got error: %v", err)
	}
	if taxID.String() != "A12345678" {
		t.Errorf("TaxID should be uppercase, got: %s", taxID.String())
	}
	if taxID.Type() != "CIF" {
		t.Errorf("TaxID type should be uppercase, got: %s", taxID.Type())
	}
}

func TestNewTaxID_Invalid(t *testing.T) {
	tests := []struct {
		value string
		typ   string
	}{
		{"", "CIF"},
		{"ABC", "CIF"},                    // too short
		{"AB12345678901234567890", "CIF"}, // too long
		{"ABC@#$%", "CIF"},                // invalid chars
		{"A12345678", ""},                 // no type
	}

	for _, test := range tests {
		_, err := NewTaxID(test.value, test.typ)
		if err == nil {
			t.Errorf("NewTaxID should error for invalid input: %s (type: %s)", test.value, test.typ)
		}
	}
}

func TestTaxID_Equals(t *testing.T) {
	taxID1, _ := NewTaxID("A12345678", "CIF")
	taxID2, _ := NewTaxID("A12345678", "CIF")
	taxID3, _ := NewTaxID("B12345678", "CIF")
	taxID4, _ := NewTaxID("A12345678", "NIF")

	if !taxID1.Equals(taxID2) {
		t.Error("Equal tax IDs should compare as equal")
	}
	if taxID1.Equals(taxID3) {
		t.Error("Tax IDs with different values should not be equal")
	}
	if taxID1.Equals(taxID4) {
		t.Error("Tax IDs with different types should not be equal")
	}
	if taxID1.Equals(nil) {
		t.Error("Tax ID should not equal nil")
	}
}

// Test Address Value Object
func TestNewAddress_Valid(t *testing.T) {
	addr, err := NewAddress("Calle Principal 123", "Madrid", "Madrid", "28001", "Spain")
	if err != nil {
		t.Errorf("NewAddress should not error for valid address, got: %v", err)
	}
	if addr.Street() != "Calle Principal 123" {
		t.Errorf("Street should be 'Calle Principal 123', got: %s", addr.Street())
	}
	if addr.City() != "Madrid" {
		t.Errorf("City should be 'Madrid', got: %s", addr.City())
	}
}

func TestNewAddress_MissingFields(t *testing.T) {
	tests := []struct {
		street     string
		city       string
		province   string
		postalCode string
		country    string
	}{
		{"", "Madrid", "Madrid", "28001", "Spain"},
		{"Calle 123", "", "Madrid", "28001", "Spain"},
		{"Calle 123", "Madrid", "Madrid", "", "Spain"},
		{"Calle 123", "Madrid", "Madrid", "28001", ""},
	}

	for _, test := range tests {
		_, err := NewAddress(test.street, test.city, test.province, test.postalCode, test.country)
		if err == nil {
			t.Error("NewAddress should error when required fields are empty")
		}
	}
}

func TestNewAddress_InvalidPostalCode(t *testing.T) {
	tests := []string{
		"12",            // too short
		"1234567890123", // too long
	}

	for _, postalCode := range tests {
		_, err := NewAddress("Calle 123", "Madrid", "Madrid", postalCode, "Spain")
		if err == nil {
			t.Errorf("NewAddress should error for invalid postal code: %s", postalCode)
		}
	}
}

func TestAddress_Equals(t *testing.T) {
	addr1, _ := NewAddress("Calle 123", "Madrid", "Madrid", "28001", "Spain")
	addr2, _ := NewAddress("Calle 123", "Madrid", "Madrid", "28001", "Spain")
	addr3, _ := NewAddress("Calle 456", "Madrid", "Madrid", "28001", "Spain")

	if !addr1.Equals(addr2) {
		t.Error("Equal addresses should compare as equal")
	}
	if addr1.Equals(addr3) {
		t.Error("Different addresses should not compare as equal")
	}
	if addr1.Equals(nil) {
		t.Error("Address should not equal nil")
	}
}
