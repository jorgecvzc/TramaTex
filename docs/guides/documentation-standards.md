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

## 3. Taxonomía del Portal de Documentación (`/docs/`)

| Carpeta | Contenido Exclusivo |
| :--- | :--- |
| `architecture/` | Visiones globales, ADRs, Glosario Ubicuo y Diagramas C4. |
| `modules/` | Guías de comportamiento, lógica de dominio y contratos de API por módulo. |
| `guides/` | Manuales de usuario, guías de ingeniería y estándares (este documento). |
| `log/` | Trazabilidad: Sprints, tareas finalizadas y bitácoras de sesión. |

---

## 4. Protocolo de Refinamiento (Poda y Síntesis)

Toda pasada de refinamiento debe seguir estos pasos:
1.  **Estudio:** Analizar la desalineación entre el documento y la realidad del sistema.
2.  **Poda:** Eliminar listas de campos técnicos y detalles de implementación de bajo nivel.
3.  **Síntesis:** Redactar la lógica de comportamiento y reglas de negocio de forma narrativa.
4.  **Validación:** Comprobar la integridad de los enlaces y la ausencia de archivos huérfanos.

---
**Versión del Plan:** 2.0 (Consolidada)
**Última Actualización:** 2026-03-07
