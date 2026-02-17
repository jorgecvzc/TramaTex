# Placeholders System for Project Scaffolding

Este documento describe el sistema de placeholders utilizado por `bootstrap.yaml` para personalizar los templates al crear nuevos proyectos.

## 📋 Variables del Sistema

Las siguientes variables son recolectadas por el sistema de bootstrapping durante el proceso de inicialización y utilizadas por el agente `placeholder-population.yaml` para personalizar los templates.

### Variables Críticas (Obligatorias)
*   **Fuente:** Extraídas de la documentación de entrada (`user-input-docs/`) o solicitadas al usuario.

| Variable | Descripción | Ejemplo | Uso en Templates |
|----------|-------------|---------|------------------|
| `PROJECT_NAME` | Nombre del proyecto | `mi-proyecto`, `TramaTex` | Reemplaza `$PROJECT_NAME` en todos los templates |
| `PROJECT_VISION` | Visión/propósito del proyecto (1-3 oraciones) | `"Sistema ERP modular para la industria textil"` | Reemplaza `$PROJECT_VISION` en templates de documentación |
| `COMPONENT_TYPES` | Tipos de componentes separados por coma | `backend,frontend` | Determina qué estructuras de directorios crear |

### Variables de Stack Tecnológico
*   **Fuente:** Extraídas de la documentación de entrada (`user-input-docs/`) o solicitadas al usuario.

| Variable | Descripción | Ejemplo | Uso en Templates |
|----------|-------------|---------|------------------|
| `DATABASE_CHOICE` | Motor de base de datos | `PostgreSQL`, `MySQL`, `MongoDB` | Reemplaza `$DATABASE_CHOICE` en config y docs |
| `BACKEND_LANGUAGE_FRAMEWORK` | Framework backend | `Go/Gin`, `Python/FastAPI`, `Node.js/Express` | Reemplaza `$BACKEND_LANGUAGE_FRAMEWORK` en templates |
| `FRONTEND_FRAMEWORK` | Framework frontend | `Vue/Vite`, `React/Next.js`, `Angular` | Reemplaza `$FRONTEND_FRAMEWORK` en templates |

### Variables de Bounded Contexts
*   **Fuente:** Extraídas de la documentación de entrada (`user-input-docs/`) o solicitadas al usuario.

| Variable | Descripción | Ejemplo | Uso en Templates |
|----------|-------------|---------|------------------|
| `BOUNDED_CONTEXTS` | Módulos/contextos del dominio | `Auth,Users,Products,Orders` | Reemplaza `$BOUNDED_CONTEXTS` en architecture.yaml |

### Variables de Nombres de Aplicaciones
*   **Fuente:** Inferidas o solicitadas al usuario por el agente `directory-creation.yaml`.

| Variable | Descripción | Default | Uso en Templates |
|----------|-------------|---------|------------------|
| `BACKEND_NAME` | Nombre del componente backend | `api` | Reemplaza `$BACKEND_NAME` en estructura de directorios |
| `FRONTEND_NAME` | Nombre del componente frontend | `web` | Reemplaza `$FRONTEND_NAME` en estructura de directorios |
| `INTERMEDIARY_NAME` | Nombre del componente intermediario | `bff` | Reemplaza `$INTERMEDIARY_NAME` en estructura de directorios |

### Variables de Comandos Make (Generadas Automáticamente)
*   **Fuente:** Inferidas por el agente `makefile-population.yaml` basado en el stack tecnológico seleccionado.

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `BACKEND_BUILD_CMD_VAR` | Comando para construir backend | `cd apps/api && go build` |
| `BACKEND_RUN_CMD_VAR` | Comando para ejecutar backend | `cd apps/api && go run cmd/api/main.go` |
| `BACKEND_TEST_CMD_VAR` | Comando para testear backend | `cd apps/api && go test ./...` |
| `BACKEND_LINT_CMD_VAR` | Comando para linter backend | `cd apps/api && golangci-lint run` |
| `FRONTEND_BUILD_CMD_VAR` | Comando para construir frontend | `cd apps/web && npm run build` |
| `FRONTEND_RUN_CMD_VAR` | Comando para ejecutar frontend | `cd apps/web && npm run dev` |
| `FRONTEND_TEST_CMD_VAR` | Comando para testear frontend | `cd apps/web && npm run test` |
| `FRONTEND_LINT_CMD_VAR` | Comando para linter frontend | `cd apps/web && npm run lint` |

---

## 📂 Templates que Requieren Reemplazo de Placeholders

### 1. Agentes de Contexto

#### `templates/agents/project/project-context-template.yaml`

**Placeholders a reemplazar:**
- `$PROJECT_NAME` → Nombre del proyecto
- `$PROJECT_VISION` → Visión del proyecto
- `TBD` → Dominio específico (ej: "E-commerce", "FinTech")
- `$DATABASE_CHOICE` → Motor de base de datos
- `$BOUNDED_CONTEXTS` → Lista de bounded contexts

**Ejemplo:**
```yaml
# Antes
project_identity:
  vision: "Your project's vision (1-2 sentences)."
  domain: "Your project's business domain (e.g., E-commerce, Healthcare, FinTech)."
  
# Después
project_identity:
  vision: "Sistema ERP modular para la industria textil"
  domain: "Manufacturing & ERP"
```

#### `templates/agents/load-session.yaml`

**Placeholders a reemplazar:**
- ❌ **NO HAY** - Este template ahora está completamente genérico y no requiere reemplazo de placeholders

---

### 2. Documentación

#### `templates/docs/adr-template.md`

**Placeholders a reemplazar:**
- `[[YYYY-MM-DD]]` → Fecha actual (generada automáticamente)
- `Bootstrap Process` → Nombre del creador del proyecto (si se recolecta)

#### ADRs Generados (`ADR-001-technology-stack-selection.md`)

**Placeholders a reemplazar:**
- `$PROJECT_NAME` → Nombre del proyecto
- `$PROJECT_VISION` → Visión del proyecto
- `$BACKEND_LANGUAGE_FRAMEWORK` → Framework backend seleccionado
- `$FRONTEND_FRAMEWORK` → Framework frontend seleccionado
- `$DATABASE_CHOICE` → Base de datos seleccionada
- `[[YYYY-MM-DD]]` → Fecha de creación

---

### 3. Archivos Raíz

#### `README.md`

**Placeholders a reemplazar:**
- `$PROJECT_NAME` → Nombre del proyecto (aparece múltiples veces)
- `$PROJECT_VISION` → Visión del proyecto

#### `AGENTS.md`

**Placeholders a reemplazar:**
- `$PROJECT_NAME` → Nombre del proyecto (aparece múltiples veces)
- `$PROJECT_VISION` → Visión del proyecto

#### `.env.example`

**Placeholders a reemplazar:**
- `$PROJECT_NAME` → Nombre del proyecto

#### `Makefile`

**Placeholders a reemplazar:**
- `$PROJECT_NAME` → Nombre del proyecto
- `$(BACKEND_NAME)` → Nombre del componente backend
- `$(FRONTEND_NAME)` → Nombre del componente frontend
- `$(INTERMEDIARY_NAME)` → Nombre del componente intermediario
- `$BACKEND_BUILD_CMD_VAR` → Comando build backend
- `$BACKEND_RUN_CMD_VAR` → Comando run backend
- `$BACKEND_TEST_CMD_VAR` → Comando test backend
- `$BACKEND_LINT_CMD_VAR` → Comando lint backend
- `$FRONTEND_BUILD_CMD_VAR` → Comando build frontend
- `$FRONTEND_RUN_CMD_VAR` → Comando run frontend
- `$FRONTEND_TEST_CMD_VAR` → Comando test frontend
- `$FRONTEND_LINT_CMD_VAR` → Comando lint frontend
- `$INTERMEDIARY_BUILD_CMD_VAR` → Comando build intermediary
- `$INTERMEDIARY_RUN_CMD_VAR` → Comando run intermediary
- `$INTERMEDIARY_TEST_CMD_VAR` → Comando test intermediary
- `$INTERMEDIARY_LINT_CMD_VAR` → Comando lint intermediary

---

### 4. Estado del Proyecto

#### `docs/log/project-status.md`

**Placeholders a reemplazar:**
- `$PROJECT_NAME` → Nombre del proyecto (título)
- `$PROJECT_VISION` → Visión del proyecto
- `[[YYYY-MM-DD]]` → Fecha de creación

---

## 🔧 Implementación del Sistema de Placeholders

El sistema de placeholders es gestionado por el agente `placeholder-population.yaml` en `agents/bootstrap_workflow/`. Este agente contiene un `python_script` que se encarga de identificar automáticamente los archivos que contienen placeholders, aplicar los reemplazos con los valores del proyecto y generar un reporte de los archivos procesados.

**Ubicación:** `agents/bootstrap_workflow/placeholder-population.yaml`

**Función:**
1.  Recolecta todas las variables del sistema (extraídas de documentos, entrada de usuario, comandos generados).
2.  Define un diccionario de reemplazos completo.
3.  Itera sobre una lista predefinida de archivos (`files_to_process`).
4.  Para cada archivo, realiza todos los reemplazos de placeholders.
5.  Reporta los archivos procesados, omitidos y errores.

**Variables disponibles para el script:** Todas las variables del sistema (`$PROJECT_NAME`, `$BACKEND_BUILD_CMD_VAR`, etc.) se pasan al script de Python.

---

## 💡 Sistema de Reemplazo Unificado Actual

El sistema de reemplazo de placeholders implementado en `agents/bootstrap_workflow/placeholder-population.yaml` se encarga de aplicar los valores finales del proyecto a todos los templates generados.

### Funcionamiento:

1.  **Centralización:** El agente `placeholder-population.yaml` contiene un `python_script` que consolida todas las variables del sistema (extraídas de documentos, entrada de usuario, comandos generados).
2.  **Identificación de Archivos:** El script tiene una lista predefinida de `files_to_process` (archivos generados que se sabe que contienen placeholders).
3.  **Aplicación de Reemplazos:** Itera sobre estos archivos, buscando y reemplazando todos los placeholders (`$VARIABLE`, `[[YYYY-MM-DD]]`) con los valores correspondientes.
4.  **Reporte:** Genera un resumen de los archivos procesados, omitidos o con errores.

---

## ✅ Checklist de Validación

Después de ejecutar `bootstrap.yaml`, verificar que:

- [ ] `agents/project/project-context.yaml` contiene el nombre y visión del proyecto
- [ ] `README.md` tiene el título correcto y la visión
- [ ] `AGENTS.md` menciona el proyecto correcto
- [ ] `docs/log/project-status.md` está inicializado con el nombre del proyecto
- [ ] ADRs generados contienen stack tecnológico correcto
- [ ] `Makefile` tiene comandos apropiados para el stack seleccionado
- [ ] `.env.example` tiene el nombre del proyecto
- [ ] No quedan placeholders sin reemplazar (buscar patrones `$[A-Z_]+` o `\[\[.*?\]\]`)

---

## 🔗 Referencias

- **bootstrap.yaml:** [../bootstrap.yaml](../bootstrap.yaml)
- **Template Project Context:** [../templates/agents/project/project-context-template.yaml](../templates/agents/project/project-context-template.yaml)
- **Template Load Session:** [../templates/agents/load-session.yaml](../templates/agents/load-session.yaml)

---

**Última Actualización:** [[YYYY-MM-DD]]
