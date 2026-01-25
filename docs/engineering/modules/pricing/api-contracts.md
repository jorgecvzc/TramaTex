# Pricing Module API Contracts

This document specifies the API contracts for the Pricing module.

---

## 1. Calculate Price

- **Endpoint:** `POST /pricing/calculate`
- **Description:** Calculates the price for a given product variant, quantity, and client. This is the core endpoint of the pricing engine.
- **Request Body:**
  ```json
  {
    "product_variant_id": "variant-uuid-456",
    "client_id": "client-uuid-789",
    "quantity": 100
  }
  ```
- **Success Response (200 OK):**
  ```json
  {
    "final_price": 175.50,
    "currency": "EUR",
    "breakdown": {
      "base_cost": 120.00,
      "margin_applied": "25%",
      "discounts": [
        {
          "name": "Volume Discount",
          "amount": 15.00
        }
      ]
    }
  }
  ```
- **Error Responses:**
  - `400 Bad Request`: If any of the IDs are invalid or quantity is zero.
  - `404 Not Found`: If the product variant or client does not exist.

---

## 2. Pricing Rules

### 2.1 List Pricing Rules
- **Endpoint:** `GET /pricing/rules`
- **Description:** Retrieves a list of all pricing rules.
- **Success Response (200 OK):** `ListResponse<PricingRuleDTO>`

### 2.2 Create Pricing Rule
- **Endpoint:** `POST /pricing/rules`
- **Description:** Creates a new pricing rule.
- **Request Body:**
  ```json
  {
    "product_category": "T-SHIRTS",
    "client_category": "WHOLESALE",
    "margin_percentage": 25.0,
    "volume_discounts": [
      { "min_quantity": 100, "discount_percentage": 5.0 },
      { "min_quantity": 500, "discount_percentage": 10.0 }
    ]
  }
  ```
- **Success Response (201 Created):** `PricingRuleDTO`

---

## 3. Client Pricing Overrides

- **Endpoint:** `POST /pricing/client-overrides`
- **Description:** Sets a specific price or discount for a specific client and product, overriding general rules.
- **Request Body:**
  ```json
  {
    "client_id": "client-uuid-789",
    "product_variant_id": "variant-uuid-456",
    "fixed_price": 160.00 
  }
  ```
- **Success Response (201 Created):** `ClientOverrideDTO`

---

## 4. Get Pricing History

- **Endpoint:** `GET /pricing/history/{product-id}`
- **Description:** Retrieves the history of price calculations for a specific product.
- **Success Response (200 OK):** `ListResponse<PriceCalculationDTO>`
