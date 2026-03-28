# TramaTex – Glosario Unificado

**Versión:** 2.2
**Fecha:** 17/03/2026
**Propósito:** Definición centralizada de todos los términos y conceptos clave del proyecto, tanto técnicos como de proceso.

---

## 📚 Términos de Dominio y Técnicos

### A

**Adapter / Adaptador**
Componente en la capa de **Infraestructura** de Clean Architecture que traduce entre el dominio y sistemas externos (BD, API, UI). Implementa las interfaces (puertos) definidas en el dominio.

**Address (Dirección - Entidad)**
Entidad de dominio vinculada a una `Party` que representa una ubicación física (calle, ciudad, provincia, código postal, país). Una Party puede tener múltiples direcciones, identificando una como `IsPrimary`.

**Application Layer / Capa de Aplicación**
Capa que orquesta los casos de uso del sistema. No contiene lógica de negocio (esta reside en el Dominio), sino que coordina la interacción entre entidades, servicios de dominio y repositorios.

**Auditoría**
Registro de cambios en el sistema (quién, qué, cuándo). Por mandato arquitectónico, los campos de auditoría (`CreatedAt`, `UpdatedAt`, etc.) se gestionan exclusivamente en la capa de persistencia/infraestructura, quedando excluidos de las entidades de dominio para mantener la pureza del modelo.

**Attribute (Atributo)**
Entidad de dominio para una característica configurable de producto (ej. "Talla", "Color") y sus valores. Define su alcance y hereda valores a productos.

**AttributeValue (Valor de Atributo)**
Entidad de dominio que representa un valor específico de un `Attribute` (ej. "Large" para "Talla", "Rojo" para "Color"). Puede contener modificadores de precio (`FIXED` o `PERCENTAGE`).

---

### B

**BaseCost (Coste Base de Variante)**
Valor calculado dinámicamente para una `ProductVariant`: `Product.BasePrice` + modificadores de los `AttributeValues` seleccionados. No se persiste en BD para mantener la flexibilidad JIT.

**BSP (Base Selling Price - Precio de Venta Base)**
Precio de una variante tras aplicar el margen comercial de la marca al `BaseCost`. Fórmula: `BaseCost * (1 + Brand.DefaultMarkupPercentage / 100)`.

**BaseSalesPriceRule**
Entidad de dominio que define el precio base de venta de un ProductVariant (coste/tarifa + incrementos). Es parte del motor de Pricing y sigue una jerarquía de precedencia por especificidad.

**Brand (Marca)**
Entidad de dominio que agrupa productos bajo una marca común. Es clave para el pricing (vía `DefaultMarkupPercentage`) y el alcance de los atributos.

**Bounded Context**
Límite explícito donde un modelo de dominio es válido. En TramaTex, módulos como Party, Product, Pricing, Sales y MES son Bounded Contexts separados pero conectados mediante interfaces de aplicación.

**Borrado Inteligente (Smart Deletion)**
Lógica de integridad referencial que impide eliminar una `Party` (contacto u operario) si existen documentos vinculados (ventas, trabajos MES) o relaciones activas en el sistema.

---

### C

**Clean Architecture**
Arquitectura por capas concéntricas (Domain, Application, Infrastructure, Interfaces) que coloca el dominio de negocio en el centro, protegido de dependencias externas.

**ContactDetails (Detalles de Contacto)**
Entidad de dominio vinculada a una organización (`OrganizationProfile`) que representa a una persona física de contacto, con su email, teléfono y cargo.

---

### D

**DDD (Domain-Driven Design)**
Metodología de diseño que prioriza el dominio de negocio. En TramaTex se aplica con **Rigor Asimétrico**: estricto en motores críticos (Pricing, Sales) y pragmático en CRUDs simples.

**DeliveryNote (Albarán)**
Documento de dominio que registra la entrega física de mercancía al cliente, derivado de un Pedido de Venta.

**DeliveryNoteNumber (Value Object)**
Objeto de Valor inmutable que representa el número único de un Albarán, encapsulando su formato y lógica de generación.

**Descuento**
Reducción de precio aplicada a un pedido. El motor de Pricing aplica una jerarquía: Override (manual) > Cliente (específico) > Global (campaña).

**Domain Layer / Capa de Dominio**
El núcleo del sistema. Contiene entidades, objetos de valor y servicios de dominio. Es "Pure Go", sin dependencias de frameworks o bases de datos.

---

### E

**Entidad (en DDD)**
Objeto del dominio con identidad única (ID), cuyo estado puede evolucionar. Ejemplo: `Product`, `Quote`.

**Especificación**
Documento que describe requisitos funcionales y reglas de negocio. En la fase final, se transforman en "Guías de Comportamiento" simplificadas.

---

### F

**FSP (Final Selling Price - Precio de Venta Final)**
Precio final ofrecido al cliente tras aplicar descuentos al `BSP`. Fórmula: `BSP * (1 - DescuentoAplicable / 100)`.

---

### I

**i18n (Internacionalización)**
Soporte multi-idioma. MVP: interfaz en español/catalán; estados y prioridades técnicos en inglés en la API para estandarización, pero **siempre mostrados en castellano** en la UI (ver `code-and-style-standards.md`, sección "Interfaz de Usuario").

**Infrastructure Layer / Capa de Infraestructura**
Capa que contiene las implementaciones técnicas: persistencia (GORM), clientes de APIs externas y adaptadores.

**Interfaces Layer / Capa de Interfaces**
Punto de entrada al sistema (Handlers HTTP, Controllers). Traduce peticiones externas a llamadas de la Capa de Aplicación.

**Invoice (Factura)**
Documento de dominio que representa la solicitud legal de pago al cliente, consolidando líneas de venta.

**InvoiceNumber (Value Object)**
Objeto de Valor inmutable que representa el número único de una Factura, encapsulando su formato y lógica de generación.

---

### J

**JIT (Just-In-Time) Variants**
Estrategia donde las variantes de producto no se pre-generan, sino que se crean dinámicamente según la combinación de atributos seleccionada en el momento de la venta.

**JWT (JSON Web Token)**
Estándar de autenticación stateless para validar usuarios en las peticiones a la API.

---

### M

**MES (Manufacturing Execution System)**
Sistema de ejecución de manufactura que gestiona la producción en taller de artículos personalizados.

**Task (Tarea)**
Proceso atómico e indivisible del flujo MES (ej. Diseñar, Imprimir, Marcar). Dato maestro reutilizable.

**Position (Posición)**
Zona de la prenda donde se realiza un trabajo (ej. Pecho izquierdo, Espalda). Dato maestro transversal.

**WorkType (Tipo de Trabajo)**
Secuencia ordenada de tareas que define un tipo de marcado/personalización. Es la "receta" del proceso. Backend: tabla `work_types`.

**WorkSetup (Configuración de Trabajo)**
Plantilla reutilizable que define la personalización de un tipo de prenda para un cliente. Combina WorkType + Position en líneas.

**WorkOrder (Orden de Trabajo)**
Instancia operativa de producción vinculada a un pedido. Realiza seguimiento de estados, tiempos y operarios en tiempo real. Backend: tabla `work_orders`.

**Modificador de precio**
Valor (`FIXED` o `PERCENTAGE`) asociado a un `AttributeValue` que afecta el cálculo del `BaseCost` de una variante.

**Moneda (Money - Value Object)**
Objeto de Valor que representa una cantidad monetaria. En el MVP, la divisa está fijada a Euro (€). Redondeo: utiliza round-half-up (redondeo comercial, `math.Floor(amount*100+0.5)/100`) para garantizar precisión en cálculos financieros, evitando banker's rounding que causaba discrepancias en descuentos.

**Monolito modular**
Estilo arquitectónico donde el sistema es un único binario pero está dividido internamente en módulos independientes (Bounded Contexts) con comunicación controlada.

---

### O

**OrganizationProfile**
Perfil de una `Party` que contiene datos corporativos (NIF, nombre fiscal). Una Party puede actuar simultáneamente como Cliente y Proveedor (Perfil Dual).

**OrderNumber (Value Object)**
Objeto de Valor inmutable que representa el número único de un Pedido de Venta, encapsulando su formato y lógica de generación.

---

### P

**Party (Tercero)**
Entidad base que representa a cualquier persona u organización que interactúa con el sistema. Se especializa mediante perfiles (`Person`, `Organization`) y roles.

**PartyServiceConfiguration**
Configuración específica de un `Party` para la prestación de servicios, permitiendo personalización por cliente/proveedor.

**Pedido (SalesOrder)**
Entidad de dominio que representa un compromiso de venta firme. Orquesta el flujo hacia Albarán y Factura.

**PersonProfile**
Perfil de una `Party` que contiene datos de personas físicas (nombre, apellidos, contacto directo).

**Porcentaje (Percentage - Value Object)**
Objeto de Valor inmutable que representa una proporción (ej. 0.21 para el 21% de IVA).

**Pricing Engine (Motor de Precios)**
Componente que calcula el `FSP` dinámicamente aplicando la jerarquía de reglas y modificadores de atributos.

**Product (Producto)**
Plantilla general de un artículo. Define el `BasePrice` y el `TaxRate` (IVA) aplicable.

**ProductVariant (Variante de Producto)**
Instancia final y vendible de un `Product`, resultante de una combinación específica de atributos.

---

### Q

**Quote (Presupuesto / Cotización)**
Oferta comercial inicial. No tiene impacto contable ni de producción hasta que se confirma y convierte en `SalesOrder`.

---

### T

**TaxID (Identificación Fiscal)**
Objeto de valor que encapsula y valida el identificador fiscal (NIF/CIF/NIE).

**TaxRate (Tipo impositivo)**
Porcentaje de IVA aplicable. Se almacena en el `Product` para asegurar la soberanía del catálogo sobre el precio.

**Terminal de taller**
Interfaz táctil para operarios que permite gestionar el avance de las `WorkOrder` (inicio, pausa, fin de tareas).

---

## 🎨 UI/UX y Sistema de Diseño

### B

**BaseEntityPage**
Plantilla maestra para la gestión integral de una Entidad de Dominio a pantalla completa. Implementa la arquitectura de tres capas para asegurar consistencia y enfoque operativo.

### C

**Context Header (Cabecera de Contexto)**
Bloque dinámico de la `BaseEntityPage` que hace scroll con el contenido. Agrupa la Toolbar de estado, la cinta de Summary y la trazabilidad (Related). Su fondo es Gris Ceniza (`#f9fafb`).

### I

**Identity Header (Cabecera de Identidad)**
Bloque superior fijo (`Sticky`) de la `BaseEntityPage`. Contiene la identidad del objeto (Título e ID) y las acciones globales. Su fondo es Blanco Puro (`#ffffff`).

### M

**Main Content (Área de Trabajo)**
Zona operativa central de la página donde residen las secciones de datos (`FormSection`) y el motor de líneas. Utiliza el fondo gris base de la aplicación para resaltar las tarjetas blancas.

**Material Symbols Outlined**
Estándar obligatorio de iconografía para TramaTex. Se integra mediante la clase CSS `.material-symbols-outlined` y utiliza nombres semánticos de iconos de Google.

### R

**Related Traceability (Consola de Trazabilidad)**
Sección dentro del `Context Header` que muestra enlaces rápidos a documentos vinculados (orígenes y destinos), permitiendo navegar por la genealogía del objeto sin perder contexto.

### S

**Summary Ribbon (Cinta de Resumen)**
Fila de tarjetas KPI situada en el `Context Header`. Proporciona lectura rápida de los 4 datos más críticos de la entidad (ej. Cliente, Fecha, Total).

---

## Abreviaturas Comunes

| Abreviatura | Significado |
|---|---|
| BSP | Base Selling Price (Precio Venta Base) |
| FSP | Final Selling Price (Precio Venta Final) |
| IAM | Identity and Access Management |
| JIT | Just-In-Time (Variantes bajo demanda) |
| MES | Manufacturing Execution System |
| VAT | Value Added Tax (IVA - Impuesto sobre el Valor Añadido) |

... (resto de abreviaturas iguales)
