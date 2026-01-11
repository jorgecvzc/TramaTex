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

- Stack: Go (backend), Vue.js 3 (frontend), PostgreSQL (persistencia), Docker Compose
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
├── backend/
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
├── frontend/
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
    ├── backend/
    └── frontend/
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
- Organización del monorepo (backend + frontend + documentación)
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
├── docs/                              # DOCUMENTACIÓN
│   ├── adr/                          # Architecture Decision Records
│   │   ├── ADR-001-stack-tecnologico.md
│   │   ├── ADR-002-clean-architecture-ddd.md
│   │   ├── ADR-003-monolito-modular.md
│   │   ├── ADR-004-ciclo-vida-mvp.md
│   │   ├── ADR-005-party-organizacion.md
│   │   ├── ADR-006-desarrollo-dirigido-dominio.md
│   │   ├── ADR-007-orden-implementacion.md
│   │   ├── ADR-008-planificacion-cronograma.md
│   │   └── ADR-009-estructura-proyecto.md
│   │
│   ├── modules/                      # Documentación de módulos
│   │   ├── _TEMPLATE.md             # Plantilla base para documentar módulos
│   │   ├── party/
│   │   │   ├── module-spec.md       # Especificación del módulo
│   │   │   ├── domain-model.md      # Modelo de dominio detallado
│   │   │   ├── use-cases.md         # Casos de uso documentados
│   │   │   └── api-contracts.md     # Contratos de API REST
│   │   ├── product/
│   │   │   ├── module-spec.md
│   │   │   ├── domain-model.md
│   │   │   ├── use-cases.md
│   │   │   └── api-contracts.md
│   │   ├── pricing/
│   │   │   └── ... (misma estructura)
│   │   ├── sales/
│   │   │   └── ... (misma estructura)
│   │   └── mes/
│   │       └── ... (misma estructura)
│   │
│   ├── diagrams/                     # Diagramas técnicos
│   │   ├── architecture/
│   │   │   ├── clean-architecture.png
│   │   │   ├── bounded-contexts.png
│   │   │   └── deployment.png
│   │   ├── domain/
│   │   │   ├── party-er-diagram.png
│   │   │   ├── product-er-diagram.png
│   │   │   ├── pricing-flow.png
│   │   │   └── mes-states.png
│   │   └── flows/
│   │       ├── order-standard-flow.png
│   │       ├── order-custom-flow.png
│   │       └── pricing-calculation-flow.png
│   │
│   ├── sessions/                     # Documentación de sesiones de desarrollo
│   │   ├── _SESSION_TEMPLATE.md    # Plantilla para documentar sesiones
│   │   ├── 2026-01-11-session-01.md # Sesión de planificación
│   │   ├── 2026-01-18-session-02.md
│   │   └── ... (una por semana de desarrollo)
│   │
│   ├── consolidated/                 # Documentos consolidados
│   │   └── tramatex-consolidated-v3.0.md
│   │
│   └── guides/                       # Guías técnicas
│       ├── setup-development.md     # Cómo configurar entorno de desarrollo
│       ├── testing-strategy.md      # Estrategia de testing
│       ├── deployment.md            # Guía de despliegue
│       └── contribution.md          # Guía de contribución (si hay equipo)
│
├── backend/                          # BACKEND (Go)
│   ├── cmd/                         # Puntos de entrada de la aplicación
│   │   └── api/
│   │       └── main.go              # Punto de entrada principal
│   │
│   ├── internal/                    # Código privado de la aplicación
│   │   │
│   │   ├── domain/                  # CAPA DE DOMINIO (DDD)
│   │   │   │                        # Lógica de negocio pura, sin dependencias externas
│   │   │   │
│   │   │   ├── party/               # Bounded Context: Party / Organización
│   │   │   │   ├── party.go         # Entidad raíz Party
│   │   │   │   ├── party_test.go    # Tests unitarios Party
│   │   │   │   ├── party_role.go    # Value Object PartyRole
│   │   │   │   ├── customer.go      # Entidad Customer (especialización)
│   │   │   │   ├── customer_test.go
│   │   │   │   ├── supplier.go      # Entidad Supplier (especialización)
│   │   │   │   ├── supplier_test.go
│   │   │   │   ├── supplier_cost.go # Entidad SupplierCost
│   │   │   │   ├── supplier_cost_test.go
│   │   │   │   └── repository.go    # Interface del repositorio (puerto)
│   │   │   │
│   │   │   ├── product/             # Bounded Context: Producto
│   │   │   │   ├── product.go       # Entidad raíz Product
│   │   │   │   ├── product_test.go
│   │   │   │   ├── variant.go       # Entidad Variant
│   │   │   │   ├── variant_test.go
│   │   │   │   ├── category.go      # Entidad Category
│   │   │   │   ├── category_test.go
│   │   │   │   └── repository.go    # Interface del repositorio
│   │   │   │
│   │   │   ├── pricing/             # Bounded Context: Tarificación
│   │   │   │   ├── pricing_engine.go      # Motor de cálculo de precios
│   │   │   │   ├── pricing_engine_test.go # Tests (cobertura ≥80%)
│   │   │   │   ├── pricing_rule.go        # Entidad PricingRule
│   │   │   │   ├── pricing_rule_test.go
│   │   │   │   ├── price_calculator.go    # Servicio de dominio
│   │   │   │   ├── price_calculator_test.go
│   │   │   │   └── repository.go          # Interface del repositorio
│   │   │   │
│   │   │   ├── sales/               # Bounded Context: Ventas
│   │   │   │   ├── order.go         # Agregado raíz Order
│   │   │   │   ├── order_test.go
│   │   │   │   ├── order_line.go    # Entidad OrderLine
│   │   │   │   ├── order_line_test.go
│   │   │   │   ├── order_state.go   # Value Object OrderState
│   │   │   │   ├── order_state_test.go
│   │   │   │   └── repository.go    # Interface del repositorio
│   │   │   │
│   │   │   ├── mes/                 # Bounded Context: MES (Producción)
│   │   │   │   ├── custom_order.go       # Agregado raíz CustomOrder
│   │   │   │   ├── custom_order_test.go
│   │   │   │   ├── production_state.go   # Value Object ProductionState
│   │   │   │   ├── production_state_test.go
│   │   │   │   ├── workshop_job.go       # Entidad WorkshopJob
│   │   │   │   ├── workshop_job_test.go
│   │   │   │   ├── design_file.go        # Entidad DesignFile
│   │   │   │   ├── design_file_test.go
│   │   │   │   └── repository.go         # Interface del repositorio
│   │   │   │
│   │   │   └── shared/              # Value Objects compartidos entre dominios
│   │   │       ├── money.go         # Value Object Money
│   │   │       ├── money_test.go
│   │   │       ├── percentage.go    # Value Object Percentage
│   │   │       ├── percentage_test.go
│   │   │       ├── email.go         # Value Object Email
│   │   │       ├── email_test.go
│   │   │       ├── nif.go           # Value Object NIF/CIF
│   │   │       ├── nif_test.go
│   │   │       ├── address.go       # Value Object Address
│   │   │       └── address_test.go
│   │   │
│   │   ├── application/             # CAPA DE APLICACIÓN (Casos de Uso)
│   │   │   │                        # Orquesta la lógica de dominio
│   │   │   │
│   │   │   ├── party/               # Casos de uso de Party
│   │   │   │   ├── create_party.go           # Caso de uso: Crear Party
│   │   │   │   ├── create_party_test.go
│   │   │   │   ├── update_party.go           # Caso de uso: Actualizar Party
│   │   │   │   ├── update_party_test.go
│   │   │   │   ├── assign_role.go            # Caso de uso: Asignar rol
│   │   │   │   ├── assign_role_test.go
│   │   │   │   ├── manage_hierarchy.go       # Caso de uso: Gestionar jerarquía
│   │   │   │   ├── manage_hierarchy_test.go
│   │   │   │   ├── register_supplier_cost.go # Caso de uso: Registrar coste
│   │   │   │   ├── register_supplier_cost_test.go
│   │   │   │   └── party_service.go          # Servicio de aplicación
│   │   │   │
│   │   │   ├── product/             # Casos de uso de Producto
│   │   │   │   ├── create_product.go
│   │   │   │   ├── create_product_test.go
│   │   │   │   ├── add_variant.go
│   │   │   │   ├── add_variant_test.go
│   │   │   │   └── product_service.go
│   │   │   │
│   │   │   ├── pricing/             # Casos de uso de Tarificación
│   │   │   │   ├── calculate_price.go
│   │   │   │   ├── calculate_price_test.go
│   │   │   │   ├── create_pricing_rule.go
│   │   │   │   ├── create_pricing_rule_test.go
│   │   │   │   └── pricing_service.go
│   │   │   │
│   │   │   ├── sales/               # Casos de uso de Ventas
│   │   │   │   ├── create_order.go
│   │   │   │   ├── create_order_test.go
│   │   │   │   ├── add_order_line.go
│   │   │   │   ├── add_order_line_test.go
│   │   │   │   ├── change_order_state.go
│   │   │   │   ├── change_order_state_test.go
│   │   │   │   └── sales_service.go
│   │   │   │
│   │   │   └── mes/                 # Casos de uso de MES
│   │   │       ├── create_custom_order.go
│   │   │       ├── create_custom_order_test.go
│   │   │       ├── change_production_state.go
│   │   │       ├── change_production_state_test.go
│   │   │       ├── assign_workshop_job.go
│   │   │       ├── assign_workshop_job_test.go
│   │   │       └── mes_service.go
│   │   │
│   │   ├── infrastructure/          # CAPA DE INFRAESTRUCTURA (Adaptadores)
│   │   │   │                        # Implementaciones concretas de puertos
│   │   │   │
│   │   │   ├── persistence/         # Persistencia de datos
│   │   │   │   ├── postgres/       # Implementación PostgreSQL
│   │   │   │   │   ├── party_repository.go     # Implementa domain/party/repository
│   │   │   │   │   ├── party_repository_test.go # Tests de integración
│   │   │   │   │   ├── product_repository.go
│   │   │   │   │   ├── product_repository_test.go
│   │   │   │   │   ├── pricing_repository.go
│   │   │   │   │   ├── pricing_repository_test.go
│   │   │   │   │   ├── order_repository.go
│   │   │   │   │   ├── order_repository_test.go
│   │   │   │   │   ├── mes_repository.go
│   │   │   │   │   ├── mes_repository_test.go
│   │   │   │   │   └── db.go                    # Configuración de DB
│   │   │   │   │
│   │   │   │   └── migrations/      # Migraciones de base de datos
│   │   │   │       ├── 000001_create_schema.up.sql
│   │   │   │       ├── 000001_create_schema.down.sql
│   │   │   │       ├── 000002_create_parties.up.sql
│   │   │   │       ├── 000002_create_parties.down.sql
│   │   │   │       ├── 000003_create_products.up.sql
│   │   │   │       ├── 000003_create_products.down.sql
│   │   │   │       └── ... (una pareja up/down por migración)
│   │   │   │
│   │   │   ├── storage/             # Almacenamiento de archivos
│   │   │   │   └── nas/
│   │   │   │       ├── file_storage.go      # Adaptador NAS
│   │   │   │       └── file_storage_test.go
│   │   │   │
│   │   │   └── security/            # Seguridad y autenticación
│   │   │       ├── jwt.go           # Generación y validación JWT
│   │   │       ├── jwt_test.go
│   │   │       ├── password.go      # Hash de passwords
│   │   │       ├── password_test.go
│   │   │       ├── rbac.go          # Control de acceso basado en roles
│   │   │       └── rbac_test.go
│   │   │
│   │   └── interfaces/              # CAPA DE INTERFACES (Controllers)
│   │       │                        # Puntos de entrada a la aplicación
│   │       │
│   │       ├── http/                # Interfaz HTTP (REST API)
│   │       │   ├── router.go        # Configuración de rutas
│   │       │   │
│   │       │   ├── middleware/      # Middlewares HTTP
│   │       │   │   ├── auth.go      # Autenticación JWT
│   │       │   │   ├── cors.go      # CORS
│   │       │   │   ├── logger.go    # Logging de requests
│   │       │   │   └── error_handler.go # Manejo de errores
│   │       │   │
│   │       │   ├── handlers/        # Handlers por módulo
│   │       │   │   ├── auth_handler.go       # Login/logout
│   │       │   │   ├── party_handler.go      # CRUD Party
│   │       │   │   ├── party_handler_test.go
│   │       │   │   ├── product_handler.go    # CRUD Producto
│   │       │   │   ├── product_handler_test.go
│   │       │   │   ├── pricing_handler.go    # Cálculo de precios
│   │       │   │   ├── pricing_handler_test.go
│   │       │   │   ├── order_handler.go      # CRUD Pedidos
│   │       │   │   ├── order_handler_test.go
│   │       │   │   ├── mes_handler.go        # Gestión MES
│   │       │   │   └── mes_handler_test.go
│   │       │   │
│   │       │   └── dto/              # Data Transfer Objects (request/response)
│   │       │       ├── auth_dto.go
│   │       │       ├── party_dto.go
│   │       │       ├── product_dto.go
│   │       │       ├── pricing_dto.go
│   │       │       ├── order_dto.go
│   │       │       └── mes_dto.go
│   │       │
│   │       └── cli/                  # Interfaz CLI (comandos opcionales)
│   │           ├── migrate.go        # Comando para migraciones
│   │           ├── seed.go           # Comando para seed de datos
│   │           └── user.go           # Comando para gestión de usuarios
│   │
│   ├── pkg/                          # Código compartido (público)
│   │   ├── logger/                  # Librería de logging
│   │   │   ├── logger.go
│   │   │   └── logger_test.go
│   │   ├── validator/               # Validaciones genéricas
│   │   │   ├── validator.go
│   │   │   └── validator_test.go
│   │   └── errors/                  # Errores personalizados
│   │       ├── errors.go
│   │       └── errors_test.go
│   │
│   ├── config/                       # Configuración de la aplicación
│   │   ├── config.go                # Carga de configuración
│   │   ├── config.yaml              # Configuración por defecto
│   │   ├── config.dev.yaml          # Configuración desarrollo
│   │   └── config.prod.yaml         # Configuración producción
│   │
│   ├── scripts/                      # Scripts de utilidad
│   │   ├── setup.sh                 # Setup inicial de proyecto
│   │   ├── migrate.sh               # Ejecutar migraciones
│   │   ├── seed.sh                  # Seed de datos de prueba
│   │   └── test.sh                  # Ejecutar todos los tests
│   │
│   ├── go.mod                        # Dependencias Go
│   ├── go.sum                        # Checksums de dependencias
│   ├── Makefile                      # Comandos de desarrollo
│   └── README.md                     # Documentación del backend
│
├── frontend/                         # FRONTEND (Vue.js 3)
│   ├── public/                      # Archivos públicos estáticos
│   │   ├── index.html
│   │   └── favicon.ico
│   │
│   ├── src/
│   │   ├── assets/                  # Assets (estilos, imágenes)
│   │   │   ├── styles/
│   │   │   │   ├── main.css
│   │   │   │   ├── tailwind.css
│   │   │   │   └── print.css        # Estilos para Web-to-Print
│   │   │   └── images/
│   │   │       └── logo.png
│   │   │
│   │   ├── components/              # Componentes reutilizables
│   │   │   ├── common/              # Componentes comunes
│   │   │   │   ├── Button.vue
│   │   │   │   ├── Input.vue
│   │   │   │   ├── Select.vue
│   │   │   │   ├── Table.vue
│   │   │   │   ├── Modal.vue
│   │   │   │   ├── Card.vue
│   │   │   │   └── Loading.vue
│   │   │   │
│   │   │   ├── party/               # Componentes de Party
│   │   │   │   ├── PartyCard.vue
│   │   │   │   ├── PartySelector.vue
│   │   │   │   └── HierarchyTree.vue
│   │   │   │
│   │   │   ├── product/             # Componentes de Producto
│   │   │   │   ├── ProductCard.vue
│   │   │   │   ├── ProductSelector.vue
│   │   │   │   └── VariantSelector.vue
│   │   │   │
│   │   │   ├── pricing/             # Componentes de Tarificación
│   │   │   │   └── PriceBreakdown.vue
│   │   │   │
│   │   │   ├── sales/               # Componentes de Ventas
│   │   │   │   ├── OrderCard.vue
│   │   │   │   ├── OrderLineItem.vue
│   │   │   │   └── OrderStateChip.vue
│   │   │   │
│   │   │   └── mes/                 # Componentes de MES
│   │   │       ├── ProductionStateChip.vue
│   │   │       ├── WorkshopJobCard.vue
│   │   │       └── DesignViewer.vue
│   │   │
│   │   ├── views/                   # Vistas (páginas)
│   │   │   ├── auth/
│   │   │   │   └── Login.vue        # Página de login
│   │   │   │
│   │   │   ├── party/               # Vistas de Party
│   │   │   │   ├── PartyList.vue    # Listado de clientes/proveedores
│   │   │   │   ├── PartyForm.vue    # Formulario crear/editar
│   │   │   │   ├── PartyDetail.vue  # Detalle de Party
│   │   │   │   └── PartyHierarchy.vue # Vista de jerarquía empresarial
│   │   │   │
│   │   │   ├── product/             # Vistas de Producto
│   │   │   │   ├── ProductList.vue   # Listado de productos
│   │   │   │   ├── ProductForm.vue   # Formulario crear/editar
│   │   │   │   ├── ProductDetail.vue # Detalle de producto
│   │   │   │   └── ProductCatalog.vue # Vista de catálogo
│   │   │   │
│   │   │   ├── pricing/             # Vistas de Tarificación
│   │   │   │   └── PricingCalculator.vue # Calculadora de precios
│   │   │   │
│   │   │   ├── sales/               # Vistas de Ventas
│   │   │   │   ├── OrderList.vue     # Listado de pedidos
│   │   │   │   ├── OrderForm.vue     # Formulario crear/editar pedido
│   │   │   │   └── OrderDetail.vue   # Detalle de pedido
│   │   │   │
│   │   │   └── mes/                 # Vistas de MES
│   │   │       ├── WorkshopTerminal.vue      # Terminal de taller (tablet)
│   │   │       ├── ProductionDashboard.vue   # Dashboard de producción
│   │   │       └── CustomOrderDetail.vue     # Detalle pedido personalizado
│   │   │
│   │   ├── stores/                  # Pinia stores (estado global)
│   │   │   ├── auth.js              # Store de autenticación
│   │   │   ├── party.js             # Store de Party
│   │   │   ├── product.js           # Store de Producto
│   │   │   ├── pricing.js           # Store de Tarificación
│   │   │   ├── sales.js             # Store de Ventas
│   │   │   └── mes.js               # Store de MES
│   │   │
│   │   ├── router/                  # Vue Router
│   │   │   └── index.js             # Configuración de rutas
│   │   │
│   │   ├── services/                # API clients (llamadas al backend)
│   │   │   ├── api.js               # Configuración base de axios
│   │   │   ├── authService.js       # Servicio de autenticación
│   │   │   ├── partyService.js      # Servicio API Party
│   │   │   ├── productService.js    # Servicio API Producto
│   │   │   ├── pricingService.js    # Servicio API Tarificación
│   │   │   ├── salesService.js      # Servicio API Ventas
│   │   │   └── mesService.js        # Servicio API MES
│   │   │
│   │   ├── composables/             # Composition API helpers
│   │   │   ├── useAuth.js           # Composable de autenticación
│   │   │   ├── useParty.js          # Composable de Party
│   │   │   ├── useProduct.js        # Composable de Producto
│   │   │   ├── usePricing.js        # Composable de Tarificación
│   │   │   ├── useOrder.js          # Composable de Pedidos
│   │   │   └── useMES.js            # Composable de MES
│   │   │
│   │   ├── utils/                   # Utilidades
│   │   │   ├── validators.js        # Validaciones frontend
│   │   │   ├── formatters.js        # Formateo de datos (fechas, moneda)
│   │   │   ├── constants.js         # Constantes de la aplicación
│   │   │   └── helpers.js           # Funciones helper generales
│   │   │
│   │   ├── App.vue                  # Componente raíz
│   │   └── main.js                  # Punto de entrada frontend
│   │
│   ├── package.json                 # Dependencias npm
│   ├── package-lock.json            # Lock de dependencias
│   ├── vite.config.js               # Configuración Vite
│   ├── tailwind.config.js           # Configuración Tailwind CSS
│   ├── .eslintrc.js                 # Configuración ESLint
│   ├── .prettierrc                  # Configuración Prettier
│   └── README.md                    # Documentación del frontend
│
├── docker/                          # DOCKER
│   ├── docker-compose.yml           # Desarrollo
│   ├── docker-compose.prod.yml      # Producción
│   ├── Dockerfile.backend           # Imagen backend
│   ├── Dockerfile.frontend          # Imagen frontend
│   └── postgres/
│       ├── init.sql                 # Script de inicialización PostgreSQL
│       └── postgres.conf            # Configuración PostgreSQL
│
└── .github/                         # CI/CD (GitHub Actions)
    └── workflows/
        ├── backend-ci.yml           # Pipeline backend (tests, build)
        └── frontend-ci.yml          # Pipeline frontend (tests, build)
```

---

### Convenciones de Nombres

#### Backend (Go)

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
backend/
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

**Backend:**

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
├── backend/
│   ├── cmd/api/
│   ├── internal/
│   │   ├── infrastructure/security/  # JWT, RBAC
│   │   └── interfaces/http/
│   │       ├── middleware/
│   │       └── handlers/
│   └── config/
├── frontend/
│   └── src/
│       ├── views/auth/
│       ├── stores/
│       └── services/
└── docker/
```

#### Fase 1 (Party + Producto + Tarificación)

```
# Se añaden:
backend/internal/domain/
├── party/
├── product/
└── pricing/

backend/internal/application/
├── party/
├── product/
└── pricing/

backend/internal/infrastructure/persistence/postgres/
├── party_repository.go
├── product_repository.go
└── pricing_repository.go

backend/internal/infrastructure/persistence/migrations/
├── 000002_create_parties.up.sql
├── 000003_create_products.up.sql
└── 000004_create_pricing.up.sql

frontend/src/
├── views/party/
├── views/product/
├── views/pricing/
├── components/party/
├── components/product/
└── components/pricing/

docs/modules/
├── party/
├── product/
└── pricing/
```

#### Fase 2 (Pedidos)

```
# Se añaden:
backend/internal/domain/sales/
backend/internal/application/sales/
backend/internal/infrastructure/persistence/postgres/order_repository.go
frontend/src/views/sales/
frontend/src/components/sales/
docs/modules/sales/
```

#### Fase 3 (MES)

```
# Se añaden:
backend/internal/domain/mes/
backend/internal/application/mes/
backend/internal/infrastructure/
├── persistence/postgres/mes_repository.go
└── storage/nas/
frontend/src/views/mes/
frontend/src/components/mes/
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