# SESSION INITIALIZATION TEMPLATE

**Copia este archivo y rellena al inicio de cada nueva sesión**

---

## 📋 INFORMACIÓN DE LA SESIÓN

| Campo | Valor |
|-------|-------|
| **Sesión** | Session-XX (completa: 2026-MM-DD-session-XX) |
| **Facilitador/LLM** | GitHub Copilot / Claude Anthropic / otro |
| **Fecha Inicio** | YYYY-MM-DD |
| **Hora Inicio** | HH:MM |
| **Duración Estimada** | X horas |
| **Participantes** | Jorge Cortés Villalba, [LLM] |

---

## 🎯 OBJETIVOS PRINCIPALES

**Enumera los 3-5 objetivos principales de esta sesión:**

1. [ ] **Objetivo 1:** [Descripción clara]
   - Subtarea 1a
   - Subtarea 1b
   - Subtarea 1c

2. [ ] **Objetivo 2:** [Descripción clara]
   - Subtarea 2a
   - Subtarea 2b

3. [ ] **Objetivo 3:** [Descripción clara]
   - Subtarea 3a

4. [ ] **Objetivo 4 (si aplica):** [Descripción clara]

5. [ ] **Objetivo 5 (si aplica):** [Descripción clara]

---

## 📊 CONTEXTO DE ENTRADA

### Estado Anterior

**Última sesión completada:** 2026-01-11-session-08

**Cambios desde última sesión:**
- [Cambio 1]
- [Cambio 2]
- [Cambio 3]

**Estado en PROJECT_STATUS.md:**
- Fase actual: [Fase X]
- Horas invertidas: [X / 782]
- Porcentaje completado: [X%]

### Bloqueadores/Dependencias

- [ ] Bloqueador 1: [Descripción] (Impacto: Alto/Medio/Bajo)
- [ ] Dependencia 1: [Qué depende de quién]
- [ ] Riesgo 1: [Qué podría salir mal]

### Prioridades Esta Sesión

**Crítica (Must Have):**
- Prioridad crítica 1
- Prioridad crítica 2

**Alta (Should Have):**
- Prioridad alta 1

**Media (Nice to Have):**
- Prioridad media 1

---

## 🛠️ PLAN DE TRABAJO

### Fase 1: Análisis Arquitectónico (30 min)

- [ ] Identificar Bounded Contexts involucrados
- [ ] Mapear dependencias con otros módulos
- [ ] Revisar ADRs aplicables
- [ ] Diseñar interfaces de dominio (si es new module)
- [ ] Confirmar estructura con usuario

**Notas:**
```
[Espacio para notas de análisis]
```

### Fase 2: Tests (TDD-First)

**Backend – Tests Dominio (30-60 min):**
- [ ] Test cases para [módulo/funcionalidad]
- [ ] Tests de validación en domain entities
- [ ] Tests de integration layer (repositories)
- [ ] Fixtures de datos

**Frontend – Tests (30-60 min si aplica):**
- [ ] Tests Vitest para composables
- [ ] Tests para stores (Pinia)
- [ ] Tests para componentes críticos

**Notas:**
```
[Tests a escribir]
```

### Fase 3: Implementación

**Backend (X horas):**
- [ ] Implementar domain entities/services
- [ ] Implementar application use cases
- [ ] Implementar infrastructure repositories
- [ ] Implementar HTTP handlers/controllers
- [ ] Integración con PostgreSQL (migraciones si es needed)

**Frontend (X horas si aplica):**
- [ ] Stores (Pinia)
- [ ] Composables
- [ ] Componentes
- [ ] Services (API clients)
- [ ] Views/Pages

**Database (X horas si aplica):**
- [ ] Diseñar schema
- [ ] Crear migraciones
- [ ] Scripts iniciales (fixtures)

**Notas de implementación:**
```
[Decisiones de implementación]
```

### Fase 4: Validación

- [ ] `go test ./...` pasa (backend)
- [ ] `npm run test` pasa (frontend si aplica)
- [ ] `golangci-lint run ./...` sin warnings
- [ ] `npm run lint` sin warnings (si aplica)
- [ ] Cobertura de tests ≥75% en módulos críticos
- [ ] Docker Compose levanta sin errores
- [ ] API funcional (test manual con curl/Postman si aplica)
- [ ] Frontend carga y funciona (manual test)

### Fase 5: Documentación

- [ ] Completar esta sesión en `/docs/sessions/2026-MM-DD-session-XX.md`
- [ ] Actualizar PROJECT_STATUS.md
- [ ] Commits descriptivos (see section below)
- [ ] Si es decisión arquitectónica → crear/actualizar ADR (docs/adr/)

---

## 📝 CHANGES MADE

### Commits Realizados

**Formato: `[TYPE]: Brief description`**

Types: feat (feature), fix (bug fix), refactor (restructuring), docs (documentation), test (tests), chore (configuration)

```
[feat]: Implement pricing domain entities and value objects
  - Created PricingService interface in domain
  - Implemented CalculatePrice method with discount logic
  - Added comprehensive unit tests (100% coverage)
  
[test]: Add pricing service integration tests
  - Test pricing with various discount scenarios
  - Test repository interactions
  
[docs]: Document pricing module architecture in ADR-XXX
  - Explain design decisions
  - Reference bounded context mapping
```

**Lista de commits:**
1. [ ] Commit 1: [hash] - [mensaje]
2. [ ] Commit 2: [hash] - [mensaje]
3. [ ] Commit 3: [hash] - [mensaje]
4. [ ] Commit 4: [hash] - [mensaje]

### Archivos Modificados

| Archivo | Tipo | Descripción |
|---------|------|------------|
| `backend/internal/domain/pricing/service.go` | NEW | Tarificación domain service |
| `backend/internal/application/pricing/use_case.go` | NEW | Pricing use case |
| `backend/internal/infrastructure/repository/pricing_repo.go` | NEW | PostgreSQL implementation |
| `backend/tests/unit/domain/pricing_test.go` | NEW | Domain tests |
| `docs/sessions/2026-MM-DD-session-XX.md` | NEW | This session documentation |
| `PROJECT_STATUS.md` | MODIFIED | Updated progress metrics |

### Métricas de Cambio

| Métrica | Valor |
|---------|-------|
| **Archivos creados** | X |
| **Archivos modificados** | X |
| **Líneas de código agregadas** | ~XXX |
| **Tests agregadas** | ~XX |
| **Commits** | X |

---

## ✅ DEFINICIÓN DE "HECHO"

La sesión se considera completada cuando:

- [x] Todos los objetivos están marcados como completados
- [x] Tests backend pasan: `go test ./...`
- [x] Tests frontend pasan: `npm run test` (si aplica)
- [x] Lint sin warnings: `golangci-lint run ./...` y `npm run lint`
- [x] Cobertura ≥75% en dominio crítico
- [x] Docker Compose levanta sin errores
- [x] Sesión documentada completamente
- [x] PROJECT_STATUS.md actualizado
- [x] Commits descriptivos y bien organizados
- [x] No hay breaking changes no documentados

---

## 🚨 BLOQUEADORES/PROBLEMAS ENCONTRADOS

### Durante la Sesión

**Problema 1:** [Descripción]
- **Impacto:** [Qué afectó]
- **Solución:** [Cómo se resolvió]
- **Tiempo invertido:** [X minutos]
- **Prevención futura:** [Qué se aprendió]

**Problema 2:** [Descripción]
- **Impacto:** [Qué afectó]
- **Solución:** [Cómo se resolvió]
- **Tiempo invertido:** [X minutos]

### Deuda Técnica Identificada

- [ ] Deuda 1: [Descripción] → Ticket: [Link o nota para Post-MVP]
- [ ] Deuda 2: [Descripción] → Ticket: [Link o nota para Post-MVP]

---

## 📚 DECISIONES ARQUITECTÓNICAS TOMADAS

### Decisión 1

**Contexto:** [Por qué tomamos esta decisión]

**Alternativas consideradas:**
- Opción A: [Descripción]
- Opción B: [Descripción]

**Decisión adoptada:** Opción X

**Justificación:** [Por qué Opción X]

**Referencia:** [ADR-XXX si aplica, o notas técnicas]

---

## 🎓 APRENDIZAJES/NOTES TÉCNICOS

```
[Notas sobre lo aprendido, patrones descubiertos, gotchas de la tecnología, etc.]

Ejemplo:
- PostgreSQL JSONB es más rápido que columnas separadas para variantes
- Vue 3 Composition API requiere cuidado con ciclo de vida
- GORM hooks tienen quirks con transacciones
```

---

## 📊 MÉTRICAS FINALES

| Métrica | Valor | Target | Status |
|---------|-------|--------|--------|
| **Horas invertidas** | X | X | ✓ |
| **Porcentaje proyecto** | X% | — | — |
| **Tests cobertura** | X% | ≥75% crítico | ✓/✗ |
| **Commits** | X | 1+ | ✓ |
| **Líneas código** | +XXX | — | — |
| **Documentación actualizada** | ✓ | ✓ | ✓ |

---

## 🚀 PRÓXIMOS PASOS (Session N+1)

**Qué continúa la próxima sesión:**

1. [ ] Paso 1: [Descripción]
   - Subtarea 1a
   - Subtarea 1b

2. [ ] Paso 2: [Descripción]
   - Subtarea 2a

3. [ ] Paso 3: [Descripción]

**Prerequisitos para próxima sesión:**
- [ ] [Prerequisito 1]
- [ ] [Prerequisito 2]

**Configuración a tener lista:**
```
[Configuración Docker, variables de entorno, etc.]
```

---

## 📝 NOTES FINALES

```
[Cualquier nota adicional, contexto para futuras sesiones, recomendaciones, etc.]
```

---

## ✍️ FIRMA

**Sesión completada:** [YYYY-MM-DD HH:MM]

**Facilitador:** [Nombre]  
**LLM:** [GitHub Copilot / Claude / otro]  
**Revisor:** [Nombre si aplica]

---

**Para documentación completa de sesión, ver:** `/docs/sessions/_SESSION_TEMPLATE.md`

