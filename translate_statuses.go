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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	fmt.Println("Actualizando estados de Español a Inglés en la DB...")

	updates := []struct {
		table string
		from  string
		to    string
	}{
		// Quotes
		{"quotes", "BORRADOR", "DRAFT"},
		{"quotes", "EMITIDA", "ISSUED"},
		{"quotes", "APROBADA", "APPROVED"},
		{"quotes", "RECHAZADA", "REJECTED"},
		{"quotes", "EXPIRADA", "EXPIRED"},
		{"quotes", "CONVERTIDA_A_PEDIDO", "CONVERTED_TO_ORDER"},
		
		// Sales Orders
		{"sales_orders", "PENDIENTE", "PENDING"},
		{"sales_orders", "EN_PREPARACION", "IN_PREPARATION"},
		{"sales_orders", "ENTREGADO_PARCIALMENTE", "PARTIALLY_DELIVERED"},
		{"sales_orders", "ENTREGADO", "DELIVERED"},
		{"sales_orders", "CANCELADO", "CANCELLED"},
		{"sales_orders", "FACTURADO_PARCIALMENTE", "PARTIALLY_INVOICED"},
		{"sales_orders", "FACTURADO_COMPLETAMENTE", "INVOICED"},

		// Delivery Notes
		{"delivery_notes", "PENDIENTE", "PENDING"},
		{"delivery_notes", "ENTREGADO", "DELIVERED"},
		{"delivery_notes", "CANCELADO", "CANCELLED"},

		// Invoices
		{"invoices", "BORRADOR", "DRAFT"},
		{"invoices", "EMITIDA", "ISSUED"},
		{"invoices", "PAGADA", "PAID"},
		{"invoices", "VENCIDA", "OVERDUE"},
		{"invoices", "ANULADA", "VOID"},
	}

	for _, u := range updates {
		result := db.Exec(fmt.Sprintf("UPDATE %s SET status = ? WHERE status::text = ?", u.table), u.to, u.from)
		if result.Error != nil {
			fmt.Printf("Error actualizando %s (%s -> %s): %v\n", u.table, u.from, u.to, result.Error)
		} else if result.RowsAffected > 0 {
			fmt.Printf("Actualizado %s: %d filas cambiadas de %s a %s\n", u.table, result.RowsAffected, u.from, u.to)
		}
	}

	fmt.Println("Migración de datos completada.")
}
