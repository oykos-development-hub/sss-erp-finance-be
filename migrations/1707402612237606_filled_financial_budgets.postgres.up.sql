CREATE TABLE IF NOT EXISTS filled_financial_budgets (
    id serial PRIMARY KEY,
    organization_unit_id INTEGER NOT NULL,
    finance_budget_id INTEGER NOT NULL REFERENCES financial_budgets(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL,
    current_year INTEGER NOT NULL,
    next_year INTEGER NOT NULL,
    year_after_next INTEGER NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
