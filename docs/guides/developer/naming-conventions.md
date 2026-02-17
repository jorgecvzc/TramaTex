# Convenciones de Nomenclatura del Proyecto TramaTex

Este documento detalla las convenciones de nomenclatura para archivos, carpetas, paquetes, y elementos de código en el proyecto TramaTex. Estas convenciones aseguran la consistencia y legibilidad del código, complementando la decisión arquitectónica definida en [ADR-009: Estructura de Carpetas y Organización del Proyecto](../../architecture/adrs/ADR-009-project-structure.md).

---

## Regla Global de Nomenclatura

Para mantener la consistencia y la compatibilidad con herramientas internacionales, se establece la siguiente política:

-   **Nombres de Archivos y Carpetas:** **Inglés**, usando `kebab-case` para archivos de documentación (ej: `01-initial-architecture.md`) y el case apropiado para código fuente (ej: `userRepository.go`, `UserCard.vue`).
-   **Contenido de Archivos:** **Español** para toda la documentación (`/docs`), **Inglés** para todo el código fuente, comentarios y mensajes de commit.

---

## Backend (Go)

### Archivos

-   **Snake case:** `party_repository.go`, `pricing_engine.go`
-   **Tests:** sufijo `_test.go`
-   **Interfaces:** sufijo `_interface.go` solo si no es obvio (ej: `repository.go` es suficiente)

### Paquetes

-   Todo en minúsculas, sin guiones bajos.
-   Singular preferido: `party`, `product`, no `parties`, `products`.

### Entidades y Value Objects

-   **PascalCase:** `Party`, `Customer`, `Money`.

### Funciones

-   **PascalCase (exportadas):** `CreateParty`, `CalculatePrice`.
-   **camelCase (privadas):** `validateNIF`, `applyDiscount`.

---

## Frontend (Vue.js)

### Archivos

-   **PascalCase para componentes:** `PartyList.vue`, `ProductCard.vue`.
-   **camelCase para servicios/stores:** `partyService.js`, `auth.js`.

### Componentes

-   Siempre multi-palabra: `PartyList` (bien), `List` (mal).
-   Prefijos descriptivos: `TheHeader`, `BaseButton`, `AppSidebar`.

### Stores (Pinia)

-   **Singular:** `party.js`, no `parties.js`.
-   **Export default:** `usePartyStore`, `useProductStore`.

---

## Base de Datos (PostgreSQL)

### Tablas

-   **Snake case, plural:** `parties`, `products`, `orders`.

### Columnas

-   **Snake case:** `party_id`, `created_at`, `supplier_cost`.

### Índices

-   **Prefijo `idx_`:** `idx_parties_nif`, `idx_orders_customer_id`.

### Constraints

-   **Foreign keys:** `fk_[tabla]_[columna]`.
-   **Unique:** `uk_[tabla]_[columna]`.
-   **Check:** `ck_[tabla]_[condicion]`.
