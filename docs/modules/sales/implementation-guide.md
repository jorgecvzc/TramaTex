# Guía de Implementación - Módulo Sales

Este documento contiene detalles sobre la implementación técnica del módulo Sales, su estructura de directorios, dependencias y flujo de trabajo sugerido.

---

## Integración con Módulo Product

### Obtención del BaseCost de Variantes

**⚠️ CRÍTICO:** El `baseCost` de una `ProductVariant` se calcula **dinámicamente** y **NO se almacena** en la base de datos.

#### Cálculo del BaseCost

```
baseCost = Product.BasePrice + Σ(AttributeValue modifiers)
```

Los modificadores de atributos se aplican secuencialmente según `Attribute.sortOrder`:
- **FIXED**: Suma/resta cantidad fija (€) → `baseCost += modifierAmount`
- **PERCENTAGE**: Aplica porcentaje sobre precio acumulado → `baseCost += baseCost × (modifierAmount / 100)`

#### Ejemplo de Cálculo

```
Producto: basePrice = 100.00€

Atributos aplicados a la variante:
1. Talla "L" (FIXED): +10.00€ → precio = 110.00€
2. Acabado "Premium" (PERCENTAGE): +15% → precio = 110.00 × 1.15 = 126.50€

baseCost final = 126.50€
```

#### Implementación en Sales

Cuando se crea una línea de pedido/presupuesto:

```go
// 1. Obtener la variante completa desde Product API
variant, err := productService.GetProductVariantByID(ctx, lineItem.ProductVariantID)
if err != nil {
    return err
}

// 2. El campo variant.BaseCost contiene el valor calculado dinámicamente
// Este es el precio base antes de aplicar markup de marca y descuentos del pricing module
baseCost := variant.BaseCost

// 3. Pasar baseCost al módulo Pricing para obtener el precio de venta
calculatedPrice, err := pricingService.CalculatePrice(ctx, PriceCalculationParams{
    BaseCost: baseCost,
    BrandID: product.BrandID,
    PartyID: order.PartyID,
    // ... otros parámetros
})

// 4. Crear el line item con el precio calculado
lineItem.ListUnitPrice = calculatedPrice  // precio de tarifa
lineItem.UnitPrice = calculatedPrice      // precio de venta (puede ser override)
```

#### Consideraciones Importantes

1. **Siempre obtener variantes frescas**: Nunca cachear el `baseCost` por períodos largos, ya que puede cambiar si se modifican:
   - El `basePrice` del producto padre
   - Los valores de `modifierAmount` en los atributos

2. **Performance**: Para operaciones batch (ej: listar 100 pedidos con sus líneas), considerar:
   - Endpoint bulk en Product API: `POST /products/variants/batch` que devuelva múltiples variantes con baseCost calculado
   - Cachear por la duración de la transacción/request

3. **Validación de precios**: Al convertir un `Quote` a `SalesOrder`, re-calcular los precios para asegurar que reflejan los valores actuales del sistema.

---

## Estructura de Directorios

_Por definir según implementación..._
