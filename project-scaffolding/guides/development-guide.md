# Guía de Desarrollo - Modificar el Sistema de Scaffolding

## 🎯 Propósito
Esta guía explica cómo modificar y mejorar el sistema de scaffolding de proyectos.

## 📋 Reglas Fundamentales

### 1. Sin Historial en Archivos
- ❌ NO usar versiones en nombres de archivos (`v2.0`, `_v1`, etc.)
- ❌ NO mantener CHANGELOG ni historial de cambios
- ✅ Siempre se ejecuta la última versión
- ✅ Versionado solo en documentación como referencia informativa

### 2. Documentación Concisa (Contextual)

La concisión es un principio clave, pero su aplicación es contextual:

-   **Para Guías de Usuario (`user-guide.md`):** Aplicar estrictamente. Máximo 40 líneas para un entendimiento rápido.
-   **Para Guías de Desarrollo (como esta):** El objetivo es ser exhaustivo pero claro. Priorizar la estructura (listas, tablas, bloques de código) sobre párrafos extensos. Se permite más detalle cuando sea necesario para la comprensión técnica.
-   **Para Documentación de Referencia (`placeholders-guide.md`):** Exhaustividad es clave. La concisión se refiere a ir al grano en cada punto, no a la longitud total del documento.

### 3. Estructura de Archivos
```
project-scaffolding/
├── bootstrap.yaml                      ← Orquestador principal del flujo
├── guides/                             ← Guías y documentación del scaffolding
│   ├── user-guide.md                   ← Para usuarios (cómo usar el scaffolding)
│   ├── development-guide.md            ← Este archivo (modificar scaffolding)
│   ├── placeholders-guide.md           ← Documentación del sistema de placeholders
│   └── input-docs/                     ← Guías para preparar documentos de entrada
│       ├── README.md                   ← Guía del directorio
│       └── EJEMPLO_README.md           ← Template de referencia
├── user-input-docs/                    ← Directorio de entrada de documentos
├── agents/                             ← Agentes del scaffolding
│   ├── bootstrap_workflow/             ← Agentes modulares del flujo de bootstrapping
│   │   ├── application-scaffolding.yaml
│   │   ├── adr-generation.yaml
│   │   ├── directory-creation.yaml
│   │   ├── finalization.yaml
│   │   ├── git-initialization.yaml
│   │   ├── input-processing.yaml
│   │   ├── makefile-population.yaml
│   │   ├── modular-context-setup.yaml
│   │   ├── placeholder-population.yaml
│   │   ├── root-files-creation.yaml
│   │   └── template-copying.yaml
│   ├── close-session.yaml              ← Agente para cerrar sesiones
│   ├── generic-rules.yaml              ← Reglas universales
│   ├── init.yaml                       ← Punto de entrada principal
│   ├── load-project-context.yaml       ← Carga contexto del proyecto
│   ├── load-session.yaml               ← Gestiona sesiones de desarrollo
│   └── project/                        ← Agentes específicos del proyecto
│       ├── project-context.yaml        ← Contexto principal del proyecto
│       ├── sprint-registry.yaml        ← Registro de sprints
│       └── context/                    ← Agentes de contexto modular
│           ├── architecture.yaml
│           ├── bounded-contexts.yaml
│           ├── code-standards.yaml
│           ├── README.md
│           └── tech-stack.yaml
├── tmp/                                ← Archivos temporales (limpiar al terminar)
│   └── README.md                       ← Reglas del directorio temporal
└── templates/                          ← Templates para nuevos proyectos
    ├── agents/                         ← Templates de agentes IA
    ├── docs/                           ← Templates de documentación
    └── .github/                        ← Templates de CI/CD workflows
```

## 🔧 Modificar el Flujo de Bootstrapping

### Archivo Principal: `bootstrap.yaml` (Orquestador)

`bootstrap.yaml` ha sido refactorizado para actuar como un **orquestador del flujo de bootstrapping**. Esto significa que su `workflow_steps` ahora consiste principalmente en cargar (load) agentes especializados de `agents/bootstrap_workflow/` en una secuencia lógica.

**Estructura Típica de `workflow_steps` en `bootstrap.yaml`:**
```yaml
workflow_steps:
  - step_id: "load_agente_x"
    title: "Cargar Agente X"
    description: "Carga el agente responsable de la Fase X del proceso."
    load: "agents/bootstrap_workflow/agente-x.yaml"
    next_step: "load_agente_y"

  - step_id: "load_agente_y"
    title: "Cargar Agente Y"
    description: "Carga el agente responsable de la Fase Y del proceso."
    load: "agents/bootstrap_workflow/agente-y.yaml"
    next_step: "load_agente_z"
  
  # ... y así sucesivamente
```

**Para Modificar el Flujo de Bootstrapping:**

1.  **Identifica la Fase:** Determina qué parte del proceso de bootstrapping necesitas modificar (ej. extracción de entrada, creación de directorios, generación de ADRs).
2.  **Localiza el Agente Responsable:** Dirígete al agente YAML correspondiente en `agents/bootstrap_workflow/` (ej. `input-processing.yaml` para modificar la extracción de entrada).
3.  **Modifica el Agente Específico:** Realiza los cambios dentro de ese archivo YAML.
4.  **Ajusta la Orquestación (si es necesario):** Si añades un nuevo agente o cambias el orden, actualiza la secuencia de `load` y `next_step` en `bootstrap.yaml`.

**Estructura de `metadata` en `bootstrap.yaml` (sin cambios):**
```yaml
metadata:
  filename: "bootstrap.yaml"
  version: "X.X"           # Solo informativo
  decision_log:            # Decisiones arquitectónicas importantes
    - date: "YYYY-MM-DD"
      decision: "..."
      rationale: "..."
```

### Fases del Flujo Actual (Modular)

El flujo de bootstrapping ha sido modularizado para mejorar la claridad y mantenibilidad. `bootstrap.yaml` ahora actúa como un orquestador que carga agentes especializados en secuencia.

1.  **Agente de Procesamiento de Entrada (`input-processing.yaml`):**
    *   **Responsabilidad:** Extrae metadata de documentos de entrada, interactúa con el usuario para obtener información faltante, y valida la configuración del proyecto.
    *   **Fase:** Obtención y Validación de Configuración.

2.  **Agente de Creación de Directorios (`directory-creation.yaml`):**
    *   **Responsabilidad:** Establece la estructura de directorios principal del proyecto, incluyendo `docs/`, `agents/`, `apps/`, y directorios específicos para cada tipo de aplicación seleccionada (backend, frontend, intermediary).
    *   **Fase:** Creación de Estructura Base.

3.  **Agente de Copia de Plantillas (`template-copying.yaml`):**
    *   **Responsabilidad:** Copia todos los agentes universales y plantillas de documentación (`.md`, `.yaml`, CI/CD) desde `project-scaffolding/templates/` a las ubicaciones adecuadas en el nuevo proyecto.
    *   **Fase:** Inyección de Plantillas.

4.  **Agente de Creación de Archivos Raíz (`root-files-creation.yaml`):**
    *   **Responsabilidad:** Crea archivos esenciales en la raíz del proyecto (`.gitignore`, `README.md`, `AGENTS.md`, `.env.example`, `Makefile`), siguiendo la "Política de Raíz Limpia".
    *   **Fase:** Creación de Archivos Raíz.

5.  **Agente de Población de Makefile (`makefile-population.yaml`):**
    *   **Responsabilidad:** Rellena el `Makefile` generado con comandos específicos del stack tecnológico seleccionado (Go, Node.js, Python, Vue, React, Angular).
    *   **Fase:** Configuración de Build/Test.

6.  **Agente de Población de Placeholders (`placeholder-population.yaml`):**
    *   **Responsabilidad:** Procesa todos los archivos generados que contienen placeholders (`$VARIABLE`, `[[YYYY-MM-DD]]`) y los reemplaza con los valores de configuración finales del proyecto.
    *   **Fase:** Personalización de Contenido.

7.  **Agente de Scaffolding de Aplicaciones (`application-scaffolding.yaml`):**
    *   **Responsabilidad:** Basado en la configuración del usuario, crea estructuras de directorios específicas para cada aplicación (backend, frontend, intermediary) dentro de `apps/` y configura sus respectivos workflows de CI/CD en `.github/workflows/`.
    *   **Fase:** Scaffolding de Componentes.

8.  **Agente de Inicialización de Git (`git-initialization.yaml`):**
    *   **Responsabilidad:** Realiza el `git add .`, el commit inicial de la estructura del proyecto y establece la rama `main`.
    *   **Fase:** Control de Versiones Inicial.

9.  **Agente de Generación de ADRs (`adr-generation.yaml`):**
    *   **Responsabilidad:** Genera ADRs fundacionales para el stack tecnológico, seguridad y estrategia de calidad/testing, y solicita al usuario su revisión.
    *   **Fase:** Establecimiento de Arquitectura.

10. **Agente de Configuración de Contexto Modular (`modular-context-setup.yaml`):**
    *   **Responsabilidad:** Genera y puebla los agentes de contexto modular (`architecture.yaml`, `bounded-contexts.yaml`, `tech-stack.yaml`, `code-standards.yaml`) a partir de la documentación y configuración del proyecto.
    *   **Fase:** Configuración de Contexto IA.

11. **Agente de Finalización (`finalization.yaml`):**
    *   **Responsabilidad:** Muestra la política de estructura del proyecto y un mensaje de finalización, proporcionando al usuario los siguientes pasos claros.
    *   **Fase:** Conclusión.

### Añadir un Nuevo Paso

```yaml
- step_id: "mi_nuevo_paso"
  title: "🎯 Título del Paso"
  order: X
  description: "Qué hace este paso"
  prompt_user:
    - question: "¿Pregunta al usuario?"
      variable_name: "MI_VARIABLE"
  next_step: "siguiente_paso"
```

### Variables del Sistema

El sistema utiliza placeholders con el formato `$VARIABLE_NAME` para inyectar valores específicos del proyecto en los templates. Estas variables se consolidan a partir de la extracción de documentos y/o la entrada manual del usuario.

**Variables extraídas de documentos (con sufijo `_EXTRACTED`):**
- `$PROJECT_NAME_EXTRACTED` - Nombre del proyecto extraído
- `$PROJECT_VISION_EXTRACTED` - Visión del proyecto extraída
- `$COMPONENT_TYPES_EXTRACTED` - Tipos de componentes extraídos
- `$DATABASE_CHOICE_EXTRACTED` - Base de datos extraída
- `$BACKEND_LANGUAGE_FRAMEWORK_EXTRACTED` - Stack backend extraído
- `$FRONTEND_FRAMEWORK_EXTRACTED` - Stack frontend extraído
- `$BOUNDED_CONTEXTS_EXTRACTED` - Bounded contexts extraídos
- `$CONFIDENCE_SCORE` - Nivel de confianza de la extracción
- `$SOURCE_FILE` - Archivo fuente de la extracción

**Variables finales (disponibles para templates después de consolidación):**
- `$PROJECT_NAME`
- `$PROJECT_VISION`
- `$COMPONENT_TYPES`
- `$DATABASE_CHOICE`
- `$BACKEND_LANGUAGE_FRAMEWORK`
- `$FRONTEND_FRAMEWORK`
- `$BOUNDED_CONTEXTS`
- `$BACKEND_NAME`
- `$FRONTEND_NAME`
- `$INTERMEDIARY_NAME`
- `$USER_NAME` (si se obtiene del usuario o contexto)
- `$CURRENT_DATE` (para fechas como `[[YYYY-MM-DD]]`)

**Variables de Comandos Make (generadas por `makefile-population.yaml`):**
- `$BACKEND_BUILD_CMD_VAR`
- `$BACKEND_RUN_CMD_VAR`
- `$BACKEND_TEST_CMD_VAR`
- `$BACKEND_LINT_CMD_VAR`
- `$FRONTEND_BUILD_CMD_VAR`
- `$FRONTEND_RUN_CMD_VAR`
- `$FRONTEND_TEST_CMD_VAR`
- `$FRONTEND_LINT_CMD_VAR`
- `$INTERMEDIARY_BUILD_CMD_VAR`
- `$INTERMEDIARY_RUN_CMD_VAR`
- `$INTERMEDIARY_TEST_CMD_VAR`
- `$INTERMEDIARY_LINT_CMD_VAR`

## 🧠 Sistema de Extracción Inteligente

### Ubicación
Paso `process_input_documents_early` en `agents/bootstrap_workflow/input-processing.yaml`

### Estrategia
1. **Lee TODOS** los archivos en `user-input-docs/`
2. **Extrae** usando múltiples patrones:
   - Búsqueda explícita ("Nombre:", "Visión:")
   - Inferencia contextual ("API REST" → backend detected)
   - Análisis estructural (H1 = nombre probable)
3. **Asigna confianza** (alta/media/baja)
4. **Detecta conflictos** entre documentos
5. **Calcula completitud** (%)

### Añadir Nuevo Campo a Extraer

En `ai_extraction_instructions` del paso `process_input_documents_early`:

```yaml
X️⃣  NUEVO_CAMPO (VARIABLE_NAME)
═══════════════════════════════════════
Buscar en:
• Secciones: "Keyword1", "Keyword2"
• Líneas que contengan: "pattern:"
• Campo "field:" en YAML/JSON

Formato esperado: tipo de dato
```

## 📦 Modificar Templates

### Ubicación: `templates/`

- **`agents/`** - Templates de agentes IA
- **`docs/`** - Templates de documentación
- **`.github/`** - Templates de CI/CD

### Añadir Nuevo Template

1.  **Crear el Archivo Template:** Crea tu nuevo archivo template en el subdirectorio adecuado dentro de `templates/` (ej. `templates/docs/mi-nuevo-template.md`).
    *   Usa variables del sistema (placeholders) como `$PROJECT_NAME`, `$PROJECT_VISION`, `[[YYYY-MM-DD]]` para contenido dinámico.
2.  **Referenciar en el Agente de Copia:** Abre `agents/bootstrap_workflow/template-copying.yaml` y añade una nueva entrada bajo `files_to_copy_relative_to_project_scaffolding` para tu nuevo template.
    *   **source:** Ruta relativa al template dentro de `project-scaffolding/`.
    *   **destination:** Ruta relativa donde se copiará en el nuevo proyecto.

**Ejemplo:**
```yaml
# agents/bootstrap_workflow/template-copying.yaml
...
files_to_copy_relative_to_project_scaffolding:
  - source: "templates/docs/mi-nuevo-template.md"
    destination: "docs/mis-documentos/mi-nuevo-documento.md"
```

## ✅ Checklist al Modificar

- [ ] Modificación reflejada en `bootstrap.yaml`
- [ ] Variables documentadas si son nuevas
- [ ] Decision log actualizado si es cambio arquitectónico
- [ ] Documentación actualizada (user-guide.md si afecta usuarios)
- [ ] Sin versiones en nombres de archivos
- [ ] Documentación concisa (máx 40 líneas cuando sea posible)

## 🔍 Testing

### Probar un Cambio

1. Coloca documentos de prueba en `user-input-docs/`
2. Ejecuta el bootstrap
3. Verifica que el cambio funciona correctamente
4. Limpia archivos generados

### Casos de Prueba Recomendados

- **Caso 1:** Documentación completa (README + specs)
- **Caso 2:** Solo README básico
- **Caso 3:** Sin documentación

## 📝 Mejores Prácticas

1. **Mantén el flujo simple** - Cada paso debe tener un propósito claro
2. **Extrae antes de preguntar** - Aprovecha documentación existente
3. **Muestra lo extraído** - Transparencia con el usuario
4. **Pregunta solo lo necesario** - Evita preguntas redundantes
5. **Valida antes de crear** - Confirma configuración final
6. **Documenta decisiones** - Usa decision_log para cambios importantes

### Mantenimiento de Agentes de Contexto

El directorio `docs/` es la **única fuente de verdad** para la mayoría de los detalles arquitectónicos y de módulos del proyecto (ej: ADRs, especificaciones de módulos). Los agentes de contexto en `agents/project/context/*.yaml` se **generan o actualizan** automáticamente a partir de esta documentación.

**Cómo funciona:**
1- El script `agents/scripts/doc_to_agent_context.py` lee los archivos relevantes en `docs/` (como ADRs y especificaciones de módulos).
2- Extrae información estructurada (ej: descripción de Bounded Contexts, capas de arquitectura).
3- Actualiza los archivos YAML correspondientes en `agents/project/context/`.

**Integración en el flujo:**
- **Creación de proyectos:** Este script se ejecuta automáticamente durante el proceso de `bootstrap.yaml` para poblar los agentes de contexto iniciales del nuevo proyecto.
- **Proyectos existentes:** Los agentes de contexto pueden re-sincronizarse manualmente ejecutando `python agents/scripts/doc_to_agent_context.py` o seleccionando la opción correspondiente en `agents/init.yaml`.

**Guía de Mantenimiento:**
- Cuando se actualice o se cree un nuevo ADR relevante o una especificación de módulo en `docs/`, **debe ejecutarse** el script de sincronización para mantener los agentes de contexto actualizados.
- Si necesitas modificar cómo se extrae o se mapea la información, actualiza `scripts/doc_to_agent_context.py`.
- Si necesitas añadir un nuevo campo a un agente de contexto a partir de la documentación, extiende el script `scripts/doc_to_agent_context.py`.

## 🆘 Resolver Problemas Comunes

**Error en extracción de documentos:**
- Verificar patrones en `ai_extraction_instructions`
- Revisar formatos soportados
- Ajustar nivel de confianza

**Flujo se salta pasos:**
- Revisar `next_step` en cada paso
- Verificar condiciones `if` en pasos condicionales

**Variables no disponibles:**
- Verificar que se generan en paso previo
- Revisar `merge_and_finalize_metadata`

## 📚 Referencias

- **Archivo principal:** `bootstrap.yaml`
- **Guía de usuario:** `guides/user-guide.md`
- **Templates:** `templates/`
- **Ejemplo de docs:** `guides/input-docs/EJEMPLO_README.md`
