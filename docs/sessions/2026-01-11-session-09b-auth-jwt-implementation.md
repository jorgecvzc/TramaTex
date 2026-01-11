# SESSION 09B – AUTENTICACIÓN JWT: IMPLEMENTACIÓN MÓDULO AUTH

**Sesión:** Session-09b (Auth JWT Implementation)  
**Facilitador/LLM:** GitHub Copilot (Claude Haiku 4.5)  
**Fecha Inicio:** 2026-01-11 (16:00 UTC aprox.)  
**Duración Estimada:** 6 horas (Fases 1-5)  
**Participantes:** Jorge Cortés Villalba, GitHub Copilot  
**Status:** 🟢 INICIADA  

---

## 📋 ESTRUCTURA DE FASES COMPLETA

### **Fase 1 – Análisis Arquitectónico** ✅ COMPLETADA

**Duración:** ~20 minutos  
**Propósito:** Defines bounded context, dependencies, interfaces  

**Bounded Context:** Auth/Security (Foundational)  
**Criticality:** ⭐⭐⭐⭐⭐ (Maximum - all modules depend on this)  
**Dependencies:** ZERO external dependencies (all modules depend on Auth)

**Key Domain Interfaces:**
```go
type Email struct { value string }              // RFC 5322 validation
type Password struct { hash string }            // Bcrypt (cost ≥10)
type User struct {                              // Root Aggregate
    ID, Email, Password, Role, Active, Timestamps
}
type TokenClaims struct {                       // JWT Standard Claims
    Subject, Email, Role, IssuedAt, ExpiresAt
}
interface UserRepository {                      // Domain contract
    ByID(), ByEmail(), Save(), Delete()
}
interface JWTService {                          // Domain contract
    GenerateAccessToken(), GenerateRefreshToken(), ValidateToken()
}
```

---

### **Fase 2 – Tests TDD-First** ✅ COMPLETADA

**Duración:** ~30 minutos  
**Propósito:** Complete test specifications BEFORE implementation (34 test cases)

**Test Distribution Rule:** 70% Unit / 25% Integration / 5% E2E (critical only)

**Planned Test Cases (34 total):**

#### Email VO Tests (7 tests, 70% of Email module)
```go
1. TestEmailNewWithValidFormat                 // Valid: user@domain.com
2. TestEmailNewWithInvalidFormat               // Invalid: no @
3. TestEmailNewWithEmptyString                 // Empty: ""
4. TestEmailNewWithTooLongAddress              // >254 chars
5. TestEmailEquals                             // Equivalence testing
6. TestEmailNotEquals                          // Inequality
7. TestEmailImmutable                          // Verify immutability
```

#### Password VO Tests (8 tests, 70% of Password module)
```go
1. TestPasswordNewWithValidLength              // 8-72 chars
2. TestPasswordNewWithTooShort                 // <8 chars
3. TestPasswordNewWithTooLong                  // >72 chars
4. TestPasswordNewWithEmptyString              // Empty
5. TestPasswordMatchesWithCorrectPassword      // Correct match
6. TestPasswordMatchesWithWrongPassword        // Wrong match
7. TestPasswordNeverStoredPlaintext            // Hash verification
8. TestPasswordBcryptCostAtLeast10             // Security requirement
```

#### User Entity Tests (9 tests, 70% of User module)
```go
1. TestUserNewWithValidData                    // Valid creation
2. TestUserNewWithMissingEmail                 // Email required
3. TestUserNewWithMissingPassword              // Password required
4. TestUserNewWithEmptyID                      // ID required
5. TestUserNewWithInvalidRole                  // Role validation
6. TestUserNewWithValidRoles                   // Roles: admin, manager, operator
7. TestUserImmutableAfterCreation              // Aggregation guarantee
8. TestUserTimestampsAutomatic                 // CreatedAt, UpdatedAt
9. TestUserActiveFlag                          // Default: true
```

#### TokenClaims VO Tests (4 tests, 70% of TokenClaims module)
```go
1. TestTokenClaimsCreation                     // Valid token claims
2. TestTokenClaimsIsExpired                    // Expired check
3. TestTokenClaimsNotExpired                   // Valid expiry
4. TestTokenClaimsAllFields                    // Subject, Email, Role, Iat, Exp
```

#### LoginUseCase Integration Tests (6 tests, 25% of Auth layer)
```go
1. TestLoginWithValidCredentials               // Success path (integration)
2. TestLoginWithUserNotFound                   // 404 scenario
3. TestLoginWithInvalidPassword                // Wrong password
4. TestLoginWithEmptyEmail                     // Input validation
5. TestLoginWithEmptyPassword                  // Input validation
6. TestLoginOutputDTOMapping                   // DTO correctness
```

**Mock Patterns Established:**
```go
type MockUserRepository struct {
    GetByEmailFunc func(ctx context.Context, email string) (*User, error)
    SaveFunc       func(ctx context.Context, user *User) error
}

type MockJWTService struct {
    GenerateAccessTokenFunc func(claims *TokenClaims) (string, error)
    ValidateTokenFunc       func(token string) (*TokenClaims, error)
}
```

---

### **Fase 3 – Implementación por Capas** ⏳ EN PROGRESO

**Duración:** ~4 horas  
**Propósito:** Write actual Go code layer by layer (TDD-driven)

**Implementation Sequence:**
1. Domain Layer - Value Objects & Entities (Email, Password, User, TokenClaims)
2. Domain Layer - Interfaces (UserRepository, JWTService)
3. Application Layer - Use Cases (LoginUseCase)
4. Application Layer - DTOs (LoginInput, LoginOutput, UserDTO)
5. Interfaces Layer - HTTP Handler (POST /api/auth/login)
6. Infrastructure Layer (Post-session) - Database & JWT implementation

**Phase 3 Milestones:**
- [ ] Email VO (validation, immutability)
- [ ] Password VO (bcrypt, cost ≥10)
- [ ] User Entity (invariants, timestamps)
- [ ] TokenClaims VO (standard JWT claims)
- [ ] Domain Interfaces (UserRepository, JWTService)
- [ ] LoginUseCase (orchestration)
- [ ] DTOs (mapping)
- [ ] HTTP Handler (Gin binding)

---

### **Fase 4 – Validación** ⏳ PENDIENTE

**Duración:** ~1 hora  
**Requerimientos:**
- `go test ./... -v` → 100% pasan
- `golangci-lint run ./...` → 0 warnings
- `go fmt` → sin cambios
- `docker-compose up` → sin errores
- Coverage: domain/user ≥90%, domain/security ≥80%, application/auth ≥80%
- HTTP endpoint: `curl -X POST http://localhost:8080/api/auth/login`

---

### **Fase 5 – Documentación & Commits** ⏳ PENDIENTE

**Duración:** ~1 hora  
**Entregables:**
- [ ] Update auth.md with implementation details
- [ ] Final git commit with all code (Phase 3)
- [ ] Update PROJECT_STATUS.md
- [ ] Close Session 09b

---

## 🎯 OBJETIVOS PRINCIPALES (5 objetivos)

### 1. [ ] Implementar User y Password Value Objects en Dominio

**Descripción:** Crear entidades de dominio para representar usuario y su contraseña hasheada.

**Subtareas:**
- [ ] Value Object `Email` con validación (RFC 5322 básico)
- [ ] Value Object `Password` con hashing bcrypt (cost ≥10)
- [ ] Entity `User` con invariantes
- [ ] Tests unitarios: ≥100% cobertura (TDD-first)
- [ ] Error types: `domain.ValidationError`, `domain.AuthError`

**Salida esperada:**
```
backend/internal/domain/user/
├── email.go          (≈50 líneas)
├── email_test.go     (≈80 líneas)
├── password.go       (≈60 líneas)
├── password_test.go  (≈100 líneas)
├── user.go           (≈80 líneas)
├── user_test.go      (≈140 líneas)
└── repository.go     (≈20 líneas)
```

**Validación:** `go test ./internal/domain/user/... -v` → 100% pasan, cobertura ≥90%

---

### 2. [ ] Implementar JWT Service Interface y TokenClaims

**Descripción:** Definir contrato de JWT en dominio (interfaz), implementación en infraestructura.

**Subtareas:**
- [ ] Value Object `TokenClaims` con standard JWT claims (sub, iat, exp, email, role)
- [ ] Interface `JWTService` con métodos GenerateAccessToken, GenerateRefreshToken, ValidateToken
- [ ] Tests para TokenClaims parsing y validación
- [ ] NO incluir implementación real (solo interfaz)

**Salida esperada:**
```
backend/internal/domain/security/
├── jwt.go            (≈80 líneas)
└── jwt_test.go       (≈100 líneas)
```

**Validación:** `go test ./internal/domain/security/... -v` → 100% pasan, cobertura ≥80%

---

### 3. [ ] Implementar Login Use Case (Aplicación)

**Descripción:** Orquestar flujo de autenticación (buscar usuario → validar password → generar tokens).

**Subtareas:**
- [ ] Use Case `LoginUseCase` en application/auth/
- [ ] DTOs: `LoginInput`, `LoginOutput`, `UserDTO`
- [ ] Lógica: validar input → buscar usuario → comparar password → generar tokens
- [ ] Tests integración con mock UserRepository
- [ ] Errores tipados (no strings)

**Salida esperada:**
```
backend/internal/application/auth/
├── login_use_case.go (≈100 líneas)
├── login_test.go     (≈150 líneas)
└── auth_dto.go       (≈50 líneas)
```

**Validación:** `go test ./internal/application/auth/... -v` → 100% pasan, cobertura ≥80%

---

### 4. [ ] HTTP Handler para POST /api/auth/login

**Descripción:** Endpoint REST que ejecuta LoginUseCase y retorna tokens.

**Subtareas:**
- [ ] DTOs de entrada/salida (LoginRequest, LoginResponse)
- [ ] Handler: `HandleLogin` en interfaces/http/handlers/
- [ ] Validación Gin binding (email, password)
- [ ] Error handling: 400, 401, 500
- [ ] Response: `{ access_token, refresh_token, user_id }`

**Salida esperada:**
```
backend/internal/interfaces/http/
├── handlers/
│   └── auth_handler.go     (≈80 líneas)
└── middleware/
    └── auth_middleware.go  (placeholder for next session)
```

**Validación:** 
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
# Expected: 200 + tokens, or 401/400 + error
```

---

### 5. [ ] Tests: Unidad, Integración, E2E (Cobertura Completa)

**Descripción:** Ejecutar suite de tests completamente (34 casos definidos en Fase 2).

**Subtareas:**
- [ ] Unit tests: Email, Password, User, TokenClaims (28 tests)
- [ ] Integration tests: LoginUseCase with mocks (6 tests)
- [ ] Coverage: Mínimo 85% en domain/, 80% en application/
- [ ] E2E: HTTP POST /api/auth/login happy path

**Validación:**
```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out  # Review in browser
```

**Acceptance Criteria:**
- ✅ All 34 test cases passing
- ✅ Coverage ≥85% overall
- ✅ No warnings from golangci-lint
- ✅ go fmt clean
- ✅ docker-compose up without errors

---

## 📊 ANÁLISIS ARQUITECTÓNICO FASE 1 (Resumen)

### Bounded Context
- **Name:** Auth/Security
- **Core Domain:** Authentication & Authorization
- **Root Aggregate:** User
- **Value Objects:** Email, Password, TokenClaims
- **Criticality:** ⭐⭐⭐⭐⭐

### Dependency Map
```
All Modules ↓
    ├── Sales (depends on Auth)
    ├── Product (depends on Auth)
    ├── Production/MES (depends on Auth)
    ├── Supplier/Party (depends on Auth)
    └── Finance (depends on Auth)
         ↓
    Auth/Security (NO external deps - foundational)
         ↓
    Infrastructure (Database, JWT, Crypto)
```

### Domain Model UML (Conceptual)
```
┌─────────────────────────────────────┐
│           User Entity               │ (Root Aggregate)
├─────────────────────────────────────┤
│ - ID: UUID                          │
│ - Email: Email (VO)                 │
│ - Password: Password (VO)           │
│ - Role: Role (enum)                 │
│ - Active: bool                      │
│ - CreatedAt: timestamp              │
│ - UpdatedAt: timestamp              │
└─────────────────────────────────────┘
        ↓ manages ↓
┌─────────────────────────────────────┐
│      TokenClaims (VO)               │
├─────────────────────────────────────┤
│ - Subject: string (user ID)         │
│ - Email: string                     │
│ - Role: string                      │
│ - IssuedAt: time                    │
│ - ExpiresAt: time                   │
│ + IsExpired() bool                  │
└─────────────────────────────────────┘

Repositories (Domain Contracts):
  interface UserRepository {
    ByID(ctx, id) (*User, error)
    ByEmail(ctx, email) (*User, error)
    Save(ctx, user) error
    Delete(ctx, id) error
  }

Services (Domain Contracts):
  interface JWTService {
    GenerateAccessToken(*TokenClaims) (string, error)
    GenerateRefreshToken(*TokenClaims) (string, error)
    ValidateToken(token) (*TokenClaims, error)
  }
```

---

## 📝 NOTAS DE DESARROLLO

### Bcrypt Requirements
```go
cost ≥ 10              // Security requirement (computational cost)
password 8-72 chars    // Input constraint
hash never plaintext   // Security invariant
hashing time: ~200ms   // Performance expectation
```

### JWT Configuration
```go
AccessToken TTL:  900 seconds (15 minutes)
RefreshToken TTL: 604800 seconds (7 days)
Algorithm:        HS256 (HMAC-SHA256)
Secret:           Env var JWT_SECRET (min 32 chars recommended)
```

### RFC 5322 Email Validation (Simplified)
```go
// Simplified validation (not full RFC 5322 - that's complex)
// Check: local@domain
// - local: alphanumeric + . + - + _
// - domain: alphanumeric + . with at least one dot
// - length: ≤254 chars total
```

### Error Handling Pattern
```go
// Domain errors (exported, no HTTP codes)
var (
    ErrInvalidEmail = errors.New("invalid email format")
    ErrInvalidPassword = errors.New("invalid password")
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
)

// HTTP handler maps to HTTP codes
func handleLoginError(err error) (code int, msg string) {
    switch {
    case errors.Is(err, ErrUserNotFound):
        return http.StatusUnauthorized, "invalid credentials"
    case errors.Is(err, ErrInvalidCredentials):
        return http.StatusUnauthorized, "invalid credentials"
    default:
        return http.StatusInternalServerError, "internal error"
    }
}
```

---

## ✅ PREREQUISITOS SATISFECHOS (from Session-09a)

✅ Go project scaffold complete  
✅ Clean Architecture base directories created  
✅ Docker stack ready (PostgreSQL + API)  
✅ Makefile targets verified  
✅ Error handling base defined  
✅ Configuration management ready  
✅ Development environment ready  

---

## 🚀 COMENCEMOS: FASE 3 – IMPLEMENTACIÓN

**Next Action:** Implementar Value Object `Email` con tests TDD-first

```bash
# Run tests (will fail until implementation)
go test ./internal/domain/user/... -v

# Then implement to make tests pass
```

---

**Referencia:** 
- [ADR-002 - Clean Architecture & DDD](../adr/ADR-002-adopcion-clean-architecture-ddd.md)
- [ADR-004 - MVP Lifecycle](../adr/ADR-004-ciclo-vida-desarrollo-mvp.md)
- [Module Template - Auth requirements](../modules/auth.md)
