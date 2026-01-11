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

---

## 2. Decisión Adoptada

Se adopta **Estructura por Módulos de Dominio con Clean Architecture** (Alternativa C).

**Justificación:**

- Refleja el dominio de negocio: Party, Producto, Tarificación, Ventas, MES
- Soporta Clean Architecture: Capas separadas dentro de cada módulo
- Facilita TDD: Tests están junto al código que prueban
- Escalabilidad: Módulos pueden extraerse a servicios sin refactoring masivo
- Navegación intuitiva: Toda la lógica de Party está en carpetas `*/party/`

---

## 3. Estructura de Carpetas

```
tramatex/
├── docs/                              # DOCUMENTACIÓN
│   ├── adr/                          # Architecture Decision Records
│   │   ├── ADR-006-*.md
│   │   ├── ADR-007-*.md
│   │   ├── ADR-008-*.md
│   │   └── ADR-009-*.md
│   ├── modules/                      # Documentación de módulos
│   │   ├── _TEMPLATE.md
│   │   ├── party/
│   │   │   ├── module-spec.md
│   │   │   ├── domain-model.md
│   │   │   ├── use-cases.md
│   │   │   └── api-contracts.md
│   │   ├── product/
│   │   ├── pricing/
│   │   ├── sales/
│   │   └── mes/
│   ├── sessions/                     # Documentación de sesiones
│   │   └── _SESSION_TEMPLATE.md
│   ├── diagrams/
│   │   ├── architecture/
│   │   ├── domain/
│   │   └── flows/
│   ├── consolidated/
│   └── guides/
├── backend/                          # BACKEND (Go)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── party/
│   │   │   ├── product/
│   │   │   ├── pricing/
│   │   │   ├── sales/
│   │   │   ├── mes/
│   │   │   └── shared/
│   │   ├── application/
│   │   │   ├── party/
│   │   │   ├── product/
│   │   │   ├── pricing/
│   │   │   ├── sales/
│   │   │   └── mes/
│   │   ├── infrastructure/
│   │   │   ├── persistence/
│   │   │   │   ├── postgres/
│   │   │   │   └── migrations/
│   │   │   ├── storage/
│   │   │   └── security/
│   │   └── interfaces/
│   │       ├── http/
│   │       │   ├── middleware/
│   │       │   ├── handlers/
│   │       │   └── dto/
│   │       └── cli/
│   ├── pkg/
│   ├── config/
│   ├── scripts/
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile
│   └── README.md
├── frontend/                         # FRONTEND (Vue.js 3)
│   ├── public/
│   ├── src/
│   │   ├── assets/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/
│   │   ├── router/
│   │   ├── services/
│   │   ├── composables/
│   │   ├── utils/
│   │   ├── App.vue
│   │   └── main.js
│   ├── package.json
│   ├── package-lock.json
│   ├── vite.config.js
│   ├── tailwind.config.js
│   ├── .eslintrc.js
│   ├── .prettierrc
│   └── README.md
├── docker/                          # DOCKER
│   ├── docker-compose.yml
│   ├── docker-compose.prod.yml
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   └── postgres/
├── .github/                         # CI/CD
│   └── workflows/
├── .gitignore
├── README.md
├── LICENSE
└── Makefile
```

---

## 4. Convenciones de Nombres

### Backend (Go)

- **Archivos:** snake_case + `_test.go` (ej: `party_repository.go`)
- **Paquetes:** minúsculas, singular (ej: `party`, `product`)
- **Entidades:** PascalCase (ej: `Party`, `Customer`, `Money`)
- **Funciones:** PascalCase exportadas, camelCase privadas

### Frontend (Vue.js)

- **Componentes:** PascalCase multi-palabra (ej: `PartyList.vue`)
- **Servicios/Stores:** camelCase (ej: `partyService.js`, `auth.js`)

### Base de Datos (PostgreSQL)

- **Tablas:** snake_case, plural (ej: `parties`, `products`)
- **Columnas:** snake_case (ej: `party_id`, `created_at`)

---

## 5. Tests

- **Tests Unitarios:** Mismo paquete que código (`party.go` + `party_test.go`)
- **Tests Integración:** Mismo paquete que adaptador
- **Tests E2E:** Carpeta `/backend/test/e2e/` (Post-MVP)

---

## Consecuencias

### Positivas

- Claridad de dominio: Fácil encontrar código por módulo
- Separación de responsabilidades: Clean Architecture respetada
- Testabilidad: Tests junto al código
- Documentación trazable: Estructura espeja al código
- Evolución controlada: Nuevos módulos sin afectar existentes
- Escalabilidad: Módulos extraíbles a servicios

### Negativas

- Profundidad de carpetas: Rutas pueden ser largas
- Duplicación aparente: Carpetas `party/` en domain, application, infrastructure
- Disciplina requerida: Fácil romper separación si no hay code reviews
- Curva de aprendizaje: Requiere entender Clean Architecture

---

## Referencias

- ADR-001, ADR-002, ADR-003: Stack, arquitectura, modularidad
- ADR-006, ADR-007: Estrategia y orden de desarrollo
- Golang Project Layout: https://github.com/golang-standards/project-layout
- Vue.js Style Guide: https://vuejs.org/style-guide/

---

**Fin del ADR-009**
