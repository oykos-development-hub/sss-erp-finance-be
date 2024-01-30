CREATE TABLE IF NOT EXISTS financial_budgets (
    id serial PRIMARY KEY,
    account_version INTEGER NOT NULL,
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
