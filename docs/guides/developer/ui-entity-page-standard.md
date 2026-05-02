# Estándar de Interfaz: BaseEntityPage (Diseño de 3 Capas)

La `BaseEntityPage` es la plantilla maestra obligatoria para la gestión integral de una **Entidad de Dominio** (Pedido, Producto, Cliente, Factura, etc.) a pantalla completa.

## 1. Jerarquía Visual y Cromática

El diseño se organiza en tres profundidades visuales para maximizar la claridad:

### Capa 1: Identidad Fija (`Identity Header`)
* **Color:** Blanco Puro (`#ffffff`).
* **Comportamiento:** `Sticky` (Fijo en la parte superior).
* **Contenido:** Identificación del objeto (Breadcrumbs + Título) y Acciones Globales (Guardar, Imprimir, Cerrar).
* **Propósito:** Autoridad y control siempre disponible.

### Capa 2: Consola de Contexto (`Context Header`)
* **Color:** Gris Ceniza Muy Suave (`#f9fafb`).
* **Comportamiento:** Flujo de scroll (desaparece al subir). Ocupa el 100% del ancho.
* **Contenido (en este orden):**
    1. **Toolbar:** Estado del proceso (Badge resaltado) y botones de flujo.
    2. **Summary:** Cinta de tarjetas KPI con datos clave para lectura rápida.
    3. **Related:** Tarjetas de trazabilidad (documentos vinculados).
* **Propósito:** Definir el "momento" y la "genealogía" del objeto.

### Capa 3: Área de Trabajo (`Main Content`)
* **Color:** Gris Base de la aplicación (`var(--color-background)`).
* **Contenido:** Bloques de información organizados en componentes `FormSection` y `DataRow`.
* **Estándar de Tarjeta:** Las tarjetas en esta zona deben tener:
    * **Borde:** `1px solid #d1d5db` (Reforzado para visibilidad en todo tipo de monitores).
    * **Sombra:** `box-shadow` multinivel para dar relieve real.
* **Propósito:** Operación técnica y entrada de datos.

## 2. Estándar de Iconografía

**REGLA INAMOVIBLE:** Se utiliza exclusivamente **Lucide Icons** registrados en la Verdad Única (`src/utils/icons.ts`).
* **Uso:** `<component :is="getIcon(entityIcon)" :size="28" />` en la cabecera.
* **Prohibición:** No utilizar Material Symbols ni fuentes tipográficas de iconos externos.

## 3. Ejemplo de Estructura de Código

```html
<BaseEntityPage>
  <template #header>
    <PageHeader title="Entidad #123" ... />
  </template>

  <template #toolbar>
    <div class="action-toolbar"> ... </div>
  </template>

  <template #summary>
    <div class="overview-tags-row"> ... </div>
  </template>

  <template #related>
    <div class="related-history-grid"> ... </div>
  </template>

  <!-- Cuerpo Principal (Slot por defecto) -->
  <FormSection title="Datos Maestros" icon="person"> ... </FormSection>
  <FormSection title="Líneas de Detalle" icon="list"> ... </FormSection>

  <template #footer>
    <div class="audit-info"> ... </div>
  </template>
</BaseEntityPage>
```
