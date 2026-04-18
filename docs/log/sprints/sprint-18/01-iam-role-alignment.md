# Tarea 18-01: IAM Role Alignment — Fix DB↔Domain↔Frontend Role Mismatch

- **Sprint:** 18
- **Estado:** ⏳ Planificado
- **Fecha de inicio:** —
- **Rama:** `fix/iam-role-alignment`
- **Facilitador:** GitHub Copilot / Claude Opus 4.6

---

## Contexto

Al intentar crear un usuario desde la UI de administración (`/admin/users`), se obtiene un error **HTTP 400**:

```
Request failed with status code 400
```

### Causa raíz

**Desalineación de roles entre la base de datos y el código:**

| Capa | Roles definidos |
|------|-----------------|
| **SQL** (`001_init_iam.sql`, constraint `chk_role`) | `admin`, `manager`, `operator`, `cashier` |
| **Go Domain** (`role.go`) | `admin`, `commercial`, `designer`, `workshop` |
| **Frontend types** (`auth.ts`) | `admin`, `commercial`, `designer`, `workshop` |
| **Frontend UI** (`UsersManagement.vue`) | `admin`, `commercial`, `designer`, `workshop` |
| **ADR-014** | `admin`, `commercial`, `designer`, `workshop` |

El domain model Go acepta `commercial` (pasa validación), pero PostgreSQL rechaza el INSERT porque el constraint `chk_role` solo permite `admin, manager, operator, cashier`.

### Reproducción (curl)

```bash
# Login como admin
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@tramatex.local","password":"admin123"}'

# Crear usuario → ERROR 400
curl -X POST http://localhost:8080/auth/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"email":"test@tramatex.local","password":"password123","role":"commercial"}'

# Respuesta:
# {"error":"failed to save user: failed to save user: ERROR: new row for relation \"users\" violates check constraint \"chk_role\" (SQLSTATE 23514)"}
```

---

## Decisión pendiente

**¿Qué set de roles es el correcto?**

**Propuesta:** Mantener los roles del ADR-014 y el domain model (`admin`, `commercial`, `designer`, `workshop`), ya que reflejan la intención de negocio de TramaTex (empresa textil con comerciales, diseñadores y taller).

Los roles de la migración SQL (`manager`, `operator`, `cashier`) parecen ser genéricos/placeholder que nunca se actualizaron al definir el ADR-014.

---

## Plan de implementación

### Paso 1: Crear nueva migración SQL
- Fichero: `apps/tramatex-api/migrations/NNN_align_iam_roles.sql`
- Actualizar constraint `chk_role` a `CHECK (role IN ('admin', 'commercial', 'designer', 'workshop'))`
- Actualizar el default role de `'operator'` a un valor válido (ej. `'commercial'`)
- UPDATE de registros existentes si hay datos con roles antiguos

### Paso 2: Verificar coherencia del domain model
- `apps/tramatex-api/internal/iam/domain/model/role.go` — ya correcto ✅
- `apps/frontend/src/types/auth.ts` — ya correcto ✅
- `apps/frontend/src/pages/admin/UsersManagement.vue` — ya correcto ✅

### Paso 3: Actualizar documentación
- `docs/modules/iam/module-spec.md` — verificar roles documentados
- `docs/architecture/adrs/adr-014-iam-module-architecture.md` — ya correcto ✅

### Paso 4: Fix bug adicional — `updateUser()` no definido
- `UsersManagement.vue` llama a `iamService.updateUser()` pero ese método no existe en `iam.ts`
- Implementar `updateUser()` en el servicio o eliminar la llamada si no hay endpoint backend

### Paso 5: Tests
- Test de creación de usuario con cada rol válido (admin, commercial, designer, workshop)
- Test que confirme rechazo de roles inválidos
- Verificar que `assignRole` también funciona con los roles correctos

---

## Archivos afectados

| Archivo | Acción |
|---------|--------|
| `apps/tramatex-api/migrations/001_init_iam.sql` | Referencia (no modificar, crear nueva migración) |
| `apps/tramatex-api/migrations/NNN_align_iam_roles.sql` | **CREAR** — Nueva migración |
| `apps/tramatex-api/internal/iam/domain/model/role.go` | Verificar (ya correcto) |
| `apps/frontend/src/types/auth.ts` | Verificar (ya correcto) |
| `apps/frontend/src/pages/admin/UsersManagement.vue` | Fix `updateUser()` call |
| `apps/frontend/src/services/iam.ts` | Añadir `updateUser()` o corregir referencia |
| `docs/modules/iam/module-spec.md` | Verificar coherencia roles |

---

## Resultado esperado

- Crear usuarios desde la UI sin error 400
- Roles consistentes en todas las capas: BD, Go, Frontend, Docs
- Tests validando creación con cada rol
