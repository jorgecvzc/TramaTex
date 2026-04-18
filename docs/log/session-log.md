# Bitácora de Sesiones de Desarrollo

---
# SESIONES ABIERTAS
---

## Documentation Review & Alignment
- **Session ID:** `documentation-review-sprint-18`
- **Status:** En Progreso
- **Sprint:** Sprint 18
- **Started:** 2026-04-17
- **Contexto:** Revisión integral de la documentación del proyecto para alinearla con los cambios recientes en la arquitectura (Pricing Engine), flujos de venta y estabilización de módulos.
- **Avance (2026-04-18):** Preservado y versionado el deck HTML de presentación (`docs/presentations/presentation.html`) en la rama `docs-consolidation`.
- **Próximos Pasos:**
  - [ ] Revisar `docs/architecture/` y actualizar diagramas o visiones obsoletas.
  - [ ] Alinear la documentación de módulos en `docs/modules/` (especialmente Pricing y Sales).
  - [ ] Actualizar guías de desarrollo en `docs/guides/` con los nuevos estándares.
  - [ ] Verificar consistencia en el glosario y visión del proyecto.
- **Archivos de Contexto:**
  - `docs/`
  - `docs/architecture/`
  - `docs/modules/`
  - `docs/guides/`

---

## Review Sprint 17 Implementation Status
- **Session ID:** `review-sprint-17-completion`
- **Status:** En Progreso
- **Sprint:** Sprint 17 / 18
- **Started:** 2026-04-12
- **Contexto:** Verificar si las tareas del Sprint 17 (Consolidación de Pricing) han sido totalmente implementadas o si quedan pendientes tareas planificadas.
- **Próximos Pasos:**
  - [x] Validar implementación de `PricingEngineService` (Tarea 17-01).
  - [x] Verificar si los endpoints admin ya fueron migrados (Tarea 17-02).
  - [x] Comprobar si el sistema legacy aún existe (Tarea 17-03).
  - [x] Verificar el tipo de dato en `Money` (Tarea 17-04).
  - [x] Fix: Recrear tablas `parties`, `products`, `brands`, `product_groups`, `attributes`, `product_variants`.
  - [x] Comprobar funcionamiento completo de Sales tras recrear tablas.
  - [x] FIX: Error `invalid quote status: EMITIDA` (Añadida tolerancia a alias en español en Quote, Invoice y DeliveryNote).
- **Archivos de Contexto:**
  - `docs/log/sprints/sprint-17/`
  - `apps/tramatex-api/internal/pricing/application/pricing_engine_service.go`
  - `apps/tramatex-api/internal/pricing/interfaces/http/handler/pricing_engine_handler.go`

---

## Fix Production Launch Errors (Sales/IAM)
- **Session ID:** `fix-prod-launch-errors`
- **Status:** En Progreso
- **Sprint:** Sprint 18
- **Started:** 2026-04-11
- **Contexto:** Se han detectado varios errores críticos al intentar "Lanzar a Producción" desde el detalle de un pedido.
- **Próximos Pasos:**
  - [x] Corregir validación de estados en el dominio (tolerancia a alias en español e inglés).
  - [x] Añadir migración SQL para estados faltantes en el ENUM de la BD.
  - [x] Corregir bug en `api.ts` que impedía obtener el token (`token` -> `accessToken`).
  - [x] Refactorizar `SalesApi.ts` y `MESApi.ts` para usar `axios` y centralizar la autenticación.
  - [ ] Verificar si el error `Authorization header is missing` persiste en producción.
  - [ ] Validar flujo completo de lanzamiento a producción.
- **Archivos de Contexto:**
  - `apps/tramatex-api/internal/sales/domain/statuses.go`
  - `apps/tramatex-api/migrations/010_align_sales_enums.sql`
  - `apps/frontend/src/services/api.ts`
  - `apps/frontend/src/services/salesApi.ts`
  - `apps/frontend/src/services/mesApi.ts`

---
# REGISTRO DE SESIONES CERRADAS
---

- **IAM Role Alignment — Fix DB↔Domain Role Mismatch (Sprint 18)** | Iniciada: 2026-04-10 | Finalizada: 2026-04-11 | ✅ Roles alineados y bug `updateUser` corregido.
- **Pricing Consolidation — Paridad Funcional + Eliminación Legacy (Sprint 17)** | Iniciada: 2026-04-08 | Finalizada: 2026-04-11 | ✅ Motor de precios optimizado y UI saneada.
- **POST-REFACTOR STABILIZATION & FINAL POLISH (Sprint 16)** | Iniciada: 2026-04-07 | Finalizada: 2026-04-08 | Mergeada a producción

