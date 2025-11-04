.PHONY: help build test clean lint fmt ci-check release-snapshot install docker-build docker-run

# Variables
BINARY_NAME=git-cc
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILT_BY?=$(shell whoami)

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 go build \
		-ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -X main.builtBy=$(BUILT_BY)" \
		-o $(BINARY_NAME) .

build-all: ## Build all platform binaries using Goreleaser
	@echo "Building all platform binaries..."
	goreleaser build --snapshot --clean

test: ## Run all tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go test -v -race ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -race -tags=integration ./...

coverage: test ## Show test coverage
	@echo "Coverage report:"
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running golangci-lint..."
	golangci-lint run --timeout=10m

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

go-mod-tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	go mod tidy

ci-check: ## Run all CI checks locally
	@echo "Running CI checks..."
	@make fmt
	@make go-mod-tidy
	@make lint
	@make test
	@make build
	@echo "All CI checks passed!"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(BINARY_NAME) *.exe
	rm -rf dist/
	rm -f coverage.out coverage.html
	go clean -cache

install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -X main.builtBy=$(BUILT_BY)"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg DATE=$(DATE) \
		--build-arg BUILT_BY=$(BUILT_BY) \
		.

docker-run: docker-build ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -it $(BINARY_NAME):$(VERSION) --version

release-snapshot: ## Create a release snapshot using Goreleaser
	@echo "Creating release snapshot..."
	goreleaser build --snapshot --clean

release-check: ## Check Goreleaser configuration
	@echo "Checking Goreleaser configuration..."
	goreleaser check

release-test: ## Test release process locally
	@echo "Testing release process..."
	goreleaser release --snapshot --clean --skip-publish
	@echo "Release test complete!"

release-notes: ## Generate release notes preview
	@echo "Generating release notes..."
	goreleaser release --snapshot --clean --skip-announce --skip-publish --skip-validate

dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	@echo "Installing tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser/v2@latest
	@echo "Development environment setup complete!"

security-scan: ## Run security scan
	@echo "Running security scan..."
	@command -v gosec >/dev/null 2>&1 || { echo "gosec not installed. Run: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; exit 1; }
	gosec ./...

all: fmt lint test build ## Run all common development tasks