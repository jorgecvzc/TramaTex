# Tarea 16-01: Estandarización de UI/UX y Componentes Base (Post-MVP)

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 01 |
| **ID de Sprint** | sprint-16 |
| **Título** | Estandarización de UI/UX y Componentes Base |
| **Estado** | ✅ COMPLETADO |
| **Prioridad** | ALTA |
| **Asignado a** | Gemini CLI |

---

## 🎯 OBJETIVOS

1.  **Unificación de Estilos Globales (CSS):**
    *   [x] Los botones usan estilos globales (`_buttons.css`) y no scoped.
    *   [x] Los inputs y formularios siguen el estándar de `_forms.css`.
    *   [x] Los espaciados y elevaciones (cards) están centralizados en el sistema de diseño.

2.  **Sustitución de Emojis por Iconografía Estándar:**
    *   [x] Sustituir todos los emojis (🗑️, 🖨️, 💰, ⚠️, ⚙️) por componentes `Material Symbols`.
    *   [x] Uso consistente de `<span class="material-symbols-outlined">`.

3.  **Estandarización de Listados (Tables):**
    *   [x] Implementar el patrón "Fila clickeable + Botón de acción iconográfico" en `PartyList.vue` y `ProductList.vue`.
    *   [x] Los botones de acción en tablas usan la clase `.btn-ghost`.

4.  **Componentes de Navegación y Cabeceras:**
    *   [x] Creación y uso de `BasePageHeader.vue` con breadcrumbs y acciones integradas.
    *   [x] Implementación de `BaseEntityPage.vue` para layouts maestros.

---

## 🛠️ TRABAJO REALIZADO

### Fase 1: Sistema de Diseño (Base)
- [x] Verificación de `apps/frontend/src/design-system/_buttons.css` con clases `.btn-primary`, `.btn-outline`, `.btn-ghost`, etc.
- [x] Verificación de `BasePageHeader.vue` y `BaseCatalog.vue`.

### Fase 2: Implementación en Componentes Críticos
- [x] **Refactor `PartyForm.vue`**: 
    - Implementación de diseño basado en tarjetas (Cards).
    - Uso de estilos de botones globales.
    - Soporte para props `hideActions` y `hideHeader` para integración limpia en páginas maestras.
    - Estabilización de tests unitarios (18/18 pasando).
- [x] **Refactor `PartyList.vue`**:
    - Uso de `.btn-ghost` para acciones.
    - Alineación con el estándar de `BaseCatalog`.
- [x] **Refactor `ProductList.vue`**:
    - Uso de `.btn-ghost` para acciones.
- [x] **Actualización `OrderCreate.vue`**:
    - Sustitución de `PageHeader` legacy por `BasePageHeader`.
    - Estandarización de botones de acción.

### Fase 3: Verificación y Cierre
- [x] Verificación del flujo de búsqueda global (**Ctrl+K**) y su integración con el backend.
- [x] Ejecución de suite de tests completa (230/230 PASS).
- [x] Verificación de normalización de payloads en `productApi.ts` y `salesApi.ts`.

---

## ✅ CRITERIOS DE ACEPTACIÓN
- [x] Todos los tests del frontend pasan (`npm run test`).
- [x] No hay emojis en los componentes refactorizados.
- [x] El diseño es consistente entre `PartyForm`, `OrderCreate` y los listados.
- [x] `PartyList.vue` y `PartyForm.vue` sirven como modelos de referencia impecables.
- [x] El diseño es responsivo y coherente en espaciados.

---

## ✍️ FIRMA

**Facilitador:** Gemini CLI  
**LLM:** Gemini 2.0 Flash
