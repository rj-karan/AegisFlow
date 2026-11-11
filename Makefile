SHELL := /bin/bash
.DEFAULT_GOAL := help

GOBIN ?= $(CURDIR)/bin
GOLANGCI_LINT_VERSION ?= v2.15.0
GOLANGCI_LINT := $(GOBIN)/golangci-lint

export GOBIN

.PHONY: help build run test lint fmt vet tidy clean dev dev-down dev-logs image

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Build all binaries into ./bin
	@mkdir -p bin
	go build -trimpath -ldflags="-s -w" -o bin/ ./cmd/...

run: ## Run the API locally
	go run ./cmd/api

test: ## Run tests with the race detector
	go test -race -count=1 -cover ./...

vet: ## go vet
	go vet ./...

fmt: ## gofmt everything
	gofmt -s -w .

tidy: ## go mod tidy
	go mod tidy

$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

lint: $(GOLANGCI_LINT) vet ## golangci-lint + vet
	$(GOLANGCI_LINT) run ./...

clean: ## Remove build output
	rm -rf bin

dev: ## Start postgres, redis and the API with Docker Compose
	docker compose up --build

dev-down: ## Stop the Compose stack and remove volumes
	docker compose down -v

dev-logs: ## Tail Compose logs
	docker compose logs -f

image: ## Build the API container image locally
	docker build --build-arg CMD=api -t aegisflow-api:dev .
