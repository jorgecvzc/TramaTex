# 🏛️ ADR-021: Estrategia de Control de Versiones y Ramas

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 22-02-2026 |
| **Autores** | Equipo de Desarrollo TramaTex |

---

## 🎯 Contexto
Con el MVP listo para producción (v1.0.0), se requiere una estrategia de branching que permita el desarrollo continuo de funcionalidades Post-MVP, un versionado semántico consistente y despliegues seguros.

---

## ✅ Decisión Adoptada
Se adopta un **GitFlow Simplificado** y el estándar **SemVer 2.0**.

### 1. Modelo de Ramas:
*   **`master`**: Código estable y desplegable. Cada merge es una versión productiva. Protegida.
*   **`develop`**: Rama de integración activa. Base para todas las nuevas tareas. Protegida.
*   **`feature/*`**: Nuevas funcionalidades (nacen de `develop`).
*   **`bugfix/*`**: Correcciones en desarrollo (nacen de `develop`).
*   **`hotfix/*`**: Errores críticos en producción (nacen de `master`).

### 2. Versionado Semántico:
Formato `MAJOR.MINOR.PATCH` (Ej: 1.2.0).
*   **MAJOR**: Cambios arquitectónicos o de API incompatibles.
*   **MINOR**: Nuevas funcionalidades compatibles.
*   **PATCH**: Corrección de errores y optimizaciones.

### 3. Seguridad de Datos:
*   **Prohibición Estricta:** Ningún archivo `.env` real debe ser trackeado por Git.
*   **Uso de Plantillas:** Obligatoriedad de mantener archivos `.env.example` actualizados para cada entorno.

---

## 📈 Consecuencias
### Positivas
*   Estabilidad garantizada en la rama de producción.
*   Desarrollo paralelo de múltiples funcionalidades sin bloqueos.
*   Historial de cambios claro y profesional mediante etiquetas (tags).

### Negativas
*   Requiere mayor disciplina en el seguimiento del flujo de trabajo por parte de los desarrolladores.
*   Aumenta ligeramente la burocracia en tareas pequeñas debido a la necesidad de PRs.

---
[Volver al Índice de ADRs](./README.md)
