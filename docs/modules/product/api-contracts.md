# Product Module: API Contracts

- **Version:** 1.1.0
- **Status:** Implementado (MVP parcial)
- **Related Documents:** [Domain Model](./domain-model.md), [Use Cases](./use-cases.md)

Este documento define la API para el módulo de `product`.

---

## Estado de Implementación

### ✅ Endpoints Implementados
- **Gestión de Attributes**: API completa con CRUD funcional
- **DTOs actualizados**: `AttributeDTO` y `AttributeValueDTO` con estructura completa

### 🚧 Pendiente
- Gestión de Products
- Gestión de ProductVariants
- Gestión de PartyServiceConfiguration

---

## Principios Generales

- La API seguirá un estilo RESTful.
- Los formatos de datos para request y response body serán `application/json`.
- Las respuestas de error seguirán un formato estándar: `{ "error": { "code": "...", "message": "..." } }`.

---

## 1. Gestión de `Attribute` ✅

Recurso base: `/attributes`

> **Estado**: ✅ Implementado y funcional (MVP)

### 1.1. DTOs

#### `AttributeDTO`
```json
{
  "id": "uuid",
  "name": "string",
  "code": "string",
  "sortOrder": "number",
  "values": [
    {
      "id": "uuid",
      "value": "string",
      "code": "string",
      "hasPriceModifier": "boolean",
      "modifierType": "FIXED | PERCENTAGE | null",
      "modifierAmount": "number"
    }
  ]
}
```

**Notas:**
- `code`: Código único del atributo usado para construcción de SKUs (ej. "T" para Talla, "C" para Color)
- `values`: Array de objetos `AttributeValueDTO` con estructura completa incluyendo price modifiers
- Los IDs son UUIDs generados automáticamente por el sistema
- Soporte completo para UTF-8 (caracteres españoles: ñ, á, é, etc.)
- Los price modifiers permiten que cada valor de atributo modifique el precio base del producto

#### `AttributeValueDTO`
```json
{
  "id": "uuid",
  "value": "string",
  "code": "string",
  "hasPriceModifier": "boolean",
  "modifierType": "FIXED | PERCENTAGE | null",
  "modifierAmount": "number"
}
```

**Notas:**
- `value`: Nombre descriptivo del valor (ej. "Large", "Rojo")
- `code`: Código corto usado en SKUs (ej. "L", "R")
- `hasPriceModifier`: Indica si este valor afecta al precio base del producto
- `modifierType`: Tipo de modificación (FIXED = cantidad fija en €, PERCENTAGE = porcentaje del precio base)
- `modifierAmount`: Cantidad del modificador. Puede ser **positiva o negativa**
  - Si FIXED: Se suma/resta directamente (ej. +5.00 € o -2.50 €)
  - Si PERCENTAGE: Se aplica como porcentaje (ej. 10 = +10%, -15 = -15%)
- Este DTO se retorna dentro de `AttributeDTO.values` con estructura completa para permitir edición

### 1.2. Endpoints ✅

#### `POST /api/attributes` ✅
Crea un nuevo atributo con sus valores.

**Request Body:**
```json
{
  "name": "Talla",
  "code": "T",
  "sortOrder": 1,
  "values": [
    {"value": "Small", "code": "S", "hasPriceModifier": false},
    {"value": "Medium", "code": "M", "hasPriceModifier": false},
    {
      "value": "Large", 
      "code": "L",
      "hasPriceModifier": true,
      "modifierType": "FIXED",
      "modifierAmount": 2.50
    }
  ]
}
```

**Notas sobre el request:**
- No incluir `id` en el payload (auto-generado)
- `values`: Estructura completa con campos opcionales de price modifier
  - `hasPriceModifier`: Obligatorio (por defecto false)
  - `modifierType` y `modifierAmount`: Solo necesarios si `hasPriceModifier=true`
  - `modifierAmount` puede ser positivo (incrementa precio) o negativo (reduce precio)
- `sortOrder` en camelCase

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Talla",
  "code": "T",
  "sortOrder": 1,
  "values": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "value": "Small",
      "code": "S",
      "hasPriceModifier": false,
      "modifierType": null,
      "modifierAmount": 0
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440002",
      "value": "Medium",
      "code": "M",
      "hasPriceModifier": false,
      "modifierType": null,
      "modifierAmount": 0
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440003",
      "value": "Large",
      "code": "L",
      "hasPriceModifier": true,
      "modifierType": "FIXED",
      "modifierAmount": 2.50
    }
  ]
}
```

---

#### `GET /api/attributes` ✅
Obtiene la lista completa de atributos.

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Talla",
    "code": "T",
    "sortOrder": 1,
    "values": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001", 
        "value": "Small", 
        "code": "S",
        "hasPriceModifier": false,
        "modifierType": null,
        "modifierAmount": 0
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440002", 
        "value": "Medium", 
        "code": "M",
        "hasPriceModifier": false,
        "modifierType": null,
        "modifierAmount": 0
      }
    ]
  }
]
```

---

#### `GET /api/attributes/:id` ✅
Obtiene un atributo específico por su ID.

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Talla",
  "code": "T",
  "sortOrder": 1,
  "values": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001", 
      "value": "Small", 
      "code": "S",
      "hasPriceModifier": false,
      "modifierType": null,
      "modifierAmount": 0
    }
  ]
}
```

---

#### `PUT /api/attributes/:id` ✅
Actualiza un atributo existente.

**Request Body:**
```json
{
  "name": "Talla Actualizada",
  "code": "T",
  "sortOrder": 2,
  "values": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001", 
      "value": "Extra Small", 
      "code": "XS",
      "hasPriceModifier": true,
      "modifierType": "FIXED",
      "modifierAmount": -1.50
    },
    {
      "value": "Extra Large", 
      "code": "XL",
      "hasPriceModifier": true,
      "modifierType": "PERCENTAGE",
      "modifierAmount": 15
    }
  ]
}
```

**Notas importantes:**
- Para **editar** un valor existente: incluir su `id` y todos los campos (value, code, price modifiers)
- Para **agregar** un nuevo valor: enviarlo sin `id` (se auto-genera)
- Para **eliminar** un valor: no incluirlo en el array
- Los valores con IDs presentes se actualizan, sin IDs se crean nuevos
- Los price modifiers son opcionales pero recomendados para variantes con diferencias de precio

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Talla Actualizada",
  "code": "T",
  "sortOrder": 2,
  "values": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001", 
      "value": "Extra Small", 
      "code": "XS",
      "hasPriceModifier": true,
      "modifierType": "FIXED",
      "modifierAmount": -1.50
    },
    {
      "id": "770e8400-e29b-41d4-a716-446655440004", 
      "value": "Extra Large", 
      "code": "XL",
      "hasPriceModifier": true,
      "modifierType": "PERCENTAGE",
      "modifierAmount": 15
    }
  ]
}
```

---

### 1.3. Errores Comunes

#### 400 Bad Request
- Payload mal formado o campos requeridos faltantes
- `code` duplicado (debe ser único por atributo)

#### 404 Not Found
- El `id` del atributo no existe

#### 500 Internal Server Error
- Error de base de datos o conexión

---

## 2. Gestión de `Product` 🚧

> **Estado**: 🚧 Pendiente de implementación

Recurso base: `/products`

### 2.1. DTOs

#### `ProductDto`
```json
{
  "id": "string",
  "sku": "string | null",
  "name": "string",
  "description": "string",
  "productType": "TANGIBLE | SERVICE",
  "brandId": "string",
  "groupIds": ["string"],
  "directAttributeIds": ["string"],
  "taxRate": "number",
  "isActive": "boolean"
}
```

**Campos:**
- `taxRate`: Tasa de impuesto (IVA) aplicable al producto como porcentaje (ej: 21.00 = 21%). Valores típicos: 21.00 (general), 10.00 (reducido), 4.00 (súper reducido), 0.00 (exento)

### 2.2. Endpoints

- `POST /products`: Crea una nueva plantilla de producto. (UC-P-003)
- `GET /products/{id}`: Obtiene la información de una plantilla de producto.
- `GET /products`: Obtiene una lista de productos, con filtros opcionales por `brandId`, `groupId`, `isActive`.
- `POST /products/{id}/direct-option-sets`: Asigna un conjunto de opciones directamente a un producto. (UC-P-004)
- `GET /products/{id}/calculated-option-sets`: Obtiene la lista completa de `ProductOptionSet` aplicables a un producto. (UC-P-005)

---

## 3. Gestión de `ProductVariant` 🚧

> **Estado**: 🚧 Pendiente de implementación

Recurso base: `/products/{productId}/variants` y `/variants`

### 3.1. DTOs

#### `ProductVariantDto`
```json
{
  "id": "string",
  "productId": "string",
  "sku": "string",
  "baseCost": "number",
  "optionConfiguration": {
    "Talla": "L",
    "Color": "Rojo"
  },
  "status": "PROVISIONAL | CONFIRMED",
  "isActive": "boolean"
}
```

**Campos:**
- `baseCost`: Precio base calculado de la variante. Se calcula automáticamente como:
  - `baseCost = producto.basePrice + sum(modificadores de atributos seleccionados)`
- Los modificadores pueden ser:
  - **FIXED**: Se suma/resta directamente (ej. +2.50€ o -1.00€)
  - **PERCENTAGE**: Se aplica como porcentaje del precio acumulado (ej. 10% = +10%, -15% = -15%)
- **Ejemplo de cálculo:**
  - Producto base: 50.00€
  - Talla L (FIXED +2.50€): 52.50€
  - Color Premium (PERCENTAGE +10%): 52.50 + (52.50 × 0.10) = 57.75€
  - **Resultado final**: baseCost = 57.75€

### 3.2. Endpoints

- `POST /products/{productId}/variants/generate`: Genera explícitamente un subconjunto o todas las posibles variantes para un producto. Es una operación opcional, útil para pre-generar variantes que se sabe que tendrán stock o alta demanda. (UC-P-006)
- **Response:** `202 Accepted`
  - **Nota:** Esta es una operación potencialmente larga, por lo que se procesa de forma asíncrona. La respuesta indica que la tarea ha comenzado.
  - **Body:** `{ "taskId": "string", "status": "PENDING" }`

- `POST /products/{productId}/variants/find-or-create`: Obtiene una `ProductVariant` existente o la crea dinámicamente si no existe, basándose en un `Product` y una `OptionConfiguration` válida. (UC-P-008)
- **Request Body:**
```json
{
  "optionConfiguration": {
    "AttributeName1": "Value1",
    "AttributeName2": "Value2"
  }
}
```
- **Response:** `200 OK`
  - **Body:** `ProductVariantDto`

- `GET /products/{productId}/variants`: Obtiene una lista de todas las variantes de un producto.
- `GET /variants/{id}`: Obtiene una variante específica por su ID.
- `GET /variants?sku={sku}`: Obtiene una variante específica por su SKU.
- `PUT /variants/{id}`: Actualiza los datos de una variante específica (ej. precio). (UC-P-007)

---

## 4. Gestión de `PartyServiceConfiguration` 🚧

> **Estado**: 🚧 Pendiente de implementación

Recurso base: `/parties/{partyId}/service-configurations`

### 4.1. DTOs

#### `PartyServiceConfigurationDto`
```json
{
  "id": "string",
  "partyId": "string",
  "serviceId": "string",
  "name": "string",
  "configurationDetails": {} // Objeto JSON flexible
}
```

### 4.2. Endpoints

- `POST /parties/{partyId}/service-configurations`: Crea una nueva configuración de servicio para un cliente.
- `GET /parties/{partyId}/service-configurations`: Obtiene todas las configuraciones de servicio de un cliente.
- `GET /parties/{partyId}/service-configurations/{id}`: Obtiene una configuración de servicio específica por su ID.
- `PUT /parties/{partyId}/service-configurations/{id}`: Actualiza una configuración de servicio existente.
- `DELETE /parties/{partyId}/service-configurations/{id}`: Elimina una configuración de servicio.

---

## 5. Cálculo de Precio de Variante

### 5.1. Concepto

El precio base de una variante (`baseCost`) se calcula dinámicamente a partir del precio base del producto más los modificadores de precio de los atributos seleccionados. Este valor es fundamental para el módulo de Pricing, que lo utilizará para calcular el precio de venta final aplicando márgenes de beneficio y descuentos.

**Flujo completo de pricing:**

```
Producto.basePrice (50.00€)
    ↓
+ Modificadores de Atributos
    ↓
= Variante.baseCost (57.75€)
    ↓
× Margen de Beneficio de Marca (+30%)
    ↓
= Precio de Venta (75.08€)
    ↓
- Descuentos de Cliente/Volumen
    ↓
= Precio Final al Cliente
```

### 5.2. Algoritmo de Cálculo

**Función:** `domain.CalculateBaseCost(basePrice float64, attributeValues []AttributeValue) float64`

**Pseudocódigo:**

```go
baseCost := producto.basePrice

for cada attributeValue en variante.attributeValues:
    if !attributeValue.hasPriceModifier:
        continue
    
    if attributeValue.modifierType == "FIXED":
        // Suma/resta cantidad fija
        baseCost += attributeValue.modifierAmount
    
    else if attributeValue.modifierType == "PERCENTAGE":
        // Aplica porcentaje sobre precio acumulado
        baseCost += baseCost × (attributeValue.modifierAmount / 100.0)

// Asegurar que no sea negativo
if baseCost < 0:
    baseCost = 0

return baseCost
```

### 5.3. Ejemplos de Cálculo

#### Ejemplo 1: Modificadores FIXED

**Producto:** Camiseta Base
- `basePrice`: 20.00€

**Atributos seleccionados:**
- Talla L: FIXED +2.50€
- Color Azul: FIXED +0.00€
- Estampado Logo: FIXED +3.50€

**Cálculo:**
```
baseCost = 20.00€
         + 2.50€  (Talla L)
         + 0.00€  (Color Azul)
         + 3.50€  (Estampado Logo)
         = 26.00€
```

#### Ejemplo 2: Modificadores PERCENTAGE

**Producto:** Sofá Base
- `basePrice`: 500.00€

**Atributos seleccionados:**
- Talla XXL: PERCENTAGE +20%
- Tejido Premium: PERCENTAGE +15%

**Cálculo:**
```
baseCost = 500.00€
         + (500.00€ × 0.20) = 600.00€  (Talla XXL)
         + (600.00€ × 0.15) = 690.00€  (Tejido Premium)
         = 690.00€
```

#### Ejemplo 3: Modificadores Mixtos (FIXED + PERCENTAGE)

**Producto:** Pantalón Base
- `basePrice`: 50.00€

**Atributos seleccionados:**
- Talla L: FIXED +2.50€
- Color Premium: PERCENTAGE +10%
- Tipo Slim Fit: FIXED +5.00€

**Cálculo:**
```
baseCost = 50.00€
         + 2.50€                    = 52.50€  (Talla L FIXED)
         + (52.50€ × 0.10)          = 57.75€  (Color PERCENTAGE)
         + 5.00€                    = 62.75€  (Tipo FIXED)
         = 62.75€
```

#### Ejemplo 4: Modificadores Negativos (Descuentos)

**Producto:** Cortina Base
- `basePrice`: 80.00€

**Atributos seleccionados:**
- Talla S: FIXED -5.00€ (más pequeña, menos tela)
- Tejido Básico: PERCENTAGE -10% (material económico)

**Cálculo:**
```
baseCost = 80.00€
         - 5.00€                    = 75.00€  (Talla S)
         - (75.00€ × 0.10)          = 67.50€  (Tejido Básico)
         = 67.50€
```

### 5.4. Consideraciones Técnicas

1. **Orden de Aplicación:** Los modificadores se aplican en el orden definido por `Attribute.sortOrder`. Los modificadores PERCENTAGE se calculan sobre el precio acumulado, no sobre el precio base original.

2. **Precisión:** El cálculo se realiza con precisión de 2 decimales (centavos de euro).

3. **Validación:** El sistema valida que:
   - Los modificadores PERCENTAGE estén entre -100% y +1000%
   - El `baseCost` final no sea negativo (se ajusta a 0 si resulta negativo)

4. **Persistencia:** El `baseCost` **NO se almacena** en la base de datos. Se calcula dinámicamente cada vez que se consulta una variante, garantizando que cualquier cambio en los modificadores de precio de los atributos se refleje automáticamente.

5. **Integración con Pricing:** El módulo de Pricing consume el `baseCost` calculado para aplicar:
   - Márgenes de beneficio de marca
   - Descuentos específicos de cliente
   - Descuentos por volumen de venta
   - Precios especiales o promociones

### 5.5. API de Cálculo

Aunque el cálculo es automático, el módulo expone una función utilizada internamente:

```go
// domain/attribute.go
func CalculateBaseCost(basePrice float64, attributeValues []AttributeValue) float64
```

Esta función es invocada por:
- El servicio de creación/actualización de variantes
- El endpoint `/products/{id}/variants/find-or-create` (JIT variant creation)
- El módulo de Pricing al solicitar costos de variantes
