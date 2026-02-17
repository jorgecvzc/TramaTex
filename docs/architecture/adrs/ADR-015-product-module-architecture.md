# ADR-015 – Arquitectura del Módulo de Product

**Fecha:** viernes, 6 de febrero de 2026  
**Última actualización:** 13 de febrero de 2026 (Correcciones de Implementación)  
**Estado:** Implementado  
**Autores:** Gemini CLI

---

## ACTUALIZACIÓN: Correcciones de Implementación (2026-02-13)

**Correcciones realizadas en la implementación del MVP:**

### Backend - DTOs y Data Models
1. **AttributeDTO mejorado:**
   - Creado `AttributeValueDTO` con estructura completa: `{ id: UUID, value: string, code: string }`
   - Modificado `AttributeDTO.Values` de `[]string` a `[]AttributeValueDTO`
   - Cambiado campo `AttributeName` → `Code` para consistencia con el dominio
   - **Impacto:** Permite edición correcta de atributos en UI al incluir IDs de valores

2. **Data Models simplificados:**
   - Eliminados campos `CreatedBy` y `ModifiedBy` de `AttributeDataModel` y `AttributeValueDataModel`
   - **Razón:** Campos no existen en el esquema de base de datos actual
   - **Solución:** Auditoría se maneja con `CreatedAt`, `UpdatedAt`, `DeletedAt` (gorm.Model)

### Frontend - API Integration
3. **Servicios API corregidos:**
   - `createAttribute()`: Payload corregido (removido `id`, usar `sortOrder` en camelCase)
   - `updateAttribute()`: Formateado correcto con estructura completa de valores
   - `AttributeForm.vue`: Preserva IDs de valores existentes al editar

### Testing Completado
4. **CRUD de Atributos verificado:**
   - ✅ GET /api/attributes → Lista completa con valores estructurados
   - ✅ POST /api/attributes → Creación exitosa con UTF-8
   - ✅ PUT /api/attributes/:id → Actualización con preservación de IDs
   - ✅ UI funcionando: crear, editar, listar con encoding correcto

**Estado actual:** Sistema de Atributos completamente funcional en MVP.

---

## NOTA IMPORTANTE: Simplificación para MVP (2026-02-12)

**Decisión:** El sistema de scope (brand/group) de atributos ha sido simplificado para el MVP.

**Cambios:**
- Los campos `scope_brand_id` y `scope_group_id` se han eliminado de la entidad `Attribute`
- Los atributos ahora son globales; la asignación correcta es responsabilidad del usuario
- Ejemplo: crear "Talla FYR" y asignarlo manualmente a productos FYR
- Sistema de scope completo diferido para post-MVP

**Impacto:**
- Reduce complejidad de implementación inicial
- No afecta la filosofía de Pricing
- Simplifica UIs y lógica de dominio
- **La documentación técnica a continuación refleja el diseño original completo, pero debe considerarse la simplificación para el MVP actual**

---

## 1. Contexto

El módulo `Product` es fundacional para TramaTex, sirviendo como la fuente de verdad para la definición de todos los bienes y servicios vendibles. La complejidad radica en la necesidad de gestionar productos con un alto grado de configuración (variantes, atributos, opciones) y un enfoque "Just-in-Time" (JIT) para la creación de variantes.

**Problemas a Resolver:**
- Gestionar un catálogo de productos flexible que soporte atributos y valores configurables.
- Permitir la creación dinámica (JIT) de `ProductVariant`s para evitar la pre-generación masiva.
- Soportar la herencia y anulación de atributos a lo largo de una jerarquía (genérico, marca, grupo, producto).
- Proporcionar una composición de SKU determinista y jerárquica para cada `ProductVariant`.
- Integrar la gestión de configuraciones de servicio específicas por `Party` (PartyServiceConfiguration).

---

## 2. Alternativas Consideradas

**Alternativa A – Modelo simple de Producto/Variante (como en `module-spec.md` inicial):**
- **Descripción:** `Product` con atributos básicos y `ProductVariant`s con campos fijos como talla y color.
- **Ventajas:** Muy simple, fácil de implementar inicialmente.
- **Desventajas:** No soporta la flexibilidad de atributos arbitrarios, la herencia, ni la generación JIT de variantes. Es insuficiente para los requisitos reales del dominio textil/EPIs.

**Alternativa B – Uso de JSON para atributos/variantes:**
- **Descripción:** `Product` y `ProductVariant` tendrían un campo `attributes_json` para almacenar configuraciones flexibles.
- **Ventajas:** Máxima flexibilidad para añadir nuevos atributos sin cambios de esquema.
- **Desventajas:** Pérdida de type-safety, lógica compleja en el dominio para interpretar el JSON, dificultad para consultas y reportes, no soporta la lógica de herencia/anulación de atributos ni la composición determinista de SKU.

**Alternativa C – Modelo de Atributos/Valores explícito con Herencia y JIT (Adoptada):**
- **Descripción:** Se introducen entidades explícitas `Attribute` y `AttributeValue` (que se refieren a `ProductOptionSet` en los contratos API). `Product` y `ProductVariant` se relacionan con estos. Se define una lógica de herencia de atributos con anulación y un mecanismo "Find or Create" (JIT) para `ProductVariant`.
- **Ventajas:** Modelo de dominio robusto y flexible, type-safe, soporta todas las complejidades de atributos/variantes, permite composición de SKU determinista, ideal para gestión de precios y MES.
- **Desventajas:** Mayor complejidad inicial en el dominio y base de datos, requiere una definición cuidadosa de la lógica de herencia y JIT.

---

## 3. Criterios de Decisión

-   **Flexibilidad del Catálogo:** Capacidad para definir y gestionar un número arbitrario de atributos y valores.
-   **Claridad del Dominio:** El modelo debe ser intuitivo y reflejar directamente los conceptos de negocio de productos configurables.
-   **Cohesión y Acoplamiento:** Mantener la cohesión dentro del módulo y minimizar el acoplamiento con otros módulos.
-   **Mantenibilidad y Extensibilidad:** Facilidad para entender, modificar y extender el módulo.
-   **Soporte JIT:** Capacidad para generar variantes bajo demanda.
-   **Integración:** Cómo el módulo interactúa con `Pricing`, `Party` y `Sales`.

---

## 4. Decisión Adoptada

Se adopta la **Alternativa C: Modelo de Atributos/Valores explícito con Herencia y JIT**. Este modelo es el que se ha implementado y documentado en `domain-model.md` y `use-cases.md`.

### Entidades de Dominio Clave:

1.  **`Attribute` (`ProductOptionSet` en API):** Gestiona una característica configurable (ej. "Talla", "Color") y sus `AttributeValue`s. Define su alcance (genérico, por marca, por grupo, etc.).
2.  **`Product`:** Plantilla del producto/servicio. Contiene `SKU` base, `Name`, `LongName`, `ProductType`, `BrandID`, `GroupIDs`. Hereda y puede sobrescribir `Attribute`s.
3.  **`ProductVariant`:** Instancia final vendible de un `Product`, definida por una combinación única de `AttributeValue`s. Tiene un `SKU` compuesto determinista y un `Status` (`PROVISIONAL`, `CONFIRMED`).
4.  **`Brand`:** Agrupa productos bajo una marca.
5.  **`ProductGroup`:** Categoría jerárquica para productos.
6.  **`PartyServiceConfiguration`:** (Decisión de ADR-013) Entidad para guardar configuraciones de servicios específicas por cliente (`PartyID`, `ServiceID`, `ConfigurationDetails`).

### Aspectos Clave del Diseño:

*   **Composición de SKU Jerárquico:** El `SKU` de `ProductVariant` se construye de forma determinista (ej. `{Product.SKU}-{Attr1.Code}.{Val1.Code}-{Attr2.Code}.{Val2.Code}...`).
*   **Herencia de Atributos con Anulación:** Los `Attribute`s aplicables a un `Product` se resuelven mediante una fusión con anulación por especificidad (Directo > Grupo+Marca > Grupo > Marca > Genérico).
*   **Creación JIT de `ProductVariant`s:** El mecanismo "Find or Create" permite generar `ProductVariant`s bajo demanda (`Status = PROVISIONAL`) si no existen previamente.
*   **Exclusión de Campos de Auditoría del Dominio:** `CreatedAt`, `UpdatedAt`, `CreatedBy`, `UpdatedBy` se gestionan en la capa de infraestructura/persistencia.

**Justificación:**
Esta decisión proporciona un modelo de dominio rico y flexible que maneja la complejidad del catálogo de productos configurable, esencial para el negocio textil y de EPIs. Soporta la generación dinámica de variantes y la gestión precisa de atributos, lo cual es fundamental para Pricing y MES. La "Filosofía No JSON" aplicada a los atributos y valores garantiza type-safety y claridad.

---

## 5. Consecuencias

### Positivas
-   **Alta Flexibilidad del Catálogo:** Soporte robusto para productos altamente configurables.
-   **Coherencia de Dominio:** Modelo que refleja con precisión los conceptos de negocio.
-   **Rendimiento Optimizado:** La creación JIT evita la sobrecarga de generar y mantener variantes innecesarias.
-   **Integración Directa:** Facilita la interacción con los módulos de `Pricing` y `Sales` al proporcionar `ProductVariantID`s y atributos claros.
-   **Claridad de SKU:** SKUs deterministas facilitan la identificación de variantes.

### Negativas
-   **Mayor Complejidad de Implementación:** El motor de herencia de atributos y la lógica JIT son complejos de implementar y testear.
-   **Curva de Aprendizaje:** Requiere que los desarrolladores comprendan la lógica de herencia de atributos y el comportamiento JIT.

---

## 6. Alcance

Esta decisión aplica al diseño y la implementación del módulo de Product del sistema TramaTex, afectando directamente al backend (Go) y a su interacción con la base de datos (PostgreSQL), así como a los módulos de `Party`, `Pricing` y `Sales`.

---

## 7. Integración con otros ADRs

-   ADR-002: Clean Architecture and DDD Adoption (Refuerza la pureza del dominio y la separación de capas).
-   ADR-006: Domain-Driven Development Strategy (Alineación con la definición explícita de entidades y Value Objects).
-   ADR-007: Module Implementation Order (Product se implementa en la Fase 1).
-   ADR-013: Manejo de Modificaciones de Producto (Arreglos/Marcajes) (Este ADR es la base de las decisiones de ProductType='SERVICE' y PartyServiceConfiguration).
-   [ADR-016: Arquitectura del Módulo de Pricing](ADR-016-pricing-module-architecture.md) (Product es una dependencia clave para Pricing).
-   [ADR-017: Arquitectura del Módulo de Sales](ADR-017-sales-module-architecture.md) (Product es una dependencia clave para Sales).

---

## 8. Notas Adicionales / Consideraciones Especiales

*   **`ProductOptionSet` vs `Attribute`:** Aunque el dominio usa `Attribute`, los contratos API (`api-contracts.md`) usan el término `ProductOptionSet`. Es crucial mantener una correspondencia clara en la implementación y la documentación para evitar confusiones.
*   **`ProductVariantDto.price`:** En los DTOs, se observó un campo `price` en `ProductVariantDto`. Se debe asegurar que `ProductVariant` **no almacene precios calculados**, ya que esa es responsabilidad del módulo de Pricing. Podría almacenar un costo base que Pricing utilizaría.
*   **`ProductDto.description`:** El `ProductDto` debe reflejar `Name` (nombre corto) y `LongName` (nombre completo) del dominio.

---

## 9. Referencias

*   Contextos de `architecture.yaml` y `bounded-contexts.yaml`.
*   Documentación del módulo de Product (`module-spec.md`, `domain-model.md`, `use-cases.md`, `api-contracts.md`).
*   Discusiones y decisiones en este hilo de conversación.
