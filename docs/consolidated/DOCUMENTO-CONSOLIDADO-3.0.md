# TramaTex – Documentación Consolidada

**Autor:** Jorge Cortés Villalba  
**Copilotos LLM:** Gemini, ChatGPT, Claude (Anthropic)  
**Fecha:** 10/01/2026  
**Estatus:** Consolidación de Documentación, MVP definido, ADRs completos  
**Versión:** 3.0

---

## 0. Antecedentes del Proyecto

Aunque la tendencia actual es el "cloud computing and cloud storage", más del 95% del tejido empresarial en España está formado por microempresas, muchas de las cuales no disponen de los recursos monetarios y/o conocimientos tecnológicos necesarios para aplicar estas metodologías, por lo que requieren soluciones integrables en sus sistemas locales y que sean fácilmente mantenibles, tanto en infraestructura como en monetario, por un tercero externo.

TramaTex surge de la necesidad de una de estas microempresas, dedicada a la venta de EPIs y vestuario laboral/personalizable, de contar con un software que:

- Permita gestionar pedidos estándar y personalizados, incluyendo variantes de talla y color de ciertos productos o, por ejemplo, marcaje de las prendas.
- Controle el estado de los pedidos, especialmente los que requieren procesos de diseño y producción (MES).
- Sea operativo en hardware limitado, sin depender de cloud computing ni grandes infraestructuras.
- Sea modular y mantenible, preparado para futuras expansiones sin generar deuda técnica.

Problema detectado:

- El mercado ofrece soluciones demasiado genéricas o complejas/costosas.
- La empresa requiere control absoluto sobre inventario, pedidos y producción personalizada.
- Necesidad de un sistema local-first, eficiente y escalable opcionalmente hacia la nube.

---

## 1. Definición del Proyecto y Visión General

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

**Stack Tecnológico:**

- Backend: Go 1.21+, Gin Gonic, GORM, PostgreSQL 15+, Docker Compose
- Persistencia: PostgreSQL 15+, GORM solo en capa de infraestructura
- Frontend: Vue.js 3, Composition API, Pinia, Tailwind CSS
- Generación de documentos: Web-to-Print en frontend
- Storage: NAS para archivos pesados, indexación en PostgreSQL
- Seguridad: JWT, roles básicos (Admin, Comercial, Diseño, Taller)
- i18n: Etiquetas estáticas básicas en frontend (MVP), go-i18n backend completo (Post-MVP)
- Despliegue: Docker Compose, sin Kubernetes en MVP

**Hardware objetivo:**

- Ordenadores i3 con 8GB RAM (clientes)
- Servidor Linux i3 con 16GB RAM, SSD 2TB en espejo
- Tablets básicas para taller
- Posible inversión de 3.000€ para optimización y SAI

**Tipo de aplicación:**

TramaTex se concibe como un **monolito modular**, diseñado bajo principios de **Domain-Driven Design y Clean Architecture con rigor asimétrico**, con una arquitectura preparada para una **evolución progresiva hacia microservicios** si el crecimiento del negocio lo justifica.

**Referencias a ADRs:**

- ADR-001: Selección del Stack Tecnológico y Estrategia Tecnológica Base
- ADR-002: Adopción de Clean Architecture y DDD con Rigor Asimétrico
- ADR-003: Tipo y Distribución de la Aplicación (Monolito Modular Local-First)
- ADR-004: Ciclo de Vida de Desarrollo e Implementación hasta MVP
- ADR-005: Gestión Unificada de Clientes y Proveedores 
- ADR-006: Estrategia de Desarrollo Dirigida por Dominio (MVP)
- ADR-007: Orden de Implementación de Módulos y Creación de Infraestructura

---

## 2. Definición del MVP

### 2.1 Alcance

**Dominio principal:**

- **Party / Organización** (fundacional)
- **Producto / Variante / Categoría** (fundacional)
- **Tarificación** (núcleo económico)
- **Ventas** (pedidos estándar)

**Subdominio especializado:**

- **MES** para pedidos personalizados

### 2.2 Funcionalidades MVP (ordenadas por prioridad de implementación)

**Fase 0 – Fundaciones Técnicas:**

- Infraestructura base (Docker, PostgreSQL, estructura Clean Architecture)
- Autenticación y autorización básica (JWT, RBAC mínimo)
- Pipeline de calidad (TDD, tests, linters)

**Fase 1 – Dominio Base para Tarificación:**

- **Gestión de Party / Organización:**
    - Clientes y proveedores como roles de una entidad común
    - Jerarquías empresariales para clientes (empresa matriz → dependientes)
    - Herencia de descuentos en jerarquías
    - Proveedores sin jerarquía (planos)
- **Gestión de productos con variantes y clasificación:**
    - Productos como entidad raíz
    - Variantes (talla, color, modificaciones) que afectan precio
    - Categorías para clasificación (no afectan tarificación)
- **Motor de tarificación MVP:**
    - Cálculo de precios de venta
    - Aplicación de márgenes
    - Gestión de descuentos (base, heredados, específicos)
    - Integración con costes de proveedor
- **Gestión de proveedores como fuente de costes:**
    - Costes base por producto/variante
    - Histórico simple de costes
    - Selección de proveedor preferente
    - **NO incluye gestión de compras formales** (Post-MVP)

**Fase 2 – Casos de Uso Core y Pedidos:**

- **Gestión de pedidos estándar:**
    - Creación, edición, seguimiento de pedidos
    - Selección de cliente (rol de Party)
    - Selección de productos/variantes
    - Cálculo automático de precios usando tarificación
    - Estados de pedido (Borrador, Confirmado, En preparación, Entregado)
    - Documentos enlazados con pedidos: presupuestos, albaranes, facturas.
    - Generación de documentos (delegado a frontend vía Web-to-Print)

**Fase 3 – Subdominio Secundario MES:**

- **Terminal de taller para pedidos personalizados:**
    - Interfaz simplificada para tablets
    - Gestión de estados de producción (Diseño, Aprobación, Marcaje, Taller, Control calidad, Listo)
    - Visualización de trabajos asignados
    - Trazabilidad de procesos productivos
- **Gestión documental mínima para diseños:**
    - Almacenamiento en NAS
    - Indexación en PostgreSQL
    - Adjuntar/visualizar archivos asociados a pedidos personalizados
- **Auditoría básica de cambios críticos:**
    - Trazabilidad de estados de producción
- **Seguridad básica y roles:**
    - Admin, Comercial, Diseño, Taller

> **Nota crítica sobre el motor de tarificación:**  
> El motor de tarificación es **lógica de dominio**, independiente del frontend. La generación de documentos (PDF/Impresión) se delega al frontend como adaptador de infraestructura mediante Web-to-Print (CSS Media Queries), descargando al servidor de tareas de renderizado pesado y garantizando compatibilidad con hardware limitado.
> 
> **Implementación incremental:** El desarrollo sigue el **orden de módulos definido en ADR-007**, asegurando que Party / Organización y Producto / Variante / Categoría se implementen primero, seguido de Tarificación, Pedidos y MES. El frontend se desarrolla en paralelo en todas las fases, de modo que cada fase deja interfaces operativas para los módulos implementados.

**Referencias ADRs:** ADR-007, ADR-006, ADR-005, ADR-004, ADR-001, ADR-002, ADR-003

### 2.3 Funcionalidades Post-MVP

- **Gestión formal de compras:**
    - Pedidos de compra a proveedores
    - Recepción de mercancía (total/parcial)
    - Actualización de stock desde compras
    - Trazabilidad proveedor → producto → coste
- **Control de stock completo:**
    - Stock físico, lotes, ubicaciones
    - Regularizaciones
    - Reservas automáticas desde pedidos
    - Alertas de stock mínimo
- **Tarificación avanzada:**
    - Descuentos complejos y campañas
    - Reglas dinámicas de tarificación
    - Promociones temporales
- **Multimoneda y facturación legal completa**
- **i18n completa backend** (go-i18n)
- **Auditoría completa y observabilidad**
- **Escalabilidad distribuida**
- **Automatizaciones inteligentes en MES**

> **Nota importante:** Estas funcionalidades se documentan para referencia futura, **pero no se implementarán en el proyecto actual**, el cual finalizará con el MVP operativo. Cualquier evolución posterior se considerará un **nuevo proyecto independiente** (ver ADR-004).

### 2.4 Criterios de éxito MVP

- Pedidos estándar y personalizados gestionados completamente dentro del sistema
- Tarificación funcional y fiable con cálculo automático de precios
- Taller puede operar con terminal sin comunicación externa
- Reducción de pérdida de información entre departamentos
- Sistema estable y operativo diariamente
- Cobertura de tests ≥75% global, ≥80% en dominio de tarificación
- Base clara para decisiones de expansión en proyecto independiente post-MVP

---

## 3. Requisitos Funcionales y No Funcionales Priorizados

**Nota:** La trazabilidad entre requisitos y módulos implementados sigue ADR-006, asegurando que el núcleo económico (tarificación) se construya sobre Party y Producto. La estrategia de desarrollo y testing se detalla en ADR-006 y ADR-004. Las decisiones de infraestructura se alinean con ADR-001 y ADR-003.

### 3.1 Requisitos Funcionales MVP

- **RF1: Gestión de Party / Organización**
    - Clientes y proveedores como roles de entidades
    - Jerarquías empresariales para clientes
    - Herencia de descuentos
- **RF2: Gestión de productos y variantes**
    - CRUD de productos, variantes, categorías
    - Relación productos → variantes → categorías
- **RF3: Motor de tarificación**
    - Cálculo de precios de venta
    - Aplicación de márgenes y descuentos
    - Integración con costes de proveedor
- **RF4: Gestión de ventas y pedidos estándar**
    - Creación, edición, seguimiento de pedidos
    - Cálculo automático de precios
    - Gestión de estados
    - Control Documental (Presupuesto → Pedido → Albarán → Factura)
- **RF5: Gestión de proveedores como fuente de costes**
    - Registro de costes base por producto/variante
    - Histórico simple
    - Selección de proveedor preferente
- **RF6: Ciclo de vida de producción (MES básico)**
    - Estados de producción para pedidos personalizados
    - Trazabilidad de procesos
- **RF7: Terminal de taller**
    - Interfaz simplificada para tablets
    - Gestión de trabajos y estados
- **RF8: Gestión documental asociada a pedidos**
    - Almacenamiento NAS para diseños
    - Indexación en PostgreSQL
- **RF9: Seguridad y control de acceso básico**
    - Autenticación JWT
    - Roles: Admin, Comercial, Diseño, Taller
- **RF10: Trazabilidad operativa mínima**
    - Auditoría de cambios críticos en tarificación
    - Log de estados de producción

### 3.2 Requisitos Funcionales Post-MVP

- **RF11: Gestión formal de compras**
    - Pedidos de compra, recepciones, actualización de stock
- **RF12: Control de stock completo**
    - Lotes, ubicaciones, regularizaciones, reservas, alertas
- **RF13: RBAC avanzado**
    - Permisos granulares, roles personalizados
- **RF14: i18n backend completa**
    - go-i18n para todos los mensajes del sistema
- **RF15: Auditoría completa**
    - Trazabilidad total de cambios, retención de datos
- **RF16: Tarificación avanzada**
    - Descuentos complejos, campañas, reglas dinámicas
- **RF17: Multimoneda**
    - Soporte múltiples divisas, conversiones, contabilidad multimoneda
- **RF18: Escalabilidad distribuida**
    - Preparación para microservicios, eventos de dominio

### 3.3 Requisitos No Funcionales MVP

- **RNF1: Eficiencia de recursos** (consumo <150MB RAM)
- **RNF2: Operativa 100% local**
- **RNF3: Mantenibilidad base** (TDD en dominio crítico, cobertura ≥75%)
- **RNF4: Integridad de datos** (ACID PostgreSQL)
- **RNF5: Seguridad básica** (hash de passwords, JWT, roles)
- **RNF6: Resiliencia operativa** (backups automáticos)

### 3.4 Requisitos No Funcionales Post-MVP

- **RNF7–RNF11:** Alta cobertura testing (≥90%), CI/CD completo, alta disponibilidad, multitenencia, observabilidad avanzada

---

## 4. Posicionamiento del Sistema

TramaTex es un ERP ligero implementado como un **monolito modular**, centrado en la gestión de ventas y tarificación, con subdominios especializados (MES, Compras, Inventario) claramente delimitados y preparados para una evolución futura hacia arquitecturas distribuidas.

### Arquitectura del Sistema

**Dominio principal (MVP):**

- **Party / Organización** (fundacional)
- **Producto / Variante / Categoría** (fundacional)
- **Tarificación** (núcleo económico)
- **Ventas** (pedidos estándar)

**Subdominio especializado (MVP):**

- **MES** – Producción personalizada

**Dominios Post-MVP:**

- **Compras** (gestión formal de compras)
- **Inventario / Stock** (control completo)

**Módulos transversales:**

- Seguridad (MVP: autenticación JWT, RBAC básico)
- Auditoría (MVP: cambios críticos; Post-MVP: completa)
- Gestión documental (MVP: NAS para diseños; Post-MVP: gestión avanzada)
- i18n (MVP: etiquetas estáticas frontend; Post-MVP: backend completo)

### Flujos de negocio MVP

**Pedidos estándar:**

```
Cliente (Party rol Cliente) 
→ Selección productos/variantes 
→ Cálculo automático precio (Tarificación) 
→ Confirmación pedido 
→ Estados (Confirmado → En preparación → Entregado)
```

**Pedidos personalizados:**

```
Cliente (Party rol Cliente) 
→ Selección productos/variantes 
→ Cálculo automático precio (Tarificación)
→ Adjuntar diseño 
→ Confirmación pedido 
→ MES (Diseño → Aprobación → Marcaje → Taller → Control calidad) 
→ Entrega
```

**Tarificación (núcleo económico):**

```
Producto/Variante 
+ Coste base (Proveedor) 
+ Margen 
- Descuento (Cliente: base, heredado de jerarquía, o específico + puntual) 
= Precio final de venta
```

**Nota:** El desarrollo sigue ADR-007, priorizando Party → Producto → Tarificación → Pedidos → MES, con frontend paralelo en todas las fases. Esto asegura que al final de cada fase los módulos implementados ya sean operativos.

### Beneficios

- Refleja flujo real del negocio
- Prioriza núcleo económico (tarificación) construido sobre bases sólidas (Party, Producto)
- MES aislado como subdominio, listo para escalar o refactorizar
- Modularidad y mantenimiento más sencillo
- Adaptadores de infraestructura (persistencia, generación de documentos) intercambiables sin afectar el dominio

**Referencias ADRs:** ADR-003, ADR-002, ADR-007, ADR-006, ADR-005

---

## 5. Módulos Canónicos y Bounded Contexts

Los módulos representan **Bounded Contexts internos** dentro de un **monolito modular**, con límites claros y contratos explícitos entre ellos. Cada módulo implementa **Clean Architecture con rigor asimétrico**:

- Dominio aislado de infraestructura y frontend
- Casos de uso testables de forma independiente
- Adaptadores de persistencia e interfaz intercambiables

**Prioridad de implementación para MVP (ADR-007):**

1. **Fase 0:** Fundaciones técnicas (infraestructura, sin dominio)
2. **Fase 1:** Party / Organización + Producto / Variante / Categoría + Tarificación
3. **Fase 2:** Pedidos estándar (ventas)
4. **Fase 3:** MES – Producción personalizada
5. **Frontend desarrollado en paralelo en todas las fases**

**Referencias ADRs:** ADR-002, ADR-003, ADR-005, ADR-006, ADR-007

---

### 5.1 Dominio principal – Party / Organización (MVP - Fundacional)

**Concepto:** Party es un patrón de modelado que representa a cualquier persona u organización que tiene relación con el sistema, independientemente del rol que desempeñe.

**Problema que resuelve:**

- Evita duplicación de entidades (Cliente y Proveedor separados)
- Permite que una entidad tenga múltiples roles simultáneamente
- Garantiza consistencia de datos (una sola fuente de verdad)

**Entidades:**

- **Party** (Entidad raíz): ID, Tipo (Persona/Organización), Nombre, NIF/CIF, Dirección, Contacto
- **PartyRole**: Tipo de rol (Cliente, Proveedor), Estado (Activo/Inactivo), Fechas
- **Customer** (Especialización para rol Cliente):
    - Descuento base
    - Empresa matriz (para jerarquías)
    - Límite de crédito
    - Condiciones de pago
- **Supplier** (Especialización para rol Proveedor):
    - Código de proveedor
    - Días de entrega
    - Pedido mínimo
- **SupplierCost**: Costes base por producto/variante (histórico)

**Reglas de negocio específicas:**

**Clientes:**

- Pueden tener jerarquía empresarial (empresa matriz → dependientes)
- Los descuentos pueden heredarse de la empresa matriz
- Pueden sobrescribir descuentos heredados con descuentos propios
- Ejemplo: "Construcciones ABC S.L." (matriz) con descuento 10%, sus obras heredan ese descuento

**Proveedores:**

- **NO tienen jerarquía** (simplificación MVP, según ADR-005)
- Proporcionan costes base para productos/variantes
- Estos costes son inputs para el motor de tarificación
- Un mismo Party puede ser Cliente Y Proveedor simultáneamente

**Funciones:**

- CRUD de Party con gestión de roles
- Activar/desactivar roles
- Gestión de jerarquías empresariales (solo clientes)
- Registro y consulta de costes de proveedor
- Histórico de cambios en costes

**Dependencias:**

- Ninguna (es fundacional)
- Provee servicios a: Tarificación, Ventas, MES

**Adaptadores:**

- Repositorio PostgreSQL para Party, Roles, Costes

**Implementación:** Fase 1 (ADR-007)

---

### 5.2 Dominio principal – Producto / Variante / Categoría (MVP - Fundacional)

**Entidades:**

- **Producto** (Entidad raíz): ID, Código, Nombre, Descripción
- **Variante**: ID, Producto (FK), Tipo (Talla, Color, Arreglo, etc.), Valor, Modificador de precio
- **Categoría**: ID, Nombre, Descripción, Categoría padre (jerarquía)

**Reglas de negocio:**

- Toda variante **puede afectar al precio** (modificador de precio)
- Las categorías **NO influyen en tarificación**, solo estructuran el catálogo
- Un producto puede tener múltiples variantes
- Las variantes son parte activa del modelo económico

**Funciones:**

- CRUD de productos, variantes, categorías
- Gestión de catálogo y clasificación
- Asignación de modificadores de precio a variantes
- Búsqueda y filtrado por categorías

**Dependencias:**

- Ninguna (es fundacional)
- Provee servicios a: Tarificación, Ventas, MES

**Adaptadores:**

- Repositorio PostgreSQL para Producto, Variante, Categoría

**Implementación:** Fase 1 (ADR-006)

---

### 5.3 Dominio principal – Tarificación (MVP - Núcleo Económico)

**Entidades:**

- **Tarifa**: Reglas de cálculo de precio
- **Regla de tarificación**: Márgenes, descuentos, condiciones
- **Precio calculado**: Resultado del motor de tarificación (no persistido, calculado on-demand)

**Reglas de negocio (Motor de tarificación):**

```
Precio final = ((Coste base + Modificadores de variante) × (1 + Margen) + (Modificadores no dependientes de proveedor)) × (1 - Descuento total)

Donde:
- Coste base: obtenido de Supplier (Party rol Proveedor)
- Modificadores de variante: de Variante
- Margen: regla de tarificación (por producto, categoría, o general)
- Modificadores no dependientes: p.e. el arreglo de una prenda no depende del proveedor.
- Descuento total: 
  - Descuento base del Cliente (Party rol Cliente)
  - O descuento heredado de empresa matriz
  - O descuento específico del pedido (override)
```

**Funciones:**

- Calcular precio de venta para producto/variante + cliente
- Aplicar márgenes según reglas
- Aplicar descuentos (base, heredados, específicos)
- Validar coherencia de precios (margen mínimo)
- Histórico de cambios en reglas de tarificación (auditoría)

**Dependencias:**

- **Requiere:** Party (para costes de proveedor y descuentos de cliente), Producto/Variante
- **Provee servicios a:** Ventas

**Adaptadores:**

- Repositorio PostgreSQL para Tarifas y Reglas
- **NO tiene adaptador de generación de documentos** (delegado a frontend)

**Implementación:** Fase 1 (ADR-007)

**Nota crítica:** El motor de tarificación es **lógica de dominio pura**, completamente testeable en aislamiento, sin dependencias de infraestructura.

---

### 5.4 Dominio principal – Ventas (MVP)

**Entidades:**

- **Pedido**: ID, Cliente (FK a Party rol Cliente), Fecha, Estado, Total
- **Línea de pedido**: ID, Pedido (FK), Producto (FK), Variante (FK), Cantidad, Precio unitario, Subtotal
- **Estado de pedido**: Borrador, Confirmado, En preparación, Entregado, Cancelado
- **Documentación mercantil**: Presupuesto, Pedido, Albarán, Factura. Enlazan conceptualmente con los estados del pedido.

**Funciones:**

- Crear/editar pedidos/documentos de venta
- Añadir/modificar/eliminar líneas de pedido/documentos de venta
- Calcular precio usando motor de tarificación
- Cambiar estado de pedido
- Visualizar historial de pedidos
- Generación de documentos físicos (presupuestos, albaranes, facturas proforma) **delegada a frontend vía Web-to-Print**

**Dependencias:**

- **Requiere:** Party (cliente), Producto/Variante, Tarificación
- **Provee servicios a:** MES (para pedidos personalizados)

**Adaptadores:**

- Repositorio PostgreSQL para Pedido, Línea de pedido
- **Frontend como adaptador de generación documental** (Web-to-Print)

**Implementación:** Fase 2 (ADR-007)

---

### 5.5 Subdominio MES – Producción Personalizada (MVP)

**Entidades:**

- **Pedido personalizado** (extensión de Pedido con requisitos de producción)
- **Estado de producción**: Diseño, Aprobación de diseño, Marcaje, Taller, Control de calidad, Listo para entrega
- **Trabajo de taller**: ID, Pedido personalizado (FK), Operario, Observaciones
- **Diseño/Archivo asociado**: ID, Pedido personalizado (FK), Ruta NAS, Tipo archivo, Tamaño

**Funciones:**

- Gestionar pedidos personalizados
- Registrar estados de producción
- Terminal de taller (interfaz simplificada para tablets)
- Adjuntar/visualizar diseños
- Trazabilidad de trabajos y cambios de estado
- Integración con inventario para consumos de materiales (Post-MVP)

**Dependencias:**

- **Requiere:** Ventas (pedido base), Producto/Variante, Gestión documental
- **Provee servicios a:** Ninguno (es subdominio final)

**Adaptadores:**

- Repositorio PostgreSQL para estados MES y trabajos
- Adaptador NAS para almacenamiento de diseños (indexación en PostgreSQL)

**Implementación:** Fase 3 (ADR-007)

---

### 5.6 Dominio principal – Compras (Post-MVP)

**Rol en MVP:** NO existe como dominio independiente. Proveedores existen **solo como fuente de costes** dentro de Party.

**Rol Post-MVP:** Cierre del ciclo económico, entrada formal de mercancía, impacto en stock real y contabilidad futura.

**Entidades previstas (Post-MVP):**

- Pedido de compra
- Línea de pedido de compra
- Recepción de mercancía
- Albarán y Factura de proveedor (no contable en MVP)

**Funciones Post-MVP:**

- Creación de pedidos de compra a proveedores
- Recepción total/parcial de mercancía
- Actualización de stock desde recepciones
- Trazabilidad proveedor → producto → coste
- Base para futuras integraciones contables

**Dependencias previstas:**

- **Requiere:** Party (proveedor), Producto/Variante, Inventario/Stock

---

### 5.7 Dominio principal – Inventario / Stock (Post-MVP)

**Rol en MVP:** NO existe. No hay control de stock en MVP.

**Entidades previstas (Post-MVP):**

- Producto (referencia a Producto/Variante)
- Stock físico
- Lote
- Ubicación
- Movimiento de stock
- Regularización

**Funciones Post-MVP:**

- Control de stock físico por producto/variante
- Gestión de lotes y ubicaciones
- Reservas automáticas desde pedidos
- Regularizaciones de inventario
- Alertas de stock mínimo
- Valoración de inventario

**Dependencias previstas:**

- **Requiere:** Producto/Variante, Compras (para entradas), Ventas (para salidas)

---

### 5.8 Módulos Transversales

#### Seguridad (MVP: básica; Post-MVP: avanzada)

**MVP:**

- Autenticación JWT
- Roles básicos: Admin, Comercial, Diseño, Taller
- Control de acceso por rol (RBAC básico)
- Hash de passwords

**Post-MVP:**

- RBAC avanzado con permisos granulares
- Roles personalizados
- Gestión de sesiones avanzada
- Autenticación multifactor (opcional)

**Implementación MVP:** Fase 0 (ADR-007)

---

#### Auditoría (MVP: mínima; Post-MVP: completa)

**MVP:**

- Log de cambios críticos en tarificación (precios, márgenes, descuentos)
- Log de cambios en estados de producción (MES)
- Registro de quién, cuándo, qué cambió

**Post-MVP:**

- Trazabilidad completa de todos los cambios
- Retención configurable de datos
- Informes de auditoría
- Cumplimiento normativo

**Implementación MVP:** Fase 1 (ADR-007) - auditoría mínima

---

#### Gestión Documental (MVP: mínima; Post-MVP: avanzada)

**MVP:**

- Almacenamiento de diseños en NAS
- Indexación en PostgreSQL (ID, ruta, tipo, tamaño, fecha)
- Adjuntar/visualizar archivos asociados a pedidos personalizados
- Trazabilidad básica de archivos

**Post-MVP:**

- Versionado de documentos
- Control de acceso granular
- Búsqueda avanzada
- OCR y extracción de metadatos
- Integración con sistemas externos

**Implementación MVP:** Fase 3 (ADR-007)

---

#### i18n - Internacionalización (MVP: frontend básico; Post-MVP: backend completo)

**MVP:**

- Etiquetas estáticas en frontend (Vue-i18n)
- Idiomas: Español (por defecto), posible catalán/valenciano
- Mensajes de interfaz traducibles
- **NO incluye backend i18n** (mensajes de error, validaciones en español)

**Post-MVP:**

- i18n completa backend (go-i18n)
- Mensajes de error, validaciones, notificaciones traducibles
- Soporte multiidioma completo
- Gestión de traducciones desde interfaz

**Implementación MVP:** Integrado en frontend desde Fase 1

---

### 5.9 Diagrama Conceptual: Bounded Contexts

```
┌─────────────────────────────────────────────────────────────────┐
│                        MONOLITO MODULAR                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐         ┌─────────────────┐              │
│  │  Party /        │────────▶│  Producto /     │              │
│  │  Organización   │         │  Variante /     │              │
│  │  (Fundacional)  │         │  Categoría      │              │
│  │                 │         │  (Fundacional)  │              │
│  │ - Clientes      │         │                 │              │
│  │ - Proveedores   │         │ - Productos     │              │
│  │ - Jerarquías    │         │ - Variantes     │              │
│  │ - Costes base   │         │ - Categorías    │              │
│  └─────────────────┘         └─────────────────┘              │
│         │                            │                         │
│         │                            │                         │
│         ▼                            ▼                         │
│  ┌──────────────────────────────────────────┐                 │
│  │        Tarificación                      │                 │
│  │        (Núcleo Económico)                │                 │
│  │                                          │                 │
│  │  Coste base + Margen - Descuento        │                 │
│  │  = Precio final de venta                │                 │
│  └──────────────────────────────────────────┘                 │
│                    │                                           │
│                    ▼                                           │
│  ┌──────────────────────────────────────────┐                 │
│  │        Ventas                            │                 │
│  │        (Pedidos Estándar)                │                 │
│  │                                          │                 │
│  │  Cliente → Productos → Precio → Pedido  │                 │
│  └──────────────────────────────────────────┘                 │
│                    │                                           │
│                    │                                           │
│                    ▼                                           │
│  ┌──────────────────────────────────────────┐                 │
│  │        MES                               │                 │
│  │        (Producción Personalizada)        │                 │
│  │                                          │                 │
│  │  Pedido → Diseño → Producción → Entrega │                 │
│  └──────────────────────────────────────────┘                 │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│                    MÓDULOS TRANSVERSALES                        │
│  Seguridad │ Auditoría │ Gestión Documental │ i18n (frontend) │
└─────────────────────────────────────────────────────────────────┘

POST-MVP (No implementado en proyecto actual):
┌──────────────────┐      ┌──────────────────┐
│    Compras       │      │  Inventario /    │
│    (Formal)      │─────▶│  Stock           │
└──────────────────┘      └──────────────────┘
```

**Leyenda:**

- **→** : Dependencia de dominio (usa servicios de)
- **Fundacional**: Debe implementarse primero (Fase 1)
- **Núcleo Económico**: Criticidad máxima (Fase 1)
- **Transversales**: Aplican a todos los dominios

---

### 5.10 Diagrama Conceptual: Flujo de Negocio MVP

```
PEDIDO ESTÁNDAR:
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Cliente  │───▶│ Selección│───▶│ Cálculo  │───▶│  Pedido  │
│ (Party)  │    │ Producto/│    │  Precio  │    │Confirmado│
│          │    │ Variante │    │(Tarif.)  │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
                                                       │
                                                       ▼
                                                ┌──────────┐
                                                │  Entrega │
                                                └──────────┘

PEDIDO PERSONALIZADO:
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Cliente  │───▶│ Selección│───▶│ Cálculo  │───▶│  Pedido  │
│ (Party)  │    │ Producto/│    │  Precio  │    │Confirmado│
│          │    │ Variante │    │(Tarif.)  │    │+ Diseño  │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
                                                       │
                                                       ▼
                                                ┌──────────┐
                                                │   MES    │
                                                │ (Fases)  │
                                                └──────────┘
                                                       │
                        ┌──────────┬──────────┬────────┴───┬──────────┐
                        ▼          ▼          ▼            ▼          ▼
                    ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
                    │ Diseño │▶│Aprob.  │▶│Marcaje │▶│ Taller │▶│Control │
                    │        │ │Diseño  │ │        │ │        │ │Calidad │
                    └────────┘ └────────┘ └────────┘ └────────┘ └────────┘
                                                                       │
                                                                       ▼
                                                                ┌──────────┐
                                                                │  Entrega │
                                                                └──────────┘
```

---

## 6. Matriz de Trazabilidad MVP

|Caso de Uso|Requisitos Funcionales|Bounded Contexts|Fase Implementación|
|---|---|---|---|
|Crear cliente/proveedor|RF1|Party/Organización|Fase 1|
|Gestionar jerarquía empresarial|RF1|Party/Organización|Fase 1|
|Registrar producto/variante|RF2|Producto/Variante/Categoría|Fase 1|
|Calcular precio de venta|RF3|Tarificación + Party + Producto|Fase 1|
|Crear pedido estándar|RF4|Ventas + Tarificación + Party + Producto|Fase 2|
|Registrar coste de proveedor|RF5|Party/Organización (Supplier)|Fase 1|
|Gestionar pedido personalizado|RF6|MES + Ventas|Fase 3|
|Operar terminal de taller|RF7|MES|Fase 3|
|Adjuntar diseño a pedido|RF8|MES + Gestión Documental|Fase 3|
|Autenticar usuario|RF9|Seguridad|Fase 0|
|Auditar cambios en precios|RF10|Auditoría + Tarificación|Fase 1|

**Referencias ADR:** ADR-007, ADR-006, ADR-005, ADR-004

---

## 7. Orden de Implementación y Fases del MVP

El desarrollo sigue una **estrategia incremental por fases**, con **frontend paralelo** en todas ellas, garantizando que al final de cada fase los módulos implementados sean operativos.

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

- **Party / Organización** (completo):
    - CRUD clientes/proveedores
    - Jerarquías empresariales
    - Herencia de descuentos
    - Costes de proveedor por producto/variante
- **Producto / Variante / Categoría** (completo):
    - CRUD productos, variantes, categorías
    - Modificadores de precio en variantes
- **Tarificación** (completo):
    - Motor de cálculo de precios
    - Aplicación de márgenes y descuentos
    - Integración con costes y descuentos
- **Auditoría mínima**: Log de cambios en tarificación
- **Frontend**:
    - CRUD Party completo
    - CRUD Producto completo
    - Calculadora de tarificación funcional

**Criterios de aceptación:**

- ✅ CRUD Party funcional (frontend + backend)
- ✅ CRUD Producto funcional
- ✅ Cálculo de tarificación con datos reales
- ✅ Cobertura tests ≥80% en dominio tarificación
- ✅ Jerarquías funcionan correctamente

**Duración estimada:** 4-6 semanas

---

### Fase 2 – Casos de Uso Core y Pedidos

**Objetivo:** Flujo de ventas completo.

**Entregables:**

- **Pedido Estándar** (completo):
    - Creación, edición, seguimiento
    - Cálculo automático de precios
    - Gestión de estados
- **Casos de uso de ventas**:
    - Flujo completo: cotización → confirmación → entrega
    - Validación de crédito
- **Frontend**:
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

- **MES** (completo):
    - Pedidos personalizados
    - Estados de producción
    - Trazabilidad de procesos
- **Terminal de taller**:
    - Interfaz tablet funcional
    - Gestión de trabajos
    - Cambio de estados
- **Gestión documental**:
    - Almacenamiento NAS
    - Adjuntar/visualizar diseños
- **Frontend**:
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

### Total Estimado MVP: 11-16 semanas

**Nota:** Estas estimaciones son orientativas. La duración real dependerá de:

- Disponibilidad del equipo
- Complejidad de casos de negocio específicos
- Descubrimiento de requisitos adicionales
- Iteraciones de feedback con usuarios

**Referencia:** ADR-007

---

## 8. Resumen General

### Alcance MVP claramente delimitado

**Lo que SÍ incluye el MVP:**

- Party / Organización con jerarquías de clientes
- Producto / Variante / Categoría
- Motor de tarificación funcional (núcleo económico)
- Pedidos estándar y personalizados
- MES básico para producción personalizada
- Terminal de taller
- Gestión documental mínima (NAS)
- Seguridad básica (JWT, RBAC)
- Auditoría mínima (cambios críticos)

**Lo que NO incluye el MVP (Post-MVP):**

- Compras formales
- Control de stock
- i18n backend
- Auditoría completa
- Tarificación avanzada
- Multimoneda
- Escalabilidad distribuida

---

### Arquitectura y Principios

- **Tipo de aplicación:** Monolito modular local-first
- **Arquitectura:** Clean Architecture + DDD con rigor asimétrico (CRUD y casos simples con flexibilidad controlada)
- **Dominio:** protegido, completamente testeable, independiente de frameworks
- **Stack:** Go + PostgreSQL + Vue.js + Docker
- **Generación de documentos:** Delegada a frontend (Web-to-Print)
- **Desarrollo:** Incremental por fases, frontend paralelo
- **Testing:** TDD obligatorio en dominio, cobertura ≥75% global, ≥80% tarificación

---

### Modelo de Dominio

**Party como concepto fundacional:**

- Un Party, múltiples roles (Cliente, Proveedor, o ambos)
- Jerarquías solo para clientes
- Herencia de descuentos en jerarquías
- Proveedores como fuente de costes (no gestión de compras)

**Tarificación como núcleo económico:**

- Motor independiente y testeable
- Construido sobre Party (costes, descuentos) y Producto (variantes)
- Cálculo on-demand, no persistido

**MES como subdominio especializado:**

- Aislado, preparado para extraer si es necesario
- Estados de producción claros
- Terminal simplificado para taller

---

### Estrategia de Implementación

**Orden obligatorio (ADR-007):**

1. Fase 0: Fundaciones (infraestructura)
2. Fase 1: Party + Producto + Tarificación
3. Fase 2: Pedidos
4. Fase 3: MES

**Frontend paralelo:** En todas las fases, garantizando operatividad progresiva.

**Criterios de éxito:** Cada fase tiene criterios de aceptación verificables.

---

### Preparación para el Futuro

- **Módulos extraíbles:** Cada Bounded Context puede convertirse en servicio independiente
- **Contratos explícitos:** Comunicación entre módulos mediante interfaces claras
- **Infraestructura sustituible:** Adaptadores intercambiables
- **Dominio protegido:** Lógica de negocio aislada de decisiones técnicas
- **Evolución controlada:** Post-MVP = nuevo proyecto independiente (ADR-004)

---

### Referencias ADR Completas

- **ADR-001:** Stack tecnológico base (Go, PostgreSQL, Vue.js, Docker, Web-to-Print)
- **ADR-002:** Clean Architecture + DDD con rigor asimétrico
- **ADR-003:** Monolito modular local-first
- **ADR-004:** Ciclo de vida hasta MVP (proyecto independiente post-MVP)
- **ADR-005**: Gestión Unificada de Clientes y Proveedores (Party / Organización)
- **ADR-006:** Estrategia de desarrollo dirigida por dominio
- **ADR-007:** Orden de implementación de módulos y frontend paralelo

---

### Criterios de Finalización del Proyecto

El proyecto se considera **completo y exitoso** cuando:

1. **MVP operativo en producción:**
    
    - Flujo completo de pedidos estándar funcional
    - Flujo completo de pedidos personalizados funcional
    - Terminal de taller en uso real
    
2. **Calidad técnica verificada:**
    
    - Cobertura de tests ≥75% global
    - Cobertura de tests ≥80% en tarificación
    - Sistema estable durante 2 semanas de uso real

3. **Infraestructura estable:**
    
    - Despliegue Docker funcional
    - Backups automáticos configurados
    - Procedimientos de recuperación documentados

4. **Usuarios operando el sistema:**
    
    - Comerciales gestionan pedidos sin soporte
    - Taller usa terminal sin incidencias
    - Administrador gestiona catálogo y tarifas

5. **Documentación completa:**
    
    - Documentación técnica actualizada
    - Manual de usuario básico
    - Procedimientos de mantenimiento

**Cualquier funcionalidad Post-MVP se tratará como un nuevo proyecto independiente, preservando la estabilidad del MVP existente.**

---

## Apéndice: Glosario de Términos

- **Party:** Patrón de modelado que representa a cualquier persona u organización con relación al sistema
- **Bounded Context:** Límite explícito de un módulo/dominio con su propio modelo
- **Clean Architecture:** Arquitectura por capas concéntricas con dependencias hacia el dominio
- **DDD (Domain-Driven Design):** Diseño dirigido por el dominio de negocio
- **Rigor asimétrico:** Aplicación de disciplina arquitectónica proporcional al valor estratégico
- **MES (Manufacturing Execution System):** Sistema de ejecución de manufactura
- **Web-to-Print:** Generación de documentos PDF desde el navegador usando HTML/CSS
- **TDD (Test-Driven Development):** Desarrollo guiado por pruebas
- **RBAC (Role-Based Access Control):** Control de acceso basado en roles
- **MVP (Minimum Viable Product):** Producto mínimo viable
- **NAS (Network Attached Storage):** Almacenamiento conectado a red

---

**Fin del Documento – Versión 3.0**
