# Plan de Ejecución Post-MVP — TramaTex

Este documento define la estrategia de implementación, priorización y flujo de trabajo para las mejoras y nuevos módulos planificados tras la finalización del MVP.

---

## 1. Fases de Ejecución

La ejecución se divide en cuatro fases lógicas, priorizando la estabilidad del núcleo y las dependencias técnicas antes de la expansión funcional masiva.

### Fase 1: Cimientos y Ergonomía (Sprints 19-21)
*Objetivo: Mejorar la experiencia de usuario y preparar la infraestructura para la asincronía.*

1.  **Unificación UI/UX (Hito 1)**: Implementación del Plan Maestro de Diseño.
2.  **Infraestructura de Mensajería (Hito 5)**: Despliegue de NATS JetStream y patrón Transactional Outbox.
3.  **Elevación de Calidad (Hito 12)**: Incremento sistemático de cobertura de tests en dominio y aplicación.

### Fase 2: Expansión de Lógica de Negocio (Sprints 22-25)
*Objetivo: Completar el ciclo financiero y mejorar la integración Sales-MES.*

1.  **Gestión de Cobros y Tesorería (Hito 3)**: Vencimientos y seguimiento de deuda.
2.  **Facturación Consolidada (Hito 2)**: Agrupación de albaranes.
3.  **Integración Profunda Sales ↔ MES (Hito 14)**: Sincronización automática de cambios.

### Fase 3: Estructura y Cumplimiento (Sprints 26-29)
*Objetivo: Cumplimiento legal y escalabilidad arquitectónica.*

1.  **Facturación Electrónica (Hito 4)**: Adaptación a la Ley Crea y Crece (XML/Veri*factu).
2.  **Extracción del MES como Microservicio (Hito 6)**: Desacoplamiento físico basado en la infraestructura de la Fase 1.
3.  **Gestión Avanzada de Archivos (Hito 13)**: Thumbnails y visualización técnica.

### Fase 4: Optimización y Analítica (Sprints 30+)
*Objetivo: Rendimiento avanzado y visión de negocio.*

1.  **Caché y Rendimiento (Hito 7)**: Implementación de Redis (especialmente en Pricing).
2.  **Búsqueda Global (Hito 9)**: Ctrl+K y Full-Text Search.
3.  **BI y Analítica (Hito 8)**: Dashboards de rentabilidad y KPIs.
4.  **Planificación de Operarios (Hito 15)**: Gestión de carga de trabajo en el MES.

---

## 2. Metodología de Implementación y Flujo de Trabajo

Los documentos de estrategia en `docs/post-mvp/` son **declaraciones de intenciones**. Para su ejecución, cada hito debe ser procesado a través del siguiente flujo de trabajo riguroso:

### 2.1. Fase de Definición (Pre-Desarrollo)
1.  **Creación del Sprint**: Definir el ID del sprint (ej: Sprint 19) y crear su directorio en `docs/log/sprints/`.
2.  **Desglose de Tareas**: Utilizar la plantilla `_task-template.md` para documentar cada tarea necesaria para completar el hito.
3.  **Evaluación de Cobertura**: Establecer los objetivos de test coverage específicos para el nuevo código (mínimo 85% en Dominio).

### 2.2. Ciclo de Implementación Atómica (Por Tarea)
Para cada tarea dentro del sprint, se seguirá estrictamente este orden:

1.  **Backend (Lógica de Negocio)**:
    *   Implementar Entidades y Servicios de Dominio.
    *   Validar mediante Tests Unitarios (TDD).
2.  **Capa de Persistencia e Infraestructura**:
    *   Comprobar si el cambio afecta a la persistencia (GORM models, esquemas PostgreSQL).
    *   Ejecutar y validar migraciones si fuera necesario.
    *   Implementar Repositorios y validar con Tests de Integración.
3.  **Frontend y UI**:
    *   Solo una vez que el backend es estable y está testeado, se procede a la implementación en Vue.js.
    *   Ajuste a los estándares de ergonomía e iconografía del Plan Maestro UI/UX.
4.  **Validación de Integración**:
    *   Comprobación del flujo completo (E2E) y verificación manual de la experiencia de usuario.

---

## 3. Matriz de Dependencias Técnicas

| Hito | Dependencia Técnica | Razón |
|------|----------------------|-------|
| **Facturación Electrónica** | Cobros y Tesorería | Requiere estados de pago precisos para reportes legales. |
| **Extracción MES** | NATS JetStream | Necesario para comunicación asíncrona entre servicios. |
| **Caché (Pricing)** | Tests de Dominio | El motor de precios debe estar 100% blindado antes de cachear. |
| **Asignación Operarios** | Integración Sales-MES | Requiere que los datos de fabricación sean consistentes. |

---
*Última actualización: 2026-04-29*
