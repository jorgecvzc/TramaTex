# Casos de Uso - Módulo Pricing

- **Version:** 1.0.0
- **Status:** Implementado

Este documento describe los casos de uso para el motor de precios y la gestión de reglas comerciales en TramaTex.

---

## 1. Cálculo de Precios de Venta

### UC-PRI-001: Calcular Precio Final de Variante

- **Actor:** Módulo Sales / Sistema Externo
- **Descripción:** Calcula el precio de venta final para una combinación de producto, cliente y cantidad, aplicando la jerarquía de reglas de negocio.
- **Flujo Lógico del Motor:**
  1. **Obtención del Coste:** Recupera el `BaseCost` de la variante desde el módulo `Product`.
  2. **Aplicación de Margen de Marca:** Si el producto tiene marca, aplica el `BrandProfitMargin` configurado para ella. Si `BrandID` es nulo, este paso se omite sin error.
  3. **Evaluación de Reglas Generales:** Busca y aplica `PricingRule` que coincidan con la categoría de la party y el producto.
  4. **Evaluación de Precios Específicos:** Busca `ClientPricing` (overrides) para ese cliente y producto específico. Si existe, tiene prioridad sobre las reglas generales.
  5. **Aplicación de Descuentos Fallback:** Si no hay reglas específicas, aplica el `DefaultDiscountPercentage` de la ficha del cliente (módulo `Party`).
  6. **Cálculo de Impuestos:** Aplica el `TaxRate` (IVA) definido en el producto sobre el precio neto resultante.
- **Resultado:** `CalculatePriceResponse` con el desglose de cada paso.

---

## 2. Gestión de Reglas Maestras

### UC-PRI-002: Configurar Margen por Marca

- **Actor:** Administrador Comercial
- **Descripción:** Define el porcentaje de beneficio por defecto que se aplica a todos los productos de una marca específica.

### UC-PRI-003: Crear Regla de Precio por Cantidad (Escalados)

- **Actor:** Administrador Comercial
- **Descripción:** Define precios o descuentos que solo se activan a partir de un volumen de compra (`MinQuantity`).

### UC-PRI-004: Establecer Precio Especial por Cliente (Override)

- **Actor:** Administrador Comercial
- **Descripción:** Define un precio fijo o descuento excepcional para un cliente y producto concreto, anulando la lógica general de márgenes.

---

## 3. Auditoría y Trazabilidad

### UC-PRI-005: Consultar Historial de Cálculos

- **Actor:** Administrador / Auditor
- **Descripción:** Recupera los registros inmutables de `PriceCalculation` para verificar qué reglas se aplicaron en una venta pasada y garantizar la transparencia de precios.
