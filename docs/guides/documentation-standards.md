# Estándares de Documentación del Proyecto

Este documento detalla los estándares y reglas a seguir para la creación y mantenimiento de la documentación en el proyecto TramaTex.

## 1. Regla de Oro: Nombres y Contenido

- **Nombres de Archivos/Carpetas:** ✅ **Inglés** (`kebab-case`).
- **Contenido de la Documentación:** ✅ **Castellano**.

**Ejemplo:**
- `docs/architecture/adrs/adr-001-technology-stack-selection.md` (Correcto)
- `docs/arquitectura/adrs/001-stack-tecnologico.md` (Incorrecto)

Esta regla es fundamental para mantener la consistencia del proyecto.

## 2. Estructura de carpetas

```text
/docs
  ├─ README.md              # Este estándar + índice general
  ├─ architecture/          # Arquitectura del sistema
  ├─ guides/                # Guías de uso (usuarios y devs)
  ├─ modules/               # Documentación detallada por Bounded Context (specs, diagramas)
  └─ log/                   # Registro de trabajo (sprints, estado)
```

### 2.1 Convenciones de nomenclatura de archivos

- **Prefijos de ordenación:** Se permiten prefijos numéricos (`01-`, `02-`) o de fecha (`YYYY-MM-DD-`) para establecer orden y flujos de información.
  - Ejemplo: `01-setup-guide.md`, `2026-01-24-final-audit-report.md`
- **Sufijo `-obsolete`:** Excepción permitida para marcar archivos inactivos. Facilita identificar documentación obsoleta y su posterior limpieza.
  - Ejemplo: `adr-pricing-domain-definition-obsolete.md`
- **Sufijos de versionado:** No permitidos (`-v1`, `-v2`, `-final`). El historial de versiones lo gestiona Git.

## 3. Reglas de estilo

- **Formato de Archivos:** Markdown (`.md`) para toda la documentación.
- **Diagramas:** Se utilizará **Mermaid** para la creación de diagramas. El código de Mermaid se incrustará directamente en los archivos `.md`.
- **Atomicidad:** Un documento por tema para evitar documentos "monstruo".
- **Enlaces:** Enlaces relativos siempre que sea posible (ej.: `../architecture/architecture-vision.md`).
- **Idioma:** Castellano para el contenido, salvo APIs o términos técnicos.
- **Formato de fechas:**
  - En contenido: `DD-MM-YYYY` (ej.: 24-01-2026)
  - En nombres de archivo: `YYYY-MM-DD` (ej.: `2026-01-24-audit-report.md`)

## 4. Flujo de mantenimiento

- Toda feature relevante debe ir acompañada de actualización de documentación.
- La documentación se revisa en las PR igual que el código.
- Al final de cada release se revisa /docs/ para:
    - Marcar como obsoleto lo que ya no aplique.
    - Dividir documentos demasiado grandes.
    - Ajustar la estructura si crecen nuevos módulos.
