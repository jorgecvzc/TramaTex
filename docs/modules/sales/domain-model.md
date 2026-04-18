# Modelo de Dominio - Módulo Sales

Este documento describe el flujo documental de ventas en TramaTex, desde la prospección comercial hasta la fiscalidad final, asegurando la trazabilidad y la integridad financiera.

---

## 1. El Ciclo de Vida Documental

La venta en TramaTex no es un evento aislado, sino una transición de estados de confianza y compromiso legal:

### El Presupuesto (`Quote`)
Representa una intención comercial no vinculante.
- **Validez Temporal:** Todo presupuesto nace con una fecha de expiración.
- **Conversión:** Una vez aceptado, el presupuesto "muere" para dar vida a un Pedido (`SalesOrder`), heredando todas sus condiciones de precio para mantener la palabra dada al cliente.

### El Pedido (`SalesOrder`)
Es el eje central de la operación.
- **Compromiso de Entrega:** Requiere obligatoriamente una fecha de entrega (`DeliveryDate`).
- **Disparador de Producción:** Un pedido nace directamente en estado **Confirmado** (`EN_PREPARACION`), lo que dispara automáticamente su visibilidad en el taller (módulo **MES**) si los productos son de fabricación propia.
- **Inmutabilidad Parcial:** Una vez que un pedido tiene albaranes o facturas asociadas, sus líneas de producto quedan bloqueadas para edición.

### El Albarán (`DeliveryNote`)
Documenta la salida física de mercancía.
- **Control de Saldos:** Un pedido puede generar múltiples albaranes (entregas parciales) hasta completar la cantidad total solicitada.

### La Factura y el Ticket (`Invoice`)
El documento de cierre legal y fiscal.
- **Dualidad de Emisión:** Soporta Facturas B2B (Completas) y Tickets (Simplificadas) para venta rápida.
- **Tickets — Venta Directa TPV:** Los tickets (facturas simplificadas) se crean directamente desde la interfaz TPV/POS sin necesidad de pedido ni albarán previo. El cajero selecciona productos y cantidades, opcionalmente elige un cliente (por defecto CONSUMIDOR FINAL), y el sistema genera la factura con serie "FT". El precio se calcula usando `BaseSalesPrice` del motor de precios; el descuento del cliente aparece separado y es editable manualmente.
- **Regla de los 3.000€:** Por normativa española, el sistema impide generar Facturas Simplificadas que superen este importe, obligando a la identificación completa del cliente.
- **Origen Exclusivo desde Albaranes (Facturas Completas, MVP):** Las facturas completas (B2B) se generan exclusivamente desde albaranes entregados. Cada línea de albarán se vincula a su línea de factura correspondiente (`invoice_line_item_id`), garantizando trazabilidad completa. En el MVP se crea una factura por albarán (relación 1:1). La infraestructura soporta relaciones N:1 para permitir en Post-MVP la consolidación de múltiples albaranes de un mismo cliente en una sola factura (ej: facturación mensual).

---

## 2. Reglas de Negocio y Comportamiento

### Estrategia de "Manual Override"
Aunque el módulo de **Pricing** sugiere precios y descuentos, el módulo de **Sales** es soberano. 
- Los agentes comerciales pueden sobreescribir manualmente el precio unitario o el descuento por línea. 
- **Trazabilidad:** El sistema almacena tanto el precio calculado como el manual para permitir análisis de rentabilidad posteriores.

### Normalización de Estados
Para facilitar la comunicación entre el backend (con estados técnicos en castellano para legibilidad de dominio) y el frontend (estándares internacionales), el sistema utiliza la siguiente normalización:

**Presupuestos (Quotes):**
- `BORRADOR` (Draft)
- `EMITIDA` (Issued)
- `APROBADA` (Approved)
- `RECHAZADA` (Rejected)
- `EXPIRADA` (Expired)
- `CONVERTIDA_A_PEDIDO` (Converted)

**Pedidos (Orders):**
- `EN_PREPARACION` (Confirmed / In Preparation) - **Estado Inicial**
- `ENTREGADO_PARCIALMENTE` (Partially Delivered)
- `ENTREGADO` (Delivered)
- `FACTURADO_PARCIALMENTE` (Partially Invoiced)
- `FACTURADO_COMPLETAMENTE` (Invoiced)
- `CANCELADO` (Cancelled)
- `PENDIENTE` (Pending) - Estado de seguridad tras una reactivación desde anulado.

**Facturas (Invoices):**
- `BORRADOR` (Draft)
- `EMITIDA` (Issued)
- `PAGADA` (Paid)
- `VENCIDA` (Overdue)
- `ANULADA` (Void)

### Integridad Referencial y Soberanía Comercial
- **Congelación de Precios:** Una vez que un documento sale del estado `BORRADOR`, el **Precio Unitario**, el **Descuento** y el **IVA** se congelan para garantizar la seguridad jurídica.
- **Hidratación Dinámica de Nombres:** El nombre del producto y su SKU se obtienen dinámicamente del módulo **Product** cada vez que se visualiza el documento. Esto asegura que cualquier refinamiento en la descripción del catálogo se refleje en los documentos vivos, priorizando la coherencia del catálogo sobre el histórico nominal.

---

## 3. Integración con MES — Trabajos Asociados (`SalesWorkSetup`)

Los presupuestos y pedidos pueden llevar asociados uno o varios **trabajos de personalización/modificación** que describen las operaciones a realizar sobre las prendas en el taller. Estos trabajos se modelan como `SalesWorkSetup`, una entidad propia del módulo Sales que actúa como puente hacia el módulo MES.

### La Entidad `SalesWorkSetup`
Cada trabajo asociado a un documento de venta contiene:
- **`ID`:** Identificador único.
- **`WorkSetupID`** (`*uuid.UUID`): Referencia a un `WorkSetup` en MES. Se auto-crea si el comercial no selecciona uno existente.
- **`WorkOrderID`** (`*uuid.UUID`): Referencia a la `WorkOrder` generada en MES para la ejecución real. Se establece automáticamente cuando el taller crea la orden: `MESService.CreateWorkOrder()` invoca `SalesOrderLinker.LinkWorkOrder(ctx, orderWorkSetupID, workOrderID)` que actualiza el registro en `order_work_setups`.
- **`Name`:** Nombre descriptivo del trabajo (ej. "Serigrafía camisetas Confecciones López").
- **`Observations`:** Campo de texto libre donde el comercial describe las características del trabajo, indicaciones especiales, o puntualizaciones sobre trabajos ya definidos. Es el principal canal de comunicación entre el departamento comercial y el taller.
- **`Sequence`:** Orden de los trabajos dentro del documento.

> **Nota:** `SalesWorkSetup` es una entidad **sin estado propio** (no tiene campo `Status`). El progreso del trabajo se determina consultando la presencia de `WorkSetupID`, `WorkOrderID` y el estado de la `WorkOrder` en MES.

### Ciclo de Vida

1. **Creación en Presupuesto o Pedido:** El comercial añade trabajos con nombre + observaciones. Opcionalmente selecciona un `WorkSetup` existente en MES. Si no selecciona uno, el sistema **auto-crea** un `WorkSetup` inactivo en MES (sin líneas) vía `WorkSetupCreatorAdapter`.
2. **Conversión Presupuesto → Pedido:** Al convertir, los `SalesWorkSetup` se copian. El método `ensureWorkSetups` garantiza que todos los trabajos copiados tengan un `WorkSetupID` válido (cierra la brecha de presupuestos antiguos sin `WorkSetupID`).
3. **Visibilidad en MES:** El Dashboard de MES consulta los pedidos en estado `EN_PREPARACION` y muestra los `SalesWorkSetup` cuyo `WorkOrderID` es nulo (solicitudes pendientes del taller).
4. **Creación de Orden:** El jefe de taller crea la `WorkOrder` desde el `WorkSetup` asociado. Al crearse, el backend vincula el `WorkOrderID` de vuelta al registro de `order_work_setups` vía `SalesOrderLinkerAdapter`.
5. **Seguimiento:** El estado del trabajo se consulta bajo demanda vía `WorkOrderQueryService` de MES. (Post-MVP: eventos de dominio `WorkOrderStarted`/`WorkOrderCompleted` para notificaciones reactivas).
6. **Suspensión / Reactivación:** Si un pedido es cancelado, `SalesService` llama a `WorkOrderSuspender.SuspendWorkOrders()` con los `WorkOrderID`s asociados. Si el pedido se reactiva (CANCELADO → PENDIENTE), llama a `ReactivateWorkOrders()`. Esta comunicación es unidireccional: Sales → MES.

### Reglas
- Un documento (`Quote` o `SalesOrder`) puede tener cero o más `SalesWorkSetup`.
- Los trabajos se **copian** del presupuesto al pedido durante la conversión (igual que las líneas de producto).
- El campo `Observations` no tiene restricción de longitud y se preserva en todas las operaciones.
- Todo `SalesWorkSetup` tiene `WorkSetupID` garantizado tras la auto-creación; el comercial no necesita seleccionar una plantilla MES manualmente.

### Consulta de Progreso de Trabajos
Sales puede consultar el estado de ejecución de los `WorkOrder`s asociados a un pedido sin conocer la lógica interna de MES. Para ello:
- La capa de aplicación de Sales define una interfaz `MESWorkLookup` con dos métodos: consultar progreso de una orden y consultar progreso de varias órdenes.
- La interfaz devuelve un DTO propio de Sales (`WorkOrderProgress`) con: estado global, total de tareas, tareas completadas, y desglose por línea.
- Toda la lógica de cálculo (qué tareas están hechas, cuáles faltan) reside en el módulo MES (`WorkOrderQueryService`). Sales es completamente transparente a esa lógica.
- El adaptador en infraestructura (`MESWorkLookupAdapter`) traduce entre los DTOs de MES y los DTOs de Sales.

---

## 4. Cálculos y Redondeo

La integridad financiera del módulo Sales reside en una lógica de cálculo centralizada y predecible, compartida por presupuestos, pedidos y facturas.

### Centralización en `SumAmounts`
Todos los cálculos de agregación (subtotales de líneas, totales de impuestos, totales de documento) delegan en la función de dominio `SumAmounts`. Esta función garantiza:
- **Consistencia de Moneda**: Valida que todos los importes sumados pertenezcan a la misma moneda (Euros por defecto en el MVP).
- **Tratamiento de Nulos**: Devuelve un objeto `Money` con valor cero si la lista de importes está vacía, evitando errores de puntero nulo.

### El Método `RecalculateTotals`
Cada Agregado de ventas (`Quote`, `SalesOrder`, `Invoice`) implementa un método `RecalculateTotals()` que orquestra el flujo de cálculo:
1. **Subtotal**: Suma de los subtotales de cada línea (Cantidad × Precio Unitario Neto).
2. **Impuestos**: Suma de los importes de IVA de cada línea. Si una línea no especifica IVA pero el documento tiene un IVA global definido, se usa este último como fallback.
3. **Total**: Suma aritmética de Subtotal + Impuestos.

### Precisión y Redondeo
- **Almacenamiento**: Los importes se manejan internamente mediante el Value Object `Money`, que utiliza `float64` para los cálculos pero asegura el redondeo a **2 decimales** en las operaciones de salida y persistencia.
- **Redondeo en Línea**: El subtotal de cada línea se redondea antes de sumarse al total del documento, minimizando discrepancias por decimales huérfanos en documentos extensos.

---
**Última Actualización:** 25-03-2026
