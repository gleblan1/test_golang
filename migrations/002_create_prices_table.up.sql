-- Создание таблицы цен
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    currency_id INTEGER NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    price DECIMAL(20, 8) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание индекса для быстрого поиска по криптовалюте и времени
CREATE INDEX IF NOT EXISTS idx_prices_currency_timestamp ON prices(currency_id, timestamp);

-- Создание индекса для поиска по времени
CREATE INDEX IF NOT EXISTS idx_prices_timestamp ON prices(timestamp);

-- Создание уникального индекса для предотвращения дублирования цен
CREATE UNIQUE INDEX IF NOT EXISTS idx_prices_currency_timestamp_unique ON prices(currency_id, timestamp);

-- Создание индекса для последних цен
CREATE INDEX IF NOT EXISTS idx_prices_currency_created_at ON prices(currency_id, created_at DESC); 