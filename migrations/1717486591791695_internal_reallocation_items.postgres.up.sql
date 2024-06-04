CREATE TABLE IF NOT EXISTS internal_reallocation_items (
    id serial PRIMARY KEY,
    reallocation_id INTEGER REFERENCES internal_reallocations(id) ON DELETE CASCADE,
    source_account_id INTEGER,
    destination_account_id INTEGER,
    amount DECIMAL (10, 2) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
