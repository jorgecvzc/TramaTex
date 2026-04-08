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

---
## REGISTRO DE SESIONES CERRADAS
---

- **POST-REFACTOR STABILIZATION & FINAL POLISH (Sprint 16)** | Iniciada: 2026-04-07 | Finalizada: 2026-04-08 | Mergeada a producción
