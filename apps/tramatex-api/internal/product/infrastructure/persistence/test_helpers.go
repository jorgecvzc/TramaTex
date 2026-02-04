package persistence

import (
	"context"
	"fmt"
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
	// Connection string for local PostgreSQL
	// Make sure PostgreSQL is running and accessible
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=tramatex_test sslmode=disable"

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
		DROP TYPE IF EXISTS product_type;
		DROP TYPE IF EXISTS variant_status;
	`

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop product schema: %w", err)
	}

	// Create enums and tables (same as migration)
	createSchema := `
		CREATE TYPE product_type AS ENUM ('TANGIBLE', 'SERVICE');
		CREATE TYPE variant_status AS ENUM ('PROVISIONAL', 'CONFIRMED');

		CREATE TABLE "brands" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"name" VARCHAR(255) NOT NULL,
			"is_active" BOOLEAN NOT NULL DEFAULT true,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE TABLE "product_groups" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"name" VARCHAR(255) NOT NULL,
			"parent_group_id" UUID,
			"is_active" BOOLEAN NOT NULL DEFAULT true,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			CONSTRAINT "fk_parent_group" FOREIGN KEY ("parent_group_id") REFERENCES "product_groups" ("id") ON DELETE SET NULL
		);

		CREATE TABLE "attributes" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"name" VARCHAR(255) NOT NULL,
			"code" VARCHAR(50) NOT NULL,
			"sort_order" INT NOT NULL DEFAULT 0,
			"scope_brand_id" UUID,
			"scope_group_id" UUID,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			CONSTRAINT "fk_scope_brand" FOREIGN KEY ("scope_brand_id") REFERENCES "brands" ("id") ON DELETE CASCADE,
			CONSTRAINT "fk_scope_product_group" FOREIGN KEY ("scope_group_id") REFERENCES "product_groups" ("id") ON DELETE CASCADE
		);

		CREATE TABLE "attribute_values" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"attribute_id" UUID NOT NULL,
			"value" VARCHAR(255) NOT NULL,
			"code" VARCHAR(50) NOT NULL,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
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
			"is_active" BOOLEAN NOT NULL DEFAULT true,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			CONSTRAINT "fk_products_brand" FOREIGN KEY ("brand_id") REFERENCES "brands" ("id") ON DELETE CASCADE
		);

		CREATE TABLE "product_to_groups" (
			"product_id" UUID NOT NULL,
			"group_id" UUID NOT NULL,
			PRIMARY KEY ("product_id", "group_id"),
			CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE,
			CONSTRAINT "fk_group" FOREIGN KEY ("group_id") REFERENCES "product_groups" ("id") ON DELETE CASCADE
		);

		CREATE TABLE "product_direct_attributes" (
			"product_id" UUID NOT NULL,
			"attribute_id" UUID NOT NULL,
			PRIMARY KEY ("product_id", "attribute_id"),
			CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE,
			CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES "attributes" ("id") ON DELETE CASCADE
		);

		CREATE TABLE "product_variants" (
			"id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			"product_id" UUID NOT NULL,
			"sku" VARCHAR(255) NOT NULL,
			"barcode" VARCHAR(255),
			"status" variant_status NOT NULL DEFAULT 'PROVISIONAL',
			"is_active" BOOLEAN NOT NULL DEFAULT true,
			"created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			"updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			CONSTRAINT "fk_product_variants_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE
		);

		CREATE TABLE "product_variant_values" (
			"variant_id" UUID NOT NULL,
			"attribute_value_id" UUID NOT NULL,
			PRIMARY KEY ("variant_id", "attribute_value_id"),
			CONSTRAINT "fk_variant" FOREIGN KEY ("variant_id") REFERENCES "product_variants" ("id") ON DELETE CASCADE,
			CONSTRAINT "fk_attribute_value" FOREIGN KEY ("attribute_value_id") REFERENCES "attribute_values" ("id") ON DELETE CASCADE
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
		DROP TABLE IF EXISTS "product_variant_values" CASCADE;
		DROP TABLE IF EXISTS "product_variants" CASCADE;
		DROP TABLE IF EXISTS "product_direct_attributes" CASCADE;
		DROP TABLE IF EXISTS "product_to_groups" CASCADE;
		DROP TABLE IF EXISTS "products" CASCADE;
		DROP TABLE IF EXISTS "attribute_values" CASCADE;
		DROP TABLE IF EXISTS "attributes" CASCADE;
		DROP TABLE IF EXISTS "product_groups" CASCADE;
		DROP TABLE IF EXISTS "brands" CASCADE;
		DROP TYPE IF EXISTS product_type;
		DROP TYPE IF EXISTS variant_status;
	`

	if err := tdb.DB.WithContext(ctx).Exec(dropSchema).Error; err != nil {
		return fmt.Errorf("failed to drop product schema: %w", err)
	}

	sqlDB, _ := tdb.DB.DB()
	return sqlDB.Close()
}
