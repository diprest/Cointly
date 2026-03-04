# 🚀 Go8 Crypto Exchange Platform

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Coverage-36%25+-success)](docs/COMPLIANCE_REPORT.md)

Высокопроизводительная микросервисная платформа для торговли криптовалютами с поддержкой лимитных ордеров, портфельного управления и бинарных опционов.

---

## 📋 Содержание

- [Обзор](#-обзор)
- [Архитектура](#-архитектура)
- [Технологический стек](#-технологический-стек)
- [Быстрый старт](#-быстрый-старт)
- [API Документация](#-api-документация)
- [Тестирование](#-тестирование)
- [CI/CD](#-cicd)
- [Отказоустойчивость](#-отказоустойчивость)
- [Структура проекта](#-структура-проекта)
- [Разработка](#-разработка)
- [Соответствие ТЗ](#-соответствие-тз)

---

## 🎯 Обзор

**Go8 Crypto Exchange** — это полнофункциональная платформа для торговли криптовалютами, построенная на микросервисной архитектуре. Проект демонстрирует современные практики разработки на Go, включая:

- ✅ Микросервисная архитектура с четким разделением ответственности
- ✅ Event-driven коммуникация через Apache Kafka
- ✅ Отказоустойчивость с Circuit Breaker паттерном
- ✅ Кэширование с Redis для высокой производительности
- ✅ Полная контейнеризация с Docker
- ✅ CI/CD пайплайн с автоматическим тестированием
- ✅ Покрытие тестами >30% (требование ТЗ)

### Основные возможности

- 📊 **Торговля**: Размещение и исполнение лимитных ордеров
- 💰 **Портфолио**: Управление балансами и активами
- 📈 **Рыночные данные**: Реал-тайм обновления цен криптовалют
- 🎲 **Бинарные опционы**: Ставки на направление движения цены
- 🔐 **Аутентификация**: JWT-based авторизация
- 🌐 **Web UI**: Интуитивный фронтенд интерфейс

---

## 🏗 Архитектура

### Диаграмма компонентов

┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTP (Port 80)
       ▼
┌─────────────────────────────────────────┐
│            Nginx (Reverse Proxy)         │
│  - Single Entry Point                    │
│  - Load Balancing                        │
└──────┬──────────────────────────────────┘
       │
       ├─────────────► Frontend (React/Vue)
       │
       ▼
┌─────────────────────────────────────────┐
│          API Gateway (Port 8000)         │
│  - JWT Validation                        │
│  - Circuit Breaker                       │
│  - Request Routing                       │
└──────┬──────────────────────────────────┘
       │
       ├──────────► Auth Service (Port 8081)
       │             - User Registration
       │             - JWT Issuance
       │
       ├──────────► Market Data Service (Port 8001)
       │             - Price Updates
       │             - Kafka Producer
       │             - Redis Cache
       │
       ├──────────► Trading Service (Port 8082)
       │             - Order Matching
       │             - Kafka Consumer
       │             - Order Book
       │
       ├──────────► Portfolio Service (Port 8083)
       │             - Balance Management
       │             - Fund Locking
       │
       └──────────► Bets Service (Port 8084)
                     - Binary Options
                     - Bet Resolution

┌─────────────────────────────────────────┐
│          Infrastructure Layer            │
├─────────────────────────────────────────┤
│  PostgreSQL (Port 5432)                  │
│  - 5 separate databases                  │
│                                          │
│  Redis (Port 6379)                       │
│  - Price caching                         │
│                                          │
│  Kafka + Zookeeper                       │
│  - Event streaming                       │
└─────────────────────────────────────────┘

### Микросервисы

| Сервис | Порт | Описание | База данных |
|--------|------|----------|-------------|
| **Nginx** | 80 | Единая точка входа, reverse proxy | - |
| **Frontend** | 3000* | Web интерфейс | - |
| **API Gateway** | 8000* | Маршрутизация, авторизация, Circuit Breaker | - |
| **Auth Service** | 8081* | Регистрация, аутентификация, JWT | `auth_db` |
| **Market Data** | 8001* | Получение и распространение цен | `market_db` |
| **Trading Service** | 8082* | Матчинг ордеров, торговая логика | `trading_db` |
| **Portfolio Service** | 8083* | Управление балансами | `portfolio_db` |
| **Bets Service** | 8084* | Бинарные опционы | `bets_db` |

*\* Порты доступны только внутри Docker-сети. Наружу проброшен только порт 80 (Nginx).*

---

## 🛠 Технологический стек

### Backend
- **Язык**: Go 1.24
- **Web Framework**: Chi Router
- **База данных**: PostgreSQL 15
- **Кэш**: Redis 7
- **Message Broker**: Apache Kafka 7.4.4
- **Authentication**: JWT (golang-jwt/jwt)
- **Decimal Math**: shopspring/decimal

### Infrastructure
- **Контейнеризация**: Docker, Docker Compose
- **Reverse Proxy**: Nginx
- **CI/CD**: GitLab CI/CD
- **Testing**: testify, pgxmock

### Паттерны и практики
- Clean Architecture
- Repository Pattern
- Circuit Breaker (gobreaker)
- Event-Driven Architecture
- Dependency Injection
- Structured Logging (slog)

---

## 🚀 Быстрый старт

### Требования

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.24+ (опционально, для локальной разработки)

### Установка и запуск

#### 1. Клонирование репозитория

```bash
git clone https://gitlab.crja72.ru/golang/2025/autumn/projects/go8/go8_project.git
cd go8_project

#### 2. Настройка переменных окружения

# Файл .env уже настроен с дефолтными значениями
# При необходимости отредактируйте его

#### 3. Запуск проекта

Production режим (только Nginx на порту 80):

docker-compose up --build -d

Windows Dev режим (с пробросом портов для отладки):

.\run_windows.ps1

#### 4. Проверка работоспособности

# Проверка здоровья сервисов
curl http://localhost/health

# Проверка API Gateway
curl http://localhost/api/v1/market/prices

### Доступ к приложению

- Frontend: http://localhost (Production) или http://localhost:3000 (Dev)
- API: http://localhost/api/v1/
- Swagger UI: Откройте docs/swagger.yaml в [Swagger Editor](https://editor.swagger.io/)

---

## 📚 API Документация

### Основные эндпоинты

#### Аутентификация

# Регистрация
POST /api/v1/auth/register
Content-Type: application/json

{
  "login": "user123",
  "password": "securepass"
}

# Вход
POST /api/v1/auth/login
Content-Type: application/json

{
  "login": "user123",
  "password": "securepass"
}

# Ответ
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user_id": 1
}

#### Рыночные данные

# Получить текущие цены
GET /api/v1/market/prices
Authorization: Bearer <token>

# Ответ
[
  {
    "symbol": "BTCUSDT",
    "price": 50000.00,
    "timestamp": "2025-12-17T18:00:00Z"
  }
]

#### Торговля

# Разместить лимитный ордер
POST /api/v1/trading/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "symbol": "BTCUSDT",
  "side": "buy",
  "price": 50000,
  "amount": 0.1
}

# Получить ордера пользователя
GET /api/v1/trading/orders
Authorization: Bearer <token>

#### Портфолио

# Получить баланс
GET /api/v1/portfolio
Authorization: Bearer <token>

# Ответ
{
  "balances": [
    {
      "asset": "USDT",
      "available": 10000.00,
      "locked": 500.00
    }
  ]
}

### Полная документация

Полная OpenAPI спецификация доступна в файле [docs/swagger.yaml](docs/swagger.yaml).

---

## 🧪 Тестирование

### Запуск всех тестов

# PowerShell
.\run_tests.ps1

# Или напрямую
.\full_test.ps1

### Запуск тестов для конкретного сервиса

# Auth Service
cd services/auth-service
go test ./... -v -cover

# Trading Service
cd services/trading-service
go test ./internal/service -v -coverprofile=coverage.out

# Просмотр покрытия
go tool cover -html=coverage.out

### Текущее покрытие тестами

| Сервис | Покрытие | Статус |
|--------|----------|--------|
| Portfolio Service | 50.0% | ✅ |
| Trading Service | 36.1% | ✅ |
| Bets Service | >30% | ✅ |
| Auth Service | >30% | ✅ |
| Market Data Service | >30% | ✅ |
| API Gateway | Исключен | ✅ |

Требование ТЗ: ≥30% — ВЫПОЛНЕНО ✅

### Типы тестов

- Unit тесты: Тестирование бизнес-логики с моками
- Integration тесты: Тестирование с реальными зависимостями (pgxmock)
- Handler тесты: Тестирование HTTP эндпоинтов

---

## 🔄 CI/CD

### GitLab CI/CD Pipeline

Проект включает полностью настроенный CI/CD пайплайн (.gitlab-ci.yml):

Stages:
  1. Lint    → golangci-lint для всех сервисов
  2. Test    → Запуск тестов + проверка покрытия ≥30%
  3. Build   → Сборка бинарников

### Автоматические проверки

- ✅ Линтинг кода (golangci-lint)
- ✅ Запуск тестов для каждого микросервиса
- ✅ Проверка покрытия (пайплайн падает если <30%)
- ✅ Сборка приложения (CGO_ENABLED=0 для статических бинарников)
- ✅ Параллельное выполнение для ускорения

### Запуск локально

# Линтинг
golangci-lint run ./...

# Тесты с проверкой покрытия
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total

# Сборка
go build -o bin/service ./cmd/main.go

---

## 🛡 Отказоустойчивость

### Circuit Breaker

API Gateway использует паттерн Circuit Breaker для предотвращения каскадных сбоев:

// services/api-gateway/internal/middleware/circuit_breaker.go
Settings:
  - MaxRequests: 5 (в half-open состоянии)
  - Interval: 60s (сброс счетчиков)
  - Timeout: 30s (время в open состоянии)
  - Threshold: 50% failure ratio

Состояния:
- Closed: Нормальная работа
- Open: Сервис недоступен, запросы отклоняются
- Half-Open: Проверка восстановления

### Retry механизмы

- Подключение к PostgreSQL с повторными попытками
- Подключение к Kafka с экспоненциальной задержкой
- GitLab CI: retry: 2 для нестабильных тестов

### Graceful Shutdown

Все сервисы корректно завершают работу при получении SIGTERM/SIGINT:
- Завершение обработки текущих запросов
- Закрытие соединений с БД
- Отключение от Kafka

---

## 📁 Структура проекта

go8_project/
├── .gitlab-ci.yml              # CI/CD конфигурация
├── docker-compose.yml          # Production конфигурация
├── docker-compose.windows.yml  # Dev конфигурация (Windows)
├── Makefile                    # Команды для разработки
├── run_windows.ps1             # Скрипт запуска (Windows)
├── run_tests.ps1               # Скрипт тестирования
├── full_test.ps1               # Полное тестирование
│
├── docs/
│   ├── swagger.yaml            # OpenAPI спецификация
│   ├── COMPLIANCE_REPORT.md    # Отчет о соответствии ТЗ
│   └── ARCHITECTURE.md         # Архитектурная документация
│
├── nginx/
│   └── nginx.conf              # Конфигурация Nginx
│
├── migrations/
│   └── init.sql                # Инициализация БД
│
├── frontend/                   # Web интерфейс
│   ├── Dockerfile
│   ├── package.json
│   └── src/
│
└── services/
    ├── api-gateway/
    │   ├── cmd/main.go
    │   ├── internal/
    │   │   ├── middleware/
    │   │   │   └── circuit_breaker.go
    │   │   ├── proxy/
    │   │   └── router/
    │   ├── Dockerfile
    │   └── go.mod
    │
    ├── auth-service/
    │   ├── cmd/main.go
    │   ├── internal/
    │   │   ├── config/
    │   │   ├── db/
    │   │   ├── models/
    │   │   ├── repository/
    │   │   ├── service/
    │   │   └── transport/http/
    │   ├── Dockerfile
    │   └── go.mod
    │
    ├── market-data-service/
    │   ├── cmd/main.go
    │   ├── internal/
    │   │   ├── broker/          # Kafka producer
    │   │   ├── service/
    │   │   ├── storage/         # Redis + PostgreSQL
    │   │   └── worker/          # Price updater
    │   ├── Dockerfile
    │   └── go.mod
    │
    ├── trading-service/
    │   ├── cmd/main.go
    │   ├── internal/
    │   │   ├── broker/          # Kafka consumer
    │   │   ├── service/
    │   │   ├── storage/
    │   │   └── worker/          # Order matcher
    │   ├── Dockerfile
    │   └── go.mod
    │
    ├── portfolio-service/
    │   ├── cmd/main.go
    │   ├── internal/
    │   │   ├── service/
    │   │   └── storage/
    │   ├── Dockerfile
    │   └── go.mod
    │
    └── bets-service/
        ├── cmd/main.go
        ├── internal/
        │   ├── clients/         # HTTP clients
        │   ├── service/
        │   ├── storage/
        │   └── worker/          # Bet resolver
        ├── Dockerfile
        └── go.mod

---

## 💻 Разработка

### Локальная разработка

# 1. Установка зависимостей
cd services/auth-service
go mod download

# 2. Запуск сервиса локально
go run cmd/main.go

# Форматирование
go fmt ./...

# Линтинг
golangci-lint run ./...

# Vet
go vet ./...

### Логирование

Все сервисы используют структурированное логирование (log/slog):

slog.Info("Order placed", 
    "user_id", userID, 
    "order_id", orderID, 
    "symbol", symbol)

slog.Error("Database error", 
    "error", err, 
    "operation", "CreateOrder")

---


## 📄 Лицензия

Этот проект создан в образовательных целях.


```