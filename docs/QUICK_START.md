# 🚀 QUICK START – REFERENCIA MÁXIMA PARA INICIAR SESIÓN

**Tamaño:** ~2 minutos lectura  
**Propósito:** De un vistazo saber qué hacer ahora

---

## ⚡ ESTOY AQUÍ → QUÉ HAGO

| Situación | Acción | Tiempo |
|-----------|--------|--------|
| **Primera vez en proyecto** | Leer [SESSION_PROMPT.md](SESSION_PROMPT.md) sección 1-3 | 20 min |
| **Antes de sesión técnica** | Leer [COPILOT_INSTRUCTIONS.md](COPILOT_INSTRUCTIONS.md) | 10 min |
| **Voy a iniciar sesión NN** | Copiar [SESSION_INIT_TEMPLATE.md](SESSION_INIT_TEMPLATE.md) → `sessions/2026-MM-DD-session-NN.md` | 1 min |
| **Necesito contexto rápido** | Leer [SESSION_PROMPT_QUICK.md](SESSION_PROMPT_QUICK.md) | 5 min |
| **Tengo la sesión abierta** | Rellena sección "Objetivos principales" en template | 10 min |
| **Voy a hablar con Copilot** | Comparte [SESSION_PROMPT.md](SESSION_PROMPT.md) o [SESSION_PROMPT_QUICK.md](SESSION_PROMPT_QUICK.md) | — |
| **Necesito ver qué hacer** | Leer última sesión en `docs/sessions/` | 5-10 min |
| **Termino sesión** | Rellena resto de template, commit, update [PROJECT_STATUS.md](../PROJECT_STATUS.md) | 15 min |
| **Me pierdo navegando** | Consultar [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md) | 5-10 min |
| **Quiero ejemplo completo** | Leer [EJEMPLO_PRACTICO.md](EJEMPLO_PRACTICO.md) | 15 min |

---

## 📋 ANTES DE CADA SESIÓN

```
☑ Paso 1: Leo SESSION_PROMPT_QUICK.md (5 min)
☑ Paso 2: Leo última sesión en docs/sessions/ (5 min)
☑ Paso 3: Copio SESSION_INIT_TEMPLATE.md → nueva sesión (1 min)
☑ Paso 4: Relleno objetivos en template (10 min)
☑ Paso 5: Comparto SESSION_PROMPT.md + COPILOT_INSTRUCTIONS.md con Copilot
☑ Paso 6: Comenzamos trabajo
```

**Total tiempo prep:** 20-25 minutos

---

## 🎯 PROYECTO EN 30 SEGUNDOS

| Aspecto | Valor |
|--------|-------|
| **Qué es** | ERP/MES para microempresas textil (pedidos + producción) |
| **Stack** | Go + Vue.js 3 + PostgreSQL |
| **Arquitectura** | Clean Architecture + DDD (rigor asimétrico) |
| **Duración** | 24 meses (782 horas), 8h/semana |
| **Fase actual** | Fase 0 (Fundaciones) |
| **Sesiones** | 8 completadas, inicia sesión 09 |

---

## 🔑 3 PRINCIPIOS CLAVE

1. **Dominio es sagrado** → Clean Architecture estricta, TDD obligatorio
2. **Rigor asimétrico** → Tarificación/Party estricto, resto flexible
3. **Local-first** → Sin cloud MVP, 100% offline

---

## 📚 DOCUMENTOS SEGÚN NECESIDAD

| Necesito | Documento |
|----------|-----------|
| Contexto completo | [SESSION_PROMPT.md](SESSION_PROMPT.md) |
| Contexto rápido | [SESSION_PROMPT_QUICK.md](SESSION_PROMPT_QUICK.md) |
| Cómo codificar | [COPILOT_INSTRUCTIONS.md](COPILOT_INSTRUCTIONS.md) |
| Template sesión | [SESSION_INIT_TEMPLATE.md](SESSION_INIT_TEMPLATE.md) |
| Guía prompts | [README_PROMPTS.md](README_PROMPTS.md) |
| Índice todo | [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md) |
| Ejemplo práctico | [EJEMPLO_PRACTICO.md](EJEMPLO_PRACTICO.md) |
| Especificación MVP | [consolidated/DOCUMENTO-CONSOLIDADO-3.0.md](consolidated/DOCUMENTO-CONSOLIDADO-3.0.md) |
| Decisiones arquitectónicas | [adr/ADR-001 a ADR-009](adr/) |
| Estado proyecto | [PROJECT_STATUS.md](../PROJECT_STATUS.md) |
| Sesiones anteriores | [sessions/](sessions/) |

---

## ✅ VALIDACIÓN PRE-COMMIT

Antes de hacer commit:

```bash
# Tests ✅
go test ./...

# Lint ✅
golangci-lint run ./...

# Documentación ✅
- Sesión registrada en docs/sessions/
- PROJECT_STATUS.md actualizado

# Commits ✅
- Mensajes descriptivos
- Una feature por commit
```

---

## 🚀 LISTA DE REPRODUCCIÓN (Orden Lectura)

**Sesión 1 (Primera vez):**
1. README.md (5 min)
2. SESSION_PROMPT.md § 1-3 (20 min)
3. ADR-002 (Clean Architecture) (15 min)

**Sesiones 2+ (Antes de cada sesión):**
1. SESSION_PROMPT_QUICK.md (5 min)
2. Última sesión completada (5 min)
3. COPILOT_INSTRUCTIONS.md (5 min)

**Profundo (Tech Lead):**
1. DOCUMENTO-CONSOLIDADO-3.0.md (60 min)
2. ADR-001 a ADR-009 (90 min)

---

## 💻 COMANDOS BÁSICOS

```bash
# Ver ayuda
make help

# Docker
make docker-up
make docker-down

# Backend tests
cd backend && go test ./...

# Backend run
cd backend && go run ./cmd/api

# Frontend dev
cd frontend && npm run dev

# Lint
golangci-lint run ./...
```

---

## 📊 ESTRUCTURA CRÍTICA

```
backend/internal/
├── domain/          ← Clean Architecture estricta
├── application/     ← Orquestación
├── infrastructure/  ← GORM, BD
└── interfaces/      ← REST, DTOs
```

**Regla de oro:** Dependencias siempre hacia DENTRO

---

## 🎓 CHECKLIST "ESTOY LISTO"

- [ ] Leí SESSION_PROMPT.md o SESSION_PROMPT_QUICK.md
- [ ] Entiendo que Clean Architecture + DDD es no-negociable
- [ ] Sé que TDD es obligatorio en dominio crítico
- [ ] Leí última sesión completada
- [ ] Tengo Copilot listo con contexto
- [ ] Docker está funcionando
- [ ] Entiendo la fase actual y qué viene

---

## 🚨 LO QUE NUNCA HACER

❌ Lógica de negocio fuera del dominio  
❌ ORM fuera de infraestructura  
❌ Código sin tests en dominio crítico  
❌ Ignorar Clean Architecture  
❌ Hardcodear valores  
❌ Commit sin mensaje descriptivo  

---

## 🔗 NAVEGACIÓN RÁPIDA

```
Estoy perdido ↓
¿Qué documento? ↓ [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md)

Necesito contexto ↓
¿Cuánto tiempo? ↓ [README_PROMPTS.md](README_PROMPTS.md)

Copilot comportándose raro ↓
Lee ↓ [COPILOT_INSTRUCTIONS.md](COPILOT_INSTRUCTIONS.md)

Quiero saber qué pasó ↓
Última sesión ↓ [docs/sessions/](sessions/)
```

---

## 📞 AYUDA RÁPIDA

**"¿Cómo empiezo?"**
→ Leer [SESSION_PROMPT.md](SESSION_PROMPT.md) y rellenar final

**"¿Qué hago durante sesión?"**
→ Rellenar [SESSION_INIT_TEMPLATE.md](SESSION_INIT_TEMPLATE.md) conforme avanzas

**"¿Copilot no entiende?"**
→ Compartir [COPILOT_INSTRUCTIONS.md](COPILOT_INSTRUCTIONS.md)

**"¿Necesito ejemplo?"**
→ Leer [EJEMPLO_PRACTICO.md](EJEMPLO_PRACTICO.md)

**"¿Cuál es el estado?"**
→ Revisar [PROJECT_STATUS.md](../PROJECT_STATUS.md) y última sesión

**"¿Qué documentos tengo?"**
→ Consultar [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md)

---

## 🎯 PRÓXIMA SESIÓN

**Sesión 09 (Próxima):**
```bash
cp docs/SESSION_INIT_TEMPLATE.md docs/sessions/2026-01-18-session-09.md
# Rellena objetivos: User entity, JWT, Login handler
# Comienza!
```

---

**TL;DR:**
1. Leer [SESSION_PROMPT.md](SESSION_PROMPT.md) (primera vez) o [SESSION_PROMPT_QUICK.md](SESSION_PROMPT_QUICK.md) (siempre)
2. Copiar template → nueva sesión
3. Rellenar objetivos
4. Compartir con Copilot
5. Trabajar
6. Documentar

---

**Última actualización:** 11/01/2026

