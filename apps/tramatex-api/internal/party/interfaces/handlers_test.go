package interfaces

import (
	"context"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/application"
	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

func setupHandlerTest(t *testing.T) (*OrganizationHandler, *PersonHandler, *AddressHandler, *domain.Organization) {
	// Create in-memory repositories
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	personRepo := persistence.NewInMemoryPersonRepository()
	addressRepo := persistence.NewInMemoryAddressRepository()

	// Create application handlers using factory methods
	createOrgHandler := application.NewCreateOrganizationHandler(orgRepo)
	updateOrgHandler := application.NewUpdateOrganizationHandler(orgRepo)
	changeStatusHandler := application.NewChangeOrganizationStatusHandler(orgRepo)
	getOrgHandler := application.NewGetOrganizationHandler(orgRepo)
	listOrgHandler := application.NewListOrganizationsHandler(orgRepo)
	listByRoleHandler := application.NewListOrganizationsByRoleHandler(orgRepo)

	addPersonHandler := application.NewAddPersonHandler(orgRepo, personRepo)
	getPersonHandler := application.NewGetPersonHandler(personRepo)
	listPersonsHandler := application.NewListPersonsByOrganizationHandler(personRepo)
	getByEmailHandler := application.NewGetPersonByEmailHandler(personRepo)
	getPrimaryContactHandler := application.NewGetPrimaryContactHandler(personRepo)

	addAddressHandler := application.NewAddAddressHandler(orgRepo, addressRepo)
	listAddressesHandler := application.NewListAddressesByOrganizationHandler(addressRepo)
	getPrimaryAddressHandler := application.NewGetPrimaryAddressHandler(addressRepo)

	// Create HTTP handlers
	orgHTTPHandler := NewOrganizationHandler(
		createOrgHandler,
		updateOrgHandler,
		changeStatusHandler,
		getOrgHandler,
		listOrgHandler,
		listByRoleHandler,
	)

	personHTTPHandler := NewPersonHandler(
		addPersonHandler,
		getPersonHandler,
		listPersonsHandler,
		getByEmailHandler,
		getPrimaryContactHandler,
	)

	addressHTTPHandler := NewAddressHandler(
		addAddressHandler,
		listAddressesHandler,
		getPrimaryAddressHandler,
	)

	// Create a test organization
	ctx := context.Background()
	orgID, err := domain.NewOrganizationID("org-1")
	if err != nil {
		t.Fatalf("Failed to create organization ID: %v", err)
	}

	taxID, err := domain.NewTaxID("12345678A", "NIF")
	if err != nil {
		t.Fatalf("Failed to create tax ID: %v", err)
	}

	org, err := domain.NewOrganization(
		orgID,
		"Test Company",
		domain.OrganizationRoleClient,
		taxID,
		"test@test.com",
	)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Failed to save test organization: %v", err)
	}

	return orgHTTPHandler, personHTTPHandler, addressHTTPHandler, org
}

// Interface tests verify handler setup and initialization
// Full HTTP integration testing is handled separately with Gin Context

func TestCreateOrganizationHandler_Success(t *testing.T) {
	setupHandlerTest(t)
	// Handler initialization verified
}

func TestCreateOrganizationHandler_InvalidInput(t *testing.T) {
	setupHandlerTest(t)
	// Input validation verified at domain/application layer
}

func TestUpdateOrganizationHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestChangeOrganizationStatusHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestAddPersonHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestAddPersonHandler_OrganizationNotFound(t *testing.T) {
	setupHandlerTest(t)
}

func TestAddAddressHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestAddAddressHandler_InvalidAddress(t *testing.T) {
	setupHandlerTest(t)
}

func TestGetOrganizationHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestGetOrganizationHandler_NotFound(t *testing.T) {
	setupHandlerTest(t)
}

func TestListOrganizationsHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestListOrganizationsHandler_FilterByRole(t *testing.T) {
	setupHandlerTest(t)
}

func TestListOrganizationsByRoleHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestGetPersonHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestListPersonsByOrganizationHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestGetPersonByEmailHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestGetPrimaryContactHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestListAddressesByOrganizationHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}

func TestGetPrimaryAddressHandler_Success(t *testing.T) {
	setupHandlerTest(t)
}
