# TramaTex - README

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8.svg)
![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D.svg)
![License](https://img.shields.io/badge/license-Proprietary-red.svg)

**Estado:** ✅ MVP Completado - Listo para Producción


## 🚀 Punto de Partida

- **[📖 Manual del Proyecto TramaTex](./docs/architecture/project-vision-and-scope.md)** ← **LEER PRIMERO**
- **[📊 Estado del Proyecto](./docs/log/project-status.md)**
- **[⚡ Quick Start](./docs/guides/quick-start.md)**

---

## vision-icon📋 Visión

Proporcionar una solución integrada que permita a microempresas:
- **Gestionar clientes y proveedores** de forma centralizada (Party)
- **Administrar catálogo de productos** con variantes (tallas, colores, modificaciones)
- **Calcular precios** de forma inteligente y automática (tarificación)
- **Gestionar pedidos** desde cotización hasta entrega
- **Controlar producción personalizada** con seguimiento estado-a-estado
- **Documentar procesos** con trazabilidad completa

---

## 🏗️ Arquitectura

### Stack Tecnológico
- **Backend:** Go 1.21+ (Clean Architecture + DDD)
- **Frontend:** Vue.js 3 + Vite + Tailwind CSS
- **Base de Datos:** PostgreSQL 15+
- **Contenedorización:** Docker + Docker Compose
- **Testing:** TDD (Go testing + Vitest)

### Estructura

TramaTex sigue un **monolito modular** basado en **Clean Architecture y Domain-Driven Design**:

```
tramatex/
├── apps/
│   ├── tramatex-api/
│   └── frontend/
├── docs/
│   ├── architecture/
│   ├── guides/
│   ├── modules/
│   └── log/
```

Más detalles en [ADR-009 – Estructura de Proyecto](docs/architecture/adrs/ADR-009-project-structure.md).
## 🚀 Quick Start

### Requisitos

- **Docker Desktop** (incluye Docker y Docker Compose)
- **Go 1.21+** (opcional, para desarrollo local)
- **Node.js 18+** (opcional, para desarrollo local)

### Arrancar el Proyecto

#### Con Docker Compose (recomendado)

```bash
cd tramatex
docker-compose up --build
```

El sistema estará disponible en:
- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080
- **PostgreSQL:** localhost:5432

#### Sin Docker (desarrollo local)

**Backend:**
```bash
cd apps/tramatex-api
go install github.com/golang-migrate/migrate/cmd/migrate@latest
make migrate-up
make run
```

**Frontend:**
```bash
cd apps/frontend
npm install
npm run dev
```

---

## � Versionado y Contribución

### Versionado Semántico

TramaTex sigue [Semantic Versioning 2.0](https://semver.org/):

- **MAJOR.MINOR.PATCH** (Ej: 1.0.0)
- **MAJOR:** Cambios rompientes en API o arquitectura
- **MINOR:** Nuevas funcionalidades compatibles
- **PATCH:** Correcciones de bugs

**Versión actual:** `v1.0.0` (MVP Completado - 2026-02-22)

### Estrategia de Branches

TramaTex utiliza **GitFlow Simplificado**:

#### Ramas Principales

- **`main`**: Código en producción, siempre estable y desplegable
  - Solo recibe merges desde `develop` o `hotfix/*`
  - Cada merge representa una release con tag vX.Y.Z
  - **Protegida:** Requiere PR aprobado + tests pasando

- **`develop`**: Rama de integración para desarrollo activo  - Base para todas las feature/bugfix branches
  - **Protegida:** Requiere PR + tests pasando

#### Ramas Temporales

- **`feature/*`**: Nuevas funcionalidades (desde `develop`)
- **`bugfix/*`**: Correcciones no críticas (desde `develop`)
- **`hotfix/*`**: Correcciones críticas en producción (desde `main`)
- **`release/*`**: Preparación de releases (desde `develop`)

### Flujo de Trabajo

#### Feature Development

```bash
# 1. Crear branch desde develop
git checkout develop
git pull origin develop
git checkout -b feature/mi-funcionalidad

# 2. Desarrollar y commitear
git add .
git commit -m "feat(module): add new functionality"

# 3. Push y abrir Pull Request
git push origin feature/mi-funcionalidad
# Abrir PR en GitHub hacia develop
```

#### Commits Convencionales

Usamos [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]
```

**Tipos:**
- `feat`: Nueva funcionalidad
- `fix`: Corrección de bug
- `docs`: Documentación
- `refactor`: Refactorización
- `test`: Tests
- `chore`: Cambios de build/config

**Scopes:** `party`, `product`, `pricing`, `sales`, `iam`, `mes`, `frontend`, `backend`, `infra`

**Ejemplos:**
```bash
feat(pricing): add volume discount rules
fix(party): resolve selector crash on empty results
docs(adr): add versioning strategy ADR-021
test(sales): increase coverage to 80%
```

### Configuración Inicial

#### 1. Clonar Repositorio

```bash
git clone git@github.com:jorgecvzc/TramaTex.git
cd TramaTex
```

#### 2. Configurar Variables de Entorno

```bash
# Backend
cp apps/tramatex-api/.env.example apps/tramatex-api/.env
# Editar apps/tramatex-api/.env con tus valores

# Frontend
cp 📄 Licencia

Proprietary - Todos los derechos reservados

---

## 👥 Autores y Equipo

- **Jorge Cortés Villalba** - Product Owner, Arquitectura
- **AI Assistant (Claude)** - Desarrollo técnico, Arquitectura, Documentación

---

## 📊 Estadísticas del Proyecto

**Líneas de Código (MVP v1.0.0):**
- Backend (Go): ~25,000 líneas
- Frontend (Vue/TS): ~15,000 líneas
- Documentación: ~8,000 líneas
- Tests: ~12,000 líneas

**Cobertura de Tests:**
- Backend: 75%+ (Pricing 85.4%, Party 86.7%, MES 86.9%)
- Frontend: 77.63%

**Módulos Completados:**
- ✅ Party (Clientes/Proveedores)
- ✅ Product (Catálogo y Variantes)
- ✅ Pricing (Motor de Tarificación)
- ✅ Sales (Ciclo Comercial)
- ✅ IAM (Autenticación/Autorización)
- ✅ MES (Producción y Taller)

---

**Última Actualización:** 22/02/2026  
**Versión:** 1.0.0 (MVP Completado)  
**Estado:** ✅ Listo para Produc

- **[ADR-021: Version Control & Branching Strategy](docs/architecture/adrs/ADR-021-version-control-and-branching-strategy.md)**
- **[Guía de Contribución](docs/guides/developer/CONTRIBUTING.md)** [Por crear]

---

## �📚 Documentación

### Architecture Decision Records (ADRs)

- [ADR-006: Estrategia de Desarrollo Dirigida por Dominio](docs/architecture/adrs/ADR-006-domain-driven-development-strategy.md)
- [ADR-007: Orden de Implementación de Módulos](docs/architecture/adrs/ADR-007-module-implementation-order.md)
- [ADR-008: Planificación y Cronograma](docs/architecture/adrs/ADR-008-mvp-timeline-planning.md)
- [ADR-009: Estructura de Carpetas](docs/architecture/adrs/ADR-009-project-structure.md)

### Estado del Proyecto

- [project-status.md](docs/log/project-status.md) - Progreso actual, hitos, timeline

### Guías

- [Setup de Desarrollo](docs/guides/developer/github-setup.md) [Pendiente]
- [Estrategia de Testing](docs/guides/developer/ci-cd.md) [Pendiente]
- [Configuración de impresión Sales (perfil fiscal emisor)](apps/frontend/README.md#configuración-de-emisor-para-impresión-sales)

### Módulos

Cada módulo tiene su documentación en `docs/modules/[modulo]/`:

- **Party:** Gestión de clientes y proveedores
- **Product:** Catálogo de productos y variantes
- **Pricing:** Motor de tarificación
- **Sales:** Gestión de pedidos
- **MES:** Producción personalizada y taller

---

## 📅 Cronograma

**Duración total:** 24 meses (Enero 2026 - Enero 2028)

### Fases

1. **Fase 0 (Q1 2026):** Fundaciones técnicas
   - Setup Docker, Git, autenticación JWT
   - **Hito:** Sistema arranca + login funcional
   - Detalle: [MVP Specification](docs/architecture/project-vision-and-scope.md)

2. **Fase 1 (Q2-Q4 2026):** Dominio base
   - Módulos: Party, Producto, Tarificación
   - **Hito:** Núcleo económico funcional

3. **Fase 2 (Q1-Q2 2027):** Ventas
   - Gestión de pedidos estándar
   - **Hito:** Flujo completo de ventas

4. **Fase 3 (Q3-Q4 2027):** MES
   - Producción personalizada, terminal taller
   - **Hito:** MVP completo

5. **Fase 4 (Q1 2028):** Estabilización
   - **Hito:** En producción estable

Más detalles en [ADR-008 – Cronograma](docs/architecture/adrs/ADR-008-mvp-timeline-planning.md) y [project-status.md](docs/log/project-status.md).

---

## 🛠️ Desarrollo

### Comandos Principales

**Backend:**
```bash
cd apps/tramatex-api
make run              # Ejecutar servidor
make test             # Tests
make test-coverage    # Cobertura
make lint             # Linter
make migrate-up       # Migraciones
make docker-up        # Stack completo en Docker
```

**Frontend:**
```bash
cd apps/frontend
npm run dev           # Servidor desarrollo
npm run build         # Build producción
npm run lint          # Linter
npm run format        # Prettier
```

### Testing

- **Backend:** Go testing + assertions manuales
- **Frontend:** Vitest (similar a Jest)
- **Cobertura objetivo:** ≥85% global, ≥90% en Pricing

### Code Style

- **Backend:** `gofmt` + `golangci-lint`
- **Frontend:** ESLint + Prettier

---

## 🤝 Contribución

[Guía de contribución](docs/guides/developer/CONTRIBUTING.md)

---

## 📄 Licencia

[Pendiente: Especificar licencia]

---

## 👥 Autores

- **Jorge Cortés Villalba** - Producto, Dominio
- **Claude (Anthropic)** - Arquitectura, Copiloto técnico

---

## 📞 Contacto

[Información de contacto según necesario]

---

**Última Actualización:** 11/01/2026  
**Versión:** 0.1.0 (Pre-Fase 0)  
**Estado:** En Planificación