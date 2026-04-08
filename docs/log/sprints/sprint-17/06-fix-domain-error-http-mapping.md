# Tarea 17-06: Fix Mapeo DomainError → HTTP Status

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 06 |
| **ID de Sprint** | sprint-17 |
| **Título** | Implementar mapeo correcto de DomainError.Code a HTTP Status en handlers |
| **Estado** | ⏳ Planificado |
| **Facilitador/LLM** | GitHub Copilot / Claude Opus 4.6 |
| **Fecha de Inicio** | pendiente |
| **Fecha de Fin** | pendiente |
| **Duración Estimada** | 1 hora |
| **Prioridad** | P2 — Medio |
| **Dependencia** | Independiente (paralelizable) |

---

## 🎯 OBJETIVOS PRINCIPALES

1. [ ] **Implementar función `mapDomainError`** en handlers de Pricing
2. [ ] **Mapear**: `VALIDATION_ERROR → 400`, `NOT_FOUND → 404`, `CONFLICT → 409`, `RULE_ERROR → 422`
3. [ ] **Aplicar** a `PricingEngineHandler` (y `PricingHandler` si aún existe en ese momento)

---

## 📊 CONTEXTO

### Estado actual

Ambos handlers devuelven `http.StatusBadRequest` (400) para **todos** los errores:
```go
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
```

### Requisito de implementation-guide.md

```
ErrCodeValidation → 400 Bad Request
ErrCodeNotFound   → 404 Not Found
ErrCodeConflict   → 409 Conflict
```

### Referencia: middleware global de errores (sprint-15)

En sprint-15 se implementó un middleware global de errores en `shared/infrastructure/middleware/`. Verificar si ya mapea DomainError y si Pricing puede reutilizarlo en lugar de implementar mapeo local.

---

## 🛠️ PLAN DE TRABAJO

### Paso 1: Verificar middleware global existente
- ¿`shared/infrastructure/middleware/` ya maneja `domain.DomainError`?
- Si sí → verificar que Pricing lo usa (handlers deben propagar error, no capturarlo)

### Paso 2: Si middleware no cubre Pricing
- Crear helper `mapDomainErrorToHTTP(err error) int` en interfaces/http/handler/
- Type-assert `*domain.DomainError` y mapear Code → HTTP status

### Paso 3: Aplicar a todos los handlers
- Reemplazar `c.JSON(400, ...)` por mapeo correcto

### Paso 4: Tests
- Test: error NOT_FOUND retorna 404
- Test: error VALIDATION retorna 400
- Test: error CONFLICT retorna 409
- Test: error desconocido retorna 500
