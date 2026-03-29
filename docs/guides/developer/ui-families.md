# Familias de UI: TramaTex Architecture

Este documento define la taxonomía oficial de las interfaces de TramaTex. Toda página añadida al proyecto **debe** pertenecer a una de estas familias para garantizar la coherencia visual, técnica y de experiencia de usuario.

---

## 🏗️ Jerarquía del Sistema
1.  **Átomos (Design System)**: Colores, tipografía y componentes básicos. Fuente de verdad: `_variables.css`.
2.  **Familias (Layout Structures)**: Planos de construcción que organizan los átomos según el propósito de la página.

---

## 1. Familia: Dashboards (Visión Estratégica)
Interfaces híbridas diseñadas para el escritorio y dispositivos móviles (Tablets). Utilizan el componente base `BaseDashboardPage.vue`.

*   **Propósito**: Consultar el estado del negocio, analizar datos y realizar cálculos rápidos.
*   **Layout**: Ancho máximo de 1300px, Sidebar de contexto (derecha), Área principal de trabajo (izquierda).
*   **Sub-Familias**:
    *   **Informativos**: Resumen de KPIs, métricas en tiempo real y listas de actividad reciente (ej. Dashboard Principal, Monitor MES).
    *   **Funcionales**: Herramientas interactivas de proyección o simulación. El input se sitúa preferiblemente en el sidebar y el resultado en el área principal (ej. Consulta de Precios).

---

## 2. Familia: Gestión de Entidades (Ciclo de Vida)
Interfaces orientadas a la manipulación precisa de datos de negocio. Siguen el estándar CRUD.

*   **Propósito**: Crear, listar, filtrar y editar los registros maestros del sistema.
*   **Sub-Familias**:
    *   **Catálogos (Listados)**: Utilizan `BaseCatalog.vue`. Tablas densas, filtros potentes y búsqueda global (ej. Listado de Productos, Clientes).
    *   **Fichas de Detalle**: Lectura profunda de un registro único, con historial y relaciones (ej. Detalle de Pedido).
    *   **Formularios Transaccionales**: Flujos de entrada de datos optimizados para velocidad y validación (ej. Nuevo Pedido).

---

## 3. Familia: Viewports Especializados (Contexto de Uso)
Interfaces que rompen el layout global para adaptarse a un entorno físico o audiencia específica.

*   **Propósito**: Maximizar la eficiencia en condiciones especiales donde el ratón/teclado no son la entrada principal.
*   **Sub-Familias**:
    *   **Terminales Operativos**: Diseñados para pantallas táctiles en fábrica o mostradores. Sin navegación global, botones gigantes y alto contraste (ej. Terminal de Taller, TPV).
    *   **Showcases de Marca**: Herramientas de venta y presentación. Layout aireado, tipografía cuidada y enfoque estético (ej. Design System).
    *   **Utilidades de Sistema**: Vistas técnicas para administradores (ej. Logs de error, Diagnóstico).

---

## 📏 Reglas de Oro
1.  **Finitud**: No se permiten "páginas libres". Si una página no encaja, se debe proponer una nueva sub-familia.
2.  **Responsive**: Los Dashboards y Entidades deben ser 100% usables en tablets. Los Terminales son específicos de su dispositivo.
3.  **Herencia**: Todas las familias consumen obligatoriamente las variables de `design-system-atoms.md`.
