# Plan Maestro de Gestión del Conocimiento y Documentación

Este documento constituye la directiva superior para la estructura, contenido y mantenimiento de la información en TramaTex. Su cumplimiento es obligatorio para todos los agentes y desarrolladores.

---

## 1. Principios de Contenido (Filosofía "Behavior-First")

Para evitar la redundancia técnica y maximizar el valor estratégico, se aplican los siguientes mandatos:

1.  **Narrativa de Negocio:** La documentación debe explicar la "personalidad" y el "propósito" de los componentes. El *qué* y el *porqué* prima sobre el *cómo* técnico.
2.  **Anti-Redundancia:** Queda prohibido listar campos, tipos de datos o estructuras que sean autoexplicativos en el código fuente. La documentación termina donde empieza la legibilidad del código.
3.  **Soberanía del Dominio:** El conocimiento se organiza por Bounded Contexts. Cada módulo es responsable de su propia verdad conceptual.

---

## 2. Política de Ámbito y Raíz Limpio (Mandato D)

**RESTRICCIÓN CRÍTICA:** Solo se permite la modificación de documentos dentro de `/docs/` y los archivos expresamente exceptuados en la raíz. El resto del repositorio (código, configuración funcional, agentes) está fuera del ámbito de las tareas de refinamiento documental.

### 2.1. Excepciones Permitidas en el Raíz
Siguiendo los estándares de la industria (GitHub Community Profile) y las necesidades operativas de TramaTex, solo estos archivos pueden residir en el directorio raíz:

**A. Documentación de Comunidad y Onboarding (Markdown):**
- `README.md`: Portal de entrada para humanos.
- `AGENTS.md`: Portal de entrada para IAs (instrucciones maestras).
- `LICENSE` / `LICENSE.md`: Términos legales y propiedad intelectual.
- `CONTRIBUTING.md`: Guía de colaboración y estándares de desarrollo.
- `CHANGELOG.md`: Historial de versiones y cambios significativos.
- `SECURITY.md`: Políticas de seguridad y reporte de vulnerabilidades.
- `CODE_OF_CONDUCT.md`: Normas de comportamiento de la comunidad.
- `SUPPORT.md`: Canales de soporte y ayuda.

**B. Manifiestos y Configuración de Ingeniería (No Markdown):**
- Construcción y Orquestación: `Makefile`, `Dockerfile`, `docker-compose.yml`.
- Dependencias: `go.mod`, `go.sum`, `package.json`, `package-lock.json`.
- Entorno: `.gitignore`, `.gitattributes`, `.editorconfig`, `.env.example`.
- Automatización: Scripts de utilidad (`.sh`, `.ps1`) debidamente documentados.

**Acción Correctiva:** Cualquier otro archivo `.md` detectado en el raíz debe ser movido a la subcarpeta correspondiente en `/docs/` durante las fases de refinamiento.

---

## 3. Arquitectura del Conocimiento (El Árbol TramaTex)

La información se organiza siguiendo una estructura jerárquica inspirada en un árbol:

1.  **El Tronco (`README.md` raíz):** Único punto de entrada profesional. Debe proyectar la visión del TFM, el stack tecnológico real, los pasos de ejecución (Scripts) y el mapa funcional del sistema.
2.  **Las Ramas (READMEs de Carpeta):** Actúan como "nudos" o índices descriptivos de cada subdirectorio (`architecture/`, `modules/`, `guides/`). Su función es facilitar la navegación y dar contexto a los documentos hijos.
3.  **Las Hojas (Documentos .md):** Contenido técnico y de negocio final. Deben seguir estrictamente la **Plantilla Maestra de Estilo**.

### 3.1. Mandatos de Calidad Obligatorios
-   **Unificación Lingüística:** Todo el contenido debe redactarse exclusivamente en **castellano profesional**.
-   **Plantilla Maestra:** Todo documento debe incluir la cabecera de metadatos (Título con emoji, Versión, Estado, Propósito).
-   **Navegabilidad (Cero Huérfanos):** Está prohibida la existencia de archivos `.md` que no sean alcanzables mediante enlaces directos desde el "Tronco" o sus "Ramas".
-   **Sincronía Técnica:** Los modelos de dominio y contratos de API deben reflejar la implementación real del código (Ej: uso de Decimal para finanzas).

## 4. Taxonomía del Árbol (`/docs/`)

| Rama | Contenido Exclusivo |
| :--- | :--- |
| `architecture/` | Decisiones estratégicas (ADRs), Glosario y Diagramas C4. |
| `modules/` | Especificaciones por Bounded Context (Dominios y APIs). |
| `guides/` | Estándares (este documento), Manuales de Usuario e Índice de Scripts. |
| `log/` | Trazabilidad histórica y registro de sesiones. |

## 5. Protocolo de Refinamiento (Poda y Síntesis)

Toda pasada de refinamiento debe seguir estos pasos:
1.  **Estudio:** Analizar la desalineación entre el documento y la realidad del sistema.
2.  **Poda:** Eliminar listas de campos técnicos y detalles de implementación de bajo nivel.
3.  **Síntesis:** Redactar la lógica de comportamiento y reglas de negocio de forma narrativa.
4.  **Validación:** Comprobar la integridad de los enlaces y la ausencia de archivos huérfanos.

---
**Versión del Plan:** 2.0 (Consolidada)
**Última Actualización:** 2026-03-07
