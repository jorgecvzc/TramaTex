# TramaTex - Convención de Nomenclatura de Ramas

**Versión:** 1.0  
**Fecha:** 2026-02-22  
**Referencia:** ADR-021 - Version Control & Branching Strategy

---

## 📍 Formato General de Ramas

```
<tipo>/<descripcion-corta-en-kebab-case>
```

**Reglas:**
- Todo en **minúsculas**
- Usar **guiones** (`-`) para separar palabras, **nunca** guiones bajos (`_`)
- **Máximo 50 caracteres** totales
- Descripción debe ser **clara y concreta**
- **No incluir** números de issue/ticket en el nombre (usar PR description)
- **No incluir** nombres de persona o máquina real

---

## 🏷️ Tipos de Ramas

### 1. `feature/*` - Nuevas Funcionalidades

**Propósito:** Desarrollo de nuevas características o mejoras funcionales.

**Formato:**
```
feature/<modulo>-<descripcion>
feature/<descripcion-general>
```

**Ejemplos:**
```bash
feature/pricing-volume-discounts
feature/product-bulk-import
feature/sales-invoice-templates
feature/ui-logic-validation
feature/advanced-search
feature/party-export-csv
```

**Cuándo usar:**
- ✅ Nueva funcionalidad completa
- ✅ Mejora significativa de feature existente
- ✅ Validación/testing que generará cambios funcionales
- ✅ Refactors grandes con cambios visibles

---

### 2. `bugfix/*` - Correcciones No Críticas

**Propósito:** Correcciones de bugs que no están en producción o no son críticos.

**Formato:**
```
bugfix/<modulo>-<descripcion-del-bug>
bugfix/<descripcion-del-bug>
```

**Ejemplos:**
```bash
bugfix/party-selector-empty-state
bugfix/pricing-calculation-rounding
bugfix/sales-form-validation
bugfix/product-variant-display
bugfix/mes-terminal-refresh
```

**Cuándo usar:**
- ✅ Bug encontrado en develop (no en producción)
- ✅ Corrección de comportamiento inesperado
- ✅ Fixes de UI/UX menores
- ✅ Validaciones faltantes

---

### 3. `hotfix/*` - Emergencias en Producción

**Propósito:** Correcciones críticas que deben aplicarse inmediatamente a producción.

**Formato:**
```
hotfix/<severidad>-<descripcion-critica>
hotfix/<descripcion-critica>
```

**Ejemplos:**
```bash
hotfix/security-jwt-validation
hotfix/critical-data-loss-sales
hotfix/pricing-wrong-calculation
hotfix/auth-bypass
hotfix/database-migration-error
```

**Cuándo usar:**
- ✅ Bug crítico en producción (master)
- ✅ Vulnerabilidad de seguridad
- ✅ Pérdida de datos
- ✅ Sistema no funcional
- ✅ Error que afecta a todos los usuarios

**⚠️ Importante:** Se crean desde `master`, no desde `develop`

---

### 4. `release/*` - Preparación de Releases

**Propósito:** Preparar una nueva versión para producción (bump version, changelog, últimos ajustes).

**Formato:**
```
release/v<MAJOR>.<MINOR>.<PATCH>
```

**Ejemplos:**
```bash
release/v1.1.0
release/v1.2.0
release/v2.0.0
release/v1.0.1
```

**Cuándo usar:**
- ✅ Cuando develop está listo para una nueva release
- ✅ Para actualizar versiones y CHANGELOG
- ✅ Para QA final antes de merge a master
- ✅ Para ajustes menores pre-release

**Flujo:**
1. Crear desde `develop`
2. Hacer últimos ajustes
3. Merge a `master` con tag
4. Back-merge a `develop`

---

### 5. `docs/*` - Solo Documentación (Opcional)

**Propósito:** Cambios exclusivos de documentación sin código funcional.

**Formato:**
```
docs/<tipo-documento>-<descripcion>
```

**Ejemplos:**
```bash
docs/adr-pricing-strategy
docs/guide-deployment
docs/api-swagger-update
docs/readme-improvement
```

**Cuándo usar:**
- ✅ Actualización de ADRs
- ✅ Guías de usuario/developer
- ✅ README improvements
- ✅ Documentación de API

**Nota:** Si los docs acompañan código, usar `feature/*` con commit `docs:`

---

### 6. `refactor/*` - Refactorizaciones Sin Cambio Funcional (Opcional)

**Propósito:** Mejoras de código sin cambiar funcionalidad externa.

**Formato:**
```
refactor/<modulo>-<descripcion>
```

**Ejemplos:**
```bash
refactor/pricing-service-simplification
refactor/product-variant-logic
refactor/sales-extract-validators
refactor/party-repository-optimization
```

**Cuándo usar:**
- ✅ Mejora de estructura de código
- ✅ Extracción de funciones/servicios
- ✅ Optimización de performance
- ✅ Reducción de complejidad ciclomática

**Alternativa:** Usar `feature/*` si el refactor es parte de una feature mayor

---

### 7. `test/*` - Mejoras de Testing (Opcional)

**Propósito:** Añadir o mejorar tests sin cambiar código funcional.

**Formato:**
```
test/<modulo>-<tipo-test>
```

**Ejemplos:**
```bash
test/pricing-coverage-increase
test/sales-integration-tests
test/product-e2e-playwright
test/party-unit-tests
```

**Cuándo usar:**
- ✅ Aumentar coverage sin cambiar lógica
- ✅ Añadir tests de integración
- ✅ Crear suite E2E
- ✅ Refactor de tests existentes

---

## 🎯 Casos Especiales

### Validación/QA que Genera Cambios
**Usar:** `feature/` con descripción de la validación

```bash
feature/ui-logic-validation      # Sprint de validación UI que generará fixes
feature/smoke-testing-fixes      # Testing que genera correcciones
feature/ux-review-improvements   # Review UX con mejoras
```

**Rationale:** La rama contendrá múltiples commits con fixes/improvements, más apropiado como feature.

---

### Múltiples Bugs del Mismo Tipo
**Usar:** `bugfix/` con descripción general

```bash
bugfix/sales-form-validations    # Varios bugs de validación en forms
bugfix/party-ui-fixes            # Varios bugs UI en Party module
bugfix/pricing-edge-cases        # Varios edge cases corregidos
```

---

### Experimentos/Spikes (No Recomendado para Merge)
**Usar:** `spike/` o `experiment/` (borrar tras validación)

```bash
spike/graphql-migration-poc
experiment/new-pricing-algorithm
```

**⚠️ Importante:** Estas ramas NO se mergean directamente. Se usan para validar y luego se crea una `feature/*` con la implementación final.

---

## 📋 Checklist de Validación de Nombre

Antes de crear una rama, verifica:

- [ ] ¿Usa un tipo válido? (`feature`, `bugfix`, `hotfix`, `release`)
- [ ] ¿Está en minúsculas?
- [ ] ¿Usa guiones (`-`) en lugar de guiones bajos (`_`)?
- [ ] ¿Es descriptivo y claro?
- [ ] ¿Tiene menos de 50 caracteres?
- [ ] ¿No incluye nombres de máquinas/personas reales?
- [ ] ¿No incluye números de issue? (van en PR description)

---

## 🔄 Flujo de Trabajo Completo

### Ejemplo: Feature Validation UI
```bash
# 1. Actualizar develop
git checkout develop
git pull origin develop

# 2. Crear rama con nombre correcto
git checkout -b feature/ui-logic-validation

# 3. Trabajar con commits convencionales
git add .
git commit -m "fix(sales): resolve invoice calculation in UI"
git commit -m "feat(product): add variant quick selector"
git commit -m "refactor(party): simplify selector logic"

# 4. Push y PR
git push origin feature/ui-logic-validation
# Abrir PR en GitHub hacia develop con descripción detallada

# 5. Tras merge, limpiar
git checkout develop
git pull origin develop
git branch -d feature/ui-logic-validation
```

---

## 🚫 Anti-Patrones (NO HACER)

### ❌ Nombres Vagos
```bash
feature/fixes              # ¿Qué fixes?
bugfix/stuff               # ¿Qué stuff?
feature/new-things         # ¿Qué things?
feature/wip                # Work in progress no es descriptivo
```

### ❌ Nombres con Datos Reales
```bash
feature/pcele-deployment   # pcele es máquina real
bugfix/jorge-issue         # nombre de persona
feature/client-acme-fix    # nombre de cliente real
```

### ❌ Formato Incorrecto
```bash
Feature/MyNewFeature       # Mayúsculas incorrectas
feature/my_new_feature     # Guiones bajos no permitidos
feature-my-new-feature     # Falta el separador /
FEATURE/my-new-feature     # Tipo en mayúsculas
```

### ❌ Demasiado Largos o Genéricos
```bash
feature/implement-the-new-advanced-pricing-algorithm-with-volume-discounts-and-seasonal-adjustments  # 104 caracteres
feature/update             # Demasiado genérico
```

---

## 📚 Referencia Rápida

| Tipo | Desde | Merge a | Ejemplo |
|------|-------|---------|---------|
| `feature/*` | develop | develop | `feature/ui-logic-validation` |
| `bugfix/*` | develop | develop | `bugfix/sales-form-validation` |
| `hotfix/*` | master | master + develop | `hotfix/security-critical` |
| `release/*` | develop | master + develop | `release/v1.1.0` |
| `docs/*` | develop | develop | `docs/adr-testing-strategy` |
| `refactor/*` | develop | develop | `refactor/pricing-simplify` |
| `test/*` | develop | develop | `test/sales-coverage-boost` |

---

## 🔗 Referencias

- **ADR-021:** Version Control & Branching Strategy
- **Version Control Guide:** Flujos de trabajo detallados
- **Conventional Commits:** https://www.conventionalcommits.org/

---

**Última Actualización:** 2026-02-22  
**Autor:** Equipo TramaTex
