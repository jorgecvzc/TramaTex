# 🏛️ ADR-013: Manejo de Modificaciones de Producto (Servicios)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 02-02-2026 |
| **Autores** | Gemini, Usuario |

---

## 🎯 Contexto
TramaTex debe gestionar tanto productos tangibles como servicios de personalización (marcajes, bordados) o arreglos (ajuste de bajos). Estos servicios pueden aplicarse a productos nuevos o a prendas que el cliente ya posee. El modelo debe ser eficiente para el inventario, la tarificación y la ejecución en taller (MES).

---

## 🔍 Alternativas Consideradas
1. **Como Variantes del Producto:** Las modificaciones serían una opción más del producto base. Complejo de gestionar cuando el producto es del cliente y el precio base es cero.
2. **Como Productos de Tipo "Servicio" (Decisión Adoptada):** Las modificaciones son entidades independientes en el catálogo.

---

## ✅ Decisión Adoptada
Se adopta la **Alternativa B: Modificaciones como Productos de Tipo "Servicio" independientes**.

### Justificación:
*   **Claridad:** Diferencia explícitamente entre bienes (físicos) y acciones (servicios).
*   **Flexibilidad:** Permite vender un "Bordado de Logo" de forma aislada o vinculado a una prenda.
*   **Simplicidad en MES:** El taller ve cada servicio como una tarea discreta y clara en la orden de trabajo.
*   **Escalabilidad:** Facilita la creación de plantillas de servicios personalizadas por cliente (`PartyServiceConfiguration`).

---

## 📈 Consecuencias
### Positivas
*   Catálogo intuitivo y fácil de navegar.
*   Gestión natural del escenario "producto propio del cliente".
*   Integración simplificada con el motor de producción (MES).

### Negativas
*   La relación producto-servicio en un pedido es implícita (múltiples líneas de pedido), requiriendo un diseño de UI cuidadoso para mantener la cohesión visual.

---
[Volver al Índice de ADRs](./README.md)
