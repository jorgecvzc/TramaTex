# Arquitectura General del Sistema

Este documento proporciona una visión de alto nivel de la arquitectura de TramaTex. Sirve como punto de entrada principal para entender la estructura técnica, los principios de diseño y cómo los diferentes componentes del sistema interactúan entre sí.

## 1. Principios Arquitectónicos

*   **Diseño Dirigido por el Dominio (DDD)**: La estructura del software refleja el dominio del negocio.
*   **Clean Architecture**: Separación estricta de capas (Dominio, Aplicación, Infraestructura, Interfaces).
*   **Comunicación por Eventos**: Para la comunicación entre Bounded Contexts desacoplados.

## 2. Documentos de Referencia

Este overview es el punto de partida. Para profundizar, consulta los siguientes recursos:

*   **Contextos Delimitados (Bounded Contexts)**: [./diagrams/C1-context.md](./diagrams/C1-context.md)
*   **Glosario Ubicuo**: [./glossary.md](./glossary.md)
*   **Decisiones Arquitectónicas (ADRs)**: [./adrs/](./adrs/). Incluye decisiones clave como:
    *   [ADR-001](./adrs/ADR-001-technology-stack-selection.md): Stack Tecnológico.
    *   [ADR-002](./adrs/ADR-002-clean-architecture-ddd-adoption.md): Clean Architecture y DDD.
    *   [ADR-003](./adrs/ADR-003-application-distribution-type.md): Monolito Modular.
    *   [ADR-004](./adrs/ADR-004-mvp-development-lifecycle.md): Ciclo de Vida del MVP.
    *   [ADR-005](./adrs/ADR-005-unified-customer-supplier-management.md): Patrón Party.
    *   [ADR-006](./adrs/ADR-006-domain-driven-development-strategy.md): Desarrollo Dirigido por Dominio.
    *   [ADR-007](./adrs/ADR-007-module-implementation-order.md): Orden de Implementación.
    *   [ADR-008](./adrs/ADR-008-mvp-timeline-planning.md): Planificación y Cronograma.
    *   [ADR-009](./adrs/ADR-009-project-structure.md): Estructura de Carpetas.
    *   [ADR-010](./adrs/ADR-010-defense-in-depth-security-strategy.md): Estrategia de Seguridad.
    *   [ADR-011](./adrs/ADR-011-testing-coverage-strategy.md): Estrategia de Testing.
    *   [ADR-012](./adrs/ADR-012-arquitectura-modulo-party.md): Arquitectura del Módulo de Party.
    *   [ADR-013](./adrs/ADR-013-manejo-de-modificaciones-de-producto.md): Manejo de Modificaciones de Producto.
    *   [ADR-014](./adrs/ADR-014-iam-module-architecture.md): Arquitectura del Módulo de IAM.
    *   [ADR-015](./adrs/ADR-015-product-module-architecture.md): Arquitectura del Módulo de Product.
    *   [ADR-016](./adrs/ADR-016-pricing-module-architecture.md): Arquitectura del Módulo de Pricing.
    *   [ADR-017](./adrs/ADR-017-sales-module-architecture.md): Arquitectura del Módulo de Sales.
*   **Diagramas Visuales**: [./diagrams/](./diagrams/)

## 3. Estructura de Módulos

El sistema está dividido en los siguientes módulos principales (Bounded Contexts):

*   **IAM (Identity and Access Management)**: Gestiona usuarios, roles y permisos.
*   **Party**: Gestiona la información de clientes y proveedores.
*   **Product**: Catálogo de productos y sus propiedades.
*   **Sales**: Gestiona el proceso de ventas.
*   **Pricing**: Calcula los precios de los productos.

Para más detalles sobre cada módulo, consulta la documentación específica en `../modules/`.
