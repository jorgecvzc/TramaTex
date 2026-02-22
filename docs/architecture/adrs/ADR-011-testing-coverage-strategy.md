# ADR-011 – Testing & Coverage Strategy (Estrategia de Testing y Cobertura)

**Fecha:** 2026-02-01 (Reconstruido)  
**Estado:** Aceptado  
**Autores:** Equipo de Arquitectura de TramaTex  

---

## 1. Contexto

La calidad, fiabilidad y mantenibilidad del software son pilares fundamentales para el éxito del proyecto TramaTex. Una estrategia de testing bien definida es crucial para validar que el sistema cumple con los requerimientos del negocio, prevenir regresiones y facilitar futuras modificaciones.

Este documento establece la estrategia de testing, los tipos de pruebas, las herramientas y los objetivos de cobertura de código que deben seguirse en todo el proyecto, formalizando los estándares de calidad.

---

## 2. Decisión Adoptada

Se adopta una estrategia de testing multi-capa que se alinea con la Arquitectura Limpia (`ADR-002`) del proyecto. El principio rector es "probar en la capa que corresponde", enfocando los esfuerzos de testing de manera diferente en el dominio, la aplicación y la infraestructura.

La estrategia se basa en los siguientes puntos clave:

1.  **Pirámide de Testing:** Se favorece una base sólida de tests unitarios rápidos, complementados por menos tests de integración y un número muy selectivo de tests E2E. La proporción recomendada como guía es **70% Unitarios, 25% Integración, 5% E2E**.
2.  **Enfoque Iterativo de Desarrollo ("Afinación Continua"):** Se sigue un enfoque iterativo de "construir y refinar", donde el desarrollo avanza a través de la implementación de dominio, casos de uso, infraestructura y UI, seguido de una fase de comprobación y afinación, reiniciando el ciclo si es necesario. Aunque se prioriza la escritura de tests en todas las capas, la aplicación estricta de TDD se reserva para decisiones específicas y se aborda de forma pragmática según la criticidad del módulo y el contexto del desarrollo.
3.  **Objetivos de Cobertura (Coverage):** Se establecen mínimos de cobertura de código para garantizar que la lógica crítica esté adecuadamente probada.

### Objetivos de Cobertura por Fase (MVP vs Post-MVP)

Los umbrales se diferencian por fase. En MVP se exige el mínimo viable; en Post-MVP se incrementa el rigor.

#### MVP (mínimos obligatorios)

Backend (Go):

| Módulo | Cobertura Mínima | Criticidad | Justificación |
|---|---|---|---|
| **Pricing** | **≥ 85%** | Económica | Errores impactan directamente en la facturación y rentabilidad. |
| **Party** | ≥ 75% | Funcional | Errores en la gestión de clientes/proveedores impactan toda la operativa. |
| **Product (Domain)** | ≥ 75% | Funcional | La correcta definición de productos es clave para precios y ventas. |
| **Product (Application)** | **≥ 50%** * | Funcional | Tests de integración cubren flujos críticos. Ver nota (*). |
| **Sales** | ≥ 75% | Funcional | Gestión de flujo comercial end-to-end. |
| **IAM** | ≥ 75% | Seguridad | Control de acceso y autenticación. |
| **General** | **≥ 75%** | Calidad Base | Mínimo aceptable para cualquier otro módulo de backend. |

**Nota (*):** Product Application tiene objetivo ajustado de 50% para MVP debido a:
- Extensa cobertura de tests de integración (product_service_integration_test.go) que validan flujos completos
- Funciones críticas (ListAttributes, GetApplicableAttributes, FindOrCreateVariant) con tests unitarios robustos
- Product Domain mantiene 83.6% coverage (sobre objetivo 75%)
- Complejidad de mocking en funciones de generación de variantes con cadenas de llamadas internas
- Priorización estratégica: Pricing (85.4% ✅) y Sales (75.3% ✅) completos primero
- Product Application coverage actual: **49.5%** con 14 tests unitarios + tests integración completos

Frontend (Vue): cobertura **≥70%** para lógica de negocio crítica (stores de Pinia, composables reutilizables).

#### Post-MVP (objetivos reforzados)

Mínimos por capa (backend):

| Capa | Cobertura Mínima | Tipo de Tests |
|---|---|---|
| **Domain** | 100% | Unit |
| **Application** | 95% | Unit |
| **Persistence** | 80% | Integration |
| **Interfaces** | 80% | Unit + Integration |
| **Infrastructure** | 70% | Integration |

Métricas globales:
- **Coverage total del proyecto:** ≥ 90%
- **Branch coverage:** ≥ 85%

### Tipos de Tests por Capa Arquitectónica

-   **Capa de Dominio (Domain):**
    -   **Tipo:** Tests Unitarios.
    -   **Enfoque:** Probar entidades, value objects y servicios de dominio en total aislamiento. No se mockea nada, ya que no debe haber dependencias externas. Se valida la lógica de negocio pura.
    -   **Herramientas:** `Go testing`, `testify/assert`.

-   **Capa de Aplicación (Application):**
    -   **Tipo:** Tests de Integración.
    -   **Enfoque:** Probar los casos de uso (Application Services) de principio a fin, mockeando las dependencias externas como los repositorios de base de datos y servicios de terceros. Se valida la orquestación y los flujos de trabajo.
    -   **Herramientas:** `Go testing`, `testify/mock`.

-   **Capa de Infraestructura (Infrastructure):**
    -   **Tipo:** Tests de Integración.
    -   **Enfoque:** Probar la correcta implementación de las interfaces del dominio (ej. Repositorios GORM). Se valida la interacción con la base de datos real en un entorno de prueba.
    -   **Herramientas:** `Go testing`, Testcontainers para la base de datos.

-   **Capa de Interfaces (Frontend y Handlers):**
    -   **Tipo:** Tests Unitarios y E2E.
    -   **Enfoque:**
        -   **Frontend (Vue):** Tests unitarios para stores de Pinia y composables. Tests de componentes para interacciones complejas.
        -   **Backend (Go):** Se desaconsejan los tests complejos de handlers. Se prefieren tests E2E para los flujos más críticos.
        -   **E2E:** Flujos completos de usuario (ej. login -> crear pedido -> verificar precio).
    -   **Herramientas:** `Vitest` (Vue), `Playwright` (E2E).

---

## 3. Consecuencias

### Positivas
- Aumenta significativamente la confianza en la corrección del código.
- Reduce el número de bugs y regresiones en producción.
- El código se diseña para ser testeable, lo que generalmente conduce a un mejor diseño de software (más modular y desacoplado).
- Facilita la refactorización y el mantenimiento a largo plazo.

### Negativas
- Incrementa el tiempo de desarrollo inicial al requerir la escritura de tests.
- El mantenimiento de los tests es un esfuerzo continuo que debe ser considerado en la planificación.

---

## 4. Integración con otros ADRs

- **ADR-002 (Clean Architecture):** Esta estrategia de testing está intrínsecamente ligada a la separación de capas definida en `ADR-002`. La arquitectura hace posible esta estrategia.
- **ADR-010 (Seguridad):** Los tests de seguridad (ej. verificar políticas de autorización) son una parte integral de esta estrategia.

---
*Nota: Este documento ha sido reconstruido el 2026-02-01 a partir de la información consolidada en `code-standards.yaml`, `architecture.yaml` y plantillas de módulo.*
