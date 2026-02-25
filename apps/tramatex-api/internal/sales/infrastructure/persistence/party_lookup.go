package persistence

import (
	"context"
	"strings"

	"github.com/google/uuid"

	party_domain "github.com/joran-cortez/tramatex/internal/party/domain"
	party_repo "github.com/joran-cortez/tramatex/internal/party/persistence"
)

type PartyLookupAdapter struct {
	repo party_repo.PartyRepository
}

func NewPartyLookupAdapter(repo party_repo.PartyRepository) *PartyLookupAdapter {
	return &PartyLookupAdapter{repo: repo}
}

func (p *PartyLookupAdapter) ExistsParty(ctx context.Context, partyID uuid.UUID) (bool, error) {
	id, err := party_domain.NewPartyID(partyID.String())
	if err != nil {
		return false, err
	}
	return p.repo.Exists(ctx, id)
}

func (p *PartyLookupAdapter) HasPartyRole(ctx context.Context, partyID uuid.UUID, role string) (bool, error) {
	id, err := party_domain.NewPartyID(partyID.String())
	if err != nil {
		return false, err
	}

	party, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	if party == nil {
		return false, nil
	}

	targetRole := strings.ToUpper(strings.TrimSpace(role))
	for _, partyRole := range party.Roles() {
		if strings.ToUpper(string(partyRole.Type())) == targetRole {
			return true, nil
		}
	}

	return false, nil
}
