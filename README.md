# Crypto Price Tracker

Микросервис для отслеживания цен криптовалют с использованием PostgreSQL и Docker Compose.

## Архитектура

Проект построен на принципах Clean Architecture с разделением на слои:

- **Domain**: Сущности, интерфейсы репозиториев
- **Application**: Бизнес-логика, use cases
- **Infrastructure**: Реализация репозиториев, внешние API
- **Delivery**: HTTP handlers, API endpoints

### Структура проекта

```
/cmd
  /api - основной HTTP сервер
  /worker - сервис обновления цен
/internal
  /domain
    /models - сущности
    /repository - интерфейсы хранилищ
  /application
    /services - бизнес-логика
    /dto - объекты передачи данных
  /infrastructure
    /postgres - реализация репозитория
    /coingecko - клиент внешнего API
  /delivery
    /http - обработчики запросов
    /middleware - промежуточное ПО
/migrations - SQL миграции
/docs - документация Swagger
/configs - конфигурационные файлы
/scripts - вспомогательные скрипты
```

## Технологический стек

- Go 1.21+
- PostgreSQL 15+
- Docker + Docker Compose
- Swagger UI
- Viper для конфигурации
- Zap для логирования
- Testify для тестирования

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)

### Запуск с Docker Compose

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd crypto-price-tracker-app
```

2. Создайте файл `.env` на основе `.env.example`:
```bash
cp env.example .env
```

3. Запустите сервисы:
```bash
make up
```

4. Примените миграции базы данных:
```bash
make migrate-up
```

5. Откройте Swagger UI: http://localhost:8080/swagger/

### Локальная разработка

1. Установите зависимости:
```bash
go mod download
```

2. Запустите PostgreSQL:
```bash
make db-up
```

3. Примените миграции:
```bash
make migrate-up
```

4. Запустите API сервер:
```bash
make run-api
```

5. В отдельном терминале запустите worker:
```bash
make run-worker
```

## API Endpoints

### Добавить криптовалюту в отслеживание
```bash
POST /currency/add
Content-Type: application/json

{
  "symbol": "bitcoin",
  "interval": 60
}
```

### Удалить криптовалюту из отслеживания
```bash
POST /currency/remove
Content-Type: application/json

{
  "symbol": "bitcoin"
}
```

### Получить цену криптовалюты
```bash
GET /currency/price?coin=bitcoin&timestamp=1640995200
```

## Управление проектом

### Makefile команды

- `make up` - запуск всех сервисов
- `make down` - остановка всех сервисов
- `make logs` - просмотр логов
- `make migrate-up` - применение миграций
- `make migrate-down` - откат миграций
- `make test` - запуск тестов
- `make build` - сборка приложения
- `make clean` - очистка сборки

### Переменные окружения

- `DB_HOST` - хост PostgreSQL (по умолчанию: localhost)
- `DB_PORT` - порт PostgreSQL (по умолчанию: 5432)
- `DB_NAME` - имя базы данных (по умолчанию: crypto_tracker)
- `DB_USER` - пользователь базы данных (по умолчанию: postgres)
- `DB_PASSWORD` - пароль базы данных (по умолчанию: password)
- `API_PORT` - порт API сервера (по умолчанию: 8080)
- `WORKER_INTERVAL` - интервал обновления цен в секундах (по умолчанию: 60)
- `COINGECKO_API_URL` - URL API CoinGecko (по умолчанию: https://api.coingecko.com/api/v3)

## Разработка

### Добавление новой функциональности

1. Создайте модель в `/internal/domain/models`
2. Определите интерфейс репозитория в `/internal/domain/repository`
3. Реализуйте бизнес-логику в `/internal/application/services`
4. Создайте HTTP handler в `/internal/delivery/http`
5. Добавьте тесты в соответствующие директории

### Тестирование

```bash
# Запуск всех тестов
make test

# Запуск тестов с покрытием
make test-coverage

# Запуск конкретного теста
go test ./internal/application/services -v
```

## Мониторинг

- Логи доступны через `docker-compose logs`
- Swagger UI: http://localhost:8080/swagger/
- Health check: http://localhost:8080/health

## Лицензия

MIT
