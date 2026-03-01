# Checklist de Cumplimiento de Seguridad

Este checklist es una herramienta para validar que los módulos y componentes del proyecto TramaTex cumplen con la estrategia de seguridad definida en [ADR-010: Estrategia de Seguridad](../../architecture/adrs/adr-010-defense-in-depth-security-strategy.md) y su [Guía de Implementación de Seguridad](../developer/security-implementation-guide.md).

---

## Checklist de Cumplimiento

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
