# 🏛️ ADR-022: Estrategia de Extracción de MES a Microservicio

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | 🚀 Propuesta Post-MVP |
| **Autor** | Gemini CLI (Arquitectura) |

---

## 🎯 Contexto
Siguiendo los principios de "Extraíble por Diseño" definidos en el [ADR-018](../architecture/adrs/adr-018-mes-module-architecture.md), se plantea la ruta técnica para separar físicamente el módulo MES del monolito modular y convertirlo en un microservicio autónomo. Esta decisión es estratégica para permitir el escalado independiente de la planta de producción y una mayor resiliencia ante picos de demanda en el taller.

---

## 🔍 Alternativas de Extracción

### A. Extracción en Caliente (Refactorización Directa)
*   **Descripción:** Separar el código y la base de datos en un único hito de desarrollo.
*   **Desventajas:** Alto riesgo de interrupción del servicio y dificultad para revertir cambios.

### B. Extracción Progresiva por "Strangler Fig" (Adoptada)
*   **Descripción:** Mantener la lógica en el monolito mientras se construye la infraestructura del microservicio, desviando tráfico progresivamente.
*   **Ventajas:** Riesgo controlado, validación continua y posibilidad de *rollback* instantáneo.

---

## ✅ Hoja de Ruta Propuesta (Post-MVP)

### 1. Migración de Comunicación
*   Sustituir las llamadas síncronas directas (Service-to-Service) por un **Message Broker** (ej: RabbitMQ o Redis Streams).
*   Introducir eventos de dominio para la sincronización asíncrona de datos (ej: `SalesOrderConfirmed` -> `ProductionOrderCreated`).

### 2. Separación Física de Datos
*   Migrar las tablas de MES a una base de datos independiente (PostgreSQL).
*   Eliminar cualquier dependencia técnica residual (foreign keys trans-contexto).

### 3. Despliegue Autónomo
*   Contenerización individual del microservicio MES.
*   Implementación de un **API Gateway** para orquestar las peticiones entre el ERP y el Taller.

---

## 📈 Consecuencias
### Positivas
*   **Escalabilidad:** El taller puede escalar sus recursos sin afectar al área comercial.
*   **Autonomía:** Ciclos de despliegue independientes para mejoras en la terminal de operario.
*   **Resiliencia:** Un fallo en el ERP no detiene la operación de las máquinas del taller (Local-First distribuido).

### Negativas
*   Aumento de la complejidad operativa (gestión de red, latencia).
*   Necesidad de gestionar la consistencia eventual entre sistemas.

---
[Volver al Índice de ADRs](../architecture/adrs/README.md)
