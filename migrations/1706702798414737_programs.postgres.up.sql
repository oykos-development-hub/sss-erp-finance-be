CREATE TABLE IF NOT EXISTS programs (
    id serial PRIMARY KEY,
    parent_id INTEGER REFERENCES programs(id),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    code text NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
