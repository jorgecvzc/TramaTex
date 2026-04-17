# ADR-017 – Arquitectura del Módulo de Sales

**Fecha:** viernes, 6 de febrero de 2026  
**Estado:** Aceptado  
**Autores:** Gemini CLI

---

## 1. Contexto

El módulo de Sales es crítico para el negocio TramaTex, ya que gestiona todo el ciclo de vida de la venta, desde la cotización inicial hasta la facturación. La complejidad surge de la necesidad de manejar diferentes tipos de documentos (cotizaciones, pedidos, albaranes, facturas), sus interrelaciones, flujos de estado, y la integración con otros módulos fundamentales como `Party` (clientes), `Product` (productos y variantes) y `Pricing` (motor de precios).

---

## 2. Alternativas Consideradas

**Alternativa A – Entidad `SalesDocument` Genérica con un `Type`:**
- Ventajas: Menos entidades y potencialmente menos tablas en la base de datos.
- Desventajas: La entidad se volvería muy compleja con muchos campos opcionales. La lógica de negocio y las transiciones de estado serían difíciles de manejar.

**Alternativa B – Entidades Separadas para Cada Tipo de Documento de Venta (Adoptada):**
- Ventajas: Claridad de dominio, cada entidad encapsula su propio ciclo de vida y atributos. Facilita la comprensión y el mantenimiento. Se alinea con Clean Architecture.
- Desventajas: Más entidades y tablas en la base de datos. Requiere más definición inicial.

---

## 3. Criterios de Decisión

- **Claridad del Dominio:** El modelo debe reflejar directamente los conceptos de negocio.
- **Cohesión y Acoplamiento:** Maximizar cohesión interna y minimizar acoplamiento externo.
- **Trazabilidad:** Capacidad para seguir el flujo completo de una venta (Quote -> Order -> Delivery -> Invoice).
- **Flexibilidad:** Soporte para "manual override" de precios y descuentos.

---

## 4. Decisión Adoptada

Se adopta la **Alternativa B: Entidades Separadas para Cada Tipo de Documento de Venta** (`Quote`, `SalesOrder`, `DeliveryNote`, `Invoice`).

### Aspectos Clave del Diseño:
*   **"Manual Override" en Ítems de Línea:** Campos para valores calculados y manuales, donde el manual prevalece si existe.
*   **Responsabilidad del Cálculo de Impuestos:** El módulo `Sales` es responsable de calcular el `TaxAmount` y `Total`.
*   **Integración con Módulos Externos:** Consume precios de `Pricing`, clientes de `Party` y variantes de `Product`.

---

## 5. Consecuencias

### Positivas
- Modelo fácil de entender y mantener.
- Cada entidad es autocontenida con su propia lógica de estado.
- Trazabilidad explícita del rastro documental.

### Negativas
- Mayor número de entidades y tablas.
- Potencial duplicación de atributos comunes (ej. PartyID) en diferentes contextos documentales.

---

## 6. Alcance

Aplica al diseño y la implementación del módulo de Sales, afectando al backend (Go), base de datos y la interacción con los módulos de `Party`, `Product` y `Pricing`.

---

## 7. Integración con otros ADRs

-   [ADR-002: Clean Architecture and DDD Adoption](adr-002-clean-architecture-ddd-adoption.md) (Pureza del dominio y separación de capas).
-   [ADR-006: Domain-Driven Development Strategy](adr-006-domain-driven-development-strategy.md) (Alineación con entidades y Value Objects).
-   [ADR-007: Module Implementation Order](adr-007-module-implementation-order.md) (Sales se implementa en la Fase 2).
-   [ADR-016: Arquitectura del Módulo de Pricing](adr-016-pricing-module-architecture.md) (Sales consume el motor de precios).

---

## 8. Notas Adicionales / Consideraciones Especiales

### ACTUALIZACIÓN: Sistema de Impuestos Implementado (2026-02-22)
Se ha evolucionado de un cálculo simplificado a un sistema completo:
1. **TaxRate por Producto:** Definido según normativa fiscal (21%, 10%, 4%, 0%).
2. **Tax Fields en Line Items:** Inclusión de `TaxRate` y `TaxAmount` en todas las entidades de línea.
3. **Cálculos Automáticos:** Métodos de dominio para recalcular totales automáticamente.

### Otras Consideraciones
*   **Generación de Números:** Value Objects específicos para `OrderNumber`, `QuoteNumber`, etc.
*   **Transiciones de Estado:** Gestionadas por lógica de dominio robusta que valida flujos permitidos.
*   **Simplificación de Flujo (Abril 2026):** Para agilizar la operativa, los Pedidos de Venta nacen directamente en estado `EN_PREPARACION` (Confirmado), eliminando el paso manual de confirmación desde `PENDIENTE` para usuarios autorizados. El estado `PENDIENTE` se reserva como estado de seguridad tras una reactivación.

---

## 9. Referencias

*   Contextos de `architecture.yaml` y `bounded-contexts.yaml`.
*   Documentación del módulo de Sales (`module-spec.md`, `domain-model.md`, `use-cases.md`, `api-contracts.md`).
