# Modelo de Dominio - Módulo Pricing

Este documento describe la lógica de cálculo económico y la jerarquía de reglas que gobiernan los precios en TramaTex. El módulo actúa como un motor de decisión que transforma costes base en precios finales de venta.

---

## 1. Conceptos Fundamentales

### Base Selling Price (BSP)
Es el precio teórico de venta de un producto antes de aplicar condiciones específicas del cliente o promociones. 
- **Composición:** Se calcula sumando al `BaseCost` (del módulo Product) un **Margen de Beneficio** que puede estar definido a nivel de Marca, Categoría o Producto.
- **Dinamicidad:** El BSP no se guarda; se recalcula si el coste de los materiales o el margen de la marca cambian.

### Final Selling Price (FSP)
Es el precio que el cliente ve en su factura o presupuesto.
- **Cálculo:** `BSP - Descuentos Aplicables + Recargos Específicos`.
- **Contexto:** El FSP siempre requiere un `Cliente` y una `Cantidad` para ser calculado, ya que las reglas de volumen o los acuerdos comerciales privados pueden alterarlo.

---

## 2. Jerarquía y Precedencia de Reglas

El motor de precios de TramaTex sigue una lógica de "la regla más específica gana". Cuando existen múltiples reglas aplicables, el sistema las evalúa en este orden:

1.  **Acuerdo Particular (Client Override):** Precios pactados específicamente con un cliente para un producto concreto. Es la prioridad máxima.
2.  **Regla de Producto:** Margen o precio fijo definido para un artículo específico.
3.  **Regla de Marca / Familia:** Márgenes comerciales estandarizados para todos los productos de un fabricante o categoría.
4.  **Regla Global:** Margen por defecto del sistema (fallback) si no hay ninguna otra definición.

### Tipos de Modificadores
- **Porcentuales:** (ej. +20% margen, -10% descuento por volumen).
- **Fijos:** (ej. +5€ por canon de reciclaje o montaje).

---

## 3. Comportamientos Críticos

### Orquestación del Cálculo
El cálculo de un precio no es una simple operación aritmética, sino un proceso de orquestación:
1.  **Resolución de Coste:** Obtiene el coste base actualizado de la variante desde el módulo `Product`.
2.  **Identificación de Reglas:** Busca en el motor de reglas todas las definiciones vigentes (por fecha y contexto).
3.  **Aplicación Cascada:** Aplica primero los márgenes para obtener el BSP y luego los descuentos para el FSP.
4.  **Validación de Moneda:** Garantiza que todas las operaciones se realicen exclusivamente en **Euros (EUR)**.

### Inmutabilidad del Cálculo (`PriceCalculation`)
Cada vez que se genera un precio para un documento oficial (Presupuesto o Factura), el sistema genera un registro de auditoría que congela las reglas que se aplicaron en ese momento. Esto garantiza que, si los márgenes cambian mañana, el documento emitido hoy mantenga su validez legal y coherencia histórica.

---
**Última Actualización:** 2026-03-07
