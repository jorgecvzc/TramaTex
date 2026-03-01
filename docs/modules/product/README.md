# Módulo de Product (Catálogo de Productos)

Este módulo es fundamental para TramaTex, ya que gestiona todo el catálogo de productos y servicios de la empresa. Proporciona la infraestructura para definir productos con atributos configurables, manejar variantes (incluida su creación Just-in-Time) y clasificar el catálogo de manera flexible.

## Diseño Arquitectónico

Para una descripción detallada de las decisiones arquitectónicas, entidades de dominio, objetos de valor, casos de uso y estrategia de gestión de variantes de este módulo, consulte el siguiente Architectural Decision Record (ADR):

*   [ADR-015: Arquitectura del Módulo de Product](../../architecture/adrs/adr-015-product-module-architecture.md)

## Componentes Clave

*   **Entidades de Dominio:**
    *   `Attribute` y `AttributeValue`: Para definir características configurables de productos.
    *   `Product`: La plantilla del producto/servicio.
    *   `ProductVariant`: La instancia final y vendible del producto, compuesta por `AttributeValue`s.
    *   `Brand` y `ProductGroup`: Para clasificar y agrupar productos.
    *   `PartyServiceConfiguration`: Para configuraciones de servicios específicos por cliente.

*   **Casos de Uso (Capa de Aplicación):**
    *   Gestión completa de `Attribute`s, `Product`s y `ProductVariant`s, incluyendo la creación "Just-in-Time" de variantes.

## Documentación Detallada

Consulte los siguientes documentos para una descripción más profunda del módulo de Product:

*   **Especificación del Módulo:** [module-spec.md](./module-spec.md)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md)
*   **Casos de Uso:** [use-cases.md](./use-cases.md)
*   **Contratos de API:** [api-contracts.md](./api-contracts.md)
*   **Diagramas Detallados del Dominio:** [diagrams/detailed-domain-models.md](./diagrams/detailed-domain-models.md)
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md)
