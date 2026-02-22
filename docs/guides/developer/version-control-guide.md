# TramaTex - Guía de Control de Versiones

**Versión:** 1.0  
**Fecha:** 2026-02-22  
**Referencia:** ADR-021

---

## 🏷️ Versionado Actual

**Release Actual:** `v1.0.0` (MVP Completado - Producción)

### Versiones Disponibles

- **v1.0.0** (2026-02-22): MVP Completado
  - 6 módulos funcionales (Party, Product, Pricing, Sales, IAM, MES)
  - Backend: 75%+ coverage, Frontend: 77.63% coverage
  - 60,000+ líneas de código
  - Listo para producción

---

## 🌿 Estrategia de Branches

### Ramas Principales

#### `master` (Producción)
- Código estable y desplegable
- Solo recibe merges desde `develop` o `hotfix/*`
- **Protegida:** Requiere PR aprobado + tests pasando
- Cada merge representa una release con tag `vX.Y.Z`

#### `develop` (Integración)
- Rama de desarrollo activo
- Base para todas las feature/bugfix branches
- **Protegida:** Requiere PR + tests pasando

### Ramas Temporales

- **`feature/*`**: Nuevas funcionalidades (desde `develop`)
- **`bugfix/*`**: Correcciones no críticas (desde `develop`)
- **`hotfix/*`**: Emergencias en producción (desde `master`)
- **`release/*`**: Preparación de releases (desde `develop`)

---

## 🔄 Flujos de Trabajo

### Desarrollar una Feature

```bash
# 1. Actualizar develop local
git checkout develop
git pull origin develop

# 2. Crear branch de feature
git checkout -b feature/mi-funcionalidad

# 3. Desarrollar y commitear con Conventional Commits
git add .
git commit -m "feat(module): add new functionality"

# 4. Push y abrir Pull Request en GitHub
git push origin feature/mi-funcionalidad
# Abrir PR hacia develop en GitHub

# 5. Tras aprobación y merge, limpiar
git checkout develop
git pull origin develop
git branch -d feature/mi-funcionalidad
```

### Crear una Release

```bash
# 1. Crear release branch desde develop
git checkout develop
git pull origin develop
git checkout -b release/v1.1.0

# 2. Actualizar versiones y CHANGELOG
# Editar archivos necesarios

# 3. Commit de preparación
git commit -m "chore(release): prepare v1.1.0"

# 4. Merge a master con tag
git checkout master
git merge --no-ff release/v1.1.0
git tag -a v1.1.0 -m "Release v1.1.0: [Descripción]"
git push origin master --tags

# 5. Back-merge a develop
git checkout develop
git merge --no-ff release/v1.1.0
git push origin develop

# 6. Eliminar release branch
git branch -d release/v1.1.0
git push origin --delete release/v1.1.0
```

### Hotfix Crítico

```bash
# 1. Crear desde master
git checkout master
git pull origin master
git checkout -b hotfix/descripcion-critica

# 2. Fix y commit
git commit -m "fix(security): patch critical vulnerability"

# 3. Merge a master con tag PATCH
git checkout master
git merge --no-ff hotfix/descripcion-critica
git tag -a v1.0.1 -m "Hotfix v1.0.1: Security patch"
git push origin master --tags

# 4. Merge a develop
git checkout develop
git merge --no-ff hotfix/descripcion-critica
git push origin develop

# 5. Eliminar hotfix branch
git branch -d hotfix/descripcion-critica
```

---

## 📝 Conventional Commits

### Formato

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Tipos

- **feat**: Nueva funcionalidad (incrementa MINOR)
- **fix**: Corrección de bug (incrementa PATCH)
- **docs**: Documentación
- **refactor**: Refactorización sin cambio funcional
- **test**: Añadir o mejorar tests
- **chore**: Cambios de build, config, deps
- **perf**: Mejoras de rendimiento
- **style**: Formato de código (no lógica)
- **ci**: Cambios en CI/CD

### Scopes

Por módulo: `party`, `product`, `pricing`, `sales`, `iam`, `mes`  
Por capa: `frontend`, `backend`, `infra`, `docs`

### Ejemplos

```bash
# Feature (MINOR bump)
feat(pricing): add volume discount rules for bulk orders

# Fix (PATCH bump)
fix(party): resolve selector crash when no results

# Breaking change (MAJOR bump)
feat(product)!: refactor variant generation API

BREAKING CHANGE: The variant generation endpoint now requires
a new 'strategy' parameter. Old API is deprecated.

# Documentation
docs(adr): add versioning strategy ADR-021

# Refactor
refactor(sales): extract order validation to separate service

# Test
test(pricing): increase coverage to 90%

# Chore
chore(deps): update Vue to 3.4.20
```

---

## 🔐 Variables de Entorno

### Archivos NO Permitidos en Git

❌ `.env`  
❌ `.env.development`  
❌ `.env.production`  
❌ `.env.*.local` (ej: `.env.development.local`)  
❌ `.env.*.remote` (ej: `.env.pcele.remote`, `.env.staging.remote`)

### Archivos Permitidos (Plantillas)

✅ `.env.example`  
✅ `.env.*.example` (ej: `.env.pcele.remote.example`)

### Configuración Inicial

```bash
# Backend
cp apps/tramatex-api/.env.example apps/tramatex-api/.env
# Editar con valores reales

# Frontend
cp apps/frontend/.env.example apps/frontend/.env.development
# Editar con valores reales

# Docker
cp docker/.env.example docker/.env
# Editar con valores reales

# Remoto (ejemplo: máquina pcele)
cp .env.pcele.remote.example .env.pcele.remote
# Editar con credenciales del servidor remoto
```

**⚠️ NUNCA commitear archivos `.env` con credenciales reales**

---

## 🏷️ Versionado Semántico (SemVer 2.0)

### Formato: `MAJOR.MINOR.PATCH`

- **MAJOR** (1.x.x → 2.x.x): Cambios rompientes (breaking changes)
  - Ejemplos: Cambio de API incompatible, refactor de BD mayor
  
- **MINOR** (1.0.x → 1.1.x): Nuevas funcionalidades compatibles
  - Ejemplos: Nuevo módulo, nuevos endpoints, features Post-MVP
  
- **PATCH** (1.0.0 → 1.0.1): Correcciones de bugs
  - Ejemplos: Bugfixes, optimizaciones, hotfixes

### Pre-releases (Opcional Post-MVP)

- `v1.1.0-alpha.1`: Versión alpha
- `v1.1.0-beta.2`: Versión beta
- `v1.1.0-rc.1`: Release candidate

---

## 🛡️ Protecciones de Rama en GitHub

### Configuración Recomendada

#### Rama `master`
- ✅ Require pull request before merging
- ✅ Require approvals (mínimo 1)  
- ✅ Require status checks to pass
  - Backend tests
  - Frontend tests
  - Code coverage check
- ✅ Require conversation resolution before merging
- ✅ Do not allow bypassing the above settings

#### Rama `develop`
- ✅ Require pull request before merging
- ✅ Require status checks to pass
  - Backend tests
  - Frontend tests
- ✅ Allow force pushes (solo para emergency fixes)

---

## 📊 Estado del Repositorio

### Ramas Actuales

```
* master (producción) → v1.0.0
* develop (desarrollo activo)
```

### Tags

```
v1.0.0 (2026-02-22) - MVP Completado
```

### Remoto

```
origin: git@github.com:jorgecvzc/TramaTex.git
```

---

## 🚀 Próximos Pasos (Post-MVP)

1. **Configurar Branch Protection Rules** en GitHub
2. **Configurar GitHub Actions** para CI/CD
3. **Primer feature Post-MVP** desde `develop`
4. **Release v1.1.0** con mejoras incrementales

---

## 📚 Referencias

- **ADR-021:** Version Control & Branching Strategy
- **SemVer 2.0:** https://semver.org/
- **Conventional Commits:** https://www.conventionalcommits.org/
- **GitFlow:** https://nvie.com/posts/a-successful-git-branching-model/

---

**Última Actualización:** 2026-02-22  
**Autor:** Equipo TramaTex
