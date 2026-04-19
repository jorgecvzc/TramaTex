# 🏛️ ADR-003: Tipo y Distribución de la Aplicación (Monolito Modular)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 07-01-2026 |
| **Autores** | Jorge Cortés Villalba, ChatGPT |

---

## 🎯 Contexto
El sistema debe operar en infraestructuras **Local-First** con hardware limitado y ser mantenido por terceros. Se requiere una solución que maximice la simplicidad operativa actual sin bloquear el escalado funcional o una posible transición a SaaS en el futuro.

---

## 🔍 Alternativas Consideradas
1. **Microservicios:** Excesiva complejidad operativa para un equipo reducido y hardware local.
2. **Monolito Tradicional:** Riesgo de acoplamiento ("monolito de barro") y difícil evolución.
3. **Monolito Modular (Decisión Adoptada):** Simplicidad de un solo proceso con la modularidad lógica de servicios independientes.

---

## ✅ Decisión Adoptada
Se adopta un **Monolito Modular Local-First** bajo los siguientes principios:

### 1. Un Solo Proceso, Múltiples Dominios
*   **Backend:** Un único binario Go (`tramatex-api`).
*   **Base de Datos:** Una única instancia de PostgreSQL con separación lógica por esquemas o convenciones.

### 2. Límites Claros (Bounded Contexts)
*   Cada módulo define su propio modelo, reglas y casos de uso.
*   La comunicación entre módulos es explícita mediante interfaces y contratos.
*   **Extraíble por diseño:** Cada módulo está preparado para convertirse en un microservicio independiente si fuera necesario.

### 3. Frontend SPA
*   Una única aplicación Vue.js 3 que actúa como cliente de la API y adaptador de presentación.

---

## 📈 Consecuencias
### Positivas
*   Coste operativo y de despliegue mínimo.
*   Dominio protegido frente a la degradación.
*   Compatible con hardware modesto (Local-First).

### Negativas
*   Requiere una disciplina estricta de gobernanza para evitar acoplamientos entre módulos.
*   Menor aislamiento en tiempo de ejecución comparado con microservicios físicos.

---
[Volver al Índice de ADRs](./README.md)
