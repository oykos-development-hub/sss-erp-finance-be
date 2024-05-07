CREATE TABLE IF NOT EXISTS accounting_entries (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
    organization_unit_id INTEGER NOT NULL,
    date_of_booking TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
