# Makefile for managing development and production environments
define IRIS_ASCII
$(BLUE)
╔══╗╔═══╗╔══╗╔═══╗ ║║ ╔═══╗╔══╗╔═══╗╔══╗
║  ║║   ║║  ║║   ║ ║║ ║   ║║  ║║   ║║  ║
║  ║║   ║║  ║╚══╗║ ║║ ║╔══╝║  ║║   ║║  ║
║  ║║   ║║  ║   ║║ ║║ ║║   ║  ║║   ║║  ║
╚══╝╚═══╝╚══╝═══╝║ ║║ ║╚═══╝╚══╝╚═══╝╚══╝
                  ║ ║║ ║
╔══╗╔═══╗╔══╗╔═══╝ ╚╝ ╚═══╗╔══╗╔═══╗╔══╗
║  ║║   ║║  ║║         ║   ║║  ║║   ║║  ║
║  ║║   ║║  ║╚══╗     ╔╝   ║║  ║║   ║║  ║
║  ║║   ║║  ║   ║     ║    ║║  ║║   ║║  ║
╚══╝╚═══╝╚══╝═══╝     ╚════╝╚══╝╚═══╝╚══╝
$(RESET)
endef
export IRIS_ASCII

# Colors for pretty output
BLUE := $(shell printf "\033[36m")
GREEN := $(shell printf "\033[32m")
YELLOW := $(shell printf "\033[33m")
RED := $(shell printf "\033[31m")
RESET := $(shell printf "\033[0m")

# Default shell
SHELL := /bin/bash

DEV_ENV_FILE := .env.development
PROD_ENV_FILE := .env.production

# Docker compose commands for dev and prod
DOCKER_COMPOSE_DEV := docker-compose --profile development --env-file $(DEV_ENV_FILE)
DOCKER_COMPOSE_PROD := docker-compose --profile production --env-file $(PROD_ENV_FILE)

# Include .env file by default
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: help
help:
	@echo "$$IRIS_ASCII"
	@echo "$(BLUE)Available commands:$(RESET)"
	@echo "$(GREEN)Setup Commands:$(RESET)"
	@echo "  make init            - Initialize development environment"
	@echo "  make deps            - Update dependencies"
	@echo "$(GREEN)Development Commands:$(RESET)"
	@echo "  make start-dev         - Start development environment"
	@echo "  make stop-dev          - Stop development environment"
	@echo "  make shell-app-dev     - Shell into development app container"
	@echo "  make shell-db-dev      - Shell into development database container"
	@echo "  make restart-dev       - Restart development environment"
	@echo "  make build-dev         - Build development Docker images"
	@echo "  make rebuild-dev       - Force rebuild development Docker images"
	@echo "  make build-dev-nocache - Build development Docker images without cache"
	@echo "$(GREEN)Production Commands:$(RESET)"
	@echo "  make start-prod         - Start production environment locally"
	@echo "  make stop-prod          - Stop production environment locally"
	@echo "  make shell-app-prod     - Shell into production app container"
	@echo "  make shell-db-prod      - Shell into production database container"
	@echo "  make restart-prod       - Restart production environment locally"
	@echo "  make build-prod         - Build production Docker images locally"
	@echo "  make rebuild-prod       - Force rebuild production Docker images locally"
	@echo "  make build-prod-nocache - Build production Docker images without cache locally"
	@echo "$(GREEN)Redis Commands:$(RESET)"
	@echo "  make shell-redis-dev    - Open Redis CLI in development container"
	@echo "  make shell-redis-prod   - Open Redis CLI in production container"
	@echo "  make logs-redis-dev     - View development Redis logs"
	@echo "  make logs-redis-prod    - View production Redis logs"
	@echo "  make redis-monitor-dev  - Monitor Redis commands in development"
	@echo "  make redis-monitor-prod - Monitor Redis commands in production"
	@echo "  make redis-info-dev     - Show Redis server info in development"
	@echo "  make redis-info-prod    - Show Redis server info in production"
	@echo "  make redis-flush-dev    - Flush development Redis database (with confirmation)"
	@echo "$(GREEN)Testing & Linting:$(RESET)"
	@echo "  make test          - Run tests"
	@echo "  make lint          - Run linter"
	@echo "  make format        - Format code"
	@echo "  make lint-nil      - Run nilaway"
	@echo "$(GREEN)Database Commands:$(RESET)"
	@echo "  make db-status-dev  - Check development database connection"
	@echo "  make db-status-prod - Check production database connection"
	@echo "  make db-console-dev - Enter development PostgreSQL console"
	@echo "  make db-console-prod - Enter production PostgreSQL console"
	@echo "$(GREEN)Logs & Monitoring:$(RESET)"
	@echo "  make logs-dev      - View all development logs"
	@echo "  make logs-prod     - View all production logs"
	@echo "  make logs-app-dev  - View development application logs"
	@echo "  make logs-app-prod - View production application logs"
	@echo "  make logs-db-dev   - View development database logs"
	@echo "  make logs-db-prod  - View production database logs"
	@echo "$(GREEN)Container Commands:$(RESET)"
	@echo "  make shell-app-dev  - Shell into development app container"
	@echo "  make shell-app-prod - Shell into production app container"
	@echo "  make shell-db-dev   - Shell into development database container"
	@echo "  make shell-db-prod  - Shell into production database container"
	@echo "$(GREEN)Cleanup Commands:$(RESET)"
	@echo "  make clean-dev      - Remove development containers"
	@echo "  make clean-prod     - Remove production containers"
	@echo "  make clean-all      - Remove all Docker artifacts"
	@echo "$(GREEN)Service Control:$(RESET)"
	@echo "  make pause-dev      - Pause development services"
	@echo "  make pause-prod     - Pause production services"
	@echo "  make unpause-dev    - Unpause development services"
	@echo "  make unpause-prod   - Unpause production services"
	@echo "$(GREEN)Mock Server Commands:$(RESET)"
	@echo "  make mock-private       - Start private endpoints mock server"
	@echo "  make mock-public        - Start public endpoints mock server"
	@echo "  make stop-mock-servers  - Stop all mock servers"
	@echo "$(GREEN)Environment:$(RESET)"
	@echo "  make pre-prod           - Start pre-production environment"

.PHONY: init
init:
	@echo "$(BLUE)== Initializing Development Environment ==$(RESET)"
	brew install go
	brew install node
	brew install pre-commit
	brew install golangci-lint
	brew upgrade golangci-lint

	@echo "$(BLUE)== Installing Pre-Commit Hooks ==$(RESET)"
	pre-commit install
	pre-commit autoupdate
	pre-commit install --install-hooks
	pre-commit install --hook-type commit-msg

.PHONY: deps
deps:
	@echo "$(BLUE)== Updating Dependencies ==$(RESET)"
	go mod tidy
	go mod verify

# Development Environment Commands
.PHONY: start-dev
start-dev:
	@echo "$(BLUE)== Starting Development Environment ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) up -d --build

.PHONY: stop-dev
stop-dev:
	@echo "$(BLUE)== Stopping Development Environment ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) stop

.PHONY: restart-dev
restart-dev:
	@echo "$(BLUE)== Restarting Development Environment ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) restart

.PHONY: build-dev
build-dev:
	@echo "$(BLUE)== Building Development Docker Images ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) build

.PHONY: rebuild-dev
rebuild-dev:
	@echo "$(BLUE)== Rebuilding Development Docker Images ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) build --force-rm

.PHONY: build-dev-nocache
build-dev-nocache:
	@echo "$(BLUE)== Building Development Docker Images Without Cache ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) build --no-cache

# Production Environment Commands
.PHONY: start-prod
start-prod:
	@echo "$(BLUE)== Starting Production Environment ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) up -d --build

.PHONY: stop-prod
stop-prod:
	@echo "$(BLUE)== Stopping Production Environment ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) stop

.PHONY: restart-prod
restart-prod:
	@echo "$(BLUE)== Restarting Production Environment ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) restart

.PHONY: build-prod
build-prod:
	@echo "$(BLUE)== Building Production Docker Images ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) build

.PHONY: rebuild-prod
rebuild-prod:
	@echo "$(BLUE)== Rebuilding Production Docker Images ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) build --force-rm

.PHONY: build-prod-nocache
build-prod-nocache:
	@echo "$(BLUE)== Building Production Docker Images Without Cache ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) build --no-cache

# Database Commands
.PHONY: db-status-dev
db-status-dev:
	@echo "$(BLUE)== Checking Development Database Connection ==$(RESET)"
	@$(DOCKER_COMPOSE_DEV) exec db pg_isready -U "$(DB_USER)" -d "$(DB_NAME)" || echo "$(RED)Database is not ready$(RESET)"

.PHONY: db-status-prod
db-status-prod:
	@echo "$(BLUE)== Checking Production Database Connection ==$(RESET)"
	@$(DOCKER_COMPOSE_PROD) exec db pg_isready -U "$(DB_USER)" -d "$(DB_NAME)" || echo "$(RED)Database is not ready$(RESET)"

.PHONY: db-console-dev
db-console-dev:
	@echo "$(BLUE)== Entering Development PostgreSQL Console ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) exec db psql -U "$(DB_USER)" -d "$(DB_NAME)"

.PHONY: db-console-prod
db-console-prod:
	@echo "$(BLUE)== Entering Production PostgreSQL Console ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) exec db psql -U "$(DB_USER)" -d "$(DB_NAME)"

# Logs & Monitoring
.PHONY: logs
.PHONY: logs-dev
logs-dev:
	@echo "$(BLUE)== Viewing Development Logs ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) logs -f

.PHONY: logs-prod
logs-prod:
	@echo "$(BLUE)== Viewing Production Logs ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) logs -f

.PHONY: logs-app-dev
logs-app-dev:
	@echo "$(BLUE)== Viewing Development Application Logs ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) logs -f app

.PHONY: logs-app-prod
logs-app-prod:
	@echo "$(BLUE)== Viewing Production Application Logs ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) logs -f app-prod

.PHONY: logs-db-dev
logs-db-dev:
	@echo "$(BLUE)== Viewing Development Database Logs ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) logs -f db

.PHONY: logs-db-prod
logs-db-prod:
	@echo "$(BLUE)== Viewing Production Database Logs ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) logs -f db

# Container Shell Access
.PHONY: shell-app-dev
shell-app-dev:
	@echo "$(BLUE)== Opening Shell in Development App Container ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) exec app sh

.PHONY: shell-app-prod
shell-app-prod:
	@echo "$(BLUE)== Opening Shell in Production App Container ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) exec app-prod sh

.PHONY: shell-db-dev
shell-db-dev:
	@echo "$(BLUE)== Opening Shell in Development Database Container ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) exec db bash

.PHONY: shell-db-prod
shell-db-prod:
	@echo "$(BLUE)== Opening Shell in Production Database Container ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) exec db bash

# Redis Commands
# Redis Commands
.PHONY: shell-redis-dev
shell-redis-dev:
	@echo "$(BLUE)== Opening Redis CLI in Development Redis Container ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) exec redis redis-cli -a $${REDIS_PASSWORD:-redis}

.PHONY: shell-redis-prod
shell-redis-prod:
	@echo "$(BLUE)== Opening Redis CLI in Production Redis Container ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) exec redis redis-cli -a $${REDIS_PASSWORD:-redis}

.PHONY: logs-redis-dev
logs-redis-dev:
	@echo "$(BLUE)== Viewing Development Redis Logs ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) logs -f redis

.PHONY: logs-redis-prod
logs-redis-prod:
	@echo "$(BLUE)== Viewing Production Redis Logs ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) logs -f redis

.PHONY: redis-monitor-dev
redis-monitor-dev:
	@echo "$(BLUE)== Monitoring Development Redis Commands ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) exec redis redis-cli -a $${REDIS_PASSWORD:-redis} monitor

.PHONY: redis-monitor-prod
redis-monitor-prod:
	@echo "$(BLUE)== Monitoring Production Redis Commands ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) exec redis redis-cli -a $${REDIS_PASSWORD:-redis} monitor

.PHONY: redis-info-dev
redis-info-dev:
	@echo "$(BLUE)== Showing Development Redis Info ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) exec redis redis-cli -a $${REDIS_PASSWORD:-redis} info

.PHONY: redis-info-prod
redis-info-prod:
	@echo "$(BLUE)== Showing Production Redis Info ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) exec redis redis-cli -a $${REDIS_PASSWORD:-redis} info

.PHONY: redis-flush-dev
redis-flush-dev:
	@echo "$(RED)== WARNING: Flushing Development Redis Database ==$(RESET)"
	@read -p "Are you sure? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		$(DOCKER_COMPOSE_DEV) exec redis redis-cli -a $${REDIS_PASSWORD:-redis} flushall; \
		echo "$(RED)Redis database flushed$(RESET)"; \
	fi

# Cleanup Commands
.PHONY: clean-dev
clean-dev:
	@echo "$(BLUE)== Removing Development Containers ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) down

.PHONY: clean-prod
clean-prod:
	@echo "$(BLUE)== Removing Production Containers ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) down

.PHONY: clean-all
clean-all:
	@echo "$(BLUE)== Removing All Docker Artifacts ==$(RESET)"
	docker-compose down -v --rmi all --remove-orphans

# Service Control
.PHONY: pause-dev
pause-dev:
	@echo "$(BLUE)== Pausing Development Services ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) pause

.PHONY: pause-prod
pause-prod:
	@echo "$(BLUE)== Pausing Production Services ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) pause

.PHONY: unpause-dev
unpause-dev:
	@echo "$(BLUE)== Unpausing Development Services ==$(RESET)"
	$(DOCKER_COMPOSE_DEV) unpause

.PHONY: unpause-prod
unpause-prod:
	@echo "$(BLUE)== Unpausing Production Services ==$(RESET)"
	$(DOCKER_COMPOSE_PROD) unpause

# Testing and Linting
.PHONY: test
test:
	@echo "$(BLUE)== Running Tests ==$(RESET)"
	go test -v ./...

.PHONY: lint
lint:
	@echo "$(BLUE)== Running Linter ==$(RESET)"
	golangci-lint run

.PHONY: lint-nilaway
lint-nil:
	@echo "Running nilaway..."
	@$(eval GOPATH=$(shell go env GOPATH))
	@$(eval PATH := $(GOPATH)/bin:$(PATH))
	@nilaway ./... || (echo "Nilaway check failed"; exit 1)
	@echo "Nilaway check passed"

.PHONY: format
format:
	@echo "$(BLUE)== Formatting Code ==$(RESET)"
	go fmt ./...

.PHONY: mock-private
mock-private:
	@echo "$(BLUE)== Starting Private Endpoints Mock Server ==$(RESET)"
	go run ./src/internal/cli/tokenctl.go

.PHONY:  mock-public
mock-public:
	@echo "$(BLUE)== Starting Public Endpoints Mock Server ==$(RESET)"
	cd ./test/mocks/mockingbyrd && bun run dev

.PHONY:  pre-prod
pre-prod:
	@echo "Running pre production environment..."
	@$(eval GOPATH=$(shell go env GOPATH))
	@$(eval PATH := $(GOPATH)/bin:$(PATH))
	@air -c .air.toml

.PHONY:  stop-mock-servers
stop-mock-servers:
	@echo "$(BLUE)== Stopping All Mock Servers ==$(RESET)"
	pkill -f "tokenctl.go" || true
	pkill -f "bun run dev" || true

# Default Target
.DEFAULT_GOAL := help
