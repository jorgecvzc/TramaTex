# Arquitectura General del Sistema

Este documento proporciona una visión de alto nivel de la arquitectura de TramaTex. Sirve como punto de entrada principal para entender la estructura técnica, los principios de diseño y cómo los diferentes componentes del sistema interactúan entre sí.

## 1. Principios Arquitectónicos

*   **Diseño Dirigido por el Dominio (DDD)**: La estructura del software refleja el dominio del negocio.
*   **Clean Architecture**: Separación estricta de capas (Dominio, Aplicación, Infraestructura, Interfaces).
*   **Comunicación por Eventos**: Para la comunicación entre Bounded Contexts desacoplados.

## 2. Documentos de Referencia

Este overview es el punto de partida. Para profundizar, consulta los siguientes recursos:

*   **Contextos Delimitados (Bounded Contexts)**: [./bounded-contexts.md](./bounded-contexts.md)
*   **Glosario Ubicuo**: [./glossary.md](./glossary.md)
*   **Decisiones Arquitectónicas (ADRs)**: [./adr/](./adr/). Incluye decisiones clave como:
    *   [ADR-001](./adr/ADR-001-seleccion-stack-tecnologico.md): Stack Tecnológico.
    *   [ADR-002](./adr/ADR-002-adopcion-clean-architecture-ddd.md): Clean Architecture y DDD.
    *   [ADR-003](./adr/ADR-003-tipo-distribucion-aplicacion.md): Monolito Modular.
    *   [ADR-007](./adr/ADR-007-orden-implementacion-modulos.md): Orden de Implementación.
    *   [ADR-010](./adr/ADR-010-estrategia-seguridad-defensa-profundidad.md): Estrategia de Seguridad.
*   **Diagramas Visuales**: [./diagrams/](./diagrams/)

## 3. Estructura de Módulos

El sistema está dividido en los siguientes módulos principales (Bounded Contexts):

*   **IAM (Identity and Access Management)**: Gestiona usuarios, roles y permisos.
*   **Party**: Gestiona la información de clientes y proveedores.
*   **Product**: Catálogo de productos y sus propiedades.
*   **Sales**: Gestiona el proceso de ventas.
*   **Pricing**: Calcula los precios de los productos.

Para más detalles sobre cada módulo, consulta la documentación específica en `../3_modules/`.
