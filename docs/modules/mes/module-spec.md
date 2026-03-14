# Módulo de MES (Manufacturing Execution System)

**Estado:** ✅ **COMPLETO (100%)**  
**Última actualización:** 9 de marzo de 2026

## 1. Propósito

*   **Visión del Módulo:** Controlar y monitorizar el proceso de fabricación y prestación de servicios técnicos en el taller de TramaTex, asegurando la trazabilidad de cada fase de producción.
*   **Objetivos Clave:**
    *   Gestionar **Grupos de Servicio** (plantillas de procesos) vinculados a tipos de productos.
    *   Orquestar la **Ejecución de Trabajos (MESWork)** con seguimiento de tiempos y operarios.
    *   Optimizar la interacción mediante un **Terminal de Taller** táctil para operarios.

---

## 2. Definiciones Clave (Nomenclatura Actualizada)

### 1) Grupo de Servicio (`ServiceGroup`)
- **Propósito:** Actúa como la "Definición" o "Receta" del proceso. Agrupa una secuencia de tareas predefinidas.
- **Alcance:** Puede estar vinculado a un `ProductGroup` tangible para automatizar la asignación de procesos.

### 2) Puesto de Trabajo (`Position`)
- **Propósito:** Representa una ubicación física o lógica en el taller (ej. "Mesa de Corte", "Bordadora 1"). Cada grupo de servicio en una ejecución se asigna a una posición.

### 3) Ejecución de Trabajo (`MESWork`)
- **Propósito:** La instancia real de producción en el taller.
- **Estructura:** Un `MESWork` contiene múltiples `ServiceGroups`, y cada uno contiene su propia lista de `Tasks` (Tareas) ejecutables.

### 4) Tarea de Trabajo (`MESWorkTask`)
- **Propósito:** La unidad mínima de ejecución que el operario inicia y finaliza en el terminal.

---

## 3. Estados de Producción

### Estados del Trabajo (`MESWork`)
- `DRAFT`: En preparación.
- `PENDING`: Listo para iniciar.
- `IN_PROGRESS`: Al menos una tarea ha comenzado.
- `COMPLETED`: Todas las tareas de todos los grupos han finalizado.
- `ON_HOLD`: Pausado por incidencia.
- `CANCELLED`: Trabajo descartado.

### Estados de la Tarea (`Task`)
- `PENDING`, `IN_PROGRESS`, `COMPLETED`, `BLOCKED`, `SKIPPED`.

---

## 4. Componentes Implementados ✅

### Backend
- **Maestros de Taller:** CRUD de Tareas (`Task`) y Puestos (`Position`).
- **Gestor de Recetas:** CRUD de `ServiceGroups` con definición de secuencias.
- **Motor de Producción:** Gestión integral de `MESWork` y transiciones de estado de tareas.
- **Persistencia de Diseño:** Soporte para rutas de archivos de diseño (`DesignFilePath`) y notas de prenda (`GarmentNotes`).

### Frontend
- **Dashboard de Producción:** Vista de seguimiento del estado de los trabajos activos.
- **Terminal de Taller:** Interfaz táctil para que los operarios gestionen el ciclo de vida de sus tareas (Start/Complete).
- **Admin MES:** Configuración de grupos de servicio y puestos de trabajo.
