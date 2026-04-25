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
- **Próximos Pasos:**
  - [x] Implementar validación de rol para descuento comercial en Backend.
  - [x] Cambiar `tx.Save` por `tx.Updates(map)` en repositorio GORM para forzar guardado de valores 0.
  - [x] Actualizar `demo-reset.yml` para usar `docker compose down -v` (borrado de volúmenes).
  - [x] Añadir verificación de carga de seed data (usuario admin) en el workflow de reseteo.
  - [ ] Investigar por qué el driver de Go ignora `PGUSER` y busca `root` en los tests de integración de la CI.
  - [ ] Validar funcionalmente el guardado del descuento 0% en el despliegue tras el reseteo.
- **Archivos de Contexto:**
  - `apps/frontend/src/components/party/PartyForm.vue`
  - `apps/tramatex-api/internal/party/application/party_commands.go`
  - `apps/tramatex-api/internal/party/persistence/gorm_party.go`
  - `.github/workflows/backend.yml`
  - `.github/workflows/demo-reset.yml`

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

