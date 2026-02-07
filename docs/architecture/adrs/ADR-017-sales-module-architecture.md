# ADR-017 – Arquitectura del Módulo de Sales

**Fecha:** viernes, 6 de febrero de 2026  
**Estado:** Aceptado  
**Autores:** Gemini CLI

---

## 1. Contexto

El módulo de Sales es crítico para el negocio TramaTex, ya que gestiona todo el ciclo de vida de la venta, desde la cotización inicial hasta la facturación. La complejidad surge de la necesidad de manejar diferentes tipos de documentos (cotizaciones, pedidos, albaranes, facturas), sus interrelaciones, flujos de estado, y la integración con otros módulos fundamentales como `Party` (clientes), `Product` (productos y variantes) y `Pricing` (motor de precios).

**Problemas a Resolver:**
- Necesidad de unificar la gestión de múltiples documentos de venta con ciclos de vida interconectados.
- Requisito de "manual override" en precios y descuentos, que deben coexistir con los valores calculados por el módulo de Pricing.
- Asegurar la trazabilidad completa del proceso de venta.
- Mantener la cohesión y baja acoplamiento con otros Bounded Contexts.
- Definir claramente las responsabilidades de cálculo (precios, impuestos) entre módulos.

---

## 2. Alternativas Consideradas

**Alternativa A – Entidad `SalesDocument` Genérica con un `Type`:**
- **Descripción:** Una única entidad `SalesDocument` con un campo `Type` (ej. `QUOTE`, `ORDER`, `DELIVERY_NOTE`, `INVOICE`). La lógica de negocio y los atributos específicos se manejarían condicionalmente según el `Type`.
- **Ventajas:** Menos entidades y potencialmente menos tablas en la base de datos.
- **Desventajas:** La entidad `SalesDocument` se volvería muy compleja y tendría muchos campos opcionales o relevantes solo para ciertos tipos. La lógica de negocio y las transiciones de estado serían difíciles de manejar, creando un "anémico" modelo de dominio que no encapsula bien el comportamiento. Iría en contra de la pureza del dominio.

**Alternativa B – Entidades Separadas para Cada Tipo de Documento de Venta (Adoptada):**
- **Descripción:** Se definen entidades de dominio distintas (`Quote`, `SalesOrder`, `DeliveryNote`, `Invoice`) para cada tipo de documento de venta, con sus propias entidades de línea (`QuoteLineItem`, `OrderLineItem`, etc.). Estas entidades están relacionadas entre sí de forma explícita.
- **Ventajas:** Claridad de dominio, cada entidad encapsula su propio ciclo de vida y atributos. Facilita la comprensión y el mantenimiento. Se alinea con los principios de DDD y Clean Architecture. Mejora la escalabilidad y la capacidad de refactorización futura.
- **Desventajas:** Más entidades y potencialmente más tablas en la base de datos. Requiere más definición inicial.

---

## 3. Criterios de Decisión

-   **Claridad del Dominio:** El modelo debe ser intuitivo y reflejar directamente los conceptos de negocio.
-   **Cohesión y Acoplamiento:** Maximizar la cohesión dentro del módulo y minimizar el acoplamiento con otros módulos.
-   **Mantenibilidad y Extensibilidad:** Facilidad para entender, modificar y extender el módulo.
-   **Trazabilidad:** Capacidad para seguir el flujo completo de una venta.
-   **Soporte de Flujos de Negocio:** El modelo debe soportar eficientemente las transiciones y operaciones entre documentos (ej. Quote a Order).
-   **Integración con Módulos Externos:** Capacidad para consumir servicios de `Party`, `Product` y `Pricing`.

---

## 4. Decisión Adoptada

Se adopta la **Alternativa B: Entidades Separadas para Cada Tipo de Documento de Venta**.

### Entidades de Dominio Clave:

1.  **`Quote` y `QuoteLineItem`:** Para presupuestos.
2.  **`SalesOrder` y `OrderLineItem`:** Para pedidos de venta.
3.  **`DeliveryNote` y `DeliveryNoteLineItem`:** Para albaranes.
4.  **`Invoice` y `InvoiceLineItem`:** Para facturas.

### Value Objects Clave:

*   `Money`, `Percentage` (del módulo Pricing).
*   `PartyID` (del módulo Party).
*   `ProductVariantID` (del módulo Product).
*   `OrderNumber`, `QuoteNumber`, `DeliveryNoteNumber`, `InvoiceNumber` (para identificadores únicos de documentos).

### Aspectos Clave del Diseño:

*   **"Manual Override" en Ítems de Línea:** Las entidades de línea (`QuoteLineItem`, `OrderLineItem`, `InvoiceLineItem`) incluyen campos `CalculatedUnitPrice`/`CalculatedDiscountPerUnit` y `ManualUnitPrice`/`ManualDiscountPerUnit`. Un campo `FinalUnitPrice`/`FinalDiscountPerUnit` determinará el valor real, permitiendo ajustes manuales desde la UI que prevalecen sobre los cálculos automáticos de Pricing.
*   **Responsabilidad del Cálculo de Impuestos:** El cálculo de `TaxAmount` y `Total` es responsabilidad del módulo `Sales` (lógica interna del módulo, probablemente basada en tasas configurables o del cliente).
*   **Flujo de Estado Detallado:** Cada entidad documental tiene un ciclo de vida (`Status`) y transiciones bien definidas.
*   **Integración con Módulos Externos:**
    *   **`Pricing`:** `Sales` consume el `CalculateFinalSalePriceUseCase` para obtener los precios unitarios finales (con descuentos aplicados) para los ítems de línea, que luego pueden ser sobreescritos.
    *   **`Party`:** Se usa `PartyID` para identificar clientes en todos los documentos.
    *   **`Product`:** Se usa `ProductVariantID` para identificar los productos en los ítems de línea.

**Justificación:**
Esta decisión conduce a un modelo de dominio más claro y expresivo, donde cada documento de venta tiene una identidad y un propósito únicos. Facilita la implementación de flujos de negocio complejos (ej. `ConvertQuoteToOrder`, `CreateDeliveryNote`, `CreateInvoiceFromOrder`) y la gestión de sus estados. La integración con `Pricing` es directa, consumiendo sus resultados y permitiendo el override manual según el requisito de negocio.

---

## 5. Consecuencias

### Positivas
-   **Claridad de Dominio:** Modelo fácil de entender para los stakeholders y desarrolladores.
-   **Cohesión Elevada:** Cada entidad de documento de venta es autocontenida y gestiona su propia lógica.
-   **Flexibilidad:** Facilita la adaptación a cambios en los flujos de documentos o requisitos específicos de cada tipo.
-   **Trazabilidad:** El rastro de los documentos (Quote a Order, Order a DeliveryNote/Invoice) es explícito.
-   **Soporte de Override Manual:** El modelo permite explícitamente el ajuste manual de precios.

### Negativas
-   **Mayor Número de Entidades y Tablas:** Más complejidad de gestión en la base de datos y en el código.
-   **Potencial Duplicación de Atributos:** Ciertos atributos (ej. `PartyID`, `Notes`) se repiten en varias entidades, aunque encapsulados dentro de su contexto.

---

## 6. Alcance

Esta decisión aplica al diseño y la implementación del módulo de Sales del sistema TramaTex, afectando directamente al backend (Go) y a su interacción con la base de datos (PostgreSQL), así como a los módulos de `Party`, `Product` y `Pricing`.

---

## 7. Integración con otros ADRs

-   ADR-002: Clean Architecture and DDD Adoption (Refuerza la pureza del dominio y la separación de capas).
-   ADR-006: Domain-Driven Development Strategy (Alineación con la definición explícita de entidades y Value Objects).
-   ADR-007: Module Implementation Order (Sales se implementa en la Fase 2, después de Party, Product y Pricing).
-   ADR-016: Arquitectura del Módulo de Pricing (Sales es un consumidor clave del servicio de Pricing).

---

## 8. Notas Adicionales / Consideraciones Especiales

*   **Generación de Números de Documento:** Los Value Objects `OrderNumber`, `QuoteNumber`, `DeliveryNoteNumber`, `InvoiceNumber` encapsularán la lógica de generación de identificadores únicos para cada tipo de documento.
*   **Transiciones de Estado:** El `Status` de cada documento (`Quote`, `SalesOrder`, `Invoice`) debe ser gestionado por una máquina de estados o una lógica de dominio robusta que valide las transiciones permitidas.
*   **Cálculo de Impuestos:** Para el MVP, el cálculo de impuestos puede ser simplificado (ej. tasa fija sobre el `Subtotal`). En futuras fases, esto podría volverse más complejo (reglas fiscales, exenciones) y requerir un servicio de dominio o incluso un módulo de `Taxation` propio.

*   **Tipo de PartyID en persistencia:** Sales usará UUID en sus tablas. Mientras Party mantenga `id` como `VARCHAR(36)`, la validación de existencia se hará en la capa de aplicación sin FK directa. Cuando Party migre a UUID, se activarán FKs explícitas.

*   **Orden de implementación:** Sales depende de Party, Product y Pricing. La implementación puede iniciar en paralelo, pero la integración completa (pricing real, validaciones cross-modulo y FKs) queda bloqueada hasta que esos módulos estén completados y testeados.

---

## 9. Referencias

*   Contextos de `architecture.yaml` y `bounded-contexts.yaml`.
*   Documentación del módulo de Sales (`module-spec.md`, `domain-model.md`, `use-cases.md`, `api-contracts.md`).
*   Discusiones y decisiones en este hilo de conversación.
