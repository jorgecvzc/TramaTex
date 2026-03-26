# ADR-010 – Estrategia de Seguridad: Defensa en Profundidad y Security by Default

**Fecha:** 2026-01-26  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, GitHub Copilot  

---

## 1. Contexto

TramaTex gestiona datos críticos de negocio (precios, márgenes, pedidos). Un compromiso de seguridad resultaría en pérdida de confianza, daño reputacional y riesgos legales (GDPR). Se requiere establecer una arquitectura de seguridad robusta alineada con OWASP y principios de Zero Trust.

---

## 2. Alternativas Consideradas

**Alternativa A – Seguridad Reactiva (Minimal Viable Security)**
- Ventajas: Desarrollo rápido.
- Desventajas: Alta probabilidad de brechas, costo elevado de remediación post-facto.

**Alternativa B – Seguridad Perimetral Única (Network-Only)**
- Desventajas: Si el perímetro se compromete, acceso total; no protege contra amenazas internas.

**Alternativa C – Defensa en Profundidad + Security by Default (Adoptada)**
- Ventajas: Resiliencia ante fallos de una capa; prevención proactiva.
- Desventajas: Mayor complejidad arquitectónica (~20-30% overhead).

---

## 3. Criterios de Decisión

- Protección de datos críticos mediante redundancia de controles.
- Cumplimiento con mejores prácticas (OWASP, MITRE ATT&CK).
- Escalabilidad y auditabilidad de cada capa.
- Principio de menor privilegio y Zero Trust.

---

## 4. Decisión Adoptada

Se adopta **Defensa en Profundidad (Defense in Depth)** combinada con **Security by Design** y **Security by Default**.

### Principios Fundamentales:
1. **Defensa en Profundidad:** Seguridad en capas (aplicación, datos, red).
2. **Security by Default:** Configuraciones seguras sin intervención.
3. **Principio de Menor Privilegio:** Acceso mínimo necesario.
4. **Zero Trust:** Validar todo, no confiar en ninguna entrada.

---

## 5. Consecuencias

### Positivas
- Alta resiliencia y facilidad de cumplimiento (GDPR).
- Trazabilidad completa y reducción de superficie de ataque.

### Negativas
- Mayor complejidad de mantenimiento.
- Overhead de performance (~5-10% latencia) y desarrollo (~20-30%).

---

## 6. Alcance

Aplica a Backend (Go), Frontend (Vue.js), Infraestructura (Docker, PostgreSQL) y Datos en tránsito/reposo. No aplica a seguridad física de centros de datos o protección DDoS a nivel ISP.

---

## 7. Integración con otros ADRs

- **ADR-002:** Clean Architecture permite aislar concerns de seguridad.
- **ADR-003:** Modular monolith facilita la aplicación consistente de políticas.
- **ADR-009:** Define la ubicación de los componentes de seguridad.

---

## 8. Notas Adicionales / Consideraciones Especiales

### Arquitectura de Seguridad por Capas
Para una descripción detallada, consulte [Guía de Implementación de Seguridad: Defensa en Profundidad](../../guides/developer/security-implementation-guide.md).

### Checklist de Cumplimiento
Consulte el [Checklist de Cumplimiento de Seguridad](../../guides/developer/security-compliance-checklist.md).

### Notas de Implementación
Incluye fases de madurez y scripts de validación. Ver [Guía de Implementación](../../guides/developer/security-implementation-guide.md).

---

## 9. Referencias

- OWASP Top 10 2021
- MITRE ATT&CK Framework
- NIST Cybersecurity Framework
- Go Security Best Practices
