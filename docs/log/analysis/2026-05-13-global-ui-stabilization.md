# Análisis de Tareas: Estabilización Global de UI y Flujos
**Fecha:** 2026-05-13
**Estado:** SESIÓN ABIERTA
**Rama:** `fix/sales-order-detail-blank-screen-and-conversion-flow`

## 🎯 Objetivo
Esta sesión de trabajo tiene como objetivo resolver las inconsistencias de UI/UX restantes, estabilizar los flujos de trabajo del módulo de Ventas y estandarizar la experiencia de usuario en toda la aplicación, siguiendo el patrón de cabecera de acciones unificada.

## 🚀 Protocolo de Finalización por Punto
Para cada bloque de tareas, se seguirá estrictamente:
1. **Implementación/Refactorización**: Aplicar los cambios de código.
2. **Test E2E**: Crear o actualizar tests de Playwright para validar la funcionalidad automáticamente.
3. **Despliegue**: Ejecutar `./scripts/rebuild-staging-remote.ps1` para actualizar `pcele`.
4. **Validación**: Verificar visualmente y mediante ejecución de tests en el entorno de Staging.

---

## ✅ Tareas Pendientes

### Fase 1: Verificación y Cierre de Issues Anteriores
- [x] **Test E2E Conversión y Factura**: Creado y validado en `sales-stabilization.spec.ts`.
- [x] **Despliegue a Staging**: Realizado con éxito.
- [x] **Valicación en pcele**: Confirmado el botón de conversión y carga de facturas camelCase.
- [x] **RE-VERIFICAR**: Reportado que los presupuestos EMITIDOS siguen sin mostrar el botón. Investigar lógica de estados en `QuoteDetail.vue`. **CORREGIDO**: Se añadió normalización robusta en `salesApi.ts` y soporte para variantes masculinas en `QuoteDetail.vue`.

### Fase 2: Estandarización de UI del Módulo de Entidades (Parties)
- [x] **Crear Test E2E para Botones**: Validado en `parties-stabilization.spec.ts`.
- [x] **Refactorizar `PartyCreate.vue` / `PartyDetail.vue`**: Acciones movidas a la cabecera.
- [x] **Despliegue y Validación**: Verificado en `pcele`.

### Fase 3: Refactorización y Unificación Final del Módulo de Ventas (Sales)
- [x] **Limpieza de `snake_case`**: Aplicada en `OrderDetail` y `QuoteDetail`.
- [x] **Unificación de Iconografía**: Migrado a Lucide en todo el módulo.
- [x] **Componente `OrderLines` en Albaranes**: Refactorización completada.
- [x] **Estandarización en Creación**: `QuoteCreate.vue` y `OrderCreate.vue` ahora usan `OrderLines`.
- [ ] **Despliegue y Validación E2E**: Pendiente desplegar Fase 3 y validar navegación por teclado en `pcele`.

---

## 🐛 Nuevos Bugs Detectados
- [x] **Facturas**: Al intentar marcar una factura como PAGADA, el backend devuelve el error `Invalid invoice status transition`. Revisar FSM de facturas en el dominio Go. **CORREGIDO**: Se habilitó la transición `DRAFT -> PAID` en el dominio y se actualizaron los tests.
