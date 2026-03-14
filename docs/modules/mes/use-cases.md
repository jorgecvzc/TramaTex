# Casos de Uso - Módulo MES

Este documento describe los casos de uso que orquestan el proceso de fabricación y gestión de servicios técnicos en el taller de TramaTex.

---

## 1. Gestión de Maestros y Recetas

### **CU-M-001: Configurar Grupo de Servicio**
*   **Propósito:** Definir una "receta" reusable de tareas para un tipo de servicio o producto.
*   **Actores:** Jefe de Taller / Administrador.
*   **Flujo:**
    1.  Crear un `ServiceGroup`.
    2.  Asociar una secuencia de `Task` existentes (ej. "Preparación", "Bordado", "Limpieza").
    3.  (Opcional) Vincular el grupo a un `ProductGroup` tangible del catálogo.
*   **Resultado:** Receta disponible para ser instanciada en trabajos reales.

---

## 2. Planificación y Ejecución

### **CU-M-002: Lanzar Trabajo de Producción (MESWork)**
*   **Propósito:** Iniciar una ejecución de trabajo real en el taller.
*   **Actores:** Jefe de Taller, Sistema (automatizado desde Sales).
*   **Flujo:**
    1.  Recibir una orden de producción (o crearla manualmente).
    2.  Asignar uno o varios `ServiceGroups` necesarios para el trabajo.
    3.  Definir la `Position` (puesto) para cada grupo.
    4.  Asociar notas técnicas (`GarmentNotes`) y archivos de diseño.
    5.  Establecer la prioridad y fecha de entrega.
*   **Resultado:** El trabajo aparece en el **Dashboard de Producción** en estado `PENDING`.

### **CU-M-003: Iniciar Tarea en Terminal**
*   **Propósito:** Registrar el comienzo físico de una actividad por parte de un operario.
*   **Actores:** Operario de Taller.
*   **Flujo:**
    1.  El operario accede al terminal de su puesto (`Position`).
    2.  Selecciona una tarea en estado `PENDING`.
    3.  Pulsa "Comenzar".
*   **Resultado:** El estado de la tarea cambia a `IN_PROGRESS`. El sistema registra la hora de inicio y el operario.

### **CU-M-004: Finalizar Tarea en Terminal**
*   **Propósito:** Registrar la terminación de una actividad.
*   **Actores:** Operario de Taller.
*   **Flujo:**
    1.  El operario selecciona su tarea en curso.
    2.  Añade notas si es necesario.
    3.  Pulsa "Finalizar".
*   **Resultado:** El estado de la tarea cambia a `COMPLETED`. Se calcula el tiempo real consumido. Si es la última tarea del trabajo, el `MESWork` pasa a `COMPLETED`.

---

## 3. Seguimiento y Control

### **CU-M-005: Monitorizar Dashboard de Producción**
*   **Propósito:** Tener una visión global del taller y el estado de los pedidos.
*   **Actores:** Jefe de Taller, Comercial (Consulta).
*   **Resultado:** Visualización en tiempo real de los trabajos por estado (`IN_PROGRESS`, `ON_HOLD`, etc.) y alertas por retrasos o bloqueos.

---
**Última Actualización:** 9 de marzo de 2026
