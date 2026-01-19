OUTPUT:=./bin/app.exe
GO_LINT_VERSION=1.64.8

GO_FILE:=./main.go

.DEFAULT_GOAL := build

.PHONY: run
run: build lint test ## Сборка, проверка и тестирование

.PHONY: up
up: ## Поднимает (запускает) окружение для работы приложения
	docker compose up -d

.PHONY: down
down: ## Отключает окружение для работы приложения
	docker compose down --remove-orphans

.PHONY: lint
lint: ## Запуск линтера
	@echo "Linter is checking..."
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v${GO_LINT_VERSION} run

.PHONY: lint-fix
lint-fix: ## Запуск линтера с фиксом
	@echo "Linter is fixing..."
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v${GO_LINT_VERSION} run --fix

.PHONY: build
build: ## Сборка приложения
	@echo "Building..."
	go build -o ${OUTPUT} ${GO_FILE}

.PHONY: test
test: ## Запуск тестов
	@echo "Running tests..."
	go test -count=1 -v ./...

.PHONY: clean
clean: ## Очистка
	@echo "Cleaning..."
	go clean
	rm -f $(OUTPUT)*

