# ADR-015 – Arquitectura del Módulo de Product

**Fecha:** viernes, 6 de febrero de 2026  
**Estado:** Implementado  
**Autores:** Gemini CLI  
**Última actualización:** 2026-02-22  

---

## 1. Contexto

El módulo `Product` es fundacional para TramaTex, sirviendo como la fuente de verdad para la definición de todos los bienes y servicios vendibles. La complejidad radica en la necesidad de gestionar productos con un alto grado de configuración (variantes, atributos, opciones) y un enfoque "Just-in-Time" (JIT) para la creación de variantes.

---

## 2. Alternativas Consideradas

**Alternativa A – Modelo simple de Producto/Variante**
- Ventajas: Muy simple, fácil de implementar inicialmente.
- Desventajas: No soporta la flexibilidad de atributos arbitrarios, la herencia, ni la generación JIT de variantes.

**Alternativa B – Uso de JSON para atributos/variantes**
- Ventajas: Máxima flexibilidad para añadir nuevos atributos sin cambios de esquema.
- Desventajas: Pérdida de type-safety, lógica compleja en el dominio, dificultad para consultas y reportes.

**Alternativa C – Modelo de Atributos/Valores explícito con Herencia y JIT (Adoptada)**
- Ventajas: Modelo robusto, type-safe, permite composición de SKU determinista, ideal para Pricing y MES.
- Desventajas: Mayor complejidad inicial en el dominio y base de datos.

---

## 3. Criterios de Decisión

- **Flexibilidad del Catálogo:** Gestión de un número arbitrario de atributos y valores.
- **Soporte JIT:** Capacidad para generar variantes bajo demanda ("Find or Create").
- **Herencia y Anulación:** Resolución de atributos por especificidad (Directo > Grupo > Marca > Genérico).
- **Type-Safety:** Evitar el uso de JSON para lógica de negocio crítica.

---

## 4. Decisión Adoptada

Se adopta el **Modelo de Atributos/Valores explícito con Herencia y JIT**.

### Entidades Clave:
- `Attribute` / `AttributeValue`: Definición de características configurables.
- `Product`: Plantilla base con SKU, marca y grupos.
- `ProductVariant`: Instancia vendible con SKU compuesto determinista.

---

## 5. Consecuencias

### Positivas
- Alta flexibilidad para productos altamente configurables (textil/EPIs).
- Rendimiento optimizado mediante creación JIT de variantes.
- SKUs deterministas que facilitan la identificación única.

### Negativas
- Motor de herencia de atributos complejo de implementar.
- Curva de aprendizaje elevada para el manejo de la lógica JIT.

---

## 6. Alcance

Aplica al diseño y la implementación del módulo de Product, afectando al backend (Go), base de datos y la integración con `Party`, `Pricing` y `Sales`.

---

## 7. Integración con otros ADRs

- [ADR-002: Clean Architecture and DDD Adoption](adr-002-clean-architecture-ddd-adoption.md)
- [ADR-013: Manejo de Modificaciones de Producto](adr-013-product-modifications-handling.md)
- [ADR-016: Arquitectura del Módulo de Pricing](adr-016-pricing-module-architecture.md)
- [ADR-017: Arquitectura del Módulo de Sales](adr-017-sales-module-architecture.md)

---

## 8. Notas Adicionales / Consideraciones Especiales

### Actualización: Sistema de Impuestos Integrado (2026-02-22)
- Campo `tax_rate` añadido a `Product` (soporte 21%, 10%, 4%, 0%).
- Integración de `Price Modifiers` en `AttributeValue` (FIXED y PERCENTAGE).
- Aplicación automática de `default_markup_percentage` de la marca.

### Correcciones de Implementación (2026-02-13)
- Mejora de `AttributeDTO` para incluir IDs de valores, permitiendo edición correcta en UI.
- Simplificación de Data Models eliminando campos de auditoría no presentes en el esquema físico.

### Simplificación para MVP (2026-02-12)
- El sistema de scope (brand/group) de atributos se ha simplificado: los atributos son globales en el MVP y su asignación correcta es responsabilidad del usuario.

---

## 9. Referencias

*   Documentación del módulo de Product (`module-spec.md`, `domain-model.md`, `use-cases.md`).
*   Contextos de `architecture.yaml` y `bounded-contexts.yaml`.
