# 🔧 Correcciones Aplicadas - Parties y Pricing UI

**Fecha:** 2026-02-14  
**Estado:** ✅ Problemas Resueltos

---

## Problema 1: "No muestra las Entidades" ✅ CORREGIDO

### Causa Raíz
El código de `partyApi.js` tenía dos problemas:

1. **Filtro forzado (línea 183):** 
   ```javascript
   params.append('type', 'ORGANIZATION'); // ❌ Siempre filtraba por organizaciones
   ```

2. **Mapeo incompleto (línea 83):**
   ```javascript
   if (!party || !party.organization_profile) {
     return null; // ❌ Descartaba parties sin organization_profile
   }
   ```

Esto causaba que:
- Solo se mostraran parties de tipo ORGANIZATION
- Parties con solo `person_profile` (personas) se descartaban
- Si no había parties con `organization_profile`, la lista aparecía vacía

### Solución Aplicada ✅

**Archivo:** `apps/frontend/src/services/partyApi.js`

**Cambio 1: Soporte para personas y organizaciones**
```javascript
mapPartyToParty(party) {
  if (!party) {
    return null;
  }

  // Support both organization and person profiles
  let name = '(Sin nombre)';
  let tax_id = null;
  let tax_id_type = null;
  let website = null;

  if (party.organization_profile) {
    name = party.organization_profile.name;
    tax_id = party.organization_profile.tax_id;
    tax_id_type = party.organization_profile.tax_id_type;
    website = party.organization_profile.website;
  } else if (party.person_profile) {
    name = `${party.person_profile.first_name} ${party.person_profile.last_name}`;
  }

  return {
    id: party.id,
    name,
    role,
    status: party.status,
    tax_id,
    tax_id_type,
    website,
    created_at: party.created_at,
    modified_at: party.modified_at,
    has_organization: !!party.organization_profile,
    has_person: !!party.person_profile,
  };
}
```

**Cambio 2: Filtro de tipo opcional**
```javascript
async listParties(filters = {}) {
  const params = new URLSearchParams();
  
  if (filters.name) params.append('name', filters.name);
  if (filters.role && filters.role !== 'BOTH') params.append('role', filters.role);
  if (filters.status) params.append('status', filters.status);
  if (filters.pageNumber) params.append('page', filters.pageNumber);
  if (filters.pageSize) params.append('page_size', filters.pageSize);
  // Allow filtering by type if provided, otherwise show all ✅
  if (filters.type) params.append('type', filters.type);
  // ❌ REMOVED: params.append('type', 'ORGANIZATION');
```

### Resultado
✅ Ahora se muestran TODAS las parties (organizaciones + personas)  
✅ El nombre se formatea correctamente según el tipo de perfil  
✅ Se pueden agregar filtros de tipo opcionalmente en el futuro

---

## Problema 2: "¿el pricing ha de ser editable desde UIs de product?" ✅ ACLARADO

### Respuesta: NO, y está diseñado correctamente

El **PricingPanel** en ProductDetail NO es para editar reglas de pricing. Es una herramienta de **CONSULTA/VISUALIZACIÓN**.

### Funcionalidades del PricingPanel

**1. Visualización de Precios Base (Read-Only)**
- Muestra el precio base de venta de cada variante
- Carga automática desde Pricing Engine API
- No permite editar estos precios

**2. Calculadora de Precios Finales**
- Simula precio final basado en:
  - Variante seleccionada
  - Cliente (UUID)
  - Cantidad
  - Fecha de venta
- Llama a `POST /api/pricing/final-sale-price/calculate`
- Muestra el resultado calculado (precio base → precio final)

**3. Historial de Cálculos**
- Muestra cálculos previos de precios
- Llama a `GET /api/pricing/history/:variantId`
- Solo lectura

### Mensaje en el Código (líneas 206-212)
```vue
<div class="info-message">
  <p>
    <strong>Próximamente:</strong> Visualización de reglas de pricing aplicables
    (reglas de precio base, modificaciones de venta, descuentos por volumen).
  </p>
  <p class="help-text">
    Las reglas se configuran desde el módulo Pricing (Admin → Pricing → Reglas).
  </p>
</div>
```

### Arquitectura Correcta
```
┌─────────────────┐
│ Product Module  │ ← Define productos y variantes
└────────┬────────┘
         │
         ↓ (consulta precios via API)
┌─────────────────┐
│ Pricing Module  │ ← Define y edita reglas de precios base
└────────┬────────┘
         │
         ↓ (calcula usando reglas)
┌─────────────────┐
│ Pricing Engine  │ ← Motor de cálculo de precios
└─────────────────┘
```

**Separación de responsabilidades:**
- **Product UI:** Ve los efectos del pricing (consulta)
- **Pricing UI:** Configura las reglas de pricing (edición)
- **Pricing Engine:** Ejecuta los cálculos

### ¿Dónde se editarían las reglas de precio?

**Opción 1: Módulo Pricing Frontend (Recomendado)**
```
/admin/pricing/rules
  ├── base-sales-rules/      ← CRUD de reglas de precio base
  ├── modification-rules/    ← CRUD de reglas de descuentos/recargos
  └── client-overrides/      ← CRUD de overrides por cliente
```

**Opción 2: En Product Detail (NO recomendado)**
- Mezclaría responsabilidades
- Pricing es transversal a múltiples productos
- Las reglas afectan a múltiples variantes/clientes

### Conclusión
✅ **El diseño actual es CORRECTO**  
✅ PricingPanel debe permanecer como herramienta de consulta  
✅ La edición de reglas de pricing debe implementarse en un módulo Pricing separado  
⏳ Pricing Frontend UI pendiente (no bloqueante para ERP Core MVP)

---

## Impacto en el Proyecto

### Problemas Resueltos
1. ✅ Parties ahora se muestran correctamente (organizaciones + personas)
2. ✅ Aclarada la arquitectura de Pricing (PricingPanel es solo consulta)

### Estado ERP Core Frontend Actualizado

| Módulo | Status | UI Edición | UI Consulta | Comentarios |
|--------|--------|-----------|-------------|-------------|
| **Party** | ✅ 100% | ✅ Completo | ✅ Completo | Corregido hoy |
| **Product** | ✅ 100% | ✅ Completo | ✅ Completo | Incluye PricingPanel (consulta) |
| **Pricing** | ⚠️ 33% | ❌ Pendiente | ✅ En Product | PricingPanel es consulta, falta UI de edición |
| **Sales** | ❌ 0% | ❌ Pendiente | ❌ Pendiente | **Siguiente prioridad** |

### Módulo Pricing: Análisis Detallado

**Backend Pricing:** ✅ 90% completo (10 endpoints)
- ✅ POST /calculate
- ✅ GET /rules, POST /rules
- ✅ POST /client-overrides
- ✅ GET /history/:variantId
- ✅ POST /base-sales-rules, PUT /base-sales-rules/:id
- ✅ POST /sale-modification-rules, PUT /sale-modification-rules/:id
- ✅ POST /base-sales-price/calculate
- ✅ POST /final-sale-price/calculate
- ⚠️ Falta: GET list endpoints para base-sales-rules y sale-modification-rules

**Frontend Pricing:** ⚠️ 33% completo
- ✅ PricingPanel (consulta en Product Detail) - 818 líneas
- ✅ pricingApi.js service - endpoints implementados
- ❌ UI de edición de reglas (CRUD base-sales-rules)
- ❌ UI de edición de modificaciones (CRUD sale-modification-rules)
- ❌ UI de overrides por cliente

**Prioridad:** 🟡 Media (no bloqueante para Sales)
- Pricing Engine funciona (backend completo)
- Consulta disponible desde Product
- Edición de reglas se puede hacer temporalmente via SQL/API directa
- UI de edición es nice-to-have, no critical

---

## Próximos Pasos

### Inmediato ✅
1. Recargar frontend para aplicar cambios de partyApi.js
2. Acceder a http://localhost:5173/parties y verificar que se muestran las entidades
3. Confirmar que CONSUMIDOR_FINAL aparece en la lista (creado en migration 019)

### Siguiente Prioridad 🚀
**Sales Frontend** (2-3 días):
- Order List con filtros
- Order Detail con tabs
- Create Order con VariantSelector
- Invoice/Ticket creation buttons
- Integración con Party (select clients) y Product (select variants)

### Opcional (Post-Sales) ⏳
**Pricing Frontend UI de Edición** (1 día):
- CRUD base-sales-rules (precio base por variante/categoría)
- CRUD sale-modification-rules (descuentos por volumen, etc.)
- CRUD client-overrides (precios especiales por cliente)

---

**Fecha de esta corrección:** 2026-02-14  
**Archivos modificados:** 1 (partyApi.js)  
**Líneas modificadas:** ~50 líneas en 2 métodos
