package user

import (
	"testing"
	"time"
)

func TestUserNewWithValidData(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("validPassword123")

	user, err := NewUser("user-123", email, password, RoleOperator)

	if err != nil {
		t.Errorf("NewUser with valid data should not fail: %v", err)
	}

	if user == nil {
		t.Error("NewUser should return a user instance")
	}

	if user.ID() != "user-123" {
		t.Errorf("ID() = %q, want %q", user.ID(), "user-123")
	}

	if !user.Email().Equals(email) {
		t.Error("Email should match input")
	}

	if !user.IsActive() {
		t.Error("User should be active by default")
	}
}

func TestUserNewWithMissingEmail(t *testing.T) {
	password, _ := NewPassword("validPassword123")

	user, err := NewUser("user-123", nil, password, RoleOperator)

	if err == nil {
		t.Error("NewUser with nil email should fail")
	}

	if user != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithMissingPassword(t *testing.T) {
	email, _ := NewEmail("user@example.com")

	user, err := NewUser("user-123", email, nil, RoleOperator)

	if err == nil {
		t.Error("NewUser with nil password should fail")
	}

	if user != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithEmptyID(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("validPassword123")

	user, err := NewUser("", email, password, RoleOperator)

	if err == nil {
		t.Error("NewUser with empty ID should fail")
	}

	if user != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithInvalidRole(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("validPassword123")

	user, err := NewUser("user-123", email, password, Role("invalid"))

	if err == nil {
		t.Error("NewUser with invalid role should fail")
	}

	if user != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithValidRoles(t *testing.T) {
	roles := []Role{RoleAdmin, RoleManager, RoleOperator}
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("validPassword123")

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			user, err := NewUser("user-123", email, password, role)

			if err != nil {
				t.Errorf("NewUser with role %q should succeed: %v", role, err)
			}

			if user.Role() != role {
				t.Errorf("Role() = %q, want %q", user.Role(), role)
			}
		})
	}
}

func TestUserImmutableAfterCreation(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("validPassword123")

	user, _ := NewUser("user-123", email, password, RoleOperator)

	// Verify ID doesn't change
	id1 := user.ID()
	id2 := user.ID()
	if id1 != id2 {
		t.Error("User ID should not change")
	}

	// Verify Email immutability through reference
	email1 := user.Email()
	email2 := user.Email()
	if !email1.Equals(email2) {
		t.Error("User Email should not change")
	}
}

func TestUserTimestampsAutomatic(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("validPassword123")

	before := time.Now()
	user, _ := NewUser("user-123", email, password, RoleOperator)
	after := time.Now()

	createdAt := user.CreatedAt()
	updatedAt := user.UpdatedAt()

	if createdAt.Before(before) || createdAt.After(after) {
		t.Error("CreatedAt should be between before and after times")
	}

	if updatedAt.Before(before) || updatedAt.After(after) {
		t.Error("UpdatedAt should be between before and after times")
	}

	// Initially, created and updated should be equal (or very close)
	if createdAt.Sub(updatedAt).Abs() > time.Millisecond {
		t.Errorf("CreatedAt and UpdatedAt should be equal: diff=%v", createdAt.Sub(updatedAt))
	}
}

func TestUserActiveFlag(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("validPassword123")

	user, _ := NewUser("user-123", email, password, RoleOperator)

	// Default to active
	if !user.IsActive() {
		t.Error("User should be active by default")
	}

	// Deactivate
	user.Deactivate()
	if user.IsActive() {
		t.Error("User should be inactive after Deactivate()")
	}

	// Activate
	user.Activate()
	if !user.IsActive() {
		t.Error("User should be active after Activate()")
	}
}

func TestUserChangePassword(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	oldPassword, _ := NewPassword("oldPassword123")

	user, _ := NewUser("user-123", email, oldPassword, RoleOperator)

	oldUpdatedAt := user.UpdatedAt()
	time.Sleep(10 * time.Millisecond) // Ensure time difference

	newPassword, _ := NewPassword("newPassword456")
	err := user.ChangePassword(newPassword)

	if err != nil {
		t.Errorf("ChangePassword should succeed: %v", err)
	}

	if !user.Password().Matches("newPassword456") {
		t.Error("New password should match")
	}

	if user.Password().Matches("oldPassword123") {
		t.Error("Old password should not match after change")
	}

	newUpdatedAt := user.UpdatedAt()
	if !newUpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated after ChangePassword")
	}
}

func TestUserChangePasswordWithNil(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("validPassword123")

	user, _ := NewUser("user-123", email, password, RoleOperator)

	err := user.ChangePassword(nil)
	if err == nil {
		t.Error("ChangePassword with nil should fail")
	}
}

func TestRoleIsValid(t *testing.T) {
	tests := []struct {
		role     Role
		expected bool
	}{
		{RoleAdmin, true},
		{RoleManager, true},
		{RoleOperator, true},
		{Role("invalid"), false},
		{Role(""), false},
	}

	for _, tt := range tests {
		result := tt.role.IsValid()
		if result != tt.expected {
			t.Errorf("IsValid(%q) = %v, want %v", tt.role, result, tt.expected)
		}
	}
}

func TestNewUserWithUUID(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("validPassword123")

	user, err := NewUserWithUUID(email, password, RoleOperator)

	if err != nil {
		t.Errorf("NewUserWithUUID failed: %v", err)
	}

	if user == nil {
		t.Error("NewUserWithUUID should return a user")
	}

	// UUID should be valid (non-empty, not "")
	if user.ID() == "" {
		t.Error("User ID should not be empty")
	}

	// UUID format check (simple: contains dashes)
	if len(user.ID()) != 36 {
		t.Errorf("UUID format unexpected: %q", user.ID())
	}
}
