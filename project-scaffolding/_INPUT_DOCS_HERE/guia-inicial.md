# Guía de Iniciación para la Creación de un Nuevo Proyecto — API MOSY

> Esta guía está destinada a usuarios humanos y asistentes inteligentes (IA/LLMs) responsables de generar scaffolding, estructura de carpetas, configuración y documentación de nuevos proyectos conforme al estándar API MOSY.

---

## 1. Propósito

La documentación y estructura de `api-mosy-docs` proporcionan una base integral y modular para generar proyectos backend robustos usando tecnologías .NET, Azure, y prácticas alineadas con la cultura y operativa de Record Go. El objetivo es estandarizar la creación, configuración y evolución de nuevos proyectos, garantizando cohesión entre equipos y automatización avanzada.

---

## 2. Punto de partida y navegación

- El archivo principal de referencia es `README.md`, que da acceso rápido a todo el mapa documental y explica el enfoque modular.
- El archivo `index.md` actúa como tabla de navegación central, enlazando todos los documentos especializados.
- Cada archivo/documento dentro de `api-mosy-docs` es autocontenible: introduce su propósito explícitamente (para humanos y para LLMs) y enlaza a otras secciones relevantes.

**LLM/IA**: inicia siempre leyendo `index.md` y navega las secciones usando los bloques “LLM USE” de cada archivo para adaptar tus prompts y acciones.
**Humano**: utiliza el índice o los títulos para acceder directamente al área en la que necesitas crear, extender o consultar (arquitectura, pipelines, APIs, equipo, etc.).

---

## 3. Flujo recomendado para la creación del proyecto

### A. Contextualización

- Lee `contexto-negocio.md` y `alcance-inicial.md` para entender el porqué, los límites y los destinatarios del proyecto.
- Consulta `guia-tecnologica.md` para visualizar el stack tecnológico, dependencias base y cómo se integran satélites heredados y legacy.

### B. Definición eficiente de arquitectura y tecnología

- Entra en `/arquitectura/README.md` para el árbol de decisiones clave.
    - Consulta `adr.md` para ver las decisiones técnicas/fundamentales tomadas y reutilizar postulados para justificar tus propias acciones/scaffolding.
    - Lee `estructura-bbdd.md` para formatear el área de persistencia/multi-país y la lógica de segregación de datos.
    - Si tienes integraciones externas, revisa `integraciones.md`.

### C. Organización y DevOps

- Estudia `operativa-devops.md` para definir la estructura de ramas git, flujos de pipeline, testing, deploy, naming y control de versiones.
- Aplica todos los patrones y convenciones descritos para asegurar que tu proyecto y sus scripts se integran sin fricciones en el ecosistema Record Go y Azure.

### D. Ejemplos y patrones contractuales

- Copia/extrae de `apis-ejemplo.md` las rutas, contratos y ejemplos de request/response para empezar controladores, contractos OpenAPI/Swagger, tests, etc.
- Usa estos ejemplos como referencia de naming, estructura y seguridad esperada.

### E. Instrumentación y logging

- Implementa logging según las recomendaciones de `logging.md` (nivel de log, formato struct alterno, integración con Azure Monitor, etc.).

### F. Equipo, colaboración y conocimiento

- Consulta o adapta `colaboracion.md` para decidir la dinámica de equipo, apoi de pull requests, estrategia de documentación viva y contribución transparente.
- Establece y comunica roles o puntos de contacto relevantes para el proyecto.

---

## 4. Convenciones y buenas prácticas

- Mantén los bloques “LLM USE” y encabezados claros en cada nuevo documento, fichero o módulo que crees.
- Si tomas cualquier decisión técnica relevante, deja constancia en `arquitectura/adr.md` (Architecture Decision Record).
- Sigue el naming y las rutas recomendadas para facilitar la navegación automática y la trazabilidad.

---

## 5. Arranque automatizado (instrucciones para IA/LLM)

- Realiza un parseo de `index.md`.
- Segmenta los prompts/acciones según la categoría encontrada en los encabezados.
- Usa los ejemplos, patrones y convenciones directamente como plantillas para la generación automatizada.
- Respetar los índices y enlaces internos para mantener la navegabilidad y el futuro mantenimiento.

---

## 6. Checklist rápido antes de iniciar nuevo proyecto

- [ ] Comprendida la justificación y alcance inicial.
- [ ] Revisadas/seleccionadas tecnologías conforme guía base.
- [ ] Adaptada la arquitectura a las líneas marcadas (o se ha registrado una nueva ADR si hay desviación).
- [ ] Aplicados los estándares de DevOps y naming de ramas.
- [ ] Copiados ejemplos API de contratos tipo.
- [ ] Instrumentado el logging básico.
- [ ] Comunicado el modelo de colaboración y registro de responsables.
- [ ] Documentación generada es modular, limpia y con los bloques de uso LLM/humano activos.

---

## 7. Nota final

Este conjunto documental y guía debe evolucionar de forma viva. Actualiza, enriquece y referencia los documentos base siempre que amplíes el proyecto, el stack o la organización. Así garantizas que tanto humanos como IAs tengan siempre la mejor fuente de verdad para arrancar y escalar un proyecto bajo el estándar API MOSY.
