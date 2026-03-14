# Manufacturing Execution System (MES) Module Documentation

**Nombre del Módulo:** MES (Manufacturing Execution System)  
**Bounded Context:** Producción personalizada, Planificación y Ejecución de Manufactura  
**Responsabilidad Principal:** Gestionar el ciclo de vida de la producción en taller mediante grupos de servicio y tareas.  
**Entidades Raíz:** MESWork, ServiceGroup, Task, Position  
**Dependencias:** ERP Core (Sales, Product, Party, IAM)  

---

## 1. Visión General del Módulo

El módulo MES de TramaTex está diseñado para controlar la fabricación de artículos personalizados (bordados, estampados, confección) y la prestación de servicios técnicos. Su arquitectura se basa en la separación de la **Receta de Proceso** (ServiceGroup) y la **Ejecución Real** (MESWork).

### Objetivos Estratégicos
*   **Trazabilidad Total:** Registro exacto de quién, cuándo y dónde se realizó cada tarea.
*   **Terminal de Taller:** Interfaz optimizada para operarios que trabajan en entornos de producción.
*   **Modularidad Extraíble:** Diseñado siguiendo Clean Architecture para permitir su futura extracción como microservicio independiente (ADR-018 y ADR-022).

---

## 2. Definiciones del Dominio

| Término | Responsabilidad |
|---|---|
| **Task (Tarea)** | Definición base de una actividad técnica (ej. "Corte", "Bordado"). |
| **Position (Puesto)** | Ubicación física o máquina donde se ejecutan las tareas. |
| **ServiceGroup (Grupo de Servicio)** | Plantilla reusable que agrupa una secuencia de tareas para un tipo de producto. |
| **MESWork (Trabajo)** | Entidad raíz que orquesta la ejecución real de una orden de producción. |
| **MESWorkTask** | Instancia ejecutable de una tarea con seguimiento de tiempos y operario. |

---

## 3. Flujo de Trabajo en el Taller

1.  **Definición:** El administrador configura las **Tareas** y los **Puestos de Trabajo**.
2.  **Receta:** Se crean **Grupos de Servicio** que encadenan tareas en una secuencia lógica.
3.  **Lanzamiento:** Al confirmar un pedido en `Sales`, se genera un **MESWork** (o se crea manualmente desde el Dashboard).
4.  **Ejecución:** Los operarios ven las tareas pendientes en el **Terminal de Taller** filtradas por su puesto de trabajo.
5.  **Control:** El jefe de taller monitoriza el avance desde el **Dashboard de Producción** y gestiona prioridades o incidencias.

---

## 4. Documentación Detallada

Para una descripción profunda de los componentes del módulo MES, consulte los siguientes documentos:

*   **Especificación del Módulo:** [module-spec.md](./module-spec.md) - Requisitos y lógica de estados.
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md) - Entidades y reglas de negocio.
*   **Contratos de API:** [api-contracts.md](./api-contracts.md) - Puntos de integración técnica.

---

## 5. Decisiones Arquitectónicas (ADRs Relacionados)

*   [ADR-018: MES Module Architecture](../../architecture/adrs/adr-018-mes-module-architecture.md)
*   [ADR-022: MES Microservice Extraction Strategy](../../post-mvp/adr-022-mes-microservice-extraction-strategy.md)

---
**Última Actualización:** 9 de marzo de 2026
