# ADR-009 – Estructura de Carpetas y Organización del Proyecto

**Fecha:** 11/01/2026  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, Claude (Anthropic)

---

## 1. Contexto

TramaTex requiere una estructura de proyecto que:

- Refleje los principios de **Clean Architecture** (ADR-002)
- Soporte un **monolito modular** con dominios claramente separados (ADR-003)
- Facilite el **desarrollo dirigido por dominio** (ADR-006, ADR-007)
- Sea **mantenible por un desarrollador en solitario** con IA copiloto
- Permita **escalabilidad futura** hacia microservicios si es necesario
- Organice **documentación técnica y de sesiones** de forma eficiente

**Restricciones técnicas:**

- Stack: Go (tramatex-api), Vue.js 3 (frontend), PostgreSQL (persistencia), Docker Compose
- Desarrollo local-first
- Un solo repositorio Git (monorepo)
- Sin herramientas complejas de monorepo (Nx, Turborepo, etc.)

**Riesgos si no se define estructura:**

- Mezcla de capas (dominio con infraestructura)
- Dificultad para aplicar TDD
- Pérdida de trazabilidad entre módulos
- Código difícil de navegar y mantener
- Imposibilidad de extraer módulos a servicios independientes en el futuro

---

## 2. Alternativas Consideradas

### Alternativa A – Estructura plana por tipo de archivo

```
tramatex/
├── controllers/
├── models/
├── services/
├── repositories/
└── views/
```

**Ventajas:**

- Simple y tradicional
- Fácil de entender para principiantes
- Menos niveles de anidamiento

**Desventajas:**

- Mezcla todos los dominios
- Dificulta aplicar Clean Architecture
- No soporta modularidad
- Imposible extraer módulos
- No escala con la complejidad

### Alternativa B – Estructura por capas técnicas

```
tramatex/
├── domain/
├── application/
├── infrastructure/
└── interfaces/
```

**Ventajas:**

- Respeta capas de Clean Architecture
- Separación técnica clara

**Desventajas:**

- No refleja dominios de negocio
- Dificulta navegación por módulo
- Mezcla Party, Producto, MES en mismas carpetas
- No facilita extracción de módulos

### Alternativa C – Estructura por módulos de dominio (Clean Architecture modular)

```
tramatex/
├── apps/tramatex-api/
│   └── internal/
│       ├── domain/
│       │   ├── party/
│       │   ├── product/
│       │   ├── pricing/
│       │   └── sales/
│       ├── application/
│       │   ├── party/
│       │   └── product/
│       └── infrastructure/
│           └── persistence/
│               └── postgres/
├── apps/frontend/
│   └── src/
│       ├── views/
│       │   ├── party/
│       │   └── product/
│       └── stores/
└── docs/
    └── modules/
        ├── party/
        └── product/
```

**Ventajas:**

- Refleja dominios de negocio claramente
- Facilita navegación por módulo completo
- Soporta Clean Architecture
- Permite extraer módulos a servicios
- Agrupa código relacionado
- Facilita desarrollo incremental por módulo

**Desventajas:**

- Mayor profundidad de carpetas
- Requiere disciplina para mantener separación

### Alternativa D – Monorepo con workspaces independientes

```
tramatex/
├── packages/
│   ├── party-service/
│   ├── product-service/
│   └── shared/
└── apps/
    ├── apps/tramatex-api/
    └── apps/frontend/
```

**Ventajas:**

- Preparado para microservicios desde día 1
- Reutilización clara de código compartido

**Desventajas:**

- Complejidad prematura para MVP
- Requiere herramientas adicionales
- Overhead de gestión de versiones entre paquetes
- No se necesita para monolito modular

---

## 3. Decisión Adoptada

Se adopta **Alternativa C: Estructura por módulos de dominio con Clean Architecture**.

**Justificación:**

- **Refleja el dominio de negocio:** Party, Producto, Tarificación, Ventas, MES son módulos visibles
- **Soporta Clean Architecture:** Capas separadas dentro de cada módulo
- **Facilita TDD:** Tests están junto al código que prueban
- **Escalabilidad controlada:** Módulos pueden extraerse a servicios sin refactoring masivo
- **Navegación intuitiva:** Toda la lógica de Party está en carpetas `*/party/`
- **Desarrollo incremental:** Cada fase del ADR-007 tiene carpetas claras donde trabajar
- **Mantenibilidad:** Un desarrollador puede entender la estructura en una sesión
- **Documentación alineada:** La estructura de `/docs/modules/` refleja la del código

---

## 4. Consecuencias

### Positivas

- **Claridad de dominio:** Fácil encontrar todo lo relacionado con un módulo específico
- **Separación de responsabilidades:** Capas de Clean Architecture respetadas
- **Testabilidad:** Tests unitarios junto a entidades, tests de integración junto a casos de uso
- **Documentación trazable:** Estructura de `/docs/modules/` espeja la del código
- **Evolución controlada:** Nuevos módulos se añaden sin afectar existentes
- **Onboarding rápido:** Estructura autodocumentada
- **Preparación para microservicios:** Módulos ya están delimitados

### Negativas

- **Profundidad de carpetas:** Rutas pueden ser largas (ej: `internal/domain/party/customer.go`)
- **Duplicación aparente:** Carpetas `party/` existen en domain, application, infrastructure
- **Disciplina requerida:** Fácil romper separación si no hay code reviews
- **Curva de aprendizaje inicial:** Requiere entender Clean Architecture primero

### Mitigaciones

- **Profundidad:** IDEs modernos (VSCode, GoLand) manejan bien rutas largas
- **Duplicación:** Es intencional, cada capa tiene responsabilidades diferentes
- **Disciplina:** TDD obligatorio ayuda a mantener separación
- **Curva:** Documentación clara de la estructura en README.md

---

## 5. Alcance

Este ADR define:

- Estructura completa de carpetas del proyecto TramaTex
- Organización del monorepo (tramatex-api + frontend + documentación)
- Convenciones de nombres de carpetas y archivos
- Ubicación de tests, migraciones, configuraciones
- Estructura de documentación técnica

**No define:**

- Convenciones de nombres de variables/funciones (estilo de código)
- Estrategia de branching Git
- Proceso de CI/CD específico
- Configuración de herramientas (linters, formatters)

---

## 6. Integración con otros ADRs

- **ADR-001:** Stack tecnológico → Define Go, Vue.js, PostgreSQL (estructura los refleja)
- **ADR-002:** Clean Architecture → Estructura por capas (domain, application, infrastructure, interfaces)
- **ADR-003:** Monolito modular → Módulos separados pero en mismo repositorio
- **ADR-006:** Desarrollo dirigido por dominio → Estructura facilita priorizar por dominio
- **ADR-007:** Orden de implementación → Carpetas se crean incrementalmente por fase
- **ADR-008:** Cronograma → Estructura soporta desarrollo incremental de 24 meses

---

## 7. Notas Adicionales / Consideraciones Especiales

### Estructura Completa del Proyecto

```
tramatex/
│
├── README.md                          # Documentación principal del proyecto
├── LICENSE                            # Licencia del software
├── .gitignore                         # Archivos ignorados por Git
├── Makefile                           # Comandos comunes (build, test, run, etc.)
│
├── docs/                              # DOCUMENTACIÓN (en Español)
│   ├── overview/                      # Contexto y visión general del proyecto
│   ├── architecture/                  # Arquitectura del sistema (ADRs, diagramas)
│   │   ├── adr/                       # Architecture Decision Records
│   │   │   └── ...
│   │   └── diagrams/                  # Diagramas (C4, ER, etc.)
│   │
│   ├── guides/                        # Guías y tutoriales (desarrolladores, usuarios)
│   │   ├── developer/
│   │   └── user/
│   │
│   ├── reference/                     # Documentación de referencia (módulos, API)
│   │   ├── _MODULE_TEMPLATE.md
│   │   ├── iam/
│   │   ├── party/
│   │   └── ...
│   │
│   └── records/                       # Registros del proyecto (sprints, hitos, gobernanza)
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

---

### Convenciones de Nombres

**Regla Global de Nomenclatura:** Para mantener la consistencia y la compatibilidad con herramientas internacionales, se establece la siguiente política:
- **Nombres de Archivos y Carpetas:** **Inglés**, usando `kebab-case` para archivos de documentación (ej: `01-initial-architecture.md`) y el case apropiado para código fuente (ej: `userRepository.go`, `UserCard.vue`).
- **Contenido de Archivos:** **Español** para toda la documentación (`/docs`), **Inglés** para todo el código fuente, comentarios y mensajes de commit.

#### tramatex-api (Go)

**Archivos:**

- Snake case: `party_repository.go`, `pricing_engine.go`
- Tests: sufijo `_test.go`
- Interfaces: sufijo `_interface.go` solo si no es obvio (ej: `repository.go` es suficiente)

**Paquetes:**

- Todo en minúsculas, sin guiones bajos
- Singular preferido: `party`, `product`, no `parties`, `products`

**Entidades y Value Objects:**

- PascalCase: `Party`, `Customer`, `Money`

**Funciones:**

- PascalCase (exportadas): `CreateParty`, `CalculatePrice`
- camelCase (privadas): `validateNIF`, `applyDiscount`

#### Frontend (Vue.js)

**Archivos:**

- PascalCase para componentes: `PartyList.vue`, `ProductCard.vue`
- camelCase para servicios/stores: `partyService.js`, `auth.js`

**Componentes:**

- Siempre multi-palabra: `PartyList` (bien), `List` (mal)
- Prefijos descriptivos: `TheHeader`, `BaseButton`, `AppSidebar`

**Stores (Pinia):**

- Singular: `party.js`, no `parties.js`
- Export default: `usePartyStore`, `useProductStore`

#### Base de Datos (PostgreSQL)

**Tablas:**

- Snake case, plural: `parties`, `products`, `orders`

**Columnas:**

- Snake case: `party_id`, `created_at`, `supplier_cost`

**Índices:**

- Prefijo `idx_`: `idx_parties_nif`, `idx_orders_customer_id`

**Constraints:**

- Foreign keys: `fk_[tabla]_[columna]`
- Unique: `uk_[tabla]_[columna]`
- Check: `ck_[tabla]_[condicion]`

---

### Ubicación de Tests

#### Tests Unitarios

**Ubicación:** Mismo paquete que el código que prueban

```
domain/party/
├── party.go
└── party_test.go
```

#### Tests de Integración

**Ubicación:** Mismo paquete que el adaptador

```
infrastructure/persistence/postgres/
├── party_repository.go
└── party_repository_test.go
```

#### Tests End-to-End

**Ubicación:** Carpeta separada (opcional, Post-MVP)

```
apps/tramatex-api/
└── test/
    └── e2e/
        └── party_e2e_test.go
```

---

### Gestión de Configuraciones

**Variables de entorno:**

- Archivo `.env` en raíz (NO subir a Git, en `.gitignore`)
- Archivo `.env.example` con plantilla (SÍ subir a Git)

**Archivos de configuración:**

- `config/config.yaml` → Configuración por defecto
- `config/config.dev.yaml` → Sobreescribe valores para desarrollo
- `config/config.prod.yaml` → Sobreescribe valores para producción

**Carga de configuración:**

```go
// Carga primero config.yaml, luego sobreescribe con config.{env}.yaml
config.Load("config.yaml")
config.LoadEnv(os.Getenv("ENV")) // dev, prod
```

---

### Comandos Make (Makefile)

**tramatex-api:**

```makefile
# Desarrollo
make run              # Ejecutar aplicación
make test             # Ejecutar tests
make test-coverage    # Cobertura de tests
make lint             # Linter
make fmt              # Formatear código

# Base de datos
make migrate-up       # Aplicar migraciones
make migrate-down     # Revertir última migración
make seed             # Seed de datos

# Build
make build            # Compilar binario
make docker-build     # Build imagen Docker
make docker-up        # Levantar stack completo
make docker-down      # Detener stack
```

**Frontend:**

```makefile
make dev              # Servidor de desarrollo
make build            # Build para producción
make preview          # Preview del build
make lint             # ESLint
make format           # Prettier
```

---

### Flujo de Creación de Carpetas por Fase

#### Fase 0 (Fundaciones)

```
tramatex/
├── docs/
│   └── adr/
├── apps/tramatex-api/
│   ├── cmd/api/
│   ├── internal/
│   │   ├── infrastructure/security/  # JWT, RBAC
│   │   └── interfaces/http/
│   │       ├── middleware/
│   │       └── handlers/
│   └── config/
├── apps/frontend/
│   └── src/
│       ├── views/auth/
│       ├── stores/
│       └── services/
└── docker/
```

#### Fase 1 (Party + Producto + Tarificación)

```
# Se añaden:
apps/tramatex-api/internal/domain/
├── party/
├── product/
└── pricing/

apps/tramatex-api/internal/application/
├── party/
├── product/
└── pricing/

apps/tramatex-api/internal/infrastructure/persistence/postgres/
├── party_repository.go
├── product_repository.go
└── pricing_repository.go

apps/tramatex-api/internal/infrastructure/persistence/migrations/
├── 000002_create_parties.up.sql
├── 000003_create_products.up.sql
└── 000004_create_pricing.up.sql

apps/frontend/src/
├── views/party/
├── views/product/
├── views/pricing/
├── components/party/
├── components/product/
└── components/pricing/

docs/reference/
├── party/
├── product/
└── pricing/
```

#### Fase 2 (Pedidos)

```
# Se añaden:
apps/tramatex-api/internal/domain/sales/
apps/tramatex-api/internal/application/sales/
apps/tramatex-api/internal/infrastructure/persistence/postgres/order_repository.go
apps/frontend/src/views/sales/
apps/frontend/src/components/sales/
docs/modules/sales/
```

#### Fase 3 (MES)

```
# Se añaden:
apps/tramatex-api/internal/domain/mes/
apps/tramatex-api/internal/application/mes/
apps/tramatex-api/internal/infrastructure/
├── persistence/postgres/mes_repository.go
└── storage/nas/
apps/frontend/src/views/mes/
apps/frontend/src/components/mes/
docs/modules/mes/
```

---

## 8. Referencias

- **Documentos internos:**
    
    - Documento Consolidado TramaTex v3.0
    - ADR-001: Stack Tecnológico
    - ADR-002: Clean Architecture y DDD
    - ADR-003: Monolito Modular
    - ADR-007: Orden de Implementación de Módulos
- **Prácticas externas:**
    
    - Clean Architecture (Robert C. Martin)
    - Domain-Driven Design (Eric Evans)
    - Golang Project Layout: https://github.com/golang-standards/project-layout
    - Vue.js Style Guide: https://vuejs.org/style-guide/
- **Diagramas:**
    
    - Ver `/docs/diagrams/architecture/project-structure.png` (a crear)
    - Ver `/docs/diagrams/architecture/clean-architecture-layers.png` (a crear)

---

**Fin del ADR-009 – Estructura de Carpetas y Organización del Proyecto**
