# Detalles de la Estructura del Proyecto TramaTex

Este documento proporciona una vista detallada de la estructura de carpetas del proyecto TramaTex, complementando la decisión arquitectónica definida en [ADR-009: Estructura de Carpetas y Organización del Proyecto](../../architecture/adrs/ADR-009-project-structure.md).

---

## Estructura Completa del Proyecto

```
tramatex/
│
├── README.md                          # Documentación principal del proyecto
├── LICENSE                            # Licencia del software
├── .gitignore                         # Archivos ignorados por Git
├── Makefile                           # Comandos comunes (build, test, run, etc.)
│
├── docs/                              # DOCUMENTACIÓN (en Español)
│   ├── architecture/                  # Arquitectura del sistema (ADRs, diagramas)
│   │   ├── adr/                       # Architecture Decision Records
│   │   │   └── ...
│   │   └── diagrams/                  # Diagramas (C4, ER, etc.)
│   │
│   ├── guides/                        # Guías y tutoriales (desarrolladores, usuarios)
│   │   ├── developer/
│   │   └── user/
│   │
│   ├── modules/                       # Documentación detallada por Bounded Context (specs, diagramas)
│   │   ├── _MODULE_TEMPLATE.md
│   │   ├── iam/
│   │   ├── party/
│   │   └── ...
│   │
│   └── log/                           # Registros del proyecto (sprints, hitos, gobernanza)
│       ├── sprints/                  # Planificación y logs de Sprints y Tareas
│       │   ├── _SPRINT_TEMPLATE.md
│       │   ├── _TASK_TEMPLATE.md
│       │   └── sprint-01/
│       │       └── ...
│       ├── milestones/                # Hitos históricos y reportes de estado
│       │   └── ...
│       └── governance/                # Políticas del proyecto
│           └── ...
│
├── apps/
│   ├── tramatex-api/                  # Backend API (Go)
│   │   ├── cmd/api/main.go
│   │   ├── internal/
│   │   │   ├── iam/                   # Bounded Context: IAM
│   │   │   │   ├── domain/
│   │   │   │   ├── application/
│   │   │   │   ├── infrastructure/
│   │   │   │   └── interfaces/
│   │   │   ├── party/                 # Bounded Context: Party
│   │   │   │   └── ...
│   │   │   ├── product/
│   │   │   ├── pricing/
│   │   │   ├── sales/
│   │   │   ├── mes/
│   │   │   └── shared/                # Código compartido entre Bounded Contexts
│   │   ├── pkg/
│   │   └── ...
│   │
│   └── frontend/                      # Frontend (Vue.js 3)
│       ├── src/
│       │   ├── components/
│       │   ├── views/
│   │   │   ├── stores/
│       │   ├── services/
│       │   ├── composables/
│       │   └── ...
│       └── ...
│
├── docker/                            # Configuración de Docker
│   └── ...
│
└── .github/                           # CI/CD (GitHub Actions)
    └── ...
```
