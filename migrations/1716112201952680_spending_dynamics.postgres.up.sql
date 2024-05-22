CREATE TABLE IF NOT EXISTS spending_dynamics (
    id serial PRIMARY KEY,
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL,
    actual DECIMAL(10,2) NOT NULL
);
