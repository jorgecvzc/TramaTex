# ❌ POLÍTICA ESTRICTA: project-scaffolding/ - MÁXIMO 1 ARCHIVO DE DOCUMENTACIÓN

**CRITICIDAD:** 🔴 MÁXIMA  
**AUDIENCIA:** AI Assistants, Developers  
**PROPÓSITO:** Prevenir generación descontrolada de archivos en project-scaffolding/

---

## 🚨 REGLA FUNDAMENTAL

### ❌ PROHIBIDO CREAR MÚLTIPLES ARCHIVOS DE DOCUMENTACIÓN

```
❌ NUNCA:
  project-scaffolding/
  ├── PROJECT_STANDARDS.md          ← NO
  ├── BOOTSTRAP_UPDATE_v1.1.md      ← NO
  ├── BOOTSTRAP_COMPLETED.md        ← NO
  ├── README_INIT_PROJECT.md        ← NO
  └── bootstrap.yaml

✅ SÍ:
  project-scaffolding/
  ├── bootstrap.yaml                ← PRINCIPAL
  ├── README.md (máximo 1)          ← OPCIONAL, solo si es necesario
  └── (templates y ejemplos)
```

---

## 📋 REGLAS ESTRICTAS

### 1. **project-scaffolding/ es para Templates y Configuración**

**PERMITIDO:**
- ✅ bootstrap.yaml (configuración principal)
- ✅ Directorio `templates/` (plantillas)
- ✅ Directorio `_INPUT_DOCS_HERE/` (entrada)
- ✅ README.md (MÁXIMO 1, solo si es absolutamente necesario)

**PROHIBIDO:**
- ❌ Archivos .md adicionales
- ❌ Changelogs
- ❌ Guías de referencia
- ❌ Resúmenes de cambios
- ❌ Múltiples documentos de soporte

### 2. **Antes de Crear CUALQUIER Archivo**

**SIEMPRE preguntar:**
1. ¿Este archivo está en project-context.yaml o project-initialization.yaml como necesario?
2. ¿Es parte del scaffolding template?
3. ¿Podría documentarse DENTRO de bootstrap.yaml?
4. ¿Es realmente necesario o solo "útil"?

**Si la respuesta es NO a cualquiera:** NO CREAR. PREGUNTAR al usuario.

### 3. **Dónde Va la Documentación**

| Documentación | Ubicación | Notas |
|---|---|---|
| Logs de sesión | **NEXT_SESSION.md (raíz)** | Para sesiones de trabajo activas o pausadas |
| Histórico de cambios | **docs/log/project-status.md** | Registro acumulativo de hitos y progreso |
| Changelog | **docs/log/** (ej. `CHANGELOG.md`) | Registro de cambios del proyecto |
| Guías y estándares | **docs/guides/developer/** | Guías para desarrolladores, incluyendo estándares de código |
| Resúmenes ejecutivos | **docs/log/milestones/** | Reportes y resúmenes de hitos |
| Arquitectura | **docs/architecture/** | ADRs, diagramas, visión arquitectónica |
| Templates | **project-scaffolding/templates/** | Plantillas para el scaffolding de proyectos |
| Políticas de gobernanza | **docs/log/governance/** | Documentos que definen reglas y procesos del proyecto |

**REGLA CLAVE: Documentación en `docs/`**
- Toda la documentación del proyecto (excepto `README.md`, `AGENTS.md`, `NEXT_SESSION.md`) DEBE residir en el directorio `docs/` bajo la subcarpeta apropiada.
- Esto asegura una estructura coherente y fácil de navegar.---

## ✅ CORRECCIÓN

De ahora en adelante:

1. **Leo project-context.yaml** antes de cualquier creación
2. **Verifico project-initialization.yaml** si existe
3. **Modifico bootstrap.yaml** si es lo requerido
4. **PREGUNTO al usuario** si necesito crear documentación
5. **Máximo 1 archivo adicional** en project-scaffolding/ (solo si necesario)
6. Documentación adicional va a **docs/ subfolders**, NO a **project-scaffolding/**. La única excepción es para `README.md` y `AGENTS.md` en la raíz del proyecto, según la "Política de Directorio Raíz Limpio".

---

## 📌 REGLA DE ORO

```
project-scaffolding/ = Template para nuevos proyectos
               = LIMPIO, MÍNIMO, REUTILIZABLE
               = NO es depósito de documentación del proyecto final
```

