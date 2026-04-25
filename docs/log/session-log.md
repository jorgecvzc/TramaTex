# Bitácora de Sesiones de Desarrollo

---
# SESIONES ABIERTAS
---

## Estabilización de CI/CD y Lógica de Party
- **Session ID:** `ci-stability-party-discount-fix`
- **Status:** En Progreso
- **Sprint:** Sprint 18
- **Started:** 2026-04-24
- **Contexto:** Resolver fallos persistentes en Backend CI y Demo Reset, y corregir la persistencia del descuento en Party.
- **Logros:**
  - ✅ Corregida la persistencia de descuento 0% en `gorm_party.go` usando `Updates(map[string]interface{})`.
  - ✅ Añadido test de integración para validar el descuento 0% en `party`.
  - ✅ Silenciado el logger de GORM en tests unitarios y de integración para limpiar el ruido en la CI.
  - ✅ Implementado `AutoMigrate` en los `test_helpers` de `sales` para garantizar coherencia de esquema.
  - ✅ Añadida política de `concurrency` en GitHub Actions para evitar ejecuciones duplicadas.
  - ✅ Corregido test de nombre de tabla en `iam` (`UserModel`).
- **Pendientes:**
  - [ ] Identificar y corregir la causa raíz del fallo persistente en Backend CI (el semáforo sigue en rojo tras limpiar logs).
  - [ ] Investigar si hay otros módulos (MES, Pricing) con tests de integración que requieran actualización de esquema.
  - [ ] Validar funcionalmente el guardado del descuento 0% en el despliegue tras el reseteo.
- **Archivos de Contexto:**
  - `apps/tramatex-api/internal/party/persistence/gorm_party.go`
  - `apps/tramatex-api/internal/sales/infrastructure/persistence/test_helpers.go`
  - `.github/workflows/backend.yml`
  - `docs/guides/developer/ci-cd.md`


# REGISTRO DE SESIONES CERRADAS
---

- **Deep Documentation Study & Gap Analysis (Sprint 18)** | Iniciada: 2026-04-19 | Finalizada: 2026-04-19 | ✅ Barrido exhaustivo de documentación completado: ramas y hojas alineadas, referencias corregidas y especificaciones sincronizadas.
- **Documentation Review & Alignment (Sprint 18)** | Iniciada: 2026-04-17 | Finalizada: 2026-04-19 | ✅ Revisión integral de documentación cerrada tras consolidar ADRs, guías, módulos y presentación corporativa.
- **Continuidad IAM + Sales — Estado Actual para Retomar (Sprint 18)** | Iniciada: 2026-04-18 | Finalizada: 2026-04-19 | ✅ Sesión de continuidad cerrada tras verificar integración en `develop`, limpieza de pendientes locales y consolidación de fixes IAM/Sales.
- **POS Polish & MES Terminal Excellence (Sprint 18)** | Iniciada: 2026-04-18 | Finalizada: 2026-04-18 | ✅ TPV industrial optimizado (sin scroll, atajos en cabecera), Terminal MES tabular con recuperación de cliente y Dashboard con KPIs reales.
- **Fix Production Launch Errors (Sales/IAM) (Sprint 18)** | Iniciada: 2026-04-11 | Finalizada: 2026-04-18 | ✅ Bug de token corregido y flujo de lanzamiento a producción validado.
- **Review Sprint 17 Implementation Status (Sprint 17 / 18)** | Iniciada: 2026-04-12 | Finalizada: 2026-04-18 | ✅ Tareas de consolidación de Pricing, Money decimal y ACL clients verificadas y completadas.
- **IAM Role Alignment — Fix DB↔Domain Role Mismatch (Sprint 18)** | Iniciada: 2026-04-10 | Finalizada: 2026-04-11 | ✅ Roles alineados y bug `updateUser` corregido.
- **Pricing Consolidation — Paridad Funcional + Eliminación Legacy (Sprint 17)** | Iniciada: 2026-04-08 | Finalizada: 2026-04-11 | ✅ Motor de precios optimizado y UI saneada.
- **POST-REFACTOR STABILIZATION & FINAL POLISH (Sprint 16)** | Iniciada: 2026-04-07 | Finalizada: 2026-04-08 | Mergeada a producción

