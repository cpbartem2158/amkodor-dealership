#!/bin/bash

# Скрипт для создания резервной копии базы данных

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== Создание резервной копии базы данных ===${NC}"

# Загрузка переменных окружения
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
    echo -e "${GREEN}✓ Переменные окружения загружены${NC}"
else
    echo -e "${RED}✗ Файл .env не найден${NC}"
    exit 1
fi

# Создание директории для бэкапов
BACKUP_DIR="backups"
mkdir -p "$BACKUP_DIR"

# Генерация имени файла с текущей датой и временем
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_backup_$TIMESTAMP.sql"

echo -e "${YELLOW}Создание бэкапа: $BACKUP_FILE${NC}"

# Создание резервной копии
PGPASSWORD="$DB_PASSWORD" pg_dump \
    -h "$DB_HOST" \
    -p "$DB_PORT" \
    -U "$DB_USER" \
    -d "$DB_NAME" \
    -F p \
    -b \
    -v \
    -f "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Старые бэкапы удалены${NC}"
fi

echo -e "${GREEN}=== Резервное копирование завершено! ===${NC}"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Резервная копия успешно создана${NC}"

    # Размер файла
    FILE_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo -e "${GREEN}Размер файла: $FILE_SIZE${NC}"

    # Сжатие бэкапа
    echo -e "${YELLOW}Сжатие бэкапа...${NC}"
    gzip "$BACKUP_FILE"

    if [ $? -eq 0 ]; then
        COMPRESSED_FILE="${BACKUP_FILE}.gz"
        COMPRESSED_SIZE=$(du -h "$COMPRESSED_FILE" | cut -f1)
        echo -e "${GREEN}✓ Бэкап сжат: $COMPRESSED_FILE${NC}"
        echo -e "${GREEN}Размер сжатого файла: $COMPRESSED_SIZE${NC}"
    fi

    # Удаление старых бэкапов (старше 30 дней)
    echo -e "${YELLOW}Удаление старых бэкапов (старше 30 дней)...${NC}"
    find "$BACKUP_DIR" -name "*.sql.gz" -type f -mtime +30 -delete

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Старые бэкапы удалены${NC}"
    fi

    # Список существующих бэкапов
    echo -e "${YELLOW}Список бэкапов:${NC}"
    ls -lh "$BACKUP_DIR"

else
    echo -e "${RED}✗ Ошибка при создании резервной копии${NC}"
    exit 1
fi