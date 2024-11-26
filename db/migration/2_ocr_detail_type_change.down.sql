-- Reverse the change: Alter the `details` column in the `ocr` table back to JSON
ALTER TABLE ocr
ALTER COLUMN details TYPE JSON
USING details::JSON;
