# ADR-009 – Estructura de Carpetas y Organización del Proyecto

**Fecha:** 11/01/2026  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, Claude (Anthropic)

---

## 1. Contexto

TramaTex requiere una estructura de proyecto que:

- Refleje los principios de **Clean Architecture** (ADR-002)
- Soporte un **monolito modular** con dominios claramente separados (ADR-003)
- Facilite el **desarrollo dirigido por dominio** (ADR-006, ADR-007)
- Sea **mantenible por un desarrollador en solitario** con IA copiloto
- Permita **escalabilidad futura** hacia microservicios si es necesario
- Organice **documentación técnica y de sesiones** de forma eficiente

**Restricciones técnicas:**

- Stack: Go (tramatex-api), Vue.js 3 (frontend), PostgreSQL (persistencia), Docker Compose
- Desarrollo local-first
- Un solo repositorio Git (monorepo)
- Sin herramientas complejas de monorepo (Nx, Turborepo, etc.)

**Riesgos si no se define estructura:**

- Mezcla de capas (dominio con infraestructura)
- Dificultad para aplicar un desarrollo guiado por pruebas estricto.
- Pérdida de trazabilidad entre módulos
- Código difícil de navegar y mantener
- Imposibilidad de extraer módulos a servicios independientes en el futuro

---

## 2. Alternativas Consideradas

### Alternativa A – Estructura plana por tipo de archivo
- Ventajas: Simple y tradicional.
- Desventajas: Mezcla todos los dominios, no escala.

### Alternativa B – Estructura por capas técnicas
- Ventajas: Respeta capas de Clean Architecture.
- Desventajas: No refleja dominios de negocio.

### Alternativa C – Estructura por módulos de dominio (Clean Architecture modular)
- Ventajas: Refleja dominios claramente, facilita extracción de módulos.
- Desventajas: Mayor profundidad de carpetas.

### Alternativa D – Monorepo con workspaces independientes
- Ventajas: Preparado para microservicios desde día 1.
- Desventajas: Complejidad prematura para MVP.

---

## 3. Criterios de Decisión

- Alineación con Clean Architecture y DDD.
- Soporte para un monolito modular.
- Facilidad de mantenimiento por un solo desarrollador.
- Preparación para futura extracción de microservicios.
- Trazabilidad entre documentación y código.

---

## 4. Decisión Adoptada

Se adopta **Alternativa C: Estructura por módulos de dominio con Clean Architecture**.

**Justificación:**
- **Refleja el dominio de negocio:** Party, Producto, Tarificación, Ventas, MES son módulos visibles.
- **Soporta Clean Architecture:** Capas separadas dentro de cada módulo.
- **Facilita el desarrollo guiado por pruebas:** Tests están junto al código que prueban.
- **Escalabilidad controlada:** Módulos pueden extraerse a servicios sin refactoring masivo.

---

## 5. Consecuencias

### Positivas
- Claridad de dominio y separación de responsabilidades.
- Testabilidad mejorada.
- Onboarding rápido y preparación para microservicios.

### Negativas
- Mayor profundidad de carpetas (rutas largas).
- Disciplina requerida para evitar romper la separación.

---

## 6. Alcance

Este ADR define la estructura completa de carpetas del proyecto TramaTex, incluyendo la organización del monorepo (api, frontend, docs) y las convenciones de nombres.

---

## 7. Integración con otros ADRs

- **ADR-001:** Stack tecnológico.
- **ADR-002:** Clean Architecture.
- **ADR-003:** Monolito modular.

---

## 8. Notas Adicionales / Consideraciones Especiales

### Estructura Completa del Proyecto
Para una vista detallada de la estructura de carpetas del proyecto, consulte [Guía de Detalles de la Estructura del Proyecto](../../guides/developer/project-structure-details.md).

---

## 9. Referencias

- Clean Architecture (Robert C. Martin)
- Domain-Driven Design (Eric Evans)
- Golang Project Layout: https://github.com/golang-standards/project-layout
