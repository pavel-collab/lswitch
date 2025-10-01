# Makefile for lswitch

# Module name (not used in build, informational)
MODULE := yourmodule

# Application version
VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")

# Binary name
BINARY_NAME := lswitch

# Directories
CMD_DIR := cmd/lswitch
BUILD_DIR := bin

# Go parameters
GO := go
GOFLAGS := -v
LDFLAGS := -X main.version=$(VERSION)

.PHONY: all build debug test clean install help

all: test build

## build: Build release binary
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

## debug: Build debug binary (no optimizations)
debug:
	@echo "Building debug version"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -gcflags="all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME)_debug ./$(CMD_DIR)

## test: Run unit tests
test:
	@echo "Running tests..."
	$(GO) test -v ./internal/translator

## test-cover: Run tests with coverage report
test-cover:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./internal/translator
	$(GO) tool cover -html=coverage.out -o coverage.html

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./internal/translator

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## install: Install into $GOPATH/bin
install:
	@echo "Installing..."
	$(GO) install -ldflags "$(LDFLAGS)" ./$(CMD_DIR)

## help: Show help for available commands
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'