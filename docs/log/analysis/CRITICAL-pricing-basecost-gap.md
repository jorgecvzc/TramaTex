# 🔴 CRÍTICO: Gap en la Gestión de Precios Base (BaseCost)

**Fecha:** 22 de febrero de 2026  
**Severidad:** CRÍTICA  
**Módulos afectados:** Product, Pricing  
**Estado:** Detectado - Pendiente de corrección

---

## 1. Resumen Ejecutivo

El sistema tiene un **gap crítico** en el flujo de gestión de precios de productos. Aunque el campo `BaseCost` existe en la entidad `ProductVariant` y se persiste en la base de datos, **NO HAY forma de ingresarlo durante el flujo principal de creación de productos**. Esto resulta en variantes con `BaseCost = 0` por defecto, lo que produce cálculos de precios incorrectos en el módulo Pricing.

---

## 2. Hallazgos Técnicos

### 2.1. Campo BaseCost en el Backend ✅

**Archivo:** `apps/tramatex-api/internal/product/domain/variant.go` (línea 17)

```go
type ProductVariant struct {
	ID              uuid.UUID
	ProductID       uuid.UUID
	SKU             string
	Barcode         *string
	BaseCost        float64  // ← EXISTE
	Status          VariantStatus
	AttributeValues []uuid.UUID
	IsActive        bool
}
```

**Persistencia:** `variant_data_model.go` (línea 18)
```go
BaseCost  float64  `gorm:"type:numeric(12,2);not null;default:0"`
```

**Schema SQL:**
```sql
"base_cost" NUMERIC(12,2) NOT NULL DEFAULT 0
```

✅ **El campo existe y se persiste correctamente**

---

### 2.2. Valor por Defecto = 0 ⚠️

**Archivo:** `variant.go` (línea 58)

```go
func NewProductVariant(...) (*ProductVariant, error) {
	return &ProductVariant{
		ID:              uuid.New(),
		ProductID:       productID,
		SKU:             sku,
		Barcode:         barcode,
		BaseCost:        0,  // ← DEFAULT ES CERO
		Status:          status,
		AttributeValues: attributeValueIDs,
		IsActive:        true,
	}, nil
}
```

⚠️ **Todas las variantes se crean con BaseCost = 0 si no se especifica**

---

### 2.3. UI - Formulario Principal de Productos ❌

**Archivo:** `ProductFormBasic.vue` (Step 1 del wizard)

**Campos disponibles:**
- ✅ Product Type (TANGIBLE/SERVICE)
- ✅ SKU
- ✅ Name / Long Name
- ✅ Description

**Campos faltantes:**
- ❌ **BaseCost** (costo base del producto)
- ❌ **Base Price** (precio base antes de variantes)

**Resultado:** NO hay forma de ingresar el costo/precio base durante la creación del producto.

---

### 2.4. UI - Modal de Variantes (Workaround parcial) ⚠️

**Archivo:** `VariantFormModal.vue` (línea 80-90)

```vue
<!-- Base Cost -->
<div class="form-group">
  <label for="baseCost">Costo base</label>
  <input
    id="baseCost"
    v-model.number="form.baseCost"
    type="number"
    step="0.01"
    min="0"
    class="form-control"
    placeholder="0.00"
  />
  <small class="form-text">Costo unitario (opcional)</small>
</div>
```

✅ **El campo SÍ está disponible en el modal de variantes**  
⚠️ **PERO** este modal es para creación manual de variantes después de crear el producto  
❌ **NO se usa en el flujo principal de creación**  
❌ **Las variantes JIT (Just-in-Time) se crean con BaseCost = 0**

---

### 2.5. Atributos NO tienen Modificadores de Precio ❌

**Archivo:** `attribute.go` (línea 17-22)

```go
type AttributeValue struct {
	ID          uuid.UUID
	AttributeID uuid.UUID
	Value       string
	Code        string
	// NO HAY campo PriceModifier
}
```

**Búsqueda en todo el dominio:**
```
grep -r "PriceModifier|priceModifier|modifier" apps/tramatex-api/internal/product/domain/
# Sin resultados
```

❌ **Los atributos NO pueden modificar el precio base**  
❌ **Contradicción con la documentación de Pricing**

---

### 2.6. Documentación vs Implementación 📄❌

#### **Documentación de Pricing** (`pricing-domain.md` líneas 11-14):

> *   **Precio Base de Producto:** El costo o precio de referencia de un producto antes de cualquier modificación de variante, margen o descuento. Se obtiene del módulo `Product`.
> *   **Modificador de Variante:** Un factor (positivo o negativo) que ajusta el `Precio Base de Producto` según las características de una variante específica (ej., talla, color, material). Se obtiene del módulo `Product`.

#### **Realidad en el Código:**
- ✅ `BaseCost` existe en `ProductVariant`
- ❌ **NO hay "Precio Base de Producto"** (nivel superior a variante)
- ❌ **NO hay "Modificador de Variante"** en los atributos
- ❌ El `BaseCost` de cada variante debe ingresarse manualmente

**Conclusión:** La documentación describe un modelo que **NO está implementado**.

---

### 2.7. Documentación del Producto ❌

**Archivo:** `product/domain-model.md` (línea 89-98)

```markdown
### 1.3. `ProductVariant` (Entity)

- **ProductVariant:**
  - `VariantID` (ID único)
  - `ProductID` (FK a Product)
  - `SKU` (SKU compuesto, ej: `FYR2040-T.L-C.R`)
  - `Barcode` (Opcional, EAN/UPC, para la variante específica)
  - `AttributeValues` (Lista de `AttributeValueID`): Describe la combinación específica.
  - `Status` (Enum: `PROVISIONAL`, `CONFIRMED`)
  - `IsActive` (Boolean)
```

❌ **El campo `BaseCost` NO está documentado**  
❌ La documentación del dominio está desactualizada

---

## 3. Impacto en el Módulo Pricing

### 3.1. Flujo de Cálculo Actual

**Archivo:** `pricing_engine_service.go` (línea 357-367)

```go
func (s *PricingEngineService) calculateBaseSalesPriceFromInfo(
    ctx context.Context, 
    variantID uuid.UUID, 
    info *ProductPricingInfo
) (domain.Money, error) {
	baseCost, err := domain.NewMoney(info.BaseCost, info.Currency)
	if err != nil {
		return domain.Money{}, err
	}

	rules, err := s.baseRuleRepo.List(ctx)
	if err != nil {
		return domain.Money{}, err
	}

	selected := selectBaseRule(rules, info.BrandID, info.GroupIDs, info.ProductID, variantID)
	baseSalesPrice := baseCost  // ← PARTE DE BaseCost
	if selected != nil {
		baseSalesPrice, err = selected.Value.Apply(baseCost)  // ← APLICA REGLA
		if err != nil {
			return domain.Money{}, err
		}
	}

	return baseSalesPrice, nil
}
```

**Fórmula:**
```
BaseSalesPrice = BaseCost + (aplicación de BaseSalesPriceRule)
```

**Ejemplo con BaseCost = 0:**
```
BaseCost = 0 EUR  (valor por defecto)
Regla: +30% markup
BaseSalesPrice = 0 * 1.30 = 0 EUR  ← INCORRECTO
```

**Ejemplo con BaseCost = 50 EUR:**
```
BaseCost = 50 EUR  (ingresado manualmente)
Regla: +30% markup
BaseSalesPrice = 50 * 1.30 = 65 EUR  ← CORRECTO
```

### 3.2. ¿Se Almacena el Precio Calculado?

**Respuesta: NO directamente en la DB**

El `BaseSalesPrice` calculado:
- ✅ Se almacena en **caché Redis** (por ProductID + VariantID)
- ❌ **NO se persiste** en la tabla `product_variants`
- 🔄 Se **recalcula dinámicamente** cuando:
  - La caché expira (TTL)
  - Se invalida manualmente
  - Se crea/modifica una variante
  - Se crea/modifica una `BaseSalesPriceRule`

**Estrategia de caché:**
```go
// Se cachea por producto completo (todas sus variantes)
cache.SetBasePrice(ctx, productID, variantID, baseSalesPrice)

// Se invalida cuando:
- Nueva variante creada o modificada
- Regla de pricing creada/modificada
- Manualmente vía API
```

---

## 4. Flujos de Trabajo Afectados

### 4.1. Flujo Actual (PROBLEMÁTICO)

```
1. Usuario crea producto en ProductFormBasic
   └─> SKU, Name, Brand, Groups
   └─> ❌ NO ingresa BaseCost

2. Usuario configura atributos en ProductFormAttributes
   └─> Selecciona Talla, Color, etc.
   └─> ❌ Los atributos NO tienen modificadores de precio

3. Sistema crea producto
   └─> Product guardado
   └─> ⚠️ NO se crean variantes aún (JIT)

4. Escenario A: Creación manual de variantes
   a) Usuario abre tab "Variantes"
   b) Click "Añadir Variante"
   c) VariantFormModal se abre
   d) ✅ Usuario PUEDE ingresar BaseCost aquí
   e) Variante guardada con BaseCost especificado
   
5. Escenario B: Creación JIT (en ventas)
   a) Usuario crea orden de venta
   b) Añade producto con configuración específica
   c) Sistema busca variante → NO existe
   d) Sistema crea variante automáticamente
   e) ❌ BaseCost = 0 (default)
   f) Pricing calcula con BaseCost = 0 → PRECIO INCORRECTO

6. Pricing calcula BaseSalesPrice
   └─> baseCost (0 o manual) + regla → resultado
```

**Problema principal:** Las variantes JIT se crean con `BaseCost = 0`, resultando en precios incorrectos.

### 4.2. Flujo Esperado (según especificación)

```
1. Usuario crea producto
   └─> Ingresa SKU, Name, Brand, Groups
   └─> Ingresa "Precio Base del Producto" (ej: 45 EUR) ← FALTA

2. Usuario configura atributos CON modificadores
   └─> Talla S: +0 EUR
   └─> Talla M: +0 EUR
   └─> Talla L: +2 EUR  ← FALTA
   └─> Talla XL: +5 EUR ← FALTA
   └─> Color Rojo: +0 EUR
   └─> Color Azul: +1 EUR ← FALTA

3. Sistema calcula BaseCost por variante automáticamente:
   └─> Variante (L, Rojo):  45 + 2 + 0 = 47 EUR
   └─> Variante (XL, Azul): 45 + 5 + 1 = 51 EUR
   └─> ✅ Todas las variantes tienen BaseCost correcto

4. Variantes JIT se crean con BaseCost calculado:
   └─> Al añadir a venta, sistema calcula BaseCost antes de crear variante
   └─> ✅ Variante creada con BaseCost correcto

5. Pricing calcula BaseSalesPrice:
   └─> baseCost (calculado) + margen/regla → ✅ CORRECTO
```

---

## 5. Comparación: Implementación vs Diseño

| Concepto | Diseño (Documentación) | Implementación Actual | Estado |
|----------|------------------------|----------------------|--------|
| Precio Base del Producto | Campo a nivel Product | ❌ NO existe | 🔴 GAP |
| BaseCost de Variante | Calculado automáticamente | ✅ Existe pero default=0 | ⚠️ PARCIAL |
| Modificador de Atributo | AttributeValue tiene precio extra | ❌ NO existe | 🔴 GAP |
| UI para BaseCost principal | En formulario de creación | ❌ NO existe | 🔴 GAP |
| UI para BaseCost variante | En modal de variantes | ✅ Existe | ✅ OK |
| Cálculo automático de BaseCost | Precio base + Σ modificadores | ❌ NO implementado | 🔴 GAP |
| Persistencia de BaseSalesPrice | Solo en caché (Redis) | ✅ Implementado | ✅ OK |
| Invalidación de caché | Al crear/modificar variante | ✅ Implementado | ✅ OK |

---

## 6. Impacto del Negocio

### 6.1. Escenarios Críticos

#### **Escenario 1: Variantes JIT en Ventas**
```
Situación: Cliente ordena "Camiseta Nike, Talla XL, Color Rojo"
- Sistema busca variante → NO existe
- Sistema crea variante JIT con BaseCost = 0
- Pricing calcula: 0 * 1.30 (margen 30%) = 0 EUR
- Cliente ve precio = 0 EUR
❌ FALLO CRÍTICO
```

#### **Escenario 2: Productos con Muchas Variantes**
```
Producto: Camiseta con 4 tallas × 5 colores = 20 variantes
- Usuario debe crear manualmente cada variante
- Para cada una, ingresar BaseCost en el modal
- ❌ Proceso manual y propenso a errores
- ⚠️ No escalable
```

#### **Escenario 3: Cambio de Costos**
```
Situación: El proveedor aumenta costos de productos XL
- Usuario debe actualizar manualmente cada variante XL
- ❌ NO hay forma de decir "Talla XL siempre cuesta +5 EUR"
- ⚠️ Mantenimiento ineficiente
```

### 6.2. Consecuencias Operacionales

1. **Imposibilidad de usar JIT para pricing correcto**
   - Las variantes JIT siempre tendrán precio incorrecto
   - Obliga a pre-crear todas las variantes manualmente
   - Pierde el beneficio principal de la estrategia JIT

2. **Trabajo manual excesivo**
   - Cada variante requiere ingreso manual de BaseCost
   - Para productos con muchas combinaciones: inviable
   - Ejemplo: 5 tallas × 6 colores = 30 ingresos manuales

3. **Inconsistencias de datos**
   - Fácil olvidar ingresar BaseCost en algunas variantes
   - Difícil mantener coheren cia en modificadores (ej: "XL siempre +5 EUR")
   - Errores humanos en cálculo de costos

4. **Experiencia de usuario pobre**
   - Flujo de creación de productos incompleto
   - Pasos ocultos (debe ir al tab Variantes y añadir después)
   - No intuitivo

---

## 7. Soluciones Propuestas

### 7.1. Solución Mínima (Corto Plazo)

**Objetivo:** Permitir ingreso de BaseCost en el flujo principal de creación.

#### **Opción A: BaseCost a Nivel Producto**
```
Añadir campo "Precio/Costo Base" en ProductFormBasic:
- Se almacena en Product (nuevo campo)
- Al crear variantes JIT, se usa este valor por defecto
- Variantes individuales pueden sobrescribirlo en el modal
```

**Cambios requeridos:**
1. `product/domain/product.go`: Añadir `BasePrice float64`
2. Schema: `ALTER TABLE products ADD COLUMN base_price NUMERIC(12,2) DEFAULT 0`
3. `ProductFormBasic.vue`: Añadir input para `base_price`
4. Lógica de variantes: Usar `product.BasePrice` como default de `variant.BaseCost`

**Pros:**
- ✅ Mínimo cambio estructural
- ✅ Compatibilidad con código existente
- ✅ Variantes JIT funcionan correctamente

**Contras:**
- ⚠️ NO resuelve modificadores por atributo
- ⚠️ Todas las variantes comparten el mismo costo base

---

#### **Opción B: BaseCost Default + Override**
```
Combinar:
1. Campo "Costo Base Default" a nivel producto (ProductFormBasic)
2. Mantener BaseCost editable por variante (VariantFormModal)
3. JIT usa el default, usuario puede override después
```

**Cambios requeridos:**
1. Igual que Opción A
2. UI clara indicando que es "default"
3. Documentación del flujo

**Pros:**
- ✅ Flexibilidad máxima
- ✅ Soluciona JIT inmediatamente
- ✅ Permite casos especiales

**Contras:**
- ⚠️ NO resuelve modificadores automáticos
- ⚠️ Sigue requiriendo trabajo manual para variantes con costos diferentes

---

### 7.2. Solución Completa (Medio Plazo)

**Objetivo:** Implementar modificadores de precio por atributo según diseño original.

#### **Cambios en AttributeValue**
```go
type AttributeValue struct {
	ID            uuid.UUID
	AttributeID   uuid.UUID
	Value         string
	Code          string
	PriceModifier *PriceModifier  // NUEVO
}

type PriceModifier struct {
	Type   ModifierType  // FIXED, PERCENTAGE
	Amount float64
}
```

#### **Flujo de Cálculo Automático**
```
1. Product tiene BasePrice: 45 EUR
2. Atributos con modificadores:
   - Talla.L: +2 EUR (FIXED)
   - Color.Azul: +1 EUR (FIXED)
3. Al crear variante (L, Azul):
   BaseCost = 45 + 2 + 1 = 48 EUR
4. Variante se crea/guarda con BaseCost = 48 EUR
```

#### **Cambios requeridos:**
1. **Backend:**
   - `attribute.go`: Añadir `PriceModifier` a `AttributeValue`
   - Schema: Nueva tabla `attribute_value_price_modifiers`
   - `variant.go`: Método `CalculateBaseCostFromModifiers()`
   - Lógica JIT: Calcular BaseCost antes de crear variante

2. **Frontend:**
   - `ProductFormBasic.vue`: Campo "Precio Base"
   - UI de Atributos: Checkbox "Tiene modificador de precio"
   - Input para modificador (tipo + cantidad)
   - ProductFormPreview: Mostrar cálculo de costos por variante

3. **Migraciones:**
   - Añadir columna `products.base_price`
   - Crear tabla `attribute_value_price_modifiers`
   - Migración de datos existentes

**Pros:**
- ✅ Alineado con diseño original
- ✅ Escalable para muchas variantes
- ✅ Mantenimiento eficiente
- ✅ Cálculo automático de costos

**Contras:**
- ⚠️ Cambio estructural significativo
- ⚠️ Requiere más tiempo de desarrollo
- ⚠️ Migración de datos existentes compleja

---

### 7.3. Recomendación Inmediata

**Implementar Opción A (BaseCost a nivel Product) AHORA:**

**Razones:**
1. **Urgencia:** Las variantes JIT actualmente producen precios incorrectos
2. **Impacto mínimo:** Cambio pequeño, bajo riesgo
3. **Valor inmediato:** Desbloquea el flujo de trabajo principal
4. **Compatible:** No bloquea la Solución Completa posterior

**Plan de acción:**
1. Sprint actual: Implementar Opción A (3-5 días)
2. Sprint +1: Validación y testing exhaustivo
3. Sprint +2: Comenzar Solución Completa (modificadores)

---

## 8. Acciones Requeridas

### 8.1. Inmediatas (Esta Semana)

- [ ] **CRÍTICO:** Implementar campo `base_price` en Product
- [ ] Añadir input en `ProductFormBasic.vue`
- [ ] Modificar lógica de creación de variantes JIT
- [ ] Actualizar documentación (`product/domain-model.md`)
- [ ] Testing exhaustivo del flujo completo

### 8.2. Corto Plazo (Próximo Sprint)

- [ ] Diseñar UI para modificadores de atributos
- [ ] Especificar cálculo automático de BaseCost
- [ ] ADR para la implementación de modificadores
- [ ] Plan de migración de datos

### 8.3. Documentación

- [ ] Actualizar `product/domain-model.md` con `BaseCost`
- [ ] Corregir `pricing-domain.md` (eliminar referencias a modificadores no implementados)
- [ ] Crear guía de usuario para gestión de precios
- [ ] Documentar casos de uso de pricing

---

## 9. Riesgos de No Actuar

1. **Datos incorrectos en producción**
   - Variantes con BaseCost = 0
   - Precios de venta calculados incorrectamente
   - Órdenes de venta con totales erróneos

2. **Imposibilidad de usar funcionalidad JIT**
   - Obligado a pre-crear todas las variantes
   - Pérdida del beneficio arquitectónico principal

3. **Pérdida de confianza del usuario**
   - Sistema que "no funciona como esperado"
   - Datos poco confiables
   - Trabajo manual excesivo

4. **Deuda técnica acumulada**
   - Workarounds en lugar de solución real
   - Código inconsistente con documentación
   - Dificulta mantenimiento futuro

---

## 10. Referencias

- `apps/tramatex-api/internal/product/domain/variant.go`
- `apps/tramatex-api/internal/pricing/application/pricing_engine_service.go`
- `apps/frontend/src/components/product/ProductFormBasic.vue`
- `apps/frontend/src/components/product/VariantFormModal.vue`
- `docs/modules/pricing/pricing-domain.md`
- `docs/modules/product/domain-model.md`
- `docs/architecture/adrs/adr-016-pricing-module-architecture.md`

---

**Autor:** Sistema de IA (Análisis automatizado)  
**Revisión requerida:** Product Owner, Tech Lead  
**Prioridad:** 🔴 CRÍTICA
