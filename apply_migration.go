package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from apps/tramatex-api/.env
	_ = godotenv.Load("apps/tramatex-api/.env")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	fmt.Println("Applying migration 010_align_sales_enums.sql...")
	
	// Execute each command separately (needed for ALTER TYPE ADD VALUE)
	commands := []string{
		"ALTER TYPE quote_status ADD VALUE IF NOT EXISTS 'DRAFT'",
		"ALTER TYPE quote_status ADD VALUE IF NOT EXISTS 'ISSUED'",
		"ALTER TYPE quote_status ADD VALUE IF NOT EXISTS 'APPROVED'",
		"ALTER TYPE quote_status ADD VALUE IF NOT EXISTS 'ACCEPTED'",
		"ALTER TYPE quote_status ADD VALUE IF NOT EXISTS 'REJECTED'",
		"ALTER TYPE quote_status ADD VALUE IF NOT EXISTS 'EXPIRED'",
		"ALTER TYPE quote_status ADD VALUE IF NOT EXISTS 'CONVERTED_TO_ORDER'",
		"ALTER TYPE sales_order_status ADD VALUE IF NOT EXISTS 'PENDING'",
		"ALTER TYPE sales_order_status ADD VALUE IF NOT EXISTS 'IN_PREPARATION'",
		"ALTER TYPE sales_order_status ADD VALUE IF NOT EXISTS 'READY_FOR_PRODUCTION'",
		"ALTER TYPE sales_order_status ADD VALUE IF NOT EXISTS 'PARTIALLY_DELIVERED'",
		"ALTER TYPE sales_order_status ADD VALUE IF NOT EXISTS 'DELIVERED'",
		"ALTER TYPE sales_order_status ADD VALUE IF NOT EXISTS 'CANCELLED'",
		"ALTER TYPE sales_order_status ADD VALUE IF NOT EXISTS 'PARTIALLY_INVOICED'",
		"ALTER TYPE sales_order_status ADD VALUE IF NOT EXISTS 'INVOICED'",
		"ALTER TYPE delivery_note_status ADD VALUE IF NOT EXISTS 'PENDING'",
		"ALTER TYPE delivery_note_status ADD VALUE IF NOT EXISTS 'DELIVERED'",
		"ALTER TYPE delivery_note_status ADD VALUE IF NOT EXISTS 'CANCELLED'",
		"ALTER TYPE invoice_status ADD VALUE IF NOT EXISTS 'DRAFT'",
		"ALTER TYPE invoice_status ADD VALUE IF NOT EXISTS 'ISSUED'",
		"ALTER TYPE invoice_status ADD VALUE IF NOT EXISTS 'PAID'",
		"ALTER TYPE invoice_status ADD VALUE IF NOT EXISTS 'OVERDUE'",
		"ALTER TYPE invoice_status ADD VALUE IF NOT EXISTS 'VOID'",
	}

	for _, cmd := range commands {
		fmt.Printf("Executing: %s\n", cmd)
		if err := db.Exec(cmd).Error; err != nil {
			fmt.Printf("Error (might be expected): %v\n", err)
		}
	}

	fmt.Println("Migration attempt completed.")
}
