
MAIN_PACKAGE_PATH = ./
GO_LINT_VERSION = 1.64.8
BINARY_NAME = export
BUILD_DIR = ./bin

.DEFAULT_GOAL := build

.PHONY: run build up down lint lint-fix test clean

run: build lint test ## Сборка, проверка и тестирование

## build: Сборка приложения
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe $(MAIN_PACKAGE_PATH)

## Поднимает (запускает) окружение для работы приложения
up:
	docker compose up -d

## Отключает окружение для работы приложения
down:
	docker compose down --remove-orphans

## lint: Запуск линтера
lint:
	@echo "Linter is checking..."
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GO_LINT_VERSION) run

## lint-fix: Запуск линтера с фиксом
lint-fix:
	@echo "Linter is fixing..."
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GO_LINT_VERSION) run --fix

## test: Запуск тестов
test:
	@echo "Running tests..."
	go test -count=1 -v ./...

## clean: Удаление бинарных файлов
clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)*

