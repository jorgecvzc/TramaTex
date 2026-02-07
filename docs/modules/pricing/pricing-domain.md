# Dominio del Módulo de Precios

Este documento describe el dominio del módulo de Precios (Bounded Context de `Pricing`) dentro de TramaTex, siguiendo los principios de Domain-Driven Design (DDD). Se enfoca en el Lenguaje Ubicuo, las Entidades, los Value Objects y los Servicios de Dominio que componen este contexto, así como las reglas de negocio clave.

## 1. Lenguaje Ubicuo y Conceptos Clave

El Lenguaje Ubicuo para el contexto de `Pricing` incluye los siguientes términos:

*   **Precio Base de Producto:** El costo o precio de referencia de un producto antes de cualquier modificación de variante, margen o descuento. Se obtiene del módulo `Product`.
*   **Modificador de Variante:** Un factor (positivo o negativo) que ajusta el `Precio Base de Producto` según las características de una variante específica (ej., talla, color, material). Se obtiene del módulo `Product`.
*   **Margen de Beneficio de Marca:** Un porcentaje o cantidad fija que se añade al costo de un producto para determinar su precio de venta, específico para una `Marca`.
*   **Precio de Venta:** El precio final de un producto o variante antes de la aplicación de cualquier descuento de venta.
*   **Descuento de Venta:** Una reducción en el `Precio de Venta`, que puede ser específica de un `Cliente` o basarse en la `Cantidad Total de la Venta`.
*   **Regla de Precios:** Un conjunto de criterios y lógica que determina cómo se calcula un precio o se aplica un descuento.
*   **Cliente:** La entidad que compra productos (del módulo `Party`).
*   **Producto/Variante:** Los bienes que se venden (del módulo `Product`).
*   **Cálculo de Precio:** Una operación que determina un precio de venta o un descuento, y también el registro inmutable de dicha operación.

## 2. Entidades del Dominio

Las Entidades son objetos que tienen una identidad única y encapsulan comportamientos y reglas de negocio, siendo mutables a lo largo del tiempo.

*   **`PricingRule`**:
    *   **Propósito:** Representa una regla general o específica que influye en el cálculo de precios (ej., "precio fijo para el producto X en el periodo Y", "descuento del 10% para la categoría Z").
    *   **Identidad:** `ID` único.
    *   **Atributos:** `Nombre`, `Descripción`, `Tipo` (margen, descuento, fijo), `Condiciones` (producto, categoría, fecha), `Valor` (porcentaje, cantidad), `Prioridad`.
    *   **Comportamiento:** `Evaluar(contextoDeCalculo)`: determina si la regla aplica y cómo modifica el precio.
*   **`ClientPricing`**:
    *   **Propósito:** Gestiona las reglas de precios o descuentos que son específicas para un cliente particular, anulando o complementando las reglas generales.
    *   **Identidad:** `ID` único.
    *   **Atributos:** `ID_Cliente` (referencia al módulo `Party`), `ID_PricingRule` (referencia a una `PricingRule`), `ValorEspecífico` (opcional).
    *   **Comportamiento:** `ObtenerReglasAplicables(cliente)`: retorna las reglas de precios/descuento específicas de ese cliente.
*   **`BrandProfitMargin`**:
    *   **Propósito:** Define el porcentaje o la cantidad de beneficio que se debe aplicar a los productos de una `Marca` específica.
    *   **Identidad:** `ID` único.
    *   **Atributos:** `ID_Marca` (referencia al módulo `Product`), `PorcentajeBeneficio` (Value Object `Percentage`), `CantidadFijaBeneficio` (Value Object `Money`), `FechaInicio`, `FechaFin`.
    *   **Comportamiento:** `ObtenerMargen(marca, fecha)`: retorna el margen aplicable.
*   **`SalesDiscountRule`**:
    *   **Propósito:** Encapsula las reglas específicas para aplicar descuentos durante una venta. Puede ser por volumen, por cliente leal, promociones, etc.
    *   **Identidad:** `ID` único.
    *   **Atributos:** `Nombre`, `Descripción`, `Tipo` (porcentaje, cantidad fija), `Condiciones` (mínimo de compra, ID_Cliente, ID_Producto), `Valor` (Value Object `Percentage` o `Money`), `Prioridad`.
    *   **Comportamiento:** `Aplicar(contextoDeVenta)`: calcula el descuento a aplicar.
*   **`PriceCalculation`**:
    *   **Propósito:** Entidad inmutable que registra el resultado de un cálculo de precio, actuando como un registro de auditoría.
    *   **Identidad:** `ID` único.
    *   **Atributos:** `ID_Producto/Variante`, `ID_Cliente`, `PrecioBase`, `ModificadoresAplicados`, `MargenBeneficioAplicado`, `PrecioVentaCalculado`, `DescuentosAplicados` (lista de `SalesDiscount` con su valor), `PrecioFinal`, `FechaCalculo`, `ID_Usuario`.
    *   **Comportamiento:** `Guardar()`: Persiste el registro inmutable.

## 3. Value Objects del Dominio

Los Value Objects son objetos inmutables que representan conceptos descriptivos del dominio, no tienen identidad y se definen por sus atributos.

*   **`Money`**:
    *   **Propósito:** Representa una cantidad monetaria con su moneda.
    *   **Atributos:** `Cantidad` (ej., `decimal.Decimal` para precisión), `Moneda` (ej., "EUR", "USD").
    *   **Inmutabilidad:** Sus valores no cambian una vez creado.
    *   **Comportamiento:** `Sumar()`, `Restar()`, `Multiplicar()`, `Dividir()`, `EsMayorQue()`, `EsMenorQue()`.
*   **`Percentage`**:
    *   **Propósito:** Representa un porcentaje.
    *   **Atributos:** `Valor` (ej., `float64` o `decimal.Decimal`).
    *   **Inmutabilidad:** Su valor no cambia una vez creado.
    *   **Comportamiento:** `AplicarA(cantidad)`: calcula un porcentaje de una cantidad.
*   **`Brand`**:
    *   **Propósito:** Identificador y descripción de una marca de producto. (Proveniente del módulo `Product`, usado aquí como VO).
    *   **Atributos:** `ID_Marca`, `Nombre`.
    *   **Inmutabilidad:** Sus valores no cambian una vez creado.
*   **`ProductCode` / `VariantCode`**:
    *   **Propósito:** Identificadores únicos y validados para productos y sus variantes. (Provenientes del módulo `Product`, usados aquí como VO).
    *   **Atributos:** `Código`.
    *   **Validación:** Asegura que el formato del código es correcto.

## 4. Servicios de Dominio

Los Servicios de Dominio encapsulan lógica de negocio que no pertenece naturalmente a ninguna Entidad ni Value Object, a menudo orquestando el comportamiento de varias de ellas.

*   **`SellingPriceCalculatorService`**:
    *   **Propósito:** Orquestar la lógica de cálculo del `Precio de Venta` para un `Producto/Variante` dado, aplicando `Modificadores de Variante` y `Márgenes de Beneficio de Marca`.
    *   **Dependencias:** `IRepository<PricingRule>`, `IRepository<BrandProfitMargin>`, Clientes del módulo `Product` (para `Precio Base`, `Modificador de Variante`, `Marca`).
    *   **Comportamiento Clave:** `CalcularPrecioVenta(producto, variante, marca, contexto)`:
        1.  Obtiene el `Precio Base de Producto` y el `Modificador de Variante` del módulo `Product`.
        2.  Aplica el `Modificador de Variante` al `Precio Base`.
        3.  Obtiene el `Margen de Beneficio de Marca` (utilizando `BrandProfitMargin`).
        4.  Aplica el `Margen de Beneficio` para obtener el `Precio de Venta`.
        5.  Retorna el `Precio de Venta` como un Value Object `Money`.
*   **`SalesDiscountCalculatorService`**:
    *   **Propósito:** Orquestar la lógica de aplicación de `Descuentos de Venta` sobre un `Precio de Venta` o una `Cantidad Total de la Venta`, considerando `Reglas de Descuento` específicas de `Cliente` y otras condiciones.
    *   **Dependencias:** `IRepository<SalesDiscountRule>`, `IRepository<ClientPricing>`, Clientes del módulo `Party` (para `Cliente`).
    *   **Comportamiento Clave:** `AplicarDescuentos(precioVenta, cliente, cantidadTotalVenta, productosEnVenta, contexto)`:
        1.  Obtiene `SalesDiscountRule` aplicables (generales, por cliente).
        2.  Evalúa `SalesDiscountRule` basadas en la `Cantidad Total de la Venta`.
        3.  Calcula el `Descuento Total` como un Value Object `Money` o `Percentage`.
        4.  Retorna el `Descuento Total` y el `Precio Final` tras el descuento.

## 5. Reglas de Negocio Centrales

*   **Cálculo del Precio de Venta:** El `Precio de Venta` se deriva del `Precio Base de Producto`, aplicando los `Modificadores de Variante` y luego un `Margen de Beneficio de Marca`.
*   **Aplicación de Descuentos:** El `Precio de Venta` puede ser modificado por `Descuentos de Venta` específicos de `Cliente` o basados en la `Cantidad Total de la Venta`.
*   **Prioridad de Reglas:** Las reglas de `ClientPricing` prevalecen sobre las `PricingRule` generales. Las reglas de `SalesDiscountRule` tienen su propia prioridad de aplicación.
*   **Auditabilidad:** Cada cálculo de precio completo (que incluye precio de venta y descuentos) debe generar un registro inmutable de `PriceCalculation`.
*   **Consistencia:** Los precios de venta obtenidos deben ser consistentes con las reglas actuales, garantizado por la estrategia de caching e invalidación (limpieza completa de caché).

## 6. Interacciones con Otros Módulos (Nivel de Dominio)

El módulo de `Pricing` interactúa con otros Bounded Contexts principalmente para obtener información de entrada necesaria para sus cálculos:

*   **Módulo `Product`:** `Pricing` consume `ProductBasePrice`, `VariantModifier`, `Brand` (o sus identificadores) para el cálculo del precio de venta.
*   **Módulo `Party`:** `Pricing` consume `Cliente` (o sus identificadores y atributos relevantes) para aplicar `ClientPricing` y `SalesDiscountRule` específicos.

Estas interacciones se realizan a través de interfaces bien definidas en el `Application Layer` de `Pricing`, que son implementadas por clientes en el `Infrastructure Layer`. El `Domain Layer` de `Pricing` permanece agnóstico a los detalles de implementación de `Product` y `Party`.
