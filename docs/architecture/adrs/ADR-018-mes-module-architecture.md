# ADR-018 – MES Module Architecture - Independent Schema and Inter-module Communication

**Fecha:** 07/02/2026
**Estado:** Propuesto
**Autores:** Gemini CLI (acting as Lead Architect)

---

## 1. Contexto

El módulo MES (Manufacturing Execution System) es identificado como un subdominio crítico para la planificación, ejecución y control de calidad de la producción. Aunque inicialmente se integrará en un monolito modular, existe un requisito claro de diseñarlo para una eventual extracción a un microservicio separado. Esto requiere una consideración cuidadosa de su esquema de base de datos y los patrones de comunicación con otros contextos delimitados (Bounded Contexts).

Este ADR se basa en los principios de Clean Architecture y DDD con rigor asimétrico (ADR-002) y la estrategia de Monolito Modular (ADR-003), buscando extender estos principios a la granularidad del diseño de la base de datos y la comunicación inter-módulos.

---

## 2. Alternativas Consideradas

1.  **MES con esquema de BD compartido y acceso directo:**
    *   Ventajas: Implementación inicial más rápida.
    *   Desventajas: Alto acoplamiento con otros módulos, dificulta la extracción futura del microservicio, riesgo de violación de la integridad del dominio MES por otros módulos, inconsistencia con los principios de Bounded Contexts.
2.  **MES con esquema de BD independiente y comunicación exclusiva por interfaces (Decisión Adoptada):**
    *   Ventajas: Bajo acoplamiento, alta cohesión del módulo MES, facilita enormemente la extracción futura a microservicio, protege la integridad del dominio MES.
    *   Desventajas: Requiere un diseño más cuidadoso de las interfaces de comunicación, posible latencia adicional en ciertas operaciones inter-módulos, requiere disciplina para evitar atajos.

---

## 3. Criterios de Decisión

-   Facilitar la transición futura a microservicios.
-   Proteger la integridad y autonomía del dominio MES.
-   Mantener un bajo acoplamiento entre módulos.
-   Alinear con los principios de DDD y Clean Architecture ya adoptados.

---

## 4. Decisión Adoptada

El módulo MES se adherirá a los siguientes principios arquitectónicos para facilitar su futura implementación independiente y potencial despliegue como microservicio:

1.  **Esquema Lógico Independiente de Base de Datos:**
    *   El módulo MES gestionará su propio esquema lógico de base de datos.
    *   Aunque físicamente residirá dentro de la base de datos PostgreSQL compartida del monolito (según ADR-003), sus tablas se agruparán lógicamente (por ejemplo, utilizando un nombre de esquema dedicado o prefijos de tabla claros) y serán gestionadas exclusivamente por el módulo MES.
2.  **No Acceso Directo a la Base de Datos Trans-Contexto:**
    *   Otros módulos no accederán directamente a las tablas de MES (a través de `JOIN`s SQL o consultas directas) y viceversa. Esta regla es estricta para garantizar la autonomía de cada Bounded Context.
3.  **Comunicación Inter-Módulos Explícita y por Contratos:**
    *   La comunicación entre MES y otros Bounded Contexts (ej., Sales, Product, Inventory) se realizará exclusivamente a través de interfaces bien definidas a nivel de aplicación. Estas interfaces pueden incluir:
        *   **Servicios de Aplicación/APIs:** Para consultas síncronas de datos o ejecución de comandos (ej., Sales consulta a MES el estado de una ProductionOrder).
        *   **Eventos de Dominio:** Para comunicación asíncrona para la sincronización de datos no críticos o que puedan ser eventualmente consistentes (ej., MES publica un `ProductionOrderCompleted` que Inventory consume).
4.  **Diseñado para Ser Extraíble (Extractable by Design):**
    *   La estructura interna del módulo (capas de Clean Architecture) y sus contratos externos (APIs, eventos) serán diseñados para permitir su despliegue independiente como un microservicio sin requerir una re-arquitectura significativa de su lógica interna o de sus consumidores/productores.

---

## 5. Consecuencias

### Positivas

-   Refuerza los límites fuertes entre módulos, reduciendo el acoplamiento.
-   Facilita enormemente la transición futura a microservicios con menos esfuerzo y riesgo.
-   Mejora la mantenibilidad y la claridad de la propiedad de los datos y la lógica de MES.
-   Aumenta la testabilidad de los módulos individuales.

### Negativas

-   Requiere un diseño más cuidadoso de los contratos de comunicación y las estrategias de sincronización de datos.
-   Potencial de mayor latencia en la recuperación de datos entre módulos si no se optimiza (ej., mediante caché, desnormalización con eventos, o APIs bien diseñadas).
-   El desarrollo inicial puede requerir un mapeo más explícito entre contextos.

---

## 6. Alcance

-   Diseño de la base de datos del módulo MES.
-   Definición de las APIs y eventos de comunicación de MES.
-   Implementación del módulo MES.

---

## 7. Integración con otros ADRs

-   ADR-002: Adopción de Clean Architecture y DDD con Rigor Asimétrico
-   ADR-003: Tipo y Distribución de la Aplicación (Monolito Modular Local-First)

---

## 8. Notas Adicionales / Consideraciones Especiales

Este ADR establece las pautas fundamentales para el diseño de la base de datos y la comunicación del módulo MES, asegurando que se construya con la visión de un futuro independiente. Cualquier interacción de otros módulos con MES (y viceversa) debe adherirse a estos principios.
