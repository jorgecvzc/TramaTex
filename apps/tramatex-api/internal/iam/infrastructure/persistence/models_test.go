package persistence

import "testing"

func TestUserModelTableName(t *testing.T) {
	var model UserModel
	if model.TableName() != "\"users\"" {
		t.Fatalf("unexpected table name: %s", model.TableName())
	}
}

func TestSeedAdminUser(t *testing.T) {
	seed := SeedAdminUser()
	if seed.Email == "" || seed.PasswordHash == "" {
		t.Fatalf("expected seed admin user to have email and password hash")
	}
	if seed.Role != "admin" {
		t.Fatalf("expected admin role, got %s", seed.Role)
	}
	if !seed.IsActive {
		t.Fatalf("expected admin user to be active")
	}
}
