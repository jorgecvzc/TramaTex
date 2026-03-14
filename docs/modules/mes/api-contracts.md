# Contratos de API - Módulo MES

Este documento detalla la interfaz de integración para la gestión de la producción y el control del taller en TramaTex.

---

## 1. Maestros y Recetas (`/api/mes/service-groups`)

Gestiona el catálogo de grupos de servicios técnicos y procesos de fabricación.
- **Grupos de Servicio (`ServiceGroup`):** Puntos de entrada para configurar secuencias de tareas predefinidas.
- **Tareas y Puestos (`/api/mes/tasks`, `/api/mes/positions`):** CRUD de los elementos base que componen las recetas.

## 2. Ejecución de Producción (`/api/mes/works`)

Es el motor operativo para el seguimiento del trabajo real en taller.
- **Lanzamiento de Trabajos (`POST /api/mes/works`):** Permite disparar una ejecución (`MESWork`) a partir de un pedido de venta o de una planificación directa.
- **Dashboard de Seguimiento (`GET /api/mes/works`):** Proporciona visibilidad en tiempo real sobre el avance global de los trabajos activos.

## 3. Terminal de Taller (`/api/mes/works/{workId}/tasks/{taskId}`)

Puntos de entrada optimizados para el uso por operarios en planta a través de tablets.
- **Ciclo de Vida de Tarea:** Endpoints para marcar el inicio (`POST /start`) y la finalización (`POST /complete`) de tareas individuales de un trabajo.
- **Incidencias y Notas:** Permite registrar bloqueos y observaciones técnicas en tiempo real.

---

## Notificaciones de Eventos

El módulo de MES emite eventos de dominio cuando un `MESWork` cambia de estado a `COMPLETED`. Estas notificaciones son consumidas por el módulo de **Sales** para automatizar el flujo comercial (ej. marcar pedido como listo para envío).

---
**Última Actualización:** 9 de marzo de 2026
