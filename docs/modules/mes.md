# Manufacturing Execution System (MES) Module Documentation

**Nombre del Módulo:** MES (Manufacturing Execution System)  
**Bounded Context:** Producción personalizada, Planificación y Ejecución de Manufactura, Control de Calidad  
**Responsabilidad Principal:** Planificación, ejecución y control de calidad de la producción. Gestión de órdenes de producción y estaciones de trabajo.  
**Entidades Raíz:** ProductionOrder, WorkCenter, QualityCheck  
**Dependencias:** ERP Core (Sales, Product, Party, IAM)  

---

## 1. Especificación del Módulo (module-spec.md)

### Objetivo
El módulo MES tiene como objetivo principal proporcionar un sistema robusto para la gestión completa de las operaciones de fabricación, desde la asignación de órdenes de producción hasta el control de calidad final. Permite el seguimiento del estado de la producción, la asignación a centros de trabajo y la verificación de la calidad, asegurando la eficiencia y la trazabilidad en el proceso de fabricación, especialmente para productos personalizados, **integrándose de manera explícita con los módulos del ERP Core para la obtención de insumos y la notificación de resultados.**

### Alcance
*   Creación y gestión de Órdenes de Producción (`ProductionOrder`).
*   Asignación de Órdenes de Producción a Centros de Trabajo (`WorkCenter`).
*   Seguimiento del estado de las Órdenes de Producción.
*   Registro y gestión de Controles de Calidad (`QualityCheck`).
*   Integración con `Sales` (del ERP Core) para el origen de las órdenes de producción.
*   Integración con `Product` (del ERP Core) para la definición de lo que se va a producir.

### Restricciones
*   Crítico para la operatividad, pero de "criticidad media" para el MVP.
*   Diseñado para ser extraíble como microservicio en el futuro, manteniendo su propio esquema lógico de base de datos y comunicándose con otros módulos del ERP Core a través de interfaces explícitas (ver ADR-018).

---

## 2. Requisitos del Módulo

### Requisitos Funcionales

| ID | Descripción | Prioridad | Fase |
|----|-----------|---------|----|
| RF-MES-001 | El sistema debe permitir crear y gestionar Órdenes de Producción. | Alta | 3 |
| RF-MES-002 | El sistema debe permitir asignar Órdenes de Producción a Centros de Trabajo. | Alta | 3 |
| RF-MES-003 | El sistema debe permitir registrar y consultar Controles de Calidad para una producción. | Media | 3 |
| RF-MES-004 | El sistema debe permitir actualizar y seguir el estado de una Orden de Producción. | Alta | 3 |
| RF-MES-005 | El sistema debe consumir información de Órdenes de Venta del módulo `Sales` (ERP Core) para generar `ProductionOrder`. | Alta | 3 |
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

#### ProductionOrder
- Responsabilidad: Representa una orden de trabajo para fabricar un producto.
- Value Objects: ProductionOrderID, ProductionStatus
- Reglas de Negocio: Debe tener un ProductID, un SalesOrderID de referencia.

#### WorkCenter
- Responsabilidad: Representa una estación o recurso donde se realiza una parte de la producción.
- Value Objects: WorkCenterID
- Reglas de Negocio: Debe tener una capacidad.

#### QualityCheck
- Responsabilidad: Representa un control de calidad realizado sobre una producción.
- Value Objects: QualityCheckID, QualityStatus
- Reglas de Negocio: Debe estar asociado a una ProductionOrder, tener un resultado.

### Value Objects
- ProductionOrderID: Identificador único de una orden de producción.
- WorkCenterID: Identificador único de un centro de trabajo.
- QualityCheckID: Identificador único de un control de calidad.
- ProductionStatus: Enumerado (e.g., Pending, In-Progress, Completed, On-Hold).
- QualityStatus: Enumerado (e.g., Passed, Failed, Rework).

### Servicios de Dominio
- ProductionSchedulerService: Encargado de asignar ProductionOrders a WorkCenters.
- QualityControlService: Encargado de validar y registrar QualityChecks.

---

## 5. Casos de Uso (use-cases.md)

### Caso de Uso 1: Crear Orden de Producción
- **Actor:** Planificador de Producción
- **Precondiciones:** La Orden de Venta existe en el módulo `Sales` (ERP Core). El `Product` y `Variant` existen en el módulo `Product` (ERP Core).
- **Flujo Normal:**
  1. El Planificador inicia la creación de una ProductionOrder.
  2. El sistema consulta `Sales` (ERP Core) para obtener detalles de la Orden de Venta.
  3. El sistema consulta `Product` (ERP Core) para obtener detalles del Producto/Variante a producir.
  4. El sistema crea una nueva `ProductionOrder` con estado 'Pending'.
- **Postcondiciones:** Una nueva `ProductionOrder` ha sido creada en estado 'Pending'.

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
