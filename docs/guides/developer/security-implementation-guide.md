# Guía de Implementación de Seguridad: Defensa en Profundidad

Este documento detalla las implementaciones específicas para la estrategia de seguridad de Defensa en Profundidad y Security by Default en el proyecto TramaTex, complementando la decisión arquitectónica definida en [ADR-010: Estrategia de Seguridad](../../architecture/adrs/ADR-010-defense-in-depth-security-strategy.md).

---

## 1. Arquitectura de Seguridad por Capas

### Capa 1: Red y Perímetro

**Controles implementados:**
-   **CORS:** Whitelist de orígenes permitidos (configurable por entorno)
-   **TLS 1.3:** Cifrado en tránsito obligatorio (excepto desarrollo local)
-   **Rate Limiting:** Protección contra ataques de fuerza bruta (post-MVP)
-   **WAF:** Web Application Firewall si despliegue en cloud (post-MVP)

**Configuración (Ejemplo):**
```yaml
# Producción
CORS_ALLOWED_ORIGINS: "https://app.tramatex.com"
TLS_ENABLED: true
TLS_MIN_VERSION: "1.3"
```

---

### Capa 2: Identidad y Acceso (IAM)

#### 2.1 Autenticación

**Mecanismo:** JSON Web Tokens (JWT) con refresh tokens

**Política de contraseñas:**
-   Mínimo 8 caracteres
-   Al menos 1 número, 1 mayúscula, 1 minúscula
-   Validación contra diccionarios de contraseñas comunes
-   Rotación obligatoria cada 365 días
-   Hashing: bcrypt con salt automático

**Gestión de tokens:**
-   Access token: corta duración (15 minutos)
-   Refresh token: larga duración (7 días)
-   Revocación: stateless en MVP; blacklist en Redis si cloud

**Gestión de secretos:**
-   **Nunca** en código fuente
-   Variables de entorno (`.env` en `.gitignore`)
-   Validación automática en CI/CD (no commits con secretos)
-   Rotación JWT_SECRET: manual en MVP; automática post-MVP

#### 2.2 Autorización (RBAC)

**Roles del sistema:**
-   `admin`: Acceso total (CRUD completo)
-   `manager`: Gestión operativa (crear, leer, actualizar)
-   `operator`: Solo lectura (listar, ver detalles)

**Regla arquitectónica:**  
Todo endpoint de escritura (POST, PUT, PATCH, DELETE) **DEBE** verificar rol antes de ejecutar.

**Configuración por defecto:**
-   Nuevas cuentas: `active: false` (requiere activación manual)
-   Acceso: organizacional completo (sin ownership granular en MVP)
-   Expansibilidad: diseño permite roles personalizables (post-MVP)

---

### Capa 3: Aplicación

#### 3.1 Validación de Entradas (Zero Trust)

**Principio:** No confiar en ninguna entrada; validar en dos capas.

**Interfaces Layer (HTTP Handlers):**
-   Tipo de datos correcto
-   Longitud máxima (prevenir buffer overflow)
-   Formato válido (regex para emails, teléfonos, tax IDs)
-   Whitelist de caracteres permitidos

**Domain Layer:**
-   Reglas de negocio (ej. precio > 0)
-   Relaciones válidas (ej. organización existe)

#### 3.2 Prevención de Inyecciones

| Vulnerabilidad    | Mitigación                                    |
| :---------------- | :-------------------------------------------- |
| **SQL Injection** | GORM prepared statements (automático)         |
| **XSS**           | Prohibido `v-html` en frontend; sanitización automática Vue.js |
| **Path Traversal**| `filepath.Clean()` + whitelist de directorios |
| **LDAP/NoSQL**    | No aplica en MVP; crítico post-MVP            |

#### 3.3 Manejo de Errores

**Producción:**
-   Errores genéricos al cliente: `{"error": "internal server error"}`
-   Detalles completos solo en logs estructurados

**Desarrollo:**
-   Errores detallados para debugging
-   Stack traces visibles

**Headers HTTP Seguros (obligatorios):**
```http
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

**Hardening:**
-   Deshabilitar endpoints `/debug` en producción
-   Remover headers que revelan versión (`Server: Go/1.23`)

---

### Capa 4: Datos

#### 4.1 Cifrado

**En tránsito:**
-   TLS 1.3 obligatorio (post-MVP)
-   Certificados válidos (Let's Encrypt o corporativos)

**En reposo:**
-   Contraseñas: bcrypt (ya implementado)
-   Datos sensibles: cifrado AES-256 (post-MVP para precios, financiero)

#### 4.2 Privacidad en Logs

**Nunca loggear:**
-   Contraseñas (hash o plaintext)
-   Tokens JWT completos (solo últimos 4 caracteres)
-   Números de tarjetas, datos bancarios

**Enmascarar:**
-   Emails: `j***@tramatex.com`
-   IPs: loggear sin hashear (trazabilidad forense)

---

### Capa 5: Integridad y Supply Chain

#### 5.1 Gestión de Dependencias

**Backend (Go):**
-   Auditoría automática: `go list -m all | nancy` en CI/CD
-   Dependabot activado (updates supervisados manualmente)
-   Política: evaluar CVE critical/high en <7 días

**Frontend (Vue.js):**
-   `npm audit` en CI/CD (falla build si critical/high)
-   Actualización de dependencias cada sprint
-   Lock file (`package-lock.json`) versionado

#### 5.2 Integridad de Build

**Controles:**
-   Dockerfiles con tags específicos (❌ `:latest`)
-   Checksums de binarios generados
-   Escaneo de imágenes Docker: Trivy o Grype
-   SBOM (Software Bill of Materials) generado
-   Firma GPG de commits: opcional en MVP; obligatorio post-MVP

---

### Capa 6: Operaciones y Monitoreo

#### 6.1 Logging Estructurado

**Biblioteca:** logrus (balance simplicidad/performance)

**Formato:** JSON estructurado con campos obligatorios:
```json
{
  "timestamp": "2026-01-26T10:30:45Z",
  "level": "WARN",
  "requestID": "uuid-1234",
  "userID": "user-uuid",
  "event": "access_denied",
  "resource": "/api/organizations/123",
  "ip": "192.168.1.100",
  "message": "Insufficient permissions"
}
```

**Eventos críticos (WARN/ERROR):**
-   Login fallido (≥3 intentos)
-   Acceso denegado (403)
-   Cambios de roles/permisos
-   Eliminación de datos
-   Cambios masivos (>100 registros)
-   Exportación de datos

**Eventos informativos (INFO):**
-   Login exitoso
-   Escrituras críticas: pedidos, precios, organizaciones
-   Modificaciones de configuración

**Almacenamiento:**
-   MVP: archivos rotados (7 días dev / 90 días prod)
-   Post-MVP: centralización (ELK, Loki, CloudWatch)

**Alertas:**
-   MVP: revisión manual de logs
-   Post-MVP: automáticas (>10 logins fallidos/min)

#### 6.2 Correlación de Requests

**Request ID:** UUID generado por middleware, propagado en:
-   Logs de toda la cadena de ejecución
-   Headers de respuesta (`X-Request-ID`)
-   Logs de errores

Permite trazar flujo completo de un request para debugging.

---

## 2. Notas de Implementación

### Fases de Madurez

**Fase Inicial (MVP):**
-   Fundaciones: JWT + RBAC + Validación + Logging básico
-   Configuración segura por defecto
-   Sin cifrado en reposo

**Fase Intermedia (Pre-Cloud):**
-   TLS 1.3 obligatorio (post-MVP)
-   Cifrado en reposo para datos sensibles (post-MVP para precios, financiero)
-   Centralización de logs (post-MVP)

**Fase Avanzada (Cloud):**
-   WAF + Rate limiting (post-MVP)
-   Revocación de tokens (Redis) (post-MVP)
-   Alertas automáticas (post-MVP)
-   Roles personalizables (post-MVP)

### Script de Validación de Configuración

Ejemplo de validación pre-startup:

```bash
#!/bin/bash
# validate-env.sh

required_vars=(
  "DB_HOST"
  "DB_PASSWORD"
  "JWT_SECRET"
  "CORS_ALLOWED_ORIGINS"
)

for var in "${required_vars[@]}"; do
  if [ -z "${!var}" ]; then
    echo "ERROR: $var is not set"
    exit 1
  fi
done

# Validar JWT_SECRET longitud mínima
if [ ${#JWT_SECRET} -lt 32 ]; then
  echo "ERROR: JWT_SECRET must be at least 32 characters"
  exit 1
fi

echo "✓ Environment validation passed"
```
