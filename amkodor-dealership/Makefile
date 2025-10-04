.PHONY: help build run test clean migrate seed docker-up docker-down docker-logs

# Переменные
APP_NAME=amkodor-dealership
MAIN_PATH=./cmd/api/main.go
BINARY_NAME=amkodor-app
MIGRATIONS_PATH=./internal/database/migrations

# Цвета для вывода
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-15s${NC} %s\n", $$1, $$2}'

build: ## Собрать приложение
	@echo "Сборка приложения..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Приложение собрано: $(BINARY_NAME)"

run: ## Запустить приложение
	@echo "Запуск приложения..."
	@go run $(MAIN_PATH)

dev: ## Запустить в режиме разработки с hot reload (требует air)
	@echo "Запуск в режиме разработки..."
	@air

test: ## Запустить тесты
	@echo "Запуск тестов..."
	@go test -v ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "Запуск тестов с покрытием..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Отчет о покрытии: coverage.html"

clean: ## Очистить build артефакты
	@echo "Очистка..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -rf ./uploads/*
	@rm -rf ./exports/*
	@echo "Очистка завершена"

deps: ## Установить зависимости
	@echo "Установка зависимостей..."
	@go mod download
	@go mod tidy
	@echo "Зависимости установлены"

migrate: ## Применить миграции к базе данных
	@echo "Применение миграций..."
	@psql -h localhost -U postgres -d amkodor_db -f $(MIGRATIONS_PATH)/001_create_tables.sql
	@psql -h localhost -U postgres -d amkodor_db -f $(MIGRATIONS_PATH)/002_create_views.sql
	@psql -h localhost -U postgres -d amkodor_db -f $(MIGRATIONS_PATH)/003_create_functions.sql
	@psql -h localhost -U postgres -d amkodor_db -f $(MIGRATIONS_PATH)/004_create_procedures.sql
	@psql -h localhost -U postgres -d amkodor_db -f $(MIGRATIONS_PATH)/005_create_triggers.sql
	@echo "Миграции применены"

seed: ## Заполнить БД тестовыми данными
	@echo "Заполнение БД тестовыми данными..."
	@psql -h localhost -U postgres -d amkodor_db -f $(MIGRATIONS_PATH)/006_seed_data.sql
	@echo "Данные добавлены"

db-create: ## Создать базу данных
	@echo "Создание базы данных..."
	@createdb -h localhost -U postgres amkodor_db || true
	@echo "База данных создана"

db-drop: ## Удалить базу данных
	@echo "Удаление базы данных..."
	@dropdb -h localhost -U postgres amkodor_db || true
	@echo "База данных удалена"

db-reset: db-drop db-create migrate seed ## Пересоздать БД с нуля
	@echo "База данных пересоздана"

docker-build: ## Собрать Docker образ
	@echo "Сборка Docker образа..."
	@docker-compose build
	@echo "Docker образ собран"

docker-up: ## Запустить контейнеры
	@echo "Запуск контейнеров..."
	@docker-compose up -d
	@echo "Контейнеры запущены"
	@echo "Приложение: http://localhost:8080"
	@echo "PgAdmin: http://localhost:5050"

docker-down: ## Остановить контейнеры
	@echo "Остановка контейнеров..."
	@docker-compose down
	@echo "Контейнеры остановлены"

docker-logs: ## Показать логи контейнеров
	@docker-compose logs -f

docker-restart: docker-down docker-up ## Перезапустить контейнеры

docker-clean: ## Удалить контейнеры и volumes
	@echo "Удаление контейнеров и volumes..."
	@docker-compose down -v
	@echo "Очистка завершена"

lint: ## Запустить линтер
	@echo "Запуск линтера..."
	@golangci-lint run ./...

format: ## Форматировать код
	@echo "Форматирование кода..."
	@gofmt -s -w .
	@go mod tidy
	@echo "Код отформатирован"

install-tools: ## Установить инструменты для разработки
	@echo "Установка инструментов..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Инструменты установлены"

setup: deps db-create migrate seed ## Первичная настройка проекта
	@echo "Проект настроен и готов к работе!"
	@echo "Запустите 'make run' для старта приложения"

# Значение по умолчанию
.DEFAULT_GOAL := help