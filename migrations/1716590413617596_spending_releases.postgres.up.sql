CREATE TABLE IF NOT EXISTS spending_releases (
    id serial PRIMARY KEY,
    current_budget_id INTEGER NOT NULL REFERENCES current_budgets ON DELETE CASCADE,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    value DECIMAL (10, 2) NOT NULL,
    username TEXT,
    created_at TIMESTAMP
);
