# ADR-008 – MVP Timeline Planning (Planificación y Cronograma del MVP)

**Fecha:** 2026-02-01 (Reconstruido)  
**Estado:** Aceptado  
**Autores:** Equipo de Arquitectura de TramaTex  

---

## 1. Contexto

Para gestionar eficazmente el desarrollo del MVP (Minimum Viable Product) y asegurar la entrega de valor de manera incremental, es necesario establecer un cronograma de alto nivel. Este plan organiza la implementación de los módulos, definidos en `ADR-007 (Module Implementation Order)`, en fases lógicas con objetivos y plazos claros.

Esta planificación debe ser lo suficientemente flexible para adaptarse a imprevistos, pero lo suficientemente estructurada para proporcionar una hoja de ruta clara al equipo de desarrollo y a los stakeholders.

---

## 2. Alternativas Consideradas

**Alternativa A – Planificación Basada en Funcionalidades (Feature-driven)**  
- Ventajas: Visibilidad directa de la funcionalidad terminada.  
- Desventajas: Riesgo de bloqueos por dependencias técnicas profundas no resueltas.

**Alternativa B – Planificación Basada en Fases Lógicas (Adoptada)**  
- Ventajas: Respeta el orden de dependencias de los módulos (ADR-007) y garantiza bases sólidas.  
- Desventajas: Valor para el usuario final más visible en fases tardías (Fase 2 y 3).

---

## 3. Criterios de Decisión

- **Respeto a Dependencias:** Alineación con el orden de implementación de módulos.
- **Minimización de Riesgos:** Fundaciones técnicas completadas antes de lógica de negocio compleja.
- **Claridad de Objetivos:** Hitos definidos para cada fase.

---

## 4. Decisión Adoptada

Se establece un plan de desarrollo por fases, alineado con el orden de implementación de módulos (`ADR-007`). Cada fase agrupa un conjunto de módulos relacionados que constituyen una entrega de valor coherente.

### Fases del Proyecto

**Fase 0: Fundaciones Técnicas**
- **Estado:** ✅ Completada
- **Plazo de Finalización:** Q1 2026
- **Módulos Incluidos:** `IAM`.

**Fase 1: Dominio Base y Lógica Crítica de Negocio**
- **Estado:** 🔄 En Progreso
- **Plazo Objetivo:** Q2 - Q4 2026
- **Módulos Incluidos:** `Party`, `Product`, `Pricing`.

**Fase 2: Flujo de Negocio Principal**
- **Estado:** ⏳ Planificada
- **Módulos Incluidos:** `Sales`.

**Fase 3: Operaciones de Taller (MVP)**
- **Estado:** ⏳ Planificada
- **Módulos Incluidos:** `MES`.

---

## 5. Consecuencias

### Positivas
- Proporciona una visión clara del progreso del proyecto y los próximos pasos.
- Permite a los stakeholders entender cuándo se espera que las funcionalidades estén disponibles.
- Facilita la planificación de recursos y la gestión de expectativas.

### Negativas
- Las fechas son estimaciones de alto nivel y requieren comunicación constante.
- Puede reducir la agilidad si se interpreta de manera demasiado rígida.

---

## 6. Alcance

Este cronograma cubre desde la Fase 0 (Fundaciones Técnicas) hasta la Fase 3 (Operaciones de Taller / MES) para el alcance del MVP.

---

## 7. Integración con otros ADRs

- **ADR-007 (Module Implementation Order):** Este cronograma es la aplicación práctica del orden de dependencias.
- **ADR-004 (MVP Development Lifecycle):** Define el enfoque iterativo dentro de estas fases.

---

## 8. Notas Adicionales / Consideraciones Especiales

*Nota: Este documento ha sido reconstruido el 2026-02-01 a partir de las referencias encontradas en los agentes de contexto y otros documentos del proyecto.*

---

## 9. Referencias

- Registro de Sprints (`agents/project/sprint-registry.yaml`)
- Definición de Contextos (`agents/project/context/bounded-contexts.yaml`)
