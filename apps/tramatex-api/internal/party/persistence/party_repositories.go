package persistence

import (
	"context"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

// PartyRepository defines methods for Party aggregate persistence
// Implementations should persist profiles, roles, and contact details together.
type PartyRepository interface {
	Save(ctx context.Context, party *domain.Party, createdBy string, modifiedBy string) error
	FindByID(ctx context.Context, id domain.PartyID) (*domain.Party, error)
	FindAll(ctx context.Context, filters *PartyFilters) ([]*domain.Party, error)
	Delete(ctx context.Context, id domain.PartyID) error
	Exists(ctx context.Context, id domain.PartyID) (bool, error)
	Count(ctx context.Context) (int64, error)
	HasContactDetailsReferences(ctx context.Context, partyID domain.PartyID) (bool, error)
	HasMESWorkReferences(ctx context.Context, partyID domain.PartyID) (bool, error)
	HasSalesReferences(ctx context.Context, partyID domain.PartyID) (bool, error)
}

// PartyRelationshipRepository defines methods for relationships persistence
type PartyRelationshipRepository interface {
	Save(ctx context.Context, relationship domain.PartyRelationship, createdBy string, modifiedBy string) error
	FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]domain.PartyRelationship, error)
	Delete(ctx context.Context, id domain.PartyRelationshipID) error
}

// ContactDetailsRepository defines methods for contact details persistence
// If contacts are persisted as part of Party aggregate, implementations can proxy to PartyRepository.
type ContactDetailsRepository interface {
	Save(ctx context.Context, partyID domain.PartyID, details *domain.ContactDetails, createdBy string, modifiedBy string) error
	FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*domain.ContactDetails, error)
	Delete(ctx context.Context, id domain.ContactDetailsID) error
}

// PartyAddressRepository defines methods for address persistence attached to Party
type PartyAddressRepository interface {
	Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error
	FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*AddressWithID, error)
	FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error)
	Delete(ctx context.Context, id domain.AddressID) error
}

// AddressWithID contains an address with its ID
type AddressWithID struct {
	ID      string
	Address *domain.Address
}

// PartyFilters defines filtering options for parties
// Type should be one of: PERSON, ORGANIZATION, BOTH (optional)
type PartyFilters struct {
	Status     *domain.PartyStatus
	Role       *domain.PartyRoleType
	Type       string
	Name       string
	TaxID      string
	PageSize   int
	PageNumber int
}
