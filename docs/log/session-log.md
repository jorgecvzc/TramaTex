# Bitácora de Sesiones de Desarrollo

---
# SESIONES ABIERTAS
---

- **Unificación UI/UX — Fase 1 (EXTENSIÓN: Refinamiento de Entidades)** | Iniciada: 2026-05-09 | Status: 🟡 **EN CURSO**
    - **Contexto**: Reabierta para validación manual y corrección de tipos de documentos.
    - **Objetivo**: Añadir tipos de documentos faltantes y filtrar opciones según el tipo de entidad.
    - **Tareas**:
        - [ ] Añadir 'NIE' y 'DNI' (o equivalente) y 'Tarjeta de residente'.
        - [ ] Implementar filtrado dinámico en `PartyForm.vue`.
        - [ ] Verificación manual por el usuario.

---
# REGISTRO DE SESIONES CERRADAS
---

- **Unificación UI/UX — Fase 1: Fundamentos y Estilos Base** | Iniciada: 2026-05-01 | Finalizada: 2026-05-08 | ✅ **Fase 1 COMPLETADA**. 
    - Implementación de sistema de diseño industrial (Lucide, Skeletons, Colores).
    - Navegación Keyboard-First en listados y tablas de líneas.
    - **Atajos Globales**: Ctrl+K (Buscador), Alt+N (Nuevo), Ctrl+Enter (Guardar), Alt+R (Refrescar).
    - Auditoría global: Ajuste de alineación en dashboards y restauración de TPV clásico de alto rendimiento.
    - Refinamiento de Entidades: Formulario dinámico corregido y pestañas restauradas.
    - Despliegue exitoso en `pcele` (Staging).

- **Estudio Integral y Planificación Post-MVP (Sprint 18)** | Iniciada: 2026-04-27 | Finalizada: 2026-04-29 | ✅ Plan de ejecución y metodología post-MVP definidos en `docs/post-mvp/post-mvp-execution-plan.md`. Estrategias 01-15 consolidadas y Roadmap sincronizado.

- **Corrección de Errores — Módulo Party (Sprint 18)** | Iniciada: 2026-04-26 | Finalizada: 2026-04-27 | Status: ✅ COMPLETADO

- **Estudio y Documentación UI/UX Post-MVP (Sprint 18)** | Iniciada: 2026-04-26 | Finalizada: 2026-04-26 | ✅ Consolidada toda la estrategia en el **Plan Maestro de Unificación UI/UX** (`docs/post-mvp/01-ui-ux-unification-master-plan.md`). Incluye navegación por teclado, iconografía Lucide, alineación de dashboards y 7 nuevas mejoras de ergonomía industrial. Creada guía de ayuda al usuario y actualizado el roadmap post-MVP.

- **Estabilización de CI/CD y Lógica de Party (Sprint 18)** | Iniciada: 2026-04-24 | Finalizada: 2026-04-25 | ✅ CI backend completamente verde. Fixes: `type:uuid` en modelos sales, enum types explícitos, tabla stub `parties`, FSM domain sales, `NewInvoice` Draft status, cleanup party test_helpers. Deploy a producción exitoso (PR #19, commit `07017b8`). Descuento 0% validado funcionalmente en producción.
