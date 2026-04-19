# 🏛️ ADR-011: Estrategia de Testing y Cobertura

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 01-02-2026 |
| **Autores** | Equipo de Arquitectura de TramaTex |

---

## 🎯 Contexto
La fiabilidad es crítica para TramaTex. Se requiere una estrategia que valide los requisitos de negocio, prevenga regresiones y facilite el mantenimiento en un ciclo de desarrollo iterativo.

---

## 🔍 Alternativas Consideradas
1. **Foco en E2E / Manual:** Valida flujos reales pero es lento y frágil ante cambios de interfaz.
2. **Pirámide de Testing (Decisión Adoptada):** Base sólida de tests unitarios rápidos combinada con tests de integración y E2E.

---

## ✅ Decisión Adoptada
Se adopta una **Pirámide de Testing (70% Unit, 25% Integración, 5% E2E)** con un enfoque de rigor asimétrico.

### Objetivos de Cobertura (MVP):
| Módulo | Cobertura Mínima | Justificación |
| :--- | :--- | :--- |
| **Pricing** | **≥ 85%** | Crítico para la rentabilidad y facturación. |
| **Party / Sales** | ≥ 75% | Flujo comercial base del sistema. |
| **General** | ≥ 75% | Mínimo aceptable para la salud del proyecto. |

---

## 📈 Consecuencias
### Positivas
*   Máxima confianza en la integridad del código.
*   Diseño modular forzado por la testabilidad.
*   Facilitación de refactorizaciones seguras.

### Negativas
*   Inversión inicial significativa en infraestructura de pruebas.
*   Esfuerzo continuo de mantenimiento de mocks y suites de test.

---
[Volver al Índice de ADRs](./README.md)
