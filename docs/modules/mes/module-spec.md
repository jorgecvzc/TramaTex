# Especificación del Módulo - MES

## Objetivo de esta revisión

Establecer una nomenclatura clara para diferenciar:

1. El trabajo definido por cliente (plantilla reusable).
2. La ejecución real de producción derivada de esa plantilla.

Esta actualización **no crea sprint ni ADR nuevos**. Se documenta para guiar el refactor progresivo de MES.

---

## Nomenclatura canónica propuesta

### 1) Trabajo Definido (Plantilla)

- **Nombre funcional (UI):** `Trabajo Definido`
- **Nombre técnico recomendado (backend/frontend):** `MESWorkDefinition`
- **Plural API recomendado:** `work-definitions`
- **Propósito:** Define *qué* hay que hacer para un cliente/tipo de prenda.
- **Incluye:**
	- cliente (`party_id`)
	- grupo tangible / familia de producto
	- secuencia de grupos de servicio
	- secuencia de tareas por grupo
	- parámetros por defecto (notas, diseño, etc.)

### 2) Trabajo Real (Ejecución)

- **Nombre funcional (UI):** `Trabajo Real`
- **Nombre técnico recomendado (backend/frontend):** `MESWorkExecution`
- **Plural API recomendado:** `work-executions`
- **Propósito:** Representa una instancia operativa de producción.
- **Incluye:**
	- referencia al `work_definition_id`
	- fechas reales (`start_date`, `due_date`, `completed_date`)
	- prendas/lotes concretos a producir
	- estado global de ejecución
	- estado de tareas ejecutables (START/PAUSE/COMPLETE/BLOCK)

---

## Regla de negocio principal

- De **un Trabajo Definido** se pueden crear **múltiples Trabajos Reales**.
- Un **Trabajo Real** siempre debe estar vinculado a un único **Trabajo Definido**.

---

## Problema actual detectado

El modelo actual usa `MESWork` / `works` para todo y mezcla en una misma entidad:

- información de definición (estructura de servicio/tareas), y
- información de ejecución (estado/fechas).

Esto dificulta la trazabilidad, la reutilización y la semántica de UI/API.

---

## Mapeo de transición (actual → propuesto)

| Actual | Propuesto | Observación |
|---|---|---|
| `MESWork` | `MESWorkExecution` | Queda para la instancia real |
| `work_name` | `execution_name` (opcional) | Nombre operativo de la ejecución |
| `service_group_assignments` en creación de work | En `MESWorkDefinition` | La estructura base debe vivir en la plantilla |
| `works` endpoint | `work-executions` | Mantener alias temporal para compatibilidad |

---

## Convención de términos en UI

- Menú y pantallas:
	- `Trabajos Definidos`
	- `Trabajos Reales`
- Evitar etiquetas ambiguas como solo `Trabajos` en pantallas donde coexistan ambos conceptos.

---

## Alcance de implementación sugerido (sin romper de golpe)

1. Introducir términos en documentación y textos UI.
2. Mantener endpoints actuales como alias de compatibilidad.
3. Incorporar nuevos endpoints/versionado para definición vs ejecución.
4. Migrar frontend gradualmente a la nueva semántica.

---

**Última Actualización:** 2026-02-23
