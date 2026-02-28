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
      "code": "string"
    }
  ]
}
```

**Notas:**
- `code`: Código único del atributo usado para construcción de SKUs (ej. "T" para Talla, "C" para Color)
- `values`: Array de objetos `AttributeValueDTO` con estructura completa incluyendo IDs
- Los IDs son UUIDs generados automáticamente por el sistema
- Soporte completo para UTF-8 (caracteres españoles: ñ, á, é, etc.)

#### `AttributeValueDTO`
```json
{
  "id": "uuid",
  "value": "string",
  "code": "string"
}
```

**Notas:**
- `value`: Nombre descriptivo del valor (ej. "Large", "Rojo")
- `code`: Código corto usado en SKUs (ej. "L", "R")
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
    {"value": "Small", "code": "S"},
    {"value": "Medium", "code": "M"},
    {"value": "Large", "code": "L"}
  ]
}
```

**Notas sobre el request:**
- No incluir `id` en el payload (auto-generado)
- `values`: Enviar solo `value` y `code` (sin IDs para creación)
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
      "code": "S"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440002",
      "value": "Medium",
      "code": "M"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440003",
      "value": "Large",
      "code": "L"
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
      {"id": "660e8400-e29b-41d4-a716-446655440001", "value": "Small", "code": "S"},
      {"id": "660e8400-e29b-41d4-a716-446655440002", "value": "Medium", "code": "M"}
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
    {"id": "660e8400-e29b-41d4-a716-446655440001", "value": "Small", "code": "S"}
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
    {"id": "660e8400-e29b-41d4-a716-446655440001", "value": "Extra Small", "code": "XS"},
    {"value": "Extra Large", "code": "XL"}
  ]
}
```

**Notas importantes:**
- Para **editar** un valor existente: incluir su `id`
- Para **agregar** un nuevo valor: enviarlo sin `id` (se auto-genera)
- Para **eliminar** un valor: no incluirlo en el array
- Los valores con IDs presentes se actualizan, sin IDs se crean nuevos

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Talla Actualizada",
  "code": "T",
  "sortOrder": 2,
  "values": [
    {"id": "660e8400-e29b-41d4-a716-446655440001", "value": "Extra Small", "code": "XS"},
    {"id": "770e8400-e29b-41d4-a716-446655440004", "value": "Extra Large", "code": "XL"}
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
  "price": "number",
  "optionConfiguration": {
    "Talla": "L",
    "Color": "Rojo"
  },
  "isActive": "boolean"
}
```

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
