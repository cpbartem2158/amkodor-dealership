#!/bin/bash

# Скрипт для выполнения миграций базы данных

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Запуск миграций базы данных ===${NC}"

# Загрузка переменных окружения
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
    echo -e "${GREEN}✓ Переменные окружения загружены${NC}"
else
    echo -e "${RED}✗ Файл .env не найден${NC}"
    exit 1
fi

# Путь к миграциям
MIGRATIONS_DIR="internal/database/migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}✗ Директория миграций не найдена: $MIGRATIONS_DIR${NC}"
    exit 1
fi

# Подключение к БД
DB_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo -e "${YELLOW}Подключение к базе данных: ${DB_HOST}:${DB_PORT}/${DB_NAME}${NC}"

# Проверка подключения
psql "$DB_URL" -c "SELECT 1" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Не удалось подключиться к базе данных${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Подключение к базе данных успешно${NC}"

# Выполнение миграций по порядку
for migration in $(ls -v $MIGRATIONS_DIR/*.sql); do
    filename=$(basename "$migration")
    echo -e "${YELLOW}Выполнение миграции: $filename${NC}"

    psql "$DB_URL" -f "$migration"

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Миграция $filename выполнена успешно${NC}"
    else
        echo -e "${RED}✗ Ошибка при выполнении миграции $filename${NC}"
        exit 1
    fi
done

echo -e "${GREEN}=== Все миграции выполнены успешно! ===${NC}"