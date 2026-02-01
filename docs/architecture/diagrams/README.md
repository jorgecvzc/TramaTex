# Arquitectura Visual del Sistema

Esta carpeta contiene todos los diagramas que describen la arquitectura de TramaTex en diferentes niveles de abstracción, siguiendo una adaptación del modelo C4 y diagramas UML/ERD específicos por módulo.

La herramienta estándar para crear estos diagramas es **Mermaid**, incrustada directamente en los archivos Markdown.

## 1. Diagramas de Alto Nivel (C1 y C2)

Estos diagramas describen la arquitectura general del sistema.

- **[C1-context.md](./C1-context.md):** (Ya existe, renombrado desde `context.md`) Diagrama de Contexto del Sistema (Nivel 1 de C4). Muestra cómo TramaTex se relaciona con sus usuarios y sistemas externos.
- **[C2-containers.md](./C2-containers.md):** Diagrama de Contenedores (Nivel 2 de C4). Descompone el sistema en sus principales bloques ejecutables (API, Frontend, Base de Datos).

## 2. Diagramas de Módulos (C3 y otros)

La subcarpeta `modules/` contiene diagramas detallados para cada Bounded Context, siguiendo una estructura organizada.

- **/modules/iam/:** Diagramas para el módulo de Identidad y Acceso.
- **/modules/party/:** Diagramas para el módulo de Clientes/Proveedores.
- **/modules/product/:** Diagramas para el catálogo de Productos.
- **/modules/pricing/:** Diagramas para el motor de Tarificación.
- **/modules/sales/:** Diagramas para el módulo de Ventas.
- **/modules/mes/:** Diagramas para el sistema de ejecución de manufactura.

### Tipos de Diagramas por Módulo

Cada carpeta de módulo puede contener, entre otros, los siguientes tipos de diagramas:

- **`domain-model.md`**: Un diagrama de Clases UML que muestra las entidades, value objects y agregados del dominio.
- **`use-cases.md`**: Un diagrama de Casos de Uso UML que muestra los actores y sus interacciones.
- **`state-machine.md`**: Un diagrama de Estado UML para entidades con un ciclo de vida complejo (ej. Pedido, Orden de Producción).
- **`sequence.md` o `flow.md`**: Un diagrama de Secuencia o Actividad UML para flujos de lógica complejos (ej. cálculo de precio).
- **`erd.md`**: Un Diagrama Entidad-Relación (formato Crow's Foot) que representa el esquema de base de datos para ese módulo.
- **`components.md`**: Un diagrama de Componentes (Nivel 3 de C4) que descompone un contenedor en sus principales bloques lógicos.
