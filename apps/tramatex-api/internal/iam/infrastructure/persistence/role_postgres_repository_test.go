package persistence

import (
	"context"
	"testing"
)

func TestPostgresRoleRepository(t *testing.T) {
	repo := NewPostgresRoleRepository(nil)

	role, err := repo.ByName(context.Background(), "admin")
	if err != nil {
		t.Fatalf("expected role lookup success, got error: %v", err)
	}
	if role != "admin" {
		t.Fatalf("unexpected role: %s", role)
	}

	if _, err := repo.ByID(context.Background(), "invalid"); err == nil {
		t.Fatalf("expected error for invalid role")
	}

	if err := repo.Save(context.Background(), "admin"); err != nil {
		t.Fatalf("expected save success, got error: %v", err)
	}
	if err := repo.Save(context.Background(), "invalid"); err == nil {
		t.Fatalf("expected error for invalid role save")
	}
}
