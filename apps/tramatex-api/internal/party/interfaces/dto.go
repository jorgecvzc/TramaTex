package interfaces

import (
	"time"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

// DTOs for REST API

// OrganizationDTO represents an organization in API responses
type OrganizationDTO struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	TaxID      string    `json:"tax_id,omitempty"`
	Website    string    `json:"website,omitempty"`
	Notes      string    `json:"notes,omitempty"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedBy string    `json:"modified_by,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

// CreateOrganizationRequest represents a request to create organization
type CreateOrganizationRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"` // CLIENT, SUPPLIER, BOTH
	TaxID     string `json:"tax_id,omitempty"`
	TaxIDType string `json:"tax_id_type,omitempty"`
	Website   string `json:"website,omitempty"`
}

// UpdateOrganizationRequest represents a request to update organization
type UpdateOrganizationRequest struct {
	Name    string `json:"name,omitempty"`
	Website string `json:"website,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

// ChangeStatusRequest represents a request to change organization status
type ChangeStatusRequest struct {
	Status string `json:"status"` // ACTIVE, INACTIVE
}

// PersonDTO represents a person in API responses
type PersonDTO struct {
	ID               string    `json:"id"`
	OrganizationID   string    `json:"organization_id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone,omitempty"`
	JobTitle         string    `json:"job_title,omitempty"`
	IsPrimaryContact bool      `json:"is_primary_contact"`
	CreatedBy        string    `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	ModifiedBy       string    `json:"modified_by,omitempty"`
	ModifiedAt       time.Time `json:"modified_at,omitempty"`
}

// CreatePersonRequest represents a request to create person
type CreatePersonRequest struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	JobTitle  string `json:"job_title,omitempty"`
	IsPrimary bool   `json:"is_primary_contact,omitempty"`
}

// UpdatePersonRequest represents a request to update person
type UpdatePersonRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	JobTitle  string `json:"job_title,omitempty"`
	IsPrimary bool   `json:"is_primary_contact,omitempty"`
}

// AddressDTO represents an address in API responses
type AddressDTO struct {
	ID         string    `json:"id"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	Province   string    `json:"province,omitempty"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateAddressRequest represents a request to create address
type CreateAddressRequest struct {
	ID         string `json:"id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country,omitempty"`
}

// ErrorResponse represents an error in API responses
type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
}

// ListResponse represents a list response with pagination
type ListResponse struct {
	Data       interface{} `json:"data"`
	PageNumber int         `json:"page_number,omitempty"`
	PageSize   int         `json:"page_size,omitempty"`
	Total      int         `json:"total,omitempty"`
}

// Mappers

// MapOrganizationToDTO converts domain organization to DTO
func MapOrganizationToDTO(org *domain.Organization) *OrganizationDTO {
	if org == nil {
		return nil
	}

	dto := &OrganizationDTO{
		ID:         org.ID().String(),
		Name:       org.Name(),
		Role:       org.Role().String(),
		Status:     org.Status().String(),
		Website:    org.Website(),
		Notes:      org.Notes(),
		CreatedBy:  org.CreatedBy(),
		CreatedAt:  org.CreatedAt(),
		ModifiedBy: org.ModifiedBy(),
		ModifiedAt: org.ModifiedAt(),
	}

	if org.TaxID() != nil {
		dto.TaxID = org.TaxID().Value()
	}

	return dto
}

// MapPersonToDTO converts domain person to DTO
func MapPersonToDTO(person *domain.Person) *PersonDTO {
	if person == nil {
		return nil
	}

	dto := &PersonDTO{
		ID:               person.ID().String(),
		OrganizationID:   person.OrganizationID().String(),
		FirstName:        person.FirstName(),
		LastName:         person.LastName(),
		Email:            person.Email().Value(),
		JobTitle:         person.JobTitle(),
		IsPrimaryContact: person.IsPrimaryContact(),
		CreatedBy:        person.CreatedBy(),
		CreatedAt:        person.CreatedAt(),
		ModifiedBy:       person.ModifiedBy(),
		ModifiedAt:       person.ModifiedAt(),
	}

	if person.Phone() != nil {
		dto.Phone = person.Phone().Value()
	}

	return dto
}

// MapAddressToDTO converts domain address to DTO
func MapAddressToDTO(addr *domain.Address) *AddressDTO {
	if addr == nil {
		return nil
	}

	return &AddressDTO{
		Street:     addr.Street(),
		City:       addr.City(),
		Province:   addr.Province(),
		PostalCode: addr.PostalCode(),
		Country:    addr.Country(),
	}
}
