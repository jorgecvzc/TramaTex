# Estrategia Técnica: Facturación Consolidada Multi-Albarán (Post-MVP)

Este documento define la especificación técnica y funcional para implementar la facturación agrupada en TramaTex. Esta funcionalidad permite optimizar la gestión administrativa al emitir una única factura para múltiples entregas (albaranes) de un mismo cliente.

---

## 1. Lógica de Dominio y Flujo de Negocio

El sistema debe garantizar la integridad de la cadena de ventas: `Pedido -> Albarán -> Factura`.

### 1.1 El Proceso de Consolidación
1.  **Selección**: El usuario filtra albaranes con estado `ENTREGADO` y `PENDIENTE_DE_FACTURAR` para un cliente específico.
2.  **Validación de Agrupación**:
    - Todos los albaranes deben pertenecer al mismo **Party ID** (Cliente).
    - **Gestión de Moneda**: Deben compartir la misma **Moneda**. En caso de clientes internacionales (Post-MVP), el sistema permitirá seleccionar la moneda de destino y aplicará el tipo de cambio oficial de la fecha de factura, recalculando los importes desde la moneda base del sistema (EUR).
    - Deben tener la misma **Condición de Pago** (opcional, configurable).
3.  **Algoritmo de Fusión de Líneas**:
    - Las líneas de los diferentes albaranes se agrupan en la factura si coinciden:
        - `ProductVariantID`
        - `UnitPrice`
        - `DiscountPercentage`
        - `TaxRate`
    - La cantidad final en la línea de factura será la `SUM(quantity)` de las líneas de albarán vinculadas.
4.  **Criterios de Ruptura (Agrupación Dinámica)**:
    - El sistema debe permitir al usuario definir "dimensiones de ruptura" al procesar múltiples albaranes:
        - **Por Mes Natural**: Generar una factura distinta para cada mes de entrega.
        - **Por Proyecto/Obra**: Si los albaranes están vinculados a proyectos distintos, generar facturas separadas.
        - **Por Dirección de Envío**: Útil si el cliente requiere facturación separada por delegación.
        - **Por Pedido de Origen**: Mantener una factura por cada pedido original del cliente.
5.  **Vinculación N:1**: Cada línea de albarán original debe guardar la referencia a la `InvoiceLineItemID` generada para mantener la trazabilidad completa.

---

## 2. Experiencia de Usuario (UI/UX)

Siguiendo el Plan Maestro de UI/UX, la interfaz debe ser eficiente y clara.

### 2.1 Nueva Pantalla: "Generar Factura Consolidada"
- **Filtros Inteligentes**: Búsqueda por Cliente, Rango de Fechas y Proyecto.
- **Selector Masivo**: Tabla con checkboxes para seleccionar los albaranes a incluir.
- **Configurador de Ruptura**: Selector para decidir si se agrupa todo en una factura o se rompe por Mes/Proyecto/etc.
- **Panel de Previsualización**: Antes de confirmar, mostrar una lista de las facturas que se van a generar (ej: "Se generarán 3 facturas: Enero, Febrero y Proyecto X").
- **Acciones Rápidas**:
    - `[Alt + F]`: Facturar seleccionados.
    - `[Alt + C]`: Cancelar selección.

### 2.2 Indicadores de Estado
- En el listado de albaranes, se añadirá un badge dinámico:
    - `Facturado`: Vinculado a una factura.
    - `Pendiente`: Listo para ser consolidado.
    - `Parcial`: (Caso borde) Algunas líneas facturadas, otras no.

---

## 3. Especificaciones Técnicas (Backend)

### 3.1 Cambios en el Modelo de Datos
La estructura actual es compatible, pero se requieren ajustes de estado:
- **Tabla `delivery_notes`**: Asegurar que el estado pase a `FACTURADO` solo cuando todas sus líneas tengan un `invoice_line_item_id`.
- **Nuevo Endpoint**: `POST /api/v1/sales/invoices/consolidate`
    - Payload: `{ "party_id": "UUID", "delivery_note_ids": ["UUID", ...] }`

### 3.2 Integridad Referencial
- **Protección contra Doble Facturación**: El endpoint de consolidación debe ejecutarse dentro de una transacción de base de datos (`DB Transaction`). Debe verificar que ningún albarán del ID listado tenga ya líneas vinculadas a una factura antes de proceder.

---

## 4. Plan de Implementación (Pasos Críticos)

1.  **Servicio de Agrupación**: Desarrollar la lógica en `sales-dom` para iterar y colapsar líneas.
2.  **Refactor de Facturación**: Adaptar el comando `CreateInvoice` para que acepte múltiples orígenes.
3.  **Interfaz de Selección**: Crear la página en el frontend usando `BaseCatalog` con soporte multiselect.
4.  **Validación de Trazabilidad**: Test de integración que verifique que desde la Factura se puede navegar a todos los Albaranes origen y viceversa.

---

*Última actualización: 2026-04-27*
