# 🔐 Módulo de IAM (Identity and Access Management)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Vigente |
| **Referencia** | [ADR-014](../../architecture/adrs/adr-014-iam-module-architecture.md) |

---

## 🎯 Propósito
Este módulo es el pilar de seguridad de TramaTex. Se encarga de la gestión de identidades, la autenticación mediante JWT y el control de acceso basado en roles (RBAC), asegurando que cada usuario acceda solo a la información que le corresponde.

---

## 📄 Ramas del Conocimiento (Documentación)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md) — Entidades de Usuario y Rol.
*   **Casos de Uso:** [use-cases.md](./use-cases.md) — Procesos de Login, Registro y gestión administrativa.
*   **Contratos de API:** [api-contracts.md](./api-contracts.md) — Definición de endpoints y DTOs de seguridad.
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md) — Detalles técnicos sobre hashing y tokens.
*   **Especificación:** [module-spec.md](./module-spec.md) — Requisitos funcionales.

---

## 🏗️ Componentes Clave
*   **Entidades:** `User` (Identidad), `Role` (Perfil de acceso).
*   **Objetos de Valor:** `Email`, `HashedPassword`, `RoleName`.
*   **Seguridad:** Implementación stateless mediante JWT (Access & Refresh tokens).

---
[Volver al Resumen de Módulos](../README.md)
