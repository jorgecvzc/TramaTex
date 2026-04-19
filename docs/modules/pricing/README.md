# 💰 Módulo de Pricing (Motor de Precios)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Vigente |
| **Referencia** | [ADR-016](../../architecture/adrs/adr-016-pricing-module-architecture.md) |

---

## 🎯 Propósito
Este módulo encapsula la lógica económica del sistema. Calcula de forma dinámica los precios base de venta (coste + margen) y aplica modificaciones en tiempo real (descuentos por cliente o volumen). Utiliza una caché de alto rendimiento para garantizar una experiencia de usuario fluida en el punto de venta.

---

## 📄 Ramas del Conocimiento (Documentación)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md) — Reglas de precio y precisión decimal.
*   **Casos de Uso:** [use-cases.md](./use-cases.md) — Cálculo de precios finales y gestión de márgenes.
*   **Contratos de API:** [api-contracts.md](./api-contracts.md) — DTOs financieros y contratos del motor.
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md) — Detalles sobre la integración con Redis.
*   **Integración:** [integration.md](./integration.md) — Relación con Sales y Product.

---

## 🏗️ Componentes Clave
*   **Entidades:** `BaseSalesPriceRule`, `SaleModificationRule`.
*   **Objetos de Valor:** `Money` (Precisión Decimal), `Percentage`.
*   **Infraestructura:** Caché distribuida mediante **Redis** para optimización de cálculos.

---
[Volver al Resumen de Módulos](../README.md)
