# 📂 Detalles de la Estructura del Proyecto

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Vigente |
| **Relacionado con** | [ADR-009](../../architecture/adrs/adr-009-project-structure.md) |

---

## 🎯 Propósito
Este documento proporciona una vista exhaustiva de la organización del repositorio TramaTex, facilitando el descubrimiento de componentes y la comprensión de la jerarquía modular tanto para desarrolladores como para asistentes de IA.

---

## 🏗️ Estructura Completa del Proyecto

```text
tramatex/
│
├── README.md                          # Documentación principal (Tronco del Árbol)
├── LICENSE.md                         # Licencia MIT del software
├── Makefile                           # Orquestador de comandos técnicos (Go, Docker)
├── start-dev.ps1                      # Orquestador de arranque del entorno
│
├── apps/                              # APLICACIONES (Código Fuente)
│   ├── tramatex-api/                  # Backend API (Go)
│   │   ├── cmd/api/main.go            # Punto de entrada de la aplicación
│   │   ├── internal/                  # Dominios de Negocio (Monolito Modular)
│   │   │   ├── iam/                   # Bounded Context: Identidad y Acceso
│   │   │   ├── party/                 # Bounded Context: Clientes y Proveedores
│   │   │   ├── product/               # Bounded Context: Catálogo y Variantes
│   │   │   ├── pricing/               # Bounded Context: Motor de Precios
│   │   │   ├── sales/                 # Bounded Context: Flujo Comercial
│   │   │   └── mes/                   # Bounded Context: Control de Producción
│   │   └── pkg/                       # Librerías compartidas e infraestructura
│   │
│   └── frontend/                      # Frontend SPA (Vue.js 3 + TypeScript)
│       ├── src/
│       │   ├── components/            # Componentes reutilizables
│       │   ├── pages/                 # Vistas principales de negocio
│       │   ├── stores/                # Gestión de estado (Pinia)
│       │   └── design-system/         # Sistema de diseño (Vanilla CSS)
│
├── docs/                              # ÁRBOL DE CONOCIMIENTO (Castellano)
│   ├── architecture/                  # Registro de decisiones y visión técnica
│   │   ├── adrs/                      # Architectural Decision Records (ADRs 001-021)
│   │   ├── design-system/             # Documentación visual y estética
│   │   └── diagrams/                  # Diagramas C4 y flujos de dominio
│   ├── guides/                        # Manuales operativos y estándares
│   │   ├── developer/                 # Guías técnicas y operativa de scripts
│   │   └── user/                      # Manuales de uso funcional
│   ├── modules/                       # Especificaciones profundas por dominio
│   └── log/                           # Registro histórico (Session Logs y Sprints)
│
├── project-scaffolding/               # ECOSISTEMA DE ESTANDARIZACIÓN
│   ├── agents/                        # Motor de orquestación IA
│   └── templates/                     # Plantillas de arquitectura y documentos
│
├── docker/                            # INFRAESTRUCTURA (Docker Compose)
│   └── ...                            # Configuraciones de red, DB y caché
│
└── agents/                            # AGENTES OPERATIVOS
    └── project/                       # Contexto actual para asistentes de IA
```

---
[Volver al README Principal](../../../README.md)
