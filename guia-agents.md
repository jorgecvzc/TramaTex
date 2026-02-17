# Guía de Uso - Sistema de Agentes

## 🎯 ¿Qué Son?

Archivos YAML/MD en `/agents` que definen reglas y contexto para asistentes de IA.

## 🚀 Inicio Rápido

**Simplemente inicia tu conversación** - El agente `init.yaml` te guiará automáticamente.

## 📚 Agentes Principales

| Agente | Propósito |
|--------|-----------|
| **generic-rules.yaml** | Reglas universales del proyecto |
| **project-context.yaml** | Contexto específico de TramaTex |
| **load-session.yaml** | Gestión de sprints y tareas |

## 🧩 Contextos Modulares (`agents/project/context/`)

Carga bajo demanda según necesidad:

- **architecture.yaml** - Estructura Clean Architecture y capas
- **bounded-contexts.yaml** - Módulos y dominios del negocio
- **code-standards.yaml** - Estándares de calidad y testing
- **tech-stack.yaml** - Tecnologías y herramientas

## 📖 Más Información

- **Uso de contextos:** `agents/project/context/README.md`
- **Desarrollar agentes:** `agents/guia-desarrollo-agents.md`
- **Referencia completa:** `AGENTS.md`

---

**Tip:** No necesitas gestionar estos archivos directamente. El asistente los usa automáticamente.