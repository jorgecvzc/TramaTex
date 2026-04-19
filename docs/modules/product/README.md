# 👕 Módulo de Product (Catálogo y Variantes)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Vigente |
| **Referencia** | [ADR-015](../../architecture/adrs/adr-015-product-module-architecture.md) |

---

## 🎯 Propósito
Este módulo gestiona el corazón del inventario textil: el catálogo de productos y sus variantes. Implementa una lógica de creación **Just-In-Time (JIT)** que permite manejar miles de combinaciones de tallas, colores y materiales sin necesidad de darlas de alta manualmente.

---

## 📄 Ramas del Conocimiento (Documentación)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md) — Atributos, Productos y Variantes JIT.
*   **Casos de Uso:** [use-cases.md](./use-cases.md) — Definición de atributos y búsqueda de variantes.
*   **Contratos de API:** [api-contracts.md](./api-contracts.md) — Definición de SKUs y metadatos de producto.
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md) — Detalles sobre el motor de resolución de variantes.
*   **Especificación:** [module-spec.md](./module-spec.md) — Reglas del catálogo modular.

---

## 🏗️ Componentes Clave
*   **Entidades:** `Product` (Base), `ProductVariant` (Instancia vendible), `Attribute` / `AttributeValue`.
*   **Metodología JIT:** Generación automática de variantes según la demanda de ventas.
*   **SKUs Dinámicos:** Generación de identificadores únicos basados en atributos.

---
[Volver al Resumen de Módulos](../README.md)
