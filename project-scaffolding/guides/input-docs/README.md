# 📂 Directorio de Entrada de Documentos

## 🎯 Propósito

Coloca aquí cualquier documentación de tu proyecto ANTES de ejecutar el bootstrap.
El sistema la analizará automáticamente y extraerá información para evitar preguntas redundantes.

## 📥 ¿Qué Documentos Colocar?

Cualquier documento que contenga información del proyecto:
- README, especificaciones técnicas, requisitos
- Documentos de arquitectura, diseño, planificación
- Cualquier otra información relevante

**Formatos soportados:**
- ✅ Markdown (`.md`)
- ✅ Texto plano (`.txt`)
- ✅ YAML (`.yaml`, `.yml`)
- ✅ JSON (`.json`)
- ✅ Word (`.docx`)
- ✅ PDF (`.pdf`)

## 🔍 ¿Qué Información Extraerá?

| Campo | Cómo lo Encuentra |
|-------|-------------------|
| **Nombre** | Títulos H1, "Nombre:", "Project Name:" |
| **Visión** | Secciones "Visión", "Propósito", primer párrafo |
| **Componentes** | Infiere de: "API", "frontend", "backend", "BFF" |
| **Tecnologías** | Detecta: Go, Python, Vue, React, PostgreSQL, etc. |
| **Base de datos** | Menciones explícitas de BD |
| **Módulos** | Listas de bounded contexts, módulos |

## 📊 Resultados

- **Documentación completa:** 0 preguntas adicionales
- **Documentación parcial:** Solo pregunta lo que falta
- **Sin documentación:** Pregunta todo (flujo tradicional)

## 💡 Ver Ejemplo

Consulta el archivo EJEMPLO_README.md para ver un template completo.