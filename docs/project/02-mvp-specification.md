# TramaTex – MVP Specification

**Versión:** 3.0  
**Fecha:** 9/01/2026  
**Autor:** Jorge Cortés Villalba  
**Copilotos LLM:** Gemini, ChatGPT, Claude (Anthropic)  
**Propósito:** Definición completa del MVP, alcance, fases y criterios de éxito

---

## 📋 Alcance del MVP

### Dominios Principales

- **Party / Organización** (fundacional)
- **Producto / Variante / Categoría** (fundacional)
- **Tarificación** (núcleo económico)
- **Ventas** (pedidos estándar)

### Subdominio Especializado

- **MES** para pedidos personalizados

---

## 🎯 Funcionalidades por Fase

### Fase 0 – Fundaciones Técnicas

**Objetivo:** Infraestructura mínima sin lógica de dominio.

**Entregables:**

- Estructura Clean Architecture
- Docker Compose funcional
- PostgreSQL configurado
- Autenticación JWT básica
- Roles: Admin, Comercial, Diseño, Taller
- Pipeline de tests (TDD)
- Frontend: Login funcional, estructura de navegación

**Criterios de aceptación:**

- ✅ Sistema arranca con `docker-compose up`
- ✅ Login funcional
- ✅ Pipeline de tests ejecutándose

**Duración estimada:** 1-2 semanas

---

### Fase 1 – Dominio Base para Tarificación

**Objetivo:** Núcleo económico funcional y testeable.

**Entregables:**

**Party / Organización (completo):**
- CRUD clientes/proveedores
- Jerarquías empresariales
- Herencia de descuentos
- Costes de proveedor por producto/variante

**Producto / Variante / Categoría (completo):**
- CRUD productos, variantes, categorías
- Modificadores de precio en variantes

**Tarificación (completo):**
- Motor de cálculo de precios
- Aplicación de márgenes y descuentos
- Integración con costes y descuentos

**Auditoría mínima:** Log de cambios en tarificación

**Frontend:**
- CRUD Party completo
- CRUD Producto completo
- Calculadora de tarificación funcional

**Criterios de aceptación:**

- ✅ CRUD Party funcional (frontend + tramatex-api)
- ✅ CRUD Producto funcional
- ✅ Cálculo de tarificación con datos reales
- ✅ Cobertura tests ≥80% en dominio tarificación
- ✅ Jerarquías funcionan correctamente

**Duración estimada:** 4-6 semanas

---

### Fase 2 – Casos de Uso Core y Pedidos

**Objetivo:** Flujo de ventas completo.

**Entregables:**

**Pedido Estándar (completo):**
- Creación, edición, seguimiento
- Cálculo automático de precios
- Gestión de estados

**Casos de uso de ventas:**
- Flujo completo: cotización → confirmación → entrega
- Validación de crédito

**Frontend:**
- Creación de pedidos completa
- Visualización y gestión de pedidos
- Generación de documentos (Web-to-Print)

**Criterios de aceptación:**

- ✅ Crear pedido estándar desde frontend
- ✅ Precio calculado automáticamente
- ✅ Pedidos persisten y visualizan correctamente
- ✅ Estados de pedido funcionan
- ✅ Documentos PDF funcionales
- ✅ Cobertura tests ≥70% en casos de uso pedidos

**Duración estimada:** 3-4 semanas

---

### Fase 3 – Subdominio Secundario MES

**Objetivo:** MVP completo y operativo.

**Entregables:**

**MES (completo):**
- Pedidos personalizados
- Estados de producción
- Trazabilidad de procesos

**Terminal de taller:**
- Interfaz tablet funcional
- Gestión de trabajos
- Cambio de estados

**Gestión documental:**
- Almacenamiento NAS
- Adjuntar/visualizar diseños

**Frontend:**
- Interfaz departamento diseño para almacenamiento de diseños 
- Terminal de taller operativa
- Gestión de pedidos personalizados
- Visualización de diseños

**Criterios de aceptación:**

- ✅ Terminal de taller operativa en tablet
- ✅ Pedidos personalizados extremo a extremo
- ✅ Almacenamiento NAS funcional
- ✅ Diseños visualizables
- ✅ **MVP listo para producción**
- ✅ Sistema probado con casos reales

**Duración estimada:** 3-4 semanas

---

## ⏱️ Total Estimado MVP

**11-16 semanas**

Nota: Estas estimaciones son orientativas. La duración real dependerá de:
- Disponibilidad del equipo
- Complejidad de casos de negocio específicos
- Descubrimiento de requisitos adicionales
- Iteraciones de feedback con usuarios

---

## 📋 Requisitos Funcionales MVP

| ID | Descripción | Fase |
|---|---|---|
| RF1 | Gestión de Party / Organización | Fase 1 |
| RF2 | Gestión de productos y variantes | Fase 1 |
| RF3 | Motor de tarificación | Fase 1 |
| RF4 | Gestión de ventas y pedidos estándar | Fase 2 |
| RF5 | Gestión de proveedores como fuente de costes | Fase 1 |
| RF6 | Ciclo de vida de producción (MES básico) | Fase 3 |
| RF7 | Terminal de taller | Fase 3 |
| RF8 | Gestión documental asociada a pedidos | Fase 3 |
| RF9 | Seguridad y control de acceso básico | Fase 0 |
| RF10 | Trazabilidad operativa mínima | Fase 1 |

---

## 📋 Requisitos No Funcionales MVP

| ID | Descripción | Métrica |
|---|---|---|
| RNF1 | Eficiencia de recursos | <150MB RAM |
| RNF2 | Operativa 100% local | No cloud required |
| RNF3 | Mantenibilidad base | TDD + cobertura ≥75% (ver ADR-011 para estrategia detallada) |
| RNF4 | Integridad de datos | ACID PostgreSQL |
| RNF5 | Seguridad básica | JWT + RBAC |
| RNF6 | Resiliencia operativa | Backups automáticos |

---

## ✅ Criterios de Éxito MVP

- Pedidos estándar y personalizados gestionados completamente dentro del sistema
- Tarificación funcional y fiable con cálculo automático de precios
- Taller puede operar con terminal sin comunicación externa
- Reducción de pérdida de información entre departamentos
- Sistema estable y operativo diariamente
- Cobertura de tests ≥75% global, ≥80% en dominio de tarificación
- Base clara para decisiones de expansión en proyecto independiente post-MVP

---

## 📅 Funcionalidades Post-MVP (NO se implementan en este proyecto)

Las siguientes funcionalidades se documentan para referencia futura:

**Gestión formal de compras:**
- Pedidos de compra a proveedores
- Recepción de mercancía (total/parcial)
- Actualización de stock desde compras
- Trazabilidad proveedor → producto → coste

**Control de stock completo:**
- Stock físico, lotes, ubicaciones
- Regularizaciones
- Reservas automáticas desde pedidos
- Alertas de stock mínimo

**Tarificación avanzada:**
- Descuentos complejos y campañas
- Reglas dinámicas de tarificación
- Promociones temporales

**Multimoneda y facturación legal completa**

**i18n completa tramatex-api** (go-i18n)

**Auditoría completa y observabilidad**

**Escalabilidad distribuida**

**Automatizaciones inteligentes en MES**

---

## 📌 Nota Importante

Estas funcionalidades Post-MVP **no se implementarán en el proyecto actual**, el cual finalizará con el MVP operativo. Cualquier evolución posterior se considerará un **nuevo proyecto independiente** (ver [ADR-004](/docs/adr/ADR-004-ciclo-vida-desarrollo-mvp.md)).

---

## 🔗 Notas Técnicas

### Motor de Tarificación (Núcleo Económico)

El motor de tarificación es **lógica de dominio**, independiente del frontend. La generación de documentos (PDF/Impresión) se delega al frontend como adaptador de infraestructura mediante Web-to-Print (CSS Media Queries), descargando al servidor de tareas de renderizado pesado y garantizando compatibilidad con hardware limitado.

**Fórmula de cálculo:**
```
Precio final = ((Coste base + Modificadores de variante) × (1 + Margen) + (Modificadores no dependientes)) × (1 - Descuento total)
```

### Implementación Incremental

El desarrollo sigue el **orden de módulos definido en ADR-007**, asegurando que Party / Organización y Producto / Variante / Categoría se implementen primero, seguido de Tarificación, Pedidos y MES. El frontend se desarrolla en paralelo en todas las fases, de modo que cada fase deja interfaces operativas para los módulos implementados.

---

**Para contextos acotados detallados, ver:** [03-bounded-contexts.md](03-bounded-contexts.md)  
**Para charter del proyecto, ver:** [01-project-charter.md](01-project-charter.md)  
**Para glosario de términos, ver:** [04-glossary.md](04-glossary.md)
