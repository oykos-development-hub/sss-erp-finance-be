CREATE TABLE IF NOT EXISTS budget_requests (
    id serial PRIMARY KEY,
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    organization_unit_id INTEGER NOT NULL,
    request_type INTEGER NOT NULL,
    status INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
