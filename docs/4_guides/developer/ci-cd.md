# CI/CD Workflow - TramaTex

## Overview

TramaTex uses **GitHub Actions** for continuous integration and continuous deployment. Our CI/CD pipeline ensures code quality, security, and reliability through automated testing, linting, and security audits.

---

## Workflows

### 1. Backend CI (`.github/workflows/backend.yml`)

**Triggers:**
- Push to `master`/`main` branches
- Pull requests to `master`/`main` branches
- Changes in `apps/tramatex-api/**`, `go.work`, or the workflow file itself

**Jobs:**

#### Test Job
- **Environment:** Ubuntu Latest with PostgreSQL 15
- **Steps:**
  1. Checkout code
  2. Set up Go 1.23 with caching
  3. Install dependencies
  4. Run tests with race detection and coverage
  5. Upload coverage to Codecov

**Environment Variables:**
```bash
DATABASE_URL=postgres://tramatex:tramatex123@localhost:5432/tramatex_test?sslmode=disable
JWT_SECRET=test-secret-key-for-ci
```

#### Lint Job
- **Linter:** golangci-lint (latest version)
- **Timeout:** 5 minutes
- **Checks:** Code style, potential bugs, complexity, security issues

#### Security Job
- **Tools:**
  - **nancy**: Dependency vulnerability scanning
  - **govulncheck**: Go vulnerability database check
- **Purpose:** Identify known vulnerabilities in dependencies

---

### 2. Frontend CI (`.github/workflows/frontend.yml`)

**Triggers:**
- Push to `master`/`main` branches
- Pull requests to `master`/`main` branches
- Changes in `apps/frontend/**` or the workflow file itself

**Jobs:**

#### Test & Build Job
- **Environment:** Ubuntu Latest with Node.js 20
- **Steps:**
  1. Checkout code
  2. Set up Node.js with npm caching
  3. Install dependencies (`npm ci`)
  4. Run unit tests
  5. Generate coverage report
  6. Upload coverage to Codecov
  7. Build production bundle
  8. Upload build artifacts (7-day retention)

#### Lint Job
- **Tools:**
  - **Prettier**: Code formatting check
  - **ESLint**: JavaScript/Vue linting
- **Note:** Gracefully skips if tools are not configured

#### Security Job
- **Tool:** npm audit
- **Level:** Moderate and above
- **Output:** JSON vulnerability report
- **Behavior:** Non-blocking (continues on error)

---

## Status Badges

The project README displays real-time CI status:

```markdown
[![Backend CI](https://github.com/joran-cortez/tramatex/actions/workflows/backend.yml/badge.svg)](https://github.com/joran-cortez/tramatex/actions/workflows/backend.yml)
[![Frontend CI](https://github.com/joran-cortez/tramatex/actions/workflows/frontend.yml/badge.svg)](https://github.com/joran-cortez/tramatex/actions/workflows/frontend.yml)
[![codecov](https://codecov.io/gh/joran-cortez/tramatex/branch/master/graph/badge.svg)](https://codecov.io/gh/joran-cortez/tramatex)
```

---

## Code Coverage

We use **Codecov** to track and visualize test coverage:

- **Backend Coverage:** Generated via `go test -coverprofile`
- **Frontend Coverage:** Generated via Vitest coverage reporter
- **Target:** Maintain > 80% coverage for business logic
- **Reporting:** Automatic upload on every CI run

---

## Best Practices

### For Developers

1. **Run Tests Locally Before Push:**
   ```bash
   # Backend
   cd apps/tramatex-api
   make test
   
   # Frontend
   cd apps/frontend
   npm run test:unit
   ```

2. **Check Linting:**
   ```bash
   # Backend
   cd apps/tramatex-api
   golangci-lint run
   
   # Frontend
   cd apps/frontend
   npx eslint src/
   ```

3. **Security Audit:**
   ```bash
   # Backend
   cd apps/tramatex-api
   go list -json -deps ./... | nancy sleuth
   
   # Frontend
   cd apps/frontend
   npm audit
   ```

### For Pull Requests

- ✅ All CI checks must pass before merge
- ✅ Coverage should not decrease
- ✅ No high/critical security vulnerabilities
- ✅ Code must pass linting rules

---

## Troubleshooting

### Backend CI Failures

**Test Failures:**
- Check database connection (PostgreSQL service health)
- Verify environment variables are set correctly
- Ensure migrations are up-to-date

**Lint Failures:**
- Run `golangci-lint run` locally
- Fix reported issues following Go best practices
- Update `.golangci.yml` if needed

**Security Failures:**
- Update vulnerable dependencies: `go get -u`
- Check nancy/govulncheck output for specific CVEs
- Create security issues for tracking

### Frontend CI Failures

**Test Failures:**
- Run `npm run test:unit` locally
- Check for missing mocks or test data
- Verify component dependencies

**Build Failures:**
- Run `npm run build` locally
- Check for TypeScript errors
- Verify all imports resolve correctly

**Security Failures:**
- Run `npm audit fix` for automatic fixes
- Review audit report for manual updates
- Document accepted risks if needed

---

## Future Enhancements

- [ ] Add E2E tests with Playwright in CI
- [ ] Deploy preview environments for PRs
- [ ] Add performance benchmarking
- [ ] Implement semantic release automation
- [ ] Add SAST/DAST security scanning

---

**Last Updated:** 2026-01-27  
**Related:** [ADR-010 Security Architecture](../../../2_architecture/adr/ADR-010-estrategia-seguridad-defensa-profundidad.md), [Backend Setup Guide](backend/README.md)
