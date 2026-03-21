# Hoja de Ruta Post-MVP — TramaTex

> Documento consolidado de todas las mejoras, funcionalidades y módulos planificados para después del MVP.
> Cada ítem incluye contexto, prioridad estimada y referencias a la documentación existente.

---

## Índice

1. [Integración Sales ↔ MES (Producto-Trabajo)](#1-integración-sales--mes-producto-trabajo)
2. [Extracción MES como Microservicio](#2-extracción-mes-como-microservicio)
3. [Mejoras de UI y Experiencia de Usuario](#3-mejoras-de-ui-y-experiencia-de-usuario)
4. [Caché de Productos/Variantes y Precios](#4-caché-de-productosvariantes-y-precios)
5. [Comunicación Asíncrona (Message Broker)](#5-comunicación-asíncrona-message-broker)
6. [Inteligencia de Negocio y Analítica](#6-inteligencia-de-negocio-y-analítica)
7. [Notificaciones en Tiempo Real](#7-notificaciones-en-tiempo-real)
8. [Búsqueda Global Avanzada](#8-búsqueda-global-avanzada)
9. [Nuevos Módulos de Negocio](#9-nuevos-módulos-de-negocio)
10. [Mejoras Técnicas de Infraestructura](#10-mejoras-técnicas-de-infraestructura)
11. [Cobertura de Tests](#11-cobertura-de-tests)
12. [Generación de Documentos PDF](#12-generación-de-documentos-pdf)

---

## 1. Integración Sales ↔ MES (Producto-Trabajo)

**Prioridad:** Alta  
**Contexto:**

En el MVP, la relación entre los documentos de Sales (presupuestos, pedidos) y los trabajos MES se implementa a **nivel de documento** de forma informativa:

- Un documento de ventas puede tener **varios trabajos MES asociados** (relación N:M a nivel de documento).
- La columna `mes_work_id` en las **líneas de producto se mantiene en BD** (nullable) pero **no se usa en la UI del MVP**.
- Los usuarios seleccionan trabajos MES a nivel de documento como **guía orientativa**. Es labor del usuario identificar qué MES corresponde a qué productos.
- La asignación granular de MES por línea de producto **no está expuesta** en la interfaz del MVP.

**Objetivo Post-MVP:**

Implementar la vinculación completa y funcional entre líneas de producto y trabajos MES:

- [ ] Selector de MES **por línea de producto** en QuoteCreate/OrderCreate (ya implementado en BD, falta activar en UI).
- [ ] Validación: al aceptar un pedido, verificar que todas las líneas tengan MES asignado (o política configurable).
- [ ] Auto-asignación inteligente de MES basada en el tipo de producto o reglas de negocio.
- [ ] Generación automática de `WorkExecution` por línea al transicionar un pedido a `ACCEPTED`.
- [ ] Trazabilidad bidireccional completa: desde un trabajo MES ver qué líneas de pedido lo originaron y viceversa.

**Decisión MVP (ADR implícito):**

> Se decidió que para el MVP el MES no se asigna a líneas de producto individuales sino al documento completo.
> Un documento puede tener varios MES asociados. Esto sirve de guía de unión de trabajos pero es labor de
> los usuarios identificar qué MES va con qué productos. Se deja para Post-MVP una unión más funcional
> y completa entre Producto y MES.

**Refs:** `migrations/005_init_sales.sql`, `domain/quote.go`, `domain/sales_order.go`

---

## 2. Extracción MES como Microservicio

**Prioridad:** Media  
**Referencia:** [ADR-022](adr-022-mes-microservice-extraction-strategy.md)

Extraer el módulo MES del monolito modular a un microservicio independiente:

- [ ] Comunicación síncrona vía **gRPC** (consultas de productos, validaciones).
- [ ] Comunicación asíncrona vía **NATS JetStream** (comandos, efectos secundarios).
- [ ] Patrón **Transactional Outbox** para garantizar entrega de mensajes.
- [ ] Base de datos PostgreSQL independiente para MES.
- [ ] Propagación de JWT en metadatos gRPC.
- [ ] Proyecciones locales (vistas materializadas) para lectura rápida.
- [ ] Pull Consumers con suscripciones duraderas para escalado horizontal.

---

## 3. Mejoras de UI y Experiencia de Usuario

**Prioridad:** Alta  
**Referencia:** [ADR-022 (notas)](adr-022-mes-microservice-extraction-strategy.md#L3-L5)

- [ ] **Aplicación dirigida por teclado:** Atajos, navegación sin ratón, flujo rápido para operadores.
- [ ] **Mejora interfaz TPV:** Mayor agilidad para ventas rápidas (tickets/facturas simplificadas).
- [ ] **Diseño responsive:** Adaptación completa a tablets de taller.
- [ ] **Modo oscuro** para terminales de producción.

---

## 4. Caché de Productos/Variantes y Precios

**Prioridad:** Media  
**Referencia:** [ADR-022 (notas)](adr-022-mes-microservice-extraction-strategy.md#L4)

- [ ] Implementar caché (Redis o en memoria) de productos/variantes y precios para consultas rápidas.
- [ ] Invalidación automática de caché al actualizar un precio base (borrar todas las entradas derivadas).
- [ ] Estrategia de warm-up al arranque del servicio.

---

## 5. Comunicación Asíncrona (Message Broker)

**Prioridad:** Media  
**Referencia:** [ADR-019](../architecture/adrs/adr-019-synchronous-inter-module-communication-mvp.md)

- [ ] Introducir **NATS JetStream** como message broker.
- [ ] Migrar comunicación Sales → MES de síncrona a asíncrona.
- [ ] Implementar patrón Transactional Outbox en el Core.
- [ ] Domain Events publicados a streams NATS.

---

## 6. Inteligencia de Negocio y Analítica

**Prioridad:** Media  
**Referencia:** Presentaciones ([slides_spec.md](../presentations/slides_spec.md#L111))

- [ ] Cuadros de mando avanzados (dashboards).
- [ ] Analítica de ventas (tendencias, márgenes, KPIs).
- [ ] Reportes exportables.
- [ ] Módulo de **Analytics/BI** independiente.

---

## 7. Notificaciones en Tiempo Real

**Prioridad:** Baja  
**Referencia:** Presentaciones ([slides_spec.md](../presentations/slides_spec.md#L113))

- [ ] WebSockets para actualizaciones en tiempo real.
- [ ] Notificaciones por email.
- [ ] Alertas de cambios de estado (pedidos, producción).

---

## 8. Búsqueda Global Avanzada

**Prioridad:** Baja  
**Referencia:** Presentaciones, [erp-core-completion.md](../log/erp-core-completion.md#L450)

- [ ] Búsqueda full-text con Elasticsearch o similar.
- [ ] Búsqueda global unificada (productos, clientes, pedidos, facturas).
- [ ] Autocompletado y sugerencias.

---

## 9. Nuevos Módulos de Negocio

**Prioridad:** Variable  
**Referencia:** Presentaciones, [erp-core-completion.md](../log/erp-core-completion.md#L456-L460)

- [ ] **Compras** (Purchases) — gestión de pedidos a proveedores.
- [ ] **Inventario Avanzado** — control de stock, movimientos, alertas.
- [ ] **Logística** — gestión de envíos y transporte.
- [ ] **Contabilidad** — integración contable completa.
- [ ] **RRHH** (HR Module) — gestión de recursos humanos.

---

## 10. Mejoras Técnicas de Infraestructura

**Prioridad:** Baja  
**Referencia:** [erp-core-completion.md](../log/erp-core-completion.md#L443-L465)

- [ ] **Caché con Redis** — sesiones, datos frecuentes.
- [ ] **Refresh tokens** — mejora de seguridad de autenticación.
- [ ] **PWA para móvil** — acceso desde dispositivos móviles.
- [ ] **Internacionalización (i18n)** — soporte multi-idioma.

---

## 11. Cobertura de Tests

**Prioridad:** Alta  
**Referencia:** [erp-core-completion.md](../log/erp-core-completion.md#L283-L289)

Objetivos Post-MVP de cobertura:

- [ ] **Domain:** 100% (actualmente ~88-97% según módulo)
- [ ] **Application:** ≥95% (actualmente ~57.5%)
- [ ] **Infrastructure:** ≥80%
- [ ] **Frontend:** Tests E2E con Playwright (actualmente 0%)

---

## 12. Generación de Documentos PDF

**Prioridad:** Alta  
**Referencia:** [erp-core-completion.md](../log/erp-core-completion.md#L448)

- [ ] Generación de PDFs para albaranes.
- [ ] Generación de PDFs para facturas.
- [ ] Generación de PDFs para presupuestos/pedidos.
- [ ] Plantillas personalizables por empresa.

## 13. Facturación Consolidada Multi-Albarán

**Prioridad:** Media  
**Referencia:** [module-spec.md](../modules/sales/module-spec.md) — Fase 6, [use-cases.md](../modules/sales/use-cases.md) — CU-S-025

Permite agrupar múltiples albaranes de un mismo cliente en una única factura (ej: facturación mensual).  
La infraestructura base (campo `invoice_line_item_id` en `delivery_note_line_items`) se implementa en el MVP (Fase 5); la lógica de consolidación es Post-MVP.

- [ ] Pantalla de selección de albaranes pendientes por cliente y rango de fechas.
- [ ] Consolidación de líneas: las líneas con mismo `ProductVariantID`, `UnitPrice`, `DiscountAmount` y `TaxRate` se fusionan en una sola línea de factura sumando cantidades (relación N:1).
- [ ] Líneas con distinto producto/precio/descuento se mantienen como líneas independientes.
- [ ] Vinculación de todas las `DeliveryNoteLineItem`s al `InvoiceLineItemID` generado (N:1).
- [ ] Validación: todos los albaranes del mismo `Party`, todos en estado `DELIVERED`, ninguna línea previamente facturada.
- [ ] Actualización automática de estados de los pedidos relacionados.

---

## 14. Gestión Avanzada de Archivos de Diseño (MES)

**Prioridad:** Media  
**Contexto:**

En el MVP, `design_file_path` almacena una ruta de texto libre en `WorkSetupLine`. El campo se muestra en el Panel Tablet y en el configurador de WorkSetup, pero la interacción es mínima (copiar ruta al portapapeles). Las limitaciones del navegador impiden acceder al sistema de archivos local más allá del nombre de archivo.

**Objetivo Post-MVP:**

- [ ] **Vista previa de archivos** en el Panel Tablet y en el diálogo de detalle de WorkOrder/tarea: mostrar thumbnail o preview embebido (imágenes: PNG/JPG/SVG; vectores: AI/PDF si el navegador lo soporta; para otros formatos, icono de tipo de archivo).
- [ ] **Vista previa** también visible en las páginas de definición de WorkSetup (List / Edit / Create).
- [ ] **Abrir con aplicación por defecto**: botón "Abrir" exclusivo del configurador de WorkSetup (`/mes/work-setups/`) que invoque el protocolo del SO para abrir el archivo con su aplicación nativa (requiere integración Electron/Tauri o un agente de escritorio local).
- [ ] Soporte de rutas absolutas completas (implica abandonar el navegador puro — Electron/Tauri o extensión de escritorio).
- [ ] Almacenamiento opcional del archivo en servidor (upload): guardar en storage y servir URL firmada para preview remoto sin depender del sistema de archivos local.
- [ ] Validación de extensiones permitidas y tamaño máximo en el configurador.

**Dependencias:**

> La función "Abrir con app por defecto" y las rutas absolutas requieren un entorno de escritorio (Electron/Tauri). El resto de funcionalidades de preview son viables en web si se opta por upload a servidor.  
> Esta sección está vinculada a la decisión arquitectónica de si TramaTex tendrá cliente de escritorio.

---

## 15. Asignación de Tareas MES a Operarios

**Prioridad:** Media  
**Contexto:**

La columna `assigned_to` (FK a `users.id`) existe en `mes_work_tasks` desde el esquema inicial del MVP, y el campo está presente en el DTO del backend (`AssignedTo *uuid.UUID`) y en el tipo frontend (`WorkOrderTask.assigned_to`). Sin embargo, no se expone en ninguna pantalla del MVP porque la lógica de asignación de operarios pertenece a una fase posterior.

**Objetivo Post-MVP:**

- [ ] Selctor de operario (usuario) al crear o editar una tarea de WorkOrder en el Panel Tablet.
- [ ] Filtro en el terminal de tablet por operario asignado (`assigned_to = yo`).
- [ ] Visualizar el nombre del operario asignado en la tabla principal del terminal (columna “Asignado”).
- [ ] Lógica de reasignación: solo supervisores pueden reasignar tareas ya iniciadas.
- [ ] Notificación al operario cuando se le asigna una tarea (depende de sección 7 — Notificaciones).
- [ ] Registrar `assigned_by` y timestamp de asignación para auditoría.

**Estado actual (MVP):**

> Infraestructura lista (BD + backend DTO + tipo frontend). La UI muestra `—` para todas las tareas.
> La columna “Asignado” fue eliminada del Panel Tablet en MVP por no aportar valor todavía.
> Restaurarla y activarla es el primer paso de esta sección.

---

*Última actualización: 2026-03-20*
