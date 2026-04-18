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

## Continuidad IAM + Sales — Estado Actual para Retomar (Sprint 18)
- **Session ID:** `sprint-18-iam-sales-continuidad-2026-04-18`
- **Status:** En Progreso
- **Sprint:** Sprint 18
- **Started:** 2026-04-18
- **Rama:** `iam-role-alignment`
- **Contexto:** Se prepara continuidad para retomar mañana con el estado actual consolidado. En `iam-role-alignment` están publicados los commits IAM y la rama está sincronizada con origen. Existen utilidades locales no versionadas de soporte DB/migraciones y un directorio `migrations_new` pendiente de decisión (integrar, mover a rama específica o descartar). Además, hay trabajo reciente en fixes de Sales (impresión, cobro de facturas y recálculo tras borrar albarán) que quedó separado en `fix/sales-bugfixes` y debe decidirse su estrategia de integración.
- **Checklist de Arranque (mañana):**
    - [ ] Verificar rama actual y estado limpio de trabajo antes de tocar código
    - [ ] Decidir tratamiento de `migrations_new` y utilidades locales no versionadas
    - [ ] Ejecutar smoke rápido IAM + Sales para confirmar punto de partida
- **Próximos Pasos:**
    - [ ] Validar si `apps/tramatex-api/migrations_new/` se conserva para una tarea formal o se limpia del working tree
    - [ ] Revisar y clasificar utilidades locales no versionadas (`*_db.go`, `*_migration.go`, validadores) para evitar ruido en commits
    - [ ] Confirmar estado de integración de `fix/sales-bugfixes` respecto a `staging`
    - [ ] Ejecutar validación final de flujo IAM (alta/edición de usuarios) y flujo Sales (cobro factura + borrado albarán)
    - [ ] Dejar rama objetivo limpia y lista para PR/merge según decisión funcional
- **Archivos de Contexto:**
    - `docs/log/sprints/sprint-18/01-iam-role-alignment.md`
    - `apps/frontend/src/services/iam.ts`
    - `apps/frontend/src/pages/admin/UsersManagement.vue`
    - `apps/frontend/src/types/auth.ts`
    - `apps/tramatex-api/internal/iam/domain/model/role.go`
    - `apps/frontend/src/assets/sales-print.css`
    - `apps/frontend/src/components/sales/PrintDocument.vue`
    - `apps/frontend/src/pages/sales/InvoiceDetail.vue`
    - `apps/frontend/src/pages/sales/DeliveryNoteDetail.vue`
    - `apps/tramatex-api/internal/sales/application/sales_service_test.go`

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

- **POS Polish & MES Terminal Excellence (Sprint 18)** | Iniciada: 2026-04-18 | Finalizada: 2026-04-18 | ✅ TPV industrial optimizado (sin scroll, atajos en cabecera), Terminal MES tabular con recuperación de cliente y Dashboard con KPIs reales.
- **Fix Production Launch Errors (Sales/IAM) (Sprint 18)** | Iniciada: 2026-04-11 | Finalizada: 2026-04-18 | ✅ Bug de token corregido y flujo de lanzamiento a producción validado.
- **Review Sprint 17 Implementation Status (Sprint 17 / 18)** | Iniciada: 2026-04-12 | Finalizada: 2026-04-18 | ✅ Tareas de consolidación de Pricing, Money decimal y ACL clients verificadas y completadas.
- **IAM Role Alignment — Fix DB↔Domain Role Mismatch (Sprint 18)** | Iniciada: 2026-04-10 | Finalizada: 2026-04-11 | ✅ Roles alineados y bug `updateUser` corregido.
- **Pricing Consolidation — Paridad Funcional + Eliminación Legacy (Sprint 17)** | Iniciada: 2026-04-08 | Finalizada: 2026-04-11 | ✅ Motor de precios optimizado y UI saneada.
- **POST-REFACTOR STABILIZATION & FINAL POLISH (Sprint 16)** | Iniciada: 2026-04-07 | Finalizada: 2026-04-08 | Mergeada a producción

