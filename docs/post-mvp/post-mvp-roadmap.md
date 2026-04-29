# Hoja de Ruta Post-MVP — TramaTex

> Enumeración y orden de las mejoras, funcionalidades y módulos planificados para después del MVP.
> Este documento actúa como una declaración de intenciones estratégica.

---

### 🛠️ Metodología de Ejecución
Para comprender cómo se transforman estos hitos en realidad técnica (gestión de Sprints, flujo Backend → Persistencia → Frontend y criterios de Test Coverage), consulte el:
👉 **[Plan de Ejecución Post-MVP](post-mvp-execution-plan.md)**

---

## Índice

1. [Unificación UI/UX y Sistema de Diseño](#1.%20Unificación%20UI/UX%20y%20Sistema%20de%20Diseño)
2. [Facturación Consolidada Multi-Albarán](#2.%20Facturación%20Consolidada%20Multi-Albarán)
3. [Gestión de Cobros, Vencimientos y Tesorería](#3.%20Gestión%20de%20Cobros,%20Vencimientos%20y%20Tesorería)
4. [Adopción de Facturación Electrónica](#4.%20Adopción%20de%20Facturación%20Electrónica)
5. [Comunicación Asíncrona (Message Broker)](#5.%20Comunicación%20Asíncrona%20(Message%20Broker))
6. [Extracción MES como Microservicio](#6.%20Extracción%20MES%20como%20Microservicio)
7. [Caché y Rendimiento](#7.%20Caché%20y%20Rendimiento)
8. [Inteligencia de Negocio y Analítica](#8.%20Inteligencia%20de%20Negocio%20y%20Analítica)
9. [Búsqueda Global Avanzada](#9.%20Búsqueda%20Global%20Avanzada)
10. [Nuevos Módulos de Negocio](#10.%20Nuevos%20Módulos%20de%20Negocio)
11. [Mejoras Técnicas de Infraestructura](#11.%20Mejoras%20Técnicas%20de%20Infraestructura)
12. [Cobertura de Tests](#12.%20Cobertura%20de%20Tests)
13. [Gestión Avanzada de Archivos de Diseño (MES)](#13.%20Gestión%20Avanzada%20de%20Archivos%20de%20Diseño%20(MES))
14. [Integración Sales ↔ MES (Producto-Trabajo)](#14.%20Integración%20Sales%20↔%20MES%20(Producto-Trabajo))
15. [Asignación de Tareas MES a Operarios](#15.%20Asignación%20de%20Tareas%20MES%20a%20Operarios)

---

## 1. Unificación UI/UX y Sistema de Diseño

**Prioridad:** Máxima (Primera tarea Post-MVP)  
**Referencia Técnica Única:**
- [Plan Maestro de Unificación UI/UX](01-ui-ux-unification-master-plan.md)

**Contexto:**
Transformación de la UI en una herramienta industrial de alta eficiencia para operarios siguiendo el **Plan Maestro**. Foco en accesibilidad "Keyboard-First", iconografía profesional (Lucide) y feedback visual avanzado.

---

## 2. Facturación Consolidada Multi-Albarán

**Prioridad:** Media  
**Referencia Técnica Única:**
- [Estrategia de Facturación Consolidada](02-consolidated-billing-strategy.md)

**Contexto:**
Permite agrupar múltiples albaranes de un mismo cliente en una única factura (ej: facturación mensual), optimizando el flujo administrativo y reduciendo el volumen de documentos.

---

## 3. Gestión de Cobros, Vencimientos y Tesorería

**Prioridad:** Alta  
**Referencia Técnica Única:**
- [Estrategia de Gestión de Cobros y Tesorería](03-accounts-receivable-and-collections-strategy.md)

**Contexto:**
Extensión del ciclo financiero para controlar cuándo y cómo se cobran las facturas. Incluye la gestión de formas de pago, cálculo automático de vencimientos, registro de cobros parciales y control de deuda de clientes.

---

## 4. Adopción de Facturación Electrónica

**Prioridad:** Media (Requisito Legal)  
**Referencia Técnica Única:**
- [Análisis de Facturación Electrónica (Ley Crea y Crece)](04-electronic-invoicing-analysis.md)

**Contexto:**
Cumplimiento normativo con la Ley Crea y Crece y el sistema Veri*factu. Implica la generación de facturas en formato XML estructurado y la comunicación con las plataformas estatales.

---

## 5. Comunicación Asíncrona (Message Broker)

**Prioridad:** Alta (Infraestructura Core)  
**Referencia Técnica Única:**
- [Estrategia de Comunicación Asíncrona](05-asynchronous-communication-strategy.md)

**Contexto:**
Introducción de **NATS JetStream** para desacoplar módulos (Sales → MES, Product → Sales). Permite la coreografía de eventos de dominio y garantiza la integridad mediante el patrón **Transactional Outbox**.

---

## 6. Extracción MES como Microservicio

**Prioridad:** Media  
**Referencia Técnica Única:**
- [Estrategia de Extracción de MES](06-mes-microservice-extraction-strategy.md)

**Contexto:**
Extraer el módulo MES a un microservicio independiente para escalado y mantenimiento autónomo. Requiere la infraestructura de mensajería del punto anterior para la sincronización de datos y eventos de fabricación.

---

## 7. Caché y Rendimiento

**Prioridad:** Media  
**Referencia Técnica Única:**
- [Estrategia de Caché y Rendimiento](07-cache-and-performance-strategy.md)

**Contexto:**
Optimización del ERP central mediante el uso de **Redis**. Foco en el motor de precios (Pricing) y la hidratación rápida del catálogo de productos para reducir la latencia de respuesta.

---

## 8. Inteligencia de Negocio y Analítica

**Prioridad:** Baja  
**Referencia Técnica Única:**
- [Estrategia de BI y Analítica](08-business-intelligence-strategy.md)

**Contexto:**
Implementación de dashboards avanzados con visualización de tendencias, márgenes de beneficio y KPIs históricos. Incluye la capacidad de generar reportes de negocio exportables y personalizados.

---

## 9. Búsqueda Global Avanzada

**Prioridad:** Baja  
**Referencia Técnica Única:**
- [Estrategia de Búsqueda Avanzada](09-advanced-search-strategy.md)

**Contexto:**
Sistema de búsqueda full-text unificada accesible vía `Ctrl+K` para localizar instantáneamente clientes, productos y pedidos. Incorpora autocompletado inteligente y sugerencias basadas en el historial.

---

## 10. Nuevos Módulos de Negocio

**Prioridad:** Variable  
**Referencia Técnica Única:**
- [Estrategia de Nuevos Módulos](10-new-business-modules-strategy.md)

**Contexto:**
Expansión funcional hacia áreas críticas como **Compras** (aprovisionamiento), **Inventario Avanzado** (movimientos de stock y alertas) y **Logística** (gestión de envíos y transportistas.

---

## 11. Mejoras Técnicas de Infraestructura

**Prioridad:** Media  
**Referencia Técnica Única:**
- [Estrategia de Evolución de Infraestructura](11-infrastructure-evolution-strategy.md)

**Contexto:**
Evolución del despliegue actual hacia una arquitectura orquestada con **Kubernetes**. Implementación de un sistema de monitorización centralizado basado en Prometheus y Grafana para observabilidad total.

---

## 12. Cobertura de Tests

**Prioridad:** Alta  
**Referencia Técnica Única:**
- [Estrategia de QA y Testing](12-quality-assurance-and-testing-strategy.md)

**Contexto:**
Blindaje del sistema elevando la cobertura de tests en Dominio al 100% y en Aplicación al ≥95%. Implementación de una suite de tests E2E con Playwright para asegurar la integridad de los flujos críticos de usuario.

---

## 13. Gestión Avanzada de Archivos de Diseño (MES)

**Prioridad:** Media  
**Referencia Técnica Única:**
- [Estrategia de Gestión de Archivos MES](13-mes-advanced-file-management.md)

**Contexto:**
Mejora de la operativa en taller mediante la generación de vistas previas (thumbnails) de archivos técnicos e integración con aplicaciones nativas para la apertura directa de diseños desde el navegador.

---

## 14. Integración Sales ↔ MES (Producto-Trabajo)

**Prioridad:** Media  
**Referencia Técnica Única:**
- [Estrategia de Integración Profunda Sales-MES](14-sales-mes-deep-integration.md)

**Contexto:**
Establecimiento de una sincronización bidireccional automática. Cualquier cambio en un pedido de venta se reflejará instantáneamente en los trabajos de fabricación activos en el MES.

---

## 15. Asignación de Tareas MES a Operarios

**Prioridad:** Media  
**Referencia Técnica Única:**
- [Estrategia de Planificación de Operarios](15-mes-operator-planning-strategy.md)

**Contexto:**
Módulo de planificación de carga de trabajo por operario y máquina. Permite el registro preciso de tiempos por tarea y la optimización de la capacidad productiva de la planta.

---

*Última actualización: 2026-04-29*
