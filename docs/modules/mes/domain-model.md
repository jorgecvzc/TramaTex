# Modelo de Dominio - Módulo MES

Este documento describe la lógica de control de fabricación en TramaTex, centrada en la trazabilidad de los procesos desde las recetas de servicio hasta la ejecución física.

---

## 1. El Ciclo de Fabricación (Estructura Jerárquica)

El módulo MES organiza el trabajo en el taller mediante una jerarquía de tres niveles:

### 1) El Trabajo de Producción (`MESWork`)
Es la entidad raíz que agrupa todo el proceso para un cliente o pedido.
- **Identidad:** Tiene un `WorkNumber` único.
- **Vínculos:** Referencia a una `Party` (cliente) y a un `ProductGroup` tangible.
- **Información Adicional:** Incluye notas específicas de prenda (`GarmentNotes`) y plazos (`StartDate`, `DueDate`).

### 2) Los Grupos de Servicio (`MESWorkServiceGroup`)
Un trabajo de producción se divide en uno o varios grupos de servicios técnicos.
- **Contexto:** Cada grupo representa una especialidad (ej. "Bordado", "Estampado").
- **Asignación:** Se vincula a un puesto de trabajo físico (`Position`) y puede incluir una ruta a archivos de diseño (`DesignFilePath`).

### 3) Las Tareas de Trabajo (`MESWorkTask`)
Son las unidades ejecutables finales dentro de cada grupo de servicio.
- **Trazabilidad:** Registra el operario asignado (`AssignedTo`) y las marcas de tiempo (`StartedAt`, `CompletedAt`).
- **Estado Atómico:** Es lo que el operario gestiona en el terminal.

---

## 2. Maestros de Proceso

### Tarea Maestro (`Task`)
Definición base de una actividad (ej. "Preparación de Máquina", "Ejecución de Diseño").
- **Tipos:** Soporta tareas únicas (`ONE_TIME`) o recurrentes (`RECURRENT`).

### Puesto de Trabajo (`Position`)
Representa la ubicación técnica donde se realiza el trabajo (ej. "Máquina Bordadora X1", "Sección de Embolsado").

### Grupo de Servicio Maestro (`ServiceGroup`)
La "Receta" reusable. Define una secuencia de tareas maestro que deben ejecutarse habitualmente juntas. Puede estar vinculado a un grupo de productos para disparar procesos automáticamente.

---

## 3. Reglas de Comportamiento Críticas

### Trazabilidad "Just-in-Time"
El sistema registra el tiempo real consumido. Un operario no puede completar una tarea sin haberla iniciado previamente. El terminal registra el usuario del módulo **IAM** que realiza la acción para asegurar la responsabilidad individual.

### Soberanía del Taller
Aunque el pedido nazca en `Sales`, una vez que el `MESWork` entra en producción, el jefe de taller tiene soberanía sobre las notas técnicas, los archivos de diseño y la prioridad del trabajo (`LOW` a `URGENT`).

### El Terminal de Taller (Operativa Táctil)
Diseñado para la máxima simplicidad:
1. El operario escanea el código del trabajo o lo selecciona del listado.
2. El sistema muestra las tareas de su posición.
3. El operario interactúa con estados claros: **START**, **COMPLETE**, **BLOCK** (incidencia).

---
**Última Actualización:** 9 de marzo de 2026
