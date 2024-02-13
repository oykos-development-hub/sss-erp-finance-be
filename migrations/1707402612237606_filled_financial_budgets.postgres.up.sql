CREATE TABLE IF NOT EXISTS filled_financial_budgets (
    id serial PRIMARY KEY,
    budget_request_id INTEGER NOT NULL REFERENCES budget_requests(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL,
    current_year INTEGER NOT NULL,
    next_year INTEGER NOT NULL,
    year_after_next INTEGER NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
