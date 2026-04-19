# 🏛️ ADR-014: Arquitectura del Módulo de IAM

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 06-02-2026 |
| **Autores** | Gemini CLI |

---

## 🎯 Contexto
El módulo de IAM (Identity and Access Management) es el componente transversal responsable de la identidad, autenticación y autorización (RBAC). Dado que el sistema maneja datos críticos, la seguridad es de máxima prioridad.

---

## 🔍 Alternativas Consideradas
1. **Basada en Sesiones:** Tradicional, pero difícil de escalar y con mayor riesgo de CSRF en APIs.
2. **Proveedor Externo (OAuth/Auth0):** Alta seguridad pero introduce dependencia de internet y complejidad de configuración para un MVP Local-First.
3. **Basada en JWT (Decisión Adoptada):** Estándar de la industria, *stateless* y flexible para diferentes clientes (Web/Móvil).

---

## ✅ Decisión Adoptada
Se adopta la **Autenticación y Autorización basada en JWT**.

### Claves del Diseño:
*   **Tokens Duales:** Uso de `access_token` (corta duración) y `refresh_token` (larga duración).
*   **RBAC (Role Based Access Control):** Roles predefinidos (`admin`, `commercial`, `designer`, `workshop`) que controlan el acceso a la API y la UI.
*   **Seguridad de Credenciales:** Hashing de contraseñas mediante `bcrypt`.
*   **Arquitectura Limpia:** Exclusión de campos de auditoría (`CreatedAt`, etc.) de las entidades de dominio, delegándolos a la infraestructura.

---

## 📈 Consecuencias
### Positivas
*   Sistema *stateless* que facilita la escalabilidad y distribución de la API.
*   Autenticación rápida sin consultas constantes a la base de datos para validar el token.
*   Control de acceso granular y predecible.

### Negativas
*   La revocación de tokens antes de su expiración requiere una estrategia de *blacklist* (planteada para Post-MVP).

---
[Volver al Índice de ADRs](./README.md)
