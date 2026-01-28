# TramaTex – Project Charter

**Versión:** 3.0  
**Fecha:** 9/01/2026  
**Autor:** Jorge Cortés Villalba  
**Copilotos LLM:** Gemini, ChatGPT, Claude (Anthropic)  
**Propósito:** Definición de visión, antecedentes y contexto estratégico del proyecto

---

## 📋 Antecedentes del Proyecto

Aunque la tendencia actual es el "cloud computing and cloud storage", más del 95% del tejido empresarial en España está formado por microempresas, muchas de las cuales no disponen de los recursos monetarios y/o conocimientos tecnológicos necesarios para aplicar estas metodologías, por lo que requieren soluciones integrables en sus sistemas locales y que sean fácilmente mantenibles, tanto en infraestructura como en monetario, por un tercero externo.

TramaTex surge de la necesidad de una de estas microempresas, dedicada a la venta de EPIs y vestuario laboral/personalizable, de contar con un software que:

- Permita gestionar pedidos estándar y personalizados, incluyendo variantes de talla y color de ciertos productos o, por ejemplo, marcaje de las prendas.
- Controle el estado de los pedidos, especialmente los que requieren procesos de diseño y producción (MES).
- Sea operativo en hardware limitado, sin depender de cloud computing ni grandes infraestructuras.
- Sea modular y mantenible, preparado para futuras expansiones sin generar deuda técnica.

### Problema Detectado

- El mercado ofrece soluciones demasiado genéricas o complejas/costosas.
- La empresa requiere control absoluto sobre inventario, pedidos y producción personalizada.
- Necesidad de un sistema local-first, eficiente y escalable opcionalmente hacia la nube.

---

## 🎯 Definición del Proyecto y Visión General

**Objetivo principal:**  
Controlar de forma fiable y trazable el ciclo completo de los pedidos, tanto estándar como personalizados, integrando gestión de ventas, tarificación, inventario, y producción (MES).

**Áreas del negocio cubiertas:**

- **Party / Organización** (clientes y proveedores como entidades con roles)
- **Ventas y Finanzas** (desktop)
- **Inventario y Almacén** (Post-MVP)
- **Producción y Taller** (tablet - MES)
- **Diseño y Pre-impresión**
- **Compras y Proveedores** (Proveedores como fuente de costes en MVP, compras formales Post-MVP)
- **Seguridad, Auditoría e Internacionalización** (transversales)

---

## 💻 Stack Tecnológico

**tramatex-api:**
- Go 1.21+, Gin Gonic, GORM, PostgreSQL 15+, Docker Compose

**Persistencia:**
- PostgreSQL 15+, GORM solo en capa de infraestructura

**Frontend:**
- Vue.js 3, Composition API, Pinia, Tailwind CSS

**Otros:**
- Generación de documentos: Web-to-Print en frontend
- Storage: NAS para archivos pesados, indexación en PostgreSQL
- Seguridad: JWT, roles básicos (Admin, Comercial, Diseño, Taller)
- i18n: Etiquetas estáticas básicas en frontend (MVP), go-i18n tramatex-api completo (Post-MVP)
- Despliegue: Docker Compose, sin Kubernetes en MVP

---

## 🖥️ Hardware Objetivo

- Ordenadores i3 con 8GB RAM (clientes)
- Servidor Linux i3 con 16GB RAM, SSD 2TB en espejo
- Tablets básicas para taller
- Posible inversión de 3.000€ para optimización y SAI

---

## 🏗️ Tipo de Aplicación

TramaTex se concibe como un **monolito modular**, diseñado bajo principios de **Domain-Driven Design y Clean Architecture con rigor asimétrico**, con una arquitectura preparada para una **evolución progresiva hacia microservicios** si el crecimiento del negocio lo justifica.

---

## 📚 Referencias a ADRs

Decisiones arquitectónicas fundamentales:

- [ADR-001: Selección del Stack Tecnológico y Estrategia Tecnológica Base](../engineering/architecture/adr/ADR-001-seleccion-stack-tecnologico.md)
- [ADR-002: Adopción de Clean Architecture y DDD con Rigor Asimétrico](../engineering/architecture/adr/ADR-002-adopcion-clean-architecture-ddd.md)
- [ADR-003: Tipo y Distribución de la Aplicación (Monolito Modular Local-First)](../engineering/architecture/adr/ADR-003-tipo-distribucion-aplicacion.md)
- [ADR-004: Ciclo de Vida de Desarrollo e Implementación hasta MVP](../engineering/architecture/adr/ADR-004-ciclo-vida-desarrollo-mvp.md)
- [ADR-005: Gestión Unificada de Clientes y Proveedores (Party / Organización)](../engineering/architecture/adr/ADR-005-gestion-unificada-clientes-proveedores.md)
- [ADR-006: Estrategia de Desarrollo Dirigida por Dominio (MVP)](../engineering/architecture/adr/ADR-006-estrategia-desarrollo-dirigido-dominio.md)
- [ADR-007: Orden de Implementación de Módulos e Infraestructura (Revisión Final)](../engineering/architecture/adr/ADR-007-orden-implementacion-modulos.md)
- [ADR-008: Planificación y Cronograma MVP Ajustado a Disponibilidad Real](../engineering/architecture/adr/ADR-008-planificacion-cronograma-mvp.md)
- [ADR-009: Estructura de Carpetas y Organización del Proyecto](../engineering/architecture/adr/ADR-009-estructura-proyecto.md)
- [ADR-010: Estrategia de Seguridad: Defensa en Profundidad y Security by Default](../engineering/architecture/adr/ADR-010-estrategia-seguridad-defensa-profundidad.md)
- [ADR-011: Estrategia de Testing y Coverage](../engineering/architecture/adr/ADR-011-estrategia-testing-coverage.md)

---

**Para especificación del MVP, ver:** [02-mvp-specification.md](02-mvp-specification.md)  
**Para contextos acotados, ver:** [03-bounded-contexts.md](03-bounded-contexts.md)  
**Para glosario de términos, ver:** [04-glossary.md](04-glossary.md)
