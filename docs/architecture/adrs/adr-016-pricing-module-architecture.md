# 🏛️ ADR-016: Arquitectura del Módulo de Pricing

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 06-02-2026 |
| **Autores** | Gemini CLI |

---

## 🎯 Contexto
El sector textil requiere una lógica de precios compleja: márgenes dinámicos, descuentos por volumen, promociones y variantes Just-In-Time. El sistema necesita un motor exacto, de alto rendimiento y desacoplado de la infraestructura.

---

## 🔍 Alternativas Consideradas
1. **Regla Genérica Única:** Simple pero incapaz de gestionar la complejidad de precedencias.
2. **Uso de JSON:** Flexible pero peligroso para cálculos financieros al perder el tipado fuerte (*type-safety*).
3. **Reglas Separadas con Value Objects (Decisión Adoptada):** Claridad de dominio, seguridad total de tipos y lógica de aplicación predecible.

---

## ✅ Decisión Adoptada
Se implementa un motor de precios basado en dos tipos de reglas y una estrategia de caché agresiva:

### Claves del Diseño:
*   **Diferenciación de Reglas:**
    *   `BaseSalesPriceRule`: Define el precio de venta estándar (coste + margen).
    *   `SaleModificationRule`: Aplica descuentos o recargos en el momento de la venta.
*   **Value Objects Financieros:** Uso de `Money` (fijado en EUR) y `Percentage` con **precisión decimal**.
*   **Filosofía "No JSON":** Todos los criterios (Clientes, Grupos, Productos) son campos explícitos y tipados.
*   **Caché con Redis:** Los precios calculados se almacenan en Redis por `ProductID` para garantizar una respuesta instantánea en la UI de ventas. La caché se invalida automáticamente ante cambios en el catálogo o variantes JIT.

---

## 📈 Consecuencias
### Positivas
*   Reducción drástica de errores en cálculos financieros.
*   Rendimiento óptimo incluso con catálogos extensos.
*   Modelo de dominio muy cercano al lenguaje del negocio.

### Negativas
*   Mayor número de tablas y entidades.
*   Complejidad inicial en la gestión de precedencias de reglas.

---
[Volver al Índice de ADRs](./README.md)
