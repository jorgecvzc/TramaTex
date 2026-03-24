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
- **Status:** ✅ Completado
- **Sprint:** N/A
- **Started:** 2026-03-10
- **Finished:** 2026-03-24
- **Rama:** `infra/pcele-staging-deploy` → Mergeada a `develop` y `master`

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
- [x] Crear Droplet en DigitalOcean y configurar GitHub Secrets. ✅ 2026-03-23
- [ ] Configurar DNS + SSL/Let's Encrypt en producción. *(pendiente — sesión `infra-production-digitalocean-2026-03-23`)*

### Bugs Corregidos en pcele (2026-03-23)

- `postgres` container tenía `restart: no` → no arrancaba tras reinicio del sistema. Fix: `docker update --restart unless-stopped` + `docker-compose.remote.yml` actualizado.
- Nginx cacheaba la IP de `api` al arrancar → 502 cuando API se reiniciaba. Fix: directiva `resolver 127.0.0.11 valid=10s` + `set $api_backend` en `nginx.conf`.
- `crypto.randomUUID()` no disponible en contextos HTTP (pcele sin HTTPS) → error al crear Party. Fix: eliminado `crypto.randomUUID`, reemplazado por generador UUID v4 propio con `Math.random()` en `partyApi.ts` y `PartyForm.vue`.
- Docker build usaba capas cacheadas → bundle JS viejo desplegado aunque el source cambiara. Fix: rebuild con `--no-cache` + recrear contenedor con `docker compose up --force-recreate` (no `docker restart`).
- Contenedor `tramatex_db` creado por docker-compose v1 (nombre con prefijo hash) → al recrear frontend con compose v2, la DB se paraba. Fix: migrar DB a docker compose v2 (`docker rm` + `docker compose up --no-deps postgres`).
- `index.html` sin headers de caché → el navegador servía el bundle antiguo. Fix: añadido bloque `location = /index.html` con `Cache-Control: no-cache, no-store, must-revalidate` en `nginx.conf`.

### Commits relevantes

- `a557f41` — feat: infraestructura de despliegue multientorno (Dockerfile, Nginx, GitHub Actions, Makefile)
- `4c04b65` — chore: eliminar archivos de tracking que están en .gitignore
- `daf0331` — feat: completar sistema de despliegue (env examples, SSL, scripts, docs)
- `80702e2` — fix: usar 127.0.0.1 en healthcheck de Dockerfiles (IPv6 en Alpine)
- `485b8b5` — fix(infra): añadir restart unless-stopped a postgres en docker-compose.remote.yml
- `249138c` — fix(infra): usar resolver DNS dinámico en Nginx para evitar 502 al reiniciar API
- `f215450` — merge: infra/pcele-staging-deploy → develop
- `48a0957` — fix(party): eliminar crypto.randomUUID, usar generador UUID propio
- `dc0ad29` — fix(nginx): no-cache para index.html, evita que el navegador sirva bundle antiguo
- `865cf0e` — ci: corregir workflows — usar develop, --no-cache y --force-recreate en deploy-production *(rama: infra/production-digitalocean)*

---

## Despliegue en Producción (DigitalOcean)

- **Session ID:** `infra-production-digitalocean-2026-03-23`
- **Status:** En Progreso (DNS + SSL pendientes)
- **Sprint:** N/A
- **Started:** 2026-03-23
- **Rama:** `infra/production-digitalocean` → Mergeada a `develop` y `master`

### Contexto

Despliegue automático de TramaTex en producción vía DigitalOcean + GitHub Actions. El workflow `deploy-production.yml` se activa en push a `master` → build de imágenes en Actions runner → push a GHCR → SSH al Droplet → `docker pull` + `docker compose up --force-recreate`.

### Infraestructura aprovisionada (2026-03-23)

| Componente | Estado | Detalle |
|---|---|---|
| Droplet DigitalOcean | ✅ Creado | IP: `46.101.188.130`, Ubuntu 24.04 |
| Docker CE | ✅ Instalado | v29.3.0, `docker compose` v5.1.1 |
| Usuario `tramatex` | ✅ Creado | Grupo `docker`, SSH configurado |
| Repo clonado | ✅ OK | `/opt/tramatex`, rama `master` @ `f73e8d7` |
| SSH Key (sin passphrase) | ✅ OK | `tmp/do-setup/deploy_final` + `.pub` |
| 4 GitHub Secrets | ✅ OK | `PROD_IP`, `SSH_USER`, `SSH_PRIVATE_KEY`, `ENV_PROD` |

### Estado del CI/CD ✅ OPERATIVO (2026-03-24)

| Workflow | Estado | Notas |
|---|---|---|
| `backend.yml` | ✅ Corregido | `main` → `develop`, Go 1.21 → 1.23 |
| `frontend.yml` | ✅ Corregido | `main` → `develop` |
| `deploy-staging.yml` | ✅ OK | CI en `staging` + `develop` |
| `deploy-production.yml` | ✅ **Operativo** | 3 jobs: `build-api` + `build-frontend` + `deploy` |

**Estrategia CI/CD actual:**
- `build-api` → `docker build -f Dockerfile` → push `ghcr.io/jorgecvzc/tramatex-api:latest`
- `build-frontend` → `docker build -f Dockerfile.frontend` → push `ghcr.io/jorgecvzc/tramatex-frontend:latest`
- `deploy` → SSH al Droplet → `docker pull` + `docker compose up --force-recreate`
- El Droplet NO hace builds (limitación 1 GB RAM evita OOM en `npm run build`)

### Estado de contenedores en producción (verificado 2026-03-24)

```
tramatex_frontend   Up (healthy)   ghcr.io/jorgecvzc/tramatex-frontend:latest
tramatex_api        Up (healthy)   ghcr.io/jorgecvzc/tramatex-api:latest
tramatex_db         Up (healthy)   postgres:15-alpine
```

- **Frontend**: `http://46.101.188.130` → HTTP 200 ✅
- **API Health**: `http://46.101.188.130/api/health` → `{"status":"healthy"}` ✅

### Bugs resueltos (2026-03-23)

| Error | Causa | Fix |
|---|---|---|
| `apt-get upgrade` cuelga interactivo | Sin `DEBIAN_FRONTEND=noninteractive` | Añadido + `--force-confdef --force-confold` |
| `docker.io` no encontrado | No está en repos Ubuntu | Instalado desde repo oficial Docker CE |
| `chown: invalid user tramatex` | Usuario no existía | `useradd -m -s /bin/bash tramatex` manual |
| `Could not resolve host: github.com` | systemd-resolved recién instalado | DNS OK tras estabilizarse |
| `Run Command Timeout` (30s) | Build Go+npm supera timeout default | Aumentado `command_timeout` |

### Bugs resueltos (2026-03-24)

| Error | Causa | Fix |
|---|---|---|
| `Permission denied (publickey)` | Secret `SSH_PRIVATE_KEY` en GitHub tenía valor incorrecto | Copiar `tmp/do-setup/deploy_final` → actualizar secret manualmente |
| `Run Command Timeout` (30 min) | Droplet 1 GB RAM — OOM al hacer `npm run build` dentro de Docker | Mover el build a GitHub Actions runner (7 GB RAM) → push imágenes a GHCR → Droplet solo hace `docker pull` |
| `Run Command Timeout` (10 min) | Primera descarga de imágenes GHCR supera el timeout | Aumentar `command_timeout` a `20m` |
| `tramatex_frontend` en `Restarting` ejecutando binario Go | BuildKit cache collision: 2 builds en mismo job reusaban capas, `tramatex-frontend:latest` recibía el binario Go API | Separar en 2 jobs independientes (`build-api`, `build-frontend`) |
| `tramatex_frontend` en `Restarting` aún después de jobs separados | `docker/build-push-action@v5` ignora el parámetro `dockerfile:` y siempre usa `Dockerfile` por defecto | Reemplazar la action por comando `docker build --no-cache -f Dockerfile.frontend` explícito |

### Próximos Pasos

- [x] **URGENTE**: `chown tramatex:tramatex /home/tramatex` — ejecutado ✅
- [x] Relanzar workflow → "Re-run all jobs" ✅
- [x] Verificar deploy completa sin errores ✅
- [x] Probar en navegador: `http://46.101.188.130` → pantalla de login ✅
- [x] Verificar API: `curl http://46.101.188.130/api/health` → `{"status":"healthy"}` ✅
- [ ] Configurar DNS del dominio apuntando a `46.101.188.130`
- [ ] Instalar certbot + SSL: `certbot --nginx -d tudominio.com`
- [ ] Activar `nginx-ssl.conf` (descomentar volumes SSL en `docker-compose.remote.yml`)

### Claves SSH (en `tmp/do-setup/`, en .gitignore)

- `deploy_final` — clave privada para GitHub Secret `SSH_PRIVATE_KEY`
- `deploy_final.pub` — clave pública (ya en `/home/tramatex/.ssh/authorized_keys` del Droplet)
- Clave pública: `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBeWnhbHBUsNJFRioqoh/d+5cUpztmViGEhDiwWuuiOE github-actions-tramatex`

### Commits relevantes

- `20dabf7` — fix(nginx-ssl): aplicar mismo fix resolver DNS y no-cache index.html que nginx.conf
- `17e4374` — fix(ci): aumentar timeout ssh-action a 30m para builds en Droplet
- `aee6896` — ci(deploy): build images in GitHub Actions, push to GHCR — avoid OOM on Droplet
- `0bdb70c` — ci(deploy): increase SSH command_timeout to 20m for first pull from GHCR
- `f936eef` — ci(deploy): split API and frontend builds into separate jobs to fix BuildKit cache collision
- `f73e8d7` — ci(deploy): replace build-push-action with explicit docker build+push to fix dockerfile selection bug

### Archivos de Contexto

- `.github/workflows/deploy-production.yml`
- `docker/docker-compose.remote.yml`
- `.dockerignore`
- `docker/nginx.conf`
- `docker/nginx-ssl.conf`
- `tmp/do-setup/env-prod.txt` (valores .env producción)
- `tmp/do-setup/deploy_final` (clave privada SSH)

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
