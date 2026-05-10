# TAREA: 02 — Toasts y Feedback Crítico

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 02 |
| **ID de Sprint** | sprint-19 |
| **Título** | Toasts y Feedback Crítico |
| **Estado** | ⏳ En Progreso |
| **Facilitador/LLM** | Gemini CLI |
| **Fecha de Inicio** | 2026-05-01 |
| **Fecha de Fin** | — |
| **Duración Estimada** | 3 horas |
| **Duración Real** | — |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [x] **Store de Notificaciones**: Crear un store de Pinia (`useToastStore`) para gestionar mensajes globales.
2. [x] **Componente Base Toast**: Diseñar un componente `BaseToast.vue` que soporte tipos (success, error, warning, info) e iconos Lucide.
3. [x] **Contenedor Global**: Integrar un contenedor de Toasts en `App.vue` para mostrar las notificaciones.
4. [x] **Validación Inline**: Refinar los estilos de validación inline para que se activen de forma vibrante al perder el foco (`on-blur`).
5. [x] **Sustitución de `alert`**: Buscar y reemplazar usos de `alert()` nativos por el nuevo sistema de Toasts en los módulos principales.

---

## 📊 CONTEXTO DE ENTRADA

### Estado Anterior

**Última tarea completada:** 19-01 — Fundamentos y Estilos Base UI/UX

**Cambios desde última tarea:**
- Instalación de Lucide Icons.
- Fix del test preexistente en `PartyForm.test.ts`.

**Estado en project-status.md:**
- Fase actual: Post-MVP Fase 1 (Cimientos y Ergonomía).

---

## 🛠️ PLAN DE TRABAJO

### Fase 1: Análisis y Diseño (20 min) ✅

- [x] Definir la estructura del objeto Toast (id, message, type, duration).
- [x] Mapear los iconos Lucide a cada tipo de mensaje.

### Fase 2: Implementación ✅

**Frontend (2 horas):**
- [x] Crear `apps/frontend/src/stores/toast.ts`.
- [x] Crear `apps/frontend/src/components/shared/BaseToast.vue`.
- [x] Crear `apps/frontend/src/components/shared/ToastContainer.vue`.
- [x] Registrar el contenedor en `App.vue`.
- [x] Refactorizar componentes de Ventas, Entidades y Administración para usar el nuevo store.

### Fase 3: Validación ✅

- [x] Crear test unitario para `toastStore`.
- [x] Verificar visualmente el comportamiento de los Toasts (auto-dismiss, hover pause, etc.).
- [x] `npm run test` para asegurar que no hay regresiones.

---

## ✅ DEFINICIÓN DE "HECHO"

La tarea se considera completada cuando:

- [x] Los mensajes aparecen y desaparecen correctamente.
- [x] El código está libre de `alert()` en los módulos principales.
- [x] La validación inline en formularios es consistente con el nuevo estilo (vibrant).
- [x] Tarea documentada en este archivo.

---

## 🚀 PRÓXIMOS PASOS

1. [ ] Tarea 03: Navegación Core (Lógica de teclado en listados).
