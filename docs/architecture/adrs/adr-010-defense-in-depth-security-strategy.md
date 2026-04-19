# 🏛️ ADR-010: Estrategia de Seguridad (Defensa en Profundidad)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 26-01-2026 |
| **Autores** | Jorge Cortés Villalba, GitHub Copilot |

---

## 🎯 Contexto
TramaTex gestiona datos críticos (precios, márgenes, pedidos). Un compromiso de seguridad resultaría en una pérdida de confianza y riesgos legales (RGPD). Se requiere una arquitectura robusta alineada con OWASP y los principios de Zero Trust.

---

## 🔍 Alternativas Consideradas
1. **Seguridad Reactiva:** Rápida pero con alto riesgo de brechas y costes de remediación elevados.
2. **Seguridad Perimetral:** Vulnerable si se compromete el acceso externo.
3. **Defensa en Profundidad (Decisión Adoptada):** Resiliencia mediante capas redundantes de control.

---

## ✅ Decisión Adoptada
Se adopta la **Defensa en Profundidad (Defense in Depth)** junto con **Security by Design** y **Security by Default**.

### Principios Fundamentales:
*   **Capas de Seguridad:** Protección en aplicación, datos y red de forma independiente.
*   **Seguridad por Defecto:** Configuraciones seguras sin necesidad de intervención manual.
*   **Mínimo Privilegio:** Acceso estrictamente limitado a lo necesario para cada rol.
*   **Zero Trust:** Validación continua de todas las entradas, sin confiar ciegamente en el origen.

---

## 📈 Consecuencias
### Positivas
*   Alta resiliencia ante ataques dirigidos.
*   Facilidad para el cumplimiento normativo (RGPD).
*   Trazabilidad completa de acciones críticas.

### Negativas
*   Incremento en la complejidad de mantenimiento.
*   Sobrecarga ligera en el rendimiento (~5-10% latencia) y tiempo de desarrollo (~20-30%).

---
[Volver al Índice de ADRs](./README.md)
