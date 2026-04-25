package persistence

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDB provides database connection for integration tests
// Uses GORM with raw SQL for setup/teardown.
type TestDB struct {
	DB *gorm.DB
	t  *testing.T
}

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) *TestDB {
	config := loadTestDBConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Could not connect to PostgreSQL: %v. Skipping integration tests.", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Skipf("Could not get sql.DB: %v. Skipping integration tests.", err)
	}

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("Could not ping PostgreSQL: %v. Skipping integration tests.", err)
	}

	return &TestDB{
		DB: db,
		t:  t,
	}
}

type testDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func loadTestDBConfig() testDBConfig {
	// 1. Start with hardcoded defaults
	config := testDBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "tramatex",
		Password: "tramatex",
		Name:     "tramatex_test",
		SSLMode:  "disable",
	}

	// 2. Override with standard environment variables (e.g., from Docker Compose or CI)
	if value := os.Getenv("DB_HOST"); value != "" {
		config.Host = value
	}
	if value := os.Getenv("DB_PORT"); value != "" {
		config.Port = value
	}
	if value := os.Getenv("DB_USER"); value != "" {
		config.User = value
	}
	if value := os.Getenv("DB_PASSWORD"); value != "" {
		config.Password = value
	}
	if value := os.Getenv("DB_NAME"); value != "" {
		config.Name = value
	}
	if value := os.Getenv("DB_SSLMODE"); value != "" {
		config.SSLMode = value
	}

	// 3. Override with .env files if present (local development takes precedence over standard env)
	if env, err := readEnvLocal(); err == nil {
		applyEnvOverrides(&config, env, "localhost")
	} else if env, err := readEnvRemote(); err == nil {
		applyEnvOverrides(&config, env, "localhost")
	}

	// 4. Finally, override with specific test environment variables (highest priority)
	if value := os.Getenv("TRAMATEX_TEST_DB_HOST"); value != "" {
		config.Host = value
	}
	if value := os.Getenv("TRAMATEX_TEST_DB_PORT"); value != "" {
		config.Port = value
	}
	if value := os.Getenv("TRAMATEX_TEST_DB_USER"); value != "" {
		config.User = value
	}
	if value := os.Getenv("TRAMATEX_TEST_DB_PASSWORD"); value != "" {
		config.Password = value
	}
	if value := os.Getenv("TRAMATEX_TEST_DB_NAME"); value != "" {
		config.Name = value
	}
	if value := os.Getenv("TRAMATEX_TEST_DB_SSLMODE"); value != "" {
		config.SSLMode = value
	}

	return config
}

func applyEnvOverrides(config *testDBConfig, env map[string]string, fallbackHost string) {
	if value := env["DB_HOST"]; value != "" {
		config.Host = value
	}
	if value := env["DB_PORT"]; value != "" {
		config.Port = value
	}
	if value := env["DB_USER"]; value != "" {
		config.User = value
	}
	if value := env["DB_PASSWORD"]; value != "" {
		config.Password = value
	}
	if value := env["DB_NAME"]; value != "" {
		config.Name = value
	}
	if value := env["DB_SSLMODE"]; value != "" {
		config.SSLMode = value
	}
	if config.Host == "postgres" {
		if sshHost := env["SSH_HOST"]; sshHost != "" {
			config.Host = sshHost
		} else {
			config.Host = fallbackHost
		}
	}
}

func readEnvLocal() (map[string]string, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..", ".."))
	path := filepath.Join(root, ".env.local")
	return readEnvFile(path)
}

func readEnvRemote() (map[string]string, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..", ".."))
	path := filepath.Join(root, ".env.remote")
	return readEnvFile(path)
}

func readEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	env := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")
		env[key] = value
	}

	return env, nil
}

// SetUp initializes database schema for tests
func (tdb *TestDB) SetUp() error {
	ctx := context.Background()

	dropSchema := `
		DROP TABLE IF EXISTS addresses CASCADE;
		DROP TABLE IF EXISTS persons CASCADE;
		DROP TABLE IF EXISTS organizations CASCADE;
		DROP TYPE IF EXISTS organization_role CASCADE;
		DROP TYPE IF EXISTS organization_status CASCADE;
	`

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}

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
			created_by UUID NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_by UUID,
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
			created_by UUID NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_by UUID,
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
			created_by UUID NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_by UUID,
			modified_at TIMESTAMP
		);

		CREATE INDEX idx_addresses_organization_id ON addresses(organization_id);
		CREATE INDEX idx_addresses_is_primary ON addresses(is_primary);
	`

	if err := tdb.DB.WithContext(ctx).Exec(createSchema).Error; err != nil {
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

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}

	sqlDB, err := tdb.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// SetUpParty initializes Party schema for tests
func (tdb *TestDB) SetUpParty() error {
	ctx := context.Background()

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

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop party schema: %w", err)
	}

	createSchema := `
		CREATE TABLE users (
			id UUID PRIMARY KEY
		);

		INSERT INTO users (id) VALUES ('00000000-0000-0000-0000-000000000001') ON CONFLICT DO NOTHING;
		INSERT INTO users (id) VALUES ('00000000-0000-0000-0000-000000000002') ON CONFLICT DO NOTHING;

		CREATE TABLE parties (
			id VARCHAR(36) PRIMARY KEY,
			status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
			default_discount_percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
			created_by UUID NOT NULL REFERENCES users(id),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_by UUID NOT NULL REFERENCES users(id),
			modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE person_profiles (
			party_id VARCHAR(36) PRIMARY KEY REFERENCES parties(id) ON DELETE CASCADE,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			phone VARCHAR(30),
			email VARCHAR(255)
		);

		CREATE TABLE organization_profiles (
			party_id VARCHAR(36) PRIMARY KEY REFERENCES parties(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			tax_id VARCHAR(50),
			tax_id_type VARCHAR(20),
			website VARCHAR(255),
			phone VARCHAR(30),
			email VARCHAR(255),
			notes TEXT
		);

		CREATE TABLE party_roles (
			party_id VARCHAR(36) NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
			role VARCHAR(30) NOT NULL,
			creation_identifier VARCHAR(100),
			PRIMARY KEY (party_id, role)
		);

		CREATE TABLE party_relationships (
			id VARCHAR(36) PRIMARY KEY,
			from_party_id VARCHAR(36) NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
			to_party_id VARCHAR(36) NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
			type VARCHAR(50) NOT NULL,
			created_by UUID NOT NULL REFERENCES users(id),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_by UUID NOT NULL REFERENCES users(id),
			modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE contact_details (
			id VARCHAR(36) PRIMARY KEY,
			organization_party_id VARCHAR(36) NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
			type_description VARCHAR(100) NOT NULL,
			phone VARCHAR(30),
			email VARCHAR(255),
			related_party_id VARCHAR(36) REFERENCES parties(id) ON DELETE SET NULL
		);

		CREATE TABLE party_addresses (
			id VARCHAR(36) PRIMARY KEY,
			party_id VARCHAR(36) NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
			street VARCHAR(255) NOT NULL,
			city VARCHAR(100) NOT NULL,
			province VARCHAR(100),
			postal_code VARCHAR(20) NOT NULL,
			country VARCHAR(100) NOT NULL,
			is_primary BOOLEAN DEFAULT FALSE,
			created_by UUID NOT NULL REFERENCES users(id),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_by UUID NOT NULL REFERENCES users(id),
			modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	if err := tdb.DB.WithContext(ctx).Exec(createSchema).Error; err != nil {
		return fmt.Errorf("failed to create party schema: %w", err)
	}

	return nil
}

// TearDownParty cleans up Party schema
func (tdb *TestDB) TearDownParty() error {
	ctx := context.Background()

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

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop party schema: %w", err)
	}

	sqlDB, err := tdb.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
