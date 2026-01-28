# ADR-011 – Estrategia de Testing y Coverage

**Fecha:** 2026-01-27
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, GitHub Copilot

---

## Contexto

TramaTex sigue una arquitectura de Clean Architecture + DDD con capas claramente separadas. Se requiere una estrategia de testing formal que garantice:

- Alta confiabilidad del código
- Cobertura consistente en módulos críticos
- Facilidad para refactorizar con confianza
- Tests como documentación ejecutable

---

## Decisión

Se establece una **estrategia de testing multinivel** y una **política de cobertura por fases** que distingue entre los objetivos pragmáticos del MVP y los objetivos ideales a largo plazo (Post-MVP).

### 1. Niveles de Testing (Común a todas las fases)

#### 1.1 Unit Tests (Obligatorios)

**Alcance:**
- Domain layer (entities, value objects, aggregates)
- Application layer (use cases, command/query handlers)
- Casos edge y validaciones

**Herramientas:**
- Backend: Go testing package + testify
- Frontend: Vitest + Vue Test Utils

#### 1.2 Integration Tests (Selectivos)

**Alcance:**
- Persistence layer con base de datos real
- Infrastructure layer (middleware, adapters)
- Casos que requieren múltiples capas

**Herramientas:**
- Backend: Testcontainers (PostgreSQL)
- Frontend: E2E con Playwright (selectivo)

#### 1.3 E2E Tests (Críticos)

**Alcance:**
- Flujos de usuario completos
- Casos de uso principales del MVP
- Smoke tests de features core

**Herramientas:**
- Frontend: Playwright
- API: Postman/Newman (colecciones)

### 2. Política de Coverage por Fases

Para alinear las expectativas de calidad con las metas del proyecto, se definen dos niveles de exigencia.

#### 2.1 Objetivos de Coverage del MVP

Estos son los mínimos requeridos para considerar el MVP como exitoso y listo para producción, según se define en `02-mvp-specification.md`.

| Métrica | Coverage Mínimo | Contexto |
|---|---|---|
| **Coverage total del proyecto** | **≥ 75%** | Global |
| **Dominio de Tarificación** | **≥ 80%** | Crítico para el negocio |
| **Casos de Uso de Pedidos** | **≥ 70%** | Flujo de ventas principal |

#### 2.2 Objetivos de Coverage Post-MVP (Ideales)

Estos son los objetivos a largo plazo para garantizar la máxima calidad y mantenibilidad del producto una vez superada la fase de MVP.

| Capa | Coverage Mínimo | Tipo de Tests |
|------|-----------------|---------------|
| **Domain** | 100% | Unit |
| **Application** | 95% | Unit |
| **Persistence** | 80% | Integration |
| **Interfaces** | 80% | Unit + Integration |
| **Infrastructure** | 70% | Integration |

#### Métricas de Calidad (Post-MVP):

- **Coverage total del proyecto:** ≥ 90%
- **Branch coverage:** ≥ 85%
- **Mutation testing:** (futuro) ≥ 80%

#### Enforcement:

- **Durante el MVP:** El pipeline de CI/CD alertará (warn) si no se cumplen los objetivos del MVP, pero no bloqueará el merge. El cumplimiento se valida como criterio de éxito de la fase.
- **Post-MVP:** El pipeline de CI/CD rechazará (fail) PRs con coverage por debajo de los objetivos ideales.
- Dashboard de coverage en README (Codecov) se usará en todas las fases.
- Revisión mensual de métricas.

### 3. Buenas Prácticas

**DO:**
- ✅ Tests legibles como documentación
- ✅ Nombres descriptivos: `TestCreateOrganization_WhenInvalidEmail_ReturnsError`
- ✅ Arrange-Act-Assert pattern
- ✅ Tests aislados (sin dependencias entre tests)

**DON'T:**
- ❌ Tests que dependen de orden de ejecución
- ❌ Tests que modifican estado global
- ❌ Tests solo para aumentar coverage (sin valor)

### 4. Excepciones

Se acepta coverage < 100% en código generado automáticamente o en fases de refactorización gradual, siempre que se documente en `technical-debt.md`.

---

## Consecuencias

**Positivas:**
- ✅ Alta confiabilidad del código
- ✅ Refactoring seguro
- ✅ Detección temprana de bugs

**Negativas:**
- ⚠️ Tiempo adicional de desarrollo (20-30%)
- ⚠️ Curva de aprendizaje para testear correctamente

---

## Referencias

- [Testing in Go - Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Vitest Documentation](https://vitest.dev/)
- [Testcontainers](https://www.testcontainers.org/)
