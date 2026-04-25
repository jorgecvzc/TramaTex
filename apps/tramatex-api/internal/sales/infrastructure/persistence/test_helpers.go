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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

func (tdb *TestDB) SetUpSales() error {
	ctx := context.Background()

	dropSchema := `
		DROP TABLE IF EXISTS "invoice_line_items" CASCADE;
		DROP TABLE IF EXISTS "invoices" CASCADE;
		DROP TABLE IF EXISTS "delivery_note_line_items" CASCADE;
		DROP TABLE IF EXISTS "delivery_notes" CASCADE;
		DROP TABLE IF EXISTS "order_line_items" CASCADE;
		DROP TABLE IF EXISTS "sales_orders" CASCADE;
		DROP TABLE IF EXISTS "quote_line_items" CASCADE;
		DROP TABLE IF EXISTS "quotes" CASCADE;
		DROP TABLE IF EXISTS "product_variants" CASCADE;
		DROP TYPE IF EXISTS quote_status;
		DROP TYPE IF EXISTS sales_order_status;
		DROP TYPE IF EXISTS delivery_note_status;
		DROP TYPE IF EXISTS invoice_status;
		DROP TYPE IF EXISTS invoice_type;
	`

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop sales schema: %w", err)
	}

	createSchema := `
		DO $$ BEGIN
			CREATE TYPE quote_status AS ENUM (
				'DRAFT',
				'ISSUED',
				'APPROVED',
				'REJECTED',
				'EXPIRED',
				'CONVERTED_TO_ORDER'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		DO $$ BEGIN
			CREATE TYPE sales_order_status AS ENUM (
				'PENDING',
				'IN_PREPARATION',
				'PARTIALLY_DELIVERED',
				'DELIVERED',
				'CANCELLED',
				'PARTIALLY_INVOICED',
				'INVOICED'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		DO $$ BEGIN
			CREATE TYPE delivery_note_status AS ENUM (
				'PENDING',
				'DELIVERED',
				'CANCELLED'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		DO $$ BEGIN
			CREATE TYPE invoice_status AS ENUM (
				'DRAFT',
				'ISSUED',
				'PAID',
				'OVERDUE',
				'VOID'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		DO $$ BEGIN
			CREATE TYPE invoice_type AS ENUM (
				'COMPLETE',
				'SIMPLIFIED'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		CREATE TABLE "product_variants" (
			"id" UUID PRIMARY KEY
		);

		CREATE TABLE "quotes" (
			"id" UUID PRIMARY KEY,
			"quote_number" VARCHAR(50) NOT NULL,
			"party_id" UUID NOT NULL,
			"quote_date" TIMESTAMP WITH TIME ZONE NOT NULL,
			"expiration_date" TIMESTAMP WITH TIME ZONE NOT NULL,
			"status" quote_status NOT NULL,
			"mes_work_refs" JSONB,
			"subtotal_amount" NUMERIC(12,2) NOT NULL,
			"subtotal_currency" VARCHAR(3) NOT NULL,
			"tax_amount" NUMERIC(12,2) NOT NULL,
			"tax_currency" VARCHAR(3) NOT NULL,
			"total_amount" NUMERIC(12,2) NOT NULL,
			"total_currency" VARCHAR(3) NOT NULL,
			"notes" TEXT,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE "quote_line_items" (
			"id" UUID PRIMARY KEY,
			"quote_id" UUID NOT NULL REFERENCES "quotes" ("id") ON DELETE CASCADE,
			"product_variant_id" UUID NOT NULL REFERENCES "product_variants" ("id") ON DELETE RESTRICT,
			"quantity" INT NOT NULL,
			"list_unit_price_amount" NUMERIC(12,2) NOT NULL,
			"list_unit_price_currency" VARCHAR(3) NOT NULL,
			"unit_price_amount" NUMERIC(12,2) NOT NULL,
			"unit_price_currency" VARCHAR(3) NOT NULL,
			"discount_percent" NUMERIC(5,2) NOT NULL DEFAULT 0,
			"discount_per_unit_amount" NUMERIC(12,2) NOT NULL,
			"discount_per_unit_currency" VARCHAR(3) NOT NULL,
			"subtotal_amount" NUMERIC(12,2) NOT NULL,
			"subtotal_currency" VARCHAR(3) NOT NULL,
			"tax_rate" NUMERIC(5,2) NOT NULL DEFAULT 21.00,
			"tax_amount" NUMERIC(10,2),
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE "sales_orders" (
			"id" UUID PRIMARY KEY,
			"order_number" VARCHAR(50) NOT NULL,
			"quote_id" UUID,
			"party_id" UUID NOT NULL,
			"order_date" TIMESTAMP WITH TIME ZONE NOT NULL,
			"delivery_date" TIMESTAMP WITH TIME ZONE NOT NULL,
			"status" sales_order_status NOT NULL,
			"mes_work_refs" JSONB,
			"subtotal_amount" NUMERIC(12,2) NOT NULL,
			"subtotal_currency" VARCHAR(3) NOT NULL,
			"tax_amount" NUMERIC(12,2) NOT NULL,
			"tax_currency" VARCHAR(3) NOT NULL,
			"total_amount" NUMERIC(12,2) NOT NULL,
			"total_currency" VARCHAR(3) NOT NULL,
			"notes" TEXT,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE "order_line_items" (
			"id" UUID PRIMARY KEY,
			"sales_order_id" UUID NOT NULL REFERENCES "sales_orders" ("id") ON DELETE CASCADE,
			"product_variant_id" UUID NOT NULL REFERENCES "product_variants" ("id") ON DELETE RESTRICT,
			"quantity" INT NOT NULL,
			"list_unit_price_amount" NUMERIC(12,2) NOT NULL,
			"list_unit_price_currency" VARCHAR(3) NOT NULL,
			"unit_price_amount" NUMERIC(12,2) NOT NULL,
			"unit_price_currency" VARCHAR(3) NOT NULL,
			"discount_percent" NUMERIC(5,2) NOT NULL DEFAULT 0,
			"discount_per_unit_amount" NUMERIC(12,2) NOT NULL,
			"discount_per_unit_currency" VARCHAR(3) NOT NULL,
			"subtotal_amount" NUMERIC(12,2) NOT NULL,
			"subtotal_currency" VARCHAR(3) NOT NULL,
			"tax_rate" NUMERIC(5,2) NOT NULL DEFAULT 21.00,
			"tax_amount" NUMERIC(10,2),
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE "delivery_notes" (
			"id" UUID PRIMARY KEY,
			"delivery_note_number" VARCHAR(50) NOT NULL,
			"sales_order_id" UUID NOT NULL REFERENCES "sales_orders" ("id") ON DELETE CASCADE,
			"party_id" UUID NOT NULL,
			"delivery_date" TIMESTAMP WITH TIME ZONE NOT NULL,
			"status" delivery_note_status NOT NULL,
			"notes" TEXT,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE "delivery_note_line_items" (
			"id" UUID PRIMARY KEY,
			"delivery_note_id" UUID NOT NULL REFERENCES "delivery_notes" ("id") ON DELETE CASCADE,
			"sales_order_line_item_id" UUID NOT NULL REFERENCES "order_line_items" ("id") ON DELETE RESTRICT,
			"product_variant_id" UUID NOT NULL REFERENCES "product_variants" ("id") ON DELETE RESTRICT,
			"delivered_quantity" INT NOT NULL,
			"invoice_line_item_id" UUID,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE "invoices" (
			"id" UUID PRIMARY KEY,
			"invoice_number" VARCHAR(50) NOT NULL,
			"type" invoice_type NOT NULL DEFAULT 'COMPLETE',
			"series_code" VARCHAR(10) NOT NULL DEFAULT 'A',
			"series_year" INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM NOW()),
			"series_prefix" VARCHAR(10) NOT NULL DEFAULT 'A',
			"party_id" UUID NOT NULL,
			"invoice_date" TIMESTAMP WITH TIME ZONE NOT NULL,
			"due_date" TIMESTAMP WITH TIME ZONE NOT NULL,
			"status" invoice_status NOT NULL,
			"payment_terms" TEXT,
			"subtotal_amount" NUMERIC(12,2) NOT NULL,
			"subtotal_currency" VARCHAR(3) NOT NULL,
			"tax_amount" NUMERIC(12,2) NOT NULL,
			"tax_currency" VARCHAR(3) NOT NULL,
			"total_amount" NUMERIC(12,2) NOT NULL,
			"total_currency" VARCHAR(3) NOT NULL,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE "invoice_line_items" (
			"id" UUID PRIMARY KEY,
			"invoice_id" UUID NOT NULL REFERENCES "invoices" ("id") ON DELETE CASCADE,
			"sales_order_line_item_id" UUID REFERENCES "order_line_items" ("id") ON DELETE SET NULL,
			"product_variant_id" UUID NOT NULL REFERENCES "product_variants" ("id") ON DELETE RESTRICT,
			"quantity" INT NOT NULL,
			"unit_price_amount" NUMERIC(12,2) NOT NULL,
			"unit_price_currency" VARCHAR(3) NOT NULL,
			"tax_rate" NUMERIC(5,2) NOT NULL DEFAULT 21.00,
			"discount_amount" NUMERIC(12,2),
			"discount_currency" VARCHAR(3),
			"subtotal_amount" NUMERIC(12,2) NOT NULL,
			"subtotal_currency" VARCHAR(3) NOT NULL,
			"tax_amount" NUMERIC(12,2),
			"tax_currency" VARCHAR(3),
			"total_amount" NUMERIC(12,2) NOT NULL,
			"total_currency" VARCHAR(3) NOT NULL,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"deleted_at" TIMESTAMP WITH TIME ZONE
		);
	`

	if err := tdb.DB.WithContext(ctx).Exec(createSchema).Error; err != nil {
		return fmt.Errorf("failed to create sales schema: %w", err)
	}

	return nil
}

func (tdb *TestDB) TearDownSales() error {
	ctx := context.Background()

	dropSchema := `
		DROP TABLE IF EXISTS "invoice_line_items" CASCADE;
		DROP TABLE IF EXISTS "invoices" CASCADE;
		DROP TABLE IF EXISTS "delivery_note_line_items" CASCADE;
		DROP TABLE IF EXISTS "delivery_notes" CASCADE;
		DROP TABLE IF EXISTS "order_line_items" CASCADE;
		DROP TABLE IF EXISTS "sales_orders" CASCADE;
		DROP TABLE IF EXISTS "quote_line_items" CASCADE;
		DROP TABLE IF EXISTS "quotes" CASCADE;
		DROP TABLE IF EXISTS "product_variants" CASCADE;
		DROP TYPE IF EXISTS quote_status;
		DROP TYPE IF EXISTS sales_order_status;
		DROP TYPE IF EXISTS delivery_note_status;
		DROP TYPE IF EXISTS invoice_status;
	`

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop sales schema: %w", err)
	}

	return nil
}
