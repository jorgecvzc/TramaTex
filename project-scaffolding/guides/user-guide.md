# Guía de Uso - Sistema de Scaffolding de Proyectos

## 🚀 Inicio Rápido

### 1. Prepara tu Documentación (Opcional)
Coloca cualquier documento del proyecto en `user-input-docs/`:
- ✅ README.md, especificaciones, arquitectura, requisitos
- ✅ Formatos: `.md`, `.txt`, `.yaml`, `.json`, `.docx`, `.pdf`

### 2. Ejecuta el Bootstrap
El sistema automáticamente:
1. Busca y analiza tus documentos
2. Extrae: nombre, visión, tecnologías, módulos, arquitectura
3. Muestra lo que encontró (con nivel de confianza)
4. Solo pregunta lo que NO pudo extraer

### 3. Resultado
- **Con documentación completa:** 0 preguntas adicionales
- **Con documentación parcial:** Solo 2-5 preguntas
- **Sin documentación:** 7 preguntas (flujo tradicional)

## 📝 ¿Qué Información Extrae?

| Campo | Dónde lo Busca |
|-------|----------------|
| **Nombre** | Títulos H1, "Nombre:", "Project Name:" |
| **Visión** | Secciones "Visión", "Propósito", primer párrafo |
| **Componentes** | Infiere de: "API", "frontend", "backend" |
| **Tecnologías** | Detecta: Go, Python, Vue, React, PostgreSQL |
| **Módulos** | Listas de módulos, bounded contexts |

## 💡 Ejemplo Rápido

Ver `guides/input-docs/EJEMPLO_README.md` para un template completo.

## ℹ️ Más Información

- **Directorio de entrada:** Ver `guides/input-docs/README.md`
- **Para desarrollar/modificar:** Ver `guides/development-guide.md`
