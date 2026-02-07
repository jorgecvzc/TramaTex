# Casos de Uso - Módulo Sales

Este documento detalla los casos de uso para el módulo de Ventas, abarcando la gestión de Cotizaciones, Pedidos de Venta, Albaranes y Facturas. Estos casos de uso orquestan la lógica de negocio, interactuando con las entidades de dominio definidas y otros módulos como Party, Product y Pricing.

---

#### **1. Casos de Uso de `Quote` (Presupuesto / Cotización)**

*   **CU-S-001: CreateQuote**
    *   **Propósito:** Crear un nuevo presupuesto para un cliente.
    *   **Actores:** Vendedor.
    *   **Entradas:** `PartyID` del cliente, lista de `{ ProductVariantID, Quantity, ManualUnitPrice (opcional), ManualDiscountPerUnit (opcional) }` para los ítems, `ExpirationDate`, `Notes`.
    *   **Flujo:**
        1.  Validar entradas.
        2.  Para cada ítem:
            a.  Obtener `BaseSalesPrice` del `Pricing` module (`CalculateBaseSalesPriceForProductVariantUseCase`).
            b.  Calcular `CalculatedUnitPrice` y `CalculatedDiscountPerUnit` usando el `CalculateFinalSalePriceUseCase` del `Pricing` module si no hay overrides manuales.
            c.  Aplicar overrides manuales si se proporcionan (`ManualUnitPrice`, `ManualDiscountPerUnit`).
            d.  Crear `QuoteLineItem` con los precios finales.
        3.  Calcular `Subtotal`, `TaxAmount` (lógica interna de Sales), `Total` para la `Quote`.
        4.  Generar `QuoteNumber`.
        5.  Persistir `Quote` y `QuoteLineItem`s.
    *   **Salida:** `Quote` (con todos sus detalles).

*   **CU-S-002: GetQuote**
    *   **Propósito:** Recuperar los detalles de un presupuesto específico.
    *   **Actores:** Vendedor.
    *   **Entradas:** `QuoteID`.
    *   **Salida:** `Quote` (con todos sus detalles).

*   **CU-S-003: ListQuotes**
    *   **Propósito:** Listar presupuestos con opciones de filtrado y paginación.
    *   **Actores:** Vendedor.
    *   **Entradas:** Filtros (ej. `PartyID`, `Status`, `QuoteDateRange`), paginación.
    *   **Salida:** Lista paginada de `Quote`s.

*   **CU-S-004: UpdateQuote**
    *   **Propósito:** Modificar un presupuesto existente (solo si `Status` es `BORRADOR`).
    *   **Actores:** Vendedor.
    *   **Entradas:** `QuoteID`, datos a actualizar (ítems, fechas, notas, etc.).
    *   **Salida:** `Quote` actualizada.

*   **CU-S-005: ChangeQuoteStatus**
    *   **Propósito:** Actualizar el estado de un presupuesto.
    *   **Actores:** Vendedor.
    *   **Entradas:** `QuoteID`, `NewStatus`.
    *   **Salida:** `Quote` con estado actualizado.

*   **CU-S-006: ConvertQuoteToOrder**
    *   **Propósito:** Convertir un presupuesto `APROBADO` en un pedido de venta.
    *   **Actores:** Vendedor.
    *   **Entradas:** `QuoteID`, `DeliveryDate`.
    *   **Flujo:**
        1.  Validar que `Quote.Status` es `APROBADA`.
        2.  Crear un `SalesOrder` copiando los datos relevantes de la `Quote` (cliente, ítems, precios finales).
        3.  Asignar el `QuoteID` al `SalesOrder`.
        4.  Establecer `SalesOrder.Status` a `PENDIENTE`.
        5.  Persistir `SalesOrder` y sus `OrderLineItem`s.
        6.  Cambiar `Quote.Status` a `CONVERTIDA_A_PEDIDO`.
    *   **Salida:** `SalesOrder` recién creado.

#### **2. Casos de Uso de `SalesOrder` (Pedido de Venta)**

*   **CU-S-007: CreateOrder**
    *   **Propósito:** Crear un nuevo pedido de venta directamente (sin cotización previa).
    *   **Actores:** Vendedor.
    *   **Entradas:** `PartyID` del cliente, lista de `{ ProductVariantID, Quantity, ManualUnitPrice (opcional), ManualDiscountPerUnit (opcional) }` para los ítems, `DeliveryDate`, `Notes`.
    *   **Flujo:** Similar a `CreateQuote`, pero establece `SalesOrder.Status` a `PENDIENTE` y genera `OrderNumber`.
    *   **Salida:** `SalesOrder`.

*   **CU-S-008: GetOrder**
    *   **Propósito:** Recuperar los detalles de un pedido específico.
    *   **Actores:** Vendedor.
    *   **Entradas:** `SalesOrderID`.
    *   **Salida:** `SalesOrder`.

*   **CU-S-009: ListOrders**
    *   **Propósito:** Listar pedidos con opciones de filtrado y paginación.
    *   **Actores:** Vendedor.
    *   **Entradas:** Filtros (ej. `PartyID`, `Status`, `OrderDateRange`), paginación.
    *   **Salida:** Lista paginada de `SalesOrder`s.

*   **CU-S-010: UpdateOrderDetails**
    *   **Propósito:** Modificar detalles de un pedido (ej. `DeliveryDate`, `Notes`, `PartyID` si el estado lo permite). No modifica ítems.
    *   **Actores:** Vendedor.
    *   **Entradas:** `SalesOrderID`, datos a actualizar.
    *   **Salida:** `SalesOrder` actualizado.

*   **CU-S-011: ChangeOrderStatus**
    *   **Propósito:** Actualizar el estado de un pedido.
    *   **Actores:** Vendedor, Sistema (ej. módulo MES).
    *   **Entradas:** `SalesOrderID`, `NewStatus`.
    *   **Flujo:** Validar transición de estado.
    *   **Salida:** `SalesOrder` con estado actualizado.

*   **CU-S-012: AddLineItemToOrder**
    *   **Propósito:** Añadir un ítem de línea a un pedido (solo si el estado lo permite, ej. `PENDIENTE`).
    *   **Actores:** Vendedor.
    *   **Entradas:** `SalesOrderID`, `{ ProductVariantID, Quantity, ManualUnitPrice (opcional), ManualDiscountPerUnit (opcional) }`.
    *   **Flujo:** Calcula precios, añade ítem, recalcula totales de `SalesOrder`.
    *   **Salida:** `SalesOrder` actualizado.

*   **CU-S-013: UpdateOrderLineItem**
    *   **Propósito:** Modificar un ítem de línea existente en un pedido.
    *   **Actores:** Vendedor.
    *   **Entradas:** `SalesOrderID`, `OrderLineItemID`, nuevos datos.
    *   **Salida:** `SalesOrder` actualizado.

*   **CU-S-014: RemoveOrderLineItem**
    *   **Propósito:** Eliminar un ítem de línea de un pedido.
    *   **Actores:** Vendedor.
    *   **Entradas:** `SalesOrderID`, `OrderLineItemID`.
    *   **Salida:** `SalesOrder` actualizado.

#### **3. Casos de Uso de `DeliveryNote` (Albarán)**

*   **CU-S-015: CreateDeliveryNote**
    *   **Propósito:** Generar un albarán para un `SalesOrder`. Puede ser parcial o completo.
    *   **Actores:** Vendedor, Personal de Almacén.
    *   **Entradas:** `SalesOrderID`, lista de `{ OrderLineItemID, DeliveredQuantity }`, `DeliveryDate`.
    *   **Flujo:**
        1.  Validar `DeliveredQuantity` contra la `Quantity` del `OrderLineItem` original.
        2.  Generar `DeliveryNoteNumber`.
        3.  Persistir `DeliveryNote` y `DeliveryNoteLineItem`s.
        4.  Actualizar `SalesOrder.Status` (ej. a `ENTREGADO_PARCIALMENTE` o `ENTREGADO`).
    *   **Salida:** `DeliveryNote`.

*   **CU-S-016: GetDeliveryNote**
    *   **Propósito:** Recuperar los detalles de un albarán.
    *   **Actores:** Vendedor, Personal de Almacén.
    *   **Entradas:** `DeliveryNoteID`.
    *   **Salida:** `DeliveryNote`.

*   **CU-S-017: ListDeliveryNotes**
    *   **Propósito:** Listar albaranes con opciones de filtrado (ej. por `SalesOrderID`, `PartyID`, `DeliveryDateRange`).
    *   **Actores:** Vendedor, Personal de Almacén.
    *   **Entradas:** Filtros, paginación.
    *   **Salida:** Lista paginada de `DeliveryNote`s.

#### **4. Casos de Uso de `Invoice` (Factura)**

*   **CU-S-018: CreateInvoiceFromOrder**
    *   **Propósito:** Generar una factura a partir de un `SalesOrder` (o un conjunto de `SalesOrder`s completados).
    *   **Actores:** Vendedor, Personal de Contabilidad.
    *   **Entradas:** `SalesOrderID` (o lista de `SalesOrderID`s), `InvoiceDate`, `DueDate`.
    *   **Flujo:**
        1.  Recuperar `SalesOrder`(s) y sus `OrderLineItem`s.
        2.  Crear `InvoiceLineItem`s a partir de los `OrderLineItem`s, utilizando sus `FinalUnitPrice` y `FinalDiscountPerUnit`.
        3.  Calcular `Subtotal`, `TaxAmount`, `Total` de la factura.
        4.  Generar `InvoiceNumber`.
        5.  Persistir `Invoice` y `InvoiceLineItem`s.
        6.  Actualizar el `SalesOrder.Status` a `FACTURADO_PARCIALMENTE` o `FACTURADO_COMPLETAMENTE`.
    *   **Salida:** `Invoice`.

*   **CU-S-019: GetInvoice**
    *   **Propósito:** Recuperar los detalles de una factura específica.
    *   **Actores:** Vendedor, Personal de Contabilidad.
    *   **Entradas:** `InvoiceID`.
    *   **Salida:** `Invoice`.

*   **CU-S-020: ListInvoices**
    *   **Propósito:** Listar facturas con opciones de filtrado (ej. por `PartyID`, `Status`, `InvoiceDateRange`, `DueDateRange`).
    *   **Actores:** Vendedor, Personal de Contabilidad.
    *   **Entradas:** Filtros, paginación.
    *   **Salida:** Lista paginada de `Invoice`s.
