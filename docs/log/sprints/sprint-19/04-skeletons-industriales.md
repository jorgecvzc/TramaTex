# TAREA: 04 — Skeletons Industriales para estados de carga

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 04 |
| **ID de Sprint** | sprint-19 |
| **Título** | Skeletons Industriales para estados de carga |
| **Estado** | ✅ Completada |
| **Facilitador/LLM** | Gemini CLI |
| **Fecha de Inicio** | 2026-05-05 |
| **Fecha de Fin** | 2026-05-05 |
| **Duración Estimada** | 2 horas |
| **Duración Real** | 1.5 horas |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [x] **Arquitectura de Skeletons**: Definir los estilos base CSS para el efecto "shimmer" industrial.
2. [x] **Componente Base Skeleton**: Crear `BaseSkeleton.vue` altamente parametrizable.
3. [x] **Implementación en Listados**: Aplicar skeletons en el estado de carga de `BaseCatalog.vue`.
4. [x] **Implementación en Dashboards**: Aplicar skeletons en componentes como `PricingPanel` y `AttributesPanel`.
5. [x] **Consistencia Visual**: Asegurar que las formas de los skeletons coincidan con el contenido final (Cards, Tablas, Textos).

---

## 📊 CONTEXTO DE ENTRADA

### Estado Anterior

**Última tarea completada:** 19-03 — Navegación Core (Lógica de teclado en listados)

**Cambios desde última tarea:**
- Navegación por teclado y auto-scroll funcional en listados.
- Sistema de Toasts y estilos base consolidados.

**Estado en project-status.md:**
- Fase actual: Post-MVP Fase 1 (Cimientos y Ergonomía).

### Bloqueadores/Dependencias

- Ninguno.

---

## 🛠️ PLAN DE TRABAJO

### Fase 1: Análisis y Diseño (20 min)

- [x] Identificar componentes que necesitan skeletons (Listados de Ventas, Productos, Dashboards).
- [x] Definir variables CSS para el shimmer (colores neutros industriales).

### Fase 2: Implementación

**Frontend:**
- [x] Crear `_skeletons.css` en `apps/frontend/src/design-system/`.
- [x] Implementar `BaseSkeleton.vue`.
- [x] Integrar en `BaseCatalog.vue` y otros componentes críticos.

### Fase 3: Validación

- [x] Simular latencia de red para verificar el feedback visual.
- [x] Asegurar que no hay "layout shift" excesivo al cargar los datos reales.
- [x] `npm run test` (regresión).

---

## ✅ DEFINICIÓN DE "HECHO"

La tarea se considera completada cuando:

- [x] El componente `BaseSkeleton.skeleton-row` es reutilizable y documentado.
- [x] Los listados principales muestran skeletons durante la carga en lugar de estar vacíos o con spinners genéricos.
- [x] El efecto visual es fluido y alineado con la estética industrial (colores `--color-surface-soft`).
- [x] Tarea documentada en este archivo.

---

## 🚀 PRÓXIMOS PASOS

1. [x] Tarea 05: Consolidación de Dashboards y BasePageHeader. (Ver `docs/log/sprints/sprint-19/05-consolidacion-dashboards.md`)
2. [ ] Tarea 06: Implementación de Atajos de Teclado Globales.
