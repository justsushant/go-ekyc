DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'file_upload_type') THEN
        CREATE TYPE FILE_UPLOAD_TYPE AS ENUM ('face', 'id_card');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS plan (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    base_cost NUMERIC(10, 2) NOT NULL,
    per_call_cost NUMERIC(10, 2) NOT NULL,
    upload_cost_per_mb NUMERIC(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS client (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    plan_id INTEGER NOT NULL,
    refresh_token VARCHAR(50),
    FOREIGN KEY (plan_id) REFERENCES plan(id)
);

CREATE TABLE IF NOT EXISTS upload (
    id SERIAL PRIMARY KEY,
    type FILE_UPLOAD_TYPE,
    client_id INTEGER NOT NULL,
    file_path VARCHAR(100) NOT NULL,
    file_size_kb BIGINT NOT NULL,
    FOREIGN KEY (client_id) REFERENCES client(id)
);

CREATE TABLE IF NOT EXISTS face_match (
    id SERIAL PRIMARY KEY,
    client_id INTEGER,
    upload_id1 INTEGER,
    upload_id2 INTEGER,
    match_score INTEGER,
    FOREIGN KEY (client_id) REFERENCES client(id),
    FOREIGN KEY (upload_id1) REFERENCES upload(id),
    FOREIGN KEY (upload_id2) REFERENCES upload(id)
);

CREATE TABLE IF NOT EXISTS ocr (
    id SERIAL PRIMARY KEY,
    client_id INTEGER,
    upload_id INTEGER,
    details JSON,
    FOREIGN KEY (client_id) REFERENCES client(id),
    FOREIGN KEY (upload_id) REFERENCES upload(id)
);

INSERT INTO plan (name, base_cost, per_call_cost, upload_cost_per_mb)
VALUES 
    ('basic', '10', '0.1', '0.1'),
    ('advanced', '15', '0.05', '0.05'),
    ('enterprise', '20', '0.1', '0.01');