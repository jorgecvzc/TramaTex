# Estándar de documentación del proyecto

## 1. Objetivos

- Centralizar toda la documentación en `/docs/`.
- Mantener una estructura predecible y fácil de navegar.
- Asegurar que la documentación esté actualizada, versionada y trazable.

## 2. Regla de Oro: Nombres y Contenido

- **Nombres de Archivos/Carpetas:** ✅ **Inglés** (`kebab-case`).
- **Contenido de la Documentación:** ✅ **Castellano**.

**Ejemplo:**
- `docs/architecture/adrs/ADR-001-technology-stack-selection.md` (Correcto)
- `docs/arquitectura/adrs/001-stack-tecnologico.md` (Incorrecto)

Esta regla es fundamental para mantener la consistencia del proyecto.

## 3. Estructura de carpetas

```text
/docs
  ├─ README.md              # Este estándar + índice general
  ├─ architecture/          # Arquitectura del sistema
  ├─ guides/                # Guías de uso (usuarios y devs)
  ├─ modules/               # Documentación detallada por Bounded Context (specs, diagramas)
  └─ log/                   # Registro de trabajo (sprints, estado)
```

## 4. Reglas de estilo

- **Formato de Archivos:** Markdown (`.md`) para toda la documentación.
- **Diagramas:** Se utilizará **Mermaid** para la creación de diagramas. El código de Mermaid se incrustará directamente en los archivos `.md`.
- **Atomicidad:** Un documento por tema para evitar documentos “monstruo”.
- **Enlaces:** Enlaces relativos siempre que sea posible (ej.: `../architecture/architecture-vision.md`).
- **Idioma:** Castellano para el contenido, salvo APIs o términos técnicos.

## 5. Flujo de mantenimiento

- Toda feature relevante debe ir acompañada de actualización de documentación.
- La documentación se revisa en las PR igual que el código.
- Al final de cada release se revisa /docs/ para:
    - Marcar como obsoleto lo que ya no aplique.
    - Dividir documentos demasiado grandes.
    - Ajustar la estructura si crecen nuevos módulos.