# Tarea 17-02: Migrar Endpoints Admin a PricingEngineHandler

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 02 |
| **ID de Sprint** | sprint-17 |
| **Título** | Migrar endpoints admin del PricingHandler al PricingEngineHandler |
| **Estado** | ✅ Completado |
| **Facilitador/LLM** | GitHub Copilot / Claude Opus 4.6 |
| **Fecha de Inicio** | 2026-04-10 |
| **Fecha de Fin** | 2026-04-12 |
| **Duración Estimada** | 2-3 horas |
| **Prioridad** | P1 — Alto |
| **Dependencia** | Requiere 17-01 completada |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [ ] **Mover CRUD de ClientPricing overrides** al `PricingEngineHandler`
2. [ ] **Mover endpoint de historial de cálculos** al `PricingEngineHandler`
3. [ ] **Mover listado de PricingRules** (o deprecar si ya no aplica con modelo ADR-016)
4. [ ] **Actualizar rutas** en el router (Gin)
5. [ ] **Verificar que el frontend** apunta a las nuevas rutas (o mantener compatibilidad)

---

## 📊 CONTEXTO

### Endpoints actuales en PricingHandler (viejo)

| Método | Ruta | Función | Destino |
|--------|------|---------|---------|
| POST | `/api/pricing/calculate` | CalculatePrice | Se elimina (reemplazado por engine) |
| GET | `/api/pricing/rules` | ListPricingRules | Evaluar: ¿migrar o deprecar? |
| POST | `/api/pricing/rules` | CreatePricingRule | Evaluar: ¿migrar o deprecar? |
| POST | `/api/pricing/client-override` | CreateClientPricingOverride | Migrar a PricingEngineHandler |
| GET | `/api/pricing/history/:variantId` | GetPricingHistory | Migrar a PricingEngineHandler |

### Endpoints actuales en PricingEngineHandler (nuevo)

- CRUD de BaseSalesPriceRule
- CRUD de SaleModificationRule
- CalculateBaseSalesPrice
- CalculateFinalSalePrice

### Decisión sobre PricingRule/BrandProfitMargin/SalesDiscountRule

Estas entidades del modelo viejo quedarán en desuso tras la consolidación. Las reglas equivalentes se expresan como `BaseSalesPriceRule` y `SaleModificationRule` con `RuleValue` tipado (ADR-016). Se debe:
- **Mantener lectura** temporal si hay datos históricos
- **No permitir creación** de nuevas reglas viejas
- **Deprecar** en la Fase 3

---

## 🛠️ PLAN DE TRABAJO

### Paso 1: Añadir métodos a PricingEngineService
- `CreateClientPricingOverride` (reutilizar lógica de PricingService)
- `GetPricingHistory` (consultar PriceCalculationRepository)

### Paso 2: Añadir handlers a PricingEngineHandler
- POST `/api/pricing/engine/client-overrides`
- GET `/api/pricing/engine/history/:variantId`

### Paso 3: Actualizar router en main.go
- Registrar nuevas rutas
- Mantener viejas como redirect temporal (o eliminar si frontend se actualiza)

### Paso 4: Actualizar frontend
- Verificar `pricingApi.ts` y ajustar URLs si cambian

### Paso 5: Tests
- Tests de handler para nuevos endpoints
- Verificar que endpoints viejos devuelven 301/410 (deprecated)
