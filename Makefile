.PHONY: help build test lint clean install-deps examples

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build all examples"
	@echo "  test         - Run all tests"
	@echo "  lint         - Run go vet and go fmt"
	@echo "  clean        - Clean build artifacts"
	@echo "  install-deps - Install development dependencies"
	@echo "  examples     - Build example applications"

# Build targets
build: examples

examples:
	@echo "Building examples..."
	@go build -o bin/web_server examples/web_server.go
	@go build -o bin/cli_example examples/cli_example.go
	@echo "Examples built in bin/ directory"

# Test targets
test:
	@echo "Running tests..."
	@go test ./pkg/civicauth/... -v

test-coverage:
	@echo "Running tests with coverage..."
	@go test ./pkg/civicauth/... -v -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality
lint:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Checking formatting..."
	@gofmt -l . | tee /tmp/gofmt.out
	@test ! -s /tmp/gofmt.out || (echo "Code is not formatted correctly. Run 'go fmt ./...' to fix." && exit 1)
	@echo "Linting complete"

fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Dependencies
install-deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Cleaning
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Development
dev: install-deps lint test

# CI targets
ci: install-deps lint test

# Run examples with environment variables
run-web:
	@echo "Starting web server example..."
	@echo "Make sure to set CIVIC_CLIENT_ID, CIVIC_CLIENT_SECRET, and CIVIC_ISSUER environment variables"
	@bin/web_server || make build && bin/web_server

run-cli:
	@echo "Running CLI example..."
	@echo "Make sure to set CIVIC_CLIENT_ID, CIVIC_CLIENT_SECRET, and CIVIC_ISSUER environment variables"
	@bin/cli_example || make build && bin/cli_example
