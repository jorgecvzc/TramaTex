# Casos de Uso - Módulo IAM (Backend + Frontend)

## 1. Propósito

Documentar los casos de uso del módulo IAM en backend (API) y los flujos de UI mínimos del frontend.

## 2. Casos de Uso Backend (API)

### Autenticación y Sesión
- **Registrar Usuario:** Crear un usuario con email y contraseña.
- **Login de Usuario:** Validar credenciales y emitir `access_token` + `refresh_token`.
- **Refresh Token:** Emitir un nuevo `access_token` con un `refresh_token` válido.
- **Logout:** Invalidar sesión/tokens.

### Gestión de Usuarios (Admin)
- **Listar Usuarios:** Obtener lista de usuarios registrados.
- **Crear Usuario (Admin):** Crear usuario con rol explícito.
- **Asignar Rol:** Cambiar el rol de un usuario existente.
- **Eliminar Usuario:** Eliminar/desactivar usuario.

### Autorización
- **Check Authorization:** Verificar si un usuario tiene roles permitidos.

## 3. Casos de Uso Frontend (UI)

### 3.1. Gestión de Usuarios (Admin)

**Ruta:** `/admin/users`

**Objetivo:** Listar usuarios, consultar detalles y asignar roles.

**Componentes principales:**
- Tabla de usuarios con búsqueda (email, rol, estado).
- Detalle lateral (drawer) o modal para edición de rol.
- Feedback de éxito/error.

**Acciones:**
- Ver listado de usuarios.
- Seleccionar usuario.
- Asignar rol.

**Permisos:** Solo `admin`.

### 3.2. Asignación de Rol (Admin)

**Ruta:** `/admin/users/:id/role` (modal o drawer desde listado)

**Objetivo:** Cambiar el rol de un usuario existente.

**Campos:**
- `role` (selector: admin, commercial, designer, workshop)

**Validaciones:**
- Rol requerido.
- Usuario existente.

## 4. Flujos UI

### 4.1. Asignar Rol
1. Admin entra en `/admin/users`.
2. Selecciona usuario.
3. Abre modal “Asignar rol”.
4. Envía rol seleccionado.
5. UI actualiza listado y muestra confirmación.

## 5. Contratos de API

- **Listar usuarios:** `GET /auth/users`
- **Crear usuario (admin):** `POST /auth/users`
- **Asignar rol:** `POST /auth/assign-role`
  - Body: `{ "user_id": "...", "role": "designer" }`

> Nota: Endpoints de gestión requieren rol `admin`.

## 6. Estados de UI

- **Loading:** spinner en tabla y botón.
- **Error:** banner con mensaje de error y retry.
- **Empty:** estado sin usuarios.

## 7. Accesibilidad

- Labels visibles en campos.
- Navegación por teclado en selector.
- Mensajes de error con aria-live.

---

**Estado:** Documento actualizado para MVP.
