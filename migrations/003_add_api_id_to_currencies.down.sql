DROP INDEX IF EXISTS idx_currencies_api_id;
ALTER TABLE currencies DROP COLUMN IF EXISTS api_id; 