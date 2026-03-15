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
- **Disparador de Producción:** Un pedido confirmado puede disparar automáticamente órdenes de trabajo en el módulo **MES** si los productos son de fabricación propia (ver sección 3: Integración con MES).
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
- `PENDIENTE` (Pending)
- `EN_PREPARACION` (In Preparation)
- `ENTREGADO_PARCIALMENTE` (Partially Delivered)
- `ENTREGADO` (Delivered)
- `FACTURADO_PARCIALMENTE` (Partially Invoiced)
- `FACTURADO_COMPLETAMENTE` (Invoiced)

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
- **`WorkSetupID`** (opcional): Referencia a un `WorkSetup` existente en MES (plantilla de personalización).
- **`WorkOrderID`** (opcional): Referencia a la `WorkOrder` generada en MES para la ejecución real.
- **`Name`:** Nombre descriptivo del trabajo (ej. "Serigrafía camisetas Confecciones López").
- **`Observations`:** Campo de texto libre donde el comercial describe las características del trabajo, indicaciones especiales, o puntualizaciones sobre trabajos ya definidos. Es el principal canal de comunicación entre el departamento comercial y el taller.
- **`Status`:** Estado del trabajo dentro del ciclo de vida comercial.
- **`Sequence`:** Orden de los trabajos dentro del documento.

### Estados del Trabajo en Sales (`SalesWorkStatus`)
| Estado | Castellano | Descripción |
|---|---|---|
| `BORRADOR` | Borrador | Trabajo definido en un presupuesto o pedido pendiente; editable. |
| `PENDIENTE` | Pendiente | Pedido confirmado; el trabajo queda pendiente de ser recogido por el taller. |
| `EN_PROCESO` | En proceso | El taller ha creado una `WorkOrder` y la ejecución está en curso. |
| `COMPLETADO` | Completado | La `WorkOrder` en MES ha finalizado. |
| `CANCELADO` | Cancelado | El trabajo ha sido descartado. |

### Ciclo de Vida

1. **Creación en Presupuesto o Pedido:** El comercial añade trabajos con nombre + observaciones. Opcionalmente selecciona un `WorkSetup` existente en MES. Estado: `BORRADOR`.
2. **Confirmación del Pedido:** Al aprobar/confirmar el pedido, todos los `SalesWorkSetup` en `BORRADOR` pasan automáticamente a `PENDIENTE`.
3. **Recogida por el Taller:** El jefe de taller consulta los trabajos pendientes de Sales. Si el trabajo tiene un `WorkSetupID`, crea la `WorkOrder` desde la plantilla. Si solo tiene observaciones, define el trabajo manualmente. El `SalesWorkSetup` pasa a `EN_PROCESO` y se registra el `WorkOrderID`.
4. **Finalización:** Cuando la `WorkOrder` se completa en MES, el sistema notifica a Sales (evento de dominio `WorkOrderCompleted`) y el `SalesWorkSetup` pasa a `COMPLETADO`.

### Reglas
- Un documento (`Quote` o `SalesOrder`) puede tener cero o más `SalesWorkSetup`.
- Los trabajos se **copian** del presupuesto al pedido durante la conversión (igual que las líneas de producto).
- El campo `Observations` no tiene restricción de longitud y se preserva en todas las transiciones.
- El `WorkSetupID` es opcional: un comercial puede definir un trabajo nuevo desde Sales solo con nombre y observaciones, sin que exista aún la plantilla en MES.

### Consulta de Progreso de Trabajos
Sales puede consultar el estado de ejecución de los `WorkOrder`s asociados a un pedido sin conocer la lógica interna de MES. Para ello:
- La capa de aplicación de Sales define una interfaz `MESWorkLookup` con dos métodos: consultar progreso de una orden y consultar progreso de varias órdenes.
- La interfaz devuelve un DTO propio de Sales (`WorkOrderProgress`) con: estado global, total de tareas, tareas completadas, y desglose por línea.
- Toda la lógica de cálculo (qué tareas están hechas, cuáles faltan) reside en el módulo MES (`WorkOrderQueryService`). Sales es completamente transparente a esa lógica.
- El adaptador en infraestructura (`MESWorkLookupAdapter`) traduce entre los DTOs de MES y los DTOs de Sales.

---
**Última Actualización:** 2026-03-14
