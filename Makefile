# ==============================================
# E-Commerce API Makefile
# ==============================================

# Variables
APP_NAME := ecommerce-api
MAIN_PATH := ./cmd/api
BINARY_PATH := ./bin/$(APP_NAME)
COVERAGE_FILE := coverage.out

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOVET := $(GOCMD) vet
GOMOD := $(GOCMD) mod
GOFMT := gofmt

# Build flags
LDFLAGS := -ldflags="-w -s"

.PHONY: all build run test test-coverage clean lint fmt vet seed docker-build docker-up docker-down help

# Default target
all: lint test build

## build: Build the application binary
build:
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p bin
	@$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "✅ Build complete: $(BINARY_PATH)"

## run: Run the application
run:
	@echo "🚀 Starting $(APP_NAME)..."
	@$(GOCMD) run $(MAIN_PATH)/main.go

## test: Run all tests
test:
	@echo "🧪 Running tests..."
	@$(GOTEST) -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@$(GOTEST) -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

## test-race: Run tests with race detector
test-race:
	@echo "🧪 Running tests with race detector..."
	@$(GOTEST) -race -v ./...

## seed: Seed the database with initial data
seed:
	@echo "🌱 Seeding database..."
	@$(GOCMD) run $(MAIN_PATH)/main.go -seed
	@echo "✅ Seed complete"

## clean: Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f $(COVERAGE_FILE) coverage.html
	@echo "✅ Clean complete"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Run: brew install golangci-lint"; \
	fi

## fmt: Format code
fmt:
	@echo "📝 Formatting code..."
	@$(GOFMT) -s -w .
	@echo "✅ Format complete"

## vet: Run go vet
vet:
	@echo "🔍 Running go vet..."
	@$(GOVET) ./...

## mod-tidy: Tidy go modules
mod-tidy:
	@echo "📦 Tidying modules..."
	@$(GOMOD) tidy
	@echo "✅ Modules tidied"

## mod-download: Download dependencies
mod-download:
	@echo "📥 Downloading dependencies..."
	@$(GOMOD) download
	@echo "✅ Dependencies downloaded"

## docker-build: Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(APP_NAME):latest .
	@echo "✅ Docker image built: $(APP_NAME):latest"

## docker-up: Start Docker Compose services
docker-up:
	@echo "🐳 Starting services..."
	@docker-compose up -d
	@echo "✅ Services started"

## docker-down: Stop Docker Compose services
docker-down:
	@echo "🐳 Stopping services..."
	@docker-compose down
	@echo "✅ Services stopped"

## docker-logs: View Docker logs
docker-logs:
	@docker-compose logs -f

## swagger: Generate Swagger documentation (requires swag)
swagger:
	@echo "📚 Generating Swagger docs..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g $(MAIN_PATH)/main.go -o ./docs; \
	else \
		echo "⚠️  swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

## help: Show this help message
help:
	@echo "E-Commerce API - Available Commands:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'