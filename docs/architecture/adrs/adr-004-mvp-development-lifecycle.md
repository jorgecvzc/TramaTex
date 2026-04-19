# 🏛️ ADR-004: Ciclo de Vida de Desarrollo e Implementación hasta MVP

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 07-01-2026 |
| **Autores** | Jorge Cortés Villalba, ChatGPT |

---

## 🎯 Contexto
TramaTex debe alcanzar un **MVP plenamente operativo** en un entorno real con un equipo reducido y recursos de hardware limitados. El MVP no es un prototipo, sino un sistema productivo que requiere una calidad técnica elevada (Clean Architecture, TDD) para asegurar la integridad de los procesos financieros y de taller.

---

## 🔍 Alternativas Consideradas
1. **Ciclo Completo (MVP + Post-MVP):** Rechazada por diluir el foco y aumentar el riesgo de no entregar el núcleo sólido a tiempo.
2. **Prototipado Rápido (Sin Rigor):** Rechazada por el alto riesgo de errores en tarificación y una deuda técnica inasumible.
3. **Ciclo Incremental Orientado a MVP (Decisión Adoptada):** Foco en incrementos funcionales completos, validación temprana en entorno real y rigor técnico desde la base.

---

## ✅ Decisión Adoptada
Se adopta un **Ciclo de Vida Incremental** dividido en fases:

### Fase 0: Cimentación Técnica
*   Preparación de infraestructura: Docker, CI/CD, pipeline de tests y linters.
*   Esqueleto de Clean Architecture y autenticación base.

### Fase 1: Implementación del Core (MVP)
*   Desarrollo de módulos críticos: Ventas, Tarificación, Terceros (Party) y flujo esencial de taller (MES).
*   Interfaces gráficas funcionales para ventas y operarios.

### Estrategia de Iteración
*   Iteraciones cortas con cada incremento funcional y testeado.
*   Enfoque ágil pragmático orientado a la entrega de valor real.

---

## 📈 Consecuencias
### Positivas
*   Foco absoluto en el valor de negocio inmediato.
*   Reducción del riesgo de sobreingeniería.
*   Obtención de un sistema estable, productivo y mantenible por terceros.

### Negativas
*   Requiere una disciplina estricta en la planificación del alcance.
*   Funcionalidades avanzadas quedan postergadas para fases Post-MVP.

---
[Volver al Índice de ADRs](./README.md)
