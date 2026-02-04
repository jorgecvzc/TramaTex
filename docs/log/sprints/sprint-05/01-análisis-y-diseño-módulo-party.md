# TAREA: 01-análisis-y-diseño-módulo-party

> **Consolidación:** Este sprint es el único que concentra todo el trabajo del módulo Party. El contenido histórico fue archivado y no se considera un sprint independiente.

- **Sprint:** 05
- **Estado:** ✅ Completado
- **Fecha de Inicio:** 2026-02-01
- **Fecha de Fin:** 2026-02-02
- **Facilitador:** Gemini, Usuario
- **Jira:**

---

## 🎯 OBJETIVOS

El objetivo de esta tarea es realizar un análisis profundo y un diseño para el módulo **party**. Este proceso servirá como una plantilla para el desarrollo de futuros módulos.

### Fases de la Tarea:
1.  [X] **Fase 1: Definición de Dominio y Casos de Uso**
    - [X] Revisar la documentación existente sobre el módulo `party`.
    - [X] Definir el modelo de dominio en `docs/modules/party/domain-model.md`.
    - [X] Crear un diagrama de dominio en `docs/modules/party/diagrams/domain-model.md`.
    - [X] Documentar los casos de uso en `docs/modules/party/use-cases.md`.
2.  [X] **Fase 2: Creación de ADR (Architectural Decision Record)**
    - [X] Crear un nuevo ADR en `docs/architecture/adrs/ADR-012-arquitectura-modulo-party.md`.
    - [X] Documentar las decisiones de diseño clave, patrones y tecnologías a utilizar.
3.  [X] **Fase 3: Comparación con la Implementación Actual**
    - [X] Analizar el código existente en `apps/tramatex-api/internal/party/` y `apps/frontend/src/modules/party/`.
    - [X] Identificar discrepancias entre el diseño propuesto y la implementación actual.
    - [X] Documentar las discrepancias y proponer un plan de refactorización si es necesario.
4.  [X] **Fase 4: Revisión de UI**
    - [X] Revisar los mockups y componentes de UI existentes o propuestos.
    - [X] Asegurar que la UI se alinee con los casos de uso y el modelo de dominio.
5.  [X] **Fase 5: Aprobación del Usuario**
    - [X] Presentar el análisis y diseño al usuario para su aprobación.

---

## 📋 INFORMACIÓN

- **Descripción:** Análisis profundo del módulo `party` para establecer un diseño robusto y un proceso reutilizable.
- **Módulos Afectados:** `party`
- **Personas Clave:**

---

## 🚨 BLOQUEADORES

- [ ]

---

## 📝 NOTAS Y REGISTRO DE TRABAJO

### 2026-02-01
- Creación de la tarea.
- Inicio de la Fase 1.

### 2026-02-01 - Finalización Fase 1 y 2
- **Fase 1 completada:** Se ha definido y actualizado la documentación de dominio (`domain-model.md`, `diagrams/domain-model.md`, `use-cases.md`) para el módulo `party` según los requisitos iniciales.
- **Fase 2 completada:** Se ha creado y finalizado el ADR `ADR-012-arquitectura-modulo-party.md`, donde se documenta la decisión de adoptar el "Modelo de Party con Roles y Relaciones" y el "Modelo de Contactos Simplificado" para el módulo `Party`. Esta decisión fue colaborativa y aprobada por el usuario.

### 2026-02-01 - Hallazgos Fase 3: Comparación con la Implementación Actual

Tras un análisis exhaustivo del código existente en comparación con el diseño aprobado en `ADR-012-arquitectura-modulo-party.md`, se han identificado las siguientes discrepancias:

**1. Capa de Dominio (`apps/tramatex-api/internal/party/domain/`):**
-   **Modelo Fundamental:** La implementación actual se basa en `Organization` como Aggregate Root, con `Person` y `Address` dependientes. El nuevo diseño requiere un `Party` Aggregate Root, con `PersonProfile` y `OrganizationProfile` (que implementan `PartyProfile`), `PartyRole`, `PartyRelationship` y `ContactDetails`.
-   **Entidades y VOs:** `Organization` y `Person` tal como existen son incompatibles. Se requiere la creación de las nuevas entidades `Party`, `PartyRelationship`, `PartyRole` y el Value Object `ContactDetails`.
-   **ID y Enums:** `OrganizationID`, `PersonID`, `OrganizationRole`, `OrganizationStatus` deben ser reemplazados o adaptados a los nuevos conceptos de `PartyID`, `PartyRole` y un estado de `Party` más genérico.

**2. Esquema de Base de Datos (`apps/tramatex-api/migrations/002_create_party_tables.sql`):**
-   **Desalineación Completa:** El esquema actual (`organizations`, `persons`, `addresses` tablas) es totalmente incompatible con el diseño aprobado.
-   **Refactorización Mayor:** Se requerirá una migración completamente nueva para eliminar las tablas existentes y crear el nuevo esquema (`parties`, `person_profiles`, `organization_profiles`, `party_roles`, `party_relationships`, `contact_details`).

**3. Contratos de API (`docs/modules/party/api-contracts.md`):**
-   **Centrado en `Organization`:** La API actual está enteramente construida alrededor de la entidad `Organization`.
-   **Rediseño Total:** Se necesitará renombrar endpoints (`/organizations` a `/parties`), reestructurar DTOs para un `PartyDTO` genérico (con perfiles de Persona/Organización), y crear nuevos endpoints para `PartyRoles`, `PartyRelationships` y `ContactDetails`.

**4. Código Frontend (`apps/frontend/src/...` relacionado con `party`):**
-   **Componentes y Rutas:** Los nombres de componentes (`OrganizationForm`, `OrganizationList`) y rutas (`/organizations`) son específicos de `Organization` y deberán cambiar a `Party`.
-   **Lógica y Datos:** Las stores de Pinia y los servicios (`partyApi.js`) deberán ser refactorizados para consumir la nueva API centrada en `Party` y adaptarse a los nuevos DTOs y lógica de negocio (diferenciando `PersonProfile` y `OrganizationProfile`, manejando roles, etc.).

**Conclusión de la Fase 3:**

Existe una **desalineación completa** entre la implementación actual del módulo `party` y el diseño aprobado en `ADR-012`. Por lo tanto, se requiere una **refactorización mayor y comprehensiva** en todas las capas de la aplicación (dominio, persistencia, API, frontend) para implementar el nuevo modelo de dominio `Party`.

### 2026-02-01 - Cierre de Tarea (Revertido)
- Se revierte el cierre y se reabre la tarea por solicitud del usuario.
- Fases 4 y 5 quedan pendientes hasta nueva validación.

### 2026-02-02 - Fase 4: Revisión de UI (Hallazgos)

**Fuentes revisadas:**
- UI actual Party (frontend):
  - [apps/frontend/src/components/party/OrganizationList.vue](apps/frontend/src/components/party/OrganizationList.vue)
  - [apps/frontend/src/components/party/OrganizationForm.vue](apps/frontend/src/components/party/OrganizationForm.vue)
  - [apps/frontend/src/components/party/OrganizationDetail.vue](apps/frontend/src/components/party/OrganizationDetail.vue)
  - [apps/frontend/src/components/party/PersonManager.vue](apps/frontend/src/components/party/PersonManager.vue)
  - [apps/frontend/src/components/party/AddressManager.vue](apps/frontend/src/components/party/AddressManager.vue)
  - [apps/frontend/src/services/partyApi.js](apps/frontend/src/services/partyApi.js)
- Diseño global (no específico de Party):
  - [docs/architecture/design-system/palette.md](docs/architecture/design-system/palette.md)
  - [docs/architecture/design-system/theme.md](docs/architecture/design-system/theme.md)
  - [docs/architecture/design-system/typography.md](docs/architecture/design-system/typography.md)
  - Mockups generales: [agents/project/context/mockups/code.html](agents/project/context/mockups/code.html), [agents/project/context/mockups/2/code.html](agents/project/context/mockups/2/code.html)

**Hallazgos principales (alineación UI ↔ dominio/ADR):**
1. **La UI actual es 100% “Organization-centric”.** Pantallas, rutas y servicio están basados en `/organizations` y entidades `Organization/Person/Address`. Esto **no se alinea con el modelo aprobado en ADR-012** (Party + Profiles + Roles + Relationships + ContactDetails).
2. **No hay soporte UI para los nuevos conceptos clave** del ADR:
    - Selección/gestión de **tipo de Party** (persona vs organización) y sus **perfiles**.
    - **Roles múltiples** por Party (Cliente, Proveedor, Empleado, etc.).
    - **Relaciones entre Party** (matriz/filial, empleado-de, etc.).
    - **ContactDetails** con `relatedPartyId` (contacto vinculado a otra Party).
3. **Flujos actuales cubren parcialmente los casos de uso antiguos**, pero no los nuevos:
    - Sí cubre: listar organizaciones, crear/editar, cambiar estado, añadir personas y direcciones.
    - No cubre: gestión de roles múltiples, relaciones entre parties, perfiles simultáneos, ni contactos simplificados.
4. **No existen mockups específicos de Party**. Solo hay lineamientos generales y mockups de login/ordenes. Se requiere diseño UI específico del módulo Party para el nuevo modelo.

**Hallazgos de consistencia UI ↔ estándares del proyecto:**
- El frontend **no usa Tailwind** en estos componentes (se usa CSS scoped). Esto contradice el estándar de frontend en [agents/project/context/code-standards.yaml](agents/project/context/code-standards.yaml) y [agents/project/context/tech-stack.yaml](agents/project/context/tech-stack.yaml).
- Componentes exceden el tamaño recomendado y mezclan lógica/UI en un solo archivo.

**Recomendación de rediseño UI (alto nivel):**
1. **Party List** (antes OrganizationList): lista única de parties con tipo, roles, estado, taxId (si aplica), y filtros por rol/tipo.
2. **Party Create/Edit**: formulario que inicia con selector de tipo (Persona/Organización) y roles, luego perfiles correspondientes.
3. **Party Detail** con secciones:
    - Perfil(es) (persona/organización)
    - Roles
    - Relaciones
    - ContactDetails (con opcional `relatedPartyId`)
    - Direcciones (si aplica)
4. **Unificar rutas y servicio** a `/parties` y DTOs alineados a ADR-012.
5. **Implementar UI con Tailwind** y layout consistente con el design system.

**Estado:** La Fase 4 queda documentada con hallazgos.

### 2026-02-02 - Fase 5: Aprobación del Usuario
- Aprobación explícita del usuario para cerrar el análisis y diseño del módulo Party.
- Se marca la tarea como completada.