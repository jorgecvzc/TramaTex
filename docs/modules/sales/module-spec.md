# Módulo de Sales (Gestión de Órdenes)

## 1. Propósito

*   **Visión del Módulo:** Gestionar cotizaciones, órdenes y seguimiento de ventas.
*   **Objetivos Clave:**
    *   Proporcionar un sistema completo para gestionar el ciclo de vida de las ventas.
    *   Manejar desde la cotización hasta la entrega de la orden.

## 2. Requisitos

### 2.1. Requisitos Funcionales

*   **RF-001:** Crear cotizaciones.
*   **RF-002:** Convertir cotizaciones a órdenes de venta.
*   **RF-003:** Gestionar líneas de orden (productos, cantidades, precios).
*   **RF-004:** Realizar seguimiento del estado de la orden.
*   **RF-005:** Mantener un historial de cambios en la orden.

## 3. Casos de Uso

### 3.1. Actores

*   **Vendedor:** Crea cotizaciones y gestiona órdenes.
*   **Cliente:** Aprueba cotizaciones y recibe órdenes.

### 3.2. Casos de Uso Principales

Para una lista completa y detallada de los casos de uso, incluyendo flujos y entradas/salidas, consulte el documento [Casos de Uso - Módulo Sales](./use-cases.md).


## 4. Historias de Usuario

*   **HU-001:** Como Vendedor, quiero crear una cotización para un cliente potencial para presentarle una oferta formal.
*   **HU-002:** Como Vendedor, quiero convertir una cotización en una orden de venta cuando el cliente la aprueba para iniciar el proceso de producción.
*   **HU-003:** Como Vendedor, quiero poder ver el estado de todas mis órdenes para hacer seguimiento con los clientes.

## 5. Criterios de Aceptación

*   **Para HU-002:**
    *   **Criterio 1:** Dado una cotización en estado 'Aprobada', cuando la convierto a orden, entonces se genera un número de orden único y el estado cambia a 'Confirmada'.

## 6. Modelo de Dominio

### SalesOrder (Raíz de Agregación)
- **ID**: UUID
- **OrderNumber**: String (único, formato: ORD-YYYYMMDD-XXXX)
- **PartyID**: UUID (cliente)
- **OrderDate**: DateTime
- **DeliveryDate**: DateTime (requerido)
- **Estado**: Enum (Cotización, Confirmada, En Preparación, Entregada, Cancelada)
- **LineItems**: List<OrderLineItem>
- **Subtotal**: Decimal
- **Tax**: Decimal
- **Total**: Decimal
- **Notes**: String

### OrderLineItem
- **ID**: UUID
- **OrderID**: UUID
- **ProductVariantID**: UUID
- **Cantidad**: Integer
- **PrecioUnitario**: Decimal
- **Descuento**: Decimal (opcional)
- **Subtotal**: Decimal

## 7. Decisiones de Diseño

*   **Flujo de Estados de Orden:**
    *   `Cotización` -> `Confirmada` -> `En Preparación` -> `Entregada`
    *   Se puede cancelar desde cualquier estado (`Cancelada`).
*   **Relaciones con Otros Módulos:**
    *   **Party**: La orden pertenece a un `Party` (cliente).
    *   **Product**: Las `LineItems` referencian a un `ProductVariant`.
    *   **Pricing**: El `PrecioUnitario` es calculado por el motor de precios.
    *   **MES (Fase 3, MVP):** El estado 'En Preparación' podría iniciar un seguimiento en el sistema de producción.
*   **Fases de Desarrollo:**
    *   [X] Fase 1 (MVP): CRUD básico de órdenes y cambio de estado.
    *   [ ] Fase 2: Gestión de cotizaciones y conversión a órdenes.
    *   [ ] Fase 3: Integración con producción (MES).