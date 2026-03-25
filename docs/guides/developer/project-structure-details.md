# Detalles de la Estructura del Proyecto TramaTex

Este documento proporciona una vista detallada de la estructura de carpetas del proyecto TramaTex, complementando la decisión arquitectónica definida en [ADR-009: Estructura de Carpetas y Organización del Proyecto](../../architecture/adrs/adr-009-project-structure.md).

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
│   │   ├── adrs/                      # Architecture Decision Records
│   │   │   └── ...
│   │   ├── design-system/             # Principios de diseño y UI
│   │   └── diagrams/                  # Diagramas (C4, ER, etc.)
│   │
│   ├── guides/                        # Guías y tutoriales (desarrolladores, usuarios)
│   │   ├── developer/
│   │   └── user/
│   │
│   ├── modules/                       # Documentación detallada por Bounded Context (specs, diagramas)
│   │   ├── _module-template.md
│   │   ├── iam/
│   │   ├── party/
│   │   └── ...
│   │
│   └── log/                           # Registros del proyecto (sprints, hitos, gobernanza)
│       ├── erp-core-completion.md     # Informe de completitud del ERP Core
│       ├── mes-completion.md          # Informe de completitud del módulo MES
│       ├── project-status.md          # Estado actual del proyecto
│       ├── session-log.md             # Registro de sesiones de agentes
│       ├── analysis/                  # Análisis técnicos y gap reports
│       ├── governance/                # Políticas del proyecto
│       └── sprints/                  # Planificación y logs de Sprints y Tareas
│           ├── _sprint-summary-template.md
│           ├── _task-template.md
│           └── sprint-01/
│               └── ...
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
│   │   │   ├── sales/                 # Bounded Context: Sales (Ventas)
│   │   │   │   ├── application/       # Servicios fragmentados (Quote, Order, Billing)
│   │   │   │   ├── domain/            # Modelo de dominio (Quote, SalesOrder, calculations)
│   │   │   │   ├── infrastructure/    # Persistencia (GORM) y adaptadores externos
│   │   │   │   └── interfaces/        # Handlers HTTP y DTOs
│   │   │   ├── mes/                   # Bounded Context: MES (Producción)
│   │   │   └── shared/                # Código compartido entre Bounded Contexts
│   │   ├── pkg/
│   │   └── ...
│   │
│   └── frontend/                      # Frontend (Vue.js 3)
│       ├── src/
│       │   ├── components/
│       │   ├── pages/                 # Páginas principales (anteriormente views)
│       │   ├── stores/
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
