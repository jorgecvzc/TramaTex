# 🏛️ ADR-018: Arquitectura del Módulo MES (Producción)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 07-02-2026 |
| **Autores** | Gemini CLI (Architect) |

---

## 🎯 Contexto
El módulo MES (Manufacturing Execution System) es crítico para la planificación y control de calidad en taller. Aunque nace en un monolito, el diseño debe permitir su futura extracción como microservicio independiente sin re-arquitectura masiva.

---

## 🔍 Alternativas Consideradas
1. **Esquema Compartido:** Rápido de implementar, pero genera un alto acoplamiento que impediría separar el servicio en el futuro.
2. **Esquema Independiente y Comunicación por Interfaces (Decisión Adoptada):** Máximo desacoplamiento y protección de la integridad del taller.

---

## ✅ Decisión Adoptada
El módulo MES se construye bajo el principio de **"Extraíble por Diseño"**:

### Reglas de Arquitectura:
*   **Autonomía de Datos:** MES gestiona sus propias tablas. No se permiten `JOINs` SQL entre tablas de MES y otros módulos (Ventas o Productos).
*   **Comunicación por Contratos:** Toda interacción con el taller se realiza mediante interfaces de aplicación (APIs internas) o eventos de dominio.
*   **Sincronización:** Sales consulta el estado de las órdenes de trabajo a través de adaptadores, manteniendo a MES como la autoridad única sobre el taller.

---

## 📈 Consecuencias
### Positivas
*   Preparación nativa para la evolución a microservicios.
*   Alta cohesión: la lógica de producción está aislada y protegida.
*   Facilidad para realizar cambios técnicos internos en el taller sin afectar al ERP.

### Negativas
*   Requiere un diseño más riguroso de las interfaces de comunicación.
*   Sobrecarga ligera en el desarrollo inicial debido a la necesidad de adaptadores.

---
[Volver al Índice de ADRs](./README.md)
