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

## Arquitectura del Módulo

Tras el refactor del Sprint 14, el módulo `Sales` se ha organizado siguiendo el patrón de **Arquitectura Hexagonal (Puertos y Adaptadores)** con una fragmentación clara de responsabilidades en la capa de aplicación para evitar el crecimiento desmedido del servicio principal.

### Fragmentación de Servicios de Aplicación

El anterior `SalesService` monolítico se ha dividido en servicios especializados:

- **`QuoteService`**: Gestión completa del ciclo de vida de presupuestos (`Quotes`). Maneja la creación, actualización de estados y la lógica de conversión automática a pedido.
- **`OrderService`**: Gestión de pedidos (`SalesOrders`). Se encarga de la persistencia de pedidos, seguimiento de estados y coordinación con el módulo MES para la visibilidad de producción.
- **`DeliveryNoteService`**: Gestión de albaranes. Controla la creación de albaranes parciales o totales vinculados a pedidos y la actualización de cantidades entregadas.
- **`BillingService`**: Orquestador de facturación. Maneja la creación de facturas (`Invoices`) a partir de albaranes (B2B) y la generación de tickets simplificados (B2C).
- **`InvoiceService`** (Post-MVP): Servicio especializado para la gestión avanzada de facturas, series y comunicación con la AEAT (Ver ADR-020).

### Capa de Dominio y Cálculos

La lógica de cálculo de totales, impuestos y subtotales se ha centralizado en `calculations.go` dentro de la capa de dominio. Esto asegura que tanto presupuestos como pedidos y facturas utilicen las mismas reglas de negocio consistentes.

- **`SalesTotals`**: Estructura de dominio que encapsula la lógica de agregación de líneas.
- **Redondeo**: Se aplica redondeo a 2 decimales en cada línea y en el total final para consistencia contable.

---

## Estructura de Directorios

La estructura sigue el estándar del proyecto para módulos Go:

```text
internal/sales/
├── application/          # Servicios de aplicación fragmentados
│   ├── billing_service.go
│   ├── order_service.go
│   ├── quote_service.go
│   └── delivery_note_service.go
├── domain/               # Entidades de dominio, value objects y lógica de cálculo
│   ├── model/            # Entidades (Quote, SalesOrder, Invoice, etc.)
│   ├── repository/       # Interfaces de repositorio
│   └── calculations.go   # Lógica compartida de cálculos comerciales
├── infrastructure/       # Implementaciones técnicas
│   └── persistence/      # Repositorios SQL (PostgreSQL)
└── interfaces/           # Adaptadores de entrada (HTTP/REST)
    └── handlers/         # Manejadores de rutas Gin

---

### Gestión de Errores Estandarizada

El módulo `Sales` delega la traducción de errores de dominio a respuestas HTTP en el `ErrorHandlerMiddleware` de la capa `shared`. 

Para que esto funcione:
1. **Definir Errores en Dominio**: Todos los errores de negocio se definen en `internal/sales/domain/errors.go`.
2. **Implementar `HTTPStatuser`**: Los errores de dominio implementan la interfaz `shared/domain.HTTPStatuser` para indicar su código HTTP correspondiente (ej. `ErrQuoteNotFound` devuelve `404`).
3. **Delegación en Handlers**: Los controladores Gin NO deben formatear respuestas de error manualmente. Deben simplemente adjuntar el error al contexto: `c.Error(err)`. El middleware se encargará de sanitizar la respuesta y registrar el log con el ID de petición.
```

---

## Integración con Módulo Product
