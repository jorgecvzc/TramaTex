# Gestión de Configuraciones en el Proyecto TramaTex

Este documento describe cómo se gestionan las configuraciones en el proyecto TramaTex, incluyendo variables de entorno y archivos de configuración, complementando la decisión arquitectónica definida en [ADR-009: Estructura de Carpetas y Organización del Proyecto](../../architecture/adrs/ADR-009-project-structure.md).

---

## Gestión de Configuraciones

### Variables de Entorno

-   **Archivo `.env`**: Ubicado en la raíz del proyecto. Contiene variables de entorno sensibles o específicas del entorno local. **NO** debe subirse a Git (está en `.gitignore`).
-   **Archivo `.env.example`**: Un archivo de plantilla con ejemplos de las variables de entorno necesarias. **SÍ** debe subirse a Git para guiar a otros desarrolladores.

### Archivos de Configuración

-   `config/config.yaml`: Contiene la configuración por defecto de la aplicación.
-   `config/config.dev.yaml`: Sobreescribe los valores de `config.yaml` para el entorno de desarrollo.
-   `config/config.prod.yaml`: Sobreescribe los valores de `config.yaml` para el entorno de producción.

### Carga de Configuración (Ejemplo Go)

La lógica de carga de configuración debe seguir un patrón que priorice la flexibilidad y la sobrescritura por entorno.

```go
// Ejemplo de carga de configuración en Go
// Carga primero config.yaml, luego sobreescribe con config.{env}.yaml
config.Load("config.yaml")
config.LoadEnv(os.Getenv("ENV")) // Carga configuraciones específicas según la variable de entorno ENV (ej: dev, prod)
```
