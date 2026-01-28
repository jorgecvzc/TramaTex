# TramaTex - Makefile with Dual Docker Support (Windows Local + Linux Remote)
# Supports: make docker-up ENV=local  OR  make docker-up ENV=remote
# Default behavior: Uses .env.remote (Linux pcele Server - PRIMARY)

.PHONY: help setup docker-build docker-up docker-down docker-logs docker-status \
        tramatex-api-build tramatex-api-run tramatex-api-test tramatex-api-test-unit tramatex-api-coverage \
        tramatex-api-lint tramatex-api-fmt tramatex-api-vet tramatex-api-deps db-migrate qa clean \
        install-tools env-init test-connectivity docker-clean

# ============================================================================
# GLOBAL CONFIGURATION
# ============================================================================

BINARY_NAME=tramatex
GO=go

# Detect environment (default: remote - PRIMARY)
ENV ?= remote
ifeq ($(ENV),remote)
    DOCKER_COMPOSE_FILE=docker-compose.remote.yml
    ENV_FILE=.env.remote
    SSH_HOST=pcele
    SSH_USER=ele
    DOCKER_HOST_FLAG=-H ssh://$(SSH_USER)@$(SSH_HOST)
else
    DOCKER_COMPOSE_FILE=docker-compose.local.yml
    ENV_FILE=.env.local
    DOCKER_HOST_FLAG=
endif

DOCKER_COMPOSE=docker-compose -f $(DOCKER_COMPOSE_FILE) --env-file $(ENV_FILE)

# ============================================================================
# HELP & INITIALIZATION
# ============================================================================

help: ## Show this help message
	@echo "🏗️  TramaTex - Dual Docker Support (Local + Remote)"
	@echo ""
	@echo "USAGE:"
	@echo "  make TARGET                     (uses REMOTE Linux pcele - DEFAULT)"
	@echo "  make TARGET ENV=local           (explicitly use Windows Docker Desktop)"
	@echo "  make TARGET ENV=remote          (explicitly use Linux pcele via SSH)"
	@echo ""
	@echo "⭐ DEFAULT: Linux pcele server (PRIMARY ENVIRONMENT)"
	@echo ""
	@echo "SETUP FIRST:"
	@echo "  make env-init                   (Initialize SSH/Docker connections)"
	@echo "  make test-connectivity          (Test both environments)"
	@echo ""
	@echo "Available Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "EXAMPLES:"
	@echo "  make docker-up                  # Start REMOTE stack (pcele) - DEFAULT"
	@echo "  make docker-up ENV=remote       # Start REMOTE stack (pcele) - Explicit"
	@echo "  make docker-up ENV=local        # Start LOCAL stack (Windows)"
	@echo "  make docker-logs                # View logs from remote (default)"
	@echo "  make docker-logs ENV=local      # View logs from local (Windows)"

env-init: ## Initialize environment files and test connectivity
	@echo "🔧 Initializing TramaTex environments..."
	@echo ""
	@echo "📍 Environment: LOCAL (Windows Docker Desktop)"
	@if [ -f .env.local ]; then \
		echo "✓ .env.local exists"; \
	else \
		echo "❌ .env.local missing"; \
	fi
	@echo ""
	@echo "📍 Environment: REMOTE (pcele Linux Server)"
	@if [ -f .env.remote ]; then \
		echo "✓ .env.remote exists"; \
	else \
		echo "❌ .env.remote missing"; \
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

docker-build: ## Build Docker image (REMOTE by default, use ENV=local for local)
	@echo "🐳 Building Docker image on $(ENV)..."
	@echo "Using: $(DOCKER_COMPOSE_FILE)"
	@$(DOCKER_COMPOSE) build

docker-up: ## Start Docker stack (REMOTE by default, use ENV=local for local)
	@echo "🚀 Starting Docker stack on $(ENV)..."
	@echo "Environment: $(ENV_FILE)"
	@echo "Compose file: $(DOCKER_COMPOSE_FILE)"
	@echo ""
	@$(DOCKER_COMPOSE) up -d
	@echo ""
	@echo "✓ Docker stack started!"
	@sleep 2
	@make docker-status ENV=$(ENV)

docker-down: ## Stop Docker stack (REMOTE by default, use ENV=local for local)
	@echo "⛔ Stopping Docker stack on $(ENV)..."
	@$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs (REMOTE by default, use ENV=local for local)
	@echo "📋 Docker logs ($(ENV)):"
	@$(DOCKER_COMPOSE) logs -f

docker-status: ## Show Docker containers status (REMOTE by default, use ENV=local for local)
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

frontend-build: ## Build frontend
	cd apps/frontend && npm run build

# ============================================================================
# DOCUMENTATION
# ============================================================================

docs-view: ## View documentation index
	@echo "Documentación TramaTex"
	@echo ""
	@echo "ADRs:"
	@ls -1 docs/2_architecture/adr/ | grep ADR
	@echo ""
	@echo "Módulos:"
	@ls -1 docs/3_modules/ | grep -v "^_"
	@echo ""
	@echo "Sprints:"
	@ls -1 docs/archive/sprints/ | grep -v "_TEMPLATE"

.DEFAULT_GOAL := help
