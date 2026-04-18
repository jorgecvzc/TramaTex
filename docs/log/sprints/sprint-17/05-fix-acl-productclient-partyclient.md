# Tarea 17-05: Fix ACL — ProductClient y PartyClient

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 05 |
| **ID de Sprint** | sprint-17 |
| **Título** | Corregir Anti-Corruption Layer: ProductClient y PartyClient |
| **Estado** | ✅ Completado |
| **Facilitador/LLM** | GitHub Copilot / Claude Opus 4.6 |
| **Fecha de Inicio** | 2026-04-11 |
| **Fecha de Fin** | 2026-04-12 |
| **Duración Estimada** | 2-3 horas |
| **Prioridad** | P2 — Medio |
| **Dependencia** | Independiente (paralelizable) |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [x] **V1 — ProductClient**: Eliminar imports directos de `product/domain` y `product/infrastructure/persistence`. Consumir capa de aplicación de Product.
2. [x] **V2 — PartyClient**: Eliminar query SQL directa sobre tabla `parties`. Consumir capa de aplicación de Party.

---

## 📊 CONTEXTO

### V1 — ProductClient rompe ACL (import cruzado)

**Archivo:** `internal/pricing/infrastructure/productclient/product_client.go`

```go
import productdomain "github.com/joran-cortez/tramatex/internal/product/domain"
import productpersistence "github.com/joran-cortez/tramatex/internal/product/infrastructure/persistence"
```

**Violación:** Un adaptador ACL de Pricing importa directamente el dominio Y la infraestructura de Product. Si Product cambia internamente, Pricing se rompe.

**Solución:** El módulo Product debe exponer una interfaz de aplicación que Pricing pueda consumir. Opciones:
- A) Product expone un `ProductPricingInfoProvider` en su capa de aplicación
- B) Pricing define la interfaz (ya la tiene: `ProductInfoProvider`) y Product implementa un adaptador

### V2 — PartyClient accede directo a tabla ajena

**Archivo:** `internal/pricing/infrastructure/partyclient/party_client.go`

```go
c.db.Table("parties").Select("default_discount_percentage").Where("id = ?", clientID)
```

**Violación:** Acoplamiento al esquema de la tabla `parties` del módulo Party.

**Solución:** Party debe exponer un método en su capa de aplicación para obtener el descuento por defecto de un cliente. Pricing lo consume vía interfaz `ClientInfoProvider`.

---

## 🛠️ PLAN DE TRABAJO

### V1 — ProductClient

#### Paso 1: Verificar qué expone Product en su application layer
- ¿Existe ya un servicio que devuelva BaseCost + TaxRate + BrandMarkup + GroupIDs?

#### Paso 2: Si no existe, crear método en ProductService
- `GetVariantPricingInfo(ctx, variantID) → ProductPricingInfo`
- Encapsula: variant lookup + product lookup + brand markup + attribute modifiers

#### Paso 3: Refactorizar ProductClient
- Eliminar imports de `product/domain` y `product/persistence`
- Inyectar interfaz del servicio de Product (o llamar directamente)
- Mantener misma firma de `ProductInfoProvider`

### V2 — PartyClient

#### Paso 1: Verificar qué expone Party en su application layer
- ¿Existe método para obtener DefaultDiscountPercentage?

#### Paso 2: Si no existe, crear método en PartyService
- `GetClientDefaultDiscount(ctx, partyID) → float64`

#### Paso 3: Refactorizar PartyClient
- Eliminar SQL directo
- Inyectar interfaz del servicio de Party
- Mantener misma firma de `ClientInfoProvider`

### Paso 4: Actualizar wiring en main.go

### Paso 5: Tests
- Verificar que Pricing sigue funcionando con los nuevos adaptadores
- Tests unitarios con mocks de las interfaces de Product/Party
