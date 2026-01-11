package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Basic health check
	fmt.Println("🚀 TramaTex API starting...")
	fmt.Printf("Environment: %s\n", os.Getenv("ENV"))

	// TODO: Initialize database
	// TODO: Setup HTTP server
	// TODO: Setup routes
}
