# Sistema de Scaffolding de Proyectos

## 🎯 Propósito

Este directorio contiene el sistema para generar automáticamente la estructura de nuevos proyectos con documentación, agentes IA y mejores prácticas integradas.

## 📚 Documentación

[LICENSE](LICENSE) - La licencia bajo la que se distribuye este proyecto.

### Para Usuarios (Crear Proyectos)
👉 **[user-guide.md](guides/user-guide.md)** - Cómo usar el scaffolding

### Para Desarrolladores (Modificar Scaffolding)
👉 **[development-guide.md](guides/development-guide.md)** - Cómo modificar y mejorar el sistema

## 🚀 Inicio Rápido

1. **Coloca documentos** (opcional) en `user-input-docs/`
2. **Ejecuta el bootstrap** - El sistema extrae info automáticamente
3. **Confirma configuración** - Solo pregunta lo que falta
4. **¡Listo!** - Proyecto creado con estructura completa

## 📁 Estructura

```
project-scaffolding/
├── bootstrap.yaml              # Configuración del flujo
├── guides/                     # Guías y documentación del scaffolding
│   ├── user-guide.md           # Para usuarios (cómo usar el scaffolding)
│   ├── development-guide.md    # Para desarrolladores (cómo modificar el scaffolding)
│   ├── placeholders-guide.md   # Documentación del sistema de placeholders
│   └── input-docs/             # Guías para preparar documentos de entrada
│       ├── README.md           # Guía del directorio
│       └── EJEMPLO_README.md   # Template de ejemplo
├── user-input-docs/            # Coloca aquí tus documentos (el sistema los analizará)
├── agents/                     # Agentes del scaffolding
│   └── scaffolding-developer.yaml
├── tmp/                        # Archivos temporales (ver README)
│   └── README.md               # Reglas del directorio temporal
└── templates/                  # Templates para proyectos
    ├── agents/
    ├── docs/
    └── .github/
```

## ✨ Características

- 🧠 **Extracción inteligente** de metadata de documentos
- 📊 **Preguntas mínimas** - Solo lo que no encuentra
- 📁 **Multi-formato** - MD, TXT, YAML, JSON, DOCX, PDF
- 🎯 **Inferencia contextual** - Detecta componentes y tecnologías
- ✅ **Transparente** - Muestra qué extrajo y de dónde

## 🔧 Para Modificar el Sistema

Consulta **[development-guide.md](guides/development-guide.md)** y el agente **[scaffolding-developer.yaml](agents/scaffolding-developer.yaml)**