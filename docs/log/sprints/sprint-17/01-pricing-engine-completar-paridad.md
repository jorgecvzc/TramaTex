# Tarea 17-01: Completar PricingEngineService (Paridad Funcional)

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 01 |
| **ID de Sprint** | sprint-17 |
| **Título** | Completar PricingEngineService: ClientPricing Overrides + Audit Trail |
| **Estado** | ✅ Completado |
| **Facilitador/LLM** | GitHub Copilot / Claude Opus 4.6 |
| **Fecha de Inicio** | 2026-04-08 |
| **Fecha de Fin** | 2026-04-10 |
| **Duración Estimada** | 3-4 horas |
| **Prioridad** | P0 — Crítico |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [x] **Inyectar `ClientPricingRepository` en `PricingEngineService`** y consultar overrides como paso prioritario en `CalculateFinalSalePrice`
2. [x] **Inyectar `PriceCalculationRepository` en `PricingEngineService`** y generar registros de auditoría inmutables al final de cada cálculo
3. [x] **Tests unitarios** para ambas funcionalidades nuevas (≥85% cobertura, módulo STRICT)

---

## 📊 CONTEXTO DE ENTRADA

### Hallazgos del Estudio de Pricing (esta sesión)

**G1 — CRÍTICO: Client Overrides ausentes del motor nuevo**
- `domain-model.md §2`: "Acuerdo Particular (Client Override) = prioridad máxima"
- `use-cases.md UC-PRI-001 paso 4`: "Si existe ClientPricing, tiene prioridad sobre reglas generales"
- `PricingService.CalculatePrice` (viejo) → SÍ consulta `clientPricingRepo.FindApplicable()`
- `PricingEngineService.CalculateFinalSalePrice` (nuevo) → **NO** inyecta ni consulta `ClientPricingRepository`
- Sales **sólo** consume `PricingEngineService` → la regla de máxima prioridad **nunca se aplica** en ventas reales

**G2 — CRÍTICO: Audit trail ausente del motor nuevo**
- `module-spec.md §7 paso 7`: "Registrar el cálculo en PriceCalculation"
- `domain-model.md §3`: "genera un registro de auditoría que congela las reglas aplicadas"
- `PricingService.CalculatePrice` → SÍ llama `calculationRepo.Save()`
- `PricingEngineService` → **NO** inyecta `PriceCalculationRepository`, **nunca** guarda auditoría

### ADRs Aplicables

- **ADR-016**: BaseSalesPriceRule + SaleModificationRule = modelo canónico
- **ADR-002**: Clean Architecture — domain puro, sin dependencias externas
- **ADR-006**: DDD con Rigor Asimétrico — Pricing = STRICT (≥85%)

### Archivos Clave a Modificar

| Archivo | Cambio |
|---------|--------|
| `internal/pricing/application/pricing_engine_service.go` | Inyectar ClientPricingRepo + PriceCalculationRepo, modificar CalculateFinalSalePrice |
| `internal/pricing/domain/repository.go` | Ya tiene las interfaces (no requiere cambios) |
| `cmd/api/main.go` | Pasar repos adicionales al constructor de PricingEngineService |
| Tests nuevos | Cubrir: override encontrado, override no encontrado, override expirado, audit trail generado |

### Arquitectura del Cambio

```
CalculateFinalSalePrice (flujo actual):
  1. Obtener BSP (cache o cálculo)
  2. Aplicar SaleModificationRules
  3. Fallback: clientDefaultDiscount
  4. Calcular impuestos y totales

CalculateFinalSalePrice (flujo corregido):
  1. Obtener BSP (cache o cálculo)
  2. ★ NUEVO: Consultar ClientPricing override para (clientID, variantID, fecha)
     → Si existe: usar precio fijo como finalPrice, saltar paso 3
  3. Aplicar SaleModificationRules
  4. Fallback: clientDefaultDiscount
  5. Calcular impuestos y totales
  6. ★ NUEVO: Guardar PriceCalculation (audit trail inmutable)
```

### Interfaz de Sales (NO modificar)

Sales define su contrato así (confirmado en `internal/sales/domain/`):
```go
type PricingEngine interface {
    CalculateFinalSalePrice(ctx context.Context, req CalculateFinalSalePriceRequest) (*CalculateFinalSalePriceResponse, error)
}
```
La firma no cambia. Solo cambia la implementación interna.

---

## 🛠️ PLAN DE TRABAJO

### Paso 1: Modificar constructor de PricingEngineService
- Añadir `clientPricingRepo domain.ClientPricingRepository` y `calculationRepo domain.PriceCalculationRepository`
- Actualizar `NewPricingEngineService()` para aceptar los nuevos repos

### Paso 2: Modificar CalculateFinalSalePrice — Client Override
- Antes de aplicar SaleModificationRules, para cada item:
  - Llamar `clientPricingRepo.FindApplicable(ctx, clientID, variantID, saleDate)`
  - Si retorna override → usar `override.FixedPrice` como finalPrice
  - Si no → continuar flujo actual (SaleModificationRules + fallback)

### Paso 3: Modificar CalculateFinalSalePrice — Audit Trail
- Al final del cálculo, crear `domain.NewPriceCalculation(...)` con las reglas aplicadas
- Guardar con `calculationRepo.Save(ctx, calc)`
- No bloquear el flujo si falla el guardado (log + continuar)

### Paso 4: Actualizar wiring en main.go
- Pasar `clientPricingRepo` y `calculationRepo` al constructor

### Paso 5: Tests unitarios
- Test: ClientPricing override encontrado → usa precio fijo
- Test: ClientPricing override no encontrado → flujo normal con SaleModificationRules
- Test: ClientPricing override expirado → flujo normal
- Test: PriceCalculation se genera correctamente
- Test: Fallo al guardar PriceCalculation no rompe el cálculo

---

## 📝 NOTAS

- La entidad `ClientPricing` y su repositorio `ClientPricingRepository` ya existen en el dominio
- La tabla `client_pricing_overrides` ya existe en la migración 004
- La entidad `PriceCalculation` y `PriceCalculationRepository` ya existen
- No se requiere migración de base de datos
- El `CalculateFinalSalePriceRequest` ya incluye `ClientID string` — todo está preparado
