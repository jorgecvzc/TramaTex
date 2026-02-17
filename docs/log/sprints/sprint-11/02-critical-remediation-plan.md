# Tarea 02 - Sprint 11: Critical Remediation Plan - ERP Core

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 02 |
| **ID de Sprint** | sprint-11 |
| **Título** | Critical Remediation Plan - ERP Core |
| **Estado** | 🔄 En Progreso |
| **Facilitador/LLM** | GitHub Copilot (Claude Sonnet 4.5) |
| **Fecha de Inicio** | 2026-02-17 |
| **Fecha de Fin** | — |
| **Duración Estimada** | 33-46 horas |
| **Duración Real** | — (completar al finalizar) |

---

## 🎯 OBJETIVOS PRINCIPALES

**Plan de remediación de bloqueadores críticos identificados en Sprint 11-01 antes de proceder con MES Module:**

1. [✅] **FASE 1: Cleanup Artifacts & .gitignore (1-2h) - COMPLETADO**
   - Mover/eliminar 30+ archivos coverage dispersos ✅ (21 files removed)
   - Arreglar .gitignore corrupto (espacios, reglas faltantes) ✅  
   - Eliminar binarios del repositorio (*.exe, party, product) ✅ (6 binaries)
   - Limpiar /tmp/ directory (20 archivos temporales) ✅ (4 MD moved, 16 deleted)

2. [🔄] **FASE 2: TypeScript Migration - API Services (8-12h)**
   - Migrar partyApi.js → partyApi.ts (579 líneas)
   - Migrar productApi.js → productApi.ts (794 líneas)
   - Migrar pricingApi.js → pricingApi.ts (296 líneas)
   - Migrar salesApi.js → salesApi.ts (523 líneas)
   - Definir tipos para DTOs request/response
   - Total: 2,192 líneas JS → TypeScript

3. [ ] **FASE 3: Frontend ERP Core Tests (24-32h)**
   - Tests Party Module (6-8h): 6 componentes + partyApi
   - Tests Product Module (10-12h): 12 componentes + productApi
   - Tests Pricing Module (2-3h): PricingPanel + pricingApi
   - Tests Sales Module (6-9h): 11 páginas + salesApi
   - Target: ≥70% coverage en componentes críticos

---

## 📊 CONTEXTO DE ENTRADA

### Hallazgos Sprint 11-01

**Bloqueadores Críticos Identificados:**

1. **Artifacts Management Caótico (1-2h):**
   - 30+ archivos coverage dispersos en `apps/tramatex-api/`
   - .gitignore corrupto con espacios (`c o v e r a g e /`)
   - Binarios versionados (api.exe, application.test.exe, main.exe, tramatex.exe, party, product)
   - /tmp/ con 20 archivos temporales versionados

2. **Services JavaScript Sin Types (8-12h):**
   - partyApi.js: 579 líneas
   - productApi.js: 794 líneas
   - pricingApi.js: 296 líneas
   - salesApi.js: 523 líneas
   - Total: 2,192 líneas sin type safety

3. **Frontend ERP Core Sin Tests (24-32h):**
   - Coverage actual: 6.6% (solo Auth testeado)
   - Party/Product/Pricing/Sales: 0% coverage
   - 71 archivos sin tests

### Decisión de Sprint 11-01

Sprint 11-01 concluyó con decisión **NO-GO para MES Module** hasta completar esta remediación crítica.

---

## 🛠️ PLAN DE TRABAJO

### FASE 1: CLEANUP ARTIFACTS & .GITIGNORE (1-2h) ✅ COMPLETADO

**Prioridad:** CRÍTICA  
**Esfuerzo:** 1-2 horas  
**Duración Real:** 1.5 horas  
**Impacto:** Alto - Previene contaminación repo, cumple generic-rules.yaml

#### 1.1 Auditoría de Archivos Coverage ✅
- [x] Listar todos los archivos coverage en `apps/tramatex-api/` (21 files found)
- [x] Identificar ubicación correcta según generic-rules.yaml
- [x] Decidir: eliminar (regenerables) → DECISIÓN: eliminar todos

#### 1.2 Limpieza de Coverage ✅
- [x] Eliminar archivos: `cov-*`, `coverage-*`, `*.coverage.out`, `*.out` (17 files)
- [x] Validado 0 files restantes con Get-ChildItem

**Archivos eliminados:**
- cov-domain, cov-product, cov-product-all, cov-product-app, cov-product-domain
- cov-sales-int, cov-sales-interfaces
- coverage, coverage-handlers, coverage-party, coverage-product-module
- coverage-sales, coverage-sales-domain, coverage-sales-infra, coverage-sales-interfaces
- coverage-total, party_coverage

#### 1.3 Arreglo de .gitignore ✅
- [x] Leer `.gitignore` actual y documentar problemas (spaces in patterns)
- [x] Eliminar caracteres corruptos (`c o v e r a g e /` → `coverage/`)
- [x] Agregar reglas completas:
  - coverage/, *.coverage.out, *.out, cov-*, coverage-*
  - *.exe, *.test, *.test.exe
  - /tmp/, $env, $out, NUL
- [x] Rewrite completo con 80 líneas de reglas

#### 1.4 Eliminación de Binarios ✅
- [x] Eliminados 6 binarios:
  - `apps/tramatex-api/api.exe` (20MB)
  - `apps/tramatex-api/application.test.exe` (17MB)
  - `apps/tramatex-api/main.exe` (22MB)
  - `apps/tramatex-api/tramatex.exe` (19MB)
  - `apps/tramatex-api/party` (binary)
  - `apps/tramatex-api/product` (binary)
- [x] Agregado a .gitignore para prevenir re-versionado

#### 1.5 Limpieza de /tmp/ ✅
- [x] Auditar archivos en `/tmp/` (20 files found)
- [x] Decisión: 4 MD analysis files moved, 16 deleted
- [x] Archivos preservados en docs/log/analysis/:
  - CORRECION_PARTIES_PRICING.md
  - ERP_CORE_COMPLETENESS_ANALYSIS.md
  - PARTY_FRONTEND_VALIDATION.md
  - PRODUCT_FRONTEND_VALIDATION.md
- [x] Ejecutar limpieza (validado 0 files remaining)
- [x] `/tmp/` ya en .gitignore

#### 1.6 Commit y Validación ✅
- [x] Commit ejecutado: "chore(sprint-11): cleanup coverage artifacts..."
- [x] Commit hash: **526b2aa**
- [x] Archivos procesados: 20 files changed (6368 insertions, 7082 deletions)
- [x] 6 coverage files deleted
- [x] 11 documentation files added (Sprint 11 + analysis)
- [x] .gitignore fixed and enhanced

**Resultado:** FASE 1 completada exitosamente, 0 artifacts restantes en repo

---

### FASE 2: TYPESCRIPT MIGRATION - API SERVICES (8-12h)

**Prioridad:** CRÍTICA  
**Esfuerzo:** 8-12 horas  
**Impacto:** Alto - Type safety en 2,192 líneas

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

#### 2.3 Migrar partyApi.js → partyApi.ts (2-3h)
- [ ] Renombrar `services/partyApi.js` → `partyApi.ts`
- [ ] Importar tipos de `types/party.ts`
- [ ] Tipar parámetros de métodos (CreatePartyDto, UpdatePartyDto, etc.)
- [ ] Tipar retornos de métodos (Promise<Party>, Promise<Party[]>, etc.)
- [ ] Actualizar imports en componentes que usan partyApi
- [ ] Ejecutar TypeScript compiler y arreglar errores
- [ ] Tests: validar que compilación pasa

#### 2.4 Migrar productApi.js → productApi.ts (3-4h)
- [ ] Renombrar `services/productApi.js` → `productApi.ts`
- [ ] Importar tipos de `types/product.ts`
- [ ] Tipar 20+ métodos del ProductApiService
- [ ] Actualizar imports en 12 componentes Product
- [ ] Resolver errores TypeScript (esperar ~20-30 errores iniciales)
- [ ] Validar compilación

#### 2.5 Migrar pricingApi.js → pricingApi.ts (1h)
- [ ] Renombrar `services/pricingApi.js` → `pricingApi.ts`
- [ ] Importar tipos de `types/pricing.ts`
- [ ] Tipar métodos (calculate, getRules, etc.)
- [ ] Actualizar imports en PricingPanel component
- [ ] Validar compilación

#### 2.6 Migrar salesApi.js → salesApi.ts (2-3h)
- [ ] Renombrar `services/salesApi.js` → `salesApi.ts`
- [ ] Importar tipos de `types/sales.ts`
- [ ] Tipar métodos de SalesApi class (~15 métodos)
- [ ] Actualizar imports en 11 páginas Sales
- [ ] Resolver errores TypeScript
- [ ] Validar compilación

#### 2.7 Validación Final TypeScript (30min)
- [ ] Ejecutar `npm run build` - debe pasar sin errores
- [ ] Ejecutar `npm run type-check` si existe script
- [ ] Validar que autocomplete funciona en VSCode
- [ ] Commit: "refactor: migrate API services to TypeScript (2,192 lines)"

---

### FASE 3: FRONTEND ERP CORE TESTS (24-32h)

**Prioridad:** CRÍTICA  
**Esfuerzo:** 24-32 horas  
**Impacto:** Alto - Coverage 6.6% → ≥70%

#### 3.1 Setup Testing Environment (1h)
- [ ] Verificar Vitest configurado (ya existe)
- [ ] Verificar @testing-library/vue instalado
- [ ] Estudiar tests Auth existentes como referencia (33 tests, 100% coverage)
- [ ] Crear estructura `__tests__/` en cada módulo si falta

#### 3.2 Tests Party Module (6-8h)
**Componentes a testear:**
- [ ] PartySelector.vue (395 líneas, autocomplete) - 2h
- [ ] PartyForm.vue (creación/edición) - 1.5h
- [ ] PartyDetail.vue (tabs, visualización) - 1.5h
- [ ] PartyList.vue (filtros, búsqueda) - 1h
- [ ] AddressManager.vue (CRUD addresses) - 1h
- [ ] PersonManager.vue (contacts) - 1h

**API Service:**
- [ ] partyApi.ts (579 líneas) - 2h
  - Tests de métodos principales (create, update, get, list, batch)
  - Mock axios responses
  - Validar error handling

**Target:** ≥70% coverage Party Module

#### 3.3 Tests Product Module (10-12h)
**Componentes a testear (12 total):**
- [ ] ProductFormBasic.vue - 1h
- [ ] ProductFormAttributes.vue - 1h
- [ ] ProductFormClassification.vue - 1h
- [ ] VariantTable.vue - 1.5h
- [ ] VariantGenerator.vue (lógica compleja) - 2h
- [ ] VariantSelector.vue (modal) - 1.5h
- [ ] PricingPanel.vue (calculadora) - 2h
- [ ] AttributeCard.vue - 30min
- [ ] Otros componentes menores (4 componentes) - 2h

**API Service:**
- [ ] productApi.ts (794 líneas) - 3h
  - Tests CRUD Products
  - Tests Variants generation
  - Tests Attributes configurables
  - Tests Brands/Groups

**Target:** ≥70% coverage Product Module

#### 3.4 Tests Pricing Module (2-3h)
**Componente:**
- [ ] PricingPanel.vue (1,030 líneas aprox según Sprint 09) - 2h
  - Tests calculadora interactiva
  - Tests tabla precios base
  - Tests modal historial

**API Service:**
- [ ] pricingApi.ts (296 líneas) - 1h
  - Tests calculate price
  - Tests get/create rules
  - Tests client overrides

**Target:** ≥70% coverage Pricing Module

#### 3.5 Tests Sales Module (6-9h)
**Páginas a testear (11 total):**
- [ ] QuoteList.vue - 1h
- [ ] QuoteCreate.vue (548 líneas) - 2h
- [ ] QuoteDetail.vue (490 líneas) - 1.5h
- [ ] OrderList.vue - 1h
- [ ] OrderCreate.vue - 1.5h
- [ ] OrderDetail.vue (1,286 líneas) - 2h
- [ ] DeliveryNoteList.vue - 45min
- [ ] DeliveryNoteDetail.vue (430 líneas) - 1h
- [ ] InvoiceList.vue - 1h
- [ ] InvoiceDetail.vue - 1h
- [ ] TicketCreate.vue - 45min

**API Service:**
- [ ] salesApi.ts (523 líneas) - 3h
  - Tests CRUD Quotes
  - Tests CRUD Orders
  - Tests CRUD DeliveryNotes
  - Tests CRUD Invoices
  - Tests workflow transitions

**Target:** ≥70% coverage Sales Module

#### 3.6 Validación Final Coverage (1h)
- [ ] Ejecutar `npm run test:coverage`
- [ ] Generar reporte HTML
- [ ] Verificar que cada módulo cumple ≥70%
- [ ] Identificar archivos críticos <70% y agregar tests adicionales
- [ ] Commit: "test: add comprehensive tests for ERP Core (Party, Product, Pricing, Sales)"
- [ ] Actualizar documentación con nuevos números de coverage

---

## 📈 CRITERIOS DE ÉXITO

### Criterios Obligatorios (Must Have)
- [ ] .gitignore corregido sin caracteres corruptos
- [ ] 0 binarios versionados en repo
- [ ] 0 archivos coverage dispersos fuera de coverage-reports/
- [ ] /tmp/ limpio o ignorado
- [ ] 4 API services migrados a TypeScript (party, product, pricing, sales)
- [ ] Frontend ERP Core ≥70% coverage
- [ ] Compilación TypeScript sin errores
- [ ] Todos los tests pasando (100% success rate)

### Criterios Deseables (Should Have)
- [ ] Frontend ERP Core ≥80% coverage (supera target)
- [ ] DTOs TypeScript reutilizables y bien documentados
- [ ] Tests con describe/it bien estructurados
- [ ] Documentación actualizada en README.md frontend

### Criterios Opcionales (Nice to Have)
- [ ] Scripts npm para coverage por módulo
- [ ] CI/CD pipeline valida coverage mínimos
- [ ] Badges de coverage en README.md

---

## 📊 MÉTRICAS DE PROGRESO

### FASE 1: Cleanup (1-2h)

| Tarea | Estado | Tiempo | Impacto |
|-------|--------|--------|---------|
| Auditoría coverage | ⏳ | — | — |
| Limpieza coverage | — | — | — |
| Arreglo .gitignore | — | — | — |
| Eliminación binarios | — | — | — |
| Limpieza /tmp/ | — | — | — |
| Commit & validación | — | — | — |

### FASE 2: TypeScript (8-12h)

| Tarea | Líneas | Estado | Tiempo | Errores |
|-------|--------|--------|--------|---------|
| DTOs TypeScript | — | — | — | — |
| partyApi.ts | 579 | — | — | — |
| productApi.ts | 794 | — | — | — |
| pricingApi.ts | 296 | — | — | — |
| salesApi.ts | 523 | — | — | — |
| **TOTAL** | **2,192** | — | — | — |

### FASE 3: Tests Frontend (24-32h)

| Módulo | Tests | Coverage | Estado | Tiempo |
|--------|-------|----------|--------|--------|
| Party | 0 | 0% → — | — | — |
| Product | 0 | 0% → — | — | — |
| Pricing | 0 | 0% → — | — | — |
| Sales | 0 | 0% → — | — | — |
| **TOTAL** | **0 → —** | **6.6% → —** | — | — |

---

## 📝 DECISIONES Y CAMBIOS

### Decisiones Tomadas

**[2026-02-17] Decisión #1: Priorizar FASE 1 (Cleanup) como fundación**

**Contexto:** Sprint 11-01 identificó artifacts management caótico como blocker crítico.

**Decisión:** Comenzar con cleanup (1-2h) antes de TypeScript/Tests.

**Razón:**
1. Rápido (1-2h) y alto impacto
2. Previene contaminar repo con nuevos artifacts durante desarrollo
3. .gitignore correcto es prerequisito para fases 2 y 3
4. Cumple generic-rules.yaml inmediatamente

**Acción:** Ejecutar FASE 1 completa antes de proceder a FASE 2.

---

## 🚀 PRÓXIMOS PASOS

### Después de Completar Esta Tarea

- [ ] Actualizar sprint-registry.yaml con task_11_02 como completed
- [ ] Actualizar session-log.md
- [ ] Ejecutar validación final:
  - [ ] `npm run build` pasa
  - [ ] `npm run test` pasa
  - [ ] `npm run test:coverage` muestra ≥70%
  - [ ] `git status` muestra repo limpio
- [ ] Decisión GO para MES Module
- [ ] Comunicar resultados al equipo

---

## 📌 NOTAS Y OBSERVACIONES

### Notas Importantes

- Esta tarea es **BLOCKER** para MES Module según Sprint 11-01
- Esfuerzo estimado: 33-46h (ajustar si es necesario)
- Prioridad FASE 1 (cleanup) por ser rápida y fundacional
- Usar tests Auth existentes como patrón de referencia (33 tests, 100% coverage)

### Observaciones Durante la Tarea

**[2026-02-17 - Inicio de Tarea]**

Tarea iniciada basada en hallazgos críticos de Sprint 11-01. Plan estructurado en 3 fases con estimaciones detalladas. Comenzando con FASE 1 (Cleanup Artifacts) que es la más rápida y tiene mayor impacto inmediato.

**[2026-02-17 - FASE 1 Completada]**

FASE 1 completada exitosamente en ~1.5h. Resultados:
- ✅ 21 coverage artifacts eliminados (17 coverage + 4 binaries from apps/tramatex-api/)
- ✅ 6 binarios eliminados (api.exe 20MB, main.exe 22MB, tramatex.exe 19MB, application.test.exe 17MB, party, product)
- ✅ .gitignore reescrito con 80 líneas (fix corruption: "c o v e r a g e /" → "coverage/")
- ✅ /tmp/ limpiado: 4 MD analysis moved to docs/log/analysis/, 16 files deleted
- ✅ Commit 526b2aa: "chore(sprint-11): cleanup coverage artifacts and complete FASE 7"
- ✅ Validación: 0 artifacts restantes en repo

Procediendo con FASE 2 (TypeScript Migration).

---

## ✅ CHECKLIST FINAL

Antes de marcar esta tarea como completada:

- [x] **FASE 1: Cleanup completada (1-2h)** ✅ 2026-02-17
- [ ] FASE 2: TypeScript migration completada (8-12h)
- [ ] FASE 3: Frontend tests completada (24-32h)
- [x] .gitignore corregido y validado ✅
- [x] 0 binarios en repo ✅
- [x] 0 coverage artifacts dispersos ✅
- [ ] 2,192 líneas migradas a TypeScript
- [ ] Frontend ERP Core ≥70% coverage
- [ ] Todos los tests pasando
- [ ] Compilación TypeScript sin errores
- [ ] Documentación actualizada
- [ ] Commit final realizado

---

**Última Actualización:** 2026-02-17 (FASE 1 Completada)  
**Estado:** 🔄 EN PROGRESO - FASE 2 (TypeScript Migration)  
**Próxima Acción:** Auditoría de archivos coverage dispersos

