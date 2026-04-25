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

type TestDB struct {
	DB *gorm.DB
	t  *testing.T
}

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

	if err := sqlDB.Ping(); err != nil {
		t.Skipf("Could not ping PostgreSQL: %v. Skipping integration tests.", err)
	}

	return &TestDB{DB: db, t: t}
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

// SetUpSales initializes database schema for sales tests using AutoMigrate
func (tdb *TestDB) SetUpSales() error {
	// Create PostgreSQL enum types required by GORM models (AutoMigrate does not create them)
	enumTypes := `
		DO $$ BEGIN CREATE TYPE quote_status AS ENUM ('DRAFT','ISSUED','APPROVED','ACCEPTED','REJECTED','EXPIRED','CONVERTED_TO_ORDER'); EXCEPTION WHEN duplicate_object THEN null; END $$;
		DO $$ BEGIN CREATE TYPE sales_order_status AS ENUM ('PENDING','IN_PREPARATION','READY_FOR_PRODUCTION','PARTIALLY_DELIVERED','DELIVERED','CANCELLED','PARTIALLY_INVOICED','INVOICED'); EXCEPTION WHEN duplicate_object THEN null; END $$;
		DO $$ BEGIN CREATE TYPE delivery_note_status AS ENUM ('PENDING','DELIVERED','CANCELLED'); EXCEPTION WHEN duplicate_object THEN null; END $$;
		DO $$ BEGIN CREATE TYPE invoice_status AS ENUM ('DRAFT','ISSUED','PAID','OVERDUE','VOID'); EXCEPTION WHEN duplicate_object THEN null; END $$;
		DO $$ BEGIN CREATE TYPE invoice_type AS ENUM ('COMPLETE','SIMPLIFIED'); EXCEPTION WHEN duplicate_object THEN null; END $$;
	`
	if err := tdb.DB.Exec(enumTypes).Error; err != nil {
		return fmt.Errorf("failed to create sales enum types: %w", err)
	}

	// AutoMigrate all sales models to ensure schema matches GORM expectations
	err := tdb.DB.AutoMigrate(
		&QuoteDataModel{},
		&QuoteLineItemDataModel{},
		&QuoteWorkRefModel{},
		&SalesOrderDataModel{},
		&OrderLineItemDataModel{},
		&OrderWorkRefModel{},
		&DeliveryNoteDataModel{},
		&DeliveryNoteLineItemDataModel{},
		&InvoiceDataModel{},
		&InvoiceLineItemDataModel{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate sales schema: %w", err)
	}

	// Create product_variants table if not exists (cross-module reference)
	if err := tdb.DB.Exec(`CREATE TABLE IF NOT EXISTS "product_variants" ("id" UUID PRIMARY KEY)`).Error; err != nil {
		return fmt.Errorf("failed to create product_variants reference table: %w", err)
	}

	return nil
}

// TearDownSales cleans up sales test database
func (tdb *TestDB) TearDownSales() error {
	ctx := context.Background()

	dropSchema := `
		DROP TABLE IF EXISTS "invoice_line_items" CASCADE;
		DROP TABLE IF EXISTS "invoices" CASCADE;
		DROP TABLE IF EXISTS "delivery_note_line_items" CASCADE;
		DROP TABLE IF EXISTS "delivery_notes" CASCADE;
		DROP TABLE IF EXISTS "order_line_items" CASCADE;
		DROP TABLE IF EXISTS "order_work_setups" CASCADE;
		DROP TABLE IF EXISTS "sales_orders" CASCADE;
		DROP TABLE IF EXISTS "quote_line_items" CASCADE;
		DROP TABLE IF EXISTS "quote_work_setups" CASCADE;
		DROP TABLE IF EXISTS "quotes" CASCADE;
		DROP TABLE IF EXISTS "product_variants" CASCADE;
	`

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop sales schema: %w", err)
	}

	return nil
}
