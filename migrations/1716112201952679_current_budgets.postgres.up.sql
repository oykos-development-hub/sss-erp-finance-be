CREATE TABLE IF NOT EXISTS current_budgets (
    id serial PRIMARY KEY,
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL,
    actual DECIMAL (10, 2) NOT NULL,
    initial_actual DECIMAL (10, 2) NOT NULL,
    balance DECIMAL (10, 2) NOT NULL,
    created_at TIMESTAMP
);
