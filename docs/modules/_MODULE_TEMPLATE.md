# Plantilla de Documentación de Módulo

**Nombre del Módulo:** [Nombre]  
**Bounded Context:** [Contexto Delimitado]  
**Responsabilidad Principal:** [Descripción breve]  
**Entidades Raíz:** [Agregados principales]  
**Dependencias:** [Módulos de los que depende]  

---

## 1. Especificación del Módulo (module-spec.md)

### Objetivo
[Describir qué problema resuelve y qué valor proporciona]

### Alcance
[Qué está incluido, qué no]

### Restricciones
[Limitaciones técnicas o de negocio]

---

## 2. Requisitos del Módulo

### Requisitos Funcionales

| ID | Descripción | Prioridad | Fase |
|----|-----------|---------|----|
| RF-[MOD]-001 | [Descripción requisito funcional] | Alta | [Fase] |
| RF-[MOD]-002 | [Descripción requisito funcional] | Media | [Fase] |

### Requisitos No Funcionales

| ID | Descripción | Métrica | Target |
|----|----------|---------|----|
| RNF-[MOD]-001 | [Descripción requisito no funcional] | [Métrica] | [Target] |
| RNF-[MOD]-002 | [Descripción requisito no funcional] | [Métrica] | [Target] |

**Trazabilidad:** Conectar con [DOCUMENTO-CONSOLIDADO-3.0.md](../consolidated/DOCUMENTO-CONSOLIDADO-3.0.md)

---

## 3. Historias de Usuario

### Formato Estándar

**US-[MOD]-001: [Título breve]**

Como [rol de usuario]  
Quiero [objetivo/acción]  
Para [beneficio/razón de negocio]

**Prioridad:** [Alta/Media/Baja]  
**Complejidad:** [S/M/L/XL]  
**Sprint:** [Sprint de asignación]

#### Criterios de Aceptación (BDD)

```gherkin
Scenario 1: [Escenario exitoso]
  Given [precondición 1]
  When [acción principal]
  Then [resultado esperado]
  And [verificación adicional]

Scenario 2: [Escenario alternativo]
  Given [precondición]
  When [acción]
  Then [resultado esperado]

Scenario 3: [Manejo de errores]
  Given [precondición de error]
  When [acción que causa error]
  Then [error manejado correctamente]
```

### Ejemplos de User Stories

**US-[MOD]-001**

Como [rol]  
Quiero [objetivo]  
Para [beneficio]

**Criterios de Aceptación:**
```gherkin
Scenario: [Descripción]
  Given [condición inicial]
  When [acción]
  Then [resultado]
```

---

## 4. Modelo de Dominio (domain-model.md)

### Entidades Principales

```
[DiagramER simplificado]
```

#### [Entidad 1]
- Responsabilidad:
- Value Objects:
- Reglas de Negocio:

#### [Entidad 2]
- ...

### Value Objects
- [VO1]: [Descripción]
- [VO2]: [Descripción]

### Servicios de Dominio
[Si aplica]

---

## 5. Casos de Uso (use-cases.md)

### Caso de Uso 1: [Nombre]
- **Actor:** [Quién lo usa]
- **Precondiciones:** [Qué debe ser verdadero]
- **Flujo Normal:**
  1. ...
  2. ...
- **Flujos Alternativos:** [Si aplica]
- **Postcondiciones:** [Qué debe ser verdadero después]

### Caso de Uso 2: [Nombre]
- ...

---

## 6. Contratos de API (api-contracts.md)

### Endpoint 1
```json
POST /api/[modulo]/[recurso]
Request: { ... }
Response: { ... }
Errores: [...]
```

### Endpoint 2
```
...
```

---

## 7. Decisiones Técnicas

### [Decisión 1]
**Alternativas Consideradas:**
**Decisión Tomada:**
**Justificación:**

---

## 9. Tests y Cobertura

### Estrategia de Testing

| Capa | Enfoque | Cobertura Target | Herramienta |
|------|---------|-----------------|------------|
| **Dominio** | TDD-first, unit tests | ≥90% | go test / vitest |
| **Aplicación** | Integration tests con mocks | ≥80% | go test / vitest |
| **Infraestructura** | Integration E2E | ≥60% | go test / vitest |
| **HTTP/Handlers** | HTTP tests | ≥70% | go test |

### Casos de Prueba Críticos

**Dominio:**
- [ ] [Entidad principal] creation with valid inputs
- [ ] [Entidad principal] rejects invalid inputs
- [ ] [Value Object] validation rules
- [ ] Business rule: [Regla crítica]

**Aplicación:**
- [ ] Use Case: [Caso exitoso]
- [ ] Use Case: [Caso error]
- [ ] Use Case: [Caso alternativo]
- [ ] Repository mock integration

**HTTP:**
- [ ] POST [endpoint] with valid data → 200
- [ ] POST [endpoint] with invalid data → 400
- [ ] POST [endpoint] unauthorized → 401
- [ ] POST [endpoint] error handling → 500

### Línea de Base (Coverage Baseline)

```bash
# Backend (Go)
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Frontend (Vue.js)
npm run test:coverage
# o
vitest --coverage
```

### Reporte de Cobertura

| Componente | Cobertura | Status |
|------------|-----------|--------|
| domain/[module]/ | XX% | ✅/⚠️/❌ |
| application/[module]/ | XX% | ✅/⚠️/❌ |
| infrastructure/[module]/ | XX% | ✅/⚠️/❌ |
| interfaces/http/[module]/ | XX% | ✅/⚠️/❌ |

---

## 10. Despliegue y Release

### Configuración de Despliegue

**Docker:**
- [ ] Dockerfile configurado
- [ ] docker-compose.yml actualizado
- [ ] Environment variables documentadas

**Base de Datos:**
- [ ] Migration scripts presentes
- [ ] Rollback plan documentado
- [ ] Data fixtures para testing

**Configuración:**
- `.env.example` actualizado
- `config/config.go` carga variables necesarias
- JWT_SECRET / secrets manejados correctamente

### Health Check

```
GET /api/[module]/health
Response: { "status": "ok", "version": "1.0.0" }
```

### Checklist Pre-Release

- [ ] Todos los tests pasan: `go test ./...` (100%)
- [ ] Lint sin warnings: `golangci-lint run ./...`
- [ ] Cobertura ≥80%: `go test -cover ./...`
- [ ] Documentación actualizada (spec, ADRs)
- [ ] Docker image build exitoso
- [ ] docker-compose up funciona sin errores
- [ ] Endpoint health check responde
- [ ] Logs estructurados configurados
- [ ] Error handling completo
- [ ] Commit messages descriptivos
- [ ] PR/sesión documentada
- [ ] No hay breaking changes no documentados

### Estrategia de Rollback

**En caso de error post-release:**

1. **Identificar problema:** Logs, monitoring, tests
2. **Revertir cambios:**
   ```bash
   git revert [commit-hash]
   docker-compose down
   docker-compose up -d
   ```
3. **Restaurar BD (si aplica):**
   ```bash
   ./scripts/rollback-migration.sh [version]
   ```
4. **Comunicar:** Documentar en sesión qué falló y por qué

### Monitoreo Post-Deploy

- [ ] Logs sin errores
- [ ] CPU/Memory dentro de límites
- [ ] Latencia API aceptable
- [ ] No hay data corruption
- [ ] Backups completados

---

## 11. Notas y Pendientes

- [ ] Tarea 1
- [ ] Tarea 2

---

**Última Actualización:** [Fecha]  
**Responsable:** [Persona/Equipo]  
**Versión:** 1.0  
**Relacionados:** [ADRs, sesiones, documentos]

