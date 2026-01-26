.PHONY: build install uninstall clean test help run dev

# Binary name
BINARY_NAME=ship
# Installation directory
INSTALL_PATH=/usr/local/bin
# Build directory
BUILD_DIR=bin
# Version information
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
all: build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## install: Build and install the binary to system PATH
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	@sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) installed successfully!"
	@echo "You can now run '$(BINARY_NAME)' from anywhere."

## uninstall: Remove the binary from system PATH
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) uninstalled successfully!"

## dev: Build and run the binary (for development)
dev: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

## run: Run the binary directly without building
run:
	@go run .

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -v ./tests/...

## test-cover: Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./tests/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete!"

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify
	@echo "✅ Dependencies ready!"

## tidy: Tidy up go.mod and go.sum
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "✅ Dependencies tidied!"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✅ Vet complete!"

## lint: Run all code quality checks
lint: fmt vet
	@echo "✅ All checks passed!"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /' | column -t -s ':' | sed 's/^/  /'
