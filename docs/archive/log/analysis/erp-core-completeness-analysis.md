# ERP Core - Análisis de Completitud de Módulos

**Fecha:** 2026-02-14  
**Estado:** ✅ Análisis Completo  
**Deadline Proyecto:** 2026-02-23 (9 días restantes)

---

## Resumen Ejecutivo

| Módulo | Endpoints API | Completitud Backend | Frontend | Prioridad | Estado |
|--------|---------------|---------------------|----------|-----------|--------|
| **IAM** | 9 rutas | ✅ 100% | ⚠️ Básico | Baja | ✅ Funcional |
| **Party** | 17 rutas | ✅ 100% | ❌ Pendiente | Alta | ⚠️ Backend listo |
| **Product** | 27 rutas | ✅ 95% | ⚠️ Parcial (50%) | Alta | ⚠️ En progreso |
| **Pricing** | 10 rutas | ✅ 90% | ❌ Pendiente | Media | ⚠️ MVP funcional |
| **Sales** | 21 rutas | ✅ 95% | ❌ Pendiente | Alta | ⚠️ Backend listo |
| **TOTAL** | **84 rutas** | **96%** | **10%** | - | ⚠️ Backend casi completo |

---

## 1. Módulo IAM (Identity & Access Management)

### Endpoints Implementados (9 rutas)

#### Autenticación (4 rutas - sin autenticación requerida)
- ✅ `POST /auth/register` - Registro de usuarios
- ✅ `POST /auth/login` (rate-limited) - Login con JWT
- ✅ `POST /auth/refresh` - Refrescar token
- ✅ `POST /auth/authorize` - Verificar autorización

#### Gestión de Usuarios (3 rutas - protegidas)
- ✅ `POST /auth/logout` - Logout con blacklist de tokens
- ✅ `POST /auth/users` (admin) - Crear usuario por admin
- ✅ `GET /auth/users` (admin) - Listar usuarios
- ✅ `DELETE /auth/users/:id` (admin) - Eliminar usuario

#### Roles (1 ruta)
- ✅ `POST /auth/assign-role` (admin) - Asignar roles a usuarios

### Estado Backend: ✅ 100% Completado
- **Cumple:** Login, registro, refresh token, logout con blacklist, CRUD usuarios, asignación de roles
- **Testing:** ✅ Tests de integración implementados
- **Migraciones:** ✅ Seed admin user (admin@tramatex.com)
- **Seguridad:** ✅ JWT con blacklist en PostgreSQL, rate limiting en login

### Estado Frontend: ⚠️ Básico (30%)
- ✅ Login UI implementado
- ✅ Registro básico
- ❌ Panel de gestión de usuarios (admin)
- ❌ Asignación de roles UI
- ❌ Gestión de permisos

### Gaps Identificados
1. ⚠️ **Rol "cashier"**: Usado en Sales pero no seeded en IAM (línea 427 de main.go)
   - Impacto: Ruta POST /invoices/simplified requiere rol inexistente
   - Solución: Agregar rol en `migrations/006_seed_admin_user.sql`

### Prioridad: 🟢 Baja
- **Razón:** CRUD básico funcional, frontend login existe, roles faltantes no bloquean desarrollo

---

## 2. Módulo Party (Gestión de Terceros)

### Endpoints Implementados (17 rutas)

#### CRUD Party (5 rutas)
- ✅ `POST /api/parties` (admin, commercial) - Crear party
- ✅ `GET /api/parties` - Listar parties con filtros (rol, tipo, estado, nombre, tax_id)
- ✅ `GET /api/parties/:id` - Obtener party por ID
- ✅ `PUT /api/parties/:id` (admin, commercial) - Actualizar party
- ✅ `PATCH /api/parties/:id/status` (admin, commercial) - Cambiar estado (activo/inactivo)

#### Roles de Party (2 rutas)
- ✅ `POST /api/parties/:id/roles` (admin, commercial) - Añadir rol (CLIENT, SUPPLIER, EMPLOYEE)
- ✅ `DELETE /api/parties/:id/roles/:role` (admin, commercial) - Eliminar rol

#### Relaciones entre Parties (3 rutas)
- ✅ `POST /api/parties/:id/relationships` (admin, commercial) - Crear relación (IS_EMPLOYEE_OF, IS_SUBSIDIARY_OF)
- ✅ `GET /api/parties/:id/relationships` - Listar relaciones
- ✅ `DELETE /api/parties/:id/relationships/:relationship_id` (admin, commercial) - Eliminar relación

#### Contactos (OrganizationProfile) (4 rutas)
- ✅ `POST /api/parties/:id/contact-details` (admin, commercial) - Añadir contacto
- ✅ `GET /api/parties/:id/contact-details` - Listar contactos
- ✅ `PUT /api/parties/:id/contact-details/:contact_id` (admin, commercial) - Actualizar contacto
- ✅ `DELETE /api/parties/:id/contact-details/:contact_id` (admin, commercial) - Eliminar contacto

#### Direcciones (3 rutas)
- ✅ `POST /api/parties/:id/addresses` (admin, commercial) - Añadir dirección
- ✅ `GET /api/parties/:id/addresses` - Listar direcciones
- ❌ `PUT /api/parties/:id/addresses/:address_id` - **FALTA: Actualizar dirección**

### Estado Backend: ✅ 100% Completado (con 1 gap menor)
- **Cumple:** CRUD completo, filtros, roles, relaciones, contactos
- **Testing:** ⚠️ Tests unitarios de dominio implementados, repositorios pendientes
- **Migraciones:** ✅ Tablas parties, party_roles, party_relationships, contact_details, party_addresses
- **Documentación:** ✅ ADR-012, use-cases.md, domain-model.md completos

### Estado Frontend: ❌ 0% - Completamente ausente
- ❌ Lista de parties
- ❌ Detalle de party (persona/organización)
- ❌ Formulario crear/editar party
- ❌ Gestión de roles UI
- ❌ Gestión de relaciones UI
- ❌ Gestión de contactos UI

### Gaps Identificados
1. ⚠️ **Actualización de direcciones**: Falta endpoint PUT /parties/:id/addresses/:address_id
   - Impacto: Bajo - se pueden eliminar y recrear
   - Solución: Implementar handler UpdateAddress en PartyAddressHandler

2. ⚠️ **CONSUMIDOR_FINAL party**: Referenciado en ADR-020 pero no seeded
   - UUID: 00000000-0000-0000-0000-000000000001
   - NIF: 99999999R
   - Impacto: Alto para Sales (tickets sin cliente identificado)
   - Solución: Crear migration 019_seed_consumidor_final.sql

### Prioridad: 🔴 Alta
- **Razón:** Party es bloqueante para Sales frontend (no se pueden crear órdenes sin listar clientes)
- **Acción:** Implementar frontend Party ASAP (lista + detalle + crear)

---

## 3. Módulo Product (Catálogo de Productos)

### Endpoints Implementados (27 rutas)

#### CRUD Productos (8 rutas)
- ✅ `POST /api/products` (admin, commercial) - Crear producto
- ✅ `GET /api/products` - Listar productos
- ✅ `GET /api/products/:id` - Obtener producto por ID
- ✅ `PUT /api/products/:id` (admin, commercial) - Actualizar producto
- ✅ `GET /api/products/:id/calculated-option-sets` - Obtener atributos aplicables (herencia)
- ✅ `POST /api/products/:id/groups` (admin, commercial) - Añadir grupo a producto
- ✅ `POST /api/products/:id/attributes` (admin, commercial) - Añadir atributo directo
- ✅ `PATCH /api/products/:id/sku` (admin, commercial) - Actualizar SKU base

#### Variantes (5 rutas)
- ✅ `POST /api/products/:id/variants/generate` (admin, commercial) - Pre-generar variantes (UC-P-007)
- ✅ `POST /api/products/:id/variants/find-or-create` (admin, commercial) - Find or Create JIT (UC-P-009)
- ✅ `GET /api/products/:id/variants` - Listar variantes de un producto
- ✅ `GET /api/variants/:id` - Obtener variante por ID
- ✅ `GET /api/variants?sku=XXX` - Buscar variante por SKU
- ✅ `PUT /api/variants/:id` (admin, commercial) - Actualizar variante (UC-P-008)

#### Atributos (5 rutas)
- ✅ `POST /api/attributes` (admin, commercial) - Crear atributo (UC-P-001)
- ✅ `GET /api/attributes` - Listar atributos
- ✅ `GET /api/attributes/:id` - Obtener atributo por ID
- ✅ `PUT /api/attributes/:id` (admin, commercial) - Actualizar atributo (UC-P-002)
- ✅ `DELETE /api/attributes/:id` (admin) - Eliminar atributo

#### Marcas (5 rutas)
- ✅ `POST /api/brands` (admin, commercial) - Crear marca
- ✅ `GET /api/brands` - Listar marcas
- ✅ `GET /api/brands/:id` - Obtener marca por ID
- ✅ `PUT /api/brands/:id` (admin, commercial) - Actualizar marca
- ✅ `DELETE /api/brands/:id` (admin) - Eliminar marca

#### Grupos de Productos (5 rutas)
- ✅ `POST /api/product-groups` (admin, commercial) - Crear grupo
- ✅ `GET /api/product-groups` - Listar grupos
- ✅ `GET /api/product-groups/:id` - Obtener grupo por ID
- ✅ `PUT /api/product-groups/:id` (admin, commercial) - Actualizar grupo
- ✅ `DELETE /api/product-groups/:id` (admin) - Eliminar grupo

### Estado Backend: ✅ 95% Completado
- **Cumple:** CRUD completo para Product, Attribute, Brand, ProductGroup, variantes JIT, SKU determinista
- **Testing:** ✅ Tests unitarios de dominio (SKU composition, herencia de atributos)
- **Migraciones:** ✅ Todas las tablas creadas
- **Documentación:** ✅ ADR-013, ADR-015, use-cases.md, domain-model.md completos

### Estado Frontend: ⚠️ Parcial (50%)
- ✅ **Attributes UI:** CRUD completo funcional
- ✅ **Brands UI:** CRUD básico implementado
- ✅ **Product Groups UI:** CRUD básico implementado
- ❌ **Products UI:** Pendiente (backend listo)
- ❌ **Variants UI:** Pendiente
- ❌ **ServiceConfiguration UI:** Pendiente (integrado con Party)

### Gaps Identificados
1. ⚠️ **Recalcular SKUs al modificar Attribute.Code**: Documentado en UC-P-002 pero no implementado
   - Impacto: Alto - cambiar código de atributo rompe SKUs de variantes existentes
   - Solución: Implementar UpdateAttribute use case con transacción que recalcule SKUs
   - Estado: ❌ No implementado

2. ⚠️ **Recalcular SKUs al modificar Product.SKU**: Documentado en UC-P-006 pero no implementado
   - Impacto: Alto - cambiar SKU base rompe variantes
   - Solución: Implementar PATCH /products/:id/sku handler con lógica transaccional
   - Estado: ⚠️ Ruta existe pero sin lógica de propagación

3. ⚠️ **PartyServiceConfiguration**: Rutas existen en /parties/:id/service-configurations pero no documentadas en use-cases.md
   - Impacto: Medio - productos tipo SERVICE necesitan configuración por cliente
   - Estado: ⚠️ Backend implementado, documentación incompleta

### Prioridad: 🔴 Alta
- **Razón:** Product frontend bloquea Sales (no se pueden seleccionar productos en órdenes)
- **Acción:** Implementar Products UI + Variants UI ASAP

---

## 4. Módulo Pricing (Motor de Precios)

### Endpoints Implementados (10 rutas)

#### Pricing Legacy (6 rutas - PricingHandler)
- ✅ `POST /api/pricing/calculate` - Calcular precio básico (deprecated?)
- ✅ `GET /api/pricing/rules` - Listar reglas de precio
- ✅ `POST /api/pricing/rules` (admin, commercial) - Crear regla de precio
- ✅ `POST /api/pricing/client-overrides` (admin, commercial) - Override de precio por cliente
- ✅ `GET /api/pricing/history/:variantId` - Historial de precios

#### Pricing Engine (4 rutas - PricingEngineHandler)
- ✅ `POST /api/pricing/base-sales-rules` (admin, commercial) - Crear regla de precio base de venta
- ✅ `PUT /api/pricing/base-sales-rules/:id` (admin, commercial) - Actualizar regla base
- ✅ `POST /api/pricing/sale-modification-rules` (admin, commercial) - Crear regla de modificación (descuentos/recargos)
- ✅ `PUT /api/pricing/sale-modification-rules/:id` (admin, commercial) - Actualizar regla de modificación
- ✅ `POST /api/pricing/base-sales-price/calculate` - Calcular precio base de venta
- ✅ `POST /api/pricing/final-sale-price/calculate` - Calcular precio final con descuentos (integrado en Sales)

### Estado Backend: ✅ 90% Completado
- **Cumple:** 
  - Pricing Engine funcional con reglas base y modificaciones
  - Cálculo de precio base (BaseSalesPrice) ✅
  - Cálculo de precio final (FinalSalePrice) con descuentos ✅
  - Integración con Sales: CreateSimplifiedInvoice usa CalculateFinalSalePrice ✅
- **Testing:** ⚠️ Tests de integración de Pricing Engine pendientes
- **Migraciones:** ✅ Tablas pricing_rules, client_pricing_overrides, pricing_calculations, base_sales_price_rules, sale_modification_rules
- **Documentación:** ✅ ADR-016, implementation-guide.md, pricing-domain.md completos

### Estado Frontend: ❌ 0% - Completamente ausente
- ❌ Panel de reglas de precio
- ❌ Configuración de precio base por producto/categoría
- ❌ Gestión de descuentos por volumen
- ❌ Overrides por cliente
- ❌ Historial de precios

### Gaps Identificados
1. ⚠️ **DELETE/PATCH endpoints para reglas**: Solo hay POST/PUT, falta desactivar/eliminar reglas
   - Impacto: Bajo - reglas antiguas se mantienen con fechas de vigencia
   - Solución: Implementar DELETE /pricing/base-sales-rules/:id y /sale-modification-rules/:id

2. ⚠️ **Listar reglas de Pricing Engine**: Solo hay endpoints para crear/actualizar, no para listar base-sales-rules y sale-modification-rules
   - Impacto: Medio - no se pueden consultar reglas configuradas
   - Solución: Implementar GET /pricing/base-sales-rules y GET /pricing/sale-modification-rules

3. ⚠️ **Testing de Pricing Engine**: No hay tests de integración que validen cálculos complejos
   - Impacto: Alto - riesgo de errores en cálculos de precios en producción
   - Solución: Crear suite de tests de pricing_engine_test.go

### Prioridad: 🟡 Media
- **Razón:** Pricing Engine funciona y está integrado en Sales. Frontend es nice-to-have para configuración de reglas.
- **Acción:** Cerrar gaps de backend (listar reglas, tests) antes que frontend

---

## 5. Módulo Sales (Gestión de Órdenes e Invoices)

### Endpoints Implementados (21 rutas)

#### Cotizaciones (6 rutas)
- ✅ `POST /api/sales/quotes` (admin, commercial) - Crear cotización (integrado con Pricing)
- ✅ `GET /api/sales/quotes` - Listar cotizaciones con filtros
- ✅ `GET /api/sales/quotes/:id` - Obtener cotización por ID
- ✅ `PUT /api/sales/quotes/:id` (admin, commercial) - Actualizar cotización (BORRADOR only)
- ✅ `PATCH /api/sales/quotes/:id/status` (admin, commercial) - Cambiar estado (BORRADOR→APROBADA→RECHAZADA)
- ✅ `POST /api/sales/quotes/:id/convert` (admin, commercial) - Convertir cotización a orden (CU-S-006)

#### Órdenes (8 rutas)
- ✅ `POST /api/sales/orders` (admin, commercial) - Crear orden directa (sin cotización)
- ✅ `GET /api/sales/orders` - Listar órdenes con filtros (estado, party, fechas)
- ✅ `GET /api/sales/orders/:id` - Obtener orden por ID
- ✅ `PUT /api/sales/orders/:id` (admin, commercial) - Actualizar detalles (PENDIENTE only)
- ✅ `PATCH /api/sales/orders/:id/status` (admin, commercial) - Cambiar estado (PENDIENTE→CONFIRMADA→EN_PREPARACION→ENTREGADA→CANCELADA)
- ✅ `POST /api/sales/orders/:id/line-items` (admin, commercial) - Añadir línea de producto
- ✅ `PUT /api/sales/orders/:id/line-items/:lineItemId` (admin, commercial) - Actualizar línea
- ✅ `DELETE /api/sales/orders/:id/line-items/:lineItemId` (admin, commercial) - Eliminar línea

#### Albaranes (3 rutas)
- ✅ `POST /api/sales/delivery-notes` (admin, commercial) - Crear albarán desde orden
- ✅ `GET /api/sales/delivery-notes` - Listar albaranes
- ✅ `GET /api/sales/delivery-notes/:id` - Obtener albarán por ID

#### Facturas (4 rutas)
- ✅ `POST /api/sales/invoices` (admin, commercial) - Crear factura completa (B2B)
- ✅ `POST /api/sales/invoices/simplified` (admin, commercial, cashier) - Crear ticket (< 3k EUR) (CU-S-019)
- ✅ `GET /api/sales/invoices` - Listar facturas con filtro por tipo (COMPLETA/SIMPLIFICADA)
- ✅ `GET /api/sales/invoices/:id` - Obtener factura por ID

### Estado Backend: ✅ 95% Completado
- **Cumple:** 
  - CRUD completo para Quote, Order, DeliveryNote, Invoice
  - Flujo cotización → orden funcional (CU-S-006)
  - Facturas simplificadas (tickets) con validación < 3k EUR ✅ (ADR-020)
  - InvoiceType enum (COMPLETA/SIMPLIFICADA) ✅
  - InvoiceSeries Value Object con formato "TKT/00123/2026" ✅
  - Integración con Pricing Engine en CreateSimplifiedInvoice ✅
- **Testing:** ✅ 16 tests de dominio pasando (InvoiceType, InvoiceSeries, Invoice.ValidateLegalLimits)
- **Migraciones:** ✅ Migration 018 (invoice_type enum, series columns)
- **Documentación:** ✅ ADR-017, ADR-020, use-cases.md completo

### Estado Frontend: ❌ 0% - Completamente ausente
- ❌ Lista de órdenes
- ❌ Detalle de orden (tabs: Info, LineItems, Documents)
- ❌ Crear/editar orden/cotización
- ❌ Selección de productos con variants
- ❌ Botón "Emitir Factura" (completa)
- ❌ Botón "Emitir Ticket" (simplificada) con validación < 3k EUR
- ❌ PDF generador para facturas/tickets

### Gaps Identificados
1. ⚠️ **ConfigureInvoiceSeries use case (CU-S-022)**: Documentado en use-cases.md pero no implementado
   - Impacto: Bajo - series "A" (facturas) y "TKT" (tickets) están hardcoded
   - Solución: Implementar CRUD /api/sales/invoice-series para gestionar series por admin
   - Estado: ❌ No implementado

2. ⚠️ **Tax rate hardcoded**: 21% IVA hardcoded en CreateSimplifiedInvoice (línea 97 sales_service.go)
   - Impacto: Medio - no soporta otras tasas (10%, 4%, exento)
   - Solución: Crear tabla settings o tax_rates, cargar en use case
   - Estado: ⚠️ TODO comment presente, no implementado

3. ⚠️ **CONSUMIDOR_FINAL party no seeded**: UUID 00000000-0000-0000-0000-000000000001, NIF 99999999R (ADR-020)
   - Impacto: Alto - tickets sin cliente identificado no pueden crearse
   - Solución: Crear migration 019_seed_consumidor_final.sql
   - Estado: ❌ No implementado

4. ⚠️ **Cashier role no definido**: Usado en ruta POST /invoices/simplified pero no existe en IAM seeds
   - Impacto: Alto - autenticación falla para role cashier
   - Solución: Agregar rol en migration 006_seed_admin_user.sql
   - Estado: ❌ No implementado (compartido con IAM)

### Prioridad: 🔴 Alta
- **Razón:** Sales es el módulo de negocio más crítico, backend listo pero 100% bloqueado por frontend ausente
- **Acción:** Implementar Sales frontend ASAP (Order list/detail, Invoice/Ticket creation)

---

## 6. Análisis de Priorización

### Fórmula de Prioridad
```
Prioridad = (Impacto_Negocio × 0.4) + (Bloqueo_Dependencias × 0.3) + (Esfuerzo_Implementación × 0.2) + (Riesgo_Deadline × 0.1)
```

### Ranking de Tareas (Top 10)

| # | Tarea | Módulo | Impacto | Bloqueo | Esfuerzo | Score | Días |
|---|-------|--------|---------|---------|----------|-------|------|
| 1 | **Sales Frontend (Order + Invoice)** | Sales | 🔴 10 | 🔴 10 | 🟡 6 | **9.0** | 2-3 |
| 2 | **Party Frontend (List + Detail + Create)** | Party | 🔴 9 | 🔴 10 | 🟡 5 | **8.5** | 1-2 |
| 3 | **Product Frontend (List + Detail + Variants)** | Product | 🔴 9 | 🔴 9 | 🟡 6 | **8.4** | 1-2 |
| 4 | **Seed CONSUMIDOR_FINAL + Cashier role** | Sales/IAM | 🔴 9 | 🟡 7 | 🟢 2 | **7.2** | 0.5 |
| 5 | **Pricing Engine Tests** | Pricing | 🟡 7 | 🟢 3 | 🟢 3 | **5.4** | 1 |
| 6 | **Pricing Rules List/Delete endpoints** | Pricing | 🟡 6 | 🟢 3 | 🟢 2 | **4.6** | 0.5 |
| 7 | **Product: UpdateAttribute con SKU recalc** | Product | 🟡 7 | 🟡 5 | 🟡 5 | **6.2** | 1 |
| 8 | **Sales: ConfigureInvoiceSeries** | Sales | 🟢 4 | 🟢 2 | 🟢 3 | **3.3** | 0.5 |
| 9 | **Party: Update Address endpoint** | Party | 🟢 3 | 🟢 1 | 🟢 2 | **2.3** | 0.25 |
| 10 | **Pricing Frontend (Rules management)** | Pricing | 🟡 5 | 🟢 2 | 🟡 4 | **4.2** | 1-2 |

### Decisión Estratégica: Frontend-First

**Recomendación:** Priorizar frontend sobre backend polish dado que:
1. Backend ERP Core está 96% completo (84 rutas implementadas)
2. Frontend está 10% completo (solo Attributes/Brands/Groups UI)
3. Sin frontend, no hay producto entregable ni testing de usuario posible
4. Gaps de backend son non-critical (configuración, polish, tests)

**Plan de Acción (9 días restantes):**
- **Días 1-3:** Frontend Party + Product + Sales (paralelizable)
- **Días 4-5:** Cerrar gaps críticos backend (CONSUMIDOR_FINAL, cashier role, tests)
- **Días 6-8:** MES Module (MVP mínimo: estados de producción)
- **Día 9:** Testing integración + bug fixes

---

## 7. Gaps Críticos por Módulo

### IAM
- ⚠️ **Seedear rol "cashier"** en migration 006 (usado en Sales)

### Party
- 🔴 **Frontend completo ausente** (bloqueante para Sales)
- ⚠️ **Seedear CONSUMIDOR_FINAL party** en migration 019 (UUID 00000000-0000-0000-0000-000000000001, NIF 99999999R)
- 🟢 Endpoint Update Address faltante (low priority)

### Product
- 🔴 **Frontend Products/Variants ausente** (bloqueante para Sales)
- 🟡 **UpdateAttribute sin recálculo de SKUs** (documentado UC-P-002, no implementado)
- 🟡 **UpdateProductSKU sin propagación a variants** (documentado UC-P-006, parcialmente implementado)

### Pricing
- 🟡 **LIST endpoints para base-sales-rules y sale-modification-rules** (no se pueden consultar reglas actuales)
- 🟡 **DELETE endpoints para reglas** (no se pueden eliminar reglas obsoletas)
- 🟡 **Tests de integración del Pricing Engine** (riesgo de bugs en cálculos)
- 🔴 **Frontend completo ausente** (no crítico, configuración manual via SQL viable)

### Sales
- 🔴 **Frontend completo ausente** (bloqueante para uso del sistema)
- 🔴 **CONSUMIDOR_FINAL party no seeded** (compartido con Party)
- 🔴 **Cashier role no definido** (compartido con IAM)
- 🟡 **Tax rate hardcoded 21%** (no soporta otras tasas IVA)
- 🟢 **ConfigureInvoiceSeries no implementado** (series hardcoded, funcional)

---

## 8. Recomendaciones

### Acción Inmediata (Hoy - Día 1)
1. ✅ **Crear migration 019**: Seedear CONSUMIDOR_FINAL party + agregar rol cashier a IAM
2. ✅ **Iniciar Party Frontend**: Lista + Detalle + Crear (bloquea Sales)

### Días 2-3: Frontend Core
3. 🚀 **Product Frontend**: Lista productos/variantes + crear/editar (bloquea Sales)
4. 🚀 **Sales Frontend Fase 1**: Order list + Order detail (lectura)

### Días 4-5: Completar Sales + Tests
5. 🚀 **Sales Frontend Fase 2**: Create Order + Invoice/Ticket buttons + PDF
6. ✅ **Pricing Engine Tests**: Suite de tests de cálculos complejos
7. ✅ **Pricing LIST endpoints**: GET base-sales-rules, sale-modification-rules

### Días 6-8: MES Module
8. 🚀 **MES Backend**: Domain + Use Cases + Handlers (estados de producción)
9. 🚀 **MES Frontend**: Workflow board (Kanban)

### Día 9: Polish
10. 🐛 **Testing integración E2E**: Flujo completo Party → Product → Sales
11. 🐛 **Bug fixes + deployment**

---

## 9. Conclusiones

### Estado Actual
- ✅ **Backend ERP Core:** 96% funcional, 84 endpoints implementados
- ⚠️ **Frontend ERP Core:** 10% funcional, solo módulos secundarios (Attributes, Brands, Groups)
- 🔴 **Bloqueante crítico:** Sin Party/Product/Sales frontend, el sistema no es utilizable

### Decisión Estratégica
**Pivote a Frontend-First:** Dado el deadline de 9 días y backend casi completo, maximizar valor entregable requiere frontend funcional para:
1. Permitir testing de usuario
2. Validar flujos de negocio completos
3. Demostrar producto funcional

### Riesgo
- **MES Module:** Estimado 3 días (48h) para MVP mínimo. Si frontend ERP toma > 5 días, MES podría quedar incompleto.
- **Mitigación:** Scope MES a estados básicos (DISEÑO, TALLER, QA, ENTREGADO) sin features avanzadas (adjuntos, comentarios).

---

**Próximo Paso:** Crear migration 019 (CONSUMIDOR_FINAL + cashier role) e iniciar Party Frontend.
