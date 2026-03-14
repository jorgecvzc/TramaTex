# Diagrama de Dominio - Módulo Pricing

## 1. Diagrama de Clases (UML)

```mermaid
classDiagram
    direction TB
    
    class BaseSalesPriceRule {
        +ID: UUID
        +Name: string
        +Precedence: int
        +Scope: Enum (Global, Brand, Product, Client)
        +Value: Money | Percentage
        +Condition: string
        +Apply(context): Money
    }

    class BrandProfitMargin {
        +ID: UUID
        +BrandID: UUID
        +MarkupPercentage: Percentage
        +IsDefault: bool
    }

    class SalesDiscountRule {
        +ID: UUID
        +Name: string
        +ClientSpecific: bool
        +ProductSpecific: bool
        +DiscountValue: Percentage
        +Priority: int
    }

    class PriceCalculation {
        <<Service>>
        +CalculateBaseCost(Product, List~AttributeValue~): Money
        +CalculateBSP(BaseCost, Brand): Money
        +CalculateFSP(BSP, Client, Context): Money
    }

    class PriceResult {
        +BaseCost: Money
        +BSP: Money
        +FSP: Money
        +AppliedRules: List~string~
        +TaxRate: Percentage
        +FinalPriceWithTax: Money
    }

    class Money {
        +Amount: float64
        +Currency: string ("EUR")
    }

    class Percentage {
        +Value: float64
    }

    BaseSalesPriceRule "1" --o "1" Money : defines
    BaseSalesPriceRule "1" --o "1" Percentage : defines
    BrandProfitMargin "1" --o "1" Percentage : defines
    SalesDiscountRule "1" --o "1" Percentage : defines
    
    PriceCalculation ..> BaseSalesPriceRule : uses
    PriceCalculation ..> BrandProfitMargin : uses
    PriceCalculation ..> SalesDiscountRule : uses
    PriceCalculation ..> PriceResult : produces

    note for BaseSalesPriceRule "Precedencia: Override > Client > Product > Brand > Global"
    note for PriceCalculation "JIT: BaseCost = Product.BasePrice + Σ(AttributeValue.PriceAdjustment)"
```

## 2. Descripción y Jerarquía de Reglas

Este diagrama visualiza las entidades y el flujo de cálculo del motor de precios de TramaTex.

### Jerarquía de Precedencia
El motor de Pricing aplica las reglas en el siguiente orden de especificidad, donde la primera regla que coincide detiene la búsqueda para ese nivel:
1.  **Override (Manual):** Precio fijado explícitamente para una operación.
2.  **Cliente Específico:** Acuerdo de precios con una `Party` concreta.
3.  **Producto Específico:** Precio base definido para un `Product` o `ProductVariant`.
4.  **Marca:** Margen de beneficio por defecto de la `Brand`.
5.  **Global:** Regla por defecto del sistema.

### Conceptos de Cálculo
*   **BaseCost (JIT):** Calculado dinámicamente sumando el `BasePrice` del producto y los modificadores de los atributos seleccionados.
*   **BSP (Base Selling Price):** `BaseCost` aplicado el margen comercial (`BrandProfitMargin`).
*   **FSP (Final Selling Price):** `BSP` tras aplicar descuentos de cliente o promociones (`SalesDiscountRule`).
