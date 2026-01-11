# PROMPT DE SESIÓN DE TRABAJO – TramaTex

**Versión:** 1.0  
**Última actualización:** 11/01/2026  
**Propósito:** Proporcionar contexto completo sobre el proyecto TramaTex a Claude (Anthropic) / GitHub Copilot para cada nueva sesión de desarrollo.

---

## 📍 CONTEXTO DEL PROYECTO

### 1.1 Descripción General

**TramaTex** es un **ERP/MES (Enterprise Resource Planning / Manufacturing Execution System)** diseñado para microempresas del sector textil y de personalización (venta de EPIs y vestuario laboral personalizable).

**Objetivo principal:** Controlar de forma fiable, trazable e integrada el ciclo completo de pedidos (estándar y personalizados), desde cotización hasta entrega, incluyendo gestión de ventas, tarificación, inventario y producción.

**Hardware objetivo:**
- Ordenadores i3 con 8GB RAM (clientes)
- Servidor Linux i3 con 16GB RAM, SSD 2TB en espejo
- Tablets básicas para taller

### 1.2 Stack Tecnológico

| Componente | Tecnología |
|-----------|-----------|
| **Backend** | Go 1.21+, Clean Architecture + DDD, Gin Gonic, GORM (solo capa infraestructura) |
| **Frontend** | Vue.js 3, Composition API, Pinia (estado global), Tailwind CSS |
| **Base de Datos** | PostgreSQL 14+, ACID compliance |
| **Contenedorización** | Docker + Docker Compose (MVP local-first, sin Kubernetes) |
| **Documentación** | Web-to-Print delegado a frontend (no servidor PDF) |
| **Testing** | TDD obligatorio en dominio crítico (Go testing + Vitest) |
| **Seguridad** | JWT, RBAC mínimo (Admin, Comercial, Diseño, Taller) |
| **Almacenamiento** | PostgreSQL indexado + NAS para archivos pesados (Post-MVP) |

### 1.3 Principios Arquitectónicos

1. **Monolito Modular Local-First:**
   - Un único proceso Go (backend)
   - Múltiples dominios modularizados (Clean Architecture)
   - Preparado para extracción futura a microservicios
   - Sin dependencias cloud en MVP

2. **Domain-Driven Design (DDD) con Rigor Asimétrico:**
   - **Rigor estricto** en dominio crítico (tarificación, Party, Producto)
   - **Rigor flexible** en infraestructura y casos de uso de bajo valor
   - Dominio completamente testeable en aislamiento
   - Entidades, Value Objects, Servicios de dominio como activos estratégicos

3. **Clean Architecture:**
   - Capa de dominio (Entities, Value Objects, Domain Services)
   - Capa de aplicación (Use Cases, Orchestration)
   - Capa de infraestructura (GORM, PostgreSQL, Adapters)
   - Capa de interfaces (Controllers REST, DTOs, Serialización)
   - Dependencias siempre hacia el interior

---

## 📦 DEFINICIÓN DEL MVP

### 2.1 Fases de Implementación

| Fase | Hito | Duración Estimada | Estado |
|------|------|------------------|--------|
| **Fase 0** | Fundaciones Técnicas (Setup, Auth JWT, CI/CD) | 4 semanas (32h) | ⏳ En Progreso |
| **Fase 1** | Dominio Base (Party, Producto, Tarificación) | 8 semanas (64h) | 📋 Próximo |
| **Fase 2** | Pedidos Estándar (Ventas, Documentación) | 8 semanas (64h) | 📋 Próximo |
| **Fase 3** | MES y Especialización (Taller, Producción) | 12 semanas (96h) | 📋 Próximo |
| **Post-MVP** | Optimización, Despliegue, Documentación Final | 4 semanas (32h) | 📋 Futuro |

**Cronograma Total:** 24 meses (Enero 2026 - Enero 2028), 8h/semana, 782 horas proyectadas.

### 2.2 Módulos Canónicos (Bounded Contexts)

| Módulo | Responsabilidad | Fase | Dependencias |
|--------|-----------------|------|------|
| **Party/Organización** | Gestión unificada de clientes, proveedores, contactos | Fase 1 | — |
| **Producto/Variante/Categoría** | Catálogo, variantes (talla, color), clasificación | Fase 1 | Party (proveedor) |
| **Tarificación** | Motor de cálculo de precios, márgenes, descuentos | Fase 1 | Party, Producto |
| **Ventas/Pedidos** | Gestión de pedidos estándar, cotizaciones, documentación | Fase 2 | Todos Fase 1 |
| **MES** | Control de producción personalizada, estados, trazabilidad | Fase 3 | Ventas |
| **Seguridad** | Autenticación (JWT), Autorización (RBAC), Auditoría | Fase 0 | — |
| **Gestión Documental** | Almacenamiento indexado, vinculación a pedidos | Fase 3 | Ventas |

### 2.3 Requisitos Funcionales Clave (RF)

**Fase 0 (Fundaciones):**
- RF.AUTH: Autenticación JWT con refresh tokens
- RF.RBAC: Control de acceso basado en roles (Admin, Comercial, Diseño, Taller)

**Fase 1 (Dominio Base):**
- RF.PARTY: CRUD Party (clientes/proveedores), jerarquías de clientes, descuentos heredados
- RF.PRODUCT: CRUD Producto, variantes (talla/color), categorías, relación proveedor-producto
- RF.PRICING: Cálculo automático de precios, márgenes, descuentos base/específicos, histórico de costes

**Fase 2 (Ventas):**
- RF.SALES: Creación/edición/seguimiento de pedidos, selección cliente/productos, cálculo precio automático
- RF.DOCS: Generación de presupuestos/albaranes/facturas (Web-to-Print frontend)
- RF.WORKFLOW: Estados de pedido (Borrador → Confirmado → Preparación → Entregado)

**Fase 3 (MES):**
- RF.MES: Estados de producción (Diseño → Aprobación → Marcaje → Taller → Control Calidad → Listo)
- RF.TERMINAL: Interfaz simplificada para tablets (taller)
- RF.DOCS_ADV: Almacenamiento y vinculación de archivos de diseño

### 2.4 Requisitos No Funcionales (RNF)

| ID | Descripción | Métrica |
|----|-------------|---------|
| RNF1 | Eficiencia | Consumo RAM < 150MB en operación normal |
| RNF2 | Local-first | Sin dependencias cloud, 100% operativo offline |
| RNF3 | Mantenibilidad | TDD en dominio, cobertura ≥75% en crítico |
| RNF4 | Integridad | ACID PostgreSQL, transacciones consistentes |
| RNF5 | Seguridad | Hash bcrypt, JWT, RBAC, auditoría cambios críticos |
| RNF6 | Performance | <100ms latencia API en operación normal |
| RNF7 | Disponibilidad | Backups automáticos, recuperación <1h |

---

## 🏗️ ESTRUCTURA DEL PROYECTO

```
tramatex/
├── backend/                 # Go, Clean Architecture
│   ├── cmd/
│   │   └── api/            # Entry point servidor REST
│   ├── internal/
│   │   ├── domain/         # Capa de dominio (Entities, Value Objects, Domain Services)
│   │   │   ├── party/      # Bounded Context: Party/Organización
│   │   │   ├── product/    # Bounded Context: Producto/Variante/Categoría
│   │   │   ├── pricing/    # Bounded Context: Tarificación
│   │   │   ├── sales/      # Bounded Context: Ventas/Pedidos
│   │   │   ├── mes/        # Bounded Context: MES/Producción
│   │   │   └── common/     # Shared Value Objects, Base Entities
│   │   ├── application/    # Capa de aplicación (Use Cases, Orchestration)
│   │   │   ├── party/
│   │   │   ├── product/
│   │   │   ├── pricing/
│   │   │   ├── sales/
│   │   │   ├── mes/
│   │   │   └── dto/        # DTOs de entrada/salida
│   │   ├── infrastructure/ # Adaptadores, Persistencia
│   │   │   ├── persistence/ # Repositories, GORM models
│   │   │   ├── repositories/
│   │   │   ├── config/
│   │   │   └── security/
│   │   └── interfaces/     # Controllers REST, Serialización
│   │       ├── http/
│   │       └── middleware/
│   ├── tests/
│   │   ├── unit/          # Tests unitarios por módulo
│   │   ├── integration/   # Tests de integración
│   │   └── fixtures/      # Datos de prueba
│   ├── migrations/        # Migraciones PostgreSQL
│   ├── config/            # Configuración centralizada
│   ├── go.mod & go.sum    # Dependencias Go
│   ├── Makefile           # Comandos backend
│   └── Dockerfile
│
├── frontend/              # Vue.js 3
│   ├── src/
│   │   ├── components/    # Componentes reutilizables
│   │   │   ├── auth/      # Login, Register
│   │   │   ├── party/     # Listados/Forms Party
│   │   │   ├── product/   # Listados/Forms Producto
│   │   │   ├── sales/     # Listados/Forms Pedidos
│   │   │   ├── mes/       # Terminal Taller
│   │   │   └── common/    # Componentes base (Modal, Table, etc.)
│   │   ├── views/         # Páginas por módulo
│   │   ├── stores/        # Pinia state management
│   │   ├── services/      # API clients
│   │   ├── composables/   # Vue composables (lógica reutilizable)
│   │   ├── utils/         # Utilitarios, formatters
│   │   ├── styles/        # Tailwind CSS customización
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── router.js      # Vue Router
│   ├── public/
│   ├── package.json
│   ├── vite.config.js
│   ├── Makefile
│   └── Dockerfile
│
├── docker/                # Configuración Docker
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   ├── docker-compose.yml
│   └── postgresql.conf
│
├── docs/                  # Documentación
│   ├── adr/              # Architecture Decision Records (ADR-001 a ADR-009)
│   ├── consolidated/     # Documento Consolidado 3.0
│   ├── modules/          # Especificaciones de módulos por Bounded Context
│   ├── diagrams/         # Diagramas C4/ER/Flujos
│   └── sessions/         # Registro de sesiones de desarrollo
│
├── Makefile              # Makefile global
├── README.md             # Descripción general
├── PROJECT_STATUS.md     # Estado actual del proyecto
├── LICENSE
└── .gitignore

```

---

## 📚 DOCUMENTACIÓN DISPONIBLE

### ADRs (Architecture Decision Records) – 9 Decisiones Documentadas

1. **ADR-001:** Selección del Stack Tecnológico
2. **ADR-002:** Adopción de Clean Architecture + DDD con Rigor Asimétrico
3. **ADR-003:** Tipo y Distribución (Monolito Modular Local-First)
4. **ADR-004:** Ciclo de Vida de Desarrollo e Implementación MVP
5. **ADR-005:** Gestión Unificada de Party (Clientes/Proveedores)
6. **ADR-006:** Estrategia de Desarrollo Dirigida por Dominio
7. **ADR-007:** Orden de Implementación de Módulos
8. **ADR-008:** Planificación y Cronograma MVP (24 meses)
9. **ADR-009:** Estructura y Organización del Proyecto

**Ubicación:** `/docs/adr/ADR-XXX-*.md`

### Documento Consolidado 3.0

Especificación técnica completa del MVP (requisitos, arquitectura, flujos, casos de uso, matriz de trazabilidad).

**Ubicación:** `/docs/consolidated/DOCUMENTO-CONSOLIDADO-3.0.md`

### Sesiones Anteriores

Registro detallado de cada sesión de trabajo, decisiones tomadas, problemas resueltos, avances.

**Ubicación:** `/docs/sessions/2026-01-DD-session-NN.md`

---

## ⚙️ CONFIGURACIÓN DEL ENTORNO

### Herramientas Requeridas

- **Go 1.21+** (backend)
- **Node.js 18+** (frontend)
- **PostgreSQL 14+** (desarrollo: via Docker)
- **Docker Desktop** (desarrollo)
- **Git**

### Comandos Esenciales

```bash
# Setup global
make help                  # Ver todos los comandos disponibles
make setup                 # Inicializar proyecto

# Docker
make docker-build         # Build imágenes
make docker-up            # Levantar stack
make docker-down          # Bajar stack

# Backend
make backend-test         # Tests backend
make backend-run          # Ejecutar servidor

# Frontend
make frontend-dev         # Dev server (http://localhost:5173)
make frontend-build       # Build producción

# Documentación
make docs-view            # Ver índice de documentación
```

### Variables de Entorno

**Backend** (`.env`):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=tramatex
DB_PASSWORD=dev_password
DB_NAME=tramatex_db
JWT_SECRET=your-secret-key-change-in-prod
LOG_LEVEL=info
```

**Frontend** (`.env.local`):
```
VITE_API_BASE_URL=http://localhost:8080/api
```

---

## 🔄 FLUJO DE TRABAJO ESTÁNDAR

1. **Inicio de Sesión:**
   - Leer este prompt (SESSION_PROMPT.md)
   - Revisar última sesión documentada en `/docs/sessions/`
   - Revisar PROJECT_STATUS.md para estado actual

2. **Durante la Sesión:**
   - Implementar funcionalidades específicas (ver sección "OBJETIVOS DE ESTA SESIÓN" abajo)
   - Escribir tests TDD primero en dominio crítico
   - Mantener Clean Architecture: dependencias hacia el interior
   - Documentar decisiones y problemas en `docs/sessions/YYYY-MM-DD-session-NN.md`

3. **Cierre de Sesión:**
   - Commits organizados y descriptivos
   - Tests verdes (backend + frontend)
   - Documentación sesión completa
   - UPDATE PROJECT_STATUS.md con avances

4. **Validación:**
   - Cobertura de tests ≥75% en dominio crítico
   - Lint sin warnings (Go: golangci-lint, JS: ESLint)
   - Documentación sincronizada con código

---

## 🎯 PATRONES Y CONVENCIONES

### Backend (Go)

**Estructura de Módulo:**
```go
// domain/party/party.go - Entidad
type Organization struct {
    ID        string
    Name      string
    // ...
}

// domain/party/repository.go - Interface (repositorio)
type OrganizationRepository interface {
    Save(ctx context.Context, org *Organization) error
    ByID(ctx context.Context, id string) (*Organization, error)
    // ...
}

// application/party/create_org.go - Use Case
type CreateOrganizationUseCase struct {
    repo OrganizationRepository
}

func (uc *CreateOrganizationUseCase) Execute(ctx context.Context, input CreateOrgInput) (*OrgOutput, error) {
    // Lógica de aplicación
}

// infrastructure/persistence/org_repository.go - Implementación
type GormOrganizationRepository struct {
    db *gorm.DB
}

// interfaces/http/org_controller.go - Handler REST
func HandleCreateOrganization(uc *CreateOrganizationUseCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        // HTTP logic
    }
}
```

**Testing:**
- Unit tests junto al código (`*_test.go`)
- Integration tests en `tests/integration/`
- Fixtures de datos en `tests/fixtures/`
- Mocks solo donde necesarios (interfaces bien definidas)

### Frontend (Vue.js 3)

**Estructura de Módulo:**
```js
// src/components/party/PartyForm.vue
<template>
  <!-- Template -->
</template>

<script setup>
// Composition API
import { ref } from 'vue'
import { usePartyStore } from '@/stores/party'

// Lógica
</script>

// src/stores/party.js - Pinia Store
import { defineStore } from 'pinia'

export const usePartyStore = defineStore('party', () => {
    // State, actions, getters
})

// src/services/partyService.js - API Client
export const getParties = async () => {
    // API call
}

// src/views/PartyListView.vue - Página
<template>
  <PartyList />
</template>
```

**Convenciones:**
- Componentes en PascalCase (PartyForm.vue)
- Stores en camelCase (usePartyStore)
- Composables en camelCase (usePartyForm.js)
- Servicios en camelCase (partyService.js)

---

## 🔗 REFERENCIAS INTERNAS RÁPIDAS

| Tópico | Ubicación |
|--------|-----------|
| Stack Tecnológico | ADR-001 |
| Clean Architecture | ADR-002 |
| Módulos/Bounded Contexts | ADR-003, ADR-005 |
| Cronograma | ADR-008 |
| Estructura Carpetas | ADR-009 |
| Especificación Completa | DOCUMENTO-CONSOLIDADO-3.0.md |
| Sesiones Anteriores | docs/sessions/ |
| Estado Actual | PROJECT_STATUS.md |

---

## ✋ IMPORTANTE: RESTRICCIONES Y PRINCIPIOS

1. **Dominio Protegido:**
   - Cero dependencias externas en capa de dominio
   - Sin ORM, sin frameworks, sin serialización
   - Tests de dominio SIEMPRE pasan (no dependencias de infraestructura)

2. **TDD Obligatorio en Crítico:**
   - Tarificación: 100% cubierto
   - Party: ≥90% cubierto
   - Producto: ≥85% cubierto

3. **Local-First:**
   - Todas las conexiones a bases de datos deben ser configurables
   - No hardcodear endpoints cloud
   - PostgreSQL es único datastore en MVP

4. **Clean Architecture:**
   - Dependencias siempre hacia el interior
   - DTOs en interfaces, modelos de dominio en aplicación
   - Inversión de control: infraestructura implementa interfaces de dominio

5. **Versionado:**
   - Commits descriptivos
   - Tags para hitos (v0.1.0-foundational, v1.0.0-domain-base, etc.)
   - Squash commits en ramas de feature antes de merge

---

## 🚀 PRÓXIMOS PASOS (POST-SESSION)

Después de completar esta sesión:

1. Commits limpios y descriptivos
2. Tests verdes: `make backend-test` y `make frontend-test`
3. Documentar sesión: copiar `/docs/sessions/_SESSION_TEMPLATE.md` → `2026-MM-DD-session-NN.md`
4. Actualizar PROJECT_STATUS.md
5. Crear issues/PRs si corresponde

---

---

---

## 📝 OBJETIVOS DE ESTA SESIÓN

**RELLENA ESTA SECCIÓN AL INICIO DE CADA NUEVA SESIÓN**

### Sesión: [YYYY-MM-DD Session-NN]

**Facilitador/LLM:** [GitHub Copilot / Claude / otro]

**Duración estimada:** [X horas]

**Objetivos principales:**
1. [Objetivo 1]
2. [Objetivo 2]
3. [Objetivo 3]

**Contexto de entrada:**
- Estado actual: [Breve resumen del estado anterior]
- Bloqueadores: [Si hay problemas pendientes]
- Prioridades: [Qué es crítico esta sesión]

**Definición de "Hecho":**
- [ ] Objetivo 1 completado
- [ ] Objetivo 2 completado
- [ ] Objetivo 3 completado
- [ ] Tests ≥75% cobertura en crítico
- [ ] Documentación sesión actualizada
- [ ] PROJECT_STATUS.md actualizado

**Notas especiales:**
- [Aclaraciones, dependencias, riesgos]

---

## 📊 PLANTILLA DE SESIÓN (Para Registro)

Ver: `/docs/sessions/_SESSION_TEMPLATE.md`

Crear archivo: `/docs/sessions/2026-MM-DD-session-NN.md`

