# Bitácora de Sesiones de Desarrollo

---
# SESIONES ABIERTAS
---

## Corrección de Errores — Módulo Party (Sprint 18)

- **Session ID:** `party-module-fixes-2026-04-26`
- **Status:** En Progreso
- **Sprint:** Sprint 18
- **Started:** 2026-04-26
- **Branch:** `fix/party-module-fixes` (desde `develop`)

**Contexto:** Sesión dedicada a identificar y corregir errores en el módulo Party (backend y/o frontend). Los tests unitarios del backend pasan correctamente. Se investigarán bugs funcionales, de integración o de UI que se detecten durante la sesión.

**Próximos Pasos:**
- [x] Identificar y documentar los errores concretos en el módulo Party
- [x] **BUG FIX**: Error "discount can only be assigned to customers" al cambiar de Proveedor a Cliente/Ambos — corregido en `apps/frontend/src/services/partyApi.ts` (`updateParty`): roles sincronizados **antes** del PUT en lugar de después
- [ ] Ejecutar tests y verificar que no hay regresiones
- [ ] Merge a `develop` una vez validado

**Archivos de Contexto:**
- `apps/tramatex-api/internal/party/`
- `apps/frontend/src/modules/party/`
- `docs/modules/party/`

# REGISTRO DE SESIONES CERRADAS
---

- **Estudio y Documentación UI/UX Post-MVP (Sprint 18)** | Iniciada: 2026-04-26 | Finalizada: 2026-04-26 | ✅ Consolidada toda la estrategia en el **Plan Maestro de Unificación UI/UX** (`docs/post-mvp/01-ui-ux-unification-master-plan.md`). Incluye navegación por teclado, iconografía Lucide, alineación de dashboards y 7 nuevas mejoras de ergonomía industrial. Creada guía de ayuda al usuario y actualizado el roadmap post-MVP.

- **Estabilización de CI/CD y Lógica de Party (Sprint 18)** | Iniciada: 2026-04-24 | Finalizada: 2026-04-25 | ✅ CI backend completamente verde. Fixes: `type:uuid` en modelos sales, enum types explícitos, tabla stub `parties`, FSM domain sales, `NewInvoice` Draft status, cleanup party test_helpers. Deploy a producción exitoso (PR #19, commit `07017b8`). Descuento 0% validado funcionalmente en producción.

- **Deep Documentation Study & Gap Analysis (Sprint 18)** | Iniciada: 2026-04-19 | Finalizada: 2026-04-19 | ✅ Barrido exhaustivo de documentación completado: ramas y hojas alineadas, referencias corregidas y especificaciones sincronizadas.
- **Documentation Review & Alignment (Sprint 18)** | Iniciada: 2026-04-17 | Finalizada: 2026-04-19 | ✅ Revisión integral de documentación cerrada tras consolidar ADRs, guías, módulos y presentación corporativa.
- **Continuidad IAM + Sales — Estado Actual para Retomar (Sprint 18)** | Iniciada: 2026-04-18 | Finalizada: 2026-04-19 | ✅ Sesión de continuidad cerrada tras verificar integración en `develop`, limpieza de pendientes locales y consolidación de fixes IAM/Sales.
- **POS Polish & MES Terminal Excellence (Sprint 18)** | Iniciada: 2026-04-18 | Finalizada: 2026-04-18 | ✅ TPV industrial optimizado (sin scroll, atajos en cabecera), Terminal MES tabular con recuperación de cliente y Dashboard con KPIs reales.
- **Fix Production Launch Errors (Sales/IAM) (Sprint 18)** | Iniciada: 2026-04-11 | Finalizada: 2026-04-18 | ✅ Bug de token corregido y flujo de lanzamiento a producción validado.
- **Review Sprint 17 Implementation Status (Sprint 17 / 18)** | Iniciada: 2026-04-12 | Finalizada: 2026-04-18 | ✅ Tareas de consolidación de Pricing, Money decimal y ACL clients verificadas y completadas.
- **IAM Role Alignment — Fix DB↔Domain Role Mismatch (Sprint 18)** | Iniciada: 2026-04-10 | Finalizada: 2026-04-11 | ✅ Roles alineados y bug `updateUser` corregido.
- **Pricing Consolidation — Paridad Funcional + Eliminación Legacy (Sprint 17)** | Iniciada: 2026-04-08 | Finalizada: 2026-04-11 | ✅ Motor de precios optimizado y UI saneada.
- **POST-REFACTOR STABILIZATION & FINAL POLISH (Sprint 16)** | Iniciada: 2026-04-07 | Finalizada: 2026-04-08 | Mergeada a producción
