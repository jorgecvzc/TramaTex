package model

import (
	"strings"
	"testing"
)

func TestEmailValidation(t *testing.T) {
	valid, err := NewEmail("test+tag@example.com")
	if err != nil {
		t.Fatalf("expected valid email, got error: %v", err)
	}
	if valid.Value() != "test+tag@example.com" {
		t.Fatalf("unexpected email value: %s", valid.Value())
	}
	if valid.String() != valid.Value() {
		t.Fatalf("expected String to return email value")
	}

	invalidCases := []string{
		"",
		"invalid",
		"missing@tld",
		strings.Repeat("a", 65) + "@example.com",
		strings.Repeat("a", 255) + "@example.com",
	}
	for _, value := range invalidCases {
		if _, err := NewEmail(value); err == nil {
			t.Fatalf("expected error for email %q", value)
		}
	}

	other, _ := NewEmail("other@example.com")
	if !valid.Equals(valid) {
		t.Fatalf("expected email to equal itself")
	}
	if valid.Equals(other) {
		t.Fatalf("expected different emails to not be equal")
	}
	if (*Email)(nil).Equals(nil) != true {
		t.Fatalf("expected nil equals nil")
	}
}

func TestPasswordValidationAndMatches(t *testing.T) {
	if _, err := NewPassword(""); err == nil {
		t.Fatalf("expected error for empty password")
	}
	if _, err := NewPassword("short"); err == nil {
		t.Fatalf("expected error for short password")
	}
	if _, err := NewPassword(strings.Repeat("a", MaxPasswordLength+1)); err == nil {
		t.Fatalf("expected error for long password")
	}

	password, err := NewPassword("strongpass")
	if err != nil {
		t.Fatalf("expected valid password, got error: %v", err)
	}
	if password.Hash() == "" {
		t.Fatalf("expected password hash to be set")
	}
	if !password.Matches("strongpass") {
		t.Fatalf("expected password to match")
	}
	if password.Matches("wrong") {
		t.Fatalf("expected password mismatch")
	}
	if password.String() != "[REDACTED]" {
		t.Fatalf("expected redacted password string")
	}

	fromHash := NewPasswordFromHash(password.Hash())
	if !fromHash.Matches("strongpass") {
		t.Fatalf("expected hash-based password to match")
	}
}

func TestRoleValidation(t *testing.T) {
	if !RoleAdmin.IsValid() {
		t.Fatalf("expected admin role to be valid")
	}
	if _, err := NewRole("admin"); err != nil {
		t.Fatalf("expected valid role, got error: %v", err)
	}
	if _, err := NewRole("unknown"); err == nil {
		t.Fatalf("expected error for invalid role")
	}
}

func TestUserLifecycle(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	password, _ := NewPassword("strongpass")

	if _, err := NewUser("", email, password, RoleAdmin); err == nil {
		t.Fatalf("expected error for empty id")
	}
	if _, err := NewUser("id", nil, password, RoleAdmin); err == nil {
		t.Fatalf("expected error for nil email")
	}
	if _, err := NewUser("id", email, nil, RoleAdmin); err == nil {
		t.Fatalf("expected error for nil password")
	}
	if _, err := NewUser("id", email, password, Role("invalid")); err == nil {
		t.Fatalf("expected error for invalid role")
	}

	user, err := NewUser("id", email, password, RoleCommercial)
	if err != nil {
		t.Fatalf("expected user creation success, got error: %v", err)
	}
	if !user.IsActive() {
		t.Fatalf("expected user to be active by default")
	}
	user.Deactivate()
	if user.IsActive() {
		t.Fatalf("expected user to be inactive")
	}
	user.Activate()
	if !user.IsActive() {
		t.Fatalf("expected user to be active")
	}
	if err := user.ChangeRole(RoleDesigner); err != nil {
		t.Fatalf("expected role change to succeed, got error: %v", err)
	}
	if err := user.ChangeRole(Role("invalid")); err == nil {
		t.Fatalf("expected error for invalid role change")
	}
	if err := user.ChangePassword(nil); err == nil {
		t.Fatalf("expected error for nil password")
	}
	newPassword, _ := NewPassword("anotherpass")
	if err := user.ChangePassword(newPassword); err != nil {
		t.Fatalf("expected password change to succeed")
	}
}

func TestPermissionValidation(t *testing.T) {
	if _, err := NewPermission("", "name"); err == nil {
		t.Fatalf("expected error for empty permission id")
	}
	if _, err := NewPermission("perm", ""); err == nil {
		t.Fatalf("expected error for empty permission name")
	}
	perm, err := NewPermission("perm", "Perm Name")
	if err != nil {
		t.Fatalf("expected permission creation success, got error: %v", err)
	}
	if perm.ID() != "perm" || perm.Name() != "Perm Name" {
		t.Fatalf("unexpected permission values")
	}
}
