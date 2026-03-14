# Contratos de API - Módulo Pricing

Este documento detalla la interfaz de integración con el motor de cálculo de precios y la gestión de reglas económicas en TramaTex.

---

## 1. Motores de Cálculo (`/api/pricing/calculate`)

El sistema ofrece endpoints especializados según la fase del proceso de venta.
- **Cálculo de Precio Base (BSP) (`POST /api/pricing/base-sales-price/calculate`):** Proporciona el precio de venta sugerido antes de aplicar condiciones de cliente. Útil para visualización en el catálogo general.
- **Cálculo de Precio Final (FSP) (`POST /api/pricing/final-sale-price/calculate`):** Es el motor orquestador que toma en cuenta al cliente, la cantidad y las promociones vigentes. Devuelve el desglose de reglas aplicadas para que el frontend pueda mostrar los descuentos de forma transparente.

## 2. Gestión de Reglas (`/api/pricing/rules`)

Controla las definiciones que alimentan los motores de cálculo.
- **Reglas Base de Venta (`/api/pricing/base-sales-rules`):** Permite definir márgenes porcentuales o fijos por Marca o Categoría.
- **Reglas de Modificación (`/api/pricing/sale-modification-rules`):** Define descuentos por volumen, promociones temporales o recargos logísticos.

## 3. Acuerdos Particulares (`/api/pricing/client-overrides`)

Gestiona la excepción a la regla general.
- **Precios Pactados:** Permite fijar un precio específico para un binomio Cliente-Variante que anula cualquier otro cálculo automático del sistema.

## 4. Auditoría de Precios (`/api/pricing/history`)

- **Trazabilidad de Cambios:** Endpoint para consultar la evolución de los precios de una variante, permitiendo entender por qué un presupuesto antiguo tenía un precio diferente al actual.

---

## Estructura de Respuesta Común: `MoneyDTO`

Todas las cantidades económicas devueltas por la API siguen una estructura de objeto de valor (Value Object) para evitar errores de precisión y moneda:
- `amount`: Valor numérico (ej. 150.50).
- `currency`: Identificador de moneda (siempre "EUR").

---
**Última Actualización:** 2026-03-07
