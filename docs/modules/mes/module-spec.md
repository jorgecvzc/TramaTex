# Módulo de MES (Manufacturing Execution System)

**Estado:** 🔄 **EN REFACTORIZACIÓN**  
**Última actualización:** 14 de marzo de 2026

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
- **UI:** "Tareas" en sección Datos Maestros.

### 2) Posición (`Position`)
- **Propósito:** Zona de la prenda donde se realiza el trabajo (ej. Pecho izquierdo, Espalda, Bajo pantalón).
- **UI:** "Posiciones" en sección Datos Maestros.

### 3) Tipo de Trabajo (`WorkType`)
- **Propósito:** Secuencia ordenada de tareas que define un tipo de marcado/personalización. Es la "receta" del proceso.
- **Ejemplo:** "Marcado por vinilo" = Diseñar → Imprimir → Marcar → Plegar → Embolsar.
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
- `DRAFT`: En preparación.
- `PENDING`: Lista para iniciar.
- `IN_PROGRESS`: Al menos una tarea ha comenzado.
- `COMPLETED`: Todas las tareas finalizadas.
- `ON_HOLD`: Pausada por incidencia.
- `CANCELLED`: Descartada.

### Estados de Tarea (`WorkOrderTask`)
- `PENDING`, `IN_PROGRESS`, `COMPLETED`, `BLOCKED`, `SKIPPED`.

### Prioridades (`WorkOrder`)
- `LOW`, `NORMAL`, `HIGH`, `URGENT`.

---

## 4. Componentes

### Backend
- **Maestros de Taller:** CRUD de Tareas (`Task`) y Posiciones (`Position`).
- **Tipos de Trabajo:** CRUD de `WorkType` con definición de secuencias de tareas.
- **Configuraciones:** CRUD de `WorkSetup` con líneas (WorkType + Position).
- **Órdenes de Trabajo:** Gestión de `WorkOrder`, instanciación desde `WorkSetup`, transiciones de estado.
- **Persistencia de Diseño:** Soporte para archivos de diseño y notas de prenda.

### Frontend
- **Dashboard de Producción:** Vista de seguimiento de órdenes activas.
- **Terminal de Taller:** Interfaz táctil para operarios (Start/Complete/Block).
- **Datos Maestros:** Gestión de Tareas, Posiciones y Tipos de Trabajo.
- **Configuraciones:** Definición de personalizaciones por cliente/prenda.
