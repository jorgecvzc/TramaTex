.PHONY: help build run test clean deps fmt lint

# Variables
BINARY_NAME=tramatex.exe
CMD_PATH=./cmd/api
MAIN_PKG=$(CMD_PATH)

help:
	@echo "TramaTex Backend - Available Commands"
	@echo "===================================="
	@echo "make build       - Build the binary"
	@echo "make run         - Build and run the server"
	@echo "make test        - Run all tests"
	@echo "make test-unit   - Run unit tests only"
	@echo "make test-integration - Run integration tests only"
	@echo "make coverage    - Run tests with coverage report"
	@echo "make coverage-pricing - Run pricing tests with coverage report"
	@echo "make deps        - Download and tidy dependencies"
	@echo "make fmt         - Format code (gofmt)"
	@echo "make lint        - Run linter (if available)"
	@echo "make clean       - Remove binary and build artifacts"
	@echo "make dev         - Watch mode for development (requires air or similar)"

# Build the binary
build:
	@echo "🔨 Building TramaTex API..."
	@go build -o $(BINARY_NAME) $(MAIN_PKG)
	@echo "✅ Build complete: $(BINARY_NAME)"

# Build and run the server
run: build
	@echo "🚀 Starting TramaTex API..."
	@./$(BINARY_NAME)

# Run all tests
test:
	@echo "🧪 Running all tests..."
	@go test ./... -v

# Run unit tests only
test-unit:
	@echo "🧪 Running unit tests..."
	@go test ./internal/shared/tests/unit/... -v

# Run integration tests only
test-integration:
	@echo "🧪 Running integration tests..."
	@go test ./internal/shared/tests/integration/... -v

# Run tests with coverage
coverage:
	@echo "📊 Running tests with coverage..."
	@go test ./... -v -coverprofile=coverage-reports/coverage.out
	@go tool cover -html=coverage-reports/coverage.out -o coverage-reports/coverage.html
	@echo "✅ Coverage report: coverage-reports/coverage.html"

# Run pricing tests with coverage
coverage-pricing:
	@echo "📊 Running pricing tests with coverage..."
	@go test ./internal/pricing/... -v -coverprofile=coverage-reports/coverage-pricing.out
	@go tool cover -html=coverage-reports/coverage-pricing.out -o coverage-reports/coverage-pricing.html
	@echo "✅ Pricing coverage report: coverage-reports/coverage-pricing.html"

# Download and tidy dependencies
deps:
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies updated"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Run linter (requires golangci-lint to be installed)
lint:
	@echo "🔍 Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1)
	@golangci-lint run ./...
	@echo "✅ Lint complete"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html coverage-reports/*.out coverage-reports/*.html
	@go clean
	@echo "✅ Clean complete"

# Development mode (requires 'air' for hot reload)
dev:
	@echo "👀 Starting development mode (requires 'air')..."
	@which air > /dev/null || (echo "air not found. Install with: go install github.com/cosmtrek/air@latest"; exit 1)
	@air

# Run go mod tidy
tidy:
	@echo "📋 Running go mod tidy..."
	@go mod tidy
	@echo "✅ go.mod tidied"

# Build for production (stripped binary)
build-prod:
	@echo "🔨 Building production binary..."
	@go build -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_PKG)
	@echo "✅ Production build complete: $(BINARY_NAME)"

# Docker commands (if Dockerfile exists)
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t tramatex:latest .
	@echo "✅ Docker image built"

docker-run:
	@echo "🚀 Running Docker container..."
	@docker run -p 8080:8080 tramatex:latest
	@echo "✅ Container running"

# Default target
.DEFAULT_GOAL := help
