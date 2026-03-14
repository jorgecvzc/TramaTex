# Contratos de API - Módulo Party

Este documento especifica los puntos de entrada de la API para el módulo Party, enfocándose en el propósito operativo de cada endpoint y la estructura general de los datos.

---

## 1. Entidades Principales (`Parties`)

La gestión de `Parties` es el punto de entrada para cualquier tercero en el sistema. Los endpoints están diseñados para manejar la creación y actualización de perfiles duales (Persona/Organización) de forma atómica.

- **Creación y Actualización (`POST /parties`, `PUT /parties/{id}`):** Permite definir la identidad base, los roles iniciales y los perfiles de datos. El sistema valida que exista al menos un perfil válido.
- **Consulta y Listado (`GET /parties`, `GET /parties/{id}`):** Soporta filtrado por roles, tipo de entidad y estado. El listado está optimizado para la visualización en tablas de gestión.
- **Consulta Masiva (`GET /parties/batch?ids=...`):** Endpoint de alto rendimiento diseñado para resolver nombres de entidades en listados comerciales (Ventas/MES), eliminando el problema de múltiples peticiones (N+1).
- **Gestión de Estado (`PATCH /parties/{id}/status`):** Permite el bloqueo operativo de la entidad.
- **Borrado Inteligente (`DELETE /parties/{id}`):** Ejecuta las validaciones de integridad referencial operativa antes de permitir la eliminación física (ver reglas en Modelo de Dominio).

---

## 2. Roles y Relaciones

Define el "qué hace" y "con quién se vincula" una `Party`.

- **Roles (`POST /parties/{id}/roles`, `DELETE /parties/{id}/roles/{role}`):** Asignación dinámica de funciones de negocio (CLIENT, SUPPLIER, EMPLOYEE).
- **Relaciones (`POST /parties/{id}/relationships`, `GET /parties/{id}/relationships`):** Gestión de vínculos estructurales como jerarquías de empresas o pertenencia de empleados a organizaciones.

---

## 3. Direcciones y Contactos

Endpoints especializados para la gestión de datos de localización y comunicación.

- **Direcciones (`POST /parties/{id}/addresses`, `GET /parties/{id}/addresses`):** Gestión del catálogo de ubicaciones físicas. Incluye la lógica para marcar direcciones como primarias de envío o facturación.
- **Detalles de Contacto (`POST /parties/{id}/contact-details`):** Específico para perfiles de organización. Permite añadir puntos de contacto operativos, con la posibilidad de vincularlos a otras `Parties` de tipo Persona.

---

## 4. Configuraciones de Servicio

- **Configuraciones (`POST /parties/{id}/service-configurations`, `GET /parties/{id}/service-configurations`):** Almacena parámetros específicos (precios especiales, reglas técnicas) vinculados a una `Party` para servicios del catálogo (ej. Pricing Engine). Utiliza una estructura JSON flexible para adaptarse a diferentes tipos de servicios.

---

## Estructura de Respuesta Estándar (`PartyDTO`)

Las respuestas de la API devuelven un objeto consolidado que refleja la identidad actual de la entidad:

```json
{
  "id": "uuid",
  "status": "ACTIVE",
  "roles": ["CLIENT", "SUPPLIER"],
  "person_profile": { /* Datos de identidad individual si existen */ },
  "organization_profile": { 
    /* Datos corporativos y fiscales si existen */,
    "contacts": [ /* Lista de detalles de contacto */ ]
  },
  "created_at": "ISO-8601",
  "modified_at": "ISO-8601"
}
```

---
**Última Actualización:** 2026-03-07
