# MÓDULO AUTH – ESPECIFICACIÓN TÉCNICA

**Versión:** 1.0  
**Fecha:** 11/01/2026  
**Fase:** 0 (Fundaciones Técnicas)  
**Bounded Context:** Security / Authentication  
**Estado:** 🔵 Especificación Completada - Listo para Desarrollo (Sesión 09)

---

## 📋 DESCRIPCIÓN GENERAL

El módulo **Auth** (Autenticación y Autorización) es el **activo de seguridad crítico** de TramaTex. Responsable de:
- Gestión de identidad de usuarios
- Generación y validación de tokens JWT
- Control de acceso basado en roles (RBAC)
- Auditoría de accesos

**Criticidad:** ⭐⭐⭐⭐⭐ (Máxima)  
**Depende de:** Ningún módulo (es foundacional)  
**Dependen de él:** Todos los módulos (Party, Product, Sales, MES)

---

## 🎯 REQUISITOS FUNCIONALES (RF)

### RF-AUTH-001: Gestión de Usuario

**Descripción:** Sistema debe permitir crear y gestionar usuarios con credenciales seguras.

**Criterios de Aceptación:**
- [ ] Usuario se crea con email único y password hasheado
- [ ] Email se valida antes de crear usuario
- [ ] Password se hashea con bcrypt (mínimo 10 rounds)
- [ ] Password nunca se almacena en plain text
- [ ] Usuario puede cambiar password (post-MVP)

**Flujo:**
```
CreateUser(email, password) 
  → Validar email único
  → Hash password con bcrypt
  → Crear User entity
  → Persistir en BD
  → Retornar user info
```

### RF-AUTH-002: Login con Credenciales

**Descripción:** Usuario puede autenticarse con email + password y recibir tokens JWT.

**Criterios de Aceptación:**
- [ ] Login con email + password retorna access + refresh tokens
- [ ] Password se valida con bcrypt.CompareHashAndPassword
- [ ] Email no encontrado → error 401 Unauthorized
- [ ] Password incorrecto → error 401 Unauthorized
- [ ] Tokens incluyen sub (user ID), email, iat, exp

**Flujo:**
```
Login(email, password)
  → Buscar usuario por email
  → Comparar password con hash
  → Generar access token (15 min)
  → Generar refresh token (7 días)
  → Log: "User X logged in at Y"
  → Retornar tokens
```

### RF-AUTH-003: Validación de Token

**Descripción:** Sistema valida token JWT en cada request autenticado.

**Criterios de Aceptación:**
- [ ] Token válido → claims extraídos correctamente
- [ ] Token expirado → error 401 Unauthorized
- [ ] Token inválido/corrupto → error 401 Unauthorized
- [ ] Token sin firmar → error 401 Unauthorized
- [ ] Middleware valida token en endpoint protegido

**Flujo:**
```
ValidateToken(token)
  → Parsear JWT
  → Verificar firma
  → Validar expiración
  → Retornar claims si OK
  → Retornar error si inválido
```

### RF-AUTH-004: Refresh Token

**Descripción:** Usuario puede obtener nuevo access token sin re-autenticarse.

**Criterios de Aceptación:**
- [ ] POST /api/auth/refresh con refresh token válido → nuevo access token
- [ ] Refresh token expirado → error 401
- [ ] Refresh token no existe en registro → error 401

**Flujo (Post-MVP):**
```
RefreshToken(refreshToken)
  → Validar refresh token
  → Generar nuevo access token
  → Retornar nuevo token
```

### RF-AUTH-005: Logout

**Descripción:** Usuario puede cerrar sesión invalidando tokens.

**Criterios de Aceptación:**
- [ ] POST /api/auth/logout invalida refresh token
- [ ] Access token restante es válido hasta expiración
- [ ] Refresh token no puede usarse nuevamente

**Flujo (Post-MVP):**
```
Logout(refreshToken)
  → Agregar refresh token a blacklist
  → Log: "User X logged out"
  → Retornar 200 OK
```

### RF-AUTH-006: RBAC – Roles y Permisos

**Descripción:** Sistema controla acceso basado en roles del usuario.

**Roles definidos (MVP):**
- **Admin:** Acceso total a sistema
- **Comercial:** Gestión de clientes, pedidos
- **Diseño:** Creación y aprobación de diseños
- **Taller:** Ejecución de trabajos en MES

**Criterios de Aceptación:**
- [ ] Usuario asignado a rol en creación
- [ ] Endpoint protegido valida rol antes de ejecutar
- [ ] Rol insuficiente → error 403 Forbidden
- [ ] Token incluye rol en claims

**Flujo:**
```
Middleware: RequireRole("Comercial")
  → Extraer rol de token
  → Comparar con rol requerido
  → Permitir (200) o rechazar (403)
```

---

## 🚫 REQUISITOS NO FUNCIONALES (RNF)

| ID | Requisito | Métrica |
|----|-----------|---------|
| RNF-AUTH-001 | Seguridad Password | bcrypt ≥10 rounds |
| RNF-AUTH-002 | Seguridad Token | JWT firmado con HS256 o RS256 |
| RNF-AUTH-003 | Expiración Access Token | 15 minutos |
| RNF-AUTH-004 | Expiración Refresh Token | 7 días |
| RNF-AUTH-005 | Rate Limiting Login | Máx 5 intentos / 15 min |
| RNF-AUTH-006 | Validación Email | RFC 5322 básico |
| RNF-AUTH-007 | Logging | Todos los eventos de auth |
| RNF-AUTH-008 | Testing | ≥90% cobertura dominio |
| RNF-AUTH-009 | Performance | Login < 100ms |
| RNF-AUTH-010 | HTTPS | Tokens solo en HTTPS (post-MVP) |

---

## 🎬 CASOS DE USO

### CU-AUTH-001: Login Exitoso

**Actor:** Usuario sin autenticar  
**Precondición:** Usuario existe en sistema  
**Flujo:**
1. Usuario abre login en frontend
2. Ingresa email y password
3. Frontend hace POST /api/auth/login
4. Backend valida credenciales
5. Backend genera tokens
6. Frontend recibe tokens, almacena en localStorage
7. Usuario autenticado accede a sistema

**Postcondición:** Tokens válidos, usuario en sesión

### CU-AUTH-002: Login Fallido – Usuario No Existe

**Actor:** Usuario sin autenticar  
**Precondición:** Email no registrado  
**Flujo:**
1. Usuario intenta login con email no existente
2. Backend busca usuario → no encuentra
3. Backend retorna 401 "Invalid credentials"
4. Frontend muestra error genérico (no revelar que email no existe)

**Postcondición:** Login rechazado, sin sesión

### CU-AUTH-003: Login Fallido – Password Incorrecto

**Actor:** Usuario sin autenticar  
**Precondición:** Usuario existe, password es incorrecto  
**Flujo:**
1. Usuario intenta login
2. Backend encuentra usuario
3. Backend compara password → no coincide
4. Backend retorna 401 "Invalid credentials"

**Postcondición:** Login rechazado, contador de intentos +1

### CU-AUTH-004: Acceso a Endpoint Protegido

**Actor:** Usuario autenticado  
**Precondición:** Tiene token válido  
**Flujo:**
1. Usuario hace request a endpoint protegido (ej: GET /api/parties)
2. Middleware valida token en header `Authorization: Bearer <token>`
3. Token válido y no expirado
4. Middleware extrae user ID del token
5. Request procede, user ID disponible en contexto

**Postcondición:** Endpoint se ejecuta con user context

### CU-AUTH-005: Acceso Denegado – Token Expirado

**Actor:** Usuario con token expirado  
**Precondición:** Access token expiró hace 5 minutos  
**Flujo:**
1. Usuario hace request con token expirado
2. Middleware valida token → expiration check falla
3. Middleware retorna 401 "Token expired"
4. Frontend detecta 401, redirige a login (o intenta refresh)

**Postcondición:** Request rechazado, usuario debe re-autenticarse

### CU-AUTH-006: Acceso Denegado – Rol Insuficiente

**Actor:** Usuario con rol "Taller"  
**Precondición:** Intenta acceder a endpoint que requiere "Admin"  
**Flujo:**
1. Usuario hace request a admin panel
2. Middleware valida token → OK
3. Handler verifica rol requerido ("Admin") vs rol de usuario ("Taller")
4. Rol no coincide
5. Handler retorna 403 "Insufficient permissions"

**Postcondición:** Request rechazado por permisos

---

## 👤 HISTORIAS DE USUARIO

### HU-001: Login como Comercial

**Como:** Usuario comercial  
**Quiero:** Autenticarme con email y password  
**Para:** Acceder al sistema de gestión de pedidos  

**Criterios de aceptación:**
- [ ] Puedo ver formulario login
- [ ] Puedo ingresar email y password
- [ ] Al hacer click "Entrar", valida credenciales
- [ ] Si son correctas, me lleva al dashboard
- [ ] Si son incorrectas, veo error amigable
- [ ] Mi sesión persiste al recargar página

### HU-002: Logout

**Como:** Usuario autenticado  
**Quiero:** Cerrar sesión  
**Para:** Que otra persona no use mi cuenta  

**Criterios de aceptación:**
- [ ] Veo botón "Logout" en navbar
- [ ] Al clickear, sesión se cierra
- [ ] Tokens se limpian del localStorage
- [ ] Me redirige a login

### HU-003: Permiso Denegado

**Como:** Usuario con rol "Taller"  
**Quiero:** Ver qué pasó si intento acceder admin panel  
**Para:** Entender mis permisos  

**Criterios de aceptación:**
- [ ] Si intento acceder URL /admin, veo "No tienes permiso"
- [ ] No me muestra contenido sensible
- [ ] Me da opción de volver a mi sección

---

## 🏗️ ENTIDADES DE DOMINIO

### Entity: User

```go
type User struct {
    ID        string        // UUID
    Email     Email         // Value Object
    Password  Password      // Value Object (hasheado)
    Role      Role          // "Admin" | "Comercial" | "Diseño" | "Taller"
    CreatedAt time.Time
    UpdatedAt time.Time
    Active    bool
}
```

**Invariantes:**
- ID no vacío
- Email válido y único
- Password nunca en plain text
- Role en lista autorizada

### Value Object: Email

```go
type Email struct {
    value string // "user@example.com"
}
```

**Validación:**
- Formato válido (RFC 5322 básico)
- No vacío
- Máximo 255 caracteres

### Value Object: Password

```go
type Password struct {
    hash string // bcrypt hash ($2a$10$...)
}
```

**Validación en creación:**
- Mínimo 8 caracteres
- Máximo 72 caracteres (bcrypt limit)
- Hash con bcrypt cost=10

**Método:**
```go
func (p *Password) Matches(plain string) bool {
    return bcrypt.CompareHashAndPassword(p.hash, plain) == nil
}
```

### Value Object: TokenClaims

```go
type TokenClaims struct {
    Subject   string    // User ID
    Email     string
    Role      string
    IssuedAt  time.Time
    ExpiresAt time.Time
}
```

### Enum: Role

```go
type Role string

const (
    RoleAdmin      Role = "Admin"
    RoleComercial  Role = "Comercial"
    RoleDiseño     Role = "Diseño"
    RoleTaller     Role = "Taller"
)
```

---

## 🔌 INTERFACES DE DOMINIO

### Repository Interface

```go
type UserRepository interface {
    // Buscar usuario por ID
    ByID(ctx context.Context, id string) (*User, error)
    
    // Buscar usuario por email
    ByEmail(ctx context.Context, email string) (*User, error)
    
    // Guardar usuario (crear o actualizar)
    Save(ctx context.Context, user *User) error
    
    // Eliminar usuario
    Delete(ctx context.Context, id string) error
}
```

**Implementación:** Infrastructure layer (PostgreSQL con GORM)

### JWT Service Interface

```go
type JWTService interface {
    // Generar access token (15 min)
    GenerateAccessToken(user *User) (string, error)
    
    // Generar refresh token (7 días)
    GenerateRefreshToken(user *User) (string, error)
    
    // Validar y extraer claims
    ValidateToken(token string) (*TokenClaims, error)
}
```

**Implementación:** Infrastructure layer (con jwt-go o similar)

---

## 📊 SCHEMA DE BASE DE DATOS

```sql
-- Tabla de usuarios
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(72) NOT NULL,  -- bcrypt hash
    role VARCHAR(50) NOT NULL,           -- Admin, Comercial, Diseño, Taller
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de refresh tokens (para blacklist/revocation)
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de auditoría de login
CREATE TABLE auth_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    event VARCHAR(50) NOT NULL,          -- login, logout, failed_login
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_auth_audit_logs_user_id ON auth_audit_logs(user_id);
CREATE INDEX idx_auth_audit_logs_created_at ON auth_audit_logs(created_at);
```

---

## 🔄 FLUJOS PRINCIPALES

### Flujo: Login

```
┌─────────────┐
│   Frontend  │
│  Form Login │
└──────┬──────┘
       │ POST /api/auth/login
       │ {email, password}
       │
       ▼
┌──────────────────────────────────────┐
│     Interfaces (HTTP Handler)        │
│  HandleLogin(email, password)        │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│      Application (Use Case)          │
│  LoginUseCase.Execute()              │
│  1. Validar input                    │
│  2. Buscar usuario por email         │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│       Domain (Business Logic)        │
│  1. User.Password.Matches(plain)     │
│  2. JWTService.GenerateTokens()      │
│  3. Return tokens                    │
└──────┬───────────────────────────────┘
       │
       ▼ {accessToken, refreshToken}
┌─────────────┐
│   Frontend  │
│   Almacena  │
│   localStorage
└─────────────┘
```

### Flujo: Validación de Token (Middleware)

```
┌──────────────────────────┐
│  Request con Token       │
│  Authorization: Bearer.. │
└──────┬───────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│     Middleware (Interfaces)          │
│  ValidateToken(Authorization header)│
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│       Domain (JWTService)            │
│  JWTService.ValidateToken(token)     │
│  - Parse JWT                         │
│  - Verify signature                  │
│  - Check expiration                  │
│  - Return TokenClaims               │
└──────┬───────────────────────────────┘
       │
    ┌──┴──┐
    │     │
   ✓      ✗
  OK    401
    │      │
    ▼      ▼
 Continua Rechaza
```

---

## 📍 API REST ENDPOINTS

### POST /api/auth/login

**Autenticación:** No requiere  
**Request:**
```json
{
  "email": "comercial@company.com",
  "password": "SecurePass123!"
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": "uuid...",
    "email": "comercial@company.com",
    "role": "Comercial"
  }
}
```

**Response 401:**
```json
{
  "error": "Invalid credentials"
}
```

### POST /api/auth/refresh

**Autenticación:** No requiere  
**Request:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGc..."
}
```

### POST /api/auth/logout

**Autenticación:** Requiere access_token  
**Request:** (Bearer token en header)  
**Response 200:** OK

### GET /api/auth/me (Ejemplo endpoint protegido)

**Autenticación:** Requiere access_token  
**Response 200:**
```json
{
  "id": "uuid...",
  "email": "user@company.com",
  "role": "Comercial"
}
```

**Response 401:** Token inválido/expirado

---

## ✅ CRITERIOS DE ÉXITO (Sesión 09)

**Backend:**
- ✅ User, Email, Password, TokenClaims implementados sin deps externas
- ✅ UserRepository interface definida
- ✅ LoginUseCase orquestan correctamente
- ✅ JWTService interface definida (implementación en infra)
- ✅ HTTP handler para POST /api/auth/login funcional
- ✅ Tests unitarios: ≥90% cobertura dominio
- ✅ Tests integración: LoginUseCase con mocks
- ✅ Docker-compose con JWT_SECRET
- ✅ go test ./... pasa 100%
- ✅ golangci-lint ./... sin warnings

**Documentación:**
- ✅ Sesión 09 registrada completa
- ✅ Commit messages descriptivos
- ✅ PROJECT_STATUS.md actualizado

---

## 🚀 PRÓXIMAS FASES (Post-MVP)

- [ ] **Sesión 10:** Frontend Login component + integración
- [ ] **Post-MVP:** Refresh token logic, logout, token revocation
- [ ] **Post-MVP:** Rate limiting, CAPTCHA
- [ ] **Post-MVP:** 2FA (Two-Factor Authentication)
- [ ] **Post-MVP:** OAuth2 / OIDC integration
- [ ] **Post-MVP:** Auditoría avanzada de accesos

---

**Documento de especificación completado.**  
**Listo para Sesión 09 (Desarrollo).**

