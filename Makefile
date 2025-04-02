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
	@go test ./...

.PHONY: run
run:
	@echo "Building and running Docker applications..."
	$(DC) up --build

.PHONY: down
down:
	@echo "Stopping and removing Docker containers..."
	$(DC) down