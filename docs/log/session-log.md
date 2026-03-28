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

## Mejora UI/UX — Unificación Estética y Componentes Base

- **Session ID:** `ui-ux-improvement-post-mvp-21-03-2026`
- **Status:** En Progreso — **Post-MVP: Estandarización estética pre-TFM**
- **Prioridad:** Alta (Consistencia Visual TFM)
- **Started:** 21-03-2026
- **Rama:** `feat/ui-ux-standardization-post-mvp` (desde `develop`)

### Contexto

Auditoría de UI/UX completada en `tmp/ui-ux-improvement-suggestions.md`. Esta sesión implementa el plan de mejora estética para dar coherencia visual antes de la presentación TFM. Se documenta en el **Sprint 16** (Tarea 16-01).

**Dependencias:** Requiere completar previamente la sesión de despliegue multientorno.

### Próximos Pasos

- [x] Crear `apps/frontend/src/design-system/_buttons.css` con estilos globales e importar en `theme.css`.
- [x] Estandarizar iconografía en módulos clave (Sales, Party).
- [x] Unificar comportamiento de listados: fila clickeable + botón de acción iconográfico.
- [x] Refactorizar componente `PageHeader` y aplicar Material Symbols.
- [x] Crear componente maestro `BaseEntityPage`.
- [x] Implementar Sticky Identity Header y Context Header.
- [x] Estandarizar secciones: Toolbar, Summary, Related (Centralizado en `_sections.css`).
- [x] Migrar `QuoteDetail.vue` e `InvoiceDetail.vue` al nuevo estándar `BaseEntityPage`.
- [ ] **ESTABILIZAR OrderDetail.vue (BLOQUEO CRÍTICO)**: 
    - [x] Corregir error "invalid id" (redirección post-creación).
    - [x] Corregir error "invalid request body" (omisión de UUIDs vacíos).
    - [x] Corregir rutas de API (cambio de `/items` a `/line-items`).
    - [x] Eliminar bucles de actualización recursiva (`Maximum recursive updates exceeded`) persistentes.
    - [x] Asegurar carga de precios y mapeo de `productVariantId` tras limpieza de reactividad.    - [x] Confirmar/Anular pedido: campo `status` → `newStatus` en `changeOrderStatus`.
    - [x] MES: `mesWorkId`/`observations` → `workSetupId`/`description` en template y `enterEditMode`.
    - [x] Reabrir pedido cancelado: función `reopenOrder` + computed `canReopen` + botón `btn-primary`.
    - [x] Listado de albaranes: query param `order_id` → `salesOrderId` + array parsing directo.
    - [x] Listado de facturas: query param `order_id` → `orderId` + array parsing directo.
    - [x] Navegación MES: botón apunta a orden de trabajo (`workOrderId`); si no existe → badge "Pendiente".
    - [x] `enterEditMode` + `saveOrder`: preservar `workOrderId` en mesWorkRefs al editar.
- **SIGUIENTE: Arreglar UI Presupuestos (QuoteDetail.vue)**:
    - [ ] Aplicar mismos fixes de integración frontend↔backend (ver `/memories/repo/orderdetail-frontend-fixes.md`).
    - [ ] Listados relacionados (albaranes, facturas vinculadas) — mismo patrón de array parsing.
    - [ ] Revisar campos MES en template de presupuesto (`workSetupId`, `description`).
    - [ ] Verificar navegación y acciones de estado (aprobar, rechazar, convertir a pedido).
### Archivos de Contexto

- `apps/frontend/src/design-system/_sections.css`
- `apps/frontend/src/pages/sales/OrderDetail.vue`
- `apps/frontend/src/router/guards.ts`

---

## Preparación TFM — Presentación Final de TramaTex

- **Session ID:** `tfm-final-presentation-21-03-2026`
- **Status:** ⏸️ En Pausa
- **Prioridad:** 2º (bloquea entrega final)
- **Sprint:** N/A
- **Started:** 21-03-2026
- **Rama:** `docs/tfm-final-preparation` (pushed a origin, pendiente merge a `develop`)

### Contexto

TramaTex se presenta como Trabajo Fin de Máster (TFM). Esta sesión cubre la preparación integral del proyecto para su entrega y defensa académica: revisión de documentación, presentación, memoria, y asegurar que el estado del código, los tests y el despliegue son coherentes y presentables. Es la **última sesión** del proyecto.

### Progreso completado (sesiones 21–25 mar 2026)

- [x] Licencias verificadas — "Jorge Cortés Villalba" aparece en `LICENSE.md`, `project-scaffolding/LICENSE.md`, `README.md` y ADRs.
- [x] Tests Go backend: 29 paquetes OK, 0 FAIL (excluido `product/persistence` que requiere PostgreSQL local).
- [x] Tests frontend Vitest: 230 tests pasados, 0 fallos.
- [x] Bug corregido en `sales_service_test.go`: mock `ListBySalesOrderID` con `.Once()` para flujo de doble llamada.
- [x] `CONTRIBUTING.md`: corregido enlace roto a ADR-011.
- [x] Presentación: reemplazado placeholder de logo por texto.
- [x] `README.md`: añadida sección Demo Pública, corregido comando de instalación, añadida sección Mantenimiento.
- [x] Creado `.github/workflows/demo-reset.yml`: reset semanal automatizado de la demo.
- [x] Auditoría de secretos: **LIMPIO** — no se exponen secretos reales. Solo credenciales de demo (`admin123`) y placeholders en `.example`.
- [x] ADR filename casing: evaluado, descartado (funciona en Windows/GitHub, riesgo alto para beneficio mínimo).
- [x] `tmp/` ya está en `.gitignore`, no se trackea.

### Próximos Pasos (para la próxima jornada)

**PRIMERO — Limpiar artefactos trackeados en frontend:**
- [ ] Añadir a `apps/frontend/.gitignore`:
  ```
  # Test artifacts
  test-results/
  playwright-report/
  test-results.txt
  build-output.txt
  ```
- [ ] Eliminar del tracking (sin borrar local): `git rm --cached apps/frontend/build-output.txt apps/frontend/test-results.txt` y `git rm --cached -r apps/frontend/test-results/ apps/frontend/playwright-report/`
- [ ] Commit: `chore: remove tracked test artifacts and update frontend .gitignore`

**DESPUÉS — Merge y tareas TFM restantes:**
- [ ] Merge rama `docs/tfm-final-preparation` → `develop` (y luego a `main` como entrega final).
- [ ] Preparar memoria/informe TFM (estructura, introducción, conclusiones, trabajo futuro).
- [ ] Verificar despliegue Docker de principio a fin.
- [ ] Revisar presentación final (`docs/presentations/tramatex-presentation.md`, `TramaTex_Presentacion_Final.pptx`).
- [ ] Decidir sobre sesión UI/UX (`ui-ux-improvement-post-mvp-21-03-2026`) — opcional pre-TFM.

### Archivos de Contexto

- `docs/presentations/tramatex-presentation.md`
- `docs/presentations/TramaTex_Presentacion_Final.pptx`
- `docs/presentations/slides_spec.md`
- `docs/architecture/architecture-vision.md`
- `README.md`
- `CONTRIBUTING.md`
- `apps/frontend/.gitignore` ← pendiente de actualizar

---

---
# REGISTRO DE SESIONES CERRADAS
---
- **Alineación Documental Post-Refactors Sprint 14** | Iniciada: 21-03-2026 | Finalizada: 25-03-2026 | Status: ✅ COMPLETADO | Rama: `doc/alignment-sprint14-cicd-verify` → `develop`, `staging`, `master`
- **QA — Verificación de Calidad Integral** | Iniciada: 21-03-2026 | Finalizada: 22-03-2026 | Status: ✅ COMPLETADO | Rama: `qa/full-verification` → `develop` | 5 commits, 4 bugs corregidos, QA manual 6/6 puntos OK
- **Refactor sort_order → DirectAttributeIDs (Producto/Atributos)** | Iniciada: 21-03-2026 | Finalizada: 21-03-2026 | Status: ✅ COMPLETADO
- **Análisis de Refinamiento Arquitectónico del MVP (Sprint 14)** | Iniciada: 12-03-2026 | Finalizada: 21-03-2026 | Status: ✅ COMPLETADO | Ver: [sprint-14](./sprints/sprint-14/sprint-14-summary.md) | PR pendiente: `mvp-arch-refinement` → `develop`
- **Análisis de Refinamiento Arquitectónico del Módulo MES** | Iniciada: 20-03-2026 | Finalizada: 20-03-2026 | Status: ✅ COMPLETADO
- **Integración MES-Sales: Terminal de Taller y Visibilidad de Producción en Pedidos** | Iniciada: 19-03-2026 | Finalizada: 19-03-2026
- **Refinamiento y Estabilización ERP Core** | Iniciada: 09-03-2026 | Finalizada: 14-03-2026
