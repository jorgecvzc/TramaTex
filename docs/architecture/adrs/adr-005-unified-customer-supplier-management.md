# ADR-005 – Gestión Unificada de Clientes y Proveedores (Party / Organización)

**Fecha:** 09/01/2026  
**Estado:** Adoptada  
**Autores:** Jorge Cortés Villalba, ChatGPT  

---

## 1. Contexto

En TramaTex, los módulos de **Compras** y **Ventas** necesitan referirse a **entidades externas**, que pueden ser clientes, proveedores o ambos. Además, existen casos de **UTEs y empresas matrices**, que requieren relaciones jerárquicas claras.  

**Problemas a resolver:**  
- Duplicación de datos si se mantienen entidades separadas para clientes y proveedores.  
- Complejidad en mantenimiento de información para entidades que son tanto clientes como proveedores.  
- Riesgo de inconsistencias en el histórico de transacciones, precios y contactos.  

Se busca una solución que unifique la gestión de estas entidades externas.

---

## 2. Alternativas Consideradas

### Alternativa A – Entidades Separadas
- **Descripción:** Mantener `Cliente` y `Proveedor` como entidades y tablas de base de datos independientes.
- **Ventajas:** 
  - Simplicidad inicial en los módulos de Compras y Ventas.
  - Roles claramente delimitados en el modelo.
- **Desventajas:** 
  - Alta duplicidad de datos (nombre, CIF, dirección).
  - Problemas de sincronización si una entidad cambia de datos.
  - Dificultad para soportar entidades con roles híbridos (cliente y proveedor a la vez).

### Alternativa B – Entidad Base Unificada (Patrón Party)
- **Descripción:** Crear una entidad base `Party` (o `Persona/Empresa`) que contenga los atributos comunes. Los roles específicos como `Cliente` o `Proveedor` se manejarían como roles o extensiones de esta entidad base.
- **Ventajas:** 
  - Evita la duplicidad de datos.
  - Mantiene un histórico unificado.
  - Soporta roles múltiples y híbridos de forma nativa.
  - Facilita el mantenimiento y la escalabilidad.
- **Desventajas:** 
  - Requiere lógica adicional en la aplicación para filtrar por rol.
  - La interfaz de usuario necesita manejar la presentación de campos específicos de cada rol.

---

## 3. Criterios de Decisión

- **Eliminación de Duplicidad de Datos:** La solución debe evitar que la misma empresa esté registrada varias veces.
- **Soporte de Roles Múltiples:** Debe ser posible que una entidad sea cliente, proveedor o ambos.
- **Mantenibilidad:** La solución debe ser fácil de mantener y extender en el futuro.
- **Consistencia de Datos:** El modelo debe garantizar la consistencia de la información de contacto y fiscal.

---

## 4. Decisión Adoptada

Se adopta la **Alternativa B: Entidad base unificada con roles múltiples (Patrón Party)**.

**Justificación:**
Esta alternativa resuelve directamente el problema principal de la duplicidad de datos y la gestión de entidades híbridas. Aunque introduce una ligera complejidad en la lógica de la aplicación, los beneficios en términos de consistencia de datos, mantenibilidad y escalabilidad a largo plazo superan con creces este coste inicial. Se alinea con el criterio de tener un modelo de dominio robusto y flexible.

---

## 5. Consecuencias

### Positivas
- Se elimina la duplicación de datos de clientes y proveedores.
- Se simplifica la gestión de entidades que cumplen ambos roles.
- Se crea un maestro de datos único para terceros, mejorando la integridad.

### Negativas
- La lógica de la aplicación y la interfaz de usuario deben ser capaces de manejar los diferentes roles y mostrar información condicional.
- La estructura de la base de datos es ligeramente más compleja, involucrando tablas de entidades, roles y relaciones.

---

## 6. Alcance

Este ADR aplica a:
- Modelado y persistencia de las entidades `Party` (Persona/Empresa).
- La integración de los módulos de Compras y Ventas con el módulo `Party`.
- El soporte para jerarquías empresariales (matrices/UTEs) y roles múltiples.
- La definición de la interfaz de usuario para la gestión unificada de estas entidades.

Cualquier cambio estructural en cómo se gestionan clientes y proveedores requerirá un nuevo ADR.

---

## 7. Integración con otros ADRs

- **ADR-006 (Desarrollo Dirigido por Dominio):** La entidad `Party` se convierte en un agregado fundamental en nuestro dominio.
- **ADR-007 (Orden de Implementación):** La implementación del módulo `Party` es una de las primeras prioridades en la Fase 1, ya que es una dependencia para los módulos de `Producto` y `Tarificación`.

---

## 8. Notas Adicionales / Consideraciones Especiales

### Diagrama Conceptual Simplificado

```md
[Persona/Empresa]──<roles>──[Cliente]  
					 │  
					 └──[Proveedor]  
[Persona/Empresa]──<matriz_rel>──[Empresa Matriz]  
```
*Este diagrama muestra cómo una entidad `Persona/Empresa` puede tener roles de Cliente y/o Proveedor, y también cómo puede estar relacionada con una empresa matriz.*

---

## 9. Referencias
