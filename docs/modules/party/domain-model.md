# Lógica de Identidad y Dominio: Módulo Party

Este documento describe la estrategia de gestión de identidades en TramaTex, centrada en la flexibilidad de roles y la integridad operativa de los terceros.

---

## 1. El Concepto de Identidad Unificada (Party)

En TramaTex, una **Identidad (Party)** es un contenedor de información biográfica y legal que trasciende su función comercial momentánea. 

### Perfiles Duales
A diferencia de un modelo rígido, una identidad puede desplegar simultáneamente o de forma evolutiva dos perfiles:
- **Perfil Individual:** Datos de persona física, críticos para operarios (post-MVP) y contactos directos.
- **Perfil Corporativo:** Datos de persona jurídica, necesarios para la fiscalidad B2B.

**Regla de Negocio:** Una identidad es válida siempre que posea al menos uno de estos perfiles. El sistema permite la coexistencia para casos donde el dueño de una microempresa opera tanto a título personal como societario.

---

## 2. Dinámica de Roles y Relaciones

### El Rol como Capa Operativa
Los roles (Cliente, Proveedor, Empleado) son etiquetas que habilitan a la identidad en diferentes motores del sistema. 
- La presencia del rol `CLIENT` habilita el cálculo de precios específicos en el motor de **Pricing**.
- La eliminación de un rol no destruye la identidad, solo revoca su permiso para participar en nuevos flujos de ese tipo.

### Jerarquías de Mercado
Las relaciones permiten mapear la complejidad del mundo real. TramaTex no ve a los clientes como islas, sino como nodos que pueden estar conectados (ej. Filiales reportando a una Matriz), permitiendo una futura agregación de riesgos y beneficios a nivel de grupo.

---

## 3. Estrategia de Integridad (Borrado Inteligente)

La eliminación física de una identidad es una operación excepcional y protegida. El dominio impone una restricción de "Actividad Viva":

- **Protección Comercial:** Si existen documentos en **Sales** (desde presupuestos hasta facturas), la identidad queda bloqueada para borrado.
- **Protección Productiva:** Si hay tareas en **MES** vinculadas, se preserva la identidad para garantizar la trazabilidad de quién fabricó cada prenda.

**Directiva:** En caso de cese de relación, el sistema promueve el estado `INACTIVE` o `BLOCKED`, preservando el histórico pero impidiendo nuevas operaciones.

---

## 4. Gestión de Localización y Contacto

### Puntos de Comunicación
Las organizaciones gestionan sus contactos de forma interna. Un contacto puede ser una simple referencia o un vínculo a otra Identidad Individual, permitiendo navegar desde la empresa hacia la ficha personal del contacto si esta existe.

### Centralización de Direcciones
Se gestiona un catálogo de ubicaciones por identidad, donde la **Dirección Primaria** actúa como el punto de anclaje para la logística y la facturación automática, reduciendo la carga administrativa en la emisión de documentos.

---
**Última Versión:** 2026-03-07
