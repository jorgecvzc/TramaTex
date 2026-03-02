# Módulo de Sales (Gestión de Órdenes)

Este módulo es el núcleo del proceso de ventas en TramaTex. Gestiona el ciclo de vida completo de la venta, desde la creación de cotizaciones hasta la emisión de facturas, pasando por los pedidos y albaranes. Se integra estrechamente con los módulos de Party, Product y Pricing para consolidar toda la información relevante.

## Diseño Arquitectónico

Para una descripción detallada de las decisiones arquitectónicas del módulo de Sales, consulte el siguiente Architectural Decision Record (ADR):

*   [ADR-017: Arquitectura del Módulo de Sales](../../architecture/adrs/adr-017-sales-module-architecture.md)

## Componentes Clave

*   **Entidades de Dominio:**
    *   `Quote` y `QuoteLineItem`
    *   `SalesOrder` y `OrderLineItem`
    *   `DeliveryNote` y `DeliveryNoteLineItem`
    *   `Invoice` y `InvoiceLineItem`
    *   Value Objects como `Money`, `Percentage`, `OrderNumber`, `QuoteNumber`, `DeliveryNoteNumber`, `InvoiceNumber`.

*   **Casos de Uso (Capa de Aplicación):**
    *   Un conjunto completo de casos de uso para la gestión de cada tipo de documento de venta y sus transiciones de estado.

## Documentación Detallada

Consulte los siguientes documentos para una descripción más profunda del módulo de Sales:

*   **Especificación del Módulo:** [module-spec.md](./module-spec.md)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md)
*   **Casos de Uso:** [use-cases.md](./use-cases.md)
*   **Contratos de API:** [api-contracts.md](./api-contracts.md)
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md)
