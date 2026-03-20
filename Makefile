.PHONY: migrate-up migrate-down migrate-down-all migrate-version migrate-create run tidy

# Configuration
DB_URL = mysql://root:secret@tcp(localhost:3307)/coffee_pos
MIGRATIONS_PATH = migrations

# Migration targets
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

migrate-down-all:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down -all

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=create_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

# Development targets
run:
	go run cmd/api/main.go

tidy:
	go mod tidy

# Build targets
build:
	go build -o bin/coffee-pos cmd/api/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/coffee-pos-linux cmd/api/main.go

# Test targets
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Docker targets
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-ps:
	docker-compose ps

# Production Docker targets
docker-prod-up:
	docker-compose -f docker-compose.prod.yml up -d

docker-prod-down:
	docker-compose -f docker-compose.prod.yml down

# Clean targets
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Help target
help:
	@echo "Coffee Shop POS - Makefile Commands"
	@echo ""
	@echo "Migration:"
	@echo "  make migrate-up          - Run all migrations"
	@echo "  make migrate-down        - Rollback one migration"
	@echo "  make migrate-down-all    - Rollback all migrations"
	@echo "  make migrate-version     - Show current migration version"
	@echo "  make migrate-create name=<name> - Create new migration"
	@echo ""
	@echo "Development:"
	@echo "  make run                 - Run application"
	@echo "  make tidy                - Run go mod tidy"
	@echo "  make build               - Build application"
	@echo "  make build-linux         - Build for Linux"
	@echo ""
	@echo "Testing:"
	@echo "  make test                - Run tests"
	@echo "  make test-coverage       - Run tests with coverage"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up           - Start development containers"
	@echo "  make docker-down         - Stop development containers"
	@echo "  make docker-logs         - View container logs"
	@echo "  make docker-ps           - List containers"
	@echo "  make docker-prod-up      - Start production containers"
	@echo "  make docker-prod-down    - Stop production containers"
	@echo ""
	@echo "Other:"
	@echo "  make clean               - Clean build artifacts"
	@echo "  make help                - Show this help"
