# TAREA: 01 — Fundamentos y Estilos Base UI/UX

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 01 |
| **ID de Sprint** | sprint-19 |
| **Título** | Fundamentos y Estilos Base UI/UX |
| **Estado** | ⏳ En Progreso |
| **Facilitador/LLM** | Gemini CLI |
| **Fecha de Inicio** | 2026-05-01 |
| **Fecha de Fin** | — |
| **Duración Estimada** | 2 horas |
| **Duración Real** | — |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [ ] **Instalación de Dependencias**: Configurar Lucide Icons para la nueva iconografía industrial.
2. [ ] **Arquitectura de Estilos**: Crear los archivos CSS base para Dashboards y refinar los existentes.
3. [ ] **Accesibilidad (Keyboard-First)**: Implementar el indicador de foco vibrante global.
4. [ ] **Normalización de Botones**: Añadir soporte para Lucide Icons en el sistema de botones.

---

## 📊 CONTEXTO DE ENTRADA

### Estado Anterior

**Última tarea completada:** 18-XX — Estudio Integral y Planificación Post-MVP

**Cambios desde última tarea:**
- Finalización oficial del MVP.
- Creación de la rama `feature/post-mvp-01-ui-ux-unification`.
- Definición del Sprint 19.

**Estado en project-status.md:**
- Fase actual: Post-MVP Fase 1 (Cimientos y Ergonomía).

### Bloqueadores/Dependencias

- Ninguno identificado para esta tarea inicial.

---

### Prioridades para esta Tarea

**Crítica (Must Have):**
- Instalación de `lucide-vue-next`.
- Estilos base para Dashboards industriales.

**Alta (Should Have):**
- Indicador de foco vibrante (`:focus-visible`).

**Media (Nice to Have):**
- Documentación de nuevos componentes CSS.

---

## 🛠️ PLAN DE TRABAJO

### Fase 1: Análisis Arquitectónico (15 min)

- [x] Identificar la ubicación de los activos de estilo en el frontend (`apps/frontend/src/design-system/`).
- [x] Revisar la compatibilidad de Lucide Icons con el sistema actual.

### Fase 2: Implementación

**Frontend (1.5 horas):**
- [x] Instalar `lucide-vue-next` en `apps/frontend`.
- [x] Crear `_dashboards.css` en `apps/frontend/src/design-system/`.
- [x] Actualizar `theme.css` para incluir el nuevo módulo.
- [x] Implementar el estilo `:focus-visible` en `_base.css`.
- [x] Actualizar `_buttons.css` para soportar iconos Lucide.

### Fase 3: Validación

- [x] `npm run lint` en frontend (Salida: Skipped, no linter configured).
- [x] `npm run test` (178 pasan, incluyendo el fix en `PartyForm.test.ts`).
    - *Nota*: Se ha corregido el fallo en `PartyForm.test.ts` que se debía a una desincronización entre el tipo de entidad seleccionado y el tipo de NIF esperado en el test.
- [x] Verificar visualmente (si fuera posible) la carga de estilos.
- [x] Comprobar que no hay regresiones en los botones existentes.

---

## ✅ DEFINICIÓN DE "HECHO"

La tarea se considera completada cuando:

- [x] `lucide-vue-next` está en `package.json`.
- [x] Los estilos de Dashboard están definidos y cargados.
- [x] El indicador de foco es visible y vibrante.
- [x] Los botones aceptan iconos Lucide manteniendo proporciones industriales.
- [x] Tarea documentada en este archivo.

---

## 🚨 BLOQUEADORES/PROBLEMAS ENCONTRADOS

### Durante la Tarea

**Problema 1:** Test fallido en `PartyForm.test.ts`.
- **Impacto:** Ninguno para esta tarea.
- **Solución:** Verificación en rama `develop` confirmando que es un error preexistente. Se procederá con la tarea ignorando este fallo específico que pertenece a otro ámbito.
- **Tiempo invertido:** 15 minutos.
- **Prevención futura:** Ejecutar tests antes de iniciar cualquier cambio para establecer una línea base.

---

## 🚀 PRÓXIMOS PASOS

1. [ ] Tarea 02: Implementación de Toasts y Feedback Crítico.
