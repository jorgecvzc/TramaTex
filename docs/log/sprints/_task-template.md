# PLANTILLA DE TAREA DE SPRINT

**Copia este archivo y rellena al inicio de cada nuevo bloque de trabajo**

## 📝 Convenciones de Nombres

**Archivo de tarea:**
- Formato: `[ID_TAREA]-nombre-descriptivo-kebab-case.md`
- ID de Tarea: Número secuencial local por sprint (01, 02, 03... se reinicia en cada sprint)
- Identificación única: Combinación sprint-XX + tarea-YY
- Ubicación: `docs/log/sprints/sprint-[XX]/`
- Ejemplos:
  - `sprint-01/01-diseño-y-arquitectura-inicial.md` (tarea 01-01)
  - `sprint-01/04-auditoria-seguridad-owasp.md` (tarea 01-04)
  - `sprint-02/01-refactorizacion-sistema-documentacion.md` (tarea 02-01)
  - `sprint-03/01-definicion-sistema-de-diseno.md` (tarea 03-01)

**Anotaciones generales:**
- Los títulos y explicaciones tienen que estar en castellano.
- Sólo existirá un archivo por cada bloque de trabajo/funcionalidad.
- Usar terminología "Tarea" (no "bitácora" ni "journal").
- La numeración de tareas se reinicia en cada sprint (facilita reabrir sprints).

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | XX |
| **ID de Sprint** | sprint-XX |
| **Título** | [Nombre Descriptivo de la Tarea o Funcionalidad] |
| **Estado** | ⏳ En Progreso / ✅ Completado / ❌ Cancelado |
| **Facilitador/LLM** | GitHub Copilot / Claude Anthropic / Gemini |
| **Fecha de Inicio** | YYYY-MM-DD |
| **Fecha de Fin** | YYYY-MM-DD |
| **Duración Estimada** | X horas |
| **Duración Real** | X horas (completar al finalizar) |

**Nota sobre IDs:**
- **ID de Tarea**: Número secuencial local del sprint (01, 02, 03... reinicia cada sprint)
- **ID de Sprint**: Sprint al que pertenece (sprint-01, sprint-02...)
- **ID Único**: La combinación sprint-XX + tarea-YY (ej: 01-04, 02-01, 04-01)

---

## 🎯 OBJETIVOS PRINCIPALES

**Enumera los 3-5 objetivos principales de esta tarea:**

1. [ ] **Objetivo 1:** [Descripción clara]
   - Subtarea 1a

---

## 📊 CONTEXTO DE ENTRADA

### Estado Anterior

**Última tarea completada:** [ID]-[descripcion]

**Cambios desde última tarea:**
- [Cambio 1]
- [Cambio 2]

**Estado en project-status.md:**
- Fase actual: [Fase X]

### Bloqueadores/Dependencias

- [ ] Bloqueador 1: [Descripción] (Impacto: Alto/Medio/Bajo)

---

### Prioridades para esta Tarea

**Crítica (Must Have):**
- Prioridad crítica 1

**Alta (Should Have):**
- Prioridad alta 1

**Media (Nice to Have):
- Prioridad media 1

---



## 🛠️ PLAN DE TRABAJO



### Fase 1: Análisis Arquitectónico (30 min)

- [ ] Identificar Bounded Contexts involucrados
- [ ] Mapear dependencias con otros módulos
- [ ] Revisar ADRs aplicables
- [ ] Diseñar interfaces de dominio (si es new module)
- [ ] Confirmar estructura con usuario



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



### Fase 4: Validación



- [ ] `go test ./...` pasa (backend)

- [ ] `npm run test` pasa (frontend si aplica)

- [ ] `golangci-lint run ./...` sin warnings

- [ ] `npm run lint` sin warnings (si aplica)

- [ ] Cobertura de tests ≥85% en módulos críticos

- [ ] Docker Compose levanta sin errores

- [ ] API funcional (test manual con curl/Postman si aplica)

- [ ] Frontend carga y funciona (manual test)



### Fase 5: Documentación



- [ ] Completar esta tarea en `../sprint-[ID_SPRINT]/[ID_TAREA]-descripcion.md`

- [ ] Actualizar project-status.md

- [ ] Commits descriptivos (see section below)

   - [ ] Si es decisión arquitectónica → crear/actualizar ADR (`../../architecture/adrs/`)

---



## 📝 CHANGES MADE



### Commits Realizados

**Formato: `[TYPE]: Brief description`**

Types: feat (feature), fix (bug fix), refactor (restructuring), docs (documentation), test (tests), chore (configuration)

**Lista de commits:**
1. [ ] Commit 1: [hash] - [mensaje]



### Archivos Modificados

| Archivo | Tipo | Descripción |
|---------|------|------------|
| `apps/component-name/internal/domain/service.go` | NEW | Tarificación domain service |

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

La tarea se considera completada cuando:

- [x] Todos los objetivos están marcados como completados
- [x] Tests backend pasan: `go test ./...`
- [x] Tests frontend pasan: `npm run test` (si aplica)
- [x] Lint sin warnings: `golangci-lint run ./...` y `npm run lint`
- [x] Cobertura ≥85% en dominio crítico
- [x] Docker Compose levanta sin errores
- [x] Tarea documentada completamente
- [x] project-status.md actualizado
- [x] Commits descriptivos y bien organizados
- [x] No hay breaking changes no documentados

---

## 🚨 BLOQUEADORES/PROBLEMAS ENCONTRADOS

### Durante la Tarea

**Problema 1:** [Descripción]
- **Impacto:** [Qué afectó]
- **Solución:** [Cómo se resolvió]
- **Tiempo invertido:** [X minutos]
- **Prevención futura:** [Qué se aprendió]

### Deuda Técnica Identificada

---

## 📚 DECISIONES ARQUITECTÓNICAS TOMADAS

### Decisión 1

**Contexto:** [Por qué tomamos esta decisión]

**Alternativas consideradas:**
- Opción A: [Descripción]

**Decisión adoptada:** Opción X

**Justificación:** [Por qué Opción X]

**Referencia:** [ADR-XXX si aplica, o notas técnicas]

## 🎓 APRENDIZAJES/NOTES TÉCNICOS

```
[Notas sobre lo aprendido, patrones descubiertos, gotchas de la tecnología, etc.]
```

---

## 📊 MÉTRICAS FINALES

| Métrica | Valor | Target | Status |
|---------|-------|--------|--------|
| **Horas invertidas** | X | X | ✓ |
| **Porcentaje proyecto** | X% | — | — |
| **Tests cobertura** | X% | ≥85% crítico | ✓/✗ |
| **Commits** | X | 1+ | ✓ |
| **Líneas código** | +XXX | — | — |
| **Documentación actualizada** | ✓ | ✓ | ✓ |

---

## 🚀 PRÓXIMOS PASOS

**Qué continúa la próxima tarea:**

1. [ ] Paso 1: [Descripción]

**Prerequisitos para próxima tarea:**
- [ ] [Prerequisito 1]

## 📝 NOTES FINALES

```
[Cualquier nota adicional, contexto para futuras tareas, recomendaciones, etc.]
```

---

## ✍️ FIRMA

**Tarea completada:** [YYYY-MM-DD HH:MM]

**Facilitador:** [Nombre]  
**LLM:** [GitHub Copilot / Claude / otro]  
**Revisor:** [Nombre si aplica]

---



