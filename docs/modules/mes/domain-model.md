# Modelo de Dominio - Módulo MES

## Entidades principales (propuesta objetivo)

### `MESWorkDefinition` (Trabajo Definido)

Representa la plantilla operativa por cliente/producto.

**Campos clave (conceptuales):**
- `id`
- `definition_code`
- `definition_name`
- `party_id`
- `tangible_group_id`
- `default_notes`
- `service_group_blueprint[]`

### `MESWorkExecution` (Trabajo Real)

Representa una ejecución real de producción derivada de una plantilla.

**Campos clave (conceptuales):**
- `id`
- `execution_number`
- `execution_name`
- `work_definition_id`
- `status`
- `priority`
- `start_date`
- `due_date`
- `completed_date`
- `garment_items[]`
- `task_instances[]`

---

## Relación entre entidades

- `MESWorkDefinition (1) -> (N) MESWorkExecution`

Una definición se reutiliza para múltiples ejecuciones. Las ejecuciones no deben redefinir la estructura base completa salvo overrides controlados.

---

## Estados recomendados

### Estado de ejecución (`MESWorkExecution.status`)
- `DRAFT`
- `PENDING`
- `IN_PROGRESS`
- `ON_HOLD`
- `COMPLETED`
- `CANCELLED`

### Estado de tareas ejecutables
- `PENDING`
- `IN_PROGRESS`
- `PAUSED`
- `COMPLETED`
- `BLOCKED`

---

## Terminología obligatoria en código y UI

- `Definition` = Trabajo Definido
- `Execution` = Trabajo Real

Evitar el uso ambiguo de `Work` sin sufijo en nuevas piezas de código/documentación.

---

## Estado actual vs objetivo

Actualmente `MESWork` contiene atributos mixtos de definición + ejecución.

Objetivo: separar responsabilidades en dos agregados para mejorar:
- claridad del modelo,
- reutilización de plantillas,
- trazabilidad de producción real,
- mantenibilidad de API/UI.

---

**Última Actualización:** 2026-02-23
