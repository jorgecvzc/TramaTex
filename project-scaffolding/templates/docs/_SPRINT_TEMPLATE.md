# PLANTILLA DE SPRINT

---

## 📋 INFORMACIÓN DEL SPRINT

| Campo | Valor |
|-------|-------|
| **ID de Sprint** | XX |
| **Título/Objetivo** | [Objetivo principal del Sprint] |
| **Estado** | ⏳ Planificado / 🔄 En Progreso / ✅ Completado |
| **Fecha de Inicio** | YYYY-MM-DD |
| **Fecha de Fin** | YYYY-MM-DD |
| **Participantes** | [Nombres del equipo], [LLM usado] |

---

## 🎯 OBJETIVO DEL SPRINT (SPRINT GOAL)

*Describe en 1-3 frases el objetivo de negocio que se persigue con este sprint. ¿Qué valor se entregará al final?*

---

## 📝 TAREAS PLANIFICADAS

*Lista de todas las tareas que se planean abordar en este sprint. El estado debe actualizarse a medida que avanza el sprint.*

| ID Tarea | Título de la Tarea | Estado | Asignado a | Estimación (h) |
|----------|--------------------|--------|------------|----------------|
| [01] | [Título de la Tarea 1] | ⏳ | [Nombre] | X |
| [02] | [Título de la Tarea 2] | ⏳ | [Nombre] | Y |
| [03] | [Título de la Tarea 3] | ⏳ | [Nombre] | Z |

**Nota sobre numeración:**
- Los IDs de tareas son locales al sprint (01, 02, 03... reinician en cada sprint)
- Identificación única: sprint-XX + tarea-YY (ej: 01-01, 02-03)
- Esto permite reabrir sprints sin conflictos de numeración

**Enlace a la plantilla de tarea:** [_TASK_TEMPLATE.md](_TASK_TEMPLATE.md)

---

## ✅ TAREAS COMPLETADAS

*Lista de tareas que se completaron con éxito durante el sprint.*

| ID Tarea | Título de la Tarea | Resultado / Enlace |
|----------|--------------------|---------------------------------|
| [01] | [Título de la Tarea 1] | [01-nombre-tarea-1.md](./01-nombre-tarea-1.md) |
| [02] | [Título de la Tarea 2] | [02-nombre-tarea-2.md](./02-nombre-tarea-2.md) |

---

## 🔍 REVISIÓN DEL SPRINT (SPRINT REVIEW)

### Demostración
*¿Qué se ha enseñado a los stakeholders? ¿Cuál fue su feedback?*

- **Funcionalidad A:** [Feedback recibido]
- **Funcionalidad B:** [Feedback recibido]

### Métricas del Sprint
| Métrica | Valor |
|---------|-------|
| **Tareas Planificadas** | X |
| **Tareas Completadas** | Y |
| **Porcentaje de Éxito**| (Y/X) % |
| **Horas Estimadas** | H |
| **Horas Reales** | R |
| **Desviación** | (R-H) horas |

---

## 🔄 RETROSPECTIVA DEL SPRINT (SPRINT RETROSPECTIVE)

### ¿Qué fue bien?
*Cosas que el equipo hizo bien y que se deben mantener.*
- [Punto positivo 1]
- [Punto positivo 2]

### ¿Qué se puede mejorar?
*Problemas o ineficiencias que el equipo encontró.*
- [Punto a mejorar 1]
- [Punto a mejorar 2]

### Acciones de Mejora
*Acciones concretas para implementar en el próximo sprint.*
1. **Acción 1:** [Descripción de la acción] (Responsable: [Responsable])
2. **Acción 2:** [Descripción de la acción] (Responsable: [Responsable])

## 🤖 Para Asistentes de IA
Una vez completado el sprint, actualiza `sprint-registry.yaml` moviendo el sprint a la sección `completed_sprints` y añadiendo los `key_outcomes` relevantes. Asegúrate de actualizar la fecha `last_updated` en el `sprint-registry.yaml`.

---

## 🔗 REFERENCIAS
- **ADRs relevantes:** [ADR-XXX](docs/architecture/adrs/ADR-XXX.md)
- **Documentos de negocio:** [Enlace]
