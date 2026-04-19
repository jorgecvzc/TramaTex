# 🏛️ Modelo de Dominio - Módulo Pricing

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Vigente |

---

## 🎯 Propósito
Este documento describe la lógica de cálculo económico y la jerarquía de reglas que gobiernan los precios en TramaTex. El módulo actúa como un motor de decisión que transforma costes base en precios finales de venta.

---

## 1. Conceptos Fundamentales

### Precio Base de Venta (PBV)
Es el precio sugerido de un producto antes de aplicar condiciones específicas del cliente o promociones.
*   **Composición:** Se obtiene sumando al coste base (procedente del módulo *Product*) un **Margen de Beneficio** definido a nivel de Marca, Familia o Producto.
*   **Precisión:** Todos los cálculos internos utilizan tipos **Decimales** para evitar errores de redondeo.

### Precio Final de Venta (PFV)
Es el precio definitivo que se aplica en el documento mercantil (Pedido o Factura).
*   **Cálculo:** `PBV - Descuentos + Recargos`.
*   **Contexto:** Siempre requiere un Cliente y una Cantidad para su determinación, ya que depende de acuerdos comerciales y reglas de volumen.

---

## 2. Jerarquía y Precedencia de Reglas
El motor de precios sigue el principio de "especificidad máxima". Las reglas se evalúan en este orden:

1.  **Acuerdo Particular:** Precios pactados manualmente con un cliente. Es la prioridad máxima.
2.  **Regla de Producto:** Margen específico para un artículo concreto.
3.  **Regla de Marca / Familia:** Márgenes estandarizados por fabricante o categoría.
4.  **Regla Global:** Margen de seguridad del sistema (fallback).

---

## 3. Comportamientos Críticos

### Inmutabilidad del Cálculo (`PriceCalculation`)
Cada vez que se genera un precio para un presupuesto o factura, el sistema congela las reglas aplicadas en ese instante. Esto garantiza coherencia histórica: si los precios suben mañana, los documentos emitidos hoy mantienen sus valores originales.

### Integridad Financiera
*   **Moneda Única:** El sistema valida que todas las operaciones se realicen exclusivamente en **Euros (EUR)** para el MVP.
*   **Caché Distribuida:** El PBV se almacena en **Redis** para optimizar el rendimiento de la interfaz de ventas, invalidándose automáticamente ante cualquier cambio en el catálogo.

---
[Volver al Módulo de Pricing](./README.md)
