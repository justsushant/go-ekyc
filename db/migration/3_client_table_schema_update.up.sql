ALTER TABLE client
ADD COLUMN access_key VARCHAR(10),
ADD COLUMN secret_key_hash VARCHAR(200),
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;