# TramaTex - Makefile with Dual Docker Support (Windows Local + Linux Remote)
# Supports: make docker-up ENV=local  OR  make docker-up ENV=remote
# Default behavior: Uses .env.local (Docker Desktop - PRIMARY)

.PHONY: help setup docker-build docker-up docker-down docker-logs docker-status \
        tramatex-api-build tramatex-api-run tramatex-api-test tramatex-api-test-unit tramatex-api-coverage \
        tramatex-api-lint tramatex-api-fmt tramatex-api-vet tramatex-api-deps db-migrate qa clean \
        install-tools env-init test-connectivity docker-clean deploy

# ============================================================================
# GLOBAL CONFIGURATION
# ============================================================================

BINARY_NAME=tramatex
GO=go

# Detect environment (default: local - PRIMARY)
ENV ?= local
ifeq ($(ENV),remote)
    DOCKER_COMPOSE_FILE=docker/docker-compose.remote.yml
    ENV_FILE=docker/.env
    SSH_HOST=pcele
    SSH_USER=ele
    DOCKER_HOST_FLAG=-H ssh://$(SSH_USER)@$(SSH_HOST)
else
    DOCKER_COMPOSE_FILE=docker/docker-compose.local.yml
    ENV_FILE=docker/.env
    DOCKER_HOST_FLAG=
endif

DOCKER_COMPOSE=docker compose -f $(DOCKER_COMPOSE_FILE) --env-file $(ENV_FILE)

# ============================================================================
# HELP & INITIALIZATION
# ============================================================================

help: ## Show this help message
	@echo "🏗️  TramaTex - Dual Docker Support (Local + Remote)"
	@echo ""
	@echo "USAGE:"
	@echo "  make TARGET                     (uses LOCAL Docker Desktop - DEFAULT)"
	@echo "  make TARGET ENV=local           (explicitly use Windows Docker Desktop)"
	@echo "  make TARGET ENV=remote          (explicitly use Linux pcele via SSH)"
	@echo ""
	@echo "⭐ DEFAULT: Local Docker Desktop (PRIMARY ENVIRONMENT)"
	@echo ""
	@echo "SETUP FIRST:"
	@echo "  make env-init                   (Initialize SSH/Docker connections)"
	@echo "  make test-connectivity          (Test both environments)"
	@echo ""
	@echo "Available Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "EXAMPLES:"
	@echo "  make docker-up                  # Start LOCAL stack (Windows) - DEFAULT"
	@echo "  make docker-up ENV=local        # Start LOCAL stack (Windows)"
	@echo "  make docker-up ENV=remote       # Start REMOTE stack (pcele) - Explicit"
	@echo "  make docker-logs                # View logs from local (default)"
	@echo "  make docker-logs ENV=remote     # View logs from remote (pcele)"
	@echo "  make deploy ENV=staging         # Deploy to pcele (staging LAN)"
	@echo "  make deploy ENV=prod            # Push to production (DigitalOcean)"

env-init: ## Initialize environment files and test connectivity
	@echo "🔧 Initializing TramaTex environments..."
	@echo ""
	@echo "📍 Environment: docker/.env"
	@if [ -f docker/.env ]; then \
		echo "✓ docker/.env exists"; \
	else \
		echo "❌ docker/.env missing — copy from docker/.env.example"; \
	fi
	@echo ""
	@echo "✓ Configuration initialized"
	@echo "➡️  Next: make test-connectivity"

test-connectivity: ## Test connectivity to both Docker environments
	@echo "🧪 Testing Docker connectivity..."
	@echo ""
	@echo "📍 LOCAL (Docker Desktop - Windows):"
	@if docker --version > /dev/null 2>&1; then \
		echo "✓ Docker Desktop: $$(docker --version)"; \
		echo "✓ Docker Compose: $$(docker-compose --version)"; \
	else \
		echo "❌ Docker not installed"; \
	fi
	@echo ""
	@echo "📍 REMOTE (pcele Linux - SSH):"
	@if ssh -o ConnectTimeout=3 -o BatchMode=no ele@pcele "docker --version" > /dev/null 2>&1; then \
		echo "✓ SSH connection successful"; \
	else \
		echo "⚠️  Cannot connect via SSH (may need credentials)"; \
	fi

# ============================================================================
# DOCKER COMMANDS - Dual Environment Support
# ============================================================================

docker-build: ## Build Docker image (LOCAL by default, use ENV=remote for remote)
	@echo "🐳 Building Docker image on $(ENV)..."
	@echo "Using: $(DOCKER_COMPOSE_FILE)"
	@$(DOCKER_COMPOSE) build

docker-up: ## Start Docker stack (LOCAL by default, use ENV=remote for remote)
	@echo "🚀 Starting Docker stack on $(ENV)..."
	@echo "Environment: $(ENV_FILE)"
	@echo "Compose file: $(DOCKER_COMPOSE_FILE)"
	@echo ""
	@$(DOCKER_COMPOSE) up -d
	@echo ""
	@echo "✓ Docker stack started!"
	@sleep 2
	@make docker-status ENV=$(ENV)

docker-down: ## Stop Docker stack (LOCAL by default, use ENV=remote for remote)
	@echo "⛔ Stopping Docker stack on $(ENV)..."
	@$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs (LOCAL by default, use ENV=remote for remote)
	@echo "📋 Docker logs ($(ENV)):"
	@$(DOCKER_COMPOSE) logs -f

docker-status: ## Show Docker containers status (LOCAL by default, use ENV=remote for remote)
	@echo "📊 Docker status ($(ENV)):"
	@echo ""
	@echo "Environment: $(ENV)"
	@echo "Compose file: $(DOCKER_COMPOSE_FILE)"
	@echo ""
	@$(DOCKER_COMPOSE) ps
	@echo ""
	@echo "Connection endpoints:"
	@echo "  API: http://localhost:8080/api/health"
	@echo "  PostgreSQL: postgresql://tramatex:tramatex@localhost:5432/tramatex"

docker-clean: ## Remove stopped containers and dangling images
	@echo "🧹 Cleaning Docker resources ($(ENV))..."
	@$(DOCKER_COMPOSE) down -v
	@echo "✓ Cleanup complete"

# ============================================================================
# SETUP & INITIALIZATION
# ============================================================================

setup: env-init test-connectivity docker-build docker-up ## Complete initial setup
	@echo ""
	@echo "✅ TramaTex setup complete!"
	@echo ""
	@echo "Access your application:"
	@echo "  API: http://localhost:8080/api/health"
	@echo "  Frontend: http://localhost:5173"
	@echo ""
	@echo "Check container status with:"
	@echo "  make docker-status"

# ============================================================================
# BACKEND COMMANDS
# ============================================================================

tramatex-api-build: ## Build Go binary
	@echo "🔨 Building $(BINARY_NAME)..."
	cd apps/tramatex-api && $(GO) build -o ../../bin/$(BINARY_NAME) .

tramatex-api-run: tramatex-api-build ## Build and run locally
	@echo "🚀 Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME)

tramatex-api-test: ## Run all tramatex-api tests
	@echo "🧪 Running all tramatex-api tests..."
	cd apps/tramatex-api && $(GO) test -v -race -coverprofile=coverage.out ./...
	cd apps/tramatex-api && $(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: tramatex-api/coverage.html"

tramatex-api-test-unit: ## Run unit tests only
	@echo "🧪 Running unit tests..."
	cd apps/tramatex-api && $(GO) test -v -run '^Test[A-Z]' -short ./...

tramatex-api-coverage: ## Generate test coverage report
	@echo "📊 Generating coverage report..."
	cd apps/tramatex-api && $(GO) test -coverprofile=coverage.out ./...
	cd apps/tramatex-api && $(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Report: tramatex-api/coverage.html"

tramatex-api-lint: ## Run Go linter
	@echo "🔍 Running golangci-lint..."
	cd apps/tramatex-api && golangci-lint run ./...

tramatex-api-fmt: ## Format Go code
	@echo "🎨 Formatting Go code..."
	cd apps/tramatex-api && $(GO) fmt ./...

tramatex-api-vet: ## Run Go vet
	@echo "🔬 Running go vet..."
	cd apps/tramatex-api && $(GO) vet ./...

tramatex-api-deps: ## Download and verify dependencies
	@echo "📦 Managing dependencies..."
	cd apps/tramatex-api && $(GO) mod tidy
	cd apps/tramatex-api && $(GO) mod download
	cd apps/tramatex-api && $(GO) mod verify

# ============================================================================
# DATABASE COMMANDS
# ============================================================================

db-migrate: ## Run database migrations
	@echo "🔄 Running migrations (via API start)..."
	@echo "Migrations run automatically on API startup"

# ============================================================================
# QUALITY ASSURANCE
# ============================================================================

qa: tramatex-api-fmt tramatex-api-vet tramatex-api-lint tramatex-api-test ## Run all QA checks
	@echo "✅ Quality assurance complete!"

# ============================================================================
# CLEANUP
# ============================================================================

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf bin/
	cd apps/tramatex-api && $(GO) clean
	cd apps/tramatex-api && rm -f coverage.out coverage.html
	@echo "✓ Cleanup complete"

# ============================================================================
# DEVELOPMENT TOOLS
# ============================================================================

install-tools: ## Install required development tools
	@echo "📦 Installing development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Tools installed"

# ============================================================================
# FRONTEND (Optional)
# ============================================================================

frontend-dev: ## Frontend development server
	cd apps/frontend && npm run dev

frontend-build: ## Build frontend for production
	cd apps/frontend && npm run build

frontend-test: ## Run frontend tests
	cd apps/frontend && npx vitest run --reporter=verbose

# ============================================================================
# DEPLOYMENT
# ============================================================================

deploy: ## Deploy to remote server (ENV=staging|prod BRANCH=branch_name)
ifeq ($(ENV),staging)
	$(eval DEPLOY_BRANCH ?= $(if $(BRANCH),$(BRANCH),staging))
	@echo "🚀 Deploying to staging (pcele LAN — 192.168.0.20) from branch $(DEPLOY_BRANCH)..."
	ssh ele@pcele "cd /opt/tramatex && git fetch origin $(DEPLOY_BRANCH) && \
		git checkout -B $(DEPLOY_BRANCH) origin/$(DEPLOY_BRANCH) && \
		docker compose -f docker/docker-compose.remote.yml --env-file docker/.env build && \
		docker compose -f docker/docker-compose.remote.yml --env-file docker/.env up -d && \
		docker image prune -f"
	@echo "✓ Deployed to pcele (staging) from branch $(DEPLOY_BRANCH)"
else ifeq ($(ENV),prod)
	@echo "🚀 Deploying to production (push to master triggers GitHub Actions)..."
	@echo "⚠️  Are you sure? This deploys to PRODUCTION (DigitalOcean)."
	@read -p "Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		git push origin staging:master; \
		echo "✓ Push to master complete — GitHub Actions will deploy automatically"; \
	else \
		echo "❌ Deployment cancelled"; \
	fi
else
	@echo "❌ Usage: make deploy ENV=staging|prod"
	@echo ""
	@echo "  staging - Deploy to pcele LAN server (SSH directo)"
	@echo "  prod    - Push staging → master (triggers GitHub Actions on DigitalOcean)"
endif

# ============================================================================
# DOCUMENTATION
# ============================================================================

docs-view: ## View documentation index
	@echo "Documentación TramaTex"
	@echo ""
	@echo "ADRs:"
	@ls -1 docs/architecture/adrs/ | grep ADR
	@echo ""
	@echo "Módulos:"
	@ls -1 docs/modules/ | grep -v "^_"
	@echo ""
	@echo "Sprints:"
	@ls -1 docs/log/sprints/ | grep -v "_TEMPLATE"

.DEFAULT_GOAL := help
