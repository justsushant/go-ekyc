-- -- Reverse the changes made in the up migration

-- -- Remove the `created_at` column from the `upload` table
-- ALTER TABLE upload
-- DROP COLUMN IF EXISTS created_at;

-- -- Revert the `status` column in `face_match` table to its previous type (VARCHAR)
-- -- The `USING status::VARCHAR` safely converts the ENUM back to VARCHAR
-- ALTER TABLE face_match
--     ALTER COLUMN status TYPE VARCHAR(10)
--     USING status::VARCHAR;

-- -- Revert the `status` column in `ocr` table to its previous type (VARCHAR)
-- -- The `USING status::VARCHAR` safely converts the ENUM back to VARCHAR
-- ALTER TABLE ocr
--     ALTER COLUMN status TYPE VARCHAR(10)
--     USING status::VARCHAR;

-- -- Drop the `STATUS_TYPE` ENUM type if it exists
-- DO $$
-- BEGIN
--     IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type') THEN
--         DROP TYPE STATUS_TYPE;
--     END IF;
-- END
-- $$;

-- Reverse the changes made in the up migration

-- Remove the `created_at` column from the `upload` table
ALTER TABLE upload
DROP COLUMN IF EXISTS created_at;

-- Revert the `status` column in `face_match` table to its previous type (VARCHAR), if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'face_match' AND column_name = 'status'
    ) THEN
        ALTER TABLE face_match
        ALTER COLUMN status TYPE VARCHAR(10)
        USING status::VARCHAR;
    END IF;
END
$$;

-- Revert the `status` column in `ocr` table to its previous type (VARCHAR), if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'ocr' AND column_name = 'status'
    ) THEN
        ALTER TABLE ocr
        ALTER COLUMN status TYPE VARCHAR(10)
        USING status::VARCHAR;
    END IF;
END
$$;

-- Drop the `STATUS_TYPE` ENUM type if it exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type') THEN
        DROP TYPE STATUS_TYPE;
    END IF;
END
$$;
