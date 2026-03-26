# Módulo de Pricing

Este módulo es el encargado de toda la lógica relacionada con el cálculo de precios en el sistema TramaTex. Se adhiere a los principios de Clean Architecture y Domain-Driven Design, separando claramente las responsabilidades entre el cálculo de precios base y las modificaciones de precio en el momento de la venta.

## Diseño Arquitectónico

Para una descripción detallada de las decisiones arquitectónicas, entidades de dominio, objetos de valor, casos de uso, estrategia de caché e infraestructura asociada a este módulo, consulte el siguiente Architectural Decision Record (ADR):

*   [ADR-016: Arquitectura del Módulo de Pricing](../../architecture/adrs/adr-016-pricing-module-architecture.md)

## Documentación Detallada

Consulte los siguientes documentos para una descripción más profunda del módulo de Pricing:

*   **Especificación del Módulo:** [module-spec.md](./module-spec.md)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md)
*   **Casos de Uso:** [use-cases.md](./use-cases.md)
*   **Contratos de API:** [api-contracts.md](./api-contracts.md)
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md)
*   **Integración con otros módulos:** [integration.md](./integration.md)
*   **Diagrama de Dominio:** [diagrams/domain-model.md](./diagrams/domain-model.md)
*   **Diagramas del Módulo:** [diagrams/pricing-module-diagrams.md](./diagrams/pricing-module-diagrams.md)

## Componentes Clave

*   **Entidades de Dominio:**
    *   `BaseSalesPriceRule`: Define cómo se construye el precio base de venta a partir de un costo/tarifa, mediante incrementos y una lógica de precedencia específica.
    *   `SaleModificationRule`: Define cómo se modifica el precio base de venta en el contexto de una transacción de venta (descuentos por cliente, monto mínimo, grupo de productos o fechas de campaña).

*   **Objetos de Valor:**
    *   `Money`: Cantidad monetaria (siempre en EUR para MVP).
    *   `Percentage`: Valor porcentual.
    *   `RuleValue`: Encapsula el tipo y el valor del efecto de una regla.

*   **Casos de Uso (Capa de Aplicación):**
    *   `CalculateBaseSalesPriceForProductVariantUseCase`: Calcula y cachea los precios base de venta.
    *   `CalculateFinalSalePriceUseCase`: Recupera precios de caché y aplica reglas de modificación de venta.

*   **Infraestructura de Caché:**
    *   **Tecnología:** Redis.
    *   **Estrategia:** Caché por `ProductID` conteniendo todos sus `ProductVariantID`s y sus `BaseSalesPrice`s. Maneja invalidación por eventos (creación/modificación de variantes) y TTL.

## Persistencia

Las reglas de pricing se persisten en una base de datos relacional (PostgreSQL) en tablas separadas:

*   `base_sales_price_rules`
*   `sale_modification_rules`
*   `rule_value_types` (tabla lookup para tipos de valores de reglas)

---

## Cobertura de Pruebas

El modulo de Pricing requiere **≥90% de cobertura unitaria** en dominio. Para generar un reporte local:

```bash
cd apps/tramatex-api
make coverage-pricing
```

El reporte HTML queda en `apps/tramatex-api/coverage-pricing.html`.

## Contratos de API (DTOs)

Los siguientes DTOs (Data Transfer Objects) definen los contratos para la interacción con el módulo de Pricing, tanto para la gestión de reglas como para los casos de uso de cálculo.

### DTOs Comunes

1.  **`MoneyDTO`**
    *   `amount`: DECIMAL (ej. 123.45)
    *   `currency`: String (fijo a "EUR" para MVP)

2.  **`PercentageDTO`**
    *   `value`: DECIMAL (ej. 0.10 para 10%)

3.  **`RuleValueDTO`**
    *   `type`: String (ENUM. Ej: `PERCENTAGE_MARKUP`, `FIXED_AMOUNT_INCREASE`, `SET_TO_FIXED_PRICE`, `APPLY_PERCENTAGE_DISCOUNT`, `APPLY_FIXED_AMOUNT_DISCOUNT`, `SET_TO_FIXED_DISCOUNTED_PRICE`)
    *   `percentageValue`: `PercentageDTO` (opcional)
    *   `moneyValue`: `MoneyDTO` (opcional)

### DTOs para Gestión de `BaseSalesPriceRule`

1.  **`CreateBaseSalesPriceRuleRequest`**
    *   `name`: String
    *   `brandId`: UUID (nullable)
    *   `productGroupId`: UUID (nullable)
    *   `productId`: UUID (nullable)
    *   `variantId`: UUID (nullable)
    *   `value`: `RuleValueDTO`
    *   `isActive`: Boolean (default `true`)

2.  **`UpdateBaseSalesPriceRuleRequest`**
    *   `id`: UUID
    *   `name`: String (opcional)
    *   `brandId`: UUID (opcional, nullable)
    *   `productGroupId`: UUID (opcional, nullable)
    *   `productId`: UUID (opcional, nullable)
    *   `variantId`: UUID (opcional, nullable)
    *   `value`: `RuleValueDTO` (opcional)
    *   `isActive`: Boolean (opcional)

3.  **`BaseSalesPriceRuleResponse`**
    *   `id`: UUID
    *   `name`: String
    *   `brandId`: UUID (nullable)
    *   `productGroupId`: UUID (nullable)
    *   `productId`: UUID (nullable)
    *   `variantId`: UUID (nullable)
    *   `value`: `RuleValueDTO`
    *   `isActive`: Boolean

### DTOs para Gestión de `SaleModificationRule`

1.  **`CreateSaleModificationRuleRequest`**
    *   `name`: String
    *   `clientIds`: List of UUIDs (nullable)
    *   `productGroupId`: UUID (nullable)
    *   `minOrderTotalAmount`: `MoneyDTO` (nullable)
    *   `value`: `RuleValueDTO`
    *   `priority`: Integer
    *   `isActive`: Boolean (default `true`)
    *   `effectiveFrom`: Timestamp (ISO 8601 string)
    *   `effectiveTo`: Timestamp (ISO 8601 string, nullable)

2.  **`UpdateSaleModificationRuleRequest`**
    *   `id`: UUID
    *   `name`: String (opcional)
    *   `clientIds`: List of UUIDs (opcional, nullable)
    *   `productGroupId`: UUID (opcional, nullable)
    *   `minOrderTotalAmount`: `MoneyDTO` (opcional, nullable)
    *   `value`: `RuleValueDTO` (opcional)
    *   `priority`: Integer (opcional)
    *   `isActive`: Boolean (opcional)
    *   `effectiveFrom`: Timestamp (ISO 8601 string, opcional)
    *   `effectiveTo`: Timestamp (ISO 8601 string, opcional, nullable)

3.  **`SaleModificationRuleResponse`**
    *   `id`: UUID
    *   `name`: String
    *   `clientIds`: List of UUIDs (nullable)
    *   `productGroupId`: UUID (nullable)
    *   `minOrderTotalAmount`: `MoneyDTO` (nullable)
    *   `value`: `RuleValueDTO`
    *   `priority`: Integer
    *   `isActive`: Boolean
    *   `effectiveFrom`: Timestamp (ISO 8601 string)
    *   `effectiveTo`: Timestamp (ISO 8601 string, nullable)

### DTOs para Casos de Uso de Cálculo

1.  **`CalculateBaseSalesPriceRequest`**
    *   `productId`: UUID
    *   `variantId`: UUID

2.  **`CalculatedBaseSalesPriceResponse`**
    *   `variantId`: UUID
    *   `baseSalesPrice`: `MoneyDTO`

3.  **`SaleItemRequest`**
    *   `productVariantId`: UUID
    *   `quantity`: Integer

4.  **`CalculateFinalSalePriceRequest`**
    *   `saleItems`: List of `SaleItemRequest`
    *   `clientId`: UUID (del cliente que realiza la venta)
    *   `saleDate`: Timestamp (ISO 8601 string)

5.  **`CalculatedSaleItemResponse`**
    *   `productVariantId`: UUID
    *   `quantity`: Integer
    *   `baseSalesPrice`: `MoneyDTO`
    *   `finalPrice`: `MoneyDTO`
    *   `taxRate`: DECIMAL (porcentaje de IVA, ej. 21.0)
    *   `finalPriceWithTax`: `MoneyDTO` (precio final con IVA incluido)

6.  **`CalculateFinalSalePriceResponse`**
    *   `calculatedItems`: List of `CalculatedSaleItemResponse`
    *   `saleTotal`: `MoneyDTO` (total sin IVA)
    *   `saleTotalWithTax`: `MoneyDTO` (total con IVA incluido)

---

## Documentación Detallada

*   **Modelo de Dominio:** [domain-model.md](./domain-model.md)
*   **Casos de Uso:** [use-cases.md](./use-cases.md)
*   **Integración Inter-modular:** [integration.md](./integration.md)
*   **Contratos de API:** [api-contracts.md](./api-contracts.md)
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md)
*   **Diagramas Detallados:** [diagrams/domain-model.md](./diagrams/domain-model.md)

---

## Interfaz de Usuario (Frontend)

El módulo Pricing se materializa en los siguientes componentes clave en `apps/frontend/`:

- **`PricingPanel.vue`**: Componente crítico que muestra el desglose del precio calculado (Base, Margen, Descuentos) en tiempo real.
- **`MasterData/Brands`**: Gestión administrativa de los márgenes de beneficio por marca.
- **`MasterData/Attributes`**: Configuración de los modificadores de precio por valor de atributo.
- **Integración con Sales**: El motor de precios se consume de forma asíncrona desde los formularios de Presupuesto y Pedido para actualizar los totales.


