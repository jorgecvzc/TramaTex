# 🏛️ Visión General de la Arquitectura del Sistema TramaTex

---

Este documento proporciona una visión de alto nivel de la arquitectura de TramaTex. Sirve como punto de entrada principal para entender la estructura técnica, los principios de diseño y cómo los diferentes componentes del sistema interactúan entre sí.

---

## 🔑 Principios Arquitectónicos

*   **Diseño Dirigido por el Dominio (DDD)**: La estructura del software refleja el dominio del negocio mediante Bounded Contexts bien definidos.
*   **Clean Architecture**: Separación estricta de capas (Dominio, Aplicación, Infraestructura, Interfaces) para asegurar la testabilidad y la independencia tecnológica.
*   **Monolito Modular**: Distribución física unificada pero con una separación lógica y funcional rigurosa entre módulos para facilitar una futura extracción a microservicios.
*   **Comunicación Inter-Módulo Síncrona (MVP)**: Para simplificar el desarrollo inicial, los módulos se comunican mediante interfaces de servicio síncronas (ver ADR-019), evitando la complejidad de un bus de eventos distribuido en esta fase.

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
    *   [ADR-018: Arquitectura del Módulo de MES](./adrs/adr-018-mes-module-architecture.md)
    *   [ADR-019: Comunicación Inter-Módulo Síncrona](./adrs/adr-019-synchronous-inter-module-communication-mvp.md)
    *   [ADR-020: Facturas Simplificadas y Series de Facturación](./adrs/adr-020-tickets-and-invoice-series.md)
    *   [ADR-021: Estrategia de Control de Versiones](./adrs/adr-021-version-control-and-branching-strategy.md)
*   **[Diagramas Visuales](./diagrams/README.md)**

---

## 🧩 Estructura de Módulos (Bounded Contexts)

El sistema está dividido en los siguientes contextos delimitados, cada uno con su propio modelo de dominio, lógica de aplicación e infraestructura:

*   **IAM (Identity and Access Management)**: Autenticación, gestión de usuarios, roles y permisos.
*   **Party**: Gestión unificada de terceros (Clientes, Proveedores, Empleados), direcciones y contactos.
*   **Product**: Catálogo de productos tangibles y servicios, gestión de atributos/variantes y marcas.
*   **Pricing**: Motor dinámico de precios que aplica márgenes comerciales y descuentos específicos.
*   **Sales**: Ciclo de vida completo de la venta (Presupuestos, Pedidos, Albaranes, Facturación).
*   **MES (Manufacturing Execution System)**: Control del taller, gestión de plantillas de procesos y seguimiento en tiempo real de la producción.

Para más detalles sobre cada módulo, consulta la documentación específica en `../modules/README.md`.
