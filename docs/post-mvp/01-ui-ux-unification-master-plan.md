# Plan Maestro: Unificación UI/UX y Accesibilidad Universal (Post-MVP)

Este documento es la **fuente única de verdad** para la transformación de la interfaz de TramaTex en una herramienta industrial de alta eficiencia. Consolida todos los estudios técnicos, especificaciones de diseño y planes de migración de la fase Post-MVP.

---

## 1. Estrategia de Accesibilidad (Keyboard-First)

El objetivo es que un operario realice el flujo completo (Ventas -> Producción) sin usar el ratón.

### 1.1 Atajos Globales y Foco
- **Búsqueda Global**: `Ctrl + K`.
- **Nuevas Entidades**: `Alt + N`.
- **Guardar/Confirmar**: `Ctrl + Enter`.
- **Refrescar Datos**: `Alt + R`.
- **Indicador de Foco**: Implementar un estilo `:focus-visible` vibrante en `--color-primary`.
- **Modales**: Implementar "Focus Trap" y cierre con `Esc`.

### 1.2 Listados Maestros (BaseCatalog)
- **Navegación**: Flechas `Up`/`Down` para mover el resaltado (`selectedIndex`).
- **Selección**: `Enter` para acceder al detalle.
- **Scroll**: Implementar `scrollToSelected` para mantener la fila activa visible.

### 1.3 Líneas de Documentos (Ventas/Pedidos)
- **Navegación Celdas**: `Tab` o `Flechas Laterales` entre Producto, Cantidad y Precio.
- **Ajuste Rápido**: Teclas `+` y `-` para incrementar/decrementar cantidades.
- **Continuidad**: `Enter` en la última celda confirma la línea y crea una nueva vacía.

---

## 2. Sistema de Diseño e Identidad Visual

### 2.1 Estandarización de Iconografía
Se migra de Material Symbols y emojis a **Lucide Icons** (`lucide-vue-next`).
- **Productos Tangibles**: `Package`.
- **Servicios**: `Settings2`.
- **Acciones**: `Trash2` (Borrar), `Printer` (Imprimir), `Euro` (Precios), `AlertTriangle` (Avisos).

### 2.2 Componentes Consolidados
- **BasePageHeader**: Sustituye a `PageHeader`. Incluye indicadores de atajos `<kbd>` estilizados.
- **BaseDashboardPage**: Estándar para todos los paneles de módulo (Ventas, Productos, Entidades, MES).
- **Dashboard Global**: Debe refactorizarse para usar `BaseDashboardPage`.

---

## 3. Ergonomía Industrial y Feedback (Nuevas Mejoras)

1.  **Notificaciones (Toasts)**: Store global de Pinia para mensajes no intrusivos (Sustituye a `alert`).
2.  **Campos Calculados**: Clase `.input-calculated` (fondo gris suave) para datos protegidos de backend.
3.  **Skeleton Loaders**: Estructuras de carga que imitan el contenido para reducir fatiga visual.
4.  **Chips de Filtros**: Etiquetas eliminables en la cabecera de las tablas para ver filtros activos.
5.  **Tooltips Dinámicos**: Mostrar el atajo de teclado (`[Alt+N]`) al pasar el ratón.
6.  **Validación Inline**: Indicadores rojo/verde inmediatamente tras perder el foco (`on-blur`).
7.  **Modo Industrial (MES)**: Versión de alto contraste para condiciones de iluminación difíciles en taller.

---

## 4. Plan de Implementación (Orden Sugerido)

1.  **Fundamentos**: Instalación de `lucide-vue-next` y creación de `_buttons.css` y `_dashboards.css`.
2.  **Feedback Crítico**: Implementación del sistema de Toasts y validación inline.
3.  **Navegación Core**: Lógica de teclado en `BaseCatalog` y líneas de venta.
4.  **Consolidación**: Refactor de Dashboards y BasePageHeader.
5.  **Especialización**: Skeletons y Modo Industrial.

---

*Última actualización: 2026-04-26*
