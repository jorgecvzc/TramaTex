# Tarea 17-03: Eliminar Sistema Legacy de Pricing

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 03 |
| **ID de Sprint** | sprint-17 |
| **Título** | Eliminar PricingService, PricingHandler y entidades obsoletas |
| **Estado** | ⏳ Planificado |
| **Facilitador/LLM** | GitHub Copilot / Claude Opus 4.6 |
| **Fecha de Inicio** | pendiente |
| **Fecha de Fin** | pendiente |
| **Duración Estimada** | 2-3 horas |
| **Prioridad** | P1 — Alto |
| **Dependencia** | Requiere 17-01 y 17-02 completadas |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [ ] **Eliminar `PricingService`** y todo código que lo referencia
2. [ ] **Eliminar `PricingHandler`** y sus rutas
3. [ ] **Eliminar entidades domain obsoletas**: `PricingRule`, `BrandProfitMargin`, `SalesDiscountRule`, servicios de dominio asociados (`SellingPriceCalculatorService`, `SalesDiscountCalculatorService`)
4. [ ] **Eliminar repositorios e infra obsoleta**: GORM repos + data models de las entidades viejas
5. [ ] **Eliminar DTOs obsoletos**: `dtos.go`, `commands.go`, `queries.go` del sistema viejo
6. [ ] **Limpiar wiring** en `main.go`
7. [ ] **Evaluar migración SQL** para tablas `pricing_rules` y `sales_discount_rules`

---

## 📊 CONTEXTO

### Archivos a eliminar (estimación)

**Domain:**
- `pricing_rule.go` — entidad PricingRule
- `brand_profit_margin.go` — entidad BrandProfitMargin
- `sales_discount_rule.go` — entidad SalesDiscountRule
- `services.go` — SellingPriceCalculatorService + SalesDiscountCalculatorService

**Application:**
- `pricing_service.go` + `pricing_service_test.go`
- `dtos.go` (parte vieja), `commands.go`, `queries.go`

**Infrastructure:**
- Repositorios GORM: PricingRule, BrandProfitMargin, SalesDiscountRule
- Data models correspondientes

**Interfaces:**
- `pricing_handler.go` + `pricing_handler_test.go`

### Archivos a conservar

- `client_pricing.go` — usado por el motor nuevo (17-01)
- `price_calculation.go` — usado por el motor nuevo (17-01)
- `money.go`, `percentage.go`, `rule_value.go` — Value Objects ADR-016
- `base_sales_price_rule.go`, `sale_modification_rule.go` — entidades ADR-016
- `pricing_engine_service.go`, `pricing_engine_handler.go` — motor canónico
- Todos los repos e infra de las entidades conservadas

### Repositorios en `repository.go` — análisis

| Interfaz | Mantener | Motivo |
|----------|----------|--------|
| `BaseSalesPriceRuleRepository` | ✅ | ADR-016 |
| `SaleModificationRuleRepository` | ✅ | ADR-016 |
| `ClientPricingRepository` | ✅ | Usado por motor nuevo (17-01) |
| `PriceCalculationRepository` | ✅ | Audit trail (17-01) |
| `PricingRuleRepository` | ❌ | Entidad vieja |
| `BrandProfitMarginRepository` | ❌ | Entidad vieja — margen de marca se obtiene de ProductClient |
| `SalesDiscountRuleRepository` | ❌ | Reemplazada por SaleModificationRule |

### Tablas en migración 004

| Tabla | Mantener | Motivo |
|-------|----------|--------|
| `rule_value_types` | ✅ | ADR-016 |
| `base_sales_price_rules` | ✅ | ADR-016 |
| `sale_modification_rules` | ✅ | ADR-016 |
| `client_pricing_overrides` | ✅ | Overrides |
| `price_calculations` | ✅ | Audit trail |
| `brand_profit_margins` | ❌ | Evaluar: archivar datos → DROP |
| `pricing_rules` | ❌ | Evaluar: archivar datos → DROP |
| `sales_discount_rules` | ❌ | Evaluar: archivar datos → DROP |

---

## 🛠️ PLAN DE TRABAJO

### Paso 1: Verificar que ningún módulo externo referencia entidades viejas
- `grep` PricingRule, BrandProfitMargin, SalesDiscountRule fuera de `/pricing/`
- Confirmar que Sales solo usa la interfaz `PricingEngine`

### Paso 2: Eliminar archivos (bottom-up)
- Handlers → Application → Infrastructure → Domain

### Paso 3: Limpiar `repository.go`
- Eliminar interfaces obsoletas

### Paso 4: Limpiar `main.go`
- Eliminar instanciación de repos y servicios viejos

### Paso 5: Migración SQL (nueva migración)
- Crear `XXX_drop_legacy_pricing_tables.sql`
- DROP tablas obsoletas (con backup previo si hay datos)

### Paso 6: Compilar y ejecutar tests
- `go build ./...` limpio
- `go test ./internal/pricing/...` verde
- Verificar que Sales compila y tests pasan
