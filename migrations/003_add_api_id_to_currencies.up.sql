ALTER TABLE currencies ADD COLUMN api_id VARCHAR(50);
UPDATE currencies SET api_id = symbol WHERE api_id IS NULL;
ALTER TABLE currencies ALTER COLUMN api_id SET NOT NULL;
CREATE INDEX idx_currencies_api_id ON currencies(api_id); 