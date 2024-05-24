CREATE TABLE IF NOT EXISTS spending_dynamics (
    id serial PRIMARY KEY,
    current_budget_id INTEGER REFERENCES current_budgets ON DELETE CASCADE,
    actual DECIMAL(10,2) NOT NULL
);
