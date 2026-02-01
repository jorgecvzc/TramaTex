# Casos de Uso - Módulo IAM (Frontend)

## 1. Propósito

Definir los flujos de interfaz necesarios para la gestión de usuarios y asignación de roles en el MVP, alineados con los endpoints IAM.

## 2. Pantallas

### 2.1. Gestión de Usuarios (Admin)

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

### 2.2. Asignación de Rol (Admin)

**Ruta:** `/admin/users/:id/role` (modal o drawer desde listado)

**Objetivo:** Cambiar el rol de un usuario existente.

**Campos:**
- `role` (selector: admin, commercial, designer, workshop)

**Validaciones:**
- Rol requerido.
- Usuario existente.

## 3. Flujos

### 3.1. Asignar Rol
1. Admin entra en `/admin/users`.
2. Selecciona usuario.
3. Abre modal “Asignar rol”.
4. Envía rol seleccionado.
5. UI actualiza listado y muestra confirmación.

## 4. Contratos de API

- **Listar usuarios:** `GET /auth/users`
- **Asignar rol:** `POST /auth/assign-role`
  - Body: `{ "user_id": "...", "role": "designer" }`

> Nota: El endpoint de listado requiere rol `admin`.

## 5. Estados de UI

- **Loading:** spinner en tabla y botón.
- **Error:** banner con mensaje de error y retry.
- **Empty:** estado sin usuarios.

## 6. Accesibilidad

- Labels visibles en campos.
- Navegación por teclado en selector.
- Mensajes de error con aria-live.

---

**Estado:** Documento inicial para MVP.
