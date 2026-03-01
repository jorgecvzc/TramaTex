# Módulo de Party (Gestión de Terceros)

Este módulo es fundamental para TramaTex, ya que gestiona la información de todos los terceros (clientes, proveedores, empleados) que interactúan con el sistema. Proporciona una visión unificada de las Party, sus roles, perfiles (persona/organización), contactos y relaciones.

## Diseño Arquitectónico

Para una descripción detallada de las decisiones arquitectónicas, entidades de dominio, objetos de valor, casos de uso y estrategia de gestión de Party de este módulo, consulte el siguiente Architectural Decision Record (ADR):

*   [ADR-012: Arquitectura del Módulo Party](../../architecture/adrs/adr-012-party-module-architecture.md)

## Componentes Clave

*   **Entidades de Dominio:**
    *   `Party`: La raíz del agregado, representando una persona, una organización, o ambas.
    *   `PersonProfile`: Datos específicos de una persona.
    *   `OrganizationProfile`: Datos específicos de una organización.
    *   `ContactDetails`: Detalles de contacto asociados a una organización.
    *   `PartyRole`: Roles que una Party puede asumir (Cliente, Proveedor, Empleado).
    *   `PartyRelationship`: Relaciones tipadas entre Party (ej. "es empleado de", "es filial de").

*   **Objetos de Valor:**
    *   `PartyID`, `ContactDetailsID`, `Email`, `Phone`, `TaxID`.
    *   Enumeraciones: `PartyStatus`, `PartyRoleType`, `RelationshipType`.

*   **Casos de Uso (Capa de Aplicación):**
    *   Gestión completa de `Party` (Crear, Listar, Obtener, Actualizar, Cambiar Estado).
    *   Gestión de `PartyRole` (Añadir, Eliminar).
    *   Gestión de `PartyRelationship` (Crear, Listar, Eliminar).
    *   Gestión de `ContactDetails` (Añadir, Listar, Actualizar, Eliminar).

## Documentación Detallada

Consulte los siguientes documentos para una descripción más profunda del módulo de Party:

*   **Especificación del Módulo:** [module-spec.md](./module-spec.md)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md)
*   **Casos de Uso:** [use-cases.md](./use-cases.md)
*   **Contratos de API:** [api-contracts.md](./api-contracts.md)
*   **Diagramas Detallados del Dominio:** [diagrams/domain-model.md](./diagrams/domain-model.md)
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md)
*   **Resumen de Implementación (histórico):** [party-module-implementation-summary.md](../../archive/log/party-module-implementation-summary.md)
