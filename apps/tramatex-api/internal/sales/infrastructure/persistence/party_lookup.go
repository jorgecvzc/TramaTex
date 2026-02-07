package persistence

import (
	"context"

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
