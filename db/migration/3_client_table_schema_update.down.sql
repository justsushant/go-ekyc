-- Remove the columns added to the `client` table in the up migration
ALTER TABLE client
DROP COLUMN IF EXISTS access_key,                -- Removes the `access_key` column if it exists
DROP COLUMN IF EXISTS secret_key_hash,           -- Removes the `secret_key_hash` column if it exists
DROP COLUMN IF EXISTS created_at;                -- Removes the `created_at` column if it exists
