package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost is the cost parameter for bcrypt hashing
	// Must match the cost in the Password domain model
	BcryptCost = 10
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <password>")
		fmt.Println("Example: go run main.go admin123")
		os.Exit(1)
	}

	password := os.Args[1]

	if len(password) < 8 {
		fmt.Println("Error: Password must be at least 8 characters")
		os.Exit(1)
	}

	if len(password) > 72 {
		fmt.Println("Error: Password must be at most 72 characters (bcrypt limit)")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Password hash generated successfully!")
	fmt.Println("\n📋 Hash (copy this):")
	fmt.Printf("%s\n", string(hash))
	fmt.Println("\n📝 SQL UPDATE command:")
	fmt.Printf("UPDATE users SET password_hash = '%s' WHERE email = 'admin@tramatex.local';\n", string(hash))
	fmt.Println("\n🔐 Test credentials:")
	fmt.Printf("Email: admin@tramatex.local\n")
	fmt.Printf("Password: %s\n", password)
}
