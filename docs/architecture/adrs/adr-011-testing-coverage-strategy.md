# ADR-011 – Testing & Coverage Strategy (Estrategia de Testing y Cobertura)

**Fecha:** 2026-02-01 (Reconstruido)  
**Estado:** Aceptado  
**Autores:** Equipo de Arquitectura de TramaTex  

---

## 1. Contexto

La calidad y fiabilidad son pilares fundamentales para TramaTex. Una estrategia de testing bien definida es crucial para validar requerimientos, prevenir regresiones y facilitar el mantenimiento en un entorno de desarrollo iterativo.

---

## 2. Alternativas Consideradas

**Alternativa A – Enfoque en Tests E2E / Manuales**  
- Ventajas: Valida flujos reales de usuario.  
- Desventajas: Lentos, frágiles ante cambios de UI, difícil diagnóstico de errores.

**Alternativa B – Pirámide de Testing (Adoptada)**  
- Ventajas: Base sólida de tests unitarios rápidos y económicos.  
- Desventajas: Requiere inversión inicial en infraestructura de testing y mocks.

---

## 3. Criterios de Decisión

- **Velocidad de ejecución:** Feedback rápido para el desarrollador.
- **Precisión:** Identificación clara del punto de fallo.
- **Rigor Asimétrico:** Mayor cobertura en módulos críticos (Pricing).
- **Mantenibilidad:** Tests desacoplados de la implementación técnica.

---

## 4. Decisión Adoptada

Se adopta una **estrategia de testing multi-capa** alineada con la Pirámide de Testing (70% Unit, 25% Integración, 5% E2E) y un enfoque iterativo de "construir y refinar".

### Objetivos de Cobertura (MVP)

| Módulo | Cobertura Mínima | Justificación |
|---|---|---|
| **Pricing** | **≥ 85%** | Crítico para facturación. |
| **Party / Sales** | ≥ 75% | Flujo comercial base. |
| **Product** | ≥ 50% * | Ver notas sobre integración. |
| **General** | ≥ 75% | Mínimo aceptable. |

*Nota: Product Application compensa cobertura unitaria con tests de integración robustos.*

---

## 5. Consecuencias

### Positivas
- Aumento de la confianza en la corrección del código y reducción de bugs.
- Diseño de software más modular (diseñado para ser testeable).
- Facilitación de refactorizaciones.

### Negativas
- Incremento del tiempo de desarrollo inicial.
- Esfuerzo continuo de mantenimiento de tests y mocks.

---

## 6. Alcance

Aplica a todo el ciclo de vida de desarrollo de los módulos Backend (Go) y Frontend (Vue). Define umbrales de cobertura obligatorios para el paso a producción del MVP.

---

## 7. Integración con otros ADRs

- **ADR-002 (Clean Architecture):** Los tests se organizan por capas (Domain, Application, Infrastructure).
- **ADR-010 (Seguridad):** Inclusión de tests de autorización y políticas.

---

## 8. Notas Adicionales / Consideraciones Especiales

### Tipos de Tests por Capa
- **Dominio:** Unitarios puros (sin mocks).
- **Aplicación:** Integración con mocks de servicios/repositorios.
- **Infraestructura:** Integración con base de datos real (Testcontainers).

---

## 9. Referencias

- Estándares de código (`agents/project/context/code-standards.yaml`)
- Guía de arquitectura de testing de TramaTex.
