# ADR-016 – Arquitectura del Módulo de Pricing

**Fecha:** viernes, 6 de febrero de 2026  
**Estado:** Aceptado  
**Autores:** Gemini CLI

---

## 1. Contexto

El sistema TramaTex necesita un módulo robusto para gestionar la lógica de pricing de productos, tanto para el cálculo de precios base de venta como para la aplicación de modificaciones en el momento de la venta (descuentos promocionales). La necesidad surge de la complejidad inherente a la determinación de precios en el sector textil, incluyendo márgenes, descuentos por volumen, promociones y la diversidad de productos y variantes.

**Problemas a Resolver:**
- Garantizar la coherencia y exactitud en el cálculo de precios.
- Optimizar el rendimiento de cálculo de precios para una experiencia de usuario fluida, especialmente en el contexto de ventas.
- Permitir la flexibilidad en la definición de reglas de pricing.
- Mantener la pureza del dominio de Pricing, desacoplándolo de preocupaciones de infraestructura.
- Soporte para la creación Just-in-Time (JIT) de ProductVariants.

---

## 2. Alternativas Consideradas

**Alternativa A – Unificar todas las reglas en una única entidad `PricingRule` con un campo `rule_category`:**
- Ventajas: Menos entidades en el dominio y persistencia inicialmente.
- Desventajas: La entidad `PricingRule` se volvería muy compleja y tendría muchos campos opcionales. La lógica de aplicación y precedencia sería difícil de manejar dentro de un único servicio.

**Alternativa B – Usar JSON flexible para `Criteria` y `Value` en las reglas:**
- Ventajas: Alta flexibilidad para definir nuevos criterios o efectos sin cambios de esquema en la base de datos.
- Desventajas: Pérdida de type-safety en el dominio, dificultad para realizar consultas SQL basadas en criterios, mayor complejidad en la lógica de interpretación en la capa de aplicación, riesgo de errores en cálculos financieros.

**Alternativa C – Separar `BaseSalesPriceRule` y `SaleModificationRule` y usar Value Objects explícitos (Adoptada):**
- Ventajas: Claridad de dominio, fuerte type-safety, lógica de aplicación y precedencia bien definida para cada tipo de regla, persistencia optimizada para cada caso, alta mantenibilidad.
- Desventajas: Más entidades en el dominio y tablas en la base de datos, requiere más definición inicial.

---

## 3. Criterios de Decisión

-   **Pureza del Dominio:** Mantener las entidades de dominio limpias y enfocadas en la lógica de negocio.
-   **Type-Safety:** Garantizar la seguridad de tipos, especialmente para cálculos financieros.
-   **Rendimiento:** Asegurar que el cálculo de precios sea eficiente.
-   **Mantenibilidad:** Facilidad para entender, modificar y extender el módulo.
-   **Escalabilidad:** Capacidad para soportar un crecimiento en el número y complejidad de reglas.
-   **Simplicidad para MVP:** Enfocarse en la funcionalidad esencial para la primera versión.

---

## 4. Decisión Adoptada

Se ha decidido implementar el módulo de Pricing con las siguientes características arquitectónicas clave:

1.  **Separación de Reglas en Dos Entidades de Dominio Distintas:**
    *   **`BaseSalesPriceRule`:** Para el cálculo del precio base de venta.
    *   **`SaleModificationRule`:** Para modificaciones de precio en el momento de la venta.
2.  **Uso de Value Objects Explícitos:**
    *   `Money`: Representación de cantidades monetarias (fijado a EUR para MVP).
    *   `Percentage`: Representación de porcentajes.
    *   `RuleValue`: Encapsula el tipo y el valor del efecto de una regla.
3.  **Filosofía "No JSON" para Criterios y Valores:** Los criterios y valores de las reglas se definen mediante campos explícitos en las entidades o Value Objects, eliminando el uso de JSON genérico para estos fines.
4.  **Exclusión de Campos de Auditoría del Dominio:** `CreatedAt`, `UpdatedAt`, `CreatedBy`, `UpdatedBy` se gestionan en la capa de infraestructura/persistencia, no en las entidades de dominio.
5.  **Estrategia de Caché NoSQL con Redis:**
    *   Los precios base de venta calculados se almacenarán en una caché Redis (por `ProductID`, conteniendo los precios de todas sus variantes).
    *   Manejo de `ProductVariant`s JIT (Just-in-Time): Una vez persistidos, se tratan como cualquier otro `ProductVariant` nuevo, invalidando y repoblando la caché del `ProductID` asociado.
6.  **Esquema de Persistencia Adaptado:**
    *   Tablas `base_sales_price_rules` y `sale_modification_rules` que reflejan el dominio.
    *   Tabla `rule_value_types` como lookup para los tipos de valor.
    *   Campos `percentage_value` y `money_value_amount` separados para type-safety.
7.  **Casos de Uso Claramente Definidos:**
    *   `CalculateBaseSalesPriceForProductVariantUseCase`: Calcula y cachea precios base de venta.
    *   `CalculateFinalSalePriceUseCase`: Recupera de caché y aplica reglas de modificación de venta.
8.  **Simplificación para MVP de `SaleModificationRule`:** La `SaleModificationRule` se define con un conjunto de criterios más limitado para el MVP (`ClientIDs`, `ProductGroupID`, `MinOrderTotalAmount`, rango de fechas).

**Justificación:** Esta decisión optimiza la claridad del dominio y la robustez de los cálculos financieros a través de la fuerte tipificación y la separación de responsabilidades. La estrategia de caché aborda el rendimiento, y la simplificación para el MVP permite una entrega ágil sin comprometer la extensibilidad futura.

---

## 5. Consecuencias

### Positivas
-   **Alta Coherencia de Dominio:** Modelo de Pricing muy cercano al lenguaje ubicuo del negocio.
-   **Seguridad de Tipos:** Reducción drástica de errores en cálculos financieros críticos.
-   **Mantenibilidad Elevada:** Cada regla y su lógica son explícitas y bien encapsuladas.
-   **Rendimiento Optimizado:** La caché de Redis minimiza la recarga de cálculos complejos.
-   **Extensibilidad Clara:** La estructura de `RuleValue` y la separación de reglas facilitan la adición de nuevas funcionalidades.

### Negativas
-   **Mayor Número de Entidades/Tablas:** Comparado con una solución JSON-céntrica o una única entidad de regla.
-   **Complejidad Inicial de Configuración:** La gestión de `BaseSalesPriceRule` y `SaleModificationRule` requiere un entendimiento profundo de sus precedencias.

---

## 6. Alcance

Esta decisión aplica al diseño y la implementación del módulo de Pricing del sistema TramaTex, afectando directamente al backend (Go) y a la interacción con la base de datos (PostgreSQL) y la caché (Redis).

---

## 7. Integración con otros ADRs

-   ADR-002: Clean Architecture and DDD Adoption (Refuerza principios de pureza de dominio y separación de capas).
-   ADR-006: Domain-Driven Development Strategy (Alineación con la definición explícita de entidades y Value Objects).
-   ADR-pricing-domain-definition (OBSOLETE): Definicion del Dominio del Modulo de Precios (Sustituido por este ADR-016).

---

## 8. Notas Adicionales / Consideraciones Especiales

-   **Regla Genérica Única de `BaseSalesPriceRule`:** Se impondrá una restricción de negocio para asegurar la existencia de una única `BaseSalesPriceRule` genérica (todos los campos de targeting de producto a `null`) para garantizar un margen base.
-   **Manejo de JIT `ProductVariant`s:** La caché se invalidará y repoblará para todo el `ProductID` cuando se cree o modifique un `ProductVariant` (incluyendo los JIT) para asegurar la consistencia.
-   **EUR para MVP:** La moneda de trabajo será exclusivamente EUR en el MVP, aunque el `Money` Value Object en caché mantiene el campo `currency` para futura extensibilidad.

---

## 9. Referencias

-   Documentos internos: Contextos de `architecture.yaml` y `bounded-contexts.yaml`.
-   Discusiones y decisiones en este hilo de conversación.
