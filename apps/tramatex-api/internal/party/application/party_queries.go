package application

import (
	"context"
	"strings"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

// GetPartyQuery represents a query to get a party by ID

type GetPartyQuery struct {
	ID string
}

// GetPartyHandler handles fetching a single party

type GetPartyHandler struct {
	partyRepo persistence.PartyRepository
}

func NewGetPartyHandler(partyRepo persistence.PartyRepository) *GetPartyHandler {
	return &GetPartyHandler{partyRepo: partyRepo}
}

func (h *GetPartyHandler) Handle(ctx context.Context, query *GetPartyQuery) (*domain.Party, error) {
	if query.ID == "" {
		return nil, domain.NewValidationError("party ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(query.ID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapNotFound("failed to fetch party", err)
	}

	return party, nil
}

// ListPartiesQuery represents a query to list parties

type ListPartiesQuery struct {
	Status     string
	Role       string
	Type       string
	Name       string
	TaxID      string
	PageSize   int
	PageNumber int
}

// ListPartiesHandler handles listing parties

type ListPartiesHandler struct {
	partyRepo persistence.PartyRepository
}

func NewListPartiesHandler(partyRepo persistence.PartyRepository) *ListPartiesHandler {
	return &ListPartiesHandler{partyRepo: partyRepo}
}

func (h *ListPartiesHandler) Handle(ctx context.Context, query *ListPartiesQuery) ([]*domain.Party, error) {
	filters := &persistence.PartyFilters{
		Name:       query.Name,
		TaxID:      query.TaxID,
		Type:       strings.ToUpper(query.Type),
		PageNumber: query.PageNumber,
		PageSize:   query.PageSize,
	}

	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}
	if filters.PageNumber <= 0 {
		filters.PageNumber = 1
	}

	if query.Status != "" {
		status := domain.PartyStatus(strings.ToUpper(query.Status))
		if !status.IsValid() {
			return nil, domain.NewValidationErrorf("invalid status: %s", query.Status)
		}
		filters.Status = &status
	}

	if query.Role != "" {
		role := domain.PartyRoleType(strings.ToUpper(query.Role))
		if !role.IsValid() {
			return nil, domain.NewValidationErrorf("invalid role: %s", query.Role)
		}
		filters.Role = &role
	}

	parties, err := h.partyRepo.FindAll(ctx, filters)
	if err != nil {
		return nil, domain.WrapPersistence("failed to list parties", err)
	}

	return parties, nil
}

// ListPartyRelationshipsQuery represents a query to list relationships for a party

type ListPartyRelationshipsQuery struct {
	PartyID string
}

// ListPartyRelationshipsHandler handles listing relationships

type ListPartyRelationshipsHandler struct {
	relRepo persistence.PartyRelationshipRepository
}

func NewListPartyRelationshipsHandler(relRepo persistence.PartyRelationshipRepository) *ListPartyRelationshipsHandler {
	return &ListPartyRelationshipsHandler{relRepo: relRepo}
}

func (h *ListPartyRelationshipsHandler) Handle(ctx context.Context, query *ListPartyRelationshipsQuery) ([]domain.PartyRelationship, error) {
	partyID, err := domain.NewPartyID(query.PartyID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	relationships, err := h.relRepo.FindByPartyID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to list relationships", err)
	}

	return relationships, nil
}

// ListContactDetailsQuery represents a query to list contact details for a party

type ListContactDetailsQuery struct {
	PartyID string
}

// ListContactDetailsHandler handles listing contact details

type ListContactDetailsHandler struct {
	partyRepo persistence.PartyRepository
}

func NewListContactDetailsHandler(partyRepo persistence.PartyRepository) *ListContactDetailsHandler {
	return &ListContactDetailsHandler{partyRepo: partyRepo}
}

func (h *ListContactDetailsHandler) Handle(ctx context.Context, query *ListContactDetailsQuery) ([]*domain.ContactDetails, error) {
	partyID, err := domain.NewPartyID(query.PartyID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapNotFound("party not found", err)
	}

	orgProfile := party.OrganizationProfile()
	if orgProfile == nil {
		return []*domain.ContactDetails{}, nil
	}

	return orgProfile.Contacts(), nil
}

// ListPartyAddressesQuery represents a query to list addresses for a party

type ListPartyAddressesQuery struct {
	PartyID string
}

// ListPartyAddressesHandler handles listing addresses

type ListPartyAddressesHandler struct {
	addressRepo persistence.PartyAddressRepository
}

func NewListPartyAddressesHandler(addressRepo persistence.PartyAddressRepository) *ListPartyAddressesHandler {
	return &ListPartyAddressesHandler{addressRepo: addressRepo}
}

func (h *ListPartyAddressesHandler) Handle(ctx context.Context, query *ListPartyAddressesQuery) ([]*persistence.AddressWithID, error) {
	partyID, err := domain.NewPartyID(query.PartyID)
	if err != nil {
		return nil, domain.WrapValidation("invalid party ID", err)
	}

	addresses, err := h.addressRepo.FindByPartyID(ctx, partyID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to list addresses", err)
	}

	return addresses, nil
}

// GetPartiesBatchQuery represents a query to get multiple parties by IDs

type GetPartiesBatchQuery struct {
	IDs []string
}

// GetPartiesBatchHandler handles batch retrieval of parties

type GetPartiesBatchHandler struct {
	partyRepo persistence.PartyRepository
}

func NewGetPartiesBatchHandler(partyRepo persistence.PartyRepository) *GetPartiesBatchHandler {
	return &GetPartiesBatchHandler{partyRepo: partyRepo}
}

func (h *GetPartiesBatchHandler) Handle(ctx context.Context, query *GetPartiesBatchQuery) ([]*domain.Party, error) {
	if len(query.IDs) == 0 {
		return []*domain.Party{}, nil
	}

	// Validate all IDs first
	partyIDs := make([]domain.PartyID, 0, len(query.IDs))
	for _, idStr := range query.IDs {
		partyID, err := domain.NewPartyID(idStr)
		if err != nil {
			return nil, domain.WrapValidation("invalid party ID", err)
		}
		partyIDs = append(partyIDs, partyID)
	}

	// Fetch all parties - use individual calls for now
	// In a future optimization, we can add FindByIDs to the repository interface
	parties := make([]*domain.Party, 0, len(partyIDs))
	for _, partyID := range partyIDs {
		party, err := h.partyRepo.FindByID(ctx, partyID)
		if err != nil {
			// Skip parties that don't exist instead of failing
			continue
		}
		parties = append(parties, party)
	}

	return parties, nil
}
