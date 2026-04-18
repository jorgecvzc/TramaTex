package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

func main() {
	fmt.Println("=== Validando Flujo de Sales (Dominio) ===")

	// 1. Crear Presupuesto
	quoteNumber, _ := domain.NewQuoteNumber("PRE-2026-0001")
	partyID := uuid.New()
	expirationDate := time.Now().Add(30 * 24 * time.Hour)
	money, _ := domain.NewMoney(100.0, "EUR")
	
	lineItem, _ := domain.NewQuoteLineItem(uuid.New(), 2, money, nil, 10.0, 21.0)
	
	quote, err := domain.NewQuote(
		quoteNumber,
		partyID,
		time.Now(),
		expirationDate,
		[]domain.QuoteLineItem{lineItem},
		money, // Tax amount placeholder
		"Presupuesto de prueba",
	)
	if err != nil {
		log.Fatalf("Error creando presupuesto: %v", err)
	}
	fmt.Printf("Presupuesto creado: %s (Status: %s)\n", quote.QuoteNumber, quote.Status)

	// 2. Cambiar Status a ISSUED
	err = quote.ChangeStatus(domain.QuoteStatusIssued)
	if err != nil {
		log.Fatalf("Error cambiando status a ISSUED: %v", err)
	}
	fmt.Printf("Presupuesto emitido (Status: %s)\n", quote.Status)

	// 3. Aprobar y Convertir a Pedido
	err = quote.ChangeStatus(domain.QuoteStatusApproved)
	if err != nil {
		log.Fatalf("Error aprobando presupuesto: %v", err)
	}
	
	orderNumber, _ := domain.NewOrderNumber("PED-2026-0001")
	order, err := quote.ConvertToOrder(orderNumber, time.Now().Add(7*24*time.Hour))
	if err != nil {
		log.Fatalf("Error convirtiendo a pedido: %v", err)
	}
	fmt.Printf("Pedido creado: %s (Status: %s)\n", order.OrderNumber, order.Status)
	fmt.Printf("Status presupuesto: %s\n", quote.Status)

	// 4. Crear Albarán
	dnNumber, _ := domain.NewDeliveryNoteNumber("ALB-2026-0001")
	dnLineItem, _ := domain.NewDeliveryNoteLineItem(order.LineItems[0].ID, order.LineItems[0].ProductVariantID, order.LineItems[0].Quantity)
	
	dn, err := domain.NewDeliveryNote(
		dnNumber,
		order.ID,
		order.PartyID,
		time.Now(),
		[]domain.DeliveryNoteLineItem{dnLineItem},
		"Albarán de prueba",
	)
	if err != nil {
		log.Fatalf("Error creando albarán: %v", err)
	}
	fmt.Printf("Albarán creado: %s (Status: %s)\n", dn.DeliveryNoteNumber, dn.Status)

	fmt.Println("=== Flujo validado correctamente en el dominio ===")
}
