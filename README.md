# Crypto Price Tracker

Микросервис для отслеживания цен криптовалют с использованием PostgreSQL и Docker Compose.

## 🏗️ Архитектура

Проект построен по принципам **Clean Architecture** с разделением на слои:

```
├── cmd/                    # Точки входа
│   ├── api/               # HTTP API сервер
│   └── worker/            # Фоновый worker для обновления цен
├── internal/
│   ├── domain/            # Бизнес-логика и модели
│   │   ├── models/        # Доменные модели
│   │   └── repository/    # Интерфейсы репозиториев
│   ├── application/       # Слой приложения
│   │   ├── dto/          # Data Transfer Objects
│   │   └── services/     # Бизнес-сервисы
│   ├── infrastructure/    # Внешние зависимости
│   │   ├── postgres/     # PostgreSQL репозитории
│   │   ├── coingecko/    # CoinGecko API клиент
│   │   └── config/       # Конфигурация
│   └── delivery/         # Слой доставки
│       ├── http/         # HTTP handlers
│       └── middleware/   # Middleware
├── migrations/           # Миграции базы данных
├── docs/                # Swagger документация
└── examples/            # Примеры использования
```

## 🚀 Быстрый старт

### С Docker Compose (рекомендуется)

```bash
# Клонировать репозиторий
git clone <repository-url>
cd crypto-price-tracker

# Запустить все сервисы
docker-compose up -d

# Применить миграции
docker-compose run --rm migrate migrate -path /migrations -database "postgres://postgres:password@postgres:5432/crypto_tracker?sslmode=disable" up

# Проверить статус
docker-compose ps
```

### Локальная разработка

```bash
# Установить зависимости
go mod download

# Запустить PostgreSQL
docker-compose up -d postgres

# Применить миграции
make migrate-up

# Запустить API
make run-api

# Запустить Worker (в отдельном терминале)
make run-worker
```

## 📊 API Endpoints

### Добавить криптовалюту
```bash
curl -X POST http://localhost:8080/api/v1/currency/add \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTC",
    "api_id": "bitcoin",
    "interval": 60
  }'
```

**Параметры:**
- `symbol` - символ валюты (BTC, ETH, USDT)
- `api_id` - идентификатор для API (bitcoin, ethereum, tether)
- `interval` - интервал обновления в секундах (мин. 30)

### Получить цену криптовалюты
```bash
curl "http://localhost:8080/api/v1/currency/price?coin=BTC&timestamp=$(date +%s)"
```

**Параметры:**
- `coin` - символ валюты (BTC, ETH)
- `timestamp` - Unix timestamp

### Получить список активных валют
```bash
curl http://localhost:8080/api/v1/currency/list
```

### Удалить криптовалюту
```bash
curl -X POST http://localhost:8080/api/v1/currency/remove \
  -H "Content-Type: application/json" \
  -d '{"symbol": "BTC"}'
```

### Проверка здоровья
```bash
curl http://localhost:8080/health
```

## 🔧 Конфигурация

### Переменные окружения

```bash
# База данных
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=crypto_tracker
DB_SSLMODE=disable

# API
API_PORT=8080

# Worker
WORKER_INTERVAL=60

# Логирование
LOG_LEVEL=info
```

### Файл конфигурации

Создайте `configs/config.yaml`:

```yaml
database:
  host: postgres
  port: 5432
  user: postgres
  password: password
  dbname: crypto_tracker
  sslmode: disable

api:
  port: 8080

worker:
  interval: 60

logging:
  level: info
```

## 🗄️ База данных

### Структура таблиц

#### currencies
```sql
CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) UNIQUE NOT NULL,
    api_id VARCHAR(50) NOT NULL,
    interval INTEGER NOT NULL DEFAULT 60,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### prices
```sql
CREATE TABLE prices (
    id SERIAL PRIMARY KEY,
    currency_id INTEGER REFERENCES currencies(id),
    price DECIMAL(20,8) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🛠️ Команды Makefile

```bash
# Запуск
make up              # Запустить все сервисы
make down            # Остановить все сервисы
make logs            # Показать логи

# База данных
make db-up           # Запустить PostgreSQL
make db-down         # Остановить PostgreSQL
make migrate-up      # Применить миграции
make migrate-down    # Откатить миграции

# Разработка
make build           # Собрать образы
make test            # Запустить тесты
make clean           # Очистить артефакты

# Локальный запуск
make run-api         # Запустить API локально
make run-worker      # Запустить Worker локально
```

## 📚 Примеры использования

### Python
```python
import requests

# Добавить Bitcoin
response = requests.post('http://localhost:8080/api/v1/currency/add', json={
    'symbol': 'BTC',
    'api_id': 'bitcoin',
    'interval': 60
})

# Получить цену
price = requests.get('http://localhost:8080/api/v1/currency/price', params={
    'coin': 'BTC',
    'timestamp': int(time.time())
})
```

### JavaScript
```javascript
// Добавить Ethereum
fetch('http://localhost:8080/api/v1/currency/add', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        symbol: 'ETH',
        api_id: 'ethereum',
        interval: 60
    })
});

// Получить цену
fetch(`http://localhost:8080/api/v1/currency/price?coin=ETH&timestamp=${Date.now()/1000}`)
    .then(response => response.json())
    .then(data => console.log(data));
```

## 🔍 Мониторинг

### Логи
```bash
# API логи
docker-compose logs -f api

# Worker логи
docker-compose logs -f worker

# База данных
docker-compose logs -f postgres
```

### Метрики
- HTTP запросы логируются с помощью Zap
- Время выполнения запросов
- Ошибки и предупреждения

## 🧪 Тестирование

```bash
# Запустить все тесты
make test

# Запустить тесты с покрытием
go test -cover ./...

# Запустить конкретный тест
go test -v ./internal/application/services
```

## 📖 Swagger документация

После запуска API, документация доступна по адресу:
- Swagger UI: http://localhost:8080/docs
- JSON: http://localhost:8080/docs/swagger.json

## 🐛 Отладка

### Проверка подключения к БД
```bash
docker-compose exec postgres psql -U postgres -d crypto_tracker -c "SELECT * FROM currencies;"
```

### Проверка логов
```bash
docker-compose logs api | grep ERROR
```

### Проверка статуса сервисов
```bash
docker-compose ps
curl http://localhost:8080/health
```

## 📄 Лицензия

MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📞 Поддержка

При возникновении проблем:
1. Проверьте логи: `docker-compose logs`
2. Убедитесь, что все сервисы запущены: `docker-compose ps`
3. Проверьте подключение к БД
4. Создайте issue в репозитории
