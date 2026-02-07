package application

import (
	"context"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

type capturingPartyRepo struct {
	*fakePartyRepo
	lastFilters *persistence.PartyFilters
}

func newCapturingPartyRepo() *capturingPartyRepo {
	return &capturingPartyRepo{fakePartyRepo: newFakePartyRepo()}
}

func (r *capturingPartyRepo) FindAll(ctx context.Context, filters *persistence.PartyFilters) ([]*domain.Party, error) {
	r.lastFilters = filters
	return r.fakePartyRepo.FindAll(ctx, filters)
}

func TestGetPartyHandler_Success(t *testing.T) {
	repo := newFakePartyRepo()
	partyID, _ := domain.NewPartyID("party-q1")
	personProfile, _ := domain.NewPersonProfile("Ana", "Perez")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, "user-1", personProfile, nil)
	_ = repo.Save(context.Background(), party)

	handler := NewGetPartyHandler(repo)
	result, err := handler.Handle(context.Background(), &GetPartyQuery{ID: "party-q1"})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if result == nil || result.ID().String() != "party-q1" {
		t.Fatalf("Expected party-q1, got %v", result)
	}
}

func TestGetPartyHandler_EmptyID(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewGetPartyHandler(repo)
	if _, err := handler.Handle(context.Background(), &GetPartyQuery{ID: ""}); err == nil {
		t.Fatalf("Expected error for empty party ID")
	}
}

func TestListPartiesHandler_DefaultPaging(t *testing.T) {
	repo := newCapturingPartyRepo()
	seedPartyWithProfiles(t, repo.fakePartyRepo, "party-q2", false)

	handler := NewListPartiesHandler(repo)
	_, err := handler.Handle(context.Background(), &ListPartiesQuery{})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if repo.lastFilters == nil {
		t.Fatalf("Expected filters to be set")
	}
	if repo.lastFilters.PageNumber != 1 || repo.lastFilters.PageSize != 10 {
		t.Fatalf("Expected defaults page 1 size 10, got %d size %d", repo.lastFilters.PageNumber, repo.lastFilters.PageSize)
	}
}

func TestListPartiesHandler_InvalidStatus(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewListPartiesHandler(repo)
	if _, err := handler.Handle(context.Background(), &ListPartiesQuery{Status: "UNKNOWN"}); err == nil {
		t.Fatalf("Expected error for invalid status")
	}
}

func TestListPartiesHandler_InvalidRole(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewListPartiesHandler(repo)
	if _, err := handler.Handle(context.Background(), &ListPartiesQuery{Role: "BAD"}); err == nil {
		t.Fatalf("Expected error for invalid role")
	}
}

func TestListContactDetailsHandler_NoOrgProfile(t *testing.T) {
	repo := newFakePartyRepo()
	seedPartyWithProfiles(t, repo, "party-q3", false)

	handler := NewListContactDetailsHandler(repo)
	contacts, err := handler.Handle(context.Background(), &ListContactDetailsQuery{PartyID: "party-q3"})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if len(contacts) != 0 {
		t.Fatalf("Expected 0 contacts, got %d", len(contacts))
	}
}

func TestListPartyRelationshipsHandler_InvalidID(t *testing.T) {
	relRepo := newFakeRelationshipRepo()
	handler := NewListPartyRelationshipsHandler(relRepo)
	if _, err := handler.Handle(context.Background(), &ListPartyRelationshipsQuery{PartyID: ""}); err == nil {
		t.Fatalf("Expected error for invalid party ID")
	}
}
