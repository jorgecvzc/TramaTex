# Estrategia de Seguridad TramaTex

**Referencia completa y autoritativa:** [ADR-010: Estrategia de Seguridad - Defensa en Profundidad](adr/ADR-010-estrategia-seguridad-defensa-profundidad.md)

---

## Principios Fundamentales

- **Defensa en Profundidad:** 6 capas de seguridad independientes
- **Security by Default:** Configuración segura sin intervención del usuario
- **Principio de Menor Privilegio:** Acceso mínimo necesario por defecto
- **Zero Trust:** No confiar en ninguna entrada; validar todo
- **Fail Secure:** Ante error, denegar acceso (no conceder)

---

## Arquitectura en Capas

### 1️⃣ Capa de Red y Perímetro
- **TLS 1.3:** Cifrado en tránsito obligatorio (post-inicial)
- **CORS:** Whitelist de orígenes permitidos (configurable por entorno)
- **Rate Limiting:** Protección contra fuerza bruta (post-inicial)
- **WAF:** Web Application Firewall si despliegue en cloud (post-inicial)

### 2️⃣ Capa de Identidad y Acceso (IAM)
- **Autenticación:** JWT con refresh tokens
- **Contraseñas:** Min 8 caracteres + número + mayúscula + minúscula, bcrypt, rotación 365 días
- **RBAC:** 3 roles (admin, manager, operator)
- **Gestión de secretos:** Variables de entorno, sin hardcoding, validación CI/CD

### 3️⃣ Capa de Aplicación
- **Validación dual:** Interfaces (formato) + Domain (negocio)
- **Prevención:** SQL Injection (GORM), XSS (sin v-html), Path Traversal (whitelist)
- **Headers seguros:** X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, HSTS
- **Errores:** Genéricos en producción, detallados solo en logs

### 4️⃣ Capa de Datos
- **Cifrado tránsito:** TLS 1.3 (post-inicial)
- **Cifrado reposo:** AES-256 para datos sensibles (post-inicial)
- **Privacidad logs:** Enmascarar emails, nunca contraseñas/tokens

### 5️⃣ Capa de Integridad y Supply Chain
- **Auditoría:** `go list | nancy`, `npm audit` en CI/CD
- **Dependabot:** Updates supervisados
- **SBOM:** Generado automáticamente
- **Escaneo:** Imágenes Docker (Trivy/Grype)
- **Build:** Tags específicos, checksums, GPG signing opcional

### 6️⃣ Capa de Operaciones y Monitoreo
- **Logging estructurado:** logrus (JSON), requestID para correlación
- **Eventos críticos:** Login fallido, 403, cambios roles, eliminación datos, exportación
- **Almacenamiento:** Rotación 7 días dev / 90 días prod
- **Alertas:** Manual en fase inicial; automáticas post-inicial

---

## Checklist de Cumplimiento por Módulo

Para validar que un módulo cumple la estrategia de seguridad:

**Autenticación y Autorización:**
- [ ] Endpoints protegidos con `AuthMiddleware`
- [ ] Endpoints de escritura verifican rol con `RequireRole()`
- [ ] Use cases validan ownership cuando aplica

**Validación de Entradas:**
- [ ] Validación en interfaces layer (formato, tipo, longitud)
- [ ] Validación en domain layer (reglas de negocio)
- [ ] Whitelist de caracteres en campos de texto libre

**Logging de Seguridad:**
- [ ] Eventos críticos loggeados (login, 403, cambios roles, eliminación, exportación)
- [ ] Structured logging con requestID
- [ ] Sin información sensible en logs (contraseñas, tokens, datos sin enmascarar)

**Gestión de Secretos:**
- [ ] Sin secretos hardcodeados en código
- [ ] Variables de entorno documentadas en README
- [ ] Script de validación `.env` implementado

**Dependencias:**
- [ ] Auditoría automática en CI/CD
- [ ] SBOM generado
- [ ] Imágenes Docker escaneadas

**Configuración Segura:**
- [ ] 4 headers HTTP seguros implementados
- [ ] Endpoints debug deshabilitados en producción
- [ ] CORS con whitelist configurada
- [ ] Headers de versión removidos

---

## Roles y Permisos

| Rol | Lectura | Escritura | Eliminación | Gestión Usuarios |
|-----|---------|-----------|-------------|------------------|
| **admin** | ✅ Todo | ✅ Todo | ✅ Todo | ✅ Sí |
| **manager** | ✅ Todo | ✅ Sí | ❌ No | ❌ No |
| **operator** | ✅ Todo | ❌ No | ❌ No | ❌ No |

**Configuración por defecto:**
- Nuevas cuentas: `active: false` (activación manual)
- Acceso: organizacional completo

---

## Eventos de Seguridad a Loggear

| Evento | Nivel | Campos Obligatorios |
|--------|-------|---------------------|
| Login exitoso | INFO | userID, email, IP, timestamp, user-agent |
| Login fallido | WARN | email, IP, timestamp, reason |
| Acceso denegado (403) | WARN | userID, resource, action, timestamp |
| Creación recurso crítico | INFO | userID, resourceType, resourceID, timestamp |
| Modificación recurso crítico | INFO | userID, resourceType, resourceID, changes, timestamp |
| Eliminación datos | WARN | userID, resourceType, resourceID, timestamp |
| Cambio roles/permisos | WARN | adminID, targetUserID, oldRole, newRole, timestamp |
| Exportación datos | WARN | userID, exportType, recordCount, timestamp |
| Error interno (500) | ERROR | endpoint, error, stackTrace, timestamp |

---

## Referencias Rápidas

- **ADR completo:** [ADR-010](adr/ADR-010-estrategia-seguridad-defensa-profundidad.md)
- **OWASP Top 10 2021:** https://owasp.org/Top10/
- **OWASP Cheat Sheets:** https://cheatsheetseries.owasp.org/
- **Go Security:** https://go.dev/doc/security/best-practices
- **Auditoría OWASP TramaTex:** `docs/project/milestones/auditoria-seguridad-owasp-2026-01-25.md`

---

**Última actualización:** 2026-01-26  
**Próxima revisión:** Anual o ante cambio arquitectónico
