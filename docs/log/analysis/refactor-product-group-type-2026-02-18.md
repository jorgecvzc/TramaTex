# Refactor: Clasificación de Grupos de Productos (Tangibles vs Servicios)

**Fecha:** 2026-02-18  
**Tipo:** Feature Refactor  
**Alcance:** Todas las capas (DDD, Backend, Frontend, UI, Tests, Documentación)

---

## 📋 Resumen Ejecutivo

Se ha implementado un refactor completo para agregar clasificación de grupos de productos, permitiendo distinguir entre **productos tangibles** (físicos) y **servicios**. Este cambio afecta todas las capas arquitectónicas siguiendo los principios de Clean Architecture y DDD.

### Valor de Negocio

- **Diferenciación clara** entre productos físicos y servicios en el sistema
- **Mejor organización** del catálogo de productos
- **Preparación** para lógica de negocio específica por tipo (inventario físico vs servicios)
- **Base** para futuras funcionalidades (gestión de stock solo para tangibles, calendarios para servicios, etc.)

---

## 🏗️ Cambios Implementados

### 1. Domain Layer (DDD) ✅

**Archivo:** `apps/tramatex-api/internal/product/domain/product.go`

**Cambios:**
```go
// Nuevo enum ProductGroupType
type ProductGroupType string

const (
    ProductGroupTypeTangible ProductGroupType = "TANGIBLE" // Productos físicos
    ProductGroupTypeService  ProductGroupType = "SERVICE"  // Servicios
)

// ProductGroup entity actualizada
type ProductGroup struct {
    ID            uuid.UUID
    Name          string
    Type          ProductGroupType // ⭐ NUEVO CAMPO
    ParentGroupID *uuid.UUID
    IsActive      bool
}

// Constructor actualizado
func NewProductGroup(name string, groupType ProductGroupType, parentID *uuid.UUID, isActive bool) (*ProductGroup, error)

// Nuevo método de validación
func (pgt ProductGroupType) IsValid() bool

// Nuevo método de actualización
func (pg *ProductGroup) UpdateType(groupType ProductGroupType) error
```

**Validaciones:**
- El tipo debe ser `TANGIBLE` o `SERVICE` (validado en `IsValid()`)
- Constructor rechaza tipos inválidos
- UpdateType valida antes de actualizar

---

### 2. Persistence Layer ✅

#### 2.1 Migration

**Archivo:** `apps/tramatex-api/migrations/020_add_product_group_type.sql`

```sql
-- Crear enum type
CREATE TYPE product_group_type AS ENUM ('TANGIBLE', 'SERVICE');

-- Agregar columna con default
ALTER TABLE "product_groups" 
ADD COLUMN IF NOT EXISTS "group_type" product_group_type NOT NULL DEFAULT 'TANGIBLE';

-- Índice para filtros
CREATE INDEX IF NOT EXISTS "idx_product_groups_group_type" 
ON "product_groups" ("group_type");

-- Comentario descriptivo
COMMENT ON COLUMN "product_groups"."group_type" IS 
'Classification of product group: TANGIBLE for physical products, SERVICE for service-based products';
```

#### 2.2 Data Model

**Archivo:** `apps/tramatex-api/internal/product/infrastructure/persistence/product_group_data_model.go`

```go
type ProductGroupDataModel struct {
    gorm.Model
    ID            uuid.UUID `gorm:"type:uuid;primary_key;"`
    Name          string    `gorm:"uniqueIndex;not null"`
    Type          string    `gorm:"type:product_group_type;not null;default:TANGIBLE"` // ⭐ NUEVO
    ParentGroupID *uuid.UUID
    IsActive      bool      `gorm:"not null;default:true"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// Conversiones actualizadas
func (pg *ProductGroupDataModel) ToDomain() *domain.ProductGroup {
    return &domain.ProductGroup{
        ID:            pg.ID,
        Name:          pg.Name,
        Type:          domain.ProductGroupType(pg.Type), // ⭐ Conversión de tipo
        ParentGroupID: pg.ParentGroupID,
        IsActive:      pg.IsActive,
    }
}
```

#### 2.3 Seed Data

**Archivo:** `apps/tramatex-api/migrations/016_seed_product_master_data.sql`

```sql
-- Todos los grupos seed ahora especifican el tipo explícitamente
INSERT INTO product_groups (id, name, group_type, is_active, created_at, updated_at) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Calzado Deportivo', 'TANGIBLE', true, NOW(), NOW()),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Ropa Deportiva', 'TANGIBLE', true, NOW(), NOW()),
...
```

---

### 3. Application Layer ✅

**Archivos:** 
- `apps/tramatex-api/internal/product/application/commands.go`
- `apps/tramatex-api/internal/product/application/product_service.go`

#### 3.1 Commands

```go
// CreateProductGroupCommand
type CreateProductGroupCommand struct {
    ActorID  string
    Name     string
    Type     string // ⭐ NUEVO: "TANGIBLE" o "SERVICE"
    ParentID *uuid.UUID
    IsActive bool
}

// UpdateProductGroupCommand
type UpdateProductGroupCommand struct {
    ActorID  string
    ID       uuid.UUID
    Name     *string
    Type     *string // ⭐ NUEVO: Opcional para updates
    ParentID *uuid.UUID
    IsActive *bool
}
```

#### 3.2 DTOs

```go
type ProductGroupDTO struct {
    ID            uuid.UUID  `json:"id"`
    Name          string     `json:"name"`
    Type          string     `json:"type"` // ⭐ NUEVO
    ParentGroupID *uuid.UUID `json:"parent_group_id,omitempty"`
    IsActive      bool       `json:"isActive"`
}
```

#### 3.3 Service Methods

**CreateProductGroup:**
```go
func (s *ProductService) CreateProductGroup(ctx context.Context, cmd CreateProductGroupCommand) (*ProductGroupDTO, error) {
    // Parse y validar tipo
    groupType := domain.ProductGroupType(cmd.Type)
    if !groupType.IsValid() {
        return nil, domain.NewValidationError("invalid group type: must be TANGIBLE or SERVICE")
    }
    
    // Crear con tipo
    group, err := domain.NewProductGroup(cmd.Name, groupType, cmd.ParentID, cmd.IsActive)
    // ... resto del método
}
```

**UpdateProductGroup:**
```go
func (s *ProductService) UpdateProductGroup(ctx context.Context, cmd UpdateProductGroupCommand) (*ProductGroupDTO, error) {
    // ... obtener grupo existente
    
    // Actualizar tipo si viene en el command
    if cmd.Type != nil {
        groupType := domain.ProductGroupType(*cmd.Type)
        if err := group.UpdateType(groupType); err != nil {
            return nil, domain.WrapValidation("failed to update product group type", err)
        }
    }
    // ... resto del método
}
```

---

### 4. Frontend Layer ✅

#### 4.1 Types

**Archivo:** `apps/frontend/src/types/product.ts`

```typescript
// Nuevo tipo exportado
export type ProductGroupType = 'TANGIBLE' | 'SERVICE'

export interface ProductGroup {
  id: string
  name: string
  type: ProductGroupType // ⭐ NUEVO CAMPO
  code: string
  description: string | null
  parent_id: string | null
  is_active: boolean
  created_at: string
  updated_at: string
}
```

#### 4.2 API Service

**Archivo:** `apps/frontend/src/services/productApi.ts`

```typescript
interface ProductGroupUI {
  id: string
  name: string
  type: string // ⭐ NUEVO: TANGIBLE or SERVICE
  is_active: boolean
  parent_group_id?: string
  description?: string
}

// Request interfaces actualizadas
async createProductGroup(data: {
  id?: string
  name: string
  type?: string // ⭐ NUEVO
  parentGroupId?: string
  isActive?: boolean
}): Promise<any>

async updateProductGroup(id: string, data: {
  name?: string
  type?: string // ⭐ NUEVO
  parentGroupId?: string | null
  isActive?: boolean
}): Promise<any>
```

**Transformaciones:**
```typescript
// En listProductGroups y getProductGroup
const groups: ProductGroupUI[] = rawGroups.map((g: any) => ({
  id: g.id,
  name: g.name,
  type: g.type || 'TANGIBLE', // ⭐ Default a TANGIBLE
  is_active: g.isActive,
  parent_group_id: g.parentGroupId,
  description: g.description,
}))
```

#### 4.3 UI Components

##### ProductGroupForm.vue

**Template:**
```vue
<!-- Nuevo campo de selección de tipo con radio buttons -->
<div class="form-group">
  <label>Tipo de categoría <span class="required">*</span></label>
  <div class="radio-group">
    <label class="radio-label">
      <input type="radio" v-model="formData.type" value="TANGIBLE" name="groupType" />
      <div class="radio-content">
        <span class="radio-title">🔧 Productos Tangibles</span>
        <span class="radio-description">Productos físicos: calzado, ropa, accesorios, equipamiento</span>
      </div>
    </label>
    <label class="radio-label">
      <input type="radio" v-model="formData.type" value="SERVICE" name="groupType" />
      <div class="radio-content">
        <span class="radio-title">⚙️ Servicios</span>
        <span class="radio-description">Servicios profesionales: consultoría, mantenimiento, instalación</span>
      </div>
    </label>
  </div>
</div>
```

**Script:**
```typescript
const formData = reactive({
  name: props.productGroup?.name || '',
  type: props.productGroup?.type || 'TANGIBLE', // ⭐ NUEVO con default
  parentGroupId: props.productGroup?.parent_group_id || '',
  isActive: props.productGroup?.is_active ?? true
})

// Validación actualizada
function validate() {
  let isValid = true
  
  if (!formData.name.trim()) {
    errors.name = 'El nombre es obligatorio'
    isValid = false
  }
  
  if (!formData.type || !['TANGIBLE', 'SERVICE'].includes(formData.type)) {
    errors.type = 'Debe seleccionar un tipo válido'
    isValid = false
  }
  
  return isValid
}

// Payload incluye tipo
const payload = {
  name: formData.name.trim(),
  type: formData.type, // ⭐ NUEVO
  parentGroupId: formData.parentGroupId || null,
  isActive: formData.isActive
}
```

**Estilos:**
```css
.radio-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.radio-label {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  background: #f8fafc;
}

.radio-label:hover {
  border-color: #94a3b8;
  background: #ffffff;
}

.radio-label:has(input:checked) {
  border-color: #3b82f6;
  background: #eff6ff;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}
```

##### List.vue

**Template:**
```vue
<table class="data-table">
  <thead>
    <tr>
      <th>Nombre</th>
      <th>Tipo</th> <!-- ⭐ NUEVA COLUMNA -->
      <th>ID</th>
      <th>Categoría Padre</th>
      <th class="text-center">Productos</th>
      <th class="text-right">Acciones</th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="group in groups" :key="group.id">
      <td><strong>{{ group.name }}</strong></td>
      <td>
        <!-- ⭐ NUEVO: Badge con icono y color según tipo -->
        <span v-if="group.type === 'TANGIBLE'" class="badge badge-tangible">🔧 Tangible</span>
        <span v-else-if="group.type === 'SERVICE'" class="badge badge-service">⚙️ Servicio</span>
        <span v-else class="badge">{{ group.type }}</span>
      </td>
      <!-- ... resto de columnas -->
    </tr>
  </tbody>
</table>
```

**Estilos:**
```css
.badge-tangible {
  background: #dbeafe;
  color: #1e40af;
  border: 1px solid #93c5fd;
}

.badge-service {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fde68a;
}
```

---

### 5. Tests ✅

#### 5.1 Frontend Tests

**Archivo:** `apps/frontend/src/__tests__/unit/productApi.test.ts`

**Cambios:**
```typescript
describe('listProductGroups', () => {
  it('should list product groups successfully', async () => {
    const mockGroups = [
      {
        id: 'group-001',
        name: 'Electronics',
        type: 'TANGIBLE', // ⭐ NUEVO
        isActive: true,
        parentGroupId: null,
        description: 'Electronic products',
      },
    ]
    // ... test actualizado para verificar type
    expect(result.data[0].type).toBe('TANGIBLE')
  })
})

describe('createProductGroup', () => {
  it('should create a SERVICE type product group', async () => {
    const newGroup = {
      name: 'Consulting Services',
      type: 'SERVICE', // ⭐ Test para tipo SERVICE
      isActive: true,
    }
    // ... test verifica que el tipo se envía correctamente en el body
    const body = JSON.parse(fetchCall[1].body)
    expect(body.type).toBe('SERVICE')
  })
})
```

**Cobertura:**
- ✅ Listado incluye campo type
- ✅ Creación con tipo TANGIBLE
- ✅ Creación con tipo SERVICE
- ✅ Actualización de tipo
- ✅ Validación de tipos inválidos

---

## 📊 Cobertura de Cambios

| Capa | Archivos Modificados | Estado |
|------|---------------------|---------|
| **Database** | 2 archivos (migration + seed) | ✅ Completo |
| **Domain** | 1 archivo (product.go) | ✅ Completo |
| **Persistence** | 1 archivo (data model) | ✅ Completo |
| **Application** | 2 archivos (commands + service) | ✅ Completo |
| **Frontend Types** | 1 archivo (product.ts) | ✅ Completo |
| **Frontend API** | 1 archivo (productApi.ts) | ✅ Completo |
| **Frontend UI** | 2 archivos (Form + List) | ✅ Completo |
| **Tests** | 1 archivo (productApi.test.ts) | ✅ Completo |

**Total:** 11 archivos modificados/creados

---

## 🔄 Migración de Datos Existentes

### Estrategia

La migration creada incluye un `DEFAULT 'TANGIBLE'` que garantiza:

1. **Datos existentes**: Se marcan automáticamente como `TANGIBLE`
2. **Nuevos registros sin tipo**: Defaultean a `TANGIBLE` (productos físicos son más comunes)
3. **Sin ruptura**: El sistema sigue funcionando sin cambios manuales

### Comando para revisar datos tras migration

```sql
-- Ver distribución de tipos
SELECT group_type, COUNT(*) FROM product_groups GROUP BY group_type;

-- Ver grupos específicos
SELECT id, name, group_type FROM product_groups ORDER BY name;

-- Actualizar manualmente si fuera necesario
UPDATE product_groups SET group_type = 'SERVICE' WHERE name LIKE '%Service%';
```

---

## ✅ Checklist de Validación

### Backend
- [x] Enum `ProductGroupType` creado y validado
- [x] Constructor `NewProductGroup` acepta tipo
- [x] Método `UpdateType` implementado con validación
- [x] Data model incluye campo `Type`
- [x] Migration crea enum PostgreSQL y columna
- [x] Seed data actualizado con tipos explícitos
- [x] Commands incluyen campo `Type`
- [x] DTOs exponen campo `type` en JSON
- [x] Service valida tipos en create/update
- [x] Handlers HTTP automáticamente parsean tipo del JSON

### Frontend
- [x] Type `ProductGroupType` exportado
- [x] Interface `ProductGroup` incluye campo `type`
- [x] API service maneja campo `type` en requests/responses
- [x] Form captura tipo con radio buttons
- [x] Form valida que tipo sea TANGIBLE o SERVICE
- [x] Form muestra descripciones claras de cada opción
- [x] List muestra tipo con badge visual
- [x] Badges tienen estilos distintivos por tipo
- [x] Tests cubren creación con ambos tipos
- [x] Tests verifican campo type en responses

### UX/UI
- [x] Campo tipo visualmente destacado con iconos (🔧 ⚙️)
- [x] Descripciones claras de cada opción
- [x] Selección por defecto (TANGIBLE) para facilitar UX
- [x] Estado hover/activo visualmente distintivo
- [x] Badges con colores semánticos (azul=tangible, amarillo=servicio)
- [x] Columna tipo agregada a tabla de lista

---

## 🚀 Próximos Pasos Recomendados

### Immediate (Opcional)
1. **Ejecutar migration**: Iniciar servidor backend para que se aplique migration 020
2. **Validar UI**: Abrir formulario de creación de grupo y verificar radio buttons
3. **Crear ejemplos**: Agregar 1 grupo tipo SERVICE para testing

### Futuro (Potenciales mejoras)
1. **Lógica diferenciada**: Implementar restricciones por tipo
   - Grupos TANGIBLE: Habilitar gestión de inventario
   - Grupos SERVICE: Habilitar calendarios/reservas
2. **Filtros**: Agregar filtro por tipo en listado de grupos
3. **Validaciones de negocio**: 
   - Productos físicos solo en grupos TANGIBLE
   - Servicios solo en grupos SERVICE
4. **Reportes**: Separar estadísticas por tipo de grupo
5. **Backend tests**: Agregar unit tests para ProductGroup con tipos

---

## 📝 Notas Técnicas

### Decisiones de Diseño

1. **Enum vs Boolean**: Se eligió enum por:
   - Extensibilidad futura (puede agregar MIXED, DIGITAL, etc.)
   - Claridad semántica (TANGIBLE vs SERVICE es más expresivo que `is_tangible`)
   - Validación explícita en base de datos

2. **Default TANGIBLE**: Se asume que:
   - Productos físicos son más comunes en contexto ERP
   - Datos existentes son principalmente productos físicos
   - Minimiza fricción en migración

3. **Radio buttons vs Select**: Se eligieron radio buttons porque:
   - Solo 2 opciones (no requiere desplegable)
   - Visualmente más claro
   - Permite mostrar descripciones completas
   - UX más amigable para decisión binaria

4. **Validación en múltiples capas**:
   - Domain: `IsValid()` method
   - Application: Pre-validación antes de llamar domain
   - Database: Constraint de enum PostgreSQL
   - Frontend: Validación en form ante de enviar

---

## 🔗 Referencias

- **ADR-002**: Clean Architecture + DDD (base arquitectónica)
- **Migration 009**: Creación original de tabla product_groups
- **Migration 016**: Seed data de master data
- **docs/modules/product/**: Documentación del módulo Product

---

**Autor:** GitHub Copilot (Claude Sonnet 4.5)  
**Revisado por:** — (Pendiente)  
**Estado:** ✅ Implementación Completa
