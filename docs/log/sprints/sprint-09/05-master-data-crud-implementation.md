# Tarea 09-05: Master Data CRUD + Product UPDATE Implementation

**Estado:** ✅ Completado  
**Fecha:** 2026-02-14  
**Duración:** 8 horas (incluyendo refactor de Atributos)  
**Sprint:** Sprint 09 (UI Definition Sprint)  
**Módulo:** Product

---

## 🎯 Objetivos

Implementar funcionalidades CRUD completas para las entidades maestras del módulo Product:
1. **Brands (Marcas)** - CREATE, READ, UPDATE, DELETE
2. **Product Groups (Categorías)** - CREATE, READ, UPDATE, DELETE con soporte jerárquico
3. **Attributes (Atributos)** - CREATE, READ, UPDATE, DELETE
4. **Product UPDATE** - Endpoint PUT completo

---

## ✅ Entregables

### Backend (Go)

#### Comandos
- `CreateBrandCommand` / `UpdateBrandCommand` / `DeleteBrandCommand`
- `CreateProductGroupCommand` / `UpdateProductGroupCommand` / `DeleteProductGroupCommand`
- `CreateAttributeCommand` / `UpdateAttributeCommand` / `DeleteAttributeCommand`
- `UpdateProductCommand` (9 campos: ActorID + ProductID + 7 opcionales)

#### Servicios
- `CreateBrand` / `UpdateBrand` / `DeleteBrand`
- `CreateProductGroup` / `UpdateProductGroup` / `DeleteProductGroup`
- `CreateAttribute` / `UpdateAttribute` / `DeleteAttribute`
- `UpdateProduct` (~50 líneas, actualización condicional)

#### Handlers HTTP
- POST/PUT/DELETE `/api/brands/:id`
- POST/PUT/DELETE `/api/product-groups/:id`
- POST/PUT/DELETE `/api/attributes/:id`
- PUT `/api/products/:id`

#### Rutas Registradas
```go
// main.go
brands.POST("/", infra_middleware.RequireRole("admin"), productHandler.CreateBrand)
brands.PUT("/:id", infra_middleware.RequireRole("admin"), productHandler.UpdateBrand)
brands.DELETE("/:id", infra_middleware.RequireRole("admin"), productHandler.DeleteBrand)

productGroups.POST("/", infra_middleware.RequireRole("admin"), productHandler.CreateProductGroup)
productGroups.PUT("/:id", infra_middleware.RequireRole("admin"), productHandler.UpdateProductGroup)
productGroups.DELETE("/:id", infra_middleware.RequireRole("admin"), productHandler.DeleteProductGroup)

attributes.POST("/", infra_middleware.RequireRole("admin"), productHandler.CreateAttribute)
attributes.PUT("/:id", infra_middleware.RequireRole("admin"), productHandler.UpdateAttribute)
attributes.DELETE("/:id", infra_middleware.RequireRole("admin"), productHandler.DeleteAttribute)

products.PUT("/:id", infra_middleware.RequireRole("admin"), productHandler.UpdateProduct)
```

---

### Frontend (Vue)

#### Componentes de Formularios
1. **BrandForm.vue** (178 líneas)
   - Campos: name
   - Validaciones inline
   - Modos: create/edit

2. **ProductGroupForm.vue** (233 líneas)
   - Campos: name, parentGroupId, isActive
   - Select dinámico para parent (jerarquía)
   - Prevención de referencias circulares
   - Validaciones inline

3. **AttributeForm.vue** (409 líneas)
   - Campos: name, code, order, values[]
   - Gestión dinámica de valores (agregar/eliminar)
   - Auto-uppercase para códigos
   - Validaciones inline

#### Páginas de Listado

1. **brands/List.vue** (476 líneas)
   - Tabla con columnas: Nombre, ID, Productos, Acciones
   - Modal para CREATE/EDIT
   - Botones: Editar, Eliminar (con confirmación)
   - Estados: loading, error, empty

2. **product-groups/List.vue** (507 líneas)
   - Tabla con columnas: Nombre, ID, Parent, Productos, Acciones
   - Muestra categoría padre por nombre
   - Modal para CREATE/EDIT
   - Botones: Editar, Eliminar (con confirmación)
   - Estados: loading, error, empty

3. **attributes/List.vue** (662 líneas)
   - Tabla con columnas: Nombre, Código, Valores, Acciones
   - Badge con contador de valores
   - Modal para CREATE/EDIT
   - Botones: Editar, Eliminar (con confirmación)
   - Estados: loading, error, empty

#### Servicios API

**productApi.js** - Métodos añadidos/corregidos:
```javascript
// Brands
createBrand(data)
updateBrand(id, data)
deleteBrand(id)

// Product Groups
createProductGroup(data) // Con transformación camelCase → snake_case
updateProductGroup(id, data) // CORREGIDO: ahora transforma correctamente
deleteProductGroup(id)

// Attributes
createAttribute(data)
updateAttribute(id, data)
deleteAttribute(id)

// Products
updateProduct(id, data) // Ya existía
```

---

## 🔧 Correcciones Realizadas

### 1. **Refactor Arquitectónico: Eliminación de Scope en Atributos**
   
**Cambio bloqueante solicitado:** Los atributos ya NO tienen concepto de Scope (Generic/Brand/ProductGroup). Son todos genéricos y se asignan directamente en productos.

**Backend:**
- ✅ `Attribute` domain ya estaba simplificado (sin scope, brandID, productGroupID)
- ✅ Comandos `CreateAttributeCommand` y `UpdateAttributeCommand` sin campos de scope
- ✅ Eliminado filtrado por `ScopeType` en `ListAttributesQuery`
- ✅ Simplificado `ListAttributes` service - devuelve todos sin filtros
- ✅ Simplificado handler `ListAttributes` - no parsea scopeType, brandId, productGroupId
- ✅ Eliminadas funciones helper `isValidScopeType()` y `attributeMatchesScopeType()`

**Frontend:**
- ✅ `AttributeForm.vue` ya estaba sin campos de Scope
- ✅ Solo tiene: name, code, order, values[]

**Tiempo:** ~1.5 horas

---

### 2. **Corrección de Transformación camelCase ↔ snake_case**

**Problema:** `updateProductGroup` enviaba payload sin transformar, causando que `parentGroupId` no se guardara como `parent_group_id`.

**Solución:**
```javascript
// Antes (❌)
body: JSON.stringify(data)

// Después (✅)
body: JSON.stringify({
  name: data.name,
  parent_group_id: data.parentGroupId || null,
  is_active: data.isActive !== undefined ? data.isActive : true,
})
```

**Tiempo:** ~15 min

---

### 3. **Implementación de Botones DELETE**

**Problema:** UI no tenía botones de eliminar visibles.

**Solución:**
- ✅ Agregado botón "Eliminar" en `brands/List.vue`
- ✅ Agregado botón "Eliminar" en `product-groups/List.vue`
- ✅ Botón ya existía en `attributes/List.vue`
- ✅ Funciones `deleteBrand()` y `deleteProductGroup()` implementadas
- ✅ Confirmación con `confirm()` antes de eliminar
- ✅ Estilo `.btn-danger` agregado (rojo)

**Tiempo:** ~30 min

---

## 📊 Métricas

| Métrica | Valor |
|---------|-------|
| **Backend** ||
| Comandos creados | 9 (3 delete + 6 existentes) |
| Métodos de servicio | 9 |
| Handlers HTTP | 10 (3 delete + 7 existentes) |
| Rutas registradas | 12 |
| **Frontend** ||
| Componentes de formulario | 3 |
| Páginas de listado | 3 |
| Métodos API añadidos | 9 |
| **Total Líneas** ||
| Backend Go | ~400 (estimado) |
| Frontend Vue | ~1,551 (formularios + listas) |
| **Total** | ~1,951 líneas |

---

## ✅ Testing Manual Completado

### Brands
- ✅ CREATE: Funcional
- ✅ UPDATE: Funcional
- ✅ DELETE: Funcional (con confirmación)

### Product Groups
- ✅ CREATE (raíz): Funcional
- ✅ CREATE (sub): Funcional (corregida transformación de parentGroupId)
- ✅ UPDATE: Funcional
- ✅ DELETE: Funcional (con confirmación)
- ⚠️ Jerarquía: Se guarda correctamente, visualización plana (sin árbol ni indentación)

### Attributes
- ✅ CREATE: Funcional (sin Scope)
- ✅ UPDATE: Funcional
- ✅ DELETE: Funcional (con confirmación)
- ✅ Valores dinámicos: Agregar/eliminar funciona

### Products
- ✅ UPDATE: Funcional (PUT completo con transformaciones correctas)

---

## 🎨 UX/UI Implementado

### Patrones Comunes
- Modal overlay para CREATE/EDIT
- Tabla responsive con hover states
- Botones con estados de loading
- Confirmación antes de eliminar (native `confirm()`)
- Mensajes de error descriptivos
- Empty states informativos

### Estilos
- Botón primario: Amarillo (#f4c430)
- Botón secundario: Blanco con borde
- **Botón peligro: Rojo (#ef4444)** ← Nuevo
- Loading spinner animado
- Badges para contadores
- Code badges para IDs (monospace)

---

## 📝 Decisiones Técnicas

1. **Atributos sin Scope**: Simplifica el dominio y UI. La asignación se hace directamente en productos (no desde marcas/categorías).

2. **Jerarquía plana en UI**: ProductGroups muestra parent por nombre, pero sin indentación visual. Suficiente para MVP.

3. **DELETE sin soft-delete visual**: Se usa confirmación nativa. No implementa lógica de "papelera" por simplicidad de MVP.

4. **Transformaciones explícitas**: Todas las actualizaciones ahora transforman explícitamente camelCase ↔ snake_case para evitar errores silenciosos.

5. **Validaciones frontend-only**: Backend tiene validaciones pero frontend previene la mayoría de errores antes de enviar.

---

## 🚧 Deuda Técnica

### Post-MVP
- [ ] Visualización jerárquica de ProductGroups (árbol con indentación)
- [ ] Soft-delete con papelera virtual
- [ ] Testing unitario frontend (Vitest)
- [ ] Testing E2E (Playwright) para flujos CRUD completos
- [ ] Confirmación modal custom (reemplazar `confirm()` nativo)
- [ ] Validaciones sincronizadas backend-frontend
- [ ] Prevención de eliminación si hay dependencias (ej: Brand con productos)

---

## 🔗 Referencias

- [Sprint 09 Summary](./sprint-09-summary.md)
- [Tarea 01: Product List UI](./01-product-ui-list-implementation.md)
- [Tarea 02: Product Detail UI](./02-product-ui-detail-implementation.md)
- [Tarea 03: Product Create Forms](./03-product-ui-create-forms-implementation.md)
- [Dominio Product - ADR-015](../../../architecture/adrs/ADR-015-product-module-architecture.md)

---

## 🎯 Conclusión

**Estado:** ✅ COMPLETADO

Todas las funcionalidades CRUD están **implementadas y testeadas manualmente**. Los issues identificados fueron corregidos:
- ✅ Atributos sin Scope (refactor arquitectónico)
- ✅ Botones DELETE agregados
- ✅ Jerarquía ProductGroups funcional (guardado correcto)
- ✅ Transformaciones camelCase↔snake_case corregidas

El módulo Product Master Data está **listo para MVP**.

---

**Siguiente:** Pricing Integration Panel (Tarea 09-04) o avanzar a módulo Sales dado timeline ajustado.

**Timeline:** 9 días restantes para completar Sales + MES antes del 23 de febrero.
