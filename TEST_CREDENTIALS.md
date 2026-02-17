# 🔐 Credenciales de Prueba - TramaTex ERP

**Última actualización:** 2026-02-05

---

## 👤 Usuario Administrador

### Credenciales de Login

```
Email:    admin@tramatex.local
Password: admin123
```

**Roles:** `admin`  
**Estado:** Activo  
**UUID:** `f47ac10b-58cc-4372-a567-0e02b2c3d479`

---

## 🔑 Hash de Contraseña (bcrypt)

```
$2a$10$Gd8.JP/L3j.vNvap81EpjuX7G4u5KKLmf10TSxmu779Mq/HdC/B9e
```

**Algoritmo:** bcrypt  
**Cost Factor:** 10  
**Método de generación:** `golang.org/x/crypto/bcrypt`

---

## 🛠️ Regenerar Hash de Contraseña

Si necesitas cambiar la contraseña o generar un hash para otro usuario:

### Opción 1: Usar el generador incluido

```bash
cd apps/tramatex-api/cmd/hashgen
go run main.go <tu-contraseña>
```

**Ejemplo:**
```bash
go run main.go nuevaPassword123
```

### Opción 2: Usar Go directamente

```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "tu-contraseña"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
    fmt.Println(string(hash))
}
```

### Opción 3: Comando online (solo para desarrollo)

Puedes usar herramientas online como:
- https://bcrypt-generator.com/ (asegúrate de usar cost factor 10)

---

## 🗄️ Actualizar en Base de Datos

Si la base de datos ya está creada y necesitas actualizar la contraseña:

```sql
UPDATE users 
SET password_hash = '$2a$10$Gd8.JP/L3j.vNvap81EpjuX7G4u5KKLmf10TSxmu779Mq/HdC/B9e' 
WHERE email = 'admin@tramatex.local';
```

O ejecuta las migraciones de nuevo:

```bash
cd apps/tramatex-api
go run cmd/migrate/main.go
```

---

## 📝 Archivos Actualizados

Los siguientes archivos contienen el hash de la contraseña del usuario admin:

1. **Migración SQL:**
   - `apps/tramatex-api/migrations/006_seed_admin_user.sql`

2. **Seed de Go:**
   - `apps/tramatex-api/internal/iam/infrastructure/persistence/models.go`
   - Función: `SeedAdminUser()`

---

## 🧪 Probar Login

### Con curl:

```bash
curl -X POST http://localhost:8080/api/v1/iam/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@tramatex.local",
    "password": "admin123"
  }'
```

### Respuesta esperada:

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "email": "admin@tramatex.local",
    "role": "admin",
    "is_active": true
  }
}
```

### Con Frontend (UI):

1. Abre: http://localhost:5173/login
2. Ingresa:
   - **Email:** `admin@tramatex.local`
   - **Password:** `admin123`
3. Click en "Iniciar sesión"

---

## 🔒 Seguridad

### ⚠️ IMPORTANTE - Solo para Desarrollo

- **NUNCA uses estas credenciales en producción**
- **NUNCA commitees contraseñas en texto plano**
- **SIEMPRE cambia las credenciales por defecto en producción**

### Mejores Prácticas para Producción

1. **Contraseñas fuertes:**
   - Mínimo 12 caracteres
   - Letras mayúsculas, minúsculas, números y símbolos
   - No usar palabras del diccionario

2. **Rotación de credenciales:**
   - Cambiar contraseñas cada 90 días
   - No reutilizar contraseñas anteriores

3. **2FA (Futuro):**
   - Implementar autenticación de dos factores
   - Ver roadmap en `docs/architecture/security-strategy.md`

---

## 📚 Referencias

- **Módulo IAM:** `docs/modules/iam/`
- **Password Domain Model:** `apps/tramatex-api/internal/iam/domain/model/password.go`
- **Security Strategy:** `docs/architecture/security-strategy.md`
- **OWASP Audit:** `docs/architecture/security/owasp-audit-report.md`

---

## 🆘 Troubleshooting

### "Invalid credentials" al hacer login

1. Verifica que la base de datos tiene el usuario seeded:
   ```sql
   SELECT id, email, role, is_active FROM users WHERE email = 'admin@tramatex.local';
   ```

2. Verifica que el hash de contraseña está actualizado:
   ```sql
   SELECT password_hash FROM users WHERE email = 'admin@tramatex.local';
   ```
   Debe ser: `$2a$10$Gd8.JP/L3j.vNvap81EpjuX7G4u5KKLmf10TSxmu779Mq/HdC/B9e`

3. Regenera las migraciones:
   ```bash
   # Elimina y recrea la base de datos
   docker-compose down -v
   docker-compose up -d postgres
   cd apps/tramatex-api
   go run cmd/migrate/main.go
   ```

### "User not found"

El usuario admin no existe en la base de datos. Ejecuta las migraciones:

```bash
cd apps/tramatex-api
go run cmd/migrate/main.go
```

---

**Última generación del hash:** 2026-02-05  
**Generador usado:** `apps/tramatex-api/cmd/hashgen/main.go`
