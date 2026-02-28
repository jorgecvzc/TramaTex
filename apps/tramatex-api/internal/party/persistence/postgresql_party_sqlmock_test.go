package persistence

import (
	"context"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/joran-cortez/tramatex/internal/party/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newPartySqlMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create gorm db: %v", err)
	}
	return gormDB, mock
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

func TestPostgreSQLPartyRepository_Save_NilParty(t *testing.T) {
	db, _ := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	err := repo.Save(context.Background(), nil, "actor", "actor")
	if err == nil {
		t.Fatalf("expected error for nil party")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestGORMPartyRepository_Save_WithPersonProfile(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	partyID, _ := domain.NewPartyID("party-1")
	personProfile, _ := domain.NewPersonProfile("Ana", "Perez", nil, nil)
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, personProfile, nil)

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "created_by"}).AddRow("party-1", createdAt, "seed")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FROM \"parties\"").WithArgs("party-1").WillReturnRows(rows)
	mock.ExpectExec("UPDATE .*\"parties\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE .*\"person_profiles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM \"contact_details\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM \"organization_profiles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM \"party_roles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Save(context.Background(), party, "user-1", "user-1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestGORMPartyRepository_Save_WithOrganizationProfileAndRole(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	partyID, _ := domain.NewPartyID("party-2")
	orgProfile, _ := domain.NewOrganizationProfile("Org", nil, "", nil, nil)
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, nil, orgProfile)
	role, _ := domain.NewPartyRole(domain.PartyRoleClient, nil)
	_ = party.AddRole(role)

	contactID, _ := domain.NewContactDetailsID("contact-2")
	contact, _ := domain.NewContactDetails(contactID, "Sales", nil, nil, nil)
	_ = orgProfile.AddContact(contact)

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "created_by"}).AddRow("party-2", createdAt, "seed")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FROM \"parties\"").WithArgs("party-2").WillReturnRows(rows)
	mock.ExpectExec("UPDATE .*\"parties\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM \"person_profiles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("(INSERT|UPDATE).*\"organization_profiles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM \"contact_details\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("(INSERT|UPDATE).*\"contact_details\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM \"party_roles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO \"party_roles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Save(context.Background(), party, "user-1", "user-1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestGORMPartyRepository_Save_WithContactDetailsFields(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	partyID, _ := domain.NewPartyID("party-3")
	orgProfile, _ := domain.NewOrganizationProfile("Org", nil, "", nil, nil)
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, nil, orgProfile)

	contactID, _ := domain.NewContactDetailsID("contact-3")
	phone, _ := domain.NewPhone("+34 600 111 222")
	email, _ := domain.NewEmail("ventas@org.local")
	relatedParty, _ := domain.NewPartyID("party-related")
	contact, _ := domain.NewContactDetails(contactID, "Ventas", phone, email, &relatedParty)
	_ = orgProfile.AddContact(contact)

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "created_by"}).AddRow("party-3", createdAt, "seed")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FROM \"parties\"").WithArgs("party-3").WillReturnRows(rows)
	mock.ExpectExec("UPDATE .*\"parties\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM \"person_profiles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("(INSERT|UPDATE).*\"organization_profiles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM \"contact_details\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("(INSERT|UPDATE).*\"contact_details\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM \"party_roles\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.Save(context.Background(), party, "user-1", "user-1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestPostgreSQLPartyRepository_FindByID_NotFound(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	rows := sqlmock.NewRows([]string{"id", "status"})
	mock.ExpectQuery("SELECT .* FROM \"parties\"").WithArgs("party-1").WillReturnRows(rows)

	partyID, _ := domain.NewPartyID("party-1")
	_, err := repo.FindByID(context.Background(), partyID)
	if err == nil {
		t.Fatalf("expected not found error")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeNotFound)
}

func TestPostgreSQLPartyRepository_LoadOrganizationProfile_WithContacts(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	partyID, _ := domain.NewPartyID("party-org")

	orgRows := sqlmock.NewRows([]string{"name", "tax_id", "tax_id_type", "website"}).
		AddRow("Org Name", "B12345678", "CIF", "https://org.local")
	mock.ExpectQuery("SELECT .* FROM \"organization_profiles\"").WithArgs(partyID.Value()).WillReturnRows(orgRows)

	contactRows := sqlmock.NewRows([]string{"id", "type_description", "phone", "email", "related_party_id"}).
		AddRow("contact-1", "Ventas", "+34 600 111 222", "ventas@org.local", "party-related")
	mock.ExpectQuery("SELECT .* FROM \"contact_details\"").WithArgs(partyID.Value()).WillReturnRows(contactRows)

	profile, err := repo.loadOrganizationProfile(context.Background(), partyID.Value())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if profile == nil || profile.Name() != "Org Name" {
		t.Fatalf("expected organization profile to be loaded")
	}
	if profile.TaxID() == nil || profile.TaxID().Value() != "B12345678" {
		t.Fatalf("expected tax ID to be loaded")
	}
	if len(profile.Contacts()) != 1 {
		t.Fatalf("expected 1 contact, got %d", len(profile.Contacts()))
	}
}

func TestPostgreSQLPartyRepository_LoadPartyRoles_InvalidRole(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	partyID, _ := domain.NewPartyID("party-roles")
	rows := sqlmock.NewRows([]string{"role"}).AddRow("BAD_ROLE")
	mock.ExpectQuery("SELECT .* FROM \"party_roles\"").WithArgs(partyID.Value()).WillReturnRows(rows)

	_, err := repo.loadPartyRoles(context.Background(), partyID.Value())
	if err == nil {
		t.Fatalf("expected error for invalid role")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestPostgreSQLPartyRepository_LoadContactDetails_InvalidEmail(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	partyID, _ := domain.NewPartyID("party-contact")
	rows := sqlmock.NewRows([]string{"id", "type_description", "phone", "email", "related_party_id"}).
		AddRow("contact-1", "Ventas", "+34 600 111 222", "bad-email", nil)
	mock.ExpectQuery("SELECT .* FROM \"contact_details\"").WithArgs(partyID.Value()).WillReturnRows(rows)

	_, err := repo.loadContactDetails(context.Background(), partyID.Value())
	if err == nil {
		t.Fatalf("expected error for invalid email")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestPostgreSQLPartyRepository_LoadContactDetails_WithEmailAndRelatedParty(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	rows := sqlmock.NewRows([]string{"id", "type_description", "phone", "email", "related_party_id"}).
		AddRow("contact-2", "Ventas", "+34 600 111 222", "ventas@org.local", "party-related")
	mock.ExpectQuery("SELECT .* FROM \"contact_details\"").WithArgs("party-2").WillReturnRows(rows)

	contacts, err := repo.loadContactDetails(context.Background(), "party-2")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(contacts) != 1 {
		t.Fatalf("expected 1 contact, got %d", len(contacts))
	}
	if contacts[0].Email() == nil || contacts[0].RelatedPartyID() == nil {
		t.Fatalf("expected email and related party to be set")
	}
}

func TestGORMPartyRepository_LoadPartyRoles_QueryError(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	mock.ExpectQuery("SELECT .* FROM \"party_roles\"").WithArgs("party-err").WillReturnError(errors.New("db error"))

	_, err := repo.loadPartyRoles(context.Background(), "party-err")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestPostgreSQLPartyAddressRepository_FindPrimary_NotFound(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyAddressRepository(db)

	partyID, _ := domain.NewPartyID("party-primary")
	rows := sqlmock.NewRows([]string{"id", "party_id", "street", "city", "province", "postal_code", "country", "is_primary", "created_by", "created_at", "modified_by", "modified_at"})
	mock.ExpectQuery("SELECT .* FROM \"party_addresses\"").WithArgs(partyID.Value()).WillReturnRows(rows)

	_, err := repo.FindPrimary(context.Background(), partyID)
	if err == nil {
		t.Fatalf("expected not found error")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeNotFound)
}

func TestPostgreSQLPartyRelationshipRepository_FindByPartyID_InvalidType(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRelationshipRepository(db)

	partyID, _ := domain.NewPartyID("party-rel")
	rows := sqlmock.NewRows([]string{"id", "from_party_id", "to_party_id", "type"}).
		AddRow("rel-1", "party-rel", "party-other", "BAD_TYPE")
	mock.ExpectQuery("SELECT .* FROM \"party_relationships\"").WithArgs(partyID.Value(), partyID.Value()).WillReturnRows(rows)

	_, err := repo.FindByPartyID(context.Background(), partyID)
	if err == nil {
		t.Fatalf("expected error for invalid relationship type")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestGORMPartyRepository_ExistsAndCount(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "parties" WHERE id = \$1`).WithArgs("party-1").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repo.Exists(context.Background(), domain.PartyID("party-1"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !exists {
		t.Fatalf("expected party to exist")
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "parties"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	count, err := repo.Count(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestGORMPartyRepository_FindByID_WithOrgAndRoles(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	partyID, _ := domain.NewPartyID("party-3")

	partyRows := sqlmock.NewRows([]string{"id", "status", "created_by", "created_at", "modified_by", "modified_at"}).
		AddRow("party-3", "ACTIVE", "seed", time.Now(), "seed", time.Now())
	mock.ExpectQuery("SELECT .* FROM \"parties\"").WithArgs("party-3").WillReturnRows(partyRows)

	personRows := sqlmock.NewRows([]string{"party_id", "first_name", "last_name"})
	mock.ExpectQuery("SELECT .* FROM \"person_profiles\"").WithArgs("party-3").WillReturnRows(personRows)

	orgRows := sqlmock.NewRows([]string{"party_id", "name", "tax_id", "tax_id_type", "website"}).
		AddRow("party-3", "Org", nil, nil, "")
	mock.ExpectQuery("SELECT .* FROM \"organization_profiles\"").WithArgs("party-3").WillReturnRows(orgRows)

	contactRows := sqlmock.NewRows([]string{"id", "organization_party_id", "type_description", "phone", "email", "related_party_id"}).
		AddRow("contact-3", "party-3", "Sales", nil, nil, nil)
	mock.ExpectQuery("SELECT .* FROM \"contact_details\"").WithArgs("party-3").WillReturnRows(contactRows)

	roleRows := sqlmock.NewRows([]string{"party_id", "role"}).AddRow("party-3", "CLIENT")
	mock.ExpectQuery("SELECT .* FROM \"party_roles\"").WithArgs("party-3").WillReturnRows(roleRows)

	party, err := repo.FindByID(context.Background(), partyID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if party.OrganizationProfile() == nil {
		t.Fatalf("expected org profile")
	}
	if len(party.Roles()) != 1 {
		t.Fatalf("expected 1 role, got %d", len(party.Roles()))
	}
}

func TestGORMPartyRepository_Delete(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"parties\"").WithArgs("party-1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.Delete(context.Background(), domain.PartyID("party-1")); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestGORMPartyRepository_Delete_Error(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"parties\"").WithArgs("party-err").WillReturnError(errors.New("db error"))
	mock.ExpectRollback()

	if err := repo.Delete(context.Background(), domain.PartyID("party-err")); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGORMPartyRepository_Exists_Error(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "parties" WHERE id = \$1`).WithArgs("party-err").WillReturnError(errors.New("db error"))

	_, err := repo.Exists(context.Background(), domain.PartyID("party-err"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestGORMPartyRepository_Count_Error(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "parties"`).WillReturnError(errors.New("db error"))

	_, err := repo.Count(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestGORMPartyRepository_Save_QueryError(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	partyID, _ := domain.NewPartyID("party-err")
	personProfile, _ := domain.NewPersonProfile("Ana", "Perez", nil, nil)
	party, _ := domain.NewParty(partyID, domain.PartyStatusActive, personProfile, nil)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FROM \"parties\"").WithArgs("party-err").WillReturnError(errors.New("db error"))
	mock.ExpectRollback()

	err := repo.Save(context.Background(), party, "user-1", "user-1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestGORMPartyRepository_FindByID_InvalidStatus(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	partyRows := sqlmock.NewRows([]string{"id", "status", "created_by", "created_at", "modified_by", "modified_at"}).
		AddRow("party-bad", "BAD", "seed", time.Now(), "seed", time.Now())
	mock.ExpectQuery("SELECT .* FROM \"parties\"").WithArgs("party-bad").WillReturnRows(partyRows)

	personRows := sqlmock.NewRows([]string{"party_id", "first_name", "last_name"})
	mock.ExpectQuery("SELECT .* FROM \"person_profiles\"").WithArgs("party-bad").WillReturnRows(personRows)

	orgRows := sqlmock.NewRows([]string{"party_id", "name", "tax_id", "tax_id_type", "website"})
	mock.ExpectQuery("SELECT .* FROM \"organization_profiles\"").WithArgs("party-bad").WillReturnRows(orgRows)

	roleRows := sqlmock.NewRows([]string{"party_id", "role"})
	mock.ExpectQuery("SELECT .* FROM \"party_roles\"").WithArgs("party-bad").WillReturnRows(roleRows)

	partyID, _ := domain.NewPartyID("party-bad")
	_, err := repo.FindByID(context.Background(), partyID)
	if err == nil {
		t.Fatalf("expected error for invalid status")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestGORMPartyRepository_FindAll_InvalidID(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow("")
	mock.ExpectQuery("SELECT DISTINCT .* FROM \"parties\"").WillReturnRows(rows)

	_, err := repo.FindAll(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error for invalid party id")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestGORMPartyRepository_LoadPersonProfile_QueryError(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	mock.ExpectQuery("SELECT .* FROM \"person_profiles\"").WithArgs("party-1").WillReturnError(errors.New("db error"))

	_, err := repo.loadPersonProfile(context.Background(), "party-1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestGORMPartyRepository_LoadOrganizationProfile_InvalidTaxID(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	orgRows := sqlmock.NewRows([]string{"party_id", "name", "tax_id", "tax_id_type", "website"}).
		AddRow("party-1", "Org", "BAD!", "NIF", "")
	mock.ExpectQuery("SELECT .* FROM \"organization_profiles\"").WithArgs("party-1").WillReturnRows(orgRows)

	_, err := repo.loadOrganizationProfile(context.Background(), "party-1")
	if err == nil {
		t.Fatalf("expected error for invalid tax id")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestGORMPartyRepository_LoadContactDetails_InvalidPhone(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	rows := sqlmock.NewRows([]string{"id", "type_description", "phone", "email", "related_party_id"}).
		AddRow("contact-1", "Ventas", "BAD", nil, nil)
	mock.ExpectQuery("SELECT .* FROM \"contact_details\"").WithArgs("party-1").WillReturnRows(rows)

	_, err := repo.loadContactDetails(context.Background(), "party-1")
	if err == nil {
		t.Fatalf("expected error for invalid phone")
	}
	assertPartyErrorCode(t, err, domain.ErrCodeValidation)
}

func TestGORMPartyRelationshipRepository_Save(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRelationshipRepository(db)

	fromID, _ := domain.NewPartyID("party-a")
	toID, _ := domain.NewPartyID("party-b")
	relID, _ := domain.NewPartyRelationshipID("rel-1")
	rel, _ := domain.NewPartyRelationship(relID, fromID, toID, domain.RelationshipIsEmployeeOf)

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "created_by"}).AddRow("rel-1", createdAt, "seed")

	mock.ExpectQuery("SELECT .* FROM \"party_relationships\"").WithArgs("rel-1").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE .*\"party_relationships\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.Save(context.Background(), rel, "user-1", "user-1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestGORMPartyAddressRepository_SaveAndDelete(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyAddressRepository(db)

	partyID, _ := domain.NewPartyID("party-addr")
	addressID, _ := domain.NewAddressID("addr-1")
	address, _ := domain.NewAddress("Street", "City", "Province", "28001", "Spain")

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "created_by"}).AddRow("addr-1", createdAt, "seed")

	mock.ExpectQuery("SELECT .* FROM \"party_addresses\"").WithArgs("addr-1").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE .*\"party_addresses\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.Save(context.Background(), address, addressID, partyID, "user-1", "user-1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"party_addresses\"").WithArgs("addr-1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.Delete(context.Background(), addressID); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func expectFindPartyMinimal(mock sqlmock.Sqlmock, partyID string, status string, withPerson bool) {
	partyRows := sqlmock.NewRows([]string{"id", "status", "created_by", "created_at", "modified_by", "modified_at"}).
		AddRow(partyID, status, "seed", time.Now(), "seed", time.Now())
	mock.ExpectQuery("SELECT .* FROM \"parties\"").WithArgs(partyID).WillReturnRows(partyRows)

	personRows := sqlmock.NewRows([]string{"party_id", "first_name", "last_name"})
	if withPerson {
		personRows.AddRow(partyID, "Ana", "Perez")
	}
	mock.ExpectQuery("SELECT .* FROM \"person_profiles\"").WithArgs(partyID).WillReturnRows(personRows)

	orgRows := sqlmock.NewRows([]string{"party_id", "name", "tax_id", "tax_id_type", "website"})
	mock.ExpectQuery("SELECT .* FROM \"organization_profiles\"").WithArgs(partyID).WillReturnRows(orgRows)

	roleRows := sqlmock.NewRows([]string{"party_id", "role"})
	mock.ExpectQuery("SELECT .* FROM \"party_roles\"").WithArgs(partyID).WillReturnRows(roleRows)
}

func TestGORMPartyRepository_FindAll_WithStatusFilter(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow("party-filtered")
	mock.ExpectQuery("SELECT DISTINCT .* FROM \"parties\"").WithArgs("ACTIVE").WillReturnRows(rows)

	expectFindPartyMinimal(mock, "party-filtered", "ACTIVE", true)

	filters := &PartyFilters{Status: func() *domain.PartyStatus { s := domain.PartyStatusActive; return &s }()}
	parties, err := repo.FindAll(context.Background(), filters)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(parties) != 1 {
		t.Fatalf("expected 1 party, got %d", len(parties))
	}
}

func TestGORMPartyRepository_FindAll_QueryError(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRepository(db)

	mock.ExpectQuery("SELECT DISTINCT .* FROM \"parties\"").WillReturnError(errors.New("db error"))

	_, err := repo.FindAll(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestGORMPartyRelationshipRepository_FindByPartyID_QueryError(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRelationshipRepository(db)

	partyID, _ := domain.NewPartyID("party-rel-err")
	mock.ExpectQuery("SELECT .* FROM \"party_relationships\"").WithArgs(partyID.Value(), partyID.Value()).WillReturnError(errors.New("db error"))

	_, err := repo.FindByPartyID(context.Background(), partyID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestGORMPartyRelationshipRepository_Save_QueryError(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRelationshipRepository(db)

	fromID, _ := domain.NewPartyID("party-a")
	toID, _ := domain.NewPartyID("party-b")
	relID, _ := domain.NewPartyRelationshipID("rel-err")
	rel, _ := domain.NewPartyRelationship(relID, fromID, toID, domain.RelationshipIsEmployeeOf)

	mock.ExpectQuery("SELECT .* FROM \"party_relationships\"").WithArgs("rel-err").WillReturnError(errors.New("db error"))

	if err := repo.Save(context.Background(), rel, "user-1", "user-1"); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGORMPartyRelationshipRepository_Delete_Error(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyRelationshipRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"party_relationships\"").WithArgs("rel-err").WillReturnError(errors.New("db error"))
	mock.ExpectRollback()

	if err := repo.Delete(context.Background(), domain.PartyRelationshipID("rel-err")); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGORMPartyAddressRepository_Save_NilAddress(t *testing.T) {
	db, _ := newPartySqlMock(t)
	repo := NewGORMPartyAddressRepository(db)

	partyID, _ := domain.NewPartyID("party-addr")
	addressID, _ := domain.NewAddressID("addr-nil")
	if err := repo.Save(context.Background(), nil, addressID, partyID, "user-1", "user-1"); err == nil {
		t.Fatalf("expected error for nil address")
	}
}

func TestGORMPartyAddressRepository_Save_QueryError(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyAddressRepository(db)

	partyID, _ := domain.NewPartyID("party-addr")
	addressID, _ := domain.NewAddressID("addr-err")
	address, _ := domain.NewAddress("Street", "City", "Province", "28001", "Spain")

	mock.ExpectQuery("SELECT .* FROM \"party_addresses\"").WithArgs("addr-err").WillReturnError(errors.New("db error"))

	if err := repo.Save(context.Background(), address, addressID, partyID, "user-1", "user-1"); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGORMPartyAddressRepository_FindByPartyID_QueryError(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyAddressRepository(db)

	partyID, _ := domain.NewPartyID("party-addr")
	mock.ExpectQuery("SELECT .* FROM \"party_addresses\"").WithArgs(partyID.Value()).WillReturnError(errors.New("db error"))

	_, err := repo.FindByPartyID(context.Background(), partyID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestGORMPartyAddressRepository_FindPrimary_QueryError(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyAddressRepository(db)

	partyID, _ := domain.NewPartyID("party-primary")
	mock.ExpectQuery("SELECT .* FROM \"party_addresses\"").WithArgs(partyID.Value()).WillReturnError(errors.New("db error"))

	_, err := repo.FindPrimary(context.Background(), partyID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertPartyErrorCode(t, err, domain.ErrCodePersistence)
}

func TestGORMPartyAddressRepository_Delete_Error(t *testing.T) {
	db, mock := newPartySqlMock(t)
	repo := NewGORMPartyAddressRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"party_addresses\"").WithArgs("addr-err").WillReturnError(errors.New("db error"))
	mock.ExpectRollback()

	if err := repo.Delete(context.Background(), domain.AddressID("addr-err")); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
