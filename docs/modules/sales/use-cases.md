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

#### **4. Casos de Uso de `Invoice` (Factura / Ticket)**

*   **CU-S-018: CreateInvoiceFromOrder** (Factura Completa)
    *   **Propósito:** Generar una factura completa (tipo B2B) a partir de un `SalesOrder` (o un conjunto de `SalesOrder`s completados).
    *   **Actores:** Vendedor, Personal de Contabilidad.
    *   **Entradas:** `SalesOrderID` (o lista de `SalesOrderID`s), `InvoiceDate`, `DueDate`, `Series` (opcional, default: serie "A").
    *   **Flujo:**
        1.  Recuperar `SalesOrder`(s) y sus `OrderLineItem`s.
        2.  Validar que `Party` tiene datos fiscales completos (NIF, dirección fiscal).
        3.  Crear `InvoiceLineItem`s a partir de los `OrderLineItem`s, utilizando sus `FinalUnitPrice` y `FinalDiscountPerUnit`.
        4.  Calcular `Subtotal`, `TaxAmount`, `Total` de la factura.
        5.  Obtener siguiente número de la serie especificada.
        6.  Crear `Invoice` con `Type = COMPLETA` y `Series` especificada.
        7.  Persistir `Invoice` y sus `InvoiceLineItem`s.
        8.  Actualizar `SalesOrder.Status` a `FACTURADO_PARCIALMENTE` o `FACTURADO_COMPLETAMENTE`.
    *   **Salida:** `Invoice` (tipo COMPLETA).

*   **CU-S-019: CreateSimplifiedInvoice** (Ticket / Factura Simplificada) **(NUEVO - ADR-020)**
    *   **Propósito:** Generar una factura simplificada (ticket) a partir de un `SalesOrder` para ventas retail.
    *   **Actores:** Vendedor, Cajero TPV.
    *   **Entradas:** `SalesOrderID`, `InvoiceDate`, `DueDate` (opcional), `Series` (opcional, default: serie "TKT").
    *   **Precondiciones:**
        -   `SalesOrder.Total` < 3.000 EUR (límite legal español)
        -   `PartyID` puede ser "CONSUMIDOR_FINAL" (Party genérico)
    *   **Flujo:**
        1.  Recuperar `SalesOrder` y validar que `Total < 3.000 EUR`.
        2.  Si `Total >= 3.000 EUR`, retornar error: "Debe emitirse factura completa para importes >= 3.000 EUR".
        3.  Crear `InvoiceLineItem`s a partir de los `OrderLineItem`s.
        4.  Calcular `Subtotal`, `TaxAmount`, `Total`.
        5.  Obtener siguiente número de la serie "TKT" (o serie especificada).
        6.  Crear `Invoice` con:
            -   `Type = SIMPLIFICADA`
            -   `Series = "TKT"` (o especificada)
            -   `PartyID` = del `SalesOrder` (puede ser "CONSUMIDOR_FINAL")
            -   `DueDate` = `InvoiceDate` (pago inmediato en retail)
        7.  Persistir `Invoice` y sus `InvoiceLineItem`s.
        8.  Actualizar `SalesOrder.Status` a `FACTURADO_COMPLETAMENTE`.
    *   **Salida:** `Invoice` (tipo SIMPLIFICADA).
    *   **Postcondiciones:** Ticket listo para impresión/envío al cliente.

*   **CU-S-020: GetInvoice**
        5.  Persistir `Invoice` y `InvoiceLineItem`s.
        6.  Actualizar el `SalesOrder.Status` a `FACTURADO_PARCIALMENTE` o `FACTURADO_COMPLETAMENTE`.
    *   **Salida:** `Invoice`.

*   **CU-S-020: GetInvoice**
    *   **Propósito:** Recuperar los detalles de una factura específica.
    *   **Actores:** Vendedor, Personal de Contabilidad.
    *   **Entradas:** `InvoiceID`.
    *   **Salida:** `Invoice`.

*   **CU-S-021: ListInvoices**
    *   **Propósito:** Listar facturas con opciones de filtrado (ej. por `PartyID`, `Status`, `InvoiceDateRange`, `DueDateRange`, `Type`).
    *   **Actores:** Vendedor, Personal de Contabilidad.
    *   **Entradas:** Filtros (incluye filtro por `InvoiceType`), paginación.
    *   **Salida:** Lista paginada de `Invoice`s.

*   **CU-S-022: ConfigureInvoiceSeries** **(NUEVO - ADR-020)**
    *   **Propósito:** Configurar o actualizar series de numeración de facturas por tipo y año.
    *   **Actores:** Administrador, Personal de Contabilidad.
    *   **Entradas:** `SeriesCode`, `Year`, `InvoiceType`, `Prefix` (opcional), `CurrentNumber` (opcional).
    *   **Flujo:**
        1.  Validar que `SeriesCode` es único para el año especificado.
        2.  Si la serie ya existe, actualizar `Prefix` o `CurrentNumber`.
        3.  Si es nueva, crear registro en tabla `invoice_series` con `CurrentNumber = 0`.
        4.  Persistir configuración.
    *   **Salida:** Configuración de serie actualizada.
    *   **Ejemplos:**
        -   Serie "A" para facturas completas B2B
        -   Serie "TKT" para tickets/facturas simplificadas
        -   Serie "B" para facturas rectificativas

---

#### **5. Casos de Uso Transversales**

*   **CU-S-023: SeleccionarVarianteDeProducto** (Componente Reutilizable)
    *   **Propósito:** Permitir la selección de una Variante de Producto (`ProductVariant`) para añadir a una línea de venta (Quote o SalesOrder), soportando dos modalidades: selección interactiva por atributos o búsqueda directa por SKU.
    *   **Actores:** Vendedor, Sistema de Ventas.
    *   **Contexto:** Este caso de uso es invocado por CU-S-001 (CreateQuote), CU-S-007 (CreateOrder), CU-S-012 (AddLineItemToOrder) y otros casos donde se necesite añadir un producto a una transacción de venta.
    *   **Modalidades de Selección:**
    
        **Modalidad A: Selección Interactiva por Atributos (Con JIT)**
        
        *   **Flujo:**
            1.  El vendedor selecciona un **Producto Base** (ej: "Zapatilla Nike Air Max") de un catálogo o mediante búsqueda.
            2.  El sistema invoca **UC-P-005** (GetApplicableAttributesForProduct) para obtener los atributos configurables aplicables al producto seleccionado.
                - Endpoint: `GET /api/products/{productId}/calculated-option-sets`
            3.  El sistema presenta al vendedor **dropdowns/selectores** para cada atributo aplicable:
                - Ejemplo: 
                  - `Talla`: [Small, Medium, Large, XL]
                  - `Color`: [Rojo, Azul, Verde, Negro]
            4.  El vendedor selecciona un valor para cada atributo requerido.
            5.  Al confirmar la selección, el sistema invoca **UC-P-009** (FindOrCreateProductVariant) con la combinación elegida:
                - Endpoint: `POST /api/products/{productId}/variants/find-or-create`
                - Payload: 
                  ```json
                  {
                    "optionConfiguration": {
                      "SIZE": "Large",
                      "COLOR": "Rojo"
                    }
                  }
                  ```
            6.  El backend ejecuta la lógica JIT (Just-in-Time):
                - **Si la variante existe:** Retorna la variante existente.
                - **Si NO existe:** 
                  - Valida que la combinación sea válida para el producto.
                  - Construye el SKU determinista (ej: `AIR-MAX-2024-SIZE.L-COLOR.R`).
                  - Crea una nueva `ProductVariant` con `Status = PROVISIONAL` y `IsActive = TRUE`.
                  - Persiste en la base de datos.
                  - Retorna la nueva variante.
            7.  El sistema añade el `ProductVariantID` obtenido a la línea de venta con la cantidad especificada.
        
        *   **Precondiciones:**
            - El producto debe tener al menos un atributo aplicable configurado.
            - El vendedor debe tener rol autorizado (`commercial`, `admin`) para disparar la creación JIT.
        
        *   **Postcondiciones:**
            - Se ha obtenido o creado una `ProductVariant` válida y persistida.
            - La variante está disponible para ser añadida a la línea de venta.
            - Si fue creada mediante JIT, su estado es `PROVISIONAL` hasta que sea confirmada manualmente o mediante un proceso de validación.
        
        **Modalidad B: Búsqueda Directa por SKU**
        
        *   **Flujo:**
            1.  El vendedor ingresa o escanea un **SKU de variante** (ej: `AIR-MAX-2024-SIZE.L-COLOR.R`) en un campo de búsqueda.
            2.  El sistema busca la variante por SKU:
                - Endpoint: `GET /api/variants?sku={sku}`
            3.  **Si se encuentra la variante:**
                - El sistema valida que `IsActive = TRUE`.
                - Añade el `ProductVariantID` a la línea de venta.
            4.  **Si NO se encuentra:**
                - El sistema muestra un mensaje de error: "Variante no encontrada o SKU inválido".
                - **NO se permite creación JIT desde esta modalidad** (solo búsqueda de variantes existentes).
        
        *   **Precondiciones:**
            - El SKU ingresado debe ser válido y estar bien formado.
        
        *   **Postcondiciones:**
            - Se ha obtenido una `ProductVariant` existente y activa para añadir a la línea.
    
    *   **Reglas de Negocio:**
        - **Permisos JIT:** Solo usuarios con roles `commercial` o `admin` pueden disparar la creación JIT de variantes (Modalidad A).
        - **Estado de Variantes Creadas JIT:** Todas las variantes creadas mediante JIT comienzan con `Status = PROVISIONAL`.
        - **Validación de Combinaciones:** El sistema debe garantizar que la combinación de atributos sea válida según los atributos aplicables al producto.
        - **SKU Determinista:** El SKU de una variante se construye automáticamente siguiendo la fórmula: `{Product.SKU}-{Attr1.Code}.{Val1.Code}-{Attr2.Code}.{Val2.Code}...` con atributos ordenados por `SortOrder`.
    
    *   **Entradas:** 
        - Modalidad A: `ProductID`, selección de valores de atributos.
        - Modalidad B: `SKU` de la variante.
    
    *   **Salidas:** 
        - `ProductVariantID` válido para añadir a la línea de venta.
        - Información adicional: `SKU`, `Name`, `BaseCost` (para cálculo de precios).
    
    *   **Referencia Cruzada:**
        - Este caso de uso consume: **UC-P-005** (GetApplicableAttributesForProduct) y **UC-P-009** (FindOrCreateProductVariant) del módulo Product.
        - Ver documentación completa en: `docs/modules/product/use-cases.md`