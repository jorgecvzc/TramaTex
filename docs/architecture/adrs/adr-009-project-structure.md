# 🏛️ ADR-009: Estructura de Carpetas y Organización del Proyecto

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 11-01-2026 |
| **Autores** | Jorge Cortés Villalba, Claude (Anthropic) |

---

## 🎯 Contexto
TramaTex requiere una estructura de repositorio único (monorepo) que refleje los principios de Clean Architecture y facilite el mantenimiento por un único desarrollador con apoyo de IA, permitiendo además una futura extracción a microservicios.

---

## 🔍 Alternativas Consideradas
1. **Estructura Plana:** Simple, pero no escala y mezcla dominios de negocio.
2. **Por Capas Técnicas:** Respeta Clean Architecture pero oculta la intención de negocio.
3. **Por Módulos de Dominio (Decisión Adoptada):** Refleja los Bounded Contexts claramente y facilita la testabilidad y escalabilidad.

---

## ✅ Decisión Adoptada
Se adopta una **Estructura por Módulos de Dominio con Clean Architecture Interna**.

### Organización del Monorepo:
*   `apps/tramatex-api/internal/`: Contiene los módulos (`iam`, `party`, `pricing`, etc.). Cada módulo tiene sus subcarpetas de `domain`, `application` e `infrastructure`.
*   `apps/frontend/`: Aplicación única en Vue.js.
*   `docs/`: El Árbol de Conocimiento del proyecto.
*   `project-scaffolding/`: Herramientas de estandarización.

### Justificación:
Este diseño permite que cada módulo sea una unidad lógica independiente, facilitando que los tests estén junto al código y preparando el sistema para una futura distribución física sin necesidad de refactorizaciones masivas.

---

## 📈 Consecuencias
### Positivas
*   Separación clara de responsabilidades y alta testabilidad.
*   Preparación nativa para la evolución hacia microservicios.
*   Navegación intuitiva basada en el lenguaje del negocio.

### Negativas
*   Mayor profundidad de carpetas (rutas de archivo más largas).
*   Exige disciplina para no romper los límites entre módulos.

---
[Volver al Índice de ADRs](./README.md)
