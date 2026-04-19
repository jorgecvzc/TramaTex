# 🏛️ ADR-015: Arquitectura del Módulo de Product

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 06-02-2026 |
| **Autores** | Gemini CLI |

---

## 🎯 Contexto
El módulo `Product` es la fuente de verdad para todos los bienes y servicios vendibles. El desafío radica en gestionar productos con alta configurabilidad (tallas, colores, tejidos) sin colapsar el catálogo con miles de combinaciones manuales.

---

## 🔍 Alternativas Consideradas
1. **Modelo Simple:** Fácil de implementar pero rígido ante atributos variables.
2. **Uso de JSON:** Flexible, pero pierde el tipado fuerte y dificulta las consultas complejas y el reporte.
3. **Modelo de Atributos con Herencia y JIT (Decisión Adoptada):** Robusto, *type-safe* y optimizado para la creación bajo demanda.

---

## ✅ Decisión Adoptada
Se adopta el **Modelo de Atributos/Valores explícito con creación Just-In-Time (JIT)**.

### Componentes Core:
*   **Atributos y Valores:** Definición flexible de características (ej: Talla -> XL).
*   **Producto Base:** Plantilla con SKU raíz y marca.
*   **Variante JIT:** Instancia vendible con SKU compuesto (ej: CAM-BLA-XL) que se crea automáticamente al ser demandada por primera vez en una venta o presupuesto.
*   **Herencia de Atributos:** Capacidad de definir reglas de precios o propiedades a nivel de Marca o Grupo que se propagan a los productos.

---

## 📈 Consecuencias
### Positivas
*   Catálogo extremadamente ágil: el usuario solo gestiona productos base.
*   Rendimiento optimizado: solo existen en la BD las variantes que realmente se mueven.
*   Identificación única y determinista de productos mediante SKUs compuestos.

### Negativas
*   Mayor complejidad lógica en el motor de resolución de atributos y SKUs.

---
[Volver al Índice de ADRs](./README.md)
