# Guía de Uso de Lucide Icons

## 1. Introducción

Este documento detalla la implementación y uso de Lucide Icons en el proyecto frontend de TramaTex, reemplazando el sistema anterior basado en emojis. Lucide Icons ha sido seleccionado por su estilo sobrio y profesional, su licencia permisiva (ISC) y su excelente integración con Vue.js.

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

### 3.3. Estilizatión con CSS

Los componentes de Lucide renderizan SVG. Para estilizarlos con CSS, puedes aplicar clases directamente al componente Lucide, o usar `:deep(svg)` en `<style scoped>` para apuntar al elemento SVG interno.

## 4. Mapeo de Emojis a Lucide

| Emoji | Lucide Icon | Uso |
| :--- | :--- | :--- |
| 🔧 | `Package` | Productos Tangibles |
| ⚙️ | `Settings2` | Servicios |
| 🗑️ | `Trash2` | Eliminar |
| 🖨️ | `Printer` | Imprimir |
| 💰 | `Euro` | Precios / Pagos |
| ⚠️ | `AlertTriangle` | Avisos / Alertas |
| 📦 | `Box` | Albaranes / Stock |
| 🚚 | `Truck` | Logística |

## 5. Referencia Oficial

Para una lista completa de iconos disponibles, consulta la documentación oficial de Lucide Icons: [https://lucide.dev/icons/](https://lucide.dev/icons/)
