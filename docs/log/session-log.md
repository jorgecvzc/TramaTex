# Bitácora de Sesiones de Desarrollo

---
# SESIONES ABIERTAS
---

## 📅 2026-05-13 — Estabilización Global de UI y Flujos
*   **Contexto:** Cierre de issues pendientes en Sales, estandarización de UI en Parties y refactorización final de componentes.
*   **Estado:** En curso (Rama: `fix/sales-order-detail-blank-screen-and-conversion-flow`).
*   **Análisis detallado:** [docs/log/analysis/2026-05-13-global-ui-stabilization.md](docs/log/analysis/2026-05-13-global-ui-stabilization.md)

---
# REGISTRO DE SESIONES CERRADAS
---

- **Refinamiento UI/UX y Regresiones Sales** | Iniciada: 2026-05-12 | Finalizada: 2026-05-13 | Status: ✅ **COMPLETADO**
    - **Contexto**: Resolución de fallos críticos en la visualización de pedidos y estabilización del flujo de conversión de presupuestos.
    - **Cambios Realizados**:
        - [x] **OrderDetail.vue**: Corregida pantalla en blanco (missing import `useLineNavigation`) y refactorizado a camelCase integral para consistencia con el resto de la UI.
        - [x] **QuoteDetail.vue**: Refactorizado para utilizar el componente unificado `OrderLines`, eliminando duplicidad de lógica.
        - [x] **OrderLines.vue**: Añadida columna de "Precio Tarifa" y estandarización de nombres de campos en camelCase.
        - [x] **salesApi.ts**: Mejorada la normalización automática de entidades (snake_case -> camelCase recursivo) para blindar el frontend ante desajustes del backend.
        - [x] **Backend (Go)**: Implementadas transacciones en `ConvertQuoteToOrder` y `AcceptAndConvertQuote` para asegurar la atomicidad de la creación de pedidos.
        - [x] **Dominio (Go)**: Corregido bug en `canTransitionInvoice` que permitía saltar de `DRAFT` a `PAID`, alineando el código con los tests de dominio.
    - **📋 PROTOCOLO DE VALIDACIÓN FINAL**:
        1.  **Visualización**: El detalle de pedido (`OrderDetail.vue`) ya no muestra pantalla en blanco y renderiza todas las líneas correctamente.
        2.  **Unificación**: Tanto presupuestos como pedidos utilizan el mismo componente de líneas con soporte completo de teclado.
        3.  **Conversión**: Probada la conversión de presupuestos en estado `ISSUED`, verificando que el pedido se crea y el presupuesto se cierra correctamente en una única operación atómica.
        4.  **Tests**: `go test ./apps/tramatex-api/internal/sales/...` devuelve `ok` en todos los subpaquetes.

---

- **Unificación UI/UX — Fase 1 (CIERRE: Experiencia "Sin Ratón" y Centro de Ayuda)** | Iniciada: 2026-05-10 | Finalizada: 2026-05-10 | Status: ✅ **Fase 1 COMPLETADA 100%**
    - **Contexto**: Finalización del plan maestro de UI/UX con enfoque en productividad industrial y auto-capacitación.
    - **Cambios Realizados**:
        - [x] **Iconografía Lucide**: Migración total de `Material Symbols` a componentes `Lucide` en Sales y Entidades, eliminando etiquetas técnicas en inglés.
        - [x] **Navegación Teclado**: Implementado composable `useLineNavigation` para gestión experta de tablas (Flechas +/- , Enter inteligente, Ctrl+Supr para borrar).
        - [x] **Módulo MES**: Estandarización de líneas técnicas en Presupuestos y Pedidos con la misma experiencia de teclado que las líneas de producto.
        - [x] **Sistema de Ayuda**:
            - Creado **Menú de Ayuda (`Alt+H`)** unificado en la Sidebar.
            - Implementada **Ayuda Contextual (`F1`)** con panel lateral dinámico según la ruta.
            - Creado **Centro de Capacitación (`/help`)** navegable por teclado (teclas 1-4).
        - [x] **Ergonomía**: Auto-foco en buscadores al entrar en listados y apertura de modales.
        - [x] **Atajos Globales**: `Ctrl+S` (Guardar), `Esc` (Atrás/Cerrar), `Ctrl+K` (Buscar), `Alt+1..5` (Módulos).
        - [x] **Estabilidad**: Corregidos crashes por imports faltantes (`onBeforeUnmount`) y persistencia de pestañas en edición de Entidades.
        - [x] **README**: Añadido reconocimiento a Marisol López Núñez por el Jingle oficial del proyecto.

    - **📋 PROTOCOLO DE VALIDACIÓN FINAL**:
        1.  **Sin Ratón**: Capacidad de crear un presupuesto completo con líneas MES y de producto usando solo teclado.
        2.  **Aprendizaje**: Comprobar que pulsar `?` muestra el mapa y `F1` la guía de la página actual.
        3.  **Resiliencia**: Verificar que al editar una entidad y cambiar de pestaña, los datos del formulario no se pierden.

---

- **Unificación UI/UX — Fase 1 (EXTENSIÓN: Refinamiento de Entidades)** | Iniciada: 2026-05-09 | Finalizada: 2026-05-09 | Status: ✅ **COMPLETADO**

    - **Contexto**: Corrección de tipos de documentos, desglose de nombres, UX de listados y normalización de terminales.
    - **Cambios Realizados**:
        - [x] **Entidades**: Desglosado 'Nombre' en 'Nombre' y 'Apellidos' para personas físicas. Añadidos DNI, Pasaporte y Tarjeta Residente con filtrado dinámico.
        - [x] **Terminales**: Normalizadas etiquetas en MES y TPV (ej: "BASE IMPONIBLE", "TOTAL A PAGAR", estados reales de tareas).
        - [x] **UX**: Eliminada selección persistente al hacer click en filas de catálogos (solo teclado).
        - [x] **Estabilidad**: Corregido bug de "empty ID" al crear entidades y error de renderizado en Detail.vue.
        - [x] **DevOps**: Script de despliegue remoto ahora auto-repara permisos de `.git` y `.docker_config`.

    - **📋 PROTOCOLO DE VALIDACIÓN FINAL**:
        1.  **Terminal Taller (MES)**: Verificar que los estados de tarea (Pendiente, En curso) aparecen traducidos y no como IDs técnicos.
        2.  **Terminal Venta (TPV)**: Confirmar el uso de términos profesionales ("BASE IMPONIBLE", "Bonificación Comercial").
        3.  **Formulario Entidades**: Probar el desdoble de nombre en "Persona Física" y el filtrado de documentos por tipo.
        4.  **Listados**: Confirmar que el click no deja la fila resaltada tras cerrar un diálogo.
        5.  **Despliegue**: Verificar que `./scripts/rebuild-staging-remote.ps1` termina sin errores de permisos.
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
