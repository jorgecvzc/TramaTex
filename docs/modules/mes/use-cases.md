# Casos de Uso - Módulo MES

Este documento describe los casos de uso del módulo MES para la personalización y modificación de prendas textiles en TramaTex.

---

## 1. Gestión de Datos Maestros

### **CU-M-001: Gestionar Tareas**
*   **Propósito:** Mantener el catálogo de tareas atómicas (Diseñar, Imprimir, Marcar, Plegar, Embolsar, etc.).
*   **Actores:** Jefe de Taller / Administrador.
*   **Operaciones:** CRUD con activación/desactivación.

### **CU-M-002: Gestionar Posiciones**
*   **Propósito:** Mantener el catálogo de zonas de prenda (Pecho izquierdo, Espalda, Bajo pantalón, etc.).
*   **Actores:** Jefe de Taller / Administrador.
*   **Operaciones:** CRUD con código único y activación/desactivación.

### **CU-M-003: Configurar Tipo de Trabajo**
*   **Propósito:** Definir una "receta" de tareas ordenadas para un tipo de marcado/personalización.
*   **Actores:** Jefe de Taller / Administrador.
*   **Flujo:**
    1.  Crear un `WorkType` con nombre descriptivo (ej. "Marcado por vinilo").
    2.  Asociar una secuencia de `Task` existentes (ej. Diseñar → Imprimir → Marcar → Plegar → Embolsar).
*   **Resultado:** Tipo de trabajo disponible para usarse en configuraciones y órdenes.

---

## 2. Configuración por Cliente

### **CU-M-004: Definir Configuración de Trabajo**
*   **Propósito:** Crear una plantilla que define la personalización de un tipo de prenda para un cliente concreto.
*   **Actores:** Jefe de Taller / Comercial.
*   **Flujo:**
    1.  Seleccionar el cliente (`Party`) y el grupo de producto tangible.
    2.  Añadir líneas: cada línea indica un `WorkType` (tipo de marcado) + `Position` (zona de la prenda).
    3.  Opcionalmente, adjuntar archivos de diseño y notas por línea.
*   **Ejemplo:** Confecciones López / Camisetas → Serigrafía en Pecho izdo. + Vinilo en Espalda.
*   **Resultado:** Configuración guardada y reutilizable para futuros pedidos del mismo cliente/prenda.

---

## 3. Planificación y Ejecución

### **CU-M-005: Crear Orden de Trabajo**
*   **Propósito:** Iniciar una ejecución real de producción en el taller.
*   **Actores:** Jefe de Taller, Sistema (automatizado desde Sales).
*   **Flujo:**
    1.  Seleccionar un `WorkSetup` existente (o crear la orden manualmente).
    2.  El sistema copia las líneas y tareas del setup como snapshot independiente.
    3.  Ajustar notas técnicas, archivos de diseño y prioridad si es necesario.
    4.  Establecer fecha de entrega.
*   **Resultado:** Orden de trabajo en estado `DRAFT` o `PENDING`, visible en el Dashboard.

### **CU-M-006: Iniciar Tarea en Terminal**
*   **Propósito:** Registrar el comienzo físico de una tarea por un operario.
*   **Actores:** Operario de Taller.
*   **Flujo:**
    1.  El operario accede al terminal.
    2.  Selecciona una tarea en estado `PENDING`.
    3.  Pulsa "Comenzar".
*   **Resultado:** Tarea cambia a `IN_PROGRESS`. Se registra hora de inicio y operario (IAM).

### **CU-M-007: Finalizar Tarea en Terminal**
*   **Propósito:** Registrar la terminación de una tarea.
*   **Actores:** Operario de Taller.
*   **Flujo:**
    1.  Selecciona su tarea en curso.
    2.  Añade notas si es necesario.
    3.  Pulsa "Finalizar".
*   **Resultado:** Tarea cambia a `COMPLETED`. Se calcula tiempo real consumido. Si es la última tarea de la orden, el `WorkOrder` pasa a `COMPLETED`.

### **CU-M-008: Reportar Incidencia en Tarea**
*   **Propósito:** Bloquear una tarea por problema o incidencia.
*   **Actores:** Operario de Taller.
*   **Resultado:** Tarea cambia a `BLOCKED`. Se notifica al jefe de taller.

---

## 4. Seguimiento y Control

### **CU-M-009: Monitorizar Dashboard de Producción**
*   **Propósito:** Visión global del taller y estado de las órdenes.
*   **Actores:** Jefe de Taller, Comercial (consulta).
*   **Resultado:** Visualización en tiempo real de órdenes por estado, alertas por retrasos o bloqueos.

---

## 5. Integración con Sales

### **CU-M-010: Consultar Trabajos Pendientes de Sales**
*   **Propósito:** Permitir al jefe de taller ver los trabajos de personalización solicitados desde el módulo de Ventas que están pendientes de ser recogidos.
*   **Actores:** Jefe de Taller.
*   **Flujo:**
    1.  El jefe de taller accede a la vista de trabajos pendientes.
    2.  El sistema consulta los `SalesWorkSetup` en estado `PENDIENTE` (a través del puerto `SalesPendingWorkPort`).
    3.  Muestra una lista con: nombre del trabajo, observaciones del comercial, cliente, referencia del pedido, y si tiene un `WorkSetupID` vinculado.
*   **Resultado:** Lista de trabajos pendientes de Sales para planificación del taller.

### **CU-M-011: Crear Orden de Trabajo desde Solicitud de Sales**
*   **Propósito:** Generar una `WorkOrder` a partir de un trabajo pendiente solicitado desde Sales.
*   **Actores:** Jefe de Taller.
*   **Flujo:**
    1.  Seleccionar un trabajo pendiente de la lista (CU-M-010).
    2.  **Si tiene `WorkSetupID`:** crear la `WorkOrder` desde la plantilla (igual que CU-M-005), copiando líneas y tareas como snapshot. El frontend incluye el `order_work_setup_id` en el payload de creación.
    3.  **Si no tiene `WorkSetupID`:** el jefe de taller define manualmente las líneas y tareas guiándose por las observaciones del comercial. Opcionalmente puede crear primero un `WorkSetup` como plantilla.
    4.  Al crear la `WorkOrder`, el backend invoca `SalesOrderLinker.LinkWorkOrder(ctx, orderWorkSetupID, workOrderID)` a través del adaptador `SalesOrderLinkerAdapter` (infraestructura de Sales), que actualiza el campo `work_order_id` en `order_work_setups`.
*   **Resultado:** `WorkOrder` creada y vinculada al trabajo de Sales. La solicitud desaparece del panel de pendientes del Dashboard MES.

### **CU-M-012: Consultar Progreso de Órdenes de Trabajo**
*   **Propósito:** Proporcionar a módulos externos (Sales) el estado de ejecución de una o varias `WorkOrder`s, incluyendo desglose por línea y tarea.
*   **Actores:** Sistema (Sales u otro módulo consumidor).
*   **Entradas:** Uno o varios `WorkOrderID`.
*   **Flujo:**
    1.  Recibir la petición de consulta.
    2.  Cargar las `WorkOrder`(s) con todas sus líneas y tareas.
    3.  Computar por cada orden: estado global, total de tareas, tareas completadas; por cada línea: total de tareas, tareas completadas.
    4.  Devolver un DTO de progreso pre-computado (`WorkOrderProgressDTO`).
*   **Resultado:** Información de progreso lista para consumo. Sales lo recibe a través de la interfaz `MESWorkLookup`. Toda la lógica de cálculo permanece en MES.

### **CU-M-013: Suspender / Reactivar Órdenes de Trabajo (desde Sales)**
*   **Propósito:** Permitir al módulo Sales suspender temporalmente o reactivar las `WorkOrder`s asociadas a un pedido cuando su estado cambia (ej. cancelación o reactivación del pedido).
*   **Actores:** Sistema (módulo Sales, vía `WorkOrderSuspenderAdapter`).
*   **Entradas:** Lista de `WorkOrderID`s.
*   **Flujo de Suspensión:**
    1.  El servicio Sales llama a `WorkOrderSuspender.SuspendWorkOrders(ctx, []WorkOrderID)`.
    2.  MES carga cada orden y transiciona: PENDING → SUSPENDED, IN_PROGRESS → SUSPENDED, ON_HOLD → SUSPENDED.
    3.  Se ignoran las órdenes en estado COMPLETED o CANCELLED.
*   **Flujo de Reactivación:**
    1.  El servicio Sales llama a `WorkOrderSuspender.ReactivateWorkOrders(ctx, []WorkOrderID)`.
    2.  MES carga cada orden y transiciona: SUSPENDED → PENDING, ON_HOLD → PENDING, CANCELLED → PENDING.
    3.  Se ignoran las órdenes en estado COMPLETED o IN_PROGRESS.
*   **Patrón arquitectónico:** Sales define la interfaz `WorkOrderSuspender`. El adaptador `WorkOrderSuspenderAdapter` (infraestructura de Sales) recibe `MESService` por setter injection y delega la llamada. MES no conoce el concepto de pedido de venta.
*   **Resultado:** Las `WorkOrder`s afectadas quedan en estado SUSPENDED (o PENDING al reactivar). El Dashboard MES las muestra en la sección correspondiente.

### **CU-M-014: Lanzar Orden de Trabajo a Taller**
*   **Propósito:** Activar una orden de trabajo ya configurada para que sea visible en el Terminal del Taller.
*   **Actores:** Jefe de Taller.
*   **Precondiciones:**
    - La orden de trabajo está en estado `PENDING`.
    - La orden tiene asignado un `WorkSetupID` (configuración técnica).
*   **Flujo:**
    1. El jefe de taller visualiza la orden en la sección "Pendientes de Inicio" del Dashboard.
    2. Identifica que la orden ya tiene configuración y pulsa el botón "Lanzar".
    3. El sistema actualiza el estado de la orden a `IN_PROGRESS`.
*   **Resultado:** La orden desaparece de la sección de pendientes y sus tareas se vuelven visibles en el Terminal de Taller (Tablet).

---
**Última Actualización:** 20 de marzo de 2026
