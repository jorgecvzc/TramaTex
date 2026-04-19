# 👥 Módulo de Party (Gestión de Terceros)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Vigente |
| **Referencia** | [ADR-012](../../architecture/adrs/adr-012-party-module-architecture.md) |

---

## 🎯 Propósito
Este módulo unifica la gestión de todas las entidades externas que interactúan con el sistema. Utiliza el patrón *Party* para permitir que clientes y proveedores sean tratados como una única entidad base, evitando duplicidades y permitiendo relaciones jerárquicas (matrices/filiales).

---

## 📄 Ramas del Conocimiento (Documentación)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md) — Personas, Organizaciones y Roles.
*   **Casos de Uso:** [use-cases.md](./use-cases.md) — Alta de clientes, gestión de proveedores y jerarquías.
*   **Contratos de API:** [api-contracts.md](./api-contracts.md) — DTOs para la gestión de terceros.
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md) — Estrategia de persistencia de perfiles.
*   **Especificación:** [module-spec.md](./module-spec.md) — Alcance funcional del maestro de datos.

---

## 🏗️ Componentes Clave
*   **Entidades:** `Party` (Raíz), `OrganizationProfile`, `IndividualProfile`.
*   **Roles:** `Customer` (Cliente), `Supplier` (Proveedor).
*   **Jerarquías:** Soporte para relaciones de propiedad entre empresas.

---
[Volver al Resumen de Módulos](../README.md)
