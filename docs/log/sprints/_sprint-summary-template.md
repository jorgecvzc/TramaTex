# Resumen del Sprint [XX]

> **Nota de Nomenclatura:** Este archivo debe guardarse como `README.md` dentro del directorio `docs/log/sprints/sprint-[XX]/`.

---

## 📋 INFORMACIÓN DEL SPRINT

| Campo | Valor |
|-------|-------|
| **ID de Sprint** | XX |
| **Título** | [Nombre/Tema del Sprint] |
| **Fecha de Inicio** | YYYY-MM-DD |
| **Fecha de Fin** | YYYY-MM-DD |
| **Duración** | X días/semanas |
| **Objetivo del Sprint** | [Descripción breve del objetivo principal] |

---

## 📝 TAREAS COMPLETADAS

| ID | Título | Estado | Duración | Enlace |
|----|--------|--------|----------|--------|
| XX-01 | [Nombre de la tarea] | ✅ Completado | X horas | [XX-nombre-tarea.md](./XX-nombre-tarea.md) |

**Total de tareas:** X completadas

---

## 📊 MÉTRICAS AGREGADAS

### Tests

| Capa/Módulo | Tests Pasando | Cobertura | Estado |
|-------------|---------------|-----------|--------|
| Dominio | XX/XX | XX% | ✅ |
| Persistencia | XX/XX | XX% | ✅ |
| Aplicación | XX/XX | XX% | ✅ |
| Interfaces | XX/XX | XX% | ✅ |
| **TOTAL** | **XX/XX** | **XX%** | ✅ |

### Código

| Métrica | Valor |
|---------|-------|
| **Archivos Creados** | XX |
| **Archivos Modificados** | XX |
| **Líneas de Código Agregadas** | ~XXX |
| **Líneas de Tests Agregadas** | ~XXX |
| **Commits Totales** | XX |

### Tiempo

| Métrica | Valor |
|---------|-------|
| **Horas Estimadas** | XX horas |
| **Horas Reales** | XX horas |
| **Variación** | +/-X% |

---

## 🎯 LOGROS PRINCIPALES

### Funcionalidad Implementada

1. **[Funcionalidad 1]**
   - [Detalle específico]

### Mejoras Técnicas

- ✅ [Mejora técnica 1]

### Decisiones Arquitectónicas

- **[ADR-XXX]**: [Descripción breve de la decisión]
- **[Patrón/Práctica]**: [Descripción breve]

---

## 🏗️ ARQUITECTURA Y PATRONES

### Capas Implementadas

```
┌─────────────────────────────────┐
│  [Capa 4]                       │ ← [Estado: Completo/Parcial/Pendiente]
├─────────────────────────────────┤
│  [Capa 3]                       │ ← [Estado]
├─────────────────────────────────┤
│  [Capa 2]                       │ ← [Estado]
├─────────────────────────────────┤
│  [Capa 1]                       │ ← [Estado]
└─────────────────────────────────┘
```

### Patrones de Diseño Aplicados

1. **[Patrón 1]**: [Dónde se aplicó y por qué]

---

## 🚨 PROBLEMAS Y SOLUCIONES

### Bloqueadores Superados

| Problema | Impacto | Solución | Tiempo Perdido |
|----------|---------|----------|----------------|
| [Problema 1] | Alto/Medio/Bajo | [Cómo se resolvió] | X horas |

### Deuda Técnica Identificada

- [ ] **[Deuda 1]**: [Descripción] → [Referencia/Ticket]

### Riesgos Aceptados

- **[Riesgo 1]**: [Descripción y justificación]

---

## 📚 APRENDIZAJES

### Técnicos

```
[Lecciones aprendidas sobre tecnología, patrones, prácticas de código]
```

### De Proceso

```
[Lecciones sobre planificación, estimación, colaboración]
```

### Mejores Prácticas Identificadas

- ✅ [Práctica 1]

---

## 📁 ARCHIVOS CREADOS/MODIFICADOS

### Nuevos Archivos

**Backend:**
```
apps/tramatex-api/internal/[modulo]/
├── domain/
│   ├── [archivo1.go]
```

**Frontend:**
```
apps/frontend/src/
├── components/
│   └── [componente1.vue]
```

### Modificaciones Importantes

| Archivo | Tipo de Cambio | Descripción |
|---------|----------------|-------------|
| [archivo1] | REFACTOR | [Descripción] |

---

## 🔗 INTEGRACIÓN CON OTROS MÓDULOS

### Dependencias

- **[Módulo A]**: [Qué se consume o expone]

### Contratos Definidos

- **API Endpoints**: [Cantidad de endpoints nuevos]
- **Eventos de Dominio**: [Si aplica]
- **Interfaces**: [Interfaces públicas]

---

## ✅ DEFINICIÓN DE "HECHO" - VERIFICACIÓN

Sprint completado cuando:

- [x] Todas las tareas del sprint están completadas
- [x] Todos los tests pasan: `go test ./...` y `npm run test`
- [x] Linters sin warnings: `golangci-lint` y `npm run lint`
- [x] Cobertura de tests ≥85% en módulos críticos
- [x] Docker Compose levanta sin errores
- [x] Funcionalidad demostrable en ambiente local
- [x] Documentación actualizada (tareas + project-status.md)
- [x] ADRs creados/actualizados si aplica
- [x] Commits bien organizados y descriptivos

---

## 🚀 PRÓXIMOS PASOS

### Sprint Siguiente

**Objetivo del próximo sprint:**
[Descripción breve del siguiente objetivo]

**Tareas planificadas:**

1. [ ] **Tarea [XX+1]-01**: [Título] (Estimación: X horas)

### Prerequisitos

- [ ] [Prerequisito 1]

### Contexto a Cargar

Para el próximo sprint, cargar:
- `agents/project/context/[contexto-relevante].yaml`
- `agents/project/context/[otro-contexto].yaml`

---

## 📊 ESTADO DEL PROYECTO

### Progreso del MVP

```
Fase Actual: [Nombre de la Fase]
Porcentaje Completado: XX%

Fases:
✅ Fase 0: Fundaciones Técnicas
🔄 Fase 1: Dominio Base (XX% completo)
⏳ Fase 2: Catálogo y Pricing
⏳ Fase 3: Pedidos y Producción
```

### Horas Invertidas

| Concepto | Horas | Porcentaje |
|----------|-------|------------|
| **Horas Sprint** | XX | X% |
| **Horas Acumuladas** | XXX | XX% |
| **Horas Totales Estimadas** | 782 | 100% |
| **Horas Restantes** | XXX | XX% |

---

## 📝 NOTAS FINALES

```
[Cualquier observación adicional, contexto importante para futuros sprints,
recomendaciones para el equipo, etc.]
```

---

## ✍️ FIRMA

**Sprint completado:** [YYYY-MM-DD]

**Facilitador:** [Nombre]  
**LLM Principal:** [GitHub Copilot / Claude / Gemini]  
**Revisor:** [Nombre si aplica]

---

**Referencia a tareas:**
