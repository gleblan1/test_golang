-- Удаление индексов
DROP INDEX IF EXISTS idx_prices_currency_created_at;
DROP INDEX IF EXISTS idx_prices_currency_timestamp_unique;
DROP INDEX IF EXISTS idx_prices_timestamp;
DROP INDEX IF EXISTS idx_prices_currency_timestamp;

-- Удаление таблицы
DROP TABLE IF EXISTS prices; 