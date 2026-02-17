package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run verify.go <hash> <password>")
		fmt.Println("Example: go run verify.go '$2a$10$...' admin123")
		os.Exit(1)
	}

	hash := os.Args[1]
	password := os.Args[2]

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err == nil {
		fmt.Println("\n✅ Password MATCHES the hash!")
		fmt.Printf("Hash: %s\n", hash)
		fmt.Printf("Password: %s\n", password)
		os.Exit(0)
	} else {
		fmt.Println("\n❌ Password DOES NOT MATCH the hash!")
		fmt.Printf("Hash: %s\n", hash)
		fmt.Printf("Password: %s\n", password)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
