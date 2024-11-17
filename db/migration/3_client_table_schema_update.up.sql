ALTER TABLE client
DROP COLUMN refresh_token,
ADD COLUMN access_key VARCHAR(10),
ADD COLUMN secret_key_hash VARCHAR(200),
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;