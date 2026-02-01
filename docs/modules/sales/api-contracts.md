# Sales Module API Contracts

This document specifies the API contracts for the Sales module, which handles quotes and orders.

---

## 1. Orders

### 1.1 List Orders
- **Endpoint:** `GET /orders`
- **Description:** Retrieves a list of sales orders, with optional filtering.
- **Query Parameters:**
  - `client_id` (string, optional): Filter orders by client ID.
  - `status` (string, optional): Filter orders by status (e.g., "CONFIRMED", "SHIPPED").
- **Success Response (200 OK):** `ListResponse<OrderDTO>`

### 1.2 Create Order
- **Endpoint:** `POST /orders`
- **Description:** Creates a new sales order from a list of items or converts a quote.
- **Request Body:**
  ```json
  {
    "client_id": "client-uuid-789",
    "quote_id": "quote-uuid-abc" // Optional: to convert a quote
    "items": [ // Optional: to create directly
      {
        "product_variant_id": "variant-uuid-456",
        "quantity": 50
      }
    ]
  }
  ```
- **Success Response (201 Created):** `OrderDTO`

### 1.3 Get Order
- **Endpoint:** `GET /orders/{id}`
- **Description:** Retrieves a single sales order by its ID.
- **Success Response (200 OK):** `OrderDTO` (with line items)

### 1.4 Update Order
- **Endpoint:** `PUT /orders/{id}`
- **Description:** Updates the details of an existing order (e.g., adds notes, changes delivery address). This does not modify line items.
- **Request Body:**
  ```json
  {
    "notes": "Updated delivery instructions."
  }
  ```
- **Success Response (200 OK):** `OrderDTO`

---

## 2. Order Items

### 2.1 Add Item to Order
- **Endpoint:** `POST /orders/{id}/items`
- **Description:** Adds a new line item to an existing order (if the order is in a modifiable state, e.g., 'DRAFT').
- **Request Body:**
  ```json
  {
    "product_variant_id": "variant-uuid-789",
    "quantity": 10
  }
  ```
- **Success Response (201 Created):** `OrderLineItemDTO`

---

## 3. Quotes

### 3.1 List Quotes
- **Endpoint:** `GET /quotes`
- **Description:** Retrieves a list of all quotes.
- **Success Response (200 OK):** `ListResponse<QuoteDTO>`

### 3.2 Create Quote
- **Endpoint:** `POST /quotes`
- **Description:** Creates a new price quote for a client.
- **Request Body:**
  ```json
  {
    "client_id": "client-uuid-789",
    "items": [
      {
        "product_variant_id": "variant-uuid-456",
        "quantity": 100
      }
    ]
  }
  ```
- **Success Response (201 Created):** `QuoteDTO`
