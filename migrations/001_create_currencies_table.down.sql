-- Удаление триггера
DROP TRIGGER IF EXISTS update_currencies_updated_at ON currencies;

-- Удаление функции
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаление индексов
DROP INDEX IF EXISTS idx_currencies_active;
DROP INDEX IF EXISTS idx_currencies_symbol;

-- Удаление таблицы
DROP TABLE IF EXISTS currencies; 