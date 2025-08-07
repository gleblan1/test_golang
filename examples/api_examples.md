# Примеры API запросов

## Добавление криптовалюты в отслеживание

### Запрос
```bash
curl -X POST http://localhost:8080/api/v1/currency/add \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "bitcoin",
    "interval": 60
  }'
```

### Ответ
```json
{
  "id": 1,
  "symbol": "bitcoin",
  "interval": 60,
  "is_active": true,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

## Удаление криптовалюты из отслеживания

### Запрос
```bash
curl -X POST http://localhost:8080/api/v1/currency/remove \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "bitcoin"
  }'
```

### Ответ
```json
{
  "message": "Currency removed successfully"
}
```

## Получение цены криптовалюты

### Запрос
```bash
curl -X GET "http://localhost:8080/api/v1/currency/price?coin=bitcoin&timestamp=1640995200"
```

### Ответ
```json
{
  "id": 1,
  "symbol": "bitcoin",
  "price": 50000.00,
  "timestamp": "2022-01-01T00:00:00Z",
  "created_at": "2022-01-01T00:00:00Z"
}
```

## Получение списка всех активных криптовалют

### Запрос
```bash
curl -X GET http://localhost:8080/api/v1/currency/list
```

### Ответ
```json
[
  {
    "id": 1,
    "symbol": "bitcoin",
    "interval": 60,
    "is_active": true,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  },
  {
    "id": 2,
    "symbol": "ethereum",
    "interval": 120,
    "is_active": true,
    "created_at": "2024-01-01T12:30:00Z",
    "updated_at": "2024-01-01T12:30:00Z"
  }
]
```

## Проверка здоровья сервиса

### Запрос
```bash
curl -X GET http://localhost:8080/health
```

### Ответ
```json
{
  "message": "Service is healthy"
}
```

## Примеры ошибок

### Криптовалюта уже существует
```bash
curl -X POST http://localhost:8080/api/v1/currency/add \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "bitcoin",
    "interval": 60
  }'
```

### Ответ
```json
{
  "error": "currency_error",
  "message": "currency already exists",
  "code": 409
}
```

### Криптовалюта не найдена
```bash
curl -X GET "http://localhost:8080/api/v1/currency/price?coin=nonexistent&timestamp=1640995200"
```

### Ответ
```json
{
  "error": "price_error",
  "message": "currency not found",
  "code": 404
}
```

### Неверные параметры
```bash
curl -X GET "http://localhost:8080/api/v1/currency/price?coin=bitcoin&timestamp=invalid"
```

### Ответ
```json
{
  "error": "validation_error",
  "message": "invalid timestamp format",
  "code": 400
}
```

## Использование с Python

```python
import requests
import json

# Базовый URL
base_url = "http://localhost:8080/api/v1"

# Добавление криптовалюты
def add_currency(symbol, interval):
    url = f"{base_url}/currency/add"
    data = {
        "symbol": symbol,
        "interval": interval
    }
    response = requests.post(url, json=data)
    return response.json()

# Получение цены
def get_price(coin, timestamp):
    url = f"{base_url}/currency/price"
    params = {
        "coin": coin,
        "timestamp": timestamp
    }
    response = requests.get(url, params=params)
    return response.json()

# Примеры использования
if __name__ == "__main__":
    # Добавляем Bitcoin
    result = add_currency("bitcoin", 60)
    print("Added currency:", result)
    
    # Получаем цену Bitcoin
    import time
    current_timestamp = int(time.time())
    price = get_price("bitcoin", current_timestamp)
    print("Current price:", price)
```

## Использование с JavaScript

```javascript
// Базовый URL
const baseUrl = 'http://localhost:8080/api/v1';

// Добавление криптовалюты
async function addCurrency(symbol, interval) {
    const response = await fetch(`${baseUrl}/currency/add`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            symbol: symbol,
            interval: interval
        })
    });
    return await response.json();
}

// Получение цены
async function getPrice(coin, timestamp) {
    const response = await fetch(`${baseUrl}/currency/price?coin=${coin}&timestamp=${timestamp}`);
    return await response.json();
}

// Примеры использования
async function main() {
    try {
        // Добавляем Bitcoin
        const result = await addCurrency('bitcoin', 60);
        console.log('Added currency:', result);
        
        // Получаем цену Bitcoin
        const currentTimestamp = Math.floor(Date.now() / 1000);
        const price = await getPrice('bitcoin', currentTimestamp);
        console.log('Current price:', price);
    } catch (error) {
        console.error('Error:', error);
    }
}

main();
``` 