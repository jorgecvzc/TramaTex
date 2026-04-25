# CI/CD Workflow - TramaTex

## Overview

TramaTex uses **GitHub Actions** for continuous integration and continuous deployment. Our CI/CD pipeline ensures code quality, security, and reliability through automated testing, linting, and security audits.

---

## Workflows

### 1. Backend CI (`.github/workflows/backend.yml`)

**Triggers:**
- Push to `develop` and `master` branches
- Pull requests to `develop` and `master` branches
- Changes in `apps/tramatex-api/**` or the workflow file itself

**Concurrency Policy:**
We use a concurrency group based on the workflow name and the git reference (`github.ref`). In-progress runs for the same branch or PR are automatically cancelled to save resources and avoid duplicate notifications.

**Jobs:**

#### Test Job
- **Environment:** Ubuntu Latest with PostgreSQL 15
- **PostgreSQL Configuration:**
  - **User:** `tramatex`
  - **Password:** `tramatex`
  - **Database:** `tramatex_db`
  - **Auth Method:** `trust` (to prevent "role root does not exist" errors)
- **Steps:**
  1. Checkout code
  2. Set up Go 1.23 with caching
  3. Install dependencies
  4. Run tests with race detection
- **Database Connection Handling:** 
  The tests use `test_helpers.go` which prioritize environment variables (`DB_USER`, `DB_HOST`, etc.) over `.env` files, ensuring compatibility with the CI environment.

**Environment Variables in CI:**
```yaml
DB_HOST: localhost
DB_PORT: 5432
DB_USER: tramatex
DB_PASSWORD: tramatex
DB_NAME: tramatex_db
PGUSER: tramatex
PGPASSWORD: tramatex
JWT_SECRET: a-very-secret-key-for-testing-ci-12345
```

---

### 2. Deploy a Producción (`.github/workflows/deploy-production.yml`)

**Triggers:**
- Push a la rama `master` (automático tras cada PR mergeado)
- Manual: Via `workflow_dispatch` con opciones configurables

**Estrategia:**
Las imágenes se construyen en GitHub Actions (7 GB RAM) y se publican en GHCR. El servidor de producción (DigitalOcean Droplet) solo hace `pull` y `up`, sin builds en producción.

**Inputs del workflow_dispatch:**
| Input | Default | Descripción |
|---|---|---|
| `rebuild_images` | `true` | Reconstruye y publica imágenes antes del deploy |
| `fresh_db` | `false` | Destruye el volumen de PostgreSQL (¡DESTRUCTIVO!) |

**Jobs:**

#### build-api / build-frontend
- Login a GHCR con `GITHUB_TOKEN`
- `docker build` desde `Dockerfile` / `Dockerfile.frontend`
- `docker push` a `ghcr.io/<owner>/tramatex-api:latest` / `tramatex-frontend:latest`

#### deploy
- SSH al Droplet vía `appleboy/ssh-action`
- `git reset --hard origin/master` para sincronizar `docker-compose.remote.yml`
- Escribe `docker/.env` desde el secreto `ENV_PROD`
- Si `fresh_db=true`: `docker compose down --volumes` + elimina volumen (⚠️ pérdida total de datos)
- `docker compose pull` + `docker compose up -d --force-recreate`
- Health-check contra `http://localhost/api/health` (30 reintentos × 2 s)
- `docker image prune -f`

**Secrets requeridos:**
- `SSH_HOST`, `SSH_USER`, `SSH_PRIVATE_KEY` — acceso al Droplet
- `ENV_PROD` — contenido completo del `docker/.env` de producción

---

### 3. Demo Weekly Reset (`.github/workflows/demo-reset.yml`)

**Triggers:**
- Scheduled: Every Sunday at 03:00 AM UTC
- Manual: Via `workflow_dispatch`

**Purpose:**
Resets the production/staging demo environment to a clean state by wiping database volumes and re-running migrations and seed data.

**Key Steps:**
1. Stop the Docker stack
2. Wipe PostgreSQL volumes (`docker compose down -v`)
3. Restart the stack
4. Verify that the API becomes healthy and seed data (e.g., admin user) is present.

---

### 4. Frontend CI (`.github/workflows/frontend.yml`)

**Triggers:**
- Push to `develop` and `master` branches
- Pull requests to `develop` and `master` branches
- Changes in `apps/frontend/**` or the workflow file itself

**Jobs:**

#### Test & Build Job
- **Environment:** Ubuntu Latest with Node.js 20
- **Steps:**
  1. Checkout code
  2. Set up Node.js with npm caching
  3. Install dependencies (`npm ci`)
  4. Run unit tests
  5. Build production bundle
  6. Upload build artifacts (7-day retention)

---

## Best Practices

### For Developers

1. **Run Tests Locally Before Push:**
   ```bash
   # Backend
   cd apps/tramatex-api
   go test -v ./...
   
   # Frontend
   cd apps/frontend
   npm run test:unit
   ```

2. **Environment Priority:**
   The `test_helpers` in the backend follow this priority for configuration:
   1. Specific test variables (`TRAMATEX_TEST_DB_*`)
   2. Local `.env.local` / `.env.remote` files
   3. Standard environment variables (`DB_USER`, `DB_HOST`, etc.)
   4. Hardcoded defaults (`localhost`, `tramatex`)

---

## Troubleshooting

### "role root does not exist" (PostgreSQL)
This typically happens when the Go driver falls back to the system user (root) because the connection string is incomplete or env variables are ignored. 
**Solution:** Ensure `PGUSER` or `DB_USER` are correctly exported in the CI step and that `POSTGRES_HOST_AUTH_METHOD: trust` is used in the service container.

### Duplicate CI Runs
If you see two identical workflows running for a single commit, it's usually because both `push` and `pull_request` events are triggered for the same branch.
**Solution:** The `concurrency` block added to our workflows automatically manages this by cancelling the redundant run.

---

**Last Updated:** 2026-04-25  
**Related:** [Testing Guidelines](testing-guidelines.md), [Deployment Guide](deployment-guide.md)

