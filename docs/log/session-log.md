# Bitácora de Sesiones de Desarrollo

<!--
Este archivo registra las sesiones de desarrollo.

SECCIONES:
1. SESIONES ABIERTAS: Contiene las sesiones de trabajo que están en progreso, pausadas o bloqueadas. El objetivo es detallar el contexto y los próximos pasos.
2. REGISTRO DE SESIONES CERRADAS: Un archivo histórico de todas las sesiones completadas, conservando solo metadatos esenciales.

ESTRUCTURA DE UNA SESIÓN ABIERTA:
- Título (##): Un H2 con un título descriptivo y único.
- Metadatos:
  - **Session ID:** `identificador-unico-kebab-case` (OBLIGATORIO Y ÚNICO)
  - **Status:** (En Progreso | En Pausa | Bloqueado)
  - **Sprint:** Sprint XX
  - **Started:** Fecha de inicio (YYYY-MM-DD).
- Contexto: Breve descripción del objetivo de la sesión.
- Próximos Pasos: Checklist de las tareas pendientes.
- Archivos de Contexto: Lista de archivos clave.

ESTRUCTURA DE UNA SESIÓN CERRADA (en el registro):
- Una línea de lista con: **[Título]** | Iniciada: [Fecha YYYY-MM-DD] | Finalizada: [Fecha YYYY-MM-DD]
-->
---
# SESIONES ABIERTAS

## Análisis de Refinamiento Arquitectónico del MVP

- **Session ID:** `mvp-refinement-analysis-2026-03-12`
- **Status:** En Progreso (Análisis de MES en curso)
- **Sprint:** N/A
- **Started:** 2026-03-12

### Contexto

Análisis modular sistemático del backend para identificar oportunidades de mejora (simplificación, desacoplamiento, rendimiento) antes del lanzamiento a producción del MVP. Se genera el archivo maestro `tmp/mvp_refinement_proposals.md`.

### Trabajo Completado

- [x] **Fase 0: Preparación**: Creación de `tmp/mvp_refinement_proposals.md` con scaffolding.
- [x] **Módulo 1: IAM**: Finalizado (11 propuestas: auditoría, inyección, roles en shared, etc.).
- [x] **Módulo 2: Party**: Finalizado (14 propuestas: integridad de direcciones, NIF/CIF, performance N+1, soft delete).
- [x] **Módulo 3: Product**: Finalizado (14 propuestas: Money, SmartSearch en BD, integridad referencial, fragmentación de servicio).
- [x] **Módulo 4: Pricing**: Finalizado (10 propuestas: reglas en SQL, bulk loading, histórico de cálculos, unificación de Percentage).
- [x] **Módulo 5: Sales (Ventas y Facturación)**: Finalizado (10 propuestas: fragmentación de servicio, unificación de Money, abstracción de cálculos).
- [x] **Módulo 6: MES (Producción)**: Finalizado (8 propuestas de refinamiento detalladas).
- [x] **Módulo 7: Shared**: Finalizado (4 propuestas: consolidación de VOs, middleware de auditoría, estandarización de paginación).
- [x] **Matriz de Dependencias**: Completada al finalizar todos los módulos.

### Próximos Pasos

- [ ] **Priorización**: Revisar el documento `tmp/mvp_refinement_proposals.md` para priorizar la implementación de refinamientos.
- [ ] **Ejecución**: Iniciar el refactoring de los módulos críticos según la prioridad establecida.

### Archivos de Contexto

- `tmp/mvp_refinement_proposals.md` (Registro maestro de propuestas)
- `docs/architecture/adrs/`
- `apps/tramatex-api/internal/`

---

## Implementación de Infraestructura de Despliegue Multientorno

- **Session ID:** `infra-multi-env-deployment-impl-2026-03-10`
- **Status:** En Pausa (Pendiente de inicio)
- **Sprint:** N/A
- **Started:** 2026-03-10

### Contexto

Implementación técnica de la estrategia de despliegue multientorno definida en el estudio previo (`tmp/estudio_despliegue_tramatex.md`). El objetivo es configurar la jerarquía de ramas (develop -> staging -> master), la orquestación con Nginx y la automatización CI/CD hacia DigitalOcean.      

### Próximos Pasos

- [ ] Crear rama `infra/multi-env-deployment`.
- [ ] Implementar `Dockerfile.frontend` (multi-stage build).
- [ ] Crear configuración `docker/nginx.conf` (Proxy inverso y SPA routing).
- [ ] Actualizar `docker/docker-compose.remote.yml` para incluir el servicio Nginx.
- [ ] Crear Workflows de GitHub Actions para despliegue automático en `staging` y `master`.
- [ ] Refinar `Makefile` para soportar perfiles de despliegue (`pcele`, `staging`, `prod`).

### Archivos de Contexto

- `tmp/estudio_despliegue_tramatex.md`
- `docker/docker-compose.remote.yml`
- `Dockerfile`
- `Makefile`

---

## Refinamiento y Estabilización ERP Core

- **Session ID:** `erp-core-refinement-2026-03-09`
- **Status:** ✅ Cerrada (2026-03-14)
- **Sprint:** N/A
- **Started:** 2026-03-09

### Contexto

Revisión funcional y estabilización de los 4 módulos del ERP Core (Party, Product, Pricing, Sales). El objetivo es verificar el correcto funcionamiento end-to-end, corregir bugs encontrados y refinar la experiencia de usuario.

### Trabajo Completado

- **Party Module:** Revisado y validado.
- **Product Module:** Revisado y validado.
- **Pricing Module:** Revisado y validado.
- **Sales Module:** Revisión completa (domain, application, infrastructure, HTTP, frontend). Se aplicaron 8 fixes en 3 fases:
  - **Fase 1 (Bugs funcionales):** TaxRate en conversión Quote→Order (`convert.go`), availableQuantity en albaranes parciales (`OrderDetail.vue`), minDeliveryDate = mañana (`OrderDetail.vue`, `OrderCreate.vue`).
  - **Fase 2 (Integridad de datos y UX):** Redondeo Money a 2 decimales (`money.go`), aviso de recálculo fiscal en edición (`QuoteDetail.vue`, `OrderDetail.vue`), botón "Crear Factura" en detalle de pedido (`OrderDetail.vue`).
  - **Fase 3 (Descubierto en E2E):** Fix #8 — Cálculo de impuestos en factura (`invoice.go`): `NewInvoice()` y `RecalculateTotals()` ahora computan `taxAmount` sumando los impuestos de cada línea en vez de usar el parámetro externo (que llegaba como 0). Se añadió helper `sumInvoiceLineItemTaxAmounts()`.
  - **Tests:** Corregidos 4 tests desalineados: 2 pre-existentes en `quote_test.go` y `sales_service_test.go`, 2 en `invoice_test.go` (actualizados para usar taxRate en líneas).
  - **Dato de seed:** Migración `012_add_pricing_rules_and_consumidor_final.sql` usa rol `CUSTOMER` pero el dominio espera `CLIENT`. Pendiente de corregir.
  - **Fase 4 (Testing manual):** Fix #9 — Import de `partyApi` en 6 páginas Sales (QuoteList, OrderList, InvoiceList, DeliveryNoteList, DeliveryNoteDetail, QuoteDetail): se usaba `import partyApi from` (default export = clase) en vez de `import { partyApi } from` (named export = instancia). La llamada a métodos sobre la clase lanzaba TypeError → catch mostraba "Error al cargar" en columna de cliente.
  - **Fase 4 (UX):** Fix #10 — Refactoring de `OrderDetail.vue` para unificar UX de edición de líneas: reemplazado sistema modal (Add/Edit Line Item) por edición inline en tabla (mismo patrón que QuoteDetail y OrderCreate). Vista lectura sin botones de acción, vista edición con inputs inline (cantidad, precio, descuento) y VariantSelector modal para agregar nuevas líneas. `saveOrderHeader()` ahora sincroniza líneas completas (add/update/remove) contra la API.
  - **Fase 4 (UX):** Fix #11 — Cálculo en tiempo real de subtotales y totales en modo edición de `QuoteDetail.vue` y `OrderDetail.vue`. Añadidas funciones `calculateEditLineSubtotal()`, computed `editCalculatedTotals`, columna "Subtotal" por línea y sección "Resumen de Totales" (Subtotal + IVA estimado 21% + Total estimado). Mismo patrón que las páginas de creación. Se captura `effectiveUnitPrice` del backend para líneas con precio automático.
- **Validación E2E (datos reales via API):** Flujo completo Quote → Order → DeliveryNote → Invoice ejecutado contra el entorno Docker local. Resultado:
  ```
  QUOTE:   CONVERTIDA_A_PEDIDO | Sub=292.50 Tax=61.43 Tot=353.93
  ORDER:   FACTURADO_COMPLETAMENTE | Sub=292.50 Tax=61.43 Tot=353.93
  INVOICE: PAGADA               | Sub=292.50 Tax=61.43 Tot=353.93
  MATCH:   ✅ Totales consistentes en todos los documentos
  ```
  Todos los tests (domain, application, HTTP) pasan.

- **Sesión 2026-03-11 — Investigación bug descuento y diagnóstico login:**
  - **Bug descuento (5 → 5,02):** Trazado completo del flujo descuento (frontend → API → dominio → persistencia → BD). Se identificaron 3 posibles causas raíz:
    1. **Cálculo inverso de porcentaje** (sospecha principal): Cuando `DiscountPercent` es nil (descuento viene del pricing engine), se calcula inversamente como `(calculatedDiscount.Amount() / finalUnit.Amount()) * 100`. Ambos valores `Money` están redondeados independientemente a 2 decimales, produciendo imprecisión. Archivos: `sales_service.go` líneas ~1108 y ~1195.
    2. **Banker's rounding en `resolveDiscountFromPercent`**: `math.Round(amount*100)/100` redondea `.5` al par más cercano (ej: 5.025 → 5.02 en lugar de 5.03). Archivo: `money.go`.
    3. **Backfill migración 015**: `SET discount_percent = ROUND((discount_per_unit_amount / unit_price_amount) * 100, 2)` pudo corromper datos existentes.
  - **Pendiente de confirmar:** Necesita reproducción concreta — ¿en qué página (Quote/Order), qué producto/precio, y si el "5,02" es en la columna porcentaje o la monetaria?
  - **Problema login admin:** Se verificó que la API responde correctamente a `POST /auth/login` con credenciales `admin@tramatex.local` / `admin123` (devuelve JWT válido). El problema era que el **frontend Vite no estaba arrancado**. Se inició el servidor de desarrollo en puerto 5173.        

- **Sesión 2026-03-11 (cont.) — Eliminación de cálculos de negocio en frontend:**
  Se aplicó la directriz arquitectónica "el front-end no debería calcular nada" a todos los módulos afectados. Todos los cálculos monetarios/fiscales ahora se derivan exclusivamente del backend.
  - **Pricing.vue:** Las funciones `getSubtotal`, `getDiscountAmount`, `getTaxBase`, `getTaxAmount` y `calcLineTotal` ahora derivan valores del `CalculatedSaleItem` del backend (`finalPrice`, `finalPriceWithTax`) en lugar de recalcular con `baseSalesPrice * qty * (1 - discount)`.
  - **Endpoint Order Preview (nuevo):** Creado `POST /api/sales/orders/preview` (full stack):
    - Backend: `PreviewOrderCommand`, `OrderPreviewDTO`, `PreviewOrderCalculation()` en `sales_service.go`, handler en `sales_handler.go`, ruta en `main.go`.
    - Frontend: `previewOrderCalculation()` en `salesApi.ts`.
  - **OrderCreate.vue:** Refactorizado con patrón debounced preview (400ms). `calculateLineSubtotal` y `calculateTotals` consultan `previewResult` del backend. Añadidos `buildPreviewItems()` y `fetchPreviewCalculation()`.
  - **OrderDetail.vue:** Mismo patrón debounced preview. Vista lectura: eliminadas columnas "IVA %" e "IVA línea". Vista edición: añadida columna "P. Tarifa", loading state en subtotales, watcher en `editLineItems`. Labels normalizados ("IVA:", "Total:").
  - **TicketCreate.vue:** Cambiado de `calculatePrice` a `calculateFinalSalePrice`. `ticketTaxAmount` ahora usa `(finalPriceWithTax - unitPrice) * qty` en lugar de recalcular con `taxRate`.
  - **Auditoría UI Order vs Quote:** Labels unificados ("IVA:" en vez de "IVA estimado:", "Total:" en vez de "Total estimado:"), loading states añadidos, columna "P. Tarifa" en edición de pedidos.
  - **Tests Go:** Todos pasan (sales handler, application, domain). Solo falla pre-existente `TestGORMRepositories_Sales` (columna `mes_work_refs`).
  - **Build frontend:** `vite build` OK — 0 errores de compilación en los 10 archivos modificados.
  - **Pendiente:** Comprobación manual del usuario en navegador (Presupuestos, Pedidos, Tickets, Pricing).

### Próximos Pasos

- [x] Revisión y refinamiento del módulo Sales.
- [x] Corrección de bugs encontrados durante la revisión.
- [x] Validación de flujos end-to-end (Quote → Order → DeliveryNote → Invoice) con datos reales.
- [x] **Bug descuento 5→5,02:** Corregido (2026-03-11). Causa raíz: banker's rounding en `money.go` + cálculo inverso de % sobre valores Money pre-redondeados en `sales_service.go`. Fix: round-half-up comercial (`math.Floor(amount*100+0.5)/100`) y cálculo de % desde floats crudos del DTO de pricing.
- [x] Corregir seed data: rol `CUSTOMER` → `CLIENT` en migración 012 (2026-03-11). Creada migración correctiva 016.
- [x] **Eliminación cálculos frontend (2026-03-11):** Pricing.vue, OrderCreate.vue, OrderDetail.vue, TicketCreate.vue refactorizados. Nuevo endpoint `POST /orders/preview`. Build OK, Go tests OK.
- [x] **Rebuild Docker con cambios (2026-03-12):** Imagen Docker reconstruida con todos los cambios (endpoint `/orders/preview`, refactoring Pricing→Sales, 18 fixes). Ambos contenedores healthy.
- [x] **Comprobación manual en navegador:** Verificado por el usuario — todo correcto.
- [x] Fase 3 (Robustez): Race conditions en `CreateDeliveryNote`, auth middleware, tipos de error consistentes, transacciones atómicas para facturación.

- **Sesión 2026-03-13 — Fase 3: Robustez y Transacciones Atómicas:**
  Implementación completa de las 4 áreas de robustez identificadas:
  - **Área 1 — Race conditions en CreateDeliveryNote:** `CreateDeliveryNote` ahora ejecuta toda la lógica dentro de una transacción (`runInTransaction`) y usa `FindByIDForUpdate` (`SELECT ... FOR UPDATE`) para adquirir bloqueo pesimista en la fila del pedido, serializando solicitudes concurrentes y evitando sobreentrega.
  - **Área 2 — Auth middleware:** Confirmado que ya está correctamente implementado (`AuthMiddleware` con JWT + `RoleMiddleware`). Sin cambios necesarios.
  - **Área 3 — Tipos de error consistentes:** Reemplazados 11 `fmt.Errorf` con errores de dominio tipados: `NewConfigurationError` (5 instancias — config/infraestructura), `NewValidationError` (6 instancias — datos). Nuevo `ErrCodeConfiguration` en `errors.go`, mapeado a HTTP 500 en `handleSalesError`.
  - **Área 4 — Transacciones atómicas en CreateInvoice:** `CreateInvoice` ahora ejecuta todo dentro de una transacción. `fetchOrdersForInvoice` usa `FindByIDForUpdate` para bloqueo. Factura se guarda primero, luego se actualizan estados de pedidos — todo atómico.
  - **Infraestructura transaccional (nuevo):**
    - `TransactionManager` interface en `application/transaction.go`
    - `GORMTransactionManager` + `getDB(ctx, db)` context-propagation en `infrastructure/persistence/transaction.go`
    - Todos los repos (`r.db.WithContext(ctx)` → `getDB(ctx, r.db)`) participan transparentemente en transacciones service-level
    - `SetTransactionManager()` setter backward-compatible (nil = no transaction wrapping para tests)
  - **Tests:** Todos pasan (application 5/5, handler 21/21). Mocks actualizados con `FindByIDForUpdate`. Fallos pre-existentes en domain (`InvoiceSeries_FormatNumber`, `CanceledToAnything`) y persistencia (requiere DB) no relacionados.
  - **Build:** `go build ./...` OK. Docker rebuild OK, ambos contenedores healthy.
  - **Archivos nuevos:** `application/transaction.go`, `infrastructure/persistence/transaction.go`
  - **Archivos modificados:** `domain/errors.go`, `domain/repository.go`, `application/sales_service.go`, `infrastructure/persistence/repositories.go`, `infrastructure/persistence/number_generator.go`, `interfaces/http/handler/sales_handler.go`, `cmd/api/main.go`, `application/sales_service_test.go`, `interfaces/http/handler/sales_handler_test.go`

- **Sesión 2026-03-12 — Refactoring Pricing como fuente única de cálculos:**
  La observación del usuario ("cuando hay un cálculo sólo tiene que existir una lógica") derivó en un análisis arquitectónico profundo y un refactoring en 3 fases para eliminar la duplicación de cálculos entre Pricing y Sales:

  - **Fase 1 — Enriquecer DTO de Pricing:** Se añadieron 6 campos pre-calculados a `CalculatedSaleItemResponse` en `adr15_dtos.go`: `DiscountPercent`, `DiscountAmount`, `TaxAmountPerUnit`, `LineSubtotal`, `LineTaxAmount`, `LineTotal`. El servicio `pricing_engine_service.go` ahora computa todos estos valores, convirtiendo a Pricing en la fuente única de verdad para cálculos monetarios.
  - **Fase 2 — Unificar Money (redondeo comercial):** Se añadió `roundTo2Decimals` (redondeo comercial `math.Floor(amount*100+0.5)/100`) al `NewMoney` de `pricing/domain/money.go`, igualando el comportamiento con Sales. Elimina la discrepancia de redondeo que causaba diferencias entre edición y vista de detalle.
  - **Fase 3 — Simplificar Sales (patrón dual-path):**
    - Nuevos constructores `NewQuoteLineItemFromCalculated()` y `NewOrderLineItemFromCalculated()` en el dominio Sales que aceptan valores pre-calculados de Pricing sin re-derivarlos.
    - `buildQuoteLineItems()` y `buildOrderLineItemsFromSeeds()` reescritos con patrón dual: si hay override de usuario (precio o descuento) → constructor clásico con recálculo; sin override → constructor `FromCalculated` con valores de Pricing (fuente única).
    - Eliminado `deriveCalculatedPrices()` (código muerto — ya no se re-derivan descuentos inversamente).
    - Eliminado `toDomainMoneyPtr()` (código muerto). Añadido `toDomainMoney()` para conversión DTO→dominio.
  - **Tests:** Compilación OK. Todos los tests pasan (pricing 100%, sales handler 21/21). Único fallo pre-existente: `TestGORMRepositories_Sales` (columna `mes_work_refs`).
  - **Archivos modificados:** `pricing/domain/money.go`, `pricing/application/adr15_dtos.go`, `pricing/application/pricing_engine_service.go`, `sales/domain/quote.go`, `sales/domain/sales_order.go`, `sales/application/sales_service.go`.

- **Sesión 2026-03-12 (cont.) — Numeración secuencial de documentos:**
  Se reemplazó el sistema de numeración basado en timestamp+UUID (ej: `Q-20260304-173913-bf79f795`) por numeración secuencial legible: `PREFIJO-AÑO-SECUENCIAL`.
  - **Prefijos:** `PRE` (presupuestos), `PED` (pedidos), `ALB` (albaranes), `FV` (facturas venta), `FT` (facturas ticket).
  - **Backend:** `SequentialNumberGenerator` reemplaza a `TimeBasedNumberGenerator`. Usa tabla `document_sequences` con `INSERT ... ON CONFLICT DO UPDATE` atómico. Interfaz `NextInvoiceNumber` ahora recibe `InvoiceSeries` para soporte de series diferenciadas.
  - **Series facturas:** Cambiadas de "A"→"FV" y "TKT"→"FT".
  - **DTO:** `InvoiceDTO` enriquecido con `invoiceType` y `seriesCode`.
  - **Migración 017:** Crea tabla `document_sequences`, inicializa contadores desde documentos existentes, migra series de facturas.
  - **Documentación:** Actualizada `module-spec.md` (tabla de prefijos), `use-cases.md` (series FV/FT).
  - **Tests:** Compilación OK, todos los tests pasan (mocks actualizados para nueva interfaz).

- **Sesión 2026-03-12 (cont.) — Correcciones funcionales (Fase 9: 3 fixes):**
  El usuario reportó 3 problemas detectados en pruebas manuales:
  - **Fix #12 — Crear factura desde albarán:** No existía botón en `DeliveryNoteDetail.vue`. Se añadió botón "📄 Crear Factura" (visible cuando `sttatus === 'DELIVERED' && !relatedInvoice`), función `createInvoiceFromDeliveryNote()` que envía `{ partyId, deliveryNoteIds: [id], invoiceDate, dueDate }` y navega a la nueva factura. CSS añadido para `btn-primary`, `btn-success`, `btn-danger`.
  - **Fix #13 — Descuentos en tickets/facturas simplificadas:** Implementación full-stack. Backend: `DiscountPercent float64` en `OrderLineItemInputSimplified` (`commands.go`), mapa `discountByVariant` y cálculo de `discountAmount` en `CreateSimplifiedInvoice` (`sales_service.go`). Frontend: columna "Dto. %" con input 0-100, helper `lineTotal()` que aplica descuento, `discountPercent` en `CreateSimplifiedInvoiceRequest` (`sales.ts`), recibo muestra descuento por línea (`TicketCreate.vue`).
  - **Fix #14 — Ticket no se imprime:** El CSS `@media print` ocultaba `.ticket-create-container` (contenedor padre), ocultando TODO incluido el modal de recibo. Fix: ocultar solo hijos específicos (`.page-header`, `.form-card`), `padding: 0` en contenedor print.
  - **Build:** Backend `go build` OK, frontend `vite build` OK (17.80s), Docker rebuild OK, ambos contenedores healthy.

- **Sesión 2026-03-12 (cont.) — Correcciones funcionales (Fase 10: 4 fixes):**
  Segunda ronda de pruebas manuales del usuario, 4 problemas adicionales:
  - **Fix #15 — Listado de facturas vacío:** El backend envía campos camelCase (`invoiceType`, `invoiceDate`, `relatedOrderIds`) pero el frontend esperaba (`type`, `issueDate`, `salesOrderIds`). Fix: añadido mapeo de campos invoice en `normalizeEntity()` de `salesApi.ts`.
  - **Fix #16 — "Notas Internas" → "Observaciones":** Renombrado en `QuoteDetail.vue` (heading + placeholder) y `QuoteCreate.vue` (label + placeholder).
  - **Fix #17 — Reactivar pedido cancelado:** Backend: añadido `case SalesOrderStatusCanceled: return to == SalesOrderStatusPending` en `canTransitionOrder()` (`statuses.go`). Frontend: botón "♻️ Reactivar Pedido" + función `reactivateOrder()` en `OrderDetail.vue`.
  - **Fix #18 — Cantidades incorrectas en albaranes desde pedido:** Dos bugs en `OrderDetail.vue`: (1) `initializeDeliveryNoteForm()` comparaba con `'CANCELADO'` (español) en vez de `'CANCELLED'` (normalizado), haciendo que albaranes cancelados restaran cantidades incorrectamente; (2) modo TOTAL enviaba items con `availableQuantity === 0`, rechazados por backend. Fix: cambio a `'CANCELLED'` y filtro `item.availableQuantity > 0` en modo TOTAL.
  - **Build:** Backend `go build` OK, frontend `vite build` OK, Docker rebuild OK, ambos contenedores healthy.
  - **Pendiente:** Comprobación manual del usuario en navegador.

### Archivos de Contexto

- `apps/tramatex-api/internal/sales/`
- `apps/tramatex-api/internal/sales/domain/money.go` (round-half-up comercial)
- `apps/tramatex-api/internal/sales/domain/invoice_series.go` (FormatNumber: PREFIX-YEAR-NNNN)
- `apps/tramatex-api/internal/sales/application/sales_service.go` (series FV/FT + Pricing refactor)
- `apps/tramatex-api/internal/sales/application/dtos.go` (InvoiceDTO con invoiceType/seriesCode)
- `apps/tramatex-api/internal/sales/infrastructure/persistence/number_generator.go` (SequentialNumberGenerator)
- `apps/tramatex-api/internal/sales/interfaces/http/handler/sales_handler.go` (PreviewOrderCalculation handler)
- `apps/tramatex-api/cmd/api/main.go` (wiring SequentialNumberGenerator con DB)
- `apps/tramatex-api/migrations/017_document_sequences.sql`
- `apps/frontend/src/types/sales.ts` (series_code en Invoice)
- `apps/frontend/src/services/salesApi.ts` (previewOrderCalculation)
- `apps/frontend/src/pages/sales/`
- `docs/modules/sales/module-spec.md` (numeración secuencial documentada)
- `docs/modules/sales/use-cases.md` (series FV/FT)
- `docs/log/erp-core-completion.md`
- `apps/frontend/src/pages/sales/DeliveryNoteDetail.vue` (botón crear factura desde albarán)
- `apps/frontend/src/pages/sales/TicketCreate.vue` (descuentos + fix print CSS)
- `apps/tramatex-api/internal/sales/application/commands.go` (DiscountPercent en simplified)
- `apps/tramatex-api/internal/sales/domain/statuses.go` (transición CANCELLED→PENDING)
- `apps/frontend/src/pages/sales/OrderDetail.vue` (reactivar pedido + fix albaranes)
- `apps/frontend/src/pages/sales/QuoteDetail.vue` (Observaciones)
- `apps/frontend/src/pages/sales/QuoteCreate.vue` (Observaciones)

- **Sesión 2026-03-14 — Fase 5: Trazabilidad Factura↔Albarán:**
  Detectado bug: al crear una factura desde un albarán parcial, no existe relación explícita entre facturas y albaranes en la BD. El frontend en `DeliveryNoteDetail.vue` busca facturas por `salesOrderId` (pedido completo), lo que provoca que un albarán parcial no facturado muestre como "relacionada" una factura que solo cubre otro albarán del mismo pedido, ocultando el botón "Crear Factura".
  - **Diagnóstico:** Las tablas `invoices` e `invoice_line_items` no tienen referencia a `delivery_notes` ni a `delivery_note_line_items`. La relación solo se infiere indirectamente vía `sales_order_line_item_id`.
  - **Reglas de negocio confirmadas:**
    1. Una línea de albarán siempre pertenece a exactamente una factura (1:1 línea DN → factura).
    2. Post-MVP: varias líneas de albaranes diferentes podrán consolidarse en una sola línea de factura (N:1 líneas DN → línea factura) cuando coincidan producto, precio y descuento.
  - **Solución en curso:** Añadir campo `invoice_line_item_id` en `delivery_note_line_items` para trazabilidad directa. Infraestructura preparada para soportar consolidación N:1 Post-MVP.
  - **Decisiones de diseño documentadas (2026-06-05):**
    - MVP: facturas solo desde albaranes (no desde pedidos). 1 factura = 1 albarán. Relación DN-line → Invoice-line 1:1.
    - Post-MVP: facturación consolidada multi-albarán (N líneas DN → 1 línea factura cuando coincidan producto+precio+descuento). Ver CU-S-025 en `use-cases.md`.
    - CU-S-018 renombrado a `CreateInvoiceFromDeliveryNote` en `use-cases.md`.
    - Añadida Fase 6 (Post-MVP) en `module-spec.md`.
    - Sección 13 "Facturación Consolidada Multi-Albarán" añadida en `post-mvp-roadmap.md`.
    - Modelo de dominio actualizado en `domain-model.md` (origen exclusivo desde albaranes).

### Próximos Pasos (Fase 5)

- [x] Crear migración: `019_add_dn_invoice_traceability.sql` — `ALTER TABLE delivery_note_line_items ADD COLUMN invoice_line_item_id UUID REFERENCES invoice_line_items(id)`.
- [x] Backend: actualizar `buildInvoiceItemsFromDeliveryNotes` para escribir el link DN-line → Invoice-line.
- [x] Backend: nuevo método repositorio `FindByDeliveryNoteID` / `ListDeliveryNoteIDsByInvoiceID` / `LinkLineItemsToInvoice`.
- [x] Backend: enriquecer `InvoiceDTO` con `RelatedDeliveryNoteIDs` y `DeliveryNoteDTO` con `InvoiceID`.
- [x] Frontend: `DeliveryNoteDetail.vue` — buscar factura por albarán concreto (`deliveryNoteId`), no por pedido.
- [x] Frontend: `InvoiceDetail.vue` — mostrar albaranes relacionados reales con links navegables.
- [ ] Tests: validar flujo de facturación parcial con múltiples albaranes (e2e/integration).

### Archivos de Contexto (Fase 5)

- `apps/tramatex-api/internal/sales/application/sales_service.go` (buildInvoiceItemsFromDeliveryNotes)
- `apps/tramatex-api/internal/sales/infrastructure/persistence/models.go` (DeliveryNoteLineItemDataModel)
- `apps/tramatex-api/internal/sales/infrastructure/persistence/repositories.go` (invoice/DN repos)
- `apps/tramatex-api/internal/sales/domain/invoice.go` (Invoice domain)
- `apps/tramatex-api/internal/sales/domain/delivery_note.go` (DeliveryNote domain)
- `apps/tramatex-api/internal/sales/application/dtos.go` (InvoiceDTO, DeliveryNoteDTO)
- `apps/tramatex-api/migrations/005_init_sales.sql` (esquema actual)
- `apps/frontend/src/pages/sales/DeliveryNoteDetail.vue` (loadRelatedInvoice)
- `apps/frontend/src/pages/sales/InvoiceDetail.vue` (albaranes relacionados)

---

## Refinamiento y Estabilización MES

- **Session ID:** `mes-refinement-2026-03-09`
- **Status:** En Progreso
- **Sprint:** N/A
- **Started:** 2026-03-09

### Contexto

Revisión funcional y estabilización del módulo MES (Manufacturing Execution System). La primera tarea es un refactoring de los nombres de las secciones del módulo (navegación, títulos de página, breadcrumbs).

**Estado actual del módulo MES:**
- **Entidades de dominio:** Task, Position, ServiceGroup, MESWork, MESWorkServiceGroup, MESWorkTask
- **Frontend URLs:** `/mes/tasks`, `/mes/positions`, `/mes/service-groups`, `/mes/work-definitions`, `/mes/terminal`
- **Breadcrumbs actuales (ES):** "MES / Datos Maestros" (tasks, positions, service-groups), "MES / Definiciones de trabajo" (works), "MES / Terminal" (tablet)
- **Backend API:** `/api/mes/tasks`, `/positions`, `/service-groups`, `/works`, `/works/{id}/tasks/{id}`, dashboard stats

**Rama git:** `mes-refactor` (creada desde `develop` el 2026-03-14)

### Próximos Pasos

- [ ] Refactoring de nombres de secciones del módulo MES (navegación, títulos, breadcrumbs).
- [ ] Validar integración Sales → MES (generación de órdenes de producción).
- [ ] Revisión funcional del módulo MES (maestros, órdenes, dashboard).
- [ ] Evaluar diseño futuro: creación de cabeceras MESWorkDefinition desde Sales.

### Archivos de Contexto

- `apps/tramatex-api/internal/mes/` (backend completo)
- `apps/frontend/src/pages/mes/` (frontend: Dashboard, Tasks, Positions, ServiceGroups, Works, Terminal)
- `apps/frontend/src/components/layout/Navbar.vue` (menú de navegación MES)
- `apps/frontend/src/services/mesApi.ts` (API client)
- `docs/modules/mes/` (documentación: README, module-spec, domain-model, use-cases, api-contracts)

---
# REGISTRO DE SESIONES CERRADAS
---
- **Refinamiento y Estabilización ERP Core** | Iniciada: 2026-03-09 | Finalizada: 2026-03-14 (Todos los módulos validados, Sales completo con tickets TPV, trazabilidad factura↔albarán, robustez transaccional, documentación alineada)
- **Generación de Entorno de Despliegue de Desarrollo (Dev)** | Iniciada: 2026-03-09 | Finalizada: 2026-03-10 (Absorbida por infra-multi-env-deployment-impl-2026-03-10)
- **Sales UX - Edición de Documentos, Creación de Pedido desde Presupuesto y Normalización de Estados** | Iniciada: 2026-03-06 | Finalizada: 2026-03-07

- **ERP Core - Cobertura Product Application + Fix Nombre Producto en Líneas de Venta** | Iniciada: 2026-03-05 | Finalizada: 2026-03-06
- **Revisión y Refinamiento de la Documentación General** | Iniciada: 2026-02-28 | Finalizada: 2026-03-01
- **MES - Revisión de nomenclatura y modelo Trabajo Definido vs Trabajo Real** | Iniciada: 2026-02-23 | Finalizada: 2026-02-25
- **Sales Module - Revisión UX y Validación Funcional** | Iniciada: 2026-03-02 | Finalizada: 2026-03-06
- **Product Module - Comprobaciones y Validación Continua** | Iniciada: 2026-03-01 | Finalizada: 2026-03-02
- **Product Module - Validación de Funcionalidad y Corrección de Bugs** | Iniciada: 2026-02-28 | Finalizada: 2026-03-01
- **Party Module - Consolidación de Migraciones y Smart Contact Deletion** | Iniciada: 2026-02-25 | Finalizada: 2026-02-28
- **Seguimiento Sprint 13 - Validación final Sales/Tax** | Iniciada: 2026-02-23 | Finalizada: 2026-02-24
- **Stabilización Party/IAM + continuidad MES** | Iniciada: 2026-02-23 | Finalizada: 2026-02-23
- **Sprint 13 / Implementación Sistema Impuestos + UX Improvements + Verificación Final** | Iniciada: 2026-02-22 | Finalizada: 2026-02-22
- **Sprint 13 / Tarea 01 - MVP Backend Coverage Compliance** | Iniciada: 2026-02-21 | Finalizada: 2026-02-22
- **Sprint 11 / Tarea 03 - ERP Core UX Testing & Validation** | Iniciada: 2026-02-18 | Finalizada: 2026-02-22
- **Sprint 12 / Tarea 01 - MES Module Foundation & Architecture** | Iniciada: 2026-02-18 | Finalizada: 2026-02-21
- **Sprint 11 / Critical Remediation + Error Cleanup + ProductGroup Refactor** | Iniciada: 2026-02-18 | Finalizada: 2026-02-18
- **UI Icons Review & Standardization** | Iniciada: 2026-02-15 | Finalizada: 2026-02-18
- **Sprint 11 FASE 7 / Metrics & Reporting** | Iniciada: 2026-02-16 | Finalizada: 2026-02-17
- **Sprint 11 / ERP Core Validation & Quality Assurance** | Iniciada: 2026-02-15 | Finalizada: 2026-02-16
- **Refactor bootstrap.yaml into Modular Agents** | Iniciada: 2026-02-15 | Finalizada: 2026-02-16
- **Scaffolding Improvements - bootstrap.yaml and load-session.yaml** | Iniciada: 2026-02-15 | Finalizada: 2026-02-15
- **Sprint 10 / Sales Module Complete - ERP CORE 100%** | Iniciada: 2026-02-15 | Finalizada: 2026-02-15
- **Sprint 10 / Sales UX Enhancement + Quotes & Delivery Notes** | Iniciada: 2026-02-15 | Finalizada: 2026-02-15
- **Sprint 10 / Sales Frontend Complete + MES Backend Base** | Iniciada: 2026-02-14 | Finalizada: 2026-02-14
- **Sprint 09 / Pricing Integration Panel** | Iniciada: 2026-02-14 | Finalizada: 2026-02-14
- **Sprint 09 / Master Data CRUD Complete + Refactor Atributos** | Iniciada: 2026-02-14 | Finalizada: 2026-02-14
- **Sprint 09 / Implementación UPDATE Product Endpoint** | Iniciada: 2026-02-14 | Finalizada: 2026-02-14
- **Sprint 09 / BUG FIX: Creación de Productos con Atributos Directos** | Iniciada: 2026-02-14 | Finalizada: 2026-02-14
- **Sprint 09 / Tarea 05 - Documentación y UI de Productos + Sistema de Variantes** | Iniciada: 2026-02-13 | Finalizada: 2026-02-13
- **Sprint 11 / Task 02 - Critical Remediation Plan COMPLETO** | Iniciada: 2026-02-17 | Finalizada: 2026-02-18
