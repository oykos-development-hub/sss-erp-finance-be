CREATE TABLE IF NOT EXISTS spending_releases (
    id serial PRIMARY KEY,
    current_budget_id INTEGER REFERENCES current_budgets ON DELETE CASCADE,
    month INTEGER NOT NULL,
    value DECIMAL (10, 2) NOT NULL,
    created_at TIMESTAMP
);
