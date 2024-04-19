CREATE TABLE IF NOT EXISTS budget_requests (
    id serial PRIMARY KEY,
    parent_id INTEGER REFERENCES budget_requests(id) ON DELETE CASCADE,
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    organization_unit_id INTEGER NOT NULL,
    request_type INTEGER NOT NULL,
    status INTEGER NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
