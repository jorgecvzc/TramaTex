# Product Module Implementation - Sprint Post-11

**Fecha inicio:** 2026-02-15
**Estado:** 🔨 En Progreso
**Objetivo:** Implementar completamente Product Module antes de retomar validación Sprint 11

---

## Contexto

Durante Sprint 11 (Validación ERP Core), se descubrió que el módulo Product está ~30-40% implementado:
- ✅ **Attributes API:** Completamente funcional
- ⚠️ **Products/Variants:** Código existe pero tests no compilan (scope refactoring incompleto)
- ❌ **PartyServiceConfiguration:** Implementación pendiente

---

## Fase 1: Resolver Errores de Compilación ⏳

### Problemas Identificados

**Tests desactualizados post-refactoring scope system:**

1. `dtos_additional_test.go`:
   - Llama `domain.NewAttribute` con parámetros obsoletos (brandID, groupID)
   - Referencia campos eliminados: `dto.AttributeName`, `dto.ScopeBrandID`, `dto.ScopeGroupID`

2. `product_service_additional_test.go`:
   - `CreateAttributeCommand` usa campos eliminados: `ScopeBrandID`, `ScopeGroupID`
   - Mocks repository faltan método `Delete()`

3. `helpers_test.go`:
   - Test `TestAttributeMatchesScopeType` usa función eliminada ✅ RESUELTO

### Plan de Corrección

- [ ] Actualizar `dtos_additional_test.go` con API actual (sin scope params)
- [ ] Actualizar `product_service_additional_test.go` (remover scope fields)
- [ ] Agregar método `Delete()` a mocks de repositories
- [ ] Verificar compilación completa: `go build ./internal/product/...`
- [ ] Ejecutar todos los tests: `go test ./internal/product/...`

---

## Fase 2: Validar Coverage Real ⏳

Una vez que compile:
- [ ] Ejecutar tests con coverage por capa
- [ ] Documentar coverage real de application/interfaces/persistence
- [ ] Identificar gaps de testing reales (no errores de compilación)

---

## Fase 3: Implementación Faltante ⏳

Según análisis de código vs documentación:

### 3.1 Products API (Prioridad: ALTA)
**Estado actual:** Handlers y service existen, necesitan validación funcional

Endpoints a verificar:
- [ ] `POST /api/products` - CreateProduct
- [ ] `GET /api/products` - ListProducts
- [ ] `GET /api/products/:id` - GetProductByID
- [ ] `PUT /api/products/:id` - UpdateProduct
- [ ] `POST /api/products/:id/groups` - AddGroupToProduct
- [ ] `POST /api/products/:id/attributes` - AddDirectAttributeToProduct

### 3.2 ProductVariants API (Prioridad: ALTA)
**Estado actual:** Domain y service existen, handlers parciales

Endpoints a implementar/validar:
- [ ] `POST /api/products/:id/variants/generate` - GenerateProductVariants (Just-in-Time)
- [ ] `POST /api/products/:id/variants/find-or-create` - FindOrCreateProductVariant
- [ ] `GET /api/products/:id/variants` - ListProductVariantsByProductID
- [ ] Sistema Just-in-Time completo (ADR-015)

### 3.3 PartyServiceConfiguration (Prioridad: MEDIA)
**Estado actual:** Parcialmente implementado

Endpoints a validar:
- [ ] `POST /api/parties/:id/service-configurations`
- [ ] `GET /api/parties/:id/service-configurations`
- [ ] `PUT /api/parties/:id/service-configurations/:configId`  
- [ ] `DELETE /api/parties/:id/service-configurations/:configId`

---

## Fase 4: Testing Completo ⏳

- [ ] Tests unitarios para nuevas funcionalidades
- [ ] Tests de integración endpoints
- [ ] Validar coverage ≥85% promedio, ≥90% domain
- [ ] Tests de sistema Just-in-Time variants

---

## Fase 5: Documentación ⏳

- [ ] Actualizar ERP_CORE_COMPLETION.md con estado real
- [ ] Marcar endpoints como ✅ en api-contracts.md
- [ ] Documentar decisiones técnicas tomadas
- [ ] Actualizar domain-model.md si hay cambios

---

## Criterios de Completitud

✅ **Must-Have (Retomar Sprint 11):**
- [x] Todo el módulo Product compila sin errores
- [ ] Coverage ≥85% promedio (application, interfaces, persistence)
- [ ] Products CRUD funcional (endpoints + tests)
- [ ] ProductVariants CRUD funcional
- [ ] Sistema Just-in-Time variants implementado (ADR-015)
- [ ] Documentación alineada con implementación real

⚠️ **Nice-to-Have (Opcional):**
- [ ] PartyServiceConfiguration completo
- [ ] Optimizaciones de performance
- [ ] Tests de carga

---

**Última actualización:** 2026-02-15 11:00
