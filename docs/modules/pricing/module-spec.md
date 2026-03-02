# Módulo de Pricing (Motor de Precios)

**Estado:** ✅ **COMPLETO (100%)**  
**Última actualización:** 1 de marzo de 2026

## Estado de Implementación

### Componentes Completos ✅
- **Reglas de Precios (Price Rules):**
  - Backend: Definición de reglas por producto, categoría de cliente y volumen.
  - Soporte para márgenes de marca (BrandProfitMargin) y descuentos.
- **Cálculo Dinámico:**
  - Endpoint `/api/pricing/calculate` funcional con desglose de costes y márgenes.
- **Historial y Auditoría:**
  - Registro completo de cada cálculo realizado para trazabilidad.
- **Frontend:**
  - Calculadora interactiva integrada en el detalle de producto.
  - Visualización de historial de precios y reglas aplicadas.

---

## 1. Propósito

*   **Visión del Módulo:** Proporcionar un motor de cálculo de precios automatizado y flexible que asegure la rentabilidad de TramaTex mediante la aplicación consistente de márgenes y descuentos.
*   **Objetivos Clave:**
    *   Proporcionar un motor de precios inteligente que calcule automáticamente el precio de venta.
    *   Calcular precios según categoría de cliente, producto, volumen y otras variables.

---

## 2. Requisitos

... [resto de requisitos permanecen igual] ...

---

## 7. Decisiones de Diseño

*   **Algoritmo de Cálculo:**
    1.  Obtener `CostoBase` del `ProductVariant`.
    2.  Aplicar `BrandProfitMargin` configurado.
    3.  Buscar y aplicar `PriceRule` específica por cliente o categoría.
    4.  Aplicar descuentos por volumen (SalesDiscountRule).
    4b. Si no hay reglas específicas, aplicar `DefaultDiscountPercentage` del cliente (del módulo Party) como descuento fallback.
    5.  Calcular precio final con IVA (`finalPriceWithTax = finalPrice × (1 + taxRate/100)`).
    6.  Validar margen de contribución mínimo.
    7.  Registrar el cálculo en `PriceCalculation`.
*   **Relaciones con Otros Módulos:**
    *   **Product**: Lee el coste base de la variante y el tipo impositivo (taxRate).
    *   **Party**: Obtiene la categoría, descuentos específicos del cliente y el `defaultDiscountPercentage` vía `PartyPricingClient` (anti-corruption layer).
    *   **Sales**: Provee los precios unitarios definitivos (con y sin IVA) para presupuestos y pedidos.

---

## 8. Fases de Desarrollo

*   [x] **Fase 1 (MVP):** Cálculo básico (CostoBase + Márgenes).
*   [x] **Fase 2:** Descuentos por volumen y categoría.
*   [x] **Fase 3:** Reglas dinámicas e historial de auditoría.
*   [x] **Fase 4:** Integración visual completa en el catálogo.

