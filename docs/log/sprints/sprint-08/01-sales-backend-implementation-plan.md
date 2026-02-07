# Tarea 01 - Plan de implementacion backend Sales

## 📋 INFORMACION DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 01 |
| **ID de Sprint** | sprint-08 |
| **Titulo** | Plan de implementacion backend del modulo Sales |
| **Estado** | ✅ Completado |
| **Facilitador/LLM** | GitHub Copilot |
| **Fecha de Inicio** | 2026-02-07 |
| **Fecha de Fin** | 2026-02-07 |
| **Duracion Estimada** | 6-10 horas |
| **Duracion Real** | 1 hora |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [x] **Objetivo 1:** Definir alcance MVP del backend de Sales y sus dependencias
   - Alinear contratos de API con casos de uso
   - Determinar integraciones con Party, Product, Pricing
   - Definir reglas de estado y transiciones

2. [x] **Objetivo 2:** Especificar plan tecnico por capas (domain, application, infrastructure, interfaces)
   - Enumerar entidades, VOs y servicios de dominio
   - Listar comandos/queries y DTOs de aplicacion
   - Definir repositorios y modelos GORM
   - Mapear endpoints y handlers

3. [x] **Objetivo 3:** Establecer estrategia de tests y cobertura
   - Definir tests de dominio y casos criticos
   - Definir tests de integracion (repos y use cases)
   - Definir criterios de cobertura minima

---

## 📊 CONTEXTO DE ENTRADA

### Estado Anterior

**Ultima tarea completada:** sprint-06 (Product backend en progreso)

**Cambios desde ultima tarea:**
- Modulo Sales definido en documentos de modulo (dominio, casos de uso, contratos)
- Decidido iniciar desarrollo de Sales

**Estado en project-status.md:**
- Fase actual: Fase 1 (Dominio Base) en progreso

### Bloqueadores/Dependencias

- [x] Dependencia 1: Integracion con Pricing confirmada (aplicacion in-process)
- [x] Dependencia 2: Estados alineados con migraciones y dominio
- [x] Riesgo 1: Descuentos manuales y overrides cubiertos en logica de precios

### Prioridades para esta Tarea

**Critica (Must Have):**
- Definir plan de migraciones y modelo de datos
- Definir plan de domain y application
- Definir plan de endpoints

**Alta (Should Have):**
- Definir estrategia de tests y fixtures

**Media (Nice to Have):**
- Proponer ADR si se requiere decision arquitectonica

---

## 🛠️ PLAN DE TRABAJO

### Fase 1: Analisis Arquitectonico (30-60 min)

- [x] Identificar Bounded Contexts involucrados (Sales, Party, Product, Pricing)
- [x] Revisar ADRs aplicables y contratos de Sales
- [x] Definir dependencias y limites entre modulos
- [x] Confirmar alcance MVP (prioridad SalesOrder y Quote)

**Notas:**
```
Bounded contexts: Sales (owner), Party (cliente), Product (variant), Pricing (prices/discounts).
Integraciones: Sales calcula precios via Pricing application service, no via HTTP.
MVP Sales: Quotes + SalesOrders completos; DeliveryNotes e Invoices en modo basico (create/get/list).
Reglas de estado:
- Quote: BORRADOR -> ENVIADA -> APROBADA/RECHAZADA/EXPIRADA; APROBADA -> CONVERTIDA_A_PEDIDO (interno).
- SalesOrder: PENDIENTE -> EN_PREPARACION -> ENTREGADO_PARCIALMENTE/ENTREGADO -> FACTURADO_*
Riesgo: tipo de PartyID en DB (VARCHAR(36) vs UUID). Decidir tipo y adaptar FKs.
```

### Fase 2: Plan de Migraciones y Modelo de Datos

**Backend – DB (60-90 min):**
- [x] Diseñar tablas para Quote y QuoteLineItem
- [x] Diseñar tablas para SalesOrder y OrderLineItem
- [x] Diseñar tablas para DeliveryNote y DeliveryNoteLineItem
- [x] Diseñar tablas para Invoice y InvoiceLineItem
- [x] Definir enums y constraints
- [x] Definir indices y FKs necesarios

**Notas:**
```
Tablas y enums:
- enums: quote_status, sales_order_status, delivery_note_status, invoice_status.
- quotes: id, quote_number, party_id, quote_date, expiration_date, status, subtotal, tax_amount, total, notes, created_at, updated_at, deleted_at.
- quote_line_items: id, quote_id, product_variant_id, quantity, calculated_unit_price_amount, calculated_unit_price_currency,
  manual_unit_price_amount, manual_unit_price_currency, final_unit_price_amount, final_unit_price_currency,
  calculated_discount_per_unit_amount, calculated_discount_per_unit_currency, manual_discount_per_unit_amount, manual_discount_per_unit_currency,
  final_discount_per_unit_amount, final_discount_per_unit_currency, subtotal_amount, subtotal_currency.
- sales_orders: id, order_number, quote_id (nullable), party_id, order_date, delivery_date, status,
  subtotal, tax_amount, total, notes, created_at, updated_at, deleted_at.
- order_line_items: id, sales_order_id, product_variant_id, quantity, calculated_unit_price_amount, calculated_unit_price_currency,
  manual_unit_price_amount, manual_unit_price_currency, final_unit_price_amount, final_unit_price_currency,
  calculated_discount_per_unit_amount, calculated_discount_per_unit_currency, manual_discount_per_unit_amount, manual_discount_per_unit_currency,
  final_discount_per_unit_amount, final_discount_per_unit_currency, subtotal_amount, subtotal_currency.
- delivery_notes: id, delivery_note_number, sales_order_id, party_id, delivery_date, status, notes, created_at, updated_at, deleted_at.
- delivery_note_line_items: id, delivery_note_id, sales_order_line_item_id, product_variant_id, delivered_quantity.
- invoices: id, invoice_number, party_id, invoice_date, due_date, status, payment_terms,
  subtotal, tax_amount, total, created_at, updated_at, deleted_at.
- invoice_line_items: id, invoice_id, sales_order_line_item_id (nullable), product_variant_id, quantity,
  unit_price_amount, unit_price_currency, discount_amount, discount_currency, subtotal_amount, subtotal_currency,
  tax_amount, tax_currency, total_amount, total_currency.

Indices:
- quotes: idx_quotes_party_id, idx_quotes_status, idx_quotes_quote_date.
- sales_orders: idx_sales_orders_party_id, idx_sales_orders_status, idx_sales_orders_order_date.
- delivery_notes: idx_delivery_notes_sales_order_id, idx_delivery_notes_party_id.
- invoices: idx_invoices_party_id, idx_invoices_status, idx_invoices_invoice_date.

FKs:
- quote_line_items.quote_id -> quotes.id.
- order_line_items.sales_order_id -> sales_orders.id.
- delivery_notes.sales_order_id -> sales_orders.id.
- delivery_note_line_items.delivery_note_id -> delivery_notes.id.
- invoice_line_items.invoice_id -> invoices.id.

Nota Party/Product:
- party_id: definir si UUID o VARCHAR(36) para compatibilidad con Party.
- product_variant_id: UUID (Product).
```

### Fase 3: Plan de Dominio (Domain)

- [x] Definir entidades y VOs a implementar
- [x] Definir invariantes y reglas de estado
- [x] Definir metodos de calculo de totales
- [x] Definir generacion de numeros (OrderNumber, QuoteNumber, etc.)

**Notas:**
```
Entidades y VOs:
- Value Objects: Money, Percentage, QuoteNumber, OrderNumber, DeliveryNoteNumber, InvoiceNumber.
- Quote (AR) con QuoteLineItem; metodos: AddLineItem, UpdateLineItem, RemoveLineItem, ChangeStatus, ConvertToOrder, RecalculateTotals.
- SalesOrder (AR) con OrderLineItem; metodos: AddLineItem, UpdateLineItem, RemoveLineItem, ChangeStatus, RecalculateTotals.
- DeliveryNote (AR) con DeliveryNoteLineItem; metodo: RegisterDelivery.
- Invoice (AR) con InvoiceLineItem; metodo: MarkPaid.

Invariantes:
- Totales siempre calculados con FinalUnitPrice y FinalDiscountPerUnit.
- Manual overrides: si manual_* existe, se usa como final; si no, usar calculated_*.
- Transiciones de estado validas segun enums definidos.

Generadores de numeros:
- QuoteNumber: QTE-YYYYMMDD-XXXX
- OrderNumber: ORD-YYYYMMDD-XXXX
- DeliveryNoteNumber: DN-YYYYMMDD-XXXX
- InvoiceNumber: INV-YYYYMMDD-XXXX
```

### Fase 4: Plan de Aplicacion (Use Cases)

- [x] Mapear CU-S-001..020 a comandos/queries
- [x] Definir DTOs de entrada/salida
- [x] Definir orquestacion con Pricing
- [x] Definir validaciones de estado y transiciones

**Notas:**
```
Mapeo de casos de uso a comandos/queries (nombres sugeridos):
- CU-S-001 CreateQuote -> CreateQuoteCommand
- CU-S-002 GetQuote -> GetQuoteQuery
- CU-S-003 ListQuotes -> ListQuotesQuery
- CU-S-004 UpdateQuote -> UpdateQuoteCommand
- CU-S-005 ChangeQuoteStatus -> ChangeQuoteStatusCommand
- CU-S-006 ConvertQuoteToOrder -> ConvertQuoteToOrderCommand

- CU-S-007 CreateOrder -> CreateOrderCommand
- CU-S-008 GetOrder -> GetOrderQuery
- CU-S-009 ListOrders -> ListOrdersQuery
- CU-S-010 UpdateOrderDetails -> UpdateOrderDetailsCommand
- CU-S-011 ChangeOrderStatus -> ChangeOrderStatusCommand
- CU-S-012 AddLineItemToOrder -> AddOrderLineItemCommand
- CU-S-013 UpdateOrderLineItem -> UpdateOrderLineItemCommand
- CU-S-014 RemoveOrderLineItem -> RemoveOrderLineItemCommand

- CU-S-015 CreateDeliveryNote -> CreateDeliveryNoteCommand
- CU-S-016 GetDeliveryNote -> GetDeliveryNoteQuery
- CU-S-017 ListDeliveryNotes -> ListDeliveryNotesQuery

- CU-S-018 CreateInvoiceFromOrder -> CreateInvoiceCommand
- CU-S-019 GetInvoice -> GetInvoiceQuery
- CU-S-020 ListInvoices -> ListInvoicesQuery

Integracion Pricing:
- PricingClient interface en Sales infra: CalculateBaseSalesPrice, CalculateFinalSalePrice.
- Application service compone precios por item y llena calculated_*.
```

### Fase 5: Plan de Infraestructura

- [x] Definir repositorios y data models GORM
- [x] Definir interfaces de servicios externos (Pricing/Party/Product)
- [x] Proponer estructura de carpetas

**Notas:**
```
Repositorios (domain interfaces):
- QuoteRepository, SalesOrderRepository, DeliveryNoteRepository, InvoiceRepository.

GORM data models:
- QuoteDataModel, QuoteLineItemDataModel
- SalesOrderDataModel, OrderLineItemDataModel
- DeliveryNoteDataModel, DeliveryNoteLineItemDataModel
- InvoiceDataModel, InvoiceLineItemDataModel

Servicios externos:
- PricingClient (infra) con implementacion directa a pricing application (in-process).
- Party/Product lookups opcionales (validacion de ids).

Estructura de carpetas:
apps/tramatex-api/internal/sales/
- domain/
- application/
- infrastructure/persistence/
- interfaces/http/handler/
```

### Fase 6: Plan de Interfaces HTTP

- [x] Mapear endpoints a handlers segun contratos
- [x] Definir request/response DTOs
- [x] Definir manejo de errores y codigos HTTP

**Notas:**
```
Endpoints segun contratos:

Quotes:
- POST /quotes
- GET /quotes/{id}
- GET /quotes
- PUT /quotes/{id}
- PATCH /quotes/{id}/status
- POST /quotes/{id}/convert-to-order

Orders:
- POST /orders
- GET /orders/{id}
- GET /orders
- PATCH /orders/{id}
- PATCH /orders/{id}/status
- POST /orders/{id}/line-items
- PATCH /orders/{id}/line-items/{lineItemId}
- DELETE /orders/{id}/line-items/{lineItemId}

DeliveryNotes:
- POST /delivery-notes
- GET /delivery-notes/{id}
- GET /delivery-notes

Invoices:
- POST /invoices
- GET /invoices/{id}
- GET /invoices

Handlers:
- QuoteHandler, SalesOrderHandler, DeliveryNoteHandler, InvoiceHandler.
```

### Fase 7: Plan de Tests y Cobertura

- [x] Tests unitarios de dominio (Price/Discount, totales, transiciones)
- [x] Tests de integracion de repositorios
- [x] Tests de application (use cases criticos)
- [x] Metas de cobertura y criterios de aceptacion

**Notas:**
```
Tests dominio:
- Money/Percentage validation and arithmetic.
- Quote and SalesOrder totals with manual overrides.
- Status transition validation (Quote, Order, DeliveryNote, Invoice).

Tests application:
- CreateQuote/CreateOrder with pricing client stub.
- ConvertQuoteToOrder transitions and data mapping.
- Add/Update/Remove line items updates totals.

Tests infraestructura:
- GORM repository CRUD for each aggregate.

Cobertura objetivo:
- Sales domain >= 85%, resto >= 75%.
```

---

## ✅ RESULTADO

- Plan ejecutado y documentado en la implementacion del backend Sales.
- Migraciones, capas (domain/application/infra), handlers y tests completados en la tarea 02.
- Ver detalles y lista de archivos en [docs/log/sprints/sprint-08/02-sales-backend-implementation.md](docs/log/sprints/sprint-08/02-sales-backend-implementation.md).

---

## 📝 CHANGES MADE

### Commits Realizados

```
[Sin commits]
```

### Archivos Modificados

| Archivo | Tipo | Descripcion |
|---------|------|-------------|
| docs/log/sprints/sprint-08/01-sales-backend-implementation-plan.md | NEW | Tarea de planificacion de Sales |

---

## ✅ DEFINICION DE "HECHO"

La tarea se considera completada cuando:

- [x] Plan de backend Sales definido por capas
- [x] Alcance MVP y dependencias documentados
- [x] Plan de migraciones y modelo de datos definido
- [x] Plan de endpoints definido
- [x] Estrategia de tests definida

---

## 🚨 BLOQUEADORES/PROBLEMAS ENCONTRADOS

**Sin bloqueadores por ahora.**
