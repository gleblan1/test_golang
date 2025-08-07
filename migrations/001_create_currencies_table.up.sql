-- Создание таблицы криптовалют
CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(50) UNIQUE NOT NULL,
    interval INTEGER NOT NULL DEFAULT 60,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание индекса для быстрого поиска по символу
CREATE INDEX IF NOT EXISTS idx_currencies_symbol ON currencies(symbol);

-- Создание индекса для активных криптовалют
CREATE INDEX IF NOT EXISTS idx_currencies_active ON currencies(is_active) WHERE is_active = true;

-- Создание триггера для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_currencies_updated_at 
    BEFORE UPDATE ON currencies 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column(); 