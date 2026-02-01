# Modelo de Dominio - Módulo IAM

### Entidades Principales

#### User
- **Responsabilidad:** Representar a un individuo que interactúa con el sistema. Contiene información personal y credenciales.
- **Value Objects:** `UserID`, `Email`, `HashedPassword`, `Role`.
- **Reglas de Negocio:**
    - El email debe ser único.
    - La contraseña debe cumplir con los requisitos mínimos de seguridad antes de ser hasheada.

#### Role
- **Responsabilidad:** Representar un conjunto de permisos que pueden ser asignados a un usuario.
- **Value Objects:** `RoleID`, `RoleName`.
- **Reglas de Negocio:**
    - El nombre del rol debe ser único (e.g., 'admin', 'commercial', 'designer', 'workshop').

#### Permission
- **Responsabilidad:** Representar una acción específica que un usuario puede realizar.
- **Value Objects:** `PermissionID`, `PermissionName`.
