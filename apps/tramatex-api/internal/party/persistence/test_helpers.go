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
	"gorm.io/gorm/logger"
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
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

// SetUp initializes database schema for tests (Legacy)
func (tdb *TestDB) SetUp() error {
	// For backward compatibility, but we should use SetUpParty
	return tdb.SetUpParty()
}

// TearDown cleans up test database (Legacy)
func (tdb *TestDB) TearDown() error {
	// For backward compatibility
	return tdb.TearDownParty()
}

// SetUpParty initializes Party schema for tests using AutoMigrate
func (tdb *TestDB) SetUpParty() error {
	// AutoMigrate all party models to ensure schema matches GORM expectations
	err := tdb.DB.AutoMigrate(
		&UserDataModel{},
		&PartyDataModel{},
		&PersonProfileDataModel{},
		&OrganizationProfileDataModel{},
		&PartyRoleDataModel{},
		&PartyRelationshipDataModel{},
		&ContactDetailDataModel{},
		&PartyAddressDataModel{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate party schema: %w", err)
	}

	// Seed required users for foreign keys
	users := []UserDataModel{
		{ID: "00000000-0000-0000-0000-000000000001"},
		{ID: "00000000-0000-0000-0000-000000000002"},
	}
	for _, u := range users {
		if err := tdb.DB.FirstOrCreate(&u, "id = ?", u.ID).Error; err != nil {
			return fmt.Errorf("failed to seed test user %s: %w", u.ID, err)
		}
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
