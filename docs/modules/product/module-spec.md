# Módulo de Product (Catálogo de Productos)

**Estado:** ✅ **COMPLETO (100%)**  
**Última actualización:** 30 de marzo de 2026

## Estado de Implementación

### Componentes Completos ✅
- **Atributos (Attributes):**
  - Backend: CRUD completo funcional.
  - DTOs con estructura completa (`AttributeValueDTO`).
  - Frontend: UI completa con gestión de valores y modificadores de precio.
  
- **Marcas (Brands) y Grupos (Product Groups):**
  - Backend: CRUD completo.
  - Frontend: UI de gestión integrada en Datos Maestros.
  - Clasificación de grupos en `TANGIBLE` vs `SERVICE`.

- **Productos (Products):**
  - Backend: Gestión completa de productos base.
  - Frontend: `ProductList.vue` y `ProductDetail.vue` con wizard de creación.
  - Integración de herencia de atributos.

- **Variantes (Product Variants):**
  - Backend: Generador de variantes bulk y creación Just-in-Time (JIT).
  - Frontend: Tabla dinámica de variantes con edición de SKUs y costes.
  - Selector de variantes interactivo para el módulo de Sales.

- **Configuraciones de Servicio:**
  - Implementación de `PartyServiceConfiguration` para vincular servicios a terceros específicos.

---

## 1. Propósito

*   **Visión del Módulo:** Gestionar el catálogo completo de productos y servicios de TramaTex, incluyendo la definición de atributos configurables, la creación dinámica de variantes y la clasificación jerárquica.
*   **Objetivos Clave:**
    *   Proveer una fuente de verdad para toda la información de productos y variantes.
    *   Soportar productos configurables mediante un sistema flexible de atributos y valores.
    *   Permitir la creación "Just-in-Time" (JIT) de `ProductVariant`s para optimizar la gestión de datos.
    *   Facilitar la integración con los módulos de `Pricing` (para cálculo de precios) y `MES` (para gestión de producción).

## 2. Requisitos

### 2.1. Requisitos Funcionales

*   **RF-P-001:** Crear y mantener atributos y sus valores.
*   **RF-P-002:** Crear y mantener productos (tangibles o servicios) con sus atributos base, marca opcional y grupos.
*   **RF-P-003:** Gestionar la asignación y herencia de atributos a productos.
*   **RF-P-004:** Generar SKUs deterministas para `ProductVariant`s basados en su configuración de atributos.
*   **RF-P-005:** Crear `ProductVariant`s de forma explícita o "Just-in-Time".
*   **RF-P-006:** Mantener el estado de las `ProductVariant`s (`PROVISIONAL`, `CONFIRMED`).
*   **RF-P-007:** Gestionar marcas y grupos de productos (categorización jerárquica).
*   **RF-P-008:** Gestionar `PartyServiceConfiguration`s asociadas a productos de tipo `SERVICE`.

## 3. Casos de Uso

Para una descripción detallada de los casos de uso, incluyendo flujos y entradas/salidas, consulte el documento [Casos de Uso - Módulo Product](./use-cases.md).

## 4. Modelo de Dominio

Para una descripción detallada del modelo de dominio, incluyendo entidades, Value Objects, agregados y sus relaciones, consulte el documento [Modelo de Dominio - Módulo Product](./domain-model.md).

## 5. Decisiones de Diseño

*   **Modelo de Atributos/Valores Explícito:** Se utiliza un modelo flexible de `Attribute` y `AttributeValue` para gestionar las características configurables de los productos.
*   **Modificadores de Precio en Atributos:** Los `AttributeValue` pueden incluir modificadores de precio (FIXED o PERCENTAGE) que ajustan dinámicamente el precio base de las variantes. Los modificadores pueden ser positivos (incrementan precio) o negativos (reducen precio).
*   **Cálculo Dinámico de Precio Base de Variante:** El `baseCost` de una variante se calcula algorítmicamente: `baseCost = producto.basePrice + sum(modificadores de atributos)`. Este valor NO se almacena, se calcula en tiempo real para mantener coherencia con los cambios en los modificadores.
*   **Herencia de Atributos con Anulación:** Los atributos se heredan de marcas y grupos, permitiendo la anulación en niveles más específicos (directo, grupo+marca, grupo, marca, genérico). Si el producto no tiene marca, solo participan los niveles directos, de grupo y genéricos.
*   **Creación JIT de `ProductVariant`s:** Las variantes se crean en la base de datos bajo demanda para evitar la pre-generación masiva, comenzando con un estado `PROVISIONAL`.
*   **Composición de SKU Determinista:** Los SKUs de las variantes se construyen algorítmicamente a partir de los códigos de atributos y valores.
*   **`ProductType` (`TANGIBLE` vs `SERVICE`):** Permite diferenciar entre bienes físicos y servicios, con un manejo especial para `PartyServiceConfiguration`s en el caso de servicios (ver `ADR-013`).
*   **Relaciones con Otros Módulos:**
    *   **Pricing:** Consume el `baseCost` calculado de las variantes para aplicar márgenes de beneficio, descuentos y reglas de pricing específicas.
    *   **Sales:** Las órdenes de venta y presupuestos referencian `ProductVariant`s. **⚠️ IMPORTANTE:** Sales debe obtener el `baseCost` actualizado de cada variante al momento de crear líneas de pedido, ya que este valor se calcula dinámicamente y puede cambiar si se modifican el precio base del producto o los modificadores de atributos. Ver [Sección 5 de API Contracts](./api-contracts.md#5-cálculo-de-precio-de-variante) para detalles del algoritmo.
    *   **Party:** Referenciado por `PartyServiceConfiguration`.
    *   **MES:** Utiliza `ProductVariant`s para la planificación de la producción.

---
