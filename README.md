# TramaTex - README

[![Backend CI](https://github.com/joran-cortez/tramatex/actions/workflows/backend.yml/badge.svg)](https://github.com/joran-cortez/tramatex/actions/workflows/backend.yml)
[![Frontend CI](https://github.com/joran-cortez/tramatex/actions/workflows/frontend.yml/badge.svg)](https://github.com/joran-cortez/tramatex/actions/workflows/frontend.yml)
[![codecov](https://codecov.io/gh/joran-cortez/tramatex/branch/master/graph/badge.svg)](https://codecov.io/gh/joran-cortez/tramatex)


## 🚀 Punto de Partida

- **[📖 Manual del Proyecto TramaTex](./docs/architecture/project-vision-and-scope.md)** ← **LEER PRIMERO**
- **[📊 Estado del Proyecto](./docs/log/project-status.md)**

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

## 📚 Documentación

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
- **Cobertura objetivo:** ≥75% global, 100% en Pricing

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