# [NOMBRE DEL MÓDULO] - Development Checklist & Documentation Hub

**Bounded Context:** `[Contexto Delimitado]`  
**Responsabilidad Principal:** `[Descripción breve]`  

---
> **Instrucciones:** Este documento es una checklist viva. Úsalo para guiar y rastrear el desarrollo del módulo a través de cada fase del **Standard Module Development Workflow (SMDW)**.
> 
> **NOTA CRÍTICA: Cada fase debe ser discutida y aprobada por el equipo de desarrollo humano antes de pasar a la siguiente.**

---

## Phase 0: Scaffolding & Initial Setup

- [ ] **1. Crear Estructura de Módulo:**
  - [ ] Backend: `apps/<backend_app_name>/internal/<module_name>/` (ej. `apps/tramatex-api/internal/<module_name>`)
  - [ ] Frontend: `apps/<frontend_app_name>/src/modules/<module_name>/` (si aplica, ej. `apps/frontend/src/modules/<module_name>`)
  - [ ] Otro: `apps/<other_app_name>/<module_name>/` (si aplica)
- [ ] **2. Crear Hub de Documentación:**
  - [ ] Carpeta: `docs/modules/<module_name>/`
  - [ ] Diagramas: `docs/modules/<module_name>/diagrams/`
- [ ] **3. Crear Tarea de Sprint:**
  - [ ] Archivo de tarea creado en `docs/log/sprints/sprint-XX/`
- [ ] **✅ Fase 0 Aprobada por el Equipo**

---

## Phase 1: Domain-Driven Design & API Specification

### 1.1. Domain Model (`domain-model.md`)
- [ ] Definir Entidades y Agregados.
- [ ] Definir Value Objects.
- [ ] Definir Servicios de Dominio (si es necesario).
- [ ] Establecer Reglas de Negocio Invariantes.

### 1.2. Domain Diagram (`diagrams/domain-model.md`)
- [ ] Crear diagrama de clases Mermaid con Entidades, VOs y relaciones.

### 1.3. Use Cases (`use-cases.md`)
- [ ] Listar todos los casos de uso de la capa de aplicación.
- [ ] Describir actores, pre-condiciones y post-condiciones para cada uno.

### 1.4. API Contracts (`api-contracts.md`)
- [ ] Definir Endpoints (rutas, métodos HTTP).
- [ ] Definir DTOs de Petición (Request).
- [ ] Definir DTOs de Respuesta (Response).
- [ ] Especificar Códigos de Error y sus significados.
- [ ] **✅ Fase 1 Aprobada por el Equipo**

---

## Phase 2: Backend Implementation (TDD)

### 2.1. Domain Layer
- [ ] **Tests Unitarios:** Escribir tests para la lógica de negocio (Entidades, VOs, Servicios de Dominio).
- [ ] **Implementación:** Escribir el código de dominio para que los tests pasen.

### 2.2. Application Layer
- [ ] **Tests de Aplicación:** Escribir tests para los Casos de Uso (mockear repositorios).
- [ ] **Implementación:** Escribir la orquestación de los casos de uso.

### 2.3. Infrastructure & Interfaces Layers
- [ ] **Implementación de Repositorios:** Crear la implementación concreta de los repositorios (ej. `[NOMBRE_HERRAMIENTA_ORM_PERSISTENCIA]`).
- [ ] **Implementación de Handlers:** Crear los manejadores HTTP (ej. `[NOMBRE_FRAMEWORK_WEB]`).
- [ ] **Migraciones DB:** Crear y aplicar las migraciones de base de datos necesarias.

### 2.4. Validation
- [ ] Todos los tests del backend pasan (`[COMANDO_TEST_BACKEND]`).
- [ ] Linter del backend pasa (`[COMANDO_LINTER_BACKEND]`).
- [ ] Cobertura de tests cumple los objetivos de calidad (ej: 100% en Dominio crítico).
- [ ] **✅ Fase 2 Aprobada por el Equipo**

---

## Phase 3: Frontend Implementation (si aplica)

- [ ] **State Management:** Crear stores de `[NOMBRE_HERRAMIENTA_GESTION_ESTADO]`.
- [ ] **Logic:** Crear servicios de API y `[TIPO_LOGICA_REUTILIZABLE]`.
- [ ] **UI:** Crear componentes de `[NOMBRE_FRAMEWORK_UI]`.
- [ ] **Tests:** Escribir tests unitarios para stores, `[TIPO_LOGICA_REUTILIZABLE]` u otra lógica.
- [ ] **Validation:** Todos los tests (`[COMANDO_TEST_FRONTEND]`) y linter (`[COMANDO_LINTER_FRONTEND]`) pasan.
- [ ] **✅ Fase 3 Aprobada por el Equipo**

---

## Phase 4: Final Integration & Review

- [ ] **E2E Testing:** Añadir tests de Playwright para los flujos críticos.
- [ ] **Actualizar Documentación Central:**
    - [ ] `docs/architecture/diagrams/C2-containers.md` (si hay cambios relevantes).
    - [ ] `agents/project/context/bounded-contexts.yaml`.
- [ ] **Revisión Humana:** Marcar la tarea como "Pendiente de Revisión".
- [ ] **Cierre:** Marcar la tarea como "Completada" y actualizar `sprint-registry.yaml` tras la aprobación.
- [ ] **✅ Fase 4 Aprobada por el Equipo**

---
## Documentación Específica del Módulo

*Esta sección contiene los enlaces a los documentos de diseño y especificación creados en la Fase 1.*

- **[Especificación del Módulo](module-spec.md)**
- **[Modelo de Dominio](domain-model.md)**
- **[Diagramas de Dominio](diagrams/domain-model.md)**
- **[Casos de Uso](use-cases.md)**
- **[Contratos de API](api-contracts.md)**
- **[Guía de Implementación](implementation-guide.md)**


