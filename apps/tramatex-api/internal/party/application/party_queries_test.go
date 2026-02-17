package application

import (
	"context"
	"errors"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

type capturingPartyRepo struct {
	*fakePartyRepo
	lastFilters *persistence.PartyFilters
}

type failingPartyRepo struct {
	*fakePartyRepo
	findAllErr error
}

func (r *failingPartyRepo) FindAll(ctx context.Context, filters *persistence.PartyFilters) ([]*domain.Party, error) {
	return nil, r.findAllErr
}

type failingRelationshipRepo struct {
	findErr error
}

func (r *failingRelationshipRepo) Save(ctx context.Context, relationship domain.PartyRelationship, createdBy string, modifiedBy string) error {
	return nil
}

func (r *failingRelationshipRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]domain.PartyRelationship, error) {
	return nil, r.findErr
}

func (r *failingRelationshipRepo) Delete(ctx context.Context, id domain.PartyRelationshipID) error {
	return nil
}

type failingAddressRepo struct {
	findErr error
}

func (r *failingAddressRepo) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error {
	return nil
}

func (r *failingAddressRepo) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*domain.Address, error) {
	return nil, r.findErr
}

func (r *failingAddressRepo) FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error) {
	return nil, nil
}

func (r *failingAddressRepo) Delete(ctx context.Context, id domain.AddressID) error {
	return nil
}

func newCapturingPartyRepo() *capturingPartyRepo {
	return &capturingPartyRepo{fakePartyRepo: newFakePartyRepo()}
}

func (r *capturingPartyRepo) FindAll(ctx context.Context, filters *persistence.PartyFilters) ([]*domain.Party, error) {
	r.lastFilters = filters
	return r.fakePartyRepo.FindAll(ctx, filters)
}

func assertPartyErrorCode(t *testing.T, err error, code domain.ErrorCode) {
	t.Helper()
	var partyErr domain.PartyError
	if !errors.As(err, &partyErr) {
		t.Fatalf("expected PartyError, got %T", err)
	}
	if partyErr.Code != code {
		t.Fatalf("expected error code %s, got %s", code, partyErr.Code)
	}
}

func TestGetPartyHandler_Success(t *testing.T) {
	repo := newFakePartyRepo()
	partyID, _ := domain.NewPartyID("party-q1")
	personProfile, _ := domain.NewPersonProfile("Ana", "Perez")
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, personProfile, nil)
	_ = repo.Save(context.Background(), party, "user-1", "seed")

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

func TestGetPartyHandler_NotFound(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewGetPartyHandler(repo)

	_, err := handler.Handle(context.Background(), &GetPartyQuery{ID: "missing"})
	if err == nil {
		t.Fatalf("expected error for missing party")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeNotFound)
}

func TestGetPartyHandler_RepoError(t *testing.T) {
	repo := newErrorPartyRepo(nil, errors.New("db error"))
	handler := NewGetPartyHandler(repo)

	_, err := handler.Handle(context.Background(), &GetPartyQuery{ID: "party-err"})
	if err == nil {
		t.Fatalf("expected error for repository failure")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeNotFound)
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

func TestListPartiesHandler_Filters(t *testing.T) {
	repo := newCapturingPartyRepo()
	seedPartyWithProfiles(t, repo.fakePartyRepo, "party-q4", false)

	handler := NewListPartiesHandler(repo)
	_, err := handler.Handle(context.Background(), &ListPartiesQuery{
		Status:     "active",
		Role:       "client",
		Type:       "person",
		Name:       "Ana",
		TaxID:      "B123",
		PageNumber: 2,
		PageSize:   5,
	})
	if err != nil {
		t.Fatalf("Handle should not error: %v", err)
	}
	if repo.lastFilters == nil {
		t.Fatalf("expected filters to be set")
	}
	if repo.lastFilters.Type != "PERSON" {
		t.Fatalf("expected type PERSON, got %s", repo.lastFilters.Type)
	}
	if repo.lastFilters.Status == nil || *repo.lastFilters.Status != domain.PartyStatusActive {
		t.Fatalf("expected status ACTIVE")
	}
	if repo.lastFilters.Role == nil || *repo.lastFilters.Role != domain.PartyRoleClient {
		t.Fatalf("expected role CLIENT")
	}
}

func TestListPartiesHandler_FindAllError(t *testing.T) {
	repo := &failingPartyRepo{fakePartyRepo: newFakePartyRepo(), findAllErr: errors.New("db error")}
	handler := NewListPartiesHandler(repo)

	_, err := handler.Handle(context.Background(), &ListPartiesQuery{})
	if err == nil {
		t.Fatalf("expected error")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
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

func TestListPartyRelationshipsHandler_RepoError(t *testing.T) {
	relRepo := &failingRelationshipRepo{findErr: errors.New("db error")}
	handler := NewListPartyRelationshipsHandler(relRepo)

	_, err := handler.Handle(context.Background(), &ListPartyRelationshipsQuery{PartyID: "party-9"})
	if err == nil {
		t.Fatalf("expected error")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestListContactDetailsHandler_InvalidPartyID(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewListContactDetailsHandler(repo)

	_, err := handler.Handle(context.Background(), &ListContactDetailsQuery{PartyID: ""})
	if err == nil {
		t.Fatalf("expected error for empty party ID")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestListContactDetailsHandler_PartyNotFound(t *testing.T) {
	repo := newFakePartyRepo()
	handler := NewListContactDetailsHandler(repo)

	_, err := handler.Handle(context.Background(), &ListContactDetailsQuery{PartyID: "missing"})
	if err == nil {
		t.Fatalf("expected error for missing party")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeNotFound)
}

func TestListPartyAddressesHandler_InvalidPartyID(t *testing.T) {
	repo := newFakePartyAddressRepo()
	handler := NewListPartyAddressesHandler(repo)

	_, err := handler.Handle(context.Background(), &ListPartyAddressesQuery{PartyID: ""})
	if err == nil {
		t.Fatalf("expected error for empty party ID")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestListPartyAddressesHandler_RepoError(t *testing.T) {
	repo := &failingAddressRepo{findErr: errors.New("db error")}
	handler := NewListPartyAddressesHandler(repo)

	_, err := handler.Handle(context.Background(), &ListPartyAddressesQuery{PartyID: "party-10"})
	if err == nil {
		t.Fatalf("expected error")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}
