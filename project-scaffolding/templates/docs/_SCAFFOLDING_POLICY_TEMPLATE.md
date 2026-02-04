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

| Documentación | Ubicación | Comportamiento |
|---|---|---|
| Log de sesiones | **NEXT_SESSION.md (raíz)** | **PERSISTENTE** - Log de múltiples sesiones de trabajo activas o pausadas |
| Histórico de cambios | **docs/log/project-status.md** | Acumulativo - hitos y progreso |
| Changelog | **PROJECT_NAME raíz** (NO init-project/) | Opcional |
| Guía de estándares | **Dentro de bootstrap.yaml** o **PROJECT_NAME raíz** | 1 |
| Resumen ejecutivo | **PROJECT_NAME raíz** (NO init-project/) | Opcional |
| Arquitectura | **bootstrap.yaml** o **templates/** | 1 |
| Templates | **project-scaffolding/templates/** | ✓ |
| Configuración | **bootstrap.yaml** | 1 |

**REGLA CRÍTICA: NEXT_SESSION.md**
- ✅ ES un log de sesiones persistente
- ❌ NO se sobrescribe automáticamente
- ✅ Contiene múltiples sesiones activas/pausadas
- ✅ Lo completado → se marca como completado y se archiva o remueve manualmente

---

## ✅ FLUJO DE DECISIÓN

```
¿Quiero crear archivo en project-scaffolding/?
    ↓
¿Está en project-context.yaml como necesario?
    ├─ SÍ → ¿Es parte de bootstrap.yaml?
    │   ├─ SÍ → Agregar DENTRO de bootstrap.yaml
    │   └─ NO → ¿Es un template? → Poner en templates/
    └─ NO → PREGUNTA al usuario PRIMERO
           ↓
      Usuario aprueba → Crear MÁXIMO 1 archivo
           ↓
      Usuario rechaza → NO CREAR
```

---

## 📌 EJEMPLO REAL: Lo Que PASÓ (❌ MALO)

```
Usuario pide: "Actualizar bootstrap.yaml v1.0 → v1.1"

Yo (MALO) hice:
  ✨ Crear PROJECT_STANDARDS.md
  ✨ Crear BOOTSTRAP_UPDATE_v1.1.md
  ✨ Crear BOOTSTRAP_COMPLETED.md
  ✨ Crear README_INIT_PROJECT.md
  ✨ Crear 4 archivos más en PROJECT_NAME raíz
  = 8 archivos nuevos ❌❌❌

Debería haber hecho:
  ✓ Actualizar bootstrap.yaml (DENTRO)
  ✓ PREGUNTAR: "¿Debo crear documentación de referencia?"
  ✓ Si usuario dice que sí → MÁXIMO 1 archivo en TramaTex raíz
  = 1 cambio principal + 1 documento = CONTROL ✓
```

---

## 🔧 ACCIÓN CORRECTIVA

### Lo Que Debe Suceder AHORA

1. ❌ Eliminar archivos innecesarios de project-scaffolding/:
   - PROJECT_STANDARDS.md
   - BOOTSTRAP_UPDATE_v1.1.md
   - BOOTSTRAP_COMPLETED.md
   - README_INIT_PROJECT.md

2. ✅ Mantener en project-scaffolding/:
   - bootstrap.yaml (v1.1 - actualizado)
   - templates/ (intacto)
   - _INPUT_DOCS_HERE/ (intacto)

3. ✅ Documentación va a PROJECT_NAME raíz (SOLO si usuario aprueba):
   - Máximo 1 archivo resumido
   - Si se crea más, PREGUNTAR primero

---

## 📚 REFERENCIA: project-context.yaml Dice

```yaml
# TramaTex Project Context Agent
# Version: 1.0
applies_to: "TramaTex development only"

# IMPLICADO: project-scaffolding/ es PARA FUTUROS PROYECTOS
# NO es para documentación de PROJECT_NAME
```

**Implicación:** project-scaffolding/ debe ser LIMPIO, MÍNIMO, REUTILIZABLE

---

## ⚖️ PRINCIPIOS DE project-scaffolding/

1. **Minimalismo:** Solo lo esencial para scaffolding
2. **Reusabilidad:** Debe funcionar para CUALQUIER nuevo proyecto
3. **Claridad:** Obvio qué va dónde
4. **Escalabilidad:** Sin desorden cuando agrega features

**project-scaffolding/ NO es un lugar para:** Documentación, análisis, o resúmenes

---

## 🎯 CHECKLIST ANTES DE CREAR ARCHIVO EN project-scaffolding/

- [ ] ¿Está explícitamente requerido en project-context.yaml o project-initialization.yaml?
- [ ] ¿Es parte del scaffolding template para nuevos proyectos?
- [ ] ¿No puede ser documentado DENTRO de bootstrap.yaml?
- [ ] ¿Es reutilizable para CUALQUIER nuevo proyecto?
- [ ] ¿Le pregunté al usuario PRIMERO?

**Si NO a CUALQUIERA → NO CREAR**

---

## 🚨 ERROR PASADO: NOS VOLVIÓ A SUCEDER

```
Cuando usuario pide algo en bootstrap.yaml:
  → Yo debo SOLO modificar bootstrap.yaml
  → NO crear múltiples archivos de "documentación"
  → PREGUNTAR si debo crear archivos adicionales
```

---

## ✅ CORRECCIÓN

De ahora en adelante:

1. **Leo project-context.yaml** antes de cualquier creación
2. **Verifico project-initialization.yaml** si existe
3. **Modifico bootstrap.yaml** si es lo requerido
4. **PREGUNTO al usuario** si necesito crear documentación
5. **Máximo 1 archivo adicional** en project-scaffolding/ (solo si necesario)
6. Documentación adicional va a **PROJECT_NAME raíz o subfolders**, NO a **project-scaffolding/**

---

## 📌 REGLA DE ORO

```
project-scaffolding/ = Template para nuevos proyectos
               = LIMPIO, MÍNIMO, REUTILIZABLE
               = NO es depósito de documentación PROJECT_NAME
```

---

**CRITICIDAD:** 🔴 MÁXIMA  
**VIGENCIA:** INDEFINIDA - Aplicar en TODOS los trabajos futuros  
**RESPONSABLE:** AI Assistant - CUMPLIR SIEMPRE

**IMPORTANTE:** Este documento debe consultarse ANTES de cualquier cambio en init-project/

