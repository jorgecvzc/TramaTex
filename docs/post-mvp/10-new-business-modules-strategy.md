# Estrategia Técnica: Nuevos Módulos de Negocio (Post-MVP)

Este documento esboza la arquitectura y los requisitos de los módulos de expansión previstos para completar el ecosistema ERP de TramaTex: Compras, Inventario Avanzado y Logística.

---

## 1. Módulo de Compras (Purchases)

Cierra el ciclo de aprovisionamiento.
- **Flujo**: Solicitud de Compra -> Pedido a Proveedor -> Recepción de Mercancía -> Factura de Proveedor.
- **Integración con E-Factura**: Capacidad de importar XML de proveedores para registrar facturas y deudas automáticamente.
- **Precios de Coste**: Actualización automática del `baseCost` de los productos en base a la última compra (o precio medio ponderado).

---

## 2. Inventario Avanzado y Almacén

Evolución del stock básico hacia una gestión logística.
- **Multialmacén**: Control de existencias por ubicación física.
- **Movimientos de Almacén**: Trazabilidad total de entradas, salidas y transferencias internas.
- **Alertas de Stock Mínimo**: Notificaciones automáticas cuando una materia prima baja de un umbral crítico para la producción.

---

## 3. Logística y Distribución

Gestión de la última milla.
- **Hojas de Ruta**: Organización de repartos por vehículo y zona.
- **Integración con Transportistas**: Conexión vía API con agencias externas para generación de etiquetas y seguimiento (Tracking).
- **Control de Entregas**: App móvil ligera para que los repartidores confirmen la entrega y recojan la firma digital.

---

## 4. Portal de Autoservicio: TramaTex Connect

Plataforma web de acceso externo para socios de negocio.

### 4.1 Área de Clientes
- **Descarga de Facturas**: Acceso al historial legal de facturas electrónicas (obligatorio por Ley Crea y Crece).
- **Seguimiento de Producción**: Visualización en tiempo real del estado de sus pedidos (conectado con el MES).
- **Gestión de Incidencias**: Sistema de tickets para reportar problemas de calidad o logística.

### 4.2 Área de Proveedores
- **Subida de Facturas**: Portal para que los proveedores carguen sus XML directamente al ERP.
- **Consulta de Pagos**: Visibilidad sobre el estado de sus facturas y fechas previstas de cobro.

---

## 4. Arquitectura de Integración

- **Principios Modulares**: Cada nuevo módulo seguirá el patrón de Arquitectura Hexagonal y se comunicará vía eventos (NATS) para no degradar el núcleo del sistema.
- **Extensiones de Dominio**: Los módulos existentes (Ventas, Productos) expondrán servicios de consulta para que los nuevos módulos puedan operar sin acceder a sus bases de datos.

---

*Última actualización: 2026-04-27*
