# Especificación del Módulo: Party (Gestión de Terceros)

El módulo Party es el núcleo de identidad de TramaTex. Su propósito es proporcionar una fuente de verdad única para cualquier actor externo o interno, eliminando la necesidad de silos de datos para clientes, proveedores o empleados.

---

## 1. Filosofía de Identidad Unificada

A diferencia de sistemas tradicionales que separan físicamente a un "Cliente" de un "Proveedor", TramaTex utiliza el patrón **Party**. Esto permite que una entidad (ej. una empresa textil que nos compra botones pero nos vende uniformes) mantenga una única ficha técnica, dirección y contactos, asumiendo múltiples roles comerciales simultáneamente.

### Objetivos de Valor
- **Integridad de Datos:** Evitar que cambios en la dirección o razón social deban replicarse en múltiples tablas.
- **Trazabilidad 360°:** Permitir al sistema conocer la relación total con un tercero (cuánto le debemos, cuánto nos debe, qué produce para nosotros) desde un único punto.
- **Flexibilidad Operativa:** Facilitar la promoción de un prospecto a cliente o la contratación de un proveedor como empleado sin pérdida de historial.

---

## 2. Capacidades Estratégicas

### Gestión de Perfiles Adaptativos
El módulo distingue entre la naturaleza de la entidad (**Persona** vs. **Organización**). El comportamiento del sistema se adapta automáticamente:
- Si es organización, activa la gestión de puntos de contacto y datos fiscales corporativos.
- Si es persona, se enfoca en la identidad individual y roles operativos (ej. operario de taller).

### Arquitectura de Roles Dinámicos
Los roles no son definiciones estáticas, sino "capas" de comportamiento que se añaden a la identidad. Una Party puede nacer como un contacto genérico y evolucionar a Cliente o Proveedor según se generen documentos en los módulos de **Sales** o **Product**.

### Gobernanza de la Relación
El módulo permite modelar el mercado real mediante relaciones jerárquicas (Filial de, Empleado de), permitiendo que las reglas de negocio (como descuentos heredados) fluyan a través de la estructura organizativa del cliente.

---

## 3. Interacción con el Ecosistema

Party actúa como un **Proveedor de Contexto** para el resto de módulos:
- **Sales:** Provee la identidad fiscal y el descuento base para presupuestos.
- **MES:** Provee la identidad de los operarios para el reporte de tiempos.
- **Pricing:** Facilita la segmentación de reglas de precio basadas en la categoría de la Party.

---
**Última Actualización:** 2026-03-07
