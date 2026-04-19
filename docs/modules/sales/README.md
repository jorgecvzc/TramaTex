# 🤝 Módulo de Sales (Ciclo Comercial)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Vigente |
| **Referencia** | [ADR-017](../../architecture/adrs/adr-017-sales-module-architecture.md), [ADR-020](../../architecture/adrs/adr-020-tickets-and-invoice-series.md) |

---

## 🎯 Propósito
Este módulo gestiona el flujo documental completo de la venta, desde el presupuesto inicial hasta la facturación fiscal. Orquesta la interacción entre clientes, productos y precios, garantizando la trazabilidad total de la operación y el cumplimiento legal (AEAT).

---

## 📄 Ramas del Conocimiento (Documentación)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md) — Ciclo de vida de Quote, Order, DeliveryNote e Invoice.
*   **Casos de Uso:** [use-cases.md](./use-cases.md) — Conversión de documentos y flujos de TPV.
*   **Contratos de API:** [api-contracts.md](./api-contracts.md) — Endpoints comerciales y DTOs de venta.
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md) — Lógica de transiciones y cálculos.

---

## 🏗️ Componentes Clave
*   **Agregados:** `Quote` (Presupuesto), `SalesOrder` (Pedido), `DeliveryNote` (Albarán), `Invoice` (Factura/Ticket).
*   **Lógica Comercial:** Manual Override de precios, gestión de series de facturación y límites legales para facturas simplificadas.
*   **Trazabilidad:** Rastro íntegro desde la oferta hasta el cobro.

---
[Volver al Resumen de Módulos](../README.md)
