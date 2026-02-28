package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"type:uuid;primaryKey"`
	Email        string
	PasswordHash string
}

func main() {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=tramatex password=tramatex dbname=tramatex port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var user User
	result := db.Table("users").Where("email = ?", "admin@tramatex.local").First(&user)
	if result.Error != nil {
		log.Fatalf("Failed to find user: %v", result.Error)
	}

	fmt.Printf("✅ User found in database:\n")
	fmt.Printf("   ID: %s\n", user.ID)
	fmt.Printf("   Email: %s\n", user.Email)
	fmt.Printf("   Password Hash: %s\n", user.PasswordHash)
	fmt.Println()

	// Test password
	password := "admin123"
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err == nil {
		fmt.Printf("✅ Password '%s' MATCHES the hash in database\n", password)
	} else {
		fmt.Printf("❌ Password '%s' DOES NOT MATCH the hash in database\n", password)
		fmt.Printf("   Error: %v\n", err)
	}
}
