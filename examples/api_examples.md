# Примеры API запросов

## Добавить Bitcoin

```bash
curl -X POST http://localhost:8080/api/v1/currency/add \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTC",
    "api_id": "bitcoin",
    "interval": 60
  }'
```

**Ответ:**
```json
{
  "id": 1,
  "symbol": "BTC",
  "api_id": "bitcoin",
  "interval": 60,
  "is_active": true,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

## Добавить Ethereum

```bash
curl -X POST http://localhost:8080/api/v1/currency/add \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "ETH",
    "api_id": "ethereum",
    "interval": 60
  }'
```

## Добавить Tether

```bash
curl -X POST http://localhost:8080/api/v1/currency/add \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "USDT",
    "api_id": "tether",
    "interval": 60
  }'
```

## Получить цену Bitcoin

```bash
curl "http://localhost:8080/api/v1/currency/price?coin=BTC&timestamp=$(date +%s)"
```

**Ответ:**
```json
{
  "id": 1,
  "symbol": "BTC",
  "price": 116388.00,
  "timestamp": "2024-01-01T12:00:00Z",
  "created_at": "2024-01-01T12:00:00Z"
}
```

## Получить цену Ethereum

```bash
curl "http://localhost:8080/api/v1/currency/price?coin=ETH&timestamp=$(date +%s)"
```

## Получить список активных валют

```bash
curl http://localhost:8080/api/v1/currency/list
```

**Ответ:**
```json
[
  {
    "id": 1,
    "symbol": "BTC",
    "api_id": "bitcoin",
    "interval": 60,
    "is_active": true,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  },
  {
    "id": 2,
    "symbol": "ETH",
    "api_id": "ethereum",
    "interval": 60,
    "is_active": true,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
]
```

## Удалить криптовалюту

```bash
curl -X POST http://localhost:8080/api/v1/currency/remove \
  -H "Content-Type: application/json" \
  -d '{"symbol": "BTC"}'
```

**Ответ:**
```json
{
  "message": "Currency removed successfully"
}
```

## Проверка здоровья сервиса

```bash
curl http://localhost:8080/health
```

**Ответ:**
```json
{
  "status": "ok",
  "service": "crypto-price-tracker"
}
```

## Python примеры

```python
import requests
import time

# Добавить Bitcoin
response = requests.post('http://localhost:8080/api/v1/currency/add', json={
    'symbol': 'BTC',
    'api_id': 'bitcoin',
    'interval': 60
})
print(response.json())

# Получить цену
timestamp = int(time.time())
price_response = requests.get('http://localhost:8080/api/v1/currency/price', params={
    'coin': 'BTC',
    'timestamp': timestamp
})
print(price_response.json())
```

## JavaScript примеры

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
})
.then(response => response.json())
.then(data => console.log(data));

// Получить цену
const timestamp = Math.floor(Date.now() / 1000);
fetch(`http://localhost:8080/api/v1/currency/price?coin=ETH&timestamp=${timestamp}`)
    .then(response => response.json())
    .then(data => console.log(data));
```

## Node.js примеры

```javascript
const axios = require('axios');

// Добавить Tether
async function addCurrency() {
    try {
        const response = await axios.post('http://localhost:8080/api/v1/currency/add', {
            symbol: 'USDT',
            api_id: 'tether',
            interval: 60
        });
        console.log(response.data);
    } catch (error) {
        console.error('Error:', error.response.data);
    }
}

// Получить цену
async function getPrice() {
    try {
        const timestamp = Math.floor(Date.now() / 1000);
        const response = await axios.get(`http://localhost:8080/api/v1/currency/price?coin=USDT&timestamp=${timestamp}`);
        console.log(response.data);
    } catch (error) {
        console.error('Error:', error.response.data);
    }
}

addCurrency();
getPrice();
``` 