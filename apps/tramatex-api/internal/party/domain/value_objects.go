package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// Email is a value object representing an email address
type Email struct {
	value string
}

// NewEmail creates a new Email value object
func NewEmail(value string) (*Email, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	// Email regex: allows plus-addressing and more flexible formats
	// Matches: test@example.com, test+tag@example.com, test.email@example.com
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.+_%\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return nil, fmt.Errorf("invalid email format: %s", value)
	}

	return &Email{value: strings.ToLower(value)}, nil
}

// String returns the email value
func (e *Email) String() string {
	return e.value
}

// Value returns the email value for database driver
func (e *Email) Value() string {
	return e.value
}

// Equals compares two Email objects
func (e *Email) Equals(other *Email) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}

// Phone is a value object representing a phone number
type Phone struct {
	value string
}

// NewPhone creates a new Phone value object
func NewPhone(value string) (*Phone, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("phone cannot be empty")
	}

	// Basic phone regex: allows +, -, spaces, and digits
	// Example: +34 666 123456 or 666123456
	phoneRegex := regexp.MustCompile(`^[\+]?[\d\s\-()]{8,}$`)
	if !phoneRegex.MatchString(value) {
		return nil, fmt.Errorf("invalid phone format: %s", value)
	}

	return &Phone{value: value}, nil
}

// String returns the phone value
func (p *Phone) String() string {
	return p.value
}

// Value returns the phone value for database driver
func (p *Phone) Value() string {
	return p.value
}

// Equals compares two Phone objects
func (p *Phone) Equals(other *Phone) bool {
	if other == nil {
		return false
	}
	return p.value == other.value
}

// TaxID is a value object representing a tax identification number (CIF, NIF, VAT, etc.)
type TaxID struct {
	value string
	typ   string // "CIF", "NIF", "VAT", etc.
}

// NewTaxID creates a new TaxID value object
func NewTaxID(value, typ string) (*TaxID, error) {
	value = strings.TrimSpace(strings.ToUpper(value))
	typ = strings.TrimSpace(strings.ToUpper(typ))

	if value == "" {
		return nil, fmt.Errorf("tax ID cannot be empty")
	}
	if typ == "" {
		return nil, fmt.Errorf("tax ID type cannot be empty")
	}

	// Basic validation: must be alphanumeric, 5-20 chars
	taxIDRegex := regexp.MustCompile(`^[A-Z0-9]{5,20}$`)
	if !taxIDRegex.MatchString(value) {
		return nil, fmt.Errorf("invalid tax ID format: %s (type: %s)", value, typ)
	}

	return &TaxID{value: value, typ: typ}, nil
}

// String returns the tax ID value
func (t *TaxID) String() string {
	return t.value
}

// Value returns the tax ID value for database driver
func (t *TaxID) Value() string {
	return t.value
}

// Type returns the tax ID type
func (t *TaxID) Type() string {
	return t.typ
}

// Equals compares two TaxID objects
func (t *TaxID) Equals(other *TaxID) bool {
	if other == nil {
		return false
	}
	return t.value == other.value && t.typ == other.typ
}

// Address is a value object representing a physical address
type Address struct {
	street     string
	city       string
	province   string
	postalCode string
	country    string
}

// NewAddress creates a new Address value object
func NewAddress(street, city, province, postalCode, country string) (*Address, error) {
	street = strings.TrimSpace(street)
	city = strings.TrimSpace(city)
	province = strings.TrimSpace(province)
	postalCode = strings.TrimSpace(postalCode)
	country = strings.TrimSpace(country)

	if street == "" || city == "" || postalCode == "" || country == "" {
		return nil, fmt.Errorf("address requires street, city, postal code, and country")
	}

	if len(postalCode) < 3 || len(postalCode) > 10 {
		return nil, fmt.Errorf("invalid postal code format: %s", postalCode)
	}

	return &Address{
		street:     street,
		city:       city,
		province:   province,
		postalCode: postalCode,
		country:    country,
	}, nil
}

// Street returns the street
func (a *Address) Street() string {
	return a.street
}

// City returns the city
func (a *Address) City() string {
	return a.city
}

// Province returns the province
func (a *Address) Province() string {
	return a.province
}

// PostalCode returns the postal code
func (a *Address) PostalCode() string {
	return a.postalCode
}

// Country returns the country
func (a *Address) Country() string {
	return a.country
}

// Equals compares two Address objects
func (a *Address) Equals(other *Address) bool {
	if other == nil {
		return false
	}
	return a.street == other.street &&
		a.city == other.city &&
		a.province == other.province &&
		a.postalCode == other.postalCode &&
		a.country == other.country
}
