package model_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
)

var testUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func TestUserNewWithValidData(t *testing.T) {
	email, _ := model.NewEmail("user@example.com")
	password, _ := model.NewPassword("validPassword123")

	u, err := model.NewUser(testUserID, email, password, model.RoleWorkshop)

	if err != nil {
		t.Errorf("NewUser with valid data should not fail: %v", err)
	}

	if u == nil {
		t.Error("NewUser should return a user instance")
	}

	if u.ID() != testUserID {
		t.Errorf("ID() = %s, want %s", u.ID(), testUserID)
	}

	if !u.Email().Equals(email) {
		t.Error("Email should match input")
	}

	if !u.IsActive() {
		t.Error("User should be active by default")
	}
}

func TestUserNewWithMissingEmail(t *testing.T) {
	password, _ := model.NewPassword("validPassword123")

	u, err := model.NewUser(testUserID, nil, password, model.RoleWorkshop)

	if err == nil {
		t.Error("NewUser with nil email should fail")
	}

	if u != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithMissingPassword(t *testing.T) {
	email, _ := model.NewEmail("user@example.com")

	u, err := model.NewUser(testUserID, email, nil, model.RoleWorkshop)

	if err == nil {
		t.Error("NewUser with nil password should fail")
	}

	if u != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithEmptyID(t *testing.T) {
	email, _ := model.NewEmail("user@example.com")
	password, _ := model.NewPassword("validPassword123")

	u, err := model.NewUser(uuid.Nil, email, password, model.RoleWorkshop)

	if err == nil {
		t.Error("NewUser with empty ID should fail")
	}

	if u != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithInvalidRole(t *testing.T) {
	email, _ := model.NewEmail("user@example.com")
	password, _ := model.NewPassword("validPassword123")

	u, err := model.NewUser(testUserID, email, password, model.Role("invalid"))

	if err == nil {
		t.Error("NewUser with invalid role should fail")
	}

	if u != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithValidRoles(t *testing.T) {
	roles := []model.Role{model.RoleAdmin, model.RoleCommercial, model.RoleDesigner, model.RoleWorkshop}
	email, _ := model.NewEmail("user@example.com")
	password, _ := model.NewPassword("validPassword123")

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			u, err := model.NewUser(testUserID, email, password, role)

			if err != nil {
				t.Errorf("NewUser with role %q should succeed: %v", role, err)
			}

			if u.Role() != role {
				t.Errorf("Role() = %q, want %q", u.Role(), role)
			}
		})
	}
}

func TestUserImmutableAfterCreation(t *testing.T) {
	email, _ := model.NewEmail("user@example.com")
	password, _ := model.NewPassword("validPassword123")

	u, _ := model.NewUser(testUserID, email, password, model.RoleWorkshop)

	// Verify ID doesn't change
	id1 := u.ID()
	id2 := u.ID()
	if id1 != id2 {
		t.Error("User ID should not change")
	}

	// Verify Email immutability through reference
	email1 := u.Email()
	email2 := u.Email()
	if !email1.Equals(email2) {
		t.Error("User Email should not change")
	}
}

func TestUserActiveFlag(t *testing.T) {
	email, _ := model.NewEmail("user@example.com")
	password, _ := model.NewPassword("validPassword123")

	u, _ := model.NewUser(testUserID, email, password, model.RoleWorkshop)

	// Default to active
	if !u.IsActive() {
		t.Error("User should be active by default")
	}

	// Deactivate
	u.Deactivate()
	if u.IsActive() {
		t.Error("User should be inactive after Deactivate()")
	}

	// Activate
	u.Activate()
	if !u.IsActive() {
		t.Error("User should be active after Activate()")
	}
}

func TestUserChangePassword(t *testing.T) {
	email, _ := model.NewEmail("user@example.com")
	oldPassword, _ := model.NewPassword("oldPassword123")

	u, _ := model.NewUser(testUserID, email, oldPassword, model.RoleWorkshop)

	newPassword, _ := model.NewPassword("newPassword456")
	err := u.ChangePassword(newPassword)

	if err != nil {
		t.Errorf("ChangePassword should succeed: %v", err)
	}

	if !u.Password().Matches("newPassword456") {
		t.Error("New password should match")
	}

	if u.Password().Matches("oldPassword123") {
		t.Error("Old password should not match after change")
	}
}

func TestUserChangePasswordWithNil(t *testing.T) {
	email, _ := model.NewEmail("user@example.com")
	password, _ := model.NewPassword("validPassword123")

	u, _ := model.NewUser(testUserID, email, password, model.RoleWorkshop)

	err := u.ChangePassword(nil)
	if err == nil {
		t.Error("ChangePassword with nil should fail")
	}
}

func TestRoleIsValid(t *testing.T) {
	tests := []struct {
		role     model.Role
		expected bool
	}{
		{model.RoleAdmin, true},
		{model.RoleCommercial, true},
		{model.RoleDesigner, true},
		{model.RoleWorkshop, true},
		{model.Role("invalid"), false},
		{model.Role(""), false},
	}

	for _, tt := range tests {
		result := tt.role.IsValid()
		if result != tt.expected {
			t.Errorf("IsValid(%q) = %v, want %v", tt.role, result, tt.expected)
		}
	}
}

func TestNewUserWithUUID(t *testing.T) {
	email, _ := model.NewEmail("user@example.com")
	password, _ := model.NewPassword("validPassword123")

	u, err := model.NewUserWithUUID(email, password, model.RoleWorkshop)

	if err != nil {
		t.Errorf("NewUserWithUUID failed: %v", err)
	}

	if u == nil {
		t.Error("NewUserWithUUID should return a user")
	}

	// UUID should be valid (non-empty, not "")
	if u.ID() == uuid.Nil {
		t.Error("User ID should not be empty")
	}

	// UUID format check (simple: contains dashes)
	if len(u.ID().String()) != 36 {
		t.Errorf("UUID format unexpected: %q", u.ID())
	}
}
