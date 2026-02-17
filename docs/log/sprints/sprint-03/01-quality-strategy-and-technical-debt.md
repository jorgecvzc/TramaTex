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
**Archivo a crear:** `docs/architecture/adrs/ADR-011-testing-coverage-strategy.md`
**Resumen:** Se documentará formalmente la estrategia de testing del proyecto, incluyendo niveles de testing (unitarios, integración, E2E), política de cobertura por capa, responsabilidades y herramientas/frameworks. Para más detalles, consulte [ADR-011: Estrategia de Testing y Cobertura](../../../architecture/adrs/ADR-011-testing-coverage-strategy.md).

### Fase 2: Registro de Deuda Técnica
**Archivo a crear:** `docs/guides/developer/technical-debt.md`
**Resumen:** Se creará un documento central para registrar y priorizar la deuda técnica, incluyendo su clasificación (seguridad, rendimiento, mantenibilidad, escalabilidad, funcionalidad) y un plan de resolución. Se migrarán los hallazgos de OWASP de baja prioridad a este registro. Para más detalles, consulte [Registro de Deuda Técnica](../../../guides/developer/technical-debt.md).

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

Ver [ADR-011](../../architecture/adrs/ADR-011-testing-coverage-strategy.md) para detalles completos.

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
- [Project Status](../../project-status.md)

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
