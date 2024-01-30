CREATE TABLE IF NOT EXISTS financial_budget_limits (
    id serial PRIMARY KEY,
    financial_budget_id INTEGER NOT NULL REFERENCES financial_budgets(id) ON DELETE CASCADE,
    organization_unit_id INTEGER NOT NULL,
    limit_value INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
