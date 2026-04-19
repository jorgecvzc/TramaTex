# 🏛️ ADR-002: Adopción de Clean Architecture y DDD con Rigor Asimétrico

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 06-01-2026 |
| **Autores** | Jorge Cortés Villalba, ChatGPT |

---

## 🎯 Contexto
TramaTex es un ERP/MES para microempresas textiles con una gestión compleja de variantes de producto y una dependencia crítica de una tarificación precisa. El riesgo principal es la degradación del modelo de dominio por acoplamientos técnicos o dispersión de la lógica de negocio.

---

## 🔍 Alternativas Consideradas
1. **Sin Arquitectura Formal:** Rápido al inicio, pero con alto riesgo de deuda técnica y difícil escalabilidad.
2. **Rigor Uniforme:** Dominio muy protegido, pero con un esfuerzo inicial elevado y lentitud en áreas no críticas (CRUD simples).
3. **Rigor Asimétrico (Decisión Adoptada):** Protege el núcleo crítico (precios, ventas) mientras permite agilidad en áreas menos complejas.

---

## ✅ Decisión Adoptada
Se adopta **Domain-Driven Design (DDD)** junto con **Clean Architecture** aplicando un **rigor asimétrico**.

### 1. Capa de Dominio (Rigor Máximo)
*   Contiene: Entidades, Value Objects y Servicios de Dominio.
*   **Regla de Oro:** No depende de frameworks, ORMs ni infraestructura. Es testeable en aislamiento total.
*   Activo estratégico principal: Motor de Tarificación.

### 2. Capa de Aplicación
*   Orquesta los casos de uso y flujos de negocio.
*   Permite abstracciones menos estrictas para CRUDs simples, evitando boilerplate innecesario donde no aporta valor real.

### 3. Capa de Infraestructura
*   Implementa los detalles técnicos: persistencia (GORM), framework web (Gin), adaptadores externos y despliegue.
*   **Prohibición:** No debe contener lógica de negocio.

---

## 📈 Consecuencias
### Positivas
*   Dominio estable, expresivo y protegido.
*   Alta mantenibilidad a largo plazo.
*   Base sólida para una transición futura a microservicios.

### Negativas
*   Mayor coste inicial de diseño estructural.
*   Curva de aprendizaje más elevada para nuevos desarrolladores.

---
[Volver al Índice de ADRs](./README.md)
