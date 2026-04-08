# Tarea 16-01: TramaTex Design System — Átomos y Estilo Global

**Estado:** ✅ COMPLETADO
**Sprint:** 16
**Fecha Inicio:** 2026-03-26
**Fecha Fin:** 2026-03-28
**Facilitador:** Gemini CLI

---

## 📝 Descripción
Definición y estandarización de los componentes básicos (átomos) de la interfaz de TramaTex para garantizar una experiencia de usuario coherente y profesional. El foco se centra en la iconografía, la paleta de colores y la tipografía técnica.

## 🎯 Objetivos
- [x] Estandarizar la iconografía exclusivamente con **Material Symbols Outlined**.
- [x] Definir 6 tipos semánticos de botones con comportamiento y estética unificados.
- [x] Sincronizar variables CSS globales para colores de estado (Success, Error, Primary, Secondary).
- [x] Reforzar la visibilidad de componentes en distintos monitores mediante el ajuste de bordes y sombras.

## 🛠️ Implementación

### Iconografía
Se ha prohibido el uso de Material Symbols Icons y Emojis en la interfaz profesional. Se ha configurado la carga de **Material Symbols Outlined** y se ha actualizado el sistema de botones para gestionar correctamente el tamaño de estos iconos (`20px` base, `18px` en compactos).

### Botones Semánticos
Se han centralizado los estilos en `_buttons.css`, eliminando definiciones locales redundantes.
- **Primary:** Amarillo (`#E6B800`).
- **Secondary:** Azul corporativo.
- **Outline:** Borde gris para acciones secundarias.
- **Danger:** Rojo para acciones críticas.
- **Success:** Verde para validaciones positivas.
- **Ghost:** Minimalista para iconos.

## ✅ Resultados
- Sistema de botones 100% coherente en toda la aplicación.
- Eliminación de inconsistencias cromáticas en el módulo MES y Ventas.
- Guía de referencia técnica creada en `docs/guides/developer/design-system-atoms.md`.

## 📂 Artefactos
- `apps/frontend/src/design-system/_buttons.css`
- `docs/guides/developer/design-system-atoms.md`
- `apps/frontend/src/design-system/_variables.css`
