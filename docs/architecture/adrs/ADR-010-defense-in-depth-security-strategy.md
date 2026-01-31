# ADR-010 – Estrategia de Seguridad: Defensa en Profundidad y Security by Default

**Fecha:** 2026-01-26  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, GitHub Copilot  

---

## 1. Contexto

TramaTex es un sistema ERP/MES que gestiona **datos críticos de negocio**: información de clientes y proveedores, catálogos de productos, cálculos de precios, pedidos y flujos de producción. Un compromiso de seguridad podría resultar en:

- Pérdida de confianza del cliente
- Robo de información comercial sensible (precios, márgenes)
- Manipulación de datos financieros
- Cumplimiento normativo comprometido (GDPR)
- Daño reputacional irreparable para microempresas

**Necesidad estratégica:**  
Establecer una arquitectura de seguridad robusta que proteja el sistema mediante múltiples capas de defensa y configuraciones seguras por defecto, alineada con OWASP Top 10 2021 y principios de Zero Trust.

**Restricciones:**
- Equipo pequeño (1-2 desarrolladores)
- Recursos limitados para auditorías externas
- Necesidad de balance entre seguridad y velocidad de desarrollo

---

## 2. Alternativas Consideradas

### Alternativa A – Seguridad Reactiva (Minimal Viable Security)
- **Enfoque:** Implementar solo autenticación básica, resolver vulnerabilidades cuando se detecten
- **Ventajas:** Desarrollo rápido, menor complejidad inicial
- **Desventajas:** Alta probabilidad de brechas, deuda técnica de seguridad, costo elevado de remediación post-facto

### Alternativa B – Seguridad Perimetral Única (Network-Only)
- **Enfoque:** Confiar en firewalls, WAF y seguridad de red como única capa
- **Desventajas:** Si el perímetro se compromete, acceso total; no protege contra amenazas internas

### Alternativa C – Defensa en Profundidad + Security by Default ✅
- **Enfoque:** Múltiples capas de seguridad independientes; configuración segura desde el diseño
- **Ventajas:** Resiliencia ante fallos de una capa; prevención proactiva; reducción de superficie de ataque
- **Desventajas:** Mayor complejidad arquitectónica; overhead de desarrollo (~20-30%)

---

## 3. Decisión Adoptada

**Adoptar Defensa en Profundidad (Defense in Depth) combinada con Security by Design y Security by Default como estrategia arquitectónica transversal a todo el sistema.**

### Principios Fundamentales:

1. **Defensa en Profundidad:** Seguridad en múltiples capas (aplicación, datos, red)
2. **Security by Default:** Configuraciones seguras sin intervención del usuario
3. **Principio de Menor Privilegio:** Acceso mínimo necesario por defecto
4. **Zero Trust:** No confiar en ninguna entrada; validar todo
5. **Fail Secure:** Ante error, denegar acceso (no conceder)

**Justificación:**
- Protege datos críticos de negocio mediante redundancia de controles
- Cumple con mejores prácticas de la industria (OWASP, MITRE ATT&CK)
- Escalable: permite añadir capas adicionales sin rediseño
- Auditable: cada capa puede verificarse independientemente

---

## 4. Consecuencias

### Positivas
- **Resiliencia:** Si una capa falla, otras siguen protegiendo
- **Compliance:** Facilita cumplimiento de normativas (GDPR, ISO 27001)
- **Auditabilidad:** Trazabilidad completa de eventos de seguridad
- **Confianza:** Demuestra compromiso con la seguridad a clientes
- **Reducción de riesgo:** Superficie de ataque minimizada

### Negativas
- **Complejidad:** Mayor número de componentes de seguridad a mantener
- **Performance:** Overhead de validaciones múltiples (~5-10% latencia)
- **Desarrollo:** Tiempo adicional por feature (~20-30%)
- **Curva de aprendizaje:** Equipo debe dominar prácticas de seguridad

---

## 5. Alcance

### Aplica a:
- Backend (Go): API REST, lógica de negocio, acceso a datos
- Frontend (Vue.js): Manejo de autenticación, validación cliente
- Infraestructura: Configuración servidores, bases de datos, CI/CD
- Datos: En tránsito y en reposo

### No aplica (delegado a otras capas):
- Seguridad física de centros de datos
- Protección DDoS a nivel ISP
- Seguridad de dispositivos de usuario final

---

## 6. Arquitectura de Seguridad por Capas

### Capa 1: Red y Perímetro

**Controles implementados:**
- **CORS:** Whitelist de orígenes permitidos (configurable por entorno)
- **TLS 1.3:** Cifrado en tránsito obligatorio (excepto desarrollo local)
- **Rate Limiting:** Protección contra ataques de fuerza bruta (post-MVP)
- **WAF:** Web Application Firewall si despliegue en cloud (post-MVP)

**Configuración:**
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
- Mínimo 8 caracteres
- Al menos 1 número, 1 mayúscula, 1 minúscula
- Validación contra diccionarios de contraseñas comunes
- Rotación obligatoria cada 365 días
- Hashing: bcrypt con salt automático

**Gestión de tokens:**
- Access token: corta duración (15 minutos)
- Refresh token: larga duración (7 días)
- Revocación: stateless en MVP; blacklist en Redis si cloud

**Gestión de secretos:**
- **Nunca** en código fuente
- Variables de entorno (`.env` en `.gitignore`)
- Validación automática en CI/CD (no commits con secretos)
- Rotación JWT_SECRET: manual en MVP; automática post-MVP

#### 2.2 Autorización (RBAC)

**Roles del sistema:**
- `admin`: Acceso total (CRUD completo)
- `manager`: Gestión operativa (crear, leer, actualizar)
- `operator`: Solo lectura (listar, ver detalles)

**Regla arquitectónica:**  
Todo endpoint de escritura (POST, PUT, PATCH, DELETE) **DEBE** verificar rol antes de ejecutar.

**Configuración por defecto:**
- Nuevas cuentas: `active: false` (requiere activación manual)
- Acceso: organizacional completo (sin ownership granular en MVP)
- Expansibilidad: diseño permite roles personalizables (post-MVP)

---

### Capa 3: Aplicación

#### 3.1 Validación de Entradas (Zero Trust)

**Principio:** No confiar en ninguna entrada; validar en dos capas.

**Interfaces Layer (HTTP Handlers):**
- Tipo de datos correcto
- Longitud máxima (prevenir buffer overflow)
- Formato válido (regex para emails, teléfonos, tax IDs)
- Whitelist de caracteres permitidos

**Domain Layer:**
- Reglas de negocio (ej. precio > 0)
- Relaciones válidas (ej. organización existe)

#### 3.2 Prevención de Inyecciones

| Vulnerabilidad | Mitigación |
|----------------|------------|
| **SQL Injection** | GORM prepared statements (automático) |
| **XSS** | Prohibido `v-html` en frontend; sanitización automática Vue.js |
| **Path Traversal** | `filepath.Clean()` + whitelist de directorios |
- **LDAP/NoSQL** | No aplica en MVP; crítico post-MVP

#### 3.3 Manejo de Errores

**Producción:**
- Errores genéricos al cliente: `{"error": "internal server error"}`
- Detalles completos solo en logs estructurados

**Desarrollo:**
- Errores detallados para debugging
- Stack traces visibles

**Headers HTTP Seguros (obligatorios):**
```http
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

**Hardening:**
- Deshabilitar endpoints `/debug` en producción
- Remover headers que revelan versión (`Server: Go/1.23`)

---

### Capa 4: Datos

#### 4.1 Cifrado

**En tránsito:**
- TLS 1.3 obligatorio (post-MVP)
- Certificados válidos (Let's Encrypt o corporativos)

**En reposo:**
- Contraseñas: bcrypt (ya implementado)
- Datos sensibles: cifrado AES-256 (post-MVP para precios, financiero)

#### 4.2 Privacidad en Logs

**Nunca loggear:**
- Contraseñas (hash o plaintext)
- Tokens JWT completos (solo últimos 4 caracteres)
- Números de tarjetas, datos bancarios

**Enmascarar:**
- Emails: `j***@tramatex.com`
- IPs: loggear sin hashear (trazabilidad forense)

---

### Capa 5: Integridad y Supply Chain

#### 5.1 Gestión de Dependencias

**Backend (Go):**
- Auditoría automática: `go list -m all | nancy` en CI/CD
- Dependabot activado (updates supervisados manualmente)
- Política: evaluar CVE critical/high en <7 días

**Frontend (Vue.js):**
- `npm audit` en CI/CD (falla build si critical/high)
- Actualización de dependencias cada sprint
- Lock file (`package-lock.json`) versionado

#### 5.2 Integridad de Build

**Controles:**
- Dockerfiles con tags específicos (❌ `:latest`)
- Checksums de binarios generados
- Escaneo de imágenes Docker: Trivy o Grype
- SBOM (Software Bill of Materials) generado
- Firma GPG de commits: opcional en MVP; obligatorio post-MVP

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
- Login fallido (≥3 intentos)
- Acceso denegado (403)
- Cambios de roles/permisos
- Eliminación de datos
- Cambios masivos (>100 registros)
- Exportación de datos

**Eventos informativos (INFO):**
- Login exitoso
- Escrituras críticas: pedidos, precios, organizaciones
- Modificaciones de configuración

**Almacenamiento:**
- MVP: archivos rotados (7 días dev / 90 días prod)
- Post-MVP: centralización (ELK, Loki, CloudWatch)

**Alertas:**
- MVP: revisión manual de logs
- Post-MVP: automáticas (>10 logins fallidos/min)

#### 6.2 Correlación de Requests

**Request ID:** UUID generado por middleware, propagado en:
- Logs de toda la cadena de ejecución
- Headers de respuesta (`X-Request-ID`)
- Logs de errores

Permite trazar flujo completo de un request para debugging.

---

## 7. Checklist de Cumplimiento

Para validar que un módulo cumple esta estrategia:

**Autenticación y Autorización:**
- [ ] Endpoints protegidos con AuthMiddleware
- [ ] Endpoints de escritura verifican rol (RoleMiddleware)
- [ ] Use cases validan ownership cuando aplica

**Validación de Entradas:**
- [ ] Validación en interfaces layer (formato)
- [ ] Validación en domain layer (negocio)
- [ ] Whitelist de caracteres en campos críticos

**Logging de Seguridad:**
- [ ] Eventos críticos loggeados (login, 403, cambios roles)
- [ ] Structured logging con requestID
- [ ] Sin información sensible en logs

**Gestión de Secretos:**
- [ ] Sin secretos hardcodeados
- [ ] Variables documentadas en README
- [ ] Script de validación `.env`

**Dependencias:**
- [ ] Auditoría automática en CI/CD
- [ ] SBOM generado
- [ ] Imágenes Docker escaneadas

**Configuración Segura:**
- [ ] Headers HTTP seguros implementados
- [ ] Endpoints debug deshabilitados en producción
- [ ] CORS con whitelist

---

## 8. Integración con otros ADRs

- **ADR-002:** Clean Architecture permite aislar concerns de seguridad en infrastructure layer
- **ADR-003:** Modular monolith facilita aplicación consistente de políticas
- **ADR-006:** DDD - Security es transversal a todos los bounded contexts
- **ADR-009:** Estructura de proyecto define ubicación de componentes de seguridad

---

## 9. Notas de Implementación

### Fases de Madurez

**Fase Inicial (MVP):**
- Fundaciones: JWT + RBAC + Validación + Logging básico
- Configuración segura por defecto
- Sin cifrado en reposo

**Fase Intermedia (Pre-Cloud):**
- TLS 1.3 obligatorio (post-MVP)
- Cifrado en reposo para datos sensibles (post-MVP para precios, financiero)
- Centralización de logs (post-MVP)

**Fase Avanzada (Cloud):**
- WAF + Rate limiting (post-MVP)
- Revocación de tokens (Redis) (post-MVP)
- Alertas automáticas (post-MVP)
- Roles personalizables (post-MVP)

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

---

## 10. Referencias

- **OWASP Top 10 2021:** https://owasp.org/Top10/
- **MITRE ATT&CK Framework:** https://attack.mitre.org/
- **OWASP Cheat Sheet Series:** https://cheatsheetseries.owasp.org/
- **NIST Cybersecurity Framework:** https://www.nist.gov/cyberframework
- **Go Security Best Practices:** https://go.dev/doc/security/best-practices
- **Zero Trust Architecture (NIST SP 800-207)**

---

**Aprobado:** 2026-01-26  
**Revisión próxima:** Anual o ante cambio arquitectónico significativo
