# Tarea 02: Pipeline CI/CD Básico con GitHub Actions

---

## 📋 INFORMACIÓN DE LA TAREA

| Campo | Valor |
|-------|-------|
| **ID de Tarea** | 02 |
| **ID de Sprint** | sprint-04 |
| **Título** | Pipeline CI/CD con GitHub Actions |
| **Estado** | ⏳ Planificado |
| **Facilitador/LLM** | GitHub Copilot (Claude Sonnet 4.5) |
| **Fecha de Inicio** | (Por determinar) |
| **Fecha de Fin** | (Por determinar) |
| **Duración Estimada** | 4 horas |
| **Duración Real** | (Completar al finalizar) |

**Nota sobre IDs:**
- **ID de Tarea**: 02 (segunda tarea del sprint-04)
- **ID de Sprint**: sprint-04
- **ID Único**: 04-02

---

## 🎯 OBJETIVOS PRINCIPALES

Establecer un pipeline de CI/CD básico y automatizado que ejecute tests, linters y validaciones en cada pull request y merge, asegurando la calidad del código antes de integrarlo.

### Subtareas

1. [ ] **GitHub Actions Workflows** (2 horas)
   - [ ] Workflow para backend (Go): tests + linters + coverage
   - [ ] Workflow para frontend (Vue): tests + linters
   - [ ] Configurar triggers (push, pull_request)
   - [ ] Badges de status en README

2. [ ] **Pre-commit Hooks** (1 hora)
   - [ ] Configurar pre-commit framework
   - [ ] Hooks para formato de código (gofmt, prettier)
   - [ ] Hooks para tests rápidos
   - [ ] Script de instalación para equipo

3. [ ] **Linters y Formatters** (1 hora)
   - [ ] golangci-lint configurado para backend
   - [ ] eslint + prettier para frontend
   - [ ] Scripts en Makefile
   - [ ] Documentar estándares en CONTRIBUTING.md

---

## 📊 CONTEXTO DE ENTRADA

### Estado Actual

**Sin CI/CD:**
- ❌ Tests no se ejecutan automáticamente
- ❌ Sin validación de código en PRs
- ❌ Riesgo de commits con código roto
- ❌ Sin verificación de estándares

**Hallazgo de Auditoría OWASP:**
- A08:2021 - Software and Data Integrity Failures
- Severidad: MEDIA
- Sin pipeline automatizado de tests

**Tests disponibles:**
- ✅ Backend: 75 tests (100% coverage en módulos core)
- ✅ Frontend: Tests con Vitest configurados
- ✅ Scripts en Makefile: `make test`, `make lint`

---

## 🛠️ PLAN DE TRABAJO

### Fase 1: GitHub Actions Workflows

#### 1.1 Workflow Backend (Go)

**Archivo a crear:** `.github/workflows/backend.yml`

```yaml
name: Backend CI

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'apps/tramatex-api/**'
      - '.github/workflows/backend.yml'
  pull_request:
    branches: [ main, develop ]
    paths:
      - 'apps/tramatex-api/**'

jobs:
  test:
    name: Test Backend
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: tramatex_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Cache Go modules
        uses: actions/cache@v4
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-
      
      - name: Install dependencies
        working-directory: apps/tramatex-api
        run: go mod download
      
      - name: Run tests
        working-directory: apps/tramatex-api
        run: go test -v -race -coverprofile=coverage.out ./...
        env:
          DATABASE_URL: postgresql://test:test@localhost:5432/tramatex_test?sslmode=disable
      
      - name: Coverage report
        working-directory: apps/tramatex-api
        run: go tool cover -func=coverage.out
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: ./apps/tramatex-api/coverage.out
          flags: backend
          name: backend-coverage

  lint:
    name: Lint Backend
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          working-directory: apps/tramatex-api
          args: --timeout=5m

  build:
    name: Build Backend
    runs-on: ubuntu-latest
    needs: [test, lint]
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Build
        working-directory: apps/tramatex-api
        run: go build -v ./cmd/api
```

#### 1.2 Workflow Frontend (Vue)

**Archivo a crear:** `.github/workflows/frontend.yml`

```yaml
name: Frontend CI

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'apps/frontend/**'
      - '.github/workflows/frontend.yml'
  pull_request:
    branches: [ main, develop ]
    paths:
      - 'apps/frontend/**'

jobs:
  test:
    name: Test Frontend
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: apps/frontend/package-lock.json
      
      - name: Install dependencies
        working-directory: apps/frontend
        run: npm ci
      
      - name: Run tests
        working-directory: apps/frontend
        run: npm run test:unit -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: ./apps/frontend/coverage/coverage-final.json
          flags: frontend
          name: frontend-coverage

  lint:
    name: Lint Frontend
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: apps/frontend/package-lock.json
      
      - name: Install dependencies
        working-directory: apps/frontend
        run: npm ci
      
      - name: Run ESLint
        working-directory: apps/frontend
        run: npm run lint

  build:
    name: Build Frontend
    runs-on: ubuntu-latest
    needs: [test, lint]
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: apps/frontend/package-lock.json
      
      - name: Install dependencies
        working-directory: apps/frontend
        run: npm ci
      
      - name: Build
        working-directory: apps/frontend
        run: npm run build
```

#### 1.3 Badges en README

**Actualizar:** `README.md` (raíz del proyecto)

```markdown
# TramaTex

![Backend CI](https://github.com/[usuario]/TramaTex/workflows/Backend%20CI/badge.svg)
![Frontend CI](https://github.com/[usuario]/TramaTex/workflows/Frontend%20CI/badge.svg)
[![codecov](https://codecov.io/gh/[usuario]/TramaTex/branch/main/graph/badge.svg)](https://codecov.io/gh/[usuario]/TramaTex)

Sistema ERP/MES para microempresas del sector textil.
```

---

### Fase 2: Pre-commit Hooks

#### 2.1 Configuración Pre-commit

**Archivo a crear:** `.pre-commit-config.yaml` (raíz del proyecto)

```yaml
repos:
  # Backend (Go)
  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
        args: [-w]
        files: ^apps/tramatex-api/
      - id: go-vet
        files: ^apps/tramatex-api/
      - id: go-imports
        args: [-w]
        files: ^apps/tramatex-api/
      - id: golangci-lint
        args: [--timeout=5m]
        files: ^apps/tramatex-api/

  # Frontend (JavaScript/Vue)
  - repo: https://github.com/pre-commit/mirrors-prettier
    rev: v3.1.0
    hooks:
      - id: prettier
        files: ^apps/frontend/
        types_or: [javascript, vue, css, json, markdown]

  - repo: https://github.com/pre-commit/mirrors-eslint
    rev: v8.56.0
    hooks:
      - id: eslint
        files: ^apps/frontend/.*\.[jt]sx?$
        types: [file]
        additional_dependencies:
          - eslint@^8.56.0
          - eslint-plugin-vue@^9.19.2

  # General
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
        args: ['--maxkb=1000']
      - id: check-merge-conflict
```

#### 2.2 Script de Instalación

**Archivo a crear:** `scripts/setup-pre-commit.sh`

```bash
#!/bin/bash

echo "🔧 Configurando pre-commit hooks..."

# Instalar pre-commit (si no está instalado)
if ! command -v pre-commit &> /dev/null; then
    echo "📦 Instalando pre-commit..."
    pip install pre-commit
fi

# Instalar hooks
echo "⚙️  Instalando hooks..."
pre-commit install

# Ejecutar contra todos los archivos (primera vez)
echo "✅ Ejecutando validación inicial..."
pre-commit run --all-files

echo "✨ ¡Pre-commit configurado exitosamente!"
```

**PowerShell para Windows:** `scripts/setup-pre-commit.ps1`

```powershell
Write-Host "🔧 Configurando pre-commit hooks..." -ForegroundColor Green

# Verificar Python
if (-not (Get-Command python -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Python no encontrado. Instálalo primero." -ForegroundColor Red
    exit 1
}

# Instalar pre-commit
Write-Host "📦 Instalando pre-commit..." -ForegroundColor Yellow
pip install pre-commit

# Instalar hooks
Write-Host "⚙️  Instalando hooks..." -ForegroundColor Yellow
pre-commit install

# Ejecutar validación inicial
Write-Host "✅ Ejecutando validación inicial..." -ForegroundColor Yellow
pre-commit run --all-files

Write-Host "✨ ¡Pre-commit configurado exitosamente!" -ForegroundColor Green
```

---

### Fase 3: Linters y Formatters

#### 3.1 golangci-lint (Backend)

**Archivo a crear:** `apps/tramatex-api/.golangci.yml`

```yaml
run:
  timeout: 5m
  tests: true
  modules-download-mode: readonly

linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck
    - unconvert
    - unparam
    - misspell
    - revive
    - gocyclo
    - goconst

linters-settings:
  gocyclo:
    min-complexity: 15
  goconst:
    min-len: 3
    min-occurrences: 3
  revive:
    rules:
      - name: exported
        disabled: false

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

**Agregar a Makefile:**

```makefile
.PHONY: lint
lint:
	@echo "🔍 Ejecutando linters..."
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	@echo "🔧 Arreglando issues de linting..."
	golangci-lint run --fix ./...
```

#### 3.2 ESLint + Prettier (Frontend)

**Archivo a actualizar:** `apps/frontend/.eslintrc.cjs`

```javascript
module.exports = {
  root: true,
  env: {
    node: true,
    'vue/setup-compiler-macros': true
  },
  extends: [
    'plugin:vue/vue3-recommended',
    'eslint:recommended',
    '@vue/eslint-config-prettier'
  ],
  parserOptions: {
    ecmaVersion: 'latest'
  },
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'vue/multi-word-component-names': 'off'
  }
}
```

**Archivo a crear:** `apps/frontend/.prettierrc`

```json
{
  "semi": false,
  "singleQuote": true,
  "trailingComma": "none",
  "printWidth": 100,
  "tabWidth": 2
}
```

**Agregar scripts a package.json:**

```json
{
  "scripts": {
    "lint": "eslint . --ext .vue,.js,.jsx,.cjs,.mjs --fix",
    "format": "prettier --write src/"
  }
}
```

---

## 🚨 BLOQUEADORES/PROBLEMAS ENCONTRADOS

*(Se actualizará durante la implementación)*

**Posibles Issues:**
1. Configuración de secrets en GitHub (CODECOV_TOKEN)
2. Compatibilidad de pre-commit en Windows
3. Cache de dependencias en GitHub Actions

---

## 🎓 APRENDIZAJES Y NOTAS

### Decisiones Técnicas

1. **GitHub Actions vs otros CI:**
   - Integración nativa con GitHub
   - Sin costo para repositorios públicos
   - Fácil configuración con acciones pre-hechas

2. **Pre-commit hooks:**
   - Validación local antes de push
   - Ahorra tiempo en CI
   - Feedback inmediato al desarrollador

3. **Codecov:**
   - Tracking de coverage a lo largo del tiempo
   - Visualización de coverage en PRs
   - Badges automáticos

---

## 📚 REFERENCIAS

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [golangci-lint](https://golangci-lint.run/)
- [pre-commit](https://pre-commit.com/)
- [Codecov](https://about.codecov.io/)

---

## ✅ CHECKLIST DE FINALIZACIÓN

- [ ] Workflows de GitHub Actions creados
- [ ] Backend CI ejecutándose correctamente
- [ ] Frontend CI ejecutándose correctamente
- [ ] Pre-commit hooks configurados
- [ ] Script de instalación de hooks disponible
- [ ] golangci-lint configurado
- [ ] ESLint + Prettier configurados
- [ ] Badges en README
- [ ] Documentación en CONTRIBUTING.md
- [ ] Tests pasando en CI

---

**Última actualización:** 2026-01-26
