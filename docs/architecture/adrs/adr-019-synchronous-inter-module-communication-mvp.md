# ADR-019 – Comunicación Síncrona entre Módulos para el MVP

**Fecha:** 07/02/2026
**Estado:** Aceptado
**Autores:** Gemini CLI (acting as Lead Architect)

---

## 1. Contexto

Durante la fase de diseño y análisis de las interacciones entre los módulos del 'ERP Core' (ej., Sales, Product) y el subdominio MES, se planteó la necesidad de una comunicación eficiente y desacoplada, especialmente para la formalización de pedidos que inician flujos de producción. Se evaluaron dos patrones principales: llamadas a API síncronas y mensajería asíncrona mediante un Message Broker.

Aunque la mensajería asíncrona ofrece ventajas significativas en términos de desacoplamiento, resiliencia y escalabilidad a largo plazo (alineándose con la preparación para microservicios del ADR-002 y ADR-003), la introducción de un Message Broker añade una capa de complejidad considerable en infraestructura, gestión operacional y curva de aprendizaje para el equipo, lo cual podría impactar la agilidad y el tiempo de entrega del Producto Mínimo Viable (MVP).

---

## 2. Alternativas Consideradas

1.  **Implementar Mensajería Asíncrona (Message Broker) desde el MVP:**
    *   Ventajas: Desacoplamiento óptimo, alta resiliencia, escalabilidad desde el inicio, mejor alineación con la visión de microservicios.
    *   Desventajas: Mayor complejidad de infraestructura, mayor curva de aprendizaje para el desarrollo, gestión operacional adicional para el MVP, posible retraso en la entrega.
2.  **Implementar Comunicación Síncrona Directa para el MVP (Decisión Adoptada):**
    *   Ventajas: Menor complejidad de infraestructura y desarrollo inicial, mayor agilidad para el MVP, simplicidad operacional para un entorno 'local-first'.
    *   Desventajas: Mayor acoplamiento temporal entre módulos, menor resiliencia frente a fallos de un componente, puede requerir refactorización post-MVP para introducir asincronía.

---

## 3. Criterios de Decisión

-   Minimizar la complejidad de infraestructura y desarrollo para el MVP.
-   Asegurar la entrega ágil del MVP.
-   Mantener la simplicidad operacional para un despliegue 'local-first'.
-   No bloquear la evolución futura hacia arquitecturas más desacopladas.

---

## 4. Decisión Adoptada

Para el Producto Mínimo Viable (MVP), la comunicación entre los módulos del 'ERP Core' y MES (y entre otros Bounded Contexts dentro del monolito modular) se realizará de manera **síncrona mediante llamadas directas a los servicios de aplicación** expuestos por cada módulo.

**Específicamente para el flujo 'Sales' -> 'MES':**
Cuando un pedido se formalice en `Sales`, el servicio de aplicación de `Sales` invocará directamente un método del `ProductionPlanningService` de MES. Ambas invocaciones ocurrirán dentro del mismo proceso del monolito.

Esta decisión se documenta con la clara intención de que la introducción de un Message Broker para la comunicación asíncrona sea una **mejora planificada para la fase post-MVP**, a medida que el proyecto madure y los requisitos de resiliencia y escalabilidad justifiquen la complejidad adicional.

---

## 5. Consecuencias

### Positivas

-   Reducción significativa de la complejidad de infraestructura para el MVP.
-   Desarrollo más rápido para la fase inicial del proyecto.
-   Facilidad de depuración debido a un flujo de ejecución más directo.
-   Menor sobrecarga operacional para un entorno 'local-first'.

### Negativas

-   Mayor acoplamiento temporal entre los módulos.
-   Menor resiliencia: la falla de un módulo llamado puede impactar directamente al módulo llamador.
-   Potencial riesgo de cuello de botella si las operaciones síncronas se vuelven costosas.
-   Necesidad de refactorización futura para introducir el Message Broker, lo que implica un "costo de transición" a largo plazo.

---

## 6. Alcance

-   Estrategia de comunicación inter-módulos para el MVP.
-   Justificación para la postergación del uso de Message Brokers.

---

## 7. Integración con otros ADRs

-   ADR-002: Adopción de Clean Architecture y DDD con Rigor Asimétrico (La decisión mantiene la separación lógica de Bounded Contexts, pero relaja el desacoplamiento temporal para el MVP).
-   ADR-003: Tipo y Distribución de la Aplicación (Monolito Modular Local-First) (La decisión apoya la simplicidad operacional del monolito para el MVP).
-   ADR-018: MES Module Architecture - Independent Schema and Inter-module Communication (La decisión complementa este ADR al definir el *cómo* se comunicarán los módulos en el MVP, manteniendo la separación lógica de esquemas y la comunicación por interfaces).

---

## 8. Notas Adicionales / Consideraciones Especiales

Se recomienda que las llamadas síncronas se realicen a través de interfaces de servicios de aplicación bien definidas, para facilitar la futura transición a un modelo asíncrono. Los DTOs de comunicación deben ser concisos y no exponer detalles internos de la implementación de cada módulo.

---

## 9. Referencias
