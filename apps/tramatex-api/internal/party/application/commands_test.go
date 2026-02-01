package application

import (
	"context"
	"testing"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"github.com/joran-cortez/tramatex/internal/party/persistence"
)

func TestCreateOrganizationHandler_Success(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	handler := NewCreateOrganizationHandler(orgRepo)
	ctx := context.Background()

	cmd := &CreateOrganizationCommand{
		ID:        "org-001",
		Name:      "Test Corp",
		Role:      "CLIENT",
		TaxID:     "12345678A",
		TaxIDType: "NIF",
		Website:   "https://test.com",
		CreatedBy: "user-1",
	}

	org, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if org.Name() != "Test Corp" {
		t.Errorf("Organization name mismatch, got: %s", org.Name())
	}
	if org.Website() != "https://test.com" {
		t.Errorf("Organization website mismatch, got: %s", org.Website())
	}
}

func TestCreateOrganizationHandler_InvalidInput(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	handler := NewCreateOrganizationHandler(orgRepo)
	ctx := context.Background()

	tests := []struct {
		name    string
		cmd     *CreateOrganizationCommand
		wantErr bool
	}{
		{
			name:    "Empty ID",
			cmd:     &CreateOrganizationCommand{Name: "Test", Role: "CLIENT", CreatedBy: "user-1"},
			wantErr: true,
		},
		{
			name:    "Empty Name",
			cmd:     &CreateOrganizationCommand{ID: "org-001", Role: "CLIENT", CreatedBy: "user-1"},
			wantErr: true,
		},
		{
			name:    "Empty CreatedBy",
			cmd:     &CreateOrganizationCommand{ID: "org-001", Name: "Test", Role: "CLIENT"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Handle(ctx, tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle wantErr %v, got error: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUpdateOrganizationHandler_Success(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	ctx := context.Background()

	// Create initial organization
	orgID, _ := domain.NewOrganizationID("org-001")
	org, _ := domain.NewOrganization(orgID, "Original Name", domain.OrganizationRoleClient, nil, "user-1")
	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Save should not error, got: %v", err)
	}

	// Update it
	handler := NewUpdateOrganizationHandler(orgRepo)
	cmd := &UpdateOrganizationCommand{
		ID:         "org-001",
		Name:       "Updated Name",
		Website:    "https://updated.com",
		Notes:      "Test notes",
		ModifiedBy: "user-2",
	}

	updated, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if updated.Name() != "Updated Name" {
		t.Errorf("Organization name should be updated, got: %s", updated.Name())
	}
}

func TestChangeOrganizationStatusHandler_Success(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	ctx := context.Background()

	// Create initial organization
	orgID, _ := domain.NewOrganizationID("org-001")
	org, _ := domain.NewOrganization(orgID, "Test Corp", domain.OrganizationRoleClient, nil, "user-1")
	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Save should not error, got: %v", err)
	}

	// Change status
	handler := NewChangeOrganizationStatusHandler(orgRepo)
	cmd := &ChangeOrganizationStatusCommand{
		ID:         "org-001",
		Status:     "INACTIVE",
		ModifiedBy: "user-2",
	}

	updated, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if updated.Status() != domain.OrganizationStatusInactive {
		t.Errorf("Organization status should be INACTIVE, got: %v", updated.Status())
	}
}

func TestAddPersonHandler_Success(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	personRepo := persistence.NewInMemoryPersonRepository()
	ctx := context.Background()

	// Create organization first
	orgID, _ := domain.NewOrganizationID("org-001")
	org, _ := domain.NewOrganization(orgID, "Test Corp", domain.OrganizationRoleClient, nil, "user-1")
	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Save should not error, got: %v", err)
	}

	// Add person
	handler := NewAddPersonHandler(orgRepo, personRepo)
	cmd := &AddPersonCommand{
		ID:             "person-001",
		OrganizationID: "org-001",
		FirstName:      "John",
		LastName:       "Doe",
		Email:          "john@example.com",
		Phone:          "+34666123456",
		JobTitle:       "Manager",
		IsPrimary:      true,
		CreatedBy:      "user-1",
	}

	person, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if person.FullName() != "John Doe" {
		t.Errorf("Person full name mismatch, got: %s", person.FullName())
	}
	if !person.IsPrimaryContact() {
		t.Error("Person should be marked as primary contact")
	}
}

func TestAddPersonHandler_OrganizationNotFound(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	personRepo := persistence.NewInMemoryPersonRepository()
	ctx := context.Background()

	handler := NewAddPersonHandler(orgRepo, personRepo)
	cmd := &AddPersonCommand{
		ID:             "person-001",
		OrganizationID: "nonexistent",
		FirstName:      "John",
		LastName:       "Doe",
		Email:          "john@example.com",
		CreatedBy:      "user-1",
	}

	_, err := handler.Handle(ctx, cmd)
	if err == nil {
		t.Error("Handle should error for nonexistent organization")
	}
}

func TestAddAddressHandler_Success(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	addressRepo := persistence.NewInMemoryAddressRepository()
	ctx := context.Background()

	// Create organization first
	orgID, _ := domain.NewOrganizationID("org-001")
	org, _ := domain.NewOrganization(orgID, "Test Corp", domain.OrganizationRoleClient, nil, "user-1")
	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Save should not error, got: %v", err)
	}

	// Add address
	handler := NewAddAddressHandler(orgRepo, addressRepo)
	cmd := &AddAddressCommand{
		ID:             "addr-001",
		OrganizationID: "org-001",
		Street:         "Calle Principal 123",
		City:           "Madrid",
		Province:       "Madrid",
		PostalCode:     "28001",
		Country:        "Spain",
		IsPrimary:      true,
		CreatedBy:      "user-1",
	}

	address, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Errorf("Handle should not error, got: %v", err)
	}
	if address.City() != "Madrid" {
		t.Errorf("Address city mismatch, got: %s", address.City())
	}
}

func TestAddAddressHandler_InvalidAddress(t *testing.T) {
	orgRepo := persistence.NewInMemoryOrganizationRepository()
	addressRepo := persistence.NewInMemoryAddressRepository()
	ctx := context.Background()

	// Create organization
	orgID, _ := domain.NewOrganizationID("org-001")
	org, _ := domain.NewOrganization(orgID, "Test Corp", domain.OrganizationRoleClient, nil, "user-1")
	if err := orgRepo.Save(ctx, org); err != nil {
		t.Fatalf("Save should not error, got: %v", err)
	}

	handler := NewAddAddressHandler(orgRepo, addressRepo)
	cmd := &AddAddressCommand{
		ID:             "addr-001",
		OrganizationID: "org-001",
		Street:         "", // Invalid: empty street
		City:           "Madrid",
		CreatedBy:      "user-1",
	}

	_, err := handler.Handle(ctx, cmd)
	if err == nil {
		t.Error("Handle should error for invalid address")
	}
}
