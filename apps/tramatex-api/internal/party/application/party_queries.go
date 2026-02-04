package application

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("party ID cannot be empty")
	}

	partyID, err := domain.NewPartyID(query.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch party: %w", err)
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
			return nil, fmt.Errorf("invalid status: %s", query.Status)
		}
		filters.Status = &status
	}

	if query.Role != "" {
		role := domain.PartyRoleType(strings.ToUpper(query.Role))
		if !role.IsValid() {
			return nil, fmt.Errorf("invalid role: %s", query.Role)
		}
		filters.Role = &role
	}

	parties, err := h.partyRepo.FindAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list parties: %w", err)
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
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	relationships, err := h.relRepo.FindByPartyID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list relationships: %w", err)
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
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	party, err := h.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("party not found: %w", err)
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

func (h *ListPartyAddressesHandler) Handle(ctx context.Context, query *ListPartyAddressesQuery) ([]*domain.Address, error) {
	partyID, err := domain.NewPartyID(query.PartyID)
	if err != nil {
		return nil, fmt.Errorf("invalid party ID: %w", err)
	}

	addresses, err := h.addressRepo.FindByPartyID(ctx, partyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}

	return addresses, nil
}
