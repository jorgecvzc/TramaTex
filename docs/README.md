# Estándar de documentación del proyecto

## 1. Objetivos

- Centralizar toda la documentación en `/docs/`.
- Mantener una estructura predecible y fácil de navegar.
- Asegurar que la documentación esté actualizada, versionada y trazable.

## 2. Regla de Oro: Nombres y Contenido

- **Nombres de Archivos/Carpetas:** ✅ **Inglés** (`kebab-case`).
- **Contenido de la Documentación:** ✅ **Castellano**.

**Ejemplo:**
- `docs/architecture/adr/001-tech-stack.md` (Correcto)
- `docs/arquitectura/adrs/001-stack-tecnologico.md` (Incorrecto)

Esta regla es fundamental para mantener la consistencia del proyecto.

## 3. Estructura de carpetas

```text
/docs
  ├─ README.md              # Este estándar + índice general
  ├─ overview/              # Contexto y visión
  ├─ architecture/          # Arquitectura del sistema
  ├─ guides/                # Guías de uso (usuarios y devs)
  ├─ reference/             # Información de referencia estable (API, modelos)
  └─ log/                   # Registro de trabajo (sprints, estado)
```

## 4. Reglas de estilo

- Formato por defecto: Markdown (.md).
- Un documento por tema: evitar documentos “monstruo”.
- Enlaces relativos siempre que sea posible (ej.: ../architecture/vision.md).
- Idioma: usar siempre el mismo (por defecto, español), salvo APIs o términos técnicos.

## 5. Flujo de mantenimiento

- Toda feature relevante debe ir acompañada de actualización de documentación.
- La documentación se revisa en las PR igual que el código.
- Al final de cada release se revisa /docs/ para:
    - Marcar como obsoleto lo que ya no aplique.
    - Dividir documentos demasiado grandes.
    - Ajustar la estructura si crecen nuevos módulos.