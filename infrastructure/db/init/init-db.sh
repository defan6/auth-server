#!/bin/bash
set -euo pipefail

export PGUSER=${POSTGRES_USER:-postgres}
export PGPASSWORD=${POSTGRES_PASSWORD:-postgres}

# Базы данных для каждого микросервиса
AUTH_DB_NAME="users"
ORDER_DB_NAME="orders"

# Функция проверки и создания БД
create_db_if_not_exists() {
    local db_name=$1
    if ! psql -lqt | cut -d \| -f 1 | grep -qw "$db_name"; then
        echo "Creating database $db_name..."
        createdb "$db_name"
    else
        echo "Database $db_name already exists"
    fi
}

# Создаем БД для auth-server
create_db_if_not_exists "$AUTH_DB_NAME"

# Создаем БД для order-service
create_db_if_not_exists "$ORDER_DB_NAME"

echo "All databases initialized"
