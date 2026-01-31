# Módulo de Pricing (Motor de Precios)

## 1. Propósito

*   **Visión del Módulo:** Calcular precios basados en Party, Product y volumen.
*   **Objetivos Clave:**
    *   Proporcionar un motor de precios inteligente que calcule automáticamente el precio de venta.
    *   Calcular precios según categoría de cliente, producto, volumen y otras variables.

## 2. Requisitos

### 2.1. Requisitos Funcionales

*   **RF-001:** Definir reglas de precio base.
*   **RF-002:** Aplicar descuentos por volumen.
*   **RF-003:** Aplicar márgenes por categoría de Party.
*   **RF-004:** Calcular precio final automáticamente.
*   **RF-005:** Auditoría de cálculos.

## 3. Casos de Uso

### 3.1. Actores

*   **Gerente Comercial:** Define las reglas de precios.
*   **Vendedor:** Utiliza el motor para cotizar.
*   **Sistema:** Interactúa con el motor para obtener precios en el e-commerce.

### 3.2. Casos de Uso Principales

*   **CU-001: CreatePriceRule**
    *   **Actor:** Gerente Comercial
    *   **Descripción:** Crear una nueva regla de precio.
*   **CU-002: UpdatePriceRule**
    *   **Actor:** Gerente Comercial
    *   **Descripción:** Modificar una regla existente.
*   **CU-003: CalculatePrice**
    *   **Actor:** Vendedor/Sistema
    *   **Descripción:** Calcular el precio para una combinación de Producto, Cliente y Cantidad.
*   **CU-004: GetApplicableRules**
    *   **Actor:** Gerente Comercial
    *   **Descripción:** Obtener las reglas de precio vigentes para un producto.
*   **CU-005: ApplyBulkDiscount**
    *   **Actor:** Sistema
    *   **Descripción:** Aplicar un descuento por volumen según las reglas.

## 4. Historias de Usuario

*   **HU-001:** Como Gerente Comercial, quiero crear reglas de precio por categoría de cliente para asegurar márgenes de ganancia adecuados.
*   **HU-002:** Como Vendedor, quiero que el sistema calcule el precio automáticamente al crear una cotización para evitar errores manuales.
*   **HU-003:** Como Gerente Comercial, quiero definir descuentos por volumen para incentivar compras más grandes.

## 5. Criterios de Aceptación

*   **Para HU-001:**
    *   **Criterio 1:** Dado que defino un margen del 20% para clientes mayoristas, cuando un vendedor cotiza a un mayorista, entonces el precio final debe incluir dicho margen sobre el costo.

## 6. Modelo de Dominio

### PriceRule (Raíz de Agregación)
- **ID**: UUID
- **ProductVariantID**: UUID
- **PartyCategory**: Enum (Mayorista, Minorista, etc.)
- **CostBase**: Decimal
- **Markup**: Percentage
- **MinimumQuantity**: Integer
- **MaximumQuantity**: Integer (nullable)
- **EffectiveFrom**: DateTime
- **EffectiveTo**: DateTime (nullable)
- **Estado**: Enum (Activo, Inactivo)

### PriceCalculation
- **ID**: UUID
- **ProductVariantID**: UUID
- **PartyID**: UUID
- **Quantity**: Integer
- **BaseCost**: Decimal
- **AppliedRules**: List<String>
- **FinalPrice**: Decimal
- **CalculatedAt**: DateTime

## 7. Decisiones de Diseño

*   **Algoritmo de Cálculo:**
    1.  Obtener `CostoBase` del `ProductVariant`.
    2.  Buscar `PriceRule` aplicable (por categoría, cantidad y fecha).
    3.  Aplicar `Markup`: `FinalPrice = CostoBase * (1 + Markup%)`.
    4.  Registrar el cálculo en `PriceCalculation`.
    5.  Retornar `FinalPrice`.
*   **Relaciones con Otros Módulos:**
    *   **Product**: Lee costo de `ProductVariant`.
    *   **Party**: Lee categoría para aplicar `markup`.
    *   **Sales**: Usa cálculo de precio para generar cotizaciones y órdenes.
*   **Fases de Desarrollo:**
    *   [X] Fase 1 (MVP): Cálculo básico (CostoBase + Markup).
    *   [ ] Fase 2: Descuentos por volumen y categoría.
    *   [ ] Fase 3: Reglas complejas (campañas, promociones).