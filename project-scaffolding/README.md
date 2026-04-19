# 🏗️ Ecosistema de Scaffolding y Metodología IA

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 2.0 |
| **Estado** | ✅ Vigente |
| **Propósito** | Motor de materialización y estandarización de proyectos |

---

## 🎯 Propósito
Este proyecto, desarrollado en paralelo a TramaTex, constituye la **infraestructura del pensamiento** del ecosistema. No es una simple colección de plantillas, sino un motor de orquestación basado en IA que permite materializar la estructura completa de un proyecto profesional (documentación, agentes y procesos) a partir de una visión inicial.

TramaTex es el primer sistema de grado industrial nacido y escalado íntegramente bajo esta metodología.

---

## 🐹 Componentes del Ecosistema

### 1. El Motor de Bootstrap (`/agents/bootstrap_workflow/`)
Flujos de trabajo automatizados en YAML que guían a la IA para:
*   Procesar requisitos y definir contextos acotados (Bounded Contexts).
*   Generar el árbol de directorios y archivos base de arquitectura.
*   Poblar automáticamente marcadores de posición (*placeholders*) para mantener la coherencia técnica.

### 2. El Cerebro de Agentes (`/templates/agents/`)
Define el comportamiento de los asistentes de IA que gestionan el proyecto:
*   **`init.yaml`**: Punto de entrada al flujo de trabajo.
*   **`load-session.yaml`**: Inteligencia para retomar tareas y cargar contextos.
*   **`end-session.yaml`**: Protocolo de cierre y persistencia de memoria.

### 3. Plantillas de Ingeniería (`/templates/docs/`)
Estándares de oro para la creación de documentos:
*   Registros de Decisión de Arquitectura (ADRs).
*   Especificaciones funcionales de módulos.
*   Estrategias de prueba y gobernanza.

---

## 📜 Metodología de Bitácora (Session Logging)
El sistema de Scaffolding inyecta en el ADN del proyecto la cultura de la **Trazabilidad Continua**. 

A través del archivo `session-log.md`, los agentes registran cada jornada de desarrollo, las decisiones tomadas y los impedimentos encontrados. Esto permite:
1.  **Continuidad Asíncrona**: Retomar tareas complejas semanas después sin pérdida de contexto.
2.  **Auditoría de Decisiones**: Entender el *porqué* de cada cambio en el código.
3.  **Calidad Certificada**: Asegurar que cada paso cumple con los estándares definidos en la Fase de Bootstrap.

---

## 🛠️ Guías de Operativa
*   **[Guía de Desarrollo](./guides/development-guide.md):** Cómo modificar y extender el propio motor de scaffolding.
*   **[Guía de Placeholders](./guides/placeholders-guide.md):** Manual técnico sobre la inyección de variables en el proceso de generación.

---
[Volver al README Principal de TramaTex](../README.md)
