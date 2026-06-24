.PHONY: help dev stop logs migrate seed lint test

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ===== DEV =====
dev: ## Start all services with docker-compose
	docker-compose up -d

stop: ## Stop all services
	docker-compose down

logs: ## View logs
	docker-compose logs -f

# ===== DATABASE =====
migrate-up: ## Run database migrations
	cd backend && go run cmd/migrate/main.go up

migrate-down: ## Rollback last migration
	cd backend && go run cmd/migrate/main.go down

migrate-create: ## Create new migration (usage: make migrate-create name=create_users)
	cd backend/migrations && migrate create -ext sql -dir . -seq $(name)

seed: ## Seed database with sample data
	cd backend && go run cmd/seed/main.go

# ===== BACKEND =====
backend-api: ## Run backend API server
	cd backend && go run cmd/api/main.go

backend-realtime: ## Run backend realtime gateway
	cd backend && go run cmd/realtime/main.go

backend-test: ## Run backend tests
	cd backend && go test ./... -v

# ===== AI SERVICE =====
ai-service: ## Run AI service
	cd ai-service && uvicorn app.main:app --reload --port 8000

ai-test: ## Run AI service tests
	cd ai-service && pytest tests/ -v

# ===== FRONTEND =====
frontend-dev: ## Run frontend dev server
	cd frontend && npm run dev

frontend-build: ## Build frontend for production
	cd frontend && npm run build

frontend-test: ## Run frontend tests
	cd frontend && npm run test

# ===== LINT =====
lint: ## Lint all projects
	cd backend && golangci-lint run ./...
	cd frontend && npm run lint
	cd ai-service && flake8 app/

# ===== CLEAN =====
clean: ## Remove all containers and volumes
	docker-compose down -v
	rm -rf frontend/node_modules
	rm -rf ai-service/.venv
