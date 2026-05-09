# Bitácora de Sesiones de Desarrollo

---
# SESIONES ABIERTAS
---

- **Unificación UI/UX — Fase 1 (EXTENSIÓN: Refinamiento de Entidades)** | Iniciada: 2026-05-09 | Status: 🟡 **EN CURSO**
    - **Contexto**: Reabierta para corrección de tipos de documentos, desglose de nombres y UX de listados.
    - **Cambios Realizados**:
        - [x] Añadidos 'DNI', 'Pasaporte' y 'Tarjeta de residente' a `TaxIdType`.
        - [x] Implementado filtrado dinámico de documentos según tipo de entidad (Persona/Empresa).
        - [x] Desglosado campo 'Nombre' en 'Nombre' y 'Apellidos' para Personas Físicas.
        - [x] Corregido bug de visualización (interfaz en blanco) en `Detail.vue` por importación faltante.
        - [x] Asegurada generación de UUID en `partyApi.ts` para evitar error "party ID cannot be empty".
        - [x] Mejorada la UX de catálogos: eliminación de selección persistente al hacer click.
        - [x] Script de despliegue remoto ahora auto-repara permisos de `.git` y `.docker_config`.

    - **📋 GUÍA DE COMPROBACIÓN MANUAL (Protocolo de Validación)**:
        1.  **Alta de Empresa**: Crear una entidad "ORGANIZATION". Verificar que solo sale un campo de nombre y que los documentos filtran (CIF, VAT).
        2.  **Alta de Persona**: Crear una entidad "PERSON". Verificar que aparecen dos campos (Nombre y Apellidos). Verificar que los documentos incluyen DNI, NIE y Tarjeta de Residente.
        3.  **Edición y Pestañas**: Editar una empresa y verificar que la pestaña "Contactos" es visible y funcional.
        4.  **Validación de Guardado**: Intentar guardar sin nombre y comprobar que el sistema lo impide. Guardar con éxito y verificar que el ID se genera correctamente (sin errores de "empty ID").
        5.  **Comportamiento de Listados**: En Marcas o Entidades, hacer click en una fila. Al volver atrás o cancelar, la fila NO debe estar resaltada. Usar flechas del teclado y verificar que ahí SÍ hay resalto visual.
        6.  **Despliegue Remoto**: Ejecutar `./scripts/rebuild-staging-remote.ps1` y confirmar que no hay errores de permisos ("Permission denied").

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
