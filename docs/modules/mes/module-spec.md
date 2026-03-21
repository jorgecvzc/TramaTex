# Módulo de MES (Manufacturing Execution System)

**Estado:** ✅ **ACTIVO (MVP)**  
**Última Actualización:** 20 de marzo de 2026

> **Convención de UI (i18n):** Todos los estados, prioridades y etiquetas de acción se muestran en **castellano** en la interfaz de usuario, independientemente de que se almacenen en inglés en la API.

## 1. Propósito

*   **Visión del Módulo:** Controlar y monitorizar el proceso de personalización y modificación de prendas textiles en el taller de TramaTex, asegurando la trazabilidad de cada fase.
*   **Objetivos Clave:**
    *   Gestionar **Tipos de Trabajo** (recetas de tareas) para cada tipo de marcado/personalización.
    *   Definir **Configuraciones de Trabajo** por cliente y tipo de prenda, reutilizables en pedidos.
    *   Orquestar **Órdenes de Trabajo** reales con seguimiento de tiempos y operarios.
    *   Optimizar la interacción mediante un **Terminal de Taller** táctil para operarios.

---

## 2. Entidades del Dominio

### 1) Tarea (`Task`)
- **Propósito:** Proceso atómico e indivisible (ej. Diseñar, Imprimir, Marcar, Plegar, Embolsar).
- **Campos clave:** `Name`, `Reference` (código corto), `Description`, `IsActive`.
- **UI:** "Tareas" en sección Datos Maestros.

### 2) Posición (`Position`)
- **Propósito:** Zona de la prenda donde se realiza el trabajo (ej. Pecho izquierdo, Espalda, Bajo pantalón).
- **UI:** "Posiciones" en sección Datos Maestros.

### 3) Tipo de Trabajo (`WorkType`)
- **Propósito:** Secuencia ordenada de tareas que define un tipo de marcado/personalización. Es la "receta" del proceso.
- **Ejemplo:** "Marcado por vinilo" = Diseñar → Imprimir → Marcar → Plegar → Embolsar.
- **Campos clave:** `Name`, `Reference` (código corto para identificación rápida), `Description`, `IsActive`, `Tasks[]`.
- **Backend:** Tabla `work_types` / Entidad `WorkType`.
- **UI:** "Tipos de Trabajo" en sección Datos Maestros.

### 4) Configuración de Trabajo (`WorkSetup`)
- **Propósito:** Plantilla reutilizable que define la personalización de un tipo de prenda para un cliente. Combina múltiples líneas de WorkType + Position.
- **Ejemplo:** Confecciones López / Camisetas → Serigrafía en Pecho + Vinilo en Espalda.
- **UI:** "Configuraciones" en sección Configuración.

### 5) Orden de Trabajo (`WorkOrder`)
- **Propósito:** Instancia real de producción vinculada a un pedido, con prendas físicas, tiempos y operarios.
- **Origen:** Se crea desde un `WorkSetup` (copiando líneas y tareas) o manualmente.
- **UI:** "Órdenes de Trabajo" en sección Producción.

---

## 3. Estados

### Estados de la Orden de Trabajo (`WorkOrder`)
| Código | Castellano (UI) | Descripción |
|---|---|---|
| `PENDING` | Pendiente | Lista para iniciar producción. |
| `IN_PROGRESS` | En progreso | Al menos una tarea ha comenzado. |
| `COMPLETED` | Completado | Todas las tareas finalizadas. |
| `ON_HOLD` | En espera | Pausada por incidencia. |
| `SUSPENDED` | Suspendida | Suspendida temporalmente desde Sales (p.ej. al cancelar un pedido). No puede suspenderse si ya está COMPLETED o CANCELLED. |
| `CANCELLED` | Cancelado | Descartada definitivamente. |

### Estados de Tarea (`WorkOrderTask`)
| Código | Castellano (UI) | Descripción |
|---|---|---|
| `PENDING` | Pendiente | No iniciada. |
| `IN_PROGRESS` | En progreso | Operario trabajando. |
| `COMPLETED` | Completada | Finalizada. |
| `BLOCKED` | Bloqueada | Incidencia o problema. |
| `SKIPPED` | Omitida | Decisión del jefe de taller. |

### Prioridades (`WorkOrder`)
| Código | Castellano (UI) |
|---|---|
| `LOW` | Baja |
| `NORMAL` | Normal |
| `HIGH` | Alta |
| `URGENT` | Urgente |

> **Convención de UI:** Los estados y prioridades se almacenan en inglés en la API y base de datos, pero **siempre se muestran en castellano** en la interfaz de usuario. El frontend mantiene un mapa de traducción centralizado para este propósito.

---

## 4. Componentes

### Backend
- **Maestros de Taller:** CRUD de Tareas (`Task`) y Posiciones (`Position`).
- **Tipos de Trabajo:** CRUD de `WorkType` con definición de secuencias de tareas.
- **Configuraciones:** CRUD de `WorkSetup` con líneas (WorkType + Position).
- **Órdenes de Trabajo:** Gestión de `WorkOrder`, instanciación desde `WorkSetup`, transiciones de estado.
- **Persistencia de Diseño:** Soporte para archivos de diseño y notas de prenda.

### Frontend
- **Dashboard de Producción:** Vista de seguimiento de órdenes activas y **panel de solicitudes pendientes de Sales** (configuraciones de pedidos confirmados sin órdenes de trabajo aún).
- **Terminal de Taller:** Interfaz táctil para operarios (Comenzar/Finalizar/Bloquear).
- **Datos Maestros:** Gestión de Tareas, Posiciones y Tipos de Trabajo.
- **Configuraciones:** Definición de personalizaciones por cliente/prenda.

---

## 5. Integración con Sales — Flujo Completo

El Dashboard de MES actúa como puente entre el departamento comercial y el taller:

### Flujo de trabajo
1. **Sales → Configuraciones:** Desde presupuestos y/o pedidos, el comercial asocia trabajos con nombre + observaciones. Opcionalmente selecciona un `WorkSetup` existente; si no lo hace, el sistema **auto-crea** uno en MES (inactivo, sin líneas) para que el taller lo configure después.
2. **Conversión Presupuesto → Pedido:** Al convertir, los `SalesWorkSetup` se copian y `ensureWorkSetups` garantiza que todos tengan un `WorkSetupID` válido.
3. **Panel MES → Solicitudes:** El Dashboard muestra las Configuraciones asociadas a **pedidos en preparación** (`EN_PREPARACION`) que aún no tienen `WorkOrderID`. Las configuraciones auto-creadas (sin líneas definidas) se señalizan visualmente para que el taller las configure primero.
4. **Creación de Órdenes:** El jefe de taller crea Órdenes de Trabajo desde las Configuraciones. Todas tienen `WorkSetupID` (gracias a la auto-creación); se instancia directamente o se personaliza.
5. **Visibilidad en el Pedido:** Las Órdenes creadas aparecen en el pedido con su estado general y un botón para navegar al detalle de tareas.
6. **Ejecución en Taller:** Los operarios marcan tareas como completadas desde el Terminal de Taller (tablet).
7. **Consulta de estado:** La forma de saber si los trabajos MES de un pedido están terminados es acceder al pedido o al Panel de Control de MES. (Post-MVP: notificaciones, eventos de dominio y dashboard mejorado).

---
**Última Actualización:** 15 de marzo de 2026
