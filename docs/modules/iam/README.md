# Módulo de IAM (Identity and Access Management)

Este módulo es fundamental para TramaTex, ya que gestiona la identidad de los usuarios, su autenticación y autorización en el sistema. Proporciona una base de seguridad robusta para toda la plataforma.

## Diseño Arquitectónico

Para una descripción detallada de las decisiones arquitectónicas, entidades de dominio, objetos de valor, casos de uso, y estrategia de seguridad de este módulo, consulte el siguiente Architectural Decision Record (ADR):

*   [ADR-014: Arquitectura del Módulo de IAM](../../architecture/adrs/adr-014-iam-module-architecture.md)

## Componentes Clave

*   **Entidades de Dominio:**
    *   `User`: Representa a un individuo que interactúa con el sistema.
    *   `Role`: Define los permisos de un usuario.
    *   `Permission`: Representa una acción específica (implícita en roles para MVP).

*   **Objetos de Valor:**
    *   `UserID`, `Email`, `HashedPassword`, `RoleID`, `RoleName`, `PermissionID`, `PermissionName`.

*   **Casos de Uso (Capa de Aplicación):**
    *   Gestión de Autenticación y Sesión (Registrar, Login, Refresh, Logout).
    *   Gestión de Usuarios (Listar, Crear, Asignar Rol, Eliminar - solo Admin).
    *   Autorización (Check Authorization).

## Documentación Detallada

Consulte los siguientes documentos para una descripción más profunda del módulo de IAM:

*   **Especificación del Módulo:** [module-spec.md](./module-spec.md)
*   **Modelo de Dominio:** [domain-model.md](./domain-model.md)
*   **Casos de Uso:** [use-cases.md](./use-cases.md)
*   **Contratos de API:** [api-contracts.md](./api-contracts.md)
*   **Guía de Implementación:** [implementation-guide.md](./implementation-guide.md)
