# Modelo de Dominio - Módulo MES

Este documento describe el modelo de dominio del módulo MES (Manufacturing Execution System) de TramaTex, centrado en la personalización y modificación de prendas textiles.

---

## 1. Entidades del Dominio

El módulo MES se estructura en **5 entidades** organizadas en dos niveles: datos maestros (configuración) y ejecución (producción real).

### Datos Maestros

#### 1) Tarea (`Task`)
Proceso atómico e indivisible dentro del flujo de personalización de una prenda. Cada tarea se ejecuta con las mismas herramientas y por una persona.
- **Ejemplos:** Diseñar, Imprimir, Marcar, Plegar, Embolsar.
- **Campos:** `ID`, `Name`, `Description`, `IsActive`.
- **Regla:** Una tarea no define orden ni contexto — es un bloque funcional reutilizable.

#### 2) Posición (`Position`)
Zona de la prenda donde se realiza un trabajo de personalización/modificación.
- **Ejemplos:** Pecho izquierdo, Pecho derecho, Espalda, Bajo pantalón, Manga izquierda.
- **Campos:** `ID`, `Name`, `Code` (único), `Description`, `IsActive`.
- **Regla:** Es un dato maestro transversal, usado tanto en `WorkSetup` como en `WorkOrder`.

#### 3) Tipo de Trabajo (`WorkType`)
Secuencia ordenada de tareas que define un tipo de marcado/personalización concreto. Es la "receta" que indica todas las tareas necesarias para un proceso completo.
- **Ejemplos:** "Marcado por vinilo" → Diseñar → Imprimir → Marcar → Plegar → Embolsar. "Bordado" → Diseñar → Bordar → Plegar → Embolsar.
- **Campos:** `ID`, `Name`, `Description`, `IsActive`, `Tasks[]` (lista ordenada de `WorkTypeTask` con `TaskID` + `Sequence`).
- **Regla:** Cada `WorkType` define exactamente qué tareas se ejecutan y en qué orden. No tiene posición — la posición se asigna al usar el tipo de trabajo en un `WorkSetup` o `WorkOrder`.

### Configuración por Cliente

#### 4) Configuración de Trabajo (`WorkSetup`)
Define la personalización completa de un tipo de prenda para un cliente concreto. Combina un cliente (`Party`) + grupo de producto tangible (`TangibleGroup`) con una o varias líneas que indican qué trabajo realizar y en qué posición.
- **Ejemplo:** Para Confecciones López, grupo "Camisetas": línea 1 → Serigrafía en Pecho izquierdo, línea 2 → Vinilo en Espalda.
- **Campos:** `ID`, `Name`, `PartyID`, `TangibleGroupID`, `Description`, `IsActive`, `Lines[]`.
- **WorkSetupLine:** `ID`, `WorkTypeID`, `PositionID`, `DesignFilePath`, `Notes`, `Sequence`.
- **Regla:** Un `WorkSetup` es una plantilla reutilizable. Cuando se recibe un pedido del cliente para ese tipo de prenda, se instancia como `WorkOrder`.

### Ejecución

#### 5) Orden de Trabajo (`WorkOrder`)
Trabajo real vinculado a un pedido de venta, con prendas físicas, tiempos de ejecución y operarios asignados. Es la instancia ejecutable de un `WorkSetup` (o creada manualmente).
- **Campos:** `ID`, `OrderNumber` (único, auto-generado), `OrderName`, `PartyID`, `TangibleGroupID`, `WorkSetupID` (opcional — referencia a la plantilla origen), `GarmentNotes`, `Status`, `Priority`, `StartDate`, `DueDate`, `CompletedDate`, `Lines[]`.
- **WorkOrderLine:** `ID`, `WorkTypeID`, `PositionID`, `DesignFilePath`, `Notes`, `Sequence`, `Tasks[]`.
- **WorkOrderTask:** `ID`, `TaskID`, `Sequence`, `Status`, `AssignedTo`, `StartedAt`, `CompletedAt`, `Notes`.
- **Regla:** Al crear un `WorkOrder` desde un `WorkSetup`, las líneas y tareas se copian (no se referencian). El `WorkOrder` vive de forma independiente para permitir modificaciones en producción sin afectar la plantilla.

---

## 2. Jerarquía y Relaciones

```
Task (Tarea)                              ← proceso atómico
  ↑ se agrupa en
WorkType (Tipo de Trabajo)                ← receta: secuencia de tareas
  ↑ se usa en líneas de
WorkSetup (Configuración de Trabajo)      ← plantilla: cliente + tipo prenda → N líneas (WorkType + Position)
  ↑ se instancia como
WorkOrder (Orden de Trabajo)              ← ejecución real: pedido + prendas + tiempos
  └── WorkOrderLine                       ← línea ejecutable (WorkType + Position)
       └── WorkOrderTask                  ← tarea con estado/operario/tiempos

Position (Posición)                       ← zona de la prenda (transversal)
```

---

## 3. Estados de Producción

### Estados de la Orden de Trabajo (`WorkOrder.Status`)
| Estado | Descripción |
|---|---|
| `DRAFT` | En preparación, editable libremente. |
| `PENDING` | Lista para iniciar producción. |
| `IN_PROGRESS` | Al menos una tarea ha comenzado. |
| `COMPLETED` | Todas las tareas de todas las líneas finalizadas. |
| `ON_HOLD` | Pausada por incidencia o decisión. |
| `CANCELLED` | Descartada. |

### Estados de Tarea (`WorkOrderTask.Status`)
| Estado | Descripción |
|---|---|
| `PENDING` | No iniciada. |
| `IN_PROGRESS` | En curso (operario trabajando). |
| `COMPLETED` | Finalizada. |
| `BLOCKED` | Incidencia o problema. |
| `SKIPPED` | Omitida (decisión del jefe de taller). |

### Prioridades (`WorkOrder.Priority`)
`LOW` · `NORMAL` · `HIGH` · `URGENT`

---

## 4. Reglas de Comportamiento

### Trazabilidad "Just-in-Time"
El sistema registra el tiempo real consumido por tarea. Un operario no puede completar una tarea sin haberla iniciado previamente. El terminal registra el usuario del módulo **IAM** que realiza la acción.

### Independencia de la Orden
Al crear un `WorkOrder` desde un `WorkSetup`, las líneas y tareas se **copian** como snapshot. Las modificaciones posteriores al `WorkSetup` no afectan a órdenes ya creadas, y el jefe de taller puede modificar la orden en producción sin alterar la plantilla.

### Soberanía del Taller
Una vez que el `WorkOrder` entra en producción, el jefe de taller tiene soberanía sobre las notas técnicas, los archivos de diseño y la prioridad.

---

## 5. Integración con Sales

El módulo MES recibe solicitudes de trabajo desde el módulo **Sales** a través de la entidad `SalesWorkSetup` (definida en el dominio de Sales). Esta integración funciona así:

### Flujo de Entrada
1. El comercial asocia trabajos a presupuestos/pedidos en Sales (con nombre, observaciones y opcionalmente un `WorkSetupID`).
2. Al confirmar un pedido, los trabajos quedan en estado `PENDIENTE` en Sales.
3. El jefe de taller consulta los trabajos pendientes de Sales (CU-M-010).
4. Si el `SalesWorkSetup` tiene un `WorkSetupID`, el taller crea la `WorkOrder` directamente desde la plantilla asociada.
5. Si solo tiene nombre y observaciones, el jefe de taller define manualmente el `WorkSetup` y/o `WorkOrder`.

### Flujo de Retorno (Eventos de Dominio)
- **`WorkOrderStarted`:** Al iniciar la primera tarea de un `WorkOrder` vinculado a Sales, se emite un evento que actualiza el `SalesWorkSetup` a `EN_PROCESO`.
- **`WorkOrderCompleted`:** Al completar todas las tareas de un `WorkOrder`, se emite un evento que actualiza el `SalesWorkSetup` a `COMPLETADO`.

### Principio de Desacoplamiento
MES no depende del dominio de Sales. La comunicación es mediante:
- **Consulta:** MES lee los trabajos pendientes a través de una interfaz de consulta (puerto en el dominio MES).
- **Notificación:** MES emite eventos de dominio que Sales consume para actualizar sus estados.

### Servicio de Consulta de Progreso (`WorkOrderQueryService`)
MES expone un servicio de consulta que permite a otros módulos conocer el estado de ejecución de las órdenes de trabajo **sin acceder a los detalles internos del dominio MES**. El servicio:
- Recibe un `WorkOrderID` (o varios) y devuelve un DTO pre-computado (`WorkOrderProgressDTO`) con: estado global, número total de tareas, tareas completadas y desglose por línea.
- Sales consume este servicio a través de su interfaz `MESWorkLookup` (capa de aplicación), implementada por un adaptador en infraestructura que delega al `WorkOrderQueryService` de MES.
- **Toda la lógica de cálculo de progreso** (conteo de tareas, estados por línea) reside en MES. Sales solo recibe el resultado.

---

### El Terminal de Taller (Operativa Táctil)
1. El operario accede al terminal y ve las tareas pendientes.
2. Selecciona una tarea → **START**.
3. Al terminar → **COMPLETE**. Si hay incidencia → **BLOCK**.

---
**Última Actualización:** 14 de marzo de 2026
