# Tarea 03: Estrategia de Calidad y Registro de Deuda Técnica

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 01 |
| **ID de Sprint** | sprint-03 |
| **Título** | Estrategia de Calidad y Registro de Deuda Técnica |
| **Estado** | ⏳ Planificado |
| **Facilitador/LLM** | GitHub Copilot (Claude Sonnet 4.5) |
| **Fecha de Inicio** | (Por determinar) |
| **Fecha de Fin** | (Por determinar) |
| **Duración Estimada** | 2-4 horas |
| **Duración Real** | (Completar al finalizar) |

**Nota sobre IDs:**
- **ID de Tarea**: 01 (primera tarea del sprint-03)
- **ID de Sprint**: sprint-03
| **ID Único** | 03-01 |

---

## 🎯 OBJETIVOS PRINCIPALES

Documentar formalmente la estrategia de calidad del proyecto y establecer un registro sistemático para gestionar la deuda técnica, asegurando que las mejoras futuras estén priorizadas y trackeadas.

### Subtareas

1. [ ] **ADR-010: Estrategia de Testing y Coverage** (1 hora)
   - [ ] Documentar política de coverage mínimo (90%+)
   - [ ] Definir niveles de testing (unit, integration, e2e)
   - [ ] Establecer responsabilidades por capa arquitectónica
   - [ ] Especificar herramientas y frameworks

- [ ] **Registro de Deuda Técnica** (1 hora)
   - [ ] Crear documento central: `/docs/guides/developer/technical-debt.md`
   - [ ] Definir template para registrar deuda
   - [ ] Clasificación (seguridad, rendimiento, mantenibilidad, escalabilidad)
   - [ ] Migrar hallazgos OWASP de baja prioridad
   - [ ] Priorización y estimación

3. [ ] **Guía de Contribución** (1-2 horas)
   - [ ] Crear `CONTRIBUTING.md` en raíz
   - [ ] Documentar proceso de PR
   - [ ] Estándares de código (Go y Vue)
   - [ ] Checklist de revisión
   - [ ] Workflow de desarrollo

---

## 📊 CONTEXTO DE ENTRADA

### Estado Actual de Documentación

**Existente:**
- ✅ ADR-001 a ADR-009: Arquitectura y decisiones estratégicas
- ✅ Módulos documentados: IAM, Party
- ✅ README.md básico

**Faltante:**
- ❌ ADR sobre estrategia de testing
- ❌ Registro sistemático de deuda técnica
- ❌ Guía de contribución formal
- ❌ Políticas de calidad documentadas

---

## 🛠️ PLAN DE TRABAJO

### Fase 1: ADR-010 - Estrategia de Testing y Coverage

**Archivo a crear:** `docs/architecture/adrs/ADR-011-estrategia-testing-coverage.md`

```markdown
# ADR-010 – Estrategia de Testing y Coverage

**Fecha:** 29/01/2026  
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

Se establece la siguiente **estrategia de testing multinivel** con políticas de coverage por capa.

### 1. Niveles de Testing

#### 1.1 Unit Tests (Obligatorios)

**Alcance:**
- Domain layer (entities, value objects, aggregates)
- Application layer (use cases, command/query handlers)
- Casos edge y validaciones

**Coverage objetivo:** 100% en domain, 95%+ en application

**Herramientas:**
- Backend: Go testing package + testify
- Frontend: Vitest + Vue Test Utils

**Ejemplos:**
```go
// Test de Value Object
func TestEmail_MustBeValid(t *testing.T) {
    _, err := NewEmail("invalid")
    assert.Error(t, err)
}

// Test de Use Case
func TestCreateOrganization_Success(t *testing.T) {
    repo := NewInMemoryRepo()
    handler := NewCreateOrganizationHandler(repo)
    
    result, err := handler.Handle(ctx, command)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

#### 1.2 Integration Tests (Selectivos)

**Alcance:**
- Persistence layer con base de datos real
- Infrastructure layer (middleware, adapters)
- Casos que requieren múltiples capas

**Coverage objetivo:** 80%+ en persistence critical paths

**Herramientas:**
- Backend: Testcontainers (PostgreSQL)
- Frontend: E2E con Playwright (selectivo)

**Ejemplos:**
```go
func TestPostgreSQLOrganizationRepository_Create(t *testing.T) {
    db := setupTestDB(t) // Testcontainer
    defer teardownTestDB(t, db)
    
    repo := NewPostgreSQLOrganizationRepository(db)
    org, err := repo.Create(ctx, organization)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, org.ID)
}
```

#### 1.3 E2E Tests (Críticos)

**Alcance:**
- Flujos de usuario completos
- Casos de uso principales del MVP
- Smoke tests de features core

**Coverage objetivo:** 100% de flujos críticos de MVP

**Herramientas:**
- Frontend: Playwright
- API: Postman/Newman (colecciones)

**Flujos críticos:**
- Login + navegación
- Crear organización completa (con contactos/direcciones)
- Búsqueda y filtrado

### 2. Política de Coverage

#### Mínimos Obligatorios por Capa:

| Capa | Coverage Mínimo | Tipo de Tests |
|------|-----------------|---------------|
| **Domain** | 100% | Unit |
| **Application** | 95% | Unit |
| **Persistence** | 80% | Integration |
| **Interfaces** | 80% | Unit + Integration |
| **Infrastructure** | 70% | Integration |

#### Métricas de Calidad:

- **Coverage total del proyecto:** ≥ 90%
- **Branch coverage:** ≥ 85%
- **Mutation testing:** (futuro) ≥ 80%

#### Enforcement:

- CI/CD rechaza PRs con coverage < mínimo
- Dashboard de coverage en README (Codecov)
- Revisión mensual de métricas

### 3. Responsabilidades

#### Por Rol:

**Desarrollador:**
- Escribir tests antes de implementar (TDD)
- Mantener coverage al crear/modificar código
- Documentar casos edge en tests

**Revisor de PR:**
- Verificar tests adecuados
- Validar coverage no disminuye
- Revisar calidad de tests (no solo cantidad)

**CI/CD:**
- Ejecutar todos los tests automáticamente
- Reportar coverage
- Bloquear merge si coverage < mínimo

### 4. Herramientas y Frameworks

#### Backend (Go):

```go
// Frameworks
- testing (stdlib)
- testify/assert
- testify/mock
- testify/suite

// Coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Frontend (Vue):

```javascript
// Frameworks
- Vitest
- Vue Test Utils
- Testing Library

// Coverage
npm run test:unit -- --coverage
```

#### CI Integration:

- GitHub Actions ejecuta tests en cada PR
- Codecov para tracking de coverage
- Pre-commit hooks ejecutan tests rápidos localmente

### 5. Buenas Prácticas

**DO:**
- ✅ Tests legibles como documentación
- ✅ Nombres descriptivos: `TestCreateOrganization_WhenInvalidEmail_ReturnsError`
- ✅ Arrange-Act-Assert pattern
- ✅ Tests aislados (sin dependencias entre tests)
- ✅ Mocks solo cuando necesario (prefer real objects)

**DON'T:**
- ❌ Tests que dependen de orden de ejecución
- ❌ Tests que modifican estado global
- ❌ Tests con sleeps o timeouts arbitrarios
- ❌ Tests solo para aumentar coverage (sin valor)
- ❌ Mockar todo (prefer integration tests)

### 6. Excepciones

**Se acepta coverage < 100% en:**
- Código generado automáticamente
- Main functions (inicialización)
- Código legacy en refactorización gradual

**Proceso:**
- Documentar razón en comentario
- Aprobar en code review
- Registrar en technical-debt.md

---

## Consecuencias

**Positivas:**
- ✅ Alta confiabilidad del código
- ✅ Refactoring seguro
- ✅ Tests como documentación viva
- ✅ Detección temprana de bugs

**Negativas:**
- ⚠️ Tiempo adicional de desarrollo (20-30%)
- ⚠️ Curva de aprendizaje para testear correctamente

**Mitigación:**
- Pair programming en primeros módulos
- Templates de tests como ejemplos
- CI automatizado reduce fricción

---

## Alternativas Consideradas

1. **Sin política formal de coverage**
   - Rechazado: Sin garantías de calidad

2. **100% coverage obligatorio**
   - Rechazado: Demasiado estricto, puede llevar a tests sin valor

3. **Solo integration tests**
   - Rechazado: Lentos, dificultan debug, no documentan dominio

---

## Referencias

- [Testing in Go - Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Vitest Documentation](https://vitest.dev/)
- [Testcontainers](https://www.testcontainers.org/)

---

**Última actualización:** 2026-01-29
```

---

### Fase 2: Registro de Deuda Técnica

**Archivo a crear:** `../../1_project/technical-debt.md`

```markdown
# 📋 Registro de Deuda Técnica - TramaTex

**Última actualización:** 2026-01-29

---

## 🎯 Propósito

Este documento registra y prioriza la deuda técnica del proyecto TramaTex, clasificando mejoras pendientes y estableciendo un plan de resolución.

---

## 📊 Clasificación de Deuda

### Categorías:

1. **Seguridad** 🔴 - Vulnerabilidades o riesgos de seguridad
2. **Rendimiento** 🟠 - Optimizaciones de velocidad o recursos
3. **Mantenibilidad** 🟡 - Código difícil de mantener o entender
4. **Escalabilidad** 🟢 - Limitaciones para crecimiento futuro
5. **Funcionalidad** 🔵 - Features incompletas o workarounds

### Prioridades:

- **P0 - Crítica**: Bloquea desarrollo o compromete sistema
- **P1 - Alta**: Debe resolverse pronto (siguiente sprint)
- **P2 - Media**: Planificar en 1-2 meses
- **P3 - Baja**: Backlog (nice-to-have)

---

## 🔴 Deuda de Seguridad

### [P1] Rate Limiting en Endpoints de Login

**Categoría:** Seguridad  
**Origen:** Auditoría OWASP (A04 - Insecure Design)  
**Descripción:** No hay rate limiting en `/api/iam/login`, permitiendo ataques de fuerza bruta.  
**Impacto:** Medio - Riesgo de compromiso de cuentas  
**Estimación:** 2 horas  
**Propuesta:**
- Implementar rate limiting con middleware (gin-limiter)
- Configurar límite: 5 intentos/minuto por IP
- Backoff exponencial tras múltiples fallos

**Referencias:**
- [auditoria-seguridad-owasp.md](../../archive/sprints/sprint-01/04-auditoria-seguridad-owasp.md#a04-insecure-design)

---

### [P2] Validación de Complejidad de Contraseñas

**Categoría:** Seguridad  
**Origen:** Auditoría OWASP (A07 - Auth Failures)  
**Descripción:** Política de contraseñas básica (solo longitud mínima 8).  
**Impacto:** Medio - Usuarios pueden usar contraseñas débiles  
**Estimación:** 2 horas  
**Propuesta:**
- Agregar validación de complejidad (mayúsculas, números, símbolos)
- Integrar con lista de contraseñas comprometidas (haveibeenpwned API)
- Warning en UI para contraseñas débiles

**Estado:** Documentado, pospuesto para Post-MVP

---

### [P2] Sistema de Recuperación de Contraseñas

**Categoría:** Seguridad + Funcionalidad  
**Origen:** Auditoría OWASP (A07 - Auth Failures)  
**Descripción:** No existe flujo de "Forgot Password".  
**Impacto:** Medio - Dependencia del administrador para resets  
**Estimación:** 4 horas  
**Propuesta:**
- Implementar flujo de recuperación por email
- Tokens temporales con expiración (15 minutos)
- Verificación de identidad antes de reset

**Estado:** Feature pendiente, planificado para Post-MVP

---

### [P3] Bloqueo de Cuenta tras Intentos Fallidos

**Categoría:** Seguridad  
**Origen:** Auditoría OWASP (A07 - Auth Failures)  
**Descripción:** No hay lockout tras múltiples intentos fallidos.  
**Impacto:** Bajo - Facilita ataques de fuerza bruta (mitigado por uso interno)  
**Estimación:** 3 horas  
**Propuesta:**
- Contador de intentos fallidos (Redis o base de datos)
- Bloqueo temporal: 5 intentos = 15 min lockout
- Notificación al usuario de intentos sospechosos

**Estado:** Riesgo aceptado para MVP, resolver antes de producción pública

---

## 🟠 Deuda de Rendimiento

### [P2] Paginación en Endpoints de Listado

**Categoría:** Rendimiento + Escalabilidad  
**Origen:** Implementación inicial del módulo Party  
**Descripción:** Endpoints como `/organizations` retornan todos los registros sin paginación.  
**Impacto:** Bajo ahora, Alto con >1000 organizaciones  
**Estimación:** 3 horas  
**Propuesta:**
- Implementar cursor-based pagination
- Parámetros: `?limit=50&cursor=xyz`
- Documentar en API

**Estado:** Planificado para cuando se necesite (>100 organizaciones)

---

### [P3] Caché de Consultas Frecuentes

**Categoría:** Rendimiento  
**Origen:** Análisis de performance futuro  
**Descripción:** Sin caché para queries repetitivos (ej: lista de organizaciones).  
**Impacto:** Bajo - No es bottleneck actual  
**Estimación:** 6 horas  
**Propuesta:**
- Redis como capa de caché
- Caché para listados con TTL 5 minutos
- Invalidación al modificar datos

**Estado:** Mejora futura (Post-MVP)

---

## 🟡 Deuda de Mantenibilidad

### [P1] Centralizar Manejo de Errores

**Categoría:** Mantenibilidad  
**Origen:** Código actual distribuido  
**Descripción:** Cada handler maneja errores de forma diferente.  
**Impacto:** Medio - Inconsistencias, dificulta debugging  
**Estimación:** 2 horas  
**Propuesta:**
- Error handler centralizado en `pkg/errors/`
- Tipos de error estándar (NotFound, Validation, Internal)
- Logging automático de errores

**Estado:** En progreso (Tarea 05-01)

---

### [P2] Refactorizar Validaciones Duplicadas

**Categoría:** Mantenibilidad  
**Origen:** Implementación del módulo Party  
**Descripción:** Validaciones de email, phone repetidas en múltiples lugares.  
**Impacto:** Bajo - Funciona pero es redundante  
**Estimación:** 1 hora  
**Propuesta:**
- Centralizar en value objects
- Reusar en todas las capas
- Tests únicos

**Estado:** Mejora futura

---

## 🟢 Deuda de Escalabilidad

### [P2] Implementar Gestión de Sesiones con Redis

**Categoría:** Escalabilidad + Seguridad  
**Origen:** Auditoría OWASP (A04 - Insecure Design)  
**Descripción:** Tokens JWT sin revocación. No se pueden invalidar sesiones.  
**Impacto:** Medio - Tokens robados no se pueden revocar  
**Estimación:** 4 horas  
**Propuesta:**
- Redis para sesiones activas
- Blacklist de tokens revocados
- Límite de sesiones simultáneas por usuario

**Estado:** Planificado para Post-MVP

---

### [P3] Migrar a Microservicios (Futuro Lejano)

**Categoría:** Escalabilidad  
**Origen:** ADR-003 (Monolito Modular)  
**Descripción:** Arquitectura actual es monolito, podría necesitar separación.  
**Impacto:** Bajo - No es necesario para MVP ni primeros años  
**Estimación:** Semanas  
**Propuesta:**
- Evaluar cuando haya >10,000 usuarios
- Bounded contexts ya separados facilitan migración
- Considerar servicios independientes por módulo

**Estado:** Decisión futura, reevaluar en 2027

---

## 🔵 Deuda de Funcionalidad

### [P1] Documentación de API (OpenAPI/Swagger)

**Categoría:** Funcionalidad  
**Origen:** Falta de especificación formal  
**Descripción:** No hay documentación generada automáticamente de la API.  
**Impacto:** Medio - Dificulta integración de frontend  
**Estimación:** 3 horas  
**Propuesta:**
- Generar especificación OpenAPI 3.0
- Usar swaggo para Go
- UI interactiva con Swagger UI

**Estado:** Planificado para Sprint 06 o 07

---

### [P2] Soporte de Internacionalización (i18n)

**Categoría:** Funcionalidad  
**Origen:** Requisito no implementado  
**Descripción:** Aplicación solo en español/inglés hardcodeado.  
**Impacto:** Bajo - No es requerido para MVP  
**Estimación:** 8 horas  
**Propuesta:**
- vue-i18n para frontend
- go-i18n para backend
- Archivos de traducción JSON

**Estado:** Postpuesto para Post-MVP

---

## 📅 Plan de Resolución

### Sprint 03 (Actual)
- ✅ Logging estructurado
- ✅ RBAC y authorization
- ✅ CORS configurable
- ✅ Error handling centralizado

### Sprint 04-05
- [ ] Documentación API (OpenAPI)
- [ ] Rate limiting
- [ ] Paginación en listados

### Q2 2026 (Post-MVP)
- [ ] Recuperación de contraseñas
- [ ] Validación complej idad contraseñas
- [ ] Gestión de sesiones con Redis
- [ ] Caché de queries

### Backlog (Evaluar en Q3 2026)
- [ ] Bloqueo de cuenta por intentos fallidos
- [ ] i18n
- [ ] Microservicios (decisión futura)

---

## 📝 Template para Nueva Deuda

```markdown
### [PX] Título Descriptivo

**Categoría:** [Seguridad|Rendimiento|Mantenibilidad|Escalabilidad|Funcionalidad]  
**Origen:** [Dónde se identificó]  
**Descripción:** [Qué es el problema]  
**Impacto:** [Crítico|Alto|Medio|Bajo] - [Explicación]  
**Estimación:** [Horas]  
**Propuesta:** [Cómo resolverlo]  
**Estado:** [Planificado|En progreso|Postpuesto|Resuelto]  
**Referencias:** [Enlaces]
```

---

## 🔄 Proceso de Actualización

1. **Al identificar nueva deuda:**
   - Documentar aquí inmediatamente
   - Asignar categoría y prioridad
   - Estimar esfuerzo

2. **Revisión mensual:**
   - Evaluar prioridades
   - Planificar resolución
   - Actualizar estado

3. **Al resolver deuda:**
   - Mover a sección "Resuelta"
   - Enlazar PR/commit
   - Documentar learnings

---

**Última revisión:** 2026-01-29  
**Próxima revisión:** 2026-02-29
```

---

### Fase 3: Guía de Contribución

**Archivo a crear:** `CONTRIBUTING.md` (raíz del proyecto)

```markdown
# 🤝 Guía de Contribución - TramaTex

¡Gracias por tu interés en contribuir a TramaTex! Esta guía te ayudará a entender el proceso de desarrollo y los estándares del proyecto.

---

## 📋 Tabla de Contenidos

1. [Código de Conducta](#código-de-conducta)
2. [Cómo Empezar](#cómo-empezar)
3. [Workflow de Desarrollo](#workflow-de-desarrollo)
4. [Estándares de Código](#estándares-de-código)
5. [Testing](#testing)
6. [Proceso de Pull Request](#proceso-de-pull-request)
7. [Checklist de Revisión](#checklist-de-revisión)

---

## 📜 Código de Conducta

Todos los contribuyentes deben adherirse a un comportamiento respetuoso y profesional. Valoramos la diversidad y promovemos un ambiente inclusivo.

---

## 🚀 Cómo Empezar

### 1. Fork y Clone

```bash
# Fork el repositorio en GitHub
# Luego clona tu fork
git clone https://github.com/tu-usuario/TramaTex.git
cd TramaTex
```

### 2. Configurar Entorno

**Backend (Go):**
```bash
cd apps/tramatex-api
go mod download
make test
```

**Frontend (Vue):**
```bash
cd apps/frontend
npm install
npm run test:unit
```

### 3. Instalar Pre-commit Hooks

```bash
# Linux/Mac
./scripts/setup-pre-commit.sh

# Windows
.\scripts\setup-pre-commit.ps1
```

---

## 🔄 Workflow de Desarrollo

### 1. Crear Branch

```bash
git checkout -b feature/nombre-descriptivo
# o
git checkout -b fix/descripcion-del-bug
```

**Convención de nombres:**
- `feature/` - Nueva funcionalidad
- `fix/` - Corrección de bugs
- `refactor/` - Refactorización de código
- `docs/` - Cambios en documentación
- `test/` - Añadir o mejorar tests

### 2. Desarrollar con TDD

```bash
# 1. Escribir test que falla
# 2. Implementar código mínimo para pasar
# 3. Refactorizar
# 4. Repetir
```

### 3. Commit Frecuente

```bash
git add .
git commit -m "feat: agregar validación de email"
```

**Formato de commits (Conventional Commits):**
```
tipo(alcance): descripción breve

[cuerpo opcional con más detalles]

[pie opcional con referencias]
```

**Tipos:**
- `feat`: Nueva funcionalidad
- `fix`: Corrección de bug
- `refactor`: Refactorización
- `test`: Añadir/modificar tests
- `docs`: Documentación
- `chore`: Tareas de mantenimiento

**Ejemplos:**
```
feat(party): implementar CRUD de organizaciones
fix(iam): corregir validación de contraseña
refactor(domain): extraer lógica de validación a value object
test(party): agregar tests de integración para repositorios
docs(adr): documentar estrategia de testing
```

---

## 📐 Estándares de Código

### Backend (Go)

**Convenciones:**
- Seguir [Effective Go](https://go.dev/doc/effective_go)
- Nombres descriptivos, no abreviaciones
- Funciones cortas (<50 líneas)
- Comentarios en inglés

**Estructura:**
```go
// Package comment
package domain

import (
	"errors"
	"strings"
)

// Email representa un email válido como Value Object
type Email struct {
	value string
}

// NewEmail crea un nuevo Email validando el formato
func NewEmail(email string) (Email, error) {
	if !isValidEmail(email) {
		return Email{}, errors.New("invalid email format")
	}
	return Email{value: strings.ToLower(email)}, nil
}
```

**Linting:**
```bash
make lint       # Ejecutar golangci-lint
make lint-fix   # Auto-fix issues
```

### Frontend (Vue 3)

**Convenciones:**
- Composition API (no Options API)
- TypeScript para nuevos componentes
- Props y emits tipados
- Nombres descriptivos

**Estructura de componente:**
```vue
<script setup lang="ts">
import { ref, computed } from 'vue'

interface Props {
  title: string
  count?: number
}

const props = withDefaults(defineProps<Props>(), {
  count: 0
})

const emit = defineEmits<{
  update: [value: number]
}>()

const localCount = ref(props.count)

const doubleCount = computed(() => localCount.value * 2)

function increment() {
  localCount.value++
  emit('update', localCount.value)
}
</script>

<template>
  <div>
    <h2>{{ title }}</h2>
    <p>Count: {{ localCount }} (Double: {{ doubleCount }})</p>
    <button @click="increment">Increment</button>
  </div>
</template>

<style scoped>
/* Estilos locales */
</style>
```

**Linting:**
```bash
npm run lint    # ESLint + Prettier
npm run format  # Solo Prettier
```

---

## 🧪 Testing

### Política de Coverage

Ver [ADR-011](../../architecture/adrs/ADR-011-estrategia-testing-coverage.md) para detalles completos.

**Mínimos:**
- Domain: 100%
- Application: 95%
- Persistence: 80%
- Interfaces: 80%

### Escribir Tests

**Backend:**
```go
func TestCreateOrganization_Success(t *testing.T) {
	// Arrange
	repo := NewInMemoryOrganizationRepository()
	handler := NewCreateOrganizationHandler(repo)
	cmd := CreateOrganizationCommand{
		Name: "Test Org",
		TaxID: "12345678A",
	}
	
	// Act
	result, err := handler.Handle(context.Background(), cmd)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Org", result.Name)
}
```

**Frontend:**
```javascript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import OrganizationCard from '@/components/OrganizationCard.vue'

describe('OrganizationCard', () => {
  it('renders organization name', () => {
    const wrapper = mount(OrganizationCard, {
      props: {
        organization: {
          id: '1',
          name: 'Test Org'
        }
      }
    })
    
    expect(wrapper.text()).toContain('Test Org')
  })
})
```

### Ejecutar Tests

```bash
# Backend
make test
make coverage

# Frontend
npm run test:unit
npm run test:unit -- --coverage
```

---

## 🔀 Proceso de Pull Request

### 1. Antes de Crear el PR

- [ ] Tests pasan localmente
- [ ] Linters sin errores
- [ ] Coverage no disminuye
- [ ] Commits limpios y descriptivos
- [ ] Branch actualizado con main

```bash
git fetch origin
git rebase origin/main
```

### 2. Crear el PR

**Título:**
```
feat(party): implementar CRUD de organizaciones
```

**Descripción:**
```markdown
## 🎯 Objetivo
Implementar operaciones CRUD para el módulo Party - Organizations.

## 📝 Cambios
- ✅ Agregar CreateOrganizationHandler
- ✅ Agregar UpdateOrganizationHandler
- ✅ Tests unitarios (100% coverage)
- ✅ Documentación actualizada

## ✅ Checklist
- [x] Tests pasan
- [x] Coverage ≥ 90%
- [x] Linters sin errores
- [x] Documentación actualizada
- [x] No hay breaking changes

## 📚 Referencias
- Issue #42
- ADR-006
```

### 3. Code Review

- Responder a comentarios constructivamente
- Hacer cambios solicitados en commits nuevos
- No hacer force-push durante review
- Pedir clarificación si algo no es claro

### 4. Merge

- El revisor hará merge cuando esté aprobado
- Se usará "Squash and Merge" para commits limpios
- El branch será eliminado automáticamente

---

## ✅ Checklist de Revisión

### Para el Autor

Antes de enviar el PR:

**Código:**
- [ ] Sigue los estándares de estilo del proyecto
- [ ] Nombres descriptivos y claros
- [ ] Sin código comentado o debug logs
- [ ] Sin secretos o credenciales

**Tests:**
- [ ] Todos los tests pasan
- [ ] Coverage ≥ mínimos establecidos
- [ ] Tests nuevos para nueva funcionalidad
- [ ] Tests legibles y bien nombrados

**Documentación:**
- [ ] README actualizado si es necesario
- [ ] ADR creado para decisiones arquitectónicas
- [ ] Comentarios en código complejo
- [ ] CHANGELOG actualizado (si aplica)

**CI/CD:**
- [ ] GitHub Actions pasa
- [ ] Pre-commit hooks configurados
- [ ] Sin warnings de linters

### Para el Revisor

Al revisar el PR:

**Arquitectura:**
- [ ] Sigue Clean Architecture + DDD
- [ ] Separación correcta de capas
- [ ] Domain independiente de infraestructura
- [ ] Usa casos de uso correctos

**Calidad:**
- [ ] Código legible y mantenible
- [ ] Sin duplicación innecesaria
- [ ] Manejo de errores adecuado
- [ ] Logging apropiado

**Seguridad:**
- [ ] Sin vulnerabilidades obvias
- [ ] Validación de inputs
- [ ] Autorización verificada
- [ ] Datos sensibles protegidos

**Tests:**
- [ ] Coverage adecuado
- [ ] Tests significativos (no solo para números)
- [ ] Edge cases cubiertos
- [ ] Tests aislados y determinísticos

---

## 📚 Recursos Adicionales

- [Architecture Overview](../../architecture/README.md)
- [ADRs](../../architecture/adrs/)
- [Technical Debt](../../guides/developer/technical-debt.md)
- [Project Status](../../records/project-status.md)

---

## 💬 ¿Preguntas?

Si tienes dudas:
1. Revisa la documentación en `/docs`
2. Busca en issues cerrados
3. Abre un issue con tu pregunta
4. Contacta al equipo

---

¡Gracias por contribuir a TramaTex! 🎉
```

---

## 🚨 BLOQUEADORES/PROBLEMAS ENCONTRADOS

*(Se actualizará durante la implementación)*

---

## 🎓 APRENDIZAJES Y NOTAS

### Decisiones de Documentación

1. **ADR-010 como fuente única:**
   - Toda la estrategia de testing en un documento
   - Referencias desde otros lugares
   - Evita duplicación

2. **technical-debt.md como registro vivo:**
   - Se actualiza continuamente
   - Proceso claro de revisión mensual
   - Template consistente

3. **CONTRIBUTING.md exhaustivo:**
   - Onboarding completo para nuevos devs
   - Ejemplos concretos
   - Checklists accionables

---

## 📚 REFERENCIAS

- [ADR Template](../../architecture/adrs/_ADR_TEMPLATE.md)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

## ✅ CHECKLIST DE FINALIZACIÓN

- [ ] ADR-010 creado y revisado
- [ ] technical-debt.md creado con deuda actual
- [ ] Hallazgos OWASP migrados a technical-debt
- [ ] CONTRIBUTING.md creado
- [ ] Templates documentados
- [ ] Proceso de PR establecido
- [ ] Documentación revisada por equipo

---

**Última actualización:** 2026-01-29
