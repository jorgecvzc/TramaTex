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
Para una descripción detallada de la arquitectura de seguridad por capas y sus controles implementados, consulte [Guía de Implementación de Seguridad: Defensa en Profundidad](../../guides/developer/security-implementation-guide.md).
---

## 7. Checklist de Cumplimiento
Para un checklist detallado de cumplimiento de seguridad, consulte [Checklist de Cumplimiento de Seguridad](../../guides/developer/security-compliance-checklist.md).

## 8. Integración con otros ADRs

- **ADR-002:** Clean Architecture permite aislar concerns de seguridad en infrastructure layer
- **ADR-003:** Modular monolith facilita aplicación consistente de políticas
- **ADR-006:** DDD - Security es transversal a todos los bounded contexts
- **ADR-009:** Estructura de proyecto define ubicación de componentes de seguridad

---

## 9. Notas de Implementación
Para notas detalladas de implementación, incluidas las fases de madurez y el script de validación de configuración, consulte [Guía de Implementación de Seguridad: Defensa en Profundidad](../../guides/developer/security-implementation-guide.md).

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
