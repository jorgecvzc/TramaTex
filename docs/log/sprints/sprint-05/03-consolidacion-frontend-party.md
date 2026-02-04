# TAREA: 03-consolidacion-frontend-party

- **Sprint:** 05
- **Estado:** ✅ Completado
- **Fecha de Inicio:** 2026-02-04
- **Fecha de Fin:** 2026-02-04
- **Facilitador:** GitHub Copilot, Usuario
- **Jira:**

---

## 🎯 OBJETIVOS

Consolidar el frontend para eliminar referencias legacy (organizations/persons) y dejar una única experiencia “Party/Parties”, manteniendo compatibilidad con `/parties`.

### Fases de la Tarea:
1. [ ] **Renombre de rutas y páginas**
   - [ ] Cambiar rutas internas de `/organizations` a `/parties` como origen principal.
   - [ ] Renombrar páginas/componentes si aplica.
2. [ ] **Ajuste de copy y UI**
   - [ ] Revisar textos visibles y labels para “Party/Parties”.
3. [ ] **Normalización de servicios**
   - [ ] Ajustar `partyApi` para exponer nombres `Party` (manteniendo adaptadores si son necesarios).
4. [ ] **Validación**
   - [ ] Ejecutar tests frontend.
   - [ ] Smoke manual de flujo Party (listado, detalle, creación).

---

## 📋 INFORMACIÓN

- **Descripción:** Consolidación del frontend Party eliminando nomenclatura legacy.
- **Módulos Afectados:** `frontend` (router, pages, components, services)
- **Personas Clave:**

---

## 🚨 BLOQUEADORES

- [ ]

---

## 📝 NOTAS Y REGISTRO DE TRABAJO

### 2026-02-04
- Tarea creada para consolidación de nomenclatura Party en frontend.
- Rutas principales migradas a `/parties` y páginas renombradas a `pages/parties`.
- `partyApi` normalizado con métodos `Party`/`Contact` y compatibilidad legacy.
- Componentes Party actualizados para usar los nuevos métodos.
- Tests frontend ejecutados: `npm run test:unit` OK.
- Smoke UI iniciado con Vite dev server y rutas `/parties` y `/parties/new` abiertas.
- Componentes renombrados a `Party*` y páginas actualizadas a usar nuevos nombres.
- Props y estado internos ajustados a `partyId`/`party` para eliminar naming legacy.
- `partyApi` internals renombrados a `Party` (alias legacy conservados).
- Re-run de tests frontend: `npm run test:unit` falló con “No test suite found” (requiere revisión).
- Re-ejecución desde `apps/frontend`: `npm run test:unit` OK (33 tests).
- Eliminadas rutas alias `/organizations` y métodos legacy de `partyApi`.
- Tests frontend re-ejecutados: `npm run test:unit` OK (33 tests).
- Smoke UI completado en `/parties` y `/parties/new`.
