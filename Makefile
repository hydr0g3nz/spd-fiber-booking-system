.PHONY: run test clean build docs help

# Default values
APP_NAME=fiber-booking-system
BUILD_DIR=./build

help:
	@echo "Fiber Booking System - Makefile commands"
	@echo "--------------------------------------"
	@echo "make run              - Run the application"
	@echo "make build            - Build the application"
	@echo "make test             - Run all tests"
	@echo "make test-coverage    - Run tests with coverage report"
	@echo "make clean            - Remove build artifacts"
	@echo "make docs             - Generate Swagger documentation"
	@echo "make lint             - Run linter"
	@echo "make deps             - Download dependencies"

# Run the application
run:
	@echo "Starting application..."
	go run cmd/main.go

# Build the application
build:
	@echo "Building application..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go

# Run all tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Generate Swagger documentation
docs:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/main.go -o docs

# Run linter
lint:
	@echo "Running linter..."
	go vet ./...
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
	else \
		echo "golint not installed. Install with: go install golang.org/x/lint/golint@latest"; \
	fi

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download