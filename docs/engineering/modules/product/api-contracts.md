# Product Module API Contracts

This document specifies the API contracts for the Product module, which handles the product catalog, variants, and categories.

---

## 1. Products

### 1.1 List Products
- **Endpoint:** `GET /products`
- **Description:** Retrieves a list of products, with optional filtering by category.
- **Query Parameters:**
  - `category` (string, optional): Filter products by a specific category name.
- **Success Response (200 OK):** `ListResponse<ProductDTO>`

### 1.2 Create Product
- **Endpoint:** `POST /products`
- **Description:** Creates a new base product in the catalog.
- **Request Body:**
  ```json
  {
    "name": "Premium T-Shirt",
    "description": "High-quality cotton t-shirt.",
    "category_id": "cat-uuid-123"
  }
  ```
- **Success Response (201 Created):** `ProductDTO`

### 1.3 Get Product
- **Endpoint:** `GET /products/{id}`
- **Description:** Retrieves a single product by its ID, including its variants.
- **Success Response (200 OK):** `ProductDTO` (with a populated `variants` array)

### 1.4 Update Product
- **Endpoint:** `PUT /products/{id}`
- **Description:** Updates the details of an existing product.
- **Request Body:**
  ```json
  {
    "name": "Premium V-Neck T-Shirt",
    "description": "An updated description for the t-shirt."
  }
  ```
- **Success Response (200 OK):** `ProductDTO`

---

## 2. Variants

### 2.1 Create Variant for Product
- **Endpoint:** `POST /products/{id}/variants`
- **Description:** Adds a new variant (e.g., size, color) to an existing product.
- **Request Body:**
  ```json
  {
    "attributes": [
      { "name": "size", "value": "L" },
      { "name": "color", "value": "blue" }
    ],
    "supplier_cost": 12.50,
    "sku_suffix": "L-BLU"
  }
  ```
- **Success Response (201 Created):** `VariantDTO`

---

## 3. Suppliers

### 3.1 Get Suppliers for Product
- **Endpoint:** `GET /products/{id}/suppliers`
- **Description:** Retrieves a list of suppliers associated with a specific product, along with their costs.
- **Success Response (200 OK):** `ListResponse<ProductSupplierDTO>`
  ```json
  // Example item in response data
  {
    "supplier_id": "supplier-uuid-456",
    "supplier_name": "Textiles ABC",
    "cost_per_variant": [
      {
        "variant_id": "variant-uuid-789",
        "sku": "TSHIRT-L-BLU",
        "cost": 12.50,
        "currency": "EUR"
      }
    ]
  }
  ```
