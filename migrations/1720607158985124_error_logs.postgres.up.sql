CREATE TABLE IF NOT EXISTS error_logs (
    id serial PRIMARY KEY,
    error VARCHAR ( 255 ) NOT NULL,
    code INTEGER,
    entity TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
