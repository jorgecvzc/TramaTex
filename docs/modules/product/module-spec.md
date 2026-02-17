# Módulo de Product (Catálogo de Productos)

**Estado:** En Implementación (MVP)  
**Última actualización:** 13 de febrero de 2026

## Estado de Implementación

### Componentes Completos ✅
- **Atributos (Attributes):**
  - Backend: CREATE, READ, UPDATE funcionales
  - DTOs con estructura completa (`AttributeValueDTO` incluye id, value, code)
  - Frontend: UI completa con CRUD funcional
  - Encoding UTF-8 verificado y funcionando
  
- **Marcas (Brands):**
  - Backend: CREATE, UPDATE funcionales
  - Frontend: UI básica implementada

- **Grupos de Productos (Product Groups):**
  - Backend: CREATE, UPDATE funcionales
  - Frontend: UI básica implementada

### Componentes Pendientes ⏳
- **Productos (Products):** UI pendiente (lógica backend lista)
- **Variantes (Product Variants):** Pendiente
- **Configuraciones de Servicio:** Pendiente

### Simplificaciones MVP
- Sistema de scope (brand/group) de atributos removido temporalmente
- Atributos son globales (asignación manual por usuario)
- Ver ADR-015 para detalles completos

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
*   **RF-P-002:** Crear y mantener productos (tangibles o servicios) con sus atributos base, marca y grupos.
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
*   **Herencia de Atributos con Anulación:** Los atributos se heredan de marcas y grupos, permitiendo la anulación en niveles más específicos (directo, grupo+marca, grupo, marca, genérico).
*   **Creación JIT de `ProductVariant`s:** Las variantes se crean en la base de datos bajo demanda para evitar la pre-generación masiva, comenzando con un estado `PROVISIONAL`.
*   **Composición de SKU Determinista:** Los SKUs de las variantes se construyen algorítmicamente a partir de los códigos de atributos y valores.
*   **`ProductType` (`TANGIBLE` vs `SERVICE`):** Permite diferenciar entre bienes físicos y servicios, con un manejo especial para `PartyServiceConfiguration`s en el caso de servicios (ver `ADR-013`).
*   **Relaciones con Otros Módulos:**
    *   **Pricing:** Consume `ProductVariantID`s y sus atributos para el cálculo de precios.
    *   **Sales:** Las órdenes de venta contienen `ProductVariant`s.
    *   **Party:** Referenciado por `PartyServiceConfiguration`.
    *   **MES:** Utiliza `ProductVariant`s para la planificación de la producción.

---
