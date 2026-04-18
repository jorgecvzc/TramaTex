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

	host := "localhost"
	port := "5432"
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	fmt.Println("Inspecting quote_status enum...")
	var enumValues []struct {
		EnumLabel string `gorm:"column:enumlabel"`
	}
	db.Raw(`
		SELECT e.enumlabel
		FROM pg_type t 
		JOIN pg_enum e ON t.oid = e.enumtypid  
		JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace 
		WHERE t.typname = 'quote_status'
	`).Scan(&enumValues)

	for _, v := range enumValues {
		fmt.Printf("- %s\n", v.EnumLabel)
	}

	fmt.Println("\nChecking for quotes with 'DRAFT' status (if any)...")
	var count int64
	// Use Raw to avoid GORM mapping issues
	db.Raw("SELECT count(*) FROM quotes WHERE status::text = 'DRAFT'").Scan(&count)
	fmt.Printf("Quotes with status 'DRAFT': %d\n", count)

	db.Raw("SELECT count(*) FROM quotes WHERE status::text = 'BORRADOR'").Scan(&count)
	fmt.Printf("Quotes with status 'BORRADOR': %d\n", count)
}
