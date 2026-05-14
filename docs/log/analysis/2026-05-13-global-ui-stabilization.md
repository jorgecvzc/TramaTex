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
- [ ] **Test E2E Conversión y Factura**: Crear test que verifique la presencia del botón "Convertir a Pedido" y la carga de `InvoiceDetail`.
- [ ] **Despliegue a Staging**: Actualizar `pcele`.
- [ ] **Validación en pcele**: Confirmar visualmente que el botón aparece para presupuestos `ISSUED` y las facturas cargan tras el refactor `camelCase`.

### Fase 2: Estandarización de UI del Módulo de Entidades (Parties)
- [ ] **Crear Test E2E para Botones**: Escribir un test de Playwright que falle al no encontrar los botones "Guardar" / "Cancelar" en la cabecera.
- [ ] **Refactorizar `PartyCreate.vue` / `PartyDetail.vue`**: Mover botones a la cabecera (`<template #actions>`).
- [ ] **Despliegue y Validación**: Ejecutar test en `pcele` tras el despliegue.

### Fase 3: Refactorización y Unificación Final del Módulo de Ventas (Sales)
- [ ] **Limpieza de `snake_case`**: Eliminar fallbacks en `OrderDetail` y `QuoteDetail`.
- [ ] **Unificación de Iconografía**: Migrar a `Lucide` en `DeliveryNoteDetail.vue`.
- [ ] **Componente `OrderLines` en Albaranes**: Refactorizar `DeliveryNoteDetail.vue`.
- [ ] **Estandarización en Creación**: Refactorizar `QuoteCreate.vue` y `OrderCreate.vue`.
- [ ] **Despliegue y Validación E2E**: Asegurar que todos los flujos de creación/detalle funcionan por teclado.
