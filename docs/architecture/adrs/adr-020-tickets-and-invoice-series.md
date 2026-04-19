# 🏛️ ADR-020: Tickets (Facturas Simplificadas) y Series de Numeración

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 14-02-2026 |
| **Autores** | Claude AI, Joran (Product Owner) |

---

## 🎯 Contexto
TramaTex debe soportar ventas retail (mostrador) además del flujo B2B tradicional. La legislación española exige diferenciar entre Facturas Completas y Facturas Simplificadas (Tickets), obligando al uso de series de numeración independientes y límites de importe (< 3.000€).

---

## 🔍 Alternativas Consideradas
1. **Entidad Separada `Ticket`:** Separación total pero genera alta duplicidad de lógica contable y fiscal.
2. **Extensión de `Invoice` con `InvoiceType` (Decisión Adoptada):** Reutiliza la infraestructura de Sales, cumple la ley AEAT y minimiza el impacto arquitectónico.

---

## ✅ Decisión Adoptada
Se adopta una implementación híbrida basada en la extensión de la entidad `Invoice` existente.

### Claves del Diseño:
*   **Nuevos Value Objects:** `InvoiceType` (Completa/Simplificada) e `InvoiceSeries`.
*   **Validaciones Legales:** Control automático del límite de 3.000€ para facturas simplificadas.
*   **Tercero Genérico:** Uso de "CONSUMIDOR_FINAL" para tickets rápidos donde no se identifique al cliente.
*   **Preparación para TPV:** Base arquitectónica lista para el desarrollo de una interfaz de caja registradora.

---

## 📈 Consecuencias
### Positivas
*   Cumplimiento legal AEAT desde el primer día.
*   Trazabilidad unificada de toda la facturación en un solo maestro.
*   Mínima deuda técnica en el módulo de Ventas.

### Negativas
*   Mayor responsabilidad para el agregado `Invoice`.
*   El flujo inicial de retail requiere pasar por el proceso de Pedido (optimizado en fases posteriores).

---
[Volver al Índice de ADRs](./README.md)
