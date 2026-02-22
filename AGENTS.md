# 🧑‍💻 Guía para el Usuario sobre el Sistema de Agentes

> **📖 Guías Rápidas:**
> - **Para usuarios:** Ver [docs/guides/user/guia-agents.md](docs/guides/user/guia-agents.md) (guía concisa ≤40 líneas)
> - **Para desarrolladores:** Ver [docs/guides/developer/agent-development-guide.md](docs/guides/developer/agent-development-guide.md) (crear/modificar agentes)

Este documento es la referencia completa del sistema de agentes de IA para usuarios y desarrolladores del proyecto.

## ¿Qué son los "Agentes"?

Los "agentes" son una colección de archivos de configuración (`.yaml` y `.md`) ubicados en el directorio `/agents`. Actúan como la **memoria a largo plazo** y el **conjunto de reglas** para el asistente de IA que trabaja en este proyecto.

Definen su personalidad, sus principios, la arquitectura del proyecto, los estándares de código y el estado actual de los sprints.

## ¿Cómo funcionan?

Al inicio de cada sesión, el asistente de IA carga el agente `init.yaml` como su punto de entrada principal. Este agente interactúa contigo para entender tu objetivo y te dirige al agente o flujo de trabajo más adecuado. Esto le permite:
-   **Tomar decisiones consistentemente** basadas en las reglas del proyecto (ej. `generic-rules.yaml`).
-   **Entender la arquitectura** y las tecnologías sin tener que re-analizar el código cada vez (`architecture.yaml`, `tech-stack.yaml`).
-   **Conocer el estado de un sprint** y las tareas pendientes (`sprint-registry.yaml`, `load-session.yaml`).

El usuario no necesita gestionar estos archivos directamente, pero puede consultarlos para entender el "pensamiento" del asistente.

## Estructura de Agentes Clave

-   `init.yaml`: El agente de entrada principal. Gestiona la interacción inicial con el usuario y redirige a otros agentes.
-   `generic-rules.yaml`: Define las reglas universales del proyecto, como los idiomas para código y documentación, la estructura de directorios y los principios de arquitectura.
-   `project-context.yaml`: Proporciona el contexto de alto nivel del proyecto TramaTex, incluyendo su visión, stack tecnológico y módulos principales.
-   `sprint-registry.yaml`: Es el registro central de todos los sprints y tareas. Define qué tareas están activas, completadas o pendientes.
-   `init-session.yaml`: Contiene la lógica que el asistente de IA usa para gestionar sesiones de desarrollo (sprints y tareas), incluyendo la lectura de `session-log.md`.
-   `agents/project/context/`: Contiene agentes modulares con información detallada sobre aspectos específicos como la arquitectura, los bounded contexts, el stack tecnológico y los estándares de código.

## Prompt para Iniciar una Nueva Sesión de Trabajo

Para iniciar el flujo de trabajo estándar y que el asistente sepa qué hacer, simplemente inicia tu conversación con el asistente. El agente `init.yaml` te guiará automáticamente.

**¿Por qué este enfoque?**

1.  **Guía Activa**: El agente `init.yaml` te presentará opciones claras para comenzar, haciendo más fácil la interacción, especialmente para nuevos usuarios.
2.  **Activación Centralizada**: `init.yaml` se encarga de cargar los contextos necesarios y de delegar la tarea al agente especializado (como `load-session.yaml` para gestionar sesiones).
3.  **Claridad y Consistencia**: Asegura que el flujo de trabajo siempre comienza de la misma manera, proporcionando una experiencia consistente.

## Documentación Técnica de la IA

Para una descripción detallada de los agentes de contexto específicos del proyecto, consulta el archivo:
-   **[agents/project/context/README.md](./agents/project/context/README.md)**

---

## 📚 Recursos Adicionales

- **[docs/guides/user/guia-agents.md](docs/guides/user/guia-agents.md)** - Guía de uso concisa (inicio rápido)
- **[docs/guides/developer/agent-development-guide.md](docs/guides/developer/agent-development-guide.md)** - Guía para desarrollar agentes
- **[agents/project/context/README.md](agents/project/context/README.md)** - Documentación de contextos modulares

---
*Última actualización: 2026-02-05*