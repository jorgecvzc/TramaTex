# Bitácora de Sesiones de Desarrollo

<!--
Este archivo registra las sesiones de desarrollo.

SECCIONES:
1. SESIONES ABIERTAS: Contiene las sesiones de trabajo que están en progreso, pausadas o bloqueadas. El objetivo es detallar el contexto y los próximos pasos.
2. REGISTRO DE SESIONES CERRADAS: Un archivo histórico de todas las sesiones completadas, conservando solo metadatos esenciales.

ESTRUCTURA DE UNA SESIÓN ABIERTA:
- Título (##): Un H2 con un título descriptivo y único.
- Metadatos:
  - **Session ID:** `identificador-unico-kebab-case` (OBLIGATORIO Y ÚNICO)
  - **Status:** (En Progreso | En Pausa | Bloqueado)
  - **Sprint:** Sprint XX
  - **Started:** Fecha de inicio (YYYY-MM-DD).
- Contexto: Breve descripción del objetivo de la sesión.
- Próximos Pasos: Checklist de las tareas pendientes.
- Archivos de Contexto: Lista de archivos clave.

ESTRUCTURA DE UNA SESIÓN CERRADA (en el registro):
- Una línea de lista con: **[Título]** | Iniciada: [Fecha YYYY-MM-DD] | Finalizada: [Fecha YYYY-MM-DD]
-->
---
# SESIONES ABIERTAS

## Implementación de Infraestructura de Despliegue Multientorno

- **Session ID:** `infra-multi-env-deployment-impl-2026-03-10`
- **Status:** En Progreso
- **Sprint:** N/A
- **Started:** 2026-03-10
- **Rama:** `infra/pcele-staging-deploy` → Merge a `develop` al finalizar

### Contexto

Implementación técnica de la estrategia de despliegue multientorno definida en el estudio previo (`tmp/estudio_despliegue_tramatex.md`). El objetivo es configurar la jerarquía de ramas (develop -> staging -> master), la orquestación con Nginx y la automatización CI/CD hacia DigitalOcean.      

### Progreso

- [x] Crear rama `infra/multi-env-deployment`.
- [x] Implementar `Dockerfile.frontend` (multi-stage build).
- [x] Crear configuración `docker/nginx.conf` (Proxy inverso y SPA routing).
- [x] Actualizar `docker/docker-compose.remote.yml` para incluir el servicio Nginx.
- [x] Crear Workflows de GitHub Actions para despliegue automático en `staging` y `master`.
- [x] Refinar `Makefile` para soportar perfiles de despliegue (`pcele`, `staging`, `prod`).
- [x] Limpiar residuos del repositorio (archivos en `.gitignore` eliminados del tracking).
- [x] Limpiar ramas obsoletas (eliminadas: `product-module-validation`, `scaffolding-refinement`, `scaffolding-system-enhancement`, `qa/full-verification`).
- [x] Corregir Makefile — paths con prefijo `docker/`, migrar a `docker compose` v2.
- [x] Crear `.env.example` por entorno: `docker/.env.example`, `.env.pcele.example`, `.env.production.example`.
- [x] Crear `docker/nginx-ssl.conf` para producción con Let's Encrypt.
- [x] Actualizar `start-dev.ps1` — soporte parámetro `-Full` para levantar frontend/Nginx.
- [x] Actualizar `stop-dev.ps1` — migrar a `docker compose` v2, `--profile full`.
- [x] Actualizar `docker-compose.remote.yml` — soporte SSL (volumes comentados).
- [x] Crear guía de despliegue: `docs/guides/developer/deployment-guide.md`.

### Próximos Pasos

- [x] Verificar build Docker local (bloqueado temporalmente por conectividad Docker Hub).
- [x] Crear rama `staging` desde `develop` (después de merge).
- [x] Primer despliegue en pcele (requiere configurar Docker + SSH). ✅ 2026-03-23
- [ ] Crear Droplet en DigitalOcean y configurar GitHub Secrets.
- [ ] Configurar DNS + SSL/Let's Encrypt en producción.

### Bugs Corregidos en pcele (2026-03-23)

- `postgres` container tenía `restart: no` → no arrancaba tras reinicio del sistema. Fix: `docker update --restart unless-stopped` + `docker-compose.remote.yml` actualizado.
- Nginx cacheaba la IP de `api` al arrancar → 502 cuando API se reiniciaba. Fix: directiva `resolver 127.0.0.11 valid=10s` + `set $api_backend` en `nginx.conf`.

### Commits

- `a557f41` — feat: infraestructura de despliegue multientorno (Dockerfile, Nginx, GitHub Actions, Makefile)
- `4c04b65` — chore: eliminar archivos de tracking que están en .gitignore
- *(pendiente)* — feat: completar sistema de despliegue (env examples, SSL, scripts, docs)

---



## QA — Verificación de Calidad Integral

- **Session ID:** `qa-full-verification-2026-03-21`
- **Status:** ✅ Completado
- **Sprint:** N/A
- **Started:** 2026-03-21
- **Finished:** 2026-03-22
- **Rama:** `qa/full-verification` → Merged a `develop`

### Contexto

Sesión de verificación de calidad completa del sistema tras todos los refactors realizados (Sprint 14, sort_order, fragmentación SalesService, estandarización errores, migración UUID IAM). El objetivo es comprobar que todo funciona correctamente end-to-end y subsanar posibles errores antes de continuar con las sesiones posteriores.

### Próximos Pasos

- [x] **Revisión Sprint 14:** Sprint 14 validado como unidad coherente.
- [x] **Planificación de sprints:** Sesiones posteriores no requieren sprints formales (son tareas post-MVP).
- [x] Ejecutar suite completa de tests Go (30 paquetes) y verificar 0 fallos. *(1 fallo intermitente en `product/infra/persistence` por contaminación de estado BD entre paquetes — pasa en aislamiento. Todos los tests de sales: OK)*
- [x] Ejecutar suite completa de tests frontend (Vitest, 230+ tests) y verificar 0 fallos. *(11 archivos, 230 tests — todos verdes)*
- [x] Levantar entorno Docker completo (`docker-compose up`) y verificar que las migraciones se ejecutan correctamente (incluida 009). *(47 tablas, 7 migraciones OK — BD `tramatex`)*
- [x] Verificar flujo completo de login (crear usuario, autenticación, token JWT). ✅
- [x] Verificar CRUD de Party (crear, editar, listar, eliminar). ✅
- [x] Verificar CRUD de Product (crear con atributos, verificar SKU generado con orden correcto de DirectAttributeIDs). ✅
- [x] Verificar generación de variantes y que el SKU refleja el orden de atributos. ✅
- [x] Verificar flujo de Sales (crear pedido, añadir líneas, calcular totales, facturar). ✅
- [x] Verificar Pricing (reglas de precio, descuentos). ✅
- [x] Verificar MES (órdenes de producción, terminal de taller). ✅
- [x] Documentar bugs encontrados y corregirlos in situ. *(3 bugs corregidos: ver sección de Bugs Encontrados abajo)*

### Bugs Encontrados y Corregidos

| Commit | Bug | Fix |
|--------|-----|-----|
| `a6983d4` | `salesApi.ts` — `normalizeEntity()` traducía estados español→inglés. Todos los componentes comparaban contra claves en inglés | Eliminada capa de traducción; `getStatusLabel/Class` y comparaciones actualizadas a claves en español en 6 componentes Vue |
| `2ba7643` | `billing_service.go` — `CreateInvoice` usaba `item.Quantity` completa sin restar ya-facturado por albaranes previos → doble facturación | Calcular mapa `alreadyInvoiced` antes de construir líneas; usar cantidad residual |
| `7288060` | `billing_service.go` — `updateOrderInvoiceStatus` llamada después de persistir la factura. `ListBySalesOrderID` ya incluía la nueva → segundo loop `newInvoiceItems` doblaba el conteo → pedido marcado `FACTURADO_COMPLETAMENTE` prematuramente | Eliminar parámetro `newInvoiceItems` y su loop |
| `141fcc9` | `QuoteCreate/OrderCreate/QuoteDetail/OrderDetail.vue` — `calculateLineSubtotal` usaba `Array.find()` por `productVariantId` → con múltiples líneas de la misma variante, `find()` siempre devolvía la primera coincidencia, repitiendo su subtotal en las demás | Cambiar a mapeo por índice posicional, teniendo en cuenta el filtro de `buildPreviewItems` |

### Archivos de Contexto

- `agents/project/sprint-registry.yaml`
- `apps/tramatex-api/run_tests.sh`
- `apps/frontend/vitest.config.ts`
- `docker/docker-compose.yml`
- `tmp/manual-testing-guide.md`

---

## Alineación Documental Post-Refactors Sprint 14

- **Session ID:** `doc-alignment-post-sprint14-2026-03-21`
- **Status:** En Pausa (Pendiente de inicio)
- **Sprint:** N/A
- **Started:** 2026-03-21
- **Rama:** Cerrar rama actual → Crear rama nueva para esta sesión → Merge a `develop` al finalizar

### Contexto

Tras los refactors del Sprint 14 (fragmentación de SalesService, estandarización de errores, migración UUID en IAM, eliminación de sort_order en Product), la documentación técnica ha quedado desalineada con el código. Se ha elaborado un roadmap de tareas documentales en `tmp/documentation-alignment-roadmap.md`.

### Próximos Pasos

- [ ] **IAM:** Actualizar `docs/modules/iam/domain-model.md` — UserID como `uuid.UUID`, eliminar menciones a `createdAt`/`updatedAt`.
- [ ] **Sales:** Actualizar `docs/modules/sales/implementation-guide.md` y `module-spec.md` con nueva estructura de servicios fragmentados.
- [ ] **Sales:** Documentar `calculations.go` en `docs/modules/sales/domain-model.md`.
- [ ] **MES/Errores:** Actualizar guías de implementación de Product, Sales y MES para indicar delegación de errores HTTP al middleware `shared`.
- [ ] **Estructura:** Actualizar `docs/guides/developer/project-structure-details.md` con nuevos paths.
- [ ] **Product:** Documentar eliminación de `sort_order` y nuevo flujo de ordenamiento por `DirectAttributeIDs`.
- [ ] **Frontend:** Verificar que docs no referencien `apps/frontend/src/pages/organizations`.
- [ ] **Agents:** Actualizar agent contexts en `agents/project/context/` (`architecture.yaml`, `bounded-contexts.yaml`, `code-standards.yaml`, `tech-stack.yaml`) para reflejar los refactors del Sprint 14.

### Archivos de Contexto

- `tmp/documentation-alignment-roadmap.md`
- `docs/modules/iam/domain-model.md`
- `docs/modules/sales/implementation-guide.md`
- `docs/modules/sales/module-spec.md`
- `docs/guides/developer/project-structure-details.md`

---

## Mejora UI/UX — Unificación Estética y Componentes Base

- **Session ID:** `ui-ux-improvement-post-mvp-2026-03-21`
- **Status:** En Pausa (Pendiente de inicio)
- **Sprint:** Post-MVP
- **Started:** 2026-03-21
- **Rama:** Cerrar rama actual → Crear rama nueva para esta sesión → Merge a `develop` al finalizar

### Contexto

Auditoría de UI/UX completada en `tmp/ui-ux-improvement-suggestions.md`. Se detectaron inconsistencias entre módulos: patrones de navegación mixtos en listados, botones sin base global, emojis en lugar de iconos Lucide, paleta de colores fragmentada y layouts con `max-width` variables. Esta sesión implementa el plan de mejora estética para dar coherencia visual antes de la presentación TFM. También sirve como validación del trabajo Post-MVP y del flujo de despliegue definido en la sesión de infraestructura (`infra-multi-env-deployment-impl-2026-03-10`).

**Dependencias:** Requiere completar previamente la sesión de despliegue multientorno.

### Próximos Pasos

- [ ] Crear `apps/frontend/src/design-system/_buttons.css` con estilos globales (`primary`, `secondary`, `outline`, `danger`) e importar en `theme.css`.
- [ ] Estandarizar iconografía: eliminar todos los emojis de la interfaz y sustituir por Lucide Icons (🗑️→Trash2, 🖨️→Printer, 💰→Euro, ⚠️→AlertTriangle, ⚙️→Settings).
- [ ] Unificar comportamiento de listados: fila clickeable + botón de acción iconográfico al final.
- [ ] Crear componente `BasePageHeader` (Breadcrumb + Título + Acciones) y aplicar a todas las páginas.
- [ ] Estandarizar `max-width` de contenedores y jerarquía de cabeceras (H1/H2/H3) entre módulos.
- [ ] Forzar uso estricto de variables de `_variables.css` (paleta, radios de borde, sombras).
- [ ] Migrar `PartyList.vue` como primer listado estandarizado de referencia.
- [ ] Refactorizar `PartyForm.vue` para sustituir `fieldset/legend` por diseño de tarjetas.
- [ ] Mejorar contraste de etiquetas de formularios (peso 500, `--color-text-secondary`).
- [ ] Estandarizar dropdowns/selects y definir patrones visuales de validación/errores en formularios.

### Archivos de Contexto

- `tmp/ui-ux-improvement-suggestions.md`
- `apps/frontend/src/assets/styles/theme.css`
- `apps/frontend/src/assets/styles/_variables.css`
- `apps/frontend/src/components/`

---

## Preparación TFM — Presentación Final de TramaTex

- **Session ID:** `tfm-final-presentation-2026-03-21`
- **Status:** En Pausa (Pendiente de inicio)
- **Sprint:** N/A
- **Started:** 2026-03-21
- **Rama:** Cerrar rama actual → Crear rama nueva para esta sesión → Merge a `develop` (y a `main` como entrega final)

### Contexto

TramaTex se presenta como Trabajo Fin de Máster (TFM). Esta sesión cubre la preparación integral del proyecto para su entrega y defensa académica: revisión de documentación, presentación, memoria, y asegurar que el estado del código, los tests y el despliegue son coherentes y presentables. Es la **última sesión** del proyecto.

### Próximos Pasos

- [ ] Revisar y actualizar la presentación existente (`docs/presentations/tramatex-presentation.md`, `TramaTex_Presentacion_Final.pptx`).
- [ ] Asegurar que `README.md` del proyecto refleja el estado final (visión, arquitectura, instrucciones de instalación/ejecución).
- [ ] Verificar que la documentación de arquitectura (`docs/architecture/`) está completa y actualizada.
- [ ] Confirmar que todos los módulos tienen documentación consistente en `docs/modules/`.
- [ ] Validar que los tests pasan limpiamente (Go + Vitest) y documentar cobertura.
- [ ] Revisar `CONTRIBUTING.md` y `LICENSE.md` para coherencia académica.
- [ ] Preparar memoria/informe TFM si es necesario (estructura, introducción, conclusiones, trabajo futuro).
- [ ] Generar diapositivas de presentación para la defensa del TFM (basarse en `docs/presentations/slides_spec.md` y `tramatex-presentation.md`).
- [ ] Limpiar archivos temporales en `tmp/` que no deban ir en la entrega final.
- [ ] Verificar que el despliegue Docker funciona correctamente de principio a fin.

### Archivos de Contexto

- `docs/presentations/tramatex-presentation.md`
- `docs/presentations/TramaTex_Presentacion_Final.pptx`
- `docs/presentations/slides_spec.md`
- `docs/architecture/architecture-vision.md`
- `README.md`
- `CONTRIBUTING.md`

---

---
# REGISTRO DE SESIONES CERRADAS
---
- **QA — Verificación de Calidad Integral** | Iniciada: 2026-03-21 | Finalizada: 2026-03-22 | Status: ✅ COMPLETADO | Rama: `qa/full-verification` → `develop` | 5 commits, 4 bugs corregidos, QA manual 6/6 puntos OK
- **Refactor sort_order → DirectAttributeIDs (Producto/Atributos)** | Iniciada: 2026-03-21 | Finalizada: 2026-03-21 | Status: ✅ COMPLETADO
- **Análisis de Refinamiento Arquitectónico del MVP (Sprint 14)** | Iniciada: 2026-03-12 | Finalizada: 2026-03-21 | Status: ✅ COMPLETADO | Ver: [sprint-14](docs/log/sprints/sprint-14/sprint-14-summary.md) | PR pendiente: `mvp-arch-refinement` → `develop`
- **Análisis de Refinamiento Arquitectónico del Módulo MES** | Iniciada: 2026-03-20 | Finalizada: 2026-03-20 | Status: ✅ COMPLETADO
- **Integración MES-Sales: Terminal de Taller y Visibilidad de Producción en Pedidos** | Iniciada: 2026-03-19 | Finalizada: 2026-03-19
- **Refinamiento y Estabilización ERP Core** | Iniciada: 2026-03-09 | Finalizada: 2026-03-14
