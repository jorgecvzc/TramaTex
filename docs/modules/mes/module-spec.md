# Módulo de MES (Manufacturing Execution System)

**Estado:** ✅ **COMPLETO (100%)**  
**Última actualización:** 1 de marzo de 2026

## 1. Propósito

*   **Visión del Módulo:** Controlar y monitorizar el proceso de fabricación en el taller de TramaTex, asegurando la trazabilidad desde la orden de venta hasta el producto terminado.
*   **Objetivos Clave:**
    *   Proporcionar visibilidad en tiempo real del estado de la producción.
    *   Gestionar plantillas de procesos reusables para diferentes tipos de trabajos.
    *   Optimizar la interacción de los operarios mediante un terminal de taller intuitivo.

---

## 2. Definiciones Clave (Nomenclatura)

### 1) Trabajo Definido (Plantilla)
- **Nombre técnico:** `MESWorkDefinition`
- **Propósito:** Define *qué* hay que hacer (plantilla reusable por cliente/producto).
- **Incluye:** Secuencia de fases, tareas y parámetros por defecto.

### 2) Trabajo Real (Ejecución)
- **Nombre técnico:** `MESWorkExecution`
- **Propósito:** Representa una instancia operativa de producción (instancia real).
- **Incluye:** Referencia a la definición, fechas reales, lotes y estados de tareas.

---

## 3. Componentes Implementados ✅

### Backend
- **Gestión de Definiciones:** CRUD completo de `WorkDefinitions` y sus fases/tareas asociadas.
- **Motor de Ejecución:** Lógica de transiciones de estado (START, PAUSE, COMPLETE, BLOCK).
- **Integración Sales:** Generación automática de `WorkExecution` al aceptar un pedido en Sales.
- **Trazabilidad:** Registro de tiempos y operarios por tarea.

### Frontend
- **Dashboard de Producción:** Vista Kanban para el seguimiento global de trabajos.
- **Terminal de Taller (Tablet):** Interfaz simplificada para operarios con botones de acción rápida.
- **Gestión de Plantillas:** UI para la creación y edición de flujos de trabajo reusables.

---

## 4. Decisiones de Diseño

*   **Arquitectura Extraíble:** Diseñado siguiendo Clean Architecture para permitir su futura extracción como microservicio independiente (ver ADR-018 y ADR-022).
*   **Independencia de Esquema:** Utiliza su propio esquema lógico en la base de datos para evitar acoplamiento con el Core ERP.
*   **Interfaz de Operario Simplificada:** Foco en minimizar la introducción de datos manual mediante el uso de estados predefinidos.

---

## 5. Fases de Desarrollo

*   [x] **Fase 1:** Cimentación y arquitectura del módulo.
*   [x] **Fase 2:** CRUD de maestros y definiciones de trabajo.
*   [x] **Fase 3:** Motor de estados y terminal de taller (Tablet).
*   [x] **Fase 4:** Dashboard de seguimiento e integración con Sales.
