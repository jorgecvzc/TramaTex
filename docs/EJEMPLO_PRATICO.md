# EJEMPLO PRÁCTICO: CÓMO USAR LOS PROMPTS

**Fecha:** 11/01/2026  
**Propósito:** Demostración paso a paso de cómo iniciar una sesión con GitHub Copilot usando los nuevos prompts

---

## 📋 ESCENARIO

Eres Jorge y quieres iniciar **Sesión 09** para implementar autenticación JWT.

---

## 🚀 PASO A PASO

### PASO 1: Prepara los documentos (5 minutos)

**Abre en VS Code:**
1. `/docs/SESSION_PROMPT_QUICK.md` ← Refresca contexto en 5 min
2. `/docs/COPILOT_INSTRUCTIONS.md` ← Expectativas para Copilot
3. `/docs/SESSION_INIT_TEMPLATE.md` ← Template para esta sesión

**¿Por qué?**
- SESSION_PROMPT_QUICK.md te recuerda el proyecto en 5 minutos
- COPILOT_INSTRUCTIONS.md te recuerda cómo debe ser Copilot
- SESSION_INIT_TEMPLATE.md es donde documentarás esta sesión

### PASO 2: Crea archivo de sesión (2 minutos)

Copia `/docs/SESSION_INIT_TEMPLATE.md` en nueva sesión:

```bash
# En terminal, ubicado en tramatex/
cp docs/SESSION_INIT_TEMPLATE.md docs/sessions/2026-01-XX-session-09.md
```

**Resultado:** Nuevo archivo `/docs/sessions/2026-01-XX-session-09.md` listo para rellenar

### PASO 3: Rellena información básica (5 minutos)

Abre `/docs/sessions/2026-01-XX-session-09.md` y completa sección inicial:

```markdown
## 📋 INFORMACIÓN DE LA SESIÓN

| Campo | Valor |
|-------|-------|
| **Sesión** | session-09 (completa: 2026-01-18-session-09) |
| **Facilitador/LLM** | GitHub Copilot |
| **Fecha Inicio** | 2026-01-18 |
| **Hora Inicio** | 09:00 |
| **Duración Estimada** | 8 horas |
| **Participantes** | Jorge Cortés Villalba, GitHub Copilot |
```

### PASO 4: Define objetivos (10 minutos)

Rellena sección "OBJETIVOS PRINCIPALES":

```markdown
## 🎯 OBJETIVOS PRINCIPALES

1. [ ] **Implementar entidad User en dominio**
   - Value Object para Password (hashing bcrypt)
   - Value Object para Email (validación)
   - Entity User con invariantes
   - Tests unitarios (100% cobertura)

2. [ ] **Implementar JWT generation y validation**
   - Service para generar tokens (access + refresh)
   - Service para validar tokens
   - Claims estándar (sub, iat, exp)
   - Tests unitarios

3. [ ] **Implementar caso de uso Login**
   - Use case con orquestación
   - Integración con UserRepository
   - Manejo de errores tipado
   - Tests integración

4. [ ] **Crear handler HTTP para Login**
   - Endpoint POST /api/auth/login
   - DTOs entrada/salida
   - Validación input
   - Error handling REST

5. [ ] **Setup Docker con JWT_SECRET**
   - Variable de entorno JWT_SECRET
   - Docker-compose actualizado
   - .env.example con template
```

Estos son los objetivos reales para Sesión 09.

### PASO 5: Contexto de entrada (5 minutos)

Rellena "CONTEXTO DE ENTRADA":

```markdown
## 📊 CONTEXTO DE ENTRADA

### Estado Anterior
- Última sesión completada: 2026-01-11-session-08
- Cambios desde última sesión: Documentación SESSION_PROMPT.md generada
- Estado en PROJECT_STATUS.md: Fase 0 (Fundaciones), setup completado

### Bloqueadores/Dependencias
- [x] Ninguno identificado, proceder

### Prioridades Esta Sesión
**Crítica (Must Have):**
- Entidad User con Password hash
- JWT generation y validation
- Handler Login funcional

**Alta (Should Have):**
- Tests >75% cobertura
- Error handling consistente

**Media (Nice to Have):**
- Frontend Login component (puede ser Sesión 10)
```

### PASO 6: Obtén contexto completo (20 minutos)

**IMPORTANTE:** Antes de compartir con Copilot, tienes dos opciones:

**Opción A: Comparte SESSION_PROMPT.md directamente**

Mensaje a Copilot:
```
Voy a iniciar una nueva sesión de desarrollo en TramaTex.

Lee este contexto completo primero:
[PEGA CONTENIDO DE SESSION_PROMPT.md AQUÍ]

Después de leerlo, confirma que entiendes:
1. El proyecto es un ERP/MES para microempresas textil
2. Usa Go + Vue.js 3 + PostgreSQL
3. Arquitectura es Clean Architecture + DDD
4. Dominio crítico requiere rigor estricto
5. TDD obligatorio en dominio

¿Confirmás?
```

**Opción B: Comparte SESSION_PROMPT_QUICK.md (más corto)**

Mensaje a Copilot:
```
Contexto rápido del proyecto TramaTex:

[PEGA SESSION_PROMPT_QUICK.md]

¿Entiende los puntos clave?
```

**¿Cuál elegir?**
- Sesión 1: Usa Opción A (contexto completo)
- Sesiones 2+: Usa Opción B (refrescar rápido)

### PASO 7: Comparte instrucciones (5 minutos)

Mensaje a Copilot:
```
Ahora lee mis instrucciones para cómo quiero que colabores:

[PEGA CONTENIDO DE COPILOT_INSTRUCTIONS.md AQUÍ]

En particular, estos puntos son críticos:
- Tests PRIMERO (TDD) en dominio
- Clean Architecture: dependencias siempre hacia adentro
- Domain sin ORM, sin frameworks
- Commits descriptivos

¿Tienes claro cómo debo verificar tu trabajo?
```

### PASO 8: Comparte objetivo de sesión (5 minutos)

Mensaje a Copilot:
```
Para esta sesión, tenemos estos objetivos:

[PEGA SECCIÓN "OBJETIVOS PRINCIPALES" del archivo sesión-09.md]

¿Confirmas que entiende qué construiremos?

Recordá:
- Empezamos con TESTS (TDD)
- User entity en domain/ SIN dependencias
- JWT Service también en domain/
- Integration tests para use cases
- HTTP handlers al final
```

### PASO 9: Comienza desarrollo (6+ horas)

Copilot comienza con:

**1. Tests para User entity:**
```go
// backend/tests/unit/domain/user/user_test.go
func TestUserCreation_Valid(t *testing.T) {
    // Test sin dependencias
    email, _ := domain.NewEmail("user@example.com")
    password, _ := domain.NewPassword("Password123!")
    
    user, err := domain.NewUser("user-1", email, password)
    assert.NoError(t, err)
    assert.Equal(t, "user-1", user.ID)
}
```

**2. Implementa User entity:**
```go
// backend/internal/domain/user/user.go
type User struct {
    ID       string
    Email    Email
    Password Password
}

func NewUser(id string, email Email, password Password) (*User, error) {
    if id == "" {
        return nil, errors.New("id required")
    }
    return &User{ID: id, Email: email, Password: password}, nil
}
```

**3. Tests JWT Service:**
```go
// backend/tests/unit/domain/security/jwt_test.go
func TestGenerateToken_Valid(t *testing.T) {
    service := security.NewJWTService("secret-key")
    user := createTestUser()
    
    token, err := service.GenerateAccessToken(user)
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}
```

**4. Implementa JWT Service:**
```go
// backend/internal/domain/security/jwt.go
type JWTService struct {
    secret string
}

func (s *JWTService) GenerateAccessToken(user *User) (string, error) {
    // Implementación
}
```

**5. Tests Use Case:**
```go
// backend/tests/unit/application/auth/login_test.go
func TestLoginUseCase_CorrectCredentials(t *testing.T) {
    // Mock repository
    mockRepo := NewMockUserRepository()
    service := NewLoginUseCase(mockRepo, jwtService)
    
    result, err := service.Execute(ctx, "user@example.com", "password")
    assert.NoError(t, err)
    assert.NotEmpty(t, result.AccessToken)
}
```

**6. Implementa Use Case:**
```go
// backend/internal/application/auth/login.go
func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (*LoginOutput, error) {
    // Buscar usuario
    user, err := uc.userRepo.ByEmail(ctx, email)
    if err != nil {
        return nil, domain.NewAuthError("invalid credentials")
    }
    
    // Validar password
    if !user.Password.Matches(password) {
        return nil, domain.NewAuthError("invalid credentials")
    }
    
    // Generar tokens
    accessToken, _ := uc.jwtService.GenerateAccessToken(user)
    
    return &LoginOutput{AccessToken: accessToken}, nil
}
```

**7. HTTP Handler:**
```go
// backend/internal/interfaces/http/auth_controller.go
func HandleLogin(uc *LoginUseCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input LoginInput
        if err := c.BindJSON(&input); err != nil {
            c.JSON(400, ErrorResponse{Message: "invalid input"})
            return
        }
        
        output, err := uc.Execute(c.Request.Context(), input.Email, input.Password)
        if err != nil {
            c.JSON(401, ErrorResponse{Message: "invalid credentials"})
            return
        }
        
        c.JSON(200, output)
    }
}
```

### PASO 10: Valida antes de commit (5 minutos)

```bash
# Tests pasan
cd backend
go test ./...

# Lint sin warnings
golangci-lint run ./...

# Coverage >75% en critical
go test -cover ./internal/domain/...

# Docker funciona
docker-compose up -d
curl http://localhost:8080/api/health
```

### PASO 11: Commits descriptivos

```bash
git add backend/internal/domain/user/
git commit -m "[feat]: Implement User entity in domain with password hashing

- Created User entity with invariants
- Added Email and Password value objects
- Implemented bcrypt hashing for password
- 100% test coverage with unit tests
- Follows DDD and Clean Architecture principles"

git add backend/internal/domain/security/
git commit -m "[feat]: Implement JWT service for token generation

- Created JWTService for access/refresh tokens
- Implemented token validation with claims
- Added expiration handling
- Tests cover normal and edge cases"

git add backend/internal/application/auth/
git commit -m "[feat]: Implement Login use case with orchestration

- Created LoginUseCase with user/password validation
- Returns access token on success
- Integration tests with mocked repositories"

git add backend/internal/interfaces/http/
git commit -m "[feat]: Add HTTP handler for login endpoint

- POST /api/auth/login with email/password
- Returns token and user info
- Proper error handling (401, 400)
- Input validation"

git add docker/
git commit -m "[chore]: Add JWT_SECRET configuration to docker-compose

- New environment variable JWT_SECRET
- Updated .env.example with template
- Documentation for configuration"
```

### PASO 12: Documenta sesión (15 minutos)

Rellena sección "CHANGES MADE" en `/docs/sessions/2026-01-XX-session-09.md`:

```markdown
## 📝 CHANGES MADE

### Commits Realizados
1. [feat]: Implement User entity in domain with password hashing
2. [feat]: Implement JWT service for token generation
3. [feat]: Implement Login use case with orchestration
4. [feat]: Add HTTP handler for login endpoint
5. [chore]: Add JWT_SECRET configuration to docker-compose
6. [docs]: Update PROJECT_STATUS.md with Phase 0.2 progress

### Archivos Modificados
| Archivo | Tipo | Descripción |
|---------|------|------------|
| backend/internal/domain/user/user.go | NEW | User entity |
| backend/internal/domain/user/password.go | NEW | Password value object |
| backend/internal/domain/user/email.go | NEW | Email value object |
| backend/internal/domain/security/jwt.go | NEW | JWT service |
| backend/internal/application/auth/login_use_case.go | NEW | Login use case |
| backend/internal/interfaces/http/auth_controller.go | NEW | Login handler |
| backend/tests/unit/domain/user/user_test.go | NEW | User tests |
| backend/tests/unit/domain/security/jwt_test.go | NEW | JWT tests |
| backend/tests/unit/application/auth/login_test.go | NEW | Use case tests |
| docker/docker-compose.yml | MODIFIED | Added JWT_SECRET |
| docs/sessions/2026-01-XX-session-09.md | NEW | This session |

### Métricas de Cambio
| Métrica | Valor |
|---------|-------|
| Archivos creados | 9 |
| Archivos modificados | 1 |
| Líneas de código | ~600 |
| Tests agregados | ~40 |
| Commits | 6 |
```

### PASO 13: Actualiza PROJECT_STATUS.md

En `/PROJECT_STATUS.md`, actualiza:

```markdown
## ⏳ En Progreso (Próxima Sesión)

### Fase 0.2: Autenticación y Setup (Sesión 09 - En Curso)
- [x] Skeleton Go con Clean Architecture
- [x] Setup Vue.js 3 + Vite
- [x] Entidad User en dominio
- [x] JWT (generación y validación)
- [ ] Componente Login.vue
- [ ] Docker Compose básico
- [ ] Tests iniciales

## Métricas

| Métrica | Valor |
|---------|-------|
| **Horas invertidas** | 12h / 782h totales |
| **Porcentaje del proyecto** | 1.5% |
| **Commits** | 8 |
| **Líneas de código** | ~600 (dominio + aplicación) |
```

---

## 🎯 RESULTADO FINAL

**Sesión 09 completada:**
- ✅ User entity con Password hashing
- ✅ JWT generation/validation
- ✅ Login use case
- ✅ HTTP handler
- ✅ Docker configurado
- ✅ Tests >75% cobertura
- ✅ Documentación sesión completa

**Próxima sesión (Sesión 10):**
- Componente Login.vue
- Integración frontend-backend
- Tests E2E básicos

---

## 💡 TIPS IMPORTANTES

1. **Copilot acelera con contexto:** Cuanto más detalles en SESSION_PROMPT.md, mejor
2. **SESSION_INIT_TEMPLATE.md es tu aliado:** Rellena conforme avanzas
3. **Tests PRIMERO:** TDD es no-negociable en dominio
4. **Commits pequeños:** Uno por feature, descriptivo
5. **Documenta decisiones:** "Por qué" es tan importante como "qué"

---

## 🔄 PRÓXIMA SESIÓN

Reutiliza el template:
```bash
cp docs/SESSION_INIT_TEMPLATE.md docs/sessions/2026-01-XX-session-10.md
```

Y repite este proceso.

---

**FIN DEL EJEMPLO PRÁCTICO**

¿Dudas? Consulta:
- `SESSION_PROMPT.md` → contexto completo
- `COPILOT_INSTRUCTIONS.md` → expectativas código
- `README_PROMPTS.md` → guía general

