#!/bin/bash
set -euo pipefail

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Running database migrations for all services${NC}"
echo -e "${GREEN}========================================${NC}"

# Параметры подключения к БД
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

# Функция для запуска миграций
run_migrations() {
    local service_name=$1
    local db_name=$2
    local migrations_path=$3
    
    echo -e "\n${YELLOW}>>> Running migrations for ${service_name} (database: ${db_name})${NC}"
    
    if [ ! -d "$migrations_path" ]; then
        echo -e "${RED}   Error: Migrations path not found: ${migrations_path}${NC}"
        return 1
    fi
    
    migrate -path "$migrations_path" \
            -storage "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
            up
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}   ✓ Migrations completed for ${service_name}${NC}"
    else
        echo -e "${RED}   ✗ Migrations failed for ${service_name}${NC}"
        return 1
    fi
}

# Получаем директорию проекта
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Запускаем миграции для каждого сервиса
echo -e "\n${GREEN}Database connection:${NC} ${DB_HOST}:${DB_PORT}"

# Auth-server миграции
run_migrations \
    "auth-server" \
    "users" \
    "${PROJECT_ROOT}/services/auth-server/cmd/db/migrations"

# Order-service миграции
run_migrations \
    "order-service" \
    "orders" \
    "${PROJECT_ROOT}/services/order-service/cmd/migrations"

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}  All migrations completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
