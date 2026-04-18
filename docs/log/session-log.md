# Bitácora de Sesiones de Desarrollo

---
## SESIONES ABIERTAS
---

## Pricing Consolidation — Paridad Funcional + Eliminación Legacy (Sprint 17)
- **Session ID:** `sprint-17-pricing-consolidation`
- **Status:** Planificado
- **Sprint:** Sprint 17
- **Started:** 2026-04-08
- **Rama:** `feature/pricing-consolidation`
- **Contexto:** Estudio en profundidad del módulo Pricing reveló dos hallazgos críticos: (G1) PricingEngineService.CalculateFinalSalePrice no consulta ClientPricing overrides (la regla de máxima prioridad del dominio) y (G2) no genera registros de auditoría PriceCalculation. Sales sólo consume PricingEngineService, por lo que estos gaps afectan directamente a producción. El sprint consolida el motor nuevo (ADR-016), elimina el sistema dual legacy y corrige violaciones de Clean Architecture.
- **Próximos Pasos:**
    - [ ] 17-01 (P0): Completar PricingEngineService — inyectar ClientPricingRepo + PriceCalculationRepo
    - [ ] 17-02 (P1): Migrar endpoints admin de PricingHandler a PricingEngineHandler
    - [ ] 17-03 (P1): Eliminar PricingService, PricingHandler, entidades y tablas obsoletas
    - [ ] 17-04 (P1): Migrar Money VO de float64 a decimal.Decimal
    - [ ] 17-05 (P2): Fix ACL — ProductClient y PartyClient
    - [ ] 17-06 (P2): Fix mapeo DomainError → HTTP Status en handlers
- **Orden de ejecución:** 17-01 → 17-02 → 17-03 (secuenciales) → 17-04 → 17-05 + 17-06 (paralelizables)
- **Archivos de Contexto:**
    - `docs/log/sprints/sprint-17/01-pricing-engine-completar-paridad.md`
    - `docs/modules/pricing/domain-model.md`
    - `docs/modules/pricing/module-spec.md`
    - `docs/architecture/adrs/adr-016-pricing-module-architecture.md`
    - `docs/modules/pricing/implementation-guide.md`
    - `apps/tramatex-api/internal/pricing/application/pricing_engine_service.go`
    - `apps/tramatex-api/internal/pricing/application/pricing_service.go`
    - `apps/tramatex-api/internal/pricing/domain/repository.go`

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
## REGISTRO DE SESIONES CERRADAS
---

- **POST-REFACTOR STABILIZATION & FINAL POLISH (Sprint 16)** | Iniciada: 2026-04-07 | Finalizada: 2026-04-08 | Mergeada a producción
