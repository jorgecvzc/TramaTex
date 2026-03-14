# Contratos de API - Módulo Sales

Este documento detalla los puntos de integración para la gestión del flujo documental de ventas en TramaTex.

---

## 1. Gestión de Presupuestos (`/api/sales/quotes`)

Puntos de entrada para la fase de preventa.
- **Creación y Edición:** Permite gestionar líneas de productos, precios manuales y fechas de expiración.
- **Transición de Estado:** Endpoint para aceptar o rechazar la cotización.
- **Conversión Automática (`POST /api/sales/quotes/{id}/convert`):** Transforma una cotización aceptada en un pedido en firme, cerrando el presupuesto original.

## 2. Gestión de Pedidos (`/api/sales/orders`)

Eje central de la integración operativa.
- **Líneas de Pedido:** Gestión de artículos y cantidades. Permite asociar una fecha de entrega (`deliveryDate`) para coordinar con el taller.
- **Sincronización con MES:** El sistema permite consultar qué órdenes de trabajo en producción están vinculadas a un pedido específico.

## 3. Emisión Fiscal (`/api/sales/invoices`)

- **Facturación Completa B2B (`POST /api/sales/invoices`):** Requiere identificación completa del cliente. En el MVP se genera exclusivamente desde albaranes entregados (no desde pedidos). Vincula cada línea de albarán con su línea de factura correspondiente (`invoice_line_item_id`).
- **Factura Simplificada / Ticket (`POST /api/sales/invoices/simplified`):** Optimizado para venta directa TPV/POS sin requerir pedido ni albarán previo. Admite cualquier Party con rol `CLIENT` (por defecto `CONSUMIDOR_FINAL`). Valida que el importe total no exceda los límites legales (3.000€). El precio unitario se calcula usando `BaseSalesPrice` (precio de catálogo) y el descuento se aplica explícitamente por línea.
  - **Roles requeridos:** `admin`, `commercial`, `cashier`.
  - **Request body:**
    ```json
    {
      "partyId": "uuid",
      "invoiceDate": "2026-03-14T...",
      "items": [
        { "productVariantId": "uuid", "quantity": 2, "discountPercent": 5.0 }
      ]
    }
    ```
  - **Response:** `Invoice` completa (tipo SIMPLIFICADA) con líneas, totales e IVA.
- **Listado con filtro por albarán (`GET /api/sales/invoices?deliveryNoteId={id}`):** Permite consultar la factura asociada a un albarán concreto. La respuesta incluye `deliveryNoteIds` con los IDs de los albaranes relacionados.

## 4. Logística y Albaranes (`/api/sales/delivery-notes`)

- **Confirmación de Salida:** Puntos de entrada para registrar las entregas físicas de mercancía vinculadas a un pedido.

---

## Estructura de Montos (`MoneyDTO`)

Para asegurar la integridad económica en todos los documentos de venta, todos los campos de precio y descuento se comunican como objetos de valor con:
- `amount`: Valor numérico de precisión (ej. 10.00).
- `currency`: Código de moneda (siempre "EUR").

---
**Última Actualización:** 2026-03-14
