# ADR-008 – MVP Timeline Planning (Planificación y Cronograma del MVP)

**Fecha:** 2026-02-01 (Reconstruido)  
**Estado:** Aceptado  
**Autores:** Equipo de Arquitectura de TramaTex  

---

## 1. Contexto

Para gestionar eficazmente el desarrollo del MVP (Minimum Viable Product) y asegurar la entrega de valor de manera incremental, es necesario establecer un cronograma de alto nivel. Este plan organiza la implementación de los módulos, definidos en `ADR-007 (Module Implementation Order)`, en fases lógicas con objetivos y plazos claros.

Esta planificación debe ser lo suficientemente flexible para adaptarse a imprevistos, pero lo suficientemente estructurada para proporcionar una hoja de ruta clara al equipo de desarrollo y a los stakeholders.

---

## 2. Decisión Adoptada

Se establece un plan de desarrollo por fases, alineado con el orden de implementación de módulos (`ADR-007 (Module Implementation Order)`). Cada fase agrupa un conjunto de módulos relacionados que constituyen una entrega de valor coherente.

El cronograma se basa en la información de estado de fases y la asignación de módulos a fases, extraída de los agentes de contexto del proyecto (`sprint-registry.yaml` y `bounded-contexts.yaml`).

### Fases del Proyecto

**Fase 0: Fundaciones Técnicas**
- **Estado:** ✅ Completada
- **Plazo de Finalización:** Q1 2026
- **Módulos Incluidos:**
  - `IAM (Identity & Access Management)`: Sistema de autenticación y autorización.
- **Otros Entregables:**
  - Configuración del entorno de desarrollo con Docker.
  - Definición del workflow de Git.
  - Estructura inicial del proyecto y sistema de documentación.
  - Configuración inicial de CI/CD.

**Fase 1: Dominio Base y Lógica Crítica de Negocio**
- **Estado:** 🔄 En Progreso
- **Plazo Objetivo:** Q2 - Q4 2026
- **Módulos Incluidos:**
  - `Party (Gestión de Entidades)`: Clientes y proveedores.
  - `Product (Catálogo de Productos)`: Productos, variantes y categorías.
  - `Pricing (Motor de Precios)`: Lógica de cálculo de precios.
- **Otros Entregables:**
  - Sistema de Diseño de la UI.
  - Auditoría de seguridad OWASP y primeras mitigaciones.

**Fase 2: Flujo de Negocio Principal**
- **Estado:** ⏳ Planificada
- **Plazo Objetivo:** Post-Fase 1
- **Módulos Incluidos:**
  - `Sales (Ventas y Pedidos)`: Creación de pedidos, cotizaciones y transacciones.
- **Descripción:** Esta fase se centra en la orquestación de los módulos base para gestionar el flujo de ventas completo.

**Fase 3: Operaciones de Taller (MVP)**
- **Estado:** ⏳ Planificada
- **Plazo Objetivo:** Post-Fase 2
- **Módulos Incluidos:**
  - `MES (Manufacturing Execution System)`: Gestión de órdenes de producción, control de calidad y seguimiento en el taller.
- **Descripción:** Extiende la funcionalidad del ERP al área de producción.

---

## 3. Consecuencias

### Positivas
- Proporciona una visión clara del progreso del proyecto y los próximos pasos.
- Permite a los stakeholders entender cuándo se espera que las diferentes funcionalidades estén disponibles.
- Ayuda al equipo de desarrollo a enfocarse en los objetivos de la fase actual.
- Facilita la planificación de recursos y la gestión de expectativas.

### Negativas
- Las fechas son estimaciones de alto nivel y pueden cambiar, lo que requiere una comunicación constante para gestionar las expectativas.
- Un enfoque por fases puede reducir la agilidad si se interpreta de manera demasiado rígida. Los sprints dentro de cada fase deben seguir siendo flexibles.

---

## 4. Integración con otros ADRs

- **ADR-007 (Module Implementation Order):** Este cronograma es la aplicación práctica del orden de dependencias definido en `ADR-007`.
- **ADR-004 (MVP Development Lifecycle):** Define el enfoque de desarrollo iterativo que se sigue dentro de las fases aquí descritas.

---
*Nota: Este documento ha sido reconstruido el 2026-02-01 a partir de las referencias encontradas en los agentes de contexto y otros documentos del proyecto.*
