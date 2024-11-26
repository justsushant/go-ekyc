-- Create a new ENUM type `STATUS_TYPE` if it does not already exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type') THEN
        CREATE TYPE STATUS_TYPE AS ENUM ('created', 'processing', 'completed', 'failed');
    END IF;
END
$$;

-- Alter the `status` column in the `face_match` table to use the `STATUS_TYPE` ENUM type
-- The `USING status::STATUS_TYPE` ensures existing data is safely cast to the new ENUM type
ALTER TABLE face_match
    ALTER COLUMN status TYPE STATUS_TYPE
    USING status::STATUS_TYPE;

-- Alter the `status` column in the `ocr` table to use the `STATUS_TYPE` ENUM type
-- The `USING status::STATUS_TYPE` ensures existing data is safely cast to the new ENUM type
ALTER TABLE ocr
    ALTER COLUMN status TYPE STATUS_TYPE
    USING status::STATUS_TYPE;

-- Add a `created_at` column to the `upload` table to track when records are created
ALTER TABLE upload
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
