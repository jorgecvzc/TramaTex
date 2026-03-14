# Modelo de Dominio - Módulo Product

Este documento describe la lógica de estructuración del catálogo de TramaTex, centrándose en el sistema de variantes dinámicas y la gestión de productos configurables.

---

## 1. Agregados y Entidades de Dominio

### El Producto como Plantilla (`Product`)
El `Product` actúa como el concepto base o "molde". Define la naturaleza del artículo (Tangible o Servicio), su marca, su categorización y, lo más importante, su **coste base de referencia** y su **tasa impositiva**. Un producto por sí solo no es vendible si requiere configuración; es la base para generar variantes.

### Atributos y Valores (`Attribute` / `AttributeValue`)
Representan las dimensiones de personalización (ej. Talla, Color, Material). 
- **Modificadores de Precio:** Los valores de los atributos tienen la capacidad de alterar el coste base del producto de forma positiva o negativa, permitiendo que variantes específicas (ej. Tallas Especiales o Acabados Premium) tengan costes distintos sin necesidad de gestionar precios manuales para cada una.

### La Variante Operativa (`ProductVariant`)
Es la unidad mínima de venta. Representa una combinación única de valores de atributos.
- **SKU Determinista:** El identificador comercial (SKU) de la variante no es aleatorio; se construye algorítmicamente a partir de los códigos del producto y sus atributos, garantizando que el nombre del artículo sea explicativo de su configuración.
- **Identidad JIT (Just-In-Time):** El sistema no requiere que todas las combinaciones posibles existan de antemano. Las variantes pueden nacer en el momento en que se necesitan (ej. al crear un presupuesto), optimizando la base de datos y permitiendo catálogos virtualmente infinitos.

---

## 2. Comportamientos y Reglas Críticas

### Estrategia de Coste Dinámico (`BaseCost`)
A diferencia de otros sistemas, TramaTex **no persiste el coste de las variantes**. El coste se calcula en tiempo real siguiendo este flujo:
1. Se toma el precio base del producto.
2. Se aplican los modificadores porcentuales sobre el base.
3. Se añaden/restan los modificadores fijos.
**Razón:** Esto garantiza que si el coste de una materia prima (un atributo) cambia, el coste de miles de variantes se actualice instantáneamente sin procesos de migración masivos.

### Resolución de Atributos por Precedencia
El sistema decide qué atributos "ve" un producto basándose en una jerarquía de herencia:
- **Asignación Directa:** Atributos específicos del producto.
- **Herencia por Grupo/Marca:** Atributos comunes a una familia de productos.
La lógica de dominio garantiza que el nivel más específico siempre anula al más general, permitiendo excepciones controladas en el catálogo.

---
**Última Actualización:** 2026-03-07
