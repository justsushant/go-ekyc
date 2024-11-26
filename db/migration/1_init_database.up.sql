-- Create the custom type `FILE_UPLOAD_TYPE` as an ENUM, but only if it doesn't already exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'file_upload_type') THEN
        CREATE TYPE FILE_UPLOAD_TYPE AS ENUM ('face', 'id_card');
    END IF;
END
$$;

-- Create the `plan` table if it does not already exist
CREATE TABLE IF NOT EXISTS plan (
    id SERIAL NOT NULL PRIMARY KEY,            -- Primary key for the table
    name VARCHAR(50) NOT NULL,                -- Name of the plan (e.g., 'basic', 'advanced')
    base_cost NUMERIC(10, 2) NOT NULL,        -- Base cost of the plan
    per_call_cost NUMERIC(10, 2) NOT NULL,    -- Cost per API call
    upload_cost_per_mb NUMERIC(10, 2) NOT NULL -- Cost per MB of upload
);

-- Create the `client` table if it does not already exist
CREATE TABLE IF NOT EXISTS client (
    id SERIAL NOT NULL PRIMARY KEY,           -- Primary key for the client
    name VARCHAR(50) NOT NULL,                -- Name of the client
    email VARCHAR(50) NOT NULL,               -- Email address of the client
    plan_id INTEGER NOT NULL,                 -- Foreign key referencing the `plan` table
    refresh_token VARCHAR(50),                -- Token used for authentication/refresh
    FOREIGN KEY (plan_id) REFERENCES plan(id) -- Enforce plan_id must exist in `plan`
);

-- Create the `upload` table if it does not already exist
CREATE TABLE IF NOT EXISTS upload (
    id SERIAL PRIMARY KEY,                    -- Primary key for the upload
    type FILE_UPLOAD_TYPE,                    -- Type of upload, referencing the ENUM
    client_id INTEGER NOT NULL,               -- Foreign key referencing the `client` table
    file_path VARCHAR(100) NOT NULL,          -- Path to the uploaded file
    file_size_kb BIGINT NOT NULL,             -- Size of the uploaded file in KB
    FOREIGN KEY (client_id) REFERENCES client(id) -- Enforce client_id must exist in `client`
);

-- Create the `face_match` table if it does not already exist
CREATE TABLE IF NOT EXISTS face_match (
    id SERIAL PRIMARY KEY,                    -- Primary key for the face match job
    client_id INTEGER,                        -- Foreign key referencing the `client` table
    upload_id1 INTEGER,                       -- Foreign key referencing the first uploaded file
    upload_id2 INTEGER,                       -- Foreign key referencing the second uploaded file
    match_score INTEGER,                      -- Score representing the similarity of the two faces
    FOREIGN KEY (client_id) REFERENCES client(id), -- Enforce client_id must exist in `client`
    FOREIGN KEY (upload_id1) REFERENCES upload(id), -- Enforce upload_id1 must exist in `upload`
    FOREIGN KEY (upload_id2) REFERENCES upload(id)  -- Enforce upload_id2 must exist in `upload`
);

-- Create the `ocr` table if it does not already exist
CREATE TABLE IF NOT EXISTS ocr (
    id SERIAL PRIMARY KEY,                    -- Primary key for the OCR job
    client_id INTEGER,                        -- Foreign key referencing the `client` table
    upload_id INTEGER,                        -- Foreign key referencing the uploaded file
    details JSON,                             -- JSON field for OCR details
    FOREIGN KEY (client_id) REFERENCES client(id), -- Enforce client_id must exist in `client`
    FOREIGN KEY (upload_id) REFERENCES upload(id)  -- Enforce upload_id must exist in `upload`
);

-- Insert default plans into the `plan` table
INSERT INTO plan (name, base_cost, per_call_cost, upload_cost_per_mb)
VALUES 
    ('basic', '10', '0.1', '0.1'),            -- Basic plan with fixed pricing
    ('advanced', '15', '0.05', '0.05'),       -- Advanced plan with cheaper per-call cost
    ('enterprise', '20', '0.1', '0.01');      -- Enterprise plan with reduced upload cost
