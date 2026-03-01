# 🏛️ Visión General de la Arquitectura del Sistema TramaTex

---

Este documento proporciona una visión de alto nivel de la arquitectura de TramaTex. Sirve como punto de entrada principal para entender la estructura técnica, los principios de diseño y cómo los diferentes componentes del sistema interactúan entre sí.

---

## 🔑 Principios Arquitectónicos

*   **Diseño Dirigido por el Dominio (DDD)**: La estructura del software refleja el dominio del negocio.
*   **Clean Architecture**: Separación estricta de capas (Dominio, Aplicación, Infraestructura, Interfaces).
*   **Comunicación por Eventos**: Para la comunicación entre Bounded Contexts desacoplados.

---

## 📚 Documentos de Referencia

Este overview es el punto de partida. Para profundizar, consulta los siguientes recursos:

*   **[Contextos Delimitados (Bounded Contexts)](./diagrams/c1-context.md)**
*   **[Glosario Ubicuo](./glossary.md)**
*   **[Decisiones Arquitectónicas (ADRs)](./adrs/README.md)**
    *   [ADR-001: Stack Tecnológico](./adrs/adr-001-technology-stack-selection.md)
    *   [ADR-002: Clean Architecture y DDD](./adrs/adr-002-clean-architecture-ddd-adoption.md)
    *   [ADR-003: Monolito Modular](./adrs/adr-003-application-distribution-type.md)
    *   [ADR-004: Ciclo de Vida del MVP](./adrs/adr-004-mvp-development-lifecycle.md)
    *   [ADR-005: Patrón Party](./adrs/adr-005-unified-customer-supplier-management.md)
    *   [ADR-006: Desarrollo Dirigido por Dominio](./adrs/adr-006-domain-driven-development-strategy.md)
    *   [ADR-007: Orden de Implementación](./adrs/adr-007-module-implementation-order.md)
    *   [ADR-008: Planificación y Cronograma](./adrs/adr-008-mvp-timeline-planning.md)
    *   [ADR-009: Estructura de Carpetas](./adrs/adr-009-project-structure.md)
    *   [ADR-010: Estrategia de Seguridad](./adrs/adr-010-defense-in-depth-security-strategy.md)
    *   [ADR-011: Estrategia de Testing](./adrs/adr-011-testing-coverage-strategy.md)
    *   [ADR-012: Arquitectura del Módulo de Party](./adrs/adr-012-party-module-architecture.md)
    *   [ADR-013: Manejo de Modificaciones de Producto](./adrs/adr-013-product-modifications-handling.md)
    *   [ADR-014: Arquitectura del Módulo de IAM](./adrs/adr-014-iam-module-architecture.md)
    *   [ADR-015: Arquitectura del Módulo de Product](./adrs/adr-015-product-module-architecture.md)
    *   [ADR-016: Arquitectura del Módulo de Pricing](./adrs/adr-016-pricing-module-architecture.md)
    *   [ADR-017: Arquitectura del Módulo de Sales](./adrs/adr-017-sales-module-architecture.md)
*   **[Diagramas Visuales](./diagrams/README.md)**

---

## 🧩 Estructura de Módulos

El sistema está dividido en los siguientes módulos principales (Bounded Contexts):

*   **IAM (Identity and Access Management)**: Gestiona usuarios, roles y permisos.
*   **Party**: Gestiona la información de clientes y proveedores.
*   **Product**: Catálogo de productos y sus propiedades.
*   **Sales**: Gestiona el proceso de ventas.
*   **Pricing**: Calcula los precios de los productos.

Para más detalles sobre cada módulo, consulta la documentación específica en `../modules/README.md`.
