# 🏛️ ADR-019: Comunicación Síncrona entre Módulos para el MVP

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 07-02-2026 |
| **Autores** | Gemini CLI (Architect) |

---

## 🎯 Contexto
Se requiere una comunicación eficiente entre el núcleo del ERP (Ventas, Productos) y el subdominio MES (Taller). Se evaluaron dos patrones: llamadas a API síncronas y mensajería asíncrona mediante un Message Broker (como RabbitMQ o Redis Streams).

---

## 🔍 Alternativas Consideradas
1. **Mensajería Asíncrona (Message Broker):** Desacoplamiento óptimo y alta resiliencia, pero introduce una complejidad operativa y de infraestructura excesiva para el MVP.
2. **Comunicación Síncrona Directa (Decisión Adoptada):** Menor complejidad inicial, mayor agilidad y simplicidad operativa para un entorno "local-first".

---

## ✅ Decisión Adoptada
Para el MVP, la comunicación entre Bounded Contexts dentro del monolito modular se realizará de forma **síncrona mediante llamadas directas a los servicios de aplicación**.

### Reglas de Implementación:
*   **Acoplamiento Controlado:** El servicio de Ventas invocará directamente métodos del `ProductionPlanningService` de MES dentro del mismo proceso.
*   **Interfaces Claras:** Las llamadas deben realizarse a través de interfaces bien definidas para facilitar una futura migración a un modelo asíncrono.
*   **Post-MVP:** La introducción de un Message Broker se planifica como una mejora futura cuando el volumen de operaciones o los requisitos de resiliencia lo justifiquen.

---

## 📈 Consecuencias
### Positivas
*   Reducción drástica de la complejidad de infraestructura.
*   Desarrollo y depuración más rápidos.
*   Menor sobrecarga operacional para el despliegue local.

### Negativas
*   Mayor acoplamiento temporal entre módulos.
*   Menor resiliencia: un fallo en el taller puede impactar directamente en la confirmación de la venta.

---
[Volver al Índice de ADRs](./README.md)
