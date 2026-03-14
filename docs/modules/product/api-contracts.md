# Product Module: API Contracts

- **Version:** 1.2.0
- **Status:** ✅ **Implementado (100% MVP)**
- **Última actualización:** 7 de marzo de 2026

Este documento define la API consolidada para el módulo de `product`, reflejando la implementación actual en el sistema TramaTex.

---

## Principios Generales

- La API sigue un estilo RESTful.
- Formato de intercambio: `application/json`.
- Respuestas de error estándar: `{ "error": "mensaje descriptivo" }`.
- Los IDs son UUIDs v4.

---

## 1. Gestión de Atributos (`/api/attributes`)

### 1.1. DTOs
- **`AttributeDTO`**: `id`, `name`, `code`, `sortOrder`, `values` ([]AttributeValueDTO).
- **`AttributeValueDTO`**: `id`, `value`, `code`, `hasPriceModifier`, `modifierType` (FIXED|PERCENTAGE), `modifierAmount`.

### 1.2. Endpoints
- `GET /api/attributes`: Lista completa con sus valores.
- `POST /api/attributes`: Crea atributo y valores (no incluir IDs en el body).
- `GET /api/attributes/:id`: Detalle de atributo.
- `PUT /api/attributes/:id`: Actualización completa (soporta añadir/editar/eliminar valores mediante el array).
- `DELETE /api/attributes/:id`: Eliminación física (solo si no tiene referencias).

---

## 2. Gestión de Productos (`/api/products`)

### 2.1. DTOs
- **`ProductDTO`**: `id`, `sku`, `name`, `longName`, `barcode`, `description`, `productType` (TANGIBLE|SERVICE), `brandId`, `groupIds` ([]uuid), `directAttributeIds` ([]uuid), `basePrice`, `taxRate`, `isActive`.

### 2.2. Endpoints
- `GET /api/products`: Listado con filtros (`search`, `brand_id`, `is_active`).
- `POST /api/products`: Creación de producto base.
- `GET /api/products/:id`: Detalle del producto.
- `PUT /api/products/:id`: Actualización de campos y asignaciones.

---

## 3. Gestión de Variantes (`/api/variants` y `/api/products/:id/variants`)

### 3.1. DTOs
- **`ProductVariantDTO`**: `id`, `productId`, `productName`, `sku`, `barcode`, `baseCost` (Calculado), `status` (PROVISIONAL|CONFIRMED), `optionConfiguration` (Map<AttrName, ValName>), `isActive`.

### 3.2. Endpoints
- `GET /api/products/:id/variants`: Lista todas las variantes de un producto.
- `POST /api/products/:id/variants/find-or-create`: **Punto central JIT**. Busca o crea una variante según su configuración.
  - Body: `{ "optionConfiguration": { "Talla": "XL", "Color": "Negro" } }`
- `GET /api/variants/:id`: Detalle de variante.
- `GET /api/variants?sku={sku}`: Búsqueda rápida por SKU exacto.

---

## 4. Otros Recursos

### 4.1. Marcas (`/api/brands`)
- CRUD completo para gestionar marcas y sus márgenes comerciales por defecto (`DefaultMarkupPercentage`).

### 4.2. Grupos de Producto (`/api/product-groups`)
- CRUD completo para jerarquías de categorías. Requiere clasificación `Type` (TANGIBLE|SERVICE).

### 4.3. Configuraciones de Servicio (`/api/parties/:id/service-configurations`)
- Gestión de reglas JSON flexibles por cliente (ej. reglas de precios especiales).

---

## 5. Lógica de Coste de Variante (`BaseCost`)

El campo `baseCost` de una variante es **calculado dinámicamente** y no se persiste.

**Algoritmo de cálculo:**
```
baseCost = Product.BasePrice
           + SUM(Product.BasePrice * (ModifierAmount / 100)) // Para PERCENTAGE
           + SUM(ModifierAmount)                             // Para FIXED
```
*Nota: Los porcentajes se aplican siempre sobre el precio base original del producto, no de forma acumulativa.*

---
**Última Actualización:** 2026-03-07
