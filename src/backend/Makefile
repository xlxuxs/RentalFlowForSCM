.PHONY: all proto build test lint clean docker-up docker-down help

# Variables
PROTO_DIR := ./proto
PKG_DIR := ./pkg
SERVICES := auth inventory booking payment notification review gateway

# Default target
all: proto build

# ============================================
# PROTO GENERATION
# ============================================
proto:
	@echo "Generating proto files..."
	@./scripts/generate-proto.sh

# ============================================
# BUILD
# ============================================
build: $(addprefix build-,$(SERVICES))

build-%:
	@echo "Building $*-service..."
	@cd services/$*-service && go build -o bin/$* ./cmd/server

# ============================================
# RUN
# ============================================
run-%:
	@echo "Running $*-service..."
	@cd services/$*-service && go run ./cmd/server

run-all:
	@echo "Starting all services..."
	@for service in $(SERVICES); do \
		(cd services/$$service-service && go run ./cmd/server &); \
	done

# ============================================
# TEST
# ============================================
test:
	@echo "Running all tests..."
	@go test ./... -v -cover

test-%:
	@echo "Running tests for $*-service..."
	@cd services/$*-service && go test ./... -v -cover

test-coverage:
	@echo "Generating coverage report..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

# ============================================
# LINT
# ============================================
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

lint-%:
	@echo "Running linter for $*-service..."
	@cd services/$*-service && golangci-lint run ./...

# ============================================
# DATABASE MIGRATIONS
# ============================================
migrate-up-%:
	@echo "Running migrations for $*-service..."
	@cd services/$*-service && go run ./cmd/migrate up

migrate-down-%:
	@echo "Reverting migrations for $*-service..."
	@cd services/$*-service && go run ./cmd/migrate down

migrate-create-%:
	@echo "Creating migration for $*-service..."
	@cd services/$*-service && go run ./cmd/migrate create $(name)

# ============================================
# DOCKER
# ============================================
docker-up:
	@echo "Starting infrastructure..."
	@docker compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@docker compose ps

docker-down:
	@echo "Stopping infrastructure..."
	@docker compose down

docker-logs:
	@docker compose logs -f

docker-clean:
	@echo "Removing all containers and volumes..."
	@docker compose down -v

docker-build:
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building $$service-service image..."; \
		docker build -t rentalflow/$$service-service:latest ./services/$$service-service; \
	done

# ============================================
# UTILITIES
# ============================================
clean:
	@echo "Cleaning build artifacts..."
	@for service in $(SERVICES); do \
		rm -rf services/$$service-service/bin; \
	done
	@rm -f coverage.out coverage.html

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

# ============================================
# HELP
# ============================================
help:
	@echo "RentalFlow Makefile Commands:"
	@echo ""
	@echo "  make proto          - Generate protobuf files"
	@echo "  make build          - Build all services"
	@echo "  make build-{svc}    - Build specific service (auth, inventory, etc.)"
	@echo "  make run-{svc}      - Run specific service"
	@echo "  make run-all        - Run all services"
	@echo ""
	@echo "  make test           - Run all tests"
	@echo "  make test-{svc}     - Run tests for specific service"
	@echo "  make test-coverage  - Generate coverage report"
	@echo ""
	@echo "  make lint           - Run linter on all code"
	@echo "  make lint-{svc}     - Run linter on specific service"
	@echo ""
	@echo "  make migrate-up-{svc}   - Run migrations"
	@echo "  make migrate-down-{svc} - Revert migrations"
	@echo ""
	@echo "  make docker-up      - Start infrastructure (DBs, Redis, RabbitMQ)"
	@echo "  make docker-down    - Stop infrastructure"
	@echo "  make docker-logs    - View infrastructure logs"
	@echo "  make docker-clean   - Remove containers and volumes"
	@echo "  make docker-build   - Build Docker images"
	@echo ""
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Download dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make help           - Show this help"
