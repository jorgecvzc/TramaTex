# Manual del Proyecto TramaTex

**Versión:** 1.0  
**Fecha:** 2026-01-27  
**Propósito:** Servir como punto de partida y guía central para entender la estructura, arquitectura, visión y estado del proyecto TramaTex.

---

## 👋 ¡Bienvenido a TramaTex!

Este documento es el **índice maestro** del proyecto. Si eres nuevo aquí (o necesitas un recordatorio), este es el mejor lugar para empezar. No duplica información, sino que te guía hacia los documentos clave donde reside el conocimiento detallado.

---

### 1. 🎯 Visión y Propósito del Proyecto

TramaTex es un sistema ERP/MES diseñado para ser una solución `local-first` para microempresas del sector textil, enfocándose en la eficiencia, mantenibilidad y un ciclo de vida largo.

Para entender el "porqué" detrás del proyecto, los objetivos de negocio y el alcance completo del MVP, consulta el Project Charter.

- **Lectura Esencial:** [**📄 Project Charter (./01-project-charter.md)**](./01-project-charter.md)

---

### 2. 🏗️ Arquitectura y Decisiones de Diseño

La arquitectura de TramaTex se basa en principios modernos y probados para garantizar la calidad y la escalabilidad. Las decisiones más importantes no se toman a la ligera, sino que se documentan formalmente.

- **Vista General de la Arquitectura:** Para una visión de alto nivel, empieza aquí:
  - [**🗺️ Architecture Overview (../engineering/architecture/architecture-overview.md)**](../engineering/architecture/architecture-overview.md)

- **Decisiones Arquitectónicas Clave (ADRs):** Estas son las reglas fundamentales que gobiernan nuestro proyecto.
  - [**ADR-001**](../engineering/architecture/adr/ADR-001-seleccion-stack-tecnologico.md): Selección del **Stack Tecnológico** (Go, Vue.js, PostgreSQL).
  - [**ADR-002**](../engineering/architecture/adr/ADR-002-adopcion-clean-architecture-ddd.md): Adopción de **Clean Architecture y DDD** con rigor asimétrico.
  - [**ADR-003**](../engineering/architecture/adr/ADR-003-tipo-distribucion-aplicacion.md): Definición como **Monolito Modular**.
  - [**ADR-007**](../engineering/architecture/adr/ADR-007-orden-implementacion-modulos.md): **Orden de Implementación** de Módulos (Fases).
  - [**ADR-010**](../engineering/architecture/adr/ADR-010-estrategia-seguridad-defensa-profundidad.md): Estrategia de **Seguridad y Defensa en Profundidad**.

---

### 3. 💻 Stack Tecnológico

A continuación se presenta un resumen de las tecnologías utilizadas.

| Componente | Tecnología Principal | Detalles Clave |
| :--- | :--- | :--- |
| **Backend API** | Go (Golang) | Framework Gin, GORM para persistencia, JWT para seguridad. |
| **Frontend** | Vue.js 3 | Composition API, Vite como build tool, Pinia para estado. |
| **Estilos (CSS)** | Tailwind CSS | Enfoque "Utility-first", sin CSS tradicional. |
| **Base de Datos** | PostgreSQL | Motor de base de datos relacional y transaccional (ACID). |
| **Infraestructura** | Docker / Docker-Compose | Containerización para desarrollo y despliegue local-first. |

---

### 4. 🛠️ Metodología y Flujo de Trabajo

El trabajo se organiza en "Sprints", que representan fases de desarrollo de varios días centradas en un objetivo estratégico. Cada Sprint se desglosa en "Tareas" técnicas.

- **Guía de Contribución:** Para entender cómo crear ramas, el formato de los commits, el proceso de Pull Request y los estándares de código, consulta la guía de contribución.
  - **Lectura Obligatoria:** [**🤝 Guía de Contribución (../../CONTRIBUTING.md)**](../../CONTRIBUTING.md) *(Nota: Planificada en Sprint 04)*

- **Trabajo en Curso:** Para ver las tareas y sprints activos, dirígete a la carpeta de sprints.
  - [**📂 Directorio de Sprints (./sprints/)**](./sprints/)

---

### 5. 📊 Estado Actual del Proyecto

Para obtener una vista rápida y actualizada del progreso del proyecto, los hitos recientes, los próximos pasos y los bloqueadores, consulta el documento de estado del proyecto.

- **Consulta Frecuente:** [**📊 Estado del Proyecto (./project-status.md)**](./project-status.md)
