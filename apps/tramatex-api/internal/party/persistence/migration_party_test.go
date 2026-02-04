package persistence

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPartyMigration_Integration(t *testing.T) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUp(); err != nil {
		t.Fatalf("Failed to set up legacy schema: %v", err)
	}

	ctx := context.Background()
	defer func() {
		_ = dropPartySchema(ctx, tdb.DB)
		if err := tdb.TearDown(); err != nil {
			t.Logf("Failed to tear down legacy schema: %v", err)
		}
	}()

	if _, err := tdb.DB.ExecContext(ctx, "ALTER TABLE organizations ADD COLUMN IF NOT EXISTS tax_id_type VARCHAR(20)"); err != nil {
		t.Fatalf("Failed to add tax_id_type column: %v", err)
	}

	if err := createUsersForMigration(ctx, tdb.DB); err != nil {
		t.Fatalf("Failed to create users: %v", err)
	}

	if err := seedV1Data(ctx, tdb.DB); err != nil {
		t.Fatalf("Failed to seed v1 data: %v", err)
	}

	if err := execMigrationFile(ctx, tdb.DB, "007_create_party_tables.sql"); err != nil {
		t.Fatalf("Failed to run migration 007: %v", err)
	}
	if err := execMigrationFile(ctx, tdb.DB, "008_migrate_party_data.sql"); err != nil {
		t.Fatalf("Failed to run migration 008: %v", err)
	}

	assertMigrationResults(t, ctx, tdb.DB)
}

func createUsersForMigration(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY)"); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, "INSERT INTO users (id) VALUES ('00000000-0000-0000-0000-000000000011'), ('00000000-0000-0000-0000-000000000012') ON CONFLICT DO NOTHING")
	return err
}

func seedV1Data(ctx context.Context, db *sql.DB) error {
	orgSQL := `
		INSERT INTO organizations (id, name, role, status, tax_id, tax_id_type, website, notes, created_by, created_at, modified_by, modified_at)
		VALUES ('org-001', 'Textiles Perez', 'BOTH', 'ACTIVE', 'B12345678', 'CIF', 'https://textiles.local', 'Notas', '00000000-0000-0000-0000-000000000011', CURRENT_TIMESTAMP, '00000000-0000-0000-0000-000000000011', CURRENT_TIMESTAMP)
	`
	if _, err := db.ExecContext(ctx, orgSQL); err != nil {
		return err
	}

	personSQL := `
		INSERT INTO persons (id, organization_id, first_name, last_name, email, phone, job_title, is_primary_contact, created_by, created_at, modified_by, modified_at)
		VALUES ('person-001', 'org-001', 'Ana', 'Perez', 'ana@textiles.local', '+34 600 000 111', 'Ventas', true, '00000000-0000-0000-0000-000000000011', CURRENT_TIMESTAMP, '00000000-0000-0000-0000-000000000012', CURRENT_TIMESTAMP)
	`
	if _, err := db.ExecContext(ctx, personSQL); err != nil {
		return err
	}

	addressSQL := `
		INSERT INTO addresses (id, organization_id, street, city, province, postal_code, country, is_primary, created_by, created_at, modified_by, modified_at)
		VALUES ('addr-001', 'org-001', 'Calle Principal 123', 'Madrid', 'Madrid', '28001', 'Spain', true, '00000000-0000-0000-0000-000000000011', CURRENT_TIMESTAMP, '00000000-0000-0000-0000-000000000011', CURRENT_TIMESTAMP)
	`
	_, err := db.ExecContext(ctx, addressSQL)
	return err
}

func execMigrationFile(ctx context.Context, db *sql.DB, filename string) error {
	sqlText, err := readMigrationFile(filename)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, sqlText)
	return err
}

func readMigrationFile(filename string) (string, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "migrations"))
	data, err := os.ReadFile(filepath.Join(migrationsDir, filename))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func dropPartySchema(ctx context.Context, db *sql.DB) error {
	dropSchema := `
        DROP TABLE IF EXISTS party_addresses CASCADE;
        DROP TABLE IF EXISTS contact_details CASCADE;
        DROP TABLE IF EXISTS party_relationships CASCADE;
        DROP TABLE IF EXISTS party_roles CASCADE;
        DROP TABLE IF EXISTS organization_profiles CASCADE;
        DROP TABLE IF EXISTS person_profiles CASCADE;
        DROP TABLE IF EXISTS parties CASCADE;
        DROP TABLE IF EXISTS users CASCADE;
    `
	_, err := db.ExecContext(ctx, dropSchema)
	return err
}

func assertMigrationResults(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()

	var partyCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM parties").Scan(&partyCount); err != nil {
		t.Fatalf("Failed to count parties: %v", err)
	}
	if partyCount != 2 {
		t.Fatalf("Expected 2 parties, got %d", partyCount)
	}

	var rolesCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM party_roles WHERE party_id = 'org-001'").Scan(&rolesCount); err != nil {
		t.Fatalf("Failed to count org roles: %v", err)
	}
	if rolesCount != 2 {
		t.Fatalf("Expected 2 roles for org, got %d", rolesCount)
	}

	var employeeRoleCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM party_roles WHERE party_id = 'person-001' AND role = 'EMPLOYEE'").Scan(&employeeRoleCount); err != nil {
		t.Fatalf("Failed to count employee roles: %v", err)
	}
	if employeeRoleCount != 1 {
		t.Fatalf("Expected EMPLOYEE role for person, got %d", employeeRoleCount)
	}

	var relType string
	if err := db.QueryRowContext(ctx, "SELECT type FROM party_relationships WHERE from_party_id = 'person-001' AND to_party_id = 'org-001'").Scan(&relType); err != nil {
		t.Fatalf("Failed to fetch relationship: %v", err)
	}
	if relType != "IS_EMPLOYEE_OF" {
		t.Fatalf("Expected IS_EMPLOYEE_OF relationship, got %s", relType)
	}

	var contactType, relatedParty string
	if err := db.QueryRowContext(ctx, "SELECT type_description, related_party_id FROM contact_details WHERE organization_party_id = 'org-001'").Scan(&contactType, &relatedParty); err != nil {
		t.Fatalf("Failed to fetch contact details: %v", err)
	}
	if contactType != "Ventas" || relatedParty != "person-001" {
		t.Fatalf("Unexpected contact details (type=%s, related=%s)", contactType, relatedParty)
	}

	var addressCount int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM party_addresses WHERE party_id = 'org-001'").Scan(&addressCount); err != nil {
		t.Fatalf("Failed to count party addresses: %v", err)
	}
	if addressCount != 1 {
		t.Fatalf("Expected 1 party address, got %d", addressCount)
	}
}
