-- Add new columns to the `client` table to store access and secret key information, along with a timestamp
ALTER TABLE client
ADD COLUMN access_key VARCHAR(10),                    -- Stores a short access key for the client
ADD COLUMN secret_key_hash VARCHAR(200),             -- Stores the hashed value of the client's secret key
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP; -- Records when the client record was created, with a default value of the current timestamp
