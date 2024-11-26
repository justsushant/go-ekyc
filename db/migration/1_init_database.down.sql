-- Drop tables in reverse order of dependencies to maintain referential integrity
DROP TABLE IF EXISTS ocr;
DROP TABLE IF EXISTS face_match;
DROP TABLE IF EXISTS upload;
DROP TABLE IF EXISTS client;
DROP TABLE IF EXISTS plan;

-- Drop the custom type
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'file_upload_type') THEN
        DROP TYPE FILE_UPLOAD_TYPE;
    END IF;
END
$$;
