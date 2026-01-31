# Guía para Desarrolladores del Backend (tramatex-api)

Este documento sirve como punto de entrada para los desarrolladores que trabajan en el backend de TramaTex (`tramatex-api`).

## 1. Arquitectura y Estructura

La arquitectura del proyecto es la base de todo el desarrollo. Antes de escribir código, es fundamental comprender las decisiones y la estructura establecidas.

-   **Fuente de Verdad Arquitectónica:** [ADR-009 – Estructura de Carpetas y Organización del Proyecto](../architecture/adrs/ADR-009-estructura-proyecto.md)
    -   Este documento contiene la estructura completa y actualizada de todo el proyecto, incluyendo el backend. Describe cómo se organizan los módulos (Bounded Contexts) y las capas de la Clean Architecture.

-   **Principios Clave:**
    -   **Clean Architecture:** El código se organiza en capas concéntricas (Dominio, Aplicación, Infraestructura, Interfaces). El dominio es el núcleo y no depende de nada externo.
    -   **Domain-Driven Design (DDD):** La estructura del código refleja los dominios del negocio (IAM, Party, Product, etc.).
    -   **Monolito Modular:** Todo el backend reside en una única aplicación (`tramatex-api`), pero está lógicamente separado en módulos que podrían extraerse en el futuro.

## 2. Documentación de Módulos (Bounded Contexts)

Cada módulo de negocio tiene su propia documentación detallada. Para entender la implementación de un módulo específico, consulta su directorio correspondiente:

-   [**Módulo IAM**](../../reference/iam/)
-   [**Módulo Party**](../../reference/party/)
-   [**Módulo Product**](../../reference/product/)
-   [**Módulo Pricing**](../../reference/pricing/)
-   [**Módulo Sales**](../../reference/sales/)

Dentro de cada uno de estos directorios, encontrarás:
-   `module-spec.md`: La especificación funcional del módulo.
-   `domain-model.md`: Detalles sobre la implementación del modelo de dominio.
-   `api-contracts.md`: La especificación de la API REST para ese módulo.

## 3. Estrategia de Testing

El proyecto sigue un enfoque de **Test-Driven Development (TDD)**, especialmente para la lógica de dominio.

-   **Ubicación de los Tests:** Los tests se encuentran junto al código que prueban, dentro de cada módulo. Por ejemplo, los tests para `internal/iam/domain/model/user.go` están en `internal/iam/domain/model/user_test.go`.
-   **Ejecución de Tests:** Utiliza los comandos definidos en el `Makefile` en la raíz del proyecto para ejecutar la suite de tests.
    -   `make test`: Ejecuta todos los tests del backend.
    -   `make test-coverage`: Ejecuta los tests y muestra un informe de cobertura.

---
**Nota:** Esta guía reemplaza los documentos `structure.md`, `testing.md`, y `validation-report.md` que quedaron obsoletos tras la refactorización del 18 de enero de 2026.
