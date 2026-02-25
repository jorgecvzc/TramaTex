# 💡 TramaTex – Project Charter y Especificación del MVP

---

**Versión:** 3.0
**Fecha:** 28/01/2026 (Consolidado)
**Autor:** Jorge Cortés Villalba
**Propósito:** Documento central que unifica la visión, el contexto estratégico y la especificación detallada del Minimum Viable Product (MVP).

---

## 📋 Antecedentes y Visión del Proyecto

Aunque la tendencia actual es el "cloud computing", más del 95% del tejido empresarial en España está formado por microempresas. Muchas carecen de los recursos o conocimientos para aplicar estas metodologías y requieren soluciones **integrables en sistemas locales** y de fácil mantenimiento.

TramaTex surge de la necesidad de una de estas microempresas, dedicada a la venta de EPIs y vestuario laboral, de contar con un software que:

- Gestione pedidos estándar y personalizados (con variantes de producto y marcaje).
- Controle el estado de los pedidos, incluyendo el proceso de producción (MES).
- Sea operativo en **hardware limitado** y no dependa de la nube (`local-first`).
- Sea **modular y mantenible**, preparado para futuras expansiones.

### Problema Detectado

- El mercado ofrece soluciones genéricas o demasiado complejas/costosas.
- La empresa necesita control absoluto sobre inventario, pedidos y producción.
- Se requiere un sistema `local-first`, eficiente y opcionalmente escalable a la nube.

---

## 🎯 Definición y Alcance del MVP

**Objetivo Principal:** Controlar de forma fiable y trazable el ciclo completo de los pedidos, tanto estándar como personalizados, integrando gestión de ventas, tarificación y producción (MES).

### Dominios Cubiertos en el MVP

- **Party / Organización** (fundacional)
- **Producto / Variante / Categoría** (fundacional)
- **Tarificación** (núcleo económico)
- **Ventas** (pedidos estándar)
- **MES** (subdominio para producción personalizada)
- **Seguridad, Auditoría e Internacionalización** (transversales, básicos en MVP)

### Funcionalidades Post-MVP (Fuera de Alcance)

Las siguientes funcionalidades se documentan para referencia futura y **no se implementarán** en este proyecto (ver [**ADR-004**](./adrs/ADR-004-mvp-development-lifecycle.md)):
- Gestión formal de compras y control de stock completo.
- Tarificación avanzada (campañas, promociones).
- Multimoneda y facturación legal completa.
- Auditoría y observabilidad avanzadas.

---

## 🏗️ Arquitectura y Stack Tecnológico

TramaTex se concibe como un **monolito modular**, diseñado bajo principios de **Domain-Driven Design y Clean Architecture con rigor asimétrico**. Esta arquitectura está preparada para una evolución progresiva hacia microservicios si el negocio lo justifica.

### Stack Tecnológico

| Componente | Tecnología | Justificación |
| :--- | :--- | :--- |
| **Backend API** | Go 1.21+, Gin Gonic | Eficiencia, binario único, bajo consumo de recursos. |
| **Persistencia** | PostgreSQL 15+, GORM | Motor de datos robusto y transaccional (ACID). |
| **Frontend** | Vue.js 3, Pinia, Vite | Curva de aprendizaje moderada, ecosistema maduro. |
| **Estilos** | Tailwind CSS | Enfoque "Utility-first" para un desarrollo rápido. |
| **Infraestructura** | Docker / Docker-Compose | Entornos de desarrollo y despliegue consistentes. |
| **Documentos** | Web-to-Print (Frontend) | Delega la carga de generación de PDFs al cliente. |

### Hardware Objetivo

- **Clientes:** PCs i3 con 8GB RAM.
- **Servidor:** Linux i3 con 16GB RAM, SSD 2TB en espejo.
- **Taller:** Tablets básicas.

---

## 🚀 Fases de Implementación del MVP

El desarrollo sigue un orden estricto definido en [**ADR-007**](./adrs/ADR-007-module-implementation-order.md) para garantizar la coherencia del dominio.

### Fase 0 – Fundaciones Técnicas
- **Objetivo:** Infraestructura mínima sin lógica de negocio.
- **Entregables:** Estructura Clean Architecture, Docker Compose, autenticación JWT básica, pipeline de tests.
- **Criterio de Aceptación:** `docker-compose up` arranca el sistema y el login es funcional.

### Fase 1 – Dominio Base para Tarificación
- **Objetivo:** Construir el núcleo económico funcional y testeable.
- **Entregables:** Módulos `Party`, `Producto` y `Tarificación` completamente funcionales en el backend y frontend.
- **Criterio de Aceptación:** Se puede calcular un precio real con datos de clientes y productos. **Cobertura de tests ≥90%** en el dominio de tarificación.

### Fase 2 – Casos de Uso Core y Pedidos
- **Objetivo:** Orquestar el flujo de ventas completo.
- **Entregables:** Sistema de creación y seguimiento de pedidos estándar, generación de documentos (presupuestos, albaranes).
- **Criterio de Aceptación:** Se puede crear un pedido estándar de principio a fin. **Cobertura de tests ≥85%** en casos de uso de pedidos.

### Fase 3 – Subdominio Secundario MES
- **Objetivo:** Completar el MVP con la gestión de producción personalizada.
- **Entregables:** Gestión de estados de producción, terminal de taller funcional para tablets, gestión documental en NAS.
- **Criterio de Aceptación:** Un pedido personalizado puede ser gestionado desde su diseño hasta su finalización. **MVP listo para producción.**

---

## ✅ Requisitos y Criterios de Éxito del MVP

### Requisitos Funcionales Clave

| ID | Descripción | Fase |
|---|---|---|
| RF1 | Gestión de Party / Organización | 1 |
| RF2 | Gestión de productos y variantes | 1 |
| RF3 | Motor de tarificación | 1 |
| RF4 | Gestión de ventas y pedidos estándar | 2 |
| RF6 | Ciclo de vida de producción (MES básico) | 3 |
| RF7 | Terminal de taller | 3 |
| RF9 | Seguridad y control de acceso básico | 0 |

### Requisitos No Funcionales

| ID | Descripción | Métrica |
|---|---|---|
| RNF1 | Eficiencia de recursos | <150MB RAM en servidor. |
| RNF2 | Operativa 100% local | No requiere conexión a internet. |
| RNF3 | Mantenibilidad base | TDD + **cobertura ≥85%** (ver [**ADR-011**](./adrs/ADR-011-testing-coverage-strategy.md) para detalles). |
| RNF4 | Integridad de datos | Transacciones ACID garantizadas por PostgreSQL. |
| RNF5 | Seguridad básica | Autenticación JWT y RBAC (Roles: Admin, Comercial, Diseño, Taller). |

### Criterios de Éxito Finales

- El sistema gestiona pedidos estándar y personalizados de principio a fin.
- La tarificación es fiable y automática.
- El taller puede operar de forma autónoma con la tablet.
- Se reduce la pérdida de información entre departamentos.
- El sistema es estable para el uso diario.
- Se alcanza una **cobertura de tests global ≥85%**, con **≥90% en el dominio de tarificación**.
