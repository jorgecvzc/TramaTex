# Guía de Implementación - Módulo MES

**Estado:** ✅ Implementado (MVP)  
**Última Actualización:** 17 de marzo de 2026

---

## 1. Estructura del Módulo (Backend)

```
apps/tramatex-api/internal/mes/
├── domain/
│   └── entities.go          # Task, Position, WorkType, WorkSetup, WorkOrder + sub-entidades
├── application/
│   ├── mes_service.go       # Servicio de aplicación (casos de uso)
│   ├── commands.go          # CreateWorkTypeCommand, CreateWorkOrderCommand, etc.
│   ├── queries.go           # GetWorkTypeByIDQuery, ListWorkOrdersQuery, etc.
│   └── dtos.go              # DTOs de respuesta (WorkTypeDTO, WorkOrderDTO, etc.)
├── infrastructure/
│   └── persistence/
│       └── gorm_repositories.go  # Implementación GORM de repositorios
└── interfaces/
    └── http/
        └── handler/
            └── mes_handler.go    # Handlers HTTP (Gin)
```

## 2. Entidades de Dominio

| Entidad | Tabla BD | Descripción |
|---------|----------|-------------|
| `Task` | `tasks` | Tarea atómica (Diseñar, Imprimir, Marcar...) |
| `Position` | `positions` | Zona de prenda (Pecho, Espalda...) |
| `WorkType` | `work_types` | Receta de tareas (secuencia ordenada) |
| `WorkTypeTask` | `work_type_tasks` | Asignación tarea→tipo con secuencia |
| `WorkSetup` | `work_setups` | Configuración por cliente/prenda |
| `WorkSetupLine` | `work_setup_lines` | Línea de configuración (WorkType + Position) |
| `WorkOrder` | `work_orders` | Orden operativa de producción |
| `WorkOrderLine` | `work_order_lines` | Línea instanciada (WorkType + Position + tareas) |
| `WorkOrderTask` | `work_order_tasks` | Tarea instanciada con estado y tiempos |

> **Nota:** Las tablas se llamaban originalmente `service_groups`, `mes_works`, etc. Fueron renombradas en `migrations/028_rename_legacy_mes_tables.sql`.

## 3. Rutas API

Todas bajo el prefijo `/api/mes/`:

| Grupo | Ruta base | Operaciones |
|-------|-----------|-------------|
| Tasks | `/tasks` | CRUD completo |
| Positions | `/positions` | CRUD completo |
| Work Types | `/work-types` | CRUD completo (incluye task assignments) |
| Work Setups | `/work-setups` | CRUD completo (incluye líneas) |
| Work Orders | `/work-orders` | CRUD + `/dashboard/stats` + `/overdue` |
| Terminal | `/work-orders/:id/tasks/:taskId` | Cambio de estado de tarea (START/COMPLETE/BLOCK) |
| Pending Setups | `/pending-work-setups` | GET — solicitudes de Sales pendientes de WorkOrder |

## 4. Flujo de Creación de WorkOrder

1. Se recibe `CreateWorkOrderCommand` con `WorkSetupID`.
2. Se obtiene el `WorkSetup` con sus líneas.
3. Para cada línea del setup, se obtiene el `WorkType` correspondiente.
4. Se generan `WorkOrderLine`s con `WorkOrderTask`s copiadas de las tareas del WorkType.
5. Se genera `WorkNumber` con formato `MES-{año}-{secuencial}`.
6. Se persiste la orden completa.

## 5. Integración con Sales

- Los pedidos (`SalesOrder`) contienen registros en `order_work_setups` que referencian `work_setup_id` y opcionalmente `work_order_id`.
- Al confirmar un pedido, se garantiza que cada entrada tenga un `work_setup_id` válido (auto-creación si necesario).
- El Dashboard de MES consume `GET /api/mes/pending-work-setups` para mostrar las solicitudes pendientes del taller.

### Interfaces Cross-Módulo (Setter Injection)

| Interfaz | Módulo que la define | Implementación | Setter en MESService |
|---|---|---|---|
| `SalesOrderLinker` | MES (application) | `SalesOrderLinkerAdapter` (Sales infra) | `SetSalesOrderLinker()` |
| `PendingWorkSetupProvider` | MES (application) | `PendingSetupProviderAdapter` (Sales infra) | `SetPendingSetupProvider()` |
| `WorkOrderSuspender` | Sales (application) | `WorkOrderSuspenderAdapter` (Sales infra) | — (inyecta `MESService` directamente) |
| `MESWorkLookup` | Sales (application) | `MESWorkLookupAdapter` (Sales infra) | — (inyecta `MESService` directamente) |

### Flujo Suspend/Reactivate

Cuando un pedido de Sales es cancelado o reactivado, `SalesService.ChangeOrderStatus()` llama a `WorkOrderSuspender.SuspendWorkOrders()` o `ReactivateWorkOrders()` con los `WorkOrderID`s asociados. `MESService` aplica la transición de estado solo a las órdenes elegibles (ver CU-M-013).

## 6. Cobertura de Tests

| Capa | Cobertura | Tests |
|------|-----------|-------|
| `mes/application` | 57.1% | 18 tests unitarios (fake repos) |
| `mes/domain` | 87.1% | 8 tests (entidades + validación) |

## 7. Gestión de Errores Estandarizada

El módulo `MES` delega la traducción de errores de dominio a respuestas HTTP en el `ErrorHandlerMiddleware` de la capa `shared`. 

Para que esto funcione:
1. **Definir Errores en Dominio**: Todos los errores de negocio se definen en `internal/mes/domain/errors.go`.
2. **Implementar `HTTPStatuser`**: Los errores de dominio implementan la interfaz `shared/domain.HTTPStatuser` para indicar su código HTTP correspondiente (ej. `ErrWorkOrderNotFound` devuelve `404`).
3. **Delegación en Handlers**: Los controladores Gin NO deben formatear respuestas de error manualmente. Deben simplemente adjuntar el error al contexto: `c.Error(err)`. El middleware se encargará de sanitizar la respuesta y registrar el log con el ID de petición.

---
**Referencia:** Ver [module-spec.md](module-spec.md) para la especificación funcional completa.

---
**Última Actualización:** 25-03-2026
