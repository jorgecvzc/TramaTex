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

	host := os.Getenv("DB_HOST")
	if host == "" { host = "localhost" }
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("--- Checking Quotes ---")
	var quotes []map[string]interface{}
	db.Raw("SELECT id, quote_number, status, party_id FROM quotes ORDER BY created_at DESC LIMIT 5").Scan(&quotes)
	for _, q := range quotes {
		fmt.Printf("Quote: %v | Number: %v | Status: %v | Party: %v\n", q["id"], q["quote_number"], q["status"], q["party_id"])
		
		var items []map[string]interface{}
		db.Raw("SELECT id, product_variant_id, quantity, subtotal_amount FROM quote_line_items WHERE quote_id = ?", q["id"]).Scan(&items)
		fmt.Printf("  Lines (%d):\n", len(items))
		for _, item := range items {
			fmt.Printf("    - ID: %v | Variant: %v | Qty: %v | Subtotal: %v\n", item["id"], item["product_variant_id"], item["quantity"], item["subtotal_amount"])
		}
	}

	fmt.Println("\n--- Checking Orders ---")
	var orders []map[string]interface{}
	db.Raw("SELECT id, order_number, status, party_id FROM sales_orders ORDER BY created_at DESC LIMIT 5").Scan(&orders)
	for _, o := range orders {
		fmt.Printf("Order: %v | Number: %v | Status: %v | Party: %v\n", o["id"], o["order_number"], o["status"], o["party_id"])
		
		var items []map[string]interface{}
		db.Raw("SELECT id, product_variant_id, quantity, subtotal_amount FROM order_line_items WHERE sales_order_id = ?", o["id"]).Scan(&items)
		fmt.Printf("  Lines (%d):\n", len(items))
		for _, item := range items {
			fmt.Printf("    - ID: %v | Variant: %v | Qty: %v | Subtotal: %v\n", item["id"], item["product_variant_id"], item["quantity"], item["subtotal_amount"])
		}
	}
}
