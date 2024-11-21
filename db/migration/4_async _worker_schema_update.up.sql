ALTER TABLE face_match
ADD COLUMN job_id VARCHAR(100),
ADD COLUMN status VARCHAR(10),
ADD COLUMN completed_at TIMESTAMP,
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN processed_at TIMESTAMP,
ADD COLUMN failed_reason VARCHAR(100), 
ADD COLUMN failed_at TIMESTAMP;

ALTER TABLE ocr
ADD COLUMN job_id VARCHAR(100),
ADD COLUMN status VARCHAR(10),
ADD COLUMN completed_at TIMESTAMP,
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN processed_at TIMESTAMP,
ADD COLUMN failed_reason VARCHAR(100), 
ADD COLUMN failed_at TIMESTAMP;