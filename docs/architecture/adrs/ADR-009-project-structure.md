# ADR-009 – Estructura de Carpetas y Organización del Proyecto

**Fecha:** 11/01/2026  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, Claude (Anthropic)

---

## 1. Contexto

TramaTex requiere una estructura de proyecto que:

- Refleje los principios de **Clean Architecture** (ADR-002)
- Soporte un **monolito modular** con dominios claramente separados (ADR-003)
- Facilite el **desarrollo dirigido por dominio** (ADR-006, ADR-007)
- Sea **mantenible por un desarrollador en solitario** con IA copiloto
- Permita **escalabilidad futura** hacia microservicios si es necesario
- Organice **documentación técnica y de sesiones** de forma eficiente

**Restricciones técnicas:**

- Stack: Go (tramatex-api), Vue.js 3 (frontend), PostgreSQL (persistencia), Docker Compose
- Desarrollo local-first
- Un solo repositorio Git (monorepo)
- Sin herramientas complejas de monorepo (Nx, Turborepo, etc.)

**Riesgos si no se define estructura:**

- Mezcla de capas (dominio con infraestructura)
- Dificultad para aplicar un desarrollo guiado por pruebas estricto.
- Pérdida de trazabilidad entre módulos
- Código difícil de navegar y mantener
- Imposibilidad de extraer módulos a servicios independientes en el futuro

---

## 2. Alternativas Consideradas

### Alternativa A – Estructura plana por tipo de archivo

**Ventajas:**

- Simple y tradicional
- Fácil de entender para principiantes
- Menos niveles de anidamiento

**Desventajas:**

- Mezcla todos los dominios
- Dificulta aplicar Clean Architecture
- No soporta modularidad
- Imposible extraer módulos
- No escala con la complejidad

### Alternativa B – Estructura por capas técnicas

**Ventajas:**

- Respeta capas de Clean Architecture
- Separación técnica clara

**Desventajas:**

- No refleja dominios de negocio
- Dificulta navegación por módulo
- Mezcla Party, Producto, MES en mismas carpetas
- No facilita extracción de módulos

### Alternativa C – Estructura por módulos de dominio (Clean Architecture modular)

**Ventajas:**

- Refleja dominios de negocio claramente
- Facilita navegación por módulo completo
- Soporta Clean Architecture
- Permite extraer módulos a servicios
- Agrupa código relacionado
- Facilita desarrollo incremental por módulo

**Desventajas:**

- Mayor profundidad de carpetas
- Requiere disciplina para mantener separación

### Alternativa D – Monorepo con workspaces independientes

**Ventajas:**

- Preparado para microservicios desde día 1
- Reutilización clara de código compartido

**Desventajas:**

- Complejidad prematura para MVP
- Requiere herramientas adicionales
- Overhead de gestión de versiones entre paquetes
- No se necesita para monolito modular

---

## 3. Decisión Adoptada

Se adopta **Alternativa C: Estructura por módulos de dominio con Clean Architecture**.

**Justificación:**

- **Refleja el dominio de negocio:** Party, Producto, Tarificación, Ventas, MES son módulos visibles
- **Soporta Clean Architecture:** Capas separadas dentro de cada módulo
- **Facilita el desarrollo guiado por pruebas:** Tests están junto al código que prueban
- **Escalabilidad controlada:** Módulos pueden extraerse a servicios sin refactoring masivo
- **Navegación intuitiva:** Toda la lógica de Party está en carpetas `*/party/`
- **Desarrollo incremental:** Cada fase del ADR-007 tiene carpetas claras donde trabajar
- **Mantenibilidad:** Un desarrollador puede entender la estructura en una sesión
- **Documentación alineada:** La estructura de `/docs/modules/` refleja la del código

---

## 4. Consecuencias

### Positivas

- **Claridad de dominio:** Fácil encontrar todo lo relacionado con un módulo específico
- **Separación de responsabilidades:** Capas de Clean Architecture respetadas
- **Testabilidad:** Tests unitarios junto a entidades, tests de integración junto a casos de uso
- **Documentación trazable:** Estructura de `/docs/modules/` espeja la del código
- **Evolución controlada:** Nuevos módulos se añaden sin afectar existentes
- **Onboarding rápido:** Estructura autodocumentada
- **Preparación para microservicios:** Módulos ya están delimitados

### Negativas

- **Profundidad de carpetas:** Rutas pueden ser largas (ej: `internal/domain/party/customer.go`)
- **Duplicación aparente:** Carpetas `party/` existen en domain, application, infrastructure
- **Disciplina requerida:** Fácil romper separación si no hay code reviews
- **Curva de aprendizaje inicial:** Requiere entender Clean Architecture primero

### Mitigaciones

- **Profundidad:** IDEs modernos (VSCode, GoLand) manejan bien rutas largas
- **Duplicación:** Es intencional, cada capa tiene responsabilidades diferentes
- **Disciplina:** El desarrollo guiado por pruebas ayuda a mantener separación
- **Curva:** Documentación clara de la estructura en README.md

---

## 5. Alcance

Este ADR define:

- Estructura completa de carpetas del proyecto TramaTex
- Organización del monorepo (tramatex-api + frontend + documentación)
- Convenciones de nombres de carpetas y archivos
- Ubicación de tests, migraciones, configuraciones
- Estructura de documentación técnica

**No define:**

- Convenciones de nombres de variables/funciones (estilo de código)
- Estrategia de branching Git
- Proceso de CI/CD específico
- Configuración de herramientas (linters, formatters)

---

## 6. Integración con otros ADRs

- **ADR-001:** Stack tecnológico → Define Go, Vue.js, PostgreSQL (estructura los refleja)
- **ADR-002:** Clean Architecture → Estructura por capas (domain, application, infrastructure, interfaces)
- **ADR-003:** Monolito modular → Módulos separados pero en mismo repositorio
- **ADR-006:** Desarrollo dirigido por dominio → Estructura facilita priorizar por dominio
- **ADR-007:** Orden de implementación → Carpetas se crean incrementalmente por fase
- **ADR-008:** Cronograma → Estructura soporta desarrollo incremental de 24 meses

---

## 7. Notas Adicionales / Consideraciones Especiales

### Estructura Completa del Proyecto
Para una vista detallada de la estructura de carpetas del proyecto, consulte [Guía de Detalles de la Estructura del Proyecto](../../guides/developer/project-structure-details.md).

### Convenciones de Nombres
Para las convenciones de nomenclatura detalladas en el proyecto, consulte [Guía de Convenciones de Nomenclatura](../../guides/developer/naming-conventions.md).

---

### Ubicación de Tests
Para la guía detallada sobre la ubicación de los tests, consulte [Guía de Ubicación de Tests](../../guides/developer/testing-guidelines.md).

---

### Gestión de Configuraciones
Para la gestión detallada de configuraciones, consulte [Guía de Gestión de Configuraciones](../../guides/developer/configuration-management.md).

---

### Comandos Make (Makefile)
Para una lista de comandos `make` comunes, consulte [Comandos Make Comunes](../../guides/developer/makefile-commands.md).

---

### Flujo de Creación de Carpetas por Fase
Para el flujo detallado de creación de carpetas por fase, consulte [Guía de Bootstrapping del Proyecto](../../guides/developer/project-bootstrapping-guide.md).

---

## 8. Referencias

- **Documentos internos:**
    
    - Documento Consolidado TramaTex v3.0
    - ADR-001: Stack Tecnológico
    - ADR-002: Clean Architecture y DDD
    - ADR-003: Monolito Modular
    - ADR-007: Orden de Implementación de Módulos
- **Prácticas externas:**
    
    - Clean Architecture (Robert C. Martin)
    - Domain-Driven Design (Eric Evans)
    - Golang Project Layout: https://github.com/golang-standards/project-layout
    - Vue.js Style Guide: https://vuejs.org/style-guide/
- **Diagramas:**
    
    - Ver `/docs/diagrams/architecture/project-structure.png` (a crear)
    - Ver `/docs/diagrams/architecture/clean-architecture-layers.png` (a crear)

---

**Fin del ADR-009 – Estructura de Carpetas y Organización del Proyecto**
