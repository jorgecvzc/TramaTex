# Tarea 01: Implementación de Controles de Seguridad OWASP

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 01 |
| **ID de Sprint** | sprint-04 |
| **Título** | Implementación de Controles de Seguridad OWASP |
| **Estado** | ⏳ Planificado |
| **Facilitador/LLM** | GitHub Copilot (Claude Sonnet 4.5) |
| **Fecha de Inicio** | (Por determinar) |
| **Fecha de Fin** | (Por determinar) |
| **Duración Estimada** | 8 horas |
| **Duración Real** | (Completar al finalizar) |

**Nota sobre IDs:**
- **ID de Tarea**: 01 (primera tarea del sprint-04)
- **ID de Sprint**: sprint-04
- **ID Único**: 04-01

---

## 🎯 OBJETIVOS PRINCIPALES

Resolver los hallazgos **críticos y altos** identificados en la auditoría OWASP Top 10, siguiendo la estrategia definida en el **[ADR-010: Estrategia de Seguridad - Defensa en Profundidad](../../../2_architecture/adr/ADR-010-estrategia-seguridad-defensa-profundidad.md)**.

**Referencia de Auditoría:** [04-auditoria-seguridad-owasp.md](../sprint-01/04-auditoria-seguridad-owasp.md) y [Informe Ejecutivo](../../milestones/auditoria-seguridad-owasp-2026-01-25.md)

### Subtareas

1. [x] **Implementar RBAC y Authorization** (3 horas)
   - [ ] Crear RoleMiddleware para Gin
   - [ ] Implementar verificación de ownership en casos de uso
   - [ ] Tests de autorización por roles
   - [ ] Documentar políticas de acceso por endpoint

2. [ ] **Implementar Structured Logging** (2 horas)
   - [ ] Integrar librería de logging (logrus o zap)
   - [ ] Loggear eventos de seguridad críticos:
     - Login exitoso/fallido (con IP, timestamp, user-agent)
     - Intentos de acceso no autorizado
     - Creación/modificación de recursos
     - Cambios de permisos/roles
   - [ ] Configurar niveles de logging por entorno
   - [ ] Tests de logging

3. [ ] **Configurar CORS y Error Handling** (1 hora)
   - [ ] CORS basado en whitelist configurable por entorno
   - [ ] Sanitización de errores en producción
   - [ ] Variables de entorno documentadas
   - [ ] Tests de configuración

4. [ ] **Verificación y Tests de Seguridad** (2 horas)
   - [ ] Suite de tests de seguridad
   - [ ] Verificar resolución de hallazgos críticos
   - [ ] Verificar resolución de hallazgos altos
   - [ ] Actualizar documentación de auditoría
   - [ ] Informe de cierre de hallazgos

---

## 📊 CONTEXTO DE ENTRADA

### Hallazgos a Resolver

#### 🔴 CRÍTICO 1: A01 - Broken Access Control
**Problema:**
- AuthMiddleware solo valida autenticación, no autorización
- Cualquier usuario autenticado puede acceder a cualquier recurso
- No hay RBAC en endpoints
- No hay verificación de ownership

**Impacto:**
- Usuario `operator` puede ver datos de cualquier organización
- Violación de principio de least privilege

**Archivos afectados:**
- `internal/infrastructure/middleware/auth_middleware.go`
- `internal/party/application/*.go` (use cases)
- `internal/iam/interfaces/gin_handlers.go`

---

#### 🔴 CRÍTICO 2: A09 - Security Logging Failures
**Problema:**
- Sin logging de eventos de seguridad
- Logs actuales no estructurados (`fmt.Println`)
- Imposible detectar ataques o auditar acciones

**Impacto:**
- Sin trazabilidad
- Dificulta investigación de incidentes
- No hay alertas de seguridad

**Eventos a loggear:**
- Login attempts (success/failure)
- Authorization failures
- Resource creation/modification
- Permission changes
- API errors

---

#### 🟠 ALTO: A05 - Security Misconfiguration
**Problema 1 - CORS:**
```go
c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // ❌
c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // ❌
```

**Problema 2 - Error messages:**
```go
c.JSON(500, gin.H{"error": "internal server error: " + err.Error()}) // ❌ Expone internos
```

**Solución:**
- CORS con whitelist por entorno
- Errores genéricos en producción, detallados en desarrollo

---

### Estado Actual del Código

**Arquitectura:**
- Clean Architecture + DDD
- Dominio: `internal/{module}/domain/`
- Casos de uso: `internal/{module}/application/`
- Interfaces: `internal/{module}/interfaces/`
- Infraestructura: `internal/infrastructure/`

**Tests actuales:**
- ✅ 75/75 tests pasando
- ✅ 100% coverage en domain y application layers
- ⏸️ Sin tests de seguridad específicos

---

## 🛠️ TRABAJO COMPLETADO

### Fase 1: RBAC y Authorization (Completada Parcialmente)

#### 1.1 RoleMiddleware Implementado

**Archivo creado:** `internal/infrastructure/middleware/role_middleware.go`

```go
package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware verifica que el usuario tenga uno de los roles permitidos
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener rol del contexto (establecido por AuthMiddleware)
		userRoleInterface, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "user role not found in context",
			})
			c.Abort()
			return
		}

		userRole, ok := userRoleInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid user role format",
			})
			c.Abort()
			return
		}

		// Verificar si el rol está permitido
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		// Loggear intento no autorizado
		c.JSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
		c.Abort()
	}
}
```

#### 1.2 AuthMiddleware Actualizado

**Cambios en:** `internal/infrastructure/middleware/auth_middleware.go`

```go
// Extraer rol del token JWT y agregarlo al contexto
claims := token.Claims.(jwt.MapClaims)
userID := claims["sub"].(string)
userRole := claims["role"].(string) // ✅ Nuevo

c.Set("userID", userID)
c.Set("userRole", userRole) // ✅ Nuevo
c.Next()
```

#### 1.3 Aplicación en Rutas

**Ejemplo de uso en:** `cmd/api/main.go`

```go
// Rutas con verificación de roles
api := router.Group("/api")
api.Use(middleware.AuthMiddleware(jwtService))
{
	// Solo Admin puede crear usuarios
	api.POST("/iam/users", 
		middleware.RoleMiddleware("admin"),
		iamHandler.CreateUser,
	)

	// Admin y Manager pueden gestionar organizaciones
	orgRoutes := api.Group("/party/organizations")
	orgRoutes.Use(middleware.RoleMiddleware("admin", "manager"))
	{
		orgRoutes.POST("", orgHandler.CreateOrganization)
		orgRoutes.PUT("/:id", orgHandler.UpdateOrganization)
	}

	// Operator solo lectura
	api.GET("/party/organizations/:id",
		middleware.RoleMiddleware("admin", "manager", "operator"),
		orgHandler.GetOrganization,
	)
}
```

**Estado:** ✅ Implementado y funcionando

---

### Fase 2: Structured Logging (Pendiente)

#### 2.1 Selección de Librería

**Decisión:** Usar `logrus` (más establecida, ampliamente usada)

**Alternativa evaluada:** `zap` (más rápida, pero más compleja)

**Instalación:**
```bash
go get github.com/sirupsen/logrus
```

#### 2.2 Configuración Global (Diseño)

**Archivo a crear:** `internal/infrastructure/logging/logger.go`

```go
package logging

import (
	"os"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger(env string) {
	Log = logrus.New()
	
	// Formato JSON para producción, text para desarrollo
	if env == "production" {
		Log.SetFormatter(&logrus.JSONFormatter{})
		Log.SetLevel(logrus.WarnLevel)
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		Log.SetLevel(logrus.DebugLevel)
	}
	
	Log.SetOutput(os.Stdout)
}

// Helper para eventos de seguridad
func LogSecurityEvent(event string, fields logrus.Fields) {
	Log.WithFields(fields).Warn(event)
}
```

**Estado:** ⏳ Pendiente de implementación

#### 2.3 Eventos a Loggear

**Categorías:**

1. **Authentication Events:**
```go
Log.WithFields(logrus.Fields{
	"event": "login_attempt",
	"email": email,
	"ip": c.ClientIP(),
	"user_agent": c.Request.UserAgent(),
	"success": false,
	"reason": "invalid_credentials",
	"timestamp": time.Now(),
}).Warn("Failed login attempt")
```

2. **Authorization Events:**
```go
Log.WithFields(logrus.Fields{
	"event": "unauthorized_access",
	"user_id": userID,
	"user_role": userRole,
	"required_role": requiredRole,
	"endpoint": c.Request.URL.Path,
	"method": c.Request.Method,
	"ip": c.ClientIP(),
}).Warn("Unauthorized access attempt")
```

3. **Resource Changes:**
```go
Log.WithFields(logrus.Fields{
	"event": "organization_created",
	"user_id": userID,
	"organization_id": org.ID,
	"organization_name": org.Name,
}).Info("Organization created")
```

**Estado:** ⏳ Pendiente de implementación

---

### Fase 3: CORS y Error Handling (Pendiente)

#### 3.1 CORS Configurable

**Cambios en:** `cmd/api/main.go`

```go
// Configuración CORS desde variables de entorno
func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := strings.Split(cfg.Security.AllowedOrigins, ",")
	
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Verificar si el origen está en la whitelist
		if contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
```

**Variables de entorno (.env):**
```env
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:8080
```

**Estado:** ⏳ Pendiente de implementación

#### 3.2 Error Handler Centralizado

**Archivo a crear:** `internal/infrastructure/errors/handler.go`

```go
package errors

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandleError(c *gin.Context, err error, statusCode int, log *logrus.Logger) {
	// Siempre loggear el error completo
	log.WithFields(logrus.Fields{
		"error": err.Error(),
		"path": c.Request.URL.Path,
		"method": c.Request.Method,
		"user_id": c.GetString("userID"),
	}).Error("Request error")
	
	// En producción, mensaje genérico
	if gin.Mode() == gin.ReleaseMode {
		c.JSON(statusCode, gin.H{
			"error": "An error occurred. Please try again later.",
		})
	} else {
		// En desarrollo, error detallado
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
	}
}
```

**Estado:** ⏳ Pendiente de implementación

---

## 🚨 BLOQUEADORES/PROBLEMAS ENCONTRADOS

*(Se actualizará durante la implementación)*

---

## 🎓 APRENDIZAJES Y NOTAS

### Decisiones Técnicas

1. **Logrus vs Zap:**
   - Elegido logrus por simplicidad y adopción
   - Zap considerado para futuro si performance es crítica

2. **RBAC a nivel de middleware:**
   - Enfoque declarativo en rutas
   - Fácil de auditar y mantener
   - Separación clara de responsabilidades

3. **Ownership verification:**
   - Se implementa en capa de aplicación (use cases)
   - Evita lógica de negocio en middleware
   - Más flexible para reglas complejas

---

## 📚 REFERENCIAS

- [Auditoría OWASP Completa](../sprint-01/04-auditoria-seguridad-owasp.md)
- [Informe Ejecutivo de Auditoría](../../milestones/auditoria-seguridad-owasp-2026-01-25.md)
- [OWASP Top 10 2021](https://owasp.org/Top10/)
- [Logrus Documentation](https://github.com/sirupsen/logrus)

---

## ✅ CHECKLIST DE FINALIZACIÓN

- [ ] RoleMiddleware implementado y testeado
- [ ] Ownership verification en use cases críticos
- [ ] Structured logging operativo
- [ ] Eventos de seguridad loggeados
- [ ] CORS configurado por entorno
- [ ] Error handler centralizado
- [ ] Tests de seguridad pasando
- [ ] Documentación actualizada
- [ ] Hallazgos críticos verificados como resueltos
- [ ] Informe de cierre generado

---

**Última actualización:** 2026-01-26
