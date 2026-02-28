# Tarea 02 - Sprint 11: Critical Remediation Plan - ERP Core

---

## ðŸ“‹ INFORMACIÃ“N DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 02 |
| **ID de Sprint** | sprint-11 |
| **TÃ­tulo** | Critical Remediation Plan - ERP Core |
| **Estado** | ðŸ”„ En Progreso |
| **Facilitador/LLM** | GitHub Copilot (Claude Sonnet 4.5) |
| **Fecha de Inicio** | 2026-02-17 |
| **Fecha de Fin** | â€” |
| **DuraciÃ³n Estimada** | 33-46 horas |
| **DuraciÃ³n Real** | â€” (completar al finalizar) |

---

## ðŸŽ¯ OBJETIVOS PRINCIPALES

**Plan de remediaciÃ³n de bloqueadores crÃ­ticos identificados en Sprint 11-01 antes de proceder con MES Module:**

1. [âœ…] **FASE 1: Cleanup Artifacts & .gitignore (1-2h) - COMPLETADO**
   - Mover/eliminar 30+ archivos coverage dispersos âœ… (21 files removed)
   - Arreglar .gitignore corrupto (espacios, reglas faltantes) âœ…  
   - Eliminar binarios del repositorio (*.exe, party, product) âœ… (6 binaries)
   - Limpiar /tmp/ directory (20 archivos temporales) âœ… (4 MD moved, 16 deleted)

2. [ðŸ”„] **FASE 2: TypeScript Migration - API Services (8-12h)**
   - Migrar partyApi.js â†’ partyApi.ts (579 lÃ­neas)
   - Migrar productApi.js â†’ productApi.ts (794 lÃ­neas)
   - Migrar pricingApi.js â†’ pricingApi.ts (296 lÃ­neas)
   - Migrar salesApi.js â†’ salesApi.ts (523 lÃ­neas)
   - Definir tipos para DTOs request/response
   - Total: 2,192 lÃ­neas JS â†’ TypeScript

3. [ ] **FASE 3: Frontend ERP Core Tests (24-32h)**
   - Tests Party Module (6-8h): 6 componentes + partyApi
   - Tests Product Module (10-12h): 12 componentes + productApi
   - Tests Pricing Module (2-3h): PricingPanel + pricingApi
   - Tests Sales Module (6-9h): 11 pÃ¡ginas + salesApi
   - Target: â‰¥70% coverage en componentes crÃ­ticos

---

## ðŸ“Š CONTEXTO DE ENTRADA

### Hallazgos Sprint 11-01

**Bloqueadores CrÃ­ticos Identificados:**

1. **Artifacts Management CaÃ³tico (1-2h):**
   - 30+ archivos coverage dispersos en `apps/tramatex-api/`
   - .gitignore corrupto con espacios (`c o v e r a g e /`)
   - Binarios versionados (api.exe, application.test.exe, main.exe, tramatex.exe, party, product)
   - /tmp/ con 20 archivos temporales versionados

2. **Services JavaScript Sin Types (8-12h):**
   - partyApi.js: 579 lÃ­neas
   - productApi.js: 794 lÃ­neas
   - pricingApi.js: 296 lÃ­neas
   - salesApi.js: 523 lÃ­neas
   - Total: 2,192 lÃ­neas sin type safety

3. **Frontend ERP Core Sin Tests (24-32h):**
   - Coverage actual: 6.6% (solo Auth testeado)
   - Party/Product/Pricing/Sales: 0% coverage
   - 71 archivos sin tests

### DecisiÃ³n de Sprint 11-01

Sprint 11-01 concluyÃ³ con decisiÃ³n **NO-GO para MES Module** hasta completar esta remediaciÃ³n crÃ­tica.

---

## ðŸ› ï¸ PLAN DE TRABAJO

### FASE 1: CLEANUP ARTIFACTS & .GITIGNORE (1-2h) âœ… COMPLETADO

**Prioridad:** CRÃTICA  
**Esfuerzo:** 1-2 horas  
**DuraciÃ³n Real:** 1.5 horas  
**Impacto:** Alto - Previene contaminaciÃ³n repo, cumple generic-rules.yaml

#### 1.1 AuditorÃ­a de Archivos Coverage âœ…
- [x] Listar todos los archivos coverage en `apps/tramatex-api/` (21 files found)
- [x] Identificar ubicaciÃ³n correcta segÃºn generic-rules.yaml
- [x] Decidir: eliminar (regenerables) â†’ DECISIÃ“N: eliminar todos

#### 1.2 Limpieza de Coverage âœ…
- [x] Eliminar archivos: `cov-*`, `coverage-*`, `*.coverage.out`, `*.out` (17 files)
- [x] Validado 0 files restantes con Get-ChildItem

**Archivos eliminados:**
- cov-domain, cov-product, cov-product-all, cov-product-app, cov-product-domain
- cov-sales-int, cov-sales-interfaces
- coverage, coverage-handlers, coverage-party, coverage-product-module
- coverage-sales, coverage-sales-domain, coverage-sales-infra, coverage-sales-interfaces
- coverage-total, party_coverage

#### 1.3 Arreglo de .gitignore âœ…
- [x] Leer `.gitignore` actual y documentar problemas (spaces in patterns)
- [x] Eliminar caracteres corruptos (`c o v e r a g e /` â†’ `coverage/`)
- [x] Agregar reglas completas:
  - coverage/, *.coverage.out, *.out, cov-*, coverage-*
  - *.exe, *.test, *.test.exe
  - /tmp/, $env, $out, NUL
- [x] Rewrite completo con 80 lÃ­neas de reglas

#### 1.4 EliminaciÃ³n de Binarios âœ…
- [x] Eliminados 6 binarios:
  - `apps/tramatex-api/api.exe` (20MB)
  - `apps/tramatex-api/application.test.exe` (17MB)
  - `apps/tramatex-api/main.exe` (22MB)
  - `apps/tramatex-api/tramatex.exe` (19MB)
  - `apps/tramatex-api/party` (binary)
  - `apps/tramatex-api/product` (binary)
- [x] Agregado a .gitignore para prevenir re-versionado

#### 1.5 Limpieza de /tmp/ âœ…
- [x] Auditar archivos en `/tmp/` (20 files found)
- [x] DecisiÃ³n: 4 MD analysis files moved, 16 deleted
- [x] Archivos preservados en docs/log/analysis/:
  - CORRECION_PARTIES_PRICING.md
  - erp-core-completeness-analysis.md
  - party-frontend-validation.md
  - product-frontend-validation.md
- [x] Ejecutar limpieza (validado 0 files remaining)
- [x] `/tmp/` ya en .gitignore

#### 1.6 Commit y ValidaciÃ³n âœ…
- [x] Commit ejecutado: "chore(sprint-11): cleanup coverage artifacts..."
- [x] Commit hash: **526b2aa**
- [x] Archivos procesados: 20 files changed (6368 insertions, 7082 deletions)
- [x] 6 coverage files deleted
- [x] 11 documentation files added (Sprint 11 + analysis)
- [x] .gitignore fixed and enhanced

**Resultado:** FASE 1 completada exitosamente, 0 artifacts restantes en repo

---

### FASE 2: TYPESCRIPT MIGRATION - API SERVICES (8-12h)

**Prioridad:** CRÃTICA  
**Esfuerzo:** 8-12 horas  
**Impacto:** Alto - Type safety en 2,192 lÃ­neas

#### 2.1 Setup TypeScript (30min)
- [ ] Verificar tsconfig.json configurado en `apps/frontend/`
- [ ] Asegurar tipos para axios, vue-router, pinia instalados
- [ ] Crear directorio `apps/frontend/src/types/` si no existe

#### 2.2 Definir DTOs TypeScript (1-2h)
- [ ] Crear `types/party.ts` con interfaces (Party, PersonProfile, OrganizationProfile, etc.)
- [ ] Crear `types/product.ts` (Product, ProductVariant, Attribute, Brand, etc.)
- [ ] Crear `types/pricing.ts` (PriceCalculation, PricingRule, etc.)
- [ ] Crear `types/sales.ts` (Quote, Order, DeliveryNote, Invoice, etc.)
- [ ] Reutilizar tipos de `auth.ts` existente como referencia

#### 2.3 Migrar partyApi.js â†’ partyApi.ts (2-3h)
- [ ] Renombrar `services/partyApi.js` â†’ `partyApi.ts`
- [ ] Importar tipos de `types/party.ts`
- [ ] Tipar parÃ¡metros de mÃ©todos (CreatePartyDto, UpdatePartyDto, etc.)
- [ ] Tipar retornos de mÃ©todos (Promise<Party>, Promise<Party[]>, etc.)
- [ ] Actualizar imports en componentes que usan partyApi
- [ ] Ejecutar TypeScript compiler y arreglar errores
- [ ] Tests: validar que compilaciÃ³n pasa

#### 2.4 Migrar productApi.js â†’ productApi.ts (3-4h)
- [ ] Renombrar `services/productApi.js` â†’ `productApi.ts`
- [ ] Importar tipos de `types/product.ts`
- [ ] Tipar 20+ mÃ©todos del ProductApiService
- [ ] Actualizar imports en 12 componentes Product
- [ ] Resolver errores TypeScript (esperar ~20-30 errores iniciales)
- [ ] Validar compilaciÃ³n

#### 2.5 Migrar pricingApi.js â†’ pricingApi.ts (1h)
- [ ] Renombrar `services/pricingApi.js` â†’ `pricingApi.ts`
- [ ] Importar tipos de `types/pricing.ts`
- [ ] Tipar mÃ©todos (calculate, getRules, etc.)
- [ ] Actualizar imports en PricingPanel component
- [ ] Validar compilaciÃ³n

#### 2.6 Migrar salesApi.js â†’ salesApi.ts (2-3h)
- [ ] Renombrar `services/salesApi.js` â†’ `salesApi.ts`
- [ ] Importar tipos de `types/sales.ts`
- [ ] Tipar mÃ©todos de SalesApi class (~15 mÃ©todos)
- [ ] Actualizar imports en 11 pÃ¡ginas Sales
- [ ] Resolver errores TypeScript
- [ ] Validar compilaciÃ³n

#### 2.7 ValidaciÃ³n Final TypeScript (30min)
- [ ] Ejecutar `npm run build` - debe pasar sin errores
- [ ] Ejecutar `npm run type-check` si existe script
- [ ] Validar que autocomplete funciona en VSCode
- [ ] Commit: "refactor: migrate API services to TypeScript (2,192 lines)"

---

### FASE 3: FRONTEND ERP CORE TESTS (24-32h)

**Prioridad:** CRÃTICA  
**Esfuerzo:** 24-32 horas  
**Impacto:** Alto - Coverage 6.6% â†’ â‰¥70%

#### 3.1 Setup Testing Environment (1h)
- [ ] Verificar Vitest configurado (ya existe)
- [ ] Verificar @testing-library/vue instalado
- [ ] Estudiar tests Auth existentes como referencia (33 tests, 100% coverage)
- [ ] Crear estructura `__tests__/` en cada mÃ³dulo si falta

#### 3.2 Tests Party Module (6-8h)
**Componentes a testear:**
- [ ] PartySelector.vue (395 lÃ­neas, autocomplete) - 2h
- [ ] PartyForm.vue (creaciÃ³n/ediciÃ³n) - 1.5h
- [ ] PartyDetail.vue (tabs, visualizaciÃ³n) - 1.5h
- [ ] PartyList.vue (filtros, bÃºsqueda) - 1h
- [ ] AddressManager.vue (CRUD addresses) - 1h
- [ ] PersonManager.vue (contacts) - 1h

**API Service:**
- [ ] partyApi.ts (579 lÃ­neas) - 2h
  - Tests de mÃ©todos principales (create, update, get, list, batch)
  - Mock axios responses
  - Validar error handling

**Target:** â‰¥70% coverage Party Module

#### 3.3 Tests Product Module (10-12h)
**Componentes a testear (12 total):**
- [ ] ProductFormBasic.vue - 1h
- [ ] ProductFormAttributes.vue - 1h
- [ ] ProductFormClassification.vue - 1h
- [ ] VariantTable.vue - 1.5h
- [ ] VariantGenerator.vue (lÃ³gica compleja) - 2h
- [ ] VariantSelector.vue (modal) - 1.5h
- [ ] PricingPanel.vue (calculadora) - 2h
- [ ] AttributeCard.vue - 30min
- [ ] Otros componentes menores (4 componentes) - 2h

**API Service:**
- [ ] productApi.ts (794 lÃ­neas) - 3h
  - Tests CRUD Products
  - Tests Variants generation
  - Tests Attributes configurables
  - Tests Brands/Groups

**Target:** â‰¥70% coverage Product Module

#### 3.4 Tests Pricing Module (2-3h)
**Componente:**
- [ ] PricingPanel.vue (1,030 lÃ­neas aprox segÃºn Sprint 09) - 2h
  - Tests calculadora interactiva
  - Tests tabla precios base
  - Tests modal historial

**API Service:**
- [ ] pricingApi.ts (296 lÃ­neas) - 1h
  - Tests calculate price
  - Tests get/create rules
  - Tests client overrides

**Target:** â‰¥70% coverage Pricing Module

#### 3.5 Tests Sales Module (6-9h)
**PÃ¡ginas a testear (11 total):**
- [ ] QuoteList.vue - 1h
- [ ] QuoteCreate.vue (548 lÃ­neas) - 2h
- [ ] QuoteDetail.vue (490 lÃ­neas) - 1.5h
- [ ] OrderList.vue - 1h
- [ ] OrderCreate.vue - 1.5h
- [ ] OrderDetail.vue (1,286 lÃ­neas) - 2h
- [ ] DeliveryNoteList.vue - 45min
- [ ] DeliveryNoteDetail.vue (430 lÃ­neas) - 1h
- [ ] InvoiceList.vue - 1h
- [ ] InvoiceDetail.vue - 1h
- [ ] TicketCreate.vue - 45min

**API Service:**
- [ ] salesApi.ts (523 lÃ­neas) - 3h
  - Tests CRUD Quotes
  - Tests CRUD Orders
  - Tests CRUD DeliveryNotes
  - Tests CRUD Invoices
  - Tests workflow transitions

**Target:** â‰¥70% coverage Sales Module

#### 3.6 ValidaciÃ³n Final Coverage (1h)
- [ ] Ejecutar `npm run test:coverage`
- [ ] Generar reporte HTML
- [ ] Verificar que cada mÃ³dulo cumple â‰¥70%
- [ ] Identificar archivos crÃ­ticos <70% y agregar tests adicionales
- [ ] Commit: "test: add comprehensive tests for ERP Core (Party, Product, Pricing, Sales)"
- [ ] Actualizar documentaciÃ³n con nuevos nÃºmeros de coverage

---

## ðŸ“ˆ CRITERIOS DE Ã‰XITO

### Criterios Obligatorios (Must Have)
- [ ] .gitignore corregido sin caracteres corruptos
- [ ] 0 binarios versionados en repo
- [ ] 0 archivos coverage dispersos fuera de coverage-reports/
- [ ] /tmp/ limpio o ignorado
- [ ] 4 API services migrados a TypeScript (party, product, pricing, sales)
- [ ] Frontend ERP Core â‰¥70% coverage
- [ ] CompilaciÃ³n TypeScript sin errores
- [ ] Todos los tests pasando (100% success rate)

### Criterios Deseables (Should Have)
- [ ] Frontend ERP Core â‰¥80% coverage (supera target)
- [ ] DTOs TypeScript reutilizables y bien documentados
- [ ] Tests con describe/it bien estructurados
- [ ] DocumentaciÃ³n actualizada en README.md frontend

### Criterios Opcionales (Nice to Have)
- [ ] Scripts npm para coverage por mÃ³dulo
- [ ] CI/CD pipeline valida coverage mÃ­nimos
- [ ] Badges de coverage en README.md

---

## ðŸ“Š MÃ‰TRICAS DE PROGRESO

### FASE 1: Cleanup (1-2h)

| Tarea | Estado | Tiempo | Impacto |
|-------|--------|--------|---------|
| AuditorÃ­a coverage | â³ | â€” | â€” |
| Limpieza coverage | â€” | â€” | â€” |
| Arreglo .gitignore | â€” | â€” | â€” |
| EliminaciÃ³n binarios | â€” | â€” | â€” |
| Limpieza /tmp/ | â€” | â€” | â€” |
| Commit & validaciÃ³n | â€” | â€” | â€” |

### FASE 2: TypeScript (8-12h)

| Tarea | LÃ­neas | Estado | Tiempo | Errores |
|-------|--------|--------|--------|---------|
| DTOs TypeScript | â€” | â€” | â€” | â€” |
| partyApi.ts | 579 | â€” | â€” | â€” |
| productApi.ts | 794 | â€” | â€” | â€” |
| pricingApi.ts | 296 | â€” | â€” | â€” |
| salesApi.ts | 523 | â€” | â€” | â€” |
| **TOTAL** | **2,192** | â€” | â€” | â€” |

### FASE 3: Tests Frontend (24-32h)

| MÃ³dulo | Tests | Coverage | Estado | Tiempo |
|--------|-------|----------|--------|--------|
| Party | 0 | 0% â†’ â€” | â€” | â€” |
| Product | 0 | 0% â†’ â€” | â€” | â€” |
| Pricing | 0 | 0% â†’ â€” | â€” | â€” |
| Sales | 0 | 0% â†’ â€” | â€” | â€” |
| **TOTAL** | **0 â†’ â€”** | **6.6% â†’ â€”** | â€” | â€” |

---

## ðŸ“ DECISIONES Y CAMBIOS

### Decisiones Tomadas

**[2026-02-17] DecisiÃ³n #1: Priorizar FASE 1 (Cleanup) como fundaciÃ³n**

**Contexto:** Sprint 11-01 identificÃ³ artifacts management caÃ³tico como blocker crÃ­tico.

**DecisiÃ³n:** Comenzar con cleanup (1-2h) antes de TypeScript/Tests.

**RazÃ³n:**
1. RÃ¡pido (1-2h) y alto impacto
2. Previene contaminar repo con nuevos artifacts durante desarrollo
3. .gitignore correcto es prerequisito para fases 2 y 3
4. Cumple generic-rules.yaml inmediatamente

**AcciÃ³n:** Ejecutar FASE 1 completa antes de proceder a FASE 2.

---

## ðŸš€ PRÃ“XIMOS PASOS

### DespuÃ©s de Completar Esta Tarea

- [ ] Actualizar sprint-registry.yaml con task_11_02 como completed
- [ ] Actualizar session-log.md
- [ ] Ejecutar validaciÃ³n final:
  - [ ] `npm run build` pasa
  - [ ] `npm run test` pasa
  - [ ] `npm run test:coverage` muestra â‰¥70%
  - [ ] `git status` muestra repo limpio
- [ ] DecisiÃ³n GO para MES Module
- [ ] Comunicar resultados al equipo

---

## ðŸ“Œ NOTAS Y OBSERVACIONES

### Notas Importantes

- Esta tarea es **BLOCKER** para MES Module segÃºn Sprint 11-01
- Esfuerzo estimado: 33-46h (ajustar si es necesario)
- Prioridad FASE 1 (cleanup) por ser rÃ¡pida y fundacional
- Usar tests Auth existentes como patrÃ³n de referencia (33 tests, 100% coverage)

### Observaciones Durante la Tarea

**[2026-02-17 - Inicio de Tarea]**

Tarea iniciada basada en hallazgos crÃ­ticos de Sprint 11-01. Plan estructurado en 3 fases con estimaciones detalladas. Comenzando con FASE 1 (Cleanup Artifacts) que es la mÃ¡s rÃ¡pida y tiene mayor impacto inmediato.

**[2026-02-17 - FASE 1 Completada]**

FASE 1 completada exitosamente en ~1.5h. Resultados:
- âœ… 21 coverage artifacts eliminados (17 coverage + 4 binaries from apps/tramatex-api/)
- âœ… 6 binarios eliminados (api.exe 20MB, main.exe 22MB, tramatex.exe 19MB, application.test.exe 17MB, party, product)
- âœ… .gitignore reescrito con 80 lÃ­neas (fix corruption: "c o v e r a g e /" â†’ "coverage/")
- âœ… /tmp/ limpiado: 4 MD analysis moved to docs/log/analysis/, 16 files deleted
- âœ… Commit 526b2aa: "chore(sprint-11): cleanup coverage artifacts and complete FASE 7"
- âœ… ValidaciÃ³n: 0 artifacts restantes en repo

Procediendo con FASE 2 (TypeScript Migration).

---

## âœ… CHECKLIST FINAL

Antes de marcar esta tarea como completada:

- [x] **FASE 1: Cleanup completada (1-2h)** âœ… 2026-02-17
- [ ] FASE 2: TypeScript migration completada (8-12h)
- [ ] FASE 3: Frontend tests completada (24-32h)
- [x] .gitignore corregido y validado âœ…
- [x] 0 binarios en repo âœ…
- [x] 0 coverage artifacts dispersos âœ…
- [ ] 2,192 lÃ­neas migradas a TypeScript
- [ ] Frontend ERP Core â‰¥70% coverage
- [ ] Todos los tests pasando
- [ ] CompilaciÃ³n TypeScript sin errores
- [ ] DocumentaciÃ³n actualizada
- [ ] Commit final realizado

---

**Ãšltima ActualizaciÃ³n:** 2026-02-17 (FASE 1 Completada)  
**Estado:** ðŸ”„ EN PROGRESO - FASE 2 (TypeScript Migration)  
**PrÃ³xima AcciÃ³n:** AuditorÃ­a de archivos coverage dispersos


