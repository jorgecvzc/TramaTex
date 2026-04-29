# Estrategia Técnica: Gestión de Cobros, Vencimientos y Tesorería (Post-MVP)

Este documento define la arquitectura para el control del ciclo financiero posterior a la facturación. El objetivo es transformar TramaTex en una herramienta capaz de gestionar no solo la venta, sino también la liquidez y el riesgo de impago.

---

## 1. Modelado Financiero y Reglas de Negocio

### 1.1 Maestro de Condiciones de Pago (Payment Terms)
Se implementará una entidad configurable para definir cómo se espera recibir el cobro:
- **Días de Vencimiento**: Lista de días (ej: `[0]` para contado, `[30, 60]` para dos plazos).
- **Día de Pago Fijo**: Capacidad de ajustar el vencimiento al día X del mes (común en grandes empresas).
- **Método por Defecto**: Vinculación opcional a una Forma de Pago (Transferencia, Recibo, Efectivo).

### 1.2 Cálculo Automático de Vencimientos
Al emitir una factura, el sistema generará automáticamente N registros en la tabla de `Vencimientos`:
- **Fórmula**: `Fecha Factura + Días Plazo -> Ajuste a Día de Pago -> Vencimiento Final`.
- **Importe**: División equitativa del total de la factura entre los plazos definidos.

---

## 2. Ciclo de Vida del Cobro y Estados

### 2.1 Estados del Vencimiento/Factura
- **PENDIENTE**: Factura emitida, vencimiento futuro.
- **VENCIDO**: Fecha actual > Fecha de vencimiento y saldo > 0. (Activación de alertas visuales).
- **COBRADO_PARCIAL**: Registro de pagos que no cubren el total del plazo.
- **COBRADO**: Saldo de la factura es 0.
- **DEVUELTO**: Gestión de impagos (especialmente para recibos bancarios).

### 2.2 Registro de Cobros (Collections)
Interfaz para imputar pagos a facturas:
- **Conciliación Rápida**: Seleccionar factura y marcar como "Cobrada hoy" con un solo clic/tecla.
- **Pagos Multifactura**: Un único ingreso bancario puede cubrir varias facturas de un cliente.

---

## 3. UX de Tesorería y Ergonomía

Siguiendo la filosofía "Keyboard-First":

### 3.1 Dashboard de Deuda (Aging Report)
- **KPIs Críticos**: "Deuda Total", "Vencido < 30 días", "Vencido > 30 días".
- **Visualización**: Colores semafóricos (Verde: Al corriente, Naranja: Próximo vencimiento, Rojo: Vencido).

### 3.2 Acciones de Alta Velocidad
- **Búsqueda por Importe**: `Ctrl + K` -> introducir importe exacto para localizar facturas pendientes de cobrar.
- **Atajo `Alt + P`**: En el detalle de factura, abrir modal de "Registrar Cobro" con el importe restante pre-rellenado.

---

## 4. Especificaciones Técnicas (Backend)

### 4.1 Nuevas Entidades de Dominio
- `PaymentTerm`: Configuración de plazos.
- `InvoiceDue`: El registro individual de cada vencimiento vinculado a una factura.
- `Collection`: El registro del flujo de caja (dinero recibido).

### 4.2 Lógica de Integridad
- **Cierre de Factura**: Una factura cambia automáticamente a estado `PAGADA` cuando la suma de sus `Collections` iguala o supera el `TotalAmount`.
- **Auditoría**: Cada registro de cobro debe guardar el `UserID` que realizó la operación y la fecha/hora real.

---

*Última actualización: 2026-04-27*
