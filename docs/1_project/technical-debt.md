# 📋 Registro de Deuda Técnica - TramaTex

**Última actualización:** 2026-01-27

---

## 🎯 Propósito

Este documento registra y prioriza la deuda técnica del proyecto TramaTex, clasificando mejoras pendientes y estableciendo un plan de resolución. La gestión proactiva de la deuda técnica es un pilar de nuestra estrategia de calidad.

---

## 📊 Clasificación de Deuda

### Categorías

1.  **Seguridad** 🔴 - Vulnerabilidades o riesgos de seguridad.
2.  **Rendimiento** 🟠 - Optimizaciones de velocidad o uso de recursos.
3.  **Mantenibilidad** 🟡 - Código difícil de mantener, entender o refactorizar (alto acoplamiento, baja cohesión).
4.  **Testing** 🔵 - Gaps en la cobertura de tests, tests frágiles o lentos.
5.  **Documentación** 🟣 - Documentación desactualizada, incompleta o confusa.

### Prioridades

- **P0 - Crítica**: Debe resolverse inmediatamente. Bloquea el desarrollo o compromete la estabilidad/seguridad del sistema.
- **P1 - Alta**: Debe resolverse en el siguiente sprint.
- **P2 - Media**: Planificar para los próximos 1-2 meses.
- **P3 - Baja**: Backlog. Resolver cuando haya oportunidad.

---

## 🔴 Deuda de Seguridad

### [P1] Rate Limiting en Endpoints de Login

-   **Categoría:** Seguridad
-   **Origen:** Auditoría OWASP (A04 - Insecure Design)
-   **Descripción:** El endpoint `/api/iam/login` no tiene límite de peticiones, lo que permite ataques de fuerza bruta.
-   **Impacto:** **Alto**. Riesgo de compromiso de cuentas de usuario.
-   **Propuesta:** Implementar un middleware de rate limiting (ej. `gin-limiter`) con un límite de 5-10 intentos por minuto por IP.
-   **Estado:** Planificado (Sprint post-MVP).

### [P2] Validación Avanzada de Contraseñas

-   **Categoría:** Seguridad
-   **Origen:** Auditoría OWASP (A07 - Auth Failures)
-   **Descripción:** La política de contraseñas actual solo valida la longitud mínima. No valida complejidad (mayúsculas, números, símbolos) ni la compara con listas de contraseñas comunes.
-   **Impacto:** Medio. Los usuarios pueden elegir contraseñas débiles y fáciles de adivinar.
-   **Propuesta:** Añadir validación de complejidad en el `Password` Value Object y, opcionalmente, una integración con un servicio como 'Have I Been Pwned'.
-   **Estado:** Planificado (Post-MVP).

### [P2] Sistema de Recuperación de Contraseñas

-   **Categoría:** Seguridad / Funcionalidad
-   **Origen:** Auditoría OWASP (A07 - Auth Failures)
-   **Descripción:** No existe un flujo para que los usuarios recuperen su contraseña si la olvidan.
-   **Impacto:** Medio. Genera una dependencia operativa del administrador del sistema.
-   **Propuesta:** Implementar un flujo de recuperación seguro basado en tokens enviados por correo electrónico con una validez corta.
-   **Estado:** Planificado (Post-MVP).

---

## 🔵 Deuda de Testing

### [P2] Gaps de Cobertura en el Módulo Party

-   **Categoría:** Testing
-   **Origen:** Revisión de Sprint 05.
-   **Descripción:** Aunque la cobertura del módulo Party es alta, existen algunos branches en los manejadores de errores de la capa de persistencia que no están cubiertos por los tests de integración.
-   **Impacto:** Bajo. Los "happy paths" y errores principales están cubiertos.
-   **Propuesta:** Añadir tests de integración que simulen fallos de la base de datos (ej. violación de constraints de unicidad) para validar el manejo de esos errores.
-   **Estado:** Planificado.

---

## 🟣 Deuda de Documentación

### [P1] Documentación de API con OpenAPI/Swagger

-   **Categoría:** Documentación / Funcionalidad
-   **Origen:** Necesidad de facilitar la integración y el testing manual.
-   **Descripción:** La API REST carece de una especificación formal y autogenerada.
-   **Impacto:** Medio. Dificulta el trabajo del equipo de frontend y la realización de pruebas.
-   **Propuesta:** Integrar `swaggo` en el proyecto Go para generar una especificación OpenAPI 3 a partir de los comentarios del código y servir una UI de Swagger.
-   **Estado:** Planificado (Sprints 06-07).

---

## Proceso de Gestión de Deuda Técnica

1.  **Identificación:** Cualquier miembro del equipo puede añadir un nuevo ítem a este documento durante una revisión de código, una sesión de desarrollo o una retrospectiva.
2.  **Clasificación:** Se debe asignar una categoría y una prioridad inicial.
3.  **Revisión:** La deuda técnica se revisa al inicio de la planificación de cada nuevo sprint para decidir qué ítems de prioridad P1 se abordarán.
4.  **Resolución:** Cuando se resuelve un ítem, el commit o PR debe hacer referencia a él. El ítem se mueve a una sección de "Deuda Resuelta" al final de este documento.
