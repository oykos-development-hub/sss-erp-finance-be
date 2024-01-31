CREATE TABLE IF NOT EXISTS activities (
    id serial PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    code text NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
