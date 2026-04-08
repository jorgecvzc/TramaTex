# Bitácora de Sesiones de Desarrollo

---
## SESIONES ABIERTAS
---

## Pricing Consolidation — Paridad Funcional + Eliminación Legacy (Sprint 17)
- **Session ID:** `sprint-17-pricing-consolidation`
- **Status:** Planificado
- **Sprint:** Sprint 17
- **Started:** 2026-04-08 (planificación)
- **Rama:** `feature/pricing-consolidation` (desde `develop`)
- **Context:** Estudio en profundidad del módulo Pricing reveló dos hallazgos críticos: (G1) `PricingEngineService.CalculateFinalSalePrice` no consulta ClientPricing overrides (la regla de máxima prioridad del dominio) y (G2) no genera registros de auditoría `PriceCalculation`. Sales sólo consume `PricingEngineService`, por lo que estos gaps afectan directamente a producción. El sprint consolida el motor nuevo (ADR-016), elimina el sistema dual legacy y corrige violaciones de Clean Architecture.
- **Tareas planificadas:**
    - **17-01** (P0): Completar PricingEngineService — inyectar ClientPricingRepo + PriceCalculationRepo
    - **17-02** (P1): Migrar endpoints admin de PricingHandler a PricingEngineHandler
    - **17-03** (P1): Eliminar PricingService, PricingHandler, entidades y tablas obsoletas
    - **17-04** (P1): Migrar Money VO de float64 a decimal.Decimal
    - **17-05** (P2): Fix ACL — ProductClient y PartyClient (imports cruzados + SQL directo)
    - **17-06** (P2): Fix mapeo DomainError → HTTP Status en handlers
- **Orden de ejecución:** 17-01 → 17-02 → 17-03 (secuenciales) → 17-04 → 17-05 + 17-06 (paralelizables)
- **Sprint siguiente previsto:** Estabilización UI/UX + implementación de TODOS.md

---

## POST-REFACTOR STABILIZATION & FINAL POLISH (Sprint 16)
- **Session ID:** `sprint-16-stabilization-final`
- **Status:** En Progreso
- **Sprint:** Sprint 16
- **Started:** 2026-04-07
- **Context:** Estabilización general tras refactor, unificación de UI/UX, corrección de integración Sales-MES y refinamiento intensivo del TPV.
- **Logros de Hoy (2026-04-07):**
    - ✅ **TPV: Rediseño horizontal compacto** de las líneas del ticket (Cant | PVP | Dto% | Total).
    - ✅ **TPV: Integración total ADR-015 (Pricing)**. Precios de venta exactos incluyendo márgenes de marca y recargos por atributos (ej. tallas) desde el primer segundo.
    - ✅ **TPV: Gestión de Descuentos**. Sincronización automática con el cliente, redondeo a 2 decimales y visibilidad total.
    - ✅ **TPV: Cliente por defecto**. Implementado **CONSUMIDOR FINAL** como cliente fijo al inicio y tras ventas.
    - ✅ **TPV: Impresión Industrial (80mm)**. Solucionada la superposición de menús y el centrado del ticket mediante Vue Teleport y reglas CSS globales.
    - ✅ **Backend Pricing:** Corregida lógica para sumar modificadores de atributos al coste base *antes* de aplicar márgenes comerciales.
    - ✅ **Backend SQL:** Eliminadas comillas dobles en el repositorio Party que causaban errores de "relation not found" en PostgreSQL.
    - ✅ **Sales API:** Exportados mapeos de estado para corregir errores de referencia en la emisión de presupuestos.
- **Próximos Pasos (Mañana):**
    - 🔍 Realizar una venta real completa en el TPV y verificar la persistencia final en la base de datos.
    - 🔍 Verificar la creación de órdenes MES desde solicitudes de ventas (flujo integral).
    - ⌨️ Iniciar estudio de navegación 100% por teclado una vez estabilizada la funcionalidad.
- **Archivos de Contexto:**
    - `apps/frontend/src/pages/sales/TicketCreate.vue`
    - `apps/tramatex-api/internal/pricing/application/pricing_engine_service.go`
    - `apps/frontend/src/services/salesApi.ts`

---
## REGISTRO DE SESIONES CERRADAS
---
