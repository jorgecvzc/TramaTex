# Contratos de API - Módulo MES

Este documento detalla la interfaz de integración para la gestión de la producción y el control del taller en TramaTex.

---

## 1. Datos Maestros

### Tareas (`/api/mes/tasks`)
CRUD del catálogo de tareas atómicas.
- `POST /api/mes/tasks` — Crear tarea.
- `GET /api/mes/tasks` — Listar tareas (filtros: `is_active`, `search`).
- `GET /api/mes/tasks/:id` — Obtener tarea.
- `PUT /api/mes/tasks/:id` — Actualizar tarea.
- `DELETE /api/mes/tasks/:id` — Eliminar tarea.

### Posiciones (`/api/mes/positions`)
CRUD del catálogo de zonas de prenda.
- `POST /api/mes/positions` — Crear posición.
- `GET /api/mes/positions` — Listar posiciones (filtros: `is_active`, `search`).
- `GET /api/mes/positions/:id` — Obtener posición.
- `PUT /api/mes/positions/:id` — Actualizar posición.
- `DELETE /api/mes/positions/:id` — Eliminar posición.

### Tipos de Trabajo (`/api/mes/work-types`)
CRUD de recetas (secuencias de tareas).
- `POST /api/mes/work-types` — Crear tipo de trabajo con secuencia de tareas.
- `GET /api/mes/work-types` — Listar tipos (filtros: `is_active`, `search`).
- `GET /api/mes/work-types/:id` — Obtener tipo con sus tareas.
- `PUT /api/mes/work-types/:id` — Actualizar tipo de trabajo.
- `DELETE /api/mes/work-types/:id` — Eliminar tipo de trabajo.

---

## 2. Configuración por Cliente (`/api/mes/work-setups`)

Gestiona las plantillas de personalización por cliente y tipo de prenda.
- `POST /api/mes/work-setups` — Crear configuración con líneas (WorkType + Position).
- `GET /api/mes/work-setups` — Listar configuraciones (filtros: `is_active`, `search`, `party_id`).
- `GET /api/mes/work-setups/:id` — Obtener configuración con sus líneas.
- `PUT /api/mes/work-setups/:id` — Actualizar configuración.
- `DELETE /api/mes/work-setups/:id` — Eliminar configuración.

---

## 3. Órdenes de Trabajo (`/api/mes/work-orders`)

Motor operativo para la ejecución real en el taller.
- `POST /api/mes/work-orders` — Crear orden (desde WorkSetup o manualmente). El body puede incluir `order_work_setup_id` (UUID de la solicitud de Sales) para vincular automáticamente la orden al pedido al crearse.
- `GET /api/mes/work-orders` — Listar órdenes (filtros: `status`, `search`, `party_id`).
- `GET /api/mes/work-orders/:id` — Obtener orden con líneas y tareas.
- `PUT /api/mes/work-orders/:id` — Actualizar orden. También es la vía indirecta para las transiciones SUSPENDED ↔ PENDING cuando Sales suspende/reactiva un pedido (ver sección de Integración).
- `GET /api/mes/work-orders/dashboard/stats` — Estadísticas del dashboard.
- `GET /api/mes/work-orders/overdue` — Órdenes con retraso.

## 4. Solicitudes Pendientes de Sales (`/api/mes/pending-work-setups`)

Endpoint consumido exclusivamente por el Dashboard de MES para mostrar las configuraciones de pedidos confirmados que aún no tienen Orden de Trabajo creada.
- `GET /api/mes/pending-work-setups` — Retorna lista de `PendingWorkSetupDTO` con: `id`, `workSetupId`, `description`, `orderId`, `orderNumber`, `deliveryDate`, `partyId`. Delega internamente en `PendingWorkSetupProvider` (implementado por Sales vía setter injection).

---

## 4. Terminal de Taller (`/api/mes/work-orders/:orderId/tasks/:taskId`)

Endpoints optimizados para operarios en planta (tablets).
- `PUT /api/mes/work-orders/:orderId/tasks/:taskId` — Cambiar estado de tarea.
  - Body: `{ "action": "START" | "COMPLETE" | "BLOCK", "notes": "..." }`
  - Header: `X-User-ID` (operario autenticado vía IAM).

---

## Integración Sales → MES (Cross-Module)

El módulo MES expone tres interfaces que Sales implementa mediante adaptadores (setter injection):

| Interfaz (definida en MES) | Implementación (en Sales infrastructure) | Propósito |
|---|---|---|
| `SalesOrderLinker` | `SalesOrderLinkerAdapter` | Al crear una WorkOrder con `order_work_setup_id`, vincula automáticamente la orden al `order_work_setups` de Sales. |
| `PendingWorkSetupProvider` | `PendingSetupProviderAdapter` | Provee al endpoint `/pending-work-setups` la lista de configuraciones sin orden creada, consultando tablas de Sales. |
| *(en Sales)* `WorkOrderSuspender` | `WorkOrderSuspenderAdapter` | Sales llama a `SuspendWorkOrders` / `ReactivateWorkOrders` en `MESService` al cancelar o reactivar un pedido. |

**Nota:** No existe ningún endpoint HTTP explícito de "suspender/reactivar" en la API MES pública. La suspensión ocurre internamente entre módulos a través de la inyección de dependencias en el proceso del servidor.

---

## Eventos de Dominio

El módulo MES no emite eventos de dominio en el MVP. La visibilidad bidireccional se logra mediante **consulta bajo demanda**:
- Sales consulta el estado de las WorkOrders vía la interfaz `MESWorkLookup` → adaptador `MESWorkLookupAdapter`.
- Los eventos `WorkOrderStarted` / `WorkOrderCompleted` están planificados para **Post-MVP**.

---
**Última Actualización:** 20 de marzo de 2026
