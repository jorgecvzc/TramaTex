# Testing E2E - Sistema de Variantes de Producto

**Fecha:** 2025-02-05  
**Sprint:** 09  
**Módulo:** Product  
**Componentes Testeados:** VariantTable, VariantFormModal, VariantSelector

---

## ✅ Preparación del Entorno

- [x] Frontend corriendo en http://localhost:5173/
- [x] Backend API corriendo en http://localhost:8080/ (contenedor tramatex_api: healthy)
- [x] Base de datos PostgreSQL activa (contenedor tramatex_db: healthy)
- [x] Compilación sin errores de los nuevos componentes Vue

---

## 🧪 Plan de Pruebas

### 1. VariantTable - Visualización de Variantes

**URL:** http://localhost:5173/products/{productId}  
**Objetivo:** Verificar que las variantes se muestran correctamente con el nuevo diseño

#### Casos de Prueba:

1.1. **Ver variantes activas**
   - [ ] Navegar al detalle del producto con ID: `5133ae73-aebd-4b7c-8b5b-8ce550435bc3`
   - [ ] Ir a la pestaña "Variantes"
   - [ ] Verificar que se listan las variantes existentes
   - [ ] Confirmar que las variantes activas NO tienen borde rojo
   - [ ] Verificar que los pills muestran estado (PROVISIONAL/CONFIRMED) y "Activo"
   - [ ] Verificar que los precios se cargan (mockup por ahora)

1.2. **Ver variantes inactivas**
   - [ ] Si no hay variantes inactivas, crear una con is_active=false
   - [ ] Verificar que la fila tiene fondo rojo suave: `rgba(220, 38, 38, 0.04)`
   - [ ] Verificar que tiene borde izquierdo rojo: `3px solid rgba(220, 38, 38, 0.4)`
   - [ ] Verificar que el SKU tiene fondo/texto rojo
   - [ ] Verificar que el pill "Inactivo" es rojo (`#dc2626`)
   - [ ] Verificar que el texto general está en gris apagado (`#94a3b8`)

1.3. **Hover sobre variantes**
   - [ ] Hover sobre variante activa → fondo `#f8fafc`
   - [ ] Hover sobre variante inactiva → fondo `rgba(220, 38, 38, 0.08)`

1.4. **Header actions**
   - [ ] Verificar que hay dos botones: "➕ Añadir Variante" y "🔄 Actualizar"
   - [ ] Ambos botones centrados en la derecha

---

### 2. VariantFormModal - Crear Variante Manualmente

**URL:** Detail de producto → Tab Variantes → Botón "Añadir Variante"  
**Objetivo:** Verificar que se pueden crear variantes JIT con metadatos

#### Casos de Prueba:

2.1. **Abrir modal en modo creación**
   - [ ] Click en "➕ Añadir Variante"
   - [ ] Modal se abre con título "Crear Variante"
   - [ ] Se muestran todos los atributos del producto con dropdowns
   - [ ] Campos de metadatos: barcode, baseCost, status, isActive

2.2. **Seleccionar combinación de atributos**
   - [ ] Seleccionar un valor para cada atributo (ej: Talla=M, Color=Azul)
   - [ ] Verificar que aparece "SKU Generado" con formato correcto:  
        `{PRODUCT_SKU}-{ATTR1_CODE}.{VAL1_CODE}-{ATTR2_CODE}.{VAL2_CODE}`
   - [ ] Ejemplo: `TST001-SIZE.M-COLOR.BLUE`

2.3. **Completar metadatos opcionales**
   - [ ] Ingresar código de barras: `7501234567890`
   - [ ] Ingresar costo base: `15.50`
   - [ ] Cambiar estado a "Confirmado"
   - [ ] Desmarcar "Variante activa" (para crear una inactiva)

2.4. **Validaciones del formulario**
   - [ ] Botón "Crear Variante" deshabilitado si faltan atributos
   - [ ] Botón "Crear Variante" habilitado cuando todos los atributos están seleccionados
   - [ ] Campo "required" marcado con asterisco rojo

2.5. **Crear variante**
   - [ ] Click en "Crear Variante"
   - [ ] Botón muestra "Guardando..." durante el proceso
   - [ ] Modal se cierra automáticamente al éxito
   - [ ] Tabla de variantes se refresca (evento @refresh)
   - [ ] La nueva variante aparece en la lista

2.6. **Manejo de errores**
   - [ ] Intentar crear variante duplicada (misma combinación de atributos)
   - [ ] Verificar que se muestra mensaje de error en banner rojo
   - [ ] Modal NO se cierra en caso de error

---

### 3. VariantFormModal - Editar Variante Existente

**URL:** Detail de producto → Tab Variantes → Botón ✏️ en fila de variante  
**Objetivo:** Verificar que solo se pueden editar metadatos (combinación inmutable)

#### Casos de Prueba:

3.1. **Abrir modal en modo edición**
   - [ ] Click en botón ✏️ de una variante existente
   - [ ] Modal se abre con título "Editar Variante"
   - [ ] SKU de la variante se muestra (no editable)
   - [ ] Mensaje: "La combinación de atributos no puede modificarse"
   - [ ] NO se muestran dropdowns de atributos

3.2. **Editar metadatos**
   - [ ] Cambiar código de barras
   - [ ] Cambiar costo base
   - [ ] Cambiar estado (PROVISIONAL → CONFIRMED o viceversa)
   - [ ] Toggle "Variante activa"

3.3. **Guardar cambios**
   - [ ] Click en "Actualizar"
   - [ ] Modal se cierra
   - [ ] Variante actualizada en la tabla con nuevos valores
   - [ ] Si se desactivó, ahora tiene estilo de inactiva

3.4. **Verificar inmutabilidad de combinación**
   - [ ] Confirmar que NO hay forma de cambiar los atributos desde este modal
   - [ ] Para crear otra combinación, se debe usar "Añadir Variante"

---

### 4. VariantSelector - Modalidad A (Por Atributos JIT)

**Objetivo:** Verificar selección interactiva con creación JIT

#### Casos de Prueba:

4.1. **Selección de producto**
   - [ ] Si productId NO está fijado, se muestra dropdown de productos
   - [ ] Seleccionar un producto activo
   - [ ] Dropdowns de atributos se cargan dinámicamente

4.2. **Selección de atributos**
   - [ ] Seleccionar valores para cada atributo
   - [ ] Verificar que aparece "SKU Generado" con formato correcto
   - [ ] Botón "Buscar o Crear Variante" se habilita solo cuando TODOS los atributos están seleccionados

4.3. **Buscar o crear variante (JIT)**
   - [ ] Click en "✨ Buscar o Crear Variante"
   - [ ] Si variante existe: se muestra tarjeta con datos existentes
   - [ ] Si NO existe: se crea automáticamente con estado PROVISIONAL
   - [ ] Tarjeta verde aparece con:
        - Badge de estado (Provisional/Confirmado)
        - Badge "Inactivo" si is_active=false
        - SKU de la variante
        - Código de barras (si existe)
        - Botones: "✓ Confirmar Selección" y "Cancelar"

4.4. **Advertencias de estado**
   - [ ] Si variante está inactiva, banner amarillo/rojo: "⚠️ Esta variante está marcada como inactiva"
   - [ ] Usuario puede decidir si confirmar o no

4.5. **Confirmar selección**
   - [ ] Click en "✓ Confirmar Selección"
   - [ ] Evento `variant-selected` emitido con:
        - variantId: string (UUID)
        - variant: objeto completo

---

### 5. VariantSelector - Modalidad B (Por SKU)

**Objetivo:** Verificar búsqueda directa sin JIT

#### Casos de Prueba:

5.1. **Cambiar a modo SKU**
   - [ ] Click en tab "🔍 Por SKU"
   - [ ] Input de búsqueda aparece
   - [ ] Labels: "SKU de Variante" con asterisco rojo (required)

5.2. **Buscar por SKU**
   - [ ] Ingresar SKU válido existente: `TST001-SIZE.M-COLOR.BLUE`
   - [ ] Click en "🔍 Buscar" o presionar Enter
   - [ ] Spinner "Procesando..." durante la búsqueda
   - [ ] Si existe: tarjeta verde aparece con datos de la variante
        - Nombre del producto
        - SKU
        - Código de barras
        - Valores de atributos (tags: "Talla: M", "Color: Azul")
        - Botones: "✓ Confirmar Selección" y "Cancelar"

5.3. **Búsqueda sin resultados**
   - [ ] Ingresar SKU inexistente: `FAKE-SKU-123`
   - [ ] Click en "Buscar"
   - [ ] Banner rojo: "No se encontró ninguna variante con ese SKU"
   - [ ] NO se crea variante automáticamente (modo B no tiene JIT)

5.4. **Confirmar selección**
   - [ ] Click en "✓ Confirmar Selección"
   - [ ] Evento `variant-selected` emitido igual que en Modalidad A

---

### 6. Cambio entre Modalidades

**Objetivo:** Verificar que los cambios de modo limpian el estado

#### Casos de Prueba:

6.1. **Cambio con datos en progreso**
   - [ ] En modo "Por Atributos", seleccionar algunos atributos
   - [ ] Cambiar a modo "Por SKU"
   - [ ] Verificar que los atributos seleccionados se limpian
   - [ ] Volver a modo "Por Atributos"
   - [ ] Verificar que dropdowns están vacíos (estado limpio)

6.2. **Cambio con variante seleccionada**
   - [ ] Buscar y seleccionar una variante en modo "Por SKU"
   - [ ] Cambiar a modo "Por Atributos"
   - [ ] Verificar que la variante seleccionada se limpia
   - [ ] NO se emite evento `variant-selected` automáticamente

---

## 📊 Resultados Esperados

### Componentes Creados/Modificados
- ✅ `VariantTable.vue` - Actualizado con estilos para inactivos + integración de modal
- ✅ `VariantFormModal.vue` - Nuevo componente para crear/editar variantes
- ✅ `VariantSelector.vue` - Nuevo componente reutilizable con 2 modalidades
- ✅ `Detail.vue` - Actualizado para pasar productSku a VariantTable

### Documentación
- ✅ UC-S-021 documentado en `docs/modules/sales/use-cases.md` (113 líneas)

### Decisiones Arquitectónicas Aplicadas
1. ✅ Hybrid approach: Dos modalidades en VariantSelector (dropdown + SKU search)
2. ✅ Combinación inmutable: Solo metadatos editables en VariantFormModal
3. ✅ Inactivos siempre visibles: Highlight rojo en VariantTable
4. ✅ Documentación centralizada: UC-S-021 en sales/use-cases.md
5. ✅ Roles requeridos: JIT creation solo para commercial/admin (lógica en backend)

---

## 🚀 Próximos Pasos (Fuera de este Sprint)

1. **Integración real con Pricing:**
   - Reemplazar mockup de precios en VariantTable
   - Conectar con API de Pricing para obtener precios base reales

2. **Vista de detalle de variante:**
   - Implementar modal/página para ver todos los detalles de una variante
   - Botón 👁️ en VariantTable

3. **Uso en Sales:**
   - Integrar VariantSelector en formulario de Order Items
   - Implementar lógica de confirmación de variantes PROVISIONAL→CONFIRMED al añadir a orden

4. **Testing Automatizado:**
   - E2E tests con Playwright para todos los flujos
   - Unit tests para lógica de componentes (Vitest)

5. **Optimizaciones:**
   - Cache de productos en VariantSelector
   - Paginación en VariantTable si hay muchas variantes
   - Búsqueda incremental en modo SKU

---

## ✍️ Notas del Desarrollador

**Producto de prueba usado:**
- ID: `5133ae73-aebd-4b7c-8b5b-8ce550435bc3`
- SKU: `TST001`
- Nombre: Producto de Prueba
- Atributos: Talla (SIZE), Color (COLOR)

**Comandos útiles:**
```bash
# Frontend
cd apps/frontend
npm run dev

# Backend logs
docker logs -f tramatex_api

# Reiniciar contenedores si es necesario
docker-compose -f docker/docker-compose.yml restart

# Ver contenedores activos
docker ps
```

---

**Estado Final:** ✅ LISTO PARA TESTING MANUAL  
**Responsable:** Equipo de Desarrollo  
**Próxima Revisión:** Al completar Sprint 09
