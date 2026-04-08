# Estándares de Código y Estilo del Proyecto TramaTex

Este documento centraliza las reglas y convenciones para el código y el estilo del proyecto TramaTex. Es la fuente de verdad para el formato, la calidad, las convenciones de nomenclatura y los estándares de desarrollo para todos los módulos y componentes, tanto de backend como de frontend.

---

## 1. Convenciones de Lenguaje General

### Documentación
- **Idioma:** Castellano (Español)
- **Alcance:** Toda la documentación en `/docs/`
- **Incluye:** Requisitos de negocio, especificaciones técnicas, decisiones de arquitectura (ADRs), documentación de módulos, registro de trabajo (sprints y tareas), guías de usuario.
- **Razón:** Documentación en idioma nativo (español) para stakeholders y usuarios finales.

### Código y Comentarios
- **Idioma:** Inglés
- **Alcance:** Código fuente, identificadores y nombres del sistema de archivos.
- **Incluye:** Nombres de archivos y directorios, nombres de variables, funciones, clases, comentarios de código explicando el "porqué" (no el "qué"), cadenas de documentación de API, mensajes de error, mensajes de commit (preferiblemente en inglés).
- **Razón:** Alinea con los estándares de código y documentación técnica para accesibilidad global y compatibilidad con herramientas.

### Interfaz de Usuario (UI)
- **Idioma:** Castellano (Español)
- **Alcance:** Toda la interfaz de usuario visible por el usuario final.
- **Incluye:** Etiquetas, estados, prioridades, mensajes, tooltips, filtros de selección, títulos de páginas y secciones.
- **Regla:** Los estados y prioridades se almacenan en inglés en la API y base de datos como estándar técnico, pero **siempre se muestran en castellano** en la UI. Cada módulo frontend mantiene un mapa de traducción centralizado (ej. `salesApi.getStatusLabel()`, `mesStatusLabel()`).
- **Razón:** Interfaz nativa en español para usuarios finales; alineado con la documentación.

---

## 2. Convenciones de Nomenclatura de Archivos y Directorios

### Regla General
- **kebab-case universal:** Todos los archivos y directorios usan `kebab-case` (minúsculas con guiones).

### Estándar
- **Descripción:** Todos los archivos usan `kebab-case` (minúsculas con guiones).
- **Ejemplos Correctos:**
  - `guia-de-uso.md`
  - `project-context.yaml`
  - `load-session.yaml`
  - `adr-template.md`
  - `bounded-contexts.yaml`
- **Ejemplos Incorrectos:**
  - `GUIA_DE_USO.md`
  - `Project_Context.yaml`
  - `LoadSession.yaml`
  - `ADR_TEMPLATE.md`
  - `Bounded_Contexts.yaml`

### Excepciones
- **Archivos estándar universales en la raíz del proyecto:** `README.md`, `LICENSE`, `CONTRIBUTING.md`, `CHANGELOG.md`. (Reconocidos universalmente en proyectos de código abierto).
- **Archivos de plantilla con prefijo de guion bajo:** `_kebab-case-template.md`. Se utiliza el prefijo `_` para distinguirlos visualmente como archivos meta/herramientas, manteniéndose el `kebab-case`.

### Razones para `kebab-case`
- Consistencia: Una única regla sencilla de recordar.
- Legibilidad: Más fácil de leer que `SCREAMING_SNAKE_CASE`.
- Multiplataforma: No hay problemas de sensibilidad a mayúsculas/minúsculas.
- Estándar: Usado por la mayoría de proyectos de código abierto.
- URLs: Funciona bien en contextos web sin codificación.

---

## 3. Estructura de Directorios Estándar

### Niveles Raíz
- `/docs/`: Toda la documentación del proyecto (en español), orientada al usuario.
- `/agents/`: Instrucciones y prompts para agentes de IA (fuera de `/docs/`), configuración interna de IA.
- `/apps/`: Todos los componentes de la aplicación (servicios, clientes, etc.), orientados al desarrollador.

### Estructura de `/docs/`
- `docs/architecture/`: Documentación técnica central (arquitectura, ADRs, diseño de módulos).
- `docs/modules/`: Documentación detallada por cada módulo o Bounded Context.
- `docs/guides/`: Guías prácticas y tutoriales (para desarrolladores y usuarios).
- `docs/log/`: Registro histórico del trabajo (sprints, hitos).

### Política de Directorio Raíz Limpio (Regla de Oro)
- El directorio raíz del proyecto debe permanecer lo más limpio posible.
- **Archivos permitidos:** `README.md`, `AGENTS.md`, `CONTRIBUTING.md`, `LICENSE`, `SECURITY.md`, `CHANGELOG.md`, y archivos de configuración esenciales (`.gitignore`, `Makefile`, `package.json`, etc.).
- **Archivos NO permitidos:** Cualquier otro archivo `.md` (debe ir en `/docs`).

---

## 4. Estándares de Calidad de Código Universal

### Principios Generales
- El código sin tests en el dominio crítico será **RECHAZADO**.
- La lógica de negocio fuera de la capa de dominio será **RECHAZADA**.
- Los valores de configuración hardcodeados serán **RECHAZADOS**.
- Las funciones gigantes sin documentación serán **RECHAZADAS**.
- El manejo de errores como cadenas de texto será **RECHAZADO**.
- Los commits sin mensajes descriptivos serán **RECHAZADOS**.

### Filosofía de Testing
- **Test-Driven Development (TDD):** Para el dominio crítico.
- **Cobertura Mínima:** ≥85% promedio, ≥90% para rutas críticas.
- Los tests deben ser claros y mantenibles.
- Tests de integración para casos de uso.
- Tests unitarios para lógica de dominio.

### Checklist de Revisión de Código (Extracto)
- ¿Tiene el código tests apropiados?
- ¿Pasan los tests localmente?
- ¿Sigue la arquitectura del proyecto?
- ¿Son los commits descriptivos?
- ¿Los comentarios explican el "porqué" y no el "qué"?
- ¿Es el manejo de errores explícito y tipado?

---

## 5. Estándares Específicos de Backend (Go)

### Criterios Aceptados (Extracto)
- **Capa de Dominio:** Sin dependencias externas (frameworks, ORM). Reglas de negocio encapsuladas. Interfaces para persistencia. Errores tipados.
- **Casos de Uso (Aplicación):** Orquestan dominio e infraestructura. Sin lógica de negocio. Inyección de dependencias clara. Transacciones en la capa de aplicación.
- **Infraestructura:** Implementa interfaces de dominio. Contiene código específico del framework (GORM). Sin lógica de negocio.
- **Interfaces (HTTP):** Handlers como adaptadores delgados. DTOs para contratos. Validación de entrada (superficial). Traducción de errores a códigos de estado HTTP.
- **Testing:** Tests unitarios para dominio, integración para casos de uso críticos.
- **Calidad de Código:** Cero warnings de linter. Formato aplicado. Sin valores hardcodeados/credenciales. Comentarios solo del "porqué". Manejo de errores adecuado.
- **Documentación:** Interfaces documentadas con contratos. Algoritmos complejos explicados. Funciones públicas con `godoc`.

### Criterios Rechazados (Extracto)
- Lógica de negocio fuera de la capa de dominio.
- Queries de ORM en aplicación o interfaces.
- DTOs definidos en el paquete de dominio.
- Tests sin fixtures o setup.
- Configuración hardcodeada o números mágicos.
- Funciones sin tests en dominio crítico.
- Manejo de errores basado en cadenas de texto.
- Comentarios obvios/redundantes.
- Funciones que exceden las 50 líneas sin subsecciones claras.
- Commits sin mensajes descriptivos.
- Problemas de linting o formato.

---

## 6. Estándares Específicos de Frontend (Vue.js 3)

### Criterios Aceptados (Extracto)
- **Uso del Framework:** Exclusivamente Composition API. Stores de Pinia para estado global. Composables para lógica reutilizable. Capa de servicios para llamadas API.
- **Diseño de Componentes:** Pequeños y enfocados (máx. 100-150 líneas). Contratos claros de props y emits. Separación de responsabilidades.
- **Sistema de Diseño (CSS):**
  - **Uso de Vanilla CSS:** El proyecto NO usa Tailwind CSS. Se utiliza CSS nativo con variables globales (`_variables.css`).
  - **Arquitectura Modular:** Los estilos se dividen en archivos de sistema (`_base.css`, `_buttons.css`, `_forms.css`, `_sections.css`).
  - **Identidad Visual:** Uso obligatorio de la tipografía corporativa y la paleta de colores definida en las variables CSS.
- **Estándares de Formularios (Crítico):**
  - **Inputs:** Deben usar siempre la clase `.form-input` (definida en `_forms.css`).
  - **Selects:** Deben usar `.form-input` o componentes especializados como `PartySelector.vue`.
  - **Textareas:** Deben usar la clase `.form-textarea`.
  - **Estados:** Se deben implementar estilos consistentes para `:focus` (borde azul oscuro/oro y sombra suave), `:disabled` y estados de error.
  - **Estructura:** Uso de `FormSection` para agrupar campos y `DataRow` para visualización mixta (lectura/edición).
- **Iconografía:** Exclusivamente `Material Symbols Outlined`. Los emojis están **PROHIBIDOS** en la interfaz final.
- **Manejo de Formularios:** Validación en composables o stores. Lógica de validación reutilizable. Mensajes de error claros al usuario.
- **Gestión de Estado:** Estado global en stores de Pinia. Acciones de store para mutaciones. Propiedades computadas para estado derivado.
- **Testing:** Tests para lógica de negocio crítica en stores y composables. Vitest configurado. ≥80% de cobertura para componentes y composables críticos.
- **Calidad de Código:** TypeScript en modo estricto. ESLint sin warnings.

### Criterios Rechazados (Extracto)
- **Tailwind CSS:** El uso de Tailwind está estrictamente **prohibido** para mantener la pureza del sistema de diseño propio.
- Options API o Vuex.
- Estilos inline (salvo casos dinámicos excepcionales).
- Llamadas directas a API en componentes (usar servicios).
- Validación de datos a nivel de plantilla exclusivamente.
- Componentes gigantes sin sub-componentes.
- Tipos de TypeScript faltantes (`any`).
- Uso de emojis como iconos.

---

## 7. Reglas Cross-Capa

### Filosofía de Testing
- **Principio:** Testear en la capa que tenga sentido.
- **Tests de Dominio:** Unitarios y aislados. Sin mocks. Prueban reglas de negocio, no infraestructura.
- **Tests de Aplicación:** De integración, ejercitan casos de uso. Mocks para dependencias externas (repositorios). Prueban orquestación de workflows.
- **Tests de Infraestructura:** Prueban implementaciones de repositorios y adaptadores de servicios externos. Usan bases de datos de test o contenedores.
- **Tests de Interfaces:** Prueban contratos HTTP (request/response), traducción de errores y middleware.
- **Tests de UI:** Prueban mutaciones de store, lógica de composables y componentes complejos. El comportamiento de UI **NO** se testea en tests de backend.

### Cálculos de Negocio: Solo en Backend
- **Principio:** Toda lógica de cálculo de negocio (monetaria, fiscal, descuentos, precios, márgenes, subtotales, totales) se ejecuta **exclusivamente en el backend**.
- **El frontend NO debe:** Replicar fórmulas de negocio, redondear valores monetarios, calcular descuentos/impuestos, aplicar porcentajes de tipo impositivo, ni derivar cantidades a partir de reglas de negocio.
- **El frontend SÍ puede:** Mostrar valores pre-calculados por la API, formatear para presentación (`toFixed`, `Intl.NumberFormat`) sin alterar el valor subyacente, y llamar a endpoints de previsualización (`/preview`) para obtener cálculos en tiempo real durante la edición.
- **Justificación:** La duplicación de lógica de cálculo entre frontend y backend genera discrepancias de redondeo, rompe el principio de fuente única de verdad, y dificulta el mantenimiento. El backend es la autoridad para cualquier valor derivado de reglas de negocio.

### Manejo de Errores
- **Principio:** Los errores fluyen hacia arriba, tipados en todas las capas.
- **Dominio:** Define tipos de error personalizados, con código y mensaje.
- **Aplicación:** Captura errores de dominio, añade contexto (caso de uso, usuario, timestamp). No suprime errores.
- **Infraestructura:** Captura errores técnicos (DB, red). Los traduce a errores de dominio si son relevantes para el negocio. Loguea errores técnicos inesperados.
- **Interfaces:** La capa HTTP traduce errores a códigos de estado. Retorna mensajes de error significativos al cliente. Loguea lo necesario.

---

## 8. Estándares de Commits Git

### Formato de Mensaje de Commit
- **Estructura:** `[tipo]([alcance]): [asunto] - [cuerpo]`

- **Tipo:**
  - `feat`: Nueva funcionalidad.
  - `fix`: Corrección de un bug.
  - `refactor`: Reestructuración de código (sin cambios de lógica).
  - `docs`: Cambios en la documentación.
  - `test`: Añadido de tests.
  - `style`: Formato, linting.

- **Alcance (opcional):** Contexto delimitado (auth, party, product, pricing, sales, mes), capa (domain, app, infra, interfaces), infraestructura (docker, config, ci-cd).

- **Asunto:**
  - Modo imperativo ("add", no "added").
  - Menos de 50 caracteres.
  - Sin punto al final.

- **Cuerpo:** Explica el QUÉ y el PORQUÉ del cambio. Referencia a issues o ADRs relacionados. Menciona breaking changes si los hay.

### Ejemplos
- **Bueno:**
  - `feat(pricing): add volume discount calculation - Implements tiered discount rules based on order quantity`
  - `fix(party): resolve organization hierarchy loop - Add cycle detection in parent-child relationships`
  - `refactor(domain): extract money operations to value object - Improves reusability and testability`
- **Malo:**
  - `update code`
  - `fix bug`
  - `changes`
  - `wip`

### Granularidad del Commit
- **Principio:** Un cambio lógico por commit.
- **Hacer commit de:** Una funcionalidad completa y testeada, una corrección de bug con tests, refactoring de un componente.
- **NO hacer commit de:** Mezcla de funcionalidades no relacionadas, implementaciones parciales, código con tests fallidos, documentación sin commit.

---

## 9. Lista de Verificación Pre-Commit

### Backend (Go)
- `go test ./...` ← Todos los tests pasan.
- `golangci-lint run ./...` ← Cero warnings.
- `go fmt ./...` ← Código formateado.
- `go vet ./...` ← Pasa.
- El mensaje de commit es descriptivo.

### Frontend (Vue)
- `npm run lint` ← Pasa.
- `npm run test` ← Pasan.
- `npm run build` ← Construcción exitosa.
- `npm run type-check` ← Pasa.
- El mensaje de commit es descriptivo.

### Siempre
- No hay secretos o credenciales hardcodeados.
- No hay logs de depuración restantes.
- No hay código comentado.
- La documentación está actualizada.
- El archivo de sesión está actualizado.

---

## 10. Enfoque de Revisión de Código

### Arquitectura
- ¿Sigue Clean Architecture?
- ¿Está la lógica de negocio en el dominio?
- ¿Se respetan los límites de las capas?
- ¿Son correctas las dependencias (sin dependencias ascendentes)?

### Testing
- ¿Se han testeado las rutas críticas?
- ¿La cobertura cumple los objetivos?
- ¿Son los tests claros y mantenibles?
- ¿Los tests validan el comportamiento, no la implementación?

### Calidad de Código
- ¿Es el código legible y los nombres descriptivos?
- ¿Se explican las secciones complejas?
- ¿Es adecuado el manejo de errores?
- ¿Hay mejoras obvias?

### Específico del Dominio
- ¿Respeta las reglas de negocio?
- ¿Se manejan los casos excepcionales?
- ¿Se utiliza el lenguaje ubicuo?
- ¿Podría verse afectado otro módulo?
