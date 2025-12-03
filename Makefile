BINARY_NAME := wms-server
GO_CMD := go
MAIN_PATH := cmd/server/main.go
BUILD_DIR := bin
BIN_PATH := $(BUILD_DIR)/$(BINARY_NAME)
ENV_FILE := .env
ENV_EXAMPLE := .env.example
GOPATH_BIN := $(shell $(GO_CMD) env GOPATH)/bin

.DEFAULT_GOAL := all

.PHONY: all build test run dev clean docker-up docker-down docker-logs docker-restart docker-build docker-deploy deps tidy fmt vet setup help

all: build test ## Build the binary and run the full test suite

build: ## Compile the Go application into ./bin/wms-server
	@mkdir -p $(BUILD_DIR)
	$(GO_CMD) build -o $(BIN_PATH) $(MAIN_PATH)

test: ## Execute Go tests with verbose output
	$(GO_CMD) test ./... -v

run: build ## Start local server after ensuring docker services and env config
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Missing $(ENV_FILE). Run 'make setup' to create it."; \
		exit 1; \
	fi
	docker-compose up -d
	./$(BIN_PATH)

dev: setup ## Hot-reload server using Air
	@if ! command -v air >/dev/null 2>&1; then \
		echo "Installing Air for hot-reload..."; \
		$(GO_CMD) install github.com/air-verse/air@latest; \
	fi
	@if [ ! -f .air.toml ]; then \
		echo ".air.toml configuration is missing."; \
		exit 1; \
	fi
	PATH="$(PATH):$(GOPATH_BIN)" air -c .air.toml

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

deps: ## Download Go module dependencies
	$(GO_CMD) mod download

tidy: ## Ensure go.mod and go.sum are in sync
	$(GO_CMD) mod tidy

fmt: ## Format Go source files
	$(GO_CMD) fmt ./...

vet: ## Run go vet for static analysis
	$(GO_CMD) vet ./...

setup: ## Copy .env.example to .env when missing
	@if [ ! -f $(ENV_FILE) ]; then \
		cp $(ENV_EXAMPLE) $(ENV_FILE); \
		echo "Created $(ENV_FILE) from $(ENV_EXAMPLE)."; \
	else \
		echo "$(ENV_FILE) already exists."; \
	fi

docker-up: ## Start docker-compose services
	docker-compose up -d

docker-down: ## Stop docker-compose services
	docker-compose down

docker-logs: ## Tail docker-compose logs
	docker-compose logs -f

docker-restart: ## Restart docker-compose services
	docker-compose restart

docker-build: ## Build Docker image
	docker build -t wms-server:latest .

docker-deploy: ## Deploy with docker-compose (build and start all services)
	docker-compose up -d --build

docker-clean: ## Stop and remove all containers, networks and volumes
	docker-compose down -v

help: ## Display available make targets
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
