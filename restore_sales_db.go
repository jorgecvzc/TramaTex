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
	_ = godotenv.Load("apps/tramatex-api/.env")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		"localhost",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Open connection WITHOUT automatic transaction for the whole thing
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	files := []string{
		"apps/tramatex-api/migrations/003_init_product.sql",
		"apps/tramatex-api/migrations/005_init_sales.sql",
		"apps/tramatex-api/migrations/010_align_sales_enums.sql",
	}

	for _, file := range files {
		fmt.Printf("\n--- Applying migration: %s ---\n", file)
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			continue
		}

		sql := string(content)
		// Run each file in its own Exec call to avoid aborted transaction propagation
		if err := db.Exec(sql).Error; err != nil {
			fmt.Printf("Error executing %s: %v\n", file, err)
		} else {
			fmt.Printf("SUCCESS: %s applied.\n", file)
		}
	}
}
