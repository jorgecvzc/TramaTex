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
	config := testDBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Name:     "tramatex_test",
		SSLMode:  "disable",
	}

	if env, err := readEnvLocal(); err == nil {
		applyEnvOverrides(&config, env)
	}

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

// SetUpProduct initializes product schema for tests
func (tdb *TestDB) SetUpProduct() error {
	ctx := context.Background()

	// Drop tables if they exist (for clean state)
	dropSchema := `
		DROP TABLE IF EXISTS "parties" CASCADE;
		DROP TABLE IF EXISTS "party_service_configurations" CASCADE;
		DROP TABLE IF EXISTS "product_variant_values" CASCADE;
		DROP TABLE IF EXISTS "product_variants" CASCADE;
		DROP TABLE IF EXISTS "product_direct_attributes" CASCADE;
		DROP TABLE IF EXISTS "product_to_groups" CASCADE;
		DROP TABLE IF EXISTS "products" CASCADE;
		DROP TABLE IF EXISTS "attribute_values" CASCADE;
		DROP TABLE IF EXISTS "attributes" CASCADE;
		DROP TABLE IF EXISTS "product_groups" CASCADE;
		DROP TABLE IF EXISTS "brands" CASCADE;
	`

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop product schema: %w", err)
	}

	// Create enums and tables (same as migration)
	createSchema := `
		DO $$ BEGIN
			CREATE TYPE product_type AS ENUM ('TANGIBLE', 'SERVICE');
		EXCEPTION
			WHEN duplicate_object OR unique_violation THEN null;
		END $$;
		DO $$ BEGIN
			CREATE TYPE variant_status AS ENUM ('PROVISIONAL', 'CONFIRMED');
		EXCEPTION
			WHEN duplicate_object OR unique_violation THEN null;
		END $$;

		CREATE TABLE "brands" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"name" VARCHAR(255) NOT NULL,
			"is_active" BOOLEAN NOT NULL DEFAULT true,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE "product_groups" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"name" VARCHAR(255) NOT NULL,
			"parent_group_id" UUID,
			"is_active" BOOLEAN NOT NULL DEFAULT true,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE,
			CONSTRAINT "fk_parent_group" FOREIGN KEY ("parent_group_id") REFERENCES "product_groups" ("id") ON DELETE SET NULL
		);

		CREATE TABLE "attributes" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"name" VARCHAR(255) NOT NULL,
			"code" VARCHAR(50) NOT NULL,
			"sort_order" INT NOT NULL DEFAULT 0,
			"scope_brand_id" UUID,
			"scope_group_id" UUID,
			"created_by" VARCHAR(255),
			"modified_by" VARCHAR(255),
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE,
			CONSTRAINT "fk_scope_brand" FOREIGN KEY ("scope_brand_id") REFERENCES "brands" ("id") ON DELETE CASCADE,
			CONSTRAINT "fk_scope_product_group" FOREIGN KEY ("scope_group_id") REFERENCES "product_groups" ("id") ON DELETE CASCADE
		);

		CREATE TABLE "attribute_values" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"attribute_id" UUID NOT NULL,
			"value" VARCHAR(255) NOT NULL,
			"code" VARCHAR(50) NOT NULL,
			"created_by" VARCHAR(255),
			"modified_by" VARCHAR(255),
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE,
			CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES "attributes" ("id") ON DELETE CASCADE
		);

		CREATE TABLE "products" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"sku" VARCHAR(255) UNIQUE,
			"name" VARCHAR(100) NOT NULL,
			"long_name" VARCHAR(255),
			"barcode" VARCHAR(255) UNIQUE,
			"description" TEXT,
			"product_type" product_type NOT NULL,
			"brand_id" UUID NOT NULL,
			"group_ids" UUID[],
			"direct_attribute_ids" UUID[],
			"created_by" VARCHAR(255),
			"modified_by" VARCHAR(255),
			"is_active" BOOLEAN NOT NULL DEFAULT true,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE,
			CONSTRAINT "fk_products_brand" FOREIGN KEY ("brand_id") REFERENCES "brands" ("id") ON DELETE CASCADE
		);


		CREATE TABLE "product_variants" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"product_id" UUID NOT NULL,
			"sku" VARCHAR(255) NOT NULL,
			"barcode" VARCHAR(255),
			"base_cost" NUMERIC(12,2) NOT NULL DEFAULT 0,
			"status" variant_status NOT NULL DEFAULT 'PROVISIONAL',
			"attribute_values" UUID[],
			"created_by" VARCHAR(255),
			"modified_by" VARCHAR(255),
			"is_active" BOOLEAN NOT NULL DEFAULT true,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE,
			CONSTRAINT "fk_product_variants_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE
		);


		CREATE TABLE "parties" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid()
		);

		CREATE TABLE "party_service_configurations" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"party_id" UUID NOT NULL,
			"service_id" VARCHAR(255) NOT NULL,
			"name" VARCHAR(255) NOT NULL,
			"configuration_details" JSONB,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE,
			CONSTRAINT "fk_party_service_configurations_party" FOREIGN KEY ("party_id") REFERENCES "parties" ("id") ON DELETE CASCADE
		);
	`

	if err := tdb.DB.WithContext(ctx).Exec(createSchema).Error; err != nil {
		return fmt.Errorf("failed to create product schema: %w", err)
	}

	return nil
}

// TearDownProduct cleans up test database
func (tdb *TestDB) TearDownProduct() error {
	ctx := context.Background()

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
