# 🏛️ ADR-017: Arquitectura del Módulo de Sales

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 06-02-2026 |
| **Autores** | Gemini CLI |

---

## 🎯 Contexto
El módulo de Ventas gestiona el ciclo de vida completo: Cotización -> Pedido -> Albarán -> Factura. La complejidad reside en mantener la trazabilidad documental y permitir la integración fluida con Pricing, Party y Product.

---

## 🔍 Alternativas Consideradas
1. **Entidad Genérica `SalesDocument`:** Menos tablas, pero lógica interna extremadamente compleja y difícil de mantener.
2. **Entidades Separadas por Tipo (Decisión Adoptada):** Claridad total de dominio. Cada documento encapsula su propio estado y reglas de negocio.

---

## ✅ Decisión Adoptada
Se adopta un modelo de **Entidades Independientes** para cada hito del flujo comercial.

### Aspectos Clave:
*   **Soberanía Comercial (Manual Override):** El sistema sugiere precios (vía Pricing), pero el comercial puede sobreescribirlos manualmente. Se almacenan ambos valores para auditoría.
*   **Cálculo de Impuestos:** Responsabilidad delegada en el módulo Sales para asegurar la precisión fiscal en el documento final.
*   **Flujo Agilizado (Actualización Abril 2026):** Los pedidos nacen directamente en estado `EN_PREPARACION` (Confirmado) para agilizar la operativa, reservando `PENDIENTE` como estado de seguridad tras reactivaciones.
*   **Consistencia Financiera:** Uso obligatorio de tipos decimales para evitar discrepancias de redondeo en facturación.

---

## 📈 Consecuencias
### Positivas
*   Trazabilidad documental impecable.
*   Facilidad para implementar reglas específicas (ej: facturación agrupada).
*   Modelo alineado con Clean Architecture y DDD.

### Negativas
*   Ligera redundancia de datos (ej: duplicación de `PartyID` en la cadena documental) aceptada en favor del desacoplamiento.

---
[Volver al Índice de ADRs](./README.md)
