# ADR-005 – Gestión Unificada de Clientes y Proveedores (Party / Organización)

**Fecha:** 09/01/2026  
**Estado:** Propuesto  
**Decisión:** Adoptada  
**Autores:** Jorge Cortés Villalba, ChatGPT  

---

## 1. Contexto

En TramaTex, los módulos de **Compras** y **Ventas** necesitan referirse a **entidades externas**, que pueden ser clientes, proveedores o ambos. Además, existen casos de **UTEs y empresas matrices**, que requieren relaciones jerárquicas claras pero sin herencia de atributos.  

Problemas detectados:  
- Duplicación de datos si se mantienen entidades separadas para clientes y proveedores.  
- Complejidad en mantenimiento de información para entidades híbridas.  
- Riesgo de inconsistencias en histórico, precios y contactos.  

Se busca una solución que:  
- Evite duplicidad de datos.  
- Permita roles múltiples (cliente/proveedor) por entidad.  
- Soporte relaciones jerárquicas (UTEs y matrices).  
- Mantenga la mantenibilidad y escalabilidad del sistema.  

---

## 2. Decisiones consideradas

**Alternativa A – Entidades separadas:**  
- Mantener `Cliente` y `Proveedor` como entidades independientes.  
- Ventajas: Simplicidad de módulos existentes, roles claramente delimitados.  
- Desventajas: Duplicidad de datos, problemas de sincronización, mayor esfuerzo de mantenimiento, difícil soporte de entidades híbridas.  

**Alternativa B – Entidad base unificada con roles múltiples (decisión adoptada):**  
- Crear entidad base `Persona/Empresa` (`Party`).  
- Atributos comunes: nombre, CIF/NIF, dirección, teléfono, email.  
- Roles específicos:  
  - Cliente → condiciones de pago, histórico de pedidos  
  - Proveedor → precios base, condiciones de suministro, contrato  
- Entidades dependientes (UTEs) referencian a la empresa matriz pero **no heredan atributos**.  
- Integración con módulos: Compras y Ventas usan la misma entidad, filtrando por rol cuando sea necesario.  
- Ventajas: Evita duplicidad, mantiene histórico unificado, soporta roles múltiples, fácil mantenimiento.  
- Desventajas: Requiere lógica adicional para filtrar por rol y mostrar atributos específicos en la UI.  

---

## 3. Consecuencias técnicas

- La entidad `Persona/Empresa` será central y referenciada por los módulos:  
  - Ventas (clientes)  
  - Compras (proveedores)  
  - MES y otros módulos si es necesario  
- UI: Una ficha única por entidad, con filtros por rol (cliente/proveedor).  
- Base de datos: Tabla `entities` con relación `roles` (`cliente`, `proveedor`), tabla `matrix_relations` para UTEs y empresas matrices.  
- Integridad: Los cambios se registran de forma unificada (histórico común).  
- Backend: Los adaptadores de Compras y Ventas deberán manejar los filtros de rol y atributos específicos.  
- Frontend: Formulario unificado, campos condicionales según rol activo.  

---

## 4. Diagrama conceptual simplificado

```md
[Persona/Empresa]──<roles>──[Cliente]  
					 │  
					 └──[Proveedor]  
[Persona/Empresa]──<matriz_rel>──[Empresa Matriz]  
```

---

## 5. Alcance

Este ADR aplica a:

- Modelado y persistencia de entidades Party (Persona/Empresa)
- Integración con módulos de Compras y Ventas
- Soporte para jerarquías empresariales y roles múltiples
- UI unificada y filtrado por rol
- Reglas de herencia de descuentos (clientes)
- Cualquier cambio estructural o desviación requiere un nuevo ADR.
