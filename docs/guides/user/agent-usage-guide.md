# Guía de Usuario: Cómo Interactuar con los Agentes de IA

**Audiencia:** Desarrolladores, Gerentes de Proyecto, Usuarios Finales (con interés técnico)
**Versión:** 1.0
**Fecha:** 2026-02-04
**Autor:** Gemini

---

## 1. 🎯 Propósito de los Agentes de IA

Los Agentes de IA son herramientas inteligentes diseñadas para asistirte en diversas tareas del proyecto. Actúan como expertos especializados que pueden guiarte, proporcionarte información contextual, ejecutar tareas de desarrollo o ayudarte a navegar por la complejidad del proyecto.

El objetivo principal es:
-   **Centralizar el conocimiento:** Los agentes encapsulan reglas, estándares y mejores prácticas del proyecto.
-   **Automatizar flujos de trabajo:** Ayudan a seguir procesos definidos (ej. creación de módulos, gestión de sprints).
-   **Proporcionar guía:** Te asisten cuando no estás seguro de cómo proceder o necesitas información específica.

---

## 2. 🤖 El Agente de Entrada: `init.yaml`

El agente `init.yaml` es tu punto de partida. Cuando inicias una sesión con la IA, este agente te preguntará qué deseas hacer y te dirigirá al agente o flujo de trabajo más adecuado.

### ¿Cómo funciona?

1.  Al iniciar, `init.yaml` cargará automáticamente las reglas generales del proyecto y el contexto de alto nivel.
2.  Te presentará un menú con opciones comunes de tareas.
3.  Basándose en tu elección, activará otro agente más especializado o te proporcionará la información necesaria.

### Opciones Comunes

Aquí tienes algunas de las opciones que `init.yaml` puede ofrecerte:

-   **Tareas de Desarrollo (Sprints y Tareas):** Te dirigirá al agente `load-session.yaml`, que te ayudará a retomar un trabajo pendiente o a empezar una nueva tarea de sprint siguiendo el `Standard Module Development Workflow (SMDW)`.
-   **Documentación del Proyecto:** Te indicará dónde encontrar la documentación principal (`docs/README.md`) y cómo navegar por ella.
-   **Información Arquitectónica o Estándares de Código:** Te guiará para cargar agentes de contexto específicos (ej. `architecture.yaml`, `code-standards.yaml`) o te ayudará a usar el `codebase_investigator` para preguntas detalladas.

---

## 3. 📚 Otros Agentes Relevantes

El sistema utiliza varios agentes auxiliares para funciones específicas:

-   **`generic-rules.yaml`**: Define las reglas fundamentales del proyecto (idioma, estructura de directorios, principios de calidad). Este agente se carga automáticamente al inicio de cualquier sesión.
-   **`project/project-context.yaml`**: Contiene el contexto de alto nivel del proyecto (visión, stack tecnológico, bounded contexts). Se carga a través de `load-project-context.yaml`.
-   **`load-project-context.yaml`**: Se encarga de cargar el contexto principal del proyecto y te ofrece cargar contextos modulares específicos (ej. `architecture.yaml`, `bounded-contexts.yaml`).
-   **`init-session.yaml`**: Gestiona las sesiones de trabajo. Te permite continuar tareas pendientes de `session-log.md` o iniciar nuevas tareas de desarrollo.
-   **Agentes de Contexto (`agents/project/context/*.yaml`)**: Proporcionan información detallada sobre aspectos específicos del proyecto (ej. `architecture.yaml`, `tech-stack.yaml`, `code-standards.yaml`). Se cargan bajo demanda.

---

## 4. 📝 La Guía Maestra: `AGENTS.md`

El archivo `AGENTS.md` (en la raíz del proyecto) es la guía maestra de todos los agentes. Contiene:
-   Tu "persona" como asistente de IA.
-   Los principios core que rigen todas las interacciones.
-   El flujo de trabajo estándar para cualquier tarea.
-   Una descripción general de cómo se relacionan los diferentes agentes.

**Siempre que tengas dudas sobre el comportamiento de la IA o cómo se configuran los agentes, consulta `AGENTS.md`.**

---

## 5. 💡 Consejos para una Interacción Efectiva

-   **Sé claro y conciso:** Formula tus preguntas y peticiones de la manera más directa posible.
-   **Proporciona contexto:** Si tu petición es sobre una parte específica del código o la documentación, menciónalo.
-   **Responde con números de opción:** Cuando se te presente un menú, responde con el número correspondiente a tu elección (ej. "1" para "Continuar una tarea").
-   **No tengas miedo de preguntar:** Si no entiendes algo o necesitas más detalles, la IA está aquí para ayudarte.

---

## 6. 🔗 Referencias

-   [AGENTS.md](../../../AGENTS.md) (Guía Maestra de Agentes)
-   [generic-rules.yaml](../../../agents/generic-rules.yaml) (Reglas Fundamentales del Proyecto)
-   [load-session.yaml](../../../agents/load-session.yaml) (Agente de Gestión de Sesiones)
-   [ADR-009: Estructura de Carpetas y Organización](../../architecture/adrs/adr-009-project-structure.md) (Estructura de `docs/` y `agents/`)

---
