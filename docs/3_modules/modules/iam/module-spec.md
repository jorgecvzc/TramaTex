# Módulo de IAM (Identity and Access Management)

## 1. Propósito

*   **Visión del Módulo:** Gestionar la identidad del usuario, autenticación, autorización y control de acceso para toda la plataforma.
*   **Objetivos Clave:**
    *   Proporcionar un sistema centralizado para gestionar todos los aspectos de la identidad del usuario y el control de acceso.
    *   Registrar usuarios, verificar credenciales, gestionar sesiones, y controlar el acceso a diferentes partes de la aplicación basado en roles y permisos.

## 2. Requisitos

### 2.1. Requisitos Funcionales

*   **RF-001:** Registro de nuevos usuarios.
*   **RF-002:** Autenticación de usuarios mediante correo electrónico y contraseña.
*   **RF-003:** Emisión y renovación de tokens de acceso (JWT).
*   **RF-004:** Gestión de la sesión del usuario.
*   **RF-005:** Modelo de control de acceso basado en roles (RBAC).
*   **RF-006:** Definición de Usuarios, Roles y Permisos.
*   **RF-007:** Protección de endpoints de la API.

### 2.2. Requisitos No Funcionales

*   **RNF-001 (Arquitectura):** La implementación debe seguir los principios de Clean Architecture como se define en `ADR-002`.
*   **RNF-002 (Arquitectura):** La estructura del módulo debe adherirse a `ADR-009`.
*   **RNF-003 (Seguridad):** Debe utilizar el servicio de seguridad de `infrastructure/security` para la gestión de JWT y hashing de contraseñas.

**Fuera del Alcance (para la primera versión):**
- Autenticación de dos factores (2FA).
- Integración con proveedores de identidad de terceros (OAuth, SAML).
- Gestión avanzada de políticas de contraseñas.

## 3. Casos de Uso

### 3.1. Actores
*   **Usuario Anónimo:** Un usuario que aún no ha iniciado sesión.
*   **Usuario Autenticado:** Un usuario que ha iniciado sesión.
*   **Administrador:** Un usuario con permisos para gestionar otros usuarios y roles.

### 3.2. Casos de Uso Principales

*   **CU-001: RegisterUser**
    *   **Actor:** Usuario Anónimo
    *   **Descripción:** Registrar un nuevo usuario en el sistema.
*   **CU-002: LoginUser**
    *   **Actor:** Usuario Anónimo
    *   **Descripción:** Autenticar a un usuario y devolver tokens de acceso.
*   **CU-003: LogoutUser**
    *   **Actor:** Usuario Autenticado
    *   **Descripción:** Invalidar la sesión de un usuario.
*   **CU-004: RefreshToken**
    *   **Actor:** Usuario Autenticado
    *   **Descripción:** Renovar un token de acceso expirado usando un token de refresco.
*   **CU-005: AssignRoleToUser**
    *   **Actor:** Administrador
    *   **Descripción:** Asignar un rol a un usuario.
*   **CU-006: CheckAuthorization**
    *   **Actor:** Sistema
    *   **Descripción:** Verificar si un usuario tiene permiso para realizar una acción.

## 4. Historias de Usuario

*   **HU-001:** Como usuario nuevo, quiero registrarme con mi correo y contraseña para poder acceder a la plataforma.
*   **HU-002:** Como usuario registrado, quiero iniciar sesión para poder utilizar las funcionalidades del sistema.
*   **HU-003:** Como administrador, quiero poder asignar roles a los usuarios para controlar sus permisos.

## 5. Criterios de Aceptación

*   **Para HU-001:**
    *   **Criterio 1:** Dado que ingreso un correo no existente y una contraseña válida, cuando envío el formulario de registro, entonces se crea un nuevo usuario en la base de datos.
*   **Para HU-002:**
    *   **Criterio 1:** Dado que soy un usuario registrado, cuando ingreso mis credenciales correctas, entonces recibo un token de acceso y soy redirigido al dashboard.

## 6. Modelo de Dominio

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
    - El nombre del rol debe ser único (e.g., 'admin', 'operador').

#### Permission
- **Responsabilidad:** Representar una acción específica que un usuario puede realizar.
- **Value Objects:** `PermissionID`, `PermissionName`.

## 7. Decisiones de Diseño

- La estructura de directorios sigue la organización de monolito modular definida en **ADR-009**.
- La implementación de las capas de dominio, aplicación e infraestructura se detalla en el `domain-model.md` de este módulo.
- Los contratos de la API REST se especifican en `api-contracts.md`.

---

## 8. Anexo: Clarificación de Terminología (IAM vs. Auth)

### Diferencia Conceptual

| Aspecto | Auth (Autenticación) | IAM (Gestión de Identidad y Acceso) |
|---|---|---|
| **Alcance** | Login, Logout, Tokens | Login + Autorización + Control de Acceso |
| **Responsabilidades** | Verificar credenciales | Identidad + Permisos + Roles |
| **Evolución** | Solo autenticación | Plataforma completa de acceso |

TramaTex utiliza **IAM (Identity and Access Management)** como el término oficial para este bounded context porque abarca no solo la autenticación (verificar quién eres) sino también la autorización (qué puedes hacer) a través de roles y permisos.

### FAQ

**P: ¿Por qué cambiar de "auth" a "iam"?**
R: Claridad conceptual. "IAM" es más descriptivo de lo que realmente hace el módulo: gestionar la identidad, la autenticación y la autorización.

**P: ¿Qué nombres deben usarse en el proyecto?**
R: Se debe usar consistentemente "Módulo IAM" o "Contexto Delimitado IAM", y las rutas de código correspondientes como `internal/iam/`. Se debe evitar el término "Módulo Auth".