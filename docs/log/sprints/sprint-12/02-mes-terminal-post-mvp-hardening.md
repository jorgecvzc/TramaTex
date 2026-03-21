# Sprint 12 / Tarea 02 - MES Terminal Hardening Post-MVP

> **⚠️ AVISO:** Las referencias a "grupo de servicio" en este documento corresponden al nombre legacy.
> Desde la sesión F5 (2026-03-16), la nomenclatura canónica es: `WorkType` (Tipo de Trabajo),
> `WorkOrder` (Orden de Trabajo), `WorkOrderLine` (Línea), `WorkOrderTask` (Tarea).

**Estado:** 📌 Pendiente (Post-MVP)  
**Fecha de Registro:** 2026-02-21  
**Origen:** Decisión de alcance MVP (Sprint 12-01)  
**Dependencia:** `01-mes-module-foundation.md` completado

---

## 🎯 Objetivo

Endurecer las reglas de transición y bloqueo del flujo operativo MES (tablet/workshop) sin cambiar el alcance funcional ya entregado en MVP.

---

## ✅ Alcance Ya Implementado en MVP (No Pendiente)

- Endpoint backend para transición de estado de tareas MES.
- UI de terminal tablet/workshop con acciones `START`, `PAUSE`, `COMPLETE`, `BLOCK`.
- Recalculo de estado de trabajo y actualización visual tras acciones.

---

## 📋 Checklist Técnico Post-MVP

### Backend (Reglas de Dominio y Application)

- [ ] Definir matriz formal de transiciones permitidas por estado actual de tarea.
- [ ] Rechazar `COMPLETE` si la tarea nunca pasó por `IN_PROGRESS`.
- [ ] Validar secuencia por grupo de servicio (no completar tarea `n` si `n-1` no está cerrada, cuando aplique).
- [ ] Exigir nota obligatoria para acción `BLOCK` con longitud mínima configurable.
- [ ] Registrar `blocked_at` y `blocked_reason` de forma explícita en persistencia o bitácora de eventos.
- [ ] Normalizar códigos de error de transición inválida para uso consistente en frontend.
- [ ] Añadir tests table-driven para transiciones válidas/ inválidas y casos límite de secuencia.

### Frontend (UX Operativa y Validación)

- [ ] Mostrar mensaje de validación específico por regla rechazada (no genérico).
- [ ] Solicitar motivo obligatorio al bloquear (`BLOCK`) antes de enviar al backend.
- [ ] Deshabilitar acciones no válidas según estado actual para reducir errores de operador.
- [ ] Refrescar fila/tarea con estado optimista controlado + rollback ante error.
- [ ] Mostrar badge/tooltip con último motivo de bloqueo cuando exista.
- [ ] Añadir pruebas unitarias de acciones de terminal y manejo de errores de transición.

### Observabilidad y Operación

- [ ] Auditar transición de tarea (quién, cuándo, acción, motivo) para trazabilidad operativa.
- [ ] Definir métricas mínimas: bloqueos por turno, tiempo medio en bloqueo, ratio de reintento.
- [ ] Añadir guía operativa corta para supervisores de taller (resolución de bloqueos).

---

## 🧪 Criterios de Aceptación

- [ ] Todas las reglas de transición definidas y cubiertas con tests backend.
- [ ] Flujo de bloqueo con motivo obligatorio activo en frontend y backend.
- [ ] Errores de transición devueltos con código/mensaje estable y documentado.
- [ ] Evidencia de trazabilidad (logs/eventos) verificable en entorno local.
- [ ] Build frontend y tests backend relacionados en verde.

---

## 📎 Referencias

- `docs/log/sprints/sprint-12/01-mes-module-foundation.md`
- `apps/frontend/src/pages/mes/terminal/Tablet.vue`
- `apps/frontend/src/services/mesApi.ts`
- `apps/tramatex-api/internal/mes/application/mes_service.go`
- `apps/tramatex-api/internal/mes/interfaces/http/handler/mes_handler.go`
- `apps/tramatex-api/cmd/api/main.go`
