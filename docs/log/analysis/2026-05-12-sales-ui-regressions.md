# Análisis de Regresiones: Módulo Sales y UI/UX
**Fecha:** 2026-05-12
**Estado:** SESIÓN ABIERTA

## 🚩 Problemas Detectados
1. **Pantalla en Blanco en OrderDetail.vue:** Al intentar ver los detalles de un pedido, la aplicación no renderiza el contenido. Posible regresión en los tipos de datos o componentes hijos (`OrderLines.vue`).
2. **Fallo en Conversión de Presupuesto:** A pesar de habilitar el botón en estado `ISSUED`, la acción de conversión no parece estar completando el flujo esperado en staging.
3. **Persistencia de UI/UX:** Verificado que el Menú de Ayuda ya usa `Teleport` y `z-index: 200000`, pero se requiere validación final tras corregir las regresiones de Sales.

## 🛠️ Trabajo Realizado
- Rama creada: `fix/sales-order-detail-blank-screen-and-conversion-flow`
- Investigada lógica de estados en `quote_service.go`: El backend soporta `DRAFT -> ISSUED -> APPROVED -> CONVERTED`.
- Identificada posible causa en `OrderDetail.vue`: Desajuste entre nombres de campos (camelCase vs snake_case) tras las últimas refactorizaciones de UI.

## 📋 Pasos Pendientes
- [ ] Depurar `OrderDetail.vue` en local para encontrar el error de renderizado (posiblemente un campo nulo en `editableOrder`).
- [ ] Validar el endpoint `/api/sales/quotes/:id/convert` con un presupuesto en estado `ISSUED`.
- [ ] Unificar definitivamente la gestión de líneas entre Pedidos y Presupuestos para evitar regresiones de visualización.

---
**Nota para la próxima sesión:** Priorizar la resolución de la pantalla en blanco ya que bloquea la validación de la conversión (no se puede ver el pedido generado).
