# Especificación: Indicadores de Atajos en BasePageHeader

Este documento define cómo el componente `BasePageHeader` debe mostrar visualmente los atajos de teclado disponibles para mejorar la accesibilidad y velocidad de los usuarios expertos.

## 1. Concepto Visual

Los atajos de teclado se representarán mediante etiquetas `<kbd>` estilizadas, situadas junto al texto de la acción o en un área dedicada de la cabecera.

## 2. Estructura de Props en `BasePageHeader`

Se propone añadir una prop `shortcuts` opcional para mapear acciones comunes:

```typescript
interface ShortcutHint {
  key: string;       // Ej: 'N'
  modifier: string;  // Ej: 'Alt'
  label: string;     // Ej: 'Nuevo'
  action: string;    // ID de la acción o ruta
}

// Props de BasePageHeader
const props = defineProps<{
  // ... props existentes
  shortcuts?: ShortcutHint[];
}>();
```

## 3. Visualización en el UI

### 3.1 Atajos de Acciones (Botones)
Para los botones en el slot `#actions`, se recomienda usar una convención de nombre o un componente dedicado `BaseShortcutIcon`.

**Ejemplo de uso en `BaseCatalog.vue`:**
```vue
<button class="btn btn-primary btn-sm">
  <span class="material-symbols-outlined">add</span>
  Nuevo
  <kbd class="shortcut-badge">Alt+N</kbd>
</button>
```

### 3.2 Estilo CSS (Global)
```css
kbd.shortcut-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  padding: 0 4px;
  font-family: var(--font-family-mono);
  font-size: 0.65rem;
  color: var(--color-text-secondary);
  margin-left: 0.5rem;
  box-shadow: 0 1px 0 rgba(0,0,0,0.1);
  min-width: 1.5rem;
}
```

## 4. Mapeo Estándar de Atajos

| Acción | Atajo | Componente / Contexto |
| :--- | :--- | :--- |
| Nueva Entidad | `Alt + N` | Listados (BaseCatalog) |
| Guardar | `Ctrl + Enter` | Formularios |
| Cancelar | `Esc` | Modales / Edición |
| Buscar | `Ctrl + K` | Global |
| Refrescar | `Alt + R` | Listados |

## 5. Implementación en `BasePageHeader`

El componente `BasePageHeader` mostrará automáticamente una leyenda de atajos si se proporcionan a través de las props, o permitirá que los componentes hijos los incluyan en sus propios botones.

```vue
<!-- En BasePageHeader.vue -->
<div v-if="shortcuts" class="header-shortcuts-legend">
  <div v-for="s in shortcuts" :key="s.key" class="shortcut-item">
    <kbd>{{ s.modifier }}+{{ s.key }}</kbd> <span>{{ s.label }}</span>
  </div>
</div>
```
