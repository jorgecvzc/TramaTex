# Project Context Agents

This directory contains modular context agents that provide detailed information about the TramaTex project. These agents are loaded on-demand by AI assistants to understand specific aspects of the project.

## 📚 Available Context Agents

### Core Context Files

| File | Purpose | When to Load |
|------|---------|--------------|
| [architecture.yaml](./architecture.yaml) | Clean Architecture layers, DDD strategy, architectural patterns | When designing modules, understanding layer boundaries |
| [bounded-contexts.yaml](./bounded-contexts.yaml) | Business domains, module dependencies, relationships | When working across modules, understanding domain boundaries |
| [tech-stack.yaml](./tech-stack.yaml) | Technology choices, frameworks, tools, build systems | When implementing features, setting up environments |
| [code-standards.yaml](./code-standards.yaml) | Code quality criteria, acceptance/rejection rules | Before committing code, during code review |
| [load-module-context.yaml](../../load-module-context.yaml) | Loads context for a specific project module | When starting work on a specific module |

### Design Context

La documentación del sistema de diseño (paleta de colores, tipografía, etc.) no es un agente de contexto cargable, sino que se encuentra en el directorio de documentación principal. Consultar `docs/architecture/design-system/` para más detalles.

## 🎯 Loading Strategy

**Don't load everything at once.** Follow this pattern:

### For New Module Development
1. Load `architecture.yaml` - Understand layer structure
2. Load `bounded-contexts.yaml` - Understand dependencies
3. Load `code-standards.yaml` - Know quality gates
4. Load `load-module-context.yaml` (with module_name) - Get module-specific context

### For Feature Implementation
1. Load `bounded-contexts.yaml` - Find relevant module
2. Load `tech-stack.yaml` - Know tools and frameworks
3. Load `code-standards.yaml` - Ensure compliance
4. Load `load-module-context.yaml` (with module_name) - Get module-specific context

### For Code Review
1. Load `code-standards.yaml` - Review criteria
2. Load `architecture.yaml` - Verify layer compliance
3. Load `load-module-context.yaml` (with module_name) - Get module-specific context

## 🔗 Parent Context

These context agents are loaded via the main project context file:
- [../project-context.yaml](../project-context.yaml)

Which in turn is loaded after:
- [../../generic-rules.yaml](../../generic-rules.yaml)

## 📝 Maintenance

**When updating context agents:**
1. Update the `last_updated` field in metadata
2. Ensure consistency with corresponding documentation in `/docs`
3. Update this README if adding/removing files

---

**Last Updated:** 2026-02-15
