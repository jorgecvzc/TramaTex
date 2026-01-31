# 🧑‍💻 Guía para el Usuario sobre el Sistema de Agentes

Este documento explica el propósito y funcionamiento del sistema de agentes de IA para los usuarios y desarrolladores del proyecto.

## ¿Qué son los "Agentes"?

Los "agentes" son una colección de archivos de configuración (`.yaml` y `.md`) ubicados en el directorio `/agents`. Actúan como la **memoria a largo plazo** y el **conjunto de reglas** para el asistente de IA que trabaja en este proyecto.

Definen su personalidad, sus principios, la arquitectura del proyecto, los estándares de código y el estado actual de los sprints.

## ¿Cómo funcionan?

Al inicio de una sesión, el asistente de IA carga estos archivos en su contexto. Esto le permite:
-   **Tomar decisiones consistentes** basadas en las reglas del proyecto (ej. `generic-rules.yaml`).
-   **Entender la arquitectura** y las tecnologías sin tener que re-analizar el código cada vez (`architecture.yaml`, `tech-stack.yaml`).
-   **Conocer el estado de un sprint** y las tareas pendientes (`sprint-registry.yaml`, `sprint-session-loader.yaml`).

El usuario no necesita gestionar estos archivos directamente, pero puede consultarlos para entender el "pensamiento" del asistente.

## Prompt para Iniciar una Nueva Sesión de Trabajo

Para iniciar el flujo de trabajo estándar y que el asistente sepa qué hacer, el prompt más efectivo es:

> **"Dónde estábamos en @NEXT_SESSION.md"**

**¿Por qué este prompt?**

1.  **Carga el archivo de continuidad**: Hace que el asistente revise primero `NEXT_SESSION.md`, que es el punto de partida prioritario.
2.  **Activa el flujo estándar**: Si `NEXT_SESSION.md` está vacío, el asistente sabe que debe proceder a cargar el estado del sprint actual (`sprint-session-loader.yaml`).
3.  **Es claro y directo**: No deja lugar a ambigüedad y le indica al asistente que se prepare para continuar con el trabajo pendiente.

## Documentación Técnica de la IA

Para una descripción detallada de la persona, los principios y el flujo de trabajo interno de la IA, consulta el archivo:
-   **[agents/README.md](./agents/README.md)**

---
*Última actualización: 2026-01-31*