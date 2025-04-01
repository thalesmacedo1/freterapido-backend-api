# Variables
DC = docker compose

# Summary of commands
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make test       Run unit tests"
	@echo "  make run        Build and run the Docker applications"
	@echo "  make down       Stop and remove Docker containers and networks"

# Run unit tests
.PHONY: test
test:
	@echo "Running unit tests..."
	@cd api && go test ./...

# Build and run Docker applications
.PHONY: run build test test-unit test-integration
run:
	@echo "Building and running Docker applications..."
	$(DC) up --build

# App variables
APP_NAME=freterapido-backend-api
MAIN_PACKAGE=./api/cmd/api

# Build settings
BUILD_DIR=bin
BINARY_NAME=$(BUILD_DIR)/$(APP_NAME)

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	$(GOTEST) -v ./api/domain/entities/... ./api/application/usecases/... ./api/interfaces/api/... ./api/infrastructure/database/...

# Run integration tests (requires database)
test-integration:
	RUN_INTEGRATION_TESTS=true $(GOTEST) -v ./api/tests/integration/...
