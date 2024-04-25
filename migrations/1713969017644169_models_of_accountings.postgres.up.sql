CREATE TABLE IF NOT EXISTS models_of_accounting (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
    type TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
