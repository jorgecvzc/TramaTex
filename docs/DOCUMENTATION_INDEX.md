# 📑 ÍNDICE COMPLETO DE DOCUMENTACIÓN TRAMATEX

**Última actualización:** 11/01/2026  
**Propósito:** Mapa de toda la documentación del proyecto para navegación rápida

---

## 🗂️ ESTRUCTURA GENERAL

```
docs/
├── SESSION_PROMPT.md              ← COMIENZA AQUÍ (contexto completo)
├── SESSION_PROMPT_QUICK.md        ← Referencia rápida (post-sesión 1)
├── COPILOT_INSTRUCTIONS.md        ← Directivas para LLM
├── SESSION_INIT_TEMPLATE.md       ← Template para nueva sesión
├── README_PROMPTS.md              ← Guía de prompts (este archivo)
├── DOCUMENTATION_INDEX.md         ← Índice de TODA documentación
│
├── adr/                           # Architecture Decision Records
│   ├── _ADR_TEMPLATE.md
│   ├── ADR-001-seleccion-stack-tecnologico.md
│   ├── ADR-002-adopcion-clean-architecture-ddd.md
│   ├── ADR-003-tipo-distribucion-aplicacion.md
│   ├── ADR-004-ciclo-vida-desarrollo-mvp.md
│   ├── ADR-005-gestion-unificada-clientes-proveedores.md
│   ├── ADR-006-estrategia-desarrollo-dirigido-dominio.md
│   ├── ADR-007-orden-implementacion-modulos.md
│   ├── ADR-008-planificacion-cronograma-mvp.md
│   └── ADR-009-estructura-proyecto.md
│
├── consolidated/
│   └── DOCUMENTO-CONSOLIDADO-3.0.md    # Especificación técnica MVP
│
├── modules/                       # Especificaciones por módulo
│   ├── _TEMPLATE.md
│   ├── [pendiente: party.md]
│   ├── [pendiente: product.md]
│   ├── [pendiente: pricing.md]
│   ├── [pendiente: sales.md]
│   └── [pendiente: mes.md]
│
├── diagrams/
│   ├── architecture/              # Diagramas C4, componentes
│   ├── domain/                    # Diagramas de dominios/bounded contexts
│   └── flows/                     # Flujos de casos de uso
│
├── guides/                        # Guías de desarrollo
│   ├── [pendiente: setup.md]
│   ├── [pendiente: clean-architecture.md]
│   ├── [pendiente: testing-strategy.md]
│   └── [pendiente: deployment.md]
│
└── sessions/                      # Registro por sesión
    ├── _SESSION_TEMPLATE.md
    ├── 2026-01-06-session-01.md
    ├── 2026-01-06-session-02.md
    ├── 2026-01-07-session-03.md
    ├── 2026-01-09-session-04.md
    ├── 2026-01-09-session-05.md
    ├── 2026-01-09-session-06.md
    ├── 2026-01-10-session-07.md
    └── 2026-01-11-session-08.md
```

---

## 🎯 DOCUMENTACIÓN POR PROPÓSITO

### 📌 PARA ENTENDER EL PROYECTO (Primera vez)

| Documento | Propósito | Tiempo | Debe leer |
|-----------|-----------|--------|----------|
| [README.md](../README.md) | Descripción general, stack, quickstart | 5 min | ✅ Todos |
| [SESSION_PROMPT.md](SESSION_PROMPT.md) | Contexto completo + arquitectura | 20-30 min | ✅ Copilot sesión 1 |
| [DOCUMENTO-CONSOLIDADO-3.0.md](consolidated/DOCUMENTO-CONSOLIDADO-3.0.md) | Especificación técnica MVP completa | 40-60 min | ✅ Si necesita detalle |

### 📐 PARA ARQUITECTURA Y DECISIONES

| Documento | Tema | Lectura Crítica |
|-----------|------|-----------------|
| [ADR-001](adr/ADR-001-seleccion-stack-tecnologico.md) | Stack tecnológico (Go, Vue.js, PostgreSQL) | ✅ |
| [ADR-002](adr/ADR-002-adopcion-clean-architecture-ddd.md) | Clean Architecture + DDD rigor asimétrico | ✅ |
| [ADR-003](adr/ADR-003-tipo-distribucion-aplicacion.md) | Monolito modular local-first | ✅ |
| [ADR-004](adr/ADR-004-ciclo-vida-desarrollo-mvp.md) | Ciclo de vida, fases de desarrollo | ✅ |
| [ADR-005](adr/ADR-005-gestion-unificada-clientes-proveedores.md) | Modelo Party (clientes/proveedores) | ✅ |
| [ADR-006](adr/ADR-006-estrategia-desarrollo-dirigido-dominio.md) | Estrategia DDD para MVP | ✅ |
| [ADR-007](adr/ADR-007-orden-implementacion-modulos.md) | Orden de implementación módulos | ✅ |
| [ADR-008](adr/ADR-008-planificacion-cronograma-mvp.md) | Cronograma 24 meses, 782 horas | ✅ |
| [ADR-009](adr/ADR-009-estructura-proyecto.md) | Estructura de carpetas + convenciones | ✅ |

### 🔄 PARA COPILOT/LLM

| Documento | Para Qué | Cuando |
|-----------|----------|--------|
| [SESSION_PROMPT.md](SESSION_PROMPT.md) | Contexto completo del proyecto | Inicio sesión 1 |
| [SESSION_PROMPT_QUICK.md](SESSION_PROMPT_QUICK.md) | Referencia rápida, refrescar | Antes de sesión 2+ |
| [COPILOT_INSTRUCTIONS.md](COPILOT_INSTRUCTIONS.md) | Cómo debe comportarse Copilot | Antes de sesión técnica |
| [SESSION_INIT_TEMPLATE.md](SESSION_INIT_TEMPLATE.md) | Template para nueva sesión | Inicio cada sesión |
| [README_PROMPTS.md](README_PROMPTS.md) | Guía de estos archivos | Primeras veces |

### 📊 PARA SEGUIMIENTO

| Documento | Propósito | Frecuencia |
|-----------|-----------|-----------|
| [../PROJECT_STATUS.md](../PROJECT_STATUS.md) | Estado actual (horas, fase, hitos) | Actualizar cada sesión |
| [sessions/2026-01-11-session-08.md](sessions/2026-01-11-session-08.md) | Última sesión completada | Revisar antes sesión N+1 |
| [sessions/](sessions/) | Historial de todas las sesiones | Referencia histórica |

### 🏗️ PARA DESARROLLO (Por Módulo)

**Próximos:** Estos documentos especifican cada módulo/Bounded Context

- `modules/party.md` - Gestión Party/Organización
- `modules/product.md` - Gestión Productos/Variantes
- `modules/pricing.md` - Motor Tarificación
- `modules/sales.md` - Ventas/Pedidos
- `modules/mes.md` - Producción/MES

**Formato:** Especificación funcional, API, schema DB, casos de uso

### 📖 PARA GUÍAS (Próximas)

- `guides/setup.md` - Configuración entorno desarrollo
- `guides/clean-architecture.md` - Cómo seguir CA en proyecto
- `guides/testing-strategy.md` - TDD, cobertura, estrategia
- `guides/deployment.md` - Docker, despliegue, CI/CD

### 📊 PARA DIAGRAMAS (Próximos)

- `diagrams/architecture/c4-*.md` - Diagramas C4 (contexto, contenedor, componente)
- `diagrams/domain/bounded-contexts.md` - Mapa dominios
- `diagrams/flows/order-flow.md` - Flujo pedidos
- `diagrams/flows/pricing-flow.md` - Flujo tarificación

---

## 🎓 RUTAS DE APRENDIZAJE

### Para Desarrollador Backend (Go)

**Orden recomendado:**
1. `README.md` (5 min)
2. `SESSION_PROMPT.md` § 1-3, 4 (20 min)
3. `ADR-002` (Clean Architecture) (15 min)
4. `ADR-009` (Estructura carpetas) (10 min)
5. `COPILOT_INSTRUCTIONS.md` (10 min)
6. `DOCUMENTO-CONSOLIDADO-3.0.md` § 2-3 (RFC/RNF) (15 min)
7. Último `sessions/YYYY-MM-DD-session-NN.md` (5 min)

**Total:** ~80 min

### Para Desarrollador Frontend (Vue.js)

**Orden recomendado:**
1. `README.md` (5 min)
2. `SESSION_PROMPT.md` § 1-3, 4 (20 min)
3. `ADR-002` (Clean Architecture aplicado a Vue) (10 min)
4. `ADR-009` (Estructura carpetas frontend) (5 min)
5. `DOCUMENTO-CONSOLIDADO-3.0.md` § 2-3 (RFC/RNF) (15 min)
6. Último `sessions/YYYY-MM-DD-session-NN.md` (5 min)

**Total:** ~60 min

### Para Architect/Tech Lead

**Orden recomendado:**
1. `DOCUMENTO-CONSOLIDADO-3.0.md` (completo) (60 min)
2. `ADR-001` a `ADR-009` (completos) (90 min)
3. `SESSION_PROMPT.md` § 6-8 (principios, restricciones) (20 min)
4. `sessions/` (últimas 3 sesiones) (30 min)

**Total:** ~200 min

### Para LLM/Copilot (Nueva Sesión)

**Ruta rápida:**
1. `SESSION_PROMPT.md` (sección llena por usuario)
2. `COPILOT_INSTRUCTIONS.md`
3. Última sesión completada en `sessions/`

**Total:** ~15 min + lectura sesión anterior

---

## 🔗 REFERENCIAS CRUZADAS

### Si necesitas entender...

| Concepto | Ver |
|----------|-----|
| **Clean Architecture** | ADR-002 § 4, SESSION_PROMPT.md § 3 |
| **DDD / Bounded Contexts** | ADR-003, ADR-005, ADR-006 |
| **Módulos MVP** | SESSION_PROMPT.md § 2.2, DOCUMENTO-CONSOLIDADO-3.0.md § 2 |
| **Stack Tecnológico** | ADR-001, SESSION_PROMPT.md § 1.2 |
| **Cronograma** | ADR-008, PROJECT_STATUS.md |
| **Estructura Carpetas** | ADR-009, SESSION_PROMPT.md § 3 |
| **Próximos Pasos** | Última sesión en sessions/, PROJECT_STATUS.md § Próximos Hitos |
| **Decisiones Previas** | ADR-001 a ADR-009 |
| **Cómo Codificar** | SESSION_PROMPT.md § 6, COPILOT_INSTRUCTIONS.md |

---

## 📋 DOCUMENTACIÓN ESPERADA (Post-MVP)

Estos documentos se generarán durante el desarrollo:

### Especificaciones de Módulo
- `modules/party.md` - Party/Organización
- `modules/product.md` - Producto/Variante
- `modules/pricing.md` - Tarificación
- `modules/sales.md` - Ventas/Pedidos
- `modules/mes.md` - Producción/MES

### Guías de Desarrollo
- `guides/setup.md` - Environment setup
- `guides/clean-architecture.md` - Cómo aplicar CA
- `guides/testing-strategy.md` - TDD + cobertura
- `guides/api-design.md` - Convenciones REST
- `guides/database-design.md` - Schema + migraciones
- `guides/deployment.md` - Docker, CI/CD, producción

### Diagramas
- `diagrams/architecture/c4-context.md`
- `diagrams/architecture/c4-container.md`
- `diagrams/domain/bounded-contexts-map.md`
- `diagrams/flows/order-lifecycle.md`
- `diagrams/flows/pricing-calculation.md`

### ADRs Futuros (si aplica)
- `adr/ADR-010-*.md` (si hay decisiones en desarrollo)
- `adr/ADR-011-*.md` (etc.)

---

## ✅ CHECKLIST: ESTÁS LISTO

- [ ] Leí README.md (proyecto general)
- [ ] Leí SESSION_PROMPT.md o SESSION_PROMPT_QUICK.md (contexto)
- [ ] Revisé ADRs relevantes (arquitectura)
- [ ] Tengo este documento como referencia (navegación)
- [ ] Leí última sesión completada
- [ ] Leí PROJECT_STATUS.md (estado actual)
- [ ] Tengo COPILOT_INSTRUCTIONS.md a mano (si uso LLM)

---

## 🚀 PRÓXIMO PASO

**Dependiendo de tu rol:**

| Eres | Comienza por | Tiempo |
|-----|--------------|--------|
| **Desarrollador Backend** | SESSION_PROMPT.md § 1-3 | 20 min |
| **Desarrollador Frontend** | SESSION_PROMPT.md § 1-3 | 20 min |
| **Tech Lead** | DOCUMENTO-CONSOLIDADO-3.0.md | 60 min |
| **GitHub Copilot** | SESSION_PROMPT.md + COPILOT_INSTRUCTIONS.md | 20 min |
| **Auditor / Investigador** | Este archivo → ADRs → Sesiones | Variable |

---

## 📞 AYUDA RÁPIDA

**¿Cómo...?**

- ...Entender la arquitectura → ADR-002 + SESSION_PROMPT.md § 4
- ...Empezar a codificar → COPILOT_INSTRUCTIONS.md + SESSION_INIT_TEMPLATE.md
- ...Saber qué hacer esta sesión → Última sesión + PROJECT_STATUS.md
- ...Conocer el estado → PROJECT_STATUS.md
- ...Agregar nueva documentación → _ADR_TEMPLATE.md o _SESSION_TEMPLATE.md
- ...Ver todas las sesiones → Sessions / (carpeta)
- ...Entender decisiones pasadas → ADRs 001-009

---

**Última actualización:** 11/01/2026  
**Mantener sincronizado:** ADRs, Sessions, PROJECT_STATUS.md

