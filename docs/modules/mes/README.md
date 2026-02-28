# Manufacturing Execution System (MES) Module Documentation

**Nombre del Módulo:** MES (Manufacturing Execution System)  
**Bounded Context:** Producción personalizada, Planificación y Ejecución de Manufactura  
**Responsabilidad Principal:** Definir plantillas de proceso, planificar trabajos y ejecutar su ciclo operativo con trazabilidad.  
**Entidades Raíz:** MESWorkDefinition, MESWorkExecution, MESTask, MESServiceTemplate, MESPosition  
**Dependencias:** ERP Core (Sales, Product, Party, IAM)  

---

## 1. Especificación del Módulo (module-spec.md)

### Objetivo
El módulo MES tiene como objetivo principal proporcionar un sistema robusto para la gestión completa de las operaciones de fabricación, desde la definición de plantillas de proceso hasta la ejecución real del trabajo en planta. Permite el seguimiento de estado, la asignación de puestos y la trazabilidad operativa en procesos personalizados, **integrándose de manera explícita con los módulos del ERP Core para la obtención de insumos y la notificación de resultados.**

### Alcance
*   Creación y gestión de Trabajos Definidos (`MESWorkDefinition`).
*   Ejecución de Trabajos Reales (`MESWorkExecution`) con tareas generadas automáticamente.
*   Gestión de Plantillas de Proceso MES (`MESServiceTemplate`) y secuencias de tareas.
*   Gestión de Tareas (`MESTask`) y Puestos (`MESPosition`) reutilizables.
*   Integración con `Sales` (del ERP Core) para el origen de trabajos.
*   Integración con `Product` (del ERP Core) para la definición de lo que se va a producir.

### Restricciones
*   Crítico para la operatividad, pero de "criticidad media" para el MVP.
*   Diseñado para ser extraíble como microservicio en el futuro, manteniendo su propio esquema lógico de base de datos y comunicándose con otros módulos del ERP Core a través de interfaces explícitas (ver ADR-018).

---

## 2. Requisitos del Módulo

### Requisitos Funcionales

| ID | Descripción | Prioridad | Fase |
|----|-----------|---------|----|
| RF-MES-001 | El sistema debe permitir crear y gestionar Trabajos Definidos (`MESWorkDefinition`). | Alta | 3 |
| RF-MES-002 | El sistema debe permitir configurar Plantillas de Proceso y asignarlas a trabajos. | Alta | 3 |
| RF-MES-003 | El sistema debe permitir ejecutar Trabajos Reales (`MESWorkExecution`) y registrar avance operativo. | Media | 3 |
| RF-MES-004 | El sistema debe permitir actualizar y seguir el estado de un Trabajo MES. | Alta | 3 |
| RF-MES-005 | El sistema debe consumir información de Órdenes de Venta del módulo `Sales` (ERP Core) para generar trabajos MES. | Alta | 3 |
| RF-MES-006 | El sistema debe consultar información de `Product` (ERP Core) para detalles de los ítems a producir. | Alta | 3 |

### Requisitos No Funcionales

| ID | Descripción | Métrica | Target |
|----|----------|---------|----|
| RNF-MES-001 | La comunicación con `Sales` y `Product` (ERP Core) debe ser a través de APIs bien definidas. | Acoplamiento | Bajo |
| RNF-MES-002 | La base de datos del módulo MES debe ser lógicamente independiente. | Dependencia DB | Nula |
| RNF-MES-003 | El módulo debe estar preparado para una futura extracción a microservicio. | Modularidad | Alta |

**Trazabilidad:** Conectar con [DOCUMENTO-CONSOLIDADO-3.0.md](../consolidated/DOCUMENTO-CONSOLIDADO-3.0.md) (Placeholder)

---

## 3. Historias de Usuario (Pendiente de Definición)

### Formato Estándar

**US-MES-001: [Título breve]**

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
```

---

## 4. Modelo de Dominio (domain-model.md)

### Entidades Principales

```
[DiagramER simplificado]
```

#### MESWorkDefinition
- Responsabilidad: Representa la definición planificada de un trabajo a fabricar.
- Value Objects: MESWorkDefinitionID, MESWorkStatus
- Reglas de Negocio: Debe tener `party_id`, `tangible_group_id` y una estructura mínima de asignaciones.

#### MESServiceTemplate
- Responsabilidad: Define la plantilla de proceso reusable para generar secuencias de trabajo.
- Value Objects: MESServiceTemplateID
- Reglas de Negocio: Debe contener tareas con secuencia válida y sin duplicados de orden.

#### MESWorkExecution
- Responsabilidad: Representa la ejecución real de un trabajo definido.
- Value Objects: MESWorkExecutionID, MESWorkStatus
- Reglas de Negocio: Debe trazarse por tarea/estado y respetar la secuencia operativa.

### Value Objects
- MESWorkDefinitionID: Identificador único de un trabajo definido.
- MESServiceTemplateID: Identificador único de una plantilla de proceso.
- MESWorkExecutionID: Identificador único de una ejecución de trabajo.
- MESWorkStatus: Enumerado (e.g., DRAFT, PENDING, IN_PROGRESS, ON_HOLD, COMPLETED).

### Servicios de Dominio
- WorkPlanningService: Encargado de transformar trabajo definido en plan operativo ejecutable.
- WorkExecutionService: Encargado de gestionar avance, transición de estado y trazabilidad.

---

## 5. Casos de Uso (use-cases.md)

### Caso de Uso 1: Crear Orden de Producción
- **Actor:** Planificador de Producción
- **Precondiciones:** La Orden de Venta existe en el módulo `Sales` (ERP Core). El `Product` y `Variant` existen en el módulo `Product` (ERP Core).
- **Flujo Normal:**
  1. El Planificador inicia la creación de un Trabajo Definido (`MESWorkDefinition`).
  2. El sistema consulta `Sales` (ERP Core) para obtener detalles de la Orden de Venta.
  3. El sistema consulta `Product` (ERP Core) para obtener detalles del Producto/Variante a producir.
  4. El sistema crea un nuevo trabajo MES con estado 'DRAFT'/'PENDING'.
- **Postcondiciones:** Un nuevo trabajo MES ha sido creado y listo para ejecución.

### Nota de compatibilidad de nomenclatura

Los términos `ProductionOrder`, `WorkCenter` y `QualityCheck` se consideran **legado documental** y deben mapearse a:

- `ProductionOrder` → `MESWorkDefinition` / `MESWorkExecution` (según contexto)
- `WorkCenter` → `MESPosition`
- `Service Group` / `Grupo de Servicio` → `MESServiceTemplate` (etiqueta UI: **Plantilla de proceso**)

---

## 6. Contratos de API (api-contracts.md) (Pendiente de Definición)

### Endpoint 1
```json
POST /api/mes/production-orders
Request: { "salesOrderId": "uuid", "productId": "uuid", "quantity": 10 }
Response: { "productionOrderId": "uuid", "status": "Pending" }
Errores: [...]
```

---

## 7. Decisiones Técnicas (Ver ADR-018)

### Base de Datos para MES
**Decisión Tomada:** El módulo MES utilizará su propio esquema lógico en la base de datos compartida, sin acceso directo de otros módulos. (Ver ADR-018)
**Justificación:** Preparación para microservicios y encapsulación del dominio.

---

## 9. Tests y Cobertura (Pendiente de Adaptación para MES)

### Estrategia de Testing

| Capa | Enfoque | Cobertura Target | Herramienta |
|------|---------|-----------------|------------|
| **Dominio** | TDD-first, unit tests | ≥90% | go test |
| **Aplicación** | Integration tests con mocks | ≥80% | go test |
| **Infraestructura** | Integration E2E | ≥60% | go test |
| **HTTP/Handlers** | HTTP tests | ≥70% | go test |

### Proporción de Tests (Regla Práctica Inicial)

**Distribución recomendada de cobertura total:**

| Tipo de Test | Proporción | Descripción | Herramienta |
|--------------|-----------|------------|------------|
| **Unit Tests** | 70% | Pruebas de componentes individuales (dominio, VOs, entidades) | go test |
| **Integration Tests** | 25% | Pruebas de orquestación (use cases, repositorios mockeados) | go test |
| **E2E Tests** | 5% | Pruebas end-to-end solo casos críticos (HTTP handlers principales) | go test |

### Casos de Prueba Críticos (Pendiente de Definición para MES)

---

## 10. Despliegue y Release (Pendiente de Adaptación para MES)

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

- [ ] Definir User Stories específicas para MES.
- [ ] Detallar contratos de API.
- [ ] Desarrollar estrategias de prueba específicas para MES.

---

**Última Actualización:** 07/02/2026
**Responsable:** Gemini CLI (acting as Lead Architect)
**Versión:** 1.0
**Relacionados:** ADR-018, ADR-002, ADR-003, bounded-contexts.yaml
