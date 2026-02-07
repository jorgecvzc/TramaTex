# Sales Module API Contracts

Este documento especifica los contratos de la API para el módulo de Ventas, abarcando la gestión de cotizaciones, pedidos, albaranes y facturas. Se basa en el modelo de dominio consolidado y los casos de uso definidos.

---

### **DTOs Comunes**

1.  **`MoneyDTO`** (Reutilizado del módulo Pricing)
    *   `amount`: DECIMAL
    *   `currency`: String (fijo a "EUR" para MVP)

---

### **1. Gestión de `Quote` (Presupuesto / Cotización)**

#### **1.1. `QuoteLineItemRequest`**
*   `productVariantId`: UUID
*   `quantity`: Integer
*   `manualUnitPrice`: `MoneyDTO` (opcional, para sobrescribir)
*   `manualDiscountPerUnit`: `MoneyDTO` (opcional, para sobrescribir)

#### **1.2. `QuoteLineItemResponse`**
*   `id`: UUID
*   `productVariantId`: UUID
*   `quantity`: Integer
*   `calculatedUnitPrice`: `MoneyDTO`
*   `manualUnitPrice`: `MoneyDTO` (nullable)
*   `finalUnitPrice`: `MoneyDTO`
*   `calculatedDiscountPerUnit`: `MoneyDTO` (nullable)
*   `manualDiscountPerUnit`: `MoneyDTO` (nullable)
*   `finalDiscountPerUnit`: `MoneyDTO` (nullable)
*   `subtotal`: `MoneyDTO`

#### **1.3. `CreateQuoteRequest` (CU-S-001)**
*   `partyId`: UUID
*   `expirationDate`: Timestamp (ISO 8601 string)
*   `notes`: String (nullable)
*   `items`: List of `QuoteLineItemRequest`

#### **1.4. `UpdateQuoteRequest` (CU-S-004)**
*   `expirationDate`: Timestamp (ISO 8601 string, opcional)
*   `notes`: String (opcional, nullable)
*   `items`: List of `QuoteLineItemRequest` (opcional, para reemplazar todas las líneas)

#### **1.5. `ChangeQuoteStatusRequest` (CU-S-005)**
*   `newStatus`: String (ENUM: `BORRADOR`, `ENVIADA`, `APROBADA`, `RECHAZADA`, `EXPIRADA`)

#### **1.6. `QuoteResponse` (CU-S-002, CU-S-003, CU-S-004, CU-S-005)**
*   `id`: UUID
*   `quoteNumber`: String
*   `partyId`: UUID
*   `quoteDate`: Timestamp (ISO 8601 string)
*   `expirationDate`: Timestamp (ISO 8601 string)
*   `status`: String (ENUM)
*   `lineItems`: List of `QuoteLineItemResponse`
*   `subtotal`: `MoneyDTO`
*   `taxAmount`: `MoneyDTO`
*   `total`: `MoneyDTO`
*   `notes`: String (nullable)

---

### **2. Gestión de `SalesOrder` (Pedido de Venta)**

#### **2.1. `OrderLineItemRequest`**
*   `productVariantId`: UUID
*   `quantity`: Integer
*   `manualUnitPrice`: `MoneyDTO` (opcional)
*   `manualDiscountPerUnit`: `MoneyDTO` (opcional)

#### **2.2. `OrderLineItemResponse`**
*   `id`: UUID
*   `productVariantId`: UUID
*   `quantity`: Integer
*   `calculatedUnitPrice`: `MoneyDTO`
*   `manualUnitPrice`: `MoneyDTO` (nullable)
*   `finalUnitPrice`: `MoneyDTO`
*   `calculatedDiscountPerUnit`: `MoneyDTO` (nullable)
*   `manualDiscountPerUnit`: `MoneyDTO` (nullable)
*   `finalDiscountPerUnit`: `MoneyDTO` (nullable)
*   `subtotal`: `MoneyDTO`

#### **2.3. `CreateOrderRequest` (CU-S-007)**
*   `partyId`: UUID
*   `quoteId`: UUID (opcional, si proviene de una Quote)
*   `deliveryDate`: Timestamp (ISO 8601 string)
*   `notes`: String (nullable)
*   `items`: List of `OrderLineItemRequest` (opcional, si no viene de Quote)

#### **2.4. `UpdateOrderDetailsRequest` (CU-S-010)**
*   `deliveryDate`: Timestamp (ISO 8601 string, opcional)
*   `notes`: String (opcional, nullable)
*   `partyId`: UUID (opcional, si el estado lo permite)

#### **2.5. `ChangeOrderStatusRequest` (CU-S-011)**
*   `newStatus`: String (ENUM: `PENDIENTE`, `EN_PREPARACION`, `ENTREGADO_PARCIALMENTE`, `ENTREGADO`, `CANCELADO`, `FACTURADO_PARCIALMENTE`, `FACTURADO_COMPLETAMENTE`)

#### **2.6. `SalesOrderResponse` (CU-S-008, CU-S-009, CU-S-010, CU-S-011)**
*   `id`: UUID
*   `orderNumber`: String
*   `quoteId`: UUID (nullable)
*   `partyId`: UUID
*   `orderDate`: Timestamp (ISO 8601 string)
*   `deliveryDate`: Timestamp (ISO 8601 string)
*   `status`: String (ENUM)
*   `lineItems`: List of `OrderLineItemResponse`
*   `subtotal`: `MoneyDTO`
*   `taxAmount`: `MoneyDTO`
*   `total`: `MoneyDTO`
*   `notes`: String (nullable)

---

### **3. Gestión de `DeliveryNote` (Albarán)**

#### **3.1. `DeliveryNoteLineItemRequest`**
*   `salesOrderLineItemId`: UUID
*   `deliveredQuantity`: Integer

#### **3.2. `DeliveryNoteLineItemResponse`**
*   `id`: UUID
*   `salesOrderLineItemId`: UUID
*   `productVariantId`: UUID
*   `deliveredQuantity`: Integer

#### **3.3. `CreateDeliveryNoteRequest` (CU-S-015)**
*   `salesOrderId`: UUID
*   `deliveryDate`: Timestamp (ISO 8601 string)
*   `notes`: String (nullable)
*   `items`: List of `DeliveryNoteLineItemRequest`

#### **3.4. `DeliveryNoteResponse` (CU-S-016, CU-S-017)**
*   `id`: UUID
*   `deliveryNoteNumber`: String
*   `salesOrderId`: UUID
*   `partyId`: UUID
*   `deliveryDate`: Timestamp (ISO 8601 string)
*   `status`: String (ENUM)
*   `lineItems`: List of `DeliveryNoteLineItemResponse`
*   `notes`: String (nullable)

---

### **4. Gestión de `Invoice` (Factura)**

#### **4.1. `InvoiceLineItemResponse`**
*   `id`: UUID
*   `salesOrderLineItemId`: UUID (nullable)
*   `productVariantId`: UUID
*   `quantity`: Integer
*   `unitPrice`: `MoneyDTO`
*   `discountAmount`: `MoneyDTO` (nullable)
*   `subtotal`: `MoneyDTO`
*   `taxAmount`: `MoneyDTO` (nullable)
*   `total`: `MoneyDTO`

#### **4.2. `CreateInvoiceRequest` (CU-S-018)**
*   `partyId`: UUID
*   `salesOrderIds`: List of UUIDs (opcional, para facturar uno o varios pedidos)
*   `deliveryNoteIds`: List of UUIDs (opcional, para facturar uno o varios albaranes)
*   `invoiceDate`: Timestamp (ISO 8601 string)
*   `dueDate`: Timestamp (ISO 8601 string)
*   `paymentTerms`: String (nullable)

#### **4.3. `InvoiceResponse` (CU-S-019, CU-S-020)**
*   `id`: UUID
*   `invoiceNumber`: String
*   `partyId`: UUID
*   `invoiceDate`: Timestamp (ISO 8601 string)
*   `dueDate`: Timestamp (ISO 8601 string)
*   `status`: String (ENUM)
*   `lineItems`: List of `InvoiceLineItemResponse`
*   `relatedOrderIds`: List of UUIDs (nullable)
*   `subtotal`: `MoneyDTO`
*   `taxAmount`: `MoneyDTO`
*   `total`: `MoneyDTO`
*   `paymentTerms`: String (nullable)