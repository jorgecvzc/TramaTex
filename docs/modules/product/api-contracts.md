# Product Module: API Contracts

- **Version:** 1.0.0
- **Status:** Borrador
- **Related Documents:** [Domain Model](./domain-model.md), [Use Cases](./use-cases.md)

Este documento define la API para el módulo de `product`.

---

## Principios Generales

- La API seguirá un estilo RESTful.
- Los formatos de datos para request y response body serán `application/json`.
- Las respuestas de error seguirán un formato estándar: `{ "error": { "code": "...", "message": "..." } }`.

---

## 1. Gestión de `ProductOptionSet`

Recurso base: `/product-option-sets`

### 1.1. DTOs

#### `ScopeDto`
```json
{
  "type": "GENERIC | BRAND | BRAND_GROUP",
  "brandId": "string | null",
  "productGroupId": "string | null"
}
```

#### `ProductOptionSetDto`
```json
{
  "id": "string",
  "name": "string",
  "attributeName": "string",
  "values": ["string"],
  "scope": { "$ref": "#/1.1.DTOs/ScopeDto" }
}
```

### 1.2. Endpoints

- `POST /product-option-sets`: Crea un nuevo conjunto de opciones. (UC-P-001)
- `GET /product-option-sets`: Obtiene una lista de conjuntos de opciones, con filtros opcionales por `scopeType`, `brandId`, `productGroupId`.
- `GET /product-option-sets/{id}`: Obtiene un conjunto de opciones específico.
- `PUT /product-option-sets/{id}`: Actualiza un conjunto de opciones. (UC-P-002)

---

## 2. Gestión de `Product`

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
  "groupId": "string",
  "directOptionSetIds": ["string"],
  "isActive": "boolean"
}
```

### 2.2. Endpoints

- `POST /products`: Crea una nueva plantilla de producto. (UC-P-003)
- `GET /products/{id}`: Obtiene la información de una plantilla de producto.
- `GET /products`: Obtiene una lista de productos.
- `POST /products/{id}/direct-option-sets`: Asigna un conjunto de opciones directamente a un producto. (UC-P-004)
- `GET /products/{id}/calculated-option-sets`: Obtiene la lista completa de `ProductOptionSet` aplicables a un producto. (UC-P-005)

---

## 3. Gestión de `ProductVariant`

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

## 4. Gestión de `PartyServiceConfiguration`

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
