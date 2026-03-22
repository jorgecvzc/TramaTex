> **Nota:** Esta tarea fue reclasificada como Sprint 15 / Tarea 01.
> Ver documentación completa en [`docs/log/sprints/sprint-15/01-mvp-architectural-refinement.md`](../../sprint-15/01-mvp-architectural-refinement.md).

# Sprint 15 / Tarea 01 — Refinamiento Arquitectónico del MVP *(reclasificado desde 14-05)*

> ⚠️ **Contenido movido**. El detalle completo de esta tarea se encuentra en:
> [`docs/log/sprints/sprint-15/01-mvp-architectural-refinement.md`](../../sprint-15/01-mvp-architectural-refinement.md)


---

## 🎯 Objetivos

1. [x] Análisis sistemático de los 6 módulos del backend para identificar deuda técnica
2. [x] **P1**: Fragmentar `SalesService` (god service de 2.232 líneas)
3. [x] **P2**: Eliminar mapeo de errores manual en handlers; activar middleware global
4. [x] **P3**: Consolidar cálculos duplicados en el dominio Sales
5. [x] **P4**: Migrar `User.id` de `string` a `uuid.UUID` en IAM
6. [x] **IAM Cleanup**: Eliminar campos de auditoría de la entidad de dominio `User`

---

## 📊 Trabajo Realizado

### Fase de Análisis (2026-03-12 a 2026-03-20)
- Análisis modular de los 6 módulos: IAM, Party, Product, Pricing, Sales, MES
- Identificación de 4 propuestas priorizadas (`tmp/mvp_refinement_proposals.md`)
- Análisis de alineación doc-código (`tmp/technical-refinement-tasks.md`)

---

### P1 — Fragmentar SalesService (🔴 Crítico)

**Problema:** `sales_service.go` con 2.232 líneas y 53 métodos gestionando 4 dominios distintos.

**Solución:** Extracción en 4 servicios especializados:

| Archivo creado | Contenido | Líneas |
|---------------|-----------|--------|
| `quote_service.go` | Cotizaciones + ConvertToOrder | ~482 |
| `order_service.go` | Pedidos + LineItems | ~824 |
| `delivery_note_service.go` | Albaranes | ~253 |
| `billing_service.go` | Facturas + helpers | ~541 |
| `sales_service.go` | Orquestador (struct + constructor) | 247 |

**Resultado:** `sales_service.go` 2.232 → 247 líneas (reducción del 89%).

---

### P2 — Handlers Ligeros (🟠 Alto)

**Problema:** 3 mappers de errores locales duplicando lógica en `mes_handler.go`, `product_handler.go`, `sales_handler.go` (~200 líneas duplicadas).

**Solución:**
- Extendido `ErrorHandlerMiddleware` en `shared/infrastructure/middleware/` para reconocer `domain.DomainError`
- Creado `shared/domain/errors.go` con interfaz `HTTPStatuser`
- Creado `internal/mes/domain/errors.go` con `MESError`
- Eliminados los 3 mappers locales: `mapServiceError`, `mapErrorToHTTP`, `handleSalesError`
- Añadido `init()` en `logging/logger.go` para evitar nil panic en tests

---

### P3 — Consolidar Cálculos Duplicados (🟠 Alto)

**Problema:** La función de suma de importes implementada 6 veces (2 por cada entidad: Quote, SalesOrder, Invoice).

**Solución:**
- Creado `internal/sales/domain/calculations.go` con `SumAmounts([]Money) (Money, error)`
- Las 6 funciones privadas (`sumLineItemSubtotals`, etc.) delegan a `SumAmounts`

---

### P4 — Migrar User.id a uuid.UUID (🟡 Medio)

**Problema:** `User.id` era `string` en IAM mientras todos los demás módulos usan `uuid.UUID`, generando conversiones dispersas.

**Archivos modificados (11 ficheros + 3 tests):**

| Capa | Archivo | Cambio |
|------|---------|--------|
| Dominio | `user.go` | `id uuid.UUID`, `ID() uuid.UUID` |
| Repositorio | `user_repository.go` | `ByID/Delete(uuid.UUID)` |
| Persistencia | `postgres_user_repository.go` | Firmas + `uuid.Parse` en `modelToDomain` |
| Usecases (×8) | `assign_role`, `check_auth`, `delete_user`, `login`, `refresh`, `register`, `create_user`, `list_users` | `uuid.Parse` en frontera string→UUID; `.String()` en outputs |

**Estrategia de frontera:** JWT Subject y contexto HTTP permanecen como `string`; el dominio interno usa `uuid.UUID`.

---

### IAM Domain Cleanup — Eliminar campos de auditoría

**Problema:** `User` tenía `createdAt`/`updatedAt` en el dominio, violando la separación de preocupaciones (esos campos pertenecen a la capa de infraestructura/persistencia).

**Cambios en `user.go`:**
- Eliminados campos `createdAt time.Time` y `updatedAt time.Time`
- Eliminados getters `CreatedAt()` y `UpdatedAt()`
- Eliminadas 4 asignaciones `u.updatedAt = time.Now()` en métodos de mutación
- Eliminado import `"time"` y variable local `now`

Los timestamps de auditoría siguen viviendo donde corresponden: en `UserModel` de la capa de persistencia (GORM `autoCreateTime`/`autoUpdateTime`).

---

## 🔗 Commits en `mvp-arch-refinement`

| Hash | Descripción |
|------|-------------|
| `3faf17b` | `refactor(shared,mes,sales,product): activate global error middleware, remove local mappers (P2)` |
| `62358f4` | `refactor(sales): consolidate duplicate domain sum calculations (P3)` |
| `fd777bd` | `refactor(sales): fragment SalesService god service into 4 specialized files (P1)` |
| `7a6a76b` | `chore: remove temporary split script` |
| `0ceb78f` | `chore(sales): remove unused imports cleaned by goimports (P1 followup)` |
| `e9f9b2f` | `docs: mark P1, P2, P3 as implemented in session log and proposals` |
| `1104825` | `refactor(iam): migrate User.id from string to uuid.UUID (P4)` |
| `15769c9` | `docs: mark P4 as implemented; update session log` |
| `2577e84` | `refactor(iam): remove createdAt/updatedAt from User domain entity` |

---

## 📊 Métricas

| Métrica | Valor |
|---------|-------|
| Archivos creados | 6 (`calculations.go`, `quote_service.go`, `order_service.go`, `delivery_note_service.go`, `billing_service.go`, `shared/domain/errors.go`, `mes/domain/errors.go`) |
| Archivos modificados | ~22 |
| Líneas eliminadas (netas) | ~2.000+ |
| Tests actualizados | 5 ficheros |
| Build post-refactor | ✅ Verde |
| Tests post-refactor | ✅ Verde (todos los módulos) |

---

## ✅ Estado Final

- Todos los tests en verde (excepto `product/infrastructure/persistence` que requiere PostgreSQL activo — fallo pre-existente)
- Build limpio: `go build ./...` sin errores
- Rama `mvp-arch-refinement` lista para PR → `develop`
