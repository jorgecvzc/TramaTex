# Modelo de Dominio - Módulo Product

Este documento describe la lógica de estructuración del catálogo de TramaTex, centrándose en el sistema de variantes dinámicas y la gestión de productos configurables.

---

## 1. Agregados y Entidades de Dominio

### El Producto como Plantilla (`Product`)
El `Product` actúa como el concepto base o "molde". Define la naturaleza del artículo (Tangible o Servicio), su categorización y, opcionalmente, su marca. También concentra su **coste base de referencia** y su **tasa impositiva**. Un producto por sí solo no es vendible si requiere configuración; es la base para generar variantes.

- **Marca Opcional:** `BrandID` puede ser nulo. Esto permite crear productos genéricos o todavía no asociados a fabricante, manteniendo operativas la creación, la venta y la búsqueda del producto.

### Atributos y Valores (`Attribute` / `AttributeValue`)
Representan las dimensiones de personalización (ej. Talla, Color, Material). 
- **Modificadores de Precio:** Los valores de los atributos tienen la capacidad de alterar el coste base del producto de forma positiva o negativa, permitiendo que variantes específicas (ej. Tallas Especiales o Acabados Premium) tengan costes distintos sin necesidad de gestionar precios manuales para cada una.

### La Variante Operativa (`ProductVariant`)
Es la unidad mínima de venta. Representa una combinación única de valores de atributos.
- **SKU Determinista:** El identificador comercial (SKU) de la variante no es aleatorio; se construye algorítmicamente a partir de los códigos del producto y sus atributos, garantizando que el nombre del artículo sea explicativo de su configuración.
- **Identidad JIT (Just-In-Time):** El sistema no requiere que todas las combinaciones posibles existan de antemano. Las variantes pueden nacer en el momento en que se necesitan (ej. al crear un presupuesto), optimizando la base de datos y permitiendo catálogos virtualmente infinitos.

---

## 2. Comportamientos y Reglas Críticas
### Ordenamiento de Atributos por Lista Determinista
El orden en que los atributos se aplican a un producto (para el cálculo de costes y la generación de SKUs) es crítico para la consistencia del catálogo.
- **Orden por `DirectAttributeIDs`**: El orden de herencia y aplicación se determina estrictamente por la posición de los IDs de atributos en la lista `DirectAttributeIDs` del Agregado `Product`.
- **Impacto en SKU**: El SKU de una variante se genera concatenando el SKU base del producto con los códigos de los valores de atributos, siguiendo exactamente el orden definido en la lista del producto.

### Estrategia de Coste Dinámico (`BaseCost`)
A diferencia de otros sistemas, TramaTex **no persiste el coste de las variantes**. El coste se calcula en tiempo real siguiendo este flujo:
1. Se toma el precio base del producto.
2. Se aplican los modificadores de atributos secuencialmente según el orden definido en `DirectAttributeIDs`.
3. Los modificadores pueden ser de tipo **FIXED** (€) o **PERCENTAGE** (%).
**Razón:** Esto garantiza que si el coste de una materia prima (un atributo) cambia, el coste de todas las variantes afectadas se actualice instantáneamente sin procesos de migración masivos.

---
**Última Actualización:** 25-03-2026
