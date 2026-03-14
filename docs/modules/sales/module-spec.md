# Módulo de Sales (Gestión Comercial)

**Estado:** ✅ **COMPLETO (100%)**  
**Última actualización:** 14 de marzo de 2026

## Estado de Implementación

### Componentes Completos ✅
- **Presupuestos (Quotes):**
  - Backend: CRUD completo, gestión de estados, conversión automática a pedido.
  - Frontend: `QuoteList.vue`, `QuoteDetail.vue` (con acciones y avisos de expiración) y `QuoteCreate.vue`.
  - Sistema de cálculo de totales en tiempo real.
  
- **Pedidos (Orders):**
  - Backend: Gestión completa del flujo de pedidos.
  - Frontend: `OrderList.vue` (optimización batch), `OrderDetail.vue` (con creación de albaranes parciales/totales).
  
- **Albaranes (Delivery Notes):**
  - Backend: Soporte para albaranes parciales/totales.
  - Frontend: `DeliveryNoteList.vue`, `DeliveryNoteDetail.vue` con firma y linkage a pedido.
  
- **Facturación (Invoices & Tickets):**
  - Backend: Facturación completa (B2B) y simplificada (Tickets B2C).
  - Cumplimiento AEAT con series diferenciadas (ADR-020).
  - Frontend: Listado y detalle de facturas y creación rápida de tickets.

---

## 1. Propósito

*   **Visión del Módulo:** Gestionar el ciclo comercial completo desde la cotización hasta la facturación definitiva, integrando el motor de precios y asegurando la trazabilidad.
*   **Objetivos Clave:**
    *   Proporcionar un sistema completo para gestionar el ciclo de vida de las ventas.
    *   Manejar desde la cotización hasta la entrega de la orden y su cobro.

---

## 2. Requisitos

... [resto de requisitos permanecen igual] ...

---

## 7. Decisiones de Diseño

*   **Numeración Secuencial de Documentos:**
    Todos los documentos comerciales usan numeración secuencial por prefijo y año fiscal.
    Formato: `PREFIJO-AÑO-SECUENCIAL` (ej: `PRE-2026-0001`).

    | Documento | Prefijo | Ejemplo |
    |-----------|---------|---------|
    | Presupuesto (Quote) | `PRE` | `PRE-2026-0001` |
    | Pedido (SalesOrder) | `PED` | `PED-2026-0003` |
    | Albarán (DeliveryNote) | `ALB` | `ALB-2026-0001` |
    | Factura de Venta | `FV` | `FV-2026-0012` |
    | Factura de Ticket | `FT` | `FT-2026-0005` |

    Los contadores se almacenan en la tabla `document_sequences` y se incrementan
    atómicamente con `INSERT ... ON CONFLICT DO UPDATE`. Las facturas usan series
    diferenciadas (`FV` para ventas, `FT` para tickets) para cumplimiento AEAT,
    con extensibilidad para futuras series (ej: `FR` para rectificativas).

*   **Trazabilidad Factura↔Albarán (Facturas Completas B2B):**
    Las facturas completas (B2B) se generan **exclusivamente desde albaranes** (1 factura por albarán). No se permite crear facturas completas directamente desde pedidos. La trazabilidad se mantiene mediante el campo `invoice_line_item_id` en `delivery_note_line_items`, que vincula cada línea de albarán con su línea de factura correspondiente (relación N:1 preparada para consolidación Post-MVP).

    Los **tickets (facturas simplificadas)** siguen un flujo independiente: se crean directamente desde la interfaz TPV/POS sin necesidad de pedido ni albarán previo (ver CU-S-019).

    | Funcionalidad | MVP | Post-MVP |
    |---|---|---|
    | Factura completa desde albarán individual | ✅ | ✅ |
    | Factura completa desde pedido | ❌ | ❌ (siempre vía albaranes) |
    | Ticket (simplificada) — venta directa TPV | ✅ | ✅ |
    | Consolidar albaranes de un cliente en 1 factura | ❌ | ✅ |
    | Consolidar líneas iguales (producto+precio+dto) | ❌ (1:1) | ✅ (N:1) |

*   **Flujo de Estados de Documentos:**
    *   **Presupuesto (Quote):** `BORRADOR` (Draft) -> `EMITIDA` (Issued) -> `APROBADA` (Approved) -> `CONVERTIDA_A_PEDIDO` (Converted).
    *   **Pedido (SalesOrder):** `PENDIENTE` (Pending) -> `EN_PREPARACION` (In Preparation) -> `ENTREGADO` (Delivered) -> `FACTURADO_COMPLETAMENTE` (Invoiced).
    *   **Factura (Invoice):** `BORRADOR` (Draft) -> `EMITIDA` (Issued) -> `PAGADA` (Paid).
*   **Relaciones con Otros Módulos:**
    *   **Party**: La orden pertenece a un `Party` (cliente).
    *   **Product**: Las `LineItems` referencian a un `ProductVariant`.
        *   **⚠️ HIDRATACIÓN DINÁMICA DE DATOS:**
            *   El `ProductName`, `VariantSKU` y `OptionConfiguration` (atributos) **NO se almacenan** físicamente en las líneas de venta.
            *   Estos datos se obtienen en tiempo real del módulo **Product** para asegurar que la visualización siempre esté actualizada con el catálogo.
            *   **Soberanía Comercial:** Solo se congelan el **Precio Unitario**, el **Descuento** y el **Tipo Impositivo** al momento de confirmar el documento.
        *   **⚠️ DEPENDENCIA CRÍTICA - Cálculo de BaseCost:**
            *   El `baseCost` de cada variante se calcula dinámicamente en el módulo Product.
            *   Sales obtiene el `baseCost` actual al momento de crear la línea para calcular el precio final sugerido.
    *   **Pricing**: El `PrecioUnitario` es sugerido por el motor de precios, pero puede ser sobreescrito manualmente (`Manual Override`).
    *   **MES:** El estado `APROBADA` de un presupuesto (al convertirse a pedido) o la creación directa de un pedido disparan la integración con producción.

---

## 8. Fases de Desarrollo

*   [x] **Fase 1 (MVP):** CRUD básico de pedidos y estados.
*   [x] **Fase 2:** Gestión de presupuestos y conversión a pedidos.
*   [x] **Fase 3:** Sistema de facturación completa y simplificada (Tickets).
*   [x] **Fase 4:** Integración automática con producción (MES).
*   [x] **Fase 5 (MVP):** Trazabilidad Factura↔Albarán (campo `invoice_line_item_id` en líneas de albarán, facturación exclusiva desde albaranes).
*   [ ] **Fase 6 (Post-MVP):** Facturación consolidada — agrupar múltiples albaranes de un cliente en una sola factura, con consolidación de líneas iguales (N:1).