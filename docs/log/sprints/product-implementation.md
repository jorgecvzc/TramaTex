# Product Module Implementation - Sprint Post-11

**Fecha inicio:** 2026-02-15
**Estado:** ðŸ”¨ En Progreso
**Objetivo:** Implementar completamente Product Module antes de retomar validaciÃ³n Sprint 11

---

## Contexto

Durante Sprint 11 (ValidaciÃ³n ERP Core), se descubriÃ³ que el mÃ³dulo Product estÃ¡ ~30-40% implementado:
- âœ… **Attributes API:** Completamente funcional
- âš ï¸ **Products/Variants:** CÃ³digo existe pero tests no compilan (scope refactoring incompleto)
- âŒ **PartyServiceConfiguration:** ImplementaciÃ³n pendiente

---

## Fase 1: Resolver Errores de CompilaciÃ³n â³

### Problemas Identificados

**Tests desactualizados post-refactoring scope system:**

1. `dtos_additional_test.go`:
   - Llama `domain.NewAttribute` con parÃ¡metros obsoletos (brandID, groupID)
   - Referencia campos eliminados: `dto.AttributeName`, `dto.ScopeBrandID`, `dto.ScopeGroupID`

2. `product_service_additional_test.go`:
   - `CreateAttributeCommand` usa campos eliminados: `ScopeBrandID`, `ScopeGroupID`
   - Mocks repository faltan mÃ©todo `Delete()`

3. `helpers_test.go`:
   - Test `TestAttributeMatchesScopeType` usa funciÃ³n eliminada âœ… RESUELTO

### Plan de CorrecciÃ³n

- [ ] Actualizar `dtos_additional_test.go` con API actual (sin scope params)
- [ ] Actualizar `product_service_additional_test.go` (remover scope fields)
- [ ] Agregar mÃ©todo `Delete()` a mocks de repositories
- [ ] Verificar compilaciÃ³n completa: `go build ./internal/product/...`
- [ ] Ejecutar todos los tests: `go test ./internal/product/...`

---

## Fase 2: Validar Coverage Real â³

Una vez que compile:
- [ ] Ejecutar tests con coverage por capa
- [ ] Documentar coverage real de application/interfaces/persistence
- [ ] Identificar gaps de testing reales (no errores de compilaciÃ³n)

---

## Fase 3: ImplementaciÃ³n Faltante â³

SegÃºn anÃ¡lisis de cÃ³digo vs documentaciÃ³n:

### 3.1 Products API (Prioridad: ALTA)
**Estado actual:** Handlers y service existen, necesitan validaciÃ³n funcional

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
- [ ] Sistema Just-in-Time completo (adr-015)

### 3.3 PartyServiceConfiguration (Prioridad: MEDIA)
**Estado actual:** Parcialmente implementado

Endpoints a validar:
- [ ] `POST /api/parties/:id/service-configurations`
- [ ] `GET /api/parties/:id/service-configurations`
- [ ] `PUT /api/parties/:id/service-configurations/:configId`  
- [ ] `DELETE /api/parties/:id/service-configurations/:configId`

---

## Fase 4: Testing Completo â³

- [ ] Tests unitarios para nuevas funcionalidades
- [ ] Tests de integraciÃ³n endpoints
- [ ] Validar coverage â‰¥85% promedio, â‰¥90% domain
- [ ] Tests de sistema Just-in-Time variants

---

## Fase 5: DocumentaciÃ³n â³

- [ ] Actualizar erp-core-completion.md con estado real
- [ ] Marcar endpoints como âœ… en api-contracts.md
- [ ] Documentar decisiones tÃ©cnicas tomadas
- [ ] Actualizar domain-model.md si hay cambios

---

## Criterios de Completitud

âœ… **Must-Have (Retomar Sprint 11):**
- [x] Todo el mÃ³dulo Product compila sin errores
- [ ] Coverage â‰¥85% promedio (application, interfaces, persistence)
- [ ] Products CRUD funcional (endpoints + tests)
- [ ] ProductVariants CRUD funcional
- [ ] Sistema Just-in-Time variants implementado (adr-015)
- [ ] DocumentaciÃ³n alineada con implementaciÃ³n real

âš ï¸ **Nice-to-Have (Opcional):**
- [ ] PartyServiceConfiguration completo
- [ ] Optimizaciones de performance
- [ ] Tests de carga

---

**Ãšltima actualizaciÃ³n:** 2026-02-15 11:00

