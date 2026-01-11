package user_test

import (
	"testing"
	"time"

	"github.com/joran-cortez/tramatex/internal/domain/user"
)

func TestUserNewWithValidData(t *testing.T) {
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")

	u, err := user.NewUser("user-123", email, password, user.RoleOperator)

	if err != nil {
		t.Errorf("NewUser with valid data should not fail: %v", err)
	}

	if u == nil {
		t.Error("NewUser should return a user instance")
	}

	if u.ID() != "user-123" {
		t.Errorf("ID() = %q, want %q", u.ID(), "user-123")
	}

	if !u.Email().Equals(email) {
		t.Error("Email should match input")
	}

	if !u.IsActive() {
		t.Error("User should be active by default")
	}
}

func TestUserNewWithMissingEmail(t *testing.T) {
	password, _ := user.NewPassword("validPassword123")

	u, err := user.NewUser("user-123", nil, password, user.RoleOperator)

	if err == nil {
		t.Error("NewUser with nil email should fail")
	}

	if u != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithMissingPassword(t *testing.T) {
	email, _ := user.NewEmail("user@example.com")

	u, err := user.NewUser("user-123", email, nil, user.RoleOperator)

	if err == nil {
		t.Error("NewUser with nil password should fail")
	}

	if u != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithEmptyID(t *testing.T) {
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")

	u, err := user.NewUser("", email, password, user.RoleOperator)

	if err == nil {
		t.Error("NewUser with empty ID should fail")
	}

	if u != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithInvalidRole(t *testing.T) {
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")

	u, err := user.NewUser("user-123", email, password, user.Role("invalid"))

	if err == nil {
		t.Error("NewUser with invalid role should fail")
	}

	if u != nil {
		t.Error("NewUser should return nil on error")
	}
}

func TestUserNewWithValidRoles(t *testing.T) {
	roles := []user.Role{user.RoleAdmin, user.RoleManager, user.RoleOperator}
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			u, err := user.NewUser("user-123", email, password, role)

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
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")

	u, _ := user.NewUser("user-123", email, password, user.RoleOperator)

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

func TestUserTimestampsAutomatic(t *testing.T) {
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")

	before := time.Now()
	u, _ := user.NewUser("user-123", email, password, user.RoleOperator)
	after := time.Now()

	createdAt := u.CreatedAt()
	updatedAt := u.UpdatedAt()

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
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")

	u, _ := user.NewUser("user-123", email, password, user.RoleOperator)

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
	email, _ := user.NewEmail("user@example.com")
	oldPassword, _ := user.NewPassword("oldPassword123")

	u, _ := user.NewUser("user-123", email, oldPassword, user.RoleOperator)

	oldUpdatedAt := u.UpdatedAt()
	time.Sleep(10 * time.Millisecond) // Ensure time difference

	newPassword, _ := user.NewPassword("newPassword456")
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

	newUpdatedAt := u.UpdatedAt()
	if !newUpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated after ChangePassword")
	}
}

func TestUserChangePasswordWithNil(t *testing.T) {
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")

	u, _ := user.NewUser("user-123", email, password, user.RoleOperator)

	err := u.ChangePassword(nil)
	if err == nil {
		t.Error("ChangePassword with nil should fail")
	}
}

func TestRoleIsValid(t *testing.T) {
	tests := []struct {
		role     user.Role
		expected bool
	}{
		{user.RoleAdmin, true},
		{user.RoleManager, true},
		{user.RoleOperator, true},
		{user.Role("invalid"), false},
		{user.Role(""), false},
	}

	for _, tt := range tests {
		result := tt.role.IsValid()
		if result != tt.expected {
			t.Errorf("IsValid(%q) = %v, want %v", tt.role, result, tt.expected)
		}
	}
}

func TestNewUserWithUUID(t *testing.T) {
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")

	u, err := user.NewUserWithUUID(email, password, user.RoleOperator)

	if err != nil {
		t.Errorf("NewUserWithUUID failed: %v", err)
	}

	if u == nil {
		t.Error("NewUserWithUUID should return a user")
	}

	// UUID should be valid (non-empty, not "")
	if u.ID() == "" {
		t.Error("User ID should not be empty")
	}

	// UUID format check (simple: contains dashes)
	if len(u.ID()) != 36 {
		t.Errorf("UUID format unexpected: %q", u.ID())
	}
}
