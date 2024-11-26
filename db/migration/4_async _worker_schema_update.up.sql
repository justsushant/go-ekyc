-- Add job-related columns to the `face_match` table to track job status, timing, and failure reasons
ALTER TABLE face_match
ADD COLUMN job_id VARCHAR(100),                          -- Unique identifier for the job
ADD COLUMN status VARCHAR(10),                           -- Current status of the job (e.g., 'processing', 'completed', 'failed')
ADD COLUMN completed_at TIMESTAMP,                       -- Timestamp indicating when the job was completed
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp indicating when the record was created
ADD COLUMN processed_at TIMESTAMP,                       -- Timestamp indicating when the job started processing
ADD COLUMN failed_reason VARCHAR(100),                   -- Reason for job failure, if applicable
ADD COLUMN failed_at TIMESTAMP;                          -- Timestamp indicating when the job failed

-- Add job-related columns to the `ocr` table to track job status, timing, and failure reasons
ALTER TABLE ocr
ADD COLUMN job_id VARCHAR(100),                          -- Unique identifier for the job
ADD COLUMN status VARCHAR(10),                           -- Current status of the job (e.g., 'processing', 'completed', 'failed')
ADD COLUMN completed_at TIMESTAMP,                       -- Timestamp indicating when the job was completed
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp indicating when the record was created
ADD COLUMN processed_at TIMESTAMP,                       -- Timestamp indicating when the job started processing
ADD COLUMN failed_reason VARCHAR(100),                   -- Reason for job failure, if applicable
ADD COLUMN failed_at TIMESTAMP;                          -- Timestamp indicating when the job failed
