# Modelo de Dominio - Módulo Sales

Este documento describe el modelo de dominio consolidado para el módulo de Ventas, incluyendo las entidades de documentos de venta (`Quote`, `SalesOrder`, `DeliveryNote`, `Invoice`), sus ítems de línea asociados y los Value Objects clave. Este modelo está alineado con la estrategia de "manual override" y la integración con el módulo de Pricing.

---

### **Value Objects Clave**

Los siguientes Value Objects son fundamentales para la consistencia del módulo de Ventas. Algunos son reutilizados de otros módulos o son específicos de Sales.

*   **`Money`:** Representa una cantidad monetaria (fijo a EUR para MVP).
*   **`Percentage`:** Valor porcentual.
*   **`PartyID`:** Identificador de una Party (del módulo Party).
*   **`ProductVariantID`:** Identificador de una ProductVariant (del módulo Product).
*   **`OrderNumber`:** String (Value Object para números de pedido únicos, encapsula formato y generación).
*   **`QuoteNumber`:** String (Value Object para números de cotización únicos).
*   **`DeliveryNoteNumber`:** String (Value Object para números de albarán únicos).
*   **`InvoiceNumber`:** String (Value Object para números de factura únicos, incluye serie).
*   **`InvoiceType`:** Enum (`COMPLETA`, `SIMPLIFICADA`). Define si es factura completa B2B o ticket (factura simplificada) según legislación española.
*   **`InvoiceSeries`:** Value Object que encapsula la serie de numeración (código, año, prefijo). Permite gestionar series diferenciadas por tipo de documento y año fiscal.

---

### **1. Entidad: `Quote` (Presupuesto / Cotización)**

*   **Agregado Raíz:** `Quote`
*   **Propósito:** Representa una oferta de precios a un cliente que aún no ha sido confirmada.
*   **Atributos:**
    *   `ID`: UUID
    *   `QuoteNumber`: `QuoteNumber` (Value Object)
    *   `PartyID`: `PartyID` (Cliente)
    *   `QuoteDate`: `DateTime` (Fecha de emisión de la cotización)
    *   `ExpirationDate`: `DateTime` (Fecha hasta la cual la cotización es válida)
    *   `Status`: Enum (`BORRADOR`, `ENVIADA`, `APROBADA`, `RECHAZADA`, `EXPIRADA`)
    *   `LineItems`: List<`QuoteLineItem`> (Colección de ítems de línea, parte del agregado)
    *   `Subtotal`: `Money` (Total de los ítems antes de impuestos)
    *   `TaxAmount`: `Money` (Monto de impuestos calculado)
    *   `Total`: `Money` (Subtotal + Impuestos)
    *   `Notes`: String (Notas internas para la cotización)
*   **Comportamiento Clave:**
    *   `ConvertirAOrden(deliveryDate)`: Crea un `SalesOrder` a partir de la `Quote` (si `Status` es `APROBADA`).
    *   `CalcularTotales()`: Recalcula `Subtotal`, `TaxAmount`, `Total` en base a los `LineItems`.

#### **1.1. Entidad: `QuoteLineItem`**

*   **Propósito:** Detalle de un producto en una cotización, con soporte para precios y descuentos manuales.
*   **Atributos:**
    *   `ID`: UUID
    *   `ProductVariantID`: `ProductVariantID`
    *   `Quantity`: Integer
    *   `CalculatedUnitPrice`: `Money` (Precio unitario sugerido por el módulo de Pricing)
    *   `ManualUnitPrice`: `Money` (opcional, si se ha ajustado manualmente)
    *   `FinalUnitPrice`: `Money` (El precio unitario que realmente se aplica: `ManualUnitPrice` si existe, o `CalculatedUnitPrice`).
    *   `CalculatedDiscountPerUnit`: `Money` (opcional, descuento por unidad sugerido por Pricing).
    *   `ManualDiscountPerUnit`: `Money` (opcional, si se ha ajustado manualmente).
    *   `FinalDiscountPerUnit`: `Money` (El descuento por unidad que realmente se aplica).
    *   `Subtotal`: `Money` (Calculado como `Quantity` * (`FinalUnitPrice` - `FinalDiscountPerUnit`)).

---

### **2. Entidad: `SalesOrder` (Pedido de Venta)**

*   **Agregado Raíz:** `SalesOrder`
*   **Propósito:** Representa un compromiso de venta firme y la base para la ejecución de la venta (entrega, facturación).
*   **Atributos:**
    *   `ID`: UUID
    *   `OrderNumber`: `OrderNumber` (Value Object)
    *   `QuoteID`: UUID (opcional, FK a `Quote` si proviene de una cotización)
    *   `PartyID`: `PartyID` (Cliente)
    *   `OrderDate`: `DateTime` (Fecha de confirmación del pedido)
    *   `DeliveryDate`: `DateTime` (Fecha de entrega acordada, requerida)
    *   `Status`: Enum (`PENDIENTE`, `EN_PREPARACION`, `ENTREGADO_PARCIALMENTE`, `ENTREGADO`, `CANCELADO`, `FACTURADO_PARCIALMENTE`, `FACTURADO_COMPLETAMENTE`)
    *   `LineItems`: List<`OrderLineItem`> (Colección de ítems de línea, parte del agregado)
    *   `Subtotal`: `Money` (Suma de los `Subtotal`s de los `LineItems`)
    *   `TaxAmount`: `Money` (Monto de impuestos calculado por `Sales`)
    *   `Total`: `Money` (`Subtotal` + `TaxAmount`)
    *   `Notes`: String (Notas internas del pedido)
*   **Comportamiento Clave:**
    *   `ActualizarEstado(newStatus)`: Cambia el estado del pedido, validando las transiciones.
    *   `CalcularTotales()`: Recalcula `Subtotal`, `TaxAmount`, `Total` en base a los `LineItems`.

#### **2.1. Entidad: `OrderLineItem`**

*   **Propósito:** Detalle de un producto en un pedido, con soporte para precios y descuentos manuales.
*   **Atributos:**
    *   `ID`: UUID
    *   `ProductVariantID`: `ProductVariantID`
    *   `Quantity`: Integer
    *   `CalculatedUnitPrice`: `Money`
    *   `ManualUnitPrice`: `Money` (opcional)
    *   `FinalUnitPrice`: `Money`
    *   `CalculatedDiscountPerUnit`: `Money` (opcional)
    *   `ManualDiscountPerUnit`: `Money` (opcional)
    *   `FinalDiscountPerUnit`: `Money`
    *   `Subtotal`: `Money` (Calculado como `Quantity` * (`FinalUnitPrice` - `FinalDiscountPerUnit`)).

---

### **3. Entidad: `DeliveryNote` (Albarán / Nota de Entrega)**

*   **Agregado Raíz:** `DeliveryNote`
*   **Propósito:** Documentar la entrega física de mercancía al cliente.
*   **Atributos:**
    *   `ID`: UUID
    *   `DeliveryNoteNumber`: `DeliveryNoteNumber` (Value Object)
    *   `SalesOrderID`: UUID (FK a `SalesOrder`)
    *   `PartyID`: `PartyID`
    *   `DeliveryDate`: `DateTime` (Fecha real de la entrega)
    *   `Status`: Enum (`PENDIENTE`, `ENTREGADO`, `CANCELADO`)
    *   `LineItems`: List<`DeliveryNoteLineItem`> (Colección de ítems, parte del agregado)
    *   `Notes`: String
*   **Comportamiento Clave:**
    *   `RegistrarEntrega()`: Confirma la entrega y puede actualizar el `SalesOrder` asociado.

#### **3.1. Entidad: `DeliveryNoteLineItem`**

*   **Propósito:** Detalle de un producto en un albarán. Contiene la información de lo que se entrega, no del precio.
*   **Atributos:**
    *   `ID`: UUID
    *   `SalesOrderLineItemID`: UUID (FK a `OrderLineItem` del pedido original)
    *   `ProductVariantID`: `ProductVariantID`
    *   `DeliveredQuantity`: Integer

---

### **4. Entidad: `Invoice` (Factura / Ticket)**

*   **Agregado Raíz:** `Invoice`
*   **Propósito:** Documento legal para solicitar el pago al cliente y registrar el aspecto financiero de la venta. Soporta facturas completas (B2B) y facturas simplificadas/tickets (retail) según legislación española (ADR-020).
*   **Atributos:**
    *   `ID`: UUID
    *   `InvoiceNumber`: `InvoiceNumber` (Value Object, incluye serie)
    *   `Type`: `InvoiceType` (`COMPLETA` | `SIMPLIFICADA`)
    *   `Series`: `InvoiceSeries` (Value Object, gestiona código de serie, año, prefijo)
    *   `PartyID`: `PartyID` (Cliente. Para tipo `SIMPLIFICADA` puede ser "CONSUMIDOR_FINAL")
    *   `InvoiceDate`: `DateTime` (Fecha de emisión de la factura)
    *   `DueDate`: `DateTime` (Fecha de vencimiento del pago)
    *   `Status`: Enum (`BORRADOR`, `EMITIDA`, `PAGADA`, `VENCIDA`, `ANULADA`)
    *   `LineItems`: List<`InvoiceLineItem`> (Colección de ítems de línea, parte del agregado)
    *   `RelatedOrderIDs`: List<UUID> (opcional, ID de `SalesOrder` o `DeliveryNote` facturado)
    *   `Subtotal`: `Money`
    *   `TaxAmount`: `Money`
    *   `Total`: `Money`
    *   `PaymentTerms`: String (Condiciones de pago).
*   **Comportamiento Clave:**
    *   `MarcarComoPagada()`: Actualiza el estado de la factura.
    *   `ValidarLimitesLegales()`: Para tipo `SIMPLIFICADA`, valida que `Total` < 3.000 EUR según normativa española.

#### **4.1. Entidad: `InvoiceLineItem`**

*   **Propósito:** Detalle de un producto en una factura, reflejando el precio final para fines contables.
*   **Atributos:**
    *   `ID`: UUID
    *   `SalesOrderLineItemID`: UUID (FK a `OrderLineItem` del pedido original)
    *   `ProductVariantID`: `ProductVariantID`
    *   `Quantity`: Integer
    *   `UnitPrice`: `Money` (El `FinalUnitPrice` del `OrderLineItem`)
    *   `DiscountAmount`: `Money` (El `FinalDiscountPerUnit` * `Quantity` del `OrderLineItem`)
    *   `Subtotal`: `Money` (`Quantity` * `UnitPrice` - `DiscountAmount`)
    *   `TaxAmount`: `Money` (Impuesto aplicado a esta línea)
    *   `Total`: `Money` (`Subtotal` + `TaxAmount` de la línea)
