# TramaTex – Glosario Unificado

**Versión:** 2.0
**Fecha:** 25/01/2026
**Propósito:** Definición centralizada de todos los términos y conceptos clave del proyecto, tanto técnicos como de proceso.

---

## 📚 Términos de Dominio y Técnicos

### A

**Adapter / Adaptador**
En Clean Architecture, componente que traduce entre el dominio y sistemas externos (base de datos, API, interfaz). Los adaptadores son intercambiables sin afectar la lógica de negocio.

**Auditoría**
Registro y seguimiento de cambios significativos en el sistema (quién, qué, cuándo). MVP incluye auditoría mínima de cambios en tarificación y MES.

**Attribute (Atributo)**
Entidad de dominio que gestiona una característica configurable de un producto (ej. "Talla", "Color") y sus posibles valores. Define su alcance de aplicación y hereda valores a productos.

**AttributeValue (Valor de Atributo)**
Entidad de dominio que representa un valor específico de un `Attribute` (ej. "Large" para "Talla", "Rojo" para "Color").

---

### B

**BaseSalesPriceCache**
Caché NoSQL (Redis) que almacena los precios base de venta (BaseSalesPrice) precalculados para todos los ProductVariants de un producto base. Optimizado para lectura rápida.

**BaseSalesPriceRule**
Entidad de dominio que define cómo se construye el precio base de venta de un ProductVariant (a partir de un coste/tarifa más incrementos) antes de aplicar cualquier modificación de venta. Se aplica en la fase de cálculo de caché, con una estricta precedencia por especificidad de producto.

**Brand (Marca)**
Entidad de dominio que agrupa productos bajo una marca común. Es clave para el pricing y el alcance de los atributos.

**Bounded Context**
Límite explícito en el que un modelo de dominio es válido. En TramaTex: Party, Producto, Tarificación, Ventas, MES son Bounded Contexts separados pero conectados.

---

### C

**Clean Architecture**
Arquitectura por capas concéntricas que coloca el dominio de negocio en el centro, protegido de dependencias externas (frameworks, BD, interfaz).

**CRU (Customer Relationship Management)**
Gestión de relaciones con clientes. En TramaTex, implementado parcialmente a través del módulo Party.

---

### D

**DDD (Domain-Driven Design)**
Metodología de diseño que prioriza la lógica de negocio (dominio) como estructura central del software.

**DeliveryNote (Albarán)**
Documento de dominio que registra la entrega física de mercancía al cliente, derivado de uno o varios Pedidos de Venta.

**DeliveryNoteNumber (Value Object)**
Objeto de Valor inmutable que representa el número único de un Albarán, encapsulando su formato y lógica de generación.

**Descuento**
Reducción de precio aplicada a un pedido. Tipos: base (cliente), heredado (jerarquía), específico (override puntual).

**Diseño (en MES)**
Etapa de producción personalizada donde se crea/aprueba el diseño del producto personalizado.

---

### E

**Entidad (en DDD)**
Objeto del dominio con identidad única, que puede cambiar de estado a lo largo del tiempo. Ejemplo: Pedido, Party.

**Especificación**
Documento que describe requisitos funcionales, reglas de negocio y contratos de un módulo.

---

### H

**Hardware objetivo**
Equipo físico en el que TramaTex debe funcionar: i3 8GB RAM (clientes), i3 16GB + SSD (servidor), tablets (taller).

---

### I

**i18n (Internacionalización)**
Soporte para múltiples idiomas. MVP: frontend en español/catalán. Post-MVP: tramatex-api completo.

**Invoice (Factura)**
Documento de dominio que representa la solicitud legal de pago al cliente, consolidando las ventas.

**InvoiceNumber (Value Object)**
Objeto de Valor inmutable que representa el número único de una Factura, encapsulando su formato y lógica de generación.

---

### J

**JWT (JSON Web Token)**
Estándar de autenticación stateless utilizado para validar usuarios. En TramaTex: tokens firmados con JWT_SECRET.

---

### L

**Local-first**
Arquitectura donde el sistema opera 100% en infraestructura local sin dependencias cloud obligatorias. Principio de TramaTex.

---

### M

**Margen**
Porcentaje añadido al coste base para calcular precio de venta. Regla de tarificación.

**MES (Manufacturing Execution System)**
Sistema de ejecución de manufactura. En TramaTex: gestión de producción personalizada con estados (Diseño, Aprobación, Marcaje, Taller, Control QA).

**Modificador de precio**
Valor (porcentaje o cantidad) asociado a una variante que afecta al cálculo del precio final.

**Moneda (Money - Value Object)**
Objeto de Valor que representa una cantidad monetaria con su divisa asociada. Para MVP, la divisa es siempre EUR.

**Monolito modular**
Aplicación única (no distribuida) con módulos internos claramente separados y reutilizables.

**MVP (Minimum Viable Product)**
El conjunto mínimo de funcionalidades necesarias para que el producto sea utilizable por los primeros usuarios y se pueda obtener feedback valioso.

---

### N

**NAS (Network Attached Storage)**
Almacenamiento conectado a red. En TramaTex: almacena diseños de pedidos personalizados.

---

### O

**Organización**
Estructura empresarial representada como "Party" en el sistema. Puede tener roles de Cliente y/o Proveedor.

**OrderNumber (Value Object)**
Objeto de Valor inmutable que representa el número único de un Pedido de Venta, encapsulando su formato y lógica de generación.

---

### P

**Party / Organización**
Patrón de modelado que representa cualquier persona u organización en el sistema. Evita duplicación Cliente/Proveedor.

**PartyServiceConfiguration**
Entidad de dominio que guarda configuraciones de servicios específicas de un `Party` para un producto de tipo `SERVICE`.

**Pedido (SalesOrder)**
Entidad de dominio que representa un compromiso de venta firme. Sirve como la base para la ejecución de la venta, la generación de Albaranes y Facturas.

**Porcentaje (Percentage - Value Object)**
Objeto de Valor inmutable que representa un porcentaje (ej. 0.10 para 10%).

**Persistencia**
Capa de almacenamiento de datos. En TramaTex: PostgreSQL + GORM (solo en capa de infraestructura).

**Product (Producto)**
Entidad de dominio que representa la plantilla o el concepto general de un artículo o servicio vendible.

**ProductGroup (Grupo de Producto)**
Entidad de dominio que representa una categoría jerárquica para productos.

**ProductType (Tipo de Producto)**
Enumeración que clasifica un `Product` como `TANGIBLE` (bien físico) o `SERVICE` (servicio).

**Proveed**
Entidad externa que proporciona productos/materiales al sistema. Representada como Party con rol Proveedor.

---

### Q

**Quote (Cotización / Presupuesto)**
Entidad de dominio que representa una oferta de precios a un cliente que aún no ha sido confirmada. Puede convertirse en un Pedido de Venta.

**QuoteNumber (Value Object)**
Objeto de Valor inmutable que representa el número único de una Cotización, encapsulando su formato y lógica de generación.

---

### R

**RBAC (Role-Based Access Control)**
Control de acceso basado en roles. MVP: Admin, Comercial, Diseño, Taller.

**Requisito Funcional (RF)**
Capacidad observable del sistema (p.e., "crear pedido", "calcular precio").

**Requisito No Funcional (RNF)**
Atributo de calidad (p.e., "operar con <150MB RAM", "ACID en BD").

**Rigor asimétrico**
Aplicación de disciplina arquitectónica proporcional al valor estratégico. En TramaTex: rigor máximo en Tarificación, flexibilidad controlada en CRUDs simples.

**RuleValue (Value Object)**
Objeto de Valor que encapsula el tipo y el valor del efecto de una regla de pricing (ej. porcentaje de incremento, cantidad monetaria fija, etc.). Implícitamente define si la operación es modificativa o anulativa.

---

### S

**SaleModificationRule**
Entidad de dominio que define cómo se modifica el precio base de venta (ya calculado) en el contexto de una venta específica (basado en cliente, monto mínimo, grupo de productos o fechas de campaña). Se aplica en el momento de la venta.

**SalesOrder (Pedido)**
Ver **Pedido**.

**Supplier (Proveedor)**
Especialización de Party con rol Proveedor. Proporciona costes base para productos.

---

### T

**Tarif de tarificación**
Conjunto de reglas para calcular precio de venta a partir de coste base, margen y descuentos.

**Terminal de taller**
Interfaz simplificada para tablets que operan en el taller. Gestiona pedidos personalizados en MES.

**TDD (Test-Driven Development)**
Desarrollo guiado por pruebas. Obligatorio en dominio crítico de TramaTex.

---

### V

**ProductVariant (Variante de Producto)**
Entidad de dominio que representa la instancia final y vendible de un `Product`, compuesta por una combinación única de `AttributeValue`s.

**ProductVariantStatus (Estado de Variante de Producto)**
Enumeración que indica el estado de una `ProductVariant`: `PROVISIONAL` (creada JIT) o `CONFIRMED` (validada/manual).

**Value Object (Objeto de Valor)**
Objeto del dominio sin identidad única, inmutable. Ejemplo: Email, Password (en Auth).

---

### W

**Web-to-Print**
Generación de documentos (PDF, impresión) desde HTML/CSS en el navegador. En TramaTex: delegado al frontend para presupuestos, albaranes, facturas.

---

## 📋 Términos de Planificación y Proceso

### Jerarquía de Planificación

La planificación del proyecto se estructura en tres niveles jerárquicos:

#### 1. Fase (Phase)

- **Descripción:** Una fase es un período de tiempo largo (varios meses) que agrupa un conjunto de épicas relacionadas que contribuyen a un objetivo estratégico general del MVP.
- **Ejemplo:** `Fase 1: Dominio Base para Tarificación`

#### 2. Épica (Epic)

- **Descripción:** Una épica representa una funcionalidad o capacidad de negocio completa y autocontenida. Es un bloque de trabajo grande que se puede entregar de forma independiente y que aporta valor al usuario. Cada épica se desglosa en varios sprints de implementación.
- **Ejemplo:** `Épica 1: Party (Gestión de Clientes y Proveedores)`

#### 3. Sprint de Implementación (Implementation Sprint)

- **Descripción:** Un sprint de implementación es un paquete de trabajo técnico, centrado en la implementación de una capa específica de la arquitectura (p. ej., Dominio, Persistencia, API, UI) para una épica determinada. No tiene una duración fija, sino que se completa cuando los objetivos técnicos de esa capa están terminados y probados.
- **Ejemplo:** `Sprint 1: Domain Layer (para la Épica de Party)`

### Otros Términos de Proceso

- **Bitácora (Journal):** Un documento que registra el trabajo realizado, las decisiones tomadas y los problemas encontrados durante un período de desarrollo. Reemplaza el concepto de "sesión de desarrollo".
- **ADR (Architecture Decision Record):** Un documento que captura una decisión arquitectónica importante junto con su contexto y consecuencias.

---

## Abreviaturas Comunes

| Abreviatura | Significado |
|---|---|
| ACID | Atomicity, Consistency, Isolation, Durability |
| ADR | Architecture Decision Record |
| DDD | Domain-Driven Design |
| i18n | Internacionalización |
| JWT | JSON Web Token |
| MES | Manufacturing Execution System |
| MVP | Minimum Viable Product |
| NAS | Network Attached Storage |
| RBAC | Role-Based Access Control |
| RF | Requisito Funcional |
| RNF | Requisito No Funcional |
| TDD | Test-Driven Development |

---
