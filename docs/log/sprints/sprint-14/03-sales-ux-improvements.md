# Sprint 14 / Tarea 03 — Mejoras UX del Módulo Sales

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 14-03 |
| **ID de Sprint** | sprint-14 |
| **Título** | Sales UX: Layout, VariantSelector y Ciclo de Facturación Simplificado |
| **Estado** | ✅ Completado |
| **Facilitador/LLM** | GitHub Copilot / Claude Sonnet |
| **Fecha de Inicio** | 2026-03-05 |
| **Fecha de Fin** | 2026-03-12 |
| **Rama** | `develop` (mergeado desde `sales-ux-validation`) |

---

## 🎯 Objetivos

1. [x] Rediseñar layout del módulo Sales a ancho completo (full-width)
2. [x] Implementar filas únicas en líneas de pedido (single-row line items)
3. [x] Extender `VariantSelector` a todos los formularios de Sales
4. [x] Mostrar totales en `OrderCreate`
5. [x] Completar ciclo de facturas simplificadas

---

## 📊 Trabajo Realizado

### Layout y UX
- **Full-width layout**: el módulo Sales ahora ocupa el ancho completo de la pantalla
- **Single-row line items**: las líneas de pedido se muestran en una sola fila para mayor densidad de información
- **Totales en OrderCreate**: visualización de subtotal, impuestos y total en tiempo real durante la creación de pedidos

### VariantSelector Universal
- El componente `VariantSelector` se extiende a todos los formularios donde se seleccionan productos: Presupuestos, Pedidos, Facturas y Albaranes

### Facturas Simplificadas
- Implementación completa del ciclo de facturas simplificadas (tickets)
- Alineación con ADR-020: formato de factura simplificada para operaciones de bajo importe
- Mejoras en formulario de precios

---

## 🔗 Commits Clave

| Hash | Descripción |
|------|-------------|
| `bceee3e` | `chore: open sales-ux-review-validation session` |
| `3d17842` | `feat(sales-ux): full-width layout, single-row line items, VariantSelector everywhere, OrderCreate totals` |
| `a06e42b` | `docs: update session-log with Sales UX improvements results` |
| `89bcf78` | `feat(sales): complete Sales UX validation - simplified invoices, pricing, docs` |
| `29c424e` | `merge: sales-ux-validation into develop - complete Sales UX validation` |
