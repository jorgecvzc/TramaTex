package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	baseURL    = "http://localhost:8080"
	adminEmail = "admin@tramatex.local"
	adminPass  = "admin123"

	// Test data IDs from seed
	nikeID      = "11111111-1111-1111-1111-111111111111"
	calzadoID   = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	tallaAttrID = "a0000000-0000-0000-0000-000000000001"
	colorAttrID = "a0000000-0000-0000-0000-000000000002"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type ProductResponse struct {
	ID                 string   `json:"id"`
	SKU                string   `json:"sku"`
	Name               string   `json:"name"`
	DirectAttributeIDs []string `json:"directAttributeIds"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	separator := strings.Repeat("=", 70)

	fmt.Println("Test E2E - Creacion de Producto con Atributos Directos")
	fmt.Println(separator)

	// Test 1: Health Check
	if !testHealthCheck() {
		os.Exit(1)
	}

	// Test 2: Login
	token := testLogin()
	if token == "" {
		os.Exit(1)
	}

	// Test 3: Crear producto con atributos directos (SKU único)
	timestamp := time.Now().Unix()
	productSKU := fmt.Sprintf("E2E-TEST-%d", timestamp)
	if !testCreateProduct(token, productSKU, true) {
		os.Exit(1)
	}

	// Test 4: Intentar crear producto con SKU duplicado (debe fallar con 409)
	if !testCreateProductDuplicateSKU(token, productSKU) {
		os.Exit(1)
	}

	fmt.Println("\nTODOS LOS TESTS PASARON EXITOSAMENTE")
	fmt.Println(separator)
	fmt.Println("\nBug Fix Verificado:")
	fmt.Println("  - Los atributos directos se envian correctamente")
	fmt.Println("  - Los atributos directos se persisten en la base de datos")
	fmt.Println("  - La validacion de SKU duplicado funciona correctamente")
}

func testHealthCheck() bool {
	fmt.Println("\n[Test 1] Health Check")
	resp, err := http.Get(baseURL + "/api/health")
	if err != nil {
		fmt.Printf("  FALLO: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("  OK - API esta respondiendo correctamente")
		return true
	}

	fmt.Printf("  FALLO: Status code %d\n", resp.StatusCode)
	return false
}

func testLogin() string {
	fmt.Println("\n[Test 2] Login")

	loginData := map[string]string{
		"email":    adminEmail,
		"password": adminPass,
	}

	resp, body := doRequest("POST", "/auth/login", "", loginData)
	if resp.StatusCode != 200 {
		fmt.Printf("  FALLO: Login fallo con status %d\n", resp.StatusCode)
		fmt.Printf("  Response: %s\n", body)
		return ""
	}

	var loginResp LoginResponse
	if err := json.Unmarshal([]byte(body), &loginResp); err != nil {
		fmt.Printf("  FALLO: No se pudo parsear la respuesta: %v\n", err)
		return ""
	}

	if loginResp.AccessToken == "" {
		fmt.Println("  FALLO: Token vacio")
		return ""
	}

	fmt.Printf("  OK - Login exitoso (token: %s...)\n", loginResp.AccessToken[:20])
	return loginResp.AccessToken
}

func testCreateProduct(token, sku string, expectSuccess bool) bool {
	fmt.Printf("\n[Test 3] Crear Producto con SKU: %s\n", sku)

	productData := map[string]interface{}{
		"sku":                  sku,
		"name":                 "Producto E2E Test",
		"long_name":            "Producto de Prueba End-to-End con Atributos Directos",
		"description":          "Test automatizado del bug fix - Atributos directos deben persistirse correctamente",
		"product_type":         "TANGIBLE",
		"brand_id":             nikeID,
		"group_ids":            []string{calzadoID},
		"direct_attribute_ids": []string{tallaAttrID, colorAttrID},
	}

	resp, body := doRequest("POST", "/api/products", token, productData)

	if expectSuccess {
		if resp.StatusCode != 201 {
			fmt.Printf("  FALLO: Se esperaba 201, recibido %d\n", resp.StatusCode)
			fmt.Printf("  Response: %s\n", body)
			return false
		}

		var productResp ProductResponse
		if err := json.Unmarshal([]byte(body), &productResp); err != nil {
			fmt.Printf("  FALLO: No se pudo parsear la respuesta: %v\n", err)
			return false
		}

		fmt.Printf("  OK - Producto creado (ID: %s)\n", productResp.ID)
		fmt.Printf("  OK - SKU: %s\n", productResp.SKU)
		fmt.Printf("  OK - Nombre: %s\n", productResp.Name)

		// Verificar atributos directos
		if len(productResp.DirectAttributeIDs) != 2 {
			fmt.Printf("  FALLO: Se esperaban 2 atributos directos, recibidos %d\n", len(productResp.DirectAttributeIDs))
			return false
		}

		fmt.Printf("  OK - Atributos Directos (%d):\n", len(productResp.DirectAttributeIDs))
		for i, attrID := range productResp.DirectAttributeIDs {
			attrName := "Unknown"
			if attrID == tallaAttrID {
				attrName = "Talla"
			} else if attrID == colorAttrID {
				attrName = "Color"
			}
			fmt.Printf("    %d. %s (%s)\n", i+1, attrName, attrID)
		}

		return true
	}

	return true
}

func testCreateProductDuplicateSKU(token, sku string) bool {
	fmt.Printf("\n[Test 4] Intentar Crear Producto con SKU Duplicado: %s\n", sku)

	productData := map[string]interface{}{
		"sku":                  sku, // Mismo SKU del test anterior
		"name":                 "Producto Duplicado",
		"product_type":         "TANGIBLE",
		"brand_id":             nikeID,
		"group_ids":            []string{calzadoID},
		"direct_attribute_ids": []string{tallaAttrID},
	}

	resp, body := doRequest("POST", "/api/products", token, productData)

	// Debe devolver 409 Conflict
	if resp.StatusCode == 409 {
		var errResp ErrorResponse
		if err := json.Unmarshal([]byte(body), &errResp); err == nil {
			fmt.Printf("  OK - Correctamente rechazado con 409 Conflict\n")
			fmt.Printf("  OK - Mensaje de error: %s\n", errResp.Error)
			return true
		}
	}

	// Si devuelve 500, el bug NO está arreglado
	if resp.StatusCode == 500 {
		fmt.Println("  FALLO: El servidor devolvio 500 en lugar de 409")
		fmt.Println("  INFO: Esto indica que la validacion de SKU duplicado no esta funcionando")
		fmt.Printf("  Response: %s\n", body)
		return false
	}

	fmt.Printf("  FALLO: Status code inesperado %d (se esperaba 409)\n", resp.StatusCode)
	fmt.Printf("  Response: %s\n", body)
	return false
}

func doRequest(method, endpoint, token string, data interface{}) (*http.Response, string) {
	var body io.Reader
	if data != nil {
		jsonData, _ := json.Marshal(data)
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	return resp, string(bodyBytes)
}
