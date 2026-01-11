# 📖 GUÍA RÁPIDA: PROMPTS Y CONTEXTO DE TRAMATEX

**Ubicación:** `/docs/SESSION_PROMPT.md` y archivos relacionados

Esta carpeta contiene toda la información necesaria para que **GitHub Copilot** (y otros LLMs) tengan contexto completo sobre el proyecto TramaTex en cada nueva sesión.

---

## 📚 ARCHIVOS DISPONIBLES

### 1. **SESSION_PROMPT.md** (Principal)
**Lectura obligatoria antes de cada sesión**

- ✅ Descripción completa del proyecto
- ✅ Stack tecnológico y arquitectura
- ✅ Fases y cronograma
- ✅ Estructura de carpetas
- ✅ Documentación disponible
- ✅ Patrones y convenciones
- ✅ Principios y restricciones
- ✅ **Sección final para rellenar objetivos de la sesión**

**Tamaño:** ~500 líneas | **Tiempo lectura:** 20-30 min la primera vez

👉 **COMIENZA AQUÍ** para primera sesión

---

### 2. **SESSION_PROMPT_QUICK.md** (Quick Reference)
**Para sesiones posteriores o cuando tienes prisa**

- ⚡ Resumen de 1 página
- ⚡ Tabla referencia rápida
- ⚡ Links a documentación completa
- ⚡ Checklist fin de sesión

**Tamaño:** ~150 líneas | **Tiempo lectura:** 5 min

👉 **USA ESTO** después de Session 1 para refrescar contexto rápidamente

---

### 3. **COPILOT_INSTRUCTIONS.md** (Instrucciones Específicas)
**Directivas para GitHub Copilot sobre cómo colaborar**

- 🤖 Rol y responsabilidades
- 🤖 Criterios de aceptación para código
- 🤖 Flujo de trabajo Copilot
- 🤖 Reglas de decisión arquitectónicas
- 🤖 Lo que NUNCA debe hacer
- 🤖 Métricas que importan

**Tamaño:** ~350 líneas | **Tiempo lectura:** 10 min

👉 **REFERENCIA** al inicio de cada sesión técnica

---

### 4. **SESSION_INIT_TEMPLATE.md** (Template para Nueva Sesión)
**Plantilla a rellenar cuando inicia nueva sesión**

- 📋 Campos de información sesión
- 📋 Sección de objetivos (3-5 objetivos)
- 📋 Contexto de entrada (estado anterior)
- 📋 Plan de trabajo por fases
- 📋 Changes made, commits, archivos modificados
- 📋 Definición de "hecho"
- 📋 Problemas/bloqueadores encontrados
- 📋 Decisiones arquitectónicas tomadas
- 📋 Métricas finales
- 📋 Próximos pasos

**Tamaño:** ~400 líneas | **Tiempo:** Rellenar durante sesión

👉 **COPIA Y RELLENA** para cada nueva sesión

---

## 🔄 FLUJO DE USO RECOMENDADO

### Primera Sesión (Session 01+)

```mermaid
1. Leer: SESSION_PROMPT.md (sección 1-5, muy importante)
   ↓
2. Revisar: Última sesión en docs/sessions/
   ↓
3. Confirmar: Objetivos con usuario (rellenar SESSION_PROMPT.md sección final)
   ↓
4. Iniciar: Código/tests
   ↓
5. Documentar: En docs/sessions/2026-MM-DD-session-NN.md
   ↓
6. Actualizar: PROJECT_STATUS.md
```

### Sesiones Posteriores (Session 2+)

```mermaid
1. Leer: SESSION_PROMPT_QUICK.md (2 min)
   ↓
2. Leer: COPILOT_INSTRUCTIONS.md (5 min)
   ↓
3. Revisar: Última sesión completada
   ↓
4. COPIAR: SESSION_INIT_TEMPLATE.md → docs/sessions/2026-MM-DD-session-NN.md
   ↓
5. Rellenar: Objetivos y contexto
   ↓
6. Trabajar: Implementación
   ↓
7. Completar: Sesión template con resultados
```

---

## 🎯 CUÁNDO USAR CADA ARCHIVO

| Situación | Usar | Tiempo |
|-----------|------|--------|
| **Primera vez en proyecto** | SESSION_PROMPT.md (completo) | 20-30 min |
| **Segunda+ sesión, refrescar contexto** | SESSION_PROMPT_QUICK.md | 5 min |
| **Copilot comportándose extraño** | COPILOT_INSTRUCTIONS.md | 10 min |
| **Inicio nueva sesión** | SESSION_INIT_TEMPLATE.md | Rellenar durante sesión |
| **Verificar decisiones arquitectónicas** | SESSION_PROMPT.md + ADRs | Según necesidad |
| **Entender qué hizo sesión anterior** | docs/sessions/YYYY-MM-DD-session-NN.md | Según necesidad |

---

## ✅ CHECKLIST: ESTÁS LISTO PARA LA SESIÓN

- [ ] Leí SESSION_PROMPT.md (al menos secciones 1-2)
- [ ] Revisé SESSION_PROMPT_QUICK.md como referencia rápida
- [ ] Leí última sesión completada en /docs/sessions/
- [ ] Tengo a mano COPILOT_INSTRUCTIONS.md
- [ ] Copié SESSION_INIT_TEMPLATE.md como base para esta sesión
- [ ] Confirmé objetivos de esta sesión con user
- [ ] Docker está funcionando (`make docker-up`)
- [ ] Backend/Frontend compilan sin errores
- [ ] Entiendo la arquitectura DDD + Clean Architecture

---

## 📊 CONTEXTO RÁPIDO

| Aspecto | Valor |
|--------|-------|
| **Proyecto** | TramaTex – ERP/MES para microempresas textil |
| **Stack** | Go + Vue.js 3 + PostgreSQL |
| **Arquitectura** | Clean Architecture + DDD con rigor asimétrico |
| **Cronograma** | 24 meses (782 horas), 8h/semana |
| **Fase Actual** | Fase 0 (Fundaciones) |
| **Sesiones Completadas** | 8 (hasta 11/01/2026) |
| **Estado Código** | 0 (puro setup), 1435+ líneas documentación |

---

## 🔗 REFERENCIAS INTERNAS

| Necesito | Ubicación |
|---------|-----------|
| **Especificación técnica completa** | `/docs/consolidated/DOCUMENTO-CONSOLIDADO-3.0.md` |
| **Stack tecnológico** | `/docs/adr/ADR-001-*.md` |
| **Arquitectura software** | `/docs/adr/ADR-002-*.md` |
| **Orden implementación** | `/docs/adr/ADR-007-*.md` |
| **Cronograma** | `/docs/adr/ADR-008-*.md` |
| **Estructura carpetas** | `/docs/adr/ADR-009-*.md` |
| **Estado actual** | `/PROJECT_STATUS.md` |
| **Todas las sesiones** | `/docs/sessions/` |

---

## 🚀 PRÓXIMOS PASOS

1. **Ahora:** Lee [SESSION_PROMPT.md](SESSION_PROMPT.md) (completo para primera sesión)
2. **Luego:** Rellena sección final "OBJETIVOS DE ESTA SESIÓN"
3. **Durante sesión:** Usa [COPILOT_INSTRUCTIONS.md](COPILOT_INSTRUCTIONS.md) como referencia
4. **Fin de sesión:** Documenta en `/docs/sessions/2026-MM-DD-session-NN.md`

---

## 💡 TIPS IMPORTANTES

- **Dominio es activo crítico:** Clean Architecture estricta, TDD obligatorio en módulos críticos
- **Local-first:** Sin cloud MVP, 100% operativo offline
- **Rigor asimétrico:** Tarificación/Party/Producto con rigor estricto, infraestructura flexible
- **Commits descriptivos:** Facilitan mantenimiento y trazabilidad
- **Documentación sincronizada:** Sesiones registradas + ADRs = trazabilidad total

---

**Última actualización:** 11/01/2026

**Mantener sincronizado con:**
- `PROJECT_STATUS.md` (avances actuales)
- `docs/sessions/` (sesiones nuevas)
- `docs/adr/` (decisiones arquitectónicas)

