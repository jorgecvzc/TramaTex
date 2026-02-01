# Contratos de API - Módulo IAM

Este documento especifica los contratos de la API para el módulo IAM (Identity and Access Management).

---

## 1. Registrar Usuario

- **Endpoint:** `POST /auth/register`
- **Descripción:** Registra un nuevo usuario en el sistema.
- **Request Body:**
  ```json
  {
    "email": "new.user@example.com",
    "password": "a-strong-password"
  }
  ```
- **Respuesta Exitosa (201 Created):**
  ```json
  {
    "id": "user-uuid-123",
    "email": "new.user@example.com"
  }
  ```
- **Errores:**
  - `400 Bad Request`: Si el email es inválido o la contraseña es débil.
  - `409 Conflict`: Si ya existe un usuario con ese email.

---

## 2. Iniciar Sesión

- **Endpoint:** `POST /auth/login`
- **Descripción:** Autentica un usuario y devuelve tokens JWT de acceso y refresco.
- **Request Body:**
  ```json
  {
    "email": "registered.user@example.com",
    "password": "user-password"
  }
  ```
- **Respuesta Exitosa (200 OK):**
  ```json
  {
    "access_token": "ey...",
    "refresh_token": "ey...",
    "user": {
        "id": "user-uuid-123",
        "email": "registered.user@example.com"
    }
  }
  ```
- **Errores:**
  - `400 Bad Request`: Si el formato de la solicitud es inválido.
  - `401 Unauthorized`: Si las credenciales son incorrectas.

---

## 3. Renovar Token

- **Endpoint:** `POST /auth/refresh`
- **Descripción:** Emite un nuevo access token usando un refresh token válido.
- **Request Body:**
  ```json
  {
    "refresh_token": "a-valid-refresh-token"
  }
  ```
- **Respuesta Exitosa (200 OK):**
  ```json
  {
    "access_token": "ey..."
  }
  ```
- **Errores:**
  - `401 Unauthorized`: Si el refresh token es inválido o está expirado.

---

## 4. Cerrar Sesión

- **Endpoint:** `POST /auth/logout`
- **Descripción:** Invalida la sesión del usuario (por ejemplo, mediante blacklist de tokens).
- **Request Body:** (requiere autenticación)
  ```json
  {}
  ```
- **Respuesta Exitosa (204 No Content):**
- **Errores:**
  - `401 Unauthorized`: Si no se proporciona un access token válido.

---

## 5. Asignar Rol a Usuario

- **Endpoint:** `POST /auth/assign-role`
- **Descripción:** Asigna un rol a un usuario existente (solo admin).
- **Request Body:**
  ```json
  {
    "user_id": "user-uuid-123",
    "role": "designer"
  }
  ```
- **Respuesta Exitosa (200 OK):**
  ```json
  {
    "user_id": "user-uuid-123",
    "role": "designer"
  }
  ```
- **Errores:**
  - `400 Bad Request`: Si el rol es inválido o la solicitud es incorrecta.
  - `401 Unauthorized`: Si no se proporciona un access token válido.
  - `403 Forbidden`: Si el usuario no es admin.
  - `404 Not Found`: Si el usuario no existe.

---

## 6. Verificar Autorización

- **Endpoint:** `POST /auth/authorize`
- **Descripción:** Verifica si un usuario tiene alguno de los roles requeridos.
- **Request Body:**
  ```json
  {
    "user_id": "user-uuid-123",
    "required_roles": ["admin", "commercial"]
  }
  ```
- **Respuesta Exitosa (200 OK):**
  ```json
  {
    "allowed": true,
    "role": "commercial"
  }
  ```
- **Errores:**
  - `400 Bad Request`: Si la solicitud es inválida.
  - `401 Unauthorized`: Si no se proporciona un access token válido.
  - `404 Not Found`: Si el usuario no existe.

---

## 7. Listar Usuarios (Admin)

- **Endpoint:** `GET /auth/users`
- **Descripción:** Devuelve el listado de usuarios registrados (solo admin).
- **Success Response (200 OK):**
  ```json
  {
    "users": [
      {
        "id": "user-uuid-123",
        "email": "user@example.com",
        "role": "commercial"
      }
    ]
  }
  ```
- **Errores:**
  - `401 Unauthorized`: Si no se proporciona un access token válido.
  - `403 Forbidden`: Si el usuario no es admin.
