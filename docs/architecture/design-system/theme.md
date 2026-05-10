# Especificación del Tema General de la UI de TramaTex (v2 - Mockup)

---

## 🎨 Tema: "Claridad y Foco Profesional"

El tema, alineado con los mockups, es **vibrante, de alto contraste y profesional**. El diseño está pensado para una aplicación de escritorio intensiva en datos, donde la claridad, la jerarquía y la eficiencia son cruciales.

---

## 🔑 Principios Clave del Diseño

-   **Jerarquía Visual Fuerte**:
    -   El uso del azul intenso (`--color-secondary`) para cabeceras y elementos de navegación establece una jerarquía clara.
    -   El amarillo (`--color-primary`) se reserva para las acciones más importantes (CTAs) para atraer la atención del usuario.

2.  **Arquitectura en 3 Capas (Master Layout)**:
    -   **Capa 1: Identidad (Sticky)**: Cabecera superior fija con título, ID y acciones globales sobre fondo blanco.
    -   **Capa 2: Contexto (Summary)**: Cinta de KPIs, trazabilidad y estado sobre fondo gris ceniza.
    -   **Capa 3: Trabajo (Main)**: Área operativa con formularios, tablas de líneas y datos maestros.

3.  **Filosofía Keyboard-First (Industrial UX)**:
    -   Diseñado para ser operado al 100% sin ratón.
    -   Navegación intensiva mediante atajos (`Alt+1-5`, `F1`, `?`) y control experto de tablas de líneas (`Enter` continuo, flechas de incremento).

4.  **Diseño Basado en Tarjetas**:
    -   El contenido principal se organiza en tarjetas (`--color-surface`) que se elevan sobre un fondo ligeramente gris (`--color-background`).

5.  **Iconografía Lucide (SVG)**:
    -   Uso exclusivo de la librería `lucide-vue-next` para una iconografía nítida, moderna y técnicamente superior a las fuentes de iconos tradicionales.

---

## 🖼️ Look & Feel General

-   **Energético y Profesional**: La combinación de azul intenso y amarillo huevo es audaz, pero se equilibra con una gran cantidad de espacio en blanco y grises neutros.
-   **Orientado a la Productividad Industrial**: El diseño no es decorativo; está optimizado para que los usuarios encuentren información y completen tareas rápidamente.
-   **Moderno y Limpio**: El uso de una fuente sans-serif moderna (`Inter`) y bordes redondeados sutiles (`--border-radius-sm`) le da un aspecto actual.


---

## 📁 Estructura de Layout (Definida por Mockups)

-   **Navegación Principal (Sidebar)**:
    -   A la izquierda, con un fondo blanco (`--color-surface`) y un borde derecho.
    -   El logo de la marca se sitúa en la parte superior sobre un fondo azul (`--color-secondary`).
    -   El elemento de menú activo se resalta con el color azul secundario.

-   **Barra Superior (Header)**:
    -   Una cabecera horizontal en la parte superior del área de contenido.
    -   Contiene "breadcrumbs" para la navegación contextual y controles de usuario (notificaciones, perfil).

-   **Área de Contenido**:
    -   El espacio principal, con un fondo gris claro (`--color-background`).
    -   El contenido se presenta en tarjetas blancas con sombras sutiles, creando una sensación de profundidad y organización.
