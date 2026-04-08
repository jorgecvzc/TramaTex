# Tarea 17-04: Migración Money float64 → decimal.Decimal

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 04 |
| **ID de Sprint** | sprint-17 |
| **Título** | Migrar Value Object Money de float64 a decimal.Decimal |
| **Estado** | ⏳ Planificado |
| **Facilitador/LLM** | GitHub Copilot / Claude Opus 4.6 |
| **Fecha de Inicio** | pendiente |
| **Fecha de Fin** | pendiente |
| **Duración Estimada** | 3-4 horas |
| **Prioridad** | P1 — Alto |
| **Dependencia** | Preferible tras 17-03 (menos código que migrar) |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [ ] **Reemplazar `float64` por `decimal.Decimal`** en `Money` Value Object
2. [ ] **Propagar cambio** por domain → application → infrastructure → interfaces
3. [ ] **Actualizar operaciones**: Add, Subtract, Multiply, comparaciones
4. [ ] **Actualizar persistencia**: mappers GORM (decimal ↔ float64/numeric)
5. [ ] **Actualizar tests** existentes al nuevo tipo

---

## 📊 CONTEXTO

### Motivación

- `architecture.yaml` documenta explícitamente: `Money struct { amount decimal.Decimal, currency string }`
- El código actual usa `float64` con `roundTo2Decimals()` (round-half-up)
- Para módulo de criticidad económica (≥85% cobertura, STRICT), float64 puede acumular errores de representación en cadenas de operaciones
- Biblioteca recomendada: `github.com/shopspring/decimal`

### Impacto estimado

| Capa | Archivos afectados |
|------|-------------------|
| Domain | `money.go`, `percentage.go`, `rule_value.go`, `price_calculation.go`, `client_pricing.go`, services |
| Application | `pricing_engine_service.go`, DTOs (conversión a float64 para JSON) |
| Infrastructure | Data models (GORM: `numeric` ↔ `decimal.Decimal`), ProductClient, PartyClient |
| Interfaces | Handlers (sin cambio directo, DTOs manejan serialización) |
| Tests | Todos los tests de pricing domain + application |

### Decisiones de diseño

- **DTOs exponen `float64`** (para serialización JSON) — la conversión se hace en el mapper
- **Domain usa `decimal.Decimal`** internamente — toda aritmética precisa
- **Persistencia**: columnas `NUMERIC(12,2)` en PostgreSQL (ya compatible)
- **No se requiere migración SQL** — las columnas ya son numéricas

---

## 🛠️ PLAN DE TRABAJO

### Paso 1: Añadir dependencia
- `go get github.com/shopspring/decimal`

### Paso 2: Migrar Money VO
- Cambiar campo interno `amount float64` → `amount decimal.Decimal`
- Actualizar constructores, getters, operaciones
- Mantener `Amount() float64` como método de conveniencia (deprecated) + añadir `Decimal() decimal.Decimal`

### Paso 3: Migrar Percentage y RuleValue
- Percentage: evaluar si internamente conviene decimal
- RuleValue.Apply(): actualizar aritmética

### Paso 4: Actualizar application layer
- DTOs: conversión decimal → float64 para respuesta JSON
- Commands: aceptar float64 en entrada → convertir a decimal

### Paso 5: Actualizar infrastructure
- Data models: mappers GORM con `decimal.Decimal`
- ProductClient: `BaseCost float64` → convertir al crear Money

### Paso 6: Actualizar tests
- Reemplazar comparaciones float64 por decimal
- Verificar precision en operaciones encadenadas

### Paso 7: Build + test completo
- `go build ./...`
- `go test ./internal/pricing/...`
- `go test ./internal/sales/...` (consumidor)
