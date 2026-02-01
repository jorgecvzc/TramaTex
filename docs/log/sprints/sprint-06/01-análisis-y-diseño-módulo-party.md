# TAREA: 01-análisis-y-diseño-módulo-party

- **Sprint:** 06
- **Estado:** ✅ Completado
- **Fecha de Inicio:** 2026-02-01
- **Fecha de Fin:** 2026-02-01
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

### 2026-02-01 - Cierre de Tarea
- Fase 4 completada: revisión de UI con base en los mockups y el sistema de diseño.
- Fase 5 completada: validación y cierre con el usuario.
- Estado final: tarea concluida y lista para planificación de refactorización.