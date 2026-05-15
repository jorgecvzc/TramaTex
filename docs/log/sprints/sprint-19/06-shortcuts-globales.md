# TAREA: 06 — Atajos de Teclado Globales

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 06 |
| **ID de Sprint** | sprint-19 |
| **Título** | Atajos de Teclado Globales |
| **Status** | Finalizada ✅ |
| **Prioridad** | Alta |
| **Asignado a** | Gemini AI |

---

## 🎯 OBJETIVO

Implementar una capa de atajos de teclado globales para mejorar la eficiencia operativa en todo el ERP, permitiendo el acceso rápido a la búsqueda, creación de elementos, guardado de formularios y actualización de datos sin usar el ratón.

---

## 🛠️ IMPLEMENTACIÓN

### 1. Atajos Definidos
- **Ctrl + K**: Abre el buscador global desde cualquier parte del sistema.
- **Alt + N**: Crea un nuevo elemento contextualmente (ej. en `/parties` va a `/parties/new`).
- **Ctrl + Enter**: Ejecuta la acción de guardado o envío del formulario activo.
- **Alt + R**: Refresca los datos del componente actual (Soft Refresh) sin recargar la página completa.

### 2. Cambios Realizados
- **`Navbar.vue`**: Centralización de la lógica de escucha de eventos `keydown`.
- **`BaseCatalog.vue`**: Soporte para el evento global `tramatex-refresh`.
- **`Detail.vue` (Parties)**: Soporte para el evento global `tramatex-refresh`.
- **`shortcuts.test.ts`**: Suite de pruebas unitarias para validar la activación de los atajos.

---

## ✅ VERIFICACIÓN

### Tests Unitarios
- Se han ejecutado las pruebas en `src/__tests__/unit/shortcuts.test.ts` con resultado exitoso (3/3 passed).

### Escenarios Operativos
1. [x] Presionar Ctrl+K abre el diálogo de búsqueda.
2. [x] Presionar Alt+N en el listado de entidades redirige a la creación.
3. [x] Presionar Alt+R en un catálogo muestra el toast de refresco y actualiza la lista.
4. [x] Presionar Ctrl+Enter en un formulario intenta disparar el botón de submit.

---

## 🚀 PRÓXIMOS PASOS

1. [x] Despliegue en `pcele` para validación final por el usuario.
2. [ ] Fase 2: Unificación UI/UX (Módulos de Producción y Stock).
