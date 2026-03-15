# Manufacturing Execution System (MES) Module Documentation

**Nombre del Módulo:** MES (Manufacturing Execution System)  
**Bounded Context:** Personalización y modificación de prendas textiles  
**Responsabilidad Principal:** Gestionar el ciclo de vida de la producción en taller: desde la definición de tipos de trabajo hasta la ejecución de órdenes reales.  
**Entidades Raíz:** WorkOrder, WorkSetup, WorkType, Task, Position  
**Dependencias:** ERP Core (Sales, Product, Party, IAM)  

---

## 1. Visión General del Módulo

El módulo MES de TramaTex controla la personalización de artículos textiles (marcado por vinilo, serigrafía, bordado, etc.). Su arquitectura separa tres niveles: **Datos Maestros** (Task, Position, WorkType), **Configuración** (WorkSetup) y **Ejecución** (WorkOrder).

### Objetivos Estratégicos
*   **Trazabilidad Total:** Registro de quién, cuándo y dónde se realizó cada tarea.
*   **Terminal de Taller:** Interfaz optimizada para operarios en planta.
*   **Modularidad Extraíble:** Clean Architecture para futura extracción como microservicio (ADR-018 y ADR-022).

---

## 2. Entidades del Dominio

| Entidad | Responsabilidad |
|---|---|
| **Task (Tarea)** | Proceso atómico (ej. Diseñar, Imprimir, Marcar, Plegar, Embolsar). |
| **Position (Posición)** | Zona de la prenda donde se realiza el trabajo (ej. Pecho izdo., Espalda). |
| **WorkType (Tipo de Trabajo)** | Secuencia ordenada de tareas para un tipo de marcado (receta). |
| **WorkSetup (Configuración de Trabajo)** | Plantilla por cliente + tipo de prenda: líneas de WorkType + Position. |
| **WorkOrder (Orden de Trabajo)** | Ejecución real vinculada a pedido, con tiempos y operarios. |

---

## 3. Flujo de Trabajo en el Taller

1.  **Datos Maestros:** El administrador configura las **Tareas**, **Posiciones** y **Tipos de Trabajo**.
2.  **Configuración:** Se crean **WorkSetups** que definen la personalización por cliente/prenda (N líneas de tipo de trabajo + posición).
3.  **Lanzamiento:** Al recibir un pedido, se genera una **Orden de Trabajo** desde un WorkSetup (o manualmente).
4.  **Ejecución:** Los operarios gestionan tareas desde el **Terminal de Taller**.
5.  **Control:** El jefe de taller monitoriza desde el **Dashboard de Producción**.

---

## 4. Documentación Detallada

*   **Especificación del Módulo:** [module-spec.md](./module-spec.md) - Entidades, estados y componentes.
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md) - Entidades, jerarquía y reglas.
*   **Casos de Uso:** [use-cases.md](./use-cases.md) - Flujos funcionales.
*   **Contratos de API:** [api-contracts.md](./api-contracts.md) - Endpoints y contratos.

---

## 5. Decisiones Arquitectónicas (ADRs Relacionados)

*   [ADR-018: MES Module Architecture](../../architecture/adrs/adr-018-mes-module-architecture.md)
*   [ADR-022: MES Microservice Extraction Strategy](../../post-mvp/adr-022-mes-microservice-extraction-strategy.md)

---
**Última Actualización:** 14 de marzo de 2026
