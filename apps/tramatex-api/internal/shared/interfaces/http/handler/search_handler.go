package handler

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SearchResult struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	URL      string `json:"url"`
	Module   string `json:"module"`
}

type SearchHandler struct {
	db *gorm.DB
}

func NewSearchHandler(db *gorm.DB) *SearchHandler {
	return &SearchHandler{db: db}
}

func (h *SearchHandler) GlobalSearch(c *gin.Context) {
	query := c.Query("q")

	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusOK, []SearchResult{})
		return
	}

	results := make([]SearchResult, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex

	searchFuncs := []func(string) []SearchResult{
		h.searchSalesOrders,
		h.searchSalesQuotes,
		h.searchInvoices,
		h.searchDeliveryNotes,
		h.searchProducts,
		h.searchParties,
		h.searchMESWorkOrders,
	}

	for _, f := range searchFuncs {
		wg.Add(1)
		go func(fn func(string) []SearchResult) {
			defer wg.Done()
			res := fn(query)
			if len(res) > 0 {
				mu.Lock()
				results = append(results, res...)
				mu.Unlock()
			}
		}(f)
	}

	wg.Wait()
	c.JSON(http.StatusOK, results)
}

func (h *SearchHandler) searchSalesOrders(q string) []SearchResult {
	var results []SearchResult
	type Row struct {
		ID          string `gorm:"column:id"`
		OrderNumber string `gorm:"column:order_number"`
		Status      string `gorm:"column:status"`
		Name        string `gorm:"column:name"`
	}
	var rows []Row

	pattern := "%" + q + "%"
	h.db.Raw(`
		SELECT so.id, so.order_number, so.status, p.name 
		FROM sales_orders so 
		JOIN parties p ON so.party_id = p.id 
		WHERE so.order_number ILIKE ? OR p.name ILIKE ? 
		LIMIT 5`, pattern, pattern).Scan(&rows)

	for _, r := range rows {
		results = append(results, SearchResult{
			ID: r.ID, Type: "order", Module: "Ventas",
			Title:    "Pedido " + r.OrderNumber,
			Subtitle: "[" + r.Status + "] · " + r.Name,
			URL:      "/sales/orders/" + r.ID,
		})
	}
	return results
}

func (h *SearchHandler) searchSalesQuotes(q string) []SearchResult {
	var results []SearchResult
	type Row struct {
		ID          string `gorm:"column:id"`
		QuoteNumber string `gorm:"column:quote_number"`
		Status      string `gorm:"column:status"`
		Name        string `gorm:"column:name"`
	}
	var rows []Row

	pattern := "%" + q + "%"
	h.db.Raw(`
		SELECT sq.id, sq.quote_number, sq.status, p.name 
		FROM sales_quotes sq 
		JOIN parties p ON sq.party_id = p.id 
		WHERE sq.quote_number ILIKE ? OR p.name ILIKE ? 
		LIMIT 5`, pattern, pattern).Scan(&rows)

	for _, r := range rows {
		results = append(results, SearchResult{
			ID: r.ID, Type: "quote", Module: "Ventas",
			Title:    "Presupuesto " + r.QuoteNumber,
			Subtitle: "[" + r.Status + "] · " + r.Name,
			URL:      "/sales/quotes/" + r.ID,
		})
	}
	return results
}

func (h *SearchHandler) searchInvoices(q string) []SearchResult {
	var results []SearchResult
	type Row struct {
		ID            string `gorm:"column:id"`
		InvoiceNumber string `gorm:"column:invoice_number"`
		Status        string `gorm:"column:status"`
		Name          string `gorm:"column:name"`
	}
	var rows []Row

	pattern := "%" + q + "%"
	h.db.Raw(`
		SELECT i.id, i.invoice_number, i.status, p.name 
		FROM sales_invoices i 
		JOIN parties p ON i.party_id = p.id 
		WHERE i.invoice_number ILIKE ? OR p.name ILIKE ? 
		LIMIT 5`, pattern, pattern).Scan(&rows)

	for _, r := range rows {
		results = append(results, SearchResult{
			ID: r.ID, Type: "invoice", Module: "Ventas",
			Title:    "Factura " + r.InvoiceNumber,
			Subtitle: "[" + r.Status + "] · " + r.Name,
			URL:      "/sales/invoices/" + r.ID,
		})
	}
	return results
}

func (h *SearchHandler) searchDeliveryNotes(q string) []SearchResult {
	var results []SearchResult
	type Row struct {
		ID                 string `gorm:"column:id"`
		DeliveryNoteNumber string `gorm:"column:delivery_note_number"`
		Status             string `gorm:"column:status"`
		Name               string `gorm:"column:name"`
	}
	var rows []Row

	pattern := "%" + q + "%"
	h.db.Raw(`
		SELECT d.id, d.delivery_note_number, d.status, p.name 
		FROM sales_delivery_notes d 
		JOIN parties p ON d.party_id = p.id 
		WHERE d.delivery_note_number ILIKE ? OR p.name ILIKE ? 
		LIMIT 5`, pattern, pattern).Scan(&rows)

	for _, r := range rows {
		results = append(results, SearchResult{
			ID: r.ID, Type: "delivery_note", Module: "Ventas",
			Title:    "Albarán " + r.DeliveryNoteNumber,
			Subtitle: "[" + r.Status + "] · " + r.Name,
			URL:      "/sales/delivery-notes/" + r.ID,
		})
	}
	return results
}

func (h *SearchHandler) searchProducts(q string) []SearchResult {
	var results []SearchResult
	type Row struct {
		ID   string `gorm:"column:id"`
		Name string `gorm:"column:name"`
		SKU  string `gorm:"column:sku"`
	}
	var rows []Row

	pattern := "%" + q + "%"
	h.db.Raw(`SELECT id, name, sku FROM products WHERE name ILIKE ? OR sku ILIKE ? LIMIT 5`, pattern, pattern).Scan(&rows)

	for _, r := range rows {
		results = append(results, SearchResult{
			ID: r.ID, Type: "product", Module: "Catálogo",
			Title:    r.Name,
			Subtitle: "SKU: " + r.SKU,
			URL:      "/products/" + r.ID,
		})
	}
	return results
}

func (h *SearchHandler) searchParties(q string) []SearchResult {
	var results []SearchResult
	type Row struct {
		ID    string `gorm:"column:id"`
		Name  string `gorm:"column:name"`
		TaxID string `gorm:"column:tax_id"`
	}
	var rows []Row

	pattern := "%" + q + "%"
	h.db.Raw(`SELECT id, name, tax_id FROM parties WHERE name ILIKE ? OR tax_id ILIKE ? LIMIT 5`, pattern, pattern).Scan(&rows)

	for _, r := range rows {
		results = append(results, SearchResult{
			ID: r.ID, Type: "party", Module: "Entidades",
			Title:    r.Name,
			Subtitle: "NIF: " + r.TaxID,
			URL:      "/parties/" + r.ID,
		})
	}
	return results
}

func (h *SearchHandler) searchMESWorkOrders(q string) []SearchResult {
	var results []SearchResult
	type Row struct {
		ID         string  `gorm:"column:id"`
		WorkNumber string  `gorm:"column:work_number"`
		Status     string  `gorm:"column:status"`
		Name       *string `gorm:"column:name"`
	}
	var rows []Row

	pattern := "%" + q + "%"
	h.db.Raw(`
		SELECT wo.id, wo.work_number, wo.status, p.name 
		FROM mes_work_orders wo 
		LEFT JOIN sales_orders so ON wo.sales_order_id = so.id 
		LEFT JOIN parties p ON so.party_id = p.id 
		WHERE wo.work_number ILIKE ? OR p.name ILIKE ? 
		LIMIT 5`, pattern, pattern).Scan(&rows)

	for _, r := range rows {
		client := "Stock"
		if r.Name != nil {
			client = *r.Name
		}
		results = append(results, SearchResult{
			ID: r.ID, Type: "mes_work", Module: "Producción",
			Title:    "Trabajo " + r.WorkNumber,
			Subtitle: "[" + r.Status + "] · Cliente: " + client,
			URL:      "/mes/work-orders/" + r.ID,
		})
	}
	return results
}
