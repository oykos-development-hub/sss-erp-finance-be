CREATE TABLE IF NOT EXISTS financial_budget_limits (
    id serial PRIMARY KEY,
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    organization_unit_id INTEGER NOT NULL,
    limit_value INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
