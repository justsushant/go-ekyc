-- Remove job-related columns from the `face_match` table
ALTER TABLE face_match
DROP COLUMN IF EXISTS job_id,                -- Removes the `job_id` column if it exists
DROP COLUMN IF EXISTS status,               -- Removes the `status` column if it exists
DROP COLUMN IF EXISTS completed_at,         -- Removes the `completed_at` column if it exists
DROP COLUMN IF EXISTS created_at,           -- Removes the `created_at` column if it exists
DROP COLUMN IF EXISTS processed_at,         -- Removes the `processed_at` column if it exists
DROP COLUMN IF EXISTS failed_reason,        -- Removes the `failed_reason` column if it exists
DROP COLUMN IF EXISTS failed_at;            -- Removes the `failed_at` column if it exists

-- Remove job-related columns from the `ocr` table
ALTER TABLE ocr
DROP COLUMN IF EXISTS job_id,                -- Removes the `job_id` column if it exists
DROP COLUMN IF EXISTS status,               -- Removes the `status` column if it exists
DROP COLUMN IF EXISTS completed_at,         -- Removes the `completed_at` column if it exists
DROP COLUMN IF EXISTS created_at,           -- Removes the `created_at` column if it exists
DROP COLUMN IF EXISTS processed_at,         -- Removes the `processed_at` column if it exists
DROP COLUMN IF EXISTS failed_reason,        -- Removes the `failed_reason` column if it exists
DROP COLUMN IF EXISTS failed_at;            -- Removes the `failed_at` column if it exists
