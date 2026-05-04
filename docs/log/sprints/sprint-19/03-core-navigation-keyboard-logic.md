# TAREA: 03 — Navegación Core (Lógica de teclado en listados)

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 03 |
| **ID de Sprint** | sprint-19 |
| **Título** | Navegación Core (Lógica de teclado en listados) |
| **Estado** | ✅ Completada (AI) |
| **Facilitador/LLM** | Gemini CLI |
| **Fecha de Inicio** | 2026-05-04 |
| **Fecha de Fin** | 2026-05-04 |
| **Duración Estimada** | 2 horas |
| **Duración Real** | 1 hora |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [x] **Navegación por Teclado**: Implementar soporte para flechas `Up`/`Down` en componentes de listado (`BaseCatalog`).
2. [x] **Selección con Enter**: Permitir la navegación al detalle de la entidad al presionar `Enter`.
3. [x] **Auto-Scroll**: Implementar `scrollToSelected` para que la fila resaltada siempre sea visible en listados largos.
4. [x] **Feedback Visual**: Refinar el estilo de la fila seleccionada para que sea indistinguible del foco industrial.

---

## 📊 CONTEXTO DE ENTRADA

### Estado Anterior

**Última tarea completada:** 19-02 — Toasts y Feedback Crítico

**Cambios desde última tarea:**
- Sistema de notificaciones global (Toasts) implementado.
- Estilos base industriales establecidos.

**Estado en project-status.md:**
- Fase actual: Post-MVP Fase 1 (Cimientos y Ergonomía).

### Bloqueadores/Dependencias

- Ninguno identificado.

---

### Prioridades para esta Tarea

**Crítica (Must Have):**
- Navegación Up/Down en listados.
- Enter para seleccionar.

**Alta (Should Have):**
- Scroll automático a la fila seleccionada.

**Media (Nice to Have):**
- Indicadores de atajos de teclado visuales (tooltips o hints).

---

## 🛠️ PLAN DE TRABAJO

### Fase 1: Análisis Arquitectónico (20 min) ✅

- [x] Identificar el componente base de listados (`BaseCatalog.vue` o similar).
- [x] Revisar la gestión actual del foco y eventos de teclado en Vue 3.

### Fase 2: Implementación ✅

**Frontend (1.5 horas):**
- [x] Implementar `selectedIndex` en el estado del componente de listado.
- [x] Añadir listener de eventos `keydown` (global o local al componente).
- [x] Implementar lógica de `scrollToSelected` usando `Element.scrollIntoView()`.
- [x] Vincular `Enter` con la acción de navegación/edición existente.

### Fase 3: Validación ✅

- [x] Verificar navegación por teclado en varios módulos (Ventas, Productos, etc.).
- [x] `npm run test` para asegurar que no hay regresiones.
- [x] Validar accesibilidad básica.

---

## ✅ DEFINICIÓN DE "HECHO"

La tarea se considera completada cuando:

- [ ] Se puede navegar por un listado completo usando solo el teclado.
- [ ] `Enter` abre el detalle de la fila seleccionada.
- [ ] La fila seleccionada siempre es visible durante la navegación.
- [ ] Tarea documentada en este archivo.

---

## 🚀 PRÓXIMOS PASOS

1. [ ] Tarea 04: Skeletons Industriales para estados de carga.
