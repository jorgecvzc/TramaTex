package persistence

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var productTestDBLock sync.Mutex

// LockProductTestDB serializes product DB setup/teardown across packages.
func LockProductTestDB() {
	productTestDBLock.Lock()
}

// UnlockProductTestDB releases the product test DB lock.
func UnlockProductTestDB() {
	productTestDBLock.Unlock()
}

// TestDB provides database connection for integration tests
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
		applyEnvOverrides(&config, env)
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

func applyEnvOverrides(config *testDBConfig, env map[string]string) {
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
		config.Host = "localhost"
	}
}

func readEnvLocal() (map[string]string, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..", "..", ".."))
	path := filepath.Join(root, ".env.local")
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

// Logf logs a formatted message using the underlying test logger.
func (tdb *TestDB) Logf(format string, args ...interface{}) {
	if tdb != nil && tdb.t != nil {
		tdb.t.Logf(format, args...)
	}
}

// SetUpProduct initializes product schema for tests using AutoMigrate
func (tdb *TestDB) SetUpProduct() error {
	// Create PostgreSQL enum types required by GORM models (AutoMigrate does not create them)
	enumTypes := `
		DO $$ BEGIN CREATE TYPE product_type AS ENUM ('TANGIBLE', 'SERVICE'); EXCEPTION WHEN duplicate_object THEN null; END $$;
		DO $$ BEGIN CREATE TYPE product_group_type AS ENUM ('TANGIBLE', 'SERVICE'); EXCEPTION WHEN duplicate_object THEN null; END $$;
		DO $$ BEGIN CREATE TYPE variant_status AS ENUM ('PROVISIONAL', 'CONFIRMED'); EXCEPTION WHEN duplicate_object THEN null; END $$;
	`
	if err := tdb.DB.Exec(enumTypes).Error; err != nil {
		return fmt.Errorf("failed to create product enum types: %w", err)
	}

	// AutoMigrate all product models to ensure schema matches GORM expectations
	err := tdb.DB.AutoMigrate(
		&BrandDataModel{},
		&ProductGroupDataModel{},
		&AttributeDataModel{},
		&AttributeValueDataModel{},
		&ProductDataModel{},
		&VariantDataModel{},
		&PartyServiceConfigurationModel{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate product schema: %w", err)
	}

	// Create parties stub table (cross-module reference: party_service_configurations tests insert into parties)
	if err := tdb.DB.Exec(`CREATE TABLE IF NOT EXISTS "parties" ("id" UUID PRIMARY KEY)`).Error; err != nil {
		return fmt.Errorf("failed to create parties reference table: %w", err)
	}

	return nil
}

// TearDownProduct cleans up test database
func (tdb *TestDB) TearDownProduct() error {
	ctx := context.Background()

	// Drop tables if they exist
	dropSchema := `
		DROP TABLE IF EXISTS "party_service_configurations" CASCADE;
		DROP TABLE IF EXISTS "parties" CASCADE;
		DROP TABLE IF EXISTS "product_variants" CASCADE;
		DROP TABLE IF EXISTS "products" CASCADE;
		DROP TABLE IF EXISTS "attribute_values" CASCADE;
		DROP TABLE IF EXISTS "attributes" CASCADE;
		DROP TABLE IF EXISTS "product_groups" CASCADE;
		DROP TABLE IF EXISTS "brands" CASCADE;
	`

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop product schema: %w", err)
	}

	sqlDB, _ := tdb.DB.DB()
	return sqlDB.Close()
}
