CREATE TABLE IF NOT EXISTS budgets (
    id serial PRIMARY KEY,
    budget_type INTEGER NOT NULL,
    budget_status INTEGER NOT NULL DEFAULT 1,
    year INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
