# Tarea 16-01: Estandarización de UI/UX y Componentes Base (Post-MVP)

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 01 |
| **ID de Sprint** | sprint-16 |
| **Título** | Estandarización de UI/UX y Componentes Base |
| **Estado** | ⏳ En Progreso |
| **Facilitador/LLM** | Gemini CLI |
| **Fecha de Inicio** | 2026-03-26 |
| **Fecha de Fin** | |
| **Duración Estimada** | 6-8 horas |
| **Duración Real** | |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [x] **Unificación de Estilos Globales (CSS):**
   - Crear `apps/frontend/src/design-system/_buttons.css` con clases `.btn-primary`, `.btn-secondary`, `.btn-outline`, `.btn-danger`.
   - Integrar en `theme.css` y asegurar el uso de variables de `_variables.css`.
2. [x] **Eliminación de Emojis e Iconografía Lucide:**
   - Sustituir todos los emojis (🗑️, 🖨️, 💰, ⚠️, ⚙️) por componentes `lucide-vue-next`.
3. [ ] **Estandarización de Listados (Tables):**
   - Implementar el patrón "Fila clickeable + Botón de acción iconográfico" en `PartyList.vue`.
4. [ ] **Creación del Componente BasePageHeader:**
   - Unificar Breadcrumbs, Títulos (H1) y Acciones en un componente reutilizable.
5. [ ] **Refactor de Formularios (Card Design):**
   - Migrar `PartyForm.vue` de `fieldset/legend` a un diseño basado en tarjetas y mejorar contraste de etiquetas.

---

## 📊 CONTEXTO DE ENTRADA

### Estado Anterior

**Última tarea completada:** 15-01 (Refinamiento Arquitectónico MVP)

**Cambios desde última tarea:**
- Auditoría de UI/UX completada en `tmp/ui-ux-improvement-suggestions.md`.
- Proyecto preparado para la fase TFM, con necesidad de coherencia visual.

**Estado en project-status.md:**
- Fase actual: Post-MVP / Preparación TFM.

---

## 🛠️ PLAN DE TRABAJO

### Fase 1: Base de Diseño (1h)
- [x] Crear `apps/frontend/src/design-system/_buttons.css`.
- [x] Actualizar `apps/frontend/src/assets/styles/theme.css` para incluir los nuevos estilos.
- [x] Validar variables de color y espaciado en `_variables.css`.

### Fase 2: Iconografía (1h)
- [x] Identificar todos los usos de emojis en `apps/frontend/src/`.
- [x] Reemplazar sistemáticamente por iconos de Lucide.

### Fase 3: Componentes y Listados (3h)
- [ ] Crear `BasePageHeader.vue` en `apps/frontend/src/components/layout/`.
- [ ] Refactorizar `PartyList.vue` para usar el nuevo estándar de tablas y cabeceras.
- [ ] Refactorizar `PartyForm.vue` para usar el diseño de tarjetas y mejorar inputs.

### Fase 4: Validación y QA (1h)
- [ ] `npm run lint` en el frontend.
- [ ] Verificación visual de los cambios en el navegador.
- [ ] Asegurar que no hay regresiones en la funcionalidad.

---

## ✅ DEFINICIÓN DE "HECHO"

La tarea se considera completada cuando:
- [ ] Los botones usan estilos globales y no scoped.
- [ ] No quedan emojis en la interfaz de usuario.
- [ ] Los listados son consistentes entre mallas (hover, click, acciones).
- [ ] `PartyList.vue` y `PartyForm.vue` sirven como modelos de referencia impecables.
- [ ] El diseño es responsivo y coherente en espaciados.

---

## ✍️ FIRMA

**Facilitador:** Gemini CLI  
**LLM:** Gemini 2.0 Flash
