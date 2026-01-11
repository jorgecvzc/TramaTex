# TramaTex - Makefile Global & Backend

.PHONY: help setup docker-build docker-up docker-down backend-build backend-test backend-run backend-test-unit backend-coverage backend-lint backend-fmt backend-vet backend-deps backend-logs db-migrate-up db-migrate-down db-seed qa clean install-tools

# Variables
BINARY_NAME=tramatex
GO=go
DOCKER_COMPOSE=docker-compose

help: ## Show this help message
	@echo "🏗️  TramaTex - Complete Build & Development Commands"
	@echo ""
	@echo "Setup & Global:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Initial project setup
	@echo "Inicializando TramaTex..."
	@git init
	@echo "✓ Repositorio Git inicializado"
	@echo "✓ Estructura de carpetas creada"
	@echo "✓ Documentación base copiada"
	@echo ""
	@echo "Próximos pasos:"
	@echo "1. docker-compose build"
	@echo "2. docker-compose up"
	@echo "3. Navegar a http://localhost:5173"

# Backend Development
backend-build: ## Build the application
	@echo "🔨 Building $(BINARY_NAME)..."
	cd backend && $(GO) build -o ../bin/$(BINARY_NAME) .

backend-run: backend-build ## Build and run the application
	@echo "🚀 Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME)

backend-test: ## Run all backend tests
	@echo "🧪 Running tests..."
	cd backend && $(GO) test -v -race -coverprofile=coverage.out ./...
	cd backend && $(GO) tool cover -html=coverage.out -o coverage.html

backend-test-unit: ## Run unit tests only
	@echo "🧪 Running unit tests..."
	cd backend && $(GO) test -v -run '^Test[A-Z]' -short ./...

backend-coverage: ## Generate coverage report
	@echo "📊 Generating coverage report..."
	cd backend && $(GO) test -v -coverprofile=coverage.out ./...
	@echo "Coverage report: coverage.html"

backend-lint: ## Run linter
	@echo "✨ Running linter..."
	cd backend && golangci-lint run ./...

backend-fmt: ## Format backend code
	@echo "📝 Formatting code..."
	cd backend && $(GO) fmt ./...

backend-vet: ## Run go vet
	@echo "🔍 Running go vet..."
	cd backend && $(GO) vet ./...

backend-deps: ## Download Go dependencies
	@echo "📦 Downloading dependencies..."
	cd backend && $(GO) mod download
	cd backend && $(GO) mod tidy

backend-logs: ## View backend logs
	@echo "📋 Backend logs:"
	$(DOCKER_COMPOSE) logs -f api

# Database
db-migrate-up: ## Run migrations up
	@echo "🔄 Running database migrations..."
	cd backend && $(GO) run ./cmd/migrate/main.go up

db-migrate-down: ## Run migrations down
	@echo "🔄 Rolling back database migrations..."
	cd backend && $(GO) run ./cmd/migrate/main.go down

db-seed: ## Seed database with test data
	@echo "🌱 Seeding database..."
	cd backend && $(GO) run ./cmd/seed/main.go

# Docker
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	$(DOCKER_COMPOSE) build

docker-up: ## Start Docker containers
	@echo "🚀 Starting containers..."
	$(DOCKER_COMPOSE) up -d
	@echo "✅ Containers started"
	@echo "📱 API: http://localhost:8080"
	@echo "🗄️  Database: postgres://tramatex_user:tramatex_pass@localhost:5432/tramatex_db"

docker-down: ## Stop Docker containers
	@echo "🛑 Stopping containers..."
	$(DOCKER_COMPOSE) down

docker-ps: ## List running containers
	@echo "🐳 Running containers:"
	$(DOCKER_COMPOSE) ps

# Quality & Testing
qa: backend-lint backend-vet backend-fmt backend-test ## Run all quality checks

# Cleanup
clean: ## Clean build artifacts and temp files
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -rf coverage.*
	cd backend && $(GO) clean
	$(DOCKER_COMPOSE) down -v

# Install development tools
install-tools: ## Install development tools
	@echo "📦 Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install gotest.tools/gotestsum@latest

# Frontend
frontend-dev: ## Frontend development server
	cd frontend && npm run dev

frontend-build: ## Build frontend
	cd frontend && npm run build

# Documentation
docs-view: ## View documentation index
	@echo "Documentación TramaTex"
	@echo ""
	@echo "ADRs:"
	@ls -1 docs/adr/ | grep ADR
	@echo ""
	@echo "Módulos:"
	@ls -1 docs/modules/ | grep -v "^_"
	@echo ""
	@echo "Sesiones:"
	@ls -1 docs/sessions/ | grep "^2026\|^2027\|^2028"
	@echo ""
	@echo "Guías:"
	@ls -1 docs/guides/
