package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

// TestDB provides database connection for integration tests
type TestDB struct {
	DB *sql.DB
	t  *testing.T
}

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) *TestDB {
	// Connection string for local PostgreSQL
	// Make sure PostgreSQL is running and accessible
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=tramatex_test sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("Could not connect to PostgreSQL: %v. Skipping integration tests.", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		t.Skipf("Could not ping PostgreSQL: %v. Skipping integration tests.", err)
	}

	return &TestDB{
		DB: db,
		t:  t,
	}
}

// SetUp initializes database schema for tests
func (tdb *TestDB) SetUp() error {
	ctx := context.Background()

	// Drop tables if they exist (for clean state)
	dropSchema := `
		DROP TABLE IF EXISTS addresses CASCADE;
		DROP TABLE IF EXISTS persons CASCADE;
		DROP TABLE IF EXISTS organizations CASCADE;
		DROP TYPE IF EXISTS organization_role CASCADE;
		DROP TYPE IF EXISTS organization_status CASCADE;
	`

	if _, err := tdb.DB.ExecContext(ctx, dropSchema); err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}

	// Create enums and tables (same as migration)
	createSchema := `
		CREATE TYPE organization_role AS ENUM ('CLIENT', 'SUPPLIER', 'BOTH');
		CREATE TYPE organization_status AS ENUM ('ACTIVE', 'INACTIVE');

		CREATE TABLE organizations (
			id VARCHAR(100) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			role organization_role NOT NULL,
			status organization_status DEFAULT 'ACTIVE',
			tax_id VARCHAR(50),
			website VARCHAR(255),
			notes TEXT,
			created_by VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_by VARCHAR(100),
			modified_at TIMESTAMP,
			UNIQUE(tax_id)
		);

		CREATE INDEX idx_organizations_role ON organizations(role);
		CREATE INDEX idx_organizations_status ON organizations(status);
		CREATE INDEX idx_organizations_tax_id ON organizations(tax_id);

		CREATE TABLE persons (
			id VARCHAR(100) PRIMARY KEY,
			organization_id VARCHAR(100) NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(20),
			job_title VARCHAR(100),
			is_primary_contact BOOLEAN DEFAULT FALSE,
			created_by VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_by VARCHAR(100),
			modified_at TIMESTAMP
		);

		CREATE INDEX idx_persons_organization_id ON persons(organization_id);
		CREATE INDEX idx_persons_email ON persons(email);
		CREATE INDEX idx_persons_is_primary_contact ON persons(is_primary_contact);

		CREATE TABLE addresses (
			id VARCHAR(100) PRIMARY KEY,
			organization_id VARCHAR(100) NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			street VARCHAR(255) NOT NULL,
			city VARCHAR(100) NOT NULL,
			province VARCHAR(100),
			postal_code VARCHAR(20),
			country VARCHAR(100) DEFAULT 'Spain',
			is_primary BOOLEAN DEFAULT FALSE,
			created_by VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_by VARCHAR(100),
			modified_at TIMESTAMP
		);

		CREATE INDEX idx_addresses_organization_id ON addresses(organization_id);
		CREATE INDEX idx_addresses_is_primary ON addresses(is_primary);
	`

	if _, err := tdb.DB.ExecContext(ctx, createSchema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// TearDown cleans up test database
func (tdb *TestDB) TearDown() error {
	ctx := context.Background()

	dropSchema := `
		DROP TABLE IF EXISTS addresses CASCADE;
		DROP TABLE IF EXISTS persons CASCADE;
		DROP TABLE IF EXISTS organizations CASCADE;
		DROP TYPE IF EXISTS organization_role CASCADE;
		DROP TYPE IF EXISTS organization_status CASCADE;
	`

	if _, err := tdb.DB.ExecContext(ctx, dropSchema); err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}

	return tdb.DB.Close()
}
