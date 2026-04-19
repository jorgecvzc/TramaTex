# 🏭 Módulo de MES (Gestión de Taller)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Vigente |
| **Referencia** | [ADR-018](../../architecture/adrs/adr-018-mes-module-architecture.md) |

---

## 🎯 Propósito
El módulo MES (Manufacturing Execution System) es el corazón operativo del taller. Transforma las necesidades comerciales en órdenes de producción ejecutables, permitiendo el seguimiento en tiempo real de las tareas de personalización, marcaje y confección, asegurando la calidad en cada etapa.

---

## 📄 Ramas del Conocimiento (Documentación)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md) — Órdenes de trabajo, estaciones y tareas.
*   **Casos de Uso:** [use-cases.md](./use-cases.md) — Planificación de producción y terminal de operario.
*   **Contratos de API:** [api-contracts.md](./api-contracts.md) — Endpoints de taller y estados de ejecución.
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md) — Desacoplamiento de base de datos y lógica de estaciones.

---

## 🏗️ Componentes Clave
*   **Entidades:** `WorkOrder` (Orden de Trabajo), `Task` (Tarea operativa), `WorkStation` (Centro de trabajo).
*   **Sincronización:** Conectado síncronamente con Sales para la recepción de demandas de fabricación.
*   **Terminal Industrial:** Interfaz optimizada para operarios con feedback visual inmediato.

---
[Volver al Resumen de Módulos](../README.md)
