# 🏛️ ADR-006: Estrategia de Desarrollo Dirigida por el Dominio

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 09-01-2026 |
| **Autores** | Jorge Cortés Villalba, ChatGPT |

---

## 🎯 Contexto
El principal riesgo del proyecto no es tecnológico, sino de ejecución. Una implementación desordenada puede diluir el dominio core (tarificación y ventas) y generar retrabajo. Se necesita garantizar la trazabilidad total entre los requisitos de negocio, el código y la documentación.

---

## 🔍 Alternativas Consideradas
1. **Desarrollo Basado en Capas Técnicas:** Enfoque tradicional rápido pero con alto riesgo de mezclar lógica técnica con reglas de negocio.
2. **Desarrollo Dirigido por el Dominio (Decisión Adoptada):** Protección del dominio core, implementación incremental validable y minimización de la deuda técnica.

---

## ✅ Decisión Adoptada
Se adopta una **estrategia de desarrollo incremental dirigida por el dominio (DDD)** para el MVP.

### Justificación y Reglas:
*   **Prioridad al Núcleo:** Los módulos críticos (Tarificación, Producto, Party) se implementan y testean antes que su infraestructura.
*   **Infraestructura Just-In-Time:** Solo se introduce tecnología (persistencia, API) cuando existen casos de uso de dominio validados.
*   **Trazabilidad:** Cada decisión en el código debe responder a una regla de negocio documentada.
*   **Testing:** Foco absoluto en el desarrollo guiado por pruebas (TDD) en las capas de dominio y aplicación.

---

## 📈 Consecuencias
### Positivas
*   Dominio económico (precios y márgenes) protegido y estable.
*   Sistema altamente testeable y con mínima deuda técnica temprana.
*   Facilidad de mantenimiento y evolución futura.

### Negativas
*   Mayor esfuerzo inicial en análisis y diseño.
*   El ritmo de desarrollo de la interfaz de usuario (UI) es más lento al inicio.

---
[Volver al Índice de ADRs](./README.md)
