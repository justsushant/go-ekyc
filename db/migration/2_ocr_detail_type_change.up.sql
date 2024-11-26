-- Modify the `details` column in the `ocr` table to use the JSONB type for better performance and indexing
-- The `USING details::JSONB` ensures existing JSON data is cast to JSONB
ALTER TABLE ocr
ALTER COLUMN details TYPE JSONB
USING details::JSONB;
