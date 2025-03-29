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
.PHONY: run
run:
	@echo "Building and running Docker applications..."
	$(DC) up --build

# Stop and remove containers
.PHONY: down
down:
	@echo "Stopping and removing Docker containers..."
	$(DC) down