CREATE TABLE IF NOT EXISTS internal_reallocations (
    id serial PRIMARY KEY,
    title TEXT,
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    organization_unit_id INTEGER NOT NULL,
    date_of_request TIMESTAMP NOT NULL,
    requested_by INTEGER NOT NULL,
    file_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
