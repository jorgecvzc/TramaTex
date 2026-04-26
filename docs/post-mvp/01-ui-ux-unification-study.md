# Estudio Técnico: Unificación UI/UX y Accesibilidad Universal (Post-MVP)

Este documento define la estrategia técnica para transformar la interfaz de TramaTex en una herramienta de alta eficiencia para operarios, priorizando la accesibilidad por teclado y la coherencia visual.

## 1. Accesibilidad por Teclado (Prioridad 1)

El objetivo es permitir que un operario realice el flujo completo de trabajo (Ventas -> Producción) sin necesidad de usar el ratón.

### 1.1 Sistema de Atajos Globales
- **Búsqueda Global:** `Ctrl + K` para abrir el modal de búsqueda rápida.
- **Creación Rápida:** `Alt + N` para nuevas entidades/documentos (contextual).
- **Guardar/Confirmar:** `Ctrl + Enter` en formularios y modales.
- **Navegación de Listados:** Soporte para flechas `Up`/`Down` y `Enter` para selección en `BaseCatalog`.

### 1.2 Gestión de Foco e Indicadores Visuales
- **Focus Rings:** Implementar un estilo de foco visible y consistente (`:focus-visible`) que use el color `--color-primary`.
- **Tab Order:** Revisar y forzar el orden de tabulación lógico en todos los formularios, especialmente en el TPV y creación de pedidos.
- **Modales:** Implementar "Focus Trap" para evitar que el foco salga del modal activo.

### 1.3 Navegación por Componentes (Keyboard Navigation)
- **Tablas y Listados (`BaseCatalog`):**
  - Las filas deben ser seleccionables con `Up`/`Down`.
  - `Enter` para acceder al detalle de la fila seleccionada.
  - El foco debe permanecer en la fila seleccionada visualmente.
- **Formularios Dinámicos:**
  - En el TPV o creación de líneas de pedido, el foco debe saltar automáticamente al campo de "Cantidad" tras añadir un producto.
  - Soporte de `Esc` para cancelar la edición actual o cerrar diálogos de selección.

## 2. Estandarización de la Iconografía

Se confirma el uso de **Lucide Icons** como la librería estándar para el proyecto, reemplazando cualquier referencia residual a Material Symbols.

**Referencia:** [Guía de Lucide Icons](../architecture/design-system/lucide-icons-guide.md)

**Mapeo de Sustitución:**
- 🔧 (Tangible) -> `Package`
- ⚙️ (Servicios) -> `Settings2`
- 🗑️ -> `Trash2`
- 🖨️ -> `Printer`
- 💰 -> `Euro`
- ⚠️ -> `AlertTriangle`

## 3. Arquitectura de Componentes Consolidada

### 3.1 Unificación de Cabeceras
- **Acción:** Eliminar `PageHeader.vue` (layout) y centralizar en `BasePageHeader.vue` (shared).
- **Mejora:** El `BasePageHeader` debe incluir indicadores de atajos de teclado y utilizar componentes Lucide para la iconografía.

### 3.2 Listados Eficientes (BaseCatalog)
- Los listados deben soportar navegación por teclado de serie.
- Las filas de la tabla (`.row-clickable`) deben ser accesibles vía teclado (usar `tr tabindex="0"` y manejar el evento `keydown.enter`).

## 4. Próximos Pasos Técnicos

1. **CSS Global:** Definir estilos de foco en `_base.css`.
2. **Refactor de Iconos:** Migrar de `material-symbols-outlined` a componentes de `lucide-vue-next` en toda la aplicación.
3. **Lógica de Teclado:** Implementar un composable `useShortcuts` para gestionar los atajos de teclado de forma centralizada.
4. **Dashboard:** Refactorizar `SalesDashboard.vue` para que sea totalmente navegable por teclado.
