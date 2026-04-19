# 🔌 Contratos de API - Módulo Pricing

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Vigente |
| **Relacionado con** | [ADR-016](../../architecture/adrs/adr-016-pricing-module-architecture.md) |

---

## 🎯 Propósito
Este documento detalla la interfaz de integración técnica con el motor de cálculo de precios y la gestión de reglas económicas del sistema TramaTex.

---

## 1. Motores de Cálculo (`/api/pricing/calculate`)
El sistema ofrece endpoints especializados según la fase del proceso de negocio:

*   **Cálculo de Precio Base de Venta (`POST /api/pricing/base-sales-price/calculate`):** 
    Proporciona el precio sugerido antes de aplicar condiciones de cliente. Esencial para la visualización del catálogo y catálogos públicos.
*   **Cálculo de Precio Final de Venta (`POST /api/pricing/final-sale-price/calculate`):** 
    Motor orquestador que evalúa al cliente, la cantidad y las promociones vigentes. Devuelve un desglose detallado de las reglas aplicadas para transparencia en el punto de venta.

---

## 2. Gestión de Reglas (`/api/pricing/rules`)
Endpoints para administrar las definiciones que alimentan el motor:

*   **Reglas Base de Venta:** Definición de márgenes porcentuales o fijos por Marca o Categoría de producto.
*   **Reglas de Modificación de Venta:** Gestión de descuentos por volumen, campañas temporales o recargos logísticos.

---

## 3. Acuerdos Particulares (`/api/pricing/client-overrides`)
Gestiona excepciones a la política de precios general:

*   **Precios Pactados:** Permite fijar un precio inamovible para una combinación específica de Cliente y Variante de Producto.

---

## 🏗️ Estructura de Datos Común: `MoneyDTO`
Para garantizar la integridad financiera, la API utiliza un Objeto de Transferencia de Datos específico para importes:
*   `amount`: Valor numérico de alta precisión (Decimal).
*   `currency`: Código de moneda (siempre "EUR" para el MVP).

---
[Volver al Módulo de Pricing](./README.md)
