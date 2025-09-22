# Makefile for layout_switcher

# Имя модуля
MODULE := yourmodule

# Версия приложения
VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")

# Имя бинарного файла
BINARY_NAME := layout_switcher

# Директории
CMD_DIR := cmd/layout_switcher
BUILD_DIR := bin

# Go параметры
GO := go
GOFLAGS := -v
LDFLAGS := -X main.version=$(VERSION)

.PHONY: all build debug test clean install help

all: test build

## build: Сборка релизной версии
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

## debug: Сборка с отладочной информацией
debug:
	@echo "Building debug version"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -gcflags="all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME)_debug ./$(CMD_DIR)

## test: Запуск юнит-тестов
test:
	@echo "Running tests..."
	$(GO) test -v ./internal/translator

## test-cover: Запуск тестов с измерением покрытия
test-cover:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./internal/translator
	$(GO) tool cover -html=coverage.out -o coverage.html

## bench: Запуск бенчмарков
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./internal/translator

## clean: Очистка артефактов сборки
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## install: Установка в $GOPATH/bin
install:
	@echo "Installing..."
	$(GO) install -ldflags "$(LDFLAGS)" ./$(CMD_DIR)

## help: Показать справку по командам
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'