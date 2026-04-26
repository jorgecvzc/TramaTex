# Plan de Refactorización: ProductGroupForm.vue (Iconografía)

Este documento detalla los pasos para eliminar el uso de emojis en `ProductGroupForm.vue` y sustituirlos por componentes de la librería **Lucide Icons**, siguiendo el nuevo estándar del sistema de diseño.

## 1. Situación Actual

El componente `ProductGroupForm.vue` utiliza emojis directamente en las etiquetas de los botones de radio para distinguir visualmente entre "Productos Tangibles" y "Servicios".

```vue
<span class="radio-title">🔧 Productos Tangibles</span>
<span class="radio-title">⚙️ Servicios</span>
```

## 2. Propuesta de Cambio

### 2.1 Dependencias
- Requiere la instalación de `lucide-vue-next` en el proyecto frontend.

### 2.2 Importación de Componentes
En la sección `<script setup>`, importar los iconos correspondientes:

```typescript
import { Package, Settings2 } from 'lucide-vue-next';
```

### 2.3 Modificación del Template
Sustituir los emojis por los componentes de Lucide con los estilos adecuados.

```vue
<div class="radio-title">
  <Package :size="18" class="icon-inline" /> 
  <span>Productos Tangibles</span>
</div>

<div class="radio-title">
  <Settings2 :size="18" class="icon-inline" /> 
  <span>Servicios</span>
</div>
```

### 2.4 Ajustes de Estilo (CSS)
Añadir clases para alinear correctamente los iconos con el texto:

```css
.icon-inline {
  vertical-align: middle;
  margin-right: 0.5rem;
  color: var(--color-primary);
}

.radio-title {
  display: flex;
  align-items: center;
}
```

## 3. Mapeo de Emojis a Iconos

| Antiguo Emoji | Nuevo Icono Lucide | Color Sugerido |
| :--- | :--- | :--- |
| 🔧 | `Package` | `--color-primary` (Oro) |
| ⚙️ | `Settings2` | `--color-secondary` (Azul Profundo) |

## 4. Beneficios
- Apariencia profesional y consistente con el resto de la aplicación.
- Mejor legibilidad en diferentes sistemas operativos (los emojis varían mucho de estilo).
- Mayor control sobre el tamaño y color de los iconos.
