# Guía de Uso de Lucide Icons

## 1. Introducción

Este documento detalla la implementación y uso de Lucide Icons en el proyecto frontend de TramaTex, reemplazando el sistema anterior basado en emojis. Lucide Icons ha sido seleccionado por su estilo sobrio y profesional, su licencia permisiva (ISC) y su excelente integración con Vue.js y Tailwind CSS.

## 2. Instalación

Lucide Icons se instala como un paquete de Node.js en el proyecto `apps/frontend`.

```bash
npm install lucide-vue-next
```

## 3. Uso en Componentes Vue

Para utilizar un icono de Lucide en un componente Vue, se debe importar el componente del icono específico y luego usarlo en la plantilla.

### 3.1. Importación

Importa los iconos necesarios desde `lucide-vue-next` en la sección `<script setup>` de tu componente:

```typescript
import { Home, Package, Users } from 'lucide-vue-next';
```

### 3.2. Uso en Plantillas

Una vez importados, los iconos pueden usarse como componentes normales de Vue. Puedes controlar su tamaño (en píxeles) con la prop `:size` y el color con la prop `:color`.

```vue
<template>
  <div>
    <Home :size="24" color="blue" />
    <Package :size="32" />
    <Users :size="20" />
  </div>
</template>
```

### 3.3. Estilización con CSS / Tailwind CSS

Los componentes de Lucide renderizan SVG. Para estilizarlos con CSS o Tailwind CSS, puedes aplicar clases directamente al componente Lucide, o usar `:deep(svg)` en `<style scoped>` para apuntar al elemento SVG interno.

**Ejemplo de CSS (como se hizo en Navbar.vue y DashboardPage.vue):**

Para `Navbar.vue`:
```css
.nav-link :deep(svg) { /* Target Lucide SVG directly */
  width: 1.5rem; /* Equivalent to 24px */
  height: 1.5rem; /* Equivalent to 24px */
}

.dropdown-item :deep(svg) { /* Target Lucide SVG directly */
  width: 1.25rem; /* Equivalent to 20px */
  height: 1.25rem; /* Equivalent to 20px */
}

@media (max-width: 768px) {
  .nav-link :deep(svg) {
    width: 1.25rem;
    height: 1.25rem;
  }
  
  .dropdown-item :deep(svg) {
    width: 1rem;
    height: 1rem;
  }
}
```

Para `DashboardPage.vue`:
```css
.area-icon :deep(svg) {
  width: 2.5rem; /* Equivalent to 40px */
  height: 2.5rem; /* Equivalent to 40px */
}

.link-primary .link-icon :deep(svg) {
  width: 1.5rem; /* Equivalent to 24px */
  height: 1.5rem; /* Equivalent to 24px */
}

.link-secondary .link-icon :deep(svg) {
  width: 1.25rem; /* Equivalent to 20px */
  height: 1.25rem; /* Equivalent to 20px */
}

@media (max-width: 768px) {
  .area-icon :deep(svg) {
    width: 2rem;
    height: 2rem;
  }
}
```
## 4. Mapeo de Iconos

Se ha realizado un mapeo de los antiguos iconos emoji a los nuevos iconos de Lucide. Este mapeo se encuentra detallado en las tareas específicas de la sesión "UI Icons Review & Standardization" en `docs/log/session-log.md`.

## 5. Próximos Pasos

Para una lista completa de iconos disponibles, consulta la documentación oficial de Lucide Icons: [https://lucide.dev/icons/](https://lucide.dev/icons/)
