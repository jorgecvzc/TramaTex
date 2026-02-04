# NEXT_SESSION.md - Session Log

This file acts as a log for active or paused development sessions.
Each session below represents a distinct line of work that can be resumed.

---

## Session: Party Module Refactoring & Implementation (2026-02-02 - 2026-02-03)

### Status: Paused
### Context/Last Known State:
The deep analysis and design phases (Fases 1-5) for the `Party` module are completed and approved. A significant misalignment was identified between the approved design (ADR-012) and the existing implementation, necessitating a major refactoring. UI review also indicated a need for redesign. The task "Refactorización/implementación del módulo Party" has been initiated.

Current state of backend implementation for Party:
- Endpoints `/parties` are wired and functional.
- PostgreSQL v2 repositories are implemented.
- Migration v1 → v2 is created.
- Unit tests for domain and basic integration tests are added.
- Party v2 documentation is updated and aligned with ADR-012.
- IAM use cases and API contracts adjusted.
- Frontend Party still uses v1 model (organizations/persons/addresses) and is not compatible with v2 endpoints.

### Tasks to Follow:
- Decide strategy for compatibility with old `/organizations` endpoints (maintain or deprecate).
- Migrate Frontend Party to v2:
    - Update `partyApi.js` to use `/parties`, roles, relationships, and contact-details.
    - Refactor components/routers for Party v2.
    - Adjust frontend tests if applicable.
- Enable local PostgreSQL (Docker) and validate migration v1 → v2 in a controlled environment; document results.
- Execute the full test suite for Party (domain + app + integration) and record coverage (integration requires PostgreSQL).
- Verify Party module coverage (domain + repositories + app) and adjust tests if needed.
- Review compatibility with old `/organizations` endpoints (decide whether to keep or deprecate).
- Refactor frontend (routes, components, stores, services) with Tailwind (defer to later).

### Relevant Context Files:
- `agents/sprint-session-loader.yaml`
- `agents/project/project-context.yaml`
- `project-scaffolding/tmp/module-deep-review/SKILL.md`
- `project-scaffolding/tmp/module-deep-review/references/module-review-guide.md`
- `docs/architecture/adrs/ADR-012-arquitectura-modulo-party.md`
- `docs/log/sprints/sprint-05/02-refactorizacion-implementacion-modulo-party.md`
- `agents/project/context/code-standards.yaml`
- `agents/project/context/tech-stack.yaml`
- `docs/architecture/design-system/`
- `docs/modules/party/use-cases.md`
- `docs/modules/party/domain-model.md`
- `apps/tramatex-api/migrations/007_create_party_v2_tables.sql`
- `apps/tramatex-api/migrations/008_migrate_party_v1_to_v2.sql`
- `apps/tramatex-api/internal/party/` (domain/application/persistence/interfaces)
- `apps/tramatex-api/cmd/api/main.go`
- `apps/tramatex-api/internal/party/interfaces/gin_handlers.go`
- `apps/tramatex-api/internal/party/application/commands.go`
- `apps/tramatex-api/internal/party/application/queries.go`
- `apps/frontend/src/components/party/`
- `apps/frontend/src/services/partyApi.js`
- `apps/frontend/src/router/`

---

## Session: Product Module Backend Implementation (Sprint 06)

### Status: Active
### Context/Last Known State:
**Phase 1: Domain Design and API Specification** is completed and approved.
**Phase 2: Backend Implementation (Test-Driven)** is in progress.

Key progress:
- DB Migration: `009_create_product_tables.sql` created.
- Domain Layer: Entities (`Attribute`, `Product`, `ProductVariant`, etc.) and unit tests implemented in `internal/product/domain/`.
- Application Layer: Initial structure (`commands.go`, `queries.go`), `product_service.go` with use cases, updated dependency injection, and unit tests (`product_service_test.go`) with mocked repositories (`AttributeRepository`, `ProductVariantRepository`). Repository interfaces defined.
- Infrastructure (Persistence): All GORM repositories (`BrandRepository`, `ProductGroupRepository`, `AttributeRepository`, `ProductRepository`, `ProductVariantRepository`) and data models created.
- Interfaces (HTTP): `product_handler.go` implemented. Repositories, services, and handlers integrated into `apps/tramatex-api/cmd/api/main.go`.

### Tasks to Follow:
- PostgreSQL for integration tests (`product_service_integration_test.go`) is active and available in Docker.
- Continue developing integration tests for remaining use cases.
- Implement remaining HTTP handlers in `internal/product/interfaces/` to expose use cases via API, following `api-contracts.md`.

### Relevant Context Files:
- `docs/modules/product/`
- `agents/project/context/architecture.yaml`
- `agents/project/context/code-standards.yaml`
- `apps/tramatex-api/internal/product/`
