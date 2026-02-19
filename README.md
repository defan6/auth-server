# Market - Микросервисная платформа

E-commerce платформа с микросервисной архитектурой на Go.

---

## 📋 Содержание

- [Архитектура](#архитектура)
- [Сервисы](#сервисы)
- [Быстрый старт](#быстрый-старт)
- [Миграции](#миграции)
- [API](#api)
- [Разработка](#разработка)

---

## 🏗️ Архитектура

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│   Frontend      │ ──▶  │   API Gateway   │ ──▶  │   Auth Server   │
│   (React/Vue)   │      │   (Kong/Traefik)│      │   (gRPC)        │
└─────────────────┘      └─────────────────┘      └─────────────────┘
                                                        │
         ┌──────────────────────────────────────────────┘
         │
         ▼
┌─────────────────┐      ┌─────────────────┐
│  Order Service  │ ──▶  │ Product Service │
│  (HTTP/REST)    │      │  (HTTP/gRPC)    │
└─────────────────┘      └─────────────────┘
```

---

## 📦 Сервисы

| Сервис | Порт | Транспорт | Описание |
|--------|------|-----------|----------|
| **auth-server** | 44056 | gRPC | Аутентификация, JWT токены, роли |
| **order-service** | 8090 | HTTP/REST | Управление заказами |
| **product-service** | TBD | HTTP/gRPC | Каталог товаров (в разработке) |

---

## 🚀 Быстрый старт

### Требования

- Go 1.25+
- PostgreSQL 15+
- Docker & Docker Compose
- golang-migrate

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd market
```

### 2. Запуск PostgreSQL

```bash
docker-compose up -d
```

### 3. Запуск миграций

```bash
# PowerShell
.\infrastructure\db\migrate-all.ps1

# Bash
./infrastructure/db/migrate-all.sh
```

### 4. Запуск сервисов

**Auth-server:**
```bash
cd services/auth-server
go run cmd/sso/main.go -config config/local.yaml
```

**Order-service:**
```bash
cd services/order-service
go run cmd/main.go -config config/local.yaml
```

---

## 🗄️ Миграции

### Установка golang-migrate

```bash
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Запуск всех миграций

```bash
# PowerShell
.\infrastructure\db\migrate-all.ps1

# Bash
./infrastructure/db/migrate-all.sh
```

### Миграции для конкретного сервиса

**Auth-server:**
```bash
migrate -path "services/auth-server/cmd/db/migrations" \
        -storage "postgres://postgres:postgres@localhost:5433/users?sslmode=disable" \
        up
```

**Order-service:**
```bash
migrate -path "services/order-service/cmd/migrations" \
        -storage "postgres://postgres:postgres@localhost:5433/orders?sslmode=disable" \
        up
```

### Откат миграций

```bash
# Откатить последнюю
migrate -path "<path>" -storage "<url>" down 1

# Откатить все
migrate -path "<path>" -storage "<url>" down -all
```

---

## 📡 API

### Auth Server (gRPC)

| Метод | Описание |
|-------|----------|
| `Register` | Регистрация пользователя |
| `Login` | Аутентификация, получение JWT |
| `IsAdmin` | Проверка роли администратора |
| `ListUsers` | Список пользователей (admin only) |

### Order Service (HTTP/REST)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/orders` | Создать заказ |
| `GET` | `/api/v1/orders/{id}` | Получить заказ |
| `GET` | `/api/v1/orders` | Список заказов |
| `DELETE` | `/api/v1/orders/{id}` | Отменить заказ |

**Пример создания заказа:**

```bash
curl -X POST http://localhost:8090/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "payment_method": "card",
    "user_id": 1,
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 5, "quantity": 1}
    ]
  }'
```

---

## 🛠️ Разработка

### Структура проекта

```
market/
├── infrastructure/
│   └── db/
│       ├── init/              # Скрипты инициализации БД
│       ├── migrate-all.sh     # Запуск миграций (bash)
│       └── migrate-all.ps1    # Запуск миграций (PowerShell)
├── services/
│   ├── auth-server/
│   │   ├── cmd/
│   │   │   ├── db/
│   │   │   │   ├── init/
│   │   │   │   └── migrations/
│   │   │   └── sso/
│   │   ├── config/
│   │   ├── internal/
│   │   ├── storage/
│   │   └── go.mod
│   ├── order-service/
│   │   ├── cmd/
│   │   │   ├── main.go
│   │   │   └── migrations/
│   │   ├── config/
│   │   ├── internal/
│   │   │   ├── app/http/
│   │   │   ├── config/
│   │   │   ├── database/
│   │   │   ├── domain/
│   │   │   ├── dto/
│   │   │   ├── handler/
│   │   │   ├── mapper/
│   │   │   ├── service/
│   │   │   └── storage/
│   │   └── go.mod
│   └── shared/
│       └── logger/
├── docker-compose.yaml
└── go.work
```

### Конфигурация

Конфигурационные файлы в `config/local.yaml` для каждого сервиса:

**Auth-server:**
```yaml
env: local
grpc:
  port: 44056
db:
  host: localhost
  port: 5432
  name: users
  user: postgres
  password: postgres
token:
  ttl: 10m
  secret: "super-secret"
  issuer: "sso-auth-server"
```

**Order-service:**
```yaml
env: local
server:
  host: localhost
  port: 8090
db:
  host: localhost
  port: 5432
  name: orders
  user: postgres
  password: postgres
```

### Тестирование

```bash
# Запуск тестов для сервиса
cd services/<service-name>
go test ./...

# Запуск с покрытием
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Добавление нового сервиса

1. Создать директорию `services/<service-name>`
2. Инициализировать модуль:
   ```bash
   go mod init github.com/defan6/market/services/<service-name>
   ```
3. Добавить в `go.work`:
   ```bash
   go work use ./services/<service-name>
   ```
4. Создать миграции в `cmd/migrations/`
5. Добавить конфиг в `config/local.yaml`

---

## 📊 Базы данных

| Сервис | БД | Порт | Таблицы |
|--------|----|----|---------|
| auth-server | users | 5432 | users |
| order-service | orders | 5432 | orders, order_items |

---

## 🔐 Безопасность

### JWT Токены

- **Алгоритм:** HMAC-SHA256
- **TTL:** 10 минут (настраивается)
- **Claims:** user_id, email, role, aud (app_id)

### Роли

| Роль | Описание |
|------|----------|
| `user` | Базовые права (создание заказов) |
| `manager` | Управление заказами |
| `admin` | Полный доступ |

---

## 📝 Changelog

### v0.1.0 (2026-02-18)
- ✅ Auth-server: регистрация, логин, JWT
- ✅ Order-service: создание, получение заказов
- ✅ Миграции для всех сервисов
- ✅ Docker Compose для PostgreSQL
- ⏳ Product-service: в разработке

---

## 📄 License

MIT
