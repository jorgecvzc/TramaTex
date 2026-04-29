# Estrategia Técnica: Integración Profunda Sales ↔ MES (Post-MVP)

Este documento detalla la sincronización bidireccional y en tiempo real entre los compromisos comerciales (Ventas) y la realidad de la planta (Producción).

---

## 1. Sincronización Automática de Cambios

### 1.1 De Sales a MES (Propagación de Cambios)
- **Modificación de Pedido**: Si un comercial cambia la cantidad o el diseño de un producto en un pedido ya confirmado:
    - El sistema detecta si la OT ya ha iniciado.
    - Si está en espera: se actualiza automáticamente la OT.
    - Si está en curso: se lanza una **Alerta de Interrupción** al terminal del operario y se requiere validación del responsable de planta.

### 1.2 De MES a Sales (Visibilidad de Entrega)
- **Recálculo de Fecha de Entrega**: Basado en el rendimiento real de la planta y la carga de trabajo, el MES devolverá a Sales una "Fecha estimada de finalización" realista.
- **Notificación de Retrasos**: Si una máquina crítica falla, el sistema notifica automáticamente al comercial responsable del pedido para que pueda informar al cliente proactivamente.

---

## 2. Integración de Costes Reales

- **Coste Real de Fabricación**: Al finalizar la OT, el MES comunica al ERP el tiempo exacto consumido y los materiales utilizados.
- **Análisis de Desviaciones**: El ERP compara el `baseCost` estimado (usado para el precio de venta) con el coste real, permitiendo ajustar las tarifas comerciales para futuros pedidos.

---

## 3. UX Unificada

- **Timeline del Pedido**: En la ficha de pedido, el usuario verá una línea de tiempo integrada que mezcla hitos comerciales (Confirmación, Facturación) con hitos de fabricación (Inicio Corte, Fin Confección, Control de Calidad).

---

*Última actualización: 2026-04-27*
