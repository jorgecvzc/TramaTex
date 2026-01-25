# Especificación de Tipografía de TramaTex (v2 - Mockup)

---

## ✒️ Filosofía de Tipografía

La tipografía se alinea con el diseño de los mockups, buscando ser **moderna, funcional y con una clara jerarquía visual**. Se utiliza una fuente de marca distintiva para el logo y una fuente de sistema altamente legible para la interfaz.

---

## 🖋️ Familias de Fuentes

| Familia | Nombre | Variable CSS | Uso Principal |
|---|---|---|---|
| **Sans-serif** | `Inter` | `var(--font-family-sans)` | La fuente principal para toda la interfaz, desde párrafos hasta etiquetas y botones. |
| **Brand** | `Calibri`, etc. | `var(--font-family-brand)` | Usada exclusivamente para el logo "TramaTex", dándole un carácter distintivo e itálico. |
| **Monospace** | `Fira Code` | `var(--font-family-mono)` | Para mostrar IDs técnicas o fragmentos de código. |

*Nota: La fuente `Inter` se importa desde Google Fonts.*

---

## 📏 Escala Tipográfica y Pesos

La escala se mantiene consistente y basada en `rem`.

### Tamaños de Fuente

| Tamaño | Variable CSS | Valor (px) | Uso Común |
|---|---|---|---|
| `xs` | `var(--font-size-xs)` | 12px | Texto muy pequeño, metadatos, etiquetas. |
| `sm` | `var(--font-size-sm)` | 14px | Texto de cuerpo secundario, descripciones. |
| `md` (base) | `var(--font-size-md)` | 16px | Párrafos y texto principal. |
| `lg` | `var(--font-size-lg)` | 18px | Encabezados de tarjetas o secciones. |
| `xl` | `var(--font-size-xl)` | 20px | Encabezados de nivel inferior. |

### Pesos de Fuente

| Peso | Variable CSS | Valor Numérico | Uso Común |
|---|---|---|---|
| `Normal` | `var(--font-weight-normal)` | 400 | Cuerpo del texto. |
| `Medium` | `var(--font-weight-medium)` | 500 | Texto con ligero énfasis. |
| `Bold` | `var(--font-weight-bold)` | 700 | Encabezados. |
| `Black` | `var(--font-weight-black)` | 900 | Títulos muy destacados o cifras importantes. |

---

## 📄 Estilos de Elementos HTML

### Encabezados (`h1`-`h4`)

-   **Fuente**: `Inter`
-   **Peso**: `Bold` (700)
-   **Color**: `var(--color-secondary)` (Azul Intenso)
-   **Uso**: Para titular páginas y secciones, siguiendo la paleta de colores de los mockups.

### Párrafos (`p`)

-   **Fuente**: `Inter`
-   **Peso**: `Normal` (400)
-   **Color**: `var(--color-text-primary)`

### Enlaces (`a`)

-   **Color**: `var(--color-secondary)` (Azul Intenso)
-   **Estado `hover`**: Un azul ligeramente más claro (`var(--color-info)`) para retroalimentación.
