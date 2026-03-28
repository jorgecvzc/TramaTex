# Guía de Átomos: TramaTex Design System

Esta guía documenta los elementos básicos (átomos) que componen la interfaz de TramaTex. Es obligatorio seguir estos estándares para garantizar la coherencia visual y técnica del proyecto.

## 1. Paleta de Colores

Se utilizan variables CSS para asegurar la consistencia y facilitar el mantenimiento.

| Color | Variable | Uso Principal |
| :--- | :--- | :--- |
| **Amarillo TramaTex** | `--color-primary` | Acciones principales, botones destacados, acentos. |
| **Azul Corporativo** | `--color-secondary` | Identidad de marca, cabeceras, enlaces de navegación. |
| **Verde Éxito** | `--color-success` | Estados positivos, facturas, validaciones. |
| **Rojo Error** | `--color-error` | Alertas, acciones destructivas, errores de validación. |
| **Gris Base** | `--color-background` | Fondo general de la aplicación. |
| **Blanco Superficie** | `--color-surface` | Fondo de tarjetas y contenedores de contenido. |
| **Gris Borde** | `--color-border` | Divisores y siluetas de componentes. |

## 2. Tipografía

El sistema utiliza una jerarquía clara para guiar la lectura.

*   **Fuente de Marca (`.font-brand`):** Reservada para logotipos y grandes títulos de bienvenida.
*   **Encabezados (`h1`, `h2`, `h3`):**
    *   `h1`: Títulos de página principales.
    *   `h2`: Títulos de sección dentro de entidades.
    *   `h3`: Títulos de tarjetas o bloques menores.
*   **Cuerpo de Texto:** Peso normal para legibilidad.
*   **Texto Muted (`.text-muted`):** Gris suave para información secundaria.
*   **Monoespaciado (`code`, `.text-mono`):** Para SKUs, IDs técnicos y referencias.

## 3. Botones (`.btn`)

Todos los botones deben heredar de la clase base `.btn`.

### Tipos Semánticos
| Clase | Propósito | Estética |
| :--- | :--- | :--- |
| `.btn-primary` | Acción principal del flujo. | Fondo Amarillo, Texto Oscuro. |
| `.btn-secondary` | Acciones de proceso/negocio. | Fondo Azul, Texto Blanco. |
| `.btn-outline` | Acciones neutras o secundarias. | Borde gris, fondo blanco. |
| `.btn-danger` | Acciones destructivas o críticas. | Texto rojo, borde rojo suave. |
| `.btn-success` | Acciones de confirmación positiva. | Fondo verde, texto blanco. |
| `.btn-ghost` | Acciones mínimas o iconos. | Sin fondo, resalta al hover. |

### Tamaños y Estados
*   **`.btn-sm`**: Para tablas y áreas compactas.
*   **`.btn-lg`**: Para pies de formulario y acciones críticas.
*   **`:disabled`**: Reduce opacidad y bloquea interacción.

## 4. Iconografía

**ESTÁNDAR OBLIGATORIO:** Material Symbols Outlined.

*   **Formato:** `<span class="material-symbols-outlined">icon_name</span>`
*   **Integración en Botones:** Se colocan antes del texto con un espacio de separación automático.
*   **Tamaño Estándar:** `24px` para UI general, `20px` dentro de botones.

## 5. Elevación y Bordes

Para evitar la "dilución" visual en distintos monitores:
*   **Bordes en Main Content:** `1px solid #d1d5db`.
*   **Sombras en Main Content:** Sombra multinivel (dispersión + contacto) para profundidad real.
*   **Radio de Borde:** `var(--border-radius-md)` (normalmente 8px-12px) para suavizar la interfaz.
