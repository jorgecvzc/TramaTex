# 🏛️ Modelo de Dominio - Módulo Product

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Vigente |

---

## 🎯 Propósito
Este documento describe la lógica de estructuración del catálogo de TramaTex, centrándose en el sistema de variantes dinámicas y la gestión de productos configurables mediante el estándar JIT (Just-In-Time).

---

## 1. Agregados y Entidades de Dominio

### El Producto como Plantilla (`Product`)
El `Product` actúa como el "molde" o concepto base. Define la naturaleza del artículo (**Bien Tangible** o **Servicio**), su categorización y marca. Concentra el **coste base de referencia** y la **tasa impositiva (IVA)**. Un producto es la base para generar variantes vendibles.

### Atributos y Valores (`Attribute` / `AttributeValue`)
Representan las dimensiones de personalización (ej: Talla, Color, Material).
*   **Modificadores de Coste:** Los valores de los atributos pueden alterar el coste base del producto de forma positiva o negativa. Esto permite que variantes específicas (ej: Tallas Especiales) tengan costes distintos de forma automática.

### La Variante Operativa (`ProductVariant`)
Es la unidad mínima de venta. Representa una combinación única de valores de atributos.
*   **SKU Determinista:** El identificador comercial (SKU) se construye algorítmicamente a partir del producto y sus atributos (ej: `POLO-AZU-XL`).
*   **Identidad Just-In-Time (JIT):** Las variantes se crean automáticamente en el momento en que se demandan por primera vez (ej: al añadirlas a un presupuesto), evitando la carga manual de miles de combinaciones posibles.

---

## 2. Comportamientos y Reglas Críticas

### Ordenamiento de Atributos Determinista
El orden en que los atributos se asocian a un producto es crítico para la generación del SKU y el cálculo de costes. El sistema respeta estrictamente la posición de los atributos definida en la ficha del producto.

### Estrategia de Coste Dinámico
TramaTex no almacena el coste final de cada variante, sino que lo calcula en tiempo real:
1.  **Coste Base:** Precio inicial del producto.
2.  **Modificadores:** Se aplican secuencialmente (Tipo **FIJO** en € o **PORCENTUAL** en %).
Esto garantiza que un cambio en el coste de una materia prima se refleje instantáneamente en todas las variantes afectadas sin migraciones de datos masivas.

---
[Volver al Módulo de Product](./README.md)
