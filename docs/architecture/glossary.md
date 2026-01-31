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

---

### B

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

---

### P

**Party / Organización**
Patrón de modelado que representa cualquier persona u organización en el sistema. Evita duplicación Cliente/Proveedor.   

**Pedido**
Solicitud de compra de productos. Tipos: estándar (producción no personalizada), personalizado (MES).

**Persistencia**
Capa de almacenamiento de datos. En TramaTex: PostgreSQL + GORM (solo en capa de infraestructura).

**Producto**
Artículo básico del catálogo. Puede tener variantes (talla, color, arreglos).

**Proveed**
Entidad externa que proporciona productos/materiales al sistema. Representada como Party con rol Proveedor.

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

---

### S

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

**Variante**
Especificidad de un producto (talla, color, arreglo). Puede afectar precio mediante modificador.

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
| MVP | Minimum Viable Product |
| MES | Manufacturing Execution System |
| DDD | Domain-Driven Design |
| RBAC | Role-Based Access Control |
| JWT | JSON Web Token |
| NAS | Network Attached Storage |
| TDD | Test-Driven Development |
| i18n | Internacionalización |
| RF | Requisito Funcional |
| RNF | Requisito No Funcional |
| ACID | Atomicity, Consistency, Isolation, Durability |

---
