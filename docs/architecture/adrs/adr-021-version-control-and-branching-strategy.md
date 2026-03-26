# ADR-021 – Version Control & Branching Strategy

**Fecha:** 2026-02-22  
**Estado:** Aceptado  
**Autores:** Equipo de Desarrollo TramaTex  

---

## 1. Contexto

Con el MVP completado y listo para producción (v1.0.0), es necesario establecer una estrategia clara de control de versiones y branching que permita:

- **Desarrollo continuo** de features Post-MVP sin afectar la versión estable
- **Versionado semántico** claro y consistente
- **Integración controlada** de cambios con revisiones
- **Despliegues seguros** a producción

El proyecto TramaTex es un ERP modular con 6 módulos funcionales (Party, Product, Pricing, Sales, IAM, MES), desarrollado por un equipo pequeño con necesidad de agilidad pero también de estabilidad en producción.

---

## 2. Decisión Adoptada

### 2.1 Estrategia de Branching: **GitFlow Simplificado**

Adoptamos una variante simplificada de GitFlow adaptada a equipos pequeños:

#### Ramas Principales

1. **`main`** (Producción)
   - Código estable, siempre desplegable
   - Solo recibe merges desde `develop` o hotfixes
   - Cada merge representa una release
   - **Protegida:** Requiere Pull Request aprobado y tests pasando

2. **`develop`** (Integración)
   - Rama de desarrollo activo
   - Base para todas las feature branches
   - Integra cambios antes de pasar a `main`
   - **Protegida:** Requiere Pull Request y tests pasando

#### Ramas Temporales

3. **`feature/*`** (Features)
   - Formato: `feature/descripcion-corta`
   - Se crean desde `develop`
   - Se mergen de vuelta a `develop` vía PR
   - Ejemplos:
     - `feature/advanced-pricing-rules`
     - `feature/mes-analytics-dashboard`
     - `feature/stock-management`

4. **`bugfix/*`** (Correcciones)
   - Formato: `bugfix/descripcion-del-bug`
   - Se crean desde `develop`
   - Se mergen a `develop` vía PR
   - Ejemplos:
     - `bugfix/party-selector-crash`
     - `bugfix/pricing-calculation-error`

5. **`hotfix/*`** (Emergencias en Producción)
   - Formato: `hotfix/descripcion-critica`
   - Se crean desde `main`
   - Se mergen a `main` Y `develop`
   - Solo para bugs críticos en producción
   - Ejemplos:
     - `hotfix/security-jwt-validation`
     - `hotfix/data-corruption-sales`

6. **`release/*`** (Pre-releases)
   - Formato: `release/v1.1.0`
   - Se crean desde `develop` cuando está lista una nueva versión
   - Se hacen ajustes finales (bumping versión, changelog)
   - Se merge a `main` con tag, luego back-merge a `develop`

### 2.2 Versionado Semántico (SemVer 2.0)

Formato: **`MAJOR.MINOR.PATCH`**

- **MAJOR** (1.x.x → 2.x.x): Cambios rompientes en API o arquitectura
  - Ejemplos: Cambio de estructura de BD incompatible, refactor de API pública
  
- **MINOR** (1.0.x → 1.1.x): Nuevas funcionalidades compatibles
  - Ejemplos: Nuevo módulo, nuevos endpoints, features Post-MVP
  
- **PATCH** (1.0.0 → 1.0.1): Correcciones de bugs
  - Ejemplos: Bugfixes, optimizaciones de rendimiento, mejoras de tests

#### Estado Inicial del Proyecto

- **v1.0.0** (2026-02-22): MVP Completado
  - 6 módulos funcionales end-to-end
  - Backend Clean Architecture + DDD
  - Frontend Vue 3 + TypeScript
  - Testing: Backend 75%+, Frontend 77.63%
  - Docker Compose para desarrollo
  - Criterios de aceptación MVP cumplidos

### 2.3 Flujo de Trabajo

#### Desarrollo de Features

```bash
# 1. Crear branch desde develop
git checkout develop
git pull origin develop
git checkout -b feature/nueva-funcionalidad

# 2. Desarrollar con commits descriptivos
git add .
git commit -m "feat(module): add new functionality"

# 3. Push y crear Pull Request
git push origin feature/nueva-funcionalidad
# Abrir PR en GitHub hacia develop

# 4. Revisión y merge (requiere aprobación)
# Tests automáticos deben pasar
# Merge squash opcional para mantener historial limpio
```

#### Release a Producción

```bash
# 1. Crear release branch desde develop
git checkout develop
git pull origin develop
git checkout -b release/v1.1.0

# 2. Bump version y actualizar CHANGELOG
# Editar VERSION file, package.json, etc.

# 3. Merge a main con tag
git checkout main
git merge --no-ff release/v1.1.0
git tag -a v1.1.0 -m "Release v1.1.0: Advanced Pricing Rules"
git push origin main --tags

# 4. Back-merge a develop
git checkout develop
git merge --no-ff release/v1.1.0
git push origin develop

# 5. Eliminar release branch
git branch -d release/v1.1.0
```

#### Hotfix Crítico

```bash
# 1. Crear desde main
git checkout main
git pull origin main
git checkout -b hotfix/security-critical

# 2. Fix y commit
git commit -m "fix(security): patch JWT validation vulnerability"

# 3. Merge a main con tag PATCH
git checkout main
git merge --no-ff hotfix/security-critical
git tag -a v1.0.1 -m "Hotfix v1.0.1: Security patch"
git push origin main --tags

# 4. Merge a develop
git checkout develop
git merge --no-ff hotfix/security-critical
git push origin develop
```

### 2.4 Protección de Archivos Sensibles

#### Regla: **Ningún archivo `.env` se sube al repositorio**

**Archivos excluidos en `.gitignore`:**
```gitignore
# Environment variables
.env
.env.*
!.env.example
```

**Archivos permitidos:**
- `.env.example` (plantillas asépticas sin secretos)

**Archivos sensibles actuales a remover:**
- `docker/.env.remote` (actualmente trackeado, debe eliminarse del historial si contiene secretos)

### 2.5 Estructura de .env.example

Cada `.env.example` debe contener:

```bash
# Example values (NO SECRETS)
# Copy to .env and fill with real values

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tramatex_dev
DB_USER=your_user
DB_PASSWORD=your_secure_password

# JWT
JWT_SECRET=generate_a_secure_random_string_here

# Server
API_PORT=8080
FRONTEND_PORT=5173
```

---

## 3. Consecuencias

### Positivas

1. ✅ **Estabilidad en `main`**: Producción siempre tiene código estable y probado
2. ✅ **Desarrollo continuo**: `develop` permite trabajar en múltiples features sin bloqueos
3. ✅ **Revisión de código**: PRs obligatorios garantizan calidad
4. ✅ **Versionado claro**: SemVer facilita entender el impacto de cada release
5. ✅ **Historial limpio**: Tags marcan releases importantes
6. ✅ **Seguridad**: .env no se filtra al repositorio

### Negativas / Mitigaciones

1. ⚠️ **Complejidad adicional** vs GitHub Flow simple
   - **Mitigación**: Para equipo pequeño, GitFlow simplificado reduce burocracia
   
2. ⚠️ **Posibles conflictos de merge** entre `develop` y `main`
   - **Mitigación**: Releases frecuentes reducen divergencia
   
3. ⚠️ **Disciplina requerida** para seguir flujo
   - **Mitigación**: Documentación clara y protecciones de rama en GitHub

---

## 4. Automatización (GitHub Actions)

### CI/CD Pipeline

#### En PRs hacia `develop`
```yaml
# .github/workflows/ci.yml
- Run backend tests (Go)
- Run frontend tests (Vitest)
- Check code coverage (fail if < 70%)
- Lint checks
```

#### En merge a `main` (con tag)
```yaml
# .github/workflows/release.yml
- Build backend binary
- Build frontend static files
- Create GitHub Release con changelog
- (Opcional) Deploy automático a staging/producción
```

---

## 5. Convenciones de Commits

### Formato: Conventional Commits

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Tipos comunes:**
- `feat`: Nueva funcionalidad
- `fix`: Corrección de bug
- `docs`: Documentación
- `refactor`: Refactorización sin cambio funcional
- `test`: Añadir o mejorar tests
- `chore`: Cambios de build, config, etc.

**Scopes por módulo:**
- `party`, `product`, `pricing`, `sales`, `iam`, `mes`
- `frontend`, `backend`, `infra`, `docs`

**Ejemplos:**
```bash
feat(pricing): add volume discount rules
fix(party): resolve selector crash on empty results
docs(adr): add versioning strategy ADR-021
refactor(product): extract variant generation logic
test(sales): increase coverage to 80%
chore(deps): update Vue to 3.4.20
```

---

## 6. Migración Inicial

### Acciones Inmediatas (v1.0.0 Setup)

1. ✅ Crear este ADR (ADR-021)
2. 🔄 Actualizar `.gitignore` global y por carpeta
3. 🔄 Crear `.env.example` para backend y frontend
4. 🔄 Eliminar archivos temporales de coverage del tracking
5. 🔄 Commit todos los cambios del MVP en `main`
6. 🔄 Crear tag `v1.0.0` con mensaje descriptivo
7. 🔄 Crear rama `develop` desde `main`
8. 🔄 Configurar branch protection rules en GitHub
9. 🔄 Push `main`, `develop` y tags a GitHub
10. 🔄 Actualizar README.md con instrucciones de contribución

### Estado de Migración
```
Estado actual: master (legacy)
Estado objetivo: main + develop (v1.0.0)

Acciones:
- Renombrar master → main (o mantener master como main)
- Crear develop
- Establecer v1.0.0 como baseline
```

---

## 7. Referencias

- **SemVer 2.0:** https://semver.org/
- **GitFlow Original:** https://nvie.com/posts/a-successful-git-branching-model/
- **GitHub Flow:** https://docs.github.com/en/get-started/using-github/github-flow
- **Conventional Commits:** https://www.conventionalcommits.org/
- **ADR-004:** MVP Development Lifecycle
- **ADR-008:** MVP Timeline Planning

---

## 8. Aprobación y Vigencia

**Estado:** ✅ Aceptado  
**Fecha de Aplicación:** 2026-02-22 (v1.0.0)  
**Revisión:** Cada 6 meses o ante cambios mayores en el equipo  

---

_Última Actualización: 2026-02-22_
