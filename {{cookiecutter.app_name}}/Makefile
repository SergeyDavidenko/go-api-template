# Makefile for {{cookiecutter.app_name}}

# Variables
APP_NAME={{cookiecutter.app_name}}
BUILD_DIR=build
DOCKER_IMAGE={{cookiecutter.docker_hub_username}}/{{cookiecutter.app_name}}
VERSION=$(shell git describe --tags --always --dirty)
COMMIT_HASH=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.CommitHash=$(COMMIT_HASH) -X main.BuildTime=$(BUILD_TIME)"

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go

.PHONY: build-linux
build-linux: ## Build the application for Linux
	@echo "Building $(APP_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux ./cmd/main.go

.PHONY: run
run: ## Run the application
	@echo "Running $(APP_NAME)..."
	$(GOCMD) run ./cmd/main.go

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 -p 1499:1499 $(DOCKER_IMAGE):latest

.PHONY: docker-push
docker-push: ## Push Docker image
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest


.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: migration-create
migration-create: ## Create a new migration file (usage: make migration-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: migration name is required. Usage: make migration-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)..."
	@timestamp=$$(date +%s); \
	mkdir -p migrations; \
	echo "-- Migration: $(name)" > migrations/$${timestamp}_$(name).up.sql; \
	echo "-- Rollback: $(name)" > migrations/$${timestamp}_$(name).down.sql; \
	echo "Created migration files:"; \
	echo "  migrations/$${timestamp}_$(name).up.sql"; \
	echo "  migrations/$${timestamp}_$(name).down.sql"

.PHONY: check
check: fmt lint test ## Run all checks (format, lint, test)

.PHONY: all
all: clean deps check build ## Run all steps (clean, deps, check, build)
