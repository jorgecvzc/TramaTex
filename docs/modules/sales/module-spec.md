# Módulo de Sales (Gestión Comercial)

**Estado:** ✅ **COMPLETO (100%)**  
**Última actualización:** 1 de marzo de 2026

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

*   **Flujo de Estados de Orden:**
    *   `DRAFT` -> `SENT` -> `ACCEPTED` -> `IN_PROGRESS` -> `COMPLETED` -> `INVOICED`
*   **Relaciones con Otros Módulos:**
    *   **Party**: La orden pertenece a un `Party` (cliente).
    *   **Product**: Las `LineItems` referencian a un `ProductVariant`.
        *   **⚠️ DEPENDENCIA CRÍTICA - Cálculo de BaseCost:**
            *   El `baseCost` de cada variante se calcula **dinámicamente** en el módulo Product como: `Product.BasePrice` + modificadores de `AttributeValue`.
            *   Los modificadores pueden ser:
                *   **FIXED**: Suma/resta cantidad fija (€) al precio base.
                *   **PERCENTAGE**: Aplica porcentaje sobre el precio acumulado.
            *   Los modificadores se aplican **secuencialmente** según `Attribute.sortOrder`.
            *   El `baseCost` **NO se almacena** en BD, siempre se calcula on-demand.
            *   Sales debe obtener el `baseCost` actual de cada variante al momento de crear líneas de pedido/presupuesto.
    *   **Pricing**: El `PrecioUnitario` es calculado por el motor de precios a partir del `baseCost` de la variante.
    *   **MES:** El estado 'ACCEPTED' genera automáticamente una orden de producción en MES.

---

## 8. Fases de Desarrollo

*   [x] **Fase 1 (MVP):** CRUD básico de pedidos y estados.
*   [x] **Fase 2:** Gestión de presupuestos y conversión a pedidos.
*   [x] **Fase 3:** Sistema de facturación completa y simplificada (Tickets).
*   [x] **Fase 4:** Integración automática con producción (MES).