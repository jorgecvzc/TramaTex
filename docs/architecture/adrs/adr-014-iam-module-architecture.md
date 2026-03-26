# ADR-014 – Arquitectura del Módulo de IAM

**Fecha:** viernes, 6 de febrero de 2026  
**Estado:** Propuesto  
**Autores:** Gemini CLI

---

## 1. Contexto

El módulo de IAM (Identity and Access Management) es un componente transversal y fundamental en TramaTex, responsable de gestionar la identidad de los usuarios, su autenticación en el sistema, y su autorización para acceder a recursos y realizar acciones. Dado que TramaTex gestiona datos críticos de negocio y opera en un entorno de múltiples roles de usuario, un diseño robusto y seguro del IAM es de máxima prioridad.

**Problemas a Resolver:**
- Proporcionar un mecanismo seguro y eficiente para que los usuarios accedan al sistema.
- Implementar un control de acceso basado en roles (RBAC) para diferentes perfiles de usuario.
- Asegurar la protección de la API contra accesos no autorizados.
- Gestionar la identidad del usuario y sus credenciales de forma segura.
- Ofrecer funcionalidades administrativas para la gestión de usuarios y roles.

---

## 2. Alternativas Consideradas

**Alternativa A – Sistema de Autenticación y Autorización Tradicional (Basado en Sesiones):**
- **Descripción:** Utilizar cookies de sesión gestionadas por el servidor para mantener el estado del usuario.
- **Ventajas:** Simplicidad conceptual para aplicaciones monolíticas tradicionales.
- **Desventajas:** Dificultad para escalar horizontalmente (problemas de afinidad de sesión), complejidad con APIs (CSRF tokens), no ideal para arquitecturas orientadas a microservicios en el futuro (aunque TramaTex es monolito modular, se prepara para ello).

**Alternativa B – Usar un Proveedor de Identidad Externo (OAuth/OpenID Connect):**
- **Descripción:** Delegar la autenticación y autorización a un servicio externo (ej. Google, Auth0, Keycloak).
- **Ventajas:** Alta seguridad, estándares probados, reduce la carga de desarrollo de IAM.
- **Desventajas:** Introduce dependencia externa (internet), mayor complejidad de infraestructura inicial (si es auto-hospedado), posible sobre-ingeniería para un MVP "local-first".

**Alternativa C – Autenticación y Autorización Basada en JWT (Adoptada):**
- **Descripción:** Utilizar JSON Web Tokens (JWT) para `access_token` (corta duración) y `refresh_token` (larga duración) para gestionar sesiones stateless. Implementar RBAC a nivel de aplicación.
- **Ventajas:** Stateless (facilita escalado), estándar de la industria, flexible para diferentes clientes (SPA, móvil), reduce la complejidad de la gestión de sesiones en el servidor. Permite un RBAC a medida.
- **Desventajas:** La revocación de tokens JWT requiere una estrategia (ej. blacklist) si es necesario invalidarlos antes de expirar.

---

## 3. Criterios de Decisión

-   **Seguridad:** Alta protección contra ataques comunes (OWASP Top 10).
-   **Rendimiento:** Autenticación rápida y eficiente.
-   **Escalabilidad:** Capacidad de escalar el sistema sin problemas de sesión.
-   **Flexibilidad:** Adaptable a diferentes tipos de clientes (web, móvil).
-   **Mantenibilidad:** Facilidad para entender, modificar y extender el módulo.
-   **Alineación con Stack:** Consistencia con el stack Go/Vue.
-   **Desarrollo Local-First:** Funcionalidad completa sin depender de servicios externos.

---

## 4. Decisión Adoptada

Se adopta la **Alternativa C: Autenticación y Autorización Basada en JWT**.

### Entidades de Dominio Clave:

1.  **`User`:** Representa a un individuo con acceso al sistema (`UserID`, `Email`, `HashedPassword`, `Role`).
2.  **`Role`:** Define los permisos de un usuario (`RoleID`, `RoleName`).
3.  **`Permission`:** Representa una acción específica (en MVP, implícita en roles).

### Value Objects Clave:

*   `UserID`, `Email`, `HashedPassword`, `RoleID`, `RoleName`, `PermissionID`, `PermissionName`.

### Aspectos Clave del Diseño:

*   **Autenticación JWT Stateless:** Uso de `access_token` y `refresh_token` con expiración.
*   **Control de Acceso Basado en Roles (RBAC):** Roles predefinidos (`admin`, `commercial`, `designer`, `workshop`) que determinan el acceso a las funcionalidades. La autorización se verifica en la capa de aplicación/interfaces.
*   **Gestión Segura de Contraseñas:** Hashing con `bcrypt` y validación de fortaleza (`MinPasswordLength`).
*   **Validación de Entradas:** Uso de Value Objects para `Email` y `HashedPassword` para asegurar la integridad.
*   **Administración de Usuarios:** Endpoints para crear usuarios (con rol), asignar roles, listar y eliminar usuarios (solo para `admin`).
*   **Exclusión de Campos de Auditoría del Dominio:** `CreatedAt`, `UpdatedAt`, `CreatedBy`, `UpdatedBy` se gestionan en la capa de infraestructura/persistencia.

**Justificación:**
Esta decisión proporciona un balance óptimo entre seguridad, rendimiento y flexibilidad para TramaTex. La naturaleza stateless de JWT se alinea con la escalabilidad futura. El RBAC implementado permite un control de acceso robusto y adaptable a las necesidades del negocio. La gestión segura de credenciales y la integración con el ADR-010 (Estrategia de Seguridad) garantizan una base de seguridad sólida.

---

## 5. Consecuencias

### Positivas
-   **Seguridad Mejorada:** Autenticación y autorización robustas, con énfasis en la seguridad de las credenciales.
-   **Rendimiento:** Autenticación rápida sin consultas a base de datos en cada request (para access tokens).
-   **Escalabilidad:** Facilita la distribución de la API sin problemas de sesión.
-   **Flexibilidad:** Adapta bien a clientes frontend y móviles.
-   **Claridad de Roles:** Roles bien definidos y control de acceso predecible.

### Negativas
-   **Gestión de Revocación de Tokens:** Requiere una estrategia específica (ej. blacklist de tokens en Redis) para invalidar tokens antes de su expiración, aumentando la complejidad para MVP (aceptado como riesgo, con solución post-MVP).
-   **Complejidad de Contraseñas:** La política de contraseñas es básica en MVP, se necesita validación de complejidad post-MVP.

---

## 6. Alcance

Esta decisión aplica al diseño y la implementación del módulo de IAM del sistema TramaTex, afectando directamente al backend (Go) y su interacción con la base de datos (PostgreSQL), así como a todos los módulos que requieren autenticación y autorización.

---

## 7. Integración con otros ADRs

-   [ADR-002: Clean Architecture and DDD Adoption](adr-002-clean-architecture-ddd-adoption.md) (Refuerza la pureza del dominio y la separación de capas).
-   [ADR-006: Domain-Driven Development Strategy](adr-006-domain-driven-development-strategy.md) (Alineación con la definición explícita de entidades y Value Objects).
-   [ADR-010: Estrategia de Seguridad: Defensa en Profundidad y Security by Default](adr-010-defense-in-depth-security-strategy.md) (IAM es una capa fundamental en esta estrategia).
-   [ADR-007: Module Implementation Order](adr-007-module-implementation-order.md) (IAM es el primer módulo en ser implementado en Fase 0).

---

## 8. Notas Adicionales / Consideraciones Especiales

*   **Rate Limiting:** Se aplicará rate limiting en el endpoint de login para mitigar ataques de fuerza bruta (implementación post-MVP para mayor rigor).
*   **Auditoría de Seguridad:** La auditoría OWASP Top 10 (ADR-010) fue un motor clave para las decisiones de seguridad del IAM, identificando puntos críticos y su mitigación.
*   **Integración con Party:** Aunque `User` y `Party` son entidades distintas (un usuario tiene acceso, una Party es un tercero de negocio), pueden estar relacionadas (ej. `PartyID` en `User` si un usuario representa a una Party). Esta relación se gestionaría en la capa de aplicación si fuera necesaria.

---

## 9. Referencias

*   Contextos de `architecture.yaml` y `bounded-contexts.yaml`.
*   Documentación del módulo de IAM (`module-spec.md`, `domain-model.md`, `use-cases.md`, `api-contracts.md`).
*   Discusiones y decisiones en este hilo de conversación.
