CREATE TABLE IF NOT EXISTS non_financial_budget_goals (
    id serial PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    activity_id INTEGER NOT NULL REFERENCES activities(id),
    non_financial_budget_id INTEGER NOT NULL REFERENCES non_financial_budgets(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
