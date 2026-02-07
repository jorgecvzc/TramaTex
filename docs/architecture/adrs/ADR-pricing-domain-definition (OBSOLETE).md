# ADR-014 – Definición del Dominio del Módulo de Precios

**Fecha:** 2026-02-05  
**Estado:** Obsoleto
**Autores:** Gemini CLI

---

## 1. Contexto

Este ADR ha sido **sustituido** por el más completo [ADR-016: Arquitectura del Módulo de Pricing](./ADR-016-pricing-module-architecture.md). Se mantiene por razones históricas.

El proyecto TramaTex está concebido como un monolito modular, aplicando los principios de Arquitectura Limpia (ADR-002) y Domain-Driven Design (ADR-006). El módulo de Precios es identificado como **CRÍTICO** debido a su naturaleza económicamente sensible y su impacto directo en la rentabilidad del negocio.

La estrategia de DDD adoptada en TramaTex requiere un rigor estricto en módulos de negocio core, y el módulo de Precios es el ejemplo más prominente, demandando un 100% de cobertura de tests unitarios y una arquitectura rigurosa para su dominio.

Este ADR tiene como objetivo establecer la definición del dominio para el módulo de Precios, sentando las bases para su implementación. Se delinearán las entidades clave, Value Objects, servicios de dominio y reglas de negocio principales, asegurando su alineación con la arquitectura general del proyecto y las necesidades de negocio identificadas.

El módulo de Precios dependerá directamente de la información proporcionada por los módulos IAM (para seguimiento de usuario), Party (referencias de clientes) y Product (costos de variantes).

---

## 2. Revisión Detallada de Alternativas
##### Alternativa A: Extender las Tablas del Módulo `Product`

Esta opción implicaría añadir columnas directamente a las tablas existentes del módulo `Product` (ej., `products`, `product_variants`) para almacenar información de precios como `margin_percentage_brand`, `client_pricing_rule_id`, o incluso `selling_price_calculated`.

**Argumentos a Favor (y por qué son generalmente débiles en TramaTex):**
*   **Simplicidad Inicial (Percepción engañosa):** Parece más fácil "añadir una columna" que crear nuevas tablas y un nuevo contexto de base de datos.
*   **Consultas Unificadas (Aparente):** En teoría, se podría obtener toda la información de producto y precio en una sola consulta al módulo `Product`.

**Argumentos en Contra (basados en principios de TramaTex y complejidad/latencia):**
*   **Violación de Bounded Contexts (CRÍTICO):** El módulo `Product` debe preocuparse *solamente* de la definición del producto (qué es, sus variantes, características). La lógica de *cómo se vende y a qué precio* pertenece al contexto de `Pricing`. Mezclar estos conceptos en las tablas del módulo `Product` rompe la cohesión del contexto `Product` y acopla fuertemente `Product` a `Pricing`. Esto va directamente en contra de `ADR-002` y `ADR-006` (Clean Architecture y DDD).
*   **Acoplamiento Fuerte y Dependencia Circular:**
    *   Si `Pricing` modifica las tablas de `Product`, entonces `Product` depende de `Pricing` (para sus datos) y `Pricing` depende de `Product` (para los datos base del producto). Esto crea un acoplamiento circular que es extremadamente difícil de mantener y evolucionar.
    *   Cualquier cambio en la lógica de precios o en cómo se guardan las reglas requeriría modificar las tablas de `Product`, forzando al módulo `Product` a saber demasiado sobre el módulo `Pricing`.
*   **Mantenibilidad y Escalabilidad Reducida:**
    *   Las tablas de `Product` se volverían "gigantes" y difusas en su propósito, dificultando entender qué dato pertenece a qué contexto.
    *   Impediría la extracción futura del módulo `Pricing` como un servicio independiente, ya que sus datos estarían esparcidos y mezclados con los de `Product`. (Principio de Monolito Modular de `ADR-003`).
*   **Complejidad de Cálculo y Latencia (Negativa):** Si el precio calculado se almacena en `ProductVariant`, cada vez que se recalcule por un cambio en las reglas de `Pricing`, se estaría modificando el registro de `ProductVariant`. Esto puede generar latencia para asegurar la consistencia transaccional y bloqueos, e introduce la *complejidad* de gestionar estas actualizaciones en un módulo ajeno. La complejidad de cálculo no se reduce, solo se desplaza y mezcla.
*   **Inmutabilidad de Datos de Producto (Ideal):** El precio de venta calculado no es una propiedad intrínseca y estática del producto.

##### Alternativa B: Crear Nuevas Tablas Dedicadas al Módulo `Pricing`

Esta opción implica diseñar un esquema de base de datos específico para el módulo `Pricing`. Estas tablas estarían bajo la propiedad exclusiva del módulo `Pricing`.

**Argumentos a Favor (basados en principios de TramaTex y complejidad/latencia):**
*   **Coherencia de Bounded Contexts (CRÍTICO):** Cada módulo (contexto acotado) es dueño de su propio modelo de dominio y de sus datos, lo que es fundamental para DDD.
*   **Acoplamiento Débil:** Mínimas dependencias.
*   **Mantenibilidad y Extensibilidad Óptimas:** Libertad para evolucionar el esquema de `Pricing` sin afectar `Product`.
*   **Claridad del Modelo de Datos:** Las tablas reflejan directamente el dominio de `Pricing`.
*   **Preparación para Microservicios (a largo plazo).**
*   **Inmutabilidad de los Cálculos de Precio.**
*   **Gestión Clara de la Complejidad de Cálculo:** Toda la lógica y los datos complejos del cálculo de precios residen *dentro* del módulo `Pricing` y sus propias tablas. Esto encapsula la complejidad, facilitando su comprensión, desarrollo y testeo de forma aislada.
*   **Base para la Optimización de Latencia:** Al tener los datos de reglas de precio bien estructurados y accesibles dentro de su propio contexto, el módulo `Pricing` puede optimizar las consultas y cálculos para reducir la latencia inherente a la complejidad del cálculo. Además, esto sienta las bases para las estrategias de caching que abordan la latencia de forma más eficaz.

---

## 3. Criterios para la "Mejor" Solución (según TramaTex)

La "mejor" solución debe maximizar:
*   **Separación de Bounded Contexts.**
*   **Bajo Acoplamiento y Alta Cohesión.**
*   **Mantenibilidad y Extensibilidad.**
*   **Asymmetric Rigor.**
*   **Optimización para la Complejidad de Cálculo:** La solución debe permitir gestionar la complejidad intrínseca del cálculo de precios de forma encapsulada y eficiente, sin que esta complejidad se propague a otros módulos.
*   **Minimización del Impacto en Latencia (para el usuario final):** Aunque el cálculo de precios puede ser complejo, el sistema debe entregar los resultados con una latencia aceptable para la experiencia del usuario, especialmente en escenarios de alta frecuencia de lectura.
*   **Coherencia de Datos.**

---

## 4. Decisión Adoptada

**Recomendación:** Crear **nuevas tablas dedicadas** y exclusivas para la persistencia de la información del módulo `Pricing`.

**Justificación Detallada (reforzando complejidad y latencia):**

Esta opción es superior porque:
1.  **Respeta Fundamentalmente los Bounded Contexts y el Bajo Acoplamiento.**
2.  **Encapsula la Complejidad de Cálculo:** La lógica de precios es inherentemente compleja (reglas, variantes, marcas, descuentos, etc.). Al tener tablas dedicadas, toda esta complejidad se maneja dentro del módulo `Pricing` y su base de datos. Esto previene que la complejidad se filtre a otros módulos, permitiendo que `Product` se mantenga simple y centrado en la definición del producto, y `Sales` en la gestión de pedidos.
3.  **Optimiza la Mantenibilidad y Extensibilidad:** Cualquier cambio en la compleja lógica de cálculo de precios se puede realizar en el módulo `Pricing` sin afectar la base de datos ni la lógica del módulo `Product`.
4.  **Complementa de Forma Crucial la Estrategia de Caching Híbrida para Abordar la Latencia:**
    *   La complejidad del cálculo de precios (múltiples reglas, variantes, marcas, descuentos) significa que el cálculo "en caliente" *puede* tener una latencia significativa si se hace en cada solicitud.
    *   La estrategia de persistencia en tablas dedicadas, combinada con la **caché en memoria (NoSQL)**, es la forma óptima de gestionar esta latencia. Las tablas de `Pricing` actúan como la fuente de verdad para las reglas y configuraciones (que cambian raramente). El `SellingPriceCalculatorService` realiza el cálculo completo solo cuando es necesario (cache miss o invalidación).
    *   Una vez calculado, el precio se almacena en la caché (BD NoSQL en memoria) para ofrecer una **latencia mínima** en las lecturas subsiguientes. La precarga de variantes de un mismo producto base en la caché mejora aún más la latencia percibida para escenarios comunes de uso.
5.  **Alineación con el "Asymmetric Rigor":** Un módulo CRÍTICO como `Pricing` requiere una solución robusta y dedicada, que le permita manejar su complejidad y optimizar su rendimiento de manera autónoma.

---

## 5. Consecuencias

### Positivas
-   **Separación de Responsabilidades clara y fuerte encapsulación del dominio de `Pricing`.**
-   **Alta testabilidad y mantenibilidad:** La complejidad del cálculo de precios está contenida.
-   **Flexibilidad para la evolución:** El esquema de `Pricing` puede cambiar independientemente.
-   **Rendimiento optimizado para el usuario final:** Gracias a la combinación de tablas dedicadas y la estrategia de caching híbrida.
-   **Alineación con los principios de Clean Architecture y DDD.**

### Negativas
-   **Aumento de la infraestructura:** Requiere gestionar un nuevo conjunto de tablas de base de datos para el módulo `Pricing`.
-   **Complejidad en la gestión de la caché:** Aunque simplificada por la invalidación completa, sigue requiriendo atención.
-   **Requiere "uniones" lógicas:** Para combinar información de `Product` con precios, se necesita orquestación en la capa de aplicación.

---

## 6. Alcance

Esta decisión aplica al diseño y desarrollo del Bounded Context de Precios dentro de la aplicación `tramatex-api`. Influye directamente en la definición de sus interfaces de dominio (repositorios y servicios) y cómo otros módulos (ej. Sales) interactuarán con él.

---

## 7. Integración con otros ADRs

-   **ADR-002: Adopción de Clean Architecture y DDD:** Este ADR es una aplicación directa de los principios establecidos en ADR-002, reforzando la elección de capas y la separación de preocupaciones.
-   **ADR-006: Estrategia de Domain-Driven Development:** Sigue la estrategia de aplicar rigor asimétrico, con un rigor estricto en este módulo crítico.
-   **ADR-011: Estrategia de Cobertura de Tests:** Establece el compromiso de 100% de cobertura de tests unitarios para el dominio de Precios.
-   **Bounded Contexts (agents/project/context/bounded-contexts.yaml):** Complementa la descripción de los elementos clave y dependencias del contexto de Precios.
-   **Arquitectura (agents/project/context/architecture.yaml):** Refuerza la estructura de capas y el patrón de repositorio para este módulo.

---

## 8. Notas Adicionales / Consideraciones Especiales

-   Es crucial definir las interfaces de los repositorios de `PricingRule`, `ClientPricing`, `Margin`, `Discount` y `PriceCalculation` en la capa de dominio.
-   Las implementaciones de estos repositorios residirán en la capa de infraestructura.
-   El `PriceCalculationService` será el principal punto de entrada para los casos de uso relacionados con el cálculo de precios.
-   Se debe considerar un mecanismo robusto para el versionado de reglas de precios si cambian frecuentemente y se requiere aplicar reglas históricas a cálculos pasados.
-   **Estrategia de Caching:** Cuando se solicite el precio de venta de una `VariantProduct`, el sistema de caché debería intentar precargar y almacenar en caché los precios de venta de *todas* las `VariantProduct` asociadas al mismo `ProductBase` para optimizar el rendimiento en escenarios de venta de múltiples variantes.
-   **Invalidación de Caché:** Dada la baja frecuencia de cambios en los precios de productos, modificadores de variantes, márgenes de marca o reglas de descuento (aproximadamente una o dos veces al año), la estrategia para asegurar la coherencia de los precios en caché será la **invalidación completa de toda la caché de precios**. Cuando se detecte un cambio en cualquiera de los datos que influyen en el cálculo del precio de venta (ej. `ProductBasePrice`, `VariantModifier`, `BrandProfitMargin`, `SalesDiscountRule`, `ClientPricing` o `PricingRule`), se activará un mecanismo que limpiará por completo la caché de precios. Esto simplifica enormemente la implementación del mecanismo de invalidación, ya que no requiere un seguimiento granular de las dependencias de cada elemento cacheado, y el impacto de rendimiento de una reconstrucción completa de la caché es aceptable dada la poca frecuencia de los cambios.

---

## 9. Referencias

-   Contexto del proyecto (`agents/project/project-context.yaml`)
-   Contexto de Arquitectura (`agents/project/context/architecture.yaml`)
-   Contexto de Bounded Contexts (`agents/project/context/bounded-contexts.yaml`)
-   Documentación de DDD y Clean Architecture.
