# Colors for pretty output
BLUE := $(shell printf "\033[36m")
GREEN := $(shell printf "\033[32m")
YELLOW := $(shell printf "\033[33m")
RED := $(shell printf "\033[31m")
RESET := $(shell printf "\033[0m")

# Default shell
SHELL := /bin/bash

# Include .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: help
help:
	@echo "$(BLUE)Available commands:$(RESET)"
	@echo "$(GREEN)Development Setup:$(RESET)"
	@echo "  make init          - Initialize development environment"
	@echo "  make deps          - Install/Update dependencies"
	@echo "$(GREEN)Development Commands:$(RESET)"
	@echo "  make start         - Start development environment"
	@echo "  make stop          - Stop development environment"
	@echo "  make restart       - Restart development environment"
	@echo "$(GREEN)Testing & Linting:$(RESET)"
	@echo "  make test          - Run tests"
	@echo "  make lint          - Run linter"
	@echo "  make format        - Format code"
	@echo "$(GREEN)Build Commands:$(RESET)"
	@echo "  make build         - Build Docker images"
	@echo "  make rebuild       - Force rebuild Docker images"
	@echo "  make build-nocache - Build Docker images without cache"
	@echo "$(GREEN)Database Commands:$(RESET)"
	@echo "  make db-status     - Check database connection"
	@echo "  make db-console    - Enter PostgreSQL console"
	@echo "  make db-migrate    - Run database migrations"
	@echo "$(GREEN)Logs & Monitoring:$(RESET)"
	@echo "  make logs          - View all logs"
	@echo "  make logs-app      - View application logs"
	@echo "  make logs-db       - View database logs"
	@echo "$(GREEN)Container Commands:$(RESET)"
	@echo "  make shell-app     - Shell into app container"
	@echo "  make shell-db      - Shell into database container"
	@echo "$(GREEN)Cleanup Commands:$(RESET)"
	@echo "  make clean         - Remove containers"
	@echo "  make clean-volumes - Remove containers and volumes"
	@echo "  make clean-all     - Remove all Docker artifacts"
	@echo "$(GREEN)Service Control:$(RESET)"
	@echo "  make pause         - Pause services"
	@echo "  make unpause       - Unpause services"

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
.PHONY: start
start:
	@echo "$(BLUE)== Starting Development Environment ==$(RESET)"
	docker-compose up -d

.PHONY: stop
stop:
	@echo "$(BLUE)== Stopping Development Environment ==$(RESET)"
	docker-compose stop

.PHONY: restart
restart:
	@echo "$(BLUE)== Restarting Development Environment ==$(RESET)"
	docker-compose restart

# Build Commands
.PHONY: build
build:
	@echo "$(BLUE)== Building Docker Images ==$(RESET)"
	docker-compose build

.PHONY: rebuild
rebuild:
	@echo "$(BLUE)== Rebuilding Docker Images ==$(RESET)"
	docker-compose build --force-rm

.PHONY: build-nocache
build-nocache:
	@echo "$(BLUE)== Building Docker Images Without Cache ==$(RESET)"
	docker-compose build --no-cache

# Database Commands
.PHONY: db-status
db-status:
	@echo "$(BLUE)== Checking Database Connection ==$(RESET)"
	@docker-compose exec db pg_isready -U "$(DB_USER)" -d "$(DB_NAME)" || echo "$(RED)Database is not ready$(RESET)"

.PHONY: db-console
db-console:
	@echo "$(BLUE)== Entering PostgreSQL Console ==$(RESET)"
	docker-compose exec db psql -U "$(DB_USER)" -d "$(DB_NAME)"

.PHONY: db-migrate
db-migrate:
	@echo "$(BLUE)== Running Database Migrations ==$(RESET)"
	docker-compose exec app go run scripts/migrations/*.go

# Logs & Monitoring
.PHONY: logs
logs:
	@echo "$(BLUE)== Viewing All Logs ==$(RESET)"
	docker-compose logs -f

.PHONY: logs-app
logs-app:
	@echo "$(BLUE)== Viewing Application Logs ==$(RESET)"
	docker-compose logs -f app

.PHONY: logs-db
logs-db:
	@echo "$(BLUE)== Viewing Database Logs ==$(RESET)"
	docker-compose logs -f db

# Container Shell Access
.PHONY: shell-app
shell-app:
	@echo "$(BLUE)== Opening Shell in App Container ==$(RESET)"
	docker-compose exec app sh

.PHONY: shell-db
shell-db:
	@echo "$(BLUE)== Opening Shell in Database Container ==$(RESET)"
	docker-compose exec db bash

# Cleanup Commands
.PHONY: clean
clean:
	@echo "$(BLUE)== Removing Containers ==$(RESET)"
	docker-compose down

.PHONY: clean-volumes
clean-volumes:
	@echo "$(BLUE)== Removing Containers and Volumes ==$(RESET)"
	docker-compose down -v

.PHONY: clean-all
clean-all:
	@echo "$(BLUE)== Removing All Docker Artifacts ==$(RESET)"
	docker-compose down -v --rmi all --remove-orphans

# Service Control
.PHONY: pause
pause:
	@echo "$(BLUE)== Pausing Services ==$(RESET)"
	docker-compose pause

.PHONY: unpause
unpause:
	@echo "$(BLUE)== Unpausing Services ==$(RESET)"
	docker-compose unpause

# Testing and Linting
.PHONY: test
test:
	@echo "$(BLUE)== Running Tests ==$(RESET)"
	go test -v ./...

.PHONY: lint
lint:
	@echo "$(BLUE)== Running Linter ==$(RESET)"
	golangci-lint run

.PHONY: format
format:
	@echo "$(BLUE)== Formatting Code ==$(RESET)"
	go fmt ./...

# Default Target
.DEFAULT_GOAL := help
