# Módulo de Party (Gestión de Terceros)

**Estado:** Aceptado

## 1. Propósito

*   **Visión del Módulo:** Gestionar la información de todos los terceros (clientes, proveedores, empleados) que interactúan con el sistema TramaTex, consolidando sus identidades, roles, perfiles (persona/organización), contactos y relaciones.
*   **Objetivos Clave:**
    *   Proporcionar un sistema centralizado y unificado para la gestión de clientes y proveedores, evitando la duplicación de datos.
    *   Soportar la complejidad de un mismo tercero asumiendo múltiples roles (ej. cliente y proveedor).
    *   Gestionar relaciones complejas entre terceros (ej. jerarquías de empresas, empleados de organizaciones).
    *   Proveer información de Party a otros módulos (Product, Pricing, Sales, IAM).

## 2. Requisitos

### 2.1. Requisitos Funcionales

*   **RF-P-001:** Crear y mantener Parties con perfil de persona, organización o ambos.
*   **RF-P-002:** Gestionar los roles de una Party (ej. Cliente, Proveedor, Empleado).
*   **RF-P-003:** Gestionar relaciones entre Parties (ej. "es empleado de", "es filial de").
*   **RF-P-004:** Gestionar puntos de contacto para perfiles de organización.
*   **RF-P-005:** Listar Parties con filtros por roles, tipo de perfil, estado y datos de perfil.
*   **RF-P-006:** Activar y desactivar Parties.

## 3. Casos de Uso

Para una lista completa y detallada de los casos de uso, incluyendo flujos y entradas/salidas, consulte el documento [Casos de Uso - Módulo Party](./use-cases.md).

## 4. Modelo de Dominio

Para una descripción detallada del modelo de dominio, incluyendo entidades, Value Objects, agregados y sus relaciones, consulte el documento [Modelo de Dominio - Módulo Party](./domain-model.md).

## 5. Decisiones de Diseño

*   **Patrón Party (Unificado):** Se adopta un modelo de `Party` unificado que puede tener un perfil de persona, de organización, o ambos. Esto elimina la duplicación de datos y permite múltiples roles por Party (`ADR-005`).
*   **Roles y Relaciones Explícitas:** Los roles (`PartyRole`) y las relaciones (`PartyRelationship`) son entidades de dominio explícitas que definen la función y el vínculo de una Party con otras.
*   **Contactos Simplificados:** El `OrganizationProfile` gestiona una lista de `ContactDetails` (Value Object), permitiendo flexibilidad sin sobre-ingeniería (`ADR-012`).
*   **Exclusión de Campos de Auditoría del Dominio:** `CreatedAt`, `UpdatedAt`, `CreatedBy`, `UpdatedBy` se gestionan en la capa de infraestructura/persistencia.
*   **Relaciones con Otros Módulos:**
    *   **IAM:** Una Party puede vincularse a un usuario de IAM.
    *   **Product:** Referenciado por `PartyServiceConfiguration` (en el módulo Product).
    *   **Pricing:** Consume `PartyID` y atributos relevantes del cliente para el cálculo de precios.
    *   **Sales:** Las órdenes y cotizaciones referencian `PartyID`.