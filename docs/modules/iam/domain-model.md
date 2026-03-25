# Modelo de Dominio - Módulo IAM

### Entidades Principales

#### User
- **Responsabilidad:** Representar a un individuo que interactúa con el sistema. Contiene información personal y credenciales.
- **Estructura:**
    - `id`: `uuid.UUID` (Identificador único universal)
    - `email`: `*Email` (Value Object con validación de formato y unicidad)
    - `password`: `*Password` (Value Object que encapsula el hash)
    - `role`: `Role` (Enum que define las capacidades del usuario)
    - `active`: `bool` (Estado de habilitación del usuario)
- **Reglas de Negocio:**
    - El email debe ser único y válido.
    - La contraseña debe cumplir con los requisitos mínimos de seguridad antes de ser hasheada.
    - El `UserID` no puede ser nulo (`uuid.Nil`).

#### Role
- **Responsabilidad:** Representar un conjunto de permisos que pueden ser asignados a un usuario.
- **Tipos soportados:** `admin`, `commercial`, `designer`, `workshop`.
- **Reglas de Negocio:**
    - El nombre del rol es inmutable y debe ser uno de los definidos.

#### Permission
- **Responsabilidad:** Representar una acción específica que un usuario puede realizar (e.g., 'create:order', 'view:mes').
- **Nota:** Actualmente implementado como una abstracción sobre roles.
