# Estrategia Técnica: Inteligencia de Negocio y Analítica (Post-MVP)

Este documento define el sistema de análisis de datos para TramaTex, transformando la información transaccional en indicadores estratégicos (KPIs) para la dirección y mandos intermedios.

---

## 1. Métricas de Valor Estratégico

El sistema se centrará en cuatro dimensiones críticas:

### 1.1 Rentabilidad y Márgenes
- **Margen por Producto**: Diferencia real entre el `PrecioVenta` y el `baseCost` dinámico obtenido de la fabricación.
- **Rentabilidad por Cliente**: Análisis de volumen vs. margen neto tras descuentos.

### 1.2 Eficiencia de Planta (MES)
- **OEE (Overall Equipment Effectiveness)**: Disponibilidad, Rendimiento y Calidad de las máquinas.
- **Cuellos de Botella**: Identificación de fases de producción con desviaciones de tiempo recurrentes.

### 1.3 Ciclo de Caja (Cash Conversion Cycle)
- **Periodo Medio de Cobro (PMC)**: Tiempo real desde la emisión de factura hasta el cobro efectivo (integrado con el hito 03).
- **Previsión de Tesorería**: Flujo de caja esperado basado en vencimientos futuros.

---

## 2. Arquitectura Analítica

Para no penalizar el rendimiento del ERP en tiempo real (OLTP), se utilizará un enfoque de **Vistas Materializadas**:

- **Vistas Materializadas en PostgreSQL**: Creación de tablas de resumen hidratadas periódicamente (ej: cada hora o cada noche).
- **Procesamiento en Segundo Plano**: Uso de Workers para recalcular KPIs complejos fuera de la ruta crítica de las peticiones de usuario.

---

## 3. UX Analítica y Dashboards

### 3.1 Interactividad (Drill-down)
- Los Dashboards permitirán hacer clic en un KPI (ej: "Márgenes Bajos") para navegar directamente al listado de facturas o productos que están causando esa métrica.

### 3.2 Reporting Programado
- **Generador de Informes**: Capacidad de programar el envío de un PDF resumen (ej: "Estado de Ventas Semanal") al correo de los responsables de forma automática.
- **Exportación Directa**: Todos los widgets analíticos tendrán botón de descarga a Excel/CSV.

---

*Última actualización: 2026-04-27*
